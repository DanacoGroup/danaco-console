import { WindowRole } from '../../../shared/contract';
import { kierunekKoordynatora, kierunekWykonawcy } from './etykiety-ukladu';
import { numerGniazda, type IdGniazda } from './identyfikatory';

/**
 * Para koordynator–wykonawca odczytana z ról okien widocznych na scenie, wraz
 * z informacją, po której stronie sceny stoi koordynator względem wykonawcy.
 */
export interface Wiez {
  koordynator: IdGniazda;
  wykonawca: IdGniazda;
  /** Prawda, gdy koordynator stoi na scenie po lewej stronie wykonawcy. */
  wPrawo: boolean;
}

/**
 * Opis kierunku zlecenia pokazywany w nagłówku jednego gniazda: tekst opisu
 * oraz kierunek i położenie grota względem napisu.
 */
export interface OpisKierunku {
  tekst: string;
  /** Prawda, gdy grot ma wskazywać w prawo — kierunek biegu zlecenia. */
  wPrawo: boolean;
  /** Prawda, gdy grot stoi za napisem — po stronie okna partnera. */
  grotPoPrawej: boolean;
}

/**
 * Ustala parę koordynator–wykonawca, wiążąc pierwszy widoczny koordynator z pierwszym
 * widocznym wykonawcą; brak którejkolwiek roli oznacza brak więzi.
 */
export function ustalWiez(
  role: ReadonlyMap<IdGniazda, WindowRole>,
  widoczne: readonly IdGniazda[],
): Wiez | null {
  const koordynator = widoczne.find((id) => role.get(id) === WindowRole.Coordinator);
  const wykonawca = widoczne.find((id) => role.get(id) === WindowRole.Executor);
  if (koordynator === undefined || wykonawca === undefined) return null;
  return {
    koordynator,
    wykonawca,
    wPrawo: numerGniazda(koordynator) < numerGniazda(wykonawca),
  };
}

/**
 * Ustala kierunek zlecenia widoczny w nagłówku danego gniazda, z napisem i położeniem
 * grota zależnym od roli gniazda względem więzi koordynator–wykonawca.
 */
export function opisKierunku(wiez: Wiez | null, id: IdGniazda): OpisKierunku | null {
  if (wiez === null) return null;
  if (id === wiez.koordynator) {
    return {
      tekst: kierunekKoordynatora(wiez.wykonawca),
      wPrawo: wiez.wPrawo,
      grotPoPrawej: wiez.wPrawo,
    };
  }
  if (id === wiez.wykonawca) {
    return {
      tekst: kierunekWykonawcy(wiez.koordynator),
      wPrawo: wiez.wPrawo,
      grotPoPrawej: !wiez.wPrawo,
    };
  }
  return null;
}

/**
 * Krańce więzi ułożone w kolejności scenicznej, od lewej do prawej strony sceny, na
 * podstawie ustalonego kierunku pary koordynator–wykonawca.
 */
export function krance(wiez: Wiez): { lewy: IdGniazda; prawy: IdGniazda } {
  return wiez.wPrawo
    ? { lewy: wiez.koordynator, prawy: wiez.wykonawca }
    : { lewy: wiez.wykonawca, prawy: wiez.koordynator };
}
