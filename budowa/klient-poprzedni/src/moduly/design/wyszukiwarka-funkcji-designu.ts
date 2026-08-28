import { poleTekstowe } from '../../modele/kontrolki-formularza';
import {
  KATALOG_FUNKCJI,
  grupyKatalogu,
  liczbaWedlugWykonania,
  type PozycjaKatalogu,
} from './katalog-funkcji-designu';
import { WARSTWY, utworzRozwiniecie, type Rozwiniecie } from './warstwy-designu';

/**
 * Wyszukiwarka funkcji modułu prowadzi do każdej pozycji katalogu warstwy
 * czwartej i podaje przy niej komendę kontraktu, czynność okna albo
 * informację o braku drogi.
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
    wyjasnienie: zdanieZasiegu(),
    znacznik: '☰',
  });

  const pytanie = poleTekstowe({
    etykieta: 'Szukaj w katalogu funkcji',
    podpowiedz: 'nazwa funkcji, grupa albo słowo z opisu',
  });

  const grupa = document.createElement('select');
  grupa.className = 'dn-pole-kontrolka md-katalog__grupa';
  grupa.setAttribute('aria-label', 'Zawężenie katalogu do grupy');
  grupa.replaceChildren(
    pozycjaWyboru('', 'wszystkie grupy'),
    ...grupyKatalogu().map((nazwa) => pozycjaWyboru(nazwa, nazwa)),
  );

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'md-katalog__podsumowanie';

  const wykaz = document.createElement('ul');
  wykaz.className = 'md-katalog';

  const filtry = document.createElement('div');
  filtry.className = 'md-katalog__filtry';
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

/**
 * Zdanie o zasięgu katalogu liczy trzy wartości z wykazu pozycji: łączną
 * liczbę, liczbę wykonywaną komendą kontraktu oraz liczbę wykonywaną samym
 * oknem bez rdzenia.
 */
function zdanieZasiegu(): string {
  return (
    `Katalog liczy ${String(KATALOG_FUNKCJI.length)} pozycji z opracowania modułu: ` +
    `${String(liczbaWedlugWykonania('komenda'))} wykonuje komenda kontraktu, ` +
    `${String(liczbaWedlugWykonania('okno'))} wykonuje samo okno bez rdzenia, ` +
    `${String(liczbaWedlugWykonania('bez-drogi'))} nie ma drogi ani tu, ani tam.`
  );
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
      : `zawężenie${szukane === '' ? '' : ` „${szukane}"`}` +
        `${wybranaGrupa === '' ? '' : `, grupa „${wybranaGrupa}"`}`;
  if (ile === 0) {
    return `Katalog nie ma pozycji spełniającej to zawężenie (${zakres}).`;
  }
  return `Pozycji w wykazie: ${String(ile)} z ${String(KATALOG_FUNKCJI.length)} (${zakres}).`;
}

/** Nazwa stanu wykonania pozycji katalogu widoczna dla operatora, jedno brzmienie wspólne dla całego modułu wyszukiwarki funkcji. */
const NAZWY_WYKONANIA: Readonly<Record<PozycjaKatalogu['wykonanie'], string>> = {
  komenda: 'komenda kontraktu',
  okno: 'czynność okna',
  'bez-drogi': 'bez drogi',
};

/**
 * Jeden wiersz katalogu funkcji: zdanie o pokryciu buduje się z pól pozycji,
 * więc dopisanie komendy do pozycji zmienia treść wiersza samo.
 */
function wierszPozycji(pozycja: PozycjaKatalogu): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'md-katalog__nazwa';
  nazwa.textContent = pozycja.nazwa;

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka md-katalog__stan';
  stan.textContent = NAZWY_WYKONANIA[pozycja.wykonanie];

  const umiejscowienie = document.createElement('p');
  umiejscowienie.className = 'md-katalog__umiejscowienie';
  umiejscowienie.textContent =
    `${pozycja.okno} · ${pozycja.grupa} · warstwa ${String(pozycja.warstwa)} — ` +
    WARSTWY[pozycja.warstwa].nazwa;

  const opis = document.createElement('p');
  opis.className = 'md-katalog__opis';
  opis.textContent = pozycja.opis;

  const naglowek = document.createElement('div');
  naglowek.className = 'md-katalog__naglowek';
  naglowek.append(nazwa, stan);

  const element = document.createElement('li');
  element.className = 'md-katalog__wiersz';
  element.dataset['wykonanie'] = pozycja.wykonanie;
  element.dataset['warstwa'] = String(pozycja.warstwa);
  element.append(naglowek, umiejscowienie, opis);

  if (pozycja.komendy.length > 0) {
    const komendy = document.createElement('p');
    komendy.className = 'md-katalog__komendy';
    komendy.textContent = `Wykonuje: ${pozycja.komendy.join(' · ')}`;
    element.append(komendy);
  }
  if (pozycja.uwaga !== undefined) {
    const uwaga = document.createElement('p');
    uwaga.className = 'md-katalog__uwaga';
    uwaga.textContent = pozycja.uwaga;
    element.append(uwaga);
  }
  return element;
}
