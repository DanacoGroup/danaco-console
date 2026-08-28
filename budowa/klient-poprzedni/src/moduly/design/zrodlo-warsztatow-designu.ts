import {
  Command,
  type DesignAsset,
  type DesignAssetListResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Droga okien warsztatowych modułu Design do rdzenia: wykaz materiału komendą
 * `design.asset.list` oraz wykonanie czynności komendą podaną przez katalog
 * czynności warsztatów. Kształt odpowiedzi sprawdza rdzeń.
 */
export interface ZrodloWarsztatowDesignu {
  /** Zasoby leżące w magazynie okna — materiał czynności warsztatów. */
  zasoby(idOkna: string): Promise<Wynik<DesignAsset[]>>;
  /** Wykonuje czynność warsztatu wskazaną komendą kontraktu. */
  wykonaj(komenda: Command, zadanie: Record<string, unknown>): Promise<Wynik<unknown>>;
}

export function utworzZrodloWarsztatowDesignu(kanal: Kanal): ZrodloWarsztatowDesignu {
  return {
    async zasoby(idOkna) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetList, { windowId: idOkna }),
        Command.DesignAssetList,
        (tresc: DesignAssetListResponse) => czyTablica(tresc.assets),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
      }
      return { udany: true, wynik: wynik.wynik.assets };
    },

    wykonaj(komenda, zadanie) {
      return wywolaj(kanal, komenda, zadanie as never) as Promise<Wynik<unknown>>;
    },
  };
}
