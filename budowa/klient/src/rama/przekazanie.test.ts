/**
 * Sprawdzian mierzy samo przejście — chwilę, w której przebieg drogi
 * wejścia oddaje sterowanie ramie — prowadząc prawdziwy przebieg przez
 * powitanie, wejście do środowiska i (osobno) odmowę rdzenia na tym kroku.
 */

import {
  Command,
  EnvelopeStatus,
  PROTOCOL_VERSION,
  type Envelope,
  type ErrorInfo,
} from '../../../shared/contract.ts';
import type { Transport } from '../polaczenie/gniazdo.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia.ts';
import { utworzKanal } from '../protokol/kanal.ts';
import { utworzSesje } from '../protokol/sesja.ts';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { bieg, poOdstepie, rowne, sprawdz } from '../sprawdzian.ts';
import { magazynWPamieci, utworzPrzebieg, type StanPrzebiegu } from '../wejscie/przebieg.ts';
import { gotowaDoPrzekazania } from './przekazanie.ts';

/* ── Rdzeń zastępczy — ten sam kształt co w sprawdzianie drogi wejścia ─────── */

type Odpowiedz = { tresc: unknown } | { blad: ErrorInfo };

interface RdzenZastepczy extends Transport {
  odpowiadaj(komenda: string, odpowiedz: Odpowiedz): void;
  ogloszStan(stan: StanPolaczenia): void;
}

function utworzRdzenZastepczy(): RdzenZastepczy {
  const sluchaczeRamek: ((ramka: string) => void)[] = [];
  const sluchaczeStanu: ((stan: StanPolaczenia) => void)[] = [];
  const odpowiedzi = new Map<string, Odpowiedz>();
  let biezacy: StanPolaczenia = 'rozlaczony';

  function oddaj(koperta: Envelope): void {
    const odpowiedz = odpowiedzi.get(koperta.type);
    if (odpowiedz === undefined) return;
    const zwrot: Envelope =
      'blad' in odpowiedz
        ? { type: koperta.type, id: koperta.id, timestamp: Date.now(), status: EnvelopeStatus.Error, error: odpowiedz.blad }
        : { type: koperta.type, id: koperta.id, timestamp: Date.now(), status: EnvelopeStatus.Ok, payload: odpowiedz.tresc };
    for (const sluchacz of [...sluchaczeRamek]) sluchacz(JSON.stringify(zwrot));
  }

  return {
    odpowiadaj(komenda, odpowiedz) {
      odpowiedzi.set(komenda, odpowiedz);
    },
    ogloszStan(stan) {
      biezacy = stan;
      for (const sluchacz of [...sluchaczeStanu]) sluchacz(stan);
    },
    polacz() {},
    rozlacz() {},
    wyslij(ramka) {
      oddaj(JSON.parse(ramka) as Envelope);
    },
    naRamke(sluchacz): Odsubskrybuj {
      sluchaczeRamek.push(sluchacz);
      return () => {
        const miejsce = sluchaczeRamek.indexOf(sluchacz);
        if (miejsce >= 0) sluchaczeRamek.splice(miejsce, 1);
      };
    },
    naStan(sluchacz): Odsubskrybuj {
      sluchacz(biezacy);
      sluchaczeStanu.push(sluchacz);
      return () => {
        const miejsce = sluchaczeStanu.indexOf(sluchacz);
        if (miejsce >= 0) sluchaczeStanu.splice(miejsce, 1);
      };
    },
    stan: () => biezacy,
    oczekujace: () => 0,
  };
}

const POWITANIE_Z_TOKENEM: Odpowiedz = {
  tresc: { serverVersion: PROTOCOL_VERSION, protocolVersion: PROTOCOL_VERSION, authenticated: true, gatewayConfigured: true },
};

const SRODOWISKA: Odpowiedz = {
  tresc: { environments: [{ id: '1', code: 'talkin', name: 'TalkIn', order: 1, navigationKind: 'modules' }] },
};

const WEJSCIE: Odpowiedz = {
  tresc: {
    environment: { id: '1', code: 'talkin', name: 'TalkIn', order: 1, navigationKind: 'modules' },
    modules: [{ id: '1', code: 'studio', name: 'Studio', order: 1, operationalWindowCodes: [], kind: 'srodowisko_robocze', configuredOnHome: false }],
    sessions: [],
  },
};

const KLIENT: TozsamoscKlienta = { id: 'klient-sprawdzianu', wersja: '1.0.0' };

function zaloz(): { przebieg: ReturnType<typeof utworzPrzebieg>; rdzen: RdzenZastepczy; gotowosci: boolean[] } {
  const rdzen = utworzRdzenZastepczy();
  const przebieg = utworzPrzebieg({
    kanal: utworzKanal(rdzen, utworzSesje()),
    transport: rdzen,
    klient: KLIENT,
    magazyn: magazynWPamieci(),
  });
  const gotowosci: boolean[] = [];
  przebieg.naZmiane((stan: StanPrzebiegu) => gotowosci.push(gotowaDoPrzekazania(stan)));
  return { przebieg, rdzen, gotowosci };
}

await bieg('rama aplikacji — przekazanie sterowania', {
  async 'sterowanie nie przechodzi do ramy, dopóki środowisko nie jest znane'() {
    const { przebieg, rdzen, gotowosci } = zaloz();
    rdzen.odpowiadaj(Command.ConnectionHello, POWITANIE_Z_TOKENEM);
    rdzen.odpowiadaj(Command.EnvironmentList, SRODOWISKA);
    // Wejście do środowiska celowo bez odpowiedzi jeszcze — sprawdza stan pośredni.
    przebieg.polacz();
    rdzen.ogloszStan('polaczony');
    await poOdstepie();
    sprawdz(!gotawoscBiezaca(przebieg), 'przejście nastąpiło bez odpowiedzi environment.enter');
    sprawdz(gotowosci.every((g) => g === false), 'gotowość zapaliła się przedwcześnie w którymś stanie pośrednim');
  },

  async 'sterowanie przechodzi do ramy dokładnie wtedy, gdy przebieg zna środowisko'() {
    const { przebieg, rdzen, gotowosci } = zaloz();
    rdzen.odpowiadaj(Command.ConnectionHello, POWITANIE_Z_TOKENEM);
    rdzen.odpowiadaj(Command.EnvironmentList, SRODOWISKA);
    rdzen.odpowiadaj(Command.EnvironmentEnter, WEJSCIE);
    przebieg.polacz();
    rdzen.ogloszStan('polaczony');
    await poOdstepie();
    await poOdstepie();

    sprawdz(gotawoscBiezaca(przebieg), 'przejście nie nastąpiło mimo udanego wejścia do środowiska');
    rowne(przebieg.stan().srodowisko?.name, 'TalkIn', 'nazwa środowiska dostępna ramie po przejściu');
    sprawdz(gotowosci.includes(false), 'nie zaobserwowano stanu sprzed przejścia');
    sprawdz(gotowosci.includes(true), 'nie zaobserwowano stanu po przejściu');
    rowne(gotowosci[gotowosci.length - 1], true, 'ostatni zaobserwowany stan jest stanem przejętym przez ramę');
  },

  async 'odmowa rdzenia przy wejściu do środowiska nie zapala przejścia'() {
    const { przebieg, rdzen, gotowosci } = zaloz();
    rdzen.odpowiadaj(Command.ConnectionHello, POWITANIE_Z_TOKENEM);
    rdzen.odpowiadaj(Command.EnvironmentList, SRODOWISKA);
    rdzen.odpowiadaj(Command.EnvironmentEnter, {
      blad: { code: 'internal_error', message: 'rdzeń zastępczy odmawia wejścia', retryable: false },
    });
    przebieg.polacz();
    rdzen.ogloszStan('polaczony');
    await poOdstepie();
    await poOdstepie();

    rowne(przebieg.stan().srodowisko, undefined, 'środowisko nie powinno być ustawione po odmowie');
    sprawdz(!gotawoscBiezaca(przebieg), 'przejście zapaliło się mimo odmowy rdzenia');
    sprawdz(gotowosci.every((g) => g === false), 'gotowość zapaliła się mimo braku środowiska');
  },
});

function gotawoscBiezaca(przebieg: ReturnType<typeof utworzPrzebieg>): boolean {
  return gotowaDoPrzekazania(przebieg.stan());
}
