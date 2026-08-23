import {
  utworzMenuDrzewo,
  type MenuDrzewo,
  type PozycjaMenu,
} from '../../komponenty/menu-drzewo';
import { GRUPY_NARZEDZI, narzedziaGrupy } from './katalog-narzedzi';

/**
 * Ster wyboru narzędzi — wykaz pogrupowany po przeznaczeniu, gałąź to obszar.
 *
 * Komponent jest sterem, nie wyświetlaczem: etykieta uchwytu niesie wskazany
 * kod, kliknięcie rozwija wybór, wybór zmienia nastawę. Nastawa jest jedna
 * i dzieli ją z polem otwartym Skills Managera — wpisanie kodu ręcznie
 * przestawia uchwyt, wskazanie w drzewie wypełnia pole.
 *
 * Pole otwarte zostaje obok drzewa: `Agent.skillIds` jest w kontrakcie listą
 * napisów bez narzuconego słownika, a serwer narzędzi rozpoznaje też kody,
 * których katalog kontraktu nie zna (melduje je jako nierozpoznane). Zamknięcie
 * wyboru do samego drzewa odcięłoby drogę do takiego kodu.
 *
 * Liście są dwojakie, bo `ekspert_wykaz.go` rozpoznaje kod jako nazwę narzędzia
 * albo nazwę grupy: pierwszy liść każdej gałęzi wskazuje cały obszar jednym
 * kodem, pozostałe — pojedyncze narzędzia. Obszar jest jednostką doboru,
 * bo pozycji jest ponad dwieście.
 *
 * Próg, od którego menu stawia pole szukania, zostaje domyślny — przy tylu
 * liściach pole stanie zawsze, więc własna liczba niczego by nie zmieniła.
 *
 * Gałąź zwinięta nie jest brakiem: widoczne są nazwy obszarów wraz z liczbą
 * narzędzi, a pozycje odsłania dopiero wejście w obszar.
 */

/** Napis na uchwycie, gdy nastawa jest pusta — uchwyt zawsze niesie zdanie. */
const BEZ_WSKAZANIA = 'wskaż narzędzie albo obszar';

/** Klucz gałęzi obszaru; osobny od kodu grupy, bo gałąź nie jest wyborem. */
function kluczGalezi(grupa: string): string {
  return `obszar:${grupa}`;
}

/** Odmiana rzeczownika „narzędzie" dla liczby stojącej przy nim. */
function odmianaNarzedzi(ile: number): string {
  if (ile === 1) return 'narzędzie';
  const dziesiatki = ile % 100;
  const jednosci = ile % 10;
  if (dziesiatki >= 12 && dziesiatki <= 14) return 'narzędzi';
  return jednosci >= 2 && jednosci <= 4 ? 'narzędzia' : 'narzędzi';
}

/**
 * Składa drzewo wyboru: gałąź na obszar, w niej kod obszaru i jego narzędzia.
 *
 * Wydzielone z widoku, żeby dało się sprawdzić bez montowania menu
 * w dokumencie.
 *
 * `przypisane` znakuje pozycje, które ekspert już ma. Znakowanie idzie zdaniem
 * w opisie, a nie samym `wybrany`: `wybrany` niesie nastawę tej kontrolki
 * (jeden kod wskazany do przypisania), a przypisanie jest czymś innym — stanem
 * eksperta. Zlanie obu w jeden znacznik kazałoby uchwytowi pokazywać naraz
 * kilkanaście wartości, czyli przestać być sterem nastawy.
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
  /**
   * Podaje kod z zewnątrz — droga dla pola otwartego.
   *
   * Nie ogłasza zmiany słuchaczom: wołający JEST źródłem tej zmiany i odbicie
   * wróciłoby do niego pętlą.
   */
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
