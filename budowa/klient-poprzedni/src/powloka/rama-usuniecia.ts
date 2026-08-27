import './usuniecie-sesji.css';

import { utworzDymekObjasnienia } from '../komponenty/dymek';

// Powierzchnia potwierdzenia usunięcia sesji — same węzły widoku, bez przebiegu czynności.

/** Sesja wskazana do usunięcia: identyfikator sesji w rdzeniu oraz jej tytuł widziany w pasie kart sesji. */
export interface WskazanieUsuniecia {
  id: string;
  tytul: string;
}

/** Części modalu potwierdzenia, po które sięga przebieg całej czynności trwałego usunięcia sesji rdzenia. */
export interface RamaUsuniecia {
  modal: HTMLDialogElement;
  /** Pas trzech stanów obowiązkowych; jego widocznością rządzi `oznaczFaze`. */
  pas: HTMLElement;
  /** Skutek powodzenia — zdanie zostaje na widoku, więc nie idzie na pas. */
  skutek: HTMLElement;
  ostrzezenie: HTMLElement;
  usun: HTMLButtonElement;
  anuluj: HTMLButtonElement;
}

// Oba zdania opisują kosz sesji: usunięcie zdejmuje sesję z historii, a zapis kasuje się trwale po 30 dniach.
const OBJASNIENIE =
  'Usunięcie zdejmuje sesję z historii wraz z wiadomościami, oknami i artefaktami; rdzeń kasuje jej zapis trwale po 30 dniach w koszu. Zamknięcie karty tylko zmienia stan sesji, a archiwizacja pozwala ją przywrócić od ręki.';

const OSTRZEZENIE =
  'Wskazane sesje znikają z historii, a ich zapis — wiadomości, okna komunikacji i artefakty — rdzeń kasuje trwale po 30 dniach w koszu. Katalogu roboczego na dysku ta czynność nie dotyczy.';

/** Wykaz tego, co zginie przy usunięciu — jedna pozycja na każde wskazanie, podpisana tytułem, nie numerem. */
function zlozWykaz(wskazania: readonly WskazanieUsuniecia[]): HTMLElement {
  const wykaz = document.createElement('ul');
  wykaz.className = 'dn-usuwanie__wykaz';
  for (const wskazanie of wskazania) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-usuwanie__pozycja';
    pozycja.textContent = wskazanie.tytul;
    pozycja.dataset['sesja'] = wskazanie.id;
    wykaz.append(pozycja);
  }
  return wykaz;
}

/** Nagłówek modalu z objaśnieniem pod znakiem zapytania — czynność mówi o sobie, zanim padnie kliknięcie. */
function zlozNaglowek(ile: number): HTMLElement {
  const naglowek = document.createElement('div');
  naglowek.className = 'dn-modal-naglowek';
  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = ile === 1 ? 'Usunąć sesję trwale?' : 'Usunąć wskazane sesje trwale?';
  naglowek.append(tytul, utworzDymekObjasnienia(OBJASNIENIE));
  return naglowek;
}

/** Ciało modalu potwierdzenia: wykaz strat, ostrzeżenie, pas stanów obowiązkowych oraz akapit ze skutkiem. */
function zlozCialo(wskazania: readonly WskazanieUsuniecia[]): {
  cialo: HTMLElement;
  ostrzezenie: HTMLElement;
  pas: HTMLElement;
  skutek: HTMLElement;
} {
  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo';

  const ostrzezenie = document.createElement('p');
  ostrzezenie.className = 'dn-usuwanie__ostrzezenie';
  ostrzezenie.textContent = OSTRZEZENIE;

  const pas = document.createElement('p');
  pas.className = 'dn-pusty-stan dn-usuwanie__stan';
  pas.hidden = true;

  const skutek = document.createElement('p');
  skutek.className = 'dn-usuwanie__skutek';
  skutek.hidden = true;
  skutek.setAttribute('role', 'status');

  cialo.append(zlozWykaz(wskazania), ostrzezenie, pas, skutek);
  return { cialo, ostrzezenie, pas, skutek };
}

/**
 * Stopka: odmowa i czynność.
 *
 * Napis „Zostaw sesje" zamiast „Anuluj" — mówi, co się stanie po naciśnięciu,
 * a nie tylko, że okno zniknie. Czynność idzie w odmianie niebezpiecznej, bo
 * zdejmuje sesję z historii od razu.
 */
function zlozStopke(): {
  stopka: HTMLElement;
  anuluj: HTMLButtonElement;
  usun: HTMLButtonElement;
} {
  const stopka = document.createElement('div');
  stopka.className = 'dn-modal-stopka';

  const anuluj = document.createElement('button');
  anuluj.type = 'button';
  anuluj.className = 'dn-btn dn-btn--zarys';
  anuluj.textContent = 'Zostaw sesje';

  const usun = document.createElement('button');
  usun.type = 'button';
  usun.className = 'dn-btn dn-btn--niebezpieczny';
  usun.textContent = 'Usuń trwale';

  stopka.append(anuluj, usun);
  return { stopka, anuluj, usun };
}

/** Budowa modalu potwierdzenia usunięcia; osadzenie w dokumencie i pokazanie należy do przebiegu czynności. */
export function zlozRameUsuniecia(wskazania: readonly WskazanieUsuniecia[]): RamaUsuniecia {
  const modal = document.createElement('dialog');
  modal.className = 'dn-modal dn-usuwanie';

  const { cialo, ostrzezenie, pas, skutek } = zlozCialo(wskazania);
  const { stopka, anuluj, usun } = zlozStopke();

  modal.append(zlozNaglowek(wskazania.length), cialo, stopka);
  return { modal, pas, skutek, ostrzezenie, usun, anuluj };
}
