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

/** Kod modułu w rdzeniu — kolumna `modul.kod`, zasilana przez
 * `migracja_007_zaczyn_slownikow.sql`. */
export const KOD_MODULU = 'studio';

/**
 * Panel akcji modułu sterowany danymi.
 *
 * Zestaw operacji Tools Panel pochodzi z katalogu akcji rdzenia: nowa operacja
 * to nowy wiersz rejestru, nie zmiana w kliencie. Dla zasięgu modułu Studio
 * `action.list` oddaje akcje okna komunikacji — `studio.message.send`,
 * `studio.message.stop`, `studio.message.list` — a nie operacje redakcyjne.
 * Wierszy operacji kontekstowych Studia w katalogu nie ma; okno pokazuje to
 * wprost, zamiast dorabiać je po stronie klienta.
 *
 * `window.action` jest drogą generyczną dla akcji bez własnej komendy. Uchwyt
 * komendy istnieje, ale akcja spoza katalogu wraca zwykłą kopertą błędu
 * `not_found` z powodem „akcja … nie istnieje w katalogu akcji". Brakuje więc
 * wiersza katalogu, a nie uchwytu komendy — i tak ma to zobaczyć Operator. Osłona
 * `wywolajUczciwie` zostaje na wypadek koperty `window.unknown`, która nie
 * niesie pola `status` i sama by się nie skorelowała.
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
