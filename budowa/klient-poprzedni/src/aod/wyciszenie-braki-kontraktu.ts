/**
 * Wykaz braków stojących przy wyciszeniu nakładki wraz ze zdaniami, którymi okno
 * nazywa je w menu wyciszania. Nazwy komend pochodzą z wyliczenia `Command`, więc
 * zmiana nazwy w kontrakcie przechodzi przez kompilację.
 */
import { Command, KOMENDY } from '../../../shared/contract';

/**
 * Jedna pozycja braku wyciszenia: czynność widziana w menu wyciszania, treść
 * brakująca oraz komenda, której wejście do kontraktu ten brak znosi.
 */
export interface BrakWyciszenia {
  /** Czynność tak, jak widzi ją Operator w menu wyciszania. */
  czynnosc: string;
  /** Czego by trzeba — zdanie o bytach kontraktu, nie o kodzie nakładki. */
  czego: string;
  /** Nazwa komendy, której wejście do kontraktu znosi ten brak. */
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

/**
 * Oddaje komendy rodziny `aod.*`, które kontrakt niesie w chwili wywołania.
 * Wykaz powstaje z odczytu stałej `KOMENDY`, więc nadąża za kontraktem bez
 * przepisywania nazw komend do kodu nakładki.
 */
export function komendyNakladki(): readonly string[] {
  return [...(KOMENDY as readonly string[])].filter((komenda) => komenda.startsWith('aod.')).sort();
}

/**
 * Oddaje braki nadal stojące. Pozycja, której komenda znosząca weszła już do
 * wykazu komend kontraktu, wypada z wyniku, a pozycja bez komendy znoszącej
 * pozostaje w nim zawsze.
 */
export function brakiCzynne(): readonly BrakWyciszenia[] {
  return BRAKI_WYCISZENIA.filter(
    (brak) =>
      brak.komendaZnoszaca === undefined ||
      !(KOMENDY as readonly string[]).includes(brak.komendaZnoszaca),
  );
}

/**
 * Składa pełne zdanie jednego braku: czynność, treść brakującą oraz stronę, po
 * której brak leży. Stronę rozstrzyga obecność komendy znoszącej w wykazie
 * komend kontraktu.
 */
export function zdanieBraku(brak: BrakWyciszenia): string {
  const strona =
    brak.komendaZnoszaca === undefined
      ? 'Brak jest po stronie nakładki; kontrakt niesie, czego trzeba.'
      : `Brak jest po stronie kontraktu — komendy ${brak.komendaZnoszaca} nie ma w wykazie ` +
        `${(KOMENDY as readonly string[]).length} komend.`;
  return `${brak.czynnosc}: ${brak.czego} ${strona}`;
}

/**
 * Składa zdanie zamykające menu wyciszania: podaje, dokąd wyciszenie sięga
 * i dokąd nie sięga, oraz po czyjej stronie leży brak. Treść zdania wynika
 * z wykazu braków czynnych, a wykaz pusty daje zdanie o zapisie w rdzeniu.
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
