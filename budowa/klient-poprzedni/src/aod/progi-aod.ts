/**
 * Progi, częstotliwość i czasy wyciszenia funkcji globalnej Always On Display.
 * Plik jest jedynym miejscem, w którym te liczby stoją; reguła rozpoznania,
 * magazyn kolejki i stan obecności biorą je stąd, zamiast powtarzać u siebie.
 */

/**
 * Milisekunda liczona minutami, żeby progi czasowe niżej dało się czytać wprost
 * w minutach, a nie przeliczać z liczby milisekund przy każdym odczycie.
 */
const MINUTA_MS = 60_000;

/**
 * Milisekunda liczona godzinami, złożona z sześćdziesięciu minut. Stała służy
 * zapisaniu czasu życia sugestii w jednostce, w której próg jest podawany.
 */
const GODZINA_MS = 60 * MINUTA_MS;

/**
 * Progi wyzwalania sugestii wraz z ich wartościami domyślnymi. Nazwy pól są
 * pełnymi nazwami progów, a nie skrótami: próg czytany w kodzie ma mówić to
 * samo, co próg czytany w tabeli wartości.
 */
export const PROGI_AOD = {
  /** Liczba sugestii ujawnianych samoczynnie w godzinie: 3; kolejne czekają. */
  liczbaDymkowNaGodzine: 3,

  /** Odstęp między dymkami: 5 minut; bliższe łączą się w pozycję zbiorczą. */
  odstepMiedzyDymkamiMs: 5 * MINUTA_MS,

  /** Próg czasu oczekiwania zadania w kolejce: 15 minut. */
  oczekiwanieWKolejceMs: 15 * MINUTA_MS,

  /** Próg liczby ponowień zadania: 2; trzecie tworzy sugestię wagi wysokiej. */
  liczbaPonowien: 2,

  /** Próg wypełnienia kolejki: 80 procent pojemności. */
  wypelnienieKolejki: 0.8,

  /** Próg powtarzalności czynności ręcznej: 3 wystąpienia w karcie sesji. */
  powtarzalnoscCzynnosci: 3,

  /** Czas życia sugestii nieprzyjętej: 24 godziny, po których zostaje odrzucona. */
  czasZyciaSugestiiMs: 24 * GODZINA_MS,
} as const;

/**
 * Rodzaje wyciszenia czasowego dostępne w menu kebab: kwadrans, godzina oraz
 * wyciszenie do końca dnia. Zakres jest zamknięty, więc rodzaj spoza wykazu nie
 * powstaje w oknie ani nie wchodzi do stanu obecności.
 */
export const WyciszenieCzasowe = {
  /** 15 minut. */
  Kwadrans: 'kwadrans',
  /** 1 godzina. */
  Godzina: 'godzina',
  /** Do końca dnia. */
  DoKoncaDnia: 'do-konca-dnia',
} as const;
export type WyciszenieCzasowe = (typeof WyciszenieCzasowe)[keyof typeof WyciszenieCzasowe];

/**
 * Nazwa wyciszenia widziana przez Operatora w menu kebab. Nazwa mówi o czasie
 * trwania, a nie o rodzaju technicznym, ponieważ to czas jest tym, co Operator
 * wybiera, przystawiając wyciszenie.
 */
export const NAZWY_WYCISZEN: Readonly<Record<WyciszenieCzasowe, string>> = {
  [WyciszenieCzasowe.Kwadrans]: 'Wycisz na 15 minut',
  [WyciszenieCzasowe.Godzina]: 'Wycisz na godzinę',
  [WyciszenieCzasowe.DoKoncaDnia]: 'Wycisz do końca dnia',
};

/**
 * Chwila, do której trwa wyciszenie czasowe.
 *
 * „Do końca dnia" liczymy zegarem maszyny Operatora — funkcja stoi na jego
 * stanowisku, a nie w strefie czasowej rdzenia.
 *
 * @param teraz chwila wydania polecenia w milisekundach epoki.
 */
export function konieCzasuWyciszenia(rodzaj: WyciszenieCzasowe, teraz: number): number {
  switch (rodzaj) {
    case WyciszenieCzasowe.Kwadrans:
      return teraz + 15 * MINUTA_MS;
    case WyciszenieCzasowe.Godzina:
      return teraz + GODZINA_MS;
    case WyciszenieCzasowe.DoKoncaDnia: {
      const koniec = new Date(teraz);
      koniec.setHours(23, 59, 59, 999);
      return koniec.getTime();
    }
  }
}
