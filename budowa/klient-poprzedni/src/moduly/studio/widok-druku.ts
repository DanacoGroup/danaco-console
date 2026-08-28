import type { MagazynNastawWidoku } from './widok-nastawy-operatora';
import { przytnijSkale } from './widok-skali';

/** Typ ZakresDruku nazywa zakres stron do wydruku dokumentu: wszystkie strony, stronę bieżącą albo zakres podany ręcznie. */
export type ZakresDruku = 'wszystkie' | 'biezaca' | 'podany';

/** Typ AdiustacjaDruku nazywa sposób ujęcia adiustacji na wydruku dokumentu: z adiustacją albo tekst po zmianach. */
export type AdiustacjaDruku = 'z-adiustacja' | 'po-zmianach';

/** Interfejs NastawyDruku niesie wszystkie nastawy wydruku dokumentu: zakres stron, liczbę kopii, tryb dwustronny, skalę i adiustację. */
export interface NastawyDruku {
  zakres: ZakresDruku;
  /** Pierwsza strona zakresu podanego. */
  odStrony: number;
  /** Ostatnia strona zakresu podanego. */
  doStrony: number;
  kopie: number;
  dwustronny: boolean;
  /** Skala wydruku w procentach. */
  skala: number;
  adiustacja: AdiustacjaDruku;
}

/** Funkcja domyslneNastawyDruku zwraca nastawy domyślne wydruku: wszystkie strony, jedna kopia, bez adiustacji. */
export function domyslneNastawyDruku(): NastawyDruku {
  return {
    zakres: 'wszystkie',
    odStrony: 1,
    doStrony: 1,
    kopie: 1,
    dwustronny: false,
    // Bez adiustacji domyślnie: najczęstszy wydruk to pismo do wysłania, korekta to wybór osobny.
    skala: 100,
    adiustacja: 'po-zmianach',
  };
}

/** Stała KLUCZ_DRUKU jest kluczem zapisu nastaw druku w magazynie; szybkie drukowanie bierze je bez pytania Operatora. */
const KLUCZ_DRUKU = 'dn.studio.druk';

function magazynDomyslny(): MagazynNastawWidoku | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}

/** Funkcja czytajNastawyDruku czyta ostatnie zapisane nastawy druku z magazynu; brak zapisu daje nastawy domyślne. */
export function czytajNastawyDruku(
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): NastawyDruku {
  const domyslne = domyslneNastawyDruku();
  if (magazyn === null) return domyslne;
  try {
    const surowy = magazyn.getItem(KLUCZ_DRUKU);
    if (surowy === null || surowy === '') return domyslne;
    const zapis = JSON.parse(surowy) as Record<string, unknown>;
    return {
      zakres:
        zapis['zakres'] === 'biezaca' || zapis['zakres'] === 'podany'
          ? zapis['zakres']
          : 'wszystkie',
      odStrony: liczba(zapis['odStrony'], 1),
      doStrony: liczba(zapis['doStrony'], 1),
      kopie: Math.max(1, Math.min(99, liczba(zapis['kopie'], 1))),
      dwustronny: zapis['dwustronny'] === true,
      skala: przytnijSkale(liczba(zapis['skala'], 100)),
      adiustacja: zapis['adiustacja'] === 'z-adiustacja' ? 'z-adiustacja' : 'po-zmianach',
    };
  } catch {
    return domyslne;
  }
}

/** Funkcja zapamietajNastawyDruku zapisuje nastawy druku w magazynie; awaria zapisu niczego w oknie nie przerywa. */
export function zapamietajNastawyDruku(
  nastawy: NastawyDruku,
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): void {
  if (magazyn === null) return;
  try {
    magazyn.setItem(KLUCZ_DRUKU, JSON.stringify(nastawy));
  } catch {
    // Nastawa druku nie jest powodem, żeby cokolwiek przerywać.
  }
}

function liczba(wartosc: unknown, domyslna: number): number {
  return typeof wartosc === 'number' && Number.isFinite(wartosc) ? wartosc : domyslna;
}

/**
 * Funkcja stronyDoDruku zwraca numery stron objętych nastawami, przycięte do liczby stron dokumentu; zakres podany odwrotnie zostaje odwrócony, a nie odrzucony.
 */
export function stronyDoDruku(
  nastawy: NastawyDruku,
  liczbaStron: number,
  kartkaBiezaca: number,
): number[] {
  const ile = Math.max(1, Math.floor(liczbaStron));
  if (nastawy.zakres === 'biezaca') {
    const numer = Math.min(Math.max(1, Math.round(kartkaBiezaca)), ile);
    return [numer];
  }
  if (nastawy.zakres === 'podany') {
    const od = Math.min(nastawy.odStrony, nastawy.doStrony);
    const doStrony = Math.max(nastawy.odStrony, nastawy.doStrony);
    const wynik: number[] = [];
    for (let numer = Math.max(1, Math.round(od)); numer <= Math.min(ile, Math.round(doStrony)); numer += 1) {
      wynik.push(numer);
    }
    return wynik;
  }
  return Array.from({ length: ile }, (_, numer) => numer + 1);
}

/** Funkcja opiszNastawyDruku zwraca zdanie o nastawach druku, wyświetlane jako podsumowanie przed naciśnięciem przycisku Drukuj. */
export function opiszNastawyDruku(nastawy: NastawyDruku, liczbaStron: number, kartka: number): string {
  const strony = stronyDoDruku(nastawy, liczbaStron, kartka);
  const zakres =
    nastawy.zakres === 'wszystkie'
      ? `wszystkie strony (${strony.length})`
      : nastawy.zakres === 'biezaca'
        ? `strona bieżąca (${strony[0] ?? 1})`
        : `strony ${strony.length === 0 ? 'poza dokumentem' : `${strony[0]}–${strony[strony.length - 1]}`}`;
  return (
    `${zakres} · kopii ${nastawy.kopie} · ${nastawy.dwustronny ? 'dwustronnie' : 'jednostronnie'} · ` +
    `skala ${nastawy.skala} % · ` +
    `${nastawy.adiustacja === 'z-adiustacja' ? 'Z adiustacją (zmiany i komentarze widoczne)' : 'tekst po zmianach, bez adiustacji'}`
  );
}

/**
 * Funkcja czyDrukDostepny sprawdza, czy droga druku jest w tym środowisku dostępna, zanim okno pokaże przycisk drukowania.
 */
export function czyDrukDostepny(): boolean {
  return typeof (globalThis as { print?: unknown }).print === 'function';
}

/** Stała POWOD_BRAKU_DRUKU niesie powód, którym okno odmawia druku, gdy droga druku nie jest w tym środowisku dostępna. */
export const POWOD_BRAKU_DRUKU =
  'Ta powłoka nie udostępnia okna drukarki, więc wydruk nie ma czym wyjść. Drukowanie jest ' +
  'czynnością Operatora, nie rdzenia: w kontrakcie nie ma ani jednej komendy drukowania, bo ' +
  'drukarka stoi na maszynie Operatora, a rdzeń na serwerze. Wydanie dokumentu do pliku jest ' +
  'czynnością INNĄ i działa niezależnie od tego braku.';

/** Interfejs OtoczenieDruku opisuje otoczenie druku: powierzchnię, która wie, jak przygotować kartki dokumentu do wydruku. */
export interface OtoczenieDruku {
  liczbaStron(): number;
  kartkaBiezaca(): number;
  /** Przygotowuje powierzchnię do wydruku i oddaje sposób przywrócenia jej stanu poprzedniego. */
  przygotuj(nastawy: NastawyDruku, strony: readonly number[]): () => void;
}

/** Interfejs WynikDruku niesie wynik próby wydruku dokumentu: czy wydruk się powiódł i zdanie opisujące jego skutek. */
export interface WynikDruku {
  udany: boolean;
  zdanie: string;
}

/**
 * Drukuje wedle nastaw: przygotowuje kartki, woła okno drukarki, przywraca stan.
 *
 * Przywrócenie idzie w `finally`, bo okno drukarki bywa przerwane przez Operatora
 * i wyjątek nie może zostawić dokumentu w postaci wydruku.
 */
export function wydrukuj(nastawy: NastawyDruku, otoczenie: OtoczenieDruku): WynikDruku {
  if (!czyDrukDostepny()) return { udany: false, zdanie: POWOD_BRAKU_DRUKU };
  const strony = stronyDoDruku(nastawy, otoczenie.liczbaStron(), otoczenie.kartkaBiezaca());
  if (strony.length === 0) {
    return {
      udany: false,
      zdanie:
        `Podany zakres stron leży poza dokumentem, który ma stron ${otoczenie.liczbaStron()}. ` +
        'Wydruk wstrzymany — drukowanie wszystkiego zamiast wskazanego zakresu byłoby ' +
        'zrobieniem czegoś innego niż polecenie.',
    };
  }
  const przywroc = otoczenie.przygotuj(nastawy, strony);
  try {
    (globalThis as { print: () => void }).print();
  } finally {
    przywroc();
  }
  zapamietajNastawyDruku(nastawy);
  return {
    udany: true,
    zdanie:
      `Wydruk oddany drukarce: ${opiszNastawyDruku(nastawy, otoczenie.liczbaStron(), otoczenie.kartkaBiezaca())}. ` +
      'Liczbę kopii i druk dwustronny zatwierdza okno drukarki systemu — strona ich nie narzuca, ' +
      'bo drukarka jednostronna nie stanie się dwustronną od nastawy w oknie.',
  };
}
