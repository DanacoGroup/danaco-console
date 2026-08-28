import { pokazKomunikat } from '../aplikacja/komunikaty';
import { wpisSystemowy } from '../okno-komunikacji/wpis';
import {
  komunikatPrzyjecia,
  komunikatWyslania,
  PRZEKAZANIE_OPIS,
  PRZEKAZANIE_TYTUL,
  TRESC_ZLECENIA_DOMYSLNA,
} from './etykiety-ukladu';
import type { GniazdoOkna } from './gniazdo-okna';
import type { PasRelacji } from './pas-relacji';
import type { StanPary } from './stan-pary';

/** Elementy sceny biorące udział w przekazaniu zlecenia — gniazda źródłowe i docelowe, pas relacji oraz treść komunikatu. */
export interface CzesciPrzekazania {
  /** Gniazdo, które zlecenie wydaje — w pętli koordynator. */
  od: GniazdoOkna;
  /** Gniazdo, które zlecenie przyjmuje — w pętli wykonawca. */
  do_: GniazdoOkna;
  /** Pas relacji, po którym biegnie żeton zlecenia. */
  pas: PasRelacji;
  /** Ustawienie stanu pętli po stronie układu. */
  ustawStan(stan: StanPary): void;
  /** Treść zlecenia widoczna w obu oknach. */
  tresc?: string;
}

/** Czas, po którym scena przechodzi z chwili przekazania zlecenia do stanu pracy wykonawcy, w milisekundach. */
const CZAS_PRZEKAZANIA = 900;

/**
 * Podgląd przekazania zlecenia z okna koordynatora do okna wykonawcy pokazuje na scenie siedem jednoczesnych sygnałów wizualnych, bo kontrakt nie ma osobnej komendy przekazania i moduł nie zapisuje niczego, tylko porusza tym, co widać.
 */
export function pokazPrzekazanieZlecenia(czesci: CzesciPrzekazania): void {
  const { od, do_, pas, ustawStan } = czesci;
  const tresc = czesci.tresc ?? TRESC_ZLECENIA_DOMYSLNA;

  od.fasada.dopisz(wpisSystemowy(komunikatWyslania(do_.id, tresc)));
  do_.fasada.dopisz(wpisSystemowy(komunikatPrzyjecia(od.id, tresc)));

  od.naglowek.blysnij();
  do_.naglowek.blysnij();

  pas.puscZeton();
  pas.oznaczPodglad();
  pokazKomunikat({ tytul: PRZEKAZANIE_TYTUL, tresc: PRZEKAZANIE_OPIS, waga: 'ostrz' });
  ustawStan('przekazanie');

  // Chwila przekazania jest zdarzeniem, nie stanem trwałym; scena wraca do stanu kolejki.
  window.setTimeout(() => ustawStan('wykonawca-pracuje'), CZAS_PRZEKAZANIA);
}
