import type { RodzajBloku } from './zapis-formatowany';

/**
 * Nastawy wizualne treści — krój, stopień, interlinia, wcięcia, odstępy
 * akapitowe, wyrównanie i barwa.
 *
 * ── Czego kontrakt nie niesie ───────────────────────────────────────────────
 * `StudioDocument` ma treść, tytuł, format i wersję. Pola stylów nie ma i nie ma
 * komendy, którą styl dojechałby do rdzenia: `studio.document.save` przyjmuje
 * `documentId`, `content`, `title` i `createVersion`. Nastawy z tego pliku żyją
 * więc PRZEZ SESJĘ OKNA i giną z jej zamknięciem — okno mówi to Operatorowi
 * wprost, zamiast udawać zapis, którego nie ma. Trwałe jest wyłącznie to, co
 * jest składnią treści (styl nazwany bloku, pogrubienie, kursywa, podkreślenie,
 * przekreślenie, lista, tabela, podział strony) oraz nastawy strony, bo te mają
 * w kontrakcie `StudioPageSetup` i profil wydania.
 *
 * Plik nie zna okna: oddaje wartości i zdania, a przypięcie ich do elementu
 * należy do powierzchni dokumentu.
 */

/** Wyrównanie akapitu. */
export type Wyrownanie = 'lewo' | 'srodek' | 'prawo' | 'obustronne';

/** Kroje oferowane w oknie — nazwy rodzin obecnych w motywie i w systemie. */
export const KROJE: readonly { wartosc: string; etykieta: string }[] = [
  { wartosc: 'var(--dn-ff-tekst)', etykieta: 'Krój tekstowy motywu' },
  { wartosc: 'var(--dn-ff-naglowek)', etykieta: 'Krój nagłówkowy motywu' },
  { wartosc: 'var(--dn-ff-mono)', etykieta: 'Krój o stałej szerokości' },
  { wartosc: 'Georgia, serif', etykieta: 'Szeryfowy (Georgia)' },
  { wartosc: 'Arial, Helvetica, sans-serif', etykieta: 'Bezszeryfowy (Arial)' },
  { wartosc: '"Times New Roman", serif', etykieta: 'Times New Roman' },
];

/** Stopnie pisma w punktach typograficznych. */
export const STOPNIE: readonly number[] = [8, 9, 10, 11, 12, 14, 16, 18, 24, 36];

/** Nastawy wizualne całego dokumentu. */
export interface NastawyWizualne {
  krój: string;
  /** Stopień pisma w punktach typograficznych. */
  stopien: number;
  /** Interlinia jako krotność stopnia. */
  interlinia: number;
  /** Wcięcie pierwszego wiersza akapitu w milimetrach. */
  wciecieMm: number;
  /** Odstęp przed i po akapicie w milimetrach. */
  odstepMm: number;
}

/** Nastawy akapitu — te, które Operator daje pojedynczemu blokowi. */
export interface NastawyAkapitu {
  wyrownanie: Wyrownanie;
  /** Barwa pisma; puste znaczy „barwa tekstu motywu". */
  barwa: string;
}

/** Nastawy domyślne: krój tekstowy motywu, 11 punktów, interlinia 1,5. */
export function domyslneNastawy(): NastawyWizualne {
  return {
    krój: 'var(--dn-ff-tekst)',
    stopien: 11,
    interlinia: 1.5,
    wciecieMm: 0,
    odstepMm: 3,
  };
}

/** Zdanie o nastawach wizualnych dla paska stanu. */
export function opiszNastawy(nastawy: NastawyWizualne): string {
  const krój = KROJE.find((pozycja) => pozycja.wartosc === nastawy.krój)?.etykieta ?? nastawy.krój;
  return (
    `${krój} · ${nastawy.stopien} pt · interlinia ${nastawy.interlinia} · ` +
    `wcięcie ${nastawy.wciecieMm} mm · odstęp akapitowy ${nastawy.odstepMm} mm`
  );
}

/** Powód, dla którego nastawy wizualne nie dojeżdżają do rdzenia. */
export const POWOD_NIETRWALOSCI =
  'Krój, stopień, interlinia, wcięcie, odstęp akapitowy, wyrównanie i barwa są nastawami WIDOKU ' +
  'i żyją przez tę sesję okna. Kontrakt nie ma pola stylu w dokumencie ani komendy, którą styl ' +
  'dojechałby do rdzenia — studio.document.save przyjmuje documentId, content, title ' +
  'i createVersion. Trwałe zostaje to, co jest składnią treści: styl nazwany bloku, pogrubienie, ' +
  'kursywa, podkreślenie, przekreślenie, listy, tabela i podział strony. Nastawy strony trwają ' +
  'osobno — mają w kontrakcie StudioPageSetup i profil wydania.';

/**
 * Nastawy akapitów trzymane numerem bloku.
 *
 * Numer bloku, a nie odwołanie do elementu: elementy powstają od nowa przy
 * każdym wyrysie, więc odwołanie przeżyłoby jedno przerysowanie. Numer przeżywa
 * wyrys, a przy zmianie liczby bloków nastawa schodzi z bloku, którego już nie
 * ma — i o tym też okno mówi wprost.
 */
export interface NastawyAkapitow {
  /** Nastawy bloku o wskazanym numerze; blok bez nastawy bierze wartości domyślne. */
  dla(numer: number): NastawyAkapitu;
  ustaw(numer: number, nastawy: Partial<NastawyAkapitu>): void;
  /** Liczba bloków z nastawą własną — do zdania paska stanu. */
  ile(): number;
  /** Zdejmuje nastawy bloków, których w treści już nie ma. */
  przytnij(liczbaBlokow: number): void;
}

export function utworzNastawyAkapitow(): NastawyAkapitow {
  const nastawy = new Map<number, NastawyAkapitu>();
  const domyslne = (): NastawyAkapitu => ({ wyrownanie: 'lewo', barwa: '' });

  return {
    dla: (numer) => nastawy.get(numer) ?? domyslne(),

    ustaw(numer, zmiana) {
      const biezace = nastawy.get(numer) ?? domyslne();
      nastawy.set(numer, { ...biezace, ...zmiana });
    },

    ile: () => nastawy.size,

    przytnij(liczbaBlokow) {
      for (const numer of [...nastawy.keys()]) {
        if (numer >= liczbaBlokow) nastawy.delete(numer);
      }
    },
  };
}

/**
 * Style nazwane, którym stopień pisma nadaje sam wyrys, a nie nastawa.
 *
 * Nagłówek ma być większy od tekstu zasadniczego niezależnie od stopnia
 * wybranego dla treści — to jest sedno stylu nazwanego. Krotność stoi tutaj,
 * a nie w arkuszu, bo arkusz nie zna stopnia wybranego przez Operatora.
 */
export function krotnoscStopnia(rodzaj: RodzajBloku): number {
  if (rodzaj === 'naglowek-1') return 1.8;
  if (rodzaj === 'naglowek-2') return 1.45;
  if (rodzaj === 'naglowek-3') return 1.2;
  return 1;
}
