import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';
import { elementIkony, type NazwaIkony } from '../../ikony/ikony';

/**
 * Trzy stany obowiązkowe okna operacyjnego: pusty, ładowania, błędu.
 *
 * Jedna odpowiedzialność: powłoka komunikatu stanu wokół treści okna. Rdzeń
 * odpowiada na każde wywołanie, także odmową, a odmowa ma być widoczna tam,
 * gdzie ją wywołano, nie wyłącznie w konsoli.
 *
 * Forma stanu pustego ma trzy części: ikona, tytuł, opis. Wariantów CSS nie
 * ma — różnicuje wyłącznie treść, a ikona jest treścią, nie wariantem. Stan
 * początkowy mówi, że okno jeszcze nie pytało rdzenia; meldunek o pustce
 * rdzenia postawiony przed jego odpowiedzią byłby zmyśleniem.
 *
 * Nazwa fazy i jej znakowanie pochodzą z biblioteki (`komponenty/faza-okna`):
 * wartość trafia do `data-faza`, po którym sięgają arkusze i sprawdziany, więc
 * własny zestaw wartości w module oznaczałby ten sam stan pod inną nazwą niż
 * w oknach sąsiadów. Tutaj zostaje wyłącznie to, czym Research różni się
 * świadomie: wskaźnik odczytu i chowanie treści na czas ładowania.
 *
 * Stan nie kasuje treści, tylko ją przesłania — nieudane odświeżenie zostawia
 * to, co już było widoczne, a powrót do `gotowe` odsłania treść nietkniętą.
 * Ładowanie niesie wskaźnik `.dn-spinner` obok opisu, nigdy samodzielnie na
 * pełnym ekranie.
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

/** Znak formy stanu pustego dobrany do fazy; `null` znaczy „bez znaku". */
const ZNAK_FAZY: Record<FazaOkna, NazwaIkony | null> = {
  puste: 'info',
  // W ładowaniu znaku nie ma: jego miejsce zajmuje wskaźnik odczytu.
  ladowanie: null,
  blad: 'blad',
  gotowe: null,
};

/** Bok znaku w pikselach — wartość arkusza `.dn-pusty-stan > svg`. */
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
    naglowek.textContent = tytul;
    naglowek.hidden = tytul === '';
    komunikat.textContent = opis;
    wskaznik.hidden = faza !== 'ladowanie';
    // Treść zostaje widoczna także w błędzie: komunikat jest trwały w układzie
    // i poprzedza zawartość, zamiast ją kasować.
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
