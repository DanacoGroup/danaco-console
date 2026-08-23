import {
  Command,
  type SessionFocusRequest,
  type SessionFocusResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLiczba, czyTekst, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * `session.focus` — przeniesienie ogniska na kartę sesji, opcjonalnie na okno
 * w jej wnętrzu.
 *
 * Ognisko jest właściwością klienta, nie konta: to samo konto otwarte na dwóch
 * urządzeniach ma dwa ogniska, dlatego treść żądania niesie `clientId`.
 * Rozgłoszenie zmiany należy do rdzenia — po wykonaniu komendy rozsyła on
 * zdarzenie `session.focus.changed`, które przechwytuje obserwator ogniska
 * w warstwie łączności.
 */
export function zadajPrzeniesienieOgniska(
  kanal: Kanal,
  zadanie: SessionFocusRequest,
): Promise<Wynik<SessionFocusResponse>> {
  return wywolaj(kanal, Command.SessionFocus, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionFocus,
      (tresc) => czyTekst(tresc.sessionId) && czyLiczba(tresc.focusedAt),
    ),
  );
}
