import type { PozycjaWyboru } from './wybor-drzewem';

/**
 * Warstwy widoczności modułu Terminal i wykaz czynności, które okna oddają
 * palecie poleceń oraz skrótom klawiszowym.
 *
 * Moduł zastępuje emulator terminala, multiplekser sesji, klienta SSH,
 * harmonogram zadań i bibliotekę skryptów naraz. Gdyby każda z tych zdolności
 * miała stale własny przycisk, przestrzeń pracy byłaby ścianą kontrolek, przez
 * którą nie widać powłoki. Dlatego każdy element sterujący należy do jednej
 * z czterech warstw, a moduł pokazuje naraz tylko tyle, ile potrzeba do
 * bieżącego zadania.
 *
 * Ukrycie nie jest blokadą i nie wolno go z blokadą mylić. Element warstwy
 * zwiniętej nie dostaje `disabled`, nie znika z drzewa dostępności przez
 * `aria-hidden` na przodku ogniskowalnym i nie przestaje działać — zostaje
 * zdjęty z pola widzenia, a droga do niego prowadzi paletą poleceń
 * (`paleta-polecen.ts`) albo skrótem (`skroty-klawiszowe.ts`). Każda czynność
 * modułu jest osiągalna jednym wskazaniem niezależnie od nastawy widoczności.
 *
 * Czynność jest opisem, nie przyciskiem: paleta i skróty nie sięgają do węzłów
 * okien, bo węzeł przerysowuje się przy każdej zmianie stanu i uchwyt trzymany
 * na zewnątrz wskazywałby po chwili element, którego nie ma w dokumencie.
 * Okno oddaje więc wywołanie, a nie kontrolkę.
 */

/**
 * Warstwa widoczności elementu sterującego.
 *
 * Nazwy są słowne, nie numerowane: numer warstwy nic nie mówi czytelnikowi
 * kodu, a nazwa mówi, kiedy element ma być widoczny.
 */
export type Warstwa = 'zawsze' | 'na-zadanie' | 'kontekstowa' | 'ekspercka';

/** Nazwa warstwy do zdań widocznych dla Operatora. */
export const NAZWY_WARSTW: Readonly<Record<Warstwa, string>> = {
  zawsze: 'zawsze widoczna',
  'na-zadanie': 'widoczna na żądanie',
  kontekstowa: 'rozwinięcie kontekstowe',
  ekspercka: 'funkcja ekspercka',
};

/**
 * Wykaz nastaw widoczności.
 *
 * Wyborem początkowym jest widoczność podstawowa, bo tak stanowi zasada modułu:
 * funkcja niepotrzebna do bieżącego zadania nie stoi na ekranie. Widoczność
 * pełna nie jest trybem administracyjnym ani ukrytym — to zwykła nastawa
 * o jedno wskazanie dalej.
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
