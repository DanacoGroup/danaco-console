/**
 * Ślad adnotacji — model rysunku operatora i jego wykreślenie na płótnie.
 * Rysunek żyje w modelu, nie w pikselach, więc wykaz śladów jest jedynym
 * trwałym zapisem adnotacji; ślad zapisuje się w ułamku ramki, nie bezwzględnie.
 */

/** Narzędzia paska adnotacji dostępne operatorowi sesji: ołówek, linia, prostokąt, elipsa i pole tekstu. */
export type NarzedzieAdnotacji = 'olowek' | 'linia' | 'prostokat' | 'elipsa' | 'tekst';

/** Punkt zapisany w ułamku ramki płótna: zero to lewa albo górna krawędź, jedynka to prawa albo dolna krawędź. */
export interface Punkt {
  x: number;
  y: number;
}

/** Jedno pociągnięcie adnotacji: narzędzie, barwa, treść napisu i wykaz punktów, którymi je poprowadzono. */
export interface Slad {
  narzedzie: NarzedzieAdnotacji;
  /** Nazwa żetonu motywu, nie rozwiązana barwa — rozwiązuje ją płótno. */
  zeton: string;
  punkty: Punkt[];
  /** Treść napisu; pusta dla każdego narzędzia poza „tekst". */
  tekst: string;
}

/**
 * Wycinek płótna 2D, którego rysowanie naprawdę używa. Produkt podaje kontekst
 * przeglądarki wprost, a sprawdzian — atrapę pisaka spisującą wywołania.
 */
export interface PisakPlotna {
  strokeStyle: string | CanvasGradient | CanvasPattern;
  fillStyle: string | CanvasGradient | CanvasPattern;
  lineWidth: number;
  lineCap: CanvasLineCap;
  lineJoin: CanvasLineJoin;
  font: string;
  beginPath(): void;
  moveTo(x: number, y: number): void;
  lineTo(x: number, y: number): void;
  stroke(): void;
  strokeRect(x: number, y: number, szerokosc: number, wysokosc: number): void;
  ellipse(
    x: number,
    y: number,
    promienX: number,
    promienY: number,
    obrot: number,
    poczatek: number,
    koniec: number,
  ): void;
  fillText(tekst: string, x: number, y: number): void;
  clearRect(x: number, y: number, szerokosc: number, wysokosc: number): void;
}

/**
 * Oprawa jednego wykreślenia: rozmiar ramki, barwa i krój. Barwa i krój
 * przychodzą już rozwiązane, bo płótno przyjmuje wyłącznie wartości gotowe.
 */
export interface OprawaSladu {
  szerokosc: number;
  wysokosc: number;
  barwa: string;
  kroj: string;
}

/** Grubość kreski i wysokość napisu w pikselach płótna — stała geometria rysunku, niezależna od motywu. */
const GRUBOSC_KRESKI = 3;
const WYSOKOSC_NAPISU = 16;

/**
 * Wykłada jeden ślad na pisaku.
 *
 * Ślad bez punktów nie rysuje nic i nie jest usterką: pociągnięcie zaczęte
 * i porzucone poza płótnem ma prawo nie mieć ani jednego punktu.
 */
export function narysujSlad(pisak: PisakPlotna, slad: Slad, oprawa: OprawaSladu): void {
  const punkty = slad.punkty.map((punkt) => ({
    x: punkt.x * oprawa.szerokosc,
    y: punkt.y * oprawa.wysokosc,
  }));
  const poczatek = punkty[0];
  if (poczatek === undefined) return;

  pisak.strokeStyle = oprawa.barwa;
  pisak.fillStyle = oprawa.barwa;
  pisak.lineWidth = GRUBOSC_KRESKI;
  pisak.lineCap = 'round';
  pisak.lineJoin = 'round';

  const koniec = punkty[punkty.length - 1] ?? poczatek;

  switch (slad.narzedzie) {
    case 'olowek':
      pisak.beginPath();
      pisak.moveTo(poczatek.x, poczatek.y);
      for (const punkt of punkty.slice(1)) pisak.lineTo(punkt.x, punkt.y);
      pisak.stroke();
      return;

    case 'linia':
      pisak.beginPath();
      pisak.moveTo(poczatek.x, poczatek.y);
      pisak.lineTo(koniec.x, koniec.y);
      pisak.stroke();
      return;

    case 'prostokat':
      pisak.strokeRect(poczatek.x, poczatek.y, koniec.x - poczatek.x, koniec.y - poczatek.y);
      return;

    case 'elipsa':
      pisak.beginPath();
      pisak.ellipse(
        (poczatek.x + koniec.x) / 2,
        (poczatek.y + koniec.y) / 2,
        Math.abs(koniec.x - poczatek.x) / 2,
        Math.abs(koniec.y - poczatek.y) / 2,
        0,
        0,
        Math.PI * 2,
      );
      pisak.stroke();
      return;

    case 'tekst':
      // Krój ustawia się przed napisem, bo `font` przestawia też metrykę
      // linii bazowej.
      pisak.font = oprawa.kroj;
      pisak.fillText(slad.tekst, poczatek.x, poczatek.y);
      return;
  }
}

/**
 * Wykłada cały wykaz śladów na czystym pisaku.
 *
 * Czyszczenie jest częścią wykreślenia, nie osobną czynnością: obraz płótna
 * ma być odbiciem wykazu, a przerysowanie bez zmazania zostawiłoby na nim
 * ślady już zdjęte z wykazu.
 */
export function narysujSlady(
  pisak: PisakPlotna,
  slady: readonly Slad[],
  ramka: { szerokosc: number; wysokosc: number },
  barwaZetonu: (zeton: string) => string,
  kroj: string,
): void {
  pisak.clearRect(0, 0, ramka.szerokosc, ramka.wysokosc);
  for (const slad of slady) {
    narysujSlad(pisak, slad, {
      szerokosc: ramka.szerokosc,
      wysokosc: ramka.wysokosc,
      barwa: barwaZetonu(slad.zeton),
      kroj,
    });
  }
}

/** Krój napisu płótna zbudowany z kroju czcionki wyliczonego dla elementu dokumentu w bieżącym motywie. */
export function krojNapisu(rodzina: string): string {
  return `${WYSOKOSC_NAPISU}px ${rodzina}`;
}
