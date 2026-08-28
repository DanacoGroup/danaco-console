import {
  Command,
  type ModuleListRequest,
  type ModuleListResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from './ksztalt-odpowiedzi';
import { wywolaj } from './wywolanie';

/** `module.list` — moduły platformy albo moduły widoczne w jednym środowisku, zależnie od podanego pola. */
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
