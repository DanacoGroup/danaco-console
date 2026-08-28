import type { ContextTransferResponse, DesignAsset, DesignBoard } from '../../../../shared/contract';

/**
 * Zdania o skutku trzech czynności modułu Design budowane z odpowiedzi rdzenia,
 * nigdy z zamówienia okna. Każde zdanie zestawia zamówienie z odpowiedzią,
 * a rozbieżność ogłasza odmową.
 */
export interface SkutekDesignu {
  zdanie: string;
  udany: boolean;
}

/**
 * Rozstrzyga, czy dwa zestawy etykiet są równe co do składu, niezależnie od
 * kolejności pozycji; rdzeń wykaz sortuje, więc kolejność nie jest miarą
 * różnicy.
 */
function tenSamZestaw(a: readonly string[], b: readonly string[]): boolean {
  if (a.length !== b.length) return false;
  const posortowane = [...b].sort();
  return [...a].sort().every((wartosc, numer) => wartosc === posortowane[numer]);
}

/**
 * Składa wykaz etykiet do zdania o skutku: pozycje rozdzielone przecinkiem,
 * a zestaw pusty ma własne słowo zamiast pustego nawiasu.
 */
function wypisz(etykiety: readonly string[]): string {
  return etykiety.length === 0 ? 'zestaw pusty' : etykiety.join(', ');
}

/**
 * Skutek `design.asset.tag.set` — z etykiet odesłanych przez rdzeń, wraz
 * z różnicą wobec zestawu zamówionego.
 *
 * Rdzeń oddaje etykiety odczytane z bazy po zapisie, nie echo żądania
 * (`adapter_modul_design_etykiety.go`), więc porównanie jest darmowe.
 */
export function skutekEtykiet(
  zamowione: readonly string[],
  zasob: DesignAsset,
  nazwa: string,
): SkutekDesignu {
  const oddane = zasob.tags ?? [];
  if (!tenSamZestaw(zamowione, oddane)) {
    return {
      zdanie:
        `Rdzeń przyjął żądanie, ale zasób „${nazwa}" ma po nim inny zestaw etykiet niż zamówiony: ` +
        `zamówiono ${wypisz(zamowione)}, rdzeń oddał ${wypisz(oddane)}. Zapisu, o który prosiłeś, nie ma.`,
      udany: false,
    };
  }
  if (oddane.length === 0) {
    return { zdanie: `Zasób „${nazwa}" nie ma już żadnej etykiety.`, udany: true };
  }
  return { zdanie: `Zasób „${nazwa}" ma etykiety: ${oddane.join(', ')}.`, udany: true };
}

/**
 * Skutek `context.transfer` — moduł docelowy wzięty z okna, które wróciło.
 *
 * `transferred` i `window.moduleId` są dwoma różnymi pytaniami: pierwsze mówi,
 * czy przeniesienie się odbyło, drugie — dokąd. Zdanie odpowiada na oba.
 */
export function skutekPrzekazania(
  kodZamowiony: string,
  odpowiedz: ContextTransferResponse,
  nazwaZasobu: string,
): SkutekDesignu {
  const okno = odpowiedz.window;
  if (okno.moduleId !== kodZamowiony) {
    return {
      zdanie:
        `Rdzeń przeniósł zasób „${nazwaZasobu}" do okna ${okno.id} modułu ${okno.moduleId}, ` +
        `a zamówiony był moduł ${kodZamowiony}. Zasób jest gdzie indziej, niż prosiłeś.`,
      udany: false,
    };
  }
  if (!odpowiedz.transferred) {
    return {
      zdanie:
        `Rdzeń oddał okno ${okno.id} modułu ${okno.moduleId}, ale pole transferred jest fałszem — ` +
        `przeniesienia zasobu „${nazwaZasobu}" nie było.`,
      udany: false,
    };
  }
  return {
    zdanie: `Zasób „${nazwaZasobu}" przeniesiony do okna ${okno.id} modułu ${okno.moduleId}.`,
    udany: true,
  };
}

/**
 * Skutek `design.board.update` — nazwa i liczba warstw z kompozycji, która
 * wróciła. Rdzeń składa odpowiedź z wierszy odczytanych po zapisie
 * (`adapter_modul_design_kompozycje.go`), więc wykaz warstw jest stanem bazy,
 * a nie powtórzeniem żądania.
 */
export function skutekZapisuKompozycji(
  nazwaZamowiona: string,
  warstwZamowionych: number,
  board: DesignBoard,
): SkutekDesignu {
  const warstwOddanych = (board.layers ?? []).length;
  const nazwaOddana = board.name ?? '';
  const roznice: string[] = [];
  if (warstwOddanych !== warstwZamowionych) {
    roznice.push(`warstw wysłano ${warstwZamowionych}, rdzeń oddał ${warstwOddanych}`);
  }
  if (nazwaOddana !== nazwaZamowiona) {
    roznice.push(`nazwę wysłano „${nazwaZamowiona}", rdzeń oddał „${nazwaOddana}"`);
  }
  if (roznice.length > 0) {
    return {
      zdanie:
        `Rdzeń zapisał kompozycję ${board.id}, ale nie tak, jak zamówiono: ${roznice.join('; ')}. ` +
        'Kanwa pokazuje układ okna, nie ten, który leży w rdzeniu.',
      udany: false,
    };
  }
  return {
    zdanie:
      `Kompozycja zapisana pod identyfikatorem ${board.id} — rdzeń oddał ${warstwOddanych} warstw` +
      `${nazwaOddana === '' ? ' i kompozycję bez nazwy' : ` pod nazwą „${nazwaOddana}"`}.`,
    udany: true,
  };
}
