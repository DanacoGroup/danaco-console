import {
  Command,
  EventType,
  type WindowUpdateRequest,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { niepowodzenie, potwierdzenie, type OdbiorcaKomunikatu } from './komunikat-zmiany';
import type { StanSterowania } from './stan-sterowania';

/** Zmiana pól okna niesionych treścią `window.update`; okno wskazuje stan. */
export type ZmianaPolOkna = Omit<WindowUpdateRequest, 'windowId'>;

/**
 * Wysyłka zmiany ustawienia okna i odbiór potwierdzenia.
 *
 * Nazwa komendy i nazwa zdarzenia pochodzą z pakietu `shared` —
 * zmiana nazwy w `contract.json` przerywa kompilację tego pliku. Zmiana idzie
 * komendą `window.update`, potwierdzeniem jest odpowiedź na nią oraz zdarzenie
 * `window.changed` rozgłaszane do wszystkich urządzeń konta.
 *
 * Zdarzenie dotyczące innego okna nie zmienia tego stanu — filtr po
 * identyfikatorze okna jest adresowaniem. Niepowodzenie zmiany nie wstrzymuje
 * sterowania: widok wraca do stanu potwierdzonego i przyjmuje kolejną zmianę.
 */
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
