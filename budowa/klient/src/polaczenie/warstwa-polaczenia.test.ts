import { GniazdoZastepcze, podstawGniazdo } from '../gniazdo-zastepcze.ts';
import { bieg, poOdstepie, rowne, sprawdz } from '../sprawdzian.ts';
import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from './adres-rdzenia.ts';
import { utworzKolejkeWychodzaca } from './kolejka-wychodzaca.ts';
import { utworzMagistrale } from './magistrala-zdarzen.ts';
import { wykladniczePonawianie } from './ponawianie.ts';
import { utworzTransport, type Transport } from './gniazdo.ts';
import type { StanPolaczenia } from './stan-polaczenia.ts';

// Warstwa gwarantuje: ramka nie ginie, ponawianie nie ustaje, subskrypcje przeżywają wymianę gniazda.

/** Zakłada transport bez odstępu ponowienia, bez losowości w odmierzaniu czasu i bez rzeczywistego czekania. */
function zalozTransport(): { transport: Transport; stany: StanPolaczenia[] } {
  const stany: StanPolaczenia[] = [];
  const transport = utworzTransport('ws://127.0.0.1:17870/ws', { opoznienie: () => 0 });
  transport.naStan((stan) => stany.push(stan));
  return { transport, stany };
}

await bieg('warstwa połączenia', {
  'kolejka wydaje elementy w kolejności dodania i zostawia się pusta'() {
    const kolejka = utworzKolejkeWychodzaca<string>();
    kolejka.dodaj('pierwsza');
    kolejka.dodaj('druga');
    kolejka.dodaj('trzecia');

    sprawdz(kolejka.rozmiar() === 3, 'kolejka nie zliczyła dołożonych elementów');
    rowne(kolejka.wydajWszystko(), ['pierwsza', 'druga', 'trzecia'], 'kolejność wydania');
    sprawdz(kolejka.rozmiar() === 0, 'kolejka nie opróżniła się po wydaniu');
    rowne(kolejka.wydajWszystko(), [], 'powtórne wydanie z pustej kolejki');
  },

  'kolejka nie oddaje odbiorcy wglądu we własne wnętrze'() {
    const kolejka = utworzKolejkeWychodzaca<string>();
    kolejka.dodaj('pierwsza');

    const wydane = kolejka.wydajWszystko();
    wydane.push('dopisana z zewnątrz');
    kolejka.dodaj('druga');

    rowne(kolejka.wydajWszystko(), ['druga'], 'zawartość kolejki po dopisaniu z zewnątrz');
  },

  'magistrala powiadamia słuchaczy i zdejmuje ich po odsubskrybowaniu'() {
    const magistrala = utworzMagistrale<number>();
    const odebrane: number[] = [];

    const odsubskrybuj = magistrala.subskrybuj((liczba) => odebrane.push(liczba));
    sprawdz(magistrala.liczbaSluchaczy() === 1, 'słuchacz nie został zarejestrowany');

    magistrala.oglos(1);
    odsubskrybuj();
    magistrala.oglos(2);

    rowne(odebrane, [1], 'odebrane po odsubskrybowaniu');
    sprawdz(magistrala.liczbaSluchaczy() === 0, 'słuchacz został po odsubskrybowaniu');
  },

  'magistrala nie przerywa rozgłaszania na błędzie jednego słuchacza'() {
    const magistrala = utworzMagistrale<string>();
    const odebrane: string[] = [];

    magistrala.subskrybuj(() => {
      throw new Error('słuchacz sprawdzianu');
    });
    magistrala.subskrybuj((dane) => odebrane.push(dane));

    magistrala.oglos('zdarzenie');

    rowne(odebrane, ['zdarzenie'], 'słuchacz za wywracającym się nie dostał zdarzenia');
  },

  'magistrala znosi odsubskrybowanie w trakcie rozgłaszania'() {
    const magistrala = utworzMagistrale<string>();
    const odebrane: string[] = [];

    const odsubskrybuj = magistrala.subskrybuj(() => odsubskrybuj());
    magistrala.subskrybuj((dane) => odebrane.push(dane));

    magistrala.oglos('zdarzenie');

    rowne(odebrane, ['zdarzenie'], 'rozgłaszanie przerwane przez odsubskrybowanie w trakcie');
  },

  'ponawianie rośnie wykładniczo i staje na pułapie'() {
    const polityka = wykladniczePonawianie(500, 15_000);
    // Rozproszenie sięga 25%, więc porównywane są widełki, nie wartość.
    for (const [proba, podstawa] of [
      [1, 500],
      [2, 1000],
      [3, 2000],
      [4, 4000],
      [5, 8000],
      [6, 15_000],
      [20, 15_000],
    ] as const) {
      const opoznienie = polityka.opoznienie(proba);
      sprawdz(
        opoznienie >= podstawa && opoznienie <= Math.round(podstawa * 1.25),
        `próba ${proba}: odstęp ${opoznienie} poza widełkami ${podstawa}–${Math.round(podstawa * 1.25)}`,
      );
    }
  },

  'ponawianie traktuje numer próby poniżej jedynki jak pierwszą próbę'() {
    const polityka = wykladniczePonawianie(500, 15_000);
    for (const proba of [0, -1]) {
      const opoznienie = polityka.opoznienie(proba);
      sprawdz(opoznienie >= 500 && opoznienie <= 625, `próba ${proba}: odstęp ${opoznienie}`);
    }
  },

  'adres gniazda powstaje z adresu HTTP rdzenia'() {
    rowne(adresGniazdaRdzenia('http://127.0.0.1:17870'), 'ws://127.0.0.1:17870/ws', 'http');
    rowne(adresGniazdaRdzenia('https://rdzen.example:8443'), 'wss://rdzen.example:8443/ws', 'https');
    rowne(adresGniazdaRdzenia('ftp://127.0.0.1'), null, 'schemat spoza HTTP');
    rowne(adresGniazdaRdzenia('to nie jest adres'), null, 'napis nierozkładalny');
    rowne(adresRdzeniaLokalnego(), 'ws://127.0.0.1:17870/ws', 'adres rdzenia lokalnego');
  },

  'transport kolejkuje ramkę wpisaną przed połączeniem i wydaje ją po otwarciu'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();

      transport.wyslij('ramka sprzed połączenia');
      sprawdz(transport.oczekujace() === 1, 'ramka zgubiona zamiast odłożona');

      transport.polacz();
      const gniazdo = GniazdoZastepcze.otwarte[0];
      sprawdz(gniazdo !== undefined, 'transport nie założył gniazda');
      rowne(gniazdo?.wyslane, [], 'ramka poszła do gniazda niegotowego');

      gniazdo?.otworz();
      sprawdz(transport.oczekujace() === 0, 'kolejka nie opróżniła się po otwarciu');
      rowne(gniazdo?.wyslane, ['ramka sprzed połączenia'], 'ramka wydana po otwarciu');
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  async 'transport odkłada ramkę wpisaną po zerwaniu i wydaje ją po wznowieniu'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();
      transport.polacz();
      GniazdoZastepcze.otwarte[0]?.otworz();

      GniazdoZastepcze.otwarte[0]?.close();
      transport.wyslij('ramka z czasu rozłączenia');
      sprawdz(transport.oczekujace() === 1, 'ramka z czasu rozłączenia zgubiona');

      await poOdstepie();
      const wznowione = GniazdoZastepcze.otwarte[1];
      sprawdz(wznowione !== undefined, 'transport nie podjął ponowienia');
      wznowione?.otworz();

      sprawdz(transport.oczekujace() === 0, 'kolejka nie opróżniła się po wznowieniu');
      rowne(wznowione?.wyslane, ['ramka z czasu rozłączenia'], 'ramka wydana po wznowieniu');
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  async 'transport przechodzi stany łączenie → połączony → ponawianie → połączony'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport, stany } = zalozTransport();

      transport.polacz();
      GniazdoZastepcze.otwarte[0]?.otworz();
      GniazdoZastepcze.otwarte[0]?.close();
      await poOdstepie();
      GniazdoZastepcze.otwarte[1]?.otworz();

      rowne(
        stany,
        ['rozlaczony', 'laczenie', 'polaczony', 'ponawianie', 'polaczony'],
        'kolejność stanów',
      );
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  async 'transport ponawia bez końca, dopóki rdzeń milczy'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();
      transport.polacz();

      for (let proba = 1; proba <= 5; proba += 1) {
        GniazdoZastepcze.otwarte[proba - 1]?.close();
        await poOdstepie();
      }

      sprawdz(GniazdoZastepcze.otwarte.length === 6, 'ponawianie ustało');
      sprawdz(transport.stan() === 'ponawianie', `stan po ponowieniach: ${transport.stan()}`);
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  'transport nie zakłada drugiego gniazda przy powtórnym wywołaniu połączenia'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();
      transport.polacz();
      transport.polacz();
      sprawdz(GniazdoZastepcze.otwarte.length === 1, 'powstało drugie gniazdo');
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  'transport rozgłasza ramki tekstowe i pomija ramki innej maści'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();
      const odebrane: string[] = [];
      transport.naRamke((ramka) => odebrane.push(ramka));

      transport.polacz();
      const gniazdo = GniazdoZastepcze.otwarte[0];
      gniazdo?.otworz();
      gniazdo?.przyjmij('{"type":"session.changed"}');
      gniazdo?.przyjmij(new ArrayBuffer(4));

      rowne(odebrane, ['{"type":"session.changed"}'], 'odebrane ramki');
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  'transport podaje nowemu słuchaczowi stan bieżący, zanim cokolwiek się zmieni'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();
      transport.polacz();
      GniazdoZastepcze.otwarte[0]?.otworz();

      const odebrane: StanPolaczenia[] = [];
      transport.naStan((stan) => odebrane.push(stan));

      rowne(odebrane, ['polaczony'], 'stan podany nowemu słuchaczowi');
      transport.rozlacz();
    } finally {
      przywroc();
    }
  },

  async 'zaniechanie połączenia wstrzymuje ponawianie'() {
    const przywroc = podstawGniazdo();
    try {
      const { transport } = zalozTransport();
      transport.polacz();
      GniazdoZastepcze.otwarte[0]?.otworz();

      transport.rozlacz();
      await poOdstepie();

      sprawdz(GniazdoZastepcze.otwarte.length === 1, 'ponawianie ruszyło mimo zaniechania');
      sprawdz(transport.stan() === 'rozlaczony', `stan po zaniechaniu: ${transport.stan()}`);
    } finally {
      przywroc();
    }
  },
});
