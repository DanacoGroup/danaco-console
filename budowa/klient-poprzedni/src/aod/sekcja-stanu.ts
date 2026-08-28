import type { AodStatus } from '../../../shared/contract';
import { dodajPole, utworzAkapit, utworzPodtytul, utworzWykazPol } from './pola-wykazu';

/**
 * Sekcja stanu nakładki (`aod.status.get`), czyli nadzorca sesji i telemetria.
 * Wypisuje to, co zmierzył rdzeń: urządzenie, kartę sesji i okno ogniskowane,
 * liczbę procesów w biegu oraz chwilę pomiaru. Pola puste zostają puste.
 */
export interface SekcjaStanu {
  element: HTMLElement;
  pokaz(status: AodStatus): void;
  /** Nanosi odmowę rdzenia w tę jedną sekcję, nie ruszając pozostałych. */
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
