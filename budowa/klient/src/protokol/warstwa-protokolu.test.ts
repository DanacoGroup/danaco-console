import {
  Command,
  ErrorCode,
  EnvelopeStatus,
  EventType,
  ZDARZENIA_NIEZNANEJ,
  type Envelope,
} from '../../../shared/contract.ts';
import { GniazdoZastepcze, podstawGniazdo } from '../gniazdo-zastepcze.ts';
import { utworzTransport, type Transport } from '../polaczenie/gniazdo.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia.ts';
import { bieg, poOdstepie, rowne, sprawdz } from '../sprawdzian.ts';
import { utworzKanal } from './kanal.ts';
import { czyOdpowiedz, tresc, zbudujKoperte } from './koperta.ts';
import { utworzKorelacje } from './korelacja.ts';
import { sprawdzKsztalt, czyTekst } from './ksztalt-odpowiedzi.ts';
import { odczytajRamke, zapiszRamke } from './ramka.ts';
import { utworzSesje } from './sesja.ts';

/**
 * Warstwa protokołu klienta.
 *
 * Mierzone są trzy rzeczy: koperta i ramka w obie strony, wiązanie odpowiedzi
 * z żądaniem oraz odbiór zdarzeń — wszystkich, jakie zna kontrakt, bo warstwa
 * nie wybiera spośród nich i żadnego nie wyróżnia nazwą wpisaną w kod.
 */

/** Wszystkie zdarzenia kontraktu poza zapasowymi `*.unknown`. */
const ZDARZENIA_KONTRAKTU: EventType[] = (() => {
  const zapasowe = new Set<EventType>(Object.values(ZDARZENIA_NIEZNANEJ));
  return Object.values(EventType).filter((zdarzenie) => !zapasowe.has(zdarzenie));
})();

/** Zdarzenia zapasowe obszarów, bez powtórzeń. */
const ZDARZENIA_ZAPASOWE: EventType[] = [...new Set<EventType>(Object.values(ZDARZENIA_NIEZNANEJ))];

/** Transport w pamięci: ramki wychodzące odkłada, przychodzące podaje na żądanie. */
interface TransportZastepczy extends Transport {
  /** Ramki, które warstwa protokołu oddała transportowi. */
  wyslane: string[];
  /** Podaje ramkę tak, jakby przyszła z rdzenia. */
  podaj(ramka: string): void;
}

function utworzTransportZastepczy(): TransportZastepczy {
  const sluchaczeRamek: ((ramka: string) => void)[] = [];
  const wyslane: string[] = [];
  return {
    wyslane,
    polacz() {},
    rozlacz() {},
    wyslij: (ramka) => void wyslane.push(ramka),
    naRamke(sluchacz): Odsubskrybuj {
      sluchaczeRamek.push(sluchacz);
      return () => {
        const miejsce = sluchaczeRamek.indexOf(sluchacz);
        if (miejsce >= 0) sluchaczeRamek.splice(miejsce, 1);
      };
    },
    naStan(sluchacz): Odsubskrybuj {
      sluchacz('polaczony' as StanPolaczenia);
      return () => {};
    },
    stan: () => 'polaczony' as StanPolaczenia,
    oczekujace: () => 0,
    podaj(ramka) {
      for (const sluchacz of [...sluchaczeRamek]) sluchacz(ramka);
    },
  };
}

/** Kanał osadzony na transporcie w pamięci. */
function zalozKanal(): { kanal: ReturnType<typeof utworzKanal>; transport: TransportZastepczy } {
  const transport = utworzTransportZastepczy();
  return { kanal: utworzKanal(transport, utworzSesje()), transport };
}

/** Koperta zdarzenia w postaci ramki, tak jak nadaje ją rdzeń. */
function ramkaZdarzenia(typ: EventType, ladunek: unknown, dodatkowe: Partial<Envelope> = {}): string {
  return JSON.stringify({ type: typ, id: '', payload: ladunek, timestamp: 1, ...dodatkowe });
}

/** Koperta odpowiedzi na żądanie o podanym identyfikatorze. */
function ramkaOdpowiedzi(typ: string, idZadania: string, ladunek: unknown): string {
  return JSON.stringify({
    type: typ,
    id: idZadania,
    payload: ladunek,
    timestamp: 1,
    status: EnvelopeStatus.Ok,
  });
}

await bieg('warstwa protokołu', {
  'koperta wychodząca niesie identyfikator, czas i treść'() {
    const koperta = zbudujKoperte(Command.ConnectionHello, '', { clientId: 'k' });
    sprawdz(koperta.type === Command.ConnectionHello, 'typ koperty');
    sprawdz(koperta.id.startsWith('zadanie-'), `identyfikator żądania: ${koperta.id}`);
    sprawdz(typeof koperta.timestamp === 'number' && koperta.timestamp > 0, 'czas nadania');
    rowne(koperta.payload, { clientId: 'k' }, 'treść koperty');
  },

  'koperta pomija sesję pustą i niesie sesję założoną'() {
    const bezSesji = zbudujKoperte(Command.ConnectionHello, '', {});
    sprawdz(!('sessionId' in bezSesji), 'sesja pusta trafiła do koperty');

    const zSesja = zbudujKoperte(Command.SessionStop, 'sesja-1', {});
    sprawdz(zSesja.sessionId === 'sesja-1', 'sesja założona nie trafiła do koperty');
  },

  'koperta z polem status jest odpowiedzią, bez niego — zdarzeniem'() {
    const zdarzenie = odczytajRamke(ramkaZdarzenia(EventType.SessionChanged, {}));
    sprawdz(!czyOdpowiedz(zdarzenie), 'zdarzenie wzięte za odpowiedź');

    const odpowiedz = odczytajRamke(ramkaOdpowiedzi(Command.ConnectionHello, 'zadanie-1', {}));
    sprawdz(czyOdpowiedz(odpowiedz), 'odpowiedź nierozpoznana');
  },

  'ramka wraca z zapisu i odczytu bez zmiany'() {
    const koperta = zbudujKoperte(Command.ModuleList, 'sesja-1', { environmentId: 'srod' });
    const odczytana = odczytajRamke(zapiszRamke(koperta));
    rowne(odczytana, koperta, 'koperta po obiegu zapis–odczyt');
  },

  'ramka nieczytelna wraca jako zdarzenie zapasowe z treścią surową'() {
    const odczytana = odczytajRamke('to nie jest JSON');
    sprawdz(
      ZDARZENIA_ZAPASOWE.includes(odczytana.type as EventType),
      `ramka nieczytelna dała typ ${odczytana.type}`,
    );
    rowne(odczytana.payload, 'to nie jest JSON', 'treść surowa zachowana');
  },

  'ramka o kształcie niezgodnym z kopertą wraca jako zdarzenie zapasowe'() {
    const odczytana = odczytajRamke('{"cos":"innego"}');
    sprawdz(
      ZDARZENIA_ZAPASOWE.includes(odczytana.type as EventType),
      `ramka bez pola type dała typ ${odczytana.type}`,
    );
    rowne(odczytana.payload, { cos: 'innego' }, 'treść zachowana');
  },

  'korelacja oddaje odpowiedź jej odbiorcy dokładnie raz'() {
    const korelacja = utworzKorelacje();
    const odebrane: Envelope[] = [];
    korelacja.zarejestruj('zadanie-1', (odpowiedz) => odebrane.push(odpowiedz));
    sprawdz(korelacja.oczekujace() === 1, 'żądanie nie czeka na odpowiedź');

    const odpowiedz = odczytajRamke(ramkaOdpowiedzi(Command.ConnectionHello, 'zadanie-1', { a: 1 }));
    sprawdz(korelacja.rozstrzygnij(odpowiedz), 'odbiorca nie został odnaleziony');
    sprawdz(!korelacja.rozstrzygnij(odpowiedz), 'odbiorca odebrał odpowiedź drugi raz');
    sprawdz(odebrane.length === 1, `odbiorca wywołany ${odebrane.length} razy`);
    sprawdz(korelacja.oczekujace() === 0, 'żądanie zostało w rejestrze po odpowiedzi');
  },

  'korelacja nie rejestruje żądania bez identyfikatora'() {
    const korelacja = utworzKorelacje();
    korelacja.zarejestruj('', () => {});
    sprawdz(korelacja.oczekujace() === 0, 'żądanie bez identyfikatora weszło do rejestru');
  },

  'kanał wysyła komendę kontraktu i wiąże z nią odpowiedź rdzenia'() {
    const { kanal, transport } = zalozKanal();
    let odebrany: unknown = null;
    const idZadania = kanal.wyslij(
      Command.ConnectionHello,
      { clientId: 'k', clientVersion: '1.0.0', protocolVersion: '1.0' },
      (wynik) => {
        odebrany = wynik;
      },
    );

    sprawdz(transport.wyslane.length === 1, 'komenda nie poszła do transportu');
    const wyslana = odczytajRamke(transport.wyslane[0] ?? '');
    sprawdz(wyslana.type === Command.ConnectionHello, `typ wysłanej komendy: ${wyslana.type}`);
    sprawdz(wyslana.id === idZadania, 'identyfikator wysłanej koperty');

    transport.podaj(ramkaOdpowiedzi(Command.ConnectionHello, idZadania, { serverVersion: '1' }));
    rowne(odebrany, { udany: true, wynik: { serverVersion: '1' } }, 'wynik komendy');
  },

  'kanał oddaje odmowę rdzenia jako wynik nieudany z błędem kontraktu'() {
    const { kanal, transport } = zalozKanal();
    let odebrany: unknown = null;
    const idZadania = kanal.wyslij(Command.ModuleList, {}, (wynik) => {
      odebrany = wynik;
    });

    transport.podaj(
      JSON.stringify({
        type: Command.ModuleList,
        id: idZadania,
        timestamp: 1,
        status: EnvelopeStatus.Error,
        error: { code: ErrorCode.ValidationFailed, message: 'brak pola', retryable: false },
      }),
    );

    rowne(
      odebrany,
      {
        udany: false,
        blad: { code: ErrorCode.ValidationFailed, message: 'brak pola', retryable: false },
      },
      'wynik odmowy',
    );
  },

  'sprawdzian kształtu zamienia odpowiedź bez pola obowiązkowego w odmowę'() {
    const zdana = sprawdzKsztalt({ udany: true, wynik: { serverVersion: '1' } }, 'próba', (t) =>
      czyTekst(t.serverVersion),
    );
    sprawdz(zdana.udany, 'odpowiedź poprawna została odrzucona');

    const odrzucona = sprawdzKsztalt<{ serverVersion?: string }>(
      { udany: true, wynik: {} },
      'próba',
      (t) => czyTekst(t.serverVersion),
    );
    sprawdz(!odrzucona.udany, 'odpowiedź bez pola obowiązkowego przeszła');
    sprawdz(
      odrzucona.blad?.code === ErrorCode.ValidationFailed,
      `kod odmowy: ${odrzucona.blad?.code}`,
    );
  },

  'kanał dostarcza każde z 75 zdarzeń kontraktu do subskrybenta jego nazwy'() {
    const { kanal, transport } = zalozKanal();
    sprawdz(
      ZDARZENIA_KONTRAKTU.length === 75,
      `kontrakt niesie ${ZDARZENIA_KONTRAKTU.length} zdarzeń poza zapasowymi, oczekiwano 75`,
    );

    const nieodebrane: string[] = [];
    for (const zdarzenie of ZDARZENIA_KONTRAKTU) {
      let odebrana: unknown = null;
      let typKoperty = '';
      const odsubskrybuj = kanal.naZdarzenie(zdarzenie, (tresc, koperta) => {
        odebrana = tresc;
        typKoperty = koperta.type;
      });
      transport.podaj(ramkaZdarzenia(zdarzenie, { znacznik: zdarzenie }));
      odsubskrybuj();
      if (JSON.stringify(odebrana) !== JSON.stringify({ znacznik: zdarzenie })) {
        nieodebrane.push(`${zdarzenie} (treść)`);
      } else if (typKoperty !== zdarzenie) {
        nieodebrane.push(`${zdarzenie} (koperta)`);
      }
    }
    sprawdz(
      nieodebrane.length === 0,
      `zdarzenia bez odbioru: ${nieodebrane.join(', ')}`,
    );
    console.log(`           odebrano ${ZDARZENIA_KONTRAKTU.length} zdarzeń kontraktu`);
  },

  'kanał nie myli zdarzenia z odpowiedzią o tej samej nazwie'() {
    const { kanal, transport } = zalozKanal();
    const odebrane: unknown[] = [];
    kanal.naZdarzenie(EventType.SessionChanged, (t) => void odebrane.push(t));

    transport.podaj(
      ramkaZdarzenia(EventType.SessionChanged, { odpowiedz: true }, {
        status: EnvelopeStatus.Ok,
      }),
    );
    transport.podaj(ramkaZdarzenia(EventType.SessionChanged, { zdarzenie: true }));

    rowne(odebrane, [{ zdarzenie: true }], 'odebrane przez subskrybenta zdarzenia');
  },

  'dziennik nieznanych zbiera każde z 68 zdarzeń zapasowych kontraktu'() {
    const { kanal, transport } = zalozKanal();
    const dziennik = kanal.dziennikNieznanych();
    sprawdz(
      ZDARZENIA_ZAPASOWE.length === 68,
      `kontrakt niesie ${ZDARZENIA_ZAPASOWE.length} zdarzeń zapasowych, oczekiwano 68`,
    );

    for (const zdarzenie of ZDARZENIA_ZAPASOWE) {
      transport.podaj(
        ramkaZdarzenia(zdarzenie, {
          requestedType: `${zdarzenie}-zrodlo`,
          requestId: 'zadanie-1',
          reason: 'nieznana komenda',
        }),
      );
    }

    sprawdz(
      dziennik.liczba() === ZDARZENIA_ZAPASOWE.length,
      `dziennik zapisał ${dziennik.liczba()} z ${ZDARZENIA_ZAPASOWE.length} zdarzeń zapasowych`,
    );
    sprawdz(dziennik.ostatni()?.powod === 'nieznana komenda', 'powód ostatniego wpisu');
    console.log(`           zapisano ${dziennik.liczba()} zdarzeń zapasowych`);
  },

  'dziennik nieznanych zbiera kopertę o typie spoza kontraktu'() {
    const { kanal, transport } = zalozKanal();
    const dziennik = kanal.dziennikNieznanych();

    transport.podaj(ramkaZdarzenia('cos.czego.nie.ma' as EventType, { pole: 1 }));

    sprawdz(dziennik.liczba() === 1, `dziennik zapisał ${dziennik.liczba()} wpisów`);
    sprawdz(
      dziennik.ostatni()?.typZdarzenia === 'cos.czego.nie.ma',
      `typ wpisu: ${dziennik.ostatni()?.typZdarzenia}`,
    );
  },

  'treść koperty wraca w kształcie zapowiedzianym przez kontrakt'() {
    const koperta = odczytajRamke(ramkaZdarzenia(EventType.StreamChunk, { text: 'porcja' }));
    rowne(tresc<{ text: string }>(koperta), { text: 'porcja' }, 'treść fragmentu strumienia');
  },

  async 'zerwanie połączenia w trakcie strumienia i powrót nie gubią zdarzeń'() {
    const przywroc = podstawGniazdo();
    try {
      const transport = utworzTransport('ws://127.0.0.1:17870/ws', { opoznienie: () => 0 });
      const kanal = utworzKanal(transport, utworzSesje());

      const numery: number[] = [];
      kanal.naZdarzenie(EventType.StreamChunk, (_, koperta) => {
        numery.push(koperta.seq ?? 0);
      });

      transport.polacz();
      const pierwsze = GniazdoZastepcze.otwarte[0];
      pierwsze?.otworz();
      pierwsze?.przyjmij(ramkaZdarzenia(EventType.StreamChunk, { text: 'a' }, { seq: 1 }));
      pierwsze?.przyjmij(ramkaZdarzenia(EventType.StreamChunk, { text: 'b' }, { seq: 2 }));

      // Zerwanie w trakcie strumienia: gniazdo pada, warstwa protokołu zostaje.
      pierwsze?.close();
      // Żądanie złożone przy rozłączeniu czeka, zamiast przepaść.
      kanal.wyslij(Command.SessionStop, { sessionId: 'sesja-1' });
      sprawdz(transport.oczekujace() === 1, 'żądanie z czasu rozłączenia zgubione');

      await poOdstepie();
      const drugie = GniazdoZastepcze.otwarte[1];
      sprawdz(drugie !== undefined, 'transport nie wznowił połączenia');
      drugie?.otworz();

      sprawdz(transport.oczekujace() === 0, 'kolejka nie opróżniła się po powrocie');
      sprawdz(drugie?.wyslane.length === 1, 'żądanie nie poszło po powrocie');

      drugie?.przyjmij(ramkaZdarzenia(EventType.StreamChunk, { text: 'c' }, { seq: 3 }));
      drugie?.przyjmij(
        ramkaZdarzenia(EventType.StreamChunk, { text: 'd' }, { seq: 4, done: true }),
      );

      rowne(numery, [1, 2, 3, 4], 'fragmenty strumienia odebrane przez subskrybenta');
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },
});
