import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { zamowienieOkna } from '../okno-komunikacji/zamowienie-okna';
import { adresRdzenia } from '../polaczenie/adres-rdzenia';
import { utworzTransport, type Transport } from '../polaczenie/gniazdo';
import { utworzKanal, type Kanal } from '../protokol/kanal';
import type { NawigacjaPlatformy } from '../protokol/nawigacja-platformy';
import { utworzSesje, type Sesja } from '../protokol/sesja';
import { zadajWejscieDoSrodowiska } from '../protokol/wejscie-do-srodowiska';
import { zadajWykazModulow } from '../protokol/wykaz-modulow';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import { utworzUzgodnienie, type Uzgodnienie } from '../protokol/uzgodnienie';
import { odczytajSesje } from '../uwierzytelnienie/sesja-bramki';

/** Droga klienta do rdzenia: transport, kanał kontraktu, uzgodnienie. */
export interface PolaczenieZRdzeniem {
  /** Adres gniazda rdzenia użyty przy złożeniu. */
  adres: string;
  /** Transport ramek WebSocket. */
  transport: Transport;
  /** Kanał komunikatów kontraktu osadzony na transporcie. */
  kanal: Kanal;
  /** Sesja klienta; identyfikator nadaje rdzeń. */
  sesja: Sesja;
  /** Uzgodnienie: powitanie, sesja, okno komunikacji. */
  uzgodnienie: Uzgodnienie;
  /**
   * Nawigacja platformy dla bocznej kolumny powłoki: `environment.enter`
   * i `module.list` osadzone na tym kanale.
   *
   * Pole służy wyłącznie powłoce: powstaje ona bez połączenia i nie ma jak
   * sięgnąć po kanał sama, więc dostaje te dwie drogi gotowe. Widoki mające
   * kanał wołają opakowania wprost — `home.enter` widok strony głównej,
   * `workspace.enter` przestrzeń modułu — bez obiektu pośredniego.
   */
  platforma: NawigacjaPlatformy;
}

/**
 * Złożenie warstw łączności w jedną drogę do rdzenia.
 *
 * Plik wyłącznie składa — nie zna ramki, nie buduje koperty i nie zna nazwy
 * żadnej komendy. Nazwy pochodzą z pakietu `shared` i żyją w warstwie
 * protokołu.
 *
 * Złożenie nie otwiera połączenia. Rozpoczęcie łączności należy do cyklu
 * życia, żeby moment jej nawiązania był jednym miejscem, a nie skutkiem
 * ubocznym budowy obiektów.
 */
export function zlozPolaczenieZRdzeniem(opis: OpisOkna): PolaczenieZRdzeniem {
  const adres = adresRdzenia();
  const transport = utworzTransport(adres);
  const sesja = utworzSesje();
  const kanal = utworzKanal(transport, sesja);
  // Token sesji bramki czytany przy każdym powitaniu, nie raz przy składaniu:
  // po ponownym nawiązaniu połączenia obowiązuje sesja bieżąca, nie ta sprzed
  // zerwania.
  const uzgodnienie = utworzUzgodnienie(kanal, zamowienieOkna(opis), tozsamoscKlienta(), () =>
    odczytajSesje()?.token,
  );

  // Nawigacja powłoki składa się tutaj, a nie w uzgodnieniu: uzgodnienie
  // odpowiada wyłącznie za powitanie, sesję i okno, i o nawigacji nic nie wie.
  // Korzeń montażu klienta jest jedynym punktem znającym jednocześnie kanał
  // i odbiorcę tych dróg.
  const platforma: NawigacjaPlatformy = {
    wejdzDoSrodowiska: (zadanie) => zadajWejscieDoSrodowiska(kanal, zadanie),
    wykazModulow: (zadanie = {}) => zadajWykazModulow(kanal, zadanie),
  };

  return { adres, transport, kanal, sesja, uzgodnienie, platforma };
}
