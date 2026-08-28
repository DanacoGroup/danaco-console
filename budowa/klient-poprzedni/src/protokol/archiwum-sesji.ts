import {
  Command,
  type SessionArchiveListRequest,
  type SessionArchiveListResponse,
  type SessionArchiveRequest,
  type SessionArchiveResponse,
  type SessionRestoreRequest,
  type SessionRestoreResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

// Archiwum historii sesji — odłożenie, przywrócenie i wgląd w trzy komendy jednego przejścia sesji.

/** `session.archive` — przeniesienie wskazanych sesji z historii bieżącej do archiwum, wraz z wykazem przeniesionych. */
export function zadajArchiwizacjeSesji(
  kanal: Kanal,
  zadanie: SessionArchiveRequest,
): Promise<Wynik<SessionArchiveResponse>> {
  return wywolaj(kanal, Command.SessionArchive, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionArchive, (tresc) => czyTablica(tresc.archivedIds)),
  );
}

/** `session.restore` — przywrócenie wskazanych sesji z archiwum z powrotem do historii bieżącej klienta. */
export function zadajPrzywrocenieSesji(
  kanal: Kanal,
  zadanie: SessionRestoreRequest,
): Promise<Wynik<SessionRestoreResponse>> {
  return wywolaj(kanal, Command.SessionRestore, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionRestore, (tresc) => czyTablica(tresc.restoredIds)),
  );
}

/** `session.archive.list` — stronicowany wykaz sesji archiwum wraz z ich łączną liczbą po stronie rdzenia. */
export function zadajWykazArchiwum(
  kanal: Kanal,
  zadanie: SessionArchiveListRequest,
): Promise<Wynik<SessionArchiveListResponse>> {
  return wywolaj(kanal, Command.SessionArchiveList, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionArchiveList,
      (tresc) => czyTablica(tresc.sessions) && czyLiczba(tresc.total),
    ),
  );
}
