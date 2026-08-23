import type { AssistantAction, AssistantActivityEntry } from '../../../../shared/contract';
import type { FazaOkna } from '../../komponenty/faza-okna';
import type { StanAssistant } from './stan-assistant';

/**
 * Stanowisko sprawdzianów modułu Assistant: stan modułu podstawiony ręcznie
 * oraz trzy odczyty złożonego okna.
 *
 * Sprawdzianów modułu są dwa — stany okien i wykonywanie pracy — a atrapa
 * `StanAssistant` jest jedna. Skopiowana do obu plików rozjeżdżałaby się przy
 * każdym nowym polu interfejsu i jeden ze sprawdzianów badałby wtedy stan,
 * którego moduł już nie ma.
 *
 * Plik nie należy do produktu: sięgają po niego wyłącznie `*.test.ts`, więc
 * `main.ts` go nie wciąga.
 *
 * Atrapa nie sięga do rdzenia i nie udaje, że sięga: każda droga wywołania
 * rzuca wyjątkiem, dopóki sprawdzian jej nie obsadzi. Sprawdzian, który
 * przypadkiem wywoła komendę, dostaje przez to błąd, a nie ciszę.
 */

/** Nastawy stanu modułu, którymi sprawdzian steruje fazą okien. */
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
    // Trzy źródła dobudowane obok rdzenia modułu. Sprawdzian nie sięga nimi do
    // rdzenia, więc każde wywołanie od razu mówi, że tędy droga nie prowadzi —
    // atrapa oddająca pustą odpowiedź udawałaby wynik, którego nie ma.
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

/** Faza odczytana z powłoki okna — tej samej, którą znakuje `oznaczFaze`. */
export function faza(element: HTMLElement): string {
  const powloka = element.querySelector<HTMLElement>('.ma-stan__powloka');
  return powloka?.dataset['faza'] ?? '(brak powłoki stanu)';
}

/** Zdanie stanu widoczne obok wskaźnika odczytu. */
export function zdanie(element: HTMLElement): string {
  return element.querySelector<HTMLElement>('.ma-stan__opis')?.textContent ?? '';
}

/** Przycisk o wskazanej treści; brak przycisku jest błędem sprawdzianu. */
export function przyciskOTresci(element: HTMLElement, tresc: string): HTMLButtonElement {
  const znaleziony = [...element.querySelectorAll('button')].find((kandydat) =>
    (kandydat.textContent ?? '').includes(tresc),
  );
  if (znaleziony === undefined) throw new Error(`brak przycisku „${tresc}"`);
  return znaleziony;
}
