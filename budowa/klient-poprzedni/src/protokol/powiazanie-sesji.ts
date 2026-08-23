import {
  Command,
  type SessionBindRequest,
  type SessionBindResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * `session.bind` — powiązanie połączenia z sesją, wykonywane po uzgodnieniu.
 *
 * Klient albo wiąże się z sesją trwającą na rdzeniu, albo wchodzi na stronę
 * główną. Rozłączenie klienta nie kończy sesji ani procesów, więc powiązanie
 * zwykle zastaje sesję żywą: pole `resumed` odróżnia odtworzenie stanu od
 * założenia go od nowa, a `windows` niesie okna, których strumienie od tej
 * chwili idą na bieżące połączenie.
 *
 * Puste `windowIds` w żądaniu wiąże wszystkie okna sesji — zawężenie jest
 * decyzją widoku, nie protokołu.
 */
export function zadajPowiazanieSesji(
  kanal: Kanal,
  zadanie: SessionBindRequest,
): Promise<Wynik<SessionBindResponse>> {
  return wywolaj(kanal, Command.SessionBind, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.SessionBind,
      (tresc) =>
        czyObiekt(tresc.session) &&
        czyTablica(tresc.windows) &&
        czyLogiczna(tresc.bound) &&
        czyLogiczna(tresc.resumed),
    ),
  );
}
