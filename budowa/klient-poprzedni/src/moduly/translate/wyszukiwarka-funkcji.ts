import { poleTekstowe } from '../../modele/kontrolki-formularza';
import {
  KATALOG_FUNKCJI,
  grupyKatalogu,
  liczbaZKomenda,
  type PozycjaKatalogu,
} from './katalog-funkcji-translate';
import { WARSTWY, utworzRozwiniecie, type Rozwiniecie } from './warstwy-translate';

/**
 * Wyszukiwarka funkcji modułu — droga warstwy czwartej do każdej pozycji
 * katalogu.
 *
 * Zasada jednego kliknięcia mówi, że każdy element modułu jest osiągalny jednym
 * kliknięciem, jednym skrótem albo jednym poleceniem. Wyszukiwarka jest tą
 * drogą dla pozycji, które nie mają własnego przycisku w oknie: nazywa je,
 * mówi, w którym oknie stoją i na której warstwie, oraz podaje przy każdej
 * jedno z dwojga — komendę kontraktu albo brak.
 *
 * Wyszukiwarka niczego nie uruchamia i nie udaje, że uruchamia. Pozycja bez
 * komendy nie dostaje tu przycisku, który po naciśnięciu przeprosi — dostaje
 * zdanie o tym, czego brakuje, czytelne bez naciskania czegokolwiek.
 *
 * Szukanie idzie środkiem nazwy, opisu i grupy, bo nazwy pozycji są w większości
 * angielskie i złożone: szukanie wyłącznie od początku nazwy nie znalazłoby
 * „Placeholder Style Mapping" po słowie „placeholder" wpisanym z pamięci.
 */
export interface WyszukiwarkaFunkcji {
  /** Element osadzany w module. */
  element: HTMLElement;
  /** Rozwija wyszukiwarkę i prowadzi do niej ognisko — wołane skrótem. */
  otworz(): void;
}

export function utworzWyszukiwarkeFunkcji(): WyszukiwarkaFunkcji {
  const rozwiniecie: Rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Wyszukiwarka funkcji modułu',
    wyjasnienie:
      `Katalog liczy ${String(KATALOG_FUNKCJI.length)} pozycji z opracowania modułu, ` +
      `z czego ${String(liczbaZKomenda())} ma komendę kontraktu.`,
    znacznik: '☰',
  });

  const pytanie = poleTekstowe({
    etykieta: 'Szukaj w katalogu funkcji',
    podpowiedz: 'nazwa funkcji, grupa albo słowo z opisu',
  });

  const grupa = document.createElement('select');
  grupa.className = 'dn-pole-kontrolka mt-katalog__grupa';
  grupa.setAttribute('aria-label', 'Zawężenie katalogu do grupy');
  grupa.replaceChildren(
    pozycjaWyboru('', 'wszystkie grupy'),
    ...grupyKatalogu().map((nazwa) => pozycjaWyboru(nazwa, nazwa)),
  );

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'mt-katalog__podsumowanie';

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-katalog';

  const filtry = document.createElement('div');
  filtry.className = 'mt-katalog__filtry';
  filtry.append(pytanie.element, grupa);

  rozwiniecie.tresc.append(filtry, podsumowanie, wykaz);

  pytanie.kontrolka.addEventListener('input', odswiez);
  grupa.addEventListener('change', odswiez);

  function odswiez(): void {
    const szukane = pytanie.kontrolka.value.trim().toLowerCase();
    const wybranaGrupa = grupa.value;
    const znalezione = KATALOG_FUNKCJI.filter(
      (pozycja) =>
        (wybranaGrupa === '' || pozycja.grupa === wybranaGrupa) && pasuje(pozycja, szukane),
    );
    podsumowanie.textContent = zdaniePodsumowania(znalezione.length, szukane, wybranaGrupa);
    wykaz.replaceChildren(...znalezione.map(wierszPozycji));
  }

  odswiez();

  return {
    element: rozwiniecie.element,
    otworz() {
      rozwiniecie.rozwin();
      pytanie.kontrolka.focus();
      pytanie.kontrolka.select();
    },
  };
}

function pozycjaWyboru(wartosc: string, etykieta: string): HTMLOptionElement {
  const element = document.createElement('option');
  element.value = wartosc;
  element.textContent = etykieta;
  return element;
}

function pasuje(pozycja: PozycjaKatalogu, szukane: string): boolean {
  if (szukane === '') return true;
  const stog = `${pozycja.nazwa} ${pozycja.grupa} ${pozycja.opis} ${pozycja.okno}`.toLowerCase();
  return stog.includes(szukane);
}

function zdaniePodsumowania(ile: number, szukane: string, wybranaGrupa: string): string {
  const zakres =
    szukane === '' && wybranaGrupa === ''
      ? 'cały katalog'
      : `zawężenie${szukane === '' ? '' : ` „${szukane}"`}${wybranaGrupa === '' ? '' : `, grupa „${wybranaGrupa}"`}`;
  if (ile === 0) {
    return `Katalog nie ma pozycji spełniającej to zawężenie (${zakres}).`;
  }
  return `Pozycji w wykazie: ${String(ile)} z ${String(KATALOG_FUNKCJI.length)} (${zakres}).`;
}

/**
 * Jeden wiersz katalogu.
 *
 * Zdanie o pokryciu buduje się z pól pozycji, nie z osobnego napisu przy każdej
 * z nich: pozycja z komendami wymienia je co do nazwy, pozycja bez komend mówi,
 * czego brakuje. Dzięki temu dopisanie komendy do pozycji zmienia zdanie samo.
 */
function wierszPozycji(pozycja: PozycjaKatalogu): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mt-katalog__nazwa';
  nazwa.textContent = pozycja.nazwa;

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka mt-katalog__stan';
  stan.textContent = pozycja.komendy.length === 0 ? 'bez komendy' : 'komenda kontraktu';

  const umiejscowienie = document.createElement('p');
  umiejscowienie.className = 'mt-katalog__umiejscowienie';
  umiejscowienie.textContent =
    `${pozycja.okno} · ${pozycja.grupa} · warstwa ${String(pozycja.warstwa)} — ` +
    WARSTWY[pozycja.warstwa].nazwa;

  const opis = document.createElement('p');
  opis.className = 'mt-katalog__opis';
  opis.textContent = pozycja.opis;

  const element = document.createElement('li');
  element.className = 'mt-katalog__wiersz';
  element.dataset['pokrycie'] = pozycja.komendy.length === 0 ? 'brak' : 'komenda';
  element.dataset['warstwa'] = String(pozycja.warstwa);

  const naglowek = document.createElement('div');
  naglowek.className = 'mt-katalog__naglowek';
  naglowek.append(nazwa, stan);

  element.append(naglowek, umiejscowienie, opis);

  if (pozycja.komendy.length > 0) {
    const komendy = document.createElement('p');
    komendy.className = 'mt-katalog__komendy';
    komendy.textContent = `Wykonuje: ${pozycja.komendy.join(' · ')}`;
    element.append(komendy);
  }
  if (pozycja.brak !== undefined) {
    const brak = document.createElement('p');
    brak.className = 'mt-katalog__brak';
    brak.textContent = pozycja.brak;
    element.append(brak);
  }
  return element;
}
