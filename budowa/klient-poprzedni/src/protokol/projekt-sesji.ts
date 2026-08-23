import {
  Command,
  type SessionProjectClearRequest,
  type SessionProjectClearResponse,
  type SessionProjectSetRequest,
  type SessionProjectSetResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyTablica, czyTekst, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * Przynależność sesji do projektu — nadanie i wyjęcie.
 *
 * Wyjęcie z projektu nie jest archiwizacją ani usunięciem: sesja zostaje
 * w historii bieżącej i traci wyłącznie przypisanie.
 *
 * Żądanie przypisania niesie albo `projectId`, albo `projectName`; przy pustym
 * identyfikatorze rdzeń zakłada projekt o podanej nazwie i oddaje jego
 * identyfikator, więc klient nie zakłada projektu osobną komendą.
 *
 * `movedIds` może być krótsze od żądania — rdzeń oddaje sesje faktycznie
 * przeniesione, więc widok melduje z wyniku, nie z treści żądania.
 */

/** `session.project.set` — przeniesienie sesji do projektu istniejącego lub nowego. */
export function zadajPrzypisanieProjektu(
  kanal: Kanal,
  zadanie: SessionProjectSetRequest,
): Promise<Wynik<SessionProjectSetResponse>> {
  return wywolaj(kanal, Command.SessionProjectSet, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionProjectSet,
      (tresc) => czyTekst(tresc.projectId) && czyTablica(tresc.movedIds),
    ),
  );
}

/** `session.project.clear` — wyjęcie sesji z projektu; sesja zostaje w historii. */
export function zadajWyjecieZProjektu(
  kanal: Kanal,
  zadanie: SessionProjectClearRequest,
): Promise<Wynik<SessionProjectClearResponse>> {
  return wywolaj(kanal, Command.SessionProjectClear, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionProjectClear, (tresc) => czyTablica(tresc.clearedIds)),
  );
}
