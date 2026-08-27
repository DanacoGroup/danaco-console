/** Siedem kategorii operacji kontekstowych Tools Panel, każda z identyfikatorem podawanym komendzie kontekstowej rdzenia. */

/** Jedna operacja kontekstowa w panelu: identyfikator akcji podawany rdzeniowi i jej nazwa widoczna w panelu. */
export interface OperacjaKontekstowa {
  /** Identyfikator akcji podawany komendzie `studio.contextual.op`. */
  id: string;
  /** Nazwa widoczna dla Operatora. */
  nazwa: string;
}

/** Kategoria panelu operacji kontekstowych wraz z kodem, nazwą widoczną i wykazem należących do niej operacji. */
export interface KategoriaOperacji {
  kod: string;
  nazwa: string;
  operacje: readonly OperacjaKontekstowa[];
}

/** Buduje operacje kategorii, nadając każdej identyfikator złożony z nazwy kategorii i kodu tej pozycji. */
function operacje(
  kategoria: string,
  pozycje: readonly (readonly [string, string])[],
): OperacjaKontekstowa[] {
  return pozycje.map(([kod, nazwa]) => ({ id: `studio.${kategoria}.${kod}`, nazwa }));
}

export const KATEGORIE_OPERACJI: readonly KategoriaOperacji[] = [
  {
    kod: 'korekta',
    nazwa: 'Korekta i jakość',
    operacje: operacje('korekta', [
      ['ortografia', 'Korekta ortograficzno-gramatyczna'],
      ['interpunkcja', 'Korekta interpunkcji'],
      ['terminologia', 'Spójność terminologii'],
      ['powtorzenia', 'Wykrycie powtórzeń'],
      ['czytelnosc', 'Kontrola czytelności'],
    ]),
  },
  {
    kod: 'styl',
    nazwa: 'Przekształcenia stylu',
    operacje: operacje('styl', [
      ['rejestr', 'Zmiana stylu (formalny/nieformalny/techniczny/marketingowy)'],
      ['ton', 'Zmiana tonu (neutralny/perswazyjny/empatyczny)'],
      ['skrocenie', 'Skrócenie'],
      ['rozwiniecie', 'Rozwinięcie treści'],
      ['uproszczenie', 'Uproszczenie języka'],
      ['podniesienie', 'Podniesienie rejestru'],
    ]),
  },
  {
    kod: 'streszczenie',
    nazwa: 'Streszczenia',
    operacje: operacje('streszczenie', [
      ['jednozdaniowe', 'Streszczenie jednozdaniowe'],
      ['akapitowe', 'Streszczenie akapitowe'],
      ['tezy', 'Wypunktowanie tez'],
      ['dzialania', 'Lista działań'],
      ['tytuly', 'Tytuł i podtytuły'],
      ['meta', 'Meta-opis'],
    ]),
  },
  {
    kod: 'tlumaczenie',
    nazwa: 'Tłumaczenia',
    operacje: operacje('tlumaczenie', [
      ['zakres', 'Tłumaczenie zaznaczenia lub całości'],
      ['dwujezyczne', 'Wersja dwujęzyczna równoległa'],
    ]),
  },
  {
    kod: 'struktura',
    nazwa: 'Przepisywanie strukturalne',
    operacje: operacje('struktura', [
      ['akapit-lista', 'Akapit na listę'],
      ['lista-tabela', 'Lista na tabelę'],
      ['spis', 'Spis treści'],
      ['zarzadcze', 'Streszczenie zarządcze'],
    ]),
  },
  {
    kod: 'wzbogacenie',
    nazwa: 'Wzbogacanie',
    operacje: operacje('wzbogacenie', [
      ['kontekst', 'Rozwinięcie o kontekst'],
      ['przyklady', 'Przykłady'],
      ['pytania', 'Pytania kontrolne'],
      ['kontrargumenty', 'Kontrargumenty'],
    ]),
  },
  {
    kod: 'weryfikacja',
    nazwa: 'Weryfikacja',
    operacje: operacje('weryfikacja', [
      ['zrodla', 'Wskazanie fragmentów wymagających źródła'],
      ['niescislosci', 'Oznaczenie nieścisłości'],
      ['library', 'Porównanie z materiałem z Library'],
    ]),
  },
];

/** Operacje pływaka kontekstowego stanowiące stan początkowy, dopóki rzeczywiste użycie nie ustali własnej kolejności. */
export const OPERACJE_PASKA: readonly OperacjaKontekstowa[] = [
  { id: 'studio.korekta.ortografia', nazwa: 'Korekta' },
  { id: 'studio.styl.rejestr', nazwa: 'Przepisz' },
  { id: 'studio.streszczenie.akapitowe', nazwa: 'Streść' },
  { id: 'studio.styl.ton', nazwa: 'Styl' },
];

/** Wszystkie operacje wykazu złożone w jeden ciąg, wspólny dla pływaka, menu pełnego, wiersza polecenia i narzędzi modelu. */
export const WSZYSTKIE_OPERACJE: readonly OperacjaKontekstowa[] = KATEGORIE_OPERACJI.flatMap(
  (kategoria) => kategoria.operacje,
);

/** Zwraca nazwę operacji wraz z nazwą jej kategorii; zwraca undefined, gdy identyfikator jest spoza wykazu operacji. */
export function nazwaOperacji(idAkcji: string): string | undefined {
  for (const kategoria of KATEGORIE_OPERACJI) {
    const operacja = kategoria.operacje.find((pozycja) => pozycja.id === idAkcji);
    if (operacja !== undefined) return `${operacja.nazwa} · ${kategoria.nazwa}`;
  }
  return undefined;
}

/** Zwraca kategorię, w której stoi operacja; zwraca undefined, gdy identyfikator jest spoza wykazu operacji. */
export function kategoriaOperacji(idAkcji: string): KategoriaOperacji | undefined {
  return KATEGORIE_OPERACJI.find((kategoria) =>
    kategoria.operacje.some((pozycja) => pozycja.id === idAkcji),
  );
}

/** Operacje pasujące do wpisanej frazy, dopasowane w nazwie operacji, w jej identyfikatorze albo w nazwie kategorii. */
export function operacjePoFrazie(fraza: string): OperacjaKontekstowa[] {
  const szukane = fraza.trim().toLowerCase();
  if (szukane === '') return [...WSZYSTKIE_OPERACJE];
  return KATEGORIE_OPERACJI.flatMap((kategoria) =>
    kategoria.operacje.filter(
      (operacja) =>
        operacja.nazwa.toLowerCase().includes(szukane) ||
        operacja.id.toLowerCase().includes(szukane) ||
        kategoria.nazwa.toLowerCase().includes(szukane),
    ),
  );
}
