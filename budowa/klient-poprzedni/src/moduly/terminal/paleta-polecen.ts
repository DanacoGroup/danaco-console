import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { NAZWY_WARSTW, type CzynnoscOkna } from './czynnosci-okna';
import { KATALOG_FUNKCJI, LICZBA_POZYCJI, type PozycjaKatalogu } from './katalog-funkcji';

/**
 * Paleta poleceń modułu Terminal daje jedno wskazanie do każdej czynności okien i do każdej pozycji katalogu funkcji, uzupełniając warstwy widoczności o dostęp spoza panelu akcji.
 */
export interface PaletaPolecen {
  /** Element osadzany w pasie modułu. */
  element: HTMLElement;
  /** Otwiera paletę — nośnik skrótu klawiszowego. */
  otworz(): void;
  /** Podaje palecie bieżący wykaz czynności okien. */
  ustawCzynnosci(czynnosci: readonly CzynnoscOkna[]): void;
  /** Zwija paletę wraz z nasłuchami dokumentu. */
  zwin(): void;
}

/**
 * Wartość uchwytu palety, dopóki operator jeszcze niczego z niej nie wywołał w bieżącej sesji terminala.
 */
const ZDANIE_BEZ_WYWOLANIA = 'Paleta poleceń — bez wywołania';

/**
 * Przedrostki kluczy pozycji palety; każdy klucz ma pozostać niepowtarzalny w obrębie całego drzewa wykazu.
 */
const KLUCZ_CZYNNOSCI = 'czynnosc:';
const KLUCZ_FUNKCJI = 'funkcja:';

export function utworzPaletePolecen(): PaletaPolecen {
  let czynnosci: readonly CzynnoscOkna[] = [];
  let uchwyt = ZDANIE_BEZ_WYWOLANIA;

  const element = document.createElement('div');
  element.className = 'dt-paleta';

  const objasnienie = document.createElement('p');
  objasnienie.className = 'dn-pole-opis dt-paleta__objasnienie';
  objasnienie.textContent =
    `Paleta zna czynności okien modułu i ${LICZBA_POZYCJI} pozycji katalogu funkcji. ` +
    'Czynność wykonuje się po wybraniu; pozycja katalogu mówi, gdzie funkcja stoi w tej budowie.';

  const menu = utworzMenuDrzewo({
    nastawa: 'Paleta poleceń modułu Terminal',
    ikona: 'polecenie',
    // Próg wyszukiwania niższy od domyślnego: katalog funkcji liczy kilkadziesiąt pozycji.
    progSzukania: 8,
    // Opis przy każdej pozycji byłby ścianą tekstu; należy się tylko tej, na którą patrzy operator.
    opisTylkoPrzyWyroznionej: true,
    naWybor: (klucz) => wybierz(klucz),
  });

  element.append(menu.element, objasnienie);

  function odrysuj(): void {
    menu.ustaw(uchwyt, drzewoPalety(czynnosci));
  }

  function wybierz(klucz: string): void {
    if (klucz.startsWith(KLUCZ_CZYNNOSCI)) {
      const czynnosc = czynnosci.find((pozycja) => kluczCzynnosci(pozycja) === klucz);
      if (czynnosc === undefined) return;
      uchwyt = `${czynnosc.okno} — ${czynnosc.nazwa}`;
      odrysuj();
      czynnosc.wykonaj();
      return;
    }
    if (!klucz.startsWith(KLUCZ_FUNKCJI)) return;
    const numer = Number(klucz.slice(KLUCZ_FUNKCJI.length));
    const pozycja = KATALOG_FUNKCJI[numer];
    if (pozycja === undefined) return;
    uchwyt = `Katalog funkcji — ${pozycja.nazwa}`;
    objasnienie.textContent = zdanieOFunkcji(pozycja);
    odrysuj();
  }

  odrysuj();

  return {
    element,

    // Otwarcie palety z przeniesieniem ogniska w jej wnętrze; skrót stawia ognisko na polu szukania.
    otworz() {
      menu.rozwin();
      menu.element.querySelector<HTMLElement>('input, button')?.focus();
    },

    ustawCzynnosci(nowe) {
      czynnosci = nowe;
      odrysuj();
    },
    zwin: menu.zwin,
  };
}

/**
 * Klucz pozycji czynności łączy nazwę okna i nazwę czynności, bo same nazwy czynności powtarzają się między oknami.
 */
function kluczCzynnosci(czynnosc: CzynnoscOkna): string {
  return `${KLUCZ_CZYNNOSCI}${czynnosc.okno}:${czynnosc.nazwa}`;
}

/**
 * Drzewo palety układa czynności okien w gałęziach według okna, a katalog funkcji w gałęziach według rodziny funkcji.
 */
function drzewoPalety(czynnosci: readonly CzynnoscOkna[]): PozycjaMenu[] {
  const okna = new Map<string, CzynnoscOkna[]>();
  for (const czynnosc of czynnosci) {
    const pozycje = okna.get(czynnosc.okno) ?? [];
    pozycje.push(czynnosc);
    okna.set(czynnosc.okno, pozycje);
  }

  const rodziny = new Map<string, Array<[number, PozycjaKatalogu]>>();
  for (const [numer, pozycja] of KATALOG_FUNKCJI.entries()) {
    const wykaz = rodziny.get(pozycja.rodzina) ?? [];
    wykaz.push([numer, pozycja]);
    rodziny.set(pozycja.rodzina, wykaz);
  }

  return [
    {
      rodzaj: 'grupa',
      nazwa: 'Czynności okien modułu',
      dzieci: [...okna.entries()].map(([okno, pozycje]) => ({
        rodzaj: 'galaz' as const,
        klucz: `okno:${okno}`,
        nazwa: okno,
        opis: `${pozycje.length} czynności tego okna.`,
        ikona: 'terminal' as const,
        dzieci: pozycje.map((czynnosc) => ({
          rodzaj: 'wybor' as const,
          klucz: kluczCzynnosci(czynnosc),
          nazwa: czynnosc.nazwa,
          nazwaKrotka: czynnosc.nazwa,
          opis: opisCzynnosci(czynnosc),
          wybrany: false,
        })),
      })),
    },
    {
      rodzaj: 'grupa',
      nazwa: 'Katalog funkcji modułu',
      dzieci: [...rodziny.entries()].map(([rodzina, pozycje]) => ({
        rodzaj: 'galaz' as const,
        klucz: `rodzina:${rodzina}`,
        nazwa: rodzina,
        opis: `${pozycje.length} pozycji katalogu.`,
        ikona: 'biblioteka' as const,
        dzieci: pozycje.map(([numer, pozycja]) => ({
          rodzaj: 'wybor' as const,
          klucz: `${KLUCZ_FUNKCJI}${numer}`,
          nazwa: pozycja.nazwa,
          opis: zdanieOFunkcji(pozycja),
          wybrany: false,
        })),
      })),
    },
  ];
}

/**
 * Objaśnienie czynności palety wraz z nazwą jej warstwy widoczności oraz przypisanym skrótem klawiszowym.
 */
function opisCzynnosci(czynnosc: CzynnoscOkna): string {
  const czesci = [czynnosc.opis, `Warstwa: ${NAZWY_WARSTW[czynnosc.warstwa]}.`];
  if (czynnosc.skrot !== undefined) czesci.push(`Skrót: ${czynnosc.skrot}.`);
  return czesci.join(' ');
}

/**
 * Zdanie o pozycji katalogu funkcji: co funkcja robi, gdzie w budowie stoi i od jakiego programu zależy.
 */
function zdanieOFunkcji(pozycja: PozycjaKatalogu): string {
  const czesci = [pozycja.coRobi, pozycja.stan];
  if (pozycja.zaleznosc !== '') {
    czesci.push(
      `Zależność od programów spoza instalki Danaco Console: ${pozycja.zaleznosc}.`,
    );
  }
  return czesci.join(' ');
}
