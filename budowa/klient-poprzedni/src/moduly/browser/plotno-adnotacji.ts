import {
  krojNapisu,
  narysujSlady,
  type NarzedzieAdnotacji,
  type PisakPlotna,
  type Punkt,
  type Slad,
} from './slady-adnotacji';

/**
 * Płótno adnotacji — `<canvas>` położony nad sceną podglądu strony.
 *
 * Jedna odpowiedzialność: element płótna, gesty wskaźnika i przerysowanie
 * wykazu śladów. Pasek narzędzi jest osobno (`pasek-adnotacji.ts`), a złożenie
 * jednego z drugim w warstwę nad podglądem — w `warstwa-adnotacji.ts`.
 *
 * DOM podglądu pozostaje nietknięty: płótno nie dopisuje ani jednego węzła do
 * `mb-podglad`, nie zmienia mu klas i niczego z niego nie czyta — leży nad nim.
 *
 * Tło jest przezroczyste, bo zrzutu strony nie ma: rdzeń pobiera stronę
 * biblioteką HTTP i zostawia `screenshotRef` pusty
 * (`przegladarka_pobieranie.go`). Rysunek powstaje więc po tym, co pokazuje
 * scena, i nie utrwala tła.
 *
 * `getContext('2d')` oddaje `null` w środowisku bez rasteryzacji — wtedy nie ma
 * czego wykreślić, ale wykaz śladów, wybór narzędzia i czyszczenie działają
 * nietknięte.
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
  /**
   * Spłaszcza płótno do PNG (`canvas.toDataURL`). Pusty napis znaczy, że
   * przeglądarka nie oddała obrazu — okno ma o tym powiedzieć, nie wysłać pustkę.
   */
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

  // Kontekst czytany raz: `getContext` przy każdym pociągnięciu kosztuje bez
  // powodu, a jego brak jest stanem trwałym środowiska, nie chwilowym.
  const pisak = element.getContext('2d') as PisakPlotna | null;

  /**
   * Barwa rozwiązana z żetonu motywu w chwili rysowania.
   *
   * Płótno 2D nie zna `var(--dn-…)`, więc żeton trzeba rozwiązać z wyliczonego
   * stylu elementu. Gdy arkusz nie jest wczytany (sprawdzian bez stylów),
   * wracamy do koloru tekstu elementu, a nie do barwy zapisanej wprost — żadna
   * nie ma prawa paść w kodzie.
   */
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
    // Ramka o zerowym boku daje dzielenie przez zero; jedynka zostawia
    // współrzędną nietkniętą zamiast wpisać do śladu `NaN`.
    return {
      x: (zdarzenie.clientX - ramka.left) / (ramka.width === 0 ? 1 : ramka.width),
      y: (zdarzenie.clientY - ramka.top) / (ramka.height === 0 ? 1 : ramka.height),
    };
  }

  function zacznij(zdarzenie: PointerEvent): void {
    const poczatek = punkt(zdarzenie);
    if (narzedzie === 'tekst') {
      // Napis powstaje jednym naciśnięciem, bez ciągnięcia: treść bierze się
      // z pola paska, więc pociągnięcie nie miałoby czego zbierać.
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
    // Ołówek zbiera całą trasę, figura tylko dwa naroża — dlatego kolejny
    // punkt raz się dokłada, a raz podmienia drugi.
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

  // Zmiana rozmiaru zeruje bufor płótna — obraz odtwarza się z wykazu śladów.
  // Obserwator bywa nieobecny (środowisko sprawdzianu), a jego brak nie jest
  // powodem, by tryb adnotacji nie działał.
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
