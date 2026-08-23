import {
  KnownModuleIds,
  NavigationKind,
  type Environment,
  type Module,
} from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/**
 * Widok wykazu środowiska — kształt, w jakim boczna nawigacja czyta odpowiedź
 * rdzenia, wraz z tym, czego kontrakt nie oddaje.
 *
 * Plik nie jest katalogiem modułów: środowiska, moduły i macierz widoczności
 * są sterowane danymi i przychodzą komendami `environment.enter`
 * oraz `module.list`. Nie ma tu ani jednej nazwy modułu, ani jednego wiersza
 * macierzy.
 *
 * Zostały dwie rzeczy, których kontrakt nie niesie:
 *  1. Ikony pozycji. `Module` nie ma pola ikony, bo rysunek jest zasobem
 *     pakietu wizualnego, nie wierszem tabeli. Mapa niżej wiąże kod modułu
 *     z nazwą ikony zestawu; kod nieznany dostaje ikonę zastępczą zamiast
 *     pustego miejsca.
 *  2. Sekcje panelu orkiestracji. Dla środowiska o nawigacji `orchestration`
 *     kontrakt zwraca pusty wykaz modułów (`Environment.moduleCodes` puste,
 *     `environment.enter` oddaje `modules: []`), bo sekcje panelu modułami nie
 *     są. Nie ma ich skąd wziąć z rdzenia, więc stoją tutaj — a że żadna nie ma
 *     modułu, żadna nie ma też zbudowanego widoku (patrz `pozycja.modul`).
 */

/** Klucz środowiska — kody znane kontraktowi, nie literały klienta. */
export type KluczSrodowiska = (typeof KnownModuleIds)[number];

/** Rodzaj wykazu w bocznej nawigacji. */
export type RodzajWykazu = 'moduly' | 'sekcje';

/** Stan wykazu: przed odpowiedzią rdzenia, po niej albo po odmowie. */
export type StanWykazu = 'ladowanie' | 'gotowe' | 'blad';

/** Pojedyncza pozycja bocznej nawigacji: moduł albo sekcja orkiestracji. */
export interface PozycjaModulu {
  /** Klucz stabilny, używany w atrybutach i zdarzeniach; kod modułu z rdzenia. */
  klucz: string;
  /** Nazwa własna, tak jak podaje ją rdzeń. */
  nazwa: string;
  /** Ikona z zestawu `ikony/` — jednolity viewBox, wymienna bez zmiany wymiarów. */
  ikona: NazwaIkony;
  /** Rola pozycji. Dla modułu — opis z rdzenia; pusty, gdy rdzeń go nie podał. */
  opis: string;
  /** Identyfikator modułu w rdzeniu; pusty dla sekcji panelu orkiestracji. */
  modul?: string;
  /** Kody okien operacyjnych modułu wprost z katalogu rdzenia. */
  okna: readonly string[];
}

/** Środowisko wraz z pełnym wykazem pozycji nawigacji. */
export interface Srodowisko {
  klucz: KluczSrodowiska;
  nazwa: string;
  motto: string;
  rodzaj: RodzajWykazu;
  stan: StanWykazu;
  pozycje: PozycjaModulu[];
  /** Treść odmowy rdzenia — wyłącznie przy stanie `blad`. */
  blad?: string;
}

/** Środowisko otwierane, gdy wywołanie nie wskazuje innego. */
export const SRODOWISKO_DOMYSLNE: KluczSrodowiska = KnownModuleIds[0];

/** Ikona pozycji, gdy kod modułu jest rdzeniowi znany, a klientowi nie. */
const IKONA_ZASTEPCZA: NazwaIkony = 'karta-okna';

/**
 * Kod modułu → ikona zestawu. Jedyny powód istnienia mapy: kontrakt nie niesie
 * rysunku. Kod spoza mapy nie jest błędem — dostaje ikonę zastępczą, bo nowy
 * moduł to nowy wiersz w bazie, nie zmiana kodu klienta.
 */
const IKONY_MODULOW: Readonly<Record<string, NazwaIkony>> = {
  studio: 'dokument',
  workspace: 'folder',
  automations: 'automatyzacja',
  browser: 'karta-okna',
  research: 'badanie',
  library: 'biblioteka',
  translate: 'tlumacz',
  roundtable: 'debata',
  design: 'paleta',
  assistant: 'mikrofon',
  terminal: 'terminal',
  developer: 'kod',
  diagnostics: 'diagnostyka',
  apps: 'aplikacje',
  agents: 'agent',
};

/**
 * Klucze sekcji panelu orkiestracji — zestaw domyślny.
 *
 * Wykaz stoi osobno i publicznie, bo znać go musi zarówno boczna nawigacja
 * (żeby narysować pozycje), jak i moduł MultitaskingAI (żeby zbudować
 * powierzchnię sekcji). Dwa wykazy rozjechałyby się przy pierwszej zmianie
 * kolejności, a rozjazd objawiłby się pozycją nawigacji bez treści.
 *
 * Kolejność i widoczność podlegają konfiguracji: układ podsekcji trzyma rdzeń
 * (`panel.sections.*`), a brak ustawienia znaczy wartość domyślną, a nie
 * niedostępność sekcji.
 */
export const KLUCZE_SEKCJI_ORKIESTRACJI = [
  'zespoly',
  'role',
  'kolejki',
  'orkiestracja',
  'harmonogram',
  'monitor',
] as const;

/**
 * Sekcje panelu orkiestracji. Nie są modułami platformy, więc rdzeń ich nie
 * zwraca — i właśnie dlatego żadna nie ma pola `modul`.
 */
const SEKCJE_ORKIESTRACJI: readonly PozycjaModulu[] = [
  sekcja('zespoly', 'Zespoły', 'agenci', 'skład zespołu wykonawców'),
  sekcja('role', 'Role', 'tarcza', 'Coordinator, Executor, Validator'),
  sekcja('kolejki', 'Kolejki', 'warstwy', 'kroki, zależności, licznik obiegów'),
  sekcja('orkiestracja', 'Orkiestracja', 'wezly', 'przydział zleceń między oknami'),
  sekcja('harmonogram', 'Harmonogram i automatyki', 'kalendarz', 'procesy czasu i zdarzenia'),
  sekcja('monitor', 'Monitor procesu', 'monitor', 'przebieg pracy zespołu na żywo'),
];

/** Skrót zapisu sekcji panelu — sekcja nigdy nie ma modułu ani okien rdzenia. */
function sekcja(klucz: string, nazwa: string, ikona: NazwaIkony, opis: string): PozycjaModulu {
  return { klucz, nazwa, ikona, opis, okna: [] };
}

/** Wykaz w drodze — stan pokazywany, dopóki rdzeń nie odpowiedział. */
export function srodowiskoWczytywane(klucz: KluczSrodowiska): Srodowisko {
  return {
    klucz,
    nazwa: 'Wczytywanie…',
    motto: 'Wykaz modułów w drodze z rdzenia.',
    rodzaj: 'moduly',
    stan: 'ladowanie',
    pozycje: [],
  };
}

/** Odmowa rdzenia — wykaz pusty z treścią odmowy, nie kopia zapasowa. */
export function srodowiskoNiedostepne(klucz: KluczSrodowiska, blad: string): Srodowisko {
  return {
    klucz,
    nazwa: 'Wykaz niedostępny',
    motto: blad,
    rodzaj: 'moduly',
    stan: 'blad',
    pozycje: [],
    blad,
  };
}

/**
 * Przekład odpowiedzi rdzenia na wykaz nawigacji.
 *
 * Kolejność pozycji jest kolejnością, w jakiej oddał je rdzeń — klient nie
 * sortuje po swojemu, bo macierz `srodowisko_modul` niesie własną kolumnę.
 */
export function srodowiskoZKontraktu(dane: Environment, moduly: readonly Module[]): Srodowisko {
  const orkiestracja = dane.navigationKind === NavigationKind.Orchestration;
  return {
    klucz: (dane.code as KluczSrodowiska) || SRODOWISKO_DOMYSLNE,
    nazwa: dane.name,
    motto: dane.description ?? '',
    rodzaj: orkiestracja ? 'sekcje' : 'moduly',
    stan: 'gotowe',
    pozycje: orkiestracja ? [...SEKCJE_ORKIESTRACJI] : moduly.map(pozycjaModulu),
  };
}

/**
 * Moduł kontraktu jako pozycja wykazu; ikona dokładana, reszta z rdzenia.
 *
 * Publiczna, bo pozycję buduje się także poza wykazem środowiska: moduł bez
 * wiersza widoczności w macierzy nie stoi w bocznej nawigacji, a mimo to daje
 * się otworzyć — kafel komponentu własnego na stronie głównej sięga po niego
 * wprost do katalogu `module.list` i składa pozycję tą samą funkcją. Drugiego
 * przekładu modułu na pozycję nie ma.
 */
export function pozycjaModulu(modul: Module): PozycjaModulu {
  return {
    klucz: modul.code,
    nazwa: modul.name,
    ikona: IKONY_MODULOW[modul.code] ?? IKONA_ZASTEPCZA,
    opis: modul.description ?? '',
    modul: modul.id,
    okna: modul.operationalWindowCodes,
  };
}

/**
 * Podpis pod nazwą środowiska: liczba pozycji wykazu w odmianie polskiej.
 * Trzy formy: „1 moduł", „4 moduły", „9 modułów".
 */
export function opisLiczbyPozycji(dane: Srodowisko): string {
  if (dane.stan !== 'gotowe') return dane.stan === 'blad' ? 'wykaz pusty' : 'wykaz w drodze';
  const ile = dane.pozycje.length;
  const formy: readonly [string, string, string] =
    dane.rodzaj === 'sekcje' ? ['sekcja', 'sekcje', 'sekcji'] : ['moduł', 'moduły', 'modułów'];
  return `${ile} ${odmiana(ile, formy)}`;
}

/** Odmiana rzeczownika po liczebniku: pojedyncza, mnoga, dopełniaczowa. */
function odmiana(ile: number, formy: readonly [string, string, string]): string {
  if (ile === 1) return formy[0];
  const setki = ile % 100;
  const jednosci = ile % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (setki < 12 || setki > 14);
  return mnoga ? formy[1] : formy[2];
}
