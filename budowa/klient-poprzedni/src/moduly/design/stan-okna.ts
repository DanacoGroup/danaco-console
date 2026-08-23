/**
 * Trzy stany obowiązkowe okna operacyjnego modułu Design: pusty, ładowanie,
 * błąd. Jedna odpowiedzialność: powłoka komunikatu stanu wokół treści okna.
 *
 * Komunikat przesłania treść, a nie zastępuje jej: nieudane odświeżenie nie
 * kasuje tego, co Operator już widział, i treść wraca nietknięta, gdy okno
 * wróci do fazy `gotowe`. Odmowa rdzenia jest widoczna w oknie, w którym
 * Operator ją wywołał, a nie wyłącznie w konsoli.
 *
 * Wygląd bierze się w całości z biblioteki `komponenty/` (`dn-pusty-stan`,
 * `dn-spinner`) — plik nie zna ani jednej barwy.
 *
 * Nazwy faz i znakowanie powłoki są wspólne dla wszystkich modułów
 * (`komponenty/faza-okna.ts`); tutaj zostaje wyłącznie to, czym Design różni
 * się świadomie: wskaźnik odczytu i chowanie treści na czas ładowania.
 */
import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/** Powłoka okna wraz z komunikatem stanu. */
export interface StanOkna {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pokazuje pustkę merytoryczną — wywołanie się udało, wyniku nie ma. */
  puste(opis: string): void;
  /** Pokazuje odmowę albo awarię wraz z jej powodem. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

export function utworzStanOkna(): StanOkna {
  let biezaca: FazaOkna = 'puste';

  // Spinner stoi obok komunikatu, nigdy samodzielnie na pełnym ekranie.
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner md-stan__wskaznik';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Odczyt w toku');

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis md-stan__opis';

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan md-stan';
  powloka.append(wskaznik, komunikat);

  const tresc = document.createElement('div');
  tresc.className = 'md-stan__tresc';

  const element = document.createElement('div');
  element.className = 'md-stan__powloka';
  element.append(powloka, tresc);

  function ustaw(faza: FazaOkna, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, powloka, faza);
    komunikat.textContent = opis;
    wskaznik.hidden = faza !== 'ladowanie';
    // Treść zostaje widoczna także w błędzie: komunikat ją poprzedza,
    // nie zabiera.
    tresc.hidden = faza === 'ladowanie' && !czyStoiOtwartyModal();
  }

  /**
   * Czy w treści okna stoi otwarty modal — jeśli stoi, treści schować nie wolno.
   *
   * Kreator Prompt Buildera stoi na natywnym `<dialog>` otwartym przez
   * `showModal()`, a modal otwarty czyni resztę dokumentu bezwładną. Schowanie
   * treści okna na czas ładowania schowałoby też ten modal — bezwładność by
   * została i strona nie miałaby ani jednej drogi wyjścia.
   */
  function czyStoiOtwartyModal(): boolean {
    return tresc.querySelector('dialog[open]') !== null;
  }

  ustaw('puste', 'Brak danych.');

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', opis),
    puste: (opis) => ustaw('puste', opis),
    blad: (opis) => ustaw('blad', opis),
    gotowe: () => ustaw('gotowe', ''),
    faza: () => biezaca,
  };
}

/** Nagłówek okna operacyjnego wraz z jego rolą z katalogu rdzenia. */
export function naglowekOkna(nazwa: string, rola: string): HTMLElement {
  const tytul = document.createElement('h3');
  tytul.className = 'md-okno__tytul';
  tytul.textContent = nazwa;

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--rola';
  plakietka.textContent = rola;

  const element = document.createElement('header');
  element.className = 'md-okno__naglowek';
  element.append(tytul, plakietka);
  return element;
}
