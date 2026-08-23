import type { Window } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import {
  utworzStanSterowania,
  type MigawkaSterowania,
} from '../sterowanie/stan-sterowania';
import { utworzZmianeOkna } from '../sterowanie/zmiana-okna';
import { utworzZmianeUstawienia } from '../sterowanie/zmiana-ustawienia';

/**
 * Odczyt bieżących ustawień okna na potrzeby podsumowania widoku.
 *
 * Podsumowanie pokazuje osiem wartości obok siebie także wtedy, gdy szuflada
 * z kontrolkami jest zwinięta — a kontrolki kompletu `sterowanie/` trzymają
 * swój stan wewnątrz `utworzPanelSterowania` i nie wystawiają go na zewnątrz
 * (`PanelSterowania` niesie wyłącznie `element`, `idOkna`, `przyjmijOkno`
 * i `rozlacz`). Ten plik jest odczytem tego samego stanu, a nie drugą jego
 * definicją: buduje `StanSterowania`, `ZmianaOkna` i `ZmianaUstawienia`
 * z katalogu `sterowanie/` i nie dopisuje ani jednej reguły protokołu.
 *
 * Obserwator wyłącznie czyta — nie wywołuje `zastosuj` ani `zapisz`, więc
 * nie ma drogi, którą mógłby zmienić okno. Komunikaty o losie zmian pomija,
 * bo pokazuje je pasek kompletu sterowania; podwójny komunikat byłby szumem.
 *
 * Gdy `PanelSterowania` wystawi `migawka()` i `naZmiane()`, ten plik zniknie,
 * a podsumowanie odczyta stan kompletu wprost.
 */
export interface ObserwatorUstawien {
  /** Bieżąca migawka ustawień okna. */
  migawka(): MigawkaSterowania;
  /** Subskrypcja zmian potwierdzonych przez rdzeń. */
  naZmiane(sluchacz: (migawka: MigawkaSterowania) => void): Odsubskrybuj;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

/** Odbiorca komunikatów o losie zmiany; obserwator nic nie zmienia, więc go pomija. */
const POMIN_KOMUNIKAT = (): void => {};

export function utworzObserwatorUstawien(
  kanal: Kanal,
  okno: Window,
): ObserwatorUstawien {
  const stan = utworzStanSterowania(okno);

  // Obie drogi zmiany podpinają subskrypcje zdarzeń `window.changed`
  // i `config.changed` filtrowane po identyfikatorze okna.
  const zmianaOkna = utworzZmianeOkna(kanal, stan, POMIN_KOMUNIKAT);
  const zmianaUstawienia = utworzZmianeUstawienia(kanal, stan, POMIN_KOMUNIKAT);

  // Odczyt ustawień poziomu okna: host wykonania, kanał zapasowy
  // i nakład rozumowania nie mają pola w `window.update`, więc bez tego
  // wywołania podsumowanie pokazywałoby wartości domyślne zamiast zapisanych.
  zmianaUstawienia.wczytaj();

  return {
    migawka: stan.migawka,
    naZmiane: stan.naZmiane,
    rozlacz() {
      zmianaOkna.rozlacz();
      zmianaUstawienia.rozlacz();
    },
  };
}
