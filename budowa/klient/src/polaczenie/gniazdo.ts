import { adresNawiazania, sekretNawiazaniaOdPowloki } from './adres-rdzenia.ts';
import { utworzKolejkeWychodzaca, type KolejkaWychodzaca } from './kolejka-wychodzaca.ts';
import { utworzMagistrale, type Magistrala, type Odsubskrybuj } from './magistrala-zdarzen.ts';
import { wykladniczePonawianie, type PolitykaPonawiania } from './ponawianie.ts';
import type { StanPolaczenia } from './stan-polaczenia.ts';

/** Powód, dla którego odłożona ramka nie została rdzeniowi wydana. */
export type PowodPorzucenia =
  /** Ramka czekała na łączność dłużej niż zapora czasu odłożenia. */
  | 'zapora-czasu'
  /** Gniazdo, do którego ramka należała, zostało zerwane. */
  | 'zerwanie'
  /** Kolejka wychodząca doszła do sufitu i ramka najstarsza ustąpiła miejsca nowej. */
  | 'przepelnienie';

/** Transport ramek tekstowych do rdzenia, ukrywający przed wołającym stan gniazda i kolejkę wychodzącą. */
export interface Transport {
  /** Rozpoczyna łączenie i utrzymuje je przez ponawianie. */
  polacz(): void;
  /** Ponawia łączenie od razu, bez czekania na zaplanowane opóźnienie. */
  wznow(): void;
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

/*
Sufit kolejki wychodzącej. Rdzeń trzyma na jedno gniazdo 256 ramek wyjściowych
(`transport/ustawienia.go`, `pojemnoscKolejkiDomyslna`) i klient odkłada tyle
samo: więcej nie wyszłoby do rdzenia jednym ciągiem. Ponad sufit najstarsza
ramka ustępuje nowej i wraca porzucona, żeby wołający dostał odmowę, nie ciszę.
*/
const SUFIT_KOLEJKI = 256;

/** Ramka, która rdzeniowi wydana nie będzie, wraz z powodem porzucenia. */
interface RamkaPorzucona {
  tresc: string;
  powod: PowodPorzucenia;
}

/*
Przeglądarka odpowiada na pingi rdzenia sama i nie pokazuje ich skryptowi, więc
o uśpieniu maszyny mówi wyłącznie skok zegara między tyknięciami licznika.
Rdzeń pinguje co 20 s i zamyka gniazdo po 10 s bez odpowiedzi
(`transport/petla_odbioru.go`): przerwa od 30 s znaczy gniazdo już zamknięte
po jego stronie, choć w przeglądarce nadal otwarte.
*/
const ODSTEP_PULSU_MS = 5_000;
const PROG_USPIENIA_MS = 30_000;

/** Połączenie WebSocket z ponawianiem i kolejkowaniem ramek, utrzymujące łączność z rdzeniem bez udziału wołającego. */
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
  /* Rdzeń wiąże sesję bramki z gniazdem dopiero w powitaniu, więc komenda
     wydana przed jego odpowiedzią wraca odmową `not_authenticated`. Każde
     gniazdo zaczyna więc wstrzymane i czeka na powitanie własne. */
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
    /* Transport bez adresu nie ma z czym się łączyć: ramki wracają porzucone
       od razu, żeby wołający dostał odmowę zamiast ciszy. */
    if (this.adres === '') {
      this.zapiszStan('rozlaczony');
      return;
    }
    this.uruchomPuls();
    this.zapiszStan(this.numerProby === 0 ? 'laczenie' : 'ponawianie');
    /* Sekret nawiązania idzie parametrem zapytania: rdzeń porównuje go przed
       uaktualnieniem gniazda (`transport/ustawienia.go`, `ParametrSekretu`). */
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

  /** Zaniechanie kończy ponawianie bezterminowe; ramki odłożone zostają w kolejce, gotowe do wysłania. */
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

  /** Zamknięcie gniazda już odciętego — przy zaniechaniu albo po uśpieniu — nie dotyczy gniazda następcy. */
  private obsluzZamkniecie(gniazdo: WebSocket): void {
    if (this.gniazdo !== gniazdo) return;
    this.gniazdo = null;
    /* Powitanie należy do gniazda, które je przyjęło: rdzeń wiąże po nim sesję
       bramki z konkretnym połączeniem, więc na nowym gnieździe jest bezużyteczne. */
    this.porzucKolejke(this.kolejkaPowitania, 'zerwanie');
    this.zapiszStan('ponawianie');
    this.zaplanujPonowienie();
  }

  /** Odcina gniazdo uznane za niepewne i łączy od razu; powitanie i żądania w locie wracają odmową jak przy zerwaniu. */
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

  /** Skok zegara między tyknięciami ponad próg znaczy uśpienie maszyny, po którym gniazdo otwarte jest niepewne. */
  private sprawdzPuls(): void {
    const teraz = Date.now();
    const przerwa = teraz - this.ostatniPuls;
    this.ostatniPuls = teraz;
    /* Karta ukryta dostaje od przeglądarki licznik dławiony do jednego
       tyknięcia na minutę; przerwa zmierzona w ukryciu mówiłaby o dławieniu,
       nie o śnie. Osąd czeka do powrotu widoczności. */
    if (document.visibilityState === 'hidden') return;
    if (przerwa < PROG_USPIENIA_MS) return;
    this.polaczOdNowa();
  }

  private obsluzWidocznosc(): void {
    if (document.visibilityState !== 'visible') return;
    this.sprawdzPuls();
  }

  /** Po powrocie sieci łączy od razu; gniazdo otwarte przez przerwę od progu rdzeń już zamknął, więc idzie do odcięcia. */
  private obsluzPowrotSieci(): void {
    const bezSieci = this.odKiedyBezSieci;
    this.odKiedyBezSieci = null;
    if (bezSieci !== null && Date.now() - bezSieci >= PROG_USPIENIA_MS) {
      this.polaczOdNowa();
      return;
    }
    if (!this.zaniechane) this.wznow();
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
    if (this.adres === '') {
      this.porzucone.oglos({ tresc: ramka, powod: 'zerwanie' });
      return;
    }
    const wyparta = kolejka.dodaj({ tresc: ramka, nadana: Date.now() });
    if (wyparta !== undefined) {
      this.porzucone.oglos({ tresc: wyparta.tresc, powod: 'przepelnienie' });
    }
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
