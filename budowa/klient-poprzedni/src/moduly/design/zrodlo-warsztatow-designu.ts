import {
  Command,
  type DesignAsset,
  type DesignAssetListResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Droga okien warsztatowych modułu Design do rdzenia: wykaz materiału i
 * wykonanie czynności.
 *
 * `wykonaj` jest jedną drogą na siedemdziesiąt pięć komend, bo wszystkie idą tak
 * samo: nazwa komendy ze stałych kontraktu, treść żądania złożona z formularza.
 * Nazwa nie jest tu nigdy napisem wpisanym z pamięci — przychodzi z katalogu
 * czynności (`czynnosci-warsztatow-designu.ts`), a ten bierze ją ze stałych.
 *
 * Kształt odpowiedzi sprawdza RDZEŃ i odsyła odmowę walidacji z nazwą pola.
 * Klient tego nie zastąpi: kontrakt rozstrzyga po stronie rdzenia, a drugie
 * sprawdzenie tutaj byłoby drugą prawdą o tym, co wolno wysłać — i rozjechałoby
 * się z pierwszą przy pierwszej zmianie kontraktu.
 *
 * Materiał bierze się z jednego magazynu zasobów (`design.asset.list`) — tego
 * samego, z którego czyta Assets Panel. Drugi wykaz materiału byłby drugim
 * miejscem, w którym ta sama treść żyje.
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
