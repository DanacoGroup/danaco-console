import type { PozycjaWyboru } from './wybor-drzewem';

/**
 * Warstwa widoczności elementu sterującego modułu Terminal; nazwy są słowne, nie numerowane,
 * bo numer nic nie mówi czytelnikowi kodu.
 */
export type Warstwa = 'zawsze' | 'na-zadanie' | 'kontekstowa' | 'ekspercka';

/** Nazwa warstwy widoczności do zdań widocznych dla Operatora w opisie pozycji palety poleceń Terminala. */
export const NAZWY_WARSTW: Readonly<Record<Warstwa, string>> = {
  zawsze: 'zawsze widoczna',
  'na-zadanie': 'widoczna na żądanie',
  kontekstowa: 'rozwinięcie kontekstowe',
  ekspercka: 'funkcja ekspercka',
};

/**
 * Wykaz nastaw widoczności. Wyborem początkowym jest widoczność podstawowa: funkcja
 * niepotrzebna do bieżącego zadania nie stoi na ekranie.
 */
export const WIDOCZNOSCI: readonly PozycjaWyboru[] = [
  [
    'podstawowa',
    'Widoczność podstawowa',
    'Na ekranie zostają elementy zawsze widoczne i przywoływane na żądanie. Rozwinięcia kontekstowe i funkcje eksperckie schodzą z pola widzenia — nadal działają i sięga po nie paleta poleceń.',
  ],
  [
    'pelna',
    'Widoczność pełna',
    'Wszystkie cztery warstwy naraz — komplet kontrolek modułu na ekranie wraz z załącznikiem skrótów klawiszowych.',
  ],
];

/**
 * Znakuje elementy warstwami widoczności.
 *
 * Nadawanie idzie jednym zapisem dla całego panelu, a nie kontrolka po
 * kontrolce: warstwa wynika z porównania z sąsiadami, a wykaz zebrany w jednym
 * miejscu daje się przeczytać jak spis treści panelu.
 */
export function oznaczWarstwy(pary: ReadonlyArray<readonly [HTMLElement, Warstwa]>): void {
  for (const [element, warstwa] of pary) element.dataset['warstwa'] = warstwa;
}

/**
 * Jedna czynność okna widziana przez paletę poleceń i skróty klawiszowe.
 *
 * Nazwa okna stoi przy czynności, bo w palecie wszystkie sześć okien modułu
 * miesza się w jednym wykazie i „Eksportuj" bez okna źródłowego nie mówi, co
 * zostanie zapisane.
 */
export interface CzynnoscOkna {
  /** Okno, do którego czynność należy — nazwa z nagłówka ramy okna. */
  okno: string;
  /** Nazwa czynności w bezokoliczniku albo trybie rozkazującym, jak na przycisku. */
  nazwa: string;
  /** Zdanie o skutku czynności — objaśnienie pozycji palety. */
  opis: string;
  warstwa: Warstwa;
  /** Skrót klawiszowy z załącznika dokumentacji modułu; pominięty znaczy „bez skrótu”. */
  skrot?: string;
  wykonaj(): void;
}
