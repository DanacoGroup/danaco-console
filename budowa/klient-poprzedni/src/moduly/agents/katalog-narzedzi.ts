import {
  NARZEDZIA_MODELU,
  SEPARATOR_OBSZARU,
  type ToolDeclaration,
} from '../../../../shared/contract';

/** Katalog narzędzi modelu pogrupowany po przeznaczeniu, na podstawie kontraktu. */

/**
 * Pozycja katalogu: jedno narzędzie modelu wraz z grupą swojego przeznaczenia. Pozycja
 * niesie zarazem nazwę widzianą przez model, komendę kontraktu i zdanie o zastosowaniu,
 * więc okno wyboru narzędzi nie sięga po nic poza tym katalogiem.
 */
export interface PozycjaNarzedzia {
  /** Nazwa narzędzia widziana przez model — kod wpisywany ekspertowi. */
  nazwa: string;
  /** Komenda kontraktu wykonywana przez narzędzie. */
  komenda: string;
  /** Zdanie z kontraktu: czym narzędzie jest i kiedy po nie sięgnąć. */
  opis: string;
  /** Obszar komendy — grupa po przeznaczeniu, zarazem gałąź drzewa wyboru. */
  grupa: string;
}

/**
 * Obszar nazwy komendy, czyli człon przed pierwszym separatorem obszaru. Nazwa bez
 * separatora daje grupę równą całej nazwie, więc każda komenda ma grupę i katalog nie
 * potrzebuje gałęzi zastępczej na komendy bez obszaru.
 */
function grupaKomendy(komenda: string): string {
  return komenda.split(SEPARATOR_OBSZARU)[0] ?? komenda;
}

function zDeklaracji(deklaracja: ToolDeclaration): PozycjaNarzedzia {
  return {
    nazwa: deklaracja.name,
    komenda: deklaracja.command,
    opis: deklaracja.description,
    grupa: grupaKomendy(deklaracja.command),
  };
}

/**
 * Komplet narzędzi kontraktu w kolejności kontraktu.
 *
 * Kolejność zostaje kolejnością kontraktu — ten sam kontrakt zawsze daje ten
 * sam wykaz, więc drzewo nie przestawia się między odsłonami.
 */
export const KATALOG_NARZEDZI: readonly PozycjaNarzedzia[] =
  NARZEDZIA_MODELU.map(zDeklaracji);

/**
 * Nazwy grup w porządku alfabetycznym, bez powtórzeń. Alfabetycznie, a nie kolejnością
 * kontraktu, ponieważ grupa jest gałęzią drzewa, a gałęzi szuka się okiem po nazwie;
 * kolejność kontraktu zostaje wewnątrz gałęzi.
 */
export const GRUPY_NARZEDZI: readonly string[] = [
  ...new Set(KATALOG_NARZEDZI.map((pozycja) => pozycja.grupa)),
].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl'));

/**
 * Narzędzia jednej grupy, w kolejności katalogu. Grupa nieznana daje wykaz pusty zamiast
 * odmowy, ponieważ definicja eksperta może wskazywać obszar, którego kontrakt tego
 * klienta jeszcze nie zna.
 */
export function narzedziaGrupy(grupa: string): readonly PozycjaNarzedzia[] {
  return KATALOG_NARZEDZI.filter((pozycja) => pozycja.grupa === grupa);
}

/**
 * Wynik rozpoznania kodów wskazanych przez definicję eksperta, rozdzielony na kody nazw,
 * kody grup oraz kody nierozpoznane. Rozdział stoi w wyniku, żeby okno mogło nazwać
 * odrzut, zamiast po cichu pomijać kod, którego nie zrozumiało.
 */
export interface RozpoznanieKodow {
  /** Kody nazywające narzędzie wprost, w kolejności definicji. */
  nazwy: readonly string[];
  /** Kody nazywające grupę, w kolejności definicji. */
  grupy: readonly string[];
  /** Kody, które nie nazwały ani narzędzia, ani grupy. */
  nierozpoznane: readonly string[];
  /** Narzędzia wskazane kodami — bez powtórzeń, w kolejności katalogu. */
  narzedzia: readonly PozycjaNarzedzia[];
}

/**
 * Czyta kody eksperta tak, jak przeczyta je serwer narzędzi.
 *
 * Kody puste i powtórzone są pomijane przed rozpoznaniem — tak samo robi
 * `DefinicjaEksperta.Kody()`. Kod nierozpoznany wraca osobną listą, żeby odrzut
 * nie przeszedł niezauważony.
 */
export function rozpoznajKody(kody: readonly string[]): RozpoznanieKodow {
  const znaneNazwy = new Set(KATALOG_NARZEDZI.map((pozycja) => pozycja.nazwa));
  const znaneGrupy = new Set(GRUPY_NARZEDZI);

  const widziane = new Set<string>();
  const nazwy: string[] = [];
  const grupy: string[] = [];
  const nierozpoznane: string[] = [];

  for (const kod of kody) {
    if (kod === '' || widziane.has(kod)) continue;
    widziane.add(kod);
    if (znaneNazwy.has(kod)) nazwy.push(kod);
    else if (znaneGrupy.has(kod)) grupy.push(kod);
    else nierozpoznane.push(kod);
  }

  const wskazaneNazwy = new Set(nazwy);
  const wskazaneGrupy = new Set(grupy);
  const narzedzia = KATALOG_NARZEDZI.filter(
    (pozycja) => wskazaneNazwy.has(pozycja.nazwa) || wskazaneGrupy.has(pozycja.grupa),
  );

  return { nazwy, grupy, nierozpoznane, narzedzia };
}
