import { pokazKomunikat } from '../aplikacja/komunikaty';
import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import {
  motywObowiazujacy,
  przelaczMotyw,
  ZDARZENIE_MOTYWU,
  type Motyw,
  type ZmianaMotywu,
} from '../motyw/motyw';
import { POZYCJE_USTAWIEN, type KodUstawienia } from '../strona-glowna/pozycje-ustawien';
import { utworzMenuProfilu, type MenuProfilu } from './menu-profilu';
import { opisNieczynnosci, wyjasnijNieczynnosc } from './nieczynne-w-pasku';
import { wykonajZaczep, type ZaczepyPaska } from './zaczepy-paska';

/**
 * Prawa strona paska górnego mieści kontrolki powiadomień, obecności globalnej, motywu i profilu, nie znając barw motywu i nie otwierając okien samodzielnie, tylko przez czynności dostarczone paskowi z zewnątrz.
 */
export interface AkcjePaska {
  /** Grupa kontrolek montowana przy prawej krawędzi paska. */
  element: HTMLElement;
  /** Ustawia licznik powiadomień; zero chowa plakietkę. */
  ustawPowiadomienia(ile: number): void;
  /** Zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu paska. */
  zamknij(): void;
}

/** Kody pozycji ustawień, które mają własną kontrolkę w pasku górnym aplikacji, obok pozostałych elementów interfejsu. */
const KODY_OBECNOSCI: readonly KodUstawienia[] = ['aod', 'mobile'];

export function utworzAkcjePaska(operator: string, zaczepy: ZaczepyPaska = {}): AkcjePaska {
  const element = document.createElement('div');
  element.className = 'dn-pasek-prawa';

  const { przycisk: dzwonek, licznik, koszyk } = utworzPowiadomienia(zaczepy);

  const motyw = przyciskPaska('ksiezyc', 'Motyw');
  motyw.addEventListener('click', () => ubierzMotyw(motyw, przelaczMotyw()));
  function naZmianeMotywu(zdarzenie: Event): void {
    ubierzMotyw(motyw, (zdarzenie as CustomEvent<ZmianaMotywu>).detail.obowiazujacy);
  }
  document.addEventListener(ZDARZENIE_MOTYWU, naZmianeMotywu);
  ubierzMotyw(motyw, motywObowiazujacy());

  const rozdzielacz = document.createElement('span');
  rozdzielacz.className = 'dn-pasek-gorny__rozdzielacz';

  const profil: MenuProfilu = utworzMenuProfilu({
    operator,
    zaczepy,
    naKomunikat: zglos,
  });

  element.append(koszyk, ...przyciskiObecnosci(zaczepy), motyw, rozdzielacz, profil.element);

  return {
    element,

    ustawPowiadomienia(ile) {
      const liczba = Math.max(0, Math.trunc(ile));
      licznik.textContent = String(liczba);
      licznik.hidden = liczba === 0;
      const stan = liczba === 0 ? 'brak nowych' : String(liczba);
      const opis =
        zaczepy.otworzPowiadomienia === undefined
          ? opisNieczynnosci('powiadomienia', `Powiadomienia: ${stan}`)
          : `Powiadomienia: ${stan}`;
      dzwonek.setAttribute('aria-label', opis);
      dzwonek.title = opis;
    },

    zamknij() {
      document.removeEventListener(ZDARZENIE_MOTYWU, naZmianeMotywu);
      profil.zamknij();
    },
  };
}

/**
 * Kontrolki obecności globalnej — Always On Display i Mobile — biorą nazwę i ikonę ze wspólnego wykazu pozycji ustawień, zasilającego też inne menu aplikacji.
 */
function przyciskiObecnosci(zaczepy: ZaczepyPaska): HTMLButtonElement[] {
  const przyciski: HTMLButtonElement[] = [];
  for (const kod of KODY_OBECNOSCI) {
    const pozycja = POZYCJE_USTAWIEN.find((p) => p.kod === kod);
    if (pozycja === undefined) continue;

    const przycisk = przyciskPaska(pozycja.ikona, `${pozycja.nazwa} — ${pozycja.wyjasnienie}`);
    przycisk.dataset['obecnosc'] = kod;
    // Zaczep przypisany paskowi jest odczytywany dopiero w tym miejscu, nie przy montażu.
    przycisk.addEventListener('click', () => wykonajZaczep(zaczepy, pozycja, zglos));
    przyciski.push(przycisk);
  }
  return przyciski;
}

/** Jedyne wyjście zdań paska do operatora; waga ostrzeżenie oznacza tu wyłącznie informację, nie że coś się zepsuło. */
function zglos(tytul: string, tresc: string): void {
  pokazKomunikat({ tytul, tresc, waga: 'ostrz' });
}

/** Dzwonek z plakietką licznika osadzoną w jego prawym górnym rogu, sygnalizującą liczbę nieprzeczytanych powiadomień. */
function utworzPowiadomienia(zaczepy: ZaczepyPaska): {
  przycisk: HTMLButtonElement;
  licznik: HTMLElement;
  koszyk: HTMLElement;
} {
  const koszyk = document.createElement('span');
  koszyk.className = 'dn-powiadomienia';

  const przycisk = przyciskPaska('dzwonek', 'Powiadomienia');
  // Zaczep czytany dopiero w chwili naciśnięcia — powłoka powstaje bez kanału do rdzenia.
  przycisk.addEventListener('click', () => {
    const otworz = zaczepy.otworzPowiadomienia;
    if (otworz === undefined) {
      wyjasnijNieczynnosc('powiadomienia');
      return;
    }
    otworz();
  });

  const licznik = document.createElement('span');
  licznik.className = 'dn-plakietka dn-plakietka--sygnal dn-powiadomienia__licznik';
  licznik.hidden = true;
  licznik.setAttribute('aria-hidden', 'true');

  koszyk.append(przycisk, licznik);
  return { przycisk, licznik, koszyk };
}

/** Przycisk ikonowy w wariancie przeznaczonym na ramę kokpitu, różniącym się rozmiarem i marginesem od wariantu zwykłego. */
function przyciskPaska(ikona: NazwaIkony, opis: string): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn-ikona dn-btn-ikona--na-ramie';
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
  przycisk.append(elementIkony(ikona, { rozmiar: 18 }));
  return przycisk;
}

/**
 * Przełącznik pokazuje motyw, w który przejdzie po naciśnięciu, a ta sama nastawa stoi też w menu profilu, czytana z jednego wspólnego miejsca.
 */
function ubierzMotyw(przycisk: HTMLButtonElement, obowiazujacy: Motyw): void {
  const ciemny = obowiazujacy === 'dark';
  const opis = ciemny ? 'Włącz motyw jasny' : 'Włącz motyw ciemny';
  przycisk.replaceChildren(elementIkony(ciemny ? 'slonce' : 'ksiezyc', { rozmiar: 18 }));
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
  przycisk.dataset.motyw = obowiazujacy;
}
