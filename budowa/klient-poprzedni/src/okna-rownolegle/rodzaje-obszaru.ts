/**
 * Rodzaje obszaru roboczego na scenie okien równoległych — kolumna rozmowy, kolumna paneli pomocniczych i widok pełnoekranowy — wraz z ich szerokościami minimalnymi wiązanymi w jednym źródle liczb.
 */
export type RodzajObszaru = 'kolumna-rozmowy' | 'kolumna-paneli' | 'widok-pelnoekranowy';

/**
 * Najwęższa dopuszczalna kolumna rozmowy.
 *
 * Poniżej tej wartości wątek przestaje być czytelny — wiersz kodu w odpowiedzi
 * modelu łamie się w każdej linii, a karta pytania nie mieści przycisków obok
 * siebie.
 */
export const SZEROKOSC_MIN_ROZMOWY = 600;

/**
 * Najwęższa dopuszczalna kolumna paneli pomocniczych.
 *
 * 320 px to próg czytelności, na którym stoi tor okien (`uklad.css`) i który
 * panele dziedziczą.
 */
export const SZEROKOSC_MIN_PANELU = 320;

/**
 * Szerokość, z jaką stos paneli otwiera się po raz pierwszy.
 *
 * Między minimum a wygodą: mieści drzewo plików i terminal bez zwijania
 * ścieżek, a rozmowie zostawia jej 600 px na scenie o typowej szerokości.
 */
export const SZEROKOSC_PANELU_DOMYSLNA = 420;

/**
 * Minimum obszaru danego rodzaju, w pikselach.
 *
 * `widok-pelnoekranowy` ma minimum zero: bierze całą scenę, więc nie konkuruje
 * o piksele z żadnym sąsiadem.
 */
export function minimumObszaru(rodzaj: RodzajObszaru): number {
  switch (rodzaj) {
    case 'kolumna-rozmowy':
      return SZEROKOSC_MIN_ROZMOWY;
    case 'kolumna-paneli':
      return SZEROKOSC_MIN_PANELU;
    case 'widok-pelnoekranowy':
      return 0;
  }
}
