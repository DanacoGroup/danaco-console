import { poleTekstowe } from '../../modele/kontrolki-formularza';
import {
  KATALOG_FUNKCJI,
  grupyKatalogu,
  liczbaWedlugWykonania,
  type PozycjaKatalogu,
} from './katalog-funkcji-designu';
import { WARSTWY, utworzRozwiniecie, type Rozwiniecie } from './warstwy-designu';

/**
 * Wyszukiwarka funkcji modułu — droga warstwy czwartej do każdej pozycji
 * katalogu.
 *
 * Opracowanie modułu stanowi, że każda ukryta funkcja jest osiągalna jednym
 * kliknięciem, jednym skrótem albo jednym poleceniem w Chat Window.
 * Wyszukiwarka jest tą drogą dla pozycji, które nie mają własnego przycisku
 * w oknie: nazywa je, mówi, w którym oknie stoją i na której warstwie, oraz
 * podaje przy każdej jedno z trojga — komendę kontraktu, czynność samego okna
 * albo brak drogi.
 *
 * Wyszukiwarka niczego nie uruchamia i nie udaje, że uruchamia. Pozycja bez
 * drogi nie dostaje tu przycisku, który po naciśnięciu przeprosi — dostaje
 * zdanie o tym, czego brakuje, czytelne bez naciskania czegokolwiek.
 *
 * Szukanie idzie środkiem nazwy, opisu, grupy i okna, bo nazwy pozycji są
 * w części angielskie i złożone: szukanie wyłącznie od początku nazwy nie
 * znalazłoby „Upscaling neuronowy" po słowie wpisanym z pamięci.
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
 * Zdanie o zasięgu katalogu — trzy liczby liczone z wykazu, żadna wpisana.
 *
 * Opracowanie podsumowuje grupy zdaniem o siedemdziesięciu trzech funkcjach,
 * a wylicza ich więcej; przepisanie tamtej liczby dałoby napis rozjeżdżający się
 * z wykazem, który stoi pod nim.
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

/** Nazwa stanu wykonania widoczna dla Operatora — jedno brzmienie na moduł. */
const NAZWY_WYKONANIA: Readonly<Record<PozycjaKatalogu['wykonanie'], string>> = {
  komenda: 'komenda kontraktu',
  okno: 'czynność okna',
  'bez-drogi': 'bez drogi',
};

/**
 * Jeden wiersz katalogu.
 *
 * Zdanie o pokryciu buduje się z pól pozycji, nie z osobnego napisu przy każdej
 * z nich: pozycja z komendami wymienia je co do nazwy, pozycja bez drogi mówi,
 * czego brakuje. Dzięki temu dopisanie komendy do pozycji zmienia zdanie samo.
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
