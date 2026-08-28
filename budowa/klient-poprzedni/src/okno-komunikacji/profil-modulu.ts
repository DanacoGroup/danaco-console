/**
 * Profil modułu niesie to, czym okno komunikacji różni się w każdym module: wygląd, możliwości, narzędzia i kontekst, przestawiane przy zmianie modułu bez dotykania historii wątku rozmowy.
 */

import { LICZBA_MAX } from '../okna-rownolegle/identyfikatory';

/** Pojedyncza pozycja paska narzędzi promptu, niosąca kod, etykietę, podpowiedź oraz treść wstawianą do pola wypowiedzi operatora. */
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
 * Postać, w jakiej moduł prowadzi rozmowę — zwykłe okno czatu, pływający awatar z dymkiem głosowym, albo brak rozmowy, gdy moduł nią nie dysponuje.
 */
export type PostacRozmowy = 'okno' | 'dymek-glosowy' | 'brak';

/**
 * Granica okien równoległych jednego modułu zawęża sufit platformy, nigdy go nie rozszerza, bo sufit pozostaje twardy niezależnie od profilu.
 */
export type GranicaOkien = 1 | 2 | 4;

/** Granice dopuszczalne, rosnąco — jedyne wartości, jakie może przyjąć pole granicy okien równoległych modułu. */
export const GRANICE: readonly GranicaOkien[] = [1, 2, 4];

/**
 * Granica modułu, dla którego liczby okien nie ustalono, istnieje po to, żeby domysł nie udawał podjętej decyzji o limicie.
 */
export const GRANICA_NIEPODANA: GranicaOkien = 4;

/** Właściwości rozmowy w module, obejmujące postać rozmowy, granicę okien i trwałość wątku, bez pasków i narzędzi. */
export interface PostacModulu {
  /** Czy i w jakiej postaci moduł prowadzi rozmowę. */
  postacRozmowy: PostacRozmowy;
  /** Ile okien równoległych moduł wykorzystuje; nigdy ponad `LICZBA_MAX`. */
  granicaOkien: GranicaOkien;
  /** Rozmowa Agents nie przeżywa zamknięcia okna — ginie przy zmianie testowanego agenta. */
  pamiecSesyjna: boolean;
  /** Okna towarzyszące obowiązkowe są podzbiorem okien kontekstu; reszta okien kontekstu to okna możliwe. */
  oknaObowiazkowe?: readonly string[];
}

/**
 * Postać stanu zastanego dla profilu bez własnej postaci: zwykłe okno, pełny sufit platformy, pamięć rozmowy zachowana.
 */
export const POSTAC_DOMYSLNA: PostacModulu = {
  postacRozmowy: 'okno',
  granicaOkien: GRANICA_NIEPODANA,
  pamiecSesyjna: true,
  oknaObowiazkowe: [],
};

/** Zestaw właściwości okna komunikacji w jednym module, obejmujący nazwę, przeznaczenie, pasek narzędzi i okna operacyjne. */
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
 * Największa granica mieszcząca się jednocześnie w żądaniu modułu i w suficie platformy, bo sufit zawęża i nigdy go nie podnosi.
 */
export function granicaWSuficie(zadana: GranicaOkien): GranicaOkien {
  const mieszczace = GRANICE.filter((g) => g <= zadana && g <= LICZBA_MAX);
  return mieszczace.length > 0 ? mieszczace[mieszczace.length - 1] : 1;
}

/**
 * Ile okien rozmowy operator może otworzyć w tym module — zero dla modułu bez rozmowy jest odpowiedzią, nie brakiem odpowiedzi.
 */
export function liczbaOkienRozmowy(p: ProfilModulu): number {
  return p.postacRozmowy === 'brak' ? 0 : p.granicaOkien;
}

/** Pozycja paska narzędzi promptu złożona z pól profilu modułu, gotowa do narysowania w interfejsie okna. */
export function narzedzie(
  kod: string,
  nazwa: string,
  opis: string,
  polecenie: string,
): NarzedziePromptu {
  return { kod, nazwa, opis, polecenie };
}

/**
 * Arsenał wspólny Chat Window to trzy pozycje obecne w każdym module, dostępne nawet modułowi nierozpoznanemu przez rejestr.
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
 * Profil modułu wraz z dołączonym arsenałem wspólnym, gdzie pominięta postać rozmowy dostaje wartość domyślną.
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
