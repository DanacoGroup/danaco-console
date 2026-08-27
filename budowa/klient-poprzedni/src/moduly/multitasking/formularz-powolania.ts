import { SubagentStatus, type Subagent, type SubagentSpawnRequest } from '../../../../shared/contract';
import {
  pole,
  poleLiczbowe,
  poleTresci,
  przyciskAkcji as przycisk,
  wiersz,
} from '../../modele/kontrolki-formularza';
import { WIERSZE_POLA } from './kontrolki';
import type { ZrodloPodagentow } from './zrodlo-podagentow';
import './formularz-powolania.css';

// Formularz powołania podagenta woła subagent.spawn; granicę i pola bierze wyłącznie z kontraktu.

/** Górna granica podagentów jednego powołania, ustalona wprost przez kontrakt komendy powołania podagenta. */
export const GRANICA_POWOLANIA = 15;

export interface FormularzPowolania {
  element: HTMLElement;
}

export interface OpcjeFormularzaPowolania {
  /** Źródło `subagent.spawn`. */
  podagenci: ZrodloPodagentow;
  /** Okno wykonawcy powołujące; puste, dopóki obsada go nie założyła. */
  okno(): string;
  potwierdz(zdanie: string, udane: boolean): void;
  /** Wołane po udanym powołaniu — wykaz żywych ma pokazać nowych od razu. */
  poPowolaniu(): void;
}

export function utworzFormularzPowolania(
  opcje: OpcjeFormularzaPowolania,
): FormularzPowolania {
  const zadanie = poleTresci(
    'Zadanie powierzane podagentom',
    WIERSZE_POLA,
    'Co podagent ma wykonać…',
  );
  const nazwa = pole('Nazwa podagenta', 'nazwa (pole nieobowiązkowe)');
  const liczba = poleLiczbowe('Liczba podagentów', '1');
  liczba.min = '1';
  liczba.max = String(GRANICA_POWOLANIA);
  const kanal = pole('Kanał modelu', 'identyfikator kanału (puste = kanał okna)');

  const powolanie = przycisk('Dodaj subagenta', 'dn-btn dn-btn--atrament dn-btn--sm');
  powolanie.addEventListener('click', () => {
    void powolaj();
  });

  const element = document.createElement('div');
  element.className = 'dm-powolanie';
  element.setAttribute('aria-label', 'Powołanie podagenta');
  element.append(
    wiersz('Zadanie', zadanie, { klasa: 'dm-powolanie__pole' }),
    wiersz('Nazwa', nazwa, { klasa: 'dm-powolanie__pole' }),
    wiersz('Liczba', liczba, {
      klasa: 'dm-powolanie__pole',
      objasnienie: `Puste znaczy jednego; kontrakt nie przyjmuje więcej niż ${GRANICA_POWOLANIA}.`,
    }),
    wiersz('Kanał modelu', kanal, {
      klasa: 'dm-powolanie__pole',
      objasnienie: 'Puste znaczy kanał okna wykonawcy — rdzeń rozstrzyga sam.',
    }),
    powolanie,
  );

  /** Powołanie przez `subagent.spawn`; każda przesłanka mówi o sobie zdaniem. */
  async function powolaj(): Promise<void> {
    const okno = opcje.okno();
    if (okno === '') {
      opcje.potwierdz(
        'Podagenta nie ma pod czym powołać: ten wykonawca nie ma jeszcze okna. Załóż okno wykonawcy w obsadzie ról, a formularz zacznie mieć adresata.',
        false,
      );
      return;
    }
    const tresc = zadanie.value.trim();
    if (tresc === '') {
      opcje.potwierdz(
        'Zadanie puste — podagent nie zostaje powołany, bo nie byłoby wiadomo, co ma wykonać.',
        false,
      );
      return;
    }
    const ile = liczbaPowolan(liczba.value);
    if (ile === null) {
      opcje.potwierdz(
        `Liczba podagentów musi być całkowita z przedziału 1–${GRANICA_POWOLANIA}; „${liczba.value}" nią nie jest, więc powołanie nie ruszyło.`,
        false,
      );
      return;
    }

    const wynik = await opcje.podagenci.powolaj(zlozZadanie(okno, tresc, nazwa.value, ile, kanal.value));
    if (!wynik.udany || wynik.wynik === undefined) {
      opcje.potwierdz(
        `Rdzeń odmówił powołania podagenta — ani jeden nie ruszył. Powód: ${wynik.blad?.message ?? 'rdzeń nie podał przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}). Powtórz przyciskiem; gdy odmowa się utrzymuje, komendy subagent.spawn nie ma w tym rdzeniu i potrzebne jest jego rozszerzenie.`,
        false,
      );
      return;
    }

    zadanie.value = '';
    nazwa.value = '';
    opcje.potwierdz(zdanieOPowolaniu(wynik.wynik, ile), wynik.wynik.length > 0);
    opcje.poPowolaniu();
  }

  return { element };
}

/**
 * Żądanie `subagent.spawn` — puste pola nieobowiązkowe nie idą do rdzenia.
 *
 * Puste `name` wysłane jako `''` byłoby nazwą pustą, a nie jej brakiem; rdzeń
 * zapisałby podagenta bez nazwy zamiast nadać mu własną.
 */
function zlozZadanie(
  okno: string,
  zadanie: string,
  nazwa: string,
  ile: number,
  kanal: string,
): SubagentSpawnRequest {
  const tresc: SubagentSpawnRequest = { windowId: okno, task: zadanie };
  const przycieta = nazwa.trim();
  if (przycieta !== '') tresc.name = przycieta;
  if (ile !== 1) tresc.count = ile;
  const kanalPrzyciety = kanal.trim();
  if (kanalPrzyciety !== '') tresc.modelChannelId = kanalPrzyciety;
  return tresc;
}

/** Liczba powołań odczytana z pola formularza; wartość spoza przedziału kontraktu zwraca pustą wartość. */
function liczbaPowolan(zapis: string): number | null {
  const przycieta = zapis.trim();
  if (przycieta === '') return 1;
  const liczba = Number(przycieta);
  if (!Number.isInteger(liczba) || liczba < 1 || liczba > GRANICA_POWOLANIA) return null;
  return liczba;
}

/**
 * Zdanie o powołaniu — liczone z odpowiedzi rdzenia, nie z prośby formularza.
 *
 * Rdzeń ma prawo powołać mniej podagentów, niż zamówiono, więc rozbieżność
 * między prośbą a odpowiedzią wchodzi do zdania, zamiast zostać przemilczana.
 */
function zdanieOPowolaniu(powolani: readonly Subagent[], proszono: number): string {
  if (powolani.length === 0) {
    return `Rdzeń przyjął powołanie, ale nie oddał ani jednego podagenta — mimo prośby o ${proszono}. Sprawdź wykaz żywych; jeśli jest pusty, praca nie ruszyła.`;
  }
  const czynni = powolani.filter(
    (podagent) =>
      podagent.status === SubagentStatus.Pending || podagent.status === SubagentStatus.Running,
  ).length;
  const nazwy = powolani.map((podagent) => podagent.name ?? podagent.id).join(', ');
  const rozbieznosc =
    powolani.length === proszono
      ? ''
      : ` Proszono o ${proszono} — rdzeń oddał ${powolani.length}.`;
  return `Rdzeń powołał ${powolani.length} podagentów (${nazwy}); czynnych: ${czynni}.${rozbieznosc}`;
}
