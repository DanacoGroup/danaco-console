import {
  Command,
  EventType,
  type AccessPoint,
  type AccessPointAddRequest,
  type AccessPointAddResponse,
  type AccessPointChangedEvent,
  type AccessPointCheckResponse,
  type AccessPointListRequest,
  type AccessPointRemoveResponse,
  type AccessPointUpdateRequest,
  type AccessPointUpdateResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Interfejs obsługuje punkty dostępu pobierane i zapisywane przez rdzeń; nazwy
 * komend i zdarzeń pochodzą wyłącznie ze stałych kontraktu. Odczyt jest
 * fail-open i przy niepowodzeniu zwraca wykaz pusty, zapis zwraca `Wynik`
 * niosący błąd.
 */
export interface ZrodloPunktow {
  /** `access.point.list` — katalog punktów w kolejności nazw. */
  lista(zadanie?: AccessPointListRequest): Promise<AccessPoint[]>;
  /** `access.point.add` — założenie punktu. */
  dodaj(zadanie: AccessPointAddRequest): Promise<Wynik<AccessPointAddResponse>>;
  /** `access.point.update` — zmiana punktu; pominięte pole zostaje. */
  zmien(zadanie: AccessPointUpdateRequest): Promise<Wynik<AccessPointUpdateResponse>>;
  /** `access.point.remove` — usunięcie punktu wraz z nadaniami. */
  usun(punktID: string): Promise<Wynik<AccessPointRemoveResponse>>;
  /** `access.point.check` — sprawdzenie, czy punkt odpowiada. */
  sprawdz(punktID: string): Promise<Wynik<AccessPointCheckResponse>>;
  /** Subskrypcja `access.point.changed`. */
  naZmiane(sluchacz: (zdarzenie: AccessPointChangedEvent) => void): Odsubskrybuj;
}

/**
 * Powiadomienie o niepowodzeniu odczytu. Odczyt jest fail-open i wraca
 * z wykazem pustym, ale o niepowodzeniu ma wiedzieć też widok, nie tylko
 * dziennik, żeby Operator odróżnił pusty wykaz od odmowy rdzenia i od
 * świeżej instalacji.
 */
export type NaNiepowodzenie = (powod: string) => void;

export function utworzZrodloPunktow(
  kanal: Kanal,
  naNiepowodzenie: NaNiepowodzenie = () => undefined,
): ZrodloPunktow {
  return {
    async lista(zadanie = {}) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AccessPointList, zadanie),
        Command.AccessPointList,
        (tresc) => czyTablica(tresc.points),
      );
      if (!wynik.udany) {
        ostrzez(Command.AccessPointList, wynik.blad?.message);
        naNiepowodzenie(powodOdczytu(Command.AccessPointList, wynik.blad?.message));
        return [];
      }
      return uporzadkujPunkty(wynik.wynik?.points ?? []);
    },

    dodaj(zadanie) {
      return sprawdzPunkt(wywolaj(kanal, Command.AccessPointAdd, zadanie), Command.AccessPointAdd);
    },

    zmien(zadanie) {
      return sprawdzPunkt(
        wywolaj(kanal, Command.AccessPointUpdate, zadanie),
        Command.AccessPointUpdate,
      );
    },

    async usun(punktID) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AccessPointRemove, { accessPointId: punktID }),
        Command.AccessPointRemove,
        (tresc) => typeof tresc.removed === 'boolean',
      );
    },

    async sprawdz(punktID) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AccessPointCheck, { accessPointId: punktID }),
        Command.AccessPointCheck,
        (tresc) => typeof tresc.status === 'string' && czyLiczba(tresc.checkedAt),
      );
    },

    naZmiane(sluchacz) {
      return kanal.naZdarzenie(EventType.AccessPointChanged, (tresc) => {
        if (czyObiekt(tresc.point as unknown)) sluchacz(tresc);
      });
    },
  };
}

/** Sprawdza kształt odpowiedzi niosącej jeden punkt dostępu; funkcja jest wspólna dla żądania założenia i żądania zmiany punktu. */
async function sprawdzPunkt<T extends { point: AccessPoint }>(
  obietnica: Promise<Wynik<T>>,
  komenda: string,
): Promise<Wynik<T>> {
  return sprawdzKsztalt(await obietnica, komenda, (tresc) =>
    czyObiekt(tresc.point as unknown),
  );
}

/**
 * Porządek wykazu: rodzaj, potem nazwa. Kontrakt zapowiada kolejność nazw,
 * ale grupowanie po rodzaju należy do widoku, więc porządkujemy sami — wykaz
 * nie może zależeć od kolejności, w jakiej rdzeń akurat odpowiedział.
 */
export function uporzadkujPunkty(punkty: readonly AccessPoint[]): AccessPoint[] {
  return [...punkty].sort((pierwszy, drugi) =>
    pierwszy.kind === drugi.kind
      ? pierwszy.name.localeCompare(drugi.name, 'pl')
      : pierwszy.kind.localeCompare(drugi.kind),
  );
}

/** Zapisuje w dzienniku niepowodzenie odczytu wykazu punktów dostępu, wraz z pełną nazwą komendy i podanym powodem. */
function ostrzez(komenda: string, powod: string | undefined): void {
  console.warn('[dostępy] wykaz nie dotarł', komenda, powod ?? '');
}

/**
 * Zdanie o niepowodzeniu odczytu pokazywane Operatorowi. Rdzeń, który nie
 * podał powodu, nie zostawia pustego miejsca — zdanie mówi przynajmniej,
 * która komenda nie odpowiedziała.
 */
export function powodOdczytu(komenda: string, powod: string | undefined): string {
  const tresc = powod !== undefined && powod !== '' ? powod : 'rdzeń nie podał powodu';
  return `Komenda ${komenda} nie odpowiedziała wykazem: ${tresc}`;
}
