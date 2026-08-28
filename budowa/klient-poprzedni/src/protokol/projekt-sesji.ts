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

// Przynależność sesji do projektu — nadanie i wyjęcie, bez wpływu na archiwizację ani usunięcie sesji.

/** `session.project.set` — przeniesienie wskazanych sesji do projektu istniejącego albo nowo zakładanego. */
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

/** `session.project.clear` — wyjęcie wskazanych sesji z projektu; sesja zostaje nadal w historii bieżącej. */
export function zadajWyjecieZProjektu(
  kanal: Kanal,
  zadanie: SessionProjectClearRequest,
): Promise<Wynik<SessionProjectClearResponse>> {
  return wywolaj(kanal, Command.SessionProjectClear, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionProjectClear, (tresc) => czyTablica(tresc.clearedIds)),
  );
}
