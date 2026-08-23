/**
 * Siedem kategorii operacji kontekstowych Tools Panel.
 *
 * Pozycje tego wykazu nie są wierszami rejestru akcji: `action.list` dla zasięgu
 * modułu Studio oddaje wyłącznie komendy okna komunikacji, a operacji
 * redakcyjnych w rejestrze nie ma. Komenda `studio.contextual.op` ma uchwyt
 * i wychodzi do kanału modelu okna; brakuje samych wierszy katalogu, które
 * nazwałyby poszczególne operacje. Dlatego panel pokazuje przy każdej pozycji,
 * skąd ona jest — z rejestru rdzenia albo z tego wykazu.
 *
 * Pozycja wykazu jest identyfikatorem akcji podawanym komendzie
 * `studio.contextual.op`, więc niczego nie udaje: naciśnięcie wychodzi do rdzenia
 * i wraca jego odpowiedzią albo odmową.
 */

/** Jedna operacja kontekstowa: identyfikator akcji i jej nazwa w panelu. */
export interface OperacjaKontekstowa {
  /** Identyfikator akcji podawany komendzie `studio.contextual.op`. */
  id: string;
  /** Nazwa widoczna dla Operatora. */
  nazwa: string;
}

/** Kategoria panelu akcji wraz z jej operacjami. */
export interface KategoriaOperacji {
  kod: string;
  nazwa: string;
  operacje: readonly OperacjaKontekstowa[];
}

/** Buduje operacje kategorii, nadając im identyfikatory `studio.<kategoria>.<kod>`. */
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

/**
 * Operacje pływaka kontekstowego, dopóki użycie nie powie własnej kolejności.
 *
 * Rozstrzygnięcie Właściciela: na wierzchu pływaka stoją czynności najczęstsze
 * **wedle użycia, nie wedle domysłu**. Użycia w chwili pierwszego uruchomienia
 * jednak nie ma, a pływak bez ani jednej czynności byłby pusty — ten wykaz jest
 * więc **stanem początkowym**, który zastępuje pierwsze użycie Operatora.
 * Kolejność liczy `przybornik-uzycie.ts`; tutaj stoi tylko punkt startu.
 */
export const OPERACJE_PASKA: readonly OperacjaKontekstowa[] = [
  { id: 'studio.korekta.ortografia', nazwa: 'Korekta' },
  { id: 'studio.styl.rejestr', nazwa: 'Przepisz' },
  { id: 'studio.streszczenie.akapitowe', nazwa: 'Streść' },
  { id: 'studio.styl.ton', nazwa: 'Styl' },
];

/**
 * Wszystkie operacje wykazu w jednym ciągu — do wyszukiwania po nazwie
 * i podpowiedzi w wierszu polecenia.
 *
 * Wykaz jest jeden dla wszystkich czterech dróg: pływaka, menu pełnego, wiersza
 * polecenia i narzędzi modelu. Druga kopia rozjechałaby się z pierwszą przy
 * pierwszym dołożeniu operacji.
 */
export const WSZYSTKIE_OPERACJE: readonly OperacjaKontekstowa[] = KATEGORIE_OPERACJI.flatMap(
  (kategoria) => kategoria.operacje,
);

/** Nazwa operacji wraz z kategorią; `undefined`, gdy identyfikator jest spoza wykazu. */
export function nazwaOperacji(idAkcji: string): string | undefined {
  for (const kategoria of KATEGORIE_OPERACJI) {
    const operacja = kategoria.operacje.find((pozycja) => pozycja.id === idAkcji);
    if (operacja !== undefined) return `${operacja.nazwa} · ${kategoria.nazwa}`;
  }
  return undefined;
}

/** Kategoria, w której stoi operacja; `undefined` dla identyfikatora spoza wykazu. */
export function kategoriaOperacji(idAkcji: string): KategoriaOperacji | undefined {
  return KATEGORIE_OPERACJI.find((kategoria) =>
    kategoria.operacje.some((pozycja) => pozycja.id === idAkcji),
  );
}

/**
 * Operacje pasujące do frazy — dopasowanie w nazwie, w identyfikatorze
 * i w nazwie kategorii.
 *
 * Fraza pusta oddaje wykaz w całości, a nie pustkę: pole szukania niewypełnione
 * nie jest zawężeniem do zera.
 */
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
