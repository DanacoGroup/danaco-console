import type { RodzajBloku } from './zapis-formatowany';

/** Wyrównanie akapitu dokumentu, jedna z nastaw wizualnych treści obok kroju, stopnia, interlinii i barwy: do lewej, do środka, do prawej albo obustronnie. */
export type Wyrownanie = 'lewo' | 'srodek' | 'prawo' | 'obustronne';

/** Kroje pisma oferowane w oknie, nazwy rodzin obecnych w motywie interfejsu oraz w systemie operacyjnym. */
export const KROJE: readonly { wartosc: string; etykieta: string }[] = [
  { wartosc: 'var(--dn-ff-tekst)', etykieta: 'Krój tekstowy motywu' },
  { wartosc: 'var(--dn-ff-naglowek)', etykieta: 'Krój nagłówkowy motywu' },
  { wartosc: 'var(--dn-ff-mono)', etykieta: 'Krój o stałej szerokości' },
  { wartosc: 'Georgia, serif', etykieta: 'Szeryfowy (Georgia)' },
  { wartosc: 'Arial, Helvetica, sans-serif', etykieta: 'Bezszeryfowy (Arial)' },
  { wartosc: '"Times New Roman", serif', etykieta: 'Times New Roman' },
];

/** Stopnie pisma dostępne operatorowi do wyboru w punktach typograficznych, od najmniejszego do największego. */
export const STOPNIE: readonly number[] = [8, 9, 10, 11, 12, 14, 16, 18, 24, 36];

/** Nastawy wizualne całego dokumentu: krój, stopień pisma, interlinia, wcięcie pierwszego wiersza i odstęp akapitowy. */
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

/** Nastawy akapitu, które operator daje pojedynczemu blokowi treści: wyrównanie tekstu oraz barwa pisma. */
export interface NastawyAkapitu {
  wyrownanie: Wyrownanie;
  /** Barwa pisma; puste znaczy „barwa tekstu motywu". */
  barwa: string;
}

/** Nastawy domyślne dokumentu: krój tekstowy motywu, jedenaście punktów wielkości i interlinia półtorej. */
export function domyslneNastawy(): NastawyWizualne {
  return {
    krój: 'var(--dn-ff-tekst)',
    stopien: 11,
    interlinia: 1.5,
    wciecieMm: 0,
    odstepMm: 3,
  };
}

/** Zdanie opisujące bieżące nastawy wizualne całego dokumentu, wyświetlane operatorowi w pasku statusu. */
export function opiszNastawy(nastawy: NastawyWizualne): string {
  const krój = KROJE.find((pozycja) => pozycja.wartosc === nastawy.krój)?.etykieta ?? nastawy.krój;
  return (
    `${krój} · ${nastawy.stopien} pt · interlinia ${nastawy.interlinia} · ` +
    `wcięcie ${nastawy.wciecieMm} mm · odstęp akapitowy ${nastawy.odstepMm} mm`
  );
}

/** Powód, dla którego nastawy wizualne dokumentu nie dojeżdżają do rdzenia i żyją wyłącznie przez sesję okna. */
export const POWOD_NIETRWALOSCI =
  'Krój, stopień, interlinia, wcięcie, odstęp akapitowy, wyrównanie i barwa są nastawami WIDOKU ' +
  'i żyją przez tę sesję okna. Kontrakt nie ma pola stylu w dokumencie ani komendy, którą styl ' +
  'dojechałby do rdzenia — studio.document.save przyjmuje documentId, content, title ' +
  'i createVersion. Trwałe zostaje to, co jest składnią treści: styl nazwany bloku, pogrubienie, ' +
  'kursywa, podkreślenie, przekreślenie, listy, tabela i podział strony. Nastawy strony trwają ' +
  'osobno — mają w kontrakcie StudioPageSetup i profil wydania.';

/** Nastawy akapitów dokumentu trzymane numerem bloku, a nie odwołaniem do jego elementu na stronie dokumentu. */
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

/** Style nazwane, którym stopień pisma nadaje sam wyrys dokumentu, a nie nastawa wybrana przez operatora. */
export function krotnoscStopnia(rodzaj: RodzajBloku): number {
  if (rodzaj === 'naglowek-1') return 1.8;
  if (rodzaj === 'naglowek-2') return 1.45;
  if (rodzaj === 'naglowek-3') return 1.2;
  return 1;
}
