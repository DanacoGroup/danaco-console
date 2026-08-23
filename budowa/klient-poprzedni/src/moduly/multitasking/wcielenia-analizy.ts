/**
 * Wcielenia roli Results Analyzer.
 *
 * Wcielenie zmienia kryteria oceny, a nie samą etykietę: każde niesie własny
 * wykaz kryteriów, po których analityk czyta wyniki wykonawców, i własny
 * nagłówek zgłoszenia niezgodności kierowanego do koordynatora.
 *
 * Wartość utrwala `config.set` na poziomie okna analityka. Rdzeń nie
 * rejestruje `role.update`, więc profilu wcielenia nie zna — pokrycie mierzy
 * `braki-kontraktu.ts`.
 */

export const Wcielenie = {
  Validator: 'validator',
  Reviewer: 'reviewer',
  SecurityAuditor: 'security-auditor',
  Architect: 'architect',
  ProductOwner: 'product-owner',
  QaLead: 'qa-lead',
  Arbitrator: 'arbitrator',
} as const;
export type Wcielenie = (typeof Wcielenie)[keyof typeof Wcielenie];

/** Klucz utrwalenia wcielenia na poziomie okna analityka. */
export const KLUCZ_WCIELENIA = 'multitasking.wcielenie';

/** Opis wcielenia: nazwa w selektorze i kryteria czytania wyników. */
export interface OpisWcielenia {
  wcielenie: Wcielenie;
  nazwa: string;
  /** Po czym to wcielenie ocenia wynik wykonawcy. */
  kryteria: readonly string[];
}

export const WCIELENIA: readonly OpisWcielenia[] = [
  {
    wcielenie: Wcielenie.Validator,
    nazwa: 'Validator',
    kryteria: ['zgodność z treścią zlecenia', 'kompletność wyniku', 'brak kroków pominiętych'],
  },
  {
    wcielenie: Wcielenie.Reviewer,
    nazwa: 'Reviewer',
    kryteria: ['czytelność rozwiązania', 'uzasadnienie decyzji', 'ślad wywołań narzędzi'],
  },
  {
    wcielenie: Wcielenie.SecurityAuditor,
    nazwa: 'Security Auditor',
    kryteria: ['dane wrażliwe w wyniku', 'zakres uprawnień użytych narzędzi', 'ślad audytowy tury'],
  },
  {
    wcielenie: Wcielenie.Architect,
    nazwa: 'Architect',
    kryteria: ['spójność z przyjętą budową', 'brak drugiego bytu tej samej rzeczy', 'granice odpowiedzialności'],
  },
  {
    wcielenie: Wcielenie.ProductOwner,
    nazwa: 'Product Owner',
    kryteria: ['pokrycie wymagania', 'wartość dla Operatora', 'to, czego w wyniku brakuje'],
  },
  {
    wcielenie: Wcielenie.QaLead,
    nazwa: 'QA Lead',
    kryteria: ['sprawdzalność wyniku', 'przypadki brzegowe', 'stan po błędzie'],
  },
  {
    wcielenie: Wcielenie.Arbitrator,
    nazwa: 'Arbitrator',
    kryteria: ['różnica między wykonawcami', 'siła uzasadnienia każdej wersji', 'wskazanie wersji przyjętej'],
  },
];

/** Opis wcielenia po wartości; nieznana schodzi na Validatora. */
export function opisWcielenia(wartosc: unknown): OpisWcielenia {
  return WCIELENIA.find((opis) => opis.wcielenie === wartosc) ?? WCIELENIA[0]!;
}

/** Pozycje selektora wcielenia. */
export function pozycjeWcielen(): ReadonlyArray<[string, string]> {
  return WCIELENIA.map((opis) => [opis.wcielenie, opis.nazwa]);
}

/**
 * Zgłoszenie niezgodności kierowane do okna koordynatora.
 *
 * Kontrakt nie ma komendy utrwalającej werdykt oceny jako stan: `monitor.status`
 * czyta stan procesów i niczego nie zapisuje. Ocena analityka jest poleceniem
 * dla koordynatora, więc zgłoszenie idzie zwykłym `message.send`.
 */
export function trescZgloszenia(
  opis: OpisWcielenia,
  werdykt: string,
  uzasadnienie: string,
): string {
  const kryteria = opis.kryteria.map((pozycja) => `- ${pozycja}`).join('\n');
  const tresc = uzasadnienie.trim() === '' ? 'Analityk nie podał uzasadnienia.' : uzasadnienie.trim();
  return `Results Analyzer (${opis.nazwa}) — ${werdykt}\n\nKryteria oceny:\n${kryteria}\n\nUzasadnienie:\n${tresc}`;
}
