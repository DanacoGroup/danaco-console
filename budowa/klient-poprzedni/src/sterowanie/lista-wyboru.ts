import { utworzNaglowekSterowania } from './naglowek-sterowania';
import { utworzWskaznikOdczytu } from './wskaznik-odczytu';

/** Pojedyncza pozycja listy wyboru. */
export interface OpcjaWyboru {
  wartosc: string;
  nazwa: string;
  /**
   * Nazwa sekcji, w której pozycja ma stanąć; pominięta znaczy „poza sekcjami".
   *
   * Sekcje służą wyborowi modelu: modele surowe i agenci stoją w jednym menu,
   * ale w dwóch sekcjach — „Modele" i „Moi agenci". Wybór pozostaje jeden
   * (okno obsługuje agent), a rodzaj pozycji jest widoczny: agent stoi na
   * modelu, więc płaska lista zatarłaby tę zależność, a przy licznych agentach
   * modele surowe utonęłyby między nimi.
   *
   * Sekcje rysują się natywnym `optgroup`, nie własnym menu: czytnik ekranu
   * czyta wtedy nazwę sekcji przy pozycji bez dodatkowych atrybutów `aria-*`.
   */
  sekcja?: string;
}

/**
 * Lista wyboru — prymityw wspólny sterowaniom wybierającym jedną wartość.
 *
 * Lista nigdy nie dostaje atrybutu `disabled`, także wtedy, gdy katalog
 * wartości jest pusty albo jeszcze nie przyszedł z rdzenia. Wartość bieżąca
 * okna pojawia się na liście również wtedy, gdy nie ma jej w katalogu —
 * interfejs pokazuje stan okna, zamiast podmieniać go na pozycję pierwszą.
 *
 * Znak [?] obok etykiety niesie zdanie z wykazu adnotacji, w tym zdanie o tym,
 * że sterowanie zapisuje wartość, której wykonanie jeszcze nie czyta.
 *
 * Wskaźnik odczytu stoi obok pola, nigdy zamiast niego: lista pozostaje
 * klikalna i przyjmuje wybór, dopóki katalog jedzie z rdzenia. Etykieta jest
 * osobnym `<label for>`, a nie obudową pola, bo znak [?] jest przyciskiem —
 * wewnątrz obudowy jego naciśnięcie przenosiłoby się na listę.
 */
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

/** Identyfikator wiążący etykietę z listą; unikalny w obrębie dokumentu. */
let licznik = 0;

/**
 * Tworzy listę wyboru z opcjonalnym polem zawężania nad listą.
 *
 * Zawężanie stoi nad listą, a nie w jej rozwinięciu, bo rozwinięcie natywnego
 * `<select>` rysuje system operacyjny, nie dokument — pola tekstowego nie da się
 * tam wstawić, a filtr wewnątrz rozwinięcia wymagałby zastąpienia `<select>`
 * własnym bytem rozwijanym. Pole nad listą przycina katalog pozycji: wpisana
 * fraza zostawia pozycje pasujące, sekcje zwężają się razem z nimi. Ceną jest
 * ruch dwutaktowy — najpierw fraza, potem rozwinięcie.
 *
 * Fraza bez trafień nie kasuje wyboru ani nie blokuje listy: wartość bieżąca
 * okna zostaje na liście zawsze (`zKatalogiem`), a pod polem staje zdanie
 * o braku trafień.
 */
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

  // Ostatni katalog i wybór trzymane po to, by przepisanie frazy odrysowało
  // listę bez pytania rejestru o cokolwiek — zawężanie jest czynnością widoku.
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

/** Katalog wartości poszerzony o wartość bieżącą, jeżeli jej w nim nie ma. */
function zKatalogiem(opcje: OpcjaWyboru[], wybrana: string): OpcjaWyboru[] {
  if (opcje.some((opcja) => opcja.wartosc === wybrana)) return opcje;
  return [{ wartosc: wybrana, nazwa: nazwaWlasna(wybrana) }, ...opcje];
}

/** Nazwa wartości spoza katalogu; pusta wartość znaczy brak wskazania. */
function nazwaWlasna(wartosc: string): string {
  return wartosc.length > 0 ? `${wartosc} (spoza wykazu)` : 'Bez wskazania';
}

/**
 * Węzły listy: pozycje bez sekcji wprost, pozycje z sekcją w `optgroup`.
 *
 * Kolejność sekcji wynika z kolejności pierwszego wystąpienia, nie z sortowania
 * po nazwie: wołający układa katalog w wybranym przez siebie porządku (modele
 * surowe przed agentami), a lista tego porządku nie przestawia.
 *
 * Sekcja pusta nie powstaje — `optgroup` bez pozycji byłby nagłówkiem nad
 * niczym.
 */
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

/** Pojedyncza pozycja listy. */
function pozycja(opcja: OpcjaWyboru): HTMLOptionElement {
  const element = document.createElement('option');
  element.value = opcja.wartosc;
  element.textContent = opcja.nazwa;
  return element;
}
