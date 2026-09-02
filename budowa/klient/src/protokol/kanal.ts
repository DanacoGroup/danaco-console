// Kanał komunikatów nad transportem; nazwy i kształty treści bierze kontrakt.
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
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia.ts';
import { czyOdpowiedz, czyUdana, tresc, zbudujKoperte } from './koperta.ts';
import { utworzKorelacje } from './korelacja.ts';
import { odczytajRamke, zapiszRamke } from './ramka.ts';
import type { Sesja } from './sesja.ts';

export interface Wynik<T> {
  udany: boolean;
  wynik?: T;
  blad?: ErrorInfo;
}

export interface Kanal {
  wyslij<K extends Command>(
    komenda: K,
    zadanie: RequestOf<K>,
    przyWyniku?: (wynik: Wynik<ResponseOf<K>>) => void,
  ): string;
  naZdarzenie<K extends EventType>(
    zdarzenie: K,
    sluchacz: (tresc: EventPayloadOf<K>, koperta: Envelope) => void,
  ): Odsubskrybuj;
  naDowolny(sluchacz: (koperta: Envelope) => void): Odsubskrybuj;
  naStan(sluchacz: (stan: StanPolaczenia) => void): Odsubskrybuj;
  wznowPolaczenie(): void;
  sesja(): Sesja;
  zapamietajKomendy(komendy: readonly string[]): void;
  dziennikNieznanych(): DziennikNieznanych;
}

export function utworzKanal(transport: Transport, sesja: Sesja): Kanal {
  const przychodzace = utworzMagistrale<Envelope>();
  const korelacja = utworzKorelacje();
  // Wykaz pusty znaczy rdzeń przed powitaniem: kanał nie odsiewa wtedy niczego.
  let komendyRdzenia: ReadonlySet<string> = new Set<string>();

  transport.naRamke((ramka) => {
    const koperta = odczytajRamke(ramka);
    if (czyOdpowiedz(koperta)) {
      korelacja.rozstrzygnij(koperta);
      // Odpowiedź po terminie nie ma wołającego: jego obietnica padła odmową.
      if (korelacja.czySpozniona(koperta.id)) {
        dziennik.odnotujSpozniona(koperta);
        return;
      }
    }
    przychodzace.oglos(koperta);
  });

  // Odmowa idzie w chwili porzucenia, nie po upływie terminu korelacji.
  transport.naPorzucona((ramka, powod) => {
    korelacja.odmow(odczytajRamke(ramka).id, bladPorzucenia(powod));
  });

  // Rdzeń wiąże żądania z połączeniem: nowe gniazdo o poprzednich nie wie.
  let bylaLacznosc = false;
  transport.naStan((stan) => {
    if (stan === 'polaczony') {
      bylaLacznosc = true;
      return;
    }
    // Stan przed pierwszym połączeniem nie jest zerwaniem: unieważnienie zdjęłoby powitanie tej samej tury.
    if (!bylaLacznosc) return;
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

    naStan: (sluchacz) => transport.naStan(sluchacz),

    wznowPolaczenie: () => transport.wznow(),

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

function zbudujWynik<T>(odpowiedz: Envelope): Wynik<T> {
  if (!czyUdana(odpowiedz)) {
    return { udany: false, blad: odpowiedz.error };
  }
  return { udany: true, wynik: tresc<T>(odpowiedz) };
}

function bladKomendyNieznanej(komenda: string): ErrorInfo {
  return {
    code: ErrorCode.NotFound,
    message: `Rdzeń nie obsługuje komendy ${komenda}`,
    retryable: false,
  };
}

const TRESC_PORZUCENIA: Readonly<Record<PowodPorzucenia, string>> = {
  'zapora-czasu': 'Żądanie czekało na łączność dłużej, niż wolno — nie zostało wysłane',
  zerwanie: 'Łączność zerwana, zanim żądanie wyszło do rdzenia',
  przepelnienie: 'Kolejka wychodząca pełna — żądanie ustąpiło miejsca nowszemu',
};

function bladPorzucenia(powod: PowodPorzucenia): ErrorInfo {
  return {
    code: ErrorCode.ChannelUnavailable,
    message: TRESC_PORZUCENIA[powod],
    retryable: true,
  };
}

function bladZerwania(): ErrorInfo {
  return {
    code: ErrorCode.ChannelUnavailable,
    message: 'Łączność z rdzeniem zerwana przed nadejściem odpowiedzi',
    retryable: true,
  };
}
