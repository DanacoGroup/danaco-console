// Transport ramek do rdzenia: gniazdo WebSocket z ponawianiem, kolejką
// wychodzącą i rozpoznaniem uśpienia maszyny.
import { adresNawiazania, sekretNawiazaniaOdPowloki } from './adres-rdzenia.ts';
import { utworzKolejkeWychodzaca, type KolejkaWychodzaca } from './kolejka-wychodzaca.ts';
import { utworzMagistrale, type Magistrala, type Odsubskrybuj } from './magistrala-zdarzen.ts';
import { wykladniczePonawianie, type PolitykaPonawiania } from './ponawianie.ts';
import type { StanPolaczenia } from './stan-polaczenia.ts';

export type PowodPorzucenia =
  | 'zapora-czasu'
  | 'zerwanie'
  | 'przepelnienie';

export interface Transport {
  polacz(): void;
  wznow(): void;
  wyslij(ramka: string): void;
  wyslijPowitanie(ramka: string): void;
  zwolnijWstrzymanie(): void;
  naRamke(sluchacz: (ramka: string) => void): Odsubskrybuj;
  naPorzucona(sluchacz: (ramka: string, powod: PowodPorzucenia) => void): Odsubskrybuj;
  naStan(sluchacz: (stan: StanPolaczenia) => void): Odsubskrybuj;
  stan(): StanPolaczenia;
  oczekujace(): number;
  rozlacz(): void;
}

// Zapora stoi poniżej terminu korelacji: odmowa ma pochodzić z porzucenia ramki.
const ZAPORA_ODLOZENIA_MS = 20_000;

interface RamkaOdlozona {
  tresc: string;
  nadana: number;
}

// Rdzeń trzyma 256 ramek na gniazdo (`transport/ustawienia.go`,
// `pojemnoscKolejkiDomyslna`); ponad sufit najstarsza ramka wraca porzucona.
const SUFIT_KOLEJKI = 256;

interface RamkaPorzucona {
  tresc: string;
  powod: PowodPorzucenia;
}

// Rdzeń pinguje co 20 s i zamyka gniazdo po 10 s ciszy
// (`transport/petla_odbioru.go`); pingów przeglądarka skryptowi nie pokazuje.
const ODSTEP_PULSU_MS = 5_000;
const PROG_USPIENIA_MS = 30_000;

class Gniazdo implements Transport {
  private readonly ramki: Magistrala<string> = utworzMagistrale<string>();
  private readonly porzucone: Magistrala<RamkaPorzucona> = utworzMagistrale<RamkaPorzucona>();
  private readonly stany: Magistrala<StanPolaczenia> = utworzMagistrale<StanPolaczenia>();
  private readonly kolejka: KolejkaWychodzaca<RamkaOdlozona> =
    utworzKolejkeWychodzaca<RamkaOdlozona>(SUFIT_KOLEJKI);
  private readonly kolejkaPowitania: KolejkaWychodzaca<RamkaOdlozona> =
    utworzKolejkeWychodzaca<RamkaOdlozona>();
  private gniazdo: WebSocket | null = null;
  private biezacyStan: StanPolaczenia = 'rozlaczony';
  private numerProby = 0;
  private zaplanowane: ReturnType<typeof setTimeout> | null = null;
  private zaniechane = false;
  // Komenda przed powitaniem wraca odmową `not_authenticated`.
  private wstrzymana = true;
  private licznikPulsu: ReturnType<typeof setInterval> | null = null;
  private ostatniPuls = 0;
  private odKiedyBezSieci: number | null = null;
  private readonly naPowrotSieci = (): void => this.obsluzPowrotSieci();
  private readonly naUtrateSieci = (): void => {
    this.odKiedyBezSieci = Date.now();
  };
  private readonly naWidocznosc = (): void => this.obsluzWidocznosc();

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
    // Bez adresu ramki wracają porzucone od razu, zamiast czekać w ciszy.
    if (this.adres === '') {
      this.zapiszStan('rozlaczony');
      return;
    }
    this.uruchomPuls();
    this.zapiszStan(this.numerProby === 0 ? 'laczenie' : 'ponawianie');
    // Rdzeń porównuje sekret przed uaktualnieniem (`ParametrSekretu`).
    const gniazdo = new WebSocket(adresNawiazania(this.adres, sekretNawiazaniaOdPowloki()));
    this.gniazdo = gniazdo;
    gniazdo.addEventListener('open', () => this.obsluzOtwarcie(gniazdo));
    gniazdo.addEventListener('message', (zdarzenie) => this.obsluzRamke(gniazdo, zdarzenie));
    gniazdo.addEventListener('close', () => this.obsluzZamkniecie(gniazdo));
    // Zamknięcie na błędzie dotyczy gniazda już otwartego; ponowne close wywołuje nawrót bez końca.
    gniazdo.addEventListener('error', () => {
      if (gniazdo.readyState === WebSocket.OPEN) gniazdo.close();
    });
  }

  wznow(): void {
    if (this.gniazdo !== null) return;
    this.numerProby = 0;
    this.polacz();
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

  rozlacz(): void {
    this.zaniechane = true;
    this.anulujPlan();
    this.zatrzymajPuls();
    const gniazdo = this.gniazdo;
    this.gniazdo = null;
    this.numerProby = 0;
    gniazdo?.close();
    this.zapiszStan('rozlaczony');
  }

  private obsluzOtwarcie(gniazdo: WebSocket): void {
    if (this.gniazdo !== gniazdo) return;
    this.numerProby = 0;
    this.wstrzymana = true;
    this.zapiszStan('polaczony');
    this.wydajKolejke(this.kolejkaPowitania);
  }

  private obsluzRamke(gniazdo: WebSocket, zdarzenie: MessageEvent<unknown>): void {
    if (this.gniazdo !== gniazdo) return;
    if (typeof zdarzenie.data === 'string') {
      this.ramki.oglos(zdarzenie.data);
    }
  }

  private obsluzZamkniecie(gniazdo: WebSocket): void {
    if (this.gniazdo !== gniazdo) return;
    this.gniazdo = null;
    // Rdzeń wiąże sesję bramki z połączeniem, więc powitanie nie przechodzi dalej.
    this.porzucKolejke(this.kolejkaPowitania, 'zerwanie');
    this.zapiszStan('ponawianie');
    this.zaplanujPonowienie();
  }

  private polaczOdNowa(): void {
    if (this.zaniechane) return;
    const stare = this.gniazdo;
    if (stare !== null) {
      this.gniazdo = null;
      stare.close();
      this.porzucKolejke(this.kolejkaPowitania, 'zerwanie');
      this.zapiszStan('ponawianie');
    }
    this.wznow();
  }

  private uruchomPuls(): void {
    if (this.licznikPulsu !== null) return;
    this.ostatniPuls = Date.now();
    this.licznikPulsu = setInterval(() => this.sprawdzPuls(), ODSTEP_PULSU_MS);
    globalThis.addEventListener('online', this.naPowrotSieci);
    globalThis.addEventListener('offline', this.naUtrateSieci);
    document.addEventListener('visibilitychange', this.naWidocznosc);
  }

  private zatrzymajPuls(): void {
    if (this.licznikPulsu === null) return;
    clearInterval(this.licznikPulsu);
    this.licznikPulsu = null;
    this.odKiedyBezSieci = null;
    globalThis.removeEventListener('online', this.naPowrotSieci);
    globalThis.removeEventListener('offline', this.naUtrateSieci);
    document.removeEventListener('visibilitychange', this.naWidocznosc);
  }

  private sprawdzPuls(): void {
    const teraz = Date.now();
    const przerwa = teraz - this.ostatniPuls;
    this.ostatniPuls = teraz;
    // W karcie ukrytej przeglądarka dławi licznik do jednego tyknięcia na minutę.
    if (document.visibilityState === 'hidden') return;
    if (przerwa < PROG_USPIENIA_MS) return;
    this.polaczOdNowa();
  }

  private obsluzWidocznosc(): void {
    if (document.visibilityState !== 'visible') return;
    this.sprawdzPuls();
  }

  private obsluzPowrotSieci(): void {
    const bezSieci = this.odKiedyBezSieci;
    this.odKiedyBezSieci = null;
    if (bezSieci !== null && Date.now() - bezSieci >= PROG_USPIENIA_MS) {
      this.polaczOdNowa();
      return;
    }
    if (!this.zaniechane) this.wznow();
  }

  private wydajAlboOdloz(
    ramka: string,
    kolejka: KolejkaWychodzaca<RamkaOdlozona>,
    wolno: boolean,
  ): void {
    if (wolno && this.gniazdo !== null && this.gniazdo.readyState === WebSocket.OPEN) {
      this.gniazdo.send(ramka);
      return;
    }
    if (this.adres === '') {
      this.porzucone.oglos({ tresc: ramka, powod: 'zerwanie' });
      return;
    }
    const wyparta = kolejka.dodaj({ tresc: ramka, nadana: Date.now() });
    if (wyparta !== undefined) {
      this.porzucone.oglos({ tresc: wyparta.tresc, powod: 'przepelnienie' });
    }
  }

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
