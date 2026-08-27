import {
  Command,
  type AccessPointListResponse,
  type AccessPointCheckResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

/**
 * Katalog punktów dostępu widziany przez Integrations Hub oraz Permissions
 * & Trust Center: komenda `access.point.list` podaje wykaz punktów,
 * a `access.point.check` odpowiada, czy most odpowiada i jakie korzenie
 * potwierdza.
 */
export interface ZrodloPunktowDostepu {
  /** `access.point.list` — katalog punktów w kolejności nazw. */
  punkty(): Promise<Wynik<AccessPointListResponse>>;
  /** `access.point.check` — sprawdzenie, czy punkt odpowiada. */
  sprawdz(idPunktu: string): Promise<Wynik<AccessPointCheckResponse>>;
}

export function utworzZrodloPunktowDostepu(kanal: Kanal): ZrodloPunktowDostepu {
  const wywolaj: WywolanieApps = utworzWywolanieApps(kanal);

  return {
    async punkty() {
      return sprawdzKsztalt(
        await wywolaj(Command.AccessPointList, {}),
        Command.AccessPointList,
        (tresc) => czyTablica(tresc.points),
      );
    },

    async sprawdz(idPunktu) {
      return sprawdzKsztalt(
        await wywolaj(Command.AccessPointCheck, { accessPointId: idPunktu }),
        Command.AccessPointCheck,
        // Stan jest polem wymaganym i niesie całą odpowiedź na pytanie okna o sprawdzenie punktu.
        (tresc) => czyTekst(tresc.status),
      );
    },
  };
}
