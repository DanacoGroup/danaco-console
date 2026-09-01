import {
  Command,
  ErrorCode,
  type Envelope,
  type ErrorInfo,
  type EventPayloadOf,
  type EventType,
  type RequestOf,
  type ResponseOf,
} from '../../../shared/contract.ts';
import {
  zalozDziennikNieznanych,
  type DziennikNieznanych,
} from '../polaczenie/dziennik-nieznanych.ts';
import type { PowodPorzucenia, Transport } from '../polaczenie/gniazdo.ts';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { czyOdpowiedz, czyUdana, tresc, zbudujKoperte } from './koperta.ts';
import { utworzKorelacje } from './korelacja.ts';
import { odczytajRamke, zapiszRamke } from './ramka.ts';
import type { Sesja } from './sesja.ts';

/** Wynik komendy widziany przez wywołującego, niosący powodzenie, zwróconą wartość albo opis napotkanego błędu. */
export interface Wynik<T> {
  udany: boolean;
  wynik?: T;
  blad?: ErrorInfo;
}

/**
 * Kanał komunikatów — cienka warstwa nad kontraktem osadzona na transporcie,
 * nieznająca treści dziedzinowej. Nazwy komend i zdarzeń oraz kształty ich
 * treści pochodzą wyłącznie ze współdzielonego kontraktu.
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
  /** Zapamiętuje wykaz komend obsługiwanych przez rdzeń, podany w powitaniu. */
  zapamietajKomendy(komendy: readonly string[]): void;
  /** Dziennik komunikatów nierozpoznanych — brama fail-open kanału. */
  dziennikNieznanych(): DziennikNieznanych;
}

export function utworzKanal(transport: Transport, sesja: Sesja): Kanal {
  const przychodzace = utworzMagistrale<Envelope>();
  const korelacja = utworzKorelacje();
  /* Wykaz komend rdzenia z powitania. Pusty znaczy: rdzeń jeszcze nie mówił,
     więc kanał niczego nie odsiewa — powitanie samo idzie przed tym wykazem. */
  let komendyRdzenia: ReadonlySet<string> = new Set<string>();

  transport.naRamke((ramka) => {
    const koperta = odczytajRamke(ramka);
    if (czyOdpowiedz(koperta)) korelacja.rozstrzygnij(koperta);
    przychodzace.oglos(koperta);
  });

  /* Ramka porzucona nie doręczy już nic: wywołujący ma dostać odmowę nazwaną
     w chwili porzucenia, nie ciszę do upływu terminu korelacji. */
  transport.naPorzucona((ramka, powod) => {
    korelacja.odmow(odczytajRamke(ramka).id, bladPorzucenia(powod));
  });

  /* Zerwane gniazdo nie przyniesie odpowiedzi na żądania, które na nim stały —
     rdzeń wiąże je z połączeniem, a nowe gniazdo o nich nie wie. */
  transport.naStan((stan) => {
    if (stan === 'polaczony') return;
    korelacja.uniewaznijWszystkie(bladZerwania());
  });

  const kanal: Kanal = {
    wyslij(komenda, zadanie, przyWyniku) {
      const koperta = zbudujKoperte(komenda, sesja.id(), zadanie);
      if (komendyRdzenia.size > 0 && !komendyRdzenia.has(komenda)) {
        przyWyniku?.({ udany: false, blad: bladKomendyNieznanej(komenda) });
        return koperta.id;
      }
      if (przyWyniku !== undefined) {
        korelacja.zarejestruj(koperta.id, komenda, (odpowiedz) =>
          przyWyniku(zbudujWynik(odpowiedz)),
        );
      }
      const ramka = zapiszRamke(koperta);
      if (komenda === Command.ConnectionHello) transport.wyslijPowitanie(ramka);
      else transport.wyslij(ramka);
      return koperta.id;
    },

    naZdarzenie(zdarzenie, sluchacz) {
      return przychodzace.subskrybuj((koperta) => {
        if (koperta.type !== zdarzenie || czyOdpowiedz(koperta)) return;
        sluchacz(
          tresc<EventPayloadOf<typeof zdarzenie>>(koperta) as EventPayloadOf<typeof zdarzenie>,
          koperta,
        );
      });
    },

    naDowolny(sluchacz) {
      return przychodzace.subskrybuj(sluchacz);
    },

    sesja: () => sesja,

    zapamietajKomendy(komendy) {
      komendyRdzenia = new Set<string>(komendy);
    },

    dziennikNieznanych: () => dziennik,
  };

  // Dziennik zakładany od razu — komunikat może przyjść przed pierwszą subskrypcją.
  const dziennik = zalozDziennikNieznanych(kanal);

  return kanal;
}

/** Przekłada kopertę odpowiedzi na wynik komendy, wyodrębniając powodzenie, treść albo błąd zgodnie z kontraktem. */
function zbudujWynik<T>(odpowiedz: Envelope): Wynik<T> {
  if (!czyUdana(odpowiedz)) {
    return { udany: false, blad: odpowiedz.error };
  }
  return { udany: true, wynik: tresc<T>(odpowiedz) };
}

/** Odmowa komendy, której rdzeń nie wymienił w powitaniu — wysłanie jej wróciłoby odmową nierozpoznania. */
function bladKomendyNieznanej(komenda: string): ErrorInfo {
  return {
    code: ErrorCode.NotFound,
    message: `Rdzeń nie obsługuje komendy ${komenda}`,
    retryable: false,
  };
}

/** Odmowa żądania porzuconego w kolejce wychodzącej, nazywająca powód porzucenia. */
function bladPorzucenia(powod: PowodPorzucenia): ErrorInfo {
  return {
    code: ErrorCode.ChannelUnavailable,
    message:
      powod === 'zapora-czasu'
        ? 'Żądanie czekało na łączność dłużej, niż wolno — nie zostało wysłane'
        : 'Łączność zerwana, zanim żądanie wyszło do rdzenia',
    retryable: true,
  };
}

/** Odmowa żądania, które stało na zerwanym połączeniu. */
function bladZerwania(): ErrorInfo {
  return {
    code: ErrorCode.ChannelUnavailable,
    message: 'Łączność z rdzeniem zerwana przed nadejściem odpowiedzi',
    retryable: true,
  };
}
