import {
  utworzMenuDrzewo,
  type MenuDrzewo,
  type PozycjaMenu,
} from '../../komponenty/menu-drzewo';
import { GRUPY_NARZEDZI, narzedziaGrupy } from './katalog-narzedzi';

/**
 * Ster wyboru narzędzi w postaci wykazu pogrupowanego po przeznaczeniu, gdzie
 * gałąź oznacza obszar, a etykieta uchwytu niesie wskazany kod.
 */
const BEZ_WSKAZANIA = 'wskaż narzędzie albo obszar';

/** Klucz gałęzi obszaru, osobny od kodu grupy, ponieważ gałąź drzewa nie jest samodzielną pozycją do wyboru w tej kontrolce. */
function kluczGalezi(grupa: string): string {
  return `obszar:${grupa}`;
}

/** Odmiana rzeczownika „narzędzie” dobierana zależnie od liczby stojącej przy nim w opisie liczby pozycji obszaru. */
function odmianaNarzedzi(ile: number): string {
  if (ile === 1) return 'narzędzie';
  const dziesiatki = ile % 100;
  const jednosci = ile % 10;
  if (dziesiatki >= 12 && dziesiatki <= 14) return 'narzędzi';
  return jednosci >= 2 && jednosci <= 4 ? 'narzędzia' : 'narzędzi';
}

/**
 * Składa drzewo wyboru: gałąź na obszar, w niej kod obszaru i jego narzędzia,
 * wydzielone z widoku, aby dało się sprawdzić bez montowania menu
 * w dokumencie.
 */
export function zbudujDrzewoNarzedzi(
  wybrany: string,
  przypisane: readonly string[],
): PozycjaMenu[] {
  const juzMa = new Set(przypisane);

  return GRUPY_NARZEDZI.map((grupa) => {
    const pozycje = narzedziaGrupy(grupa);
    const przypisanychTu = pozycje.filter((pozycja) => juzMa.has(pozycja.nazwa)).length;
    const calyObszar = juzMa.has(grupa);

    const dzieci: PozycjaMenu[] = [
      {
        rodzaj: 'wybor',
        klucz: grupa,
        nazwa: grupa,
        opis:
          `Cały obszar jednym kodem — ekspert dostanie wszystkie ${pozycje.length} ` +
          `${odmianaNarzedzi(pozycje.length)} tego obszaru.` +
          (calyObszar ? ' Ten kod ekspert już ma.' : ''),
        wybrany: wybrany === grupa,
      },
      ...pozycje.map(
        (pozycja): PozycjaMenu => ({
          rodzaj: 'wybor',
          klucz: pozycja.nazwa,
          nazwa: pozycja.nazwa,
          // Druga nazwa to komenda kontraktu — ta sama rzecz w notacji z wykazu
          // komend. Szukanie obejmuje obie.
          nazwaKrotka: pozycja.komenda,
          opis:
            (juzMa.has(pozycja.nazwa) ? 'Już przypisane. ' : '') +
            (calyObszar && !juzMa.has(pozycja.nazwa)
              ? `Wchodzi w kodzie obszaru ${grupa}, który ekspert już ma. `
              : '') +
            pozycja.opis,
          wybrany: wybrany === pozycja.nazwa,
        }),
      ),
    ];

    return {
      rodzaj: 'galaz',
      klucz: kluczGalezi(grupa),
      nazwa: grupa,
      opis:
        `${pozycje.length} ${odmianaNarzedzi(pozycje.length)}` +
        (przypisanychTu > 0 ? ` · przypisanych: ${przypisanychTu}` : '') +
        (calyObszar ? ' · obszar przypisany w całości' : ''),
      dzieci,
    };
  });
}

export interface SterNarzedzi {
  /** Element montowany w oknie obok pola otwartego. */
  element: HTMLElement;
  /** Kod wskazany w tej chwili — nazwa narzędzia albo nazwa obszaru. */
  kod(): string;
  /** Podaje kod z zewnątrz — droga dla pola otwartego, bez ogłaszania zmiany słuchaczom. */
  ustawKod(kod: string): void;
  /** Nanosi kody, które ekspert już ma — drzewo je znakuje. */
  ustawPrzypisane(kody: readonly string[]): void;
  /** Nasłuch wskazania w drzewie; woła się tylko przy zmianie na inny kod. */
  naZmiane(sluchacz: (kod: string) => void): void;
  /** Zwija menu wraz z gałęziami i zdejmuje nasłuchy dokumentu. */
  zwin(): void;
}

export function utworzSterNarzedzi(): SterNarzedzi {
  const sluchacze = new Set<(kod: string) => void>();
  let wskazany = '';
  let przypisane: readonly string[] = [];

  const menu: MenuDrzewo = utworzMenuDrzewo({
    nastawa: 'Narzędzie albo obszar do przypisania',
    naWybor: (klucz) => {
      if (klucz === wskazany) return;
      wskazany = klucz;
      odrysuj();
      for (const sluchacz of [...sluchacze]) sluchacz(wskazany);
    },
  });

  function odrysuj(): void {
    menu.ustaw(
      wskazany === '' ? BEZ_WSKAZANIA : wskazany,
      zbudujDrzewoNarzedzi(wskazany, przypisane),
    );
  }

  odrysuj();

  return {
    element: menu.element,

    kod: () => wskazany,

    ustawKod(kod) {
      if (kod === wskazany) return;
      wskazany = kod;
      odrysuj();
    },

    ustawPrzypisane(kody) {
      przypisane = [...kody];
      odrysuj();
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
    },

    zwin: () => menu.zwin(),
  };
}
