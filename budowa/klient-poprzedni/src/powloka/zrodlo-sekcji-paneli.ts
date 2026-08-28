import {
  Command,
  type PanelSection,
  type PanelSectionsGetRequest,
  type PanelSectionsSetRequest,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { przenies } from '../protokol/wynik-czastkowy';
import { wywolaj } from '../protokol/wywolanie';

/** Źródło układu sekcji paneli — jedyne miejsce w kliencie, które zna nazwy komend odczytu i zapisu układu. */
export interface ZrodloSekcjiPaneli {
  /** `panel.sections.get` — układ sekcji panelu zapisany w rdzeniu. */
  uklad(zadanie: PanelSectionsGetRequest): Promise<Wynik<PanelSection[]>>;
  /** `panel.sections.set` — zapis układu; oddaje układ obowiązujący po zmianie. */
  zapiszUklad(zadanie: PanelSectionsSetRequest): Promise<Wynik<PanelSection[]>>;
}

export function utworzZrodloSekcjiPaneli(kanal: Kanal): ZrodloSekcjiPaneli {
  return {
    async uklad(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.PanelSectionsGet, zadanie),
        Command.PanelSectionsGet,
        (tresc) => czyTablica(tresc.sections),
      );
      return przenies(wynik, (tresc) => tresc.sections);
    },

    async zapiszUklad(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.PanelSectionsSet, zadanie),
        Command.PanelSectionsSet,
        (tresc) => czyTablica(tresc.sections),
      );
      return przenies(wynik, (tresc) => tresc.sections);
    },
  };
}
