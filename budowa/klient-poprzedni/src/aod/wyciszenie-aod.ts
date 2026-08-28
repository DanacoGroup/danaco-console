import { konieCzasuWyciszenia, WyciszenieCzasowe } from './progi-aod';
import { PowodDecyzji } from './rozpoznanie-decyzji';

/**
 * Wyciszenie Always On Display: plik trzyma trzy rodzaje jako stan (czasowe, kontekstowe,
 * klasy zdarzeń) i rozstrzyga, które obejmuje sugestię. Poniżej klasa zdarzeń wyzwalających
 * nakładkę.
 */
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

/**
 * Nazwa klasy zdarzeń widziana przez Operatora w menu wyciszeń — pełne brzmienie
 * zamiast identyfikatora technicznego.
 */
export const NAZWY_KLAS: Readonly<Record<KlasaZdarzen, string>> = {
  [KlasaZdarzen.StanPetliWykonawczej]: 'stan pętli wykonawczej',
  [KlasaZdarzen.StanKolejkiZadan]: 'stan kolejki zadań',
  [KlasaZdarzen.WynikKontroliJakosci]: 'wynik kontroli jakości',
  [KlasaZdarzen.ZdarzeniaModulow]: 'zdarzenia modułów',
  [KlasaZdarzen.Harmonogram]: 'harmonogram',
  [KlasaZdarzen.KontekstPracyOperatora]: 'kontekst pracy Operatora',
};

/**
 * Zdarzenie wyzwalające klasy — pełny opis, widoczny przy tej klasie w menu
 * wyciszeń zdarzeń Operatora.
 */
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
 * Klasa zdarzeń każdego powodu rozpoznanego przez nakładkę: reguła rozpoznania
 * (`rozpoznanie-decyzji.ts`) daje pięć powodów, każdy wpada w jedną z dwóch klas zdarzeń.
 */
export const KLASA_POWODU: Readonly<Record<PowodDecyzji, KlasaZdarzen>> = {
  [PowodDecyzji.BiegStanal]: KlasaZdarzen.StanPetliWykonawczej,
  [PowodDecyzji.Wstrzymany]: KlasaZdarzen.StanPetliWykonawczej,
  [PowodDecyzji.Zatrzymany]: KlasaZdarzen.StanKolejkiZadan,
  [PowodDecyzji.Usterka]: KlasaZdarzen.StanKolejkiZadan,
  [PowodDecyzji.BezRuchu]: KlasaZdarzen.StanKolejkiZadan,
};

/**
 * Klasy, których zdarzenia nakładka dziś rozpoznaje: pozostałe klasy nie mają jeszcze
 * nośnika sygnału w kontrakcie, więc ich wyciszenie zadziała, gdy sygnał wejdzie.
 */
export const KLASY_ROZPOZNAWANE: ReadonlySet<KlasaZdarzen> = new Set(
  Object.values(KLASA_POWODU),
);

/**
 * Zakres wyciszenia kontekstowego: bieżący moduł albo bieżąca karta sesji, do której
 * sugestia dziś należy.
 */
export const ZakresKontekstu = {
  Modul: 'modul',
  KartaSesji: 'karta-sesji',
} as const;
export type ZakresKontekstu = (typeof ZakresKontekstu)[keyof typeof ZakresKontekstu];

/**
 * Nazwa zakresu widziana przez Operatora w menu wyciszeń — moduł albo karta sesji,
 * zapisana pełnym słowem.
 */
export const NAZWY_ZAKRESOW: Readonly<Record<ZakresKontekstu, string>> = {
  [ZakresKontekstu.Modul]: 'moduł',
  [ZakresKontekstu.KartaSesji]: 'karta sesji',
};

/**
 * Wyciszenie czasowe: milczą wszystkie sugestie do wskazanej chwili, niezależnie od
 * modułu i klasy zdarzenia.
 */
export interface WyciszenieCzasem {
  rodzaj: 'czasowe';
  /** Wartość czasu wybrana z menu: kwadrans, godzina, do końca dnia. */
  czas: WyciszenieCzasowe;
  /** Chwila, do której wyciszenie trwa, w milisekundach epoki. */
  doChwili: number;
}

/**
 * Wyciszenie kontekstowe: milczą sugestie jednego modułu albo jednej karty sesji,
 * pozostałe działają bez zmian.
 */
export interface WyciszenieKontekstu {
  rodzaj: 'kontekstowe';
  zakres: ZakresKontekstu;
  /** Identyfikator modułu albo karty sesji, którego wyciszenie dotyczy. */
  wartosc: string;
  /** Nazwa bytu widziana przez Operatora — pełna, nie kod. */
  nazwa: string;
}

/**
 * Wyciszenie klasy zdarzeń: milczą sugestie jednej klasy zdarzenia wyzwalającego,
 * pozostałe klasy działają bez zmian.
 */
export interface WyciszenieKlasy {
  rodzaj: 'klasa-zdarzen';
  klasa: KlasaZdarzen;
}

/** Jedno wyciszenie czynne — dokładnie jeden z trzech rodzajów: czasowe, kontekstowe albo klasy zdarzeń. */
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
 * Zdanie mówiące, co jest wyciszone i do kiedy: cisza bez tej informacji jest gorsza
 * niż brak wyciszenia.
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
 * Sugestia opisana tym, co reguła wyciszenia musi o niej wiedzieć: moduł i karta sesji
 * są opcjonalne, bo nie każda sugestia je zna, a sugestia bez modułu nie wpada
 * w wyciszenie modułu.
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
 * Które wyciszenie obejmuje tę sugestię; `null`, gdy żadne. Funkcja nie zna wyjątku
 * wagi krytycznej — rozstrzyga to reguła ujawniania w `tryb-obecnosci.ts`.
 */
export function wyciszenieObejmujace(
  czynne: readonly WyciszenieCzynne[],
  opis: OpisSugestiiWobecWyciszenia,
): WyciszenieCzynne | null {
  for (const wyciszenie of czynne) {
    switch (wyciszenie.rodzaj) {
      case 'czasowe':
        // Wyciszenie czasowe obejmuje wszystkie sugestie, niezależnie od modułu i klasy.
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

/**
 * Magazyn stanu wyciszeń — podawany wołaczowi, nie brany z globalnej przestrzeni
 * na sztywno w kodzie klienta.
 */
export interface MagazynWyciszen {
  getItem(klucz: string): string | null;
  setItem(klucz: string, wartosc: string): void;
}

/**
 * Klucz zapisu w magazynie — jeden na stanowisko, bo funkcja wyciszeń jest globalna,
 * nie zależna od konta.
 */
export const KLUCZ_ZAPISU_WYCISZEN = 'danaco.aod.wyciszenia';

/**
 * Wykaz wyciszeń czynnych wraz z czynnościami, które go zmieniają — jedyne miejsce
 * dostępu do stanu wyciszeń.
 */
export interface StanWyciszen {
  /** Wyciszenia czynne teraz — przeterminowane czasowe znosi się samo, bez odklikiwania. */
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

/**
 * Magazyn domyślny: zapis miejscowy przeglądarki, gdy jest dostępny, inaczej brak
 * trwałości między sesjami.
 */
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
      // Zapis miejscowy bywa wyłączony ustawieniem przeglądarki — wyciszenie działa dalej bez trwałości.
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
 * Odczyt wykazu z magazynu: każda pozycja jest sprawdzana osobno, zapis starszy,
 * cudzy albo uszkodzony daje wykaz uboższy, nie wyjątek.
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

/**
 * Jedna pozycja zapisu sprowadzona do wyciszenia; `null` przy każdej niezgodności
 * kształtu albo wartości pola.
 */
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
