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

/**
 * Droga klienta do rdzenia zebrana w jednym obiekcie: adres gniazda, transport
 * ramek, kanał komunikatów kontraktu, sesja nadana przez rdzeń, uzgodnienie
 * oraz nawigacja platformy podawana powłoce.
 */
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
  /** Nawigacja platformy dla bocznej kolumny powłoki, osadzona na tym kanale. */
  platforma: NawigacjaPlatformy;
}

/**
 * Złożenie warstw łączności w jedną drogę do rdzenia: adres, transport, sesja,
 * kanał, uzgodnienie i nawigacja platformy. Plik wyłącznie składa i nie otwiera
 * połączenia; rozpoczęcie łączności należy do cyklu życia.
 */
export function zlozPolaczenieZRdzeniem(opis: OpisOkna): PolaczenieZRdzeniem {
  const adres = adresRdzenia();
  const transport = utworzTransport(adres);
  const sesja = utworzSesje();
  const kanal = utworzKanal(transport, sesja);
  // Token sesji bramki czytany przy każdym powitaniu, żeby po zerwaniu
  // obowiązywała sesja bieżąca.
  const uzgodnienie = utworzUzgodnienie(kanal, zamowienieOkna(opis), tozsamoscKlienta(), () =>
    odczytajSesje()?.token,
  );

  // Nawigacja powłoki składa się tutaj, bo uzgodnienie odpowiada tylko za
  // powitanie, sesję i okno.
  const platforma: NawigacjaPlatformy = {
    wejdzDoSrodowiska: (zadanie) => zadajWejscieDoSrodowiska(kanal, zadanie),
    wykazModulow: (zadanie = {}) => zadajWykazModulow(kanal, zadanie),
  };

  return { adres, transport, kanal, sesja, uzgodnienie, platforma };
}
