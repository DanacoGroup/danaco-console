import {
  Command,
  ConfigScope,
  type IsolationPolicyPreviewRequest,
  type IsolationPolicyPreviewResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

/**
 * Polityka izolacji obowiązująca oknu albo modułowi: komenda
 * `isolation.policy.preview` oddaje warunki wykonania rozstrzygnięte przez
 * platformę po ośmiu poziomach zasięgu, a nie uprawnienia nadane pojedynczemu
 * rozszerzeniu.
 */
export interface ZrodloIzolacjiApps {
  /** `isolation.policy.preview` — polityka obowiązująca oknu albo modułowi. */
  polityka(idOkna: string): Promise<Wynik<IsolationPolicyPreviewResponse>>;
}

/**
 * Źródło pyta o najwęższy poziom zasięgu, jaki zna: podanie okna każe rdzeniowi
 * rozstrzygnąć dziedziczenie aż do niego, a kod modułu wskazuje ten sam byt
 * poziomu `module`, którym rdzeń zna moduł.
 */
export function utworzZrodloIzolacjiApps(kanal: Kanal, kodModulu: string): ZrodloIzolacjiApps {
  const wywolaj: WywolanieApps = utworzWywolanieApps(kanal);

  return {
    async polityka(idOkna) {
      const zadanie: IsolationPolicyPreviewRequest =
        idOkna === ''
          ? { scope: ConfigScope.Module, scopeId: kodModulu }
          : { scope: ConfigScope.Window, scopeId: idOkna, windowId: idOkna };
      return sprawdzKsztalt(
        await wywolaj(Command.IsolationPolicyPreview, zadanie),
        Command.IsolationPolicyPreview,
        (tresc) => czyObiekt(tresc.policy),
      );
    },
  };
}
