import type { MagazynNastawWidoku } from './widok-nastawy-operatora';
import { przytnijSkale } from './widok-skali';

/**
 * Drukowanie — czynność OPERATORA, nie rdzenia.
 *
 * ── Dlaczego w kliencie ─────────────────────────────────────────────────────
 * Rozstrzygnięcie Właściciela: „póki co tylko funkcja dla Operatora". W całym
 * kontrakcie nie ma ani jednej komendy drukowania i nie jest to przeoczenie —
 * drukarka stoi na maszynie Operatora, a rdzeń pracuje na serwerze. Rdzeń nie ma
 * czym drukować i nie udaje, że ma; tak samo odmawia uczciwie skaner
 * (`studio.ingest.device.scan`). Wydanie do pliku to CZYNNOŚĆ INNA i jedno nie
 * zastępuje drugiego.
 *
 * ── Drukuje się to, co pokazuje podgląd ─────────────────────────────────────
 * Nie surowy tekst. Wydruk różniący się od podglądu byłby usterką gorszą niż brak
 * funkcji, więc drukowanie bierze kartki powierzchni takimi, jakimi są: nośnik,
 * orientacja, marginesy, paginacja, nagłówek i stopka. Ten plik nie rysuje własnej
 * kartki i nie zna DOM powierzchni — składa nastawy i woła drogę druku, a
 * przygotowanie kartek należy do powierzchni.
 *
 * ── Czego nie da się tu obiecać ─────────────────────────────────────────────
 * Liczby kopii i druku dwustronnego nie rozstrzyga strona, tylko okno drukarki
 * systemu. Nastawy są więc przenoszone jako **życzenie wpisane w podsumowanie**,
 * a okno mówi wprost, że zatwierdza je drukarka — obietnica, że strona ustawi
 * dupleks, byłaby nieprawdą przy pierwszej drukarce jednostronnej.
 */

/** Zakres stron do wydruku. */
export type ZakresDruku = 'wszystkie' | 'biezaca' | 'podany';

/** Co robić z adiustacją na wydruku. */
export type AdiustacjaDruku = 'z-adiustacja' | 'po-zmianach';

/** Nastawy druku. */
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

/** Nastawy domyślne: wszystkie strony, jedna kopia, bez adiustacji. */
export function domyslneNastawyDruku(): NastawyDruku {
  return {
    zakres: 'wszystkie',
    odStrony: 1,
    doStrony: 1,
    kopie: 1,
    dwustronny: false,
    // Bez adiustacji, bo najczęstszy wydruk to pismo do wysłania. Korekta
    // z adiustacją jest wyborem świadomym i stoi jedno naciśnięcie obok.
    skala: 100,
    adiustacja: 'po-zmianach',
  };
}

/** Klucz zapisu nastaw druku — szybkie drukowanie bierze je bez pytania. */
const KLUCZ_DRUKU = 'dn.studio.druk';

function magazynDomyslny(): MagazynNastawWidoku | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}

/** Czyta ostatnie nastawy druku; brak zapisu daje domyślne. */
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

/** Zapisuje nastawy druku; awaria zapisu niczego nie przerywa. */
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
 * Numery stron objętych nastawami, przycięte do liczby stron dokumentu.
 *
 * Zakres podany odwrotnie (od 8 do 3) jest odwracany, a nie odrzucany: Operator
 * miał na myśli strony od trzeciej do ósmej i odmowa byłaby tu formalizmem. Zakres
 * całkowicie poza dokumentem oddaje wykaz pusty — i wołający ma wtedy odmówić
 * wydruku, zamiast drukować wszystko.
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

/** Zdanie o nastawach druku — podsumowanie przed naciśnięciem „Drukuj". */
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
 * Czy droga druku jest w tym środowisku dostępna.
 *
 * Sprawdzane PRZED próbą, bo przycisk, który po naciśnięciu milczy, jest gorszy
 * od nazwanego braku. Powłoka bez okna drukarki (osadzenie w ramce bez uprawnień,
 * środowisko sprawdzianu) nie ma funkcji `print` i wtedy okno mówi to wprost.
 */
export function czyDrukDostepny(): boolean {
  return typeof (globalThis as { print?: unknown }).print === 'function';
}

/** Powód, którym okno odmawia druku, gdy droga nie jest dostępna. */
export const POWOD_BRAKU_DRUKU =
  'Ta powłoka nie udostępnia okna drukarki, więc wydruk nie ma czym wyjść. Drukowanie jest ' +
  'czynnością Operatora, nie rdzenia: w kontrakcie nie ma ani jednej komendy drukowania, bo ' +
  'drukarka stoi na maszynie Operatora, a rdzeń na serwerze. Wydanie dokumentu do pliku jest ' +
  'czynnością INNĄ i działa niezależnie od tego braku.';

/** Otoczenie druku — powierzchnia, która wie, jak przygotować kartki. */
export interface OtoczenieDruku {
  liczbaStron(): number;
  kartkaBiezaca(): number;
  /**
   * Przygotowuje powierzchnię do wydruku i oddaje sposób przywrócenia jej.
   *
   * Przywrócenie jest oddawane, a nie domyślane: wydruk zostawiający dokument
   * w postaci przygotowanej do druku byłby usterką widoczną dopiero wtedy, gdy
   * Operator wróci do pisania.
   */
  przygotuj(nastawy: NastawyDruku, strony: readonly number[]): () => void;
}

/** Wynik próby wydruku. */
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
