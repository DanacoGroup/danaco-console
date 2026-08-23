/**
 * Odczyt treści nietekstowej fragmentu strumienia.
 *
 * Pole `data` kontraktu jest typu `unknown` — rdzeń pakuje w nie strukturę
 * właściwą rodzajowi fragmentu. Odczyt jest w całości tolerancyjny: pole
 * brakujące, pole innego typu ani ładunek nieznanego kształtu nie przerywają
 * strumienia. Brak wartości znaczy „nie wiem", nie „błąd".
 */

/** Ładunek fragmentu sprowadzony do zbioru pól albo `null`. */
export function obiekt(dane: unknown): Record<string, unknown> | null {
  if (typeof dane !== 'object' || dane === null || Array.isArray(dane)) return null;
  return dane as Record<string, unknown>;
}

/** Pole tekstowe ładunku; brak albo inny typ daje pusty napis. */
export function tekst(zrodlo: Record<string, unknown> | null, klucz: string): string {
  const wartosc = zrodlo?.[klucz];
  return typeof wartosc === 'string' ? wartosc : '';
}

/** Pole liczbowe ładunku; brak albo inny typ daje zero. */
export function liczba(zrodlo: Record<string, unknown> | null, klucz: string): number {
  const wartosc = zrodlo?.[klucz];
  return typeof wartosc === 'number' && Number.isFinite(wartosc) ? wartosc : 0;
}

/** Pole logiczne ładunku; brak albo inny typ daje fałsz. */
export function prawda(zrodlo: Record<string, unknown> | null, klucz: string): boolean {
  return zrodlo?.[klucz] === true;
}

/** Pole listy napisów; pozycje innego typu są pomijane. */
export function listaTekstow(zrodlo: Record<string, unknown> | null, klucz: string): string[] {
  const wartosc = zrodlo?.[klucz];
  if (!Array.isArray(wartosc)) return [];
  return wartosc.filter((pozycja): pozycja is string => typeof pozycja === 'string');
}

/**
 * Pole dowolnego kształtu w postaci czytelnej dla człowieka.
 * Napis zostaje napisem; struktura idzie wcięciem dwóch spacji.
 */
export function zapisCzytelny(wartosc: unknown): string {
  if (wartosc === undefined || wartosc === null) return '';
  if (typeof wartosc === 'string') return wartosc;
  try {
    return JSON.stringify(wartosc, null, 2) ?? '';
  } catch {
    return String(wartosc);
  }
}
