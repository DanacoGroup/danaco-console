import { elementIkony } from '../ikony/ikony';
import type { PozycjaModuluStrony } from './pozycje-modulow';
import { utworzStrefeZwijana } from './strefa-zwijana';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/** Strefa kafli modułów, które nie mają pozycji w bocznej nawigacji, rysuje siatkę i zgłasza wybór, chowając się całkowicie, gdy wykaz z rdzenia jest pusty. */
const ETYKIETA = 'Moduły poza nawigacją środowisk';
const WYJASNIENIE =
  'Rdzeń nie pokazuje tych modułów w bocznej nawigacji żadnego środowiska — ' +
  'jedyna droga do ich okien prowadzi stąd.';

/** Ikona kafla w skali średniej zestawu, tak samo jak w strefie komponentów własnych, ustalona jednym stałym rozmiarem. */
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

  // Strefa jest zwinięta domyślnie: wejście zapasowe do modułów spoza nawigacji.
  const strefa = utworzStrefeZwijana({
    etykieta: ETYKIETA,
    wyjasnienie: WYJASNIENIE,
    klucz: 'strona.moduly',
    domyslnieRozwiniete: false,
  });
  const element = strefa.element;
  element.classList.add('dn-strona__strefa--moduly');
  element.setAttribute('aria-label', ETYKIETA);
  // Strefa startuje ukryta: przed odpowiedzią rdzenia nie wiadomo, czy moduł jest poza nawigacją.
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

/** Jeden kafel modułu; wskazanie wyboru oddaje strefie, sam kafel nie zna jej bieżącego stanu ani wykazu. */
function kafel(
  pozycja: PozycjaModuluStrony,
  przyWyborze: (pozycja: PozycjaModuluStrony) => void,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-karta dn-karta--klikalna dn-karta--komponent dn-strona__kafel';
  // Znacznik własny, nie wspólny ze sceną modułu, żeby kafel dał się odróżnić od widoku docelowego.
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

/** Podaje, ile okien moduł otwiera, liczbą z katalogu rdzenia, nie z widoku, mówiąc operatorowi, co za kaflem stoi. */
function opisOkien(ile: number): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dn-strona__kafel-metadane';
  element.textContent = `${ile} ${odmianaOkien(ile)}`;
  return element;
}

/** Odmienia liczebnik okien operacyjnych według reguł liczby mnogiej polszczyzny, dla wartości od zera wzwyż. */
function odmianaOkien(ile: number): string {
  if (ile === 1) return 'okno operacyjne';
  const setki = ile % 100;
  const jednosci = ile % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (setki < 12 || setki > 14);
  return mnoga ? 'okna operacyjne' : 'okien operacyjnych';
}
