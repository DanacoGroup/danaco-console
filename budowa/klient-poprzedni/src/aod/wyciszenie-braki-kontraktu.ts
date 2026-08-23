import { Command, KOMENDY } from '../../../shared/contract';

/**
 * Czego kontrakt nie niesie przy wyciszeniu nakładki — nazwane wprost, w oknie,
 * a nie zasłonięte uprzejmym milczeniem.
 *
 * Stan zmierzony po dobudowie wyciszenia w rdzeniu (17.08.2026): kontrakt niesie
 * już wyciszenie nakładki jako BYT RDZENIA — `AodMute` wraz z rodzajem, zakresem
 * i chwilą końca, komendy `aod.mute.get` i `aod.mute.set`, zdarzenie
 * `aod.mute.changed` rozgłaszane na pozostałe powłoki, nośnik sygnału klas
 * zdarzeń (`aod.signal.report`, `aod.signal.list`), moduł przy `AodStatus`
 * i `AodSuggestion` oraz kategorię ustawień „Always On Display" z pozycją
 * „reguły wyciszania".
 *
 * Zostaje JEDEN brak i jest po stronie nakładki, nie kontraktu: wołacze tego
 * okna nadal piszą do magazynu stanowiska, więc wyciszenie założone tutaj nie
 * dojdzie do drugiej powłoki, dopóki magazyn nie zostanie przełożony na rdzeń.
 * Magazyn jest podawany (`MagazynWyciszen`) właśnie po to, żeby to przełożenie
 * nie ruszyło ani jednego wołacza.
 *
 * Wzorem `moduly/apps/braki-kontraktu.ts` i `aod/odmowy-aod.ts` zdanie mówi,
 * CZEGO brakuje i PO CZYJEJ stronie. Nazwy komend biorą się z `Command`, nie
 * z literału, więc zmiana nazwy w kontrakcie przechodzi przez kompilację.
 *
 * Ten plik NICZEGO nie obchodzi: nie wysyła komendy zastępczej i nie udaje
 * zapisu w rdzeniu. Wyciszenie jest dziś stanem okna i tak jest nazwane.
 */

/** Jedna pozycja braku wyciszenia. */
export interface BrakWyciszenia {
  /** Czynność tak, jak widzi ją Operator w menu wyciszania. */
  czynnosc: string;
  /** Czego by trzeba — zdanie o bytach kontraktu, nie o kodzie nakładki. */
  czego: string;
  /**
   * Nazwa komendy, której wejście do kontraktu znosi ten brak. Pozycja z komendą
   * już obecną w `KOMENDY` nie trafia do wykazu wcale. Pominięta znaczy brak
   * leżący po stronie nakładki — takiego nie zniesie żadna komenda, tylko robota
   * w tym drzewie.
   */
  komendaZnoszaca?: string;
}

/**
 * Wykaz braków stojących przy wyciszeniu nakładki.
 *
 * Po dobudowie wyciszenia w rdzeniu została jedna pozycja i jest to brak po
 * stronie nakładki: wołacze piszą do magazynu stanowiska, a nie do rdzenia.
 */
export const BRAKI_WYCISZENIA: readonly BrakWyciszenia[] = [
  {
    czynnosc: 'Wyciszenie wspólne dla wszystkich powłok Operatora',
    czego:
      'Przełożenia magazynu wyciszeń tego okna na rdzeń. Kontrakt niesie już wyciszenie jako byt ' +
      `rdzenia — \`${Command.AodMuteGet}\` czyta wykaz, \`${Command.AodMuteSet}\` zakłada je ` +
      'i znosi jednym ruchem, a zmiana rozgłasza się na pozostałe powłoki. Zapis tego okna idzie ' +
      'jednak nadal do magazynu stanowiska, więc wyciszenie założone tutaj nie dojdzie do drugiej ' +
      'powłoki. Brak jest po stronie nakładki.',
  },
];

/** Komendy rodziny `aod.*`, które kontrakt niesie dziś — odczytane w czasie działania. */
export function komendyNakladki(): readonly string[] {
  return [...(KOMENDY as readonly string[])].filter((komenda) => komenda.startsWith('aod.')).sort();
}

/** Braki nadal stojące — pozycja z komendą już wniesioną do kontraktu wypada. */
export function brakiCzynne(): readonly BrakWyciszenia[] {
  return BRAKI_WYCISZENIA.filter(
    (brak) =>
      brak.komendaZnoszaca === undefined ||
      !(KOMENDY as readonly string[]).includes(brak.komendaZnoszaca),
  );
}

/** Pełne zdanie jednego braku: czynność, czego brakuje, po czyjej stronie. */
export function zdanieBraku(brak: BrakWyciszenia): string {
  const strona =
    brak.komendaZnoszaca === undefined
      ? 'Brak jest po stronie nakładki; kontrakt niesie, czego trzeba.'
      : `Brak jest po stronie kontraktu — komendy ${brak.komendaZnoszaca} nie ma w wykazie ` +
        `${(KOMENDY as readonly string[]).length} komend.`;
  return `${brak.czynnosc}: ${brak.czego} ${strona}`;
}

/**
 * Zdanie zamykające menu wyciszania — jedno, krótkie, prawdziwe.
 *
 * Mówi, dokąd wyciszenie sięga (to okno) i dokąd nie sięga (pozostałe powłoki),
 * oraz po czyjej stronie jest brak. Zdanie zmienia się samo w chwili, w której
 * wołacze przestaną pisać do magazynu stanowiska: wykaz braków jest wtedy pusty.
 */
export function zdanieGranicyWyciszenia(): string {
  const braki = brakiCzynne();
  if (braki.length === 0) {
    return 'Wyciszenie zapisuje się w rdzeniu i obowiązuje na wszystkich urządzeniach Operatora.';
  }
  return (
    'Wyciszenie jest dziś stanem TEGO okna, choć kontrakt niesie już wyciszenie nakładki jako byt ' +
    `rdzenia — rodzina aod.* ma ${komendyNakladki().length} komend, wśród nich odczyt i zapis ` +
    'wyciszenia wraz z rozgłoszeniem na pozostałe powłoki. Zapis tego okna idzie jeszcze do ' +
    'magazynu stanowiska, więc wyciszenie nie przechodzi na pozostałe urządzenia Operatora. Brak ' +
    'jest po stronie nakładki, nie kontraktu.'
  );
}
