import type { Kanal } from '../protokol/kanal';
import type { NazwaIkony } from '../ikony/ikony';
import { utworzSekcjeKonto } from './sekcja-konto';
import { utworzSekcjeKontaModeli } from './sekcja-konta-modeli';
import { utworzSekcjePowiadomienia } from './sekcja-powiadomienia';
import { utworzSekcjeUrzadzenia } from './sekcja-urzadzenia';
import { utworzSekcjeUwierzytelnianie } from './sekcja-uwierzytelnianie';
import { utworzSekcjeWyglad } from './sekcja-wyglad';

/** Kod sekcji — sześć pozycji projektu okna ustawień, w kolejności prezentacji na ekranie tego klienta. */
export type KodSekcjiUstawien =
  | 'konto'
  | 'uwierzytelnianie'
  | 'wyglad'
  | 'urzadzenia'
  | 'powiadomienia'
  | 'konta-modeli';

/** Kontrakt wpięcia — kształt, który spełnia każdy plik sekcji, otrzymując kanał raz, przy budowie, i zarządzając odtąd własnym stanem. */
export interface SekcjaUstawien {
  /** Element osadzany w ciele okna, gdy sekcja jest czynna. */
  element: HTMLElement;
  /** Zleca odczyt danych sekcji, wołane przy pierwszym pokazaniu i przy odczycie ponownym w ramie okna. */
  odswiez(): void;
  /** Odłącza subskrypcje kanału. Wołane przy rozłączeniu całego okna. */
  rozlacz(): void;
}

/** Metadane sekcji potrzebne nawigacji: kolejność, nazwa, ikona, oraz sposób zbudowania sekcji z kanału. */
export interface OpisSekcjiUstawien {
  kod: KodSekcjiUstawien;
  nazwa: string;
  ikona: NazwaIkony;
  utworz(kanal: Kanal): SekcjaUstawien;
}

/** Rejestr w kolejności prezentacji stoi tutaj, bo sekcje są stałe i nie przychodzą katalogiem z rdzenia jak kategorie okna konfiguracji. */
export const REJESTR_SEKCJI_USTAWIEN: readonly OpisSekcjiUstawien[] = [
  {
    kod: 'konto',
    nazwa: 'Konto Operatora',
    ikona: 'uzytkownik',
    utworz: () => utworzSekcjeKonto(),
  },
  {
    kod: 'uwierzytelnianie',
    nazwa: 'Uwierzytelnianie',
    ikona: 'tarcza',
    utworz: utworzSekcjeUwierzytelnianie,
  },
  {
    kod: 'wyglad',
    nazwa: 'Wygląd i język',
    ikona: 'paleta',
    utworz: utworzSekcjeWyglad,
  },
  {
    kod: 'urzadzenia',
    nazwa: 'Urządzenia',
    ikona: 'telefon',
    utworz: utworzSekcjeUrzadzenia,
  },
  {
    kod: 'powiadomienia',
    nazwa: 'Powiadomienia',
    ikona: 'dzwonek',
    utworz: utworzSekcjePowiadomienia,
  },
  {
    kod: 'konta-modeli',
    nazwa: 'Konta modeli',
    ikona: 'agenci',
    utworz: () => utworzSekcjeKontaModeli(),
  },
];
