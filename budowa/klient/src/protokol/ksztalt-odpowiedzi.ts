import { ErrorCode, type ErrorInfo } from '../../../shared/contract.ts';
import type { Wynik } from './kanal.ts';

/**
 * Sprawdzian kształtu odpowiedzi w warstwie protokołu. Kanał nie waliduje
 * ładunku i rzutuje go na typ zapowiedziany przez kontrakt; sprawdzian
 * zamienia niezgodny kształt w zwykłe niepowodzenie wywołania.
 */
export function sprawdzKsztalt<T>(
  wynik: Wynik<T>,
  komenda: string,
  sprawdzian: (tresc: T) => boolean,
): Wynik<T> {
  if (!wynik.udany) return wynik;
  const tresc = wynik.wynik;
  if (tresc !== undefined && bezpiecznieSprawdz(tresc, sprawdzian)) return wynik;
  console.warn('[protokół] odpowiedź o niespodziewanym kształcie', komenda, tresc);
  return { udany: false, blad: bladKsztaltu(komenda) };
}

/** Czy wartość jest napisem — prosty strażnik typu wykorzystywany przez sprawdziany kształtu odpowiedzi. */
export function czyTekst(wartosc: unknown): wartosc is string {
  return typeof wartosc === 'string';
}

/** Wykonanie sprawdzianu odporne na jego własny błąd; wyjątek sprawdzianu zamienia się w wynik odmowny, a nie w awarię wywołania. */
function bezpiecznieSprawdz<T>(tresc: T, sprawdzian: (tresc: T) => boolean): boolean {
  try {
    return sprawdzian(tresc);
  } catch (blad) {
    console.warn('[protokół] sprawdzian kształtu przerwany', blad);
    return false;
  }
}

/** Błąd zgłaszany wywołującemu, gdy odpowiedź nie ma kształtu zapowiedzianego w kontrakcie, z nazwą komendy i kodem odmowy walidacji. */
function bladKsztaltu(komenda: string): ErrorInfo {
  return {
    code: ErrorCode.ValidationFailed,
    message: `Odpowiedź komendy ${komenda} nie ma kształtu zapowiedzianego w kontrakcie`,
    retryable: false,
  };
}
