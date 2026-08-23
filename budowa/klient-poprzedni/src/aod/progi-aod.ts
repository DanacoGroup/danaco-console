/**
 * Progi, częstotliwość i czasy wyciszenia funkcji globalnej Always On Display.
 *
 * Każda wartość tego pliku pochodzi z opracowania
 * `docs/funkcje-globalne/always-on-display.md` — rozdz. 3.4 (progi
 * i częstotliwość) oraz rozdz. 3.5 (wyciszanie). Żadnej wartości tu nie
 * wymyślono; przy każdej stoi rozdział, z którego jest wzięta.
 *
 * Plik jest jednym miejscem, w którym te liczby stoją. Reguła rozpoznania
 * (`rozpoznanie-decyzji.ts`), magazyn kolejki (`kolejka-decyzji.ts`) i stan
 * obecności (`tryb-obecnosci.ts`) biorą je stąd, zamiast powtarzać u siebie.
 *
 * Opracowanie mówi (rozdz. 3.4, zdanie zamykające), że wszystkie te progi są
 * ustawieniami zasięgu „Always On Display" i podlegają zmianie z okna
 * konfiguracji, z okna Ustawień i poleceniem języka naturalnego. Kontrakt nie
 * niesie ani kategorii ustawień „Always On Display", ani komendy zapisującej
 * te wartości — dlatego wartości domyślne stoją tutaj jako stałe, a nie jako
 * odczyt `config.get`. Rozjazd jest zgłoszony, nie zasypany atrapą.
 */

/** Milisekunda liczona minutami — żeby liczby niżej dało się czytać wprost z opracowania. */
const MINUTA_MS = 60_000;

/** Milisekunda liczona godzinami. */
const GODZINA_MS = 60 * MINUTA_MS;

/**
 * Progi wyzwalania sugestii — rozdz. 3.4 opracowania, kolumna „wartość domyślna".
 *
 * Nazwy pól są nazwami wierszy tabeli, a nie skrótami: próg czytany w kodzie ma
 * mówić to samo, co próg czytany w opracowaniu.
 */
export const PROGI_AOD = {
  /**
   * Liczba sugestii ujawnianych samoczynnie w godzinie — 3.
   * Po przekroczeniu kolejne trafiają do listy oczekujących bez dymka.
   */
  liczbaDymkowNaGodzine: 3,

  /**
   * Odstęp między dymkami — 5 minut.
   * Dwie sugestie bliżej siebie łączą się w jedną pozycję zbiorczą.
   */
  odstepMiedzyDymkamiMs: 5 * MINUTA_MS,

  /**
   * Próg czasu oczekiwania zadania w kolejce — 15 minut.
   * Przekroczenie tworzy sugestię klasy „stan kolejki zadań".
   */
  oczekiwanieWKolejceMs: 15 * MINUTA_MS,

  /**
   * Próg liczby ponowień zadania — 2.
   * Trzecie ponowienie tego samego zadania tworzy sugestię o wadze wysokiej.
   */
  liczbaPonowien: 2,

  /**
   * Próg wypełnienia kolejki — 80% pojemności.
   * Przekroczenie tworzy sugestię klasy „stan kolejki zadań".
   */
  wypelnienieKolejki: 0.8,

  /**
   * Próg powtarzalności czynności ręcznej — 3 wystąpienia w karcie sesji.
   * Przekroczenie tworzy sugestię konfiguracji.
   */
  powtarzalnoscCzynnosci: 3,

  /**
   * Czas życia sugestii nieprzyjętej — 24 godziny.
   * Po upływie sugestia otrzymuje status odrzuconej i znika z listy oczekujących.
   */
  czasZyciaSugestiiMs: 24 * GODZINA_MS,
} as const;

/** Rodzaje wyciszenia czasowego — rozdz. 3.5, kolumna „zakres". */
export const WyciszenieCzasowe = {
  /** 15 minut. */
  Kwadrans: 'kwadrans',
  /** 1 godzina. */
  Godzina: 'godzina',
  /** Do końca dnia. */
  DoKoncaDnia: 'do-konca-dnia',
} as const;
export type WyciszenieCzasowe = (typeof WyciszenieCzasowe)[keyof typeof WyciszenieCzasowe];

/** Nazwa wyciszenia widziana przez Operatora w menu kebab. */
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
