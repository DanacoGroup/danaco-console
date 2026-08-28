import {
  Command,
  EventType,
  type WindowUpdateRequest,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { niepowodzenie, potwierdzenie, type OdbiorcaKomunikatu } from './komunikat-zmiany';
import type { StanSterowania } from './stan-sterowania';

/** Zmiana pól okna niesionych treścią żądania aktualizacji okna, bez identyfikatora — okno wskazuje stan. */
export type ZmianaPolOkna = Omit<WindowUpdateRequest, 'windowId'>;

/** Wysyłka zmiany ustawienia okna i odbiór potwierdzenia komendą aktualizacji okna oraz zdarzeniem zmiany. */
export interface ZmianaOkna {
  /** Wysyła zmianę wskazanych pól okna. */
  zastosuj(nazwa: string, zmiana: ZmianaPolOkna): void;
  /** Odłącza subskrypcję zdarzeń. */
  rozlacz(): void;
}

export function utworzZmianeOkna(
  kanal: Kanal,
  stan: StanSterowania,
  zglos: OdbiorcaKomunikatu,
): ZmianaOkna {
  const odsubskrybuj: Odsubskrybuj = kanal.naZdarzenie(
    EventType.WindowChanged,
    (tresc) => {
      if (tresc.window.id !== stan.idOkna()) return;
      stan.przyjmijOkno(tresc.window);
    },
  );

  return {
    zastosuj(nazwa, zmiana) {
      kanal.wyslij(
        Command.WindowUpdate,
        { ...zmiana, windowId: stan.idOkna() },
        (wynik) => {
          const okno = wynik.wynik?.window;
          if (wynik.udany && okno !== undefined) {
            stan.przyjmijOkno(okno);
            zglos(potwierdzenie(nazwa));
            return;
          }
          stan.odswiez();
          zglos(niepowodzenie(nazwa, wynik.blad));
        },
      );
    },

    rozlacz: odsubskrybuj,
  };
}
