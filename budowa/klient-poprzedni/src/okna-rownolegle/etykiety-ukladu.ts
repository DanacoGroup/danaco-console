import { nazwaRoli } from '../okno-komunikacji/etykiety-okna';
import type { WindowRole } from '../../../shared/contract';
import type { PostacRozmowy } from '../okno-komunikacji/profil-modulu';
import { LICZBA_MAX, numerGniazda, type IdGniazda } from './identyfikatory';
import type { StanPary } from './stan-pary';

/**
 * Napisy układu okien równoległych, wszystkie w jednym miejscu.
 */

/**
 * Etykieta przełącznika liczby okien widoczna w pasku sterowania układem okien równoległych sceny. Dotyczy sceny okien równoległych.
 */
export const ETYKIETA_PRZELACZNIKA = 'Okna komunikacji';

/**
 * Objaśnienie pod przełącznikiem liczby okien: co naprawdę zmienia się w scenie po zmianie tej liczby.
 */
export const OPIS_PRZELACZNIKA =
  'Każde okno ma własny model, katalogi robocze, tryb uprawnień i rolę.';

/**
 * Treść dymka pomocy przy przełączniku liczby okien, wyjaśniająca skutek zmiany tej liczby okien. Dotyczy sceny okien równoległych.
 */
export const OBJASNIENIE_PRZELACZNIKA =
  'Wprowadzenie okna na scenę zamawia dla niego okno w rdzeniu komendą '
  + 'window.create. Ustawienia okna są jego własne — zmiana w jednym oknie '
  + 'nie przenosi się na pozostałe.';

/**
 * Nagłówek pasa relacji, który pokazuje powiązanie koordynator–wykonawca między dwoma oknami sceny. Dotyczy sceny okien równoległych.
 */
export const ETYKIETA_RELACJI = 'Relacje';

/**
 * Treść dymka pomocy przy pasie relacji, wyjaśniająca, co pas relacji pokazuje operatorowi sceny. Dotyczy sceny okien równoległych.
 */
export const OBJASNIENIE_RELACJI =
  'Pas rysuje jedną parę koordynator–wykonawca: kto komu zleca i w jakim '
  + 'położeniu jest pętla. Położenie przychodzi ze stanu kolejki rdzenia '
  + '(zdarzenie queue.changed), a rola okna ze zdarzenia window.changed.';

/**
 * Tytuł stanu pustego pasa relacji, pokazywany, gdy para koordynator–wykonawca jeszcze nie istnieje na scenie.
 */
export const PUSTKA_RELACJI_TYTUL = 'Brak pary koordynator–wykonawca';

/**
 * Opis stanu pustego pasa relacji: co operator ma zrobić, żeby para koordynator–wykonawca powstała na scenie.
 */
export const PUSTKA_RELACJI_OPIS =
  'Pętla biegnie między dwoma oknami — jednym w roli koordynatora i jednym '
  + 'w roli wykonawcy. Ustaw co najmniej dwa okna przełącznikiem powyżej.';

/**
 * Kontrakt nie ma komendy przekazania zlecenia koordynator-wykonawca; scena pokazuje wyłącznie jego wygląd.
 */
export const PRZEKAZANIE_TYTUL = 'Przekazanie zlecenia — podgląd układu';

/**
 * Treść uwagi o przekazaniu zlecenia; wymienia powód braku obsługi, nie samą niemożność wykonania jej.
 */
export const PRZEKAZANIE_OPIS =
  'Nic nie zostało wysłane do rdzenia: kontrakt nie ma komendy przekazania '
  + 'zlecenia koordynator→wykonawca. Wpisy w obu oknach i ruch żetonu pokazują '
  + 'wyłącznie układ tej chwili.';

/**
 * Tytuł gniazda widoczny w jego własnym nagłówku na scenie okien równoległych tej budowy. Dotyczy sceny okien równoległych.
 */
export function tytulGniazda(id: IdGniazda): string {
  return `Okno ${numerGniazda(id)}`;
}

/**
 * Nazwa modułu przy tytule gniazda jest wartością bieżącą, nie napisem rodzajowym tego gniazda sceny. Dotyczy sceny okien równoległych.
 */
export function nazwaModuluGniazda(nazwaModulu: string): string {
  return `— ${nazwaModulu}`;
}

/**
 * Nazwa gniazda wraz z jego rolą, używana na obu końcach pasa relacji koordynator–wykonawca sceny. Dotyczy sceny okien równoległych.
 */
export function opisGniazdaZRola(id: IdGniazda, rola: WindowRole): string {
  return `${tytulGniazda(id)} · ${nazwaRoli(rola)}`;
}

/**
 * Nazwa stanu pary koordynator–wykonawca widoczna operatorowi na scenie okien równoległych aplikacji. Dotyczy sceny okien równoległych.
 */
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

/**
 * Kierunek widoczny w nagłówku koordynatora: komu ten koordynator zleca pracę na scenie okien. Dotyczy sceny okien równoległych.
 */
export function kierunekKoordynatora(wykonawca: IdGniazda): string {
  return `Zleca ${celownik(wykonawca)}`;
}

/**
 * Kierunek widoczny w nagłówku wykonawcy: od kogo ten wykonawca bierze zlecenia na scenie okien. Dotyczy sceny okien równoległych.
 */
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

/**
 * Wpis historii w oknie wykonawcy w chwili przyjęcia zlecenia od koordynatora sceny okien. Dotyczy sceny okien równoległych.
 */
export function komunikatPrzyjecia(koordynator: IdGniazda, tresc: string): string {
  return `Podgląd układu — zlecenie od ${dopelniacz(koordynator)}: ${tresc}. `
    + 'Nic nie przyjęto z rdzenia: kontrakt nie ma komendy przekazania.';
}

/**
 * Treść zlecenia przyjmowana, gdy samo wywołanie zlecenia nie podaje treści tego zlecenia wprost. Dotyczy sceny okien równoległych.
 */
export const TRESC_ZLECENIA_DOMYSLNA = 'treść zastępcza podglądu, nie zlecenie rdzenia';

/**
 * Celownik nazwy gniazda, na przykład „Oknu 2”, używany w zdaniach kierowanych wprost do gniazda. Dotyczy sceny okien równoległych.
 */
function celownik(id: IdGniazda): string {
  return `Oknu ${numerGniazda(id)}`;
}

/**
 * Dopełniacz nazwy gniazda, na przykład „Okna 2”, używany w zdaniach opisujących to gniazdo sceny. Dotyczy sceny okien równoległych.
 */
function dopelniacz(id: IdGniazda): string {
  return `Okna ${numerGniazda(id)}`;
}

/**
 * Gniazdo wraz z rolą, jaką to gniazdo bierze przy budowie figury rozmowy modułu na scenie. Dotyczy sceny okien równoległych.
 */
export interface GniazdoZRola {
  id: IdGniazda;
  rola: WindowRole;
}

/**
 * Tyle informacji o figurze modułu, ile potrzeba do zbudowania pełnego zdania dla operatora sceny. Dotyczy sceny okien równoległych.
 */
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
 * Zdanie o figurze rozmowy modułu mówi, ile okien i w jakich rolach ta figura modułu obejmuje. Dotyczy sceny okien równoległych.
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

/**
 * Zdanie o pozycjach przełącznika ponad figurą modułu — wyjaśnia, czemu okna się nie dokładają same. Dotyczy sceny okien równoległych.
 */
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

/**
 * Liczba okien w mianowniku, na przykład „dwa okna”; poza zakresem sceny zostaje sama cyfra liczby. Dotyczy sceny okien równoległych.
 */
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

/**
 * Liczba okien w dopełniaczu, na przykład „czterech okien”, używana w zdaniach o liczbie okien sceny. Dotyczy sceny okien równoległych.
 */
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
