import type { ErrorInfo } from '../../../shared/contract';

/**
 * Stan treści okna operacyjnego — jeden nośnik dla okien wszystkich modułów.
 *
 * Trzy stany obowiązkowe każdego okna: ładowanie, pustka, błąd — plus stan
 * czwarty, czyli treść.
 *
 * To nie jest to samo co `faza-okna.ts`. Tamten byt niesie fazy ramy
 * (`data-faza` na powłoce, rola dostępności pasa); ten niesie miejsce treści
 * wraz z pasem komunikatu i potwierdzeniem czynności, a stan trzyma
 * w `data-stan` na akapicie komunikatu — po tym atrybucie sięgają arkusze
 * modułów.
 *
 * Stan błędu niesie treść błędu z kontraktu (kod i komunikat), bo okno
 * pokazujące po odmowie pusty wykaz mówiłoby „nic nie ma” zamiast „nie udało
 * się zapytać”. Stan pusty ma własne zdanie, bo pustka bywa poprawna:
 * instalacja bez automatyk czy repozytorium bez zmian do zatwierdzenia nie są
 * usterkami.
 *
 * Wygląd zostaje w module: klasy noszą przedrostek modułu i pokrywa je arkusz
 * modułu, dlatego przedrostek jest wartością wejściową, a nie stałą.
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
 * Buduje nośnik stanu treści dla okna modułu.
 *
 * @param przedrostek przedrostek klas modułu bez kreski — `da`, `mdev`, `dg`,
 *   `dm`, `dr`, `dt`, `dw`. Z niego powstają `…-stany`, `…-stan`,
 *   `…-potwierdzenie` i `…-tresc`, czyli dokładnie te cztery klasy, które
 *   pokrywa arkusz modułu.
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
 * Treść odmowy pokazywana w oknie. Kod kontraktu zostaje w nawiasie, bo to po
 * nim rozpoznaje się komendę bez obsługiwacza (`*.unknown`), odróżnia odmowę
 * uprawnienia (`permission_denied`) od usterki rdzenia i jedną i drugą od
 * odmowy merytorycznej.
 *
 * `komponenty/odmowa.ts` składa zdanie w innym porządku (`czynność: powód
 * (kod)`) i wymaga nazwy czynności, której te okna nie podają.
 */
function opisBledu(powod?: ErrorInfo): string {
  if (powod === undefined) return '';
  const tresc = powod.message === '' ? 'rdzeń nie podał przyczyny' : powod.message;
  return `Powód: ${tresc} (kod ${powod.code}).`;
}
