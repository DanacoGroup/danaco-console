import type { ErrorInfo, Queue } from '../../../../shared/contract';
import { pole, przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import type { WykazKomendRdzenia } from './braki-kontraktu';
import { PRZEDROSTEK } from './kontrolki';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloBiegu } from './zrodlo-biegu';

/**
 * Panel planu etapów koordynatora dzieli pracę na kolejki jednego silnika. Etap
 * zakłada komenda `queue.create` z nazwą oraz oknami wykonawców w polu
 * `windowIds`, a wykaz odświeża zdarzenie `queue.changed`.
 */
export interface PlanEtapow {
  element: HTMLElement;
  /** Przerysowuje wykaz etapów po zmianie stanu. */
  odswiez(): void;
}

export interface OpcjePlanu {
  zrodlo: ZrodloBiegu;
  stan: StanMultitaskingu;
  /** Wykaz komend rdzenia — stąd bierze się powód nieczynnej kontrolki. */
  komendy: WykazKomendRdzenia;
  potwierdz(zdanie: string, udane: boolean): void;
}

export function utworzPlanEtapow(opcje: OpcjePlanu): PlanEtapow {
  const { zrodlo, stan } = opcje;

  const nazwa = pole('Nazwa etapu', 'np. Analiza wymagań');
  const dodaj = przycisk('Dodaj etap', 'dn-btn dn-btn--atrament');
  const zaleznosci = opcje.komendy.przyciskNieczynny(
    'Zależności etapów',
    'orchestration.dependency.set',
  );

  dodaj.addEventListener('click', () => {
    void zalozEtap();
  });

  const pasek = document.createElement('div');
  pasek.className = 'dm-plan__pasek';
  pasek.append(nazwa, dodaj, zaleznosci);

  const lista = wykaz('Etapy planu', 'dm-wykaz');

  const element = document.createElement('div');
  element.className = 'dm-plan';
  element.setAttribute('aria-label', 'Plan etapów koordynatora');
  element.append(pasek, lista);

  /** Zakłada etap jako kolejkę obsługującą okna wykonawców tej obsady. */
  async function zalozEtap(): Promise<void> {
    const tytul = nazwa.value.trim();
    if (tytul === '') {
      opcje.potwierdz('Etap bez nazwy nie powstaje — plan ma być czytelny.', false);
      return;
    }
    const okna = stan.obsada().wykonawcy.map((okno) => okno.id);
    const wynik = await zrodlo.zalozKolejke({
      sessionId: stan.sesja(),
      name: tytul,
      ...(okna.length === 0 ? {} : { windowIds: okna }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      opcje.potwierdz(`Rdzeń odmówił założenia etapu. ${powod(wynik.blad)}`, false);
      return;
    }
    nazwa.value = '';
    const zalozony = wynik.wynik;
    stan.zapiszKolejke(zalozony);
    // Nazwa, stan i wykaz okien pochodzą z odpowiedzi rdzenia, a nie z wysłanego żądania.
    opcje.potwierdz(
      `Rdzeń założył kolejkę etapu „${zalozony.name ?? zalozony.id}" (${zalozony.id}) w stanie ${zalozony.status}, okien wykonawców: ${zalozony.windowIds?.length ?? 0} z ${okna.length}.`,
      true,
    );
  }

  function odswiez(): void {
    const etapy = stan.kolejki();
    lista.replaceChildren();
    if (etapy.length === 0) {
      const pusto = document.createElement('li');
      pusto.className = 'dm-pozycja dm-pozycja--pusta';
      pusto.textContent = 'Plan pusty — koordynator nie podzielił jeszcze pracy na etapy.';
      lista.append(pusto);
      return;
    }
    const biezaca = stan.kolejkaBiezaca();
    for (const etap of etapy) lista.append(wierszEtapu(etap, etap.id === biezaca?.id));
  }

  /** Wiersz etapu: nazwa, stan kolejki, licznik obiegów i wybór etapu bieżącego. */
  function wierszEtapu(etap: Queue, biezacy: boolean): HTMLElement {
    const { element: wiersz, akcje } = pozycjaWykazu(
      etap.name ?? etap.id,
      `${etap.status} · obiegów ${etap.cycle ?? 0} · okien ${etap.windowIds?.length ?? 0}`,
      PRZEDROSTEK,
    );
    wiersz.dataset['etap'] = etap.id;
    wiersz.dataset['biezacy'] = String(biezacy);

    if (biezacy) {
      const plakietka = document.createElement('span');
      plakietka.className = 'dn-plakietka';
      plakietka.textContent = 'Etap bieżący';
      akcje.append(plakietka);
      return wiersz;
    }
    const wybierz = przycisk('Ustaw bieżący', 'dn-btn dn-btn--sm dn-btn--zarys');
    wybierz.addEventListener('click', () => stan.ustawKolejke(etap.id));
    akcje.append(wybierz);
    return wiersz;
  }

  odswiez();
  return { element, odswiez };
}

/**
 * Składa treść odmowy rdzenia: podaje komunikat oraz kod błędu z pola `ErrorInfo`
 * kontraktu, a przy odmowie bez opisu mówi wprost, że rdzeń przyczyny nie podał.
 */
function powod(blad?: ErrorInfo): string {
  if (blad === undefined) return 'Rdzeń nie podał przyczyny.';
  return `Powód: ${blad.message} (kod ${blad.code}).`;
}
