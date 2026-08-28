import {
  Command,
  type SessionCopyRequest,
  type SessionCopyResponse,
  type SessionRenameRequest,
  type SessionRenameResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLiczba, czyObiekt, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

// Tożsamość sesji w historii — nazwa i kopia, dwie komendy zmieniające to, czym sesja jest w wykazie.

/** `session.rename` — zmiana nazwy sesji widocznej w historii, bez żadnego wpływu na jej zapis i przebieg. */
export function zadajZmianeNazwySesji(
  kanal: Kanal,
  zadanie: SessionRenameRequest,
): Promise<Wynik<SessionRenameResponse>> {
  return wywolaj(kanal, Command.SessionRename, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionRename, (tresc) => czyObiekt(tresc.session)),
  );
}

/** `session.copy` — kopia sesji wraz z jej całym zapisem, oddawana zawsze jako nowa, osobna sesja w historii. */
export function zadajKopieSesji(
  kanal: Kanal,
  zadanie: SessionCopyRequest,
): Promise<Wynik<SessionCopyResponse>> {
  return wywolaj(kanal, Command.SessionCopy, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionCopy,
      (tresc) => czyObiekt(tresc.session) && czyLiczba(tresc.copiedMessages),
    ),
  );
}
