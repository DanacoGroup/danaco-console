import {
  Command,
  type ContextBundle,
  type ContextTransferResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Kod modułu Library w rdzeniu — miejsce docelowe przekazania kontekstu dokumentu z okna Preview Window. */
export const KOD_LIBRARY = 'library';

/**
 * Przekazanie dokumentu do Library jedną komendą `context.transfer`; komplet niesie identyfikator
 * dokumentu, a polecenie wyjściowe zostaje puste, bo podgląd nie zleca pracy, tylko przenosi wynik.
 */
export interface ZrodloPrzekazania {
  doLibrary(
    idOknaZrodlowego: string,
    komplet: ContextBundle,
  ): Promise<Wynik<ContextTransferResponse>>;
}

export function utworzZrodloPrzekazania(kanal: Kanal): ZrodloPrzekazania {
  return {
    async doLibrary(idOknaZrodlowego, komplet) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ContextTransfer, {
          sourceWindowId: idOknaZrodlowego,
          targetModuleId: KOD_LIBRARY,
          bundle: komplet,
        }),
        Command.ContextTransfer,
        (tresc) => czyObiekt(tresc.window),
      );
    },
  };
}
