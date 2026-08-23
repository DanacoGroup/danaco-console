import type { LibraryFile, Module } from '../../../../shared/contract';
import { TRESC_NIEZNANA, type StanTresci } from './dostepnosc-tresci';

/**
 * Magazyn stanu modułu Library — same dane i jedno ogłoszenie zmiany.
 * Czynności, które ten stan zmieniają, stoją w `stan-biblioteki.ts`.
 *
 * Magazyn nie zna kontraktu ani rdzenia: nie woła komend, nie subskrybuje
 * zdarzeń i nie buduje elementów.
 */
export type FazaWykazu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

/**
 * Forma prezentacji tego samego wykazu plików w oknie Library Explorer.
 *
 * Pięć postaci opisuje dokumentacja modułu: siatka miniatur jest domyślna
 * (warstwa 1), pozostałe otwiera przełącznik widoku (warstwa 2). Wybór nie
 * dotyka rdzenia — przelicza prezentację tej samej odpowiedzi, więc mieszka
 * w pamięci modułu, nie w żądaniu.
 */
export type WidokWykazu = 'siatka' | 'lista' | 'galeria' | 'os-czasu' | 'mapa';

/**
 * Tryb wyszukiwania Library Explorera.
 *
 * `pelnotekstowy` idzie komendą `library.file.search` (dopasowanie słów
 * w indeksie repozytorium), `semantyczny` — komendą `knowledge.search`
 * w zakresie biblioteki (dopasowanie znaczenia we wskaźniku osadzeń).
 * `hybrydowy` nie jest komendą kontraktu: to złożenie obu odpowiedzi po stronie
 * okna, słowa przed znaczeniem.
 */
export type TrybWyszukiwania = 'pelnotekstowy' | 'semantyczny' | 'hybrydowy';

/** Trafienie wskaźnika znaczenia przypisane plikowi wykazu. */
export interface TrafienieZnaczenia {
  /** Trafność w setnych, tak jak oddaje ją rdzeń; brak, gdy jej nie podał. */
  trafnosc: number | null;
  /** Fragment, na którym oparte jest trafienie — podstawa do sprawdzenia. */
  fragment: string;
}

/**
 * Zawężenie wykazu do wskazanego zbioru plików.
 *
 * Powstaje z raportu higieny (duplikaty, pliki osierocone) i z trafień
 * wyszukiwania po znaczeniu, które nie zmieściły się w odczycie wykazu. Niesie
 * własne zdanie, bo wykaz zawężony wygląda jak wykaz krótki — a to dwie różne
 * rzeczy, i Operator ma widzieć, która zachodzi.
 */
export interface ZawezenieWykazu {
  kody: string[];
  opis: string;
}

export interface MagazynBiblioteki {
  zbior: LibraryFile[];
  faza: FazaWykazu;
  powod: string;
  sciezka: string;
  widok: WidokWykazu;
  tryb: TrybWyszukiwania;
  /** Zawężenie wykazu do zbioru wskazanego przez raport albo wyszukiwanie. */
  zawezenie: ZawezenieWykazu | null;
  /** Trafienia wskaźnika znaczenia po identyfikatorze pliku. */
  znaczenia: Map<string, TrafienieZnaczenia>;
  zaznaczenie: string[];
  wskazany: string | null;
  okno: string;
  /** Kolekcje założone w tej sesji; kontrakt nie ma komendy ich katalogu. */
  zalozone: string[];
  /**
   * Katalog modułów z `module.list` — pozycje stera modułu docelowego.
   *
   * Pusty wykaz znaczy „nieodczytany albo odmówiony"; powód stoi w polu obok.
   */
  moduly: Module[];
  /** Powód, dla którego katalog modułów jest pusty; pusty napis = odczyt się udał. */
  powodModulow: string;
  /**
   * Jedna nastawa modułu docelowego na cały moduł.
   *
   * Ster stoi w dwóch oknach — Library Explorer (czynność zbiorcza „Otwórz
   * w module źródłowym") i File Preview (ta sama czynność dla pojedynczego
   * pliku). Wspólne pole trzyma oba okna zgodne; osobne pokazywałyby po zmianie
   * dwie różne wartości. Pusty napis znaczy „moduł wytwórcy pliku"
   * (`sourceModuleId`).
   */
  modulDocelowy: string;
  /**
   * Werdykty rdzenia o treści plików, po identyfikatorze pliku.
   *
   * Osobno od zbioru, bo to nie jest pole kontraktu: `LibraryFile` nie niesie
   * żadnej wartości mówiącej, czy repozytorium ma treść (patrz
   * `dostepnosc-tresci.ts`).
   */
  tresci: Map<string, StanTresci>;
  /** Powiadamia widoki o zmianie pamięci. */
  oglos(): void;
  obserwuj(sluchacz: () => void): () => void;
  /** Zwalnia słuchaczy przy odpięciu modułu. */
  zamknij(): void;
}

export function utworzMagazyn(): MagazynBiblioteki {
  const sluchacze = new Set<() => void>();

  return {
    zbior: [],
    faza: 'spoczynek',
    powod: '',
    sciezka: '',
    // Siatka miniatur jest widokiem domyślnym warstwy 1; pozostałe cztery
    // otwiera przełącznik widoku.
    widok: 'siatka',
    tryb: 'pelnotekstowy',
    zawezenie: null,
    znaczenia: new Map(),
    zaznaczenie: [],
    wskazany: null,
    okno: '',
    zalozone: [],
    moduly: [],
    powodModulow: '',
    modulDocelowy: '',
    tresci: new Map(),

    oglos() {
      for (const sluchacz of [...sluchacze]) sluchacz();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
    },
  };
}

/**
 * Wciąga plik po własnej zmianie albo po zdarzeniu rdzenia.
 *
 * Werdykt o treści znika wraz ze zmianą pliku: dołożenie wersji, przywrócenie
 * wcześniejszej i zdarzenie `library.file.changed` z innego modułu zmieniają
 * dokument, więc poprzednia odpowiedź rdzenia przestaje go opisywać. Stan
 * „nieznana" jest tu bezpieczniejszy niż nieaktualna pewność. Z tego samego
 * powodu znika trafienie wskaźnika znaczenia: fragment pochodzi z treści
 * sprzed zmiany.
 */
export function wchlonPlik(magazyn: MagazynBiblioteki, plik: LibraryFile): void {
  const pozycja = magazyn.zbior.findIndex((wpis) => wpis.id === plik.id);
  magazyn.zbior =
    pozycja === -1
      ? [...magazyn.zbior, plik]
      : magazyn.zbior.map((wpis) => (wpis.id === plik.id ? plik : wpis));
  magazyn.tresci.delete(plik.id);
  magazyn.znaczenia.delete(plik.id);
  if (magazyn.wskazany === null) magazyn.wskazany = plik.id;
  magazyn.faza = 'gotowe';
  magazyn.oglos();
}

/** Usuwa plik skasowany w rdzeniu wraz z jego zaznaczeniem i wskazaniem. */
export function usunPlik(magazyn: MagazynBiblioteki, idPliku: string): void {
  magazyn.zbior = magazyn.zbior.filter((wpis) => wpis.id !== idPliku);
  magazyn.zaznaczenie = magazyn.zaznaczenie.filter((wpis) => wpis !== idPliku);
  magazyn.tresci.delete(idPliku);
  magazyn.znaczenia.delete(idPliku);
  // Zawężenie wskazujące skasowany plik przestałoby zgadzać się z wykazem:
  // licznik zbioru mówiłby o pozycji, której rdzeń już nie zna.
  if (magazyn.zawezenie !== null) {
    const kody = magazyn.zawezenie.kody.filter((wpis) => wpis !== idPliku);
    magazyn.zawezenie = kody.length === 0 ? null : { ...magazyn.zawezenie, kody };
  }
  if (magazyn.wskazany === idPliku) magazyn.wskazany = magazyn.zbior[0]?.id ?? null;
  magazyn.oglos();
}

/** Odkłada odpowiedź rdzenia o treści pliku i pokazuje ją oknom. */
export function zapiszTresc(
  magazyn: MagazynBiblioteki,
  idPliku: string,
  stan: StanTresci,
): void {
  magazyn.tresci.set(idPliku, stan);
  magazyn.oglos();
}

/** Werdykt o treści pliku; plik nieodpytany zwraca stan „nieznana". */
export function odczytajTresc(magazyn: MagazynBiblioteki, idPliku: string): StanTresci {
  return magazyn.tresci.get(idPliku) ?? TRESC_NIEZNANA;
}

/** Zbiór wartości pól wielokrotnych plików, uporządkowany po polsku. */
export function zbierz(
  pliki: readonly LibraryFile[],
  pole: (plik: LibraryFile) => readonly string[],
): readonly string[] {
  const znane = new Set<string>();
  for (const plik of pliki) for (const wartosc of pole(plik)) znane.add(wartosc);
  return [...znane].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl'));
}
