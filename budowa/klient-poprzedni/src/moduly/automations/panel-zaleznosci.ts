import {
  AutomationDependencyKind,
  type AutomationDependency,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przyciskAkcji as przycisk, pole, wybor } from '../../modele/kontrolki-formularza';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Panel zależności układu — wejście do rodziny `orchestration.*`. Okno
 * Orchestratora zapisuje układ w całości komendą
 * `automation.orchestrator.define`; ta rodzina mówi o tym samym grafie inaczej:
 * czyta go bez zapisu (`dependency.list`), zmienia pojedynczą krawędź
 * (`dependency.set`) i sprawdza układ bez dotykania go (`validate`).
 *
 * Zależności są panelem, a nie osobnym oknem, bo rejestr okien operacyjnych
 * rdzenia nie ma dla nich kodu; panel wchodzi w pas narzędzi kontekstowych
 * Orchestratora.
 *
 * Rdzeń przyjmuje także krawędź tworzącą cykl i oddaje zastrzeżenia w wyniku,
 * więc niepowodzenie czynności znaczy co innego niż niepoprawny układ. Panel
 * trzyma te dwa zdania osobno: czynność udała się albo nie, a układ jest
 * poprawny albo ma zastrzeżenia.
 */
export interface PanelZaleznosci {
  element: HTMLElement;
}

/** Rodzaje krawędzi — ten sam wykaz co w oknie Orchestratora. */
const RODZAJE: ReadonlyArray<[string, string]> = [
  [AutomationDependencyKind.Sequential, 'sekwencyjna'],
  [AutomationDependencyKind.Parallel, 'równoległa'],
  [AutomationDependencyKind.Conditional, 'warunkowa'],
];

export function utworzPanelZaleznosci(
  zrodlo: ZrodloAutomations,
  czynnyUklad: () => string,
): PanelZaleznosci {
  const odKroku = pole('Krok poprzedzający', 'identyfikator kroku');
  const doKroku = pole('Krok następujący', 'identyfikator kroku');
  const rodzaj = wybor('Rodzaj zależności', RODZAJE);

  const zapisz = przycisk('Zapisz zależność', 'dn-btn dn-btn--sm dn-btn--atrament');
  const wczytaj = przycisk('Odczytaj układ', 'dn-btn dn-btn--sm');
  const sprawdz = przycisk('Sprawdź układ', 'dn-btn dn-btn--sm');

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dm-odpowiedz da-zaleznosci__odpowiedz';
  odpowiedz.hidden = true;

  const wykazKrawedzi = document.createElement('ul');
  wykazKrawedzi.className = 'da-zaleznosci__wykaz';

  const element = document.createElement('div');
  element.className = 'da-zaleznosci';
  element.setAttribute('aria-label', 'Zależności układu — odczyt, zapis krawędzi i sprawdzenie');
  element.append(
    odKroku,
    doKroku,
    rodzaj,
    zapisz,
    wczytaj,
    sprawdz,
    odpowiedz,
    wykazKrawedzi,
  );

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  /** Układ czynny albo zdanie, dlaczego czynności nie ma na czym wykonać. */
  function uklad(): string | null {
    const kod = czynnyUklad();
    if (kod === '') {
      powiedz('Wskaż automatykę w wykazie — zależność należy do jednego układu.', false);
      return null;
    }
    return kod;
  }

  function pokazKrawedzie(krawedzie: readonly AutomationDependency[]): void {
    wykazKrawedzi.replaceChildren(
      ...krawedzie.map((k) => {
        const wpis = document.createElement('li');
        wpis.className = 'da-zaleznosci__krawedz';
        wpis.textContent = `${k.fromStepId} → ${k.toStepId} (${k.kind})`;
        return wpis;
      }),
    );
  }

  /** Zastrzeżenia układu jednym zdaniem; brak zastrzeżeń też jest odpowiedzią. */
  function zdanieOUkladzie(poprawny: boolean, zastrzezenia?: readonly string[]): string {
    const lista = zastrzezenia ?? [];
    if (poprawny && lista.length === 0) return 'Układ bez zastrzeżeń.';
    if (lista.length === 0) return 'Rdzeń uznał układ za niepoprawny, nie podając zastrzeżeń.';
    return `Zastrzeżenia do układu (${lista.length}): ${lista.join(' · ')}`;
  }

  async function odczytaj(): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    powiedz('Odczyt zależności układu…', true);
    const wynik = await zrodlo.wykazZaleznosci({ workflowId: kod });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Odczyt zależności', wynik.blad), false);
      return;
    }
    const krawedzie = wynik.wynik.dependencies;
    pokazKrawedzie(krawedzie);
    powiedz(
      krawedzie.length === 0
        ? 'Rdzeń oddał układ bez ani jednej zależności.'
        : `Układ ma ${krawedzie.length} zależności.`,
      true,
    );
  }

  async function zapiszKrawedz(): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    const od = odKroku.value.trim();
    const doo = doKroku.value.trim();
    if (od === '' || doo === '') {
      powiedz('Wskaż oba kroki — krawędź bez końca nie jest zależnością.', false);
      return;
    }
    powiedz('Zapis zależności…', true);
    const wynik = await zrodlo.zapiszZaleznosc({
      workflowId: kod,
      dependency: {
        fromStepId: od,
        toStepId: doo,
        kind: rodzaj.value as AutomationDependency['kind'],
      },
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Zapis zależności', wynik.blad), false);
      return;
    }
    pokazKrawedzie(wynik.wynik.dependencies);
    // Powodzenie czynności i ocena układu to dwa osobne zdania — ocena dochodzi
    // za potwierdzeniem zapisu, nie zamiast niego.
    powiedz(
      `Zależność ${od} → ${doo} zapisana. ${zdanieOUkladzie(wynik.wynik.valid, wynik.wynik.issues)}`,
      true,
    );
  }

  async function sprawdzUklad(): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    powiedz('Sprawdzanie układu…', true);
    const wynik = await zrodlo.sprawdzUklad(kod);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Sprawdzenie układu', wynik.blad), false);
      return;
    }
    const sciezka = wynik.wynik.criticalPathStepIds ?? [];
    const oSciezce =
      sciezka.length === 0
        ? 'Rdzeń nie wskazał ścieżki krytycznej.'
        : `Ścieżka krytyczna: ${sciezka.join(' → ')}.`;
    powiedz(`${zdanieOUkladzie(wynik.wynik.valid, wynik.wynik.issues)} ${oSciezce}`, true);
  }

  zapisz.addEventListener('click', () => void zapiszKrawedz());
  wczytaj.addEventListener('click', () => void odczytaj());
  sprawdz.addEventListener('click', () => void sprawdzUklad());

  return { element };
}
