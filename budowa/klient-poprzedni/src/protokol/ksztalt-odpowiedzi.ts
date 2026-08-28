import { ErrorCode, type ErrorInfo } from '../../../shared/contract';
import type { Wynik } from './kanal';

/** Sprawdzian kształtu odpowiedzi w warstwie protokołu, zamieniający odpowiedź niezgodną w zwykłe niepowodzenie. */
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

/** Czy wartość jest tablicą; brak pola tablicowego w odpowiedzi jest niezgodny z zapowiedzią kontraktu. */
export function czyTablica(wartosc: unknown): wartosc is unknown[] {
  return Array.isArray(wartosc);
}

/** Czy wartość jest napisem — dowolnym ciągiem znaków, obojętnie jak długim, bez dalszych warunków treści. */
export function czyTekst(wartosc: unknown): wartosc is string {
  return typeof wartosc === 'string';
}

/** Czy wartość jest liczbą skończoną — nie jest wartością nieskończoną ani wynikiem błędnego działania. */
export function czyLiczba(wartosc: unknown): wartosc is number {
  return typeof wartosc === 'number' && Number.isFinite(wartosc);
}

/** Czy wartość jest wartością logiczną — jednym z dwóch stanów, prawda albo fałsz, bez wartości pustej. */
export function czyLogiczna(wartosc: unknown): wartosc is boolean {
  return typeof wartosc === 'boolean';
}

/** Czy wartość jest obiektem w rozumieniu tego sprawdzianu — tablica oraz wartość pusta obiektem nie są. */
export function czyObiekt(wartosc: unknown): wartosc is Record<string, unknown> {
  return typeof wartosc === 'object' && wartosc !== null && !Array.isArray(wartosc);
}

/** Wykonanie sprawdzianu kształtu odporne na jego własny błąd, zgłaszane dziennikowi zamiast wywołującemu. */
function bezpiecznieSprawdz<T>(tresc: T, sprawdzian: (tresc: T) => boolean): boolean {
  try {
    return sprawdzian(tresc);
  } catch (blad) {
    console.warn('[protokół] sprawdzian kształtu przerwany', blad);
    return false;
  }
}

/** Błąd zgłaszany wywołującemu, gdy odpowiedź rdzenia nie ma kształtu zapowiedzianego w tym kontrakcie. */
function bladKsztaltu(komenda: string): ErrorInfo {
  return {
    code: ErrorCode.ValidationFailed,
    message: `Odpowiedź komendy ${komenda} nie ma kształtu zapowiedzianego w kontrakcie`,
    retryable: false,
  };
}
