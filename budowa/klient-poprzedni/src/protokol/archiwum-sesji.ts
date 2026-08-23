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

/**
 * Archiwum historii sesji — odłożenie, przywrócenie i wgląd.
 *
 * Jedna odpowiedzialność: trzy komendy jednego przejścia — sesja wychodzi
 * z historii bieżącej i wraca do niej. Zapis zostaje w całości, więc
 * archiwizacja nie jest usunięciem i widok nie nazywa jej tak.
 *
 * Wszystkie trzy działają na wielu sesjach: żądania niosą `sessionIds`,
 * a odpowiedzi oddają wykaz faktycznie przeniesionych. Rdzeń może przenieść
 * część zbioru, więc zwrócony wykaz idzie do widoku nietknięty.
 *
 * Wgląd jest stronicowany: `session.archive.list` przyjmuje `offset`/`limit`
 * i oddaje `total` — liczbę wszystkich sesji archiwum, nie długość strony.
 */

/** `session.archive` — przeniesienie sesji do archiwum. */
export function zadajArchiwizacjeSesji(
  kanal: Kanal,
  zadanie: SessionArchiveRequest,
): Promise<Wynik<SessionArchiveResponse>> {
  return wywolaj(kanal, Command.SessionArchive, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionArchive, (tresc) => czyTablica(tresc.archivedIds)),
  );
}

/** `session.restore` — przywrócenie sesji z archiwum do historii bieżącej. */
export function zadajPrzywrocenieSesji(
  kanal: Kanal,
  zadanie: SessionRestoreRequest,
): Promise<Wynik<SessionRestoreResponse>> {
  return wywolaj(kanal, Command.SessionRestore, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionRestore, (tresc) => czyTablica(tresc.restoredIds)),
  );
}

/** `session.archive.list` — wykaz sesji archiwum wraz z ich łączną liczbą. */
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
