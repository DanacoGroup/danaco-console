import { konieCzasuWyciszenia, WyciszenieCzasowe } from './progi-aod';
import { PowodDecyzji } from './rozpoznanie-decyzji';

/**
 * Pięć rodzajów wyciszenia Always On Display — rozdz. 3.5 opracowania
 * `docs/funkcje-globalne/always-on-display.md`.
 *
 * Tabela rozdz. 3.5 wymienia pięć wierszy i ten plik trzyma trzy z nich jako
 * stan (wyciszenie czasowe, kontekstowe, klasy zdarzeń); czwarty — tryb cichy —
 * jest trybem obecności i mieszka w `tryb-obecnosci.ts`; piąty — wyjątek wagi
 * krytycznej — nie jest wyciszeniem, tylko regułą przebijającą wszystkie
 * pozostałe, i stoi w regule ujawniania.
 *
 * Plik nie dotyka dokumentu i nie zna kanału. Trzyma wykaz wyciszeń czynnych
 * i rozstrzyga z niego jedno pytanie: KTÓRE wyciszenie obejmuje tę sugestię.
 * Odpowiedzią jest samo wyciszenie, nie „prawda/fałsz" — bo Operator ma
 * usłyszeć, co dokładnie milczy i do kiedy, a nie wyłącznie że milczy.
 *
 * GDZIE TEN STAN MIESZKA. Kontrakt nie niesie wyciszenia nakładki ani jednym
 * polem: `grep -i wycisz shared/contract.go` znajduje wyłącznie wyciszenie
 * uczestnika tury w module Roundtable i wyciszone wyzwolenia reguł alarmowych —
 * nic z rodziny `aod.*`. Stan wyciszenia jest więc dziś stanem okna. Magazyn
 * jest PODAWANY wołaczowi (`MagazynWyciszen`), a nie brany z globalnej
 * przestrzeni na sztywno — wzorem `moduly/studio/widok-nastawy-operatora.ts` —
 * żeby przełożenie zapisu na rdzeń nie ruszyło ani jednego wołacza. Braki
 * kontraktu nazywa wprost `wyciszenie-braki-kontraktu.ts`.
 */

/** Klasa zdarzeń wyzwalających — rozdz. 3.2, kolumna „Klasa". */
export const KlasaZdarzen = {
  /** Pętla wstrzymana, przerwana, zadanie ponawiane powyżej progu, zlecenie oczekujące. */
  StanPetliWykonawczej: 'stan-petli-wykonawczej',
  /** Kolejka zatrzymana, zadanie w błędzie, oczekiwanie dłużej niż próg, wypełnienie powyżej progu. */
  StanKolejkiZadan: 'stan-kolejki-zadan',
  /** Negatywny wynik kontroli, powtórzone niepowodzenie, rozbieżność wyniku ze zleceniem. */
  WynikKontroliJakosci: 'wynik-kontroli-jakosci',
  /** Zakończenie długiego zadania, wynik do przeglądu, konflikt zasobu, brak konfiguracji. */
  ZdarzeniaModulow: 'zdarzenia-modulow',
  /** Termin reguły czasowej, cykliczne uruchomienie automatyki, przypomnienie Operatora. */
  Harmonogram: 'harmonogram',
  /** Powtarzalność czynności, moduł bez wykonawcy, zasób powiązany z zadaniem oczekującym. */
  KontekstPracyOperatora: 'kontekst-pracy-operatora',
} as const;
export type KlasaZdarzen = (typeof KlasaZdarzen)[keyof typeof KlasaZdarzen];

/** Nazwa klasy widziana przez Operatora — brzmienie wiersza tabeli rozdz. 3.2. */
export const NAZWY_KLAS: Readonly<Record<KlasaZdarzen, string>> = {
  [KlasaZdarzen.StanPetliWykonawczej]: 'stan pętli wykonawczej',
  [KlasaZdarzen.StanKolejkiZadan]: 'stan kolejki zadań',
  [KlasaZdarzen.WynikKontroliJakosci]: 'wynik kontroli jakości',
  [KlasaZdarzen.ZdarzeniaModulow]: 'zdarzenia modułów',
  [KlasaZdarzen.Harmonogram]: 'harmonogram',
  [KlasaZdarzen.KontekstPracyOperatora]: 'kontekst pracy Operatora',
};

/** Zdarzenie wyzwalające klasy — rozdz. 3.2, kolumna „Zdarzenie wyzwalające". */
export const ZDARZENIA_KLAS: Readonly<Record<KlasaZdarzen, string>> = {
  [KlasaZdarzen.StanPetliWykonawczej]:
    'Pętla wstrzymana, pętla przerwana, zadanie ponawiane powyżej progu, zlecenie oczekujące ' +
    'na zatwierdzenie Użytkownika.',
  [KlasaZdarzen.StanKolejkiZadan]:
    'Kolejka zatrzymana, zadanie w stanie błędu, zadanie oczekujące dłużej niż próg czasu, ' +
    'kolejka wypełniona powyżej progu.',
  [KlasaZdarzen.WynikKontroliJakosci]:
    'Negatywny wynik kontroli, powtórzone niepowodzenie tego samego zadania, rozbieżność ' +
    'wyniku ze zleceniem.',
  [KlasaZdarzen.ZdarzeniaModulow]:
    'Zakończenie długiego zadania modułu, dostępny wynik do przeglądu, konflikt zasobu, brak ' +
    'konfiguracji potrzebnej do wykonania polecenia.',
  [KlasaZdarzen.Harmonogram]:
    'Nadejście terminu reguły czasowej, cykliczne uruchomienie automatyki, przypomnienie ' +
    'ustanowione przez Operatora.',
  [KlasaZdarzen.KontekstPracyOperatora]:
    'Powtarzalność ręcznie wykonywanej czynności, praca w module bez skonfigurowanego ' +
    'wykonawcy, otwarty zasób powiązany z zadaniem oczekującym.',
};

/**
 * Klasa zdarzeń każdego powodu rozpoznanego przez nakładkę.
 *
 * Rozdz. 3.2 przypisuje „pętlę wstrzymaną i przerwaną" klasie stanu pętli
 * wykonawczej, a „kolejkę zatrzymaną, zadanie w stanie błędu i zadanie
 * oczekujące dłużej niż próg" — klasie stanu kolejki zadań. Reguła rozpoznania
 * (`rozpoznanie-decyzji.ts`) daje pięć powodów i każdy z nich wpada w jedną
 * z tych dwóch klas.
 */
export const KLASA_POWODU: Readonly<Record<PowodDecyzji, KlasaZdarzen>> = {
  [PowodDecyzji.BiegStanal]: KlasaZdarzen.StanPetliWykonawczej,
  [PowodDecyzji.Wstrzymany]: KlasaZdarzen.StanPetliWykonawczej,
  [PowodDecyzji.Zatrzymany]: KlasaZdarzen.StanKolejkiZadan,
  [PowodDecyzji.Usterka]: KlasaZdarzen.StanKolejkiZadan,
  [PowodDecyzji.BezRuchu]: KlasaZdarzen.StanKolejkiZadan,
};

/**
 * Klasy, których zdarzenia nakładka dziś rozpoznaje.
 *
 * Pozostałe cztery klasy rozdz. 3.2 nie mają dziś w kontrakcie nośnika sygnału
 * (nie ma zdarzenia kontroli jakości, harmonogramu ani powtarzalności czynności
 * Operatora), więc ich wyciszenie zapisze się i zadziała z chwilą, w której
 * sygnał wejdzie. Menu mówi to wprost, zamiast udawać, że wycisza coś, co i tak
 * milczy.
 */
export const KLASY_ROZPOZNAWANE: ReadonlySet<KlasaZdarzen> = new Set(
  Object.values(KLASA_POWODU),
);

/** Zakres wyciszenia kontekstowego — rozdz. 3.5: bieżący moduł albo bieżąca karta sesji. */
export const ZakresKontekstu = {
  Modul: 'modul',
  KartaSesji: 'karta-sesji',
} as const;
export type ZakresKontekstu = (typeof ZakresKontekstu)[keyof typeof ZakresKontekstu];

/** Nazwa zakresu widziana przez Operatora. */
export const NAZWY_ZAKRESOW: Readonly<Record<ZakresKontekstu, string>> = {
  [ZakresKontekstu.Modul]: 'moduł',
  [ZakresKontekstu.KartaSesji]: 'karta sesji',
};

/** Wyciszenie czasowe — rozdz. 3.5, wiersz pierwszy. */
export interface WyciszenieCzasem {
  rodzaj: 'czasowe';
  /** Wartość czasu wybrana z menu: kwadrans, godzina, do końca dnia. */
  czas: WyciszenieCzasowe;
  /** Chwila, do której wyciszenie trwa, w milisekundach epoki. */
  doChwili: number;
}

/** Wyciszenie kontekstowe — rozdz. 3.5, wiersz drugi. */
export interface WyciszenieKontekstu {
  rodzaj: 'kontekstowe';
  zakres: ZakresKontekstu;
  /** Identyfikator modułu albo karty sesji, którego wyciszenie dotyczy. */
  wartosc: string;
  /** Nazwa bytu widziana przez Operatora — pełna, nie kod. */
  nazwa: string;
}

/** Wyciszenie klasy zdarzeń — rozdz. 3.5, wiersz trzeci. */
export interface WyciszenieKlasy {
  rodzaj: 'klasa-zdarzen';
  klasa: KlasaZdarzen;
}

/** Jedno wyciszenie czynne. */
export type WyciszenieCzynne = WyciszenieCzasem | WyciszenieKontekstu | WyciszenieKlasy;

/**
 * Tożsamość wyciszenia — podstawa zniesienia jednym kliknięciem.
 *
 * Wyciszenie czasowe jest jedno: włączenie drugiego czasu zastępuje poprzednie,
 * bo dwa czasy naraz nie dałyby Operatorowi jednej odpowiedzi na pytanie
 * „do kiedy".
 */
export function kluczWyciszenia(wyciszenie: WyciszenieCzynne): string {
  switch (wyciszenie.rodzaj) {
    case 'czasowe':
      return 'czasowe';
    case 'kontekstowe':
      return `kontekstowe:${wyciszenie.zakres}:${wyciszenie.wartosc}`;
    case 'klasa-zdarzen':
      return `klasa-zdarzen:${wyciszenie.klasa}`;
  }
}

/**
 * Zdanie mówiące, CO jest wyciszone i DO KIEDY.
 *
 * Zasada zlecenia: cisza, po której Operator nie wie, że coś jest wyłączone,
 * jest gorsza od braku wyciszenia. Dlatego każde wyciszenie ma tu swoje zdanie,
 * a wyciszenie bez końca czasowego mówi wprost, że trwa do zniesienia ręką
 * Operatora.
 */
export function zdanieWyciszenia(wyciszenie: WyciszenieCzynne): string {
  switch (wyciszenie.rodzaj) {
    case 'czasowe':
      return (
        `Wyciszenie czasowe: wszystkie sugestie gromadzą się bez dymka i bez plakietki do ` +
        `${new Date(wyciszenie.doChwili).toLocaleTimeString()}.`
      );
    case 'kontekstowe':
      return (
        `Wyciszenie kontekstowe: sugestie, których zakres to ` +
        `${NAZWY_ZAKRESOW[wyciszenie.zakres]} „${wyciszenie.nazwa}", nie ujawniają się. ` +
        'Pozostałe zachowują pełne działanie. Trwa do zniesienia ręką Operatora.'
      );
    case 'klasa-zdarzen':
      return (
        `Wyciszenie klasy zdarzeń: klasa „${NAZWY_KLAS[wyciszenie.klasa]}" nie tworzy sugestii. ` +
        'Pozostałe klasy działają bez zmian. Trwa do zniesienia ręką Operatora.'
      );
  }
}

/**
 * Sugestia opisana tym, co reguła wyciszenia musi o niej wiedzieć.
 *
 * Moduł i karta sesji są opcjonalne, bo nie każda sugestia je zna: telemetria
 * rdzenia niesie okno i sesję, modułu nie niesie wcale, a `AodSuggestion`
 * kontraktu niesie samo okno. Sugestia bez modułu NIE WPADA w wyciszenie
 * modułu — milczenie sugestii, o której nie wiemy, czy dotyczy wyciszonego
 * bytu, byłoby ciszą bez podstawy.
 */
export interface OpisSugestiiWobecWyciszenia {
  /** Klasa zdarzenia wyzwalającego; pominięta znaczy „klasa nierozpoznana". */
  klasa?: KlasaZdarzen;
  /** Moduł, którego sugestia dotyczy. */
  idModulu?: string;
  /** Karta sesji, której sugestia dotyczy. */
  idSesji?: string;
}

/**
 * Które wyciszenie obejmuje tę sugestię; `null`, gdy żadne.
 *
 * Funkcja NIE zna wyjątku wagi krytycznej i nie ma go znać: wyjątek nie znosi
 * wyciszenia, tylko przepuszcza sugestię plakietką mimo niego. Rozstrzyga to
 * reguła ujawniania w `tryb-obecnosci.ts`, w jednym miejscu dla wszystkich
 * trzech rodzajów.
 */
export function wyciszenieObejmujace(
  czynne: readonly WyciszenieCzynne[],
  opis: OpisSugestiiWobecWyciszenia,
): WyciszenieCzynne | null {
  for (const wyciszenie of czynne) {
    switch (wyciszenie.rodzaj) {
      case 'czasowe':
        // Wyciszenie czasowe obejmuje wszystko — rozdz. 3.5, kolumna „Zakres".
        return wyciszenie;
      case 'kontekstowe': {
        const bytSugestii =
          wyciszenie.zakres === ZakresKontekstu.Modul ? opis.idModulu : opis.idSesji;
        if (bytSugestii !== undefined && bytSugestii === wyciszenie.wartosc) return wyciszenie;
        break;
      }
      case 'klasa-zdarzen':
        if (opis.klasa === wyciszenie.klasa) return wyciszenie;
        break;
    }
  }
  return null;
}

/** Magazyn stanu wyciszeń — podawany wołaczowi, nie brany na sztywno. */
export interface MagazynWyciszen {
  getItem(klucz: string): string | null;
  setItem(klucz: string, wartosc: string): void;
}

/** Klucz zapisu. Jeden na stanowisko — funkcja jest globalna. */
export const KLUCZ_ZAPISU_WYCISZEN = 'danaco.aod.wyciszenia';

/** Wykaz wyciszeń czynnych wraz z czynnościami, które go zmieniają. */
export interface StanWyciszen {
  /**
   * Wyciszenia czynne w tej chwili, w kolejności włączenia.
   *
   * Wyciszenie czasowe przeterminowane znosi się samo i nie wchodzi do wykazu —
   * Operator nie ma go odklikiwać.
   */
  czynne(teraz: number): readonly WyciszenieCzynne[];
  /** Wyciszenie czasowe albo `null`; wyłącznie ono ma koniec liczony zegarem. */
  czasowe(teraz: number): WyciszenieCzasem | null;
  /** Włącza wyciszenie czasowe; drugie wywołanie zastępuje poprzednie. */
  wyciszCzasem(czas: WyciszenieCzasowe, teraz: number): void;
  /** Przełącza wyciszenie czasowe tej wartości — podstawa `Ctrl/Cmd + Shift + M`. */
  przelaczCzas(czas: WyciszenieCzasowe, teraz: number): void;
  /** Włącza wyciszenie kontekstowe wskazanego bytu; powtórzenie nic nie zmienia. */
  wyciszKontekst(zakres: ZakresKontekstu, wartosc: string, nazwa: string): void;
  /** Przełącza wyciszenie kontekstowe wskazanego bytu. */
  przelaczKontekst(zakres: ZakresKontekstu, wartosc: string, nazwa: string): void;
  /** Przełącza wyciszenie klasy zdarzeń. */
  przelaczKlase(klasa: KlasaZdarzen): void;
  czyKlasaWyciszona(klasa: KlasaZdarzen): boolean;
  czyKontekstWyciszony(zakres: ZakresKontekstu, wartosc: string): boolean;
  /** Znosi jedno wyciszenie — jedno kliknięcie, bez pytania o potwierdzenie. */
  znies(klucz: string): void;
  /** Znosi wszystkie wyciszenia naraz — również jednym kliknięciem. */
  zniesWszystkie(): void;
  /** Zgłasza obserwatora zmiany; zwraca odsubskrybowanie. */
  obserwuj(sluchacz: () => void): () => void;
}

/** Magazyn domyślny: zapis miejscowy przeglądarki, gdy jest dostępny. */
export function magazynWyciszenDomyslny(): MagazynWyciszen | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}

/**
 * Buduje stan wyciszeń nad podanym magazynem.
 *
 * @param magazyn pominięty znaczy zapis miejscowy przeglądarki, `null` znaczy
 *   „bez zapisu" — stan żyje wtedy w pamięci okna i nie przeżywa przeładowania.
 */
export function utworzStanWyciszen(
  magazyn: MagazynWyciszen | null = magazynWyciszenDomyslny(),
): StanWyciszen {
  let wykaz: WyciszenieCzynne[] = wczytajWyciszenia(magazyn);

  const sluchacze = new Set<() => void>();

  function rozglos(): void {
    for (const sluchacz of sluchacze) sluchacz();
  }

  function zapisz(): void {
    if (magazyn === null) return;
    try {
      magazyn.setItem(KLUCZ_ZAPISU_WYCISZEN, JSON.stringify(wykaz));
    } catch {
      // Zapis miejscowy bywa wyłączony ustawieniem przeglądarki. Wyciszenie
      // działa dalej, tracąc wyłącznie pamięć między przeładowaniami strony.
    }
  }

  /** Zdejmuje wyciszenia czasowe, których czas upłynął; zwraca `true`, gdy zdjęto. */
  function odsiejPrzeterminowane(teraz: number): boolean {
    const przed = wykaz.length;
    wykaz = wykaz.filter(
      (wyciszenie) => wyciszenie.rodzaj !== 'czasowe' || wyciszenie.doChwili > teraz,
    );
    if (wykaz.length === przed) return false;
    zapisz();
    return true;
  }

  function zmien(nowy: WyciszenieCzynne[]): void {
    wykaz = nowy;
    zapisz();
    rozglos();
  }

  function bezKlucza(klucz: string): WyciszenieCzynne[] {
    return wykaz.filter((wyciszenie) => kluczWyciszenia(wyciszenie) !== klucz);
  }

  return {
    czynne(teraz) {
      if (odsiejPrzeterminowane(teraz)) rozglos();
      return wykaz;
    },

    czasowe(teraz) {
      if (odsiejPrzeterminowane(teraz)) rozglos();
      return wykaz.find((wyciszenie) => wyciszenie.rodzaj === 'czasowe') ?? null;
    },

    wyciszCzasem(czas, teraz) {
      zmien([
        ...bezKlucza('czasowe'),
        { rodzaj: 'czasowe', czas, doChwili: konieCzasuWyciszenia(czas, teraz) },
      ]);
    },

    przelaczCzas(czas, teraz) {
      odsiejPrzeterminowane(teraz);
      const biezace = wykaz.find((wyciszenie) => wyciszenie.rodzaj === 'czasowe');
      if (biezace !== undefined) {
        zmien(bezKlucza('czasowe'));
        return;
      }
      zmien([...wykaz, { rodzaj: 'czasowe', czas, doChwili: konieCzasuWyciszenia(czas, teraz) }]);
    },

    wyciszKontekst(zakres, wartosc, nazwa) {
      const klucz = `kontekstowe:${zakres}:${wartosc}`;
      zmien([...bezKlucza(klucz), { rodzaj: 'kontekstowe', zakres, wartosc, nazwa }]);
    },

    przelaczKontekst(zakres, wartosc, nazwa) {
      const klucz = `kontekstowe:${zakres}:${wartosc}`;
      const bylo = wykaz.some((wyciszenie) => kluczWyciszenia(wyciszenie) === klucz);
      zmien(bylo ? bezKlucza(klucz) : [...wykaz, { rodzaj: 'kontekstowe', zakres, wartosc, nazwa }]);
    },

    przelaczKlase(klasa) {
      const klucz = `klasa-zdarzen:${klasa}`;
      const bylo = wykaz.some((wyciszenie) => kluczWyciszenia(wyciszenie) === klucz);
      zmien(bylo ? bezKlucza(klucz) : [...wykaz, { rodzaj: 'klasa-zdarzen', klasa }]);
    },

    czyKlasaWyciszona: (klasa) =>
      wykaz.some(
        (wyciszenie) => wyciszenie.rodzaj === 'klasa-zdarzen' && wyciszenie.klasa === klasa,
      ),

    czyKontekstWyciszony: (zakres, wartosc) =>
      wykaz.some(
        (wyciszenie) =>
          wyciszenie.rodzaj === 'kontekstowe' &&
          wyciszenie.zakres === zakres &&
          wyciszenie.wartosc === wartosc,
      ),

    znies(klucz) {
      if (!wykaz.some((wyciszenie) => kluczWyciszenia(wyciszenie) === klucz)) return;
      zmien(bezKlucza(klucz));
    },

    zniesWszystkie() {
      if (wykaz.length === 0) return;
      zmien([]);
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
}

/**
 * Odczyt wykazu z magazynu.
 *
 * Każda pozycja jest sprawdzana osobno: zapis starszy, cudzy albo uszkodzony
 * ma dać wykaz uboższy, a nie wyjątek i nie ciszę bez podstawy. Pozycja
 * nieznanego kształtu wypada — wyciszenie, którego nie umiemy nazwać
 * Operatorowi, nie ma prawa wyciszać.
 */
function wczytajWyciszenia(magazyn: MagazynWyciszen | null): WyciszenieCzynne[] {
  if (magazyn === null) return [];

  let zapis: unknown;
  try {
    const surowy = magazyn.getItem(KLUCZ_ZAPISU_WYCISZEN);
    if (surowy === null || surowy === '') return [];
    zapis = JSON.parse(surowy);
  } catch {
    return [];
  }

  if (!Array.isArray(zapis)) return [];

  const teraz = Date.now();
  const wykaz: WyciszenieCzynne[] = [];

  for (const pozycja of zapis) {
    const wyciszenie = odczytajWyciszenie(pozycja, teraz);
    if (wyciszenie !== null) wykaz.push(wyciszenie);
  }

  return wykaz;
}

/** Jedna pozycja zapisu sprowadzona do wyciszenia; `null` przy każdej niezgodności. */
function odczytajWyciszenie(pozycja: unknown, teraz: number): WyciszenieCzynne | null {
  if (typeof pozycja !== 'object' || pozycja === null) return null;
  const pola = pozycja as Record<string, unknown>;

  switch (pola['rodzaj']) {
    case 'czasowe': {
      const czas = pola['czas'];
      const doChwili = pola['doChwili'];
      if (!czyCzasWyciszenia(czas)) return null;
      if (typeof doChwili !== 'number' || doChwili <= teraz) return null;
      return { rodzaj: 'czasowe', czas, doChwili };
    }
    case 'kontekstowe': {
      const zakres = pola['zakres'];
      const wartosc = pola['wartosc'];
      const nazwa = pola['nazwa'];
      if (!czyZakresKontekstu(zakres)) return null;
      if (typeof wartosc !== 'string' || wartosc === '') return null;
      return {
        rodzaj: 'kontekstowe',
        zakres,
        wartosc,
        nazwa: typeof nazwa === 'string' && nazwa !== '' ? nazwa : wartosc,
      };
    }
    case 'klasa-zdarzen': {
      const klasa = pola['klasa'];
      if (!czyKlasaZdarzen(klasa)) return null;
      return { rodzaj: 'klasa-zdarzen', klasa };
    }
    default:
      return null;
  }
}

function czyCzasWyciszenia(wartosc: unknown): wartosc is WyciszenieCzasowe {
  return (
    typeof wartosc === 'string' &&
    (Object.values(WyciszenieCzasowe) as string[]).includes(wartosc)
  );
}

function czyZakresKontekstu(wartosc: unknown): wartosc is ZakresKontekstu {
  return (
    typeof wartosc === 'string' && (Object.values(ZakresKontekstu) as string[]).includes(wartosc)
  );
}

function czyKlasaZdarzen(wartosc: unknown): wartosc is KlasaZdarzen {
  return (
    typeof wartosc === 'string' && (Object.values(KlasaZdarzen) as string[]).includes(wartosc)
  );
}
