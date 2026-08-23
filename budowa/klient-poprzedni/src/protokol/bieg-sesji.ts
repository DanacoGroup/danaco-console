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

/**
 * Bieg sesji — wznowienie i zatrzymanie tur.
 *
 * Jedna odpowiedzialność: dwie komendy działające na tym, co się w sesji dzieje,
 * a nie na jej miejscu w historii.
 *
 * Nie są swoimi odwrotnościami. `session.stop` przerywa tury biegnące w oknach
 * sesji, a sesja zostaje otwarta i gotowa na kolejną turę; `session.resume`
 * otwiera sesję zamkniętą wraz z jej oknami. Widok, który postawiłby je jako
 * parę „start/stop", obiecałby coś, czego rdzeń nie robi, więc każda z nich
 * pojawia się osobno i pod własnym warunkiem.
 *
 * Obie oddają zakres skutku. `stoppedWindowIds` mówi, w ilu oknach turę
 * naprawdę zatrzymano — zero znaczy, że nic nie biegło, i to jest wynik
 * poprawny, nie odmowa. `windows` po wznowieniu niesie okna otwarte razem
 * z sesją, więc powrót do niej nie musi ich odpytywać drugi raz.
 */

/** `session.resume` — wznowienie zamkniętej sesji wraz z jej oknami. */
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

/** `session.stop` — zatrzymanie tur biegnących we wszystkich oknach sesji. */
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
