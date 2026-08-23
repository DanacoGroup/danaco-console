/**
 * Układ okien równoległych — interfejs katalogu.
 *
 * Od jednego do czterech okien komunikacji obok siebie (`LICZBA_MIN`,
 * `LICZBA_MAX`), przełącznik liczby, nagłówki ról i widoczna więź
 * koordynator–wykonawca.
 *
 * Montaż na scenie powłoki:
 *
 *   const uklad = zamontujUkladOkien(korzen, rdzen, opis);
 *   uklad.ustawLiczbe(2);
 *   uklad.nadajRole('okno-1', WindowRole.Coordinator);
 *   uklad.nadajRole('okno-2', WindowRole.Executor);
 *   uklad.pokazPrzekazanie('okno-1', 'okno-2');
 *
 * Sam układ powstaje przez `utworzUkladOkien()` i nie wymaga rdzenia —
 * montaż dokłada wyłącznie podpięcie pierwszego gniazda do łączności.
 */
export { utworzUkladOkien, type OpcjeUkladu, type UkladOkien } from './uklad-okien';

export { zamontujUkladOkien } from './montaz-ukladu';

export { figuraModulu, rolePrzyLiczbie, type FiguraModulu } from './figura-modulu';

export {
  ID_GNIAZD,
  LICZBA_MAX,
  LICZBA_MIN,
  czyIdGniazda,
  numerGniazda,
  type IdGniazda,
} from './identyfikatory';

export { czyRolaWPetli, stanZKolejki, type StanPary } from './stan-pary';

export type { GniazdoOkna } from './gniazdo-okna';

// Układ przestrzeni roboczej: kolumna rozmowy, stos paneli, widok pełnoekranowy.
// Panele należą do rozmowy, nie do ekranu — sterowanie stoi w nagłówku gniazda,
// a każde gniazdo ma własny zestaw.
export {
  SZEROKOSC_MIN_PANELU,
  SZEROKOSC_MIN_ROZMOWY,
  SZEROKOSC_PANELU_DOMYSLNA,
  minimumObszaru,
  type RodzajObszaru,
} from './rodzaje-obszaru';

export {
  kolumnyGniazda,
  ograniczSzerokoscPaneli,
  ustalPodzial,
  type PodzialGniazda,
} from './szerokosci-gniazda';

export { utworzPaneleGniazda, type OpcjePaneliGniazda, type PaneleGniazda } from './panele-gniazda';

export { utworzStanPaneli, type StanPaneli } from './stan-paneli';

// Łączność w oknie roboczym: stan łącza widać w samym gnieździe, nie tylko
// w pasku górnym — okno rozwinięte na pełny ekran przykrywa pasek.
export {
  OPIS_PONAWIANIA,
  ZDANIE_BRAKU_PRZEBIEGU,
  ZDANIE_KOLEJKI,
  czyWidocznaLacznosc,
  napisLacznosci,
  podpowiedzLacznosci,
  type OdczytLacznosci,
  type PortPonawiania,
  type StanLacznosciOkna,
} from './lacznosc-okna';

export {
  utworzPlakietkeLacznosci,
  type OpcjePlakietkiLacznosci,
  type PlakietkaLacznosci,
} from './plakietka-lacznosci';

export { portPonawiania, zwiazLacznoscUkladu } from './zrodlo-lacznosci';

// Wzorce pasa poziomego: przewijanie przepełnienia i przeciąganie kolejności.
// Nie wiedzą nic o treści pasa — wołający podaje pozycje i skutek przestawienia.
export {
  zwiazPrzepelnieniePasa,
  type OpcjePrzepelnienia,
  type StronaPrzepelnienia,
  type WiezPrzepelnienia,
} from './przepelnienie-pasa';

export {
  miejsceUpuszczenia,
  przestawWTablicy,
  zwiazPrzeciaganieKolejnosci,
  type OpcjePrzeciagania,
  type WiezPrzeciagania,
} from './przeciaganie-kolejnosci';

export {
  ZDANIE_KOLEJNOSCI_MIEJSCOWEJ,
  ZDANIE_KOLEJNOSCI_RDZENIA,
  ZDANIE_KOLEJNOSCI_ULOTNEJ,
  utworzKolejnoscMiejscowa,
  type KolejnoscMiejscowa,
  type OpcjeKolejnosci,
} from './kolejnosc-miejscowa';
