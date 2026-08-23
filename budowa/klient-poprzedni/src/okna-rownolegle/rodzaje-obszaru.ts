/**
 * Rodzaje obszaru roboczego na scenie okien równoległych i ich szerokości
 * minimalne.
 *
 * Scena dzieli szerokość między trzy różne byty: kolumnę rozmowy, kolumnę
 * paneli pomocniczych i widok pełnoekranowy. Każdy z nich ma inne minimum,
 * a minimum wiąże: kolumna rozmowy nie schodzi poniżej 600 px, kolumny panelowe
 * trzymają swoje.
 *
 * Te same liczby czyta rachunek podziału gniazda (`szerokosci-gniazda.ts`),
 * uchwyt szerokości i arkusz toru (`uklad.css`). Trzy kopie tej samej liczby to
 * trzy różne progi czytelności po pierwszej poprawce; jedno źródło znaczy, że
 * podniesienie minimum jest jedną zmianą.
 *
 * Moduł nie zna DOM, nie zna paneli po nazwie i nie rozstrzyga, ile paneli
 * wolno otworzyć. Oddaje liczby — decyzję podejmuje ten, kto pyta. Gdy miejsca
 * zabraknie, szerokość zmienia się uchwytem; panel nie chowa się sam.
 */

/**
 * Rodzaj obszaru roboczego wewnątrz gniazda albo na całej scenie.
 *
 * `kolumna-rozmowy` — wątek rozmowy; `kolumna-paneli` — stos paneli
 * pomocniczych obok rozmowy; `widok-pelnoekranowy` — obszar biorący całą scenę,
 * bez sąsiada, z którym miałby się dzielić.
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
