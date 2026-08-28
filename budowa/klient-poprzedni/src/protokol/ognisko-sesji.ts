import {
  Command,
  type SessionFocusRequest,
  type SessionFocusResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLiczba, czyTekst, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/** `session.focus` — przeniesienie ogniska na kartę sesji, opcjonalnie na okno w jej wnętrzu, dla klienta żądania. */
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
