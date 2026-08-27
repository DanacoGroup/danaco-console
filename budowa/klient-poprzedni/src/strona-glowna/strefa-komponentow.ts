import { utworzKafelKomponentu } from './kafel-komponentu';
import { POZYCJE_KOMPONENTOW, type PozycjaKomponentu } from './pozycje-komponentow';
import { utworzStrefeZwijana } from './strefa-zwijana';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/** Strefa druga pokazuje siatkę kafli komponentów: kafli modułów nastawianych na stronie głównej i kafli komponentów zbudowanych oraz nazwanych przez operatora. */
const ETYKIETA = 'Strefa 2 · Kafle komponentów własnych';
const WYJASNIENIE =
  'Kliknięcie kafla nie otwiera przestrzeni roboczej — otwiera okno konfiguracji, w którym powstaje nazwany wytwór.';

export interface StrefaKomponentow {
  /** Element `<details>` — strefa przywoływana, nie rysowana z urzędu. */
  element: HTMLDetailsElement;
  /** Miejsce na formularz zakładania — wpina je warstwa, która ma kanał. */
  przybornik: HTMLElement;
  naWybor(sluchacz: SluchaczWyboru<PozycjaKomponentu>): void;
  /** Ustawia kafle modułów strony głównej; wykaz pusty nie zdejmuje kafli, zostaje czwórka zastana. */
  ustawRodzaje(pozycje: readonly PozycjaKomponentu[]): void;
  /** Ustawia kafle komponentów zbudowanych przez operatora; wykaz pusty jest stanem poprawnym. */
  ustawPersonalizowane(pozycje: readonly PozycjaKomponentu[]): void;
}

export function utworzStrefeKomponentow(): StrefaKomponentow {
  const sygnal = utworzSygnalWyboru<PozycjaKomponentu>();

  // Strefa zwijalna, rozwinięta domyślnie: to jedna z dróg wejścia w pracę.
  const strefa = utworzStrefeZwijana({
    etykieta: ETYKIETA,
    wyjasnienie: WYJASNIENIE,
    klucz: 'strona.komponenty',
    domyslnieRozwiniete: true,
  });
  const element = strefa.element;
  element.classList.add('dn-strona__strefa--komponenty');
  element.setAttribute('aria-label', ETYKIETA);

  // Jedna siatka na oba rodzaje kafli; rodzaj rozpoznaje się po treści, nie po osobnej sekcji.
  const siatka = document.createElement('div');
  siatka.className = 'dn-strona__siatka dn-strona__siatka--komponenty';

  const przybornik = document.createElement('div');
  przybornik.className = 'dn-strona__przybornik';

  function kafel(pozycja: PozycjaKomponentu): HTMLElement {
    return utworzKafelKomponentu(pozycja, (wybrany) => sygnal.nadaj(wybrany)).element;
  }

  // Czwórka zastana stoi od pierwszej klatki, żeby strefa nie była pusta przed odpowiedzią rdzenia.
  let rodzaje: readonly PozycjaKomponentu[] = POZYCJE_KOMPONENTOW;
  let personalizowane: readonly PozycjaKomponentu[] = [];

  function przerysuj(): void {
    siatka.replaceChildren(...rodzaje.map(kafel), ...personalizowane.map(kafel));
    strefa.ustawDopisek(String(rodzaje.length + personalizowane.length));
  }

  przerysuj();

  // Przybornik stoi pod siatką, nie w wierszu nagłówka, żeby nie konkurował z etykietą strefy o uwagę.
  strefa.cialo.append(siatka, przybornik);

  return {
    element,
    przybornik,
    naWybor: sygnal.sluchaj,

    ustawRodzaje(pozycje) {
      if (pozycje.length === 0) return;
      rodzaje = pozycje;
      przerysuj();
    },

    ustawPersonalizowane(pozycje) {
      personalizowane = pozycje;
      przerysuj();
    },
  };
}
