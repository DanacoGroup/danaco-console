import {
  Command,
  type SessionBindRequest,
  type SessionBindResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/** `session.bind` — powiązanie połączenia z sesją, wykonywane po uzgodnieniu, po stronie tego samego klienta. */
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
