import { ETYKIETA_PRZELACZNIKA, OBJASNIENIE_PRZELACZNIKA, OPIS_PRZELACZNIKA } from './etykiety-ukladu';
import { LICZBA_MAX, LICZBA_MIN, ograniczLiczbe } from './identyfikatory';
import { utworzObjasnienie } from './objasnienie-ukladu';

/** Identyfikator zdania o figurze modułu, przez który pozycje ponad figurą wskazują je jako swój opis dostępności. */
const ID_UWAGI = 'dn-okna-przelacznik-figura';

/** Przełącznik liczby okien komunikacji na scenie — grupa pól jednokrotnego wyboru z zawsze czynnymi pozycjami. */
export interface PrzelacznikLiczby {
  element: HTMLElement;
  /** Odznacza wskazaną liczbę bez wywoływania słuchacza. */
  pokaz(liczba: number): void;
  /** Wypowiada figurę rozmowy modułu; liczba mniejsza od granicy nie wyłącza pozycji wyższych. */
  ustawFigure(liczbaOkien: number, zdanie: string): void;
}

/**
 * Przełącznik liczby okien komunikacji ma wszystkie pozycje czynne zawsze, jako grupę pól jednokrotnego wyboru z dymkiem objaśnienia i opisem skutku wyboru.
 */
export function utworzPrzelacznikLiczby(
  liczbaPoczatkowa: number,
  naWybor: (liczba: number) => void,
): PrzelacznikLiczby {
  const element = document.createElement('div');
  element.className = 'dn-okna__przelacznik';

  const etykieta = document.createElement('span');
  etykieta.className = 'dn-okna__przelacznik-etykieta';
  etykieta.id = 'dn-okna-przelacznik-etykieta';
  etykieta.textContent = ETYKIETA_PRZELACZNIKA;

  const grupa = document.createElement('div');
  grupa.className = 'dn-zakladki dn-okna__pigulki';
  grupa.setAttribute('role', 'radiogroup');
  grupa.setAttribute('aria-labelledby', etykieta.id);

  const pozycje = new Map<number, HTMLButtonElement>();

  for (let liczba = LICZBA_MIN; liczba <= LICZBA_MAX; liczba += 1) {
    const pozycja = document.createElement('button');
    pozycja.type = 'button';
    pozycja.className = 'dn-zakladka dn-okna__pigulka';
    pozycja.setAttribute('role', 'radio');
    pozycja.textContent = String(liczba);
    pozycja.setAttribute(
      'aria-label',
      liczba === 1 ? 'Jedno okno komunikacji' : `${liczba} okna komunikacji`,
    );
    pozycja.addEventListener('click', () => naWybor(liczba));
    pozycje.set(liczba, pozycja);
    grupa.append(pozycja);
  }

  const pomoc = document.createElement('span');
  pomoc.className = 'dn-pole-opis dn-okna__pomoc';
  pomoc.textContent = OPIS_PRZELACZNIKA;

  // Zdanie o figurze rozmowy modułu stoi pod grupą i mówi o module tyle, ile niesie jego profil.
  const figura = document.createElement('p');
  figura.className = 'dn-pole-opis dn-okna__figura';
  figura.id = ID_UWAGI;
  figura.hidden = true;

  element.append(etykieta, utworzObjasnienie(OBJASNIENIE_PRZELACZNIKA), grupa, pomoc, figura);

  function pokaz(liczba: number): void {
    const wybrana = ograniczLiczbe(liczba);
    for (const [wartosc, pozycja] of pozycje) {
      const czynna = wartosc === wybrana;
      pozycja.setAttribute('aria-checked', String(czynna));
      // Ognisko wędruje wraz z wyborem — grupa zachowuje jedno wejście z Tab.
      pozycja.tabIndex = czynna ? 0 : -1;
    }
  }

  function ustawFigure(liczbaOkien: number, zdanie: string): void {
    figura.textContent = zdanie;
    figura.hidden = zdanie.length === 0;

    for (const [wartosc, pozycja] of pozycje) {
      // Pozycja ponad figurą modułu zostaje czynna i bierze opis, by odczyt podał powód razem z liczbą.
      const ponadFigura = zdanie.length > 0 && wartosc > liczbaOkien;
      if (ponadFigura) pozycja.setAttribute('aria-describedby', ID_UWAGI);
      else pozycja.removeAttribute('aria-describedby');
    }
  }

  pokaz(liczbaPoczatkowa);
  return { element, pokaz, ustawFigure };
}
