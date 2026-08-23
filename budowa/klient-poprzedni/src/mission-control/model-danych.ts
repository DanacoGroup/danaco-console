import type {
  ProgressStatus,
  QueueStatus,
  SessionStatus,
  WindowRole,
} from '../../../shared/contract';

/**
 * Model danych pulpitu operacyjnego (Mission Control).
 *
 * Jedna odpowiedzialność: kształt danych, które widok pulpitu umie wyrysować.
 * Wartości buduje `zlozenie-danych.ts` wyłącznie z odczytów i zdarzeń rdzenia.
 *
 * Pole `null` znaczy brak źródła: miara, której kontrakt nie niesie, ma w modelu
 * typ `X | null`, a widok wypisuje przy niej etykietę „brak źródła danych
 * w kontrakcie" zamiast liczby.
 *
 * Nazwy stanów pochodzą z kontraktu (`shared/contract`) zamiast z powtarzanych
 * literałów.
 */

/** Skąd pochodzą liczby pokazane na pulpicie. */
export const ZrodloDanych = {
  /** Odczyt z rdzenia jeszcze nie nadszedł — pulpit czeka na dane. */
  Oczekiwanie: 'oczekiwanie',
  /** Dane odczytane z rdzenia przez kanał kontraktu. */
  Rdzen: 'rdzen',
} as const;
export type ZrodloDanych = (typeof ZrodloDanych)[keyof typeof ZrodloDanych];

/**
 * Kod środowiska Danaco Console — wartość kolumny `srodowisko.kod` z rdzenia.
 *
 * Typ jest napisem, a nie unią wywiedzioną z `KnownModuleIds`, bo wykaz środowisk
 * należy do rdzenia jako dane i `environment.list` niesie go w całości. Unia
 * zamykałaby matrycę na kody znane klientowi, a środowisko spoza niej trafiałoby
 * do wykazu „poza środowiskami" jako sesja bez wskazania środowiska.
 *
 * `KnownModuleIds` pilnuje kodów tam, gdzie klient sam je wymienia
 * (`strona-glowna/pozycje-srodowisk.ts`); kolumna matrycy przepisuje to,
 * co przyszło z rdzenia.
 */
export type IdSrodowiska = string;

/** Rodzaj kafla w rzędzie „Utwórz". */
export const RodzajUtworzenia = {
  Sesja: 'sesja',
  Projekt: 'projekt',
  Agent: 'agent',
  Automatyka: 'automatyka',
  Kolejka: 'kolejka',
  Zespol: 'zespol',
} as const;
export type RodzajUtworzenia = (typeof RodzajUtworzenia)[keyof typeof RodzajUtworzenia];

/** Cztery liczby sekcji „Aktywność AI"; dwie ostatnie bez źródła — `null`. */
export interface AktywnoscAI {
  /** Procesy w stanie `running` z telemetrii `progress.changed`. */
  aktywneProcesy: number;
  /** Okna komunikacji, w których biegnie proces telemetrii. */
  agenciPracuja: number;
  /** Brak źródła: kontrakt nie niesie liczby kanałów przetwarzających zlecenie. */
  modeleAnalizuja: number | null;
  /** Brak źródła: kontrakt nie niesie liczby wyników czekających na walidatora. */
  walidatoryOczekuja: number | null;
}

/** Pas eskalacji koordynatora — przepływy czekające na człowieka. */
export interface PasDecyzji {
  /** Ile przepływów wstrzymano do decyzji użytkownika. */
  przeplywyDoDecyzji: number;
  /** Nazwa przepływu, który czeka najdłużej. */
  najstarszyPrzeplyw: string;
  /** Ile czasu czeka, w minutach. */
  czekaMinut: number;
}

/**
 * Plakietka jednego kanału modelu w sekcji „Operacje AI" — odwzorowanie wiersza
 * `channel.list`. Miary wysycenia, kolejki, limitu i kosztu nie mają dziś
 * źródła w kontrakcie, więc plakietka ich nie niesie — widok wypisuje przy nich
 * etykietę braku źródła.
 */
export interface KanalOperacyjny {
  id: string;
  /** Nazwa kanału widoczna dla operatora — `Channel.name`. */
  nazwa: string;
  /** Rodzaj kanału — `Channel.kind`; wartość danych, nie typ kodu. */
  rodzaj: string;
  /** Identyfikator modelu — `Channel.model`; `null`, gdy wiersz go nie podaje. */
  model: string | null;
  /** Czy kanał czynny — `Channel.enabled`. */
  czynny: boolean;
}

/** Jedna sesja w kolumnie środowiska. */
export interface SesjaMatrycy {
  /** Identyfikator sesji kontraktu — trafia do zdarzenia wejścia. */
  id: string;
  /** Tytuł sesji; sesja bez tytułu pokazuje identyfikator. */
  tytul: string;
  status: SessionStatus;
  /** Rola okna wiodącego; `null`, gdy odczyt okien sesji nie zna żadnego. */
  rolaOkna: WindowRole | null;
  /** Liczba okien komunikacji sesji; sesja wiąże wiele okien. */
  okna: number;
  /** Czy w sesji biegnie praca, gdy operator patrzy gdzie indziej. */
  pracaWTle: boolean;
}

/** Kolumna matrycy — jedno środowisko wraz z sesjami. */
export interface KolumnaSrodowiska {
  id: IdSrodowiska;
  nazwa: string;
  motto: string;
  sesje: SesjaMatrycy[];
}

/** Rodzaj powiązania między dwiema sesjami w pasie RELACJE. */
export const RodzajRelacji = {
  /** Wymiana dwustronna: obie sesje zasilają się nawzajem. */
  Wymiana: 'wymiana',
  /** Przekazanie jednostronne: wynik pierwszej zasila drugą. */
  Przekazanie: 'przekazanie',
} as const;
export type RodzajRelacji = (typeof RodzajRelacji)[keyof typeof RodzajRelacji];

/** Powiązanie dwóch sesji, także z różnych środowisk. */
export interface RelacjaSesji {
  id: string;
  zrodloId: string;
  zrodloNazwa: string;
  celId: string;
  celNazwa: string;
  rodzaj: RodzajRelacji;
  /** Czym rzeczy się wymieniają: kontekst, artefakt, werdykt. */
  przedmiot: string;
}

/** Proces biegnący w tle — kolumna „Procesy w tle" (zdarzenie `progress.changed`). */
export interface ProcesWTle {
  /** Identyfikator procesu z telemetrii — `processId`. */
  id: string;
  /** Nazwa okna procesu; bez znanego okna — identyfikator procesu. */
  nazwa: string;
  /** Sesja okna procesu; `null`, gdy telemetria nie wskazała okna. */
  sesjaId: string | null;
  etapBiezacy: number;
  /** Liczba etapów; 0 znaczy „nieznana" (kontrakt, `totalSteps`). */
  etapowRazem: number;
  /** Nazwa etapu bieżącego; `null`, gdy telemetria jej nie podała. */
  etykietaEtapu: string | null;
  status: ProgressStatus;
}

/**
 * Kolejka na pulpicie — odwzorowanie bytu `Queue` ze zdarzeń `queue.changed`.
 *
 * Kontrakt nie niesie roli kolejki ani liczby zadań czekających; nie ma też
 * odczytu `queue.list`, więc przed pierwszym zdarzeniem wykaz jest pusty.
 */
export interface KolejkaPulpitu {
  /** Identyfikator kolejki kontraktu — trafia do `queue.action`. */
  id: string;
  /** Nazwa kolejki; kolejka bez nazwy pokazuje identyfikator. */
  nazwa: string;
  status: QueueStatus;
  /** Liczba okien obsługiwanych; `null`, gdy zdarzenie ich nie wymieniło. */
  okna: number | null;
  /** Licznik obiegów naprawczych; bez limitu; `null` = nie podano. */
  obiegi: number | null;
}

/** Agent zespołu — kolumna „Zespół agentów"; odwzorowanie otwartego okna. */
export interface AgentZespolu {
  /** Identyfikator okna komunikacji. */
  id: string;
  nazwa: string;
  /** Etap bieżący z telemetrii procesu okna; `null`, gdy telemetrii brak. */
  zajecie: string | null;
  rola: WindowRole;
  /** Brak źródła: kontrakt nie niesie liczby podagentów okna. */
  podagenci: number | null;
  /** Czy w oknie biegnie proces telemetrii. */
  czynny: boolean;
}

/** Komplet danych jednego wyrysowania pulpitu. */
export interface DanePulpitu {
  /** Rozstrzyga, czy odczyt z rdzenia już nadszedł. */
  zrodloDanych: ZrodloDanych;
  aktywnosc: AktywnoscAI;
  /** `null` = brak źródła: kontrakt nie niesie przepływów do decyzji. */
  decyzje: PasDecyzji | null;
  kanaly: KanalOperacyjny[];
  matryca: KolumnaSrodowiska[];
  /** Sesje, których środowiska odczyt nie wskazał (`presence.environmentCode`). */
  pozaSrodowiskami: SesjaMatrycy[];
  /** `null` = brak źródła: kontrakt nie niesie powiązań między sesjami. */
  relacje: RelacjaSesji[] | null;
  procesy: ProcesWTle[];
  kolejki: KolejkaPulpitu[];
  zespol: AgentZespolu[];
}
