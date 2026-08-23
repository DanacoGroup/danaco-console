import { ErrorCode, type ErrorInfo } from '../../../shared/contract';
import type { Wynik } from './kanal';

/**
 * Sprawdzian kształtu odpowiedzi w warstwie protokołu.
 *
 * Kanał nie waliduje ładunku: rzutuje go na typ zapowiedziany przez kontrakt
 * i oddaje wywołującemu. Rzutowanie jest obietnicą kompilatora, nie rdzenia —
 * rdzeń starszej wersji albo pośrednik może przysłać treść bez pola
 * obowiązkowego, a widok dostałby `undefined` w miejscu, w którym typ obiecuje
 * tablicę.
 *
 * Sprawdzian zamienia taką odpowiedź w zwykłe niepowodzenie wywołania: wpis do
 * dziennika i `Wynik` z błędem `validation_failed`.
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

/** Czy wartość jest tablicą; brak pola tablicowego jest niezgodny z kontraktem. */
export function czyTablica(wartosc: unknown): wartosc is unknown[] {
  return Array.isArray(wartosc);
}

/** Czy wartość jest napisem. */
export function czyTekst(wartosc: unknown): wartosc is string {
  return typeof wartosc === 'string';
}

/** Czy wartość jest liczbą skończoną. */
export function czyLiczba(wartosc: unknown): wartosc is number {
  return typeof wartosc === 'number' && Number.isFinite(wartosc);
}

/** Czy wartość jest wartością logiczną. */
export function czyLogiczna(wartosc: unknown): wartosc is boolean {
  return typeof wartosc === 'boolean';
}

/** Czy wartość jest obiektem — tablica i `null` obiektem nie są. */
export function czyObiekt(wartosc: unknown): wartosc is Record<string, unknown> {
  return typeof wartosc === 'object' && wartosc !== null && !Array.isArray(wartosc);
}

/** Wykonanie sprawdzianu odporne na jego własny błąd. */
function bezpiecznieSprawdz<T>(tresc: T, sprawdzian: (tresc: T) => boolean): boolean {
  try {
    return sprawdzian(tresc);
  } catch (blad) {
    console.warn('[protokół] sprawdzian kształtu przerwany', blad);
    return false;
  }
}

/** Błąd zgłaszany wywołującemu, gdy odpowiedź nie ma kształtu z kontraktu. */
function bladKsztaltu(komenda: string): ErrorInfo {
  return {
    code: ErrorCode.ValidationFailed,
    message: `Odpowiedź komendy ${komenda} nie ma kształtu zapowiedzianego w kontrakcie`,
    retryable: false,
  };
}
