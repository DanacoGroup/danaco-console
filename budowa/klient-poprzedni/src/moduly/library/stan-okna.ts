import { elementIkony } from '../../ikony/ikony';
import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/**
 * Cztery stany okna modułu Library: puste, ładowanie, błąd, gotowe.
 *
 * Stany są rozdzielne z powodu: „jeszcze nie pytałem", „pytam" i „rdzeń nie ma
 * ani jednego pliku" to trzy różne rzeczy, a zlanie ich w jedno kazałoby
 * Operatorowi zgadywać, czy czekać, czy działać.
 *
 * Stan nie kasuje treści, tylko ją przesłania: nieudane odświeżenie zostawia to,
 * co Operator już widział, i dokłada powód. Stan pusty niesie dokładnie jeden
 * przycisk pierwszej akcji.
 *
 * Forma stanu pustego jest jedna: ikona · tytuł · opis, treść wyśrodkowana, bez
 * wariantów klasy — różnicuje ją sama treść. Stopnie pisma i szerokość łamania
 * niesie wspólny `komponenty/drobne.css` (`.dn-pusty-stan-tytul`,
 * `.dn-pusty-stan-opis`), nie arkusz modułu; nadpisanie ich u siebie dałoby dwie
 * formy tego samego stanu.
 *
 * Ikona należy wyłącznie do stanu pustego. Ładowanie ma własny nośnik
 * (`.dn-spinner`), a błąd — kreskę po lewej; trzeci rysunek nad nimi niczego by
 * nie dopowiedział. Widoczność rozstrzyga arkusz po `data-faza`, tak samo jak
 * kreskę błędu.
 *
 * Nazwy faz i znakowanie powłoki pochodzą ze wspólnego `komponenty/faza-okna`,
 * żeby ten sam stan nazywał się w całym drzewie tak samo.
 */

/** Powłoka okna wraz z komunikatem stanu. */
export interface StanOkna {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pustka merytoryczna: wywołanie się udało, wyniku nie ma. */
  puste(tytul: string, opis: string): void;
  /** Odmowa albo awaria wraz z powodem podanym przez rdzeń. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do decyzji widoku i do sprawdzianu. */
  faza(): FazaOkna;
  /** Przycisk pierwszej akcji stanu pustego; ukryty, dopóki nie ustawiony. */
  pierwszaAkcja(etykieta: string, czynnosc: () => void): void;
}

export function utworzStanOkna(): StanOkna {
  let biezaca: FazaOkna = 'puste';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis';

  const akcja = document.createElement('button');
  akcja.type = 'button';
  akcja.className = 'dn-btn dn-btn--sm dn-btn--zarys ml-stan__akcja';
  akcja.hidden = true;

  const ikona = elementIkony('biblioteka', { rozmiar: 24, klasa: 'dn-ikona ml-stan__ikona' });

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan ml-stan';
  powloka.append(ikona, tytul, komunikat, akcja);

  const spinner = document.createElement('span');
  spinner.className = 'dn-spinner ml-stan__spinner';
  spinner.setAttribute('role', 'status');
  spinner.setAttribute('aria-label', 'Trwa odczyt z rdzenia');
  spinner.hidden = true;
  tytul.append(spinner);

  const tresc = document.createElement('div');
  tresc.className = 'ml-stan__tresc';

  const element = document.createElement('div');
  element.className = 'ml-stan__powloka';
  element.append(powloka, tresc);

  function ustaw(faza: FazaOkna, naglowek: string, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, powloka, faza);
    tytul.replaceChildren(naglowek, spinner);
    komunikat.textContent = opis;
    spinner.hidden = faza !== 'ladowanie';
    akcja.hidden = akcja.textContent === '' || faza !== 'puste';
  }

  // Zdanie początkowe mówi, co się zaraz stanie, zamiast nazywać brak danych.
  // Okno stoi w tym stanie przez chwilę między zbudowaniem układu a pierwszą
  // odpowiedzią rdzenia i jest to wtedy jedyne zdanie, które Operator czyta.
  ustaw(
    'puste',
    'Okno czeka na pierwszy odczyt',
    'Moduł Library pyta rdzeń o repozytorium zaraz po wejściu — do tej chwili okno nie ' +
      'ma o czym mówić.',
  );

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', 'Odczyt z rdzenia', opis),
    puste: (naglowek, opis) => ustaw('puste', naglowek, opis),
    blad: (opis) => ustaw('blad', 'Rdzeń odmówił', opis),
    gotowe: () => ustaw('gotowe', '', ''),
    faza: () => biezaca,

    pierwszaAkcja(etykieta, czynnosc) {
      akcja.textContent = etykieta;
      akcja.hidden = biezaca !== 'puste';
      akcja.addEventListener('click', czynnosc);
    },
  };
}
