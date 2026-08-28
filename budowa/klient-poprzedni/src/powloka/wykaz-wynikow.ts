import './wykaz-wynikow.css';

import { utworzMenuDrzewo, type MenuDrzewo, type PozycjaMenu } from '../komponenty/menu-drzewo';
import { utworzBrakWynikow, type BrakWynikow } from './brak-wynikow';
import {
  NAZWY_GRUP,
  PORZADEK_GRUP,
  ZAKRES_PRZESZUKANIA,
  type RodzajWpisu,
  type WpisWyszukiwania,
  type ZrodloWyszukiwania,
} from './wyszukiwanie-globalne';

// Wykaz wyników wyszukiwania globalnego — obsada biblioteki menu drzewa.

/** Nazwa nastawy dla czytnika ekranu; na ekranie nie staje, bo uchwytu tego menu tutaj w ogóle wcale nie ma. */
const NASTAWA = 'Wyszukiwanie globalne';

/** Ile ostatnich wyborów pamięta ten wykaz w historii, trzymanej wyłącznie w pamięci tego samego klienta. */
const ILE_OSTATNICH = 5;

/** Nagłówek grupy historii pokazywanej przy pustej frazie w polu wyszukiwania paska globalnego klienta. */
const GRUPA_OSTATNICH = 'Ostatnio wybrane';

export interface WykazWynikow {
  /** Nośnik montowany pod polem paska; sam nie zajmuje wysokości. */
  element: HTMLElement;
  /** Czy wykaz jest w tej chwili rozwinięty. */
  otwarty(): boolean;
  /** Rozwija wykaz bez zabierania ogniska polu paska. */
  otworz(): void;
  /** Podaje frazę z pola paska — to jest całe filtrowanie. */
  ustawFraze(fraza: string): void;
  /** Przesuwa wyróżnienie; zawija się na końcach. Bez ruszania ogniska. */
  przesun(krok: number): void;
  /** Zatwierdza wyróżnioną. Fałsz znaczy „nie było czego zatwierdzić". */
  zatwierdz(): boolean;
  /** Zwija wykaz wraz ze stanem pustki. */
  zamknij(): void;
}

export interface OpcjeWykazuWynikow {
  /** Skąd biorą się wpisy. */
  zrodlo: ZrodloWyszukiwania;
  /** Wołane po każdym wyborze — pasek czyści wtedy pole. */
  naWybor?(wpis: WpisWyszukiwania): void;
}

export function utworzWykazWynikow(opcje: OpcjeWykazuWynikow): WykazWynikow {
  /** Klucze ostatnio wybranych, najświeższy z przodu. Ginie z zamknięciem okna. */
  const ostatnie: string[] = [];
  let fraza = '';
  let czynny = false;

  const menu: MenuDrzewo = utworzMenuDrzewo({
    nastawa: NASTAWA,
    bezUchwytu: true,
    // Opis przy każdej pozycji byłby ścianą tekstu: wykaz bywa
    // kilkudziesięciopozycyjny i płaski.
    opisTylkoPrzyWyroznionej: true,
    kierunek: 'dol',
    naWybor: (klucz) => wykonaj(klucz),
  });

  const pustka: BrakWynikow = utworzBrakWynikow();
  pustka.element.classList.add('dn-wyniki__pustka');
  pustka.element.hidden = true;

  const element = document.createElement('div');
  element.className = 'dn-wyniki';
  element.append(menu.element, pustka.element);

  /** Wpis po kluczu; wykaz nie trzyma własnej kopii materiału. */
  function wpisPoKluczu(klucz: string): WpisWyszukiwania | undefined {
    return opcje.zrodlo.wpisy().find((wpis) => wpis.klucz === klucz);
  }

  /** Jeden wiersz wykazu. `wybrany` zostaje fałszem: to wykaz, nie nastawa. */
  function lisc(wpis: WpisWyszukiwania): PozycjaMenu {
    return {
      rodzaj: 'wybor',
      klucz: wpis.klucz,
      nazwa: wpis.nazwa,
      ...(wpis.opis === '' ? {} : { opis: wpis.opis }),
      ikona: wpis.ikona,
      wybrany: false,
    };
  }

  /** Wpisy jednego rodzaju pod wspólnym nagłówkiem; grupa pusta nie wchodzi. */
  function grupa(rodzaj: RodzajWpisu, wpisy: readonly WpisWyszukiwania[]): PozycjaMenu[] {
    const swoje = wpisy.filter((wpis) => wpis.rodzaj === rodzaj);
    if (swoje.length === 0) return [];
    return [{ rodzaj: 'grupa', nazwa: NAZWY_GRUP[rodzaj], dzieci: swoje.map(lisc) }];
  }

  // Drzewo pokazywane przy pustej frazie — historia, nie cały katalog, krótka lista pozycji.
  function drzewoHistorii(wpisy: readonly WpisWyszukiwania[]): PozycjaMenu[] {
    const znalezione = ostatnie
      .map((klucz) => wpisy.find((wpis) => wpis.klucz === klucz))
      .filter((wpis): wpis is WpisWyszukiwania => wpis !== undefined);
    if (znalezione.length === 0) return [];
    return [{ rodzaj: 'grupa', nazwa: GRUPA_OSTATNICH, dzieci: znalezione.map(lisc) }];
  }

  /** Ile pozycji mechanizm naprawdę narysował — pytanie o wynik, nie o regułę. */
  function ileNarysowanych(): number {
    return menu.element.querySelectorAll('[role="menuitemradio"]').length;
  }

  /** Pokazuje wzorzec pustki zamiast wykazu albo chowa go z powrotem. */
  function ustawPustke(pokaz: boolean): void {
    pustka.element.hidden = !pokaz;
    element.dataset['pustka'] = pokaz ? 'tak' : 'nie';
    if (pokaz) menu.zwin();
  }

  function odrysuj(): void {
    const wpisy = opcje.zrodlo.wpisy();
    const szukane = fraza.trim();

    // Stan odczytu i odmowa idą przed dopasowaniem: wykaz modułów bywa jeszcze w drodze z rdzenia.
    const stan = opcje.zrodlo.stan();
    if (stan === 'blad') {
      pustka.odmowa(opcje.zrodlo.odmowa());
      ustawPustke(true);
      return;
    }
    if (stan === 'ladowanie') {
      pustka.wOdczycie('wykazu modułów środowiska');
      ustawPustke(true);
      return;
    }

    if (szukane === '') {
      const historia = drzewoHistorii(wpisy);
      if (historia.length === 0) {
        pustka.pusto(
          'ostatnich wyborów',
          'Wpisz choć jeden znak — wykaz przeszuka wtedy '
          + `${ZAKRES_PRZESZUKANIA}.`,
        );
        ustawPustke(true);
        return;
      }
      ustawPustke(false);
      menu.ustawFraze('');
      menu.ustaw(NASTAWA, historia);
      if (czynny) menu.rozwin();
      return;
    }

    ustawPustke(false);
    menu.ustaw(NASTAWA, PORZADEK_GRUP.flatMap((rodzaj) => grupa(rodzaj, wpisy)));
    menu.ustawFraze(szukane);
    if (czynny) menu.rozwin();

    if (ileNarysowanych() === 0) {
      pustka.bezTrafien(szukane, ZAKRES_PRZESZUKANIA);
      ustawPustke(true);
    }
  }

  function wykonaj(klucz: string): void {
    const wpis = wpisPoKluczu(klucz);
    czynny = false;
    // Fraza ginie razem z wyborem: pole paska jest po wyborze puste, filtr by już nie miał treści.
    fraza = '';
    ustawPustke(false);
    menu.zwin();
    if (wpis === undefined) return;
    zapamietaj(klucz);
    opcje.naWybor?.(wpis);
    wpis.wykonaj();
  }

  function zapamietaj(klucz: string): void {
    const gdzie = ostatnie.indexOf(klucz);
    if (gdzie >= 0) ostatnie.splice(gdzie, 1);
    ostatnie.unshift(klucz);
    if (ostatnie.length > ILE_OSTATNICH) ostatnie.length = ILE_OSTATNICH;
  }

  return {
    element,

    otwarty: () => czynny,

    otworz() {
      czynny = true;
      odrysuj();
    },

    ustawFraze(nowa) {
      fraza = nowa;
      if (!czynny) czynny = true;
      odrysuj();
    },

    przesun(krok) {
      if (!czynny) return;
      menu.przesunWyroznienie(krok);
    },

    zatwierdz() {
      if (!czynny) return false;
      if (!pustka.element.hidden) return false;
      return menu.wybierzWyrozniona();
    },

    zamknij() {
      czynny = false;
      fraza = '';
      ustawPustke(false);
      menu.zwin();
    },
  };
}
