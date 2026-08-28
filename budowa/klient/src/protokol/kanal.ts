import type {
  Command,
  Envelope,
  ErrorInfo,
  EventPayloadOf,
  EventType,
  RequestOf,
  ResponseOf,
} from '../../../shared/contract.ts';
import {
  zalozDziennikNieznanych,
  type DziennikNieznanych,
} from '../polaczenie/dziennik-nieznanych.ts';
import type { Transport } from '../polaczenie/gniazdo.ts';
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
