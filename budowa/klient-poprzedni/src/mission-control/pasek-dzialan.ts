import { elementIkony } from '../ikony/ikony';

/**
 * Pasek ostatniego działania — potwierdzenie, że naciśnięcie zostało przyjęte.
 *
 * Jedna odpowiedzialność: wypisanie, co pulpit właśnie nadał powłoce. Żaden
 * przycisk pulpitu nie jest wyszarzony, a przycisk czynny bez widocznego skutku
 * nie mówi operatorowi, czy zamiar poszedł dalej — pasek pokazuje więc fakt
 * nadania zamiaru, nie wynik wykonania.
 *
 * Pasek pisze „zamiar nadany", nigdy „wykonano": powłoka mogła jeszcze nie
 * wysłać komendy, a rdzeń mógł jej nie potwierdzić.
 */
export interface PasekDzialan {
  element: HTMLElement;
  /** Wypisuje zamiar wraz z godziną nadania. */
  zapisz(tresc: string): void;
}

/** Buduje pasek ostatniego działania. */
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

/** Godzina nadania w postaci `gg:mm:ss`. */
function godzina(): string {
  return new Date().toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
