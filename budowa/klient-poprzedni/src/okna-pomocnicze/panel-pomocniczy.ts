import type { Kanal } from '../protokol/kanal';

/**
 * Panel pomocniczy jest wspólną postacią wszystkiego, co można otworzyć obok rozmowy: umowa niesie element do osadzenia, odświeżenie i zamknięcie zamykające subskrypcje kanału założone przez panel.
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
