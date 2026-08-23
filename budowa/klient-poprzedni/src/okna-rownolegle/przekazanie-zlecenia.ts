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

/** Elementy sceny biorące udział w przekazaniu zlecenia. */
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

/** Czas, po którym scena przechodzi z chwili przekazania do pracy wykonawcy. */
const CZAS_PRZEKAZANIA = 900;

/**
 * Podgląd przekazania zlecenia z okna koordynatora do okna wykonawcy.
 *
 * To jest podgląd układu, nie wykonana praca. Kontrakt
 * (`shared/contract.json`) nie ma komendy przekazania zlecenia
 * koordynator→wykonawca: pole `coordinatorWindowId` istnieje po obu stronach,
 * ale samego przekazania nie wywołuje żadna komenda. Moduł nie wysyła więc ani
 * jednej ramki i niczego nie zapisuje — porusza wyłącznie tym, co widać na
 * scenie, i mówi o tym wprost, zamiast zostawiać Operatora z animacją do
 * zinterpretowania.
 *
 * Zdarzenie jest widoczne w obu oknach. Przekazanie odbywa się bez udziału
 * operatora, więc gdyby ślad został tylko u wykonawcy, koordynator pokazywałby
 * rozmowę z dziurą — nie dałoby się odtworzyć, co i kiedy zostało zlecone.
 *
 * Na scenie dzieje się jednocześnie siedem rzeczy:
 *   1. wpis w historii koordynatora — z czym i dokąd, wraz z zastrzeżeniem;
 *   2. wpis w historii wykonawcy — od kogo, wraz z tym samym zastrzeżeniem;
 *   3. błysk obu nagłówków — powiązanie widać, zanim wpisy zostaną przeczytane;
 *   4. żeton biegnący po pasie relacji — kierunek przekazania;
 *   5. dymek — natychmiastowa odpowiedź na czynność Operatora;
 *   6. trwała uwaga na pasie relacji — odpowiedź także kwadrans później;
 *   7. przejście stanu: „przekazanie zlecenia" → „wykonawca pracuje".
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

  // Chwila przekazania jest zdarzeniem, nie stanem trwałym — po przebiegu
  // żetonu scena wraca do stanu, który pokazuje kolejka rdzenia.
  window.setTimeout(() => ustawStan('wykonawca-pracuje'), CZAS_PRZEKAZANIA);
}
