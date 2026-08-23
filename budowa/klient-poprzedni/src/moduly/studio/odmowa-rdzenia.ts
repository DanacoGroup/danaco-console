import {
  ErrorCode,
  zdarzenieNieznanej,
  type Command,
  type ErrorInfo,
  type RequestOf,
  type ResponseOf,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';

/**
 * Wywołanie komendy odporne na odmowę `<obszar>.unknown`.
 *
 * Rdzeń odmawia komendy bez uchwytu kopertą typu `<obszar>.unknown`, która nie
 * niesie pola `status` (`server/internal/protocol/zadanie.go`). Korelacja klienta
 * rozstrzyga wyłącznie koperty ze statusem (`protokol/koperta.ts`, `czyOdpowiedz`),
 * więc zwykłe `wywolaj` na komendzie bez uchwytu zostałoby obietnicą
 * nierozstrzygniętą, a okno stałoby w stanie ładowania bez końca.
 *
 * Wywołanie łączy więc dwie drogi w jeden wynik: odpowiedź skorelowaną oraz
 * zdarzenie odmowy o tym samym `requestId`. Odmowa wraca jako zwykłe
 * niepowodzenie `Wynik.blad`, więc okno obsługuje ją tą samą ścieżką co każdy
 * inny błąd i nie buduje drugiego mechanizmu.
 *
 * Nazwa zdarzenia pochodzi z kontraktu: `zdarzenieNieznanej` tnie typ po
 * separatorze obszaru i sięga do mapy `ZDARZENIA_NIEZNANEJ`, w której stoją
 * obszary `studio` i `window`. Obie drogi odmowy modułu są więc rozpoznawalne
 * bez literału nazwy.
 */
export function wywolajUczciwie<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    let idZadania = '';
    let zamkniete = false;

    function zakoncz(wynik: Wynik<ResponseOf<K>>): void {
      if (zamkniete) return;
      zamkniete = true;
      odsubskrybuj();
      rozstrzygnij(wynik);
    }

    // Subskrypcja przed wysyłką: rdzeń może odmówić natychmiast, a odmowa
    // wyprzedziłaby założenie nasłuchu.
    const odsubskrybuj = kanal.naZdarzenie(zdarzenieNieznanej(komenda), (tresc: unknown) => {
      const odmowa = odczytajOdmowe(tresc);
      if (odmowa === null) return;
      if (odmowa.idZadania !== '' && odmowa.idZadania !== idZadania) return;
      zakoncz({ udany: false, blad: bladOdmowy(komenda, odmowa) });
    });

    idZadania = kanal.wyslij(komenda, zadanie, zakoncz);
  });
}

/** Odmowa rdzenia odczytana z ładunku zdarzenia `<obszar>.unknown`. */
interface OdmowaRdzenia {
  /** Typ, którego rdzeń nie rozpoznał; pusty, gdy rdzeń go nie podał. */
  zadanyTyp: string;
  /** Identyfikator żądania, którego odmowa dotyczy. */
  idZadania: string;
  /** Powód podany przez rdzeń. */
  powod: string;
}

/**
 * Zdanie odmowy nazywające żądany typ.
 *
 * Treść nazywa komendę, której rdzeń nie zna. Kod `not_found` jest tu adekwatny,
 * bo brakuje uchwytu, a nie treści żądania; `retryable: false` powstrzymuje widok
 * przed ponawianiem czegoś, czego rdzeń nie nabędzie przed wdrożeniem nowej wersji.
 */
function bladOdmowy(komenda: string, odmowa: OdmowaRdzenia): ErrorInfo {
  const typ = odmowa.zadanyTyp === '' ? komenda : odmowa.zadanyTyp;
  const powod = odmowa.powod === '' ? '' : ` (${odmowa.powod})`;
  return {
    code: ErrorCode.NotFound,
    message:
      `Rdzeń nie obsługuje jeszcze komendy ${typ}${powod}. ` +
      'Komenda jest w kontrakcie, uchwytu w rdzeniu nie ma — okno nie ma czym wykonać tej czynności.',
    retryable: false,
  };
}

/** Odczyt ładunku odmowy odporny na jego brak i na inny kształt. */
function odczytajOdmowe(ladunek: unknown): OdmowaRdzenia | null {
  if (typeof ladunek !== 'object' || ladunek === null) return null;
  const pola = ladunek as Record<string, unknown>;
  if (typeof pola['requestedType'] !== 'string') return null;
  return {
    zadanyTyp: tekst(pola['requestedType']),
    idZadania: tekst(pola['requestId']),
    powod: tekst(pola['reason']),
  };
}

/** Napis albo pusty łańcuch — pole opcjonalne kontraktu bywa nieobecne. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}
