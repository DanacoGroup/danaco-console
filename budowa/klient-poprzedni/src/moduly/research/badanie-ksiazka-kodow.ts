import type { ResearchCode } from '../../../../shared/contract';

/**
 * Książka kodów badania składana z tekstu okna — `research.codebook.set`.
 *
 * Komenda zapisuje książkę kodów W CAŁOŚCI: pole `codes` to „kody po zmianie",
 * więc żądanie złożone z samych kodów wpisanych teraz WYMAZAŁOBY wszystkie
 * dotychczasowe. Dlatego zapis idzie po odczycie (`research.codebook.get`)
 * i jest złączeniem, nie podmianą — a zdanie odpowiedzi mówi osobno, ile kodów
 * zostało zachowanych i ile dołożonych, bo to jedyny sposób, żeby Operator
 * rozpoznał wymazanie, gdyby rdzeń zapisał co innego.
 *
 * Kod istniejący rozpoznaje się po nazwie bez względu na wielkość liter, bo
 * nazwa jest tym, co Operator wpisuje; identyfikator kodu nadaje rdzeń
 * i Operator go nie zna. Kod dopisany jedzie z pustym identyfikatorem — tak samo
 * jak nowe pytanie badawcze w `research.workspace.question.set`.
 *
 * Definicja jest w wierszu za znakiem `|`. Znak rozdzielający jest wyborem
 * okna i okno mówi o nim wprost przy chwycie; nazwa kodu bez definicji jest
 * poprawna, bo kontrakt ma pole `description` nieobowiązkowe.
 */

/** Kod wpisany przez Operatora: nazwa i nieobowiązkowa definicja. */
export interface KodZTekstu {
  nazwa: string;
  definicja: string;
}

/** Rozkłada tekst okna na kody; wiersz pusty jest pomijany, nie zgłaszany. */
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

/** Skutek złączenia: co jedzie do rdzenia i co się w nim zmienia. */
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
 * Łączy kody wpisane z książką kodów zastaną w rdzeniu.
 *
 * Kod zastany, którego Operator teraz nie wpisał, ZOSTAJE. Wpisanie trzech
 * kodów nie jest oświadczeniem, że badanie ma trzy kody — jest dopisaniem
 * trzech, a wymazanie reszty byłoby skutkiem, którego nikt nie zamówił. Zdjęcie
 * kodu z książki jest osobną czynnością i okno mówi, że jej dziś nie ma.
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
