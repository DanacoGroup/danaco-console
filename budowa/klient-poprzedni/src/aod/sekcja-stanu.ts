import type { AodStatus } from '../../../shared/contract';
import { dodajPole, utworzAkapit, utworzPodtytul, utworzWykazPol } from './pola-wykazu';

/**
 * Sekcja stanu nakładki (`aod.status.get`) — nadzorca sesji i telemetria.
 *
 * Wypisuje to, co rdzeń zmierzył: urządzenie, kartę sesji i okno ogniskowane,
 * liczbę procesów w biegu i chwilę pomiaru. Wykaz procesów przypiętych niesie
 * sekcja obecności, bo tam Operator przypina i odpina — wykaz stoi przy
 * czynności, która go zmienia.
 *
 * Pola puste zostają puste — `(brak)` zamiast wartości zmyślonej po stronie
 * widoku. Świeża nakładka bez ogniska jest stanem poprawnym.
 */
export interface SekcjaStanu {
  element: HTMLElement;
  pokaz(status: AodStatus): void;
  /**
   * Nanosi odmowę rdzenia w tę jedną sekcję.
   *
   * Odmowa jednego odczytu jest faktem o jednym odczycie, nie o całej
   * nakładce; pozostałe sekcje stoją na innych komendach i zostają widoczne.
   */
  odmowa(zdanie: string): void;
}

export function utworzSekcjeStanu(): SekcjaStanu {
  const element = document.createElement('section');
  element.className = 'ao-sekcja';

  const miejsce = document.createElement('div');
  miejsce.className = 'ao-sekcja__tresc';

  element.append(utworzPodtytul('Stan nakładki (aod.status.get)'), miejsce);

  return {
    element,

    pokaz(status) {
      const lista = utworzWykazPol();
      dodajPole(lista, 'Urządzenie', status.deviceId ?? '');
      dodajPole(lista, 'Karta sesji ogniskowana', status.activeSessionId ?? '');
      dodajPole(lista, 'Okno ogniskowane', status.activeWindowId ?? '');
      dodajPole(lista, 'Procesy w biegu', String(status.runningProcessCount));
      dodajPole(lista, 'Zmierzono', new Date(status.updatedAt).toLocaleString());
      miejsce.replaceChildren(lista);
    },

    odmowa(zdanie) {
      miejsce.replaceChildren(utworzAkapit('ao-odmowa', zdanie));
    },
  };
}
