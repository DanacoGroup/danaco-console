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

/**
 * Tożsamość sesji w historii — nazwa i kopia.
 *
 * Jedna odpowiedzialność: dwie komendy, które zmieniają to, czym sesja jest
 * w wykazie, a nie to, co się w niej dzieje. Bieg sesji niesie `bieg-sesji`,
 * przynależność `projekt-sesji`, a odłożenie `archiwum-sesji`.
 *
 * Kopia jest osobnym bytem: `session.copy` oddaje nową sesję wraz z liczbą
 * przeniesionych wiadomości — kopiowanie nie jest odnośnikiem do źródła,
 * więc widok mówi wprost, ile zapisu naprawdę poszło.
 *
 * Pusta nazwa kopii nie jest błędem. Kontrakt pozwala pominąć `title`
 * i wtedy rdzeń bierze nazwę źródła z dopiskiem. Klient nie składa tego
 * dopisku sam — nazwę nadaje ta strona, która ją zna.
 */

/** `session.rename` — zmiana nazwy sesji w historii. */
export function zadajZmianeNazwySesji(
  kanal: Kanal,
  zadanie: SessionRenameRequest,
): Promise<Wynik<SessionRenameResponse>> {
  return wywolaj(kanal, Command.SessionRename, zadanie).then((wynik) =>
    sprawdzKsztalt(wynik, Command.SessionRename, (tresc) => czyObiekt(tresc.session)),
  );
}

/** `session.copy` — kopia sesji wraz z jej zapisem. */
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
