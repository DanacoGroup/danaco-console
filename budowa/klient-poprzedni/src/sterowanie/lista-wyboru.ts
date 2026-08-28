import { utworzNaglowekSterowania } from './naglowek-sterowania';
import { utworzWskaznikOdczytu } from './wskaznik-odczytu';

/** Pojedyncza pozycja listy wyboru: wartość, nazwa widoczna i opcjonalna nazwa sekcji, w której pozycja stoi. */
export interface OpcjaWyboru {
  wartosc: string;
  nazwa: string;
  /**
   * Nazwa sekcji pozycji w menu; pominięta znaczy pozycję poza sekcjami.
   */
  sekcja?: string;
}

/** Lista wyboru: prymityw wspólny sterowaniom wybierającym jedną wartość, z wartością bieżącą zawsze widoczną na liście. */
export interface ListaWyboru {
  /** Element montowany w kompletach sterowania. */
  element: HTMLElement;
  /** Odrysowuje pozycje i zaznacza wartość bieżącą. */
  pokaz(opcje: OpcjaWyboru[], wybrana: string): void;
  /** Zapala albo gasi wskaźnik odczytu katalogu obok pola. */
  ustawOdczyt(trwa: boolean): void;
  /** Zdanie pod polem; puste je zdejmuje. Nośnik stanu pustego katalogu. */
  ustawUwage(tresc: string): void;
}

/** Identyfikator wiążący etykietę z polem listy; liczony rosnąco, więc pozostaje unikalny w obrębie dokumentu. */
let licznik = 0;

/** Tworzy listę wyboru z opcjonalnym polem zawężania nad listą, przycinającym katalog pozycji do frazy wpisanej przez operatora. */
export function utworzListeWyboru(
  etykieta: string,
  przyWyborze: (wartosc: string) => void,
  zawezanie = false,
): ListaWyboru {
  licznik += 1;
  const identyfikator = `dc-ster-lista-${licznik}`;

  const element = document.createElement('div');
  element.className = 'dc-ster-pole';

  const wskaznik = utworzWskaznikOdczytu(`Odczytuję katalog: ${etykieta}`);
  const naglowek = utworzNaglowekSterowania(etykieta, identyfikator, wskaznik.element);

  const lista = document.createElement('select');
  lista.id = identyfikator;
  lista.className = 'dc-ster-pole__lista';
  lista.addEventListener('change', () => przyWyborze(lista.value));

  const uwaga = document.createElement('p');
  uwaga.className = 'dc-ster-pole__uwaga';
  uwaga.hidden = true;

  const fraza = document.createElement('input');
  fraza.type = 'search';
  fraza.className = 'dc-ster-pole__wpis dc-ster-pole__szukanie';
  fraza.placeholder = 'Zawęź wykaz';
  fraza.setAttribute('aria-label', `Zawężanie wykazu: ${etykieta}`);
  fraza.hidden = !zawezanie;

  const brakTrafien = document.createElement('p');
  brakTrafien.className = 'dc-ster-pole__uwaga';
  brakTrafien.hidden = true;

  element.append(naglowek, fraza, lista, brakTrafien, uwaga);

  // Ostatni katalog i wybór trzymane, by przepisanie frazy odrysowało listę bez pytania rejestru.
  let katalog: OpcjaWyboru[] = [];
  let wybor = '';

  function odrysuj(): void {
    const szukane = fraza.value.trim().toLocaleLowerCase('pl-PL');
    const zawezone = szukane === '' ? katalog : katalog.filter((opcja) => pasuje(opcja, szukane));
    lista.replaceChildren(...wezly(zKatalogiem(zawezone, wybor)));
    lista.value = wybor;
    brakTrafien.textContent =
      szukane !== '' && zawezone.length === 0
        ? `Fraza „${fraza.value.trim()}" nie pasuje do żadnej pozycji wykazu. Wybór bieżący zostaje na liście.`
        : '';
    brakTrafien.hidden = brakTrafien.textContent === '';
  }

  fraza.addEventListener('input', odrysuj);

  return {
    element,

    pokaz(opcje, wybrana) {
      katalog = opcje;
      wybor = wybrana;
      odrysuj();
    },

    ustawOdczyt: wskaznik.ustaw,

    ustawUwage(tresc) {
      uwaga.textContent = tresc;
      uwaga.hidden = tresc === '';
    },
  };
}

/**
 * Trafienie frazy: nazwa pozycji albo nazwa jej sekcji.
 *
 * Sekcja liczy się razem z nazwą, bo „agent" i „model" są frazami wyszukiwania
 * tak samo użytecznymi jak nazwa własna pozycji — nazwy sekcji są widoczne
 * w wykazie.
 */
function pasuje(opcja: OpcjaWyboru, szukane: string): boolean {
  const tekst = `${opcja.nazwa} ${opcja.sekcja ?? ''}`.toLocaleLowerCase('pl-PL');
  return tekst.includes(szukane);
}

/** Katalog wartości poszerzony o wartość bieżącą okna, jeżeli katalog przysłany z rejestru jej nie zawiera. */
function zKatalogiem(opcje: OpcjaWyboru[], wybrana: string): OpcjaWyboru[] {
  if (opcje.some((opcja) => opcja.wartosc === wybrana)) return opcje;
  return [{ wartosc: wybrana, nazwa: nazwaWlasna(wybrana) }, ...opcje];
}

/** Nazwa wartości spoza katalogu pokazywana na liście; pusta wartość znaczy brak wskazania w tym oknie. */
function nazwaWlasna(wartosc: string): string {
  return wartosc.length > 0 ? `${wartosc} (spoza wykazu)` : 'Bez wskazania';
}

/** Węzły listy: pozycje bez sekcji wprost, pozycje z sekcją zebrane w grupę, w kolejności pierwszego wystąpienia sekcji. */
function wezly(opcje: OpcjaWyboru[]): Node[] {
  const wynik: Node[] = [];
  const grupy = new Map<string, HTMLOptGroupElement>();

  for (const opcja of opcje) {
    const nazwaSekcji = opcja.sekcja ?? '';
    if (nazwaSekcji === '') {
      wynik.push(pozycja(opcja));
      continue;
    }
    let grupa = grupy.get(nazwaSekcji);
    if (grupa === undefined) {
      grupa = document.createElement('optgroup');
      grupa.label = nazwaSekcji;
      grupy.set(nazwaSekcji, grupa);
      wynik.push(grupa);
    }
    grupa.append(pozycja(opcja));
  }

  return wynik;
}

/** Pojedyncza pozycja listy jako element wyboru, z wartością i nazwą przepisanymi wprost z opcji katalogu. */
function pozycja(opcja: OpcjaWyboru): HTMLOptionElement {
  const element = document.createElement('option');
  element.value = opcja.wartosc;
  element.textContent = opcja.nazwa;
  return element;
}
