import {
  Command,
  type WorkspaceEnterRequest,
  type WorkspaceEnterResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * `workspace.enter` — przeładowanie przestrzeni roboczej karty sesji na moduł.
 *
 * Żądanie idzie z identyfikatorem żywego okna rozmowy: okno wskazane przestawia
 * moduł i zachowuje historię, a dopiero jego brak zakłada okno nowe. Okno
 * zakładane po stronie rdzenia nie dostaje kanału modelu, więc pierwsze
 * `message.send` skończyłoby się odmową.
 *
 * Odpowiedź niesie sesję, moduł, okno i kody okien operacyjnych modułu —
 * sprawdzian kształtu pyta o wszystkie cztery, bo widok przestawia planszę
 * na ich podstawie i pusty komplet dałby planszę bez treści.
 */
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
