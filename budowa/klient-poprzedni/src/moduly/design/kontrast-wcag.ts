import progi from '../../../../shared/kontrasty-progi.json';
import { wartoscZetonu } from './zetony-systemu';

// Rachunek kontrastu par żetonów wedle WCAG 2.1, liczony na barwach z motywu obowiązującego

/** Wynik pomiaru jednej pary żetonów wraz z wymaganym progiem oraz zdaniem o powodzie, gdy pomiar się nie udał. */
export interface PomiarPary {
  /** Nazwa pary z wykazu progów. */
  readonly para: string;
  /** Żetony pary: pierwszy na drugim. */
  readonly zetony: readonly string[];
  /** Próg wymagany przez wykaz. */
  readonly prog: number;
  /** Współczynnik kontrastu; `null`, gdy którejś barwy nie dało się odczytać. */
  readonly wspolczynnik: number | null;
  /** Czy pomiar sięga progu; `false` przy braku odczytu. */
  readonly spelnia: boolean;
  /** Zdanie o powodzie braku pomiaru; puste, gdy pomiar się udał. */
  readonly powod: string;
}

/** Składowe barwy zapisanej w modelu sRGB, każda z liczbowego zakresu od zera do dwustu pięćdziesięciu pięciu. */
export interface Barwa {
  r: number;
  g: number;
  b: number;
}

/** Kształt wykazu progów kontrastu prowadzonego przez produkt — tylko te pola, których moduł rzeczywiście używa. */
interface WykazProgow {
  readonly prog_domyslny: number;
  readonly pomiary: readonly { para: string; zetony: string[]; prog: number }[];
}

const WYKAZ = progi as WykazProgow;

/** Próg domyślny kontrastu wzięty wprost z wykazu progów prowadzonego przez produkt, nigdy z pamięci ani domysłu. */
export const PROG_DOMYSLNY = WYKAZ.prog_domyslny;

/** Liczba par kontrastu wymienionych w wykazie progów prowadzonym przez sam produkt, a nie przez ten moduł. */
export function liczbaPar(): number {
  return WYKAZ.pomiary.length;
}

/**
 * Barwa z zapisu CSS. Obsługuje zapis szesnastkowy trzy- i sześcioznakowy oraz zapis `rgb`/`rgba`.
 * Zapis, którego nie umiemy odczytać, daje `null`, a nie czerń, bo podstawienie czerni
 * zafałszowałoby całą tabelę.
 */
export function odczytajBarwe(zapis: string): Barwa | null {
  const tekst = zapis.trim().toLowerCase();
  if (tekst === '') return null;

  if (tekst.startsWith('#')) {
    const znaki = tekst.slice(1);
    if (znaki.length === 3) {
      const [r, g, b] = [...znaki].map((znak) => Number.parseInt(`${znak}${znak}`, 16));
      return zlozBarwe(r, g, b);
    }
    if (znaki.length === 6 || znaki.length === 8) {
      return zlozBarwe(
        Number.parseInt(znaki.slice(0, 2), 16),
        Number.parseInt(znaki.slice(2, 4), 16),
        Number.parseInt(znaki.slice(4, 6), 16),
      );
    }
    return null;
  }

  const skladowe = /^rgba?\(([^)]+)\)$/.exec(tekst);
  if (skladowe === null) return null;
  const liczby = (skladowe[1] ?? '')
    .split(/[\s,/]+/)
    .filter((czesc) => czesc !== '')
    .map((czesc) => Number.parseFloat(czesc));
  return zlozBarwe(liczby[0], liczby[1], liczby[2]);
}

function zlozBarwe(r?: number, g?: number, b?: number): Barwa | null {
  if (r === undefined || g === undefined || b === undefined) return null;
  if (!Number.isFinite(r) || !Number.isFinite(g) || !Number.isFinite(b)) return null;
  return { r, g, b };
}

/** Luminancja względna barwy wedle normy WCAG 2.1, kryterium 1.4.3, liczona z jej trzech składowych modelu sRGB. */
export function luminancja(barwa: Barwa): number {
  const skladowe = [barwa.r, barwa.g, barwa.b].map((wartosc) => {
    const udzial = wartosc / 255;
    return udzial <= 0.03928 ? udzial / 12.92 : ((udzial + 0.055) / 1.055) ** 2.4;
  });
  const [r, g, b] = skladowe as [number, number, number];
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

/** Współczynnik kontrastu dwóch barw wedle WCAG 2.1, jako wynik z przedziału od jeden do dwadzieścia jeden. */
export function wspolczynnikKontrastu(pierwsza: Barwa, druga: Barwa): number {
  const jasniejsza = Math.max(luminancja(pierwsza), luminancja(druga));
  const ciemniejsza = Math.min(luminancja(pierwsza), luminancja(druga));
  return (jasniejsza + 0.05) / (ciemniejsza + 0.05);
}

/**
 * Pomiar wszystkich par wykazu na motywie obowiązującym.
 *
 * Barwa nieodczytana nie wypada z tabeli — zostaje w niej z nazwanym powodem.
 * Para pominięta w wyniku wyglądałaby jak para zgodna.
 */
export function zmierzPary(): readonly PomiarPary[] {
  return WYKAZ.pomiary.map((wpis) => {
    const [pierwszyZeton, drugiZeton] = wpis.zetony;
    const pierwsza = pierwszyZeton === undefined ? null : odczytajBarwe(wartoscZetonu(pierwszyZeton));
    const druga = drugiZeton === undefined ? null : odczytajBarwe(wartoscZetonu(drugiZeton));

    if (pierwsza === null || druga === null) {
      const nieodczytane = [
        pierwsza === null ? pierwszyZeton : undefined,
        druga === null ? drugiZeton : undefined,
      ].filter((nazwa): nazwa is string => nazwa !== undefined);
      return {
        para: wpis.para,
        zetony: wpis.zetony,
        prog: wpis.prog,
        wspolczynnik: null,
        spelnia: false,
        powod:
          `Nie odczytano barwy żetonu: ${nieodczytane.join(', ')}. Motyw albo tego żetonu nie ` +
          'definiuje, albo oddaje go w zapisie, którego panel nie rozkłada na składowe — ' +
          'pomiaru nie zmyślamy.',
      };
    }

    const wspolczynnik = wspolczynnikKontrastu(pierwsza, druga);
    return {
      para: wpis.para,
      zetony: wpis.zetony,
      prog: wpis.prog,
      wspolczynnik,
      spelnia: wspolczynnik >= wpis.prog,
      powod: '',
    };
  });
}

/** Pary poniżej wymaganego progu kontrastu — sprawdzian dostępności kolorem osiągalny jednym wywołaniem funkcji. */
export function poniejProgu(pomiary: readonly PomiarPary[]): readonly PomiarPary[] {
  return pomiary.filter((pomiar) => !pomiar.spelnia);
}

/** Współczynnik przygotowany do pokazania: dwie cyfry po przecinku, bo w takiej postaci podaje go norma WCAG. */
export function zapisWspolczynnika(wspolczynnik: number | null): string {
  return wspolczynnik === null ? 'bez pomiaru' : `${wspolczynnik.toFixed(2)}:1`;
}
