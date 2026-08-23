import {
  Command,
  type ModuleListRequest,
  type ModuleListResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/**
 * `module.list` — moduły platformy albo moduły widoczne w jednym środowisku.
 *
 * Komenda pełni dwie role:
 *  • z `environmentId` — drugie podejście bocznej nawigacji, gdy
 *    `environment.enter` oddał wykaz pusty mimo nawigacji modułowej;
 *  • bez pola — komplet modułów platformy, niezależny od macierzy
 *    `srodowisko_modul`. Tą drogą otwiera się moduł, którego macierz nie
 *    pokazuje w żadnym środowisku.
 *
 * Odpowiedź niesie `total` obok wykazu, więc sprawdzian kształtu pyta o oba:
 * wykaz bez licznika znaczy, że rdzeń odpowiedział czymś innym niż `module.list`.
 */
export function zadajWykazModulow(
  kanal: Kanal,
  zadanie: ModuleListRequest = {},
): Promise<Wynik<ModuleListResponse>> {
  return wywolaj(kanal, Command.ModuleList, zadanie).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.ModuleList,
      (tresc) => czyTablica(tresc.modules) && czyLiczba(tresc.total),
    ),
  );
}
