/**
 * Odczyt treści nietekstowej fragmentu strumienia, tolerancyjny wobec pola brakującego,
 * innego typu albo nieznanego kształtu.
 */

/**
 * Ładunek fragmentu strumienia sprowadzony do zbioru odczytanych pól o znanym kształcie
 * albo do wartości `null`.
 */
export function obiekt(dane: unknown): Record<string, unknown> | null {
  if (typeof dane !== 'object' || dane === null || Array.isArray(dane)) return null;
  return dane as Record<string, unknown>;
}

/**
 * Pole tekstowe ładunku fragmentu strumienia; brak pola albo wartość innego typu daje w
 * wyniku pusty napis.
 */
export function tekst(zrodlo: Record<string, unknown> | null, klucz: string): string {
  const wartosc = zrodlo?.[klucz];
  return typeof wartosc === 'string' ? wartosc : '';
}

/**
 * Pole liczbowe ładunku fragmentu strumienia; brak pola albo wartość innego typu daje w
 * wyniku zero.
 */
export function liczba(zrodlo: Record<string, unknown> | null, klucz: string): number {
  const wartosc = zrodlo?.[klucz];
  return typeof wartosc === 'number' && Number.isFinite(wartosc) ? wartosc : 0;
}

/**
 * Pole logiczne ładunku fragmentu strumienia; brak pola albo wartość innego typu daje w
 * wyniku fałsz.
 */
export function prawda(zrodlo: Record<string, unknown> | null, klucz: string): boolean {
  return zrodlo?.[klucz] === true;
}

/**
 * Pole listy napisów w ładunku fragmentu strumienia — pozycje o innym typie są z tej listy
 * pomijane.
 */
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
