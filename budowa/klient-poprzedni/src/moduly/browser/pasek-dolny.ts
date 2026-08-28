import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { przycisk } from '../../modele/kontrolki-formularza';
import type { CzynnosciPaskaDolnego } from './czynnosci-paska-dolnego';
import { KLASY_DYMKA, OBJASNIENIA } from './etykiety-browser';
import { KLASA_PRZYCISKU } from './przyciski-browser';
import { utworzSterWyboru } from './ster-wyboru';

/**
 * Pasek dolny Browser Window — podział ekranu, tryb czytnika, tryb adnotacji,
 * zrzut ekranu, wyodrębnienie danych, tłumaczenie, zakładka, makro i pobrania.
 * Jedna odpowiedzialność: złożenie kontrolek paska.
 */
export interface PasekDolny {
  element: HTMLElement;
  /** Nanosi stan trybu adnotacji na przełącznik paska — lustro warstwy, nie jej właściciel. */
  ustawTrybAdnotacji(wlaczony: boolean): void;
}

/** Czego pasek dolny potrzebuje od ramy okna: ujścia zdarzeń dla każdego z jego przycisków sterujących. */
export interface UjsciaPaska {
  naPodzial(wlaczony: boolean): void;
  naCzytnik(wlaczony: boolean): void;
  /** Dopisuje wyodrębnioną treść do panelu wyodrębnień okna. */
  naWyodrebnienie(rodzaj: string): void;
  /** Włącza albo wyłącza płótno adnotacji nad sceną podglądu. */
  naAdnotacje(wlaczony: boolean): void;
  powiedz(tresc: string, powodzenie: boolean): void;
}

/**
 * Trzy rodzaje wyodrębnienia danych z migawki — nastawa pozycji „Wyodrębnij
 * dane". Opis przy pozycji mówi, co wyjdzie z wyodrębnienia, bo same nazwy
 * tego nie rozstrzygają.
 */
const RODZAJE_WYODREBNIENIA = [
  {
    wartosc: 'tekst',
    etykieta: 'Treść renderowana strony',
    opis: 'Sam tekst strony, tak jak czyta go Operator — bez znaczników.',
  },
  {
    wartosc: 'zrodlo',
    etykieta: 'Źródło strony (HTML)',
    opis: 'Znaczniki pobrane przez rdzeń, bez uruchamiania skryptów strony.',
  },
  {
    wartosc: 'adres',
    etykieta: 'Adres i tytuł strony',
    opis: 'Dwa wiersze: adres migawki i tytuł odczytany ze strony.',
  },
];

/** Przełącznik stanu na przycisku: stan niesie atrybut dostępności, a nie osobna klasa stylu wizualnego. */
function przelacz(kontrolka: HTMLButtonElement, oddaj: (wlaczony: boolean) => void): void {
  const wlaczony = kontrolka.getAttribute('aria-pressed') !== 'true';
  kontrolka.setAttribute('aria-pressed', String(wlaczony));
  oddaj(wlaczony);
}

/** Przycisk-przełącznik paska dolnego wraz z ustawionym stanem początkowym przed pierwszym naciśnięciem. */
function przyciskPrzelacznik(
  nazwa: string,
  oddaj: (wlaczony: boolean) => void,
): HTMLButtonElement {
  const element = przycisk(nazwa, KLASA_PRZYCISKU.zarys);
  element.setAttribute('aria-pressed', 'false');
  element.addEventListener('click', () => przelacz(element, oddaj));
  return element;
}

/**
 * Czynności strony przychodzą z zewnątrz, bo tę samą trójkę wywołań ma pasek
 * pływający zaznaczenia: „Tłumacz" na pasku dolnym i „Tłumacz" nad zaznaczeniem
 * to jedno wywołanie `context.transfer`, a nie dwa podobne.
 */
export function utworzPasekDolny(
  czynnosci: CzynnosciPaskaDolnego,
  ujscia: UjsciaPaska,
): PasekDolny {
  const podzial = przyciskPrzelacznik('Podziel ekran', ujscia.naPodzial);
  const czytnik = przyciskPrzelacznik('Tryb czytnika', ujscia.naCzytnik);
  const zrzut = przycisk('Zrzut ekranu', KLASA_PRZYCISKU.zarys);
  const tlumacz = przycisk('Tłumacz', KLASA_PRZYCISKU.zarys);
  const wyodrebnij = przycisk('Wyodrębnij dane', KLASA_PRZYCISKU.zarys);
  // Ster, nie wyświetlacz: na uchwycie stoi rodzaj wybrany teraz, nie stała nazwa nastawy.
  const rodzaj = utworzSterWyboru({
    nastawa: 'Rodzaj wyodrębnienia',
    pozycje: RODZAJE_WYODREBNIENIA,
    klasa: 'mb-pasek-dolny__ster',
    podpis: false,
  });
  const podswietl = przycisk('Podświetl i adnotuj', KLASA_PRZYCISKU.zarys);
  const adnotacja = przyciskPrzelacznik('Tryb adnotacji', ujscia.naAdnotacje);

  // Zakładka, nagrywarka i pobrania mają czynności w panelu rodzin — pasek ich nie dubluje.
  const bezObslugi: HTMLButtonElement[] = [];

  zrzut.addEventListener('click', () => void czynnosci.zrzutEkranu());
  tlumacz.addEventListener('click', () => void czynnosci.doTlumaczenia());
  podswietl.addEventListener('click', () => void czynnosci.adnotujFragment());
  wyodrebnij.addEventListener('click', () => ujscia.naWyodrebnienie(rodzaj.wartosc()));

  const element = document.createElement('footer');
  element.className = 'mb-pasek-dolny';
  element.setAttribute('aria-label', 'Narzędzia okna przeglądarki');
  element.append(
    podzial,
    utworzDymekObjasnienia(OBJASNIENIA.podzialEkranu, KLASY_DYMKA),
    czytnik,
    utworzDymekObjasnienia(OBJASNIENIA.trybCzytnika, KLASY_DYMKA),
    adnotacja,
    utworzDymekObjasnienia(OBJASNIENIA.adnotacja, KLASY_DYMKA),
    zrzut,
    utworzDymekObjasnienia(OBJASNIENIA.zrzut, KLASY_DYMKA),
    rodzaj.element,
    wyodrebnij,
    tlumacz,
    podswietl,
    ...bezObslugi,
  );

  return {
    element,
    ustawTrybAdnotacji(wlaczony) {
      adnotacja.setAttribute('aria-pressed', String(wlaczony));
    },
  };
}
