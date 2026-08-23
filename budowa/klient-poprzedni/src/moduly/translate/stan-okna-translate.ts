/**
 * Trzy stany obowiązkowe okna modułu Translate: pusty, ładowania, błędu.
 *
 * Jedna forma, różnicowana wyłącznie treścią — stąd brak wariantów klasy
 * pustego stanu. Faza siedzi w `data-faza`, żeby arkusz mógł odróżnić błąd od
 * pustki bez mnożenia klas, a sprawdzian mógł zapytać o stan okna, nie o jego
 * wygląd.
 *
 * Forma jest ta sama co w pozostałych modułach: ikona, tytuł, opis, treść
 * wyśrodkowana. Stopnie pisma i szerokość łamania niesie wspólny
 * `komponenty/drobne.css` (`.dn-pusty-stan-tytul`, `.dn-pusty-stan-opis`) —
 * moduł ich u siebie nie nadpisuje, bo dwie definicje tej samej formy to dwie
 * formy. Ikona należy do stanu pustego i tylko do niego (arkusz modułu chowa ją
 * w pozostałych fazach): ładowanie ma własny nośnik `.dn-spinner`, błąd —
 * kreskę po lewej.
 *
 * Komunikat przesłania treść, nie kasuje jej: nieudane odświeżenie zostawia to,
 * co już było widoczne, więc powrót do treści nie wymaga ponownego odczytu.
 *
 * Ładowanie nie chowa treści i jest to decyzja modułu. Treścią każdego z trzech
 * okien jest formularz, do którego Operator właśnie pisze (tekst źródłowy, kod
 * języka, termin glosariusza) — schowanie go na czas zapisu zabrałoby z oczu
 * to, co wpisane. Kontrolka zostaje więc widoczna i edytowalna, a wskaźnik
 * odczytu stoi obok niej, nigdy zamiast niej. Wspólny `komponenty/faza-okna`
 * tej decyzji nie przesądza — moduły różnią się tu świadomie.
 *
 * Bez własnego „Spróbuj ponownie": wszystkie trzy okna stoją na zapisach
 * wyzwalanych z formularza (zapis źródła, dodanie języka, zapis terminu), więc
 * ponowieniem jest ten sam przycisk, który czynność wywołał — zostaje klikalny
 * także po odmowie, a drugi przycisk o tym samym skutku byłby dwiema drogami
 * do jednej czynności. Jedyny odczyt modułu, kontekst okna, ma ponowienie
 * w pasie kontekstu.
 *
 * Podział z bytem wspólnym: zestaw faz i znakowanie powłoki (`data-faza`, rola
 * pasa, jego widoczność) przychodzą z `komponenty/faza-okna`. Tutaj zostaje to,
 * czego wspólny byt nie przesądza: klasy `mt-…`, treść komunikatu, nośnik ikony
 * i spinnera oraz pozostawienie treści widocznej na czas ładowania.
 */

import { elementIkony } from '../../ikony/ikony';
import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/** Zdanie stanu pustego: tytuł nad opisem. */
export interface ZdanieStanu {
  /** Tytuł — nazywa sytuację, w której okno stoi. */
  readonly tytul: string;
  /** Opis — mówi, czym okno jest i jak je zapełnić. */
  readonly opis: string;
}

/** Powłoka okna wraz z pasem stanu. */
export interface StanOkna {
  /** Element osadzany w oknie; niesie pas stanu i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pokazuje pustkę merytoryczną — wywołanie się udało, wyniku nie ma. */
  puste(zdanie: ZdanieStanu): void;
  /** Pokazuje odmowę albo awarię wraz z jej powodem. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

/** Tytuł pasa w fazach, które nie są pustką — jedno brzmienie na trzy okna. */
const TYTUL_LADOWANIA = 'Rdzeń pracuje';
const TYTUL_BLEDU = 'Rdzeń odmówił';

export function utworzStanOkna(poczatkowe: ZdanieStanu): StanOkna {
  let biezaca: FazaOkna = 'puste';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul mt-stan__tytul';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis mt-stan__opis';

  const spinner = document.createElement('span');
  spinner.className = 'dn-spinner mt-stan__spinner';
  spinner.setAttribute('role', 'status');
  spinner.setAttribute('aria-label', 'Trwa wywołanie rdzenia');
  spinner.hidden = true;
  tytul.append(spinner);

  const ikona = elementIkony('tlumacz', { rozmiar: 24, klasa: 'dn-ikona mt-stan__ikona' });

  const pas = document.createElement('div');
  pas.className = 'dn-pusty-stan mt-stan';
  pas.append(ikona, tytul, komunikat);

  const tresc = document.createElement('div');
  tresc.className = 'mt-stan__tresc';

  const element = document.createElement('div');
  element.className = 'mt-stan__powloka';
  element.append(pas, tresc);

  function ustaw(faza: FazaOkna, naglowek: string, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, pas, faza);
    tytul.replaceChildren(naglowek, spinner);
    komunikat.textContent = opis;
    spinner.hidden = faza !== 'ladowanie';
  }

  ustaw('puste', poczatkowe.tytul, poczatkowe.opis);

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', TYTUL_LADOWANIA, opis),
    puste: (zdanie) => ustaw('puste', zdanie.tytul, zdanie.opis),
    blad: (opis) => ustaw('blad', TYTUL_BLEDU, opis),
    gotowe: () => ustaw('gotowe', '', ''),
    faza: () => biezaca,
  };
}
