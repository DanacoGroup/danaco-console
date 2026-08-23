/**
 * Kanwa grafu kroków — rysunek układu zależności w oknie Orchestratora.
 *
 * Opracowanie modułu opisuje graf zależności jako widok własny, oddzielny od
 * wykazu: węzły krokami, krawędzie zależnościami, podświetlona ścieżka
 * krytyczna, powiększanie i przewijanie. Wykaz mówi, co z czym jest związane;
 * rysunek mówi, jak długi jest łańcuch i gdzie tory się rozchodzą — a tego
 * z wykazu wierszy nie widać.
 *
 * Rysunek nie jest edytorem. Zapis układu idzie komendami kontraktu z panelu
 * akcji i z wykazu; kanwa pokazuje stan po zapisie. Układ węzłów liczy się tu
 * z samych zależności, więc jest powtarzalny, lecz nie jest układem, który
 * Operator mógłby ułożyć myszą: przenoszenie węzłów wymaga zapisania ich
 * położenia osobną komendą, a tej rdzeń jeszcze nie obsługuje. Do tego czasu
 * kanwa układa graf sama i nie obiecuje, że zapamięta cudze ułożenie.
 *
 * Barwy i grubości nie stoją tutaj: węzły i krawędzie noszą klasy rodziny
 * modułu, a wygląd niesie `automations.css` na żetonach `--dn-*`.
 */

/** Węzeł grafu — jeden krok automatyki. */
export interface WezelGrafu {
  /** Identyfikator kroku; po nim idą krawędzie i ścieżka krytyczna. */
  kod: string;
  /** Nazwa kroku; pusta zastępowana identyfikatorem. */
  nazwa: string;
}

/** Krawędź grafu — jedna zależność między krokami. */
export interface KrawedzGrafu {
  odKroku: string;
  doKroku: string;
  /** Rodzaj zależności albo jej warunek — podpis przy krawędzi. */
  podpis: string;
}

/** Układ podany kanwie do narysowania. */
export interface OpisGrafu {
  wezly: readonly WezelGrafu[];
  krawedzie: readonly KrawedzGrafu[];
  /** Kroki ścieżki krytycznej wyliczonej przez rdzeń. */
  sciezkaKrytyczna: ReadonlySet<string>;
}

/** Wymiary rysunku w jednostkach kanwy; wygląd bierze je przez `viewBox`. */
const SZEROKOSC_WEZLA = 176;
const WYSOKOSC_WEZLA = 40;
const ODSTEP_POZIOMY = 72;
const ODSTEP_PIONOWY = 24;
const MARGINES = 16;

/** Granice powiększenia — poniżej dolnej podpis znika, powyżej górnej graf ucieka z ekranu. */
const POWIEKSZENIE_NAJMNIEJSZE = 0.5;
const POWIEKSZENIE_NAJWIEKSZE = 2;
const KROK_POWIEKSZENIA = 0.25;

/** Przestrzeń nazw rysunku wektorowego; `createElement` nie zna tych znaczników. */
const PRZESTRZEN_SVG = 'http://www.w3.org/2000/svg';

export interface GrafKrokow {
  /** Element osadzany w ciele okna. */
  element: HTMLElement;
  /** Przerysowuje kanwę na nowym układzie. */
  pokaz(opis: OpisGrafu): void;
  /** Zmienia powiększenie o zadaną liczbę kroków; wartość ujemna pomniejsza. */
  powieksz(oIle: number): void;
  /** Rysunek w postaci tekstu — treść pliku eksportu. */
  zapisWektorowy(): string;
}

/**
 * Rozkłada węzły na warstwy według najdłuższej drogi od kroku bez poprzednika.
 *
 * Warstwa mówi, ile zależności trzeba przejść, zanim krok może ruszyć, więc
 * układ warstwowy pokazuje kolejność wykonania wprost. Węzeł stojący w cyklu
 * nie ma najdłuższej drogi — rachunek zatrzymuje się wtedy na węźle już
 * odwiedzonym i węzeł zostaje na warstwie, do której doszedł. Cykl jest
 * zastrzeżeniem układu, nie powodem, żeby nie narysować niczego.
 *
 * Funkcja jest czysta i nie dotyka dokumentu — sprawdzian czyta ją wprost.
 */
export function warstwyGrafu(opis: OpisGrafu): string[][] {
  const poprzednicy = new Map<string, string[]>();
  for (const wezel of opis.wezly) poprzednicy.set(wezel.kod, []);
  for (const krawedz of opis.krawedzie) {
    if (!poprzednicy.has(krawedz.doKroku) || !poprzednicy.has(krawedz.odKroku)) continue;
    poprzednicy.get(krawedz.doKroku)?.push(krawedz.odKroku);
  }

  const numery = new Map<string, number>();
  const wTrakcie = new Set<string>();

  function numer(kod: string): number {
    const znany = numery.get(kod);
    if (znany !== undefined) return znany;
    if (wTrakcie.has(kod)) return 0;
    wTrakcie.add(kod);
    let najdalszy = 0;
    for (const poprzednik of poprzednicy.get(kod) ?? []) {
      najdalszy = Math.max(najdalszy, numer(poprzednik) + 1);
    }
    wTrakcie.delete(kod);
    numery.set(kod, najdalszy);
    return najdalszy;
  }

  const warstwy: string[][] = [];
  for (const wezel of opis.wezly) {
    const miejsce = numer(wezel.kod);
    while (warstwy.length <= miejsce) warstwy.push([]);
    warstwy[miejsce]?.push(wezel.kod);
  }
  return warstwy;
}

export function utworzGrafKrokow(): GrafKrokow {
  let powiekszenie = 1;
  let szerokosc = 0;
  let wysokosc = 0;

  const rysunek = document.createElementNS(PRZESTRZEN_SVG, 'svg');
  rysunek.setAttribute('class', 'da-graf__rysunek');
  rysunek.setAttribute('role', 'img');

  const element = document.createElement('div');
  element.className = 'da-graf';
  element.append(rysunek);

  /** Ustawia rozmiar rysunku na ekranie — przewijanie należy do gospodarza. */
  function zastosujPowiekszenie(): void {
    rysunek.setAttribute('width', String(Math.round(szerokosc * powiekszenie)));
    rysunek.setAttribute('height', String(Math.round(wysokosc * powiekszenie)));
  }

  return {
    element,

    pokaz(opis) {
      const warstwy = warstwyGrafu(opis);
      const polozenia = polozeniaWezlow(warstwy);
      const najwyzszaWarstwa = warstwy.reduce((ile, warstwa) => Math.max(ile, warstwa.length), 0);
      szerokosc = MARGINES * 2 + warstwy.length * SZEROKOSC_WEZLA +
        Math.max(0, warstwy.length - 1) * ODSTEP_POZIOMY;
      wysokosc = MARGINES * 2 + najwyzszaWarstwa * WYSOKOSC_WEZLA +
        Math.max(0, najwyzszaWarstwa - 1) * ODSTEP_PIONOWY;
      rysunek.setAttribute('viewBox', `0 0 ${szerokosc} ${wysokosc}`);
      rysunek.setAttribute(
        'aria-label',
        `Graf zależności: ${opis.wezly.length} kroków, ${opis.krawedzie.length} zależności`,
      );
      rysunek.replaceChildren(
        grotStrzalki(),
        ...opis.krawedzie
          .map((krawedz) => narysujKrawedz(krawedz, polozenia, opis.sciezkaKrytyczna))
          .filter((krawedz): krawedz is SVGElement => krawedz !== null),
        ...opis.wezly
          .map((wezel) => narysujWezel(wezel, polozenia, opis.sciezkaKrytyczna))
          .filter((wezel): wezel is SVGElement => wezel !== null),
      );
      zastosujPowiekszenie();
    },

    powieksz(oIle) {
      powiekszenie = Math.min(
        POWIEKSZENIE_NAJWIEKSZE,
        Math.max(POWIEKSZENIE_NAJMNIEJSZE, powiekszenie + oIle * KROK_POWIEKSZENIA),
      );
      zastosujPowiekszenie();
    },

    zapisWektorowy: () => new XMLSerializer().serializeToString(rysunek),
  };
}

/** Środek węzła w jednostkach kanwy. */
interface Polozenie {
  x: number;
  y: number;
}

/** Rozkłada węzły warstw na współrzędne; warstwa idzie kolumną, węzeł wierszem. */
function polozeniaWezlow(warstwy: readonly (readonly string[])[]): Map<string, Polozenie> {
  const polozenia = new Map<string, Polozenie>();
  warstwy.forEach((warstwa, kolumna) => {
    warstwa.forEach((kod, wiersz) => {
      polozenia.set(kod, {
        x: MARGINES + kolumna * (SZEROKOSC_WEZLA + ODSTEP_POZIOMY),
        y: MARGINES + wiersz * (WYSOKOSC_WEZLA + ODSTEP_PIONOWY),
      });
    });
  });
  return polozenia;
}

/**
 * Grot strzałki dla krawędzi. Jeden na rysunek, przywoływany znacznikiem —
 * powielanie kształtu przy każdej krawędzi rozdęłoby plik eksportu.
 */
function grotStrzalki(): SVGElement {
  const grot = document.createElementNS(PRZESTRZEN_SVG, 'marker');
  grot.setAttribute('id', 'da-graf-grot');
  grot.setAttribute('viewBox', '0 0 8 8');
  grot.setAttribute('refX', '7');
  grot.setAttribute('refY', '4');
  grot.setAttribute('markerWidth', '6');
  grot.setAttribute('markerHeight', '6');
  grot.setAttribute('orient', 'auto-start-reverse');
  const ksztalt = document.createElementNS(PRZESTRZEN_SVG, 'path');
  ksztalt.setAttribute('d', 'M 0 0 L 8 4 L 0 8 z');
  ksztalt.setAttribute('class', 'da-graf__grot');
  grot.append(ksztalt);
  const zbior = document.createElementNS(PRZESTRZEN_SVG, 'defs');
  zbior.append(grot);
  return zbior;
}

/** Węzeł kroku: prostokąt z podpisem; krok ścieżki krytycznej dostaje znacznik. */
function narysujWezel(
  wezel: WezelGrafu,
  polozenia: ReadonlyMap<string, Polozenie>,
  sciezka: ReadonlySet<string>,
): SVGElement | null {
  const polozenie = polozenia.get(wezel.kod);
  if (polozenie === undefined) return null;
  const grupa = document.createElementNS(PRZESTRZEN_SVG, 'g');
  grupa.setAttribute('class', 'da-graf__wezel');
  grupa.dataset['krok'] = wezel.kod;
  if (sciezka.has(wezel.kod)) grupa.dataset['sciezka'] = 'krytyczna';

  const ksztalt = document.createElementNS(PRZESTRZEN_SVG, 'rect');
  ksztalt.setAttribute('x', String(polozenie.x));
  ksztalt.setAttribute('y', String(polozenie.y));
  ksztalt.setAttribute('width', String(SZEROKOSC_WEZLA));
  ksztalt.setAttribute('height', String(WYSOKOSC_WEZLA));
  ksztalt.setAttribute('rx', '6');

  const podpis = document.createElementNS(PRZESTRZEN_SVG, 'text');
  podpis.setAttribute('x', String(polozenie.x + SZEROKOSC_WEZLA / 2));
  podpis.setAttribute('y', String(polozenie.y + WYSOKOSC_WEZLA / 2 + 4));
  podpis.setAttribute('text-anchor', 'middle');
  podpis.textContent = skroc(wezel.nazwa === '' ? wezel.kod : wezel.nazwa);

  const pelnaNazwa = document.createElementNS(PRZESTRZEN_SVG, 'title');
  pelnaNazwa.textContent = wezel.nazwa === '' ? wezel.kod : `${wezel.nazwa} (${wezel.kod})`;

  grupa.append(ksztalt, podpis, pelnaNazwa);
  return grupa;
}

/** Krawędź zależności: odcinek od prawej krawędzi poprzednika do lewej następnika. */
function narysujKrawedz(
  krawedz: KrawedzGrafu,
  polozenia: ReadonlyMap<string, Polozenie>,
  sciezka: ReadonlySet<string>,
): SVGElement | null {
  const od = polozenia.get(krawedz.odKroku);
  const doo = polozenia.get(krawedz.doKroku);
  if (od === undefined || doo === undefined) return null;
  const grupa = document.createElementNS(PRZESTRZEN_SVG, 'g');
  grupa.setAttribute('class', 'da-graf__lacznik');
  if (sciezka.has(krawedz.odKroku) && sciezka.has(krawedz.doKroku)) {
    grupa.dataset['sciezka'] = 'krytyczna';
  }

  const odcinek = document.createElementNS(PRZESTRZEN_SVG, 'line');
  odcinek.setAttribute('x1', String(od.x + SZEROKOSC_WEZLA));
  odcinek.setAttribute('y1', String(od.y + WYSOKOSC_WEZLA / 2));
  odcinek.setAttribute('x2', String(doo.x));
  odcinek.setAttribute('y2', String(doo.y + WYSOKOSC_WEZLA / 2));
  odcinek.setAttribute('marker-end', 'url(#da-graf-grot)');

  const podpis = document.createElementNS(PRZESTRZEN_SVG, 'title');
  podpis.textContent = `${krawedz.odKroku} → ${krawedz.doKroku}: ${krawedz.podpis}`;

  grupa.append(odcinek, podpis);
  return grupa;
}

/** Podpis węzła przycięty do szerokości prostokąta; pełna nazwa idzie w dymku. */
function skroc(nazwa: string): string {
  const granica = 22;
  return nazwa.length <= granica ? nazwa : `${nazwa.slice(0, granica - 1)}…`;
}
