/** Pochodzenie wniesionego fragmentu opisuje, skąd fragment jest: adres, plik, wersja i moment wniesienia, zapisywane trwale przez rdzeń. */

/** Skąd fragment pochodzi: z pliku biblioteki, ze strony internetowej albo z migawki przeglądarki użytkownika. */
export const ZrodloWniesienia = {
  Biblioteka: 'biblioteka',
  Strona: 'strona',
  MigawkaPrzegladarki: 'migawka',
} as const;
export type ZrodloWniesienia = (typeof ZrodloWniesienia)[keyof typeof ZrodloWniesienia];

/** Jeden zapis pochodzenia wniesionego fragmentu, wraz z jego źródłem, wskazaniem i miejscem wstawienia w treści. */
export interface Pochodzenie {
  kod: string;
  zrodlo: ZrodloWniesienia;
  /** Nazwa źródła widoczna dla Operatora: nazwa pliku albo tytuł strony. */
  nazwa: string;
  /** Wskazanie dokładne: identyfikator pliku Biblioteki albo adres strony. */
  wskazanie: string;
  /** Wersja pliku Biblioteki albo czas pobrania strony; puste, gdy nieznana. */
  wersja: string;
  /** Suma kontrolna pliku, o ile Biblioteka ją oddała. */
  sumaKontrolna: string;
  /** Strona podglądu, z której fragment wzięto; zero znaczy „bez stron". */
  strona: number;
  /** Liczba znaków wniesionych. */
  znakow: number;
  /** Miejsce w treści, w które fragment wszedł. */
  wstawionoNaZnaku: number;
  czas: number;
}

/** Zbiór pochodzeń fragmentów wniesionych w bieżącej sesji okna, podręczny wobec zapisu trwałego rdzenia. */
export interface WykazPochodzen {
  wykaz(): readonly Pochodzenie[];
  /** Zapisuje pochodzenie i oddaje je. */
  zapisz(pochodzenie: Omit<Pochodzenie, 'kod' | 'czas'>): Pochodzenie;
  /** Zdejmuje wszystkie zapisy — po zamknięciu dokumentu. */
  wyczysc(): void;
}

export function utworzWykazPochodzen(): WykazPochodzen {
  const zapisy: Pochodzenie[] = [];
  let licznik = 0;

  return {
    wykaz: () => zapisy,

    zapisz(pochodzenie) {
      licznik += 1;
      const zapis: Pochodzenie = {
        ...pochodzenie,
        kod: `pochodzenie-${licznik}-${Date.now()}`,
        czas: Date.now(),
      };
      zapisy.push(zapis);
      return zapis;
    },

    wyczysc() {
      zapisy.length = 0;
    },
  };
}

/**
 * Wiersz pochodzenia wnoszony do treści dokumentu razem z fragmentem.
 *
 * Jest pełnym zdaniem, nie kodem: zakaz numeracji wymyślonej obowiązuje i tutaj,
 * więc nie ma tu ani `[1]`, ani `§`. Wiersz da się przeczytać bez klucza.
 */
export function osadzenieWierszPochodzenia(pochodzenie: Pochodzenie): string {
  const czesci: string[] = [];
  if (pochodzenie.zrodlo === ZrodloWniesienia.Biblioteka) {
    czesci.push(`Źródło: plik Biblioteki „${pochodzenie.nazwa}"`);
    if (pochodzenie.wskazanie !== '') czesci.push(`identyfikator ${pochodzenie.wskazanie}`);
    if (pochodzenie.wersja !== '') czesci.push(`wersja ${pochodzenie.wersja}`);
    if (pochodzenie.sumaKontrolna !== '') czesci.push(`suma kontrolna ${pochodzenie.sumaKontrolna}`);
    if (pochodzenie.strona > 0) czesci.push(`strona ${pochodzenie.strona}`);
  } else {
    czesci.push(`Źródło: strona „${pochodzenie.nazwa}"`);
    if (pochodzenie.wskazanie !== '') czesci.push(`adres ${pochodzenie.wskazanie}`);
    czesci.push(
      pochodzenie.zrodlo === ZrodloWniesienia.Strona
        ? 'pobrana i oczyszczona komendą studio.ingest.url'
        : 'migawka z modułu Browser',
    );
  }
  czesci.push(`pobrano ${new Date(pochodzenie.czas).toLocaleString('pl-PL')}`);
  czesci.push(`${pochodzenie.znakow} znaków`);
  return `${czesci.join(', ')}.`;
}

/** Zdanie opisujące wykaz pochodzeń fragmentów, wyświetlane w pasku stanu oraz w panelu redaktora tekstu. */
export function osadzenieOpiszPochodzenia(zapisy: readonly Pochodzenie[]): string {
  if (zapisy.length === 0) {
    return 'Do tego dokumentu nie wniesiono jeszcze ani jednego fragmentu z Biblioteki ani ze ' +
      'strony, więc pochodzenia nie ma czego dotyczyć.';
  }
  const zBiblioteki = zapisy.filter(
    (zapis) => zapis.zrodlo === ZrodloWniesienia.Biblioteka,
  ).length;
  const znakow = zapisy.reduce((suma, zapis) => suma + zapis.znakow, 0);
  return (
    `Fragmentów wniesionych ${zapisy.length} (z Biblioteki ${zBiblioteki}, ze stron ` +
    `${zapisy.length - zBiblioteki}), razem ${znakow} znaków. Każdy niesie wiersz pochodzenia ` +
    'w treści — ten przeżywa wydanie do formatu, który zapisu pochodzenia nie niesie. Zapis ' +
    'trwały prowadzi rdzeń: wniesienie oddaje pochodzenie, a studio.provenance.list oddaje ' +
    'pochodzenie fragmentów całego dokumentu. Ten wykaz jest podręczny i ginie z kartą.'
  );
}
