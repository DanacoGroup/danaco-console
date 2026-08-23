import type { Kanal } from '../protokol/kanal';

/**
 * Panel pomocniczy — wspólna postać wszystkiego, co można otworzyć obok rozmowy.
 *
 * Umowa niesie dokładnie to, czego potrzebuje gospodarz panelu — pas modułu
 * albo kolumna paneli sceny: element do osadzenia, odświeżenie i zamknięcie.
 * Tytuł, czynności i stany treści należą do samego panelu.
 *
 * `zamknij()` jest wymagane, bo panel trzyma subskrypcje kanału
 * (`okno-podglad-bash.ts` słucha `stream.chunk`) i po zejściu ze sceny
 * zostawiłby je żywe.
 *
 * Plik nie zna kodów paneli i niczego nie buduje: mapowanie kodu na wytwórnię
 * stoi w `wytwornia-paneli.ts`, a spis opisowy — w `rejestr-pomocniczych.ts`.
 *
 * Umowa nie wymaga ramy okna z `komponenty/rama-okna.ts`; panel wolno w nią
 * ubrać — `okno-podglad-bash.ts` tak robi — ale oddaje wyłącznie element.
 */
export interface PanelPomocniczy {
  /** Element osadzany przez gospodarza — pas modułu albo kolumna paneli sceny. */
  element: HTMLElement;
  /** Ponawia odczyt panelu z rdzenia. */
  odswiez(): void;
  /** Zamyka subskrypcje. Obowiązkowe — panel bez tego zostawia je żywe. */
  zamknij(): void;
}

/**
 * Zależności, które każdy panel dostaje od gospodarza.
 *
 * Skład odpowiada pole w pole `OpcjePodgladuBash` — jedynemu panelowi, który
 * dziś istnieje. Pole, którego żaden panel nie czyta, byłoby atrapą umowy.
 */
export interface OpcjePanelu {
  kanal: Kanal;
  /** Okno wykonania gospodarza; pusty napis = rdzeń nie dał okna. */
  okno: string;
  /** Nazwa modułu w etykietach dostępności — `Developer`, `Diagnostics`. */
  modul: string;
  /** Przedrostek klas modułu: `mdev`, `dg`. */
  przedrostek: string;
}
