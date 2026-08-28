import {
  KnownModuleIds,
  NavigationKind,
  type Environment,
  type Module,
} from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

// Widok wykazu środowiska — kształt, w jakim boczna nawigacja czyta odpowiedź rdzenia.

/** Klucz środowiska — kody znane kontraktowi rdzenia, nie literały wymyślone samodzielnie przez klienta. */
export type KluczSrodowiska = (typeof KnownModuleIds)[number];

/** Rodzaj wykazu w bocznej nawigacji — moduły tego środowiska albo sekcje panelu orkiestracji zespołowej. */
export type RodzajWykazu = 'moduly' | 'sekcje';

/** Stan wykazu bocznej nawigacji: przed odpowiedzią rdzenia, zaraz po niej albo po jego pełnej odmowie. */
export type StanWykazu = 'ladowanie' | 'gotowe' | 'blad';

/** Pojedyncza pozycja bocznej nawigacji: moduł tego środowiska albo sekcja panelu orkiestracji zespołowej. */
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

/** Środowisko wraz z pełnym wykazem pozycji bocznej nawigacji przygotowanym dla tego środowiska klienta. */
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

/** Środowisko otwierane, gdy wywołanie tworzące powłokę środowiska nie wskazuje żadnego innego środowiska. */
export const SRODOWISKO_DOMYSLNE: KluczSrodowiska = KnownModuleIds[0];

/** Ikona pozycji, gdy kod modułu jest rdzeniowi znany, a temu klientowi jeszcze zupełnie nie jest znany. */
const IKONA_ZASTEPCZA: NazwaIkony = 'karta-okna';

/** Kod modułu przekładany na ikonę zestawu wizualnego — jedyny powód istnienia tej pomocniczej mapy kodów. */
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

/** Klucze sekcji panelu orkiestracji — zestaw domyślny, stojący osobno i publicznie dla wielu odbiorców. */
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

/** Skrót zapisu sekcji panelu orkiestracji — sekcja nigdy nie ma własnego modułu ani żadnych okien rdzenia. */
function sekcja(klucz: string, nazwa: string, ikona: NazwaIkony, opis: string): PozycjaModulu {
  return { klucz, nazwa, ikona, opis, okna: [] };
}

/** Wykaz w drodze — stan pokazywany bocznej nawigacji, dopóki rdzeń jeszcze wcale nie odpowiedział na żądanie. */
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

/** Odmowa rdzenia — wykaz pusty z treścią odmowy tego samego rdzenia, nigdy kopia zapasowa tego wykazu. */
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

/** Moduł kontraktu jako pozycja wykazu bocznej nawigacji; ikona dokładana lokalnie, reszta wprost z rdzenia. */
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

/** Odmiana rzeczownika po liczebniku wskazań: forma pojedyncza, forma mnoga oraz forma dopełniaczowa liczby. */
function odmiana(ile: number, formy: readonly [string, string, string]): string {
  if (ile === 1) return formy[0];
  const setki = ile % 100;
  const jednosci = ile % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (setki < 12 || setki > 14);
  return mnoga ? formy[1] : formy[2];
}
