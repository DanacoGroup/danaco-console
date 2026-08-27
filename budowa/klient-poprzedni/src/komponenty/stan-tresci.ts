import type { ErrorInfo } from '../../../shared/contract';

/**
 * Stan treści okna operacyjnego, będący jednym nośnikiem dla okien wszystkich
 * modułów. Niesie trzy stany obowiązkowe każdego okna — ładowanie, pustkę
 * i błąd — oraz stan czwarty, którym jest treść.
 */
export interface StanTresci {
  /** Element montowany w miejscu treści okna. */
  element: HTMLElement;
  /** Zapowiedź odczytu — zostaje do przyjścia odpowiedzi. */
  ladowanie(opis?: string): void;
  /** Odczyt zakończony brakiem pozycji. */
  pusto(zdanie: string): void;
  /** Odczyt nieudany — treść odmowy wprost z rdzenia. */
  blad(zdanie: string, powod?: ErrorInfo): void;
  /** Odczyt udany — miejsce na treść oddane wywołującemu. */
  tresc(): HTMLElement;
  /** Krótkie potwierdzenie czynności; nie zastępuje treści. */
  potwierdzenie(zdanie: string, udane: boolean): void;
  /** Rodzaj stanu bieżącego — do sprawdzianów i decyzji widoku. */
  rodzaj(): string;
}

/**
 * Buduje nośnik stanu treści dla okna modułu. Przedrostek klas modułu podawany
 * jest bez kreski; powstają z niego cztery klasy stanów, komunikatu,
 * potwierdzenia i treści, które pokrywa arkusz modułu.
 */
export function utworzStanTresci(przedrostek: string): StanTresci {
  const komunikat = document.createElement('p');
  komunikat.className = `${przedrostek}-stan`;
  komunikat.setAttribute('aria-live', 'polite');

  const potwierdzenie = document.createElement('p');
  potwierdzenie.className = `${przedrostek}-potwierdzenie`;
  potwierdzenie.hidden = true;
  potwierdzenie.setAttribute('role', 'status');

  const miejsce = document.createElement('div');
  miejsce.className = `${przedrostek}-tresc`;

  const element = document.createElement('div');
  element.className = `${przedrostek}-stany`;
  element.append(komunikat, potwierdzenie, miejsce);

  function ustaw(rodzaj: string, zdanie: string): void {
    komunikat.dataset['stan'] = rodzaj;
    komunikat.textContent = zdanie;
    komunikat.hidden = zdanie === '';
    miejsce.replaceChildren();
  }

  return {
    element,

    ladowanie(opis = 'Odczyt w toku…') {
      ustaw('ladowanie', opis);
      potwierdzenie.hidden = true;
    },

    pusto(zdanie) {
      ustaw('pusto', zdanie);
    },

    blad(zdanie, powod) {
      ustaw('blad', `${zdanie} ${opisBledu(powod)}`.trim());
    },

    tresc() {
      ustaw('tresc', '');
      return miejsce;
    },

    potwierdzenie(zdanie, udane) {
      potwierdzenie.textContent = zdanie;
      potwierdzenie.dataset['udane'] = String(udane);
      potwierdzenie.hidden = zdanie === '';
    },

    rodzaj: () => komunikat.dataset['stan'] ?? '',
  };
}

/**
 * Treść odmowy pokazywana w oknie. Kod kontraktu zostaje w nawiasie, ponieważ
 * to po nim rozpoznaje się komendę bez obsługiwacza i odróżnia odmowę
 * uprawnienia od usterki rdzenia oraz od odmowy merytorycznej.
 */
function opisBledu(powod?: ErrorInfo): string {
  if (powod === undefined) return '';
  const tresc = powod.message === '' ? 'rdzeń nie podał przyczyny' : powod.message;
  return `Powód: ${tresc} (kod ${powod.code}).`;
}
