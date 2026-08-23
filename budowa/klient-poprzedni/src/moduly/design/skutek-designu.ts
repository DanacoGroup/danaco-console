import type { ContextTransferResponse, DesignAsset, DesignBoard } from '../../../../shared/contract';

/**
 * Zdania o skutku trzech czynności modułu — budowane z odpowiedzi rdzenia,
 * nigdy z zamówienia okna. Wydzielone z okien, bo okna składają kontrolki,
 * a to jest ocena odpowiedzi (wzór: `assistant/skutek-sterowania.ts`).
 *
 * Każde zdanie porównuje zamówienie z odpowiedzią, a rozbieżność ogłasza
 * odmową. Zdanie zbudowane z tego, co okno wysłało, bywa prawdziwe przypadkiem
 * i skłamie, gdy rdzeń zapisze co innego.
 *
 * Trzy miejsca, w których to ma znaczenie:
 *
 *   (1) Przekazanie. `context.transfer` nie sprawdza katalogu modułów —
 *       przepisuje `targetModuleId` do okna docelowego jak leci i oddaje
 *       `transferred: true` także dla modułu, którego nie ma. Jedynym polem
 *       mówiącym, gdzie zasób wylądował, jest `window.moduleId` z odpowiedzi.
 *
 *   (2) Zapis kompozycji. Odpowiedź niesie całą kompozycję odczytaną po
 *       zapisie — z nazwą i wykazem warstw. Liczba warstw, która wróciła, jest
 *       jedyną miarą tego, ile ich leży w rdzeniu.
 *
 *   (3) Nadanie etykiet. Zestaw zastępuje poprzedni, więc różnica wobec
 *       zamówionego zestawu znaczy, że zapis nie jest tym, o który Operator
 *       prosił — i musi być odmową, nie milczeniem.
 */
export interface SkutekDesignu {
  zdanie: string;
  udany: boolean;
}

/** Zestawy równe co do składu, niezależnie od kolejności — rdzeń sortuje wykaz. */
function tenSamZestaw(a: readonly string[], b: readonly string[]): boolean {
  if (a.length !== b.length) return false;
  const posortowane = [...b].sort();
  return [...a].sort().every((wartosc, numer) => wartosc === posortowane[numer]);
}

/** Wykaz etykiet do zdania; pusty zestaw ma własne słowo, nie pusty nawias. */
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
 * wróciła.
 *
 * Rdzeń składa odpowiedź z wierszy odczytanych po zapisie
 * (`adapter_modul_design_kompozycje.go`, `zlozBoard`), więc wykaz warstw
 * odpowiedzi jest stanem bazy, a nie powtórzeniem żądania.
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
