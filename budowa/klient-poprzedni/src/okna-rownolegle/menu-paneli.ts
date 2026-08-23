import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { BRAK_KANALU, WYKAZ_NAGLOWEK, etykietaPozycji, zdanieOBrakach } from './etykiety-paneli';
import './menu.css';

/**
 * Treść menu `⋮` w nagłówku okna rozmowy — sekcja „Panele".
 *
 * Wiersz wykazu ma stały układ: znacznik otwarcia, ikona pozycji, nazwa
 * z przeznaczeniem, skrót klawiaturowy po prawej. Same pozycje przychodzą
 * z zewnątrz przez `OpcjeMenuPaneli.pozycje`, więc zmiana spisu paneli nie
 * dotyka tego pliku.
 *
 * W menu stoją wyłącznie pozycje, które realnie się otworzą, i każda jest
 * klikalna — brak dostępnej pozycji skraca wykaz, zamiast stawiać w nim wiersz
 * nieczynny. Żeby brak nie zniknął po cichu, pod wykazem stoi jedno zdanie
 * mówiące, ilu pozycji spisu nie da się otworzyć; przy `nieotwieralne === 0`
 * zdania nie ma wcale.
 *
 * Miejsce na znacznik otwarcia jest zajęte zawsze — przełączenie panelu zmienia
 * jego widoczność, nie szerokość wiersza, więc nazwy nie skaczą w poziomie.
 *
 * Skrót klawiaturowy jest tu wyłącznie napisem: ten plik niczego nie nasłuchuje
 * i nie rejestruje żadnego klawisza globalnie.
 *
 * Działania sesji nie są panelami i nie powstają tutaj. `sekcjeDalsze` doklejają
 * się za kreską, a kreska rysuje się tylko wtedy, gdy jest co za nią postawić.
 *
 * `PozycjaMenu` jest typem własnym, zgodnym strukturalnie
 * z `okna-pomocnicze/panele-otwieralne.ts`; rozjazd obu kształtów zatrzyma
 * kompilator na przypisaniu u odbiorcy.
 */

/** Pozycja, którą menu potrafi otworzyć. Zgodna strukturalnie z `okna-pomocnicze/panele-otwieralne.ts`. */
export interface PozycjaMenu {
  kod: string;
  nazwa: string;
  przeznaczenie: string;
  /** Ikona własna pozycji, po lewej przy nazwie. */
  ikona?: NazwaIkony;
  /** Skrót klawiaturowy do pokazania, np. `Ctrl+⇧+F`. Nic go tu nie podpina. */
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

    // Trzy różne prawdy, trzy różne zdania: nie ma czego otworzyć w ogóle,
    // część spisu czeka na kanał, albo spis jest domknięty i zdania nie ma.
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

/** Jeden wiersz wykazu: ptaszek · ikona · nazwa i przeznaczenie · skrót. */
function zbudujWiersz(pozycja: PozycjaMenu, naWybor: () => void): HTMLElement {
  const wiersz = document.createElement('button');
  wiersz.type = 'button';
  wiersz.className = 'dn-menu-paneli__pozycja';
  wiersz.setAttribute('role', 'menuitem');
  wiersz.dataset.panel = pozycja.kod;
  wiersz.dataset.nazwa = pozycja.nazwa;
  wiersz.title = pozycja.przeznaczenie;

  // Miejsce znacznika zajęte zawsze — inaczej wiersze skakałyby w poziomie
  // przy każdym przełączeniu panelu.
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
