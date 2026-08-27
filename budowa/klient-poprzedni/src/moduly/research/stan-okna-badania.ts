import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';
import { elementIkony, type NazwaIkony } from '../../ikony/ikony';

/**
 * Trzy stany obowiązkowe okna operacyjnego — pusty, ładowania, błędu — jako powłoka komunikatu
 * wokół treści okna.
 */
export interface StanOknaBadania {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pustka merytoryczna — wywołanie się udało, wyniku nie ma. */
  puste(tytul: string, opis: string): void;
  /** Odmowa albo awaria wraz z powodem. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

/** Znak formy stanu pustego dobrany do fazy okna; wartość null znaczy brak znaku, bo ładowanie ma własny wskaźnik. */
const ZNAK_FAZY: Record<FazaOkna, NazwaIkony | null> = {
  puste: 'info',
  ladowanie: null,
  blad: 'blad',
  gotowe: null,
};

/** Bok znaku formy stanu pustego w pikselach, dopasowany do wartości arkusza stylu, jaki okno stosuje dla tego elementu. */
const BOK_ZNAKU = 28;

export function utworzStanOknaBadania(): StanOknaBadania {
  let biezaca: FazaOkna = 'puste';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner mr-stan__wskaznik';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Trwa odczyt z rdzenia');

  const naglowek = document.createElement('p');
  naglowek.className = 'dn-pusty-stan-tytul';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis mr-stan__opis';

  const wiersz = document.createElement('div');
  wiersz.className = 'mr-stan__wiersz';
  wiersz.append(wskaznik, komunikat);

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan mr-stan';
  powloka.append(naglowek, wiersz);

  const tresc = document.createElement('div');
  tresc.className = 'mr-stan__tresc';

  const element = document.createElement('div');
  element.className = 'mr-stan__powloka';
  element.append(powloka, tresc);

  // Znak jest wymieniany, nie chowany: znacznik hidden nie działa na elemencie SVG jako dziecku.
  let znak: SVGSVGElement | null = null;
  let nazwaZnaku: NazwaIkony | null = null;

  function ustawZnak(nazwa: NazwaIkony | null): void {
    if (nazwa === nazwaZnaku) return;
    nazwaZnaku = nazwa;
    if (nazwa === null) {
      znak?.remove();
      znak = null;
      return;
    }
    const nowy = elementIkony(nazwa, { rozmiar: BOK_ZNAKU });
    if (znak === null) powloka.prepend(nowy);
    else znak.replaceWith(nowy);
    znak = nowy;
  }

  function ustaw(faza: FazaOkna, tytul: string, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, powloka, faza);
    ustawZnak(ZNAK_FAZY[faza]);
    naglowek.textContent = tytul;
    naglowek.hidden = tytul === '';
    komunikat.textContent = opis;
    wskaznik.hidden = faza !== 'ladowanie';
    // Treść zostaje widoczna w błędzie: komunikat jest trwały w układzie i poprzedza ją, zamiast kasować.
    tresc.hidden = faza === 'ladowanie';
  }

  ustaw(
    'puste',
    'Okno jeszcze nie pytało rdzenia',
    'Moduł Research prowadzi jedno badanie: zakres, źródła, ustalenia, raport i jego wydanie. ' +
      'To jest stan sprzed pierwszego wywołania, nie twierdzenie o tym, co rdzeń ma.',
  );

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', 'Odczyt z rdzenia', opis),
    puste: (tytul, opis) => ustaw('puste', tytul, opis),
    blad: (opis) => ustaw('blad', 'Rdzeń odmówił', opis),
    gotowe: () => ustaw('gotowe', '', ''),
    faza: () => biezaca,
  };
}
