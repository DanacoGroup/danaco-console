import { Command, type ErrorInfo, type Window, type WindowRole } from '../../../shared/contract';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import { zamowienieOkna } from '../okno-komunikacji/zamowienie-okna';
import type { Kanal } from '../protokol/kanal';

/**
 * Wynik zamówienia okna: rdzeń otworzył okno albo podał przyczynę odmowy.
 * Oba pola są rozłączne, więc odbiorca rozstrzyga powodzenie po polu okna,
 * a treść odmowy czyta z pola błędu.
 */
export interface OdpowiedzNaZamowienie {
  /** Okno otwarte przez rdzeń; `null`, gdy komenda się nie powiodła. */
  okno: Window | null;
  /** Przyczyna niepowodzenia, jeżeli rdzeń ją podał. */
  blad?: ErrorInfo;
}

/**
 * Zamówienie kolejnego okna komunikacji w rdzeniu: wysyła komendę
 * `window.create` dla okna innego niż uzgodnione przy nawiązaniu połączenia.
 * Rola nadpisuje rolę z opisu okna, ponieważ rozstrzyga o niej miejsce okna
 * w figurze koordynator–wykonawca.
 */
export function zamowOknoRdzenia(
  kanal: Kanal,
  opis: OpisOkna,
  rola: WindowRole,
  tytul: string,
  przyOdpowiedzi: (odpowiedz: OdpowiedzNaZamowienie) => void,
): void {
  kanal.wyslij(
    Command.WindowCreate,
    {
      ...zamowienieOkna(opis),
      sessionId: kanal.sesja().id(),
      windowRole: rola,
      title: tytul,
    },
    (wynik) => {
      przyOdpowiedzi({ okno: wynik.wynik?.window ?? null, blad: wynik.blad });
    },
  );
}
