import { Command, ConfigScope, type Action } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';

/**
 * Źródło katalogu akcji jednego modułu.
 *
 * Oddaje pozycje katalogu albo powód, dla którego ich nie ma. Powód jest
 * potrzebny widokowi: panel pusty bez wyjaśnienia byłby atrapą, panel pusty
 * z powodem jest stanem opisanym.
 */
export type ZrodloAkcjiModulu = (
  kodModulu: string,
  oddaj: (akcje: Action[], powod: string) => void,
) => void;

/**
 * Katalog akcji modułu pochodzi z rejestru rdzenia komendą `action.list` o zasięgu `module` i kluczu równym kodowi modułu, więc przestawienie okna na inny moduł jest jednym zapytaniem, nie inną gałęzią widoku.
 */
export function zrodloAkcjiKanalu(kanal: Kanal): ZrodloAkcjiModulu {
  return (kodModulu, oddaj) => {
    if (kodModulu.length === 0) {
      oddaj([], 'Moduł okna nie został jeszcze potwierdzony przez rdzeń.');
      return;
    }
    kanal.wyslij(
      Command.ActionList,
      { scope: ConfigScope.Module, scopeId: kodModulu, enabledOnly: true },
      (wynik) => {
        if (!wynik.udany) {
          oddaj([], `Katalog akcji nieodczytany — ${opisBledu(wynik.blad)}`);
          return;
        }
        const akcje = [...(wynik.wynik?.actions ?? [])].sort((a, b) => a.order - b.order);
        oddaj(
          akcje,
          akcje.length > 0 ? '' : `Rejestr rdzenia nie ma akcji zasięgu modułu ${kodModulu}.`,
        );
      },
    );
  };
}

/** Opis błędu kontraktu przedstawiany operatorowi, gdy katalog akcji modułu nie daje się odczytać z rejestru rdzenia. */
function opisBledu(blad: { message: string; code: string } | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}
