import type { DesignAsset, DesignBoardLayer } from '../../../../shared/contract';
import { nowyIdentyfikator } from '../../protokol/identyfikator';

/** Zapis kompozycji Design Board wraz z czynnościami zmieniającymi układ warstw i zaznaczenie. */

/** Widok kanwy: powiększenie, przesunięcie swobodne oraz siatka pomocnicza widoczne w oknie modułu Design. */
export interface WidokKanwy {
  /** Powiększenie; 1 znaczy skalę naturalną. */
  powiekszenie: number;
  przesuniecieX: number;
  przesuniecieY: number;
  siatka: boolean;
}

/** Oś wyrównania warstw zaznaczonych na kanwie kompozycji: krawędź, środek albo przeciwna krawędź układu. */
export type OsWyrownania = 'lewo' | 'srodek' | 'prawo' | 'gora' | 'srodek-pion' | 'dol';

/**
 * Szablon układu z panelu akcji modułu Design.
 *
 * Cztery, bo tyle wymienia opracowanie: tablica nastroju, siatka porównawcza,
 * arkusz brandingowy i formaty społecznościowe.
 */
export type SzablonUkladu =
  | 'tablica-nastroju'
  | 'siatka-porownawcza'
  | 'arkusz-brandingowy'
  | 'formaty-spolecznosciowe';

/** Pełny zapis kompozycji — jedyny nośnik prawdy o układzie warstw, zaznaczeniu i widoku kanwy modułu Design. */
export interface ZapisKompozycji {
  warstwy: DesignBoardLayer[];
  wybor: string[];
  nazwa: string;
  identyfikator: string;
  widok: WidokKanwy;
  licznik: number;
}

/** Bok wyjściowy nowej warstwy dokładanej na kanwie, wyrażony w jednostkach kompozycji modułu Design okna. */
export const BOK_WYJSCIOWY = 240;

/** Odstęp w jednostkach kompozycji między warstwami układanymi automatycznie wybranym szablonem układu. */
const ODSTEP = 24;

export function pustyZapis(): ZapisKompozycji {
  return {
    warstwy: [],
    wybor: [],
    nazwa: '',
    identyfikator: '',
    widok: { powiekszenie: 1, przesuniecieX: 0, przesuniecieY: 0, siatka: true },
    licznik: 0,
  };
}

/**
 * Wstawia kompozycję odczytaną z rdzenia na miejsce układu bieżącego,
 * zastępując go pełnym stanem planszy przeniesionym z rdzenia.
 */
export function wczytajUklad(
  zapis: ZapisKompozycji,
  identyfikator: string,
  nazwa: string,
  warstwy: readonly DesignBoardLayer[],
): void {
  zapis.identyfikator = identyfikator;
  zapis.nazwa = nazwa;
  zapis.warstwy = [...warstwy];
  zapis.wybor = [];
  zapis.licznik = Math.max(zapis.licznik, warstwy.length);
}

/** Dokłada nową warstwę kompozycji zbudowaną z zasobu wybranego w panelu Assets Panel modułu Design okna. */
export function dolozZasob(zapis: ZapisKompozycji, zasob: DesignAsset): void {
  dolozWarstwe(zapis, {
    assetId: zasob.id,
    note: zasob.name ?? zasob.id,
    ...(zasob.width === undefined ? {} : { width: zasob.width }),
    ...(zasob.height === undefined ? {} : { height: zasob.height }),
  });
}

/** Dokłada warstwę pomocniczą bez zasobu, czyli element wzięty z biblioteki pomocniczej modułu Design okna. */
export function dolozElement(zapis: ZapisKompozycji, nazwa: string): void {
  dolozWarstwe(zapis, { note: nazwa });
}

/** Identyfikator warstwy — kod nadawany przez okno bytowi rdzenia, unikalny w tabeli warstw kompozycji. */
function nowyKodWarstwy(numer: number): string {
  return nowyIdentyfikator(`warstwa-${numer}`);
}

function dolozWarstwe(zapis: ZapisKompozycji, czesc: Partial<DesignBoardLayer>): void {
  zapis.licznik += 1;
  const kolejnosc = zapis.warstwy.length + 1;
  zapis.warstwy = [
    ...zapis.warstwy,
    {
      id: nowyKodWarstwy(zapis.licznik),
      x: ODSTEP * kolejnosc,
      y: ODSTEP * kolejnosc,
      width: BOK_WYJSCIOWY,
      height: BOK_WYJSCIOWY,
      order: kolejnosc,
      locked: false,
      ...czesc,
    },
  ];
}

/** Nanosi zmianę na jedną warstwę kompozycji wskazaną kodem, pozostawiając pozostałe warstwy bez zmian. */
export function zmienWarstwe(
  zapis: ZapisKompozycji,
  idWarstwy: string,
  zmiana: Partial<DesignBoardLayer>,
): void {
  zapis.warstwy = zapis.warstwy.map((w) => (w.id === idWarstwy ? { ...w, ...zmiana } : w));
}

export function usunWarstwe(zapis: ZapisKompozycji, idWarstwy: string): void {
  zapis.warstwy = zapis.warstwy.filter((w) => w.id !== idWarstwy);
  zapis.wybor = zapis.wybor.filter((kod) => kod !== idWarstwy);
}

/**
 * Przełącza zaznaczenie warstwy.
 *
 * `dolacz` odpowiada `Ctrl/Cmd + Klik`, czyli zaznaczeniu wielokrotnemu warstw
 * na kanwie Design Board.
 */
export function przestawZaznaczenie(
  zapis: ZapisKompozycji,
  idWarstwy: string,
  dolacz: boolean,
): void {
  if (dolacz) {
    zapis.wybor = zapis.wybor.includes(idWarstwy)
      ? zapis.wybor.filter((kod) => kod !== idWarstwy)
      : [...zapis.wybor, idWarstwy];
    return;
  }
  const juzSam = zapis.wybor.length === 1 && zapis.wybor[0] === idWarstwy;
  zapis.wybor = juzSam ? [] : [idWarstwy];
}

/** Warstwy objęte działaniem narzędzi kompozycji: zaznaczone na kanwie i przy tym niezablokowane w oknie. */
export function objeteWarstwy(zapis: ZapisKompozycji): DesignBoardLayer[] {
  return zapis.warstwy.filter((w) => zapis.wybor.includes(w.id) && w.locked !== true);
}

/** Nanosi wynik działania narzędzia na zapis kompozycji, zachowując kolejność warstw w wykazie zapisu kanwy. */
export function naniesWarstwy(
  zapis: ZapisKompozycji,
  zmienione: readonly DesignBoardLayer[],
): void {
  const wedlugKodu = new Map(zmienione.map((w) => [w.id, w]));
  zapis.warstwy = zapis.warstwy.map((w) => wedlugKodu.get(w.id) ?? w);
}

/** Wyrównanie warstw zaznaczonych do skrajnej albo środkowej pozycji całego zbioru wybranych warstw kanwy. */
export function wyrownajWarstwy(
  wybrane: readonly DesignBoardLayer[],
  os: OsWyrownania,
): readonly DesignBoardLayer[] {
  if (wybrane.length < 2) return wybrane;
  const lewe = wybrane.map((w) => w.x ?? 0);
  const gorne = wybrane.map((w) => w.y ?? 0);
  const prawe = wybrane.map((w) => (w.x ?? 0) + (w.width ?? BOK_WYJSCIOWY));
  const dolne = wybrane.map((w) => (w.y ?? 0) + (w.height ?? BOK_WYJSCIOWY));
  const srodekX = (Math.min(...lewe) + Math.max(...prawe)) / 2;
  const srodekY = (Math.min(...gorne) + Math.max(...dolne)) / 2;

  return wybrane.map((warstwa) => {
    const bokX = warstwa.width ?? BOK_WYJSCIOWY;
    const bokY = warstwa.height ?? BOK_WYJSCIOWY;
    if (os === 'lewo') return { ...warstwa, x: Math.min(...lewe) };
    if (os === 'prawo') return { ...warstwa, x: Math.max(...prawe) - bokX };
    if (os === 'srodek') return { ...warstwa, x: srodekX - bokX / 2 };
    if (os === 'gora') return { ...warstwa, y: Math.min(...gorne) };
    if (os === 'dol') return { ...warstwa, y: Math.max(...dolne) - bokY };
    return { ...warstwa, y: srodekY - bokY / 2 };
  });
}

/** Równe odstępy między warstwami zaznaczonymi na kanwie, rozłożone wzdłuż jednej wybranej osi układu warstw. */
export function rozmiescWarstwy(
  wybrane: readonly DesignBoardLayer[],
  pionowo: boolean,
): readonly DesignBoardLayer[] {
  if (wybrane.length < 3) return wybrane;
  const uporzadkowane = [...wybrane].sort((a, b) =>
    pionowo ? (a.y ?? 0) - (b.y ?? 0) : (a.x ?? 0) - (b.x ?? 0),
  );
  const pierwsza = uporzadkowane[0] as DesignBoardLayer;
  const ostatnia = uporzadkowane[uporzadkowane.length - 1] as DesignBoardLayer;
  const poczatek = pionowo ? (pierwsza.y ?? 0) : (pierwsza.x ?? 0);
  const krok = ((pionowo ? (ostatnia.y ?? 0) : (ostatnia.x ?? 0)) - poczatek)
    / (uporzadkowane.length - 1);

  return uporzadkowane.map((warstwa, numer) =>
    pionowo
      ? { ...warstwa, y: poczatek + krok * numer }
      : { ...warstwa, x: poczatek + krok * numer },
  );
}

/** Trzy role formatu społecznościowego — post, relacja i okładka — wyrażone proporcją, a nie rozmiarem. */
const PROPORCJE_SPOLECZNOSCIOWE: readonly (readonly [number, number])[] = [
  [1, 1], // post kwadratowy
  [9, 16], // relacja pionowa
  [16, 9], // okładka pozioma
];

/** Cztery szablony układu wymienione w panelu akcji modułu Design: tablica nastroju, siatka porównawcza, arkusz brandingowy i formaty społecznościowe. */
export function ulozWedlugSzablonu(
  wszystkie: readonly DesignBoardLayer[],
  szablon: SzablonUkladu,
): DesignBoardLayer[] {
  if (szablon === 'formaty-spolecznosciowe') return ulozFormaty(wszystkie);
  const kolumny = szablon === 'tablica-nastroju' ? 3 : 2;
  const bok = szablon === 'siatka-porownawcza' ? BOK_WYJSCIOWY * 1.5 : BOK_WYJSCIOWY;
  return wszystkie.map((warstwa, numer) => ({
    ...warstwa,
    x: (numer % kolumny) * (bok + ODSTEP),
    y: Math.floor(numer / kolumny) * (bok + ODSTEP),
    width: bok,
    height: szablon === 'arkusz-brandingowy' ? bok / 2 : bok,
    order: numer + 1,
  }));
}

/** Formaty społecznościowe: każda warstwa dostaje kolejną z trzech proporcji, w jednym rzędzie do wspólnej linii górnej. */
function ulozFormaty(wszystkie: readonly DesignBoardLayer[]): DesignBoardLayer[] {
  let przesuniecie = 0;
  return wszystkie.map((warstwa, numer) => {
    const proporcja = PROPORCJE_SPOLECZNOSCIOWE[
      numer % PROPORCJE_SPOLECZNOSCIOWE.length
    ] as readonly [number, number];
    const dluzszy = BOK_WYJSCIOWY * 1.5;
    const [szerokoscUdzial, wysokoscUdzial] = proporcja;
    const skala = dluzszy / Math.max(szerokoscUdzial, wysokoscUdzial);
    const szerokosc = szerokoscUdzial * skala;
    const ulozona: DesignBoardLayer = {
      ...warstwa,
      x: przesuniecie,
      y: 0,
      width: szerokosc,
      height: wysokoscUdzial * skala,
      order: numer + 1,
    };
    przesuniecie += szerokosc + ODSTEP;
    return ulozona;
  });
}

/**
 * Dokłada ramkę obszaru roboczego o zadanej szerokości, zachowując wysokość
 * równą bokowi wyjściowemu kompozycji.
 */
export function dolozRamke(zapis: ZapisKompozycji, nazwa: string, szerokosc: number): void {
  dolozWarstwe(zapis, { note: nazwa, width: szerokosc });
}
