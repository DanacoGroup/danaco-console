import {
  Command,
  EventType,
  type DeviceChangedEvent,
  type DeviceListResponse,
  type DeviceRevokeResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLogiczna, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/** Czynności na urządzeniach powiązanych z kontem obejmują dwie komendy: wykaz urządzeń i unieważnienie tokenu jednego urządzenia, bo przemianowania ani odłączenia kontrakt nie zna. */
export interface ZrodloUrzadzen {
  /** Wykaz urządzeń powiązanych z kontem właściciela; jedyna droga do wykazu przy wejściu do sekcji. */
  wykaz(): Promise<Wynik<DeviceListResponse>>;
  /** Unieważnienie tokenu wskazanego urządzenia; wynik niesie potwierdzenie, wykaz idzie zdarzeniem. */
  uniewaznij(idUrzadzenia: string): Promise<Wynik<DeviceRevokeResponse>>;
  /** Subskrypcja zmiany urządzeń: jedyna droga, którą wykaz nadąża za czynnością poza tym oknem. */
  naZmianeUrzadzen(sluchacz: (zmiana: DeviceChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloUrzadzen(kanal: Kanal): ZrodloUrzadzen {
  return {
    async wykaz() {
      // Kształt sprawdzamy, bo rzutowanie kanału jest obietnicą kompilatora, nie rdzenia.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeviceList, {}),
        Command.DeviceList,
        (tresc) => czyTablica(tresc.devices),
      );
    },

    async uniewaznij(idUrzadzenia) {
      // Brak potwierdzenia bez sprawdzianu wyglądałby jak czynność udana.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeviceRevoke, { deviceId: idUrzadzenia }),
        Command.DeviceRevoke,
        (tresc) => czyLogiczna(tresc.revoked),
      );
    },

    naZmianeUrzadzen(sluchacz) {
      return kanal.naZdarzenie(EventType.DeviceChanged, (tresc) => sluchacz(tresc));
    },
  };
}
