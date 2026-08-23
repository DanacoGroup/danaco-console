import type { SettingDefinition } from '../../../shared/contract';
import {
  odczytajWybor,
  zapiszWybor,
  zastosujMotyw,
  ZDARZENIE_MOTYWU,
  type Motyw,
  type ZmianaMotywu,
} from '../motyw/motyw';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import {
  naZmianeKlucza,
  odczytajNastawe,
  poziomZapisu,
  zapiszNastawe,
} from './zrodlo-nastaw';

/**
 * Most motywu — jedyny właściciel nastawy motywu po stronie klienta.
 *
 * Motyw ma dwa miejsca zapisu: klucz `personalizacja.motyw` w katalogu ustawień
 * rdzenia oraz wybór trzymany przez `motyw/motyw.ts` w `localStorage`
 * i rozgłaszany zdarzeniem DOM, którego słuchają przełączniki paska. Most zwija
 * je w jedną nastawę: prawdą jest rdzeń, a przeglądarka zostaje pamięcią
 * podręczną pierwszego rysunku.
 *
 *   (a) przy podłączeniu czyta `config.get` i stosuje wynik `zastosujMotyw`;
 *   (b) subskrybuje `config.changed` i stosuje zmianę — także tę wykonaną
 *       w drugim oknie albo na drugim urządzeniu, bo rdzeń rozgłasza zdarzenie
 *       do obu połączeń, również do tego, które zapisało;
 *   (c) nasłuchuje `ZDARZENIE_MOTYWU`, więc przełącznik paska jest drugim
 *       sterem tej samej nastawy: jego kliknięcie idzie `config.set` do rdzenia;
 *   (d) wartość pusta (`''`) znaczy preferencję systemu i zostaje trzecim
 *       stanem, bo katalog rdzenia go ma — most jej nie spłaszcza do „jasny".
 *
 * Zapis lokalny zostaje, choć prawdą jest rdzeń: okno rysuje się, zanim rdzeń
 * odpowie na pierwszy `config.get`, a bez niego każde uruchomienie zaczynałoby
 * się mignięciem motywu preferowanego przez system. Most zapisuje lokalnie
 * każdą wartość potwierdzoną przez rdzeń — to odbicie prawdy, nie druga prawda:
 * zapis lokalny nigdy nie jedzie z powrotem do rdzenia jako nastawa.
 *
 * Na kanał przypada jeden most. Dwa mosty na tym samym kanale byłyby dwoma
 * właścicielami jednej nastawy i odbiłyby sobie nawzajem każde zdarzenie, więc
 * `podepnijMostMotywu` oddaje most już podpięty, gdy kanał się zgadza.
 */

/**
 * Klucz nastawy motywu w katalogu rdzenia.
 *
 * Jedyna rzecz z katalogu zaszyta w tym pliku. Etykiet, opcji, poziomów
 * i rodzaju kontrolki most nie zna — czyta je z `settings.definition.list`.
 * Nazwę klucza znać musi, bo most z definicji dotyczy jednej nastawy.
 */
export const KLUCZ_MOTYWU = 'personalizacja.motyw';

/** Wartość nastawy motywu: pusta znaczy „preferencja systemu". */
export type WyborMotywu = '' | Motyw;

export interface MostMotywu {
  /** Wartość bieżąca nastawy — to, co ma stanąć na uchwycie steru. */
  wybor(): WyborMotywu;
  /** Definicja katalogu: etykiety i opcje wyboru; `undefined` przed odczytem. */
  definicja(): SettingDefinition | undefined;
  /** Zdanie odmowy z ostatniego odczytu albo zapisu; puste znaczy „bez odmowy". */
  odmowa(): string;
  /** Odczytuje katalog i wartość z rdzenia, po czym stosuje motyw. */
  odczytaj(): Promise<void>;
  /** Zapisuje wybór do rdzenia. Oddaje zdanie odmowy albo puste przy powodzeniu. */
  ustaw(wybor: WyborMotywu): Promise<string>;
  /** Subskrypcja zmiany — dla każdego miejsca pokazującego tę nastawę. */
  naZmiane(sluchacz: (wybor: WyborMotywu) => void): () => void;
  /** Odpina nasłuch kanału i nasłuch dokumentu. */
  rozlacz(): void;
}

/** Most zbudowany dla danego kanału — jeden na kanał (patrz nagłówek). */
let most: MostMotywu | null = null;
let kanalMostu: Kanal | null = null;

/**
 * Podpina most motywu do kanału i oddaje go.
 *
 * Wołanie z tym samym kanałem oddaje most już stojący — bez drugiego odczytu
 * i bez drugiego nasłuchu. Wołanie z kanałem innym (po ponownym połączeniu
 * z rdzeniem) rozłącza poprzedni i buduje nowy, żeby zapisy nie szły przez
 * transport, którego już nie ma.
 */
export function podepnijMostMotywu(kanal: Kanal): MostMotywu {
  if (most !== null && kanalMostu === kanal) return most;
  if (most !== null) most.rozlacz();
  most = zbudujMostMotywu(kanal);
  kanalMostu = kanal;
  void most.odczytaj();
  return most;
}

/** Zdejmuje most bieżący — droga wyjścia dla sprawdzianów i dla rozłączenia. */
export function odepnijMostMotywu(): void {
  if (most !== null) most.rozlacz();
  most = null;
  kanalMostu = null;
}

function zbudujMostMotywu(kanal: Kanal): MostMotywu {
  /** Wartość potwierdzona przez rdzeń; przed pierwszym odczytem — zapis lokalny. */
  let wybor: WyborMotywu = odczytajWybor() ?? '';
  let definicja: SettingDefinition | undefined;
  let odmowa = '';
  const sluchacze: Array<(wybor: WyborMotywu) => void> = [];

  /**
   * Zapora pętli zwrotnej.
   *
   * `zastosujMotyw` rozgłasza `ZDARZENIE_MOTYWU`, a most tego zdarzenia
   * słucha, żeby wyłapać kliknięcie przełącznika na pasku. Bez zapory własne
   * zastosowanie wartości z rdzenia natychmiast odesłałoby ją z powrotem do
   * rdzenia jako „zmiana Operatora" — zapis w kółko na każde zdarzenie.
   * Rozgłoszenie jest synchroniczne, więc zapora zdejmuje się zaraz po nim.
   */
  let stosujeZRdzenia = false;

  /** Odczyt w toku — żeby zapis nie wyprzedził poznania poziomu zasięgu. */
  let wToku: Promise<void> | null = null;

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz(wybor);
  }

  /**
   * Stosuje wartość: atrybut dokumentu, odbicie lokalne, ogłoszenie słuchaczom.
   *
   * Jedna droga dla obu źródeł: odczytu pierwszego i zdarzenia `config.changed`.
   * Dwie drogi rozjechałyby się przy pierwszej poprawce.
   */
  function zastosuj(nowy: WyborMotywu): void {
    wybor = nowy;
    stosujeZRdzenia = true;
    try {
      zastosujMotyw(nowy === '' ? null : nowy);
    } finally {
      stosujeZRdzenia = false;
    }
    // Odbicie prawdy rdzenia, żeby następne uruchomienie nie migało cudzym
    // motywem, zanim przyjdzie pierwsza odpowiedź (patrz nagłówek).
    zapiszWybor(nowy === '' ? null : nowy);
    oglos();
  }

  /**
   * Wartość z rdzenia sprowadzona do trzech stanów nastawy.
   *
   * Rdzeń oddaje `unknown` — wartość idzie przez bazę jako tekst i kontrakt jej
   * nie zawęża. Wartość spoza katalogu (literówka, wpis ręką w bazie) znaczy
   * „bez wskazania", a nie awarię: nierozpoznane degraduje się do preferencji
   * systemu, nigdy do blokady.
   */
  function jakoWybor(wartosc: unknown): WyborMotywu {
    return wartosc === 'light' || wartosc === 'dark' ? wartosc : '';
  }

  const odsubskrybuj: Odsubskrybuj = naZmianeKlucza(kanal, KLUCZ_MOTYWU, (wartosc) => {
    // `undefined` znaczy `config.reset` — wpis zdjęty, obowiązuje wartość
    // domyślna katalogu. Wartości z wpisu użyć tu nie wolno: przy usunięciu
    // rdzeń niesie w nim wartość właśnie zdjętą (`zrodlo-nastaw.ts`).
    zastosuj(wartosc === undefined ? jakoWybor(definicja?.defaultValue) : jakoWybor(wartosc));
  });

  /** Kliknięcie przełącznika na pasku — drugi ster tej samej nastawy. */
  function odZdarzeniaMotywu(zdarzenie: Event): void {
    if (stosujeZRdzenia) return;
    const zmiana = (zdarzenie as CustomEvent<ZmianaMotywu>).detail;
    const zadany: WyborMotywu = zmiana.wybor ?? '';
    // Zmiana preferencji systemu przy braku wyboru Operatora też idzie tym
    // zdarzeniem (`sledzPreferencjeSystemu`). Nastawy nie zmienia, więc
    // porównanie z wartością bieżącą wycisza ją bez osobnego warunku.
    if (zadany === wybor) return;
    void ustaw(zadany);
  }

  document.addEventListener(ZDARZENIE_MOTYWU, odZdarzeniaMotywu);

  async function odczytaj(): Promise<void> {
    const odczyt = odczytajNastawe(kanal, KLUCZ_MOTYWU).then((wynik) => {
      definicja = wynik.definicja;
      odmowa = wynik.odmowa;
      zastosuj(jakoWybor(wynik.wpis?.value ?? wynik.definicja?.defaultValue));
    });
    wToku = odczyt;
    await odczyt;
    wToku = null;
  }

  async function ustaw(nowy: WyborMotywu): Promise<string> {
    // Poziom zapisu bierze się z definicji katalogu, nigdy z domysłu: rdzeń
    // zapisu na poziomie niedozwolonym nie odrzuca, więc pomyłka byłaby cicha.
    if (definicja === undefined && wToku !== null) await wToku;
    const poziom = poziomZapisu(definicja);
    if (poziom === undefined) {
      odmowa =
        'Nie ma gdzie zapisać motywu: katalog ustawień rdzenia nie oddał definicji ' +
        `klucza „${KLUCZ_MOTYWU}" ani jego dozwolonego poziomu zasięgu.`;
      oglos();
      return odmowa;
    }

    // Wartość stosuje się od razu, nie czekając na odpowiedź: Operator ma
    // zobaczyć skutek kliknięcia natychmiast, a rdzeń i tak przyśle
    // `config.changed`, który wartość potwierdzi albo poprawi.
    zastosuj(nowy);

    const wynik = await zapiszNastawe(kanal, KLUCZ_MOTYWU, nowy, poziom);
    odmowa = wynik.udany ? '' : (wynik.blad?.message ?? 'Rdzeń odmówił zapisu bez podania powodu.');
    if (odmowa !== '') oglos();
    return odmowa;
  }

  return {
    wybor: () => wybor,
    definicja: () => definicja,
    odmowa: () => odmowa,
    odczytaj,
    ustaw,

    naZmiane(sluchacz) {
      sluchacze.push(sluchacz);
      return () => {
        const miejsce = sluchacze.indexOf(sluchacz);
        if (miejsce >= 0) sluchacze.splice(miejsce, 1);
      };
    },

    rozlacz() {
      odsubskrybuj();
      document.removeEventListener(ZDARZENIE_MOTYWU, odZdarzeniaMotywu);
      sluchacze.length = 0;
    },
  };
}
