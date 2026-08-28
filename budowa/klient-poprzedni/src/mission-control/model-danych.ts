import type {
  ProgressStatus,
  QueueStatus,
  SessionStatus,
  WindowRole,
} from '../../../shared/contract';

// Model danych pulpitu — kształt danych, budowany wyłącznie z odczytów rdzenia.

/** Skąd pochodzą liczby pokazane na pulpicie: odczyt jeszcze nie nadszedł albo dane odczytane z rdzenia. */
export const ZrodloDanych = {
  /** Odczyt z rdzenia jeszcze nie nadszedł — pulpit czeka na dane. */
  Oczekiwanie: 'oczekiwanie',
  /** Dane odczytane z rdzenia przez kanał kontraktu. */
  Rdzen: 'rdzen',
} as const;
export type ZrodloDanych = (typeof ZrodloDanych)[keyof typeof ZrodloDanych];

/**
 * Kod środowiska Danaco Console — wartość kolumny `srodowisko.kod` z rdzenia.
 * Typ jest napisem, nie unią, bo wykaz środowisk należy do rdzenia jako dane,
 * a środowisko spoza znanej unii trafiałoby do wykazu „poza środowiskami".
 */
export type IdSrodowiska = string;

/** Rodzaj kafla w rzędzie „Utwórz" pulpitu operacyjnego — jedna z dostępnych pozycji szybkiego tworzenia. */
export const RodzajUtworzenia = {
  Sesja: 'sesja',
  Projekt: 'projekt',
  Agent: 'agent',
  Automatyka: 'automatyka',
  Kolejka: 'kolejka',
  Zespol: 'zespol',
} as const;
export type RodzajUtworzenia = (typeof RodzajUtworzenia)[keyof typeof RodzajUtworzenia];

/** Cztery liczby sekcji „Aktywność AI" pulpitu; dwie ostatnie nie mają dziś źródła w kontrakcie — `null`. */
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

/** Pas eskalacji koordynatora na pulpicie operacyjnym — przepływy sesji czekające na decyzję użytkownika. */
export interface PasDecyzji {
  /** Ile przepływów wstrzymano do decyzji użytkownika. */
  przeplywyDoDecyzji: number;
  /** Nazwa przepływu, który czeka najdłużej. */
  najstarszyPrzeplyw: string;
  /** Ile czasu czeka, w minutach. */
  czekaMinut: number;
}

/**
 * Plakietka jednego kanału modelu w sekcji „Operacje AI" — odwzorowanie
 * wiersza `channel.list`. Miary wysycenia, kolejki, limitu i kosztu nie mają
 * dziś źródła w kontrakcie, więc plakietka ich nie niesie.
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

/** Jedna sesja odwzorowana w kolumnie środowiska macierzy pulpitu operacyjnego aplikacji Mission Control. */
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

/** Kolumna matrycy pulpitu operacyjnego — jedno środowisko Danaco Console wraz ze wszystkimi jego sesjami. */
export interface KolumnaSrodowiska {
  id: IdSrodowiska;
  nazwa: string;
  motto: string;
  sesje: SesjaMatrycy[];
}

/** Rodzaj powiązania między dwiema sesjami w pasie „Relacje" pulpitu operacyjnego aplikacji Mission Control. */
export const RodzajRelacji = {
  /** Wymiana dwustronna: obie sesje zasilają się nawzajem. */
  Wymiana: 'wymiana',
  /** Przekazanie jednostronne: wynik pierwszej zasila drugą. */
  Przekazanie: 'przekazanie',
} as const;
export type RodzajRelacji = (typeof RodzajRelacji)[keyof typeof RodzajRelacji];

/** Powiązanie dwóch sesji pulpitu operacyjnego, także pochodzących z dwóch różnych środowisk Danaco Console. */
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

/** Proces biegnący w tle — kolumna „Procesy w tle" pulpitu operacyjnego, zasilana zdarzeniem `progress.changed`. */
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

/** Agent zespołu w kolumnie „Zespół agentów" pulpitu operacyjnego — odwzorowanie jednego otwartego okna. */
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

/** Komplet danych potrzebny do jednego wyrysowania pulpitu operacyjnego Mission Control całej aplikacji. */
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
