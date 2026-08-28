import type { JednostkaMiary } from './nastawy-strony';

/** Typ TrybDwochDokumentow nazywa dwa równorzędne tryby pracy nad dwoma dokumentami: zakładki i podział powierzchni. */
export type TrybDwochDokumentow = 'zakladki' | 'podzial';

/** Typ KierunekPodzialu nazywa kierunek podziału powierzchni pracy nad dwoma dokumentami: pionowy albo poziomy. */
export type KierunekPodzialu = 'pionowy' | 'poziomy';

/** Typ TrybPrzewijania nazywa sposób przewijania powierzchni dokumentu przy pracy: ciągłe albo strona po stronie. */
export type TrybPrzewijania = 'ciagle' | 'strona-po-stronie';

/** Typ UkladKartek nazywa układ kartek na powierzchni dokumentu widoku: jedną, obok siebie albo rozkładówkę. */
export type UkladKartek = 'jedna' | 'obok' | 'rozkladowka';

/** Interfejs NastawyOperatoraWidoku niesie wszystkie nastawy widoku należące do Operatora: tryb dokumentów, podział, przewijanie, układ kartek, jednostkę, linijki, tryb źródłowy i skalę. */
export interface NastawyOperatoraWidoku {
  /** Zakładki albo podział powierzchni — dwa tryby, żaden zapasowy. */
  trybDokumentow: TrybDwochDokumentow;
  kierunekPodzialu: KierunekPodzialu;
  /** Udział pierwszego pola w podziale, od 0,15 do 0,85. */
  udzialPodzialu: number;
  przewijanie: TrybPrzewijania;
  ukladKartek: UkladKartek;
  /** Liczba kartek w rzędzie przy układzie „obok siebie". */
  kartekWRzedzie: number;
  jednostka: JednostkaMiary;
  linijkiWidoczne: boolean;
  graniceMarginesow: boolean;
  /** Tryb źródłowy ze znacznikami jest PRZEŁĄCZNIKIEM, nie widokiem domyślnym. */
  trybZrodlowy: boolean;
  /** Skala widoku w procentach — nastawa ostatnia, pamiętana między wejściami. */
  skala: number;
}

/** Funkcja domyslneNastawyWidoku zwraca nastawy domyślne: zakładki, przewijanie ciągłe, jedna kartka w rzędzie, jednostka milimetrowa. */
export function domyslneNastawyWidoku(): NastawyOperatoraWidoku {
  return {
    trybDokumentow: 'zakladki',
    kierunekPodzialu: 'pionowy',
    udzialPodzialu: 0.5,
    przewijanie: 'ciagle',
    ukladKartek: 'jedna',
    kartekWRzedzie: 2,
    jednostka: 'mm',
    linijkiWidoczne: true,
    graniceMarginesow: true,
    trybZrodlowy: false,
    skala: 100,
  };
}

/** Interfejs MagazynNastawWidoku opisuje magazyn nastaw widoku; pominięty magazyn znaczy localStorage, a wartość null znaczy brak zapisu. */
export interface MagazynNastawWidoku {
  getItem(klucz: string): string | null;
  setItem(klucz: string, wartosc: string): void;
}

/** Stała KLUCZ_ZAPISU jest kluczem magazynu, pod którym nastawy widoku nie mieszają się z danymi sesji dokumentu. */
const KLUCZ_ZAPISU = 'dn.studio.widok';

/** Stałe UDZIAL_DOLNY i UDZIAL_GORNY wyznaczają granice udziału podziału powierzchni; pole węższe niż 15 procent nie jest polem pracy. */
const UDZIAL_DOLNY = 0.15;
const UDZIAL_GORNY = 0.85;

/** Stałe SKALA_DOLNA i SKALA_GORNA wyznaczają granice skali widoku, wzorowane na granicach pakietu biurowego. */
export const SKALA_DOLNA = 10;
export const SKALA_GORNA = 500;

function magazynDomyslny(): MagazynNastawWidoku | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}

/**
 * Funkcja czytajNastawyWidoku czyta nastawy Operatora z magazynu; brak zapisu, zapis uszkodzony albo pole niezgodne z typem dają wartość domyślną tego pola.
 */
export function czytajNastawyWidoku(
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): NastawyOperatoraWidoku {
  const domyslne = domyslneNastawyWidoku();
  if (magazyn === null) return domyslne;
  let zapis: unknown;
  try {
    const surowy = magazyn.getItem(KLUCZ_ZAPISU);
    if (surowy === null || surowy === '') return domyslne;
    zapis = JSON.parse(surowy);
  } catch {
    return domyslne;
  }
  if (typeof zapis !== 'object' || zapis === null) return domyslne;
  const pola = zapis as Record<string, unknown>;

  return {
    trybDokumentow: wybor(pola['trybDokumentow'], ['zakladki', 'podzial'], domyslne.trybDokumentow),
    kierunekPodzialu: wybor(
      pola['kierunekPodzialu'],
      ['pionowy', 'poziomy'],
      domyslne.kierunekPodzialu,
    ),
    udzialPodzialu: liczba(pola['udzialPodzialu'], UDZIAL_DOLNY, UDZIAL_GORNY, domyslne.udzialPodzialu),
    przewijanie: wybor(pola['przewijanie'], ['ciagle', 'strona-po-stronie'], domyslne.przewijanie),
    ukladKartek: wybor(pola['ukladKartek'], ['jedna', 'obok', 'rozkladowka'], domyslne.ukladKartek),
    kartekWRzedzie: Math.round(liczba(pola['kartekWRzedzie'], 1, 8, domyslne.kartekWRzedzie)),
    jednostka: wybor(pola['jednostka'], ['mm', 'cal'], domyslne.jednostka),
    linijkiWidoczne: prawda(pola['linijkiWidoczne'], domyslne.linijkiWidoczne),
    graniceMarginesow: prawda(pola['graniceMarginesow'], domyslne.graniceMarginesow),
    trybZrodlowy: prawda(pola['trybZrodlowy'], domyslne.trybZrodlowy),
    skala: liczba(pola['skala'], SKALA_DOLNA, SKALA_GORNA, domyslne.skala),
  };
}

/** Funkcja zapamietajNastawyWidoku zapisuje nastawy Operatora w magazynie; awaria zapisu niczego w oknie nie przerywa. */
export function zapamietajNastawyWidoku(
  nastawy: NastawyOperatoraWidoku,
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): void {
  if (magazyn === null) return;
  try {
    magazyn.setItem(KLUCZ_ZAPISU, JSON.stringify(nastawy));
  } catch {
    // Nastawa widoku nie jest powodem, żeby cokolwiek przerywać.
  }
}

function wybor<T extends string>(wartosc: unknown, dozwolone: readonly T[], domyslna: T): T {
  return typeof wartosc === 'string' && (dozwolone as readonly string[]).includes(wartosc)
    ? (wartosc as T)
    : domyslna;
}

function liczba(wartosc: unknown, dol: number, gora: number, domyslna: number): number {
  if (typeof wartosc !== 'number' || !Number.isFinite(wartosc)) return domyslna;
  return Math.min(Math.max(wartosc, dol), gora);
}

function prawda(wartosc: unknown, domyslna: boolean): boolean {
  return typeof wartosc === 'boolean' ? wartosc : domyslna;
}

/** Funkcja przytnijUdzialPodzialu przycina udział podziału powierzchni do granic dolnej i górnej pola pracy. */
export function przytnijUdzialPodzialu(udzial: number): number {
  if (!Number.isFinite(udzial)) return 0.5;
  return Math.min(Math.max(udzial, UDZIAL_DOLNY), UDZIAL_GORNY);
}

/**
 * Nastawy widoku wraz z zapisem — jedno miejsce, przez które okno je przestawia.
 *
 * Każde przestawienie zapisuje całość: nastawy widoku są jednym zapisem, a nie
 * jedenastoma, więc częściowy zapis nie może się rozjechać z odczytem.
 */
export interface PamiecNastawWidoku {
  biezace(): NastawyOperatoraWidoku;
  /** Przestawia wskazane pola i zapisuje całość; oddaje nastawy po zmianie. */
  przestaw(zmiana: Partial<NastawyOperatoraWidoku>): NastawyOperatoraWidoku;
}

export function utworzPamiecNastawWidoku(
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): PamiecNastawWidoku {
  let biezace = czytajNastawyWidoku(magazyn);
  return {
    biezace: () => biezace,
    przestaw(zmiana) {
      biezace = {
        ...biezace,
        ...zmiana,
        ...(zmiana.udzialPodzialu === undefined
          ? {}
          : { udzialPodzialu: przytnijUdzialPodzialu(zmiana.udzialPodzialu) }),
      };
      zapamietajNastawyWidoku(biezace, magazyn);
      return biezace;
    },
  };
}

/** Funkcja opiszNastawyWidoku zwraca zdanie o bieżących nastawach widoku, przeznaczone do paska stanu i do objaśnień. */
export function opiszNastawyWidoku(nastawy: NastawyOperatoraWidoku): string {
  const uklad =
    nastawy.ukladKartek === 'rozkladowka'
      ? 'rozkładówka'
      : nastawy.ukladKartek === 'obok'
        ? `${nastawy.kartekWRzedzie} kartki w rzędzie`
        : 'jedna kartka';
  return (
    `${nastawy.trybDokumentow === 'podzial' ? 'podział powierzchni' : 'zakładki'} · ${uklad} · ` +
    `przewijanie ${nastawy.przewijanie === 'ciagle' ? 'ciągłe' : 'strona po stronie'} · ` +
    `skala ${nastawy.skala} % · linijki ${nastawy.linijkiWidoczne ? 'widoczne' : 'ukryte'} ` +
    `(${nastawy.jednostka === 'cal' ? 'cale' : 'milimetry'})` +
    `${nastawy.trybZrodlowy ? ' · tryb źródłowy włączony' : ''}`
  );
}
