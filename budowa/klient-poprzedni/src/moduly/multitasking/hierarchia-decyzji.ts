/**
 * Hierarchia decyzji MultitaskingAI i dwa tryby nakładki Always On Display.
 * Poziomy decyzji są opisem po stronie klienta, a nie stanem rdzenia: kontrakt
 * niesie role okien i podagentów, lecz nie niesie wykazu poziomów ani kolejności.
 */

/**
 * Klucz utrwalenia trybu nakładki na poziomie karty sesji. Tryb zapisuje się
 * ustawieniem `config.set`, bo kontrakt nie ma komendy przestawiającej tryb,
 * a nakładka ma go pamiętać między otwarciami karty.
 */
export const KLUCZ_TRYBU_AOD = 'multitasking.tryb_aod';

/**
 * Kształt żądania, którego kontrakt nie ma — do wypisania w meldunku braków.
 * Napis podaje nazwę komendy wraz z polami żądania i odpowiedzi, więc meldunek
 * mówi wprost, czego brakuje, zamiast nazywać brak ogólnie.
 */
export const KSZTALT_TRYBU_AOD =
  'aod.mode.set { deviceId?: string, sessionId?: string, mode: "observer" | "operator" } → { status: AodStatus }';

export const TrybNakladki = {
  /** Podgląd pętli, statusów i hierarchii; sugestie i ostrzeżenia, bez ingerencji. */
  Obserwator: 'obserwator',
  /** Jak obserwator plus zatwierdzanie i wstrzymywanie kroków z dowolnego miejsca. */
  Operator: 'operator',
} as const;
export type TrybNakladki = (typeof TrybNakladki)[keyof typeof TrybNakladki];

/**
 * Nazwy trybów w przełączniku sekcji Monitor. Każda pozycja niesie kod trybu wraz
 * ze zdaniem o jego zakresie, żeby wybór mówił, co się zmieni, jeszcze przed
 * przestawieniem przełącznika.
 */
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

/**
 * Jeden poziom hierarchii decyzji: numer, nazwa i zakres rozstrzygnięć. Zakres
 * jest zdaniem, a nie zbiorem uprawnień, ponieważ hierarchia tłumaczy podział
 * pracy, a nie steruje dostępem do czynności.
 */
export interface PoziomDecyzji {
  numer: number;
  nazwa: string;
  zakres: string;
}

/**
 * Poziomy w kolejności od nadrzędnego do najwęższego. Wykaz tłumaczy, kto może
 * przerwać kogo, i nie jest bramą — interwencja pozostaje możliwa na dowolnym
 * poziomie i w dowolnej chwili.
 */
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
