import type { Kanal } from '../protokol/kanal';
import { utworzOknoKonfiguracji, type OknoKonfiguracji } from './okno-konfiguracji';

/**
 * Okno konfiguracji — punkt zbiorczy katalogu.
 *
 * Reszta aplikacji zna stąd jedną czynność:
 *
 *   import { otworzOknoKonfiguracji } from '../konfiguracja/indeks';
 *   otworzOknoKonfiguracji(rdzen.kanal);
 *
 * Okno jest jedno na klienta i żyje między otwarciami. Powód jest ten sam,
 * dla którego router nie porzuca widoku opuszczonej trasy: powtórne otwarcie
 * wraca do kategorii, na której się skończyło, a nie zaczyna od początku.
 * Przy okazji nie gubi się subskrypcja zdarzenia `config.changed`,
 * więc zapis dokonany gdzie indziej dociera także wtedy, gdy okno jest
 * zamknięte.
 *
 * Kanał podajemy przy pierwszym otwarciu. Wywołanie z innym kanałem —
 * po ponownym połączeniu z rdzeniem — buduje okno na nowo, żeby komendy nie
 * szły przez transport, którego już nie ma.
 */

/** Okno zbudowane przy pierwszym otwarciu; `null` przed nim. */
let okno: OknoKonfiguracji | null = null;

/** Kanał, na którym zbudowano okno — podstawa rozpoznania zmiany połączenia. */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera okno konfiguracji i zleca wczytanie katalogu.
 *
 * Otwarcie nie czeka na rdzeń: okno pojawia się od razu, a katalog dojeżdża
 * do niego odpowiedzią.
 */
export function otworzOknoKonfiguracji(kanal: Kanal): OknoKonfiguracji {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoKonfiguracji(kanal);
    osadzonyKanal = kanal;
  }

  okno.otworz();
  return okno;
}

export type { OknoKonfiguracji } from './okno-konfiguracji';
export type { StanKonfiguracji } from './stan-konfiguracji';
export type { PunktWidzenia, Rozstrzygniecie } from './rozstrzygniecie';
export type { AdresUstawienia } from './adres-ustawienia';
