import { profilModulu } from '../okno-komunikacji/rejestr-profilow';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Ulotność rozmowy: moduł bez pamięci sesyjnej, którego zakres kończy się na kliencie.
 */

/**
 * Kontekst roboczy rozmowy ulotnej, czyli to, czego zmiana kończy dotychczasową rozmowę i
 * zaczyna zupełnie nową.
 */
export interface KontekstRoboczy {
  /** Klucz kontekstu; pusty znaczy „kontekstu jeszcze nie wskazano". */
  klucz(): string;
  // Subskrypcja zmiany kontekstu, niosąca gotowe zdanie o tym, co się zmieniło.
  naZmiane(sluchacz: (zdanie: string) => void): Odsubskrybuj;
}

/**
 * Reguła pamięci obowiązująca w rozmowie jednego okna, niosąca wraz z sobą gotowe zdania
 * przeznaczone dla Operatora.
 */
export interface PolitykaUlotnosci {
  /** Kod modułu, którego polityka dotyczy. */
  modul: string;
  /** Czy rozmowa przeżywa zamknięcie okna; fałsz znaczy „rozmowa ulotna". */
  pamiecSesyjna: boolean;
  // Zdania stanu wypisywane przy otwarciu rozmowy ulotnej i po każdym czyszczeniu.
  zapowiedz: readonly string[];
  /** Kontekst, którego zmiana czyści rozmowę; `null` = takiego kontekstu nie ma. */
  kontekst: KontekstRoboczy | null;
}

/**
 * Zdanie mówiące, dokąd sięga ulotność, a dokąd nie. Bez niego określenie
 * „czat roboczy" czytałoby się jako obietnica kasowania, którego rdzeń nie
 * wykonuje.
 */
export const ZDANIE_O_ZAPISIE_RDZENIA =
  'Zasięg: okno nie odtwarza i nie pokazuje tej rozmowy ponownie. Rdzeń zapisuje ' +
  'wiadomości w swojej bazie mimo to — kontrakt nie ma dziś komendy ich kasowania ' +
  '(obszar message.* to send, stop, list).';

/**
 * Polityka rozmowy z pamięcią — stan zwykły okna rozmowy, bez ani jednego dodatkowego
 * zdania widocznego w interfejsie.
 */
export function politykaTrwala(modul: string): PolitykaUlotnosci {
  return { modul, pamiecSesyjna: true, zapowiedz: [], kontekst: null };
}

/**
 * Polityka okna pracującego we wskazanym module roboczym, wyprowadzona wprost z profilu
 * przypisanego temu modułowi.
 */
export function politykaModulu(
  modul: string,
  kontekst: KontekstRoboczy | null = null,
): PolitykaUlotnosci {
  const profil = profilModulu(modul);
  if (profil.pamiecSesyjna) return politykaTrwala(modul);

  const zapowiedz = [
    `Moduł ${profil.nazwa}: czat roboczy BEZ pamięci sesyjnej. Ten wątek nie zostanie ` +
      'odtworzony po zamknięciu okna' +
      (kontekst === null ? '.' : ' ani po zmianie kontekstu testowania.'),
    ZDANIE_O_ZAPISIE_RDZENIA,
  ];

  return { modul, pamiecSesyjna: false, zapowiedz, kontekst };
}

/**
 * Czy obowiązująca w tym oknie polityka pamięci każe bieżącej rozmowie być ulotną, czy
 * pozostać trwałą.
 */
export function czyUlotna(polityka: PolitykaUlotnosci | null): boolean {
  return polityka !== null && !polityka.pamiecSesyjna;
}
