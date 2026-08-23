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
 * Droga okna warsztatu do rdzenia: wykaz materiału i wykonanie czynności.
 *
 * Materiał bierze się z magazynu zasobów (`design.asset.list`), a nie z obszaru
 * `studio`: warsztat pracuje na dokumentach wniesionych do okna, a te leżą
 * w jednym magazynie całego produktu. Drugi magazyn dokumentów byłby drugim
 * miejscem, w którym ta sama treść żyje.
 *
 * `wykonaj` jest jedną drogą na piętnaście komend, bo wszystkie idą tak samo:
 * nazwa komendy ze stałych kontraktu, treść żądania złożona z formularza. Nazwa
 * nie jest tu nigdy napisem wpisanym z pamięci — przychodzi z katalogu czynności,
 * a ten bierze ją ze stałych.
 */
export interface ZrodloWarsztatuDokumentu {
  /** Dokumenty leżące w magazynie okna — materiał czynności warsztatu. */
  dokumenty(idOkna: string): Promise<Wynik<DesignAsset[]>>;
  /** Zasoby okna bez zawężenia rodzaju — pieczęcie graficzne i certyfikaty. */
  zasoby(idOkna: string): Promise<Wynik<DesignAsset[]>>;
  /** Wykonuje czynność warsztatu wskazaną komendą kontraktu. */
  wykonaj(komenda: Command, zadanie: Record<string, unknown>): Promise<Wynik<unknown>>;
}

/** Odczyt wykazu zasobów okna z zawężeniem rodzaju albo bez niego. */
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
      // Treść żądania powstaje z formularza, więc jej kształt znany jest dopiero
      // w czasie działania. Sprawdzenie kształtu robi rdzeń i odsyła odmowę
      // walidacji z nazwą pola — a to jest sprawdzenie, którego klient i tak nie
      // zastąpi, bo kontrakt rozstrzyga po stronie rdzenia.
      return wywolaj(kanal, komenda, zadanie as never) as Promise<Wynik<unknown>>;
    },
  };
}
