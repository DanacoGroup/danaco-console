import './podsumowanie-ustawien.css';

import { utworzDymekObjasnienia, type KlasyDymka } from '../komponenty/dymek';
import { elementIkony } from '../ikony/ikony';
import type { RejestrKanalow } from '../sterowanie/rejestr-kanalow';
import { pozycjePodsumowania, type PozycjaPodsumowania } from './nazwy-ustawien';
import type { ObserwatorUstawien } from './obserwator-ustawien';

/**
 * Podsumowanie ośmiu ustawień okna — odczyt, nie kontrolka.
 *
 * Widoczne wtedy, gdy szuflada z kontrolkami jest zwinięta: operator zamyka
 * sterowanie i nadal wie, na czym okno pracuje. Wartość każdego ustawienia
 * stoi obok jego nazwy, więc osiem odpowiedzi widać jednym spojrzeniem, bez
 * rozwijania czegokolwiek.
 *
 * Podsumowanie nie jest bramą: kliknięcie wiersza otwiera szufladę i prowadzi
 * do kontrolki, a nie odmawia dostępu. Ikona przy każdym wierszu
 * sprawia, że rozróżnienie nie opiera się na samej barwie (dostępność).
 */
/**
 * Klasy własne dymka [?] podawane bibliotecznej fabryce.
 *
 * Powłoka niesie `flex: none` i odstęp z prawej, bez których dymek kurczy się
 * w rzędzie nazwy; znak zastępuje biblioteczne `.dn-btn-ikona`, bo jest
 * pierścieniem 14 px ze wskaźnikiem `help`, a nie kwadratowym przyciskiem
 * ikonowym. Reguły stoją w `podsumowanie-ustawien.css`.
 */
const KLASY_DYMKA: KlasyDymka = {
  powloka: 'dc-widok-ster__dymek',
  znak: 'dc-widok-ster__dymek-znak',
};

export interface PodsumowanieUstawien {
  /** Element montowany w widoku sterowania. */
  element: HTMLElement;
  /** Odłącza subskrypcje stanu i rejestru. */
  rozlacz(): void;
}

export function utworzPodsumowanie(
  obserwator: ObserwatorUstawien,
  rejestr: RejestrKanalow,
  przyWskazaniu: (klucz: string) => void,
): PodsumowanieUstawien {
  const element = document.createElement('dl');
  element.className = 'dc-widok-ster__podsumowanie';
  element.setAttribute('aria-label', 'Bieżące ustawienia okna');

  function odrysuj(): void {
    element.replaceChildren();
    for (const pozycja of pozycjePodsumowania(obserwator.migawka(), rejestr)) {
      element.append(wiersz(pozycja, przyWskazaniu));
    }
  }

  const odsubskrybujStan = obserwator.naZmiane(odrysuj);
  const odsubskrybujRejestr = rejestr.naZmiane(odrysuj);
  odrysuj();

  return {
    element,
    rozlacz() {
      odsubskrybujStan();
      odsubskrybujRejestr();
    },
  };
}

/**
 * Pozycja listy opisowej: nazwa ustawienia i jego bieżąca wartość.
 *
 * Para `dt`–`dd` mieszka we własnym bloku — dopuszcza to budowa listy
 * opisowej, a układ zyskuje jedną komórkę siatki na ustawienie zamiast
 * dwóch niezależnych, które przy zmianie liczby kolumn rozjechałyby się
 * względem siebie.
 *
 * Nazwa jest przyciskiem prowadzącym do kontrolki — nie oznakowaniem, którego
 * nie da się nacisnąć. Wartość zostaje tekstem, bo jest odczytem stanu.
 * Obok nazwy stoi dymek [?] z objaśnieniem ustawienia, a wiersz czekający na
 * dane rdzenia niesie wskaźnik ładowania obok wartości — nigdy zamiast niej.
 */
function wiersz(
  pozycja: PozycjaPodsumowania,
  przyWskazaniu: (klucz: string) => void,
): HTMLElement {
  const nazwa = document.createElement('dt');
  nazwa.className = 'dc-widok-ster__nazwa';

  const przejscie = document.createElement('button');
  przejscie.type = 'button';
  przejscie.className = 'dc-widok-ster__przejscie';
  przejscie.title = `Otwórz sterowanie: ${pozycja.etykieta}`;
  przejscie.append(
    elementIkony(pozycja.ikona, { rozmiar: 16 }),
    tekst('span', 'dc-widok-ster__etykieta', pozycja.etykieta),
  );
  przejscie.addEventListener('click', () => przyWskazaniu(pozycja.klucz));
  nazwa.append(przejscie, utworzDymekObjasnienia(pozycja.objasnienie, KLASY_DYMKA));

  const wartosc = document.createElement('dd');
  wartosc.className = 'dc-widok-ster__wartosc';
  wartosc.dataset.ustawienie = pozycja.klucz;
  wartosc.textContent = pozycja.wartosc;
  wartosc.title = pozycja.wartosc;
  if (pozycja.ladowanie === true) wartosc.append(wskaznikLadowania());

  const blok = document.createElement('div');
  blok.className = 'dc-widok-ster__pozycja';
  blok.dataset.ustawienie = pozycja.klucz;
  blok.append(nazwa, wartosc);
  return blok;
}

/** Wskaźnik ładowania obok wartości — wykaz kanałów w drodze z rdzenia. */
function wskaznikLadowania(): HTMLElement {
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Wykaz kanałów w drodze z rdzenia');
  wskaznik.title = 'Wykaz kanałów w drodze z rdzenia';
  return wskaznik;
}

/** Element tekstowy o wskazanej klasie. */
function tekst(rodzaj: string, klasa: string, tresc: string): HTMLElement {
  const element = document.createElement(rodzaj);
  element.className = klasa;
  element.textContent = tresc;
  return element;
}
