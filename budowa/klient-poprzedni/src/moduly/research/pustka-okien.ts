import type { StanBadania } from './stan-badania';
import type { StanOknaBadania } from './stan-okna-badania';

/** Zdania stanu pustego siedmiu okien modułu Research, w formie inwentarza z tytułem i opisem, tłumaczą, czym okno jest i jak je zapełnić. */
export interface ZdaniePustki {
  tytul: string;
  opis: string;
}

// Przed pierwszym odczytem moduł nie pytał jeszcze rdzenia o nic — pustka nie znaczy, że rdzeń nic nie ma, tylko że modułu jeszcze nie zapytał.
const PRZED_ODCZYTEM: ZdaniePustki = {
  tytul: 'Okno jeszcze nie pytało rdzenia',
  opis:
    'Moduł czeka na kartę sesji z powłoki; po jej wskazaniu przeczyta okno badania i jego stan. ' +
    'Dopóki nie przeczytał, nie ma o czym twierdzić, że jest puste.',
};

/** Rdzeń odpowiedział, ale okna modułu Research nie dał w tej sesji — pustka bez miejsca na treść, nie pusty wykaz gotowego okna. */
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

/** Research Workspace trzyma zakres badania i jego etapy: jedno zdanie o tym, co i po co badamy, oraz wykaz kroków do wypełnienia. */
export const PUSTKA_ZAKRESU: ZdaniePustki = {
  tytul: 'Badanie bez zestawionego zakresu',
  opis:
    'Research Workspace trzyma zakres badania i jego etapy: jedno zdanie o tym, co i po co ' +
    'badamy, oraz wykaz kroków. Wpisz je w polach poniżej i naciśnij „Zapisz zakres badania" ' +
    '— pozostałe okna modułu biorą zakres stąd.',
};

/** Discovery Panel wyszukuje i odkrywa źródła materiału badania spośród wiedzy Operatora oraz treści zasobów repozytorium. */
export const PUSTKA_ODKRYWANIA: ZdaniePustki = {
  tytul: 'Panel jeszcze o nic nie pytał',
  opis:
    'Discovery Panel szuka materiału do badania. Dwa tryby mają dziś komendę kontraktu i idą ' +
    'nią naprawdę: „własne semantyczne" przeszukuje wiedzę Operatora po znaczeniu, ' +
    '„pełnotekstowe" — treść zasobów repozytorium. Tryb webowy i naukowy komendy nie mają ' +
    'i kończą się odmową rdzenia, wypisaną tu wprost. Wpisz zapytanie i naciśnij „Szukaj".',
};

/** Reading View czyta materiał wskazanego źródła i zamienia zaznaczony fragment w ustalenie z cytatem powiązanym z tym źródłem. */
export const PUSTKA_LEKTURY: ZdaniePustki = {
  tytul: 'Nie wskazano, co czytać',
  opis:
    'Reading View czyta materiał źródła i zamienia zaznaczony fragment w ustalenie z cytatem. ' +
    'Naciśnij „Czytaj" przy pozycji w Sources Manager. Treść dochodzi jedną drogą kontraktu — ' +
    'podglądem zasobu repozytorium — więc czytelne jest źródło wskazujące dokument ' +
    'repozytorium; przy pozostałych okno powie to wprost, zamiast pokazać pusty czytnik.',
};

/** Sources Manager jest katalogiem materiału badania: tytuł, typ, pochodzenie i ocena wiarygodności każdego zapisanego źródła. */
export const PUSTKA_ZRODEL: ZdaniePustki = {
  tytul: 'Katalog źródeł jeszcze pusty',
  opis:
    'Sources Manager jest katalogiem materiału badania: tytuł, typ, pochodzenie i ocena ' +
    'wiarygodności każdego źródła. Skataloguj pierwsze formularzem powyżej. Wykaz składa się ' +
    'ze źródeł potwierdzonych przez rdzeń w tej sesji: komenda odczytu wykazu jest już ' +
    'w kontrakcie, ale rdzeń nie ma dla niej uchwytu, więc materiał zebrany wcześniej nie ' +
    'ma jeszcze którędy dojść.',
};

/** Findings Panel zbiera ustalenia narastające w toku badania: jedno ustalenie na wpis, wiązane z zaznaczonymi źródłami materiału. */
export const PUSTKA_USTALEN: ZdaniePustki = {
  tytul: 'Bez zapisanego jeszcze ustalenia',
  opis:
    'Findings Panel zbiera to, co z materiału wynika: jedno ustalenie na wpis, wiązane ' +
    'z zaznaczonymi źródłami. Zapisz pierwsze formularzem powyżej. Wykaz narasta wyłącznie ' +
    'z zapisów tej sesji: komenda odczytu i zdarzenie zmiany ustalenia są już w kontrakcie, ' +
    'a rdzeń nie ma jeszcze uchwytu ani nie rozgłasza zdarzenia, więc ustalenie zapisane na ' +
    'innym urządzeniu konta tu nie dojdzie.',
};

/** Report Builder składa raport z ustaleń badania, zaznaczonych w Findings Panel, w kompozycję gotową do wydania dokumentowego. */
export const PUSTKA_RAPORTU: ZdaniePustki = {
  tytul: 'Raport jeszcze nie złożony',
  opis:
    'Report Builder składa raport z ustaleń. Redakcja pusta kieruje budowę na kanał modelu — ' +
    'rdzeń redaguje streszczenie zaznaczonych ustaleń; sekcja wypełniona ma pierwszeństwo ' +
    'i składa raport z niej samej. Zaznacz ustalenia w Findings Panel i otwórz kreator.',
};

/** Export Panel wydaje złożony raport badania w formacie dokumentowym, pod miejscem docelowym wybranym przez Operatora. */
export const PUSTKA_EKSPORTU: ZdaniePustki = {
  tytul: 'Nie ma jeszcze czego wydać',
  opis:
    'Export Panel wydaje złożony raport w formacie dokumentowym i pod wskazanym miejscem. ' +
    'Złóż raport w Report Builderze. Wybór formatu i miejsca docelowego zostaje czynny już ' +
    'teraz — to warunek merytoryczny, nie blokada.',
};

/**
 * Wybiera i pokazuje zdanie pustki właściwe dla stanu modułu: nie pytałem, pytałem bez okna,
 * albo miejsce jest, a treści nie ma.
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
