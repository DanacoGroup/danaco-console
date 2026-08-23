/**
 * Profil modułu — to, czym okno komunikacji różni się w każdym z 15 modułów.
 *
 * Chat Window jest jedynym oknem występującym we wszystkich modułach i przy
 * każdej zmianie modułu przestawia wygląd, możliwości, dostępne narzędzia
 * i kontekst. Profil niesie to, co się przestawia; historia wątku do profilu
 * nie należy, bo przy zmianie modułu zostaje nietknięta.
 *
 * Narzędzie paska promptu nie jest osobną komendą kontraktu — jest gotowym
 * poleceniem wstawianym do pola wypowiedzi, które model wykonuje wywołaniem
 * narzędzia platformy. Dzięki temu pasek nie zawiera ani jednego
 * przycisku bez działania: każda pozycja czymś kończy — treścią
 * w polu wypowiedzi, którą Operator zatwierdza albo poprawia.
 */

import { LICZBA_MAX } from '../okna-rownolegle/identyfikatory';

/** Pojedyncza pozycja paska narzędzi promptu. */
export interface NarzedziePromptu {
  /** Kod pozycji; przedrostkiem jest kod modułu albo `wspolne`. */
  kod: string;
  /** Etykieta przycisku. */
  nazwa: string;
  /** Po co pozycja jest — podpowiedź i opis dla czytnika ekranu. */
  opis: string;
  /** Treść wstawiana do pola wypowiedzi Operatora. */
  polecenie: string;
}

/**
 * Postać, w jakiej moduł prowadzi rozmowę.
 *
 * - `okno` — zwykłe okno czatu; większość modułów.
 * - `dymek-glosowy` — pływający awatar z oknem dymkowym, głos przed tekstem;
 *   dziś wyłącznie Assistant („Okno komunikacji: NIE. Zamiast niego: pływający
 *   awatar, okno dymkowe, komunikacja głosowa, komunikacja tekstowa").
 * - `brak` — moduł nie prowadzi rozmowy w ogóle; dziś wyłącznie Library
 *   („To nie jest środowisko pracy z AI, tylko menedżer zasobów").
 *
 * `brak` nie znaczy „pusto". Okno, które trafi na moduł bez rozmowy, ma
 * powiedzieć to Operatorowi wprost i pokazać, czym ten moduł jest zamiast
 * czatu — pusta lista czyta się jak „nic tu nie ma" i jest błędem.
 */
export type PostacRozmowy = 'okno' | 'dymek-glosowy' | 'brak';

/**
 * Granica okien równoległych jednego modułu.
 *
 * Zawęża sufit platformy `okna-rownolegle/identyfikatory.ts` → `LICZBA_MAX`,
 * nigdy go nie rozszerza. Sufit zostaje twardy i obowiązuje niezależnie od
 * profilu; granica modułu mówi tylko, ile z niego moduł wykorzystuje.
 */
export type GranicaOkien = 1 | 2 | 4;

/** Granice dopuszczalne, rosnąco — jedyne wartości `GranicaOkien`. */
export const GRANICE: readonly GranicaOkien[] = [1, 2, 4];

/**
 * Granica modułu, dla którego liczby okien nie ustalono.
 *
 * Ustalona jest dla: Apps (4), MultitaskingAI (4), Developer (do 4),
 * Automations (2), Translate (2), Agents (1), Assistant (1 dymek). Browser,
 * Research, Roundtable, Workspace, Design i Studio dostają pełny sufit
 * platformy. Ta stała istnieje po to, żeby domysł nie udawał decyzji: gdzie
 * stoi `GRANICA_NIEPODANA`, tam nikt liczby nie ustalił.
 */
export const GRANICA_NIEPODANA: GranicaOkien = 4;

/** Właściwości rozmowy w module — bez pasków i narzędzi. */
export interface PostacModulu {
  /** Czy i w jakiej postaci moduł prowadzi rozmowę. */
  postacRozmowy: PostacRozmowy;
  /** Ile okien równoległych moduł wykorzystuje; nigdy ponad `LICZBA_MAX`. */
  granicaOkien: GranicaOkien;
  /**
   * Czy rozmowa przeżywa zamknięcie okna.
   *
   * Agents ma tu `false`: czat jest roboczy i testowy, ginie przy zamknięciu
   * okna i przy zmianie testowanego agenta („To nie jest pełnoprawna sesja").
   */
  pamiecSesyjna: boolean;
  /**
   * Okna towarzyszące, bez których moduł nie jest sobą — podzbiór
   * `oknaKontekstu`. Reszta `oknaKontekstu` to okna możliwe.
   */
  oknaObowiazkowe?: readonly string[];
}

/**
 * Postać stanu zastanego: zwykłe okno, pełny sufit platformy, pamięć jest.
 *
 * Dostaje ją profil, który postaci nie podał, oraz moduł spoza rejestru.
 */
export const POSTAC_DOMYSLNA: PostacModulu = {
  postacRozmowy: 'okno',
  granicaOkien: GRANICA_NIEPODANA,
  pamiecSesyjna: true,
  oknaObowiazkowe: [],
};

/** Zestaw właściwości okna komunikacji w jednym module. */
export interface ProfilModulu {
  /** Kod modułu; odpowiednik kolumny `modul.kod` i pola `moduleId` kontraktu. */
  kod: string;
  /** Nazwa modułu widoczna we wskaźniku okna. */
  nazwa: string;
  /** Przeznaczenie modułu — pierwszy wiersz panelu kontekstu. */
  przeznaczenie: string;
  /** Pasek narzędzi promptu tego modułu wraz z arsenałem wspólnym. */
  narzedzia: readonly NarzedziePromptu[];
  /** Okna operacyjne modułu — źródła kontekstu okna komunikacji. */
  oknaKontekstu: readonly string[];
  /** Postać rozmowy modułu; `brak` znaczy „ten moduł rozmowy nie prowadzi". */
  postacRozmowy: PostacRozmowy;
  /** Granica okien równoległych modułu, już zmieszczona w suficie platformy. */
  granicaOkien: GranicaOkien;
  /** Czy rozmowa modułu przeżywa zamknięcie okna. */
  pamiecSesyjna: boolean;
  /** Okna towarzyszące obowiązkowe — podzbiór `oknaKontekstu`. */
  oknaObowiazkowe: readonly string[];
}

/**
 * Największa granica mieszcząca się jednocześnie w żądaniu modułu i w suficie
 * platformy. Sufit zawęża; podniesienia granicy ponad `LICZBA_MAX` nie ma jak
 * wyrazić, bo wynik zawsze pochodzi z `GRANICE`.
 */
export function granicaWSuficie(zadana: GranicaOkien): GranicaOkien {
  const mieszczace = GRANICE.filter((g) => g <= zadana && g <= LICZBA_MAX);
  return mieszczace.length > 0 ? mieszczace[mieszczace.length - 1] : 1;
}

/**
 * Ile okien rozmowy Operator może otworzyć w tym module.
 *
 * Zero dla modułu bez rozmowy — i to jest odpowiedź, nie brak odpowiedzi.
 * Kto pyta o liczbę okien, pyta tą funkcją, a nie polem `granicaOkien`:
 * pole niesie granicę, funkcja niesie prawo do otwarcia.
 */
export function liczbaOkienRozmowy(p: ProfilModulu): number {
  return p.postacRozmowy === 'brak' ? 0 : p.granicaOkien;
}

/** Pozycja paska narzędzi promptu. */
export function narzedzie(
  kod: string,
  nazwa: string,
  opis: string,
  polecenie: string,
): NarzedziePromptu {
  return { kod, nazwa, opis, polecenie };
}

/**
 * Arsenał wspólny Chat Window — pozycje obecne w każdym module.
 *
 * Trzy pozycje, trzy komendy kontraktu, które wywoła model: `action.list`,
 * `window.state.get`, `context.transfer`. Moduł nierozpoznany dostaje sam ten
 * arsenał, więc pasek nigdy nie jest pusty.
 */
export const NARZEDZIA_WSPOLNE: readonly NarzedziePromptu[] = [
  narzedzie(
    'wspolne.katalog-akcji',
    'Akcje modułu',
    'Katalog akcji bieżącego modułu z rejestru rdzenia (action.list)',
    'Wypisz akcje dostępne w module tego okna i wskaż, którą warto wykonać teraz.',
  ),
  narzedzie(
    'wspolne.stan-okna',
    'Stan okna',
    'Parametry wykonania okna: moduł, kanał modelu, katalogi, zasięg (window.state.get)',
    'Podaj stan tego okna komunikacji: moduł, kanał modelu, katalogi robocze, ' +
      'zasięg wykonania i tryb uprawnień.',
  ),
  narzedzie(
    'wspolne.przenies-kontekst',
    'Przenieś kontekst',
    'Przeniesienie kompletu kontekstu okna do innego modułu (context.transfer)',
    'Przenieś komplet kontekstu tego okna do modułu: ',
  ),
];

/**
 * Profil modułu wraz z dołączonym arsenałem wspólnym.
 *
 * `postac` pominięta znaczy „postaci tego modułu nie ustalono" — profil dostaje
 * wtedy `POSTAC_DOMYSLNA`. Okna obowiązkowe spoza `oknaKontekstu` są odrzucane:
 * profil nie może wymagać okna, którego sam nie wymienia.
 */
export function profil(
  kod: string,
  nazwa: string,
  przeznaczenie: string,
  oknaKontekstu: readonly string[],
  narzedzia: readonly NarzedziePromptu[],
  postac: PostacModulu = POSTAC_DOMYSLNA,
): ProfilModulu {
  return {
    kod,
    nazwa,
    przeznaczenie,
    oknaKontekstu,
    narzedzia: [...narzedzia, ...NARZEDZIA_WSPOLNE],
    postacRozmowy: postac.postacRozmowy,
    granicaOkien: granicaWSuficie(postac.granicaOkien),
    pamiecSesyjna: postac.pamiecSesyjna,
    oknaObowiazkowe: (postac.oknaObowiazkowe ?? []).filter((o) => oknaKontekstu.includes(o)),
  };
}
