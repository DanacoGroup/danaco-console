import './nawigacja-modulow.css';

import { elementIkony } from '../ikony/ikony';
import {
  opisLiczbyPozycji,
  srodowiskoWczytywane,
  SRODOWISKO_DOMYSLNE,
  type KluczSrodowiska,
  type PozycjaModulu,
  type Srodowisko,
} from './srodowiska';
import type { ZrodloNawigacji } from './zrodlo-nawigacji';

/** Słuchacz wyboru pozycji nawigacji, wywoływany z pozycją wybraną przez Operatora oraz danymi bieżącego środowiska. */
export type SluchaczModulu = (pozycja: PozycjaModulu, dane: Srodowisko) => void;

/** Słuchacz wskazania, które nie ma odpowiednika w wykazie bieżącego środowiska — otwarcia modułu spoza kolumny nawigacji. */
export type SluchaczBrakuPozycji = (kluczPozycji: string) => void;

export interface NawigacjaModulow {
  /** Kolumna montowana w korpusie powłoki. */
  element: HTMLElement;
  /** Podłącza źródło wykazu i przeładowuje kolumnę z rdzenia. */
  podlaczZrodlo(zrodlo: ZrodloNawigacji): void;
  /** Przestawia wykaz na wskazane środowisko i wybiera jego pierwszą pozycję. */
  pokaz(klucz: KluczSrodowiska): void;
  /** Wskazuje pozycję po kluczu; wykaz w drodze zapamiętuje wskazanie. */
  wybierz(kluczPozycji: string): void;
  /** Otwiera moduł, którego nie ma w wykazie środowiska; moduł bez wiersza macierzy otwiera się tą drogą. */
  wskazPozaWykazem(pozycja: PozycjaModulu): void;
  /** Pozycja wybrana albo brak, gdy wykaz jest pusty. */
  wybrana(): PozycjaModulu | undefined;
  /** Środowisko obecnie pokazywane. */
  srodowisko(): Srodowisko;
  naWybor(sluchacz: SluchaczModulu): void;
  /** Zgłasza wskazanie bez odpowiednika w wykazie — obsługa należy do warstwy znającej katalog modułów. */
  naBrakPozycji(sluchacz: SluchaczBrakuPozycji): void;
}

/** Puste elementy kolumny nawigacji — powstają raz przy budowie, a ich treść wymienia się przy każdym kolejnym wykazie środowiska. */
interface SzkieletKolumny {
  element: HTMLElement;
  nazwa: HTMLElement;
  motto: HTMLElement;
  licznik: HTMLElement;
  lista: HTMLElement;
  stan: HTMLElement;
}

/**
 * Buduje szkielet kolumny. Stoi osobno od wytwórni, bo to jedna odpowiedzialność:
 * kształt kolumny w drzewie dokumentu, niezależny od tego, skąd bierze się wykaz
 * i jak rozstrzyga się wskazanie pozycji.
 */
function zbudujSzkielet(): SzkieletKolumny {
  const element = document.createElement('nav');
  element.className = 'dn-nawigacja';

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-nawigacja__naglowek';

  const nazwa = document.createElement('h2');
  nazwa.className = 'dn-nawigacja__srodowisko';

  const motto = document.createElement('p');
  motto.className = 'dn-nawigacja__motto';

  const licznik = document.createElement('p');
  licznik.className = 'dn-nawigacja__licznik';

  naglowek.append(nazwa, motto, licznik);

  const lista = document.createElement('ul');
  lista.className = 'dn-nawigacja__lista';

  const stan = document.createElement('p');
  stan.className = 'dn-nawigacja__stan';
  stan.setAttribute('role', 'status');

  element.append(naglowek, lista, stan);
  return { element, nazwa, motto, licznik, lista, stan };
}

/** Jeden wiersz wykazu środowiska, budowany z pozycji modułu. Wskazanie oddaje wytwórni — sam wiersz nie zna stanu kolumny. */
function wierszPozycji(pozycja: PozycjaModulu, naWskazanie: (klucz: string) => void): HTMLLIElement {
  const punkt = document.createElement('li');

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-nawigacja__pozycja';
  przycisk.dataset.pozycja = pozycja.klucz;
  przycisk.title = pozycja.opis === '' ? pozycja.nazwa : `${pozycja.nazwa} — ${pozycja.opis}`;
  przycisk.addEventListener('click', () => naWskazanie(pozycja.klucz));

  const napis = document.createElement('span');
  napis.className = 'dn-nawigacja__nazwa';
  napis.textContent = pozycja.nazwa;

  przycisk.append(elementIkony(pozycja.ikona, { rozmiar: 18 }), napis);
  punkt.append(przycisk);
  return punkt;
}

export function utworzNawigacjeModulow(): NawigacjaModulow {
  const sluchacze: SluchaczModulu[] = [];
  const sluchaczeBraku: SluchaczBrakuPozycji[] = [];
  let zrodlo: ZrodloNawigacji | null = null;
  let dane: Srodowisko = srodowiskoWczytywane(SRODOWISKO_DOMYSLNE);
  let wybranaPozycja: PozycjaModulu | undefined;
  /** Wskazanie z adresu, które przyszło, zanim rdzeń oddał wykaz. */
  let oczekiwanyWybor: string | undefined;
  /** Numer zapytania — odpowiedź przedawniona nie nadpisuje świeższej. */
  let numerZapytania = 0;

  const { element, nazwa, motto, licznik, lista, stan } = zbudujSzkielet();

  /** Pyta rdzeń o wykaz środowiska; do czasu odpowiedzi rysuje stan w drodze. */
  function pokaz(klucz: KluczSrodowiska): void {
    const numer = ++numerZapytania;
    narysuj(srodowiskoWczytywane(klucz));
    if (zrodlo === null) return;
    void zrodlo(klucz).then((wykaz) => {
      if (numer === numerZapytania) narysuj(wykaz);
    });
  }

  /** Przerysowuje kolumnę wykazem, który właśnie stał się prawdą powłoki. */
  function narysuj(wykaz: Srodowisko): void {
    dane = wykaz;
    element.dataset.srodowisko = dane.klucz;
    element.dataset.stan = dane.stan;
    element.setAttribute(
      'aria-label',
      dane.rodzaj === 'sekcje' ? 'Panel orkiestracji' : 'Moduły środowiska',
    );
    nazwa.textContent = dane.nazwa;
    motto.textContent = dane.motto;
    licznik.textContent = opisLiczbyPozycji(dane);

    lista.replaceChildren(...dane.pozycje.map((pozycja) => wierszPozycji(pozycja, wybierz)));
    stan.textContent = opisStanu(dane);
    stan.hidden = stan.textContent === '';
    wybranaPozycja = undefined;

    const zadana = oczekiwanyWybor;
    oczekiwanyWybor = undefined;
    const pierwsza = dane.pozycje[0];
    if (zadana !== undefined && dane.pozycje.some((wpis) => wpis.klucz === zadana)) {
      wybierz(zadana);
      return;
    }
    // Pierwsza pozycja jest wartością zastępczą, więc wykaz nigdy nie zostaje bez wskazania.
    if (pierwsza !== undefined) wybierz(pierwsza.klucz);
    if (zadana !== undefined) zglosBrak(zadana);
  }

  /** Oddaje wskazanie warstwie, która zna katalog modułów spoza środowiska. */
  function zglosBrak(kluczPozycji: string): void {
    for (const sluchacz of sluchaczeBraku) sluchacz(kluczPozycji);
  }

  function wybierz(kluczPozycji: string): void {
    const pozycja = dane.pozycje.find((wpis) => wpis.klucz === kluczPozycji);
    if (pozycja === undefined) {
      // Wskazanie z adresu przychodzi przed odpowiedzią rdzenia — czeka na
      // wykaz zamiast przepadać.
      if (dane.stan === 'ladowanie') oczekiwanyWybor = kluczPozycji;
      else zglosBrak(kluczPozycji);
      return;
    }

    wybranaPozycja = pozycja;
    for (const przycisk of lista.querySelectorAll<HTMLButtonElement>('.dn-nawigacja__pozycja')) {
      const czynna = przycisk.dataset.pozycja === kluczPozycji;
      // Znacznik aria-current niesie wskazanie także wtedy, gdy barwa nie dociera do Operatora.
      if (czynna) przycisk.setAttribute('aria-current', 'page');
      else przycisk.removeAttribute('aria-current');
    }

    for (const sluchacz of sluchacze) sluchacz(pozycja, dane);
  }

  narysuj(dane);

  return {
    element,

    podlaczZrodlo(nowe) {
      zrodlo = nowe;
      pokaz(dane.klucz);
    },

    pokaz,
    wybierz,

    wskazPozaWykazem(pozycja) {
      wybranaPozycja = pozycja;
      // Żaden wiersz kolumny nie odpowiada temu modułowi, więc żaden nie niesie znacznika aria-current.
      for (const przycisk of lista.querySelectorAll<HTMLButtonElement>('.dn-nawigacja__pozycja')) {
        przycisk.removeAttribute('aria-current');
      }
      for (const sluchacz of sluchacze) sluchacz(pozycja, dane);
    },

    wybrana: () => wybranaPozycja,
    srodowisko: () => dane,
    naWybor: (sluchacz) => void sluchacze.push(sluchacz),
    naBrakPozycji: (sluchacz) => void sluchaczeBraku.push(sluchacz),
  };
}

/** Zdanie pod wykazem, gdy pozycji wykazu nie ma — opis rzeczywistego stanu wykazu, a nie pusta atrapa listy pozycji modułów. */
function opisStanu(dane: Srodowisko): string {
  if (dane.stan === 'ladowanie') return 'Wykaz modułów wczytywany z rdzenia…';
  if (dane.stan === 'blad') return dane.blad ?? 'Rdzeń nie oddał wykazu modułów.';
  if (dane.pozycje.length === 0) return 'Rdzeń nie przypisał temu środowisku żadnego modułu.';
  return '';
}
