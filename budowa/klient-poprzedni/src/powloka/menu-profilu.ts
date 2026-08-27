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

/** Klucz przełącznika motywu w menu profilu jest odrębny od kodów pozycji ustawień i nie koliduje z żadnym z nich w drzewie menu. */
export const KLUCZ_MOTYWU = 'widok:motyw-ciemny';

/** Napis widoczny na uchwycie menu, gdy imię i nazwisko Operatora nie dostarczają ani jednej litery do zbudowania inicjałów. */
const ZNAK_ZASTEPCZY = 'O';

/** Kod pozycji ustawień, która trafia do stopki menu profilu jako droga prowadząca do pełnego rejestru ustawień aplikacji. */
const KOD_STOPKI: KodUstawienia = 'konfiguracja';

/** Podział pozycji ustawień na grupy wyświetlane w menu profilu. Kolejność grup oraz kolejność pozycji wewnątrz każdej grupy odpowiada porządkowi zapisanemu w tej strukturze. */
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

/** Zdanie graniczne wyjaśniające, dlaczego menu profilu nie zawiera konta ani polecenia wylogowania z aplikacji. */
const ZDANIE_GRANICY =
  'Platforma nie prowadzi kont ani profili osobowych — bramka jest jedna, '
  + 'a Operator bezimienny. Nie ma tu „Wyloguj": po zalogowaniu nie ma bram.';

/** Uchwyt menu profilu Operatora wraz z metodą jego zamknięcia, zwracany po zbudowaniu menu i montowany w pasku powłoki środowiska. */
export interface MenuProfilu {
  /** Uchwyt wraz z menu; montowany w grupie akcji paska. */
  element: HTMLElement;
  /** Zwija menu i zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu paska. */
  zamknij(): void;
}

/** Opcje przekazywane przy budowie menu profilu — podpis Operatora, czynności paska sterujące wyborem oraz miejsce, gdzie trafia zdanie komunikatu. */
export interface OpcjeMenuProfilu {
  /** Podpis Operatora — źródło inicjałów na uchwycie. */
  operator: string;
  /** Czynności paska czytane przy wyborze pozycji, bo zaczep bywa podłączony dopiero po montażu powłoki. */
  zaczepy?: ZaczepyPaska;
  /** Gdzie postawić zdanie, gdy drogi nie podano. Bez odbiorcy zdanie przepada. */
  naKomunikat?: OdbiorcaZdania;
}

/** Buduje menu profilu Operatora w pasku powłoki wraz z obsługą wyboru pozycji, przełącznika motywu oraz zamknięcia po zdjęciu nasłuchów. */
export function utworzMenuProfilu(opcje: OpcjeMenuProfilu): MenuProfilu {
  const pozycjaKonfiguracji = POZYCJE_USTAWIEN.find((p) => p.kod === KOD_STOPKI);

  const menu: MenuDrzewo = utworzMenuDrzewo({
    nastawa: 'Operator',
    // Pole szukania jest zbędne przy tej liczbie pozycji: próg biblioteki i tak by go nie pokazał.
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
      // Przełączenie rozgłasza zdarzenie motywu, a nasłuch niżej przerysowuje menu.
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
      // Żadna pozycja nie jest nastawą jednokrotną, więc znacznik wyboru zostaje pusty.
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

    // Grupa „Widok" ma jedną nastawę zamiast wejścia do okna, więc składa się osobno.
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
 * Wyznacza inicjały z dwóch pierwszych liter kolejnych członów nazwy Operatora. Pusty wynik
 * zastępuje stały znak zastępczy przechowywany w stałej modułu.
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
