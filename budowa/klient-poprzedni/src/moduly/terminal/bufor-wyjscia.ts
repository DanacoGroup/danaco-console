import type { ChunkKind } from '../../../../shared/contract';

/**
 * Bufor pierścieniowy wyjścia procesów dla Output Console: znosi duże objętości
 * bez zamulenia okna.
 *
 * Pierścień zamiast rosnącej tablicy, bo jedno `npm install` albo jedno `find /`
 * daje setki tysięcy wierszy — tablica bez granicy rośnie do wyczerpania pamięci
 * karty, a odrysowanie jej za każdym fragmentem zatrzymuje wątek interfejsu.
 * Bufor trzyma stałą liczbę wierszy najnowszych i podaje, ile wierszy przepadło,
 * żeby było widać, że w oknie stoi ogon, a nie całość.
 *
 * Fragment to nie wiersz: rdzeń wysyła porcje bajtów, więc jeden fragment bywa
 * połową wiersza, a inny trzema wierszami naraz. Bufor skleja ogon
 * niedokończonego wiersza z początkiem następnego fragmentu, inaczej wyjście
 * rozsypałoby się na przypadkowych granicach odczytu.
 */
export interface WierszWyjscia {
  /** Proces rejestru rdzenia, z którego pochodzi wiersz. */
  proces: string;
  /** Karta powłoki procesu; pusta, gdy proces przyszedł spoza znanych kart. */
  karta: string;
  /** Rodzaj fragmentu kontraktu — `text` albo `error`. */
  rodzaj: ChunkKind;
  /** Treść wiersza bez znaku końca linii. */
  tresc: string;
  /** Czas przyjęcia w milisekundach epoki. */
  chwila: number;
}

export interface BuforWyjscia {
  /** Dokłada porcję wyjścia procesu; zwraca liczbę nowych wierszy. */
  dopisz(proces: string, karta: string, rodzaj: ChunkKind, tresc: string, chwila: number): number;
  /** Domyka niedokończony wiersz procesu — wywołuje się po jego zakończeniu. */
  domknij(proces: string): void;
  /** Wiersze bufora w kolejności przyjęcia. */
  wiersze(): readonly WierszWyjscia[];
  /** Liczba wierszy, które wypadły z pierścienia. */
  utracone(): number;
  /** Czyści bufor wraz z licznikiem utraconych wierszy. */
  wyczysc(): void;
}

/**
 * Pojemność bufora. Wartość dobrana tak, żeby przewijanie historii miało sens
 * (kilkadziesiąt ekranów), a odrysowanie całości pozostało wykonalne.
 */
export const POJEMNOSC_BUFORA = 5000;

export function utworzBuforWyjscia(pojemnosc = POJEMNOSC_BUFORA): BuforWyjscia {
  const wiersze: WierszWyjscia[] = [];
  const ogony = new Map<string, string>();
  let utracone = 0;

  function dodaj(wiersz: WierszWyjscia): void {
    wiersze.push(wiersz);
    if (wiersze.length > pojemnosc) {
      wiersze.splice(0, wiersze.length - pojemnosc);
      utracone += 1;
    }
  }

  return {
    dopisz(proces, karta, rodzaj, tresc, chwila) {
      const pelna = (ogony.get(proces) ?? '') + tresc.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
      const czesci = pelna.split('\n');
      // Ostatnia część nie ma znaku końca linii — zostaje ogonem do sklejenia
      // z następnym fragmentem tego samego procesu.
      ogony.set(proces, czesci.pop() ?? '');
      for (const czesc of czesci) dodaj({ proces, karta, rodzaj, tresc: czesc, chwila });
      return czesci.length;
    },

    domknij(proces) {
      const ogon = ogony.get(proces);
      ogony.delete(proces);
      if (ogon === undefined || ogon === '') return;
      const ostatni = wiersze[wiersze.length - 1];
      dodaj({
        proces,
        karta: ostatni?.karta ?? '',
        rodzaj: ostatni?.rodzaj ?? 'text',
        tresc: ogon,
        chwila: Date.now(),
      });
    },

    wiersze: () => wiersze,

    utracone: () => utracone,

    wyczysc() {
      wiersze.length = 0;
      ogony.clear();
      utracone = 0;
    },
  };
}
