import type {
  Command,
  Envelope,
  ErrorInfo,
  EventPayloadOf,
  EventType,
  RequestOf,
  ResponseOf,
} from '../../../shared/contract';
import type { Transport } from '../polaczenie/gniazdo';
import { zalozDziennikNieznanych, type DziennikNieznanych } from '../polaczenie/indeks';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { czyOdpowiedz, czyUdana, tresc, zbudujKoperte } from './koperta';
import { utworzKorelacje } from './korelacja';
import { odczytajRamke, zapiszRamke } from './ramka';
import type { Sesja } from './sesja';

/** Wynik komendy widziany przez wywołującego. */
export interface Wynik<T> {
  udany: boolean;
  wynik?: T;
  blad?: ErrorInfo;
}

/**
 * Kanał komunikatów — cienka warstwa nad kontraktem osadzona na transporcie.
 *
 * Kanał nie zna treści dziedzinowej. Nazwy komend i zdarzeń oraz kształty ich
 * treści pochodzą wyłącznie z `shared/contract.ts`: zmiana nazwy
 * w `contract.json` przerywa kompilację klienta. Komunikat nierozpoznany nie
 * jest odrzucany.
 */
export interface Kanal {
  /** Wysyła komendę kontraktu; zwraca identyfikator żądania. */
  wyslij<K extends Command>(
    komenda: K,
    zadanie: RequestOf<K>,
    przyWyniku?: (wynik: Wynik<ResponseOf<K>>) => void,
  ): string;
  /** Subskrypcja zdarzeń jednego typu wraz z ich treścią. */
  naZdarzenie<K extends EventType>(
    zdarzenie: K,
    sluchacz: (tresc: EventPayloadOf<K>, koperta: Envelope) => void,
  ): Odsubskrybuj;
  /** Subskrypcja całego ruchu przychodzącego. */
  naDowolny(sluchacz: (koperta: Envelope) => void): Odsubskrybuj;
  /** Sesja nadawana kopertom wychodzącym. */
  sesja(): Sesja;
  /** Dziennik komunikatów nierozpoznanych — brama fail-open kanału. */
  dziennikNieznanych(): DziennikNieznanych;
}

export function utworzKanal(transport: Transport, sesja: Sesja): Kanal {
  const przychodzace = utworzMagistrale<Envelope>();
  const korelacja = utworzKorelacje();

  transport.naRamke((ramka) => {
    const koperta = odczytajRamke(ramka);
    if (czyOdpowiedz(koperta)) korelacja.rozstrzygnij(koperta);
    przychodzace.oglos(koperta);
  });

  const kanal: Kanal = {
    wyslij(komenda, zadanie, przyWyniku) {
      const koperta = zbudujKoperte(komenda, sesja.id(), zadanie);
      if (przyWyniku !== undefined) {
        korelacja.zarejestruj(koperta.id, (odpowiedz) => przyWyniku(zbudujWynik(odpowiedz)));
      }
      transport.wyslij(zapiszRamke(koperta));
      return koperta.id;
    },

    naZdarzenie(zdarzenie, sluchacz) {
      return przychodzace.subskrybuj((koperta) => {
        if (koperta.type !== zdarzenie || czyOdpowiedz(koperta)) return;
        sluchacz(tresc<EventPayloadOf<typeof zdarzenie>>(koperta) as EventPayloadOf<typeof zdarzenie>, koperta);
      });
    },

    naDowolny(sluchacz) {
      return przychodzace.subskrybuj(sluchacz);
    },

    sesja: () => sesja,

    dziennikNieznanych: () => dziennik,
  };

  // Dziennik zakładamy od razu, bo komunikat nierozpoznany może przyjść przed
  // pierwszą subskrypcją widoku. Zapis i wpis do konsoli są jedyną reakcją:
  // ani zdarzenie `*.unknown`, ani koperta o typie spoza kontraktu nie zrywa
  // połączenia i nie blokuje sesji.
  const dziennik = zalozDziennikNieznanych(kanal);

  return kanal;
}

/** Przekłada kopertę odpowiedzi na wynik komendy. */
function zbudujWynik<T>(odpowiedz: Envelope): Wynik<T> {
  if (!czyUdana(odpowiedz)) {
    return { udany: false, blad: odpowiedz.error };
  }
  return { udany: true, wynik: tresc<T>(odpowiedz) };
}
