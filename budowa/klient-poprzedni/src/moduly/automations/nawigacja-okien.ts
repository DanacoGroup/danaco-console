import type { NazwaIkony } from '../../ikony/ikony';
import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { KODY_OKIEN, KOD_WYKAZU_PETLI } from './kody-okien';

/**
 * Przywołanie okna modułu — jedno wskazanie zamiast wędrówki po siatce.
 *
 * Nawigacja skraca drogę już istniejącą, nie otwiera nowej: wszystkie sześć
 * okien modułu stoi w jednej przewijanej siatce (`indeks.ts`). Bez niej
 * Operator dochodzi do Execution Monitora przewijaniem, z nią — jednym
 * wskazaniem.
 *
 * Mechanizm jest z biblioteki: rozwijanie, wędrówkę strzałkami, ślad wyboru na
 * gałęzi, pole szukania po progu i zdanie o pustym wykazie niesie
 * `komponenty/menu-drzewo.ts`. Tutaj nie ma własnej listy rozwijanej, nakładki
 * wyboru ani `<select>`.
 *
 * Ujawnianie jest stopniowe: okno wiodące stoi na pierwszym poziomie, bo jest
 * wejściem do modułu, a dwie rodziny ról — kreatory i pas wykonania —
 * odsłaniają się dopiero po wskazaniu gałęzi. Siatka pod menu nie zmienia się:
 * żadne okno nie znika i żadne nie dochodzi, więc przywołanie jest skokiem,
 * nie otwarciem.
 *
 * Na uchwycie stoi bieżąca wartość nastawy, nie jej nazwa rodzajowa. Przed
 * pierwszym przywołaniem żadna wartość nie jest prawdziwa — Operator niczego
 * nie przywołał, a wszystkie okna są na ekranie — więc uchwyt mówi to wprost
 * (`ZDANIE_BEZ_PRZYWOLANIA`) i żaden liść nie niesie wtedy znacznika wyboru.
 */

/** Jedno okno modułu widziane przez nawigację. */
export interface PozycjaOkna {
  /** Kod znacznika `data-okno` w układzie modułu — po nim idzie skok. */
  kod: string;
  /** Nazwa okna — ta sama, co w jego nagłówku (`utworzRameOkna`). */
  nazwa: string;
  /** Zdanie „po co to okno” — objaśnienie `[?]` w postaci pozycji menu. */
  opis: string;
  ikona: NazwaIkony;
}

/** Okno wiodące — wejście do modułu, więc stoi na pierwszym poziomie menu. */
export const OKNO_WIODACE: PozycjaOkna = {
  kod: KOD_WYKAZU_PETLI,
  nazwa: 'Wykaz gotowych pętli',
  opis: 'Pętle zapisane wcześniej, uruchamiane jednym kliknięciem — bez budowania definicji.',
  ikona: 'uruchom',
};

/** Dwa kreatory: budowa struktury automatyki i układ zależności jej kroków. */
export const OKNA_KREATOROW: readonly PozycjaOkna[] = [
  {
    kod: KODY_OKIEN.workflowBuilder,
    nazwa: 'Workflow Builder',
    opis: 'Kroki automatyki, ich rodzaje i zapis definicji w rdzeniu.',
    ikona: 'olowek',
  },
  {
    kod: KODY_OKIEN.orchestrator,
    nazwa: 'Orchestrator',
    opis: 'Zależności między krokami, ocena układu i ścieżka krytyczna.',
    ikona: 'wezly',
  },
];

/** Pas wykonania: dwaj zarządcy i monitor — praca po zbudowaniu automatyki. */
export const OKNA_WYKONANIA: readonly PozycjaOkna[] = [
  {
    kod: KODY_OKIEN.scheduler,
    nazwa: 'Scheduler',
    opis: 'Cykliczność cron, strefa czasowa i wyzwalacze automatyki.',
    ikona: 'zegar',
  },
  {
    kod: KODY_OKIEN.queueManager,
    nazwa: 'Queue Manager',
    opis: 'Kolejka wykonująca automatykę i sześć działań silnika kolejek.',
    ikona: 'warstwy',
  },
  {
    kod: KODY_OKIEN.executionMonitor,
    nazwa: 'Execution Monitor',
    opis: 'Przebiegi na żywo, telemetria etapów i interwencja przy błędzie.',
    ikona: 'monitor',
  },
];

/**
 * Wszystkie okna modułu w kolejności układu — jedno źródło prawdy dla menu
 * i dla sprawdzianu, że menu nie zgubiło ani nie dorobiło pozycji.
 */
export const OKNA_MODULU: readonly PozycjaOkna[] = [
  OKNO_WIODACE,
  ...OKNA_KREATOROW,
  ...OKNA_WYKONANIA,
];

/** Nazwa rodzajowa nastawy — idzie do `aria-label`, nie na ekran. */
export const NASTAWA_PRZYWOLANIA = 'Przywołane okno';

/** Wartość uchwytu, dopóki Operator niczego nie przywołał. Mówi prawdę. */
export const ZDANIE_BEZ_PRZYWOLANIA = 'Bez przywołania — cały układ modułu';

/** Klucze gałęzi; są kluczami menu, nie kodami okien rejestru rdzenia. */
export const GALAZ_KREATORY = 'galaz-kreatory';
export const GALAZ_WYKONANIE = 'galaz-wykonanie';

export interface NawigacjaOkien {
  /** Pas nawigacji osadzany nad układem okien modułu. */
  element: HTMLElement;
  /**
   * Przestawia wskazanie bez skoku.
   *
   * Woła się nią po przywołaniu okna z innej drogi — na przykład skokiem
   * z Workflow Buildera do Orchestratora — żeby uchwyt niósł okno, przy którym
   * Operator naprawdę jest, a nie ostatnie wybrane w menu. Kod spoza wykazu
   * okien modułu jest pomijany: menu nie zaczyna twierdzić, że przywołało coś,
   * czego w module nie ma.
   */
  wskaz(kodOkna: string): void;
  /** Kod okna wskazanego; pusty znaczy „nie przywołano żadnego”. */
  wskazane(): string;
  /** Zwija menu wraz z gałęziami i zdejmuje nasłuchy dokumentu. */
  zwin(): void;
}

export interface OpcjeNawigacji {
  /** Skok do okna wskazanego kodem — `pokaz` układu modułu. */
  przywolaj(kodOkna: string): void;
}

/**
 * Drzewo pozycji dla mechanizmu biblioteki.
 *
 * Czysta funkcja danych: bierze kod okna wskazanego, oddaje wykaz pozycji.
 * Nie zna dokumentu, więc sprawdzian czyta ją wprost, bez montażu menu.
 */
export function drzewoPrzywolania(wskazany: string): PozycjaMenu[] {
  return [
    lisc(OKNO_WIODACE, wskazany),
    {
      rodzaj: 'galaz',
      klucz: GALAZ_KREATORY,
      nazwa: 'Kreatory automatyki',
      opis: 'Okna budujące definicję: kroki i układ zależności między nimi.',
      ikona: 'galaz',
      dzieci: OKNA_KREATOROW.map((okno) => lisc(okno, wskazany)),
    },
    {
      rodzaj: 'galaz',
      klucz: GALAZ_WYKONANIE,
      nazwa: 'Pas wykonania',
      opis: 'Okna uruchamiające i obserwujące: harmonogram, kolejka, przebiegi.',
      ikona: 'aktywnosc',
      dzieci: OKNA_WYKONANIA.map((okno) => lisc(okno, wskazany)),
    },
  ];
}

/** Liść wyboru jednokrotnego — jedno okno modułu. */
function lisc(okno: PozycjaOkna, wskazany: string): PozycjaMenu {
  return {
    rodzaj: 'wybor',
    klucz: okno.kod,
    nazwa: okno.nazwa,
    opis: okno.opis,
    ikona: okno.ikona,
    wybrany: okno.kod === wskazany,
  };
}

/** Nazwa okna o danym kodzie; pusty kod i kod nieznany dają zdanie o braku. */
export function nazwaOkna(kod: string): string {
  return OKNA_MODULU.find((okno) => okno.kod === kod)?.nazwa ?? ZDANIE_BEZ_PRZYWOLANIA;
}

export function utworzNawigacjeOkien(opcje: OpcjeNawigacji): NawigacjaOkien {
  let wskazany = '';

  const element = document.createElement('nav');
  element.className = 'da-modul__nawigacja';
  element.setAttribute('aria-label', 'Okna modułu Automations');

  const menu = utworzMenuDrzewo({
    nastawa: NASTAWA_PRZYWOLANIA,
    ikona: 'automatyzacja',
    // Próg szukania niżej niż domyślne dwanaście: sześć okien mieści się na
    // ekranie bez filtra, a pole nad wykazem sześciu zabrałoby wiersz i nie
    // skróciło ani jednego ruchu.
    progSzukania: 8,
    naWybor: (klucz) => {
      // Mechanizm oddaje klucz pozycji; gałęzie kluczy nie oddają, bo nie są
      // wyborem. Wskazanie idzie przed skokiem, żeby uchwyt niósł nową wartość
      // także wtedy, gdy skok nie znajdzie kafla (układ przebudowany).
      przestaw(klucz);
      opcje.przywolaj(klucz);
    },
  });

  element.append(menu.element);

  /** Zapisuje wskazanie i przerysowuje uchwyt wraz z drzewem. */
  function przestaw(kod: string): void {
    wskazany = OKNA_MODULU.some((okno) => okno.kod === kod) ? kod : '';
    menu.ustaw(nazwaOkna(wskazany), drzewoPrzywolania(wskazany));
  }

  przestaw('');

  return {
    element,
    wskaz: przestaw,
    wskazane: () => wskazany,
    zwin: menu.zwin,
  };
}
