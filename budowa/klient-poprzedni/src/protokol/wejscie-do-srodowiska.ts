import {
  Command,
  type EnvironmentEnterRequest,
  type EnvironmentEnterResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * `environment.enter` — wejście do środowiska wraz z jego nawigacją i kartami
 * sesji.
 *
 * Jedno wywołanie oddaje komplet bocznej nawigacji: środowisko, jego moduły
 * z katalogiem okien operacyjnych oraz karty sesji otwarte w tym środowisku.
 * Dlatego wykaz nawigacji nie składa się z dwóch odczytów — `module.list` jest
 * tu wyłącznie drugim podejściem, gdy wejście oddało wykaz pusty mimo nawigacji
 * modułowej (fail-open zamiast pustej kolumny).
 *
 * Warstwa nie buduje ani jednego elementu widoku: oddaje treść odpowiedzi
 * w kształcie z `shared/contract.ts` i zostawia widokowi rozstrzygnięcie, co
 * z nią zrobić. Nazwa komendy pochodzi wyłącznie ze stałej `Command.*`.
 */
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
