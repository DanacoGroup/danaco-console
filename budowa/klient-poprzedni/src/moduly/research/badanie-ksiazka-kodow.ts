import type { ResearchCode } from '../../../../shared/contract';

/** Książka kodów badania składana z tekstu okna komendą `research.codebook.set`. */

/**
 * Kod wpisany przez Operatora: nazwa oraz nieobowiązkowa definicja. Definicja
 * jest pusta, gdy wiersz nie niesie znaku rozdzielającego albo nie ma za nim
 * treści.
 */
export interface KodZTekstu {
  nazwa: string;
  definicja: string;
}

/**
 * Rozkłada tekst okna na kody, wiersz po wierszu. Wiersz pusty oraz wiersz bez
 * nazwy jest pomijany bez zgłoszenia, ponieważ pisanie w polu tekstowym
 * przechodzi przez stany niepełne.
 */
export function kodyZTekstu(tekst: string): readonly KodZTekstu[] {
  const kody: KodZTekstu[] = [];
  for (const linia of tekst.split('\n')) {
    const przycieta = linia.trim();
    if (przycieta === '') continue;
    const granica = przycieta.indexOf('|');
    const nazwa = (granica < 0 ? przycieta : przycieta.slice(0, granica)).trim();
    const definicja = granica < 0 ? '' : przycieta.slice(granica + 1).trim();
    if (nazwa === '') continue;
    kody.push({ nazwa, definicja });
  }
  return kody;
}

/**
 * Skutek złączenia: pełna książka kodów jadąca do rdzenia wraz z licznikami
 * kodów dopisanych, poprawionych oraz zachowanych bez zmiany.
 */
export interface ZlaczenieKodow {
  /** Pełna książka kodów po zmianie — treść pola `codes`. */
  kody: readonly ResearchCode[];
  /** Kody dopisane. */
  dopisane: number;
  /** Kody zastane, którym zmieniono definicję. */
  poprawione: number;
  /** Kody zastane pozostawione bez zmiany. */
  zachowane: number;
}

/**
 * Łączy kody wpisane z książką kodów zastaną w rdzeniu. Kod zastany, którego
 * Operator teraz nie wpisał, zostaje: wpisanie trzech kodów jest dopisaniem
 * trzech, a nie oświadczeniem, że badanie ma dokładnie tyle kodów.
 */
export function zlaczKody(
  zastane: readonly ResearchCode[],
  wpisane: readonly KodZTekstu[],
): ZlaczenieKodow {
  const kody: ResearchCode[] = zastane.map((kod) => ({ ...kod }));
  let dopisane = 0;
  let poprawione = 0;

  for (const kod of wpisane) {
    const klucz = kod.nazwa.toLocaleLowerCase('pl-PL');
    const zastany = kody.find((wpis) => wpis.name.trim().toLocaleLowerCase('pl-PL') === klucz);
    if (zastany === undefined) {
      kody.push({
        id: '',
        name: kod.nazwa,
        ...(kod.definicja === '' ? {} : { description: kod.definicja }),
      });
      dopisane += 1;
      continue;
    }
    if (kod.definicja !== '' && (zastany.description ?? '') !== kod.definicja) {
      zastany.description = kod.definicja;
      poprawione += 1;
    }
  }

  return { kody, dopisane, poprawione, zachowane: zastane.length - poprawione };
}
