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
 * Prawa strona paska górnego: powiadomienia, obecność globalna, motyw i profil.
 *
 * Jedna odpowiedzialność: kontrolki akcji paska. Barw ani wartości motywu ten
 * plik nie zna — przełączenie wykonuje warstwa `motyw/`. Okien nie
 * otwiera sam: dostaje czynności w `ZaczepyPaska` (`zaczepy-paska.ts`).
 *
 * Ikona `dzwonek` należy do powiadomień, nie do Always On Display: dzwonek nosi
 * plakietkę licznika (`dn-powiadomienia__licznik`), a AOD nie ma czego liczyć,
 * bo jest jednym pływającym podglądem, a nie zbiorem zdarzeń. AOD stoi przy
 * ikonie `oko`.
 *
 * Kontrolka, której ten pasek nie dostał w `ZaczepyPaska`, mówi wprost, że to
 * pasek nie ma drogi do okna — nie że okna nie ma w produkcie. Oba okna
 * obecności są zbudowane (`aod/kolumna-aod.ts`, `mobile/okno-mobile.ts`), a menu
 * Operatora stoi na drugim pasku (`aplikacja/menu-operatora.ts`).
 *
 * Dzwonek otwiera centrum powiadomień, gdy pasek dostał do niego zaczep
 * (`otworzPowiadomienia`). Bez zaczepu mówi wprost, że to ten pasek nie ma drogi
 * do kolumny — rodzina `notification.*` jest w kontrakcie, a rejestr w rdzeniu.
 */
export interface AkcjePaska {
  /** Grupa kontrolek montowana przy prawej krawędzi paska. */
  element: HTMLElement;
  /** Ustawia licznik powiadomień; zero chowa plakietkę. */
  ustawPowiadomienia(ile: number): void;
  /** Zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu paska. */
  zamknij(): void;
}

/** Kody pozycji ustawień, które mają własną kontrolkę w pasku. */
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
 * Kontrolki obecności globalnej — Always On Display i Mobile.
 *
 * Nazwa i ikona biorą się z `POZYCJE_USTAWIEN`, nie z literałów tego pliku. Ten
 * sam wykaz zasila listwę strony głównej, menu aplikacji, menu Operatora i menu
 * profilu — powtórzenie tu nazwy albo doboru ikony dałoby dwie prawdy o jednej
 * pozycji.
 */
function przyciskiObecnosci(zaczepy: ZaczepyPaska): HTMLButtonElement[] {
  const przyciski: HTMLButtonElement[] = [];
  for (const kod of KODY_OBECNOSCI) {
    const pozycja = POZYCJE_USTAWIEN.find((p) => p.kod === kod);
    if (pozycja === undefined) continue;

    const przycisk = przyciskPaska(pozycja.ikona, `${pozycja.nazwa} — ${pozycja.wyjasnienie}`);
    przycisk.dataset['obecnosc'] = kod;
    // Zaczep czytany dopiero tutaj — patrz `zaczepy-paska.ts`.
    przycisk.addEventListener('click', () => wykonajZaczep(zaczepy, pozycja, zglos));
    przyciski.push(przycisk);
  }
  return przyciski;
}

/** Jedyne wyjście zdań paska do Operatora; waga „ostrzeżenie" — nic się nie zepsuło. */
function zglos(tytul: string, tresc: string): void {
  pokazKomunikat({ tytul, tresc, waga: 'ostrz' });
}

/** Dzwonek z plakietką licznika osadzoną w jego prawym górnym rogu. */
function utworzPowiadomienia(zaczepy: ZaczepyPaska): {
  przycisk: HTMLButtonElement;
  licznik: HTMLElement;
  koszyk: HTMLElement;
} {
  const koszyk = document.createElement('span');
  koszyk.className = 'dn-powiadomienia';

  const przycisk = przyciskPaska('dzwonek', 'Powiadomienia');
  // Zaczep czytany dopiero w chwili naciśnięcia — powłoka powstaje bez kanału do
  // rdzenia, więc odczyt przy montażu zastałby go pustym. Pasek bez zaczepu mówi
  // wprost, że to on nie ma drogi do centrum — nie że centrum nie ma w produkcie.
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

/** Przycisk ikonowy w wariancie przeznaczonym na ramę kokpitu. */
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
 * Przełącznik pokazuje motyw, w który przejdzie po naciśnięciu — słońce przy
 * motywie ciemnym, księżyc przy jasnym. Oba motywy są równoprawne, więc
 * kontrolka żadnego nie wyróżnia.
 *
 * Ta sama nastawa stoi też w menu profilu; obie kontrolki czytają ją z `motyw/`
 * i żadna nie trzyma własnej kopii.
 */
function ubierzMotyw(przycisk: HTMLButtonElement, obowiazujacy: Motyw): void {
  const ciemny = obowiazujacy === 'dark';
  const opis = ciemny ? 'Włącz motyw jasny' : 'Włącz motyw ciemny';
  przycisk.replaceChildren(elementIkony(ciemny ? 'slonce' : 'ksiezyc', { rozmiar: 18 }));
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
  przycisk.dataset.motyw = obowiazujacy;
}
