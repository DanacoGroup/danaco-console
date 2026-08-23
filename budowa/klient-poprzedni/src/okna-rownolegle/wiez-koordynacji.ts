import { WindowRole } from '../../../shared/contract';
import { kierunekKoordynatora, kierunekWykonawcy } from './etykiety-ukladu';
import { numerGniazda, type IdGniazda } from './identyfikatory';

/** Para koordynator–wykonawca odczytana z ról okien na scenie. */
export interface Wiez {
  koordynator: IdGniazda;
  wykonawca: IdGniazda;
  /** Prawda, gdy koordynator stoi na scenie po lewej stronie wykonawcy. */
  wPrawo: boolean;
}

/** Opis kierunku pokazywany w nagłówku jednego gniazda. */
export interface OpisKierunku {
  tekst: string;
  /** Prawda, gdy grot ma wskazywać w prawo — kierunek biegu zlecenia. */
  wPrawo: boolean;
  /** Prawda, gdy grot stoi za napisem — po stronie okna partnera. */
  grotPoPrawej: boolean;
}

/**
 * Ustalenie pary koordynator–wykonawca.
 *
 * Wiąże się pierwszy widoczny koordynator z pierwszym widocznym wykonawcą.
 * Trzecie okno w roli wykonawcy pozostaje poza tą więzią — pas relacji
 * pokazuje jedno powiązanie naraz, żeby obraz dał się ogarnąć wzrokiem.
 * Brak którejkolwiek z ról oznacza brak więzi, nie błąd.
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
 * Kierunek zlecenia widoczny w nagłówku danego gniazda.
 *
 * Koordynator niesie napis „Zleca Oknu N", wykonawca „Zlecenia z Okna N".
 * Groty obu napisów wskazują tę samą stronę sceny — kierunek biegu zlecenia —
 * i stoją po stronie okna partnera. U koordynatora grot wychodzi z napisu ku
 * wykonawcy, u wykonawcy wchodzi w napis od strony koordynatora, więc kierunek
 * pętli czyta się bez czytania słów. Gniazdo spoza więzi nie dostaje napisu.
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

/** Krańce więzi ułożone w kolejności scenicznej, od lewej. */
export function krance(wiez: Wiez): { lewy: IdGniazda; prawy: IdGniazda } {
  return wiez.wPrawo
    ? { lewy: wiez.koordynator, prawy: wiez.wykonawca }
    : { lewy: wiez.wykonawca, prawy: wiez.koordynator };
}
