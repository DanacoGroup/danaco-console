import {
  Command,
  type HomeEnterRequest,
  type HomeEnterResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * `home.enter` — wejście na jedyną, platformową stronę główną.
 *
 * Pierwsze ogniwo łańcucha nawigacji: `home.enter` → `environment.list` →
 * `environment.enter` → `module.list` → `workspace.enter`. Wejście mówi
 * rdzeniowi, który klient stanął na stronie głównej; rdzeń wiąże ognisko
 * z klientem, więc bez tego wywołania nie ma komu oddać `focusedSessionId`.
 *
 * Ponad `environment.list` odpowiedź niesie środowiska wraz z kodami modułów,
 * sesje czynne konta, sesję ostatnio ogniskowaną oraz stan sesji trwających
 * w tle. `environment.list` bez `includeModules` nie daje żadnej z tych rzeczy.
 */
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
