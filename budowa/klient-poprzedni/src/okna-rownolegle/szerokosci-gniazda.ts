/**
 * Rachunek podziału szerokości jednego gniazda ustala, ile dostaje rozmowa i ile stos paneli pomocniczych, bez odwołania do DOM, sprowadzając wartości bezsensowne do najbliższych sensownych i pozwalając gniazdu przewijać się w ciasnocie.
 */

import {
  SZEROKOSC_MIN_PANELU,
  SZEROKOSC_MIN_ROZMOWY,
} from './rodzaje-obszaru';

/** Podział szerokości jednego gniazda na kolumnę rozmowy i kolumnę paneli, wyrażony dwiema liczbami w pikselach. */
export interface PodzialGniazda {
  rozmowa: number;
  panele: number;
}

/**
 * Liczba sprowadzona do wartości nieujemnej i skończonej; wartość spoza zakresu nie wywraca układu, tylko zostaje przycięta do najbliższej sensownej.
 */
function liczbaNieujemna(wartosc: number, gdyBezsensu: number): number {
  if (!Number.isFinite(wartosc)) return gdyBezsensu;
  return Math.max(0, Math.round(wartosc));
}

/**
 * Szerokość, o którą prosi wołający, sprowadzona do wartości sensownej — NaN i wartość ujemna znaczą brak żądania, a nieskończoność znaczy żądanie maksimum.
 */
function szerokoscZadana(zadana: number): number {
  if (Number.isNaN(zadana)) return 0;
  if (zadana === Number.POSITIVE_INFINITY) return Number.POSITIVE_INFINITY;
  if (zadana === Number.NEGATIVE_INFINITY) return 0;
  return Math.max(0, Math.round(zadana));
}

/**
 * Szerokość kolumny paneli przycięta tak, by rozmowie zostało minimum; zero na wejściu zostaje zerem, a każda wartość dodatnia dostaje co najmniej minimum panelu.
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
 * Podział przy danej szerokości gniazda, gdzie zero paneli znaczy pusty stos; rozmowa dostaje resztę, nigdy mniej niż swoje minimum, nawet kosztem przewijania.
 */
export function ustalPodzial(zadanaSzerokoscPaneli: number, dostepna: number): PodzialGniazda {
  const scena = liczbaNieujemna(dostepna, SZEROKOSC_MIN_ROZMOWY);
  const panele = ograniczSzerokoscPaneli(zadanaSzerokoscPaneli, scena);
  const rozmowa = Math.max(SZEROKOSC_MIN_ROZMOWY, scena - panele);

  return { rozmowa, panele };
}

/**
 * Wartość układu kolumn siatki gniazda, gdzie brak wartości znaczy brak kolumny paneli, a rozmowa rośnie z gniazdem, nie schodząc poniżej minimum.
 */
export function kolumnyGniazda(podzial: PodzialGniazda | null): string {
  if (podzial === null || podzial.panele <= 0) {
    return `minmax(${SZEROKOSC_MIN_ROZMOWY}px, 1fr)`;
  }

  return `minmax(${SZEROKOSC_MIN_ROZMOWY}px, 1fr) auto ${podzial.panele}px`;
}
