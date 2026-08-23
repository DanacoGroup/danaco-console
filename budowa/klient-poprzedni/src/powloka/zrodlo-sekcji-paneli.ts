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

/**
 * Źródło układu sekcji paneli — jedyne miejsce w kliencie, które zna nazwy
 * `panel.sections.get` i `panel.sections.set`.
 *
 * Źródło stoi w powłoce, nie w module, bo sekcje panelu są własnością okna,
 * a nie dziedziny modułu: każde okno operacyjne dzieli treść na sekcje i każde
 * ma prawo je zwinąć, przestawić i zdjąć z widoku. Gdyby źródło stało w module,
 * drugi moduł zbudowałby drugie takie samo. Powłoka niesie wzorzec okna
 * (`komponenty/rama-okna`), więc niesie też wzorzec jego sekcji; moduł podaje
 * wyłącznie własne sekcje i identyfikator panelu.
 *
 * Zapis oddaje układ obowiązujący, nie żądany: `panel.sections.set` zwraca pole
 * `sections`, więc widok przerysowuje się z odpowiedzi, a nie z własnego
 * żądania. Sama zgoda rdzenia nie jest dowodem skutku — dowodem jest treść
 * odpowiedzi.
 *
 * Kształt odpowiedzi jest sprawdzany tak samo jak wszędzie: rdzeń, który
 * odpowie bez tablicy `sections`, nie dowiódł, że układ zna; wynik idzie wtedy
 * jako odmowa, a nie jako pusty układ udający „panel bez sekcji".
 */
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
