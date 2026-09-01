import { utworzKolejkeWychodzaca, type KolejkaWychodzaca } from './kolejka-wychodzaca.ts';
import { utworzMagistrale, type Magistrala, type Odsubskrybuj } from './magistrala-zdarzen.ts';
import { wykladniczePonawianie, type PolitykaPonawiania } from './ponawianie.ts';
import type { StanPolaczenia } from './stan-polaczenia.ts';

/** Powód, dla którego odłożona ramka nie została rdzeniowi wydana. */
export type PowodPorzucenia =
  /** Ramka czekała na łączność dłużej niż zapora czasu odłożenia. */
  | 'zapora-czasu'
  /** Gniazdo, do którego ramka należała, zostało zerwane. */
  | 'zerwanie';

/** Transport ramek tekstowych do rdzenia, ukrywający przed wołającym stan gniazda i kolejkę wychodzącą. */
export interface Transport {
  /** Rozpoczyna łączenie i utrzymuje je przez ponawianie. */
  polacz(): void;
  /** Wysyła ramkę; przy braku połączenia albo przy wstrzymaniu trafia ona do kolejki wychodzącej. */
  wyslij(ramka: string): void;
  /** Wysyła ramkę powitania, która wstrzymania nie podlega — to ono je zdejmuje. */
  wyslijPowitanie(ramka: string): void;
  /** Zdejmuje wstrzymanie kolejki i wydaje ją do gniazda. */
  zwolnijWstrzymanie(): void;
  /** Subskrypcja ramek przychodzących. */
  naRamke(sluchacz: (ramka: string) => void): Odsubskrybuj;
  /** Subskrypcja ramek porzuconych, którym nie ma już czego doręczyć. */
  naPorzucona(sluchacz: (ramka: string, powod: PowodPorzucenia) => void): Odsubskrybuj;
  /** Subskrypcja zmian stanu połączenia. */
  naStan(sluchacz: (stan: StanPolaczenia) => void): Odsubskrybuj;
  /** Bieżący stan połączenia. */
  stan(): StanPolaczenia;
  /** Liczba ramek oczekujących w kolejce wychodzącej. */
  oczekujace(): number;
  /** Zamyka połączenie i wstrzymuje ponawianie. */
  rozlacz(): void;
}

/*
Zapora czasu odłożenia. Ramka wydana rdzeniowi po tym czasie zamawia pracę,
której odpowiedzi nikt już nie czeka — wołający dostał odmowę terminu — a przy
komendzie zmieniającej stan zamawia ją powtórnie. Zapora stoi poniżej terminu
odpowiedzi korelacji, żeby odmowa pochodziła z porzucenia, nie z ciszy.
*/
const ZAPORA_ODLOZENIA_MS = 20_000;

/** Ramka odłożona na czas rozłączenia wraz z chwilą nadania, po której liczy się zapora czasu. */
interface RamkaOdlozona {
  tresc: string;
  nadana: number;
}

/** Ramka, która rdzeniowi wydana nie będzie, wraz z powodem porzucenia. */
interface RamkaPorzucona {
  tresc: string;
  powod: PowodPorzucenia;
}

/** Połączenie WebSocket z ponawianiem i kolejkowaniem ramek, utrzymujące łączność z rdzeniem bez udziału wołającego. */
class Gniazdo implements Transport {
  private readonly ramki: Magistrala<string> = utworzMagistrale<string>();
  private readonly porzucone: Magistrala<RamkaPorzucona> = utworzMagistrale<RamkaPorzucona>();
  private readonly stany: Magistrala<StanPolaczenia> = utworzMagistrale<StanPolaczenia>();
  private readonly kolejka: KolejkaWychodzaca<RamkaOdlozona> =
    utworzKolejkeWychodzaca<RamkaOdlozona>();
  private readonly kolejkaPowitania: KolejkaWychodzaca<RamkaOdlozona> =
    utworzKolejkeWychodzaca<RamkaOdlozona>();
  private gniazdo: WebSocket | null = null;
  private biezacyStan: StanPolaczenia = 'rozlaczony';
  private numerProby = 0;
  private zaplanowane: ReturnType<typeof setTimeout> | null = null;
  private zaniechane = false;
  /* Rdzeń wiąże sesję bramki z gniazdem dopiero w powitaniu, więc komenda
     wydana przed jego odpowiedzią wraca odmową `not_authenticated`. Każde
     gniazdo zaczyna więc wstrzymane i czeka na powitanie własne. */
  private wstrzymana = true;

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
    // Zamknięcie na błędzie dotyczy gniazda już otwartego; ponowne close wywołuje nawrót bez końca.
    gniazdo.addEventListener('error', () => {
      if (gniazdo.readyState === WebSocket.OPEN) gniazdo.close();
    });
  }

  wyslij(ramka: string): void {
    this.wydajAlboOdloz(ramka, this.kolejka, !this.wstrzymana);
  }

  wyslijPowitanie(ramka: string): void {
    this.wydajAlboOdloz(ramka, this.kolejkaPowitania, true);
  }

  zwolnijWstrzymanie(): void {
    this.wstrzymana = false;
    this.wydajKolejke(this.kolejka);
  }

  naRamke(sluchacz: (ramka: string) => void): Odsubskrybuj {
    return this.ramki.subskrybuj(sluchacz);
  }

  naPorzucona(sluchacz: (ramka: string, powod: PowodPorzucenia) => void): Odsubskrybuj {
    return this.porzucone.subskrybuj((porzucona) => sluchacz(porzucona.tresc, porzucona.powod));
  }

  naStan(sluchacz: (stan: StanPolaczenia) => void): Odsubskrybuj {
    sluchacz(this.biezacyStan);
    return this.stany.subskrybuj(sluchacz);
  }

  stan(): StanPolaczenia {
    return this.biezacyStan;
  }

  oczekujace(): number {
    return this.kolejka.rozmiar() + this.kolejkaPowitania.rozmiar();
  }

  /** Zaniechanie kończy ponawianie bezterminowe; ramki odłożone zostają w kolejce, gotowe do wysłania. */
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
    this.wstrzymana = true;
    this.zapiszStan('polaczony');
    this.wydajKolejke(this.kolejkaPowitania);
  }

  private obsluzRamke(zdarzenie: MessageEvent<unknown>): void {
    if (typeof zdarzenie.data === 'string') {
      this.ramki.oglos(zdarzenie.data);
    }
  }

  private obsluzZamkniecie(): void {
    if (this.zaniechane) return;
    this.gniazdo = null;
    /* Powitanie należy do gniazda, które je przyjęło: rdzeń wiąże po nim sesję
       bramki z konkretnym połączeniem, więc na nowym gnieździe jest bezużyteczne. */
    this.porzucKolejke(this.kolejkaPowitania, 'zerwanie');
    this.zapiszStan('ponawianie');
    this.zaplanujPonowienie();
  }

  /** Wydaje ramkę do otwartego gniazda albo odkłada ją w podanej kolejce, gdy gniazdo jest zamknięte lub kolejka wstrzymana. */
  private wydajAlboOdloz(
    ramka: string,
    kolejka: KolejkaWychodzaca<RamkaOdlozona>,
    wolno: boolean,
  ): void {
    if (wolno && this.gniazdo !== null && this.gniazdo.readyState === WebSocket.OPEN) {
      this.gniazdo.send(ramka);
      return;
    }
    kolejka.dodaj({ tresc: ramka, nadana: Date.now() });
  }

  /** Wydaje kolejkę do otwartego gniazda; ramka po zaporze czasu jest porzucana, a niewysłana wraca do kolejki. */
  private wydajKolejke(kolejka: KolejkaWychodzaca<RamkaOdlozona>): void {
    const teraz = Date.now();
    for (const odlozona of kolejka.wydajWszystko()) {
      if (teraz - odlozona.nadana >= ZAPORA_ODLOZENIA_MS) {
        this.porzucone.oglos({ tresc: odlozona.tresc, powod: 'zapora-czasu' });
        continue;
      }
      if (this.gniazdo !== null && this.gniazdo.readyState === WebSocket.OPEN) {
        this.gniazdo.send(odlozona.tresc);
        continue;
      }
      kolejka.dodaj(odlozona);
    }
  }

  /** Opróżnia kolejkę, ogłaszając każdą jej ramkę jako porzuconą z podanego powodu. */
  private porzucKolejke(kolejka: KolejkaWychodzaca<RamkaOdlozona>, powod: PowodPorzucenia): void {
    for (const odlozona of kolejka.wydajWszystko()) {
      this.porzucone.oglos({ tresc: odlozona.tresc, powod });
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
