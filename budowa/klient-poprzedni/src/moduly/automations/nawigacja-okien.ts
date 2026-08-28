import type { NazwaIkony } from '../../ikony/ikony';
import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { KODY_OKIEN, KOD_WYKAZU_PETLI } from './kody-okien';

/** Przywołanie okna modułu jest jednym wskazaniem zamiast wędrówki po siatce sześciu okien. */

/** Jedno okno modułu widziane przez nawigację, niosące kod znacznika, nazwę i objaśnienie widoczne w pozycji menu. */
export interface PozycjaOkna {
  /** Kod znacznika data-okno w układzie modułu — po nim idzie skok. */
  kod: string;
  /** Nazwa okna — ta sama, co w jego nagłówku. */
  nazwa: string;
  /** Zdanie „po co to okno” — objaśnienie w postaci pozycji menu. */
  opis: string;
  ikona: NazwaIkony;
}

/** Okno wiodące jest wejściem do modułu, więc stoi na pierwszym poziomie menu, przed dwiema rodzinami gałęzi. */
export const OKNO_WIODACE: PozycjaOkna = {
  kod: KOD_WYKAZU_PETLI,
  nazwa: 'Wykaz gotowych pętli',
  opis: 'Pętle zapisane wcześniej, uruchamiane jednym kliknięciem — bez budowania definicji.',
  ikona: 'uruchom',
};

/** Dwa kreatory: budowa struktury automatyki i układ zależności jej kroków, odsłaniane po wskazaniu gałęzi kreatorów. */
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

/** Pas wykonania: dwaj zarządcy i monitor, praca po zbudowaniu automatyki, odsłaniana po wskazaniu gałęzi wykonania. */
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

/** Wszystkie okna modułu w kolejności układu są jednym źródłem prawdy dla menu i dla sprawdzianu pokrycia pozycji. */
export const OKNA_MODULU: readonly PozycjaOkna[] = [
  OKNO_WIODACE,
  ...OKNA_KREATOROW,
  ...OKNA_WYKONANIA,
];

/** Nazwa rodzajowa nastawy idzie do atrybutu dostępności menu, nie pokazuje się nigdy na ekranie modułu. */
export const NASTAWA_PRZYWOLANIA = 'Przywołane okno';

/** Wartość uchwytu, dopóki Operator niczego nie przywołał, mówi prawdę: cały układ modułu jest wtedy na ekranie. */
export const ZDANIE_BEZ_PRZYWOLANIA = 'Bez przywołania — cały układ modułu';

/** Klucze gałęzi kreatorów i wykonania są kluczami samego menu, nie kodami okien rejestru rdzenia platformy. */
export const GALAZ_KREATORY = 'galaz-kreatory';
export const GALAZ_WYKONANIE = 'galaz-wykonanie';

export interface NawigacjaOkien {
  /** Pas nawigacji osadzany nad układem okien modułu. */
  element: HTMLElement;
  /** Przestawia wskazanie bez skoku; woła się po przywołaniu okna z innej drogi niż menu. */
  wskaz(kodOkna: string): void;
  /** Kod okna wskazanego; pusty znaczy, że nie przywołano żadnego. */
  wskazane(): string;
  /** Zwija menu wraz z gałęziami i zdejmuje nasłuchy dokumentu. */
  zwin(): void;
}

export interface OpcjeNawigacji {
  /** Skok do okna wskazanego kodem, w układzie modułu. */
  przywolaj(kodOkna: string): void;
}

/** Drzewo pozycji dla mechanizmu biblioteki jest czystą funkcją danych, bierze kod okna, oddaje wykaz pozycji. */
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

/** Liść wyboru jednokrotnego reprezentuje jedno okno modułu w drzewie menu nawigacji przywołania okien. */
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

/** Nazwa okna o danym kodzie; pusty kod i kod nieznany dają wspólne zdanie o braku żadnego przywołania. */
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
    // Próg szukania niżej niż domyślne dwanaście: sześć okien mieści się na ekranie bez filtra.
    progSzukania: 8,
    naWybor: (klucz) => {
      // Wskazanie idzie przed skokiem, żeby uchwyt niósł nową wartość nawet po przebudowie układu.
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
