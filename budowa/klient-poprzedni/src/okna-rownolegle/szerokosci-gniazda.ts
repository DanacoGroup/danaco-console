/**
 * Rachunek podziału szerokości jednego gniazda: ile dostaje rozmowa, ile stos
 * paneli pomocniczych.
 *
 * Rachunek stoi osobno, bez ani jednego odwołania do DOM, bo podział szerokości
 * myli się po cichu: gdy panel weźmie za dużo, rozmowa nie znika — robi się
 * nieczytelna.
 *
 * Reguły podziału:
 *  - rozmowa nie schodzi poniżej `SZEROKOSC_MIN_ROZMOWY`; panel proszący o całą
 *    scenę dostaje tyle, ile zostaje po rozmowie;
 *  - kolumna paneli ma albo dokładnie 0 (stosu nie ma), albo co najmniej
 *    `SZEROKOSC_MIN_PANELU` — panel węższy od minimum jest paskiem, na którym
 *    nic nie widać;
 *  - gdy dostępna szerokość nie mieści sumy minimów, rozmowa dostaje swoje
 *    minimum, panele swoje, a gniazdo przewija się w poziomie, tak jak tor okien
 *    (`uklad.css`, `overflow-x: auto`). Ciasnota nie odmawia otwarcia panelu —
 *    szerokość ustawia uchwyt (`uchwyt-szerokosci.ts`);
 *  - wartość bezsensowna (NaN, nieskończoność, liczba ujemna) nie wstrzymuje
 *    rachunku, tylko zostaje sprowadzona do najbliższej sensownej, wzorem
 *    `ograniczLiczbe` z `identyfikatory.ts`.
 *
 * Moduł nie tworzy elementów, nie czyta `getBoundingClientRect` i nie wie, jakie
 * panele są otwarte ani ile ich jest. Dostaje dwie liczby, oddaje podział albo
 * gotowy napis dla `grid-template-columns`.
 */

import {
  SZEROKOSC_MIN_PANELU,
  SZEROKOSC_MIN_ROZMOWY,
} from './rodzaje-obszaru';

/** Podział szerokości jednego gniazda na kolumnę rozmowy i kolumnę paneli. */
export interface PodzialGniazda {
  rozmowa: number;
  panele: number;
}

/**
 * Liczba sprowadzona do wartości nieujemnej i skończonej.
 *
 * Wzorzec `ograniczLiczbe` (`identyfikatory.ts`): wartość spoza zakresu nie
 * wywraca układu, tylko zostaje przycięta. Nieskończoność w górę to „tyle, ile
 * się da" — sufit nakłada dopiero wołający, bo tylko on zna scenę.
 */
function liczbaNieujemna(wartosc: number, gdyBezsensu: number): number {
  if (!Number.isFinite(wartosc)) return gdyBezsensu;
  return Math.max(0, Math.round(wartosc));
}

/**
 * Szerokość, o którą prosi wołający, sprowadzona do wartości sensownej.
 *
 * Trzy przypadki bezsensowne mają trzy różne znaczenia i nie wolno ich zlepić
 * w jedno: NaN i wartość ujemna to „nie ma o co prosić" (stos pusty),
 * a nieskończoność w górę to „tyle, ile się da" — sufit nakłada minimum rozmowy
 * kilka wierszy niżej, nie ta funkcja.
 */
function szerokoscZadana(zadana: number): number {
  if (Number.isNaN(zadana)) return 0;
  if (zadana === Number.POSITIVE_INFINITY) return Number.POSITIVE_INFINITY;
  if (zadana === Number.NEGATIVE_INFINITY) return 0;
  return Math.max(0, Math.round(zadana));
}

/**
 * Szerokość kolumny paneli przycięta tak, by rozmowie zostało minimum.
 *
 * Zero na wejściu znaczy „stosu paneli nie ma" i zostaje zerem — funkcja nie
 * otwiera panelu sama z siebie. Każda wartość dodatnia dostaje co najmniej
 * `SZEROKOSC_MIN_PANELU`, także wtedy, gdy scena jest za ciasna; gniazdo
 * przewija się wtedy w poziomie.
 */
export function ograniczSzerokoscPaneli(zadana: number, dostepna: number): number {
  const chciana = szerokoscZadana(zadana);
  if (chciana === 0) return 0;

  const scena = liczbaNieujemna(dostepna, 0);
  // Ile zostaje panelom po odjęciu minimum rozmowy. Ujemne znaczy ciasnotę.
  const zapas = scena - SZEROKOSC_MIN_ROZMOWY;
  if (zapas < SZEROKOSC_MIN_PANELU) return SZEROKOSC_MIN_PANELU;

  return Math.min(Math.max(chciana, SZEROKOSC_MIN_PANELU), zapas);
}

/**
 * Podział przy danej szerokości gniazda; `panele = 0` znaczy „stos pusty".
 *
 * Rozmowa dostaje resztę, ale nigdy mniej niż `SZEROKOSC_MIN_ROZMOWY` — przy
 * ciasnocie suma podziału bywa większa od `dostepna` i to jest zamierzone:
 * nadmiar zjada poziome przewijanie, a nie minimum rozmowy.
 */
export function ustalPodzial(zadanaSzerokoscPaneli: number, dostepna: number): PodzialGniazda {
  const scena = liczbaNieujemna(dostepna, SZEROKOSC_MIN_ROZMOWY);
  const panele = ograniczSzerokoscPaneli(zadanaSzerokoscPaneli, scena);
  const rozmowa = Math.max(SZEROKOSC_MIN_ROZMOWY, scena - panele);

  return { rozmowa, panele };
}

/**
 * Wartość `grid-template-columns` gniazda; `null` = brak kolumny paneli.
 *
 * Rozmowa idzie jako `minmax(<minimum>px, 1fr)` — rośnie z gniazdem, ale nie
 * schodzi poniżej minimum. Środkowe `auto` to miejsce na uchwyt szerokości;
 * gdy paneli nie ma, nie ma też uchwytu, bo nie byłoby czego rozsuwać.
 */
export function kolumnyGniazda(podzial: PodzialGniazda | null): string {
  if (podzial === null || podzial.panele <= 0) {
    return `minmax(${SZEROKOSC_MIN_ROZMOWY}px, 1fr)`;
  }

  return `minmax(${SZEROKOSC_MIN_ROZMOWY}px, 1fr) auto ${podzial.panele}px`;
}
