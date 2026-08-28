import type { Kanal } from '../protokol/kanal';
import { utworzOknoKonfiguracji, type OknoKonfiguracji } from './okno-konfiguracji';

/**
 * Punkt zbiorczy katalogu konfiguracji: reszta aplikacji otwiera stąd jedyne
 * okno ustawień klienta. Okno powstaje przy pierwszym otwarciu i trwa między
 * kolejnymi, a wartość `null` obowiązuje, dopóki żadne otwarcie nie nastąpiło.
 */
let okno: OknoKonfiguracji | null = null;

/**
 * Kanał, na którym zbudowano okno. Porównanie z kanałem podanym przy otwarciu
 * rozpoznaje ponowne połączenie z rdzeniem i nakazuje przebudowę okna, zanim
 * pierwsza komenda pójdzie przez transport już nieczynny.
 */
let osadzonyKanal: Kanal | null = null;

/**
 * Otwiera okno konfiguracji i zleca wczytanie katalogu ustawień. Otwarcie nie
 * czeka na rdzeń: okno pojawia się od razu, a katalog dociera do niego osobną
 * odpowiedzią na zlecony odczyt.
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
