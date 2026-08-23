import type { Kanal } from '../protokol/kanal';
import type { NazwaIkony } from '../ikony/ikony';
import { utworzSekcjeKonto } from './sekcja-konto';
import { utworzSekcjeKontaModeli } from './sekcja-konta-modeli';
import { utworzSekcjePowiadomienia } from './sekcja-powiadomienia';
import { utworzSekcjeUrzadzenia } from './sekcja-urzadzenia';
import { utworzSekcjeUwierzytelnianie } from './sekcja-uwierzytelnianie';
import { utworzSekcjeWyglad } from './sekcja-wyglad';

/**
 * Rejestr sekcji Okna Ustawień oraz kontrakt, który spełnia każda z nich: rama
 * nie niesie ani listy pól, ani żadnej decyzji o treści.
 *
 * Okno ma sześć sekcji: Konto, Uwierzytelnianie, Wygląd i język, Urządzenia,
 * Powiadomienia, Konta modeli. Jedna z nich (Konto) stoi jako miejsce nazwane:
 * mówi, czego w rdzeniu nie ma, i nie niesie ani jednego pola, przycisku czy
 * przełącznika bez pokrycia.
 *
 * Ujawnianie stopniowe obowiązuje wewnątrz sekcji tak samo jak w oknie: na
 * nastawę przypada jeden wiersz, na wierszu uchwyt z wartością bieżącą, a wybór
 * rozwija się dopiero pod kliknięciem (`wiersz-nastawy.ts`). Samo okno jest
 * przywoływane z listwy Centrum dowodzenia.
 *
 * Sekcji zbudowanych na `account.*` / `identity.*` tu nie ma: te same komendy
 * niesie już `budowa/client/src/modele/` (rejestr kont, katalog kategorii
 * tożsamości, dokumenty per oś, nakładka obowiązująca), a druga rama do tych
 * samych danych byłaby drugą prawdą. Pozycja „Konta modeli" nie niesie ani
 * jednego pola tamtych komend, tylko zdanie i drogę do tamtego okna
 * (`sekcja-konta-modeli.ts`).
 *
 * Okno Ustawień istnieje, bo pozycja `okno-ustawien` stoi w katalogu rdzenia
 * (`migracja_030_rejestr_okien_operacyjnych.sql`, kategoria `konfiguracja`),
 * a `auth.*` i motyw nie mają dokąd pójść — `modele/` niesie tożsamość modelu,
 * nie nastawy Operatora.
 *
 * Stan sześciu sekcji:
 *
 *  1. `konto` — miejsce nazwane. Bramka nie zna encji konta
 *     (`migracja_071_uwierzytelnienie_bramki.sql` zakłada wyłącznie metody
 *     i sesje), kontrakt nie ma rodziny profilu.
 *  2. `uwierzytelnianie` — cztery komendy `auth.*` po zalogowaniu, wykaz metod
 *     odświeżany na żywo zdarzeniem `auth.changed`.
 *  3. `wyglad` — motyw ma jeden ster i jedną prawdę z rdzeniem; język
 *     interfejsu nazwany jako brak, bez kontrolki.
 *  4. `urzadzenia` — wykaz urządzeń powiązanych z kontem (`device.list`)
 *     z unieważnieniem tokenu (`device.revoke`), odświeżany na żywo zdarzeniem
 *     `device.changed`.
 *  5. `powiadomienia` — macierz nastaw katalogu: przełącznik główny, siedem klas
 *     zdarzeń i kanały dostarczenia, wniesione migracją 377. Rodziny kontraktu
 *     ta sekcja nie potrzebuje — model danych (rozdz. 18.4) mówi wprost, że
 *     zakres klas i kanał dostarczenia SĄ ustawieniami konfiguracyjnymi, więc
 *     droga jest ta sama, co do każdej innej nastawy platformy. Samego
 *     doręczania nie ma jeszcze czym wykonać i sekcja mówi to wprost.
 *  6. `konta-modeli` — pozycja odsyłająca do `modele/`.
 *
 * Kolumnę nawigacji osadza `okno-ustawien.ts`, gdy rejestr niesie więcej niż
 * jedną pozycję.
 */

/** Kod sekcji — sześć pozycji projektu okna, w kolejności prezentacji. */
export type KodSekcjiUstawien =
  | 'konto'
  | 'uwierzytelnianie'
  | 'wyglad'
  | 'urzadzenia'
  | 'powiadomienia'
  | 'konta-modeli';

/**
 * Kontrakt wpięcia — kształt, który spełnia każdy plik sekcji.
 *
 * Sekcja dostaje kanał raz, przy budowie, i sama zarządza swoim stanem oraz
 * subskrypcjami. Rama okna nie zna pól, komend ani stanu wewnątrz sekcji —
 * wyłącznie ten kontrakt.
 */
export interface SekcjaUstawien {
  /** Element osadzany w ciele okna, gdy sekcja jest czynna. */
  element: HTMLElement;
  /**
   * Zleca odczyt / ponowny odczyt danych sekcji. Wołane przy pierwszym
   * pokazaniu sekcji i przy naciśnięciu „Odczytaj ponownie" w ramie.
   * Sekcja bez pokrycia w rdzeniu (Konto) tę metodę
   * niesie także — nanosi jawny stan „tego nie ma", nie milczy: milczenie
   * wyglądałoby jak odczyt, który nic nie znalazł.
   */
  odswiez(): void;
  /** Odłącza subskrypcje kanału. Wołane przy rozłączeniu całego okna. */
  rozlacz(): void;
}

/** Metadane sekcji potrzebne nawigacji: kolejność, nazwa, ikona. */
export interface OpisSekcjiUstawien {
  kod: KodSekcjiUstawien;
  nazwa: string;
  ikona: NazwaIkony;
  utworz(kanal: Kanal): SekcjaUstawien;
}

/**
 * Rejestr w kolejności prezentacji.
 *
 * Kolejność stoi tutaj, bo nawigacja okna nie ma skąd wziąć jej znikąd indziej:
 * sekcje są stałe, nie przychodzą katalogiem rdzenia (inaczej niż kategorie
 * okna konfiguracji, które rdzeń oddaje `settings.category.list`).
 *
 * Kontrakt `OpisSekcjiUstawien` jest jeden dla wszystkich, więc wytwórnia, która
 * kanału nie potrzebuje, po prostu go nie czyta. Drugi kształt wpisu „bez
 * kanału" rozdzieliłby rejestr na dwa rodzaje pozycji i zmusił ramę do
 * rozróżniania ich przy montażu.
 */
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
