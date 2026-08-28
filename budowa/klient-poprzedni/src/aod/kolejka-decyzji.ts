import { PROGI_AOD } from './progi-aod';
import {
  PRZEDROSTEK_KLUCZA_OKNA,
  rozpoznajDecyzje,
  ZrodloSygnalu,
  type DecyzjaCzekajaca,
  type ObserwacjaProcesu,
} from './rozpoznanie-decyzji';

/** Wpis kolejki decyzji czekających, widziany przez sekcję okna, z czasem oczekiwania liczonym przy odczycie. */
export interface WpisKolejkiDecyzji {
  decyzja: DecyzjaCzekajaca;
  /** Ile milisekund wpis czeka — liczone przy odczycie, nie przechowywane. */
  czekaMs: number;
}

/** Magazyn decyzji czekających, przyjmujący obserwacje procesów i odczyty oraz oddający uszeregowany wykaz. */
export interface KolejkaDecyzji {
  /** Wnosi jedną obserwację; `true`, gdy wykaz się zmienił. */
  nanies(obserwacja: ObserwacjaProcesu, teraz: number): boolean;
  // Odczyt jest rozstrzygający dla wpisów o procesach: wpis bez procesu w odczycie znika.
  nanieOdczyt(obserwacje: readonly ObserwacjaProcesu[], teraz: number): boolean;
  /** Zdejmuje wpis wskazany kluczem; `true`, gdy było co zdjąć. */
  zdejmij(klucz: string): boolean;
  // Wpis odrzucony znika z wykazu i nie wraca do końca życia okna nakładki.
  odrzuc(klucz: string): boolean;
  /** Czy wpis o tym kluczu został odrzucony i nie wróci. */
  czyOdrzucona(klucz: string): boolean;
  // Wykaz uszeregowany od najdłużej czekającego; przeterminowany wpis znika i staje się odrzuconym.
  wykaz(teraz: number): WpisKolejkiDecyzji[];
  /** Liczba decyzji czekających. */
  liczba(): number;
  /** Czy kolejka jest pusta — stan poprawny, nie awaria. */
  pusta(): boolean;
}

export function utworzKolejkeDecyzji(progBezRuchuMs?: number): KolejkaDecyzji {
  /** Wpisy po kluczu; kolejność wynika z `czekaOd`, nie z kolejności wstawiania. */
  const wpisy = new Map<string, DecyzjaCzekajaca>();

  // Klucze odrzucone przez Operatora — analogiczna sugestia nie wraca.
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

  // Wpis przeterminowany nie wraca — klucz idzie do zbioru odrzuconych, a nie tylko z mapy.
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

/** Rozstrzyga, czy dwie decyzje o tym samym kluczu i powodzie różnią się polem rysowanym, a nie całą strukturą. */
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

/** Zdanie „skąd wiadomo" dla całej kolejki, wymieniające tylko sygnały, którymi faktycznie coś złapano. */
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
