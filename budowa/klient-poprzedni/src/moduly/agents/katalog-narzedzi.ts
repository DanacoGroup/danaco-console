import {
  NARZEDZIA_MODELU,
  SEPARATOR_OBSZARU,
  type ToolDeclaration,
} from '../../../../shared/contract';

/**
 * Katalog narzędzi modelu pogrupowany po przeznaczeniu — jedyne miejsce
 * w kliencie, które wie, z czego składa się wyposażenie eksperta.
 *
 * Źródłem jest `NARZEDZIA_MODELU` z `shared/contract.ts`. Grupą jest obszar nazwy
 * komendy — człon przed pierwszym `SEPARATOR_OBSZARU` (`agent.list` → `agent`);
 * własny wykaz „narzędzie → grupa" byłby drugą prawdą, rozjeżdżającą się przy
 * pierwszym narzędziu dopisanym do kontraktu. Tę samą regułę stosuje `grupaKomendy`
 * w `server/internal/narzedzia/grupa.go`.
 *
 * Kod z pól `Agent.skillIds` i `Agent.connectorIds` czyta się dwiema drogami: jako
 * nazwę narzędzia (jedna pozycja) albo jako nazwę grupy (wszystkie pozycje obszaru),
 * tak jak rozstrzyga `server/internal/narzedzia/ekspert_wykaz.go`.
 */

/** Pozycja katalogu: jedno narzędzie modelu wraz z grupą swojego przeznaczenia. */
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

/** Obszar nazwy komendy; nazwa bez separatora daje grupę równą całej nazwie. */
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
 * Nazwy grup w porządku alfabetycznym, bez powtórzeń.
 *
 * Alfabetycznie, a nie kolejnością kontraktu: grupa jest gałęzią drzewa,
 * a gałęzi szuka się okiem po nazwie. Kolejność kontraktu zostaje tam, gdzie
 * niesie znaczenie — wewnątrz gałęzi. Tak samo rozstrzyga `Grupy` w `grupa.go`.
 */
export const GRUPY_NARZEDZI: readonly string[] = [
  ...new Set(KATALOG_NARZEDZI.map((pozycja) => pozycja.grupa)),
].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl'));

/** Narzędzia jednej grupy, w kolejności katalogu. Grupa nieznana daje pustkę. */
export function narzedziaGrupy(grupa: string): readonly PozycjaNarzedzia[] {
  return KATALOG_NARZEDZI.filter((pozycja) => pozycja.grupa === grupa);
}

/** Wynik rozpoznania kodów wskazanych przez definicję eksperta. */
export interface RozpoznanieKodow {
  /** Kody nazywające narzędzie wprost, w kolejności definicji. */
  nazwy: readonly string[];
  /** Kody nazywające grupę, w kolejności definicji. */
  grupy: readonly string[];
  /** Kody, które nie nazwały ani narzędzia, ani grupy. */
  nierozpoznane: readonly string[];
  /**
   * Narzędzia wskazane kodami — bez powtórzeń, w kolejności katalogu.
   *
   * Kolejność katalogu, a nie kodów: dwa kody wskazujące tę samą pozycję
   * (nazwa narzędzia i jego grupa naraz) dają ją raz. Tak samo przesiewa
   * `przesiej` w `ekspert_wykaz.go`.
   */
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
