import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';
import { elementIkony, type NazwaIkony } from '../../ikony/ikony';

/**
 * Cztery stany okna operacyjnego Studio: pusty, ładowania, błędu, gotowy.
 *
 * Rozróżnienie jest treścią: „jeszcze nie pytałem", „pytam" i „rdzeń odmówił"
 * są osobnymi stanami, bo zlanie ich kazałoby Operatorowi zgadywać, czy czekać,
 * czy działać. W Studio różnica jest ostra, bo rdzeń oddaje tu i treść,
 * i pustkę, i odmowę: `document.open` ze ścieżką oddaje dokument pusty,
 * `repository.list` przed pierwszym zapisem — wykaz pusty, `diff.compare` bez
 * wskazanej strony — pustkę, a `document.open` z nieznanym dokumentem odmawia
 * kodem `not_found`. Pusty widok bez powodu byłby nie do odróżnienia od żadnego
 * z tych czterech.
 *
 * Komunikat nie kasuje treści, tylko ją przesłania; wyjątkiem jest ładowanie,
 * które treść ukrywa, bo pokazuje się w jej miejscu.
 *
 * Stan pusty ma jedną formę: znak, tytuł, opis. Znak dochodzi z zestawu
 * (`ikony/ikony`), a nie z własnego rysunku, i jest treścią, nie wariantem —
 * biblioteka `dn-pusty-stan` wariantów CSS nie ma. W fazie ładowania znaku nie
 * ma, bo jego miejsce zajmuje wskaźnik odczytu.
 *
 * Zestaw faz i ich znakowanie pochodzą z `komponenty/faza-okna`: ta sama nazwa
 * fazy trafia do `data-faza` w każdym module, więc wspólny selektor i wspólny
 * sprawdzian mają się o co oprzeć. Tutaj zostaje wyłącznie to, czym Studio
 * różni się świadomie: wskaźnik odczytu i chowanie treści.
 */

/** Powłoka okna wraz z pasem stanu i miejscem na treść. */
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

/** Znak formy stanu pustego dobrany do fazy; `null` znaczy „bez znaku". */
const ZNAK_FAZY: Record<FazaOkna, NazwaIkony | null> = {
  puste: 'info',
  ladowanie: null,
  blad: 'blad',
  gotowe: null,
};

/** Bok znaku w pikselach — wartość arkusza `.dn-pusty-stan > svg`. */
const BOK_ZNAKU = 28;

export function utworzStanOknaStudio(): StanOknaStudio {
  let biezaca: FazaOkna = 'puste';

  const naglowekStanu = document.createElement('p');
  naglowekStanu.className = 'dn-pusty-stan-tytul';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis';

  // Wskaźnik odczytu stoi obok zdania, nie zamiast niego: Operator ma wiedzieć,
  // na co czeka, a nie patrzeć na sam obrót.
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

  // Znak jest wymieniany, a nie chowany: `.dn-pusty-stan > svg` wymaga dziecka
  // bezpośredniego, a `hidden` nie jest własnością elementu SVG.
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
