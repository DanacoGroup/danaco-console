import type { StanBadania } from './stan-badania';
import type { StanOknaBadania } from './stan-okna-badania';

/**
 * Zdania stanu pustego siedmiu okien modułu Research — wyjęte z plików czynności.
 *
 * Stan pusty tłumaczy, czym okno jest i jak je zapełnić, a nie melduje brak
 * pozycji. Takie zdania są tekstem produktu, nie logiką widoku: mają dać się
 * przeczytać obok siebie i poprawić w jednym miejscu, gdy kontrakt dostanie
 * komendę, której dziś nie ma. Ten sam wzorzec niesie
 * `moduly/library/etykiety-biblioteki.ts`.
 *
 * Brak okna badania to co innego niż pusty wykaz, więc rozróżnienie stoi w
 * jednym miejscu, nie w siedmiu oknach. Komendy zapisu źródła, ustalenia
 * i budowy raportu mają `windowId` w polach obowiązkowych, więc bez okna
 * zaproszenie „skataloguj pierwsze źródło" prowadzi Operatora wprost pod odmowę.
 *
 * Żadne z tych zdań nie odbiera klikalności ani jednej kontrolce — nazywa
 * warunek merytoryczny, po czym formularze i przyciski zostają czynne,
 * a odmowa rdzenia wraca własnymi słowami rdzenia.
 */

/** Stan pusty w formie inwentarza: tytuł i zdanie wyjaśniające pod nim. */
export interface ZdaniePustki {
  tytul: string;
  opis: string;
}

/**
 * Przed pierwszym odczytem — moduł nie pytał jeszcze rdzenia o nic.
 *
 * „Jeszcze nie pytałem" nie jest tym samym, co „rdzeń nic nie ma": twierdzenie
 * o zawartości rdzenia postawione bez jego odpowiedzi jest zgadywaniem.
 */
const PRZED_ODCZYTEM: ZdaniePustki = {
  tytul: 'Okno jeszcze nie pytało rdzenia',
  opis:
    'Moduł czeka na kartę sesji z powłoki; po jej wskazaniu przeczyta okno badania i jego stan. ' +
    'Dopóki nie przeczytał, nie ma o czym twierdzić, że jest puste.',
};

/** Rdzeń odpowiedział, ale okna modułowi nie dał — pustka bez miejsca na treść. */
function bezOknaBadania(powod: string): ZdaniePustki {
  return {
    tytul: 'Badanie bez okna w rdzeniu',
    opis:
      (powod === '' ? 'Rdzeń nie wskazał okna modułu Research w tej sesji. ' : `${powod} `) +
      'Komendy obszaru research wymagają pola windowId, więc nie jest to pusty wykaz, tylko ' +
      'brak miejsca, w którym wykaz miałby powstać. Kontrolki zostają czynne — zapis wróci ' +
      'odmową rdzenia, a nie ciszą.',
  };
}

/** Research Workspace — zakres badania i wejście do pozostałych okien. */
export const PUSTKA_ZAKRESU: ZdaniePustki = {
  tytul: 'Badanie bez zestawionego zakresu',
  opis:
    'Research Workspace trzyma zakres badania i jego etapy: jedno zdanie o tym, co i po co ' +
    'badamy, oraz wykaz kroków. Wpisz je w polach poniżej i naciśnij „Zapisz zakres badania" ' +
    '— pozostałe okna modułu biorą zakres stąd.',
};

/** Discovery Panel — wyszukiwanie i odkrywanie źródeł. */
export const PUSTKA_ODKRYWANIA: ZdaniePustki = {
  tytul: 'Panel jeszcze o nic nie pytał',
  opis:
    'Discovery Panel szuka materiału do badania. Dwa tryby mają dziś komendę kontraktu i idą ' +
    'nią naprawdę: „własne semantyczne" przeszukuje wiedzę Operatora po znaczeniu, ' +
    '„pełnotekstowe" — treść zasobów repozytorium. Tryb webowy i naukowy komendy nie mają ' +
    'i kończą się odmową rdzenia, wypisaną tu wprost. Wpisz zapytanie i naciśnij „Szukaj".',
};

/** Reading View — lektura materiału i wypisy z niego. */
export const PUSTKA_LEKTURY: ZdaniePustki = {
  tytul: 'Nie wskazano, co czytać',
  opis:
    'Reading View czyta materiał źródła i zamienia zaznaczony fragment w ustalenie z cytatem. ' +
    'Naciśnij „Czytaj" przy pozycji w Sources Manager. Treść dochodzi jedną drogą kontraktu — ' +
    'podglądem zasobu repozytorium — więc czytelne jest źródło wskazujące dokument ' +
    'repozytorium; przy pozostałych okno powie to wprost, zamiast pokazać pusty czytnik.',
};

/** Sources Manager — katalog materiału badania. */
export const PUSTKA_ZRODEL: ZdaniePustki = {
  tytul: 'Katalog źródeł jeszcze pusty',
  opis:
    'Sources Manager jest katalogiem materiału badania: tytuł, typ, pochodzenie i ocena ' +
    'wiarygodności każdego źródła. Skataloguj pierwsze formularzem powyżej. Wykaz składa się ' +
    'ze źródeł potwierdzonych przez rdzeń w tej sesji: komenda odczytu wykazu jest już ' +
    'w kontrakcie, ale rdzeń nie ma dla niej uchwytu, więc materiał zebrany wcześniej nie ' +
    'ma jeszcze którędy dojść.',
};

/** Findings Panel — ustalenia narastające w toku badania. */
export const PUSTKA_USTALEN: ZdaniePustki = {
  tytul: 'Bez zapisanego jeszcze ustalenia',
  opis:
    'Findings Panel zbiera to, co z materiału wynika: jedno ustalenie na wpis, wiązane ' +
    'z zaznaczonymi źródłami. Zapisz pierwsze formularzem powyżej. Wykaz narasta wyłącznie ' +
    'z zapisów tej sesji: komenda odczytu i zdarzenie zmiany ustalenia są już w kontrakcie, ' +
    'a rdzeń nie ma jeszcze uchwytu ani nie rozgłasza zdarzenia, więc ustalenie zapisane na ' +
    'innym urządzeniu konta tu nie dojdzie.',
};

/** Report Builder — kompozycja raportu z ustaleń. */
export const PUSTKA_RAPORTU: ZdaniePustki = {
  tytul: 'Raport jeszcze nie złożony',
  opis:
    'Report Builder składa raport z ustaleń. Redakcja pusta kieruje budowę na kanał modelu — ' +
    'rdzeń redaguje streszczenie zaznaczonych ustaleń; sekcja wypełniona ma pierwszeństwo ' +
    'i składa raport z niej samej. Zaznacz ustalenia w Findings Panel i otwórz kreator.',
};

/** Export Panel — wydanie raportu w formacie dokumentowym. */
export const PUSTKA_EKSPORTU: ZdaniePustki = {
  tytul: 'Nie ma jeszcze czego wydać',
  opis:
    'Export Panel wydaje złożony raport w formacie dokumentowym i pod wskazanym miejscem. ' +
    'Złóż raport w Report Builderze. Wybór formatu i miejsca docelowego zostaje czynny już ' +
    'teraz — to warunek merytoryczny, nie blokada.',
};

/**
 * Wybiera i pokazuje zdanie pustki właściwe dla stanu, w jakim moduł naprawdę
 * stoi. Trzy przypadki, bo trzy różne rzeczy: nie pytałem · pytałem i nie mam
 * gdzie · pytałem, miejsce mam, treści nie ma.
 *
 * Fazy `odczyt` i `blad` tu nie dochodzą — rozstrzygają je czynności okien
 * wcześniej, wskaźnikiem odczytu i komunikatem odmowy.
 */
export function pokazPustke(
  okno: StanOknaBadania,
  stan: StanBadania,
  wlasne: ZdaniePustki,
): void {
  const zdanie =
    stan.faza() === 'spoczynek'
      ? PRZED_ODCZYTEM
      : stan.idOkna() === ''
        ? bezOknaBadania(stan.powod())
        : wlasne;
  okno.puste(zdanie.tytul, zdanie.opis);
}
