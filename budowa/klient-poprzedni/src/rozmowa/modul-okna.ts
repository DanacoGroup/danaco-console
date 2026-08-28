import { Command, EventType } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/** Interfejs opisuje śledzenie modułu, w którym w danej chwili pracuje okno komunikacji z rdzeniem systemu. */
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

/** Funkcja śledzi moduł okna czytany wprost z rdzenia komendą window.state.get oraz zdarzeniem window.changed, zamiast go zgadywać po stronie widoku. */
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
