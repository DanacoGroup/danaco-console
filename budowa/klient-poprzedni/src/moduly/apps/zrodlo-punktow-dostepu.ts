import {
  Command,
  type AccessPointListResponse,
  type AccessPointCheckResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

/**
 * Katalog punktów dostępu widziany przez Integrations Hub i Permissions & Trust
 * Center.
 *
 * Serwer MCP jest w kontrakcie pozycją katalogu rozszerzeń, a most, którym się
 * do niego dochodzi — punktem dostępu (`AccessPoint`). Oba byty łączy pole
 * `Extension.accessPointId`, więc okno integracji potrzebuje obu wykazów:
 * z pierwszego bierze rodzaj i stan włączenia, z drugiego adres mostu, korzenie
 * i wynik ostatniego sprawdzenia.
 *
 * `access.point.check` jest jedyną drogą kontraktu do „testu połączenia”
 * i „monitora zdrowia”: odpowiada, czy punkt odpowiada, jakie korzenie
 * potwierdza i co poszło nie tak. Nie jest to sprawdzenie kondycji samego
 * rozszerzenia — sprawdzeniu podlega most, nie usługa za nim.
 *
 * Źródło nie ma własnego stanu ani nie zna okna: wykazy trzyma
 * `stan-rozszerzen.ts`, żeby cztery okna patrzyły na jeden zbiór.
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
        // Stan jest polem wymaganym i niesie całą odpowiedź na pytanie okna;
        // korzenie i szczegół bywają puste także przy sprawdzeniu udanym.
        (tresc) => czyTekst(tresc.status),
      );
    },
  };
}
