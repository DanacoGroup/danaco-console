import {
  Command,
  type WorkspaceEnterRequest,
  type WorkspaceEnterResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/** `workspace.enter` — przeładowanie przestrzeni roboczej karty sesji na wskazany moduł tego środowiska. */
export function zadajWejscieDoPrzestrzeni(
  kanal: Kanal,
  zadanie: WorkspaceEnterRequest,
): Promise<Wynik<WorkspaceEnterResponse>> {
  return wywolaj(kanal, Command.WorkspaceEnter, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.WorkspaceEnter,
      (tresc) =>
        czyObiekt(tresc.session) &&
        czyObiekt(tresc.module) &&
        czyObiekt(tresc.window) &&
        czyTablica(tresc.operationalWindowCodes),
    ),
  );
}
