import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';
import { elementIkony, type NazwaIkony } from '../../ikony/ikony';

/**
 * Powłoka okna operacyjnego Studio wraz z pasem stanu i miejscem na treść;
 * pokrywa cztery fazy — pustą, ładowania, błędu i gotową — każdą osobnym
 * komunikatem, znakiem i wskaźnikiem odczytu, zgodnie z fazami wspólnymi.
 */
export interface StanOknaStudio {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia; treść ustępuje wskaźnikowi. */
  ladowanie(opis: string): void;
  /** Pustka merytoryczna: wywołanie się udało, wyniku nie ma. */
  puste(tytul: string, opis: string): void;
  /** Odmowa albo awaria wraz z jej powodem; komunikat zostaje w układzie. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

/**
 * Znak formy stanu pustego dobrany do bieżącej fazy okna; wartość `null`
 * oznacza fazę bez znaku, w której jego miejsce zajmuje inny element.
 */
const ZNAK_FAZY: Record<FazaOkna, NazwaIkony | null> = {
  puste: 'info',
  ladowanie: null,
  blad: 'blad',
  gotowe: null,
};

/**
 * Bok znaku formy stanu pustego w pikselach; wartość zgodna z regułą arkusza
 * stylów dla znaku osadzonego bezpośrednio w formie `.dn-pusty-stan`.
 */
const BOK_ZNAKU = 28;

export function utworzStanOknaStudio(): StanOknaStudio {
  let biezaca: FazaOkna = 'puste';

  const naglowekStanu = document.createElement('p');
  naglowekStanu.className = 'dn-pusty-stan-tytul';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis';

  // Wskaźnik stoi obok zdania, nie zamiast niego — operator widzi, na co czeka.
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Odczyt z rdzenia w toku');
  wskaznik.hidden = true;

  const wiersz = document.createElement('div');
  wiersz.className = 'ms-stan__wiersz';
  wiersz.append(wskaznik, komunikat);

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan ms-stan';
  powloka.append(naglowekStanu, wiersz);

  const tresc = document.createElement('div');
  tresc.className = 'ms-stan__tresc';

  const element = document.createElement('div');
  element.className = 'ms-stan__powloka';
  element.append(powloka, tresc);

  // Znak jest wymieniany, nie chowany: reguła stylu wymaga dziecka bezpośredniego.
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
    naglowekStanu.textContent = tytul;
    naglowekStanu.hidden = tytul === '';
    komunikat.textContent = opis;
    wskaznik.hidden = faza !== 'ladowanie';
    tresc.hidden = faza === 'ladowanie';
  }

  ustaw(
    'puste',
    'Okno jeszcze nie pytało rdzenia',
    'Studio prowadzi jeden dokument sesji. To jest stan sprzed pierwszego wywołania — ' +
      'nie twierdzenie o tym, co rdzeń ma, bo rdzeń nie został jeszcze o nic zapytany.',
  );

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', '', opis),
    puste: (tytul, opis) => ustaw('puste', tytul, opis),
    blad: (opis) => ustaw('blad', 'Rdzeń odmówił', opis),
    gotowe: () => ustaw('gotowe', '', ''),
    faza: () => biezaca,
  };
}
