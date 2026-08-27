import {
  Command,
  type SessionResumeRequest,
  type SessionResumeResponse,
  type SessionStopRequest,
  type SessionStopResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

// Bieg sesji — wznowienie i zatrzymanie tur biegnących w oknach sesji, dwie osobne komendy protokołu.

/** `session.resume` — wznowienie zamkniętej sesji wraz z oknami otwartymi razem z nią po stronie rdzenia. */
export function zadajWznowienieSesji(
  kanal: Kanal,
  zadanie: SessionResumeRequest,
): Promise<Wynik<SessionResumeResponse>> {
  return wywolaj(kanal, Command.SessionResume, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionResume,
      (tresc) => czyObiekt(tresc.session) && czyTablica(tresc.windows),
    ),
  );
}

/** `session.stop` — zatrzymanie tur biegnących we wszystkich oknach wskazanej sesji, bez jej zamknięcia. */
export function zadajZatrzymanieSesji(
  kanal: Kanal,
  zadanie: SessionStopRequest,
): Promise<Wynik<SessionStopResponse>> {
  return wywolaj(kanal, Command.SessionStop, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionStop,
      (tresc) => czyTekst(tresc.sessionId) && czyTablica(tresc.stoppedWindowIds),
    ),
  );
}
