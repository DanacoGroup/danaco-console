import { nazwaRoli } from '../okno-komunikacji/etykiety-okna';
import type { WindowRole } from '../../../shared/contract';
import type { PostacRozmowy } from '../okno-komunikacji/profil-modulu';
import { LICZBA_MAX, numerGniazda, type IdGniazda } from './identyfikatory';
import type { StanPary } from './stan-pary';

/**
 * Napisy układu okien równoległych.
 *
 * Jedna odpowiedzialność: wszystkie ciągi widoczne dla operatora w jednym
 * miejscu. Nazwy ról nie powstają tutaj po raz drugi — pochodzą z etykiet
 * okna komunikacji, żeby ta sama rola nie nazywała się w dwóch miejscach
 * inaczej.
 */

/** Etykieta przełącznika liczby okien. */
export const ETYKIETA_PRZELACZNIKA = 'Okna komunikacji';

/** Objaśnienie pod przełącznikiem: co naprawdę zmienia liczba okien. */
export const OPIS_PRZELACZNIKA =
  'Każde okno ma własny model, katalogi robocze, tryb uprawnień i rolę.';

/** Treść dymka [?] przy przełączniku liczby okien. */
export const OBJASNIENIE_PRZELACZNIKA =
  'Wprowadzenie okna na scenę zamawia dla niego okno w rdzeniu komendą '
  + 'window.create. Ustawienia okna są jego własne — zmiana w jednym oknie '
  + 'nie przenosi się na pozostałe.';

/** Nagłówek pasa relacji — pas pokazuje powiązanie między oknami. */
export const ETYKIETA_RELACJI = 'Relacje';

/** Treść dymka [?] przy pasie relacji. */
export const OBJASNIENIE_RELACJI =
  'Pas rysuje jedną parę koordynator–wykonawca: kto komu zleca i w jakim '
  + 'położeniu jest pętla. Położenie przychodzi ze stanu kolejki rdzenia '
  + '(zdarzenie queue.changed), a rola okna ze zdarzenia window.changed.';

/** Tytuł stanu pustego pasa relacji — pary nie ma. */
export const PUSTKA_RELACJI_TYTUL = 'Brak pary koordynator–wykonawca';

/** Opis stanu pustego pasa relacji: co zrobić, żeby para powstała. */
export const PUSTKA_RELACJI_OPIS =
  'Pętla biegnie między dwoma oknami — jednym w roli koordynatora i jednym '
  + 'w roli wykonawcy. Ustaw co najmniej dwa okna przełącznikiem powyżej.';

/**
 * Uwaga przy przekazaniu zlecenia.
 *
 * Kontrakt (`shared/contract.json`, klucz `komendy`) nie ma komendy przekazania
 * zlecenia koordynator→wykonawca; scena umie pokazać wyłącznie, jak takie
 * przekazanie wygląda. Komunikat mówi to wprost, żeby ruch na pasie nie uchodził
 * za wykonaną pracę.
 */
export const PRZEKAZANIE_TYTUL = 'Przekazanie zlecenia — podgląd układu';

/** Treść uwagi o przekazaniu; wymienia powód, nie samą niemożność. */
export const PRZEKAZANIE_OPIS =
  'Nic nie zostało wysłane do rdzenia: kontrakt nie ma komendy przekazania '
  + 'zlecenia koordynator→wykonawca. Wpisy w obu oknach i ruch żetonu pokazują '
  + 'wyłącznie układ tej chwili.';

/** Tytuł gniazda widoczny w jego nagłówku. */
export function tytulGniazda(id: IdGniazda): string {
  return `Okno ${numerGniazda(id)}`;
}

/**
 * Nazwa modułu przy tytule gniazda — wartość bieżąca, nie napis rodzajowy.
 *
 * Kreska średnia oddziela dwie odpowiedzi na dwa różne pytania: „które to
 * miejsce na scenie" (numer) i „co w nim pracuje" (moduł). Bez niej czytałoby
 * się to jak jedna, dłuższa nazwa okna.
 */
export function nazwaModuluGniazda(nazwaModulu: string): string {
  return `— ${nazwaModulu}`;
}

/** Nazwa gniazda wraz z rolą — używana na końcach pasa relacji. */
export function opisGniazdaZRola(id: IdGniazda, rola: WindowRole): string {
  return `${tytulGniazda(id)} · ${nazwaRoli(rola)}`;
}

/** Nazwa stanu pary koordynator–wykonawca. */
export function nazwaStanuPary(stan: StanPary): string {
  switch (stan) {
    case 'brak-pary':
      return 'Brak pary koordynator–wykonawca';
    case 'gotowa':
      return 'Para gotowa';
    case 'wykonawca-pracuje':
      return 'Wykonawca pracuje';
    case 'koordynator-wybudzony':
      return 'Koordynator wybudzony';
    case 'kolejka-wstrzymana':
      return 'Kolejka wstrzymana';
    case 'przekazanie':
      return 'Przekazanie zlecenia';
  }
}

/** Kierunek widoczny w nagłówku koordynatora: komu zleca. */
export function kierunekKoordynatora(wykonawca: IdGniazda): string {
  return `Zleca ${celownik(wykonawca)}`;
}

/** Kierunek widoczny w nagłówku wykonawcy: od kogo bierze zlecenia. */
export function kierunekWykonawcy(koordynator: IdGniazda): string {
  return `Zlecenia z ${dopelniacz(koordynator)}`;
}

/**
 * Wpis historii w oknie koordynatora w chwili przekazania.
 *
 * Wpis nie melduje wykonanej czynności, bo czynności nie było — nazywa podgląd
 * podglądem.
 */
export function komunikatWyslania(wykonawca: IdGniazda, tresc: string): string {
  return `Podgląd układu — zlecenie do ${dopelniacz(wykonawca)}: ${tresc}. `
    + 'Nic nie wysłano do rdzenia: kontrakt nie ma komendy przekazania.';
}

/** Wpis historii w oknie wykonawcy w chwili przyjęcia zlecenia. */
export function komunikatPrzyjecia(koordynator: IdGniazda, tresc: string): string {
  return `Podgląd układu — zlecenie od ${dopelniacz(koordynator)}: ${tresc}. `
    + 'Nic nie przyjęto z rdzenia: kontrakt nie ma komendy przekazania.';
}

/** Treść zlecenia przyjmowana, gdy wywołanie własnej nie podaje. */
export const TRESC_ZLECENIA_DOMYSLNA = 'treść zastępcza podglądu, nie zlecenie rdzenia';

/** „Oknu 2" — celownik nazwy gniazda. */
function celownik(id: IdGniazda): string {
  return `Oknu ${numerGniazda(id)}`;
}

/** „Okna 2" — dopełniacz nazwy gniazda. */
function dopelniacz(id: IdGniazda): string {
  return `Okna ${numerGniazda(id)}`;
}

/** Gniazdo wraz z rolą, jaką bierze przy figurze modułu. */
export interface GniazdoZRola {
  id: IdGniazda;
  rola: WindowRole;
}

/** Tyle o figurze modułu, ile potrzeba do zbudowania zdania dla Operatora. */
export interface FiguraDoNapisu {
  /** Nazwa modułu widoczna dla Operatora. */
  nazwa: string;
  /** Postać rozmowy modułu z profilu (`okno-komunikacji/profil-modulu.ts`). */
  postac: PostacRozmowy;
  /** Ile okien rozmowy moduł prowadzi; zero znaczy „rozmowy nie prowadzi". */
  liczbaOkien: number;
  /** Role gniazd przy tej liczbie okien — z `role-domyslne.ts`. */
  role: readonly GniazdoZRola[];
  /** Czy moduł ma ustaloną parę koordynator–wykonawca. */
  paraKoordynatorWykonawca: boolean;
}

/**
 * Zdanie sceny, dopóki rdzeń nie nazwał modułu okna.
 *
 * Pełny sufit widoczny wtedy na przełączniku jest stanem zastanym platformy,
 * a nie liczbą podaną o module.
 */
export const ZDANIE_MODUL_NIEUSTALONY =
  'Moduł okna nie został jeszcze potwierdzony przez rdzeń. Scena pracuje na '
  + 'stanie zastanym, a nie na rozstrzygnięciu o module.';

/**
 * Zdanie o figurze rozmowy modułu — ile okien i w jakich rolach.
 *
 * Moduł bez rozmowy (`brak`) i moduł mówiący dymkiem dostają zdanie wprost, bo
 * pusty przełącznik czytałby się jak awaria sceny.
 *
 * Przy `LICZBA_MAX` zdanie mówi o suficie platformy, nie o liczbie okien
 * modułu: dla części modułów czwórka jest stanem zastanym
 * (`GRANICA_NIEPODANA`), a nie liczbą podaną przez rdzeń.
 */
export function zdanieFiguryModulu(figura: FiguraDoNapisu): string {
  const modul = `Moduł ${figura.nazwa}`;

  if (figura.postac === 'brak') {
    return `${modul} nie prowadzi rozmowy — to rozstrzygnięcie o module, `
      + 'nie brak odpowiedzi rdzenia. Scena trzyma jedno okno, bo nigdy nie jest pusta.';
  }

  if (figura.postac === 'dymek-glosowy') {
    return `${modul} prowadzi rozmowę pływającym dymkiem głosowym, `
      + 'a nie oknem na scenie.';
  }

  if (figura.liczbaOkien >= LICZBA_MAX) {
    return `${modul} korzysta z pełnego sufitu platformy — do ${okienDopelniacz(LICZBA_MAX)} rozmowy.`;
  }

  if (figura.liczbaOkien <= 1) return `${modul} prowadzi jedno okno rozmowy.`;

  const obsada = figura.role.map((g) => opisGniazdaZRola(g.id, g.rola)).join(', ');
  return figura.paraKoordynatorWykonawca
    ? `${modul} prowadzi dokładnie ${okienMianownik(figura.liczbaOkien)} rozmowy: ${obsada}.`
    : `${modul} prowadzi ${okienMianownik(figura.liczbaOkien)} rozmowy. `
      + `Role są domyślne dla tej liczby okien (${obsada}) i Operator może je przestawić.`;
}

/** Zdanie o pozycjach przełącznika ponad figurą modułu — czemu nie dokładają okna. */
export function zdanieNadwyzkiPozycji(liczbaOkien: number): string {
  return `Pozycja wyższa niż ${liczbaOkien} nie dokłada tu okna.`;
}

/**
 * Zdanie o scenie trzymającej więcej okien, niż moduł prowadzi.
 *
 * Układ sam ich nie zdejmuje: zdjęcie okna ze sceny zamyka je również w rdzeniu
 * (`scena-sesji.ts` → `zamknijOstatnie`), więc należy do Operatora, a nie do
 * przestawienia modułu.
 */
export function zdanieNadmiaruSceny(naScenie: number, liczbaOkien: number): string {
  return `Scena trzyma teraz ${okienMianownik(naScenie)} — więcej, niż moduł `
    + `prowadzi (${okienMianownik(liczbaOkien)}). Układ ich nie zdejmuje: zamknięcie `
    + 'okna zamyka je również w rdzeniu i należy do Operatora.';
}

/** „dwa okna" — liczba okien w mianowniku; poza zakresem sceny zostaje cyfra. */
function okienMianownik(liczba: number): string {
  switch (liczba) {
    case 1:
      return 'jedno okno';
    case 2:
      return 'dwa okna';
    case 3:
      return 'trzy okna';
    case 4:
      return 'cztery okna';
    default:
      return `${liczba} okien`;
  }
}

/** „czterech okien" — liczba okien w dopełniaczu. */
function okienDopelniacz(liczba: number): string {
  switch (liczba) {
    case 1:
      return 'jednego okna';
    case 2:
      return 'dwóch okien';
    case 3:
      return 'trzech okien';
    case 4:
      return 'czterech okien';
    default:
      return `${liczba} okien`;
  }
}
