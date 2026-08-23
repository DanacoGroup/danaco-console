import {
  AutomationStepKind,
  type AutomationStep,
} from '../../../../shared/contract';
import {
  pole,
  poleTresci,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { ZastrzezenieDefinicji } from './walidacja-definicji';

/**
 * Edytor kroków automatyki — treść okna Workflow Builder. Obsługuje dodanie
 * kroku, ustalenie warunku i ustalenie kolejności wykonania; zapis definicji
 * należy do okna, bo to ono rozmawia z rdzeniem.
 *
 * Kolejność zmieniają przyciski „w górę/w dół”, a nie przeciąganie myszą:
 * kontrakt niesie ją liczbą (`AutomationStep.order`), a przyciski są dostępne
 * z klawiatury.
 */
export interface EdytorKrokow {
  element: HTMLElement;
  /** Kroki w kolejności wykonania — treść żądania zapisu. */
  kroki(): AutomationStep[];
  /** Wypełnia edytor krokami odczytanymi z rdzenia. */
  wczytaj(kroki: readonly AutomationStep[]): void;
  /** Dokłada krok pusty na koniec. */
  dodajKrok(): void;
  /**
   * Sygnalizuje zastrzeżenia walidacji przy krokach, których dotyczą.
   *
   * Zastrzeżenie jest ostrzeżeniem, nie bramą: wiersz dostaje znacznik
   * `data-zastrzezenie` i zdanie w dymku, ale zostaje w pełni edytowalny,
   * a zapis definicji pozostaje możliwy.
   */
  oznaczZastrzezenia(zastrzezenia: readonly ZastrzezenieDefinicji[]): void;
}

/**
 * Nazwy rodzajów kroku na ekranie.
 *
 * Mapa zupełna po wyliczeniu, nie wykaz przepisany ręcznie: gdy
 * `AutomationStepKind` urośnie, kompilacja zatrzyma się tutaj i nowy rodzaj
 * dostanie nazwę, zamiast zniknąć z pola wyboru bez śladu.
 */
const NAZWY_RODZAJOW: Readonly<Record<AutomationStepKind, string>> = {
  [AutomationStepKind.Command]: 'komenda',
  [AutomationStepKind.Model]: 'model',
  [AutomationStepKind.Condition]: 'warunek',
  [AutomationStepKind.Branch]: 'rozgałęzienie',
  [AutomationStepKind.Wait]: 'oczekiwanie',
  [AutomationStepKind.Http]: 'wywołanie zewnętrzne (HTTP)',
  [AutomationStepKind.Loop]: 'pętla',
  [AutomationStepKind.Subflow]: 'podprzepływ',
  [AutomationStepKind.Transform]: 'transformacja danych',
  [AutomationStepKind.Checkpoint]: 'punkt kontrolny',
  [AutomationStepKind.Extension]: 'rozszerzenie',
  [AutomationStepKind.Script]: 'skrypt w piaskownicy',
};

/** Rodzaje kroku w kolejności kontraktu. */
export const RODZAJE_KROKU: ReadonlyArray<[string, string]> = Object.values(AutomationStepKind).map(
  (rodzaj) => [rodzaj, NAZWY_RODZAJOW[rodzaj]],
);

export function utworzEdytorKrokow(przyZmianie: () => void): EdytorKrokow {
  const lista = document.createElement('ol');
  lista.className = 'da-kroki';
  lista.setAttribute('aria-label', 'Kroki automatyki w kolejności wykonania');

  const element = document.createElement('div');
  element.className = 'da-edytor';
  element.append(lista);

  /** Zależności kroku są własnością Orchestratora — edytor je przenosi. */
  const zaleznosci = new Map<HTMLElement, string[]>();

  /**
   * Wypełnianie edytora definicją nie jest zmianą wprowadzoną w oknie, więc na
   * czas `wczytaj` powiadomienie milknie. Bez tego każdy dokładany krok zgłasza
   * stan częściowy i zaśmieca stos cofnięć okna.
   */
  let wczytywanie = false;

  function powiadom(): void {
    if (!wczytywanie) przyZmianie();
  }

  function wierszeKrokow(): HTMLElement[] {
    return [...lista.querySelectorAll<HTMLElement>('.da-krok')];
  }

  function przenumeruj(): void {
    wierszeKrokow().forEach((rzad, numer) => {
      const podpis = rzad.querySelector<HTMLElement>('.da-krok__numer');
      if (podpis !== null) podpis.textContent = `${numer + 1}.`;
    });
  }

  function przesun(rzad: HTMLElement, oIle: number): void {
    const rzedy = wierszeKrokow();
    const miejsce = rzedy.indexOf(rzad);
    const cel = miejsce + oIle;
    if (cel < 0 || cel >= rzedy.length) return;
    if (oIle < 0) rzedy[cel]?.before(rzad);
    else rzedy[cel]?.after(rzad);
    przenumeruj();
    powiadom();
  }

  function dodaj(krok?: AutomationStep): void {
    const rzad = utworzWierszKroku(krok, powiadom, {
      wGore: (element) => przesun(element, -1),
      wDol: (element) => przesun(element, 1),
      usun: (element) => {
        zaleznosci.delete(element);
        element.remove();
        przenumeruj();
        powiadom();
      },
    });
    zaleznosci.set(rzad, [...(krok?.dependsOn ?? [])]);
    lista.append(rzad);
    przenumeruj();
    powiadom();
  }

  return {
    element,

    kroki() {
      return wierszeKrokow().map((rzad, numer) => odczytajKrok(rzad, numer, zaleznosci.get(rzad)));
    },

    wczytaj(kroki) {
      wczytywanie = true;
      try {
        lista.replaceChildren();
        zaleznosci.clear();
        for (const krok of kroki) dodaj(krok);
      } finally {
        wczytywanie = false;
      }
      przenumeruj();
    },

    dodajKrok: () => dodaj(),

    oznaczZastrzezenia(zastrzezenia) {
      const rzedy = wierszeKrokow();
      for (const rzad of rzedy) {
        delete rzad.dataset['zastrzezenie'];
        rzad.title = '';
      }
      for (const zastrzezenie of zastrzezenia) {
        const rzad = rzedy[zastrzezenie.miejsce - 1];
        if (rzad === undefined) continue;
        // Waga poważniejsza ma pierwszeństwo: krok z jednym błędem i jednym
        // ostrzeżeniem jest krokiem z błędem, a nie krokiem ostrzeżonym.
        if (rzad.dataset['zastrzezenie'] !== 'blad') {
          rzad.dataset['zastrzezenie'] = zastrzezenie.waga;
        }
        rzad.title = rzad.title === '' ? zastrzezenie.zdanie : `${rzad.title} ${zastrzezenie.zdanie}`;
      }
    },
  };
}

/** Wywołania zwrotne wiersza kroku — porządkowanie i usunięcie. */
interface CzynnosciWiersza {
  wGore(rzad: HTMLElement): void;
  wDol(rzad: HTMLElement): void;
  usun(rzad: HTMLElement): void;
}

/** Buduje jeden wiersz kroku wraz z jego polami i czynnościami. */
function utworzWierszKroku(
  krok: AutomationStep | undefined,
  przyZmianie: () => void,
  czynnosci: CzynnosciWiersza,
): HTMLElement {
  const rzad = document.createElement('li');
  rzad.className = 'da-krok';

  const numer = document.createElement('span');
  numer.className = 'da-krok__numer';

  const identyfikator = pole('Identyfikator kroku', 'np. krok-zbierz');
  identyfikator.value = krok?.id ?? '';
  const nazwa = pole('Nazwa kroku', 'np. Zbierz materiał');
  nazwa.value = krok?.name ?? '';
  const rodzaj = wybor('Rodzaj kroku', RODZAJE_KROKU);
  rodzaj.value = krok?.kind ?? AutomationStepKind.Command;
  const komenda = pole('Komenda kontraktu', 'np. library.file.list');
  komenda.value = krok?.command ?? '';
  const parametry = poleTresci('Treść żądania kroku (JSON)', 3);
  parametry.value = krok?.params === undefined ? '' : JSON.stringify(krok.params);
  const warunek = pole('Warunek wykonania kroku', 'np. wynik.count > 0');
  warunek.value = krok?.condition ?? '';

  for (const kontrolka of [identyfikator, nazwa, komenda, parametry, warunek]) {
    kontrolka.addEventListener('input', przyZmianie);
  }
  rodzaj.addEventListener('change', przyZmianie);

  const wGore = przycisk('↑', 'dn-btn dn-btn--zarys');
  wGore.title = 'Wcześniej w kolejności wykonania';
  wGore.addEventListener('click', () => czynnosci.wGore(rzad));
  const wDol = przycisk('↓', 'dn-btn dn-btn--zarys');
  wDol.title = 'Później w kolejności wykonania';
  wDol.addEventListener('click', () => czynnosci.wDol(rzad));
  const usun = przycisk('Usuń krok', 'dn-btn dn-btn--zarys');
  usun.addEventListener('click', () => czynnosci.usun(rzad));

  const pasek = document.createElement('div');
  pasek.className = 'da-krok__akcje';
  pasek.append(wGore, wDol, usun);

  const pola = document.createElement('div');
  pola.className = 'da-krok__pola';
  pola.append(
    wiersz('Identyfikator', identyfikator, {
      klasa: 'da-wiersz',
      objasnienie: 'Po nim Orchestrator ustala zależności.',
    }),
    wiersz('Nazwa', nazwa, { klasa: 'da-wiersz' }),
    wiersz('Rodzaj', rodzaj, { klasa: 'da-wiersz' }),
    wiersz('Komenda', komenda, {
      klasa: 'da-wiersz',
      objasnienie: 'Komenda kontraktu wywoływana przez krok.',
    }),
    wiersz('Treść żądania', parametry, {
      klasa: 'da-wiersz',
      objasnienie: 'Zapis JSON; puste pole znaczy żądanie bez treści.',
    }),
    wiersz('Warunek', warunek, {
      klasa: 'da-wiersz',
      objasnienie: 'Krok bez warunku wykonuje się zawsze.',
    }),
  );

  rzad.append(numer, pola, pasek);
  return rzad;
}

/** Odczytuje krok z wiersza edytora; kolejność bierze z miejsca w wykazie. */
function odczytajKrok(rzad: HTMLElement, numer: number, poprzednicy?: string[]): AutomationStep {
  const wartosci = [...rzad.querySelectorAll<HTMLInputElement | HTMLTextAreaElement>(
    'input, textarea',
  )];
  const [identyfikator, nazwa, komenda, parametry, warunek] = wartosci;
  const rodzaj = rzad.querySelector<HTMLSelectElement>('select');
  const krok: AutomationStep = {
    id: identyfikator?.value.trim() ?? '',
    kind: (rodzaj?.value ?? AutomationStepKind.Command) as AutomationStep['kind'],
    order: numer + 1,
  };
  if ((nazwa?.value ?? '') !== '') krok.name = nazwa?.value;
  if ((komenda?.value ?? '') !== '') krok.command = komenda?.value;
  if ((warunek?.value ?? '') !== '') krok.condition = warunek?.value;
  if (poprzednicy !== undefined && poprzednicy.length > 0) krok.dependsOn = poprzednicy;
  const tresc = (parametry?.value ?? '').trim();
  if (tresc !== '') krok.params = odczytajTresc(tresc);
  return krok;
}

/**
 * Treść żądania kroku. Zapis nieczytelny jako JSON trafia do rdzenia jako napis,
 * dzięki czemu wpisana treść nie ginie przy zapisie.
 */
function odczytajTresc(tresc: string): unknown {
  try {
    return JSON.parse(tresc);
  } catch {
    return tresc;
  }
}
