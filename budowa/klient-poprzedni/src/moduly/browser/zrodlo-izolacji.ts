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
 * Dwa odczyty izolacji, których żąda macierz izolacji sesji modułu Browser:
 * polityka obowiązująca oraz poziomy zasięgu. Źródło niesie wyłącznie komendy
 * odczytu, ponieważ zapis przełączników i przypisanie profilu należą do okna
 * punktów izolacji.
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
