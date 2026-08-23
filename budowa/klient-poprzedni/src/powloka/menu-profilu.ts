import { utworzMenuDrzewo, type MenuDrzewo, type PozycjaMenu } from '../komponenty/menu-drzewo';
import {
  motywObowiazujacy,
  przelaczMotyw,
  ZDARZENIE_MOTYWU,
} from '../motyw/motyw';
import {
  POZYCJE_USTAWIEN,
  type KodUstawienia,
  type PozycjaUstawienia,
} from '../strona-glowna/pozycje-ustawien';
import { wykonajZaczep, type OdbiorcaZdania, type ZaczepyPaska } from './zaczepy-paska';

/**
 * Menu profilu Operatora w pasku powłoki środowiska: siedem pozycji w trzech
 * grupach oraz stopka.
 *
 *   Tożsamość i dostęp
 *     Modele, konta i tożsamość     → `modele`
 *     Dostępy i katalog roboczy     → `dostepy`
 *     Punkty izolacji               → `punkty-izolacji`
 *     Ustawienia (hasło, wejście)   → `ustawienia`
 *   Widok
 *     Motyw ciemny                  → `PrzelacznikMenu`, nastawa dwustanowa
 *   Obecność globalna
 *     Always On Display             → `aod`
 *     Mobile                        → `mobile`
 *   Stopka
 *     Okno konfiguracji             → `konfiguracja`
 *
 * Pozycje nie są nowym wykazem: powstają z `POZYCJE_USTAWIEN`, z którego żyje
 * także listwa strony głównej, menu aplikacji i menu Operatora
 * (`aplikacja/menu-operatora.ts`). Tutaj leży wyłącznie podział na grupy oraz
 * to, że `konfiguracja` idzie do stopki — stopka jest drogą do rejestru, czyli
 * miejscem rzeczy pełniejszej niż pozycje nad nią. Zmiana nazwy albo ikony
 * w tamtym wykazie zmienia to menu sama.
 *
 * Mechanizm jest w całości z biblioteki: `komponenty/menu-drzewo` niesie
 * grupowanie (`GrupaMenu`), nastawę dwustanową (`PrzelacznikMenu`), stopkę,
 * znacznik wyboru i wędrówkę klawiaturą.
 *
 * Awatar jest uchwytem i niesie wartość: stoją na nim inicjały, a nie napis
 * rodzajowy „Profil". Znak „O" nie jest inicjałem człowieka — to pierwsza
 * litera słowa „Operator", jedynej nazwy, jaką produkt zna.
 *
 * Motyw ma jedną prawdę: przełącznik w menu i przycisk w pasku czytają
 * `motyw/motywObowiazujacy()` i nasłuchują `ZDARZENIE_MOTYWU`, więc zmiana
 * w jednym miejscu przerysowuje drugie.
 *
 * Nie ma tu konta, adresu, przełączania tożsamości ani „Wyloguj": platforma
 * kont osobowych nie prowadzi (kontrakt przy `auth.register`: „Konta NIE
 * zaklada i adresu e-mail nie przyjmuje: bramka jest jedna, Operator
 * bezimienny"), a po zalogowaniu nie ma bram. Zdanie granicy stoi na ekranie,
 * w stopce, tym samym brzmieniem co w `aplikacja/menu-operatora.ts`.
 */

/** Klucz przełącznika motywu; nie jest kodem ustawienia i nie może się z nim zderzyć. */
export const KLUCZ_MOTYWU = 'widok:motyw-ciemny';

/** Napis na uchwycie, gdy nazwa Operatora nie daje ani jednej litery. */
const ZNAK_ZASTEPCZY = 'O';

/** Kod pozycji, która idzie do stopki jako droga do rejestru. */
const KOD_STOPKI: KodUstawienia = 'konfiguracja';

/** Podział pozycji na grupy. Kolejność grup i kolejność w grupie — jak tu. */
const GRUPY: readonly { nazwa: string; kody: readonly KodUstawienia[] }[] = [
  {
    nazwa: 'Tożsamość i dostęp',
    kody: ['modele', 'dostepy', 'punkty-izolacji', 'ustawienia'],
  },
  {
    nazwa: 'Obecność globalna',
    kody: ['aod', 'mobile'],
  },
];

/** Zdanie granicy — dlaczego nie ma tu konta ani wylogowania. */
const ZDANIE_GRANICY =
  'Platforma nie prowadzi kont ani profili osobowych — bramka jest jedna, '
  + 'a Operator bezimienny. Nie ma tu „Wyloguj": po zalogowaniu nie ma bram.';

export interface MenuProfilu {
  /** Uchwyt wraz z menu; montowany w grupie akcji paska. */
  element: HTMLElement;
  /** Zwija menu i zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu paska. */
  zamknij(): void;
}

export interface OpcjeMenuProfilu {
  /** Podpis Operatora — źródło inicjałów na uchwycie. */
  operator: string;
  /**
   * Czynności paska. Czytane w chwili WYBORU pozycji, nie przy budowie menu —
   * zaczep bywa podłączony po montażu powłoki (`zaczepy-paska.ts`).
   */
  zaczepy?: ZaczepyPaska;
  /** Gdzie postawić zdanie, gdy drogi nie podano. Bez odbiorcy zdanie przepada. */
  naKomunikat?: OdbiorcaZdania;
}

export function utworzMenuProfilu(opcje: OpcjeMenuProfilu): MenuProfilu {
  const pozycjaKonfiguracji = POZYCJE_USTAWIEN.find((p) => p.kod === KOD_STOPKI);

  const menu: MenuDrzewo = utworzMenuDrzewo({
    nastawa: 'Operator',
    // Pole szukania nad ośmioma pozycjami zabrałoby wiersz i nie skróciło
    // ani jednego ruchu. Próg biblioteki (12 liści) i tak by go nie postawił;
    // zapis jawny mówi, że to wybór, a nie przypadek.
    progSzukania: Number.POSITIVE_INFINITY,
    naWybor: (klucz) => wykonaj(klucz),
    ...(pozycjaKonfiguracji === undefined
      ? {}
      : {
          stopka: {
            nazwa: pozycjaKonfiguracji.nazwa,
            opis: `${pozycjaKonfiguracji.wyjasnienie}. ${ZDANIE_GRANICY}`,
            ikona: pozycjaKonfiguracji.ikona,
            wykonaj: () => otworz(pozycjaKonfiguracji),
          },
        }),
  });

  menu.element.classList.add('dn-pasek-gorny__profil');

  /** Skutek wyboru pozycji ustawień — albo prawda o braku drogi. */
  function otworz(pozycja: PozycjaUstawienia): void {
    wykonajZaczep(opcje.zaczepy ?? {}, pozycja, opcje.naKomunikat);
  }

  function wykonaj(klucz: string): void {
    if (klucz === KLUCZ_MOTYWU) {
      // Przełączenie rozgłasza ZDARZENIE_MOTYWU, a nasłuch niżej przerysowuje
      // menu. Stanu motywu ten plik nie trzyma — jedna prawda o nastawie.
      przelaczMotyw();
      return;
    }
    const pozycja = POZYCJE_USTAWIEN.find((p) => p.kod === klucz);
    if (pozycja === undefined) return;
    otworz(pozycja);
  }

  function lisc(pozycja: PozycjaUstawienia): PozycjaMenu {
    return {
      rodzaj: 'wybor',
      klucz: pozycja.kod,
      nazwa: pozycja.nazwa,
      opis: pozycja.wyjasnienie,
      ikona: pozycja.ikona,
      // Żadna pozycja nie jest nastawą jednokrotną — wszystkie otwierają okno.
      // Haczyk zostaje pusty i to jest prawda: menu nie pamięta „ostatnio
      // otwartego okna" i nie ma czego zaznaczyć.
      wybrany: false,
    };
  }

  function drzewo(): PozycjaMenu[] {
    const grupy: PozycjaMenu[] = [];
    for (const grupa of GRUPY) {
      const dzieci = grupa.kody
        .map((kod) => POZYCJE_USTAWIEN.find((p) => p.kod === kod))
        .filter((p): p is PozycjaUstawienia => p !== undefined)
        .map(lisc);
      if (dzieci.length === 0) continue;
      grupy.push({ rodzaj: 'grupa', nazwa: grupa.nazwa, dzieci });
    }

    // Grupa „Widok" stoi między tożsamością a obecnością i ma dokładnie jedną
    // pozycję — nastawę, nie wejście do okna. Dlatego składa się osobno.
    const ciemny = motywObowiazujacy() === 'dark';
    grupy.splice(1, 0, {
      rodzaj: 'grupa',
      nazwa: 'Widok',
      dzieci: [
        {
          rodzaj: 'przelacznik',
          klucz: KLUCZ_MOTYWU,
          nazwa: 'Motyw ciemny',
          opis: ciemny
            ? 'Włączony. Wyłączenie wraca do motywu jasnego — ta sama nastawa, '
              + 'co przycisk motywu w pasku.'
            : 'Wyłączony — obowiązuje motyw jasny. Ta sama nastawa, co przycisk '
              + 'motywu w pasku.',
          ikona: ciemny ? 'ksiezyc' : 'slonce',
          wlaczony: ciemny,
        },
      ],
    });

    return grupy;
  }

  function odrysuj(): void {
    menu.ustaw(inicjaly(opcje.operator), drzewo());
  }

  function naZmianeMotywu(): void {
    odrysuj();
  }

  document.addEventListener(ZDARZENIE_MOTYWU, naZmianeMotywu);
  odrysuj();

  return {
    element: menu.element,

    zamknij() {
      document.removeEventListener(ZDARZENIE_MOTYWU, naZmianeMotywu);
      menu.zwin();
    },
  };
}

/**
 * Do dwóch pierwszych liter członów nazwy Operatora. Pusty wynik zastępuje
 * `ZNAK_ZASTEPCZY`.
 */
function inicjaly(operator: string): string {
  const znaki = operator
    .split(/\s+/u)
    .filter((czlon) => czlon.length > 0)
    .slice(0, 2)
    .map((czlon) => czlon.charAt(0).toLocaleUpperCase('pl-PL'))
    .join('');
  return znaki === '' ? ZNAK_ZASTEPCZY : znaki;
}
