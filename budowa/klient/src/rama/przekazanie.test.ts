/**
 * Sprawdzian mierzy samo przejście, prowadząc prawdziwy przebieg przez
 * powitanie i wejście do środowiska, a potem wykonując przekazanie na
 * dokumencie zbudowanym z prawdziwego `index.html` — mierzy stan węzłów
 * po przejściu, nie samą obecność funkcji rozstrzygającej.
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
import { zbudujDokument, type DokumentZastepczy } from './dom-zastepczy.ts';
import { gotowaDoPrzekazania, wykonajPrzekazanie } from './przekazanie.ts';

/* Odczyt pliku bez typów środowiska: klient nie zaciąga deklaracji Node,
   a specyfikator spoza literału zostawia moduł nieopisanym — ten sam
   obejście co w sprawdzianie katalogu treści. */
const nazwaModulu = 'node:fs';
const pliki = (await import(nazwaModulu)) as unknown as {
  readFileSync(sciezka: string, kodowanie: string): string;
};

/** Prawdziwy dokument wydania — sprawdzian mierzy przejście na jego rzeczywistym kształcie, nie na jego odwzorowaniu. */
const INDEKS_HTML = pliki.readFileSync(
  new URL('../../index.html', import.meta.url).pathname,
  'utf8',
);

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

/** Prowadzi przebieg zastępczy do stanu, w którym środowisko jest znane — punkt startowy sprawdzianów samego przekazania. */
async function stanZeSrodowiskiem() {
  const { przebieg, rdzen } = zaloz();
  rdzen.odpowiadaj(Command.ConnectionHello, POWITANIE_Z_TOKENEM);
  rdzen.odpowiadaj(Command.EnvironmentList, SRODOWISKA);
  rdzen.odpowiadaj(Command.EnvironmentEnter, WEJSCIE);
  przebieg.polacz();
  rdzen.ogloszStan('polaczony');
  await poOdstepie();
  await poOdstepie();
  const stan = przebieg.stan();
  if (!gotowaDoPrzekazania(stan)) throw new Error('stan-testowy-bez-srodowiska');
  return stan;
}

/** Podmienia `document` globalny na dokument zastępczy na czas budowy węzłów przez montaż ramy. */
function jakoGlobalny(dokument: DokumentZastepczy): void {
  (globalThis as unknown as { document: DokumentZastepczy }).document = dokument;
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
    sprawdz(!gotowoscBiezaca(przebieg), 'przejście nastąpiło bez odpowiedzi environment.enter');
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

    sprawdz(gotowoscBiezaca(przebieg), 'przejście nie nastąpiło mimo udanego wejścia do środowiska');
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
    sprawdz(!gotowoscBiezaca(przebieg), 'przejście zapaliło się mimo odmowy rdzenia');
    sprawdz(gotowosci.every((g) => g === false), 'gotowość zapaliła się mimo braku środowiska');
  },

  async 'przekazanie zdejmuje scenę wejścia, montuje ramę i odkrywa jej miejsce'() {
    const stan = await stanZeSrodowiskiem();
    const dokument = zbudujDokument(INDEKS_HTML);
    jakoGlobalny(dokument);
    let zdjeta = 0;

    const udalo = wykonajPrzekazanie(stan, {
      dokument: dokument as unknown as Document,
      zdejmijOknoWejscia: () => {
        zdjeta += 1;
      },
    });

    sprawdz(udalo, 'przekazanie zgłosiło niepowodzenie mimo obecności obu węzłów montażu');
    rowne(zdjeta, 1, 'okno wejścia nie zostało zdjęte dokładnie raz');

    const scenaWejscia = dokument.querySelector('[data-wejscie]');
    const miejsceRamy = dokument.querySelector('[data-rama-aplikacji]');
    sprawdz(scenaWejscia !== null, 'scena wejścia zniknęła z dokumentu zamiast zostać ukryta');
    sprawdz(miejsceRamy !== null, 'miejsce ramy zniknęło z dokumentu');
    sprawdz(scenaWejscia!.hidden, 'scena wejścia nie zeszła z ekranu');
    sprawdz(!miejsceRamy!.hidden, 'miejsce ramy zostało ukryte zamiast odkryte');

    const tytulBelki = miejsceRamy!.querySelector('[data-belka-tytul]');
    sprawdz(tytulBelki !== null, 'rama nie wystawiła tytułu belki');
    rowne(tytulBelki!.textContent, 'Danaco Console › TalkIn', 'tytuł belki nie niesie nazwy środowiska po montażu');

    const licznikSesji = miejsceRamy!.querySelector('[data-stan-sesje]');
    sprawdz(licznikSesji !== null, 'pas stanu nie wystawił licznika sesji');
    rowne(licznikSesji!.textContent, '0', 'licznik sesji nie odpowiada liczbie kart odtworzonych przez rdzeń');

    const przyciskModulu = miejsceRamy!.querySelector('.dn-szyna-poz--modul');
    sprawdz(przyciskModulu !== null, 'szyna nie wystawiła przycisku modułu');
    przyciskModulu!.dispatchEvent({ type: 'click', target: przyciskModulu! });
    rowne(przyciskModulu!.getAttribute('aria-current'), 'true', 'kliknięty moduł nie został oznaczony jako bieżący');
    rowne(tytulBelki!.textContent, 'Danaco Console › Studio', 'kliknięcie modułu nie zmieniło tytułu belki na jego nazwę');
  },

  async 'brak węzła montażu w dokumencie jest odmową nazwaną w dzienniku, nie cichym zaniechaniem'() {
    const stan = await stanZeSrodowiskiem();
    const dokument = zbudujDokument(INDEKS_HTML.replace('<div data-rama-aplikacji hidden></div>', ''));
    jakoGlobalny(dokument);
    let zdjeta = 0;
    const oryginalnyBlad = console.error;
    const zapisane: unknown[][] = [];
    console.error = (...argumenty: unknown[]) => {
      zapisane.push(argumenty);
    };

    let udalo: boolean;
    try {
      udalo = wykonajPrzekazanie(stan, {
        dokument: dokument as unknown as Document,
        zdejmijOknoWejscia: () => {
          zdjeta += 1;
        },
      });
    } finally {
      console.error = oryginalnyBlad;
    }

    sprawdz(!udalo, 'przekazanie zgłosiło powodzenie mimo braku węzła montażu ramy');
    rowne(zdjeta, 0, 'okno wejścia zostało zdjęte mimo nieudanego przekazania');
    rowne(zapisane.length, 1, 'brak węzła montażu nie trafił do dziennika dokładnie raz');
    sprawdz(
      typeof zapisane[0]?.[0] === 'string' && (zapisane[0]![0] as string).startsWith('[rama]'),
      'odmowa w dzienniku nie niesie znacznika warstwy',
    );
  },

  async 'zatrzask nie blokuje na stałe: przekazanie udaje się przy kolejnej zmianie, gdy węzeł montażu się pojawi'() {
    const stan = await stanZeSrodowiskiem();
    let przekazano = false;
    function naZmianeAplikacji(dokument: DokumentZastepczy): boolean {
      if (przekazano) return przekazano;
      jakoGlobalny(dokument);
      przekazano = wykonajPrzekazanie(stan, {
        dokument: dokument as unknown as Document,
        zdejmijOknoWejscia: () => {},
      });
      return przekazano;
    }

    const dokumentBezRamy = zbudujDokument(INDEKS_HTML.replace('<div data-rama-aplikacji hidden></div>', ''));
    sprawdz(!naZmianeAplikacji(dokumentBezRamy), 'przekazanie udało się mimo braku węzła montażu ramy');
    sprawdz(!przekazano, 'zatrzask zapadł mimo nieudanego przekazania');

    const dokumentZRama = zbudujDokument(INDEKS_HTML);
    sprawdz(naZmianeAplikacji(dokumentZRama), 'przekazanie nie udało się przy kolejnej zmianie, choć węzeł montażu już stał w dokumencie');
    sprawdz(przekazano, 'zatrzask nie zapadł mimo udanego przekazania');
    sprawdz(dokumentZRama.querySelector('[data-belka-tytul]') !== null, 'rama nie zamontowała się przy powtórnej próbie');
  },
});

function gotowoscBiezaca(przebieg: ReturnType<typeof utworzPrzebieg>): boolean {
  return gotowaDoPrzekazania(przebieg.stan());
}
