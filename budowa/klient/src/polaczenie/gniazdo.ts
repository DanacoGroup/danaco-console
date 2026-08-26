import { utworzKolejkeWychodzaca, type KolejkaWychodzaca } from './kolejka-wychodzaca.ts';
import { utworzMagistrale, type Magistrala, type Odsubskrybuj } from './magistrala-zdarzen.ts';
import { wykladniczePonawianie, type PolitykaPonawiania } from './ponawianie.ts';
import type { StanPolaczenia } from './stan-polaczenia.ts';

/** Transport ramek tekstowych do rdzenia. */
export interface Transport {
  /** Rozpoczyna łączenie i utrzymuje je przez ponawianie. */
  polacz(): void;
  /** Wysyła ramkę; przy braku połączenia trafia ona do kolejki wychodzącej. */
  wyslij(ramka: string): void;
  /** Subskrypcja ramek przychodzących. */
  naRamke(sluchacz: (ramka: string) => void): Odsubskrybuj;
  /** Subskrypcja zmian stanu połączenia. */
  naStan(sluchacz: (stan: StanPolaczenia) => void): Odsubskrybuj;
  /** Bieżący stan połączenia. */
  stan(): StanPolaczenia;
  /** Liczba ramek oczekujących w kolejce wychodzącej. */
  oczekujace(): number;
  /** Zamyka połączenie i wstrzymuje ponawianie. */
  rozlacz(): void;
}

/** Połączenie WebSocket z ponawianiem i kolejkowaniem ramek. */
class Gniazdo implements Transport {
  private readonly ramki: Magistrala<string> = utworzMagistrale<string>();
  private readonly stany: Magistrala<StanPolaczenia> = utworzMagistrale<StanPolaczenia>();
  private readonly kolejka: KolejkaWychodzaca<string> = utworzKolejkeWychodzaca<string>();
  private gniazdo: WebSocket | null = null;
  private biezacyStan: StanPolaczenia = 'rozlaczony';
  private numerProby = 0;
  private zaplanowane: ReturnType<typeof setTimeout> | null = null;
  private zaniechane = false;

  private readonly adres: string;
  private readonly polityka: PolitykaPonawiania;

  constructor(adres: string, polityka: PolitykaPonawiania) {
    this.adres = adres;
    this.polityka = polityka;
  }

  polacz(): void {
    if (this.gniazdo !== null) return;
    this.zaniechane = false;
    this.anulujPlan();
    this.zapiszStan(this.numerProby === 0 ? 'laczenie' : 'ponawianie');
    const gniazdo = new WebSocket(this.adres);
    this.gniazdo = gniazdo;
    gniazdo.addEventListener('open', () => this.obsluzOtwarcie());
    gniazdo.addEventListener('message', (zdarzenie) => this.obsluzRamke(zdarzenie));
    gniazdo.addEventListener('close', () => this.obsluzZamkniecie());
    // Zamknięcie na błędzie dotyczy wyłącznie gniazda już otwartego. Gniazdo,
    // które błędem kończy samo nawiązywanie, zamyka się bez niczyjej pomocy
    // i ogłasza to zdarzeniem `close`; wywołanie `close` na takim gnieździe
    // wywołuje kolejny błąd i wpada w nawrót bez końca.
    gniazdo.addEventListener('error', () => {
      if (gniazdo.readyState === WebSocket.OPEN) gniazdo.close();
    });
  }

  wyslij(ramka: string): void {
    if (this.gniazdo !== null && this.gniazdo.readyState === WebSocket.OPEN) {
      this.gniazdo.send(ramka);
      return;
    }
    this.kolejka.dodaj(ramka);
  }

  naRamke(sluchacz: (ramka: string) => void): Odsubskrybuj {
    return this.ramki.subskrybuj(sluchacz);
  }

  naStan(sluchacz: (stan: StanPolaczenia) => void): Odsubskrybuj {
    sluchacz(this.biezacyStan);
    return this.stany.subskrybuj(sluchacz);
  }

  stan(): StanPolaczenia {
    return this.biezacyStan;
  }

  oczekujace(): number {
    return this.kolejka.rozmiar();
  }

  /**
   * Zaniechanie połączenia na żądanie wołającego.
   *
   * Ponawianie jest bezterminowe, więc bez tej drogi proces klienta nie miałby
   * jak dojść do końca: zaplanowana próba trzymałaby go przy życiu, a każde
   * zamknięcie gniazda planowałoby następną. Ramki odłożone zostają w kolejce
   * — zaniechanie dotyczy połączenia, nie treści, która czeka na wysłanie.
   */
  rozlacz(): void {
    this.zaniechane = true;
    this.anulujPlan();
    const gniazdo = this.gniazdo;
    this.gniazdo = null;
    this.numerProby = 0;
    gniazdo?.close();
    this.zapiszStan('rozlaczony');
  }

  private obsluzOtwarcie(): void {
    this.numerProby = 0;
    this.zapiszStan('polaczony');
    this.oproznijKolejke();
  }

  private obsluzRamke(zdarzenie: MessageEvent<unknown>): void {
    if (typeof zdarzenie.data === 'string') {
      this.ramki.oglos(zdarzenie.data);
    }
  }

  private obsluzZamkniecie(): void {
    if (this.zaniechane) return;
    this.gniazdo = null;
    this.zapiszStan('ponawianie');
    this.zaplanujPonowienie();
  }

  /** Wydaje kolejkę do otwartego gniazda; ramka niewysłana wraca do kolejki. */
  private oproznijKolejke(): void {
    for (const ramka of this.kolejka.wydajWszystko()) {
      this.wyslij(ramka);
    }
  }

  private zaplanujPonowienie(): void {
    this.numerProby += 1;
    this.anulujPlan();
    this.zaplanowane = setTimeout(() => {
      this.zaplanowane = null;
      this.polacz();
    }, this.polityka.opoznienie(this.numerProby));
  }

  private anulujPlan(): void {
    if (this.zaplanowane !== null) {
      clearTimeout(this.zaplanowane);
      this.zaplanowane = null;
    }
  }

  private zapiszStan(stan: StanPolaczenia): void {
    if (this.biezacyStan === stan) return;
    this.biezacyStan = stan;
    this.stany.oglos(stan);
  }
}

export function utworzTransport(
  adres: string,
  polityka: PolitykaPonawiania = wykladniczePonawianie(),
): Transport {
  return new Gniazdo(adres, polityka);
}
