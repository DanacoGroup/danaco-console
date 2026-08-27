import {
  ProgressStatus,
  type LoopState,
  type MonitorStatus,
  type ProgressChangedEvent,
  type WindowStateChangedEvent,
} from '../../../shared/contract';
import { PROGI_AOD } from './progi-aod';
import { RodzajSugestii, WagaUjawnienia } from './rodzaje-sugestii';

/**
 * Rozpoznanie zdarzenia wymagającego decyzji — czysta reguła nakładki, bez dotyku
 * dokumentu i kanału; budzi to, przy czym proces stoi. Poniżej sygnał, z którego
 * rozpoznanie wzięło obserwację.
 */
export const ZrodloSygnalu = {
  /** `monitor.status` — odczyt nadrabiający; działa bez argumentów. */
  Monitor: 'monitor.status',
  /** `progress.changed` — zdarzenie rozgłaszane do każdego połączenia. */
  Postep: 'progress.changed',
  /** `window.state.changed` — jedyny nośnik stanu biegu naprawczego. */
  StanOkna: 'window.state.changed',
} as const;
export type ZrodloSygnalu = (typeof ZrodloSygnalu)[keyof typeof ZrodloSygnalu];

/** Powód, dla którego nakładka uznała, że dany proces czeka teraz na decyzję Operatora, wraz z jego pewnością. */
export const PowodDecyzji = {
  /** Bieg naprawczy koordynatora ogłosił zatrzymanie (`loop.stopped`). */
  BiegStanal: 'bieg-stanal',
  /** Proces wstrzymany (`paused`). */
  Wstrzymany: 'wstrzymany',
  /** Proces zatrzymany (`stopped`). */
  Zatrzymany: 'zatrzymany',
  /** Proces zakończony błędem (`failed`). */
  Usterka: 'usterka',
  /** Proces oczekuje (`pending`) dłużej niż próg braku ruchu. */
  BezRuchu: 'bez-ruchu',
} as const;
export type PowodDecyzji = (typeof PowodDecyzji)[keyof typeof PowodDecyzji];

/** Czyja to ocena: stanu rdzenia mówiącego to wprost, czy nakładki wnioskującej o zastoju danego procesu. */
export const WagaDecyzji = {
  /** Proces stoi — stan rdzenia mówi to wprost. */
  Pewna: 'pewna',
  /** To ocena nakładki, a nie stan orzeczony przez rdzeń. */
  Sporna: 'sporna',
} as const;
export type WagaDecyzji = (typeof WagaDecyzji)[keyof typeof WagaDecyzji];

/**
 * Próg braku ruchu dla stanu `pending`: 15 minut z `PROGI_AOD` — kontrakt progu nie
 * niesie, więc jest argumentem reguły, nie stałą wewnętrzną.
 */
export const PROG_BEZ_RUCHU_MS = PROGI_AOD.oczekiwanieWKolejceMs;

/** Przedrostek klucza wpisu, którego sygnał nie niósł identyfikatora procesu, tylko identyfikator okna. */
export const PRZEDROSTEK_KLUCZA_OKNA = 'okno:';

/** Jedna obserwacja procesu sprowadzona do wspólnego kształtu, niezależnie od sygnału, z którego przyszła. */
export interface ObserwacjaProcesu {
  /** Tożsamość wpisu: identyfikator procesu albo `okno:<id>`, gdy procesu brak. */
  klucz: string;
  /** Identyfikator procesu; pusty dla sygnału niosącego samo okno. */
  idProcesu?: string;
  /** Okno, którego proces dotyczy. */
  idOkna?: string;
  /** Karta sesji procesu. */
  idSesji?: string;
  /** Nazwa procesu do wyświetlenia. */
  nazwa?: string;
  /** Stan procesu wprost z telemetrii. */
  stan: ProgressStatus;
  /** Nazwa etapu bieżącego. */
  etap?: string;
  /** Numer etapu bieżącego, liczony od 1. */
  numerEtapu?: number;
  /** Liczba etapów; 0 albo brak, gdy nieznana. */
  liczbaEtapow?: number;
  /** Stopień ukończenia w procentach. */
  ukonczenie?: number;
  /** Stan biegu naprawczego koordynatora, gdy sygnał go niósł. */
  bieg?: LoopState;
  /** Czas ostatniej zmiany w milisekundach epoki. */
  zmienionoO: number;
  /** Sygnał, z którego obserwacja pochodzi. */
  zrodlo: ZrodloSygnalu;
}

/** Zdarzenie rozpoznane jako wymagające decyzji Operatora, wraz z powodem, wagą i zdaniem dla Operatora. */
export interface DecyzjaCzekajaca {
  /** Tożsamość wpisu — ta sama, co w obserwacji. */
  klucz: string;
  idProcesu?: string;
  idOkna?: string;
  idSesji?: string;
  nazwa?: string;
  /** Dlaczego nakładka uznała, że proces czeka. */
  powod: PowodDecyzji;
  /** Czyja to ocena. */
  waga: WagaDecyzji;
  /** Rodzaj sugestii wg katalogu rodzajów sugestii AOD. */
  rodzaj: RodzajSugestii;
  /** Waga ujawnienia — jak głośno sugestia ma wejść do Operatora. */
  wagaUjawnienia: WagaUjawnienia;
  /** Czy to punkt decyzyjny wstrzymujący proces; wyjątek krytyczny ujawnia się mimo wyciszenia. */
  krytyczna: boolean;
  /** Zdanie „co czeka" pisane dla Operatora, nie dla dziennika. */
  zdanie: string;
  /** Stan procesu w chwili rozpoznania. */
  stan: ProgressStatus;
  etap?: string;
  numerEtapu?: number;
  liczbaEtapow?: number;
  ukonczenie?: number;
  /** Okno prowadzące bieg naprawczy; wyłącznie dla powodu `bieg-stanal`. */
  idOknaKoordynatora?: string;
  /** Powód zatrzymania biegu podany przez koordynatora. */
  powodZatrzymaniaBiegu?: string;
  /** Chwila, od której proces stoi — czas ostatniej zmiany telemetrii. */
  czekaOd: number;
  /** Skąd nakładka o tym wie. */
  zrodlo: ZrodloSygnalu;
}

/**
 * Rozstrzyga, czy obserwacja jest zdarzeniem wymagającym decyzji: `null` dla procesu,
 * który biegnie albo skończył.
 * @param teraz chwila odczytu w milisekundach epoki.
 * @param progBezRuchuMs próg powodu `bez-ruchu`; domyślnie {@link PROG_BEZ_RUCHU_MS}.
 */
export function rozpoznajDecyzje(
  obserwacja: ObserwacjaProcesu,
  teraz: number,
  progBezRuchuMs: number = PROG_BEZ_RUCHU_MS,
): DecyzjaCzekajaca | null {
  const powod = ustalPowod(obserwacja, teraz, progBezRuchuMs);
  if (powod === null) return null;

  return {
    klucz: obserwacja.klucz,
    idProcesu: obserwacja.idProcesu,
    idOkna: obserwacja.idOkna,
    idSesji: obserwacja.idSesji,
    nazwa: obserwacja.nazwa,
    powod,
    waga: WAGI[powod],
    rodzaj: RODZAJE[powod],
    wagaUjawnienia: WAGI_UJAWNIENIA[powod],
    krytyczna: KRYTYCZNE.has(powod),
    zdanie: zdanieDecyzji(obserwacja, powod, teraz, progBezRuchuMs),
    stan: obserwacja.stan,
    etap: obserwacja.etap,
    numerEtapu: obserwacja.numerEtapu,
    liczbaEtapow: obserwacja.liczbaEtapow,
    ukonczenie: obserwacja.ukonczenie,
    idOknaKoordynatora: obserwacja.bieg?.coordinatorWindowId,
    powodZatrzymaniaBiegu: obserwacja.bieg?.stopReason,
    czekaOd: obserwacja.zmienionoO,
    zrodlo: obserwacja.zrodlo,
  };
}

/** Waga przypisana każdemu powodowi: pewna, gdy stan rdzenia mówi to wprost, sporna, gdy to ocena nakładki. */
const WAGI: Readonly<Record<PowodDecyzji, WagaDecyzji>> = {
  [PowodDecyzji.BiegStanal]: WagaDecyzji.Pewna,
  [PowodDecyzji.Wstrzymany]: WagaDecyzji.Pewna,
  [PowodDecyzji.Zatrzymany]: WagaDecyzji.Pewna,
  [PowodDecyzji.Usterka]: WagaDecyzji.Sporna,
  [PowodDecyzji.BezRuchu]: WagaDecyzji.Sporna,
};

/**
 * Rodzaj sugestii przypisany każdemu powodowi: cztery powody, przy których proces
 * stoi, są wskazaniem problemu, a `bez-ruchu` jest kolejnym krokiem — nic się nie
 * zepsuło, czeka na ruszenie.
 */
const RODZAJE: Readonly<Record<PowodDecyzji, RodzajSugestii>> = {
  [PowodDecyzji.BiegStanal]: RodzajSugestii.Problem,
  [PowodDecyzji.Wstrzymany]: RodzajSugestii.Problem,
  [PowodDecyzji.Zatrzymany]: RodzajSugestii.Problem,
  [PowodDecyzji.Usterka]: RodzajSugestii.Problem,
  [PowodDecyzji.BezRuchu]: RodzajSugestii.KolejnyKrok,
};

/**
 * Waga ujawnienia każdego powodu: wysoka dla punktu decyzyjnego wstrzymującego
 * proces, średnia dla `bez-ruchu`.
 */
const WAGI_UJAWNIENIA: Readonly<Record<PowodDecyzji, WagaUjawnienia>> = {
  [PowodDecyzji.BiegStanal]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.Wstrzymany]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.Zatrzymany]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.Usterka]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.BezRuchu]: WagaUjawnienia.Srednia,
};

/**
 * Powody będące punktem decyzyjnym wstrzymującym proces, ujawniane mimo wyciszenia:
 * bieg naprawczy, który stanął, i proces wstrzymany, bo czekają wprost na
 * rozstrzygnięcie człowieka.
 */
const KRYTYCZNE: ReadonlySet<PowodDecyzji> = new Set([
  PowodDecyzji.BiegStanal,
  PowodDecyzji.Wstrzymany,
]);

/** Nazwa powodu widziana przez Operatora w oknie decyzji, zamiast wewnętrznego identyfikatora technicznego. */
export const NAZWY_POWODOW: Readonly<Record<PowodDecyzji, string>> = {
  [PowodDecyzji.BiegStanal]: 'bieg naprawczy stanął',
  [PowodDecyzji.Wstrzymany]: 'proces wstrzymany',
  [PowodDecyzji.Zatrzymany]: 'proces zatrzymany',
  [PowodDecyzji.Usterka]: 'proces stanął usterką',
  [PowodDecyzji.BezRuchu]: 'proces czeka bez ruchu',
};

/** Nazwa sygnału widziana przez Operatora, mówiąca skąd wiadomo, że proces wymaga jego uwagi w tej chwili. */
export const NAZWY_ZRODEL: Readonly<Record<ZrodloSygnalu, string>> = {
  [ZrodloSygnalu.Monitor]: 'odczyt nadrabiający monitor.status',
  [ZrodloSygnalu.Postep]: 'zdarzenie progress.changed',
  [ZrodloSygnalu.StanOkna]: 'zdarzenie window.state.changed',
};

/** Pierwszy pasujący powód rozpoznania albo `null`, gdy proces nie stoi i żaden powód nakładki nie zachodzi. */
function ustalPowod(
  obserwacja: ObserwacjaProcesu,
  teraz: number,
  progBezRuchuMs: number,
): PowodDecyzji | null {
  if (obserwacja.bieg?.stopped === true) return PowodDecyzji.BiegStanal;

  switch (obserwacja.stan) {
    case ProgressStatus.Paused:
      return PowodDecyzji.Wstrzymany;
    case ProgressStatus.Stopped:
      return PowodDecyzji.Zatrzymany;
    case ProgressStatus.Failed:
      return PowodDecyzji.Usterka;
    case ProgressStatus.Pending:
      return teraz - obserwacja.zmienionoO >= progBezRuchuMs ? PowodDecyzji.BezRuchu : null;
    default:
      return null;
  }
}

/** Składa zdanie „co czeka" dla Operatora — jedno pełne zdanie, bez skrótów rodem z dziennika technicznego rdzenia. */
function zdanieDecyzji(
  obserwacja: ObserwacjaProcesu,
  powod: PowodDecyzji,
  teraz: number,
  progBezRuchuMs: number,
): string {
  const byt = obserwacja.nazwa ?? obserwacja.idProcesu ?? obserwacja.idOkna ?? obserwacja.klucz;
  const etap = obserwacja.etap === undefined || obserwacja.etap === '' ? '' : ` na etapie „${obserwacja.etap}"`;

  switch (powod) {
    case PowodDecyzji.BiegStanal:
      return (
        `Bieg naprawczy okna ${obserwacja.bieg?.coordinatorWindowId ?? byt} stanął` +
        (obserwacja.bieg?.stopReason === undefined
          ? ' i czeka na rozstrzygnięcie.'
          : ` z powodem „${obserwacja.bieg.stopReason}".`)
      );
    case PowodDecyzji.Wstrzymany:
      return `Proces ${byt} jest wstrzymany${etap} i sam nie ruszy.`;
    case PowodDecyzji.Zatrzymany:
      return `Proces ${byt} został zatrzymany${etap}; wznowienie należy do Operatora.`;
    case PowodDecyzji.Usterka:
      return `Proces ${byt} stanął usterką${etap}. Czy usterka jest decyzją, rozstrzyga Właściciel — nakładka budzi, bo proces stoi.`;
    case PowodDecyzji.BezRuchu:
      return (
        `Proces ${byt} oczekuje${etap} od ${opiszOdstep(teraz - obserwacja.zmienionoO)} — ` +
        `dłużej niż próg ${opiszOdstep(progBezRuchuMs)} czasu oczekiwania zadania w kolejce.`
      );
  }
}

/**
 * Odstęp czasu opisany słowem — sekundy, minuty, godziny.
 *
 * Odstęp ujemny (zegar rdzenia przed zegarem klienta) opisujemy jako „przed
 * chwilą", a nie liczbą ujemną: zegar rozjechany jest stanem możliwym i nie
 * jest powodem do rysowania nieprawdy.
 */
export function opiszOdstep(ms: number): string {
  const sekundy = Math.floor(ms / 1000);
  if (sekundy <= 0) return 'przed chwilą';
  if (sekundy < 60) return `${sekundy} s`;
  const minuty = Math.floor(sekundy / 60);
  if (minuty < 60) return `${minuty} min`;
  const godziny = Math.floor(minuty / 60);
  return `${godziny} h ${minuty % 60} min`;
}

/** Sprowadza wiersz `monitor.status` do obserwacji o wspólnym kształcie, niezależnym od źródła jej sygnału. */
export function obserwacjaZMonitora(status: MonitorStatus): ObserwacjaProcesu {
  return {
    klucz: status.processId,
    idProcesu: status.processId,
    idOkna: status.windowId,
    idSesji: status.sessionId,
    nazwa: status.label,
    stan: status.status,
    etap: status.stage,
    numerEtapu: status.stageIndex,
    liczbaEtapow: status.stageCount,
    ukonczenie: status.completion,
    zmienionoO: status.updatedAt,
    zrodlo: ZrodloSygnalu.Monitor,
  };
}

/**
 * Sprowadza treść `progress.changed` do obserwacji: zdarzenie nie niesie chwili
 * zmiany, więc podaje ją wywołujący, a pola `loop` tu nie czytamy, bo jest martwe
 * w kontrakcie.
 */
export function obserwacjaZPostepu(
  tresc: ProgressChangedEvent,
  odebranoO: number,
): ObserwacjaProcesu {
  return {
    klucz: tresc.processId,
    idProcesu: tresc.processId,
    idOkna: tresc.windowId,
    stan: tresc.status,
    etap: tresc.stepLabel,
    numerEtapu: tresc.currentStep,
    liczbaEtapow: tresc.totalSteps,
    ukonczenie: tresc.percent,
    zmienionoO: odebranoO,
    zrodlo: ZrodloSygnalu.Postep,
  };
}

/**
 * Sprowadza treść `window.state.changed` do obserwacji: pole `state` jest odpisem
 * `WindowStateGetResponse`, jedynym żywym nośnikiem `LoopState`, więc klucz wpisu
 * jest tu oknem, nie procesem.
 */
export function obserwacjaZeStanuOkna(
  tresc: WindowStateChangedEvent,
  odebranoO: number,
): ObserwacjaProcesu | null {
  const idOkna = tresc.windowId;
  if (typeof idOkna !== 'string' || idOkna === '') return null;

  const stanOkna = odczytajStanOkna(tresc.state);
  if (stanOkna === null) return null;

  return {
    klucz: `${PRZEDROSTEK_KLUCZA_OKNA}${idOkna}`,
    idOkna,
    idSesji: stanOkna.idSesji,
    nazwa: stanOkna.nazwa,
    stan: stanOkna.stan,
    bieg: stanOkna.bieg,
    zmienionoO: stanOkna.bieg?.updatedAt ?? odebranoO,
    zrodlo: ZrodloSygnalu.StanOkna,
  };
}

/** To, co nakładce potrzebne z odpisu stanu okna — bieg naprawczy koordynatora, jego powód oraz identyfikator okna. */
interface StanOknaZeZdarzenia {
  stan: ProgressStatus;
  idSesji?: string;
  nazwa?: string;
  bieg?: LoopState;
}

/** Czyta odpis stanu okna z treści nieznanego kształtu; `null` wraca przy każdej niezgodności pola albo typu. */
function odczytajStanOkna(state: unknown): StanOknaZeZdarzenia | null {
  if (typeof state !== 'object' || state === null) return null;
  const zapis = state as Record<string, unknown>;

  const stan = zapis['processStatus'];
  if (!czyStanProcesu(stan)) return null;

  const okno = typeof zapis['window'] === 'object' && zapis['window'] !== null
    ? (zapis['window'] as Record<string, unknown>)
    : undefined;

  return {
    stan,
    idSesji: typeof okno?.['sessionId'] === 'string' ? (okno['sessionId'] as string) : undefined,
    nazwa: typeof okno?.['title'] === 'string' ? (okno['title'] as string) : undefined,
    bieg: czyStanBiegu(zapis['loop']) ? zapis['loop'] : undefined,
  };
}

/** Czy wartość jest jedną z sześciu wartości `ProgressStatus` — pending, running, paused, stopped, done, failed. */
function czyStanProcesu(wartosc: unknown): wartosc is ProgressStatus {
  return (
    typeof wartosc === 'string' &&
    (Object.values(ProgressStatus) as string[]).includes(wartosc)
  );
}

/**
 * Czy wartość jest stanem biegu naprawczego: sprawdza pola, na których stoi reguła,
 * nie komplet `LoopState`.
 */
function czyStanBiegu(wartosc: unknown): wartosc is LoopState {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const zapis = wartosc as Record<string, unknown>;
  return typeof zapis['coordinatorWindowId'] === 'string' && typeof zapis['stopped'] === 'boolean';
}
