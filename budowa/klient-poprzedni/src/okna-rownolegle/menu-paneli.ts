import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { BRAK_KANALU, WYKAZ_NAGLOWEK, etykietaPozycji, zdanieOBrakach } from './etykiety-paneli';
import './menu.css';

/**
 * Treść menu w nagłówku okna rozmowy dla sekcji Panele pokazuje wyłącznie pozycje, które realnie się otworzą, każda klikalna, ze skrótem klawiaturowym wyłącznie jako napisem, a typ pozycji jest zgodny strukturalnie z rejestrem paneli otwieralnych.
 */
export interface PozycjaMenu {
  kod: string;
  nazwa: string;
  przeznaczenie: string;
  /** Ikona własna pozycji, po lewej przy nazwie. */
  ikona?: NazwaIkony;
  /** Skrót klawiaturowy pokazywany jako napis obok pozycji menu; ten plik go w ogóle nie podpina. */
  skrot?: string;
}

export interface OpcjeMenuPaneli {
  pozycje: readonly PozycjaMenu[];
  /** Ile pozycji spisu nie da się dziś otworzyć — do zdania pod wykazem. */
  nieotwieralne: number;
  czyOtwarty(kod: string): boolean;
  naWybor(kod: string): void;
  /** Sekcje doklejane za kreską — działania sesji, budowane poza tym pakietem. */
  sekcjeDalsze?: readonly HTMLElement[];
}

export interface MenuPaneli {
  element: HTMLElement;
  odswiez(): void;
  zamknij(): void;
}

export function utworzMenuPaneli(opcje: OpcjeMenuPaneli): MenuPaneli {
  const element = document.createElement('div');
  element.className = 'dn-menu-paneli';
  element.setAttribute('role', 'none');

  const sekcja = document.createElement('div');
  sekcja.className = 'dn-menu-paneli__sekcja';
  sekcja.setAttribute('role', 'group');
  sekcja.setAttribute('aria-label', WYKAZ_NAGLOWEK);

  const naglowek = document.createElement('p');
  naglowek.className = 'dn-menu-paneli__naglowek';
  naglowek.textContent = WYKAZ_NAGLOWEK;

  const wykaz = document.createElement('div');
  wykaz.className = 'dn-menu-paneli__wykaz';

  const braki = document.createElement('p');
  braki.className = 'dn-menu-paneli__braki';

  sekcja.append(naglowek, wykaz, braki);
  element.append(sekcja);

  // Kreska i sekcje dalsze powstają tylko wtedy, gdy jest co za nią postawić.
  const dalsze = opcje.sekcjeDalsze ?? [];
  if (dalsze.length > 0) {
    const kreska = document.createElement('hr');
    kreska.className = 'dn-menu-paneli__kreska';
    element.append(kreska, ...dalsze);
  }

  const wiersze = new Map<string, HTMLElement>();

  for (const pozycja of opcje.pozycje) {
    const wiersz = zbudujWiersz(pozycja, () => opcje.naWybor(pozycja.kod));
    wiersze.set(pozycja.kod, wiersz);
    wykaz.append(wiersz);
  }

  function odswiez(): void {
    for (const [kod, wiersz] of wiersze) {
      const otwarty = opcje.czyOtwarty(kod);
      wiersz.setAttribute('aria-checked', String(otwarty));
      const nazwa = wiersz.dataset.nazwa ?? kod;
      wiersz.setAttribute('aria-label', etykietaPozycji(nazwa, otwarty));
    }

    // Trzy różne prawdy, trzy zdania: nic czego otworzyć, część czeka na kanał, spis jest zamknięty.
    if (opcje.pozycje.length === 0) {
      braki.hidden = false;
      braki.textContent = BRAK_KANALU;
      return;
    }
    if (opcje.nieotwieralne > 0) {
      braki.hidden = false;
      braki.textContent = zdanieOBrakach(opcje.nieotwieralne);
      return;
    }
    braki.hidden = true;
    braki.textContent = '';
  }

  odswiez();

  return {
    element,
    odswiez,
    zamknij() {
      wiersze.clear();
      element.replaceChildren();
    },
  };
}

/** Jeden wiersz wykazu menu: znacznik otwarcia, ikona, nazwa z przeznaczeniem oraz skrót klawiaturowy po prawej. */
function zbudujWiersz(pozycja: PozycjaMenu, naWybor: () => void): HTMLElement {
  const wiersz = document.createElement('button');
  wiersz.type = 'button';
  wiersz.className = 'dn-menu-paneli__pozycja';
  wiersz.setAttribute('role', 'menuitem');
  wiersz.dataset.panel = pozycja.kod;
  wiersz.dataset.nazwa = pozycja.nazwa;
  wiersz.title = pozycja.przeznaczenie;

  // Miejsce znacznika jest zajęte zawsze, inaczej wiersze skakałyby w poziomie przy każdym przełączeniu.
  const znacznik = document.createElement('span');
  znacznik.className = 'dn-menu-paneli__ptaszek';
  znacznik.append(elementIkony('ptaszek', { rozmiar: 14 }));

  const tresc = document.createElement('span');
  tresc.className = 'dn-menu-paneli__tresc';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-menu-paneli__nazwa';
  nazwa.textContent = pozycja.nazwa;

  const przeznaczenie = document.createElement('span');
  przeznaczenie.className = 'dn-menu-paneli__przeznaczenie';
  przeznaczenie.textContent = pozycja.przeznaczenie;

  tresc.append(nazwa, przeznaczenie);
  wiersz.append(znacznik);

  if (pozycja.ikona !== undefined) {
    const ikona = document.createElement('span');
    ikona.className = 'dn-menu-paneli__ikona';
    ikona.append(elementIkony(pozycja.ikona, { rozmiar: 16 }));
    wiersz.append(ikona);
  }

  wiersz.append(tresc);

  if (pozycja.skrot !== undefined) {
    const skrot = document.createElement('kbd');
    skrot.className = 'dn-menu-paneli__skrot';
    skrot.textContent = pozycja.skrot;
    wiersz.append(skrot);
  }

  wiersz.addEventListener('click', naWybor);
  return wiersz;
}
