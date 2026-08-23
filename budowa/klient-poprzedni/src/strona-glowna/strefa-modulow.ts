import { elementIkony } from '../ikony/ikony';
import type { PozycjaModuluStrony } from './pozycje-modulow';
import { utworzStrefeZwijana } from './strefa-zwijana';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/**
 * Strefa kafli modułów, które nie mają pozycji w bocznej nawigacji.
 *
 * Rysuje siatkę kafli i zgłasza wybór. Strefa nie pyta rdzenia i nie otwiera
 * modułu — wykaz podaje jej `wpiecie-modulow`, a skutek wyboru należy do
 * warstwy, która stronę zamontowała.
 *
 * Wykaz pusty chowa całą strefę, zamiast zostawiać nagłówek nad pustym
 * prostokątem: gdy rdzeń przypnie moduły do środowisk, kafle znikną razem ze
 * swoim powodem.
 *
 * Waga wizualna jak w strefie drugiej — ta sama karta, ten sam krój bazowy
 * półgruby, ta sama siatka. Kafel prowadzi do pracy w module, a nie do
 * zbudowania rzeczy, więc tak jak kafel komponentu nie nosi akcentu
 * zarezerwowanego dla kart środowisk.
 */

const ETYKIETA = 'Moduły poza nawigacją środowisk';
const WYJASNIENIE =
  'Rdzeń nie pokazuje tych modułów w bocznej nawigacji żadnego środowiska — ' +
  'jedyna droga do ich okien prowadzi stąd.';

/** Ikona kafla w skali średniej zestawu — tak jak w strefie komponentów. */
const ROZMIAR_IKONY = 18;

export interface StrefaModulow {
  /** Element `<details>` — strefa przywoływana, nie rysowana z urzędu. */
  element: HTMLDetailsElement;
  naWybor(sluchacz: SluchaczWyboru<PozycjaModuluStrony>): void;
  /** Ustawia kafle; wykaz pusty chowa całą strefę. */
  ustaw(pozycje: readonly PozycjaModuluStrony[]): void;
}

export function utworzStrefeModulow(): StrefaModulow {
  const sygnal = utworzSygnalWyboru<PozycjaModuluStrony>();

  // Strefa jest zwinięta domyślnie: to wejście zapasowe do modułów, których
  // rdzeń nie postawił w żadnej nawigacji — potrzebne, ale nie codzienne.
  // Zapowiedź z liczbą modułów stoi zawsze, więc zwinięcie niczego nie ukrywa
  // przed Operatorem.
  const strefa = utworzStrefeZwijana({
    etykieta: ETYKIETA,
    wyjasnienie: WYJASNIENIE,
    klucz: 'strona.moduly',
    domyslnieRozwiniete: false,
  });
  const element = strefa.element;
  element.classList.add('dn-strona__strefa--moduly');
  element.setAttribute('aria-label', ETYKIETA);
  // Strefa startuje ukryta: przed odpowiedzią rdzenia nie wiadomo, czy
  // jakikolwiek moduł jest poza nawigacją, a zapowiedź nad pustką mówiłaby
  // o wykazie, którego nikt jeszcze nie odczytał.
  element.hidden = true;

  const siatka = document.createElement('div');
  siatka.className = 'dn-strona__siatka dn-strona__siatka--moduly';
  strefa.cialo.append(siatka);

  return {
    element,
    naWybor: sygnal.sluchaj,

    ustaw(pozycje) {
      siatka.replaceChildren(...pozycje.map((p) => kafel(p, sygnal.nadaj)));
      element.hidden = pozycje.length === 0;
      strefa.ustawDopisek(pozycje.length === 0 ? '' : `${pozycje.length}`);
    },
  };
}

/** Jeden kafel modułu. Wskazanie oddaje strefie — sam nie zna jej stanu. */
function kafel(
  pozycja: PozycjaModuluStrony,
  przyWyborze: (pozycja: PozycjaModuluStrony) => void,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-karta dn-karta--klikalna dn-karta--komponent dn-strona__kafel';
  // Znacznik własny, nie `data-modul`: ten atrybut nosi już scena modułu
  // (`moduly/terminal/indeks.ts`, `moduly/multitasking/indeks.ts`), więc kafel
  // pod tą samą nazwą byłby nie do odróżnienia od widoku, do którego prowadzi.
  element.dataset['kafelModulu'] = pozycja.kod;
  element.title = pozycja.opis === '' ? pozycja.nazwa : `${pozycja.nazwa} — ${pozycja.opis}`;

  const tresc = document.createElement('span');
  tresc.className = 'dn-karta-tresc dn-strona__kafel-tresc';

  const ikona = document.createElement('span');
  ikona.className = 'dn-strona__ikona-kafla';
  ikona.append(elementIkony(pozycja.ikona, { rozmiar: ROZMIAR_IKONY }));

  const napisy = document.createElement('span');
  napisy.className = 'dn-strona__kafel-napisy';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-karta-tytul dn-strona__kafel-nazwa';
  nazwa.textContent = pozycja.nazwa;

  const wezwanie = document.createElement('span');
  wezwanie.className = 'dn-strona__kafel-wezwanie';
  wezwanie.textContent = pozycja.wezwanie;

  napisy.append(nazwa, wezwanie, opisOkien(pozycja.okien));
  tresc.append(ikona, napisy);
  element.append(tresc);

  element.addEventListener('click', () => przyWyborze(pozycja));
  return element;
}

/**
 * Ile okien moduł otwiera — liczba z katalogu rdzenia, nie z widoku.
 *
 * Mówi Operatorowi, co za kaflem stoi, zanim w niego wejdzie. Zero jest
 * stanem możliwym i mówi się je wprost, zamiast chować wiersz.
 */
function opisOkien(ile: number): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dn-strona__kafel-metadane';
  element.textContent = `${ile} ${odmianaOkien(ile)}`;
  return element;
}

/** Odmiana po liczebniku: „1 okno operacyjne", „4 okna", „6 okien". */
function odmianaOkien(ile: number): string {
  if (ile === 1) return 'okno operacyjne';
  const setki = ile % 100;
  const jednosci = ile % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (setki < 12 || setki > 14);
  return mnoga ? 'okna operacyjne' : 'okien operacyjnych';
}
