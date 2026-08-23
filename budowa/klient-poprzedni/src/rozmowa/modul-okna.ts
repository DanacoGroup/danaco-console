import { Command, EventType } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/** Śledzenie modułu, w którym pracuje okno komunikacji. */
export interface SledzenieModulu {
  /** Kod modułu znany w tej chwili; pusty do pierwszej odpowiedzi rdzenia. */
  biezacy(): string;
  /** Subskrypcja zmiany modułu okna. */
  naZmiane(sluchacz: (kod: string) => void): Odsubskrybuj;
  /** Ustawia moduł wprost — droga dla powłoki po `workspace.enter`. */
  ustaw(kod: string): void;
  /** Odłącza subskrypcję zdarzeń kanału. */
  rozlacz(): void;
}

/**
 * Moduł okna czytany z rdzenia, nie zgadywany po stronie widoku.
 *
 * Dwie drogi, obie z kontraktu:
 *   1. `window.state.get` przy złożeniu — okno pyta, w czym pracuje, zamiast
 *      czekać na pierwszą zmianę. Bez tego okno wznowione po odświeżeniu
 *      klienta stałoby na module nieustalonym mimo znanego stanu w rdzeniu.
 *   2. `window.changed` w toku pracy — `workspace.enter` przestawia to samo
 *      okno na inny moduł zamiast je zamykać, a rdzeń rozgłasza zmianę
 *      zdarzeniem.
 *
 * Niepowodzenie zapytania nie jest błędem okna: moduł zostaje nieustalony,
 * a okno pracuje na arsenale wspólnym.
 */
export function sledzModulOkna(kanal: Kanal, idOkna: string, poczatkowy = ''): SledzenieModulu {
  const zmiany = utworzMagistrale<string>();
  let biezacy = poczatkowy;

  function ustaw(kod: string): void {
    if (kod.length === 0 || kod === biezacy) return;
    biezacy = kod;
    zmiany.oglos(kod);
  }

  kanal.wyslij(Command.WindowStateGet, { windowId: idOkna }, (wynik) => {
    if (!wynik.udany) return;
    ustaw(wynik.wynik?.window.moduleId ?? '');
  });

  const odsubskrybuj = kanal.naZdarzenie(EventType.WindowChanged, ({ window }) => {
    if (window.id !== idOkna) return;
    ustaw(window.moduleId);
  });

  return {
    biezacy: () => biezacy,
    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
    ustaw,
    rozlacz: odsubskrybuj,
  };
}
