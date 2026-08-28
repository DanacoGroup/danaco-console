import {
  Command,
  DesignAssetKind,
  type DesignAsset,
  type DesignAssetListResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Droga okna warsztatu do rdzenia: wykaz materiału i wykonanie czynności; materiał bierze się
 * z magazynu zasobów, nie z obszaru `studio`.
 */
export interface ZrodloWarsztatuDokumentu {
  /** Dokumenty leżące w magazynie okna — materiał czynności warsztatu. */
  dokumenty(idOkna: string): Promise<Wynik<DesignAsset[]>>;
  /** Zasoby okna bez zawężenia rodzaju — pieczęcie graficzne i certyfikaty. */
  zasoby(idOkna: string): Promise<Wynik<DesignAsset[]>>;
  /** Wykonuje czynność warsztatu wskazaną komendą kontraktu. */
  wykonaj(komenda: Command, zadanie: Record<string, unknown>): Promise<Wynik<unknown>>;
}

/** Odczyt wykazu zasobów okna warsztatu z zawężeniem rodzaju zasobu albo bez takiego zawężenia rodzaju. */
async function wykazZasobow(
  kanal: Kanal,
  idOkna: string,
  rodzaj?: DesignAssetKind,
): Promise<Wynik<DesignAsset[]>> {
  const wynik = sprawdzKsztalt(
    await wywolaj(kanal, Command.DesignAssetList, {
      windowId: idOkna,
      ...(rodzaj === undefined ? {} : { kind: rodzaj }),
    }),
    Command.DesignAssetList,
    (tresc: DesignAssetListResponse) => czyTablica(tresc.assets),
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
  }
  return { udany: true, wynik: wynik.wynik.assets };
}

export function utworzZrodloWarsztatuDokumentu(kanal: Kanal): ZrodloWarsztatuDokumentu {
  return {
    dokumenty(idOkna) {
      return wykazZasobow(kanal, idOkna, DesignAssetKind.Document);
    },

    zasoby(idOkna) {
      return wykazZasobow(kanal, idOkna);
    },

    wykonaj(komenda, zadanie) {
      // Kształt treści żądania jest znany dopiero w czasie działania; sprawdzenie robi rdzeń.
      return wywolaj(kanal, komenda, zadanie as never) as Promise<Wynik<unknown>>;
    },
  };
}
