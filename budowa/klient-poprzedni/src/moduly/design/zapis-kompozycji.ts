import type { DesignAsset, DesignBoardLayer } from '../../../../shared/contract';
import { nowyIdentyfikator } from '../../protokol/identyfikator';

/**
 * Zapis kompozycji Design Board wraz z czynnościami, które go zmieniają.
 *
 * Odpowiada wyłącznie za układ warstw i zaznaczenie — bez powiadamiania widoku
 * i bez wywołań rdzenia; rozgłaszaniem zmian zajmuje się `stan-kompozycji.ts`.
 * Zmiana układu jest rachunkiem na liczbach, rozgłoszenie obsługą obserwatorów;
 * rozdzielone czytają się i sprawdzają osobno.
 */

/** Widok kanwy: powiększenie, przesunięcie swobodne, siatka pomocnicza. */
export interface WidokKanwy {
  /** Powiększenie; 1 znaczy skalę naturalną. */
  powiekszenie: number;
  przesuniecieX: number;
  przesuniecieY: number;
  siatka: boolean;
}

/** Oś wyrównania warstw zaznaczonych. */
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

/** Pełny zapis kompozycji — jedyny nośnik prawdy o układzie. */
export interface ZapisKompozycji {
  warstwy: DesignBoardLayer[];
  wybor: string[];
  nazwa: string;
  identyfikator: string;
  widok: WidokKanwy;
  licznik: number;
}

/** Bok wyjściowy warstwy na kanwie, w jednostkach kompozycji. */
export const BOK_WYJSCIOWY = 240;

/** Odstęp między warstwami układanymi szablonem. */
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
 * Wstawia kompozycję odczytaną z rdzenia na miejsce układu bieżącego.
 *
 * Zastępuje, nie scala: kompozycja z rdzenia jest pełnym stanem planszy, a nie
 * jej dokładką — scalenie dałoby układ, którego nie ma ani na ekranie, ani
 * w bazie, a zapisany z powrotem założyłby byt trzeci. Okno mówi Operatorowi
 * wprost, że odczyt nadpisuje to, co ma na kanwie.
 *
 * Identyfikator wchodzi razem z układem, bo bez niego kolejny zapis założyłby
 * kompozycję nową zamiast zaktualizować odczytaną i plansza rozmnażałaby się
 * w bazie po jednej sztuce na każde wejście w moduł.
 *
 * Licznik warstw przesuwamy ponad wczytane, żeby kod warstwy dokładanej po
 * odczycie nie zderzył się z kodem warstwy, która przyszła z rdzenia
 * (uzasadnienie unikalności przy `nowyKodWarstwy` niżej).
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

/** Dokłada warstwę zbudowaną z zasobu Assets Panel. */
export function dolozZasob(zapis: ZapisKompozycji, zasob: DesignAsset): void {
  dolozWarstwe(zapis, {
    assetId: zasob.id,
    note: zasob.name ?? zasob.id,
    ...(zasob.width === undefined ? {} : { width: zasob.width }),
    ...(zasob.height === undefined ? {} : { height: zasob.height }),
  });
}

/** Dokłada warstwę pomocniczą bez zasobu — element biblioteki pomocniczej. */
export function dolozElement(zapis: ZapisKompozycji, nazwa: string): void {
  dolozWarstwe(zapis, { note: nazwa });
}

/**
 * Identyfikator warstwy — kod, który okno nadaje bytowi rdzenia.
 *
 * Sam licznik okna nie wystarcza. Kolumna
 * `warstwa_kompozycji_design.identyfikator_zewnetrzny` jest UNIQUE w całej
 * tabeli, nie w obrębie kompozycji (`migracja_048_design.sql`),
 * a `design.board.update` wstawia warstwy zwykłym INSERT-em bez ON CONFLICT.
 * Kod z licznika wracałby po każdym przeładowaniu okna do „warstwa-1"
 * i zderzał się z warstwą kompozycji zapisanej wcześniej, co rdzeń odrzuca
 * naruszeniem ograniczenia UNIQUE. Człon losowy bierzemy
 * z `protokol/identyfikator.ts`, a numer porządkowy zostaje z przodu, żeby kod
 * warstwy dało się przeczytać w panelu warstw.
 */
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

/** Nanosi zmianę na jedną warstwę. */
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

/** Warstwy objęte działaniem narzędzi: zaznaczone i niezablokowane. */
export function objeteWarstwy(zapis: ZapisKompozycji): DesignBoardLayer[] {
  return zapis.warstwy.filter((w) => zapis.wybor.includes(w.id) && w.locked !== true);
}

/** Nanosi wynik działania narzędzia na zapis, zachowując kolejność warstw. */
export function naniesWarstwy(
  zapis: ZapisKompozycji,
  zmienione: readonly DesignBoardLayer[],
): void {
  const wedlugKodu = new Map(zmienione.map((w) => [w.id, w]));
  zapis.warstwy = zapis.warstwy.map((w) => wedlugKodu.get(w.id) ?? w);
}

/** Wyrównanie warstw do skrajnej albo środkowej pozycji zbioru. */
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

/** Równe odstępy między warstwami wzdłuż jednej osi. */
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

/**
 * Trzy role formatu społecznościowego wymienione w opracowaniu — post, relacja
 * i okładka — wyrażone PROPORCJĄ, nie rozmiarem konkretnej platformy.
 *
 * Proporcja jest tym, co opracowanie rozstrzyga („rozmiary postów, relacji
 * i okładek"); wpisanie tu pikseli którejś platformy byłoby wniesieniem do
 * produktu liczby, której nie ma w żadnym dokumencie projektu i która starzeje
 * się razem z cudzym interfejsem. Bok dłuższy bierze się z boku wyjściowego
 * kompozycji, więc szablon skaluje się razem z kanwą.
 */
const PROPORCJE_SPOLECZNOSCIOWE: readonly (readonly [number, number])[] = [
  [1, 1], // post kwadratowy
  [9, 16], // relacja pionowa
  [16, 9], // okładka pozioma
];

/** Cztery szablony układu wymienione w panelu akcji modułu Design. */
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

/**
 * Formaty społecznościowe: każda warstwa dostaje kolejną z trzech proporcji,
 * a wszystkie stoją w jednym rzędzie wyrównane do wspólnej linii górnej.
 *
 * Rząd, a nie siatka: formaty ogląda się obok siebie, żeby porównać kadr tego
 * samego materiału w trzech proporcjach.
 */
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
 * Dokłada ramkę obszaru roboczego o zadanej szerokości.
 *
 * Szerokość pochodzi z punktu łamania produktu, a WYSOKOŚĆ zostaje przy boku
 * wyjściowym, bo punkty łamania opisują wyłącznie szerokość. Podstawienie tu
 * wysokości byłoby wniesieniem wymiaru, którego kierunek projektowy nie
 * rozstrzyga; Operator ustawia ją w inspektorze właściwości.
 */
export function dolozRamke(zapis: ZapisKompozycji, nazwa: string, szerokosc: number): void {
  dolozWarstwe(zapis, { note: nazwa, width: szerokosc });
}
