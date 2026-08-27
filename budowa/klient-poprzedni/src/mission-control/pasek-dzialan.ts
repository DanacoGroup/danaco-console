import { elementIkony } from '../ikony/ikony';

/**
 * Pasek ostatniego działania potwierdza, że naciśnięcie zostało przyjęte:
 * wypisuje, co pulpit nadał powłoce. Treść mówi „zamiar nadany", nigdy
 * „wykonano", ponieważ pasek podaje fakt nadania zamiaru, a nie wynik jego
 * wykonania.
 */
export interface PasekDzialan {
  element: HTMLElement;
  /** Wypisuje zamiar wraz z godziną nadania. */
  zapisz(tresc: string): void;
}

/**
 * Buduje pasek ostatniego działania jako stopkę o roli `status` i obszarze
 * powiadamiania w trybie `polite`, ze znakiem informacyjnym oraz treścią
 * początkową mówiącą, że żaden zamiar nie został jeszcze nadany.
 */
export function utworzPasekDzialan(): PasekDzialan {
  const element = document.createElement('footer');
  element.className = 'mc-dzialania';
  element.setAttribute('role', 'status');
  element.setAttribute('aria-live', 'polite');

  const znak = document.createElement('span');
  znak.className = 'mc-dzialania__znak';
  znak.append(elementIkony('info', { rozmiar: 16 }));

  const tresc = document.createElement('span');
  tresc.className = 'mc-dzialania__tresc';
  tresc.textContent = 'Pulpit gotowy. Żaden zamiar nie został jeszcze nadany.';

  element.append(znak, tresc);

  return {
    element,
    zapisz(opis) {
      tresc.textContent = `${godzina()} — zamiar nadany: ${opis}`;
    },
  };
}

/**
 * Podaje godzinę nadania w postaci `gg:mm:ss`, złożoną według ustawień
 * regionalnych `pl-PL` z dwucyfrowej godziny, minuty i sekundy.
 */
function godzina(): string {
  return new Date().toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
