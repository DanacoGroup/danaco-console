import {
  type AccessGrant,
  type AccessGrantAddRequest,
  type AccessGrantAddResponse,
  type AccessGrantUpdateRequest,
  type AccessGrantUpdateResponse,
  type AccessMode,
} from '../../../shared/contract';
import type { Wynik } from '../protokol/kanal';
import { odmowaWlasna } from './komunikat-czynnosci';
import type { ZrodloNadan } from './zrodlo-nadan';

/**
 * Trzy komendy zmieniające zbiór nadań okna wraz z uzgodnieniem stanu widoku
 * z odpowiedzią rdzenia.
 *
 * Odczyt i zapis to dwie różne odpowiedzialności: stan sekcji pilnuje tego, co
 * widok wie o punktach i nadaniach, ten plik pilnuje tego, co dzieje się po
 * zapisie.
 *
 * Komplet zastępuje zbiór w całości. Trzy komendy zwracają nie tylko zmienione
 * nadanie, lecz komplet nadań okna po zmianie — przestawienie kolejności albo
 * oznaczenia głównego dotyka pozostałych wierszy, więc widok bierze komplet
 * zamiast składać go z domysłów.
 *
 * Nadanie bez okna rozmowy nie idzie do rdzenia: nadanie żyje per okno, więc
 * bez okna nie ma czego nadać. Odmowa ma kształt wyniku komendy, żeby widok nie
 * potrzebował drugiej ścieżki obsługi.
 */
export interface ZapisyNadan {
  /** `access.grant.add` — nadaje oknu dostęp do punktu. */
  nadaj(
    punktID: string,
    tryb: AccessMode,
    korzenie: readonly string[],
  ): Promise<Wynik<AccessGrantAddResponse>>;
  /** `access.grant.update` — tryb, korzenie, kolejność albo oznaczenie głównego. */
  zmien(zadanie: AccessGrantUpdateRequest): Promise<Wynik<AccessGrantUpdateResponse>>;
  /** `access.grant.remove` — odbiera oknu nadanie. */
  odbierz(nadanieID: string): Promise<Wynik<{ removed: boolean; grants: AccessGrant[] }>>;
}

/** Wejścia zapisów: źródło komend, okno czynne i przyjęcie kompletu nadań. */
export interface ZaleznosciZapisow {
  zrodlo: ZrodloNadan;
  /** Okno rozmowy, którego dotyczy zapis; puste znaczy brak wiązania. */
  oknoID(): string;
  /** Przyjmuje komplet nadań okna z odpowiedzi rdzenia. */
  przyjmijKomplet(komplet: readonly AccessGrant[] | undefined): void;
}

export function utworzZapisyNadan(zaleznosci: ZaleznosciZapisow): ZapisyNadan {
  const { zrodlo, oknoID, przyjmijKomplet } = zaleznosci;

  return {
    async nadaj(punktID, tryb, korzenie) {
      const okno = oknoID();
      if (okno === '') return odmowaWlasna(ODMOWA_BEZ_OKNA);

      const zadanie: AccessGrantAddRequest = {
        windowId: okno,
        accessPointId: punktID,
        mode: tryb,
      };
      if (korzenie.length > 0) zadanie.roots = [...korzenie];

      const wynik = await zrodlo.nadaj(zadanie);
      przyjmijKomplet(wynik.wynik?.grants);
      return wynik;
    },

    async zmien(zadanie) {
      const wynik = await zrodlo.zmien(zadanie);
      przyjmijKomplet(wynik.wynik?.grants);
      return wynik;
    },

    async odbierz(nadanieID) {
      const wynik = await zrodlo.odbierz(nadanieID);
      przyjmijKomplet(wynik.wynik?.grants);
      return wynik;
    },
  };
}

/** Odmowa nadania bez okna: nadanie żyje per okno rozmowy. */
const ODMOWA_BEZ_OKNA =
  'Sekcja nie jest związana z oknem rozmowy — nadanie dostępu żyje per okno, nie per sesja.';
