import {
  Command,
  ConfigScope,
  type Action,
  type WindowActionResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { wywolajUczciwie } from './odmowa-rdzenia';

/** Kod modułu w rdzeniu — kolumna `modul.kod`, zasilana przez migrację `migracja_007_zaczyn_slownikow.sql`. */
export const KOD_MODULU = 'studio';

/**
 * Panel akcji modułu sterowany danymi: zestaw operacji Tools Panel pochodzi z katalogu akcji
 * rdzenia, nie ze zmian w kliencie.
 */
export interface ZrodloAkcjiStudio {
  /** Katalog akcji zasięgu modułu Studio. */
  katalog(): Promise<Wynik<{ actions: Action[] }>>;
  /** Wykonuje akcję panelu akcji okna operacyjnego. */
  wykonaj(
    idOkna: string,
    idAkcji: string,
    parametry?: unknown,
  ): Promise<Wynik<WindowActionResponse>>;
}

export function utworzZrodloAkcjiStudio(kanal: Kanal): ZrodloAkcjiStudio {
  return {
    async katalog() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ActionList, {
          scope: ConfigScope.Module,
          scopeId: KOD_MODULU,
          enabledOnly: true,
        }),
        Command.ActionList,
        (tresc) => czyTablica(tresc.actions),
      );
    },

    async wykonaj(idOkna, idAkcji, parametry) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.WindowAction, {
          windowId: idOkna,
          actionId: idAkcji,
          ...(parametry === undefined ? {} : { parameters: parametry }),
        }),
        Command.WindowAction,
        (tresc) => typeof tresc.actionId === 'string',
      );
    },
  };
}
