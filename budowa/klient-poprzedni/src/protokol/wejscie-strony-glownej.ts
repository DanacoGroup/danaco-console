import {
  Command,
  type HomeEnterRequest,
  type HomeEnterResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/** `home.enter` — wejście na jedyną, platformową stronę główną, pierwsze ogniwo łańcucha nawigacji klienta. */
export function zadajWejscieNaStroneGlowna(
  kanal: Kanal,
  zadanie: HomeEnterRequest,
): Promise<Wynik<HomeEnterResponse>> {
  return wywolaj(kanal, Command.HomeEnter, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.HomeEnter,
      (tresc) => czyTablica(tresc.environments) && czyTablica(tresc.sessions),
    ),
  );
}
