import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { NAZWY_WARSTW, type CzynnoscOkna } from './czynnosci-okna';
import { KATALOG_FUNKCJI, LICZBA_POZYCJI, type PozycjaKatalogu } from './katalog-funkcji';

/**
 * Paleta poleceń modułu Terminal — jedno wskazanie do każdej czynności okien
 * i do każdej pozycji katalogu funkcji.
 *
 * Paleta jest drugą połową zasady warstw widoczności. Warstwy zdejmują
 * z ekranu to, czego bieżące zadanie nie wymaga; paleta pilnuje, żeby zdjęcie
 * z ekranu nie stało się schowaniem. Czynność warstwy eksperckiej, której
 * w panelu akcji nie widać, jest tu o jedno wskazanie i o jedno wpisane słowo.
 *
 * Wykaz ma dwie części, bo odpowiada na dwa różne pytania. Czynności okien
 * odpowiadają „co mogę teraz zrobić" i po wybraniu wykonują się od razu.
 * Katalog funkcji odpowiada „czy moduł to potrafi" i po wybraniu mówi, gdzie
 * funkcja stoi w tej budowie albo dlaczego jej nie ma oraz od jakiego programu
 * spoza instalki zależy. Pozycja katalogu niczego nie uruchamia i paleta tego
 * nie udaje.
 *
 * Mechanizm rozwijania, wędrówkę klawiaturą, pole szukania i znacznik wyboru
 * niesie `komponenty/menu-drzewo.ts`; tutaj powstają wyłącznie dane drzewa.
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

/** Wartość uchwytu, dopóki Operator niczego nie wywołał. */
const ZDANIE_BEZ_WYWOLANIA = 'Paleta poleceń — bez wywołania';

/** Przedrostki kluczy; klucz ma być niepowtarzalny w obrębie całego drzewa. */
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
    // Próg niższy od domyślnego: wykaz liczy kilkadziesiąt liści już przy
    // samym katalogu funkcji, więc pole szukania jest tu potrzebne od razu.
    progSzukania: 8,
    // Opis przy każdej z kilkudziesięciu pozycji byłby ścianą tekstu zamiast
    // objaśnienia; należy się tej pozycji, na którą Operator właśnie patrzy.
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

    /**
     * Otwarcie palety wraz z przeniesieniem ogniska w jej wnętrze.
     *
     * Mechanizm menu rozwija wykaz, nie ruszając ogniska — słusznie, bo jego
     * drugim odbiorcą jest wykaz otwierany ukośnikiem w polu wpisywania, gdzie
     * ognisko ma zostać u wołającego. Paleta wywoływana skrótem ma wymóg
     * odwrotny: Operator naciska skrót po to, żeby zacząć pisać nazwę czynności.
     * Ognisko idzie więc na pierwszą kontrolkę wnętrza — pole szukania, gdy stoi
     * na ekranie, w przeciwnym razie uchwyt.
     */
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

/** Klucz pozycji czynności — okno i nazwa razem, bo nazwy powtarzają się między oknami. */
function kluczCzynnosci(czynnosc: CzynnoscOkna): string {
  return `${KLUCZ_CZYNNOSCI}${czynnosc.okno}:${czynnosc.nazwa}`;
}

/** Drzewo palety: czynności okien w gałęziach po oknie, katalog funkcji po rodzinie. */
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

/** Objaśnienie czynności wraz z jej warstwą i skrótem. */
function opisCzynnosci(czynnosc: CzynnoscOkna): string {
  const czesci = [czynnosc.opis, `Warstwa: ${NAZWY_WARSTW[czynnosc.warstwa]}.`];
  if (czynnosc.skrot !== undefined) czesci.push(`Skrót: ${czynnosc.skrot}.`);
  return czesci.join(' ');
}

/** Zdanie o pozycji katalogu: co robi, gdzie stoi i od czego zależy. */
function zdanieOFunkcji(pozycja: PozycjaKatalogu): string {
  const czesci = [pozycja.coRobi, pozycja.stan];
  if (pozycja.zaleznosc !== '') {
    czesci.push(
      `Zależność od programów spoza instalki Danaco Console: ${pozycja.zaleznosc}.`,
    );
  }
  return czesci.join(' ');
}
