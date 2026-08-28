import {
  Command,
  type EnvironmentEnterRequest,
  type EnvironmentEnterResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/** `environment.enter` — wejście do środowiska wraz z jego nawigacją, katalogiem modułów i kartami sesji. */
export function zadajWejscieDoSrodowiska(
  kanal: Kanal,
  zadanie: EnvironmentEnterRequest,
): Promise<Wynik<EnvironmentEnterResponse>> {
  return wywolaj(kanal, Command.EnvironmentEnter, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.EnvironmentEnter,
      (tresc) =>
        czyObiekt(tresc.environment) && czyTablica(tresc.modules) && czyTablica(tresc.sessions),
    ),
  );
}
