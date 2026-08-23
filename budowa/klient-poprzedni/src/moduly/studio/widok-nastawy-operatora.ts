import type { JednostkaMiary } from './nastawy-strony';

/**
 * Nastawy widoku należące do Operatora — jawne, odwracalne i pamiętane.
 *
 * ── Skąd ten plik ───────────────────────────────────────────────────────────
 * Zasada ogólna zlecenia: gdzie da się zrobić dwojako i obie drogi mają sens,
 * wybór należy do Operatora jako jawne, odwracalne i pamiętane ustawienie, a nie
 * rozstrzygnięcie wykonawcy zapisane w kodzie. Ten plik jest jednym miejscem, w
 * którym takie wybory stoją — zakładki kontra podział powierzchni, przewijanie,
 * jednostka linijki, tryb adiustacji, tryb źródłowy, układ kartek. Rozsypane po
 * kilku plikach rozjechałyby się, a część z nich zostałaby na twardo.
 *
 * ── Dlaczego `localStorage`, a nie rdzeń ────────────────────────────────────
 * Dotyczą powierzchni na TYM urządzeniu i nie mają swojego bytu w kontrakcie:
 * `panel.sections.*` zapisuje układ sekcji panelu, nie skalę widoku ani jednostkę
 * linijki. Ten sam wzorzec nosi `motyw/motyw.ts`. Gdy kontrakt dostanie pozycję
 * na nastawy widoku dokumentu, zapis przejdzie do rdzenia bez zmiany wołaczy —
 * dlatego magazyn jest podawany, a nie brany z globalnej przestrzeni na sztywno.
 *
 * Awaria magazynu (tryb prywatny, osadzenie w ramce) zostawia wartości domyślne
 * i nie jest zgłaszana jako błąd: nastawa widoku nie jest powodem, żeby okno nie
 * wstało.
 */

/** Dwa równorzędne tryby pracy nad dwoma dokumentami. */
export type TrybDwochDokumentow = 'zakladki' | 'podzial';

/** Kierunek podziału powierzchni. */
export type KierunekPodzialu = 'pionowy' | 'poziomy';

/** Przewijanie powierzchni — ciągłe albo strona po stronie. */
export type TrybPrzewijania = 'ciagle' | 'strona-po-stronie';

/** Układ kartek na powierzchni. */
export type UkladKartek = 'jedna' | 'obok' | 'rozkladowka';

/** Nastawy widoku Operatora. */
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

/** Nastawy domyślne: zakładki, przewijanie ciągłe, jedna kartka, milimetry. */
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

/** Magazyn nastaw — pominięty znaczy `localStorage`, `null` znaczy „bez zapisu". */
export interface MagazynNastawWidoku {
  getItem(klucz: string): string | null;
  setItem(klucz: string, wartosc: string): void;
}

/** Klucz zapisu — nastawy widoku nie mieszają się z danymi sesji. */
const KLUCZ_ZAPISU = 'dn.studio.widok';

/** Granice udziału podziału — pole węższe niż 15 % nie jest polem pracy. */
const UDZIAL_DOLNY = 0.15;
const UDZIAL_GORNY = 0.85;

/** Granice skali widoku, wzorem pakietu biurowego. */
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
 * Czyta nastawy Operatora; brak zapisu i zapis uszkodzony dają domyślne.
 *
 * Każde pole jest sprawdzane osobno, a nie cały zapis naraz: zapis z wersji
 * wcześniejszej, w którym brakuje pola nowego, ma dać nastawę domyślną tego pola,
 * a nie unieważnić wszystkie pozostałe wybory Operatora.
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

/** Zapisuje nastawy Operatora; awaria zapisu niczego nie przerywa. */
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

/** Przycina udział podziału do granic pola pracy. */
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

/** Zdanie o nastawach widoku — do paska stanu i do objaśnień. */
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
