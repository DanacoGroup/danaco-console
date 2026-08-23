import {
  Command,
  type IsolationPolicyPreviewRequest,
  type IsolationPolicyPreviewResponse,
  type IsolationScopeListRequest,
  type IsolationScopeListResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Dwa odczyty izolacji, których żąda macierz izolacji sesji modułu Browser.
 *
 * Panel modułu jest wyłącznie odczytem, więc źródło niesie wyłącznie komendy
 * odczytu. Zapis przełączników i przypisanie profilu należą do okna
 * konfiguracji punktów izolacji — dwa miejsca zapisujące tę samą politykę
 * dawałyby dwa różne zdania o tym, co obowiązuje.
 *
 * Rozdzielenie odpowiada temu, co o izolacji mówi opracowanie modułu
 * (rozdz. 6.3): moduł Browser wskazuje punkty izolacji właściwe przeglądaniu —
 * dostęp sieciowy procesu sesji, kontenery tożsamości, zakres pętli — a Operator
 * personalizuje je w oknie konfiguracji.
 */
export interface ZrodloIzolacji {
  /** `isolation.policy.preview` — polityka obowiązująca, bez zapisu. */
  polityka(
    zadanie: IsolationPolicyPreviewRequest,
  ): Promise<Wynik<IsolationPolicyPreviewResponse>>;
  /** `isolation.scope.list` — poziomy zasięgu w kolejności rozstrzygania. */
  poziomy(zadanie: IsolationScopeListRequest): Promise<Wynik<IsolationScopeListResponse>>;
}

export function utworzZrodloIzolacji(kanal: Kanal): ZrodloIzolacji {
  return {
    async polityka(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.IsolationPolicyPreview, zadanie),
        Command.IsolationPolicyPreview,
        (tresc) => czyObiekt(tresc.policy),
      );
    },

    async poziomy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.IsolationScopeList, zadanie),
        Command.IsolationScopeList,
        (tresc) => czyTablica(tresc.scopes),
      );
    },
  };
}
