/**
 * Hierarchia decyzji MultitaskingAI i dwa tryby Always On Display.
 *
 * Poziomy decyzji są opisem po stronie klienta, nie stanem rdzenia: kontrakt
 * niesie role okien (`WindowRole`: coordinator, executor, standalone)
 * i podagentów (`subagent.*`), ale nie niesie wykazu poziomów ani ich
 * kolejności. Wykaz tłumaczy w sekcji Monitor, kto może przerwać kogo, i nie
 * jest bramą — hierarchię można konfigurować i pominąć, a interwencja jest
 * możliwa na dowolnym poziomie w dowolnej chwili.
 *
 * Tryby nakładki są dwa: obserwator patrzy i doradza, operator dodatkowo
 * zatwierdza i wstrzymuje kroki. Przełącznik stoi w sekcji Monitor procesu.
 *
 * Kontrakt nie zna trybu nakładki: `AodStatus` niesie urządzenie, sesję, okno,
 * procesy przypięte i licznik procesów w biegu, a komendy `aod.mode.set` nie ma
 * wcale. Tryb utrwala się więc ustawieniem na poziomie sesji (`config.set`),
 * a sekcja wypisuje tę lukę kontraktu w meldunku braków.
 */

/** Klucz utrwalenia trybu nakładki na poziomie karty sesji. */
export const KLUCZ_TRYBU_AOD = 'multitasking.tryb_aod';

/** Kształt żądania, którego kontrakt nie ma — do wypisania w meldunku braków. */
export const KSZTALT_TRYBU_AOD =
  'aod.mode.set { deviceId?: string, sessionId?: string, mode: "observer" | "operator" } → { status: AodStatus }';

export const TrybNakladki = {
  /** Podgląd pętli, statusów i hierarchii; sugestie i ostrzeżenia, bez ingerencji. */
  Obserwator: 'obserwator',
  /** Jak obserwator plus zatwierdzanie i wstrzymywanie kroków z dowolnego miejsca. */
  Operator: 'operator',
} as const;
export type TrybNakladki = (typeof TrybNakladki)[keyof typeof TrybNakladki];

/** Nazwy trybów w przełączniku sekcji Monitor. */
export const NAZWY_TRYBOW_AOD: ReadonlyArray<[TrybNakladki, string]> = [
  [TrybNakladki.Obserwator, 'Obserwator — podgląd, sugestie, bez ingerencji'],
  [TrybNakladki.Operator, 'Operator — zatwierdzanie i wstrzymywanie kroków'],
];

/**
 * Czy wartość jest jednym z dwóch trybów.
 *
 * Wartość nieznana schodzi na obserwatora, bo to tryb węższy: nakładka, o której
 * trybie nic nie wiadomo, nie ma uchodzić za uprawnioną do wstrzymywania kroków.
 */
export function trybNakladkiZWartosci(wartosc: unknown): TrybNakladki {
  return NAZWY_TRYBOW_AOD.some(([tryb]) => tryb === wartosc)
    ? (wartosc as TrybNakladki)
    : TrybNakladki.Obserwator;
}

/** Jeden poziom hierarchii decyzji: numer, nazwa i zakres rozstrzygnięć. */
export interface PoziomDecyzji {
  numer: number;
  nazwa: string;
  zakres: string;
}

/** Poziomy w kolejności od nadrzędnego do najwęższego. */
export const POZIOMY_DECYZJI: readonly PoziomDecyzji[] = [
  { numer: 1, nazwa: 'Użytkownik (Operator)', zakres: 'decyzja nadrzędna, bez ograniczeń, w dowolnym momencie' },
  { numer: 2, nazwa: 'Always On Display (operator)', zakres: 'interwencja z dowolnego miejsca; kanałem wykonania jest Mobile' },
  { numer: 3, nazwa: 'Coordinator', zakres: 'planowanie, podział pracy, sterowanie, kolejka; nie tworzy produktu końcowego' },
  { numer: 4, nazwa: 'Executor 3 / Validator', zakres: 'ocena jakości, retry, rozstrzyganie konfliktów (Arbitrator)' },
  { numer: 5, nazwa: 'Executor 1 / Executor 2', zakres: 'realizacja; brak uprawnień do zmiany planu procesu' },
  { numer: 6, nazwa: 'Subagent Network', zakres: 'podzadania; wynik agregowany i oceniany wyżej' },
];

/**
 * Wcielenia walidatora — konfiguracja jednej roli, nie osobne byty.
 *
 * Stoją napisami, a nie wierszami rejestru, bo rdzeń zna jedną rolę `executor`
 * i jedno `role.update`.
 */
export const WCIELENIA_WALIDATORA: readonly string[] = [
  'Validator',
  'Reviewer',
  'Security Auditor',
  'Architect',
  'Product Owner',
  'QA Lead',
  'Arbitrator',
];
