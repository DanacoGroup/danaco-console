import { PROGI_AOD } from './progi-aod';
import {
  PRZEDROSTEK_KLUCZA_OKNA,
  rozpoznajDecyzje,
  ZrodloSygnalu,
  type DecyzjaCzekajaca,
  type ObserwacjaProcesu,
} from './rozpoznanie-decyzji';

/**
 * Kolejka decyzji czekających — magazyn bez DOM i bez kanału. Trzyma, co czeka,
 * od kiedy, czego dotyczy i skąd o tym wiadomo; drogi wyjścia zostają poza
 * magazynem, bo zależą od stanu rdzenia w chwili czynności (`cztery-stery.ts`).
 *
 * Dosypują dwa źródła: odczyt nadrabiający `monitor.status` wnosi to, co stanęło
 * przed otwarciem okna, a sygnały `progress.changed` i `window.state.changed` —
 * to, co staje przy otwartym oknie.
 *
 * Powtórzenie tego samego powodu nie zeruje `czekaOd`; zegar rusza od nowa
 * dopiero przy zmianie powodu. Kluczem jest identyfikator procesu, a gdy sygnał
 * go nie niesie (`window.state.changed`) — okno z przedrostkiem `okno:`. Sygnał
 * o biegnącym procesie zdejmuje wpis natychmiast. Wykaz idzie po `czekaOd`
 * rosnąco, a przy równych chwilach rozstrzyga klucz, żeby nie migotał.
 */

/** Wpis kolejki widziany przez sekcję okna. */
export interface WpisKolejkiDecyzji {
  decyzja: DecyzjaCzekajaca;
  /** Ile milisekund wpis czeka — liczone przy odczycie, nie przechowywane. */
  czekaMs: number;
}

/** Magazyn decyzji czekających. */
export interface KolejkaDecyzji {
  /**
   * Wnosi jedną obserwację.
   *
   * @returns `true`, gdy wykaz się zmienił — sekcja przerysowuje się tylko wtedy.
   */
  nanies(obserwacja: ObserwacjaProcesu, teraz: number): boolean;
  /**
   * Wnosi cały odczyt nadrabiający `monitor.status`.
   *
   * Odczyt jest ROZSTRZYGAJĄCY dla wpisów o procesach: wpis, którego proces
   * w odczycie nie wystąpił, znika — rdzeń już tego procesu nie zna. Wpisy
   * kluczowane oknem zostają, bo `monitor.status` o nich nie mówi.
   */
  nanieOdczyt(obserwacje: readonly ObserwacjaProcesu[], teraz: number): boolean;
  /** Zdejmuje wpis wskazany kluczem; `true`, gdy było co zdjąć. */
  zdejmij(klucz: string): boolean;
  /**
   * Odrzuca wpis — status `odrzucona` cyklu życia sugestii (rozdz. 9.4).
   *
   * Wpis znika z wykazu i NIE WRACA: analogiczne zdarzenie nie tworzy sugestii
   * do końca życia okna nakładki (rozdz. 4.1: „do końca karty sesji").
   * Odrzucenie idzie bez pytania „czy na pewno" — decyzja Operatora jest
   * decyzją, nie propozycją.
   */
  odrzuc(klucz: string): boolean;
  /** Czy wpis o tym kluczu został odrzucony i nie wróci. */
  czyOdrzucona(klucz: string): boolean;
  /**
   * Wykaz uszeregowany od najdłużej czekającego.
   *
   * Wpis czekający dłużej niż czas życia sugestii nieprzyjętej (rozdz. 3.4 —
   * 24 godziny) znika z wykazu; opracowanie nadaje mu wtedy status odrzuconej.
   */
  wykaz(teraz: number): WpisKolejkiDecyzji[];
  /** Liczba decyzji czekających. */
  liczba(): number;
  /** Czy kolejka jest pusta — stan poprawny, nie awaria. */
  pusta(): boolean;
}

export function utworzKolejkeDecyzji(progBezRuchuMs?: number): KolejkaDecyzji {
  /** Wpisy po kluczu; kolejność wynika z `czekaOd`, nie z kolejności wstawiania. */
  const wpisy = new Map<string, DecyzjaCzekajaca>();

  /** Klucze odrzucone przez Operatora — sugestia nie wraca (rozdz. 4, „Odrzuć"). */
  const odrzucone = new Set<string>();

  function nanies(obserwacja: ObserwacjaProcesu, teraz: number): boolean {
    if (odrzucone.has(obserwacja.klucz)) return false;

    const rozpoznana = rozpoznajDecyzje(obserwacja, teraz, progBezRuchuMs);

    if (rozpoznana === null) {
      // Proces ruszył albo skończył — wpis przestaje być prawdą.
      return wpisy.delete(obserwacja.klucz);
    }

    const poprzednia = wpisy.get(obserwacja.klucz);
    if (poprzednia !== undefined && poprzednia.powod === rozpoznana.powod) {
      const scalona: DecyzjaCzekajaca = { ...rozpoznana, czekaOd: poprzednia.czekaOd };
      wpisy.set(obserwacja.klucz, scalona);
      return czyRozne(poprzednia, scalona);
    }

    wpisy.set(obserwacja.klucz, rozpoznana);
    return true;
  }

  /**
   * Zdejmuje wpisy starsze niż czas życia sugestii nieprzyjętej (rozdz. 3.4).
   *
   * Wpis nie wraca po przeterminowaniu — opracowanie nadaje mu status
   * odrzuconej — więc klucz idzie do zbioru odrzuconych, a nie tylko z mapy.
   */
  function przeterminuj(teraz: number): boolean {
    let zmiana = false;
    for (const [klucz, decyzja] of wpisy) {
      if (teraz - decyzja.czekaOd < PROGI_AOD.czasZyciaSugestiiMs) continue;
      wpisy.delete(klucz);
      odrzucone.add(klucz);
      zmiana = true;
    }
    return zmiana;
  }

  return {
    nanies,

    nanieOdczyt(obserwacje, teraz) {
      let zmiana = false;

      const widziane = new Set<string>();
      for (const obserwacja of obserwacje) {
        widziane.add(obserwacja.klucz);
        if (nanies(obserwacja, teraz)) zmiana = true;
      }

      for (const klucz of [...wpisy.keys()]) {
        if (klucz.startsWith(PRZEDROSTEK_KLUCZA_OKNA)) continue;
        if (widziane.has(klucz)) continue;
        wpisy.delete(klucz);
        zmiana = true;
      }

      return zmiana;
    },

    zdejmij: (klucz) => wpisy.delete(klucz),

    odrzuc(klucz) {
      odrzucone.add(klucz);
      return wpisy.delete(klucz);
    },

    czyOdrzucona: (klucz) => odrzucone.has(klucz),

    wykaz(teraz) {
      przeterminuj(teraz);
      return [...wpisy.values()]
        .sort((a, b) => a.czekaOd - b.czekaOd || a.klucz.localeCompare(b.klucz))
        .map((decyzja) => ({ decyzja, czekaMs: Math.max(0, teraz - decyzja.czekaOd) }));
    },

    liczba: () => wpisy.size,

    pusta: () => wpisy.size === 0,
  };
}

/**
 * Czy dwie decyzje o tym samym kluczu i powodzie różnią się czymś widocznym.
 *
 * Porównujemy pola RYSOWANE, a nie całe struktury: `monitor.status` przychodzi
 * co odczyt i przy identycznej treści nie ma powodu przerysowywać wykazu —
 * przerysowanie gubiłoby ognisko klawiatury na sterach.
 */
function czyRozne(a: DecyzjaCzekajaca, b: DecyzjaCzekajaca): boolean {
  return (
    a.zdanie !== b.zdanie ||
    a.stan !== b.stan ||
    a.etap !== b.etap ||
    a.ukonczenie !== b.ukonczenie ||
    a.idOkna !== b.idOkna ||
    a.idSesji !== b.idSesji ||
    a.zrodlo !== b.zrodlo ||
    a.powodZatrzymaniaBiegu !== b.powodZatrzymaniaBiegu
  );
}

/**
 * Zdanie „skąd wiadomo" dla całej kolejki — jedna linia pod wykazem.
 *
 * Wymienia sygnały, którymi kolejka faktycznie coś złapała, a nie te, które
 * subskrybuje.
 */
export function opiszZrodlaKolejki(wpisy: readonly WpisKolejkiDecyzji[]): string {
  const zrodla = new Set(wpisy.map((wpis) => wpis.decyzja.zrodlo));
  if (zrodla.size === 0) return '';

  const kolejnosc: ZrodloSygnalu[] = [
    ZrodloSygnalu.Monitor,
    ZrodloSygnalu.Postep,
    ZrodloSygnalu.StanOkna,
  ];
  return kolejnosc.filter((zrodlo) => zrodla.has(zrodlo)).join(', ');
}
