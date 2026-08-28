import {
  krojNapisu,
  narysujSlady,
  type NarzedzieAdnotacji,
  type PisakPlotna,
  type Punkt,
  type Slad,
} from './slady-adnotacji';

/**
 * Płótno adnotacji — element rysunkowy położony nad sceną podglądu strony.
 * Jedna odpowiedzialność: element płótna, gesty wskaźnika i przerysowanie
 * wykazu śladów; DOM podglądu pozostaje przy tym całkiem nietknięty.
 */
export interface PlotnoAdnotacji {
  /** Element `<canvas>` osadzany w warstwie nad sceną. */
  element: HTMLCanvasElement;
  narzedzie(): NarzedzieAdnotacji;
  ustawNarzedzie(nowe: NarzedzieAdnotacji): void;
  /** Nazwa żetonu barwy bieżącej — rozwiązanie należy do chwili rysowania. */
  zeton(): string;
  ustawZeton(nowy: string): void;
  /** Treść stawiana narzędziem „tekst"; puste nie stawia nic. */
  ustawNapis(tresc: string): void;
  /** Opróżnia wykaz śladów i przerysowuje — „Wyczyść adnotacje". */
  wyczysc(): void;
  /** Ile śladów niesie rysunek; zero znaczy „nie ma czego spłaszczać". */
  liczbaSladow(): number;
  /** Spłaszcza płótno do PNG; napis pusty znaczy, że przeglądarka nie oddała obrazu. */
  doPng(): string;
  /** Dopasowuje bufor do ramki i odtwarza rysunek z wykazu śladów. */
  dopasuj(): void;
  /** Odpina obserwatora rozmiaru — warstwa kończy pracę. */
  rozlacz(): void;
}

export function utworzPlotnoAdnotacji(zetonPoczatkowy: string): PlotnoAdnotacji {
  const element = document.createElement('canvas');
  element.className = 'mb-adnotacja__plotno';
  element.setAttribute('aria-label', 'Płótno adnotacji nad podglądem strony');

  const slady: Slad[] = [];
  let narzedzie: NarzedzieAdnotacji = 'olowek';
  let zeton = zetonPoczatkowy;
  let napis = '';
  let biezacy: Slad | null = null;

  // Kontekst czytany raz — jego brak jest stanem trwałym środowiska, nie chwilowym.
  const pisak = element.getContext('2d') as PisakPlotna | null;

  // Barwa rozwiązana z żetonu motywu w chwili rysowania, z wyliczonego stylu elementu płótna.
  function barwaZetonu(nazwa: string): string {
    const styl = getComputedStyle(element);
    const wartosc = styl.getPropertyValue(nazwa).trim();
    return wartosc === '' ? styl.color : wartosc;
  }

  function odrysuj(): void {
    if (pisak === null) return;
    narysujSlady(
      pisak,
      biezacy === null ? slady : [...slady, biezacy],
      { szerokosc: element.width, wysokosc: element.height },
      barwaZetonu,
      krojNapisu(getComputedStyle(element).fontFamily),
    );
  }

  /** Punkt zdarzenia przełożony na ułamek ramki płótna. */
  function punkt(zdarzenie: PointerEvent): Punkt {
    const ramka = element.getBoundingClientRect();
    // Ramka o zerowym boku daje dzielenie przez zero; jedynka zostawia współrzędną nietkniętą.
    return {
      x: (zdarzenie.clientX - ramka.left) / (ramka.width === 0 ? 1 : ramka.width),
      y: (zdarzenie.clientY - ramka.top) / (ramka.height === 0 ? 1 : ramka.height),
    };
  }

  function zacznij(zdarzenie: PointerEvent): void {
    const poczatek = punkt(zdarzenie);
    if (narzedzie === 'tekst') {
      // Napis powstaje jednym naciśnięciem, bez ciągnięcia — treść bierze się z pola paska.
      if (napis !== '') {
        slady.push({ narzedzie, zeton, punkty: [poczatek], tekst: napis });
        odrysuj();
      }
      return;
    }
    biezacy = { narzedzie, zeton, punkty: [poczatek], tekst: '' };
    if (typeof element.setPointerCapture === 'function') {
      element.setPointerCapture(zdarzenie.pointerId);
    }
    odrysuj();
  }

  function prowadz(zdarzenie: PointerEvent): void {
    if (biezacy === null) return;
    const kolejny = punkt(zdarzenie);
    // Ołówek zbiera całą trasę, figura tylko dwa naroża — punkt raz się dokłada, a raz podmienia drugi.
    if (biezacy.narzedzie === 'olowek') biezacy.punkty.push(kolejny);
    else biezacy.punkty = [biezacy.punkty[0] ?? kolejny, kolejny];
    odrysuj();
  }

  function skoncz(): void {
    if (biezacy === null) return;
    slady.push(biezacy);
    biezacy = null;
    odrysuj();
  }

  element.addEventListener('pointerdown', zacznij);
  element.addEventListener('pointermove', prowadz);
  element.addEventListener('pointerup', skoncz);
  element.addEventListener('pointercancel', skoncz);
  element.addEventListener('pointerleave', skoncz);

  // Zmiana rozmiaru zeruje bufor; brak obserwatora w środowisku sprawdzianu nie blokuje trybu adnotacji.
  const obserwator =
    typeof ResizeObserver === 'function' ? new ResizeObserver(() => dopasuj()) : null;
  obserwator?.observe(element);

  function dopasuj(): void {
    const ramka = element.getBoundingClientRect();
    const szerokosc = Math.max(Math.round(ramka.width), 1);
    const wysokosc = Math.max(Math.round(ramka.height), 1);
    if (element.width !== szerokosc) element.width = szerokosc;
    if (element.height !== wysokosc) element.height = wysokosc;
    odrysuj();
  }

  return {
    element,
    narzedzie: () => narzedzie,
    ustawNarzedzie(nowe) {
      narzedzie = nowe;
    },
    zeton: () => zeton,
    ustawZeton(nowy) {
      zeton = nowy;
    },
    ustawNapis(tresc) {
      napis = tresc;
    },
    wyczysc() {
      slady.splice(0, slady.length);
      biezacy = null;
      odrysuj();
    },
    liczbaSladow: () => slady.length,
    doPng: () => element.toDataURL('image/png'),
    dopasuj,
    rozlacz() {
      obserwator?.disconnect();
    },
  };
}
