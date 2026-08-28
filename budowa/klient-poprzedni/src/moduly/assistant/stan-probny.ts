import type { AssistantAction, AssistantActivityEntry } from '../../../../shared/contract';
import type { FazaOkna } from '../../komponenty/faza-okna';
import type { StanAssistant } from './stan-assistant';

// Stanowisko sprawdzianów Assistanta: stan podstawiony ręcznie i trzy odczyty okna.

/**
 * Nastawy stanu modułu, którymi sprawdzian steruje fazą okien: wykazy zleceń
 * własnych i obcych, wpisy dziennika, oznaczenia okna i sesji, zlecenie wybrane
 * oraz znacznik zapytania o okno.
 */
export interface Nastawy {
  zlecenia: readonly AssistantAction[];
  zleceniaObce: readonly AssistantAction[];
  wpisy: readonly AssistantActivityEntry[];
  okno: string;
  sesja: string;
  wybrane: string;
  pytanoOOkno: boolean;
  pytanoOZlecenia: boolean;
  pytanoODziennik: boolean;
  faza: FazaOkna;
  powod: string;
  fazaDziennika: FazaOkna;
  powodDziennika: string;
}

export function nastawyDomyslne(): Nastawy {
  return {
    zlecenia: [],
    zleceniaObce: [],
    wpisy: [],
    okno: '',
    sesja: '',
    wybrane: '',
    pytanoOOkno: false,
    pytanoOZlecenia: false,
    pytanoODziennik: false,
    faza: 'puste',
    powod: '',
    fazaDziennika: 'puste',
    powodDziennika: '',
  };
}

/**
 * Stan modułu na potrzeby sprawdzianu.
 *
 * Okna czytają z niego wyłącznie przez interfejs `StanAssistant`, więc
 * podstawienie jest wymianą źródła, a nie przeróbką okna.
 */
export function stanZ(nastawy: Nastawy): StanAssistant {
  const nieuzywane = (): never => {
    throw new Error('sprawdzian nie sięga tą drogą do rdzenia');
  };
  return {
    zrodlo: {
      polecenie: nieuzywane,
      zlecenia: nieuzywane,
      steruj: nieuzywane,
      dziennik: nieuzywane,
      oznaczWpis: nieuzywane,
      naZmianeZlecenia: () => () => undefined,
    },
    // Trzy źródła dobudowane obok rdzenia modułu; każde ich wywołanie rzuca wyjątkiem.
    mowa: {
      przeslijNagranie: nieuzywane,
      pobierzNagranie: nieuzywane,
      nastawaWybudzania: nieuzywane,
      zapiszWybudzanie: nieuzywane,
      uruchomNasluch: nieuzywane,
      zatrzymajNasluch: nieuzywane,
      naOdcinek: () => () => undefined,
      naWybudzenie: () => () => undefined,
    },
    konteksty: {
      konteksty: nieuzywane,
      zapiszKontekst: nieuzywane,
      uaktywnijKontekst: nieuzywane,
      usunKontekst: nieuzywane,
      zasadyRetencji: nieuzywane,
      zapiszZasadeRetencji: nieuzywane,
      zajetoscKontekstu: nieuzywane,
    },
    schowek: {
      historia: nieuzywane,
      dopisz: nieuzywane,
      przypnij: nieuzywane,
      usun: nieuzywane,
      skroty: nieuzywane,
      zapiszSkrot: nieuzywane,
      usunSkrot: nieuzywane,
      skrotGlobalny: nieuzywane,
      zapiszSkrotGlobalny: nieuzywane,
    },
    zaplecze: {
      okna: nieuzywane,
      akcje: nieuzywane,
      profile: nieuzywane,
      silnikMowy: nieuzywane,
      transkrypcja: nieuzywane,
    },
    zlecenia: () => nastawy.zlecenia,
    zleceniaObce: () => nastawy.zleceniaObce,
    wpisy: () => nastawy.wpisy,
    idOkna: () => nastawy.okno,
    idSesji: () => nastawy.sesja,
    wybrane: () => nastawy.wybrane,
    wybierz: () => Promise.resolve(),
    faza: () => nastawy.faza,
    powod: () => nastawy.powod,
    fazaDziennika: () => nastawy.fazaDziennika,
    powodDziennika: () => nastawy.powodDziennika,
    pytanoOZlecenia: () => nastawy.pytanoOZlecenia,
    pytanoODziennik: () => nastawy.pytanoODziennik,
    pytanoOOkno: () => nastawy.pytanoOOkno,
    ustalOkno: () => Promise.resolve(),
    odswiezZlecenia: () => Promise.resolve(),
    odswiezDziennik: () => Promise.resolve(),
    wchlon: () => undefined,
    obserwuj: () => () => undefined,
    rozlacz: () => undefined,
  };
}

/**
 * Odczytuje fazę z powłoki okna, czyli z elementu, który znakuje `oznaczFaze`.
 * Brak powłoki daje zdanie o jej braku zamiast wartości pustej, żeby sprawdzian
 * odróżnił okno bez fazy od okna bez powłoki.
 */
export function faza(element: HTMLElement): string {
  const powloka = element.querySelector<HTMLElement>('.ma-stan__powloka');
  return powloka?.dataset['faza'] ?? '(brak powłoki stanu)';
}

/**
 * Odczytuje zdanie stanu widoczne obok wskaźnika odczytu. Brak elementu opisu daje
 * napis pusty, ponieważ sprawdzian porównuje treść zdania, a nie obecność samego
 * elementu.
 */
export function zdanie(element: HTMLElement): string {
  return element.querySelector<HTMLElement>('.ma-stan__opis')?.textContent ?? '';
}

/**
 * Odnajduje w oknie przycisk o wskazanej treści. Brak takiego przycisku rzuca
 * wyjątkiem, ponieważ sprawdzian sięgający po nieistniejącą kontrolkę bada okno
 * inne niż to, które zamierzał zbadać.
 */
export function przyciskOTresci(element: HTMLElement, tresc: string): HTMLButtonElement {
  const znaleziony = [...element.querySelectorAll('button')].find((kandydat) =>
    (kandydat.textContent ?? '').includes(tresc),
  );
  if (znaleziony === undefined) throw new Error(`brak przycisku „${tresc}"`);
  return znaleziony;
}
