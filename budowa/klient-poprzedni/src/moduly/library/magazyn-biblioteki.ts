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
 * Forma prezentacji tego samego wykazu plików w oknie Library Explorer: siatka
 * miniatur jest domyślna, pozostałe cztery otwiera przełącznik widoku, bez
 * udziału rdzenia.
 */
export type WidokWykazu = 'siatka' | 'lista' | 'galeria' | 'os-czasu' | 'mapa';

/**
 * Tryb wyszukiwania Library Explorera: pełnotekstowy dopasowuje słowa w indeksie
 * repozytorium, semantyczny dopasowuje znaczenie we wskaźniku osadzeń, a
 * hybrydowy składa obie odpowiedzi po stronie okna.
 */
export type TrybWyszukiwania = 'pelnotekstowy' | 'semantyczny' | 'hybrydowy';

/** Trafienie wskaźnika znaczenia przypisane plikowi wykazu, wraz z fragmentem treści, na którym to trafienie się opiera. */
export interface TrafienieZnaczenia {
  /** Trafność w setnych, tak jak oddaje ją rdzeń; brak, gdy jej nie podał. */
  trafnosc: number | null;
  /** Fragment, na którym oparte jest trafienie — podstawa do sprawdzenia. */
  fragment: string;
}

/**
 * Zawężenie wykazu do wskazanego zbioru plików powstaje z raportu higieny
 * i z trafień wyszukiwania po znaczeniu; niesie własne zdanie, bo wykaz
 * zawężony wygląda jak wykaz krótki.
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
  /** Katalog modułów pochodzi z odczytu rdzenia; pusty wykaz znaczy nieodczytany albo odmówiony. */
  moduly: Module[];
  /** Powód, dla którego katalog modułów jest pusty; pusty napis = odczyt się udał. */
  powodModulow: string;
  /** Ster stoi w dwóch oknach i trzyma je zgodne; pusty napis znaczy moduł wytwórcy pliku. */
  modulDocelowy: string;
  /** Werdykty rdzenia o treści plików, osobno od zbioru, bo to nie jest pole kontraktu. */
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
 * Wciąga plik po własnej zmianie albo po zdarzeniu rdzenia: werdykt o treści
 * i trafienie wskaźnika znaczenia znikają, bo poprzednia odpowiedź rdzenia
 * przestaje opisywać zmieniony dokument.
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

/** Usuwa plik skasowany w rdzeniu wraz z jego zaznaczeniem, wskazaniem i zawężeniem wykazu, jeśli go dotyczyły. */
export function usunPlik(magazyn: MagazynBiblioteki, idPliku: string): void {
  magazyn.zbior = magazyn.zbior.filter((wpis) => wpis.id !== idPliku);
  magazyn.zaznaczenie = magazyn.zaznaczenie.filter((wpis) => wpis !== idPliku);
  magazyn.tresci.delete(idPliku);
  magazyn.znaczenia.delete(idPliku);
  // Zawężenie wskazujące skasowany plik przestałoby zgadzać się z wykazem po jego usunięciu.
  if (magazyn.zawezenie !== null) {
    const kody = magazyn.zawezenie.kody.filter((wpis) => wpis !== idPliku);
    magazyn.zawezenie = kody.length === 0 ? null : { ...magazyn.zawezenie, kody };
  }
  if (magazyn.wskazany === idPliku) magazyn.wskazany = magazyn.zbior[0]?.id ?? null;
  magazyn.oglos();
}

/** Odkłada odpowiedź rdzenia o treści pliku w magazynie i ogłasza zmianę wszystkim obserwującym ją oknom modułu. */
export function zapiszTresc(
  magazyn: MagazynBiblioteki,
  idPliku: string,
  stan: StanTresci,
): void {
  magazyn.tresci.set(idPliku, stan);
  magazyn.oglos();
}

/** Werdykt o treści pliku odczytany z magazynu stanu modułu; plik dotąd nieodpytany zwraca stan nieznana. */
export function odczytajTresc(magazyn: MagazynBiblioteki, idPliku: string): StanTresci {
  return magazyn.tresci.get(idPliku) ?? TRESC_NIEZNANA;
}

/** Zbiór wartości pól wielokrotnych zestawu plików, uporządkowany alfabetycznie według reguł języka polskiego. */
export function zbierz(
  pliki: readonly LibraryFile[],
  pole: (plik: LibraryFile) => readonly string[],
): readonly string[] {
  const znane = new Set<string>();
  for (const plik of pliki) for (const wartosc of pole(plik)) znane.add(wartosc);
  return [...znane].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl'));
}
