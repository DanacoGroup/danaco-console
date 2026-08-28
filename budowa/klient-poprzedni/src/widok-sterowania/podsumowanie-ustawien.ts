import './podsumowanie-ustawien.css';

import { utworzDymekObjasnienia, type KlasyDymka } from '../komponenty/dymek';
import { elementIkony } from '../ikony/ikony';
import type { RejestrKanalow } from '../sterowanie/rejestr-kanalow';
import { pozycjePodsumowania, type PozycjaPodsumowania } from './nazwy-ustawien';
import type { ObserwatorUstawien } from './obserwator-ustawien';

/** Klasy własne dymka pomocy podawane bibliotecznej fabryce; powłoka niesie własny odstęp, a znak zastępuje styl biblioteczny pierścieniem ze wskaźnikiem pomocy. */
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

/** Buduje wiersz listy opisowej: nazwę jako przycisk prowadzący do kontrolki, dymek objaśnienia oraz wartość tekstową ze wskaźnikiem ładowania, gdy dane są w drodze. */
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

/** Buduje wskaźnik ładowania stawiany obok wartości wiersza, gdy wykaz kanałów modelu jest jeszcze w drodze z rdzenia. */
function wskaznikLadowania(): HTMLElement {
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Wykaz kanałów w drodze z rdzenia');
  wskaznik.title = 'Wykaz kanałów w drodze z rdzenia';
  return wskaznik;
}

/** Buduje element tekstowy o wskazanym rodzaju znacznika i klasie, niosący podaną treść jako zawartość tekstową. */
function tekst(rodzaj: string, klasa: string, tresc: string): HTMLElement {
  const element = document.createElement(rodzaj);
  element.className = klasa;
  element.textContent = tresc;
  return element;
}
