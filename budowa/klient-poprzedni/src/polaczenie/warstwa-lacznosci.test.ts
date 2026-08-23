import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { utworzTransport, type Transport } from './gniazdo';
import { utworzKolejkeWychodzaca } from './kolejka-wychodzaca';
import { utworzMagistrale } from './magistrala-zdarzen';
import { wykladniczePonawianie } from './ponawianie';
import type { StanPolaczenia } from './stan-polaczenia';

/**
 * Warstwa łączności klienta.
 *
 * Jej obietnica wobec Operatora jest jedna i mocna: okno komunikacji pozostaje
 * użyteczne w każdym stanie, a wiadomość wpisana przy rozłączeniu nie ginie.
 * Sprawdziany niżej mierzą dokładnie to — kolejkowanie przy braku rdzenia,
 * opróżnienie kolejki po powrocie oraz ponawianie, które nie ustaje.
 */

describe('kolejka wychodząca', () => {
  it('wydaje elementy w kolejności dodania i zostawia się pusta', () => {
    const kolejka = utworzKolejkeWychodzaca<string>();
    kolejka.dodaj('pierwsza');
    kolejka.dodaj('druga');
    kolejka.dodaj('trzecia');

    expect(kolejka.rozmiar()).toBe(3);
    expect(kolejka.wydajWszystko()).toEqual(['pierwsza', 'druga', 'trzecia']);
    expect(kolejka.rozmiar()).toBe(0);
    expect(kolejka.wydajWszystko()).toEqual([]);
  });

  it('nie oddaje odbiorcy wglądu we własne wnętrze', () => {
    const kolejka = utworzKolejkeWychodzaca<string>();
    kolejka.dodaj('pierwsza');

    const wydane = kolejka.wydajWszystko();
    wydane.push('dopisana z zewnątrz');
    kolejka.dodaj('druga');

    expect(kolejka.wydajWszystko()).toEqual(['druga']);
  });
});

describe('magistrala zdarzeń', () => {
  it('powiadamia słuchaczy i zdejmuje ich po odsubskrybowaniu', () => {
    const magistrala = utworzMagistrale<number>();
    const odebrane: number[] = [];

    const odsubskrybuj = magistrala.subskrybuj((liczba) => odebrane.push(liczba));
    expect(magistrala.liczbaSluchaczy()).toBe(1);

    magistrala.oglos(1);
    odsubskrybuj();
    magistrala.oglos(2);

    expect(odebrane).toEqual([1]);
    expect(magistrala.liczbaSluchaczy()).toBe(0);
  });

  it('nie przerywa rozgłaszania na błędzie jednego słuchacza', () => {
    const magistrala = utworzMagistrale<string>();
    const podglad = vi.spyOn(console, 'error').mockImplementation(() => {});
    const odebrane: string[] = [];

    magistrala.subskrybuj(() => {
      throw new Error('słuchacz sprawdzianu');
    });
    magistrala.subskrybuj((dane) => odebrane.push(dane));

    magistrala.oglos('zdarzenie');

    expect(odebrane, 'słuchacz za wywracającym się nie dostał zdarzenia').toEqual(['zdarzenie']);
    expect(podglad).toHaveBeenCalled();
  });

  it('znosi odsubskrybowanie w trakcie rozgłaszania', () => {
    const magistrala = utworzMagistrale<string>();
    const odebrane: string[] = [];

    const odsubskrybuj = magistrala.subskrybuj(() => odsubskrybuj());
    magistrala.subskrybuj((dane) => odebrane.push(dane));

    expect(() => magistrala.oglos('zdarzenie')).not.toThrow();
    expect(odebrane).toEqual(['zdarzenie']);
  });
});

describe('polityka ponawiania', () => {
  it('rośnie wykładniczo i staje na pułapie', () => {
    const polityka = wykladniczePonawianie(500, 15_000);

    // Rozproszenie sięga 25%, więc porównywane są widełki, nie wartość.
    const widelki = (podstawa: number): [number, number] => [podstawa, Math.round(podstawa * 1.25)];

    for (const [proba, podstawa] of [
      [1, 500],
      [2, 1000],
      [3, 2000],
      [4, 4000],
      [5, 8000],
      [6, 15_000],
      [20, 15_000],
    ] as const) {
      const [dol, gora] = widelki(podstawa);
      const opoznienie = polityka.opoznienie(proba);
      expect(opoznienie, `próba ${proba}`).toBeGreaterThanOrEqual(dol);
      expect(opoznienie, `próba ${proba}`).toBeLessThanOrEqual(gora);
    }
  });

  it('traktuje numer próby poniżej jedynki jak pierwszą próbę', () => {
    const polityka = wykladniczePonawianie(500, 15_000);
    for (const proba of [0, -1]) {
      const opoznienie = polityka.opoznienie(proba);
      expect(opoznienie).toBeGreaterThanOrEqual(500);
      expect(opoznienie).toBeLessThanOrEqual(625);
    }
  });
});

/**
 * Gniazdo sprawdzianu zastępuje WebSocket przeglądarki. Zastępstwo, a nie
 * przeglądarka bez okna: sprawdzane jest zachowanie klienta wobec gniazda —
 * kolejkowanie, ponawianie, kolejność stanów — a nie sama biblioteka gniazda.
 */
class GniazdoSprawdzianu {
  static readonly OPEN = 1;
  static otwarte: GniazdoSprawdzianu[] = [];

  readyState = 0;
  wyslane: string[] = [];
  private sluchacze = new Map<string, ((zdarzenie: unknown) => void)[]>();

  constructor(readonly adres: string) {
    GniazdoSprawdzianu.otwarte.push(this);
  }

  addEventListener(nazwa: string, sluchacz: (zdarzenie: unknown) => void): void {
    const wpisani = this.sluchacze.get(nazwa) ?? [];
    wpisani.push(sluchacz);
    this.sluchacze.set(nazwa, wpisani);
  }

  send(ramka: string): void {
    this.wyslane.push(ramka);
  }

  close(): void {
    this.readyState = 3;
    this.ogloś('close', {});
  }

  /** Doprowadza gniazdo do stanu otwartego, tak jak zrobiłaby to przeglądarka. */
  otworz(): void {
    this.readyState = GniazdoSprawdzianu.OPEN;
    this.ogloś('open', {});
  }

  /** Podaje ramkę przychodzącą. */
  przyjmij(dane: unknown): void {
    this.ogloś('message', { data: dane });
  }

  private ogloś(nazwa: string, zdarzenie: unknown): void {
    for (const sluchacz of this.sluchacze.get(nazwa) ?? []) sluchacz(zdarzenie);
  }
}

describe('transport gniazda', () => {
  let pierwotneGniazdo: unknown;

  beforeEach(() => {
    vi.useFakeTimers();
    GniazdoSprawdzianu.otwarte = [];
    pierwotneGniazdo = (globalThis as Record<string, unknown>).WebSocket;
    (globalThis as Record<string, unknown>).WebSocket = GniazdoSprawdzianu;
  });

  afterEach(() => {
    (globalThis as Record<string, unknown>).WebSocket = pierwotneGniazdo;
    vi.useRealTimers();
  });

  /** Zakłada transport ze stałym odstępem ponowienia — bez losowości. */
  function zalozTransport(): { transport: Transport; stany: StanPolaczenia[] } {
    const stany: StanPolaczenia[] = [];
    const transport = utworzTransport('ws://127.0.0.1:17870/ws', { opoznienie: () => 1000 });
    transport.naStan((stan) => stany.push(stan));
    return { transport, stany };
  }

  it('kolejkuje ramkę wpisaną przed połączeniem i wydaje ją po otwarciu', () => {
    const { transport } = zalozTransport();

    transport.wyslij('ramka sprzed połączenia');
    expect(transport.oczekujace(), 'ramka zgubiona zamiast odłożona').toBe(1);

    transport.polacz();
    const gniazdo = GniazdoSprawdzianu.otwarte[0];
    expect(gniazdo.wyslane, 'ramka poszła do gniazda niegotowego').toEqual([]);

    gniazdo.otworz();
    expect(transport.oczekujace()).toBe(0);
    expect(gniazdo.wyslane).toEqual(['ramka sprzed połączenia']);
  });

  it('odkłada ramkę wpisaną po zerwaniu i wydaje ją po wznowieniu', () => {
    const { transport } = zalozTransport();
    transport.polacz();
    GniazdoSprawdzianu.otwarte[0].otworz();

    GniazdoSprawdzianu.otwarte[0].close();
    transport.wyslij('ramka z czasu rozłączenia');
    expect(transport.oczekujace()).toBe(1);

    vi.advanceTimersByTime(1000);
    const wznowione = GniazdoSprawdzianu.otwarte[1];
    expect(wznowione, 'transport nie podjął ponowienia').toBeDefined();
    wznowione.otworz();

    expect(transport.oczekujace()).toBe(0);
    expect(wznowione.wyslane).toEqual(['ramka z czasu rozłączenia']);
  });

  it('przechodzi stany w kolejności łączenie → połączony → ponawianie → połączony', () => {
    const { transport, stany } = zalozTransport();

    transport.polacz();
    GniazdoSprawdzianu.otwarte[0].otworz();
    GniazdoSprawdzianu.otwarte[0].close();
    vi.advanceTimersByTime(1000);
    GniazdoSprawdzianu.otwarte[1].otworz();

    expect(stany).toEqual([
      'rozlaczony',
      'laczenie',
      'polaczony',
      'ponawianie',
      'polaczony',
    ]);
  });

  it('ponawia bez końca, dopóki rdzeń milczy', () => {
    const { transport } = zalozTransport();
    transport.polacz();

    for (let proba = 1; proba <= 5; proba += 1) {
      GniazdoSprawdzianu.otwarte[proba - 1].close();
      vi.advanceTimersByTime(1000);
    }

    expect(GniazdoSprawdzianu.otwarte, 'ponawianie ustało').toHaveLength(6);
    expect(transport.stan()).toBe('ponawianie');
  });

  it('nie zakłada drugiego gniazda przy powtórnym wywołaniu połączenia', () => {
    const { transport } = zalozTransport();
    transport.polacz();
    transport.polacz();
    expect(GniazdoSprawdzianu.otwarte).toHaveLength(1);
  });

  it('rozgłasza ramki tekstowe i pomija ramki innej maści', () => {
    const { transport } = zalozTransport();
    const odebrane: string[] = [];
    transport.naRamke((ramka) => odebrane.push(ramka));

    transport.polacz();
    const gniazdo = GniazdoSprawdzianu.otwarte[0];
    gniazdo.otworz();
    gniazdo.przyjmij('{"type":"session.changed"}');
    gniazdo.przyjmij(new ArrayBuffer(4));

    expect(odebrane).toEqual(['{"type":"session.changed"}']);
  });

  it('podaje nowemu słuchaczowi stan bieżący, zanim cokolwiek się zmieni', () => {
    const { transport } = zalozTransport();
    transport.polacz();
    GniazdoSprawdzianu.otwarte[0].otworz();

    const odebrane: StanPolaczenia[] = [];
    transport.naStan((stan) => odebrane.push(stan));

    expect(odebrane).toEqual(['polaczony']);
  });
});
