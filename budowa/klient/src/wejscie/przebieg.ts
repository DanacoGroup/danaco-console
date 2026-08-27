/**
 * DROGA WEJŚCIA — przebieg.
 *
 * Maszyna stanów trzech etapów wejścia: łączenie z rdzeniem, dostęp do konta,
 * przygotowanie środowiska. Przebieg nie dotyka dokumentu i nie zna żadnego
 * składnika — rozstrzyga wyłącznie, który etap i która odsłona obowiązuje oraz
 * co wysłać do rdzenia. Dzięki temu każda odsłona, także odsłona błędu
 * i wstrzymania, jest osiągalna w sprawdzianie bez przeglądarki.
 *
 * Nazwy komend i kształty ich treści pochodzą wyłącznie z kontraktu; przebieg
 * nie powtarza ani jednego literału nazwy komendy.
 */

import {
  Command,
  ErrorCode,
  PROTOCOL_VERSION,
  type AuthMethodKind,
  type AuthSession,
  type ConnectionHelloResponse,
  type Environment,
  type ErrorInfo,
  type Module,
  type RequestOf,
  type ResponseOf,
  type Session,
} from '../../../shared/contract.ts';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia.ts';
import type { Transport } from '../polaczenie/gniazdo.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { zadajPowitanie } from '../protokol/powitanie.ts';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/* ── Stan widziany przez widok ───────────────────────────────────────────── */

/** Etap wejścia. Trzy, tak jak w prototypie. */
export type Etap = 'uruchomienie' | 'dostep' | 'przygotowanie';

/** Odsłona etapu uruchomienia. */
export type OdslonaUruchomienia = 'w-laczenie' | 'w-token' | 'w-blad';

/**
 * Odsłona etapu dostępu.
 *
 * `konto-bez-potwierdzenia` nie pochodzi z prototypu — wymusza ją pozycja 11
 * rejestru decyzji: rejestracja bez konta nadawczego kończy się wejściem
 * hasłem, a okno ma wtedy NAZWAĆ niepotwierdzony adres.
 */
export type OdslonaDostepu =
  | 'logowanie'
  | 'logowanie-blad'
  | 'logowanie-wstrzymane'
  | 'rejestracja'
  | 'kod'
  | 'konto-bez-potwierdzenia'
  | 'odzyskiwanie-adres'
  | 'odzyskiwanie-kod'
  | 'odzyskiwanie-haslo'
  | 'odzyskiwanie-wstrzymane';

/** Odsłona etapu przygotowania. Jedna — nie ma tu czego przełączać. */
export type OdslonaPrzygotowania = 'przygotowanie';

export type Odslona = OdslonaUruchomienia | OdslonaDostepu | OdslonaPrzygotowania;

/** Stan jednego z etapów przygotowania środowiska. */
export type StanEtapu = 'gotowy' | 'pracuje' | 'blad' | 'oczekuje';

/** Usterka pokazywana nad formularzem: klucz katalogu albo opis od rdzenia. */
export interface Usterka {
  /** Klucz rozpoznania z katalogu treści; puste dla odmowy rdzenia. */
  klucz?: string;
  /** Opis wprost od rdzenia — rdzeń wie o powodzie odmowy więcej niż okno. */
  odRdzenia?: ErrorInfo;
}

/** Miara etapu przygotowania odczytana z odpowiedzi rdzenia. */
export interface MiaraEtapu {
  /** Klucz miary z katalogu treści. */
  klucz: string;
  /** Dane podstawiane w miarę. */
  dane?: Record<string, string | number>;
}

/** Wykaz etapów przygotowania wraz ze stanem i miarą każdego z nich. */
export interface PostepPrzygotowania {
  stany: StanEtapu[];
  miary: MiaraEtapu[];
  /** Postęp przywracania w procentach, liczony z etapów zakończonych. */
  wartosc: number;
}

/** Pełny stan przebiegu widziany przez widok. */
export interface StanPrzebiegu {
  etap: Etap;
  odslona: Odslona;
  /** Stany czterech etapów łączenia; wykaz jest daną, nie rozgałęzieniem. */
  etapyLaczenia: StanEtapu[];
  /** Miary czterech etapów łączenia, kluczami katalogu. */
  miaryLaczenia: string[];
  /** Czy trwa wywołanie, na które przebieg czeka. */
  wToku: boolean;
  /** Usterki nad formularzem bieżącej odsłony. */
  usterki: Usterka[];
  /**
   * Zwłoka nałożona przez rdzeń na kolejną próbę, w sekundach — ZMIERZONA
   * czasem trwania próby poprzedniej. Odmowa rdzenia nie niesie tej wartości
   * (pole `details` jest puste), więc okno nie ma jej skąd odczytać; mierzy
   * więc to, ile rdzeń kazał czekać naprawdę, i tego nie zgaduje.
   */
  zwlokaS: number;
  /** Adres, na który poszedł list — okno wpisuje go w zdanie odsłony. */
  adres: string;
  /** Odpowiedź powitania; niesie wersję rdzenia i stan uwierzytelnienia. */
  powitanie?: ConnectionHelloResponse;
  /** Sesja bramki wydana logowaniem, potwierdzeniem albo tokenem. */
  sesjaBramki?: AuthSession;
  /** Postęp trzeciego etapu. */
  przygotowanie: PostepPrzygotowania;
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko?: Environment;
  /** Moduły środowiska; liczba zasila miarę etapu przygotowania. */
  moduly: Module[];
  /** Karty sesji odtworzone przez rdzeń. */
  sesje: Session[];
}

/* ── Nastawy przebiegu ───────────────────────────────────────────────────── */

/**
 * Od jakiej zwłoki okno nazywa ją Operatorowi, w milisekundach.
 *
 * Progu prób NIE MA. Kontrakt stanowi przy `auth.login` wprost, że progu i
 * odmowy „za dużo prób” po stronie rdzenia nie ma i nie będzie — nieudana
 * próba nakłada na następną rosnącą zwłokę, nic więcej. Pomiar to potwierdza:
 * osiem kolejnych nieudanych prób oddaje osiem razy `not_authenticated`, a po
 * nich hasło poprawne wpuszcza. Odsłona z prototypu zostaje i obsługuje
 * ZWŁOKĘ, nie zaporę.
 *
 * Próg bierze się z jednostki wyświetlania, nie z reguły produktu: licznik
 * odmierza sekundy, więc zwłoka krótsza od sekundy nie ma czego pokazać
 * i odsłona tylko mignęłaby.
 */
const PROG_NAZWANIA_ZWLOKI_MS = 1000;

/** Ile znaków ma mieć hasło. Reguła wspólna miernikowi siły i sprawdzeniu. */
const NAJKROTSZE_HASLO = 12;

/**
 * Cztery warunki hasła — jedna reguła dla miernika siły w oknie i dla
 * sprawdzenia przed wysłaniem. Dwie osobne rozjechałyby się przy pierwszej
 * zmianie wymagań, a Operator zobaczyłby miernik zielony przy haśle odrzuconym.
 */
export function ocenHaslo(wartosc: string): Record<string, boolean> {
  return {
    dlugosc: wartosc.length >= NAJKROTSZE_HASLO,
    wielkosc: /[a-ząćęłńóśźż]/.test(wartosc) && /[A-ZĄĆĘŁŃÓŚŹŻ]/.test(wartosc),
    cyfra: /[0-9]/.test(wartosc),
    znak: /[^0-9A-Za-zĄĆĘŁŃÓŚŹŻąćęłńóśźż]/.test(wartosc),
  };
}

/** Czy hasło spełnia wszystkie cztery warunki. */
export function hasloSpelnia(wartosc: string): boolean {
  return Object.values(ocenHaslo(wartosc)).every(Boolean);
}

/**
 * Budowa adresu, nie jego istnienie. Że adres istnieje, rozstrzyga dopiero
 * list wysłany na niego — tutaj wyłapuje się to, co widać bez wysyłania: brak
 * znaku małpy, brak nazwy, brak domeny, odstęp w środku.
 */
export function adresPoprawny(wartosc: string): boolean {
  return /^[^\s@]+@[^\s@.]+(\.[^\s@.]+)+$/.test(wartosc.trim());
}

/** Metoda wejścia, którą niesie okno bramki. Kontrakt zna trzy; PIN i klucz
 * urządzenia zakłada się dopiero w Ustawieniach, czyli za bramką. */
const METODA_HASLEM: AuthMethodKind = 'password';

/**
 * Magazyn tokenu bramki między uruchomieniami programu.
 *
 * Przebieg czyta token przy starcie, żeby rozpoznać zaufane urządzenie, i
 * zapisuje go po wejściu. GDZIE token mieszka, nie rozstrzyga żaden ze
 * źródeł tego terenu — należy to do powłoki. Przebieg opisuje więc wyłącznie
 * potrzebę, a wołający wskazuje magazyn.
 */
export interface MagazynTokenu {
  odczytaj(): string | undefined;
  zapisz(token: string): void;
  wyczysc(): void;
}

/**
 * Magazyn trzymający token wyłącznie w pamięci procesu.
 *
 * Wartość domyślna, bo magazyn trwały nie jest rozstrzygnięty. Token ginie
 * wraz z procesem — i to jest zachowanie uczciwe: magazyn, który udawałby
 * trwałość, obiecywałby rozpoznanie urządzenia, którego nie ma.
 */
export function magazynWPamieci(): MagazynTokenu {
  let token: string | undefined;
  return {
    odczytaj: () => token,
    zapisz(nowy) {
      token = nowy;
    },
    wyczysc() {
      token = undefined;
    },
  };
}

export interface ZaleznosciPrzebiegu {
  kanal: Kanal;
  transport: Transport;
  klient: TozsamoscKlienta;
  magazyn?: MagazynTokenu;
}

/* ── Dane formularzy ─────────────────────────────────────────────────────── */

export interface DaneLogowania {
  login: string;
  haslo: string;
  niewylogowuj?: boolean;
}

export interface DaneRejestracji {
  login: string;
  email: string;
  haslo: string;
  hasloPowtorzone: string;
  niewylogowuj?: boolean;
}

export interface DanePotwierdzenia {
  droga: string;
  niewylogowuj?: boolean;
}

export interface DaneNowegoHasla {
  droga: string;
  haslo: string;
  hasloPowtorzone: string;
}

/* ── Przebieg ────────────────────────────────────────────────────────────── */

export interface Przebieg {
  /** Bieżący stan; widok czyta go po każdej zmianie. */
  stan(): StanPrzebiegu;
  /** Subskrypcja zmian stanu. Wywołanie zwrócone odłącza słuchacza. */
  naZmiane(sluchacz: (stan: StanPrzebiegu) => void): () => void;

  /** Etap 1: nawiązanie połączenia i powitanie rdzenia. */
  polacz(): void;
  /** Etap 1: powtórzenie próby po błędzie połączenia. */
  ponow(): void;

  /** Etap 2: przełączenie odsłony bez wysyłania czegokolwiek. */
  przejdzDo(odslona: OdslonaDostepu): void;
  /** Etap 2: logowanie hasłem. */
  zaloguj(dane: DaneLogowania): Promise<void>;
  /** Etap 2: założenie konta Operatora. */
  zarejestruj(dane: DaneRejestracji): Promise<void>;
  /** Etap 2: potwierdzenie adresu drogą z listu. */
  potwierdzAdres(dane: DanePotwierdzenia): Promise<void>;
  /** Etap 2: prośba o drogę odzyskania konta. */
  poprosOOdzyskanie(email: string): Promise<void>;
  /** Etap 2: ustawienie nowego hasła drogą z listu. */
  ustawNoweHaslo(dane: DaneNowegoHasla): Promise<void>;
  /** Etap 2: wejście do platformy po rejestracji bez konta nadawczego. */
  wejdzBezPotwierdzenia(): Promise<void>;
  /** Etap 2: zdjęcie wstrzymania po upływie godziny. */
  zdejmijWstrzymanie(): void;

  /** Etap 3: przygotowanie środowiska pracy. */
  przygotujSrodowisko(): Promise<void>;
}

/** Cztery etapy łączenia — wykaz jest daną, nie rozgałęzieniem w kodzie. */
const MIARY_LACZENIA = ['nawiazane', 'zgodna', 'zaufane', 'wToku'] as const;
const MIARY_NIEUDANE = ['nieudane', 'oczekuje', 'oczekuje', 'oczekuje'] as const;

/**
 * Ile etapów ma przygotowanie środowiska.
 *
 * Tyle, ile droga wejścia potrafi zmierzyć: uwierzytelnienie i przywrócenie
 * kart sesji. Etap bez komendy w kontrakcie nie jest etapem czekającym —
 * jest obietnicą, której nikt nie wykona, a postęp liczony razem z nim nie
 * dobiegłby końca nigdy.
 */
const ETAPOW_PRZYGOTOWANIA = 2;

export function utworzPrzebieg(zaleznosci: ZaleznosciPrzebiegu): Przebieg {
  const { kanal, transport, klient } = zaleznosci;
  const magazyn = zaleznosci.magazyn ?? magazynWPamieci();
  const sluchacze = new Set<(stan: StanPrzebiegu) => void>();

  let stan: StanPrzebiegu = stanPoczatkowy();
  let odsubskrybujStan: (() => void) | undefined;

  function stanPoczatkowy(): StanPrzebiegu {
    return {
      etap: 'uruchomienie',
      odslona: 'w-laczenie',
      etapyLaczenia: ['pracuje', 'oczekuje', 'oczekuje', 'oczekuje'],
      miaryLaczenia: ['wToku', 'oczekuje', 'oczekuje', 'oczekuje'],
      wToku: false,
      usterki: [],
      zwlokaS: 0,
      adres: '',
      moduly: [],
      sesje: [],
      przygotowanie: {
        stany: nowyWykazStanow(),
        miary: nowyWykazMiar(),
        wartosc: 0,
      },
    };
  }

  function nowyWykazStanow(): StanEtapu[] {
    return Array.from({ length: ETAPOW_PRZYGOTOWANIA }, () => 'oczekuje' as StanEtapu);
  }

  function nowyWykazMiar(): MiaraEtapu[] {
    return Array.from({ length: ETAPOW_PRZYGOTOWANIA }, () => ({ klucz: 'oczekuje' }));
  }

  function zmien(czesc: Partial<StanPrzebiegu>): void {
    stan = { ...stan, ...czesc };
    for (const sluchacz of [...sluchacze]) {
      try {
        sluchacz(stan);
      } catch (blad) {
        console.error('[wejście] błąd słuchacza przebiegu', blad);
      }
    }
  }

  /* ── Etap 1: łączenie ──────────────────────────────────────────────────── */

  function polacz(): void {
    zmien({
      etap: 'uruchomienie',
      odslona: 'w-laczenie',
      etapyLaczenia: ['pracuje', 'oczekuje', 'oczekuje', 'oczekuje'],
      miaryLaczenia: ['wToku', 'oczekuje', 'oczekuje', 'oczekuje'],
      usterki: [],
    });
    odsubskrybujStan?.();
    odsubskrybujStan = transport.naStan(obsluzStanPolaczenia);
    transport.polacz();
  }

  function obsluzStanPolaczenia(stanGniazda: StanPolaczenia): void {
    if (stanGniazda === 'polaczony') {
      zmien({
        etapyLaczenia: ['gotowy', 'pracuje', 'oczekuje', 'oczekuje'],
        miaryLaczenia: ['nawiazane', 'wToku', 'oczekuje', 'oczekuje'],
      });
      void przywitaj();
      return;
    }
    // Ponawianie jest stanem błędu widzianym przez Operatora: gniazdo próbuje
    // dalej samo, a okno mówi, że serwer nie odpowiada. Do etapu 1 wracamy
    // wyłącznie stąd — zerwanie po wejściu nie cofa Operatora przed bramkę.
    if (stanGniazda === 'ponawianie' && stan.etap === 'uruchomienie') {
      zmien({
        odslona: 'w-blad',
        etapyLaczenia: ['blad', 'oczekuje', 'oczekuje', 'oczekuje'],
        miaryLaczenia: [...MIARY_NIEUDANE],
      });
    }
  }

  async function przywitaj(): Promise<void> {
    const zapamietany = magazyn.odczytaj();
    const wynik = await zadajPowitanie(kanal, klient, zapamietany);
    if (!wynik.udany || wynik.wynik === undefined) {
      zmien({
        odslona: 'w-blad',
        etapyLaczenia: ['gotowy', 'blad', 'oczekuje', 'oczekuje'],
        miaryLaczenia: ['nawiazane', 'nieudane', 'oczekuje', 'oczekuje'],
        usterki: [{ odRdzenia: wynik.blad }],
      });
      return;
    }
    const powitanie = wynik.wynik;
    // Wersja protokołu jest jedyną rzeczą, którą klient uzgadnia z rdzeniem
    // przed czymkolwiek innym. Rozjazd zatrzymuje wejście tutaj: dalsza
    // rozmowa szłaby po omacku, a odmowy nie dałoby się odróżnić od usterki.
    if (powitanie.protocolVersion !== PROTOCOL_VERSION) {
      zmien({
        powitanie,
        odslona: 'w-blad',
        etapyLaczenia: ['gotowy', 'blad', 'oczekuje', 'oczekuje'],
        miaryLaczenia: ['nawiazane', 'nieudane', 'oczekuje', 'oczekuje'],
        usterki: [{ klucz: 'uruchomienie.wersja' }],
      });
      return;
    }
    if (powitanie.authenticated === true) {
      zmien({
        powitanie,
        odslona: 'w-token',
        etapyLaczenia: ['gotowy', 'gotowy', 'gotowy', 'pracuje'],
        miaryLaczenia: [...MIARY_LACZENIA],
      });
      void przygotujSrodowisko();
      return;
    }
    zmien({
      powitanie,
      etapyLaczenia: ['gotowy', 'gotowy', 'gotowy', 'gotowy'],
      miaryLaczenia: ['nawiazane', 'zgodna', 'zaufane', 'gotowe'],
      etap: 'dostep',
      odslona: 'logowanie',
    });
  }

  function ponow(): void {
    magazyn.wyczysc();
    polacz();
  }

  /* ── Etap 2: dostęp do konta ───────────────────────────────────────────── */

  function przejdzDo(odslona: OdslonaDostepu): void {
    zmien({ etap: 'dostep', odslona, usterki: [] });
  }

  /**
   * Wywołanie komendy z jedną obsługą oczekiwania dla wszystkich formularzy.
   *
   * Obietnica warstwy protokołu nie jest odrzucana nigdy, a przy zerwanym
   * połączeniu nie rozstrzyga się wcale. Stan `wToku` jest więc jedynym
   * miejscem, w którym okno mówi, że czeka. Typ treści żądania i odpowiedzi
   * bierze się z kontraktu — okno nie deklaruje ani jednego kształtu własnego.
   */
  async function doRdzenia<K extends Command>(
    komenda: K,
    zadanie: RequestOf<K>,
  ): Promise<Wynik<ResponseOf<K>>> {
    zmien({ wToku: true, usterki: [] });
    const poczatek = Date.now();
    const wynik = await wywolaj(kanal, komenda, zadanie);
    ostatniCzasMs = Date.now() - poczatek;
    zmien({ wToku: false });
    return wynik;
  }

  /**
   * Ile trwało ostatnie wywołanie. Rdzeń nakłada zwłokę PRZED odpowiedzią, więc
   * czas trwania wywołania JEST zwłoką — innego jej pomiaru okno nie ma.
   */
  let ostatniCzasMs = 0;

  /** Zwłoka w sekundach, zaokrąglona w górę; zero, gdy nie ma czego nazywać. */
  function zmierzonaZwlokaS(): number {
    if (ostatniCzasMs < PROG_NAZWANIA_ZWLOKI_MS) return 0;
    return Math.ceil(ostatniCzasMs / 1000);
  }

  function odmowa(wynik: Wynik<unknown>): Usterka[] {
    return [{ odRdzenia: wynik.blad }];
  }

  async function zaloguj(dane: DaneLogowania): Promise<void> {
    const braki = brakiLogowania(dane);
    if (braki.length > 0) {
      zmien({ usterki: braki });
      return;
    }
    const wynik = await doRdzenia(Command.AuthLogin, {
      method: METODA_HASLEM,
      login: dane.login.trim(),
      secret: dane.haslo,
      keepSignedIn: dane.niewylogowuj === true,
    });
    if (!wynik.udany) {
      nazwijOdmoweLogowania(wynik);
      return;
    }
    const odpowiedz = wynik.wynik;
    if (odpowiedz === undefined) {
      zmien({ usterki: [{ klucz: 'usterki.brakOdpowiedzi' }] });
      return;
    }
    przyjmijSesje(odpowiedz.session);
    await przygotujSrodowisko();
  }

  function brakiLogowania(dane: DaneLogowania): Usterka[] {
    const bezLoginu = dane.login.trim().length === 0;
    const bezHasla = dane.haslo.length === 0;
    // Oba pola puste to JEDNA sprawa — formularz nie został wypełniony.
    if (bezLoginu && bezHasla) return [{ klucz: 'brakDanych' }];
    if (bezLoginu) return [{ klucz: 'brakLoginu' }];
    if (bezHasla) return [{ klucz: 'brakHasla' }];
    return [];
  }

  /**
   * Odmowa logowania. Progu prób nie ma — każda kolejna próba jest przyjmowana,
   * tylko czeka dłużej. Gdy zwłoka urosła na tyle, że da się ją nazwać, okno
   * przechodzi do odsłony zwłoki i samo z niej wraca; przy zwłoce krótszej
   * zostaje na odsłonie niepowodzenia, bo nie ma czego odmierzać.
   */
  function nazwijOdmoweLogowania(wynik: Wynik<unknown>): void {
    const zwlokaS = wynik.blad?.code === ErrorCode.NotAuthenticated ? zmierzonaZwlokaS() : 0;
    if (zwlokaS > 0) {
      zmien({ zwlokaS, odslona: 'logowanie-wstrzymane', usterki: [] });
      return;
    }
    zmien({ zwlokaS: 0, odslona: 'logowanie-blad', usterki: odmowa(wynik) });
  }

  /** Zwłoka minęła — okno wraca tam, skąd Operator ją zastał. */
  function zdejmijWstrzymanie(): void {
    if (stan.odslona === 'odzyskiwanie-wstrzymane') {
      zmien({ zwlokaS: 0, odslona: 'odzyskiwanie-adres', usterki: [] });
      return;
    }
    zmien({ zwlokaS: 0, odslona: 'logowanie', usterki: [] });
  }

  async function zarejestruj(dane: DaneRejestracji): Promise<void> {
    const braki = brakiRejestracji(dane);
    if (braki.length > 0) {
      zmien({ usterki: braki });
      return;
    }
    const wynik = await doRdzenia(Command.AuthRegister, {
      login: dane.login.trim(),
      email: dane.email.trim(),
      password: dane.haslo,
    });
    if (!wynik.udany) {
      zmien({ usterki: odmowa(wynik) });
      return;
    }
    const odpowiedz = wynik.wynik;
    if (odpowiedz === undefined || !odpowiedz.registered) {
      zmien({ usterki: [{ klucz: 'usterki.brakOdpowiedzi' }] });
      return;
    }
    // Dwie gałęzie pozycji 11 rejestru decyzji, rozstrzygane odpowiedzią rdzenia.
    // Z kontem nadawczym list poszedł i okno prowadzi do jego przepisania. Bez
    // konta nadawczego wejście działa hasłem, a okno MUSI nazwać adres, którego
    // nikt nie potwierdził — bo adres jest jedyną drogą odzyskania konta.
    zmien({
      adres: dane.email.trim(),
      odslona: odpowiedz.pendingVerification ? 'kod' : 'konto-bez-potwierdzenia',
      usterki: [],
    });
  }

  function brakiRejestracji(dane: DaneRejestracji): Usterka[] {
    const braki: Usterka[] = [];
    if (dane.login.trim().length === 0) braki.push({ klucz: 'brakLoginu' });
    if (dane.email.trim().length === 0) braki.push({ klucz: 'brakAdresu' });
    else if (!adresPoprawny(dane.email)) braki.push({ klucz: 'email-bledny' });
    if (dane.haslo.length === 0) braki.push({ klucz: 'brakHasla' });
    else if (!hasloSpelnia(dane.haslo)) braki.push({ klucz: 'haslo-slabe' });
    if (dane.haslo !== dane.hasloPowtorzone) braki.push({ klucz: 'hasla-rozne' });
    return braki;
  }

  async function potwierdzAdres(dane: DanePotwierdzenia): Promise<void> {
    if (dane.droga.trim().length === 0) {
      zmien({ usterki: [{ klucz: 'brakDrogi' }] });
      return;
    }
    const wynik = await doRdzenia(Command.AuthVerify, {
      token: dane.droga.trim(),
      keepSignedIn: dane.niewylogowuj === true,
    });
    if (!wynik.udany) {
      zmien({ usterki: odmowa(wynik) });
      return;
    }
    const odpowiedz = wynik.wynik;
    if (odpowiedz === undefined || !odpowiedz.verified) {
      zmien({ usterki: [{ klucz: 'usterki.brakOdpowiedzi' }] });
      return;
    }
    przyjmijSesje(odpowiedz.session);
    await przygotujSrodowisko();
  }

  async function poprosOOdzyskanie(email: string): Promise<void> {
    if (email.trim().length === 0) {
      zmien({ usterki: [{ klucz: 'brakAdresu' }] });
      return;
    }
    if (!adresPoprawny(email)) {
      zmien({ usterki: [{ klucz: 'email-bledny' }] });
      return;
    }
    // Progu wysyłek nie ma — zmierzone: siedem kolejnych wysłań oddaje siedem
    // razy `sent: true` i siedem listów. Gdyby rdzeń kiedyś nałożył zwłokę,
    // okno nazwie ją tym samym pomiarem, którym nazywa zwłokę logowania.
    const wynik = await doRdzenia(Command.AuthRecover, { email: email.trim() });
    if (!wynik.udany) {
      zmien({ usterki: odmowa(wynik) });
      return;
    }
    const odpowiedz = wynik.wynik;
    // Odpowiedź nie zdradza, czy adres pasuje do konta — tak stanowi kontrakt.
    // Okno nie ma więc czego z niej odczytać poza tym, że żądanie przyjęto.
    if (odpowiedz === undefined || !odpowiedz.sent) {
      zmien({ usterki: [{ klucz: 'usterki.brakOdpowiedzi' }] });
      return;
    }
    const zwlokaS = zmierzonaZwlokaS();
    zmien({
      adres: email.trim(),
      zwlokaS,
      odslona: zwlokaS > 0 ? 'odzyskiwanie-wstrzymane' : 'odzyskiwanie-kod',
    });
  }

  async function ustawNoweHaslo(dane: DaneNowegoHasla): Promise<void> {
    const braki: Usterka[] = [];
    if (dane.droga.trim().length === 0) braki.push({ klucz: 'brakDrogi' });
    if (dane.haslo.length === 0) braki.push({ klucz: 'brakHasla' });
    else if (!hasloSpelnia(dane.haslo)) braki.push({ klucz: 'haslo-slabe' });
    if (dane.haslo !== dane.hasloPowtorzone) braki.push({ klucz: 'hasla-rozne' });
    if (braki.length > 0) {
      zmien({ usterki: braki });
      return;
    }
    const wynik = await doRdzenia(Command.AuthReset, {
      token: dane.droga.trim(),
      newPassword: dane.haslo,
    });
    if (!wynik.udany) {
      zmien({ usterki: odmowa(wynik) });
      return;
    }
    const odpowiedz = wynik.wynik;
    if (odpowiedz === undefined || !odpowiedz.changed) {
      zmien({ usterki: [{ klucz: 'usterki.brakOdpowiedzi' }] });
      return;
    }
    // Nowe hasło unieważnia tokeny wydane wcześniej, więc okno wraca do
    // logowania: sesji ta komenda nie wydaje i wydać nie może.
    magazyn.wyczysc();
    zmien({ zwlokaS: 0, odslona: 'logowanie', usterki: [] });
  }

  /**
   * Wejście po rejestracji bez konta nadawczego.
   *
   * Rejestracja sesji nie zakłada — tak stanowi kontrakt — więc wejście idzie
   * przez zwykłe logowanie hasłem, które Operator przed chwilą ustawił.
   * Okno prowadzi go do logowania, bo hasła nie przechowuje.
   */
  async function wejdzBezPotwierdzenia(): Promise<void> {
    zmien({ odslona: 'logowanie', usterki: [] });
    return Promise.resolve();
  }

  function przyjmijSesje(sesja: AuthSession): void {
    magazyn.zapisz(sesja.token);
    zmien({ sesjaBramki: sesja });
  }

  /* ── Etap 3: przygotowanie środowiska ──────────────────────────────────── */

  /**
   * Trzeci etap: wejście do środowiska i odczytanie tego, co rdzeń odtworzył.
   *
   * Wykaz niesie te etapy, dla których droga wejścia ma komendę:
   * uwierzytelnienie i przywracanie kart sesji. Etap oznaczony jako gotowy
   * bez pomiaru mówiłby nieprawdę, a etap czekający na komendę, której nie ma,
   * zatrzymywałby postęp na zawsze — więc wykazu nie ma dłuższego niż pomiar.
   *
   * Wywołanie powtórne jest ponowieniem: wykaz stanów rusza od nowa, więc
   * odsłona po nieudanym wejściu wraca do stanu sprzed próby.
   */
  async function przygotujSrodowisko(): Promise<void> {
    const stany = nowyWykazStanow();
    const miary = nowyWykazMiar();
    stany[0] = 'gotowy';
    miary[0] = { klucz: 'rozpoznane' };
    stany[1] = 'pracuje';
    zmien({
      etap: 'przygotowanie',
      odslona: 'przygotowanie',
      usterki: [],
      przygotowanie: { stany, miary, wartosc: postep(stany) },
    });

    const wykaz = await doRdzenia(Command.EnvironmentList, {});
    const pierwsze = pierwszeSrodowisko(wykaz.wynik?.environments);
    if (!wykaz.udany || pierwsze === undefined) {
      stany[1] = 'blad';
      zmien({
        usterki: odmowa(wykaz),
        przygotowanie: { stany, miary, wartosc: postep(stany) },
      });
      return;
    }

    const wejscie = await doRdzenia(Command.EnvironmentEnter, {
      environmentId: pierwsze.id,
      clientId: klient.id,
    });
    if (!wejscie.udany) {
      stany[1] = 'blad';
      zmien({
        usterki: odmowa(wejscie),
        przygotowanie: { stany, miary, wartosc: postep(stany) },
      });
      return;
    }
    const odpowiedz = wejscie.wynik;
    if (odpowiedz === undefined) {
      stany[1] = 'blad';
      zmien({
        usterki: [{ klucz: 'usterki.brakOdpowiedzi' }],
        przygotowanie: { stany, miary, wartosc: postep(stany) },
      });
      return;
    }
    stany[1] = 'gotowy';
    miary[1] = {
      klucz: 'karty',
      dane: { odtworzone: odpowiedz.sessions.length, wszystkie: odpowiedz.sessions.length },
    };
    zmien({
      srodowisko: odpowiedz.environment,
      moduly: odpowiedz.modules,
      sesje: odpowiedz.sessions,
      przygotowanie: { stany, miary, wartosc: postep(stany) },
    });
  }

  /**
   * Środowisko pierwsze w kolejności wyświetlania. Kolejność niesie kontrakt
   * polem `order`, liczonym od jedynki — nie wybiera go okno.
   */
  function pierwszeSrodowisko(srodowiska?: Environment[]): Environment | undefined {
    if (srodowiska === undefined || srodowiska.length === 0) return undefined;
    return [...srodowiska].sort((a, b) => a.order - b.order)[0];
  }

  /** Postęp liczony z etapów zakończonych, nie z upływu czasu. */
  function postep(stany: StanEtapu[]): number {
    const gotowe = stany.filter((s) => s === 'gotowy').length;
    return Math.round((gotowe / stany.length) * 100);
  }

  return {
    stan: () => stan,
    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => {
        sluchacze.delete(sluchacz);
      };
    },
    polacz,
    ponow,
    przejdzDo,
    zaloguj,
    zarejestruj,
    potwierdzAdres,
    poprosOOdzyskanie,
    ustawNoweHaslo,
    wejdzBezPotwierdzenia,
    zdejmijWstrzymanie,
    przygotujSrodowisko,
  };
}

/** Ile znaków wymaga hasło — miernik siły wypisuje tę liczbę w warunku. */
export const DLUGOSC_HASLA = NAJKROTSZE_HASLO;

/** Od jakiej zwłoki okno ją nazywa. Sprawdzian mierzy tę granicę. */
export const PROG_ZWLOKI_MS = PROG_NAZWANIA_ZWLOKI_MS;
