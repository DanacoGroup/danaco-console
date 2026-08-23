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
 * Rozpoznanie zdarzenia wymagającego decyzji — czysta reguła nakładki.
 *
 * Plik nie dotyka dokumentu i nie zna kanału. Bierze trzy sygnały, które rdzeń
 * wypuszcza, i oddaje jedno rozstrzygnięcie: czy przy tym procesie Operator ma
 * być obudzony, czy nie.
 *
 * Pojęcia „zdarzenie wymagające decyzji" nie ma w kontrakcie ani w rdzeniu: nie
 * ma rodziny `decision.*`, nie ma pola `awaitingDecision`, a `ProgressStatus`
 * ma sześć wartości (pending, running, paused, stopped, done, failed) i żadna
 * z nich nie znaczy „czekam na Ciebie". Regułę wywodzi więc nakładka i tak też
 * jest nazwana Operatorowi w oknie (`sekcja-decyzji.ts`).
 *
 * Sedno reguły: budzi to, przy czym proces stoi — `running` i `done` nie budzą
 * nigdy. Powody dzielą się na pewne (stan rdzenia mówi wprost, że proces stoi)
 * i sporne (to ocena nakładki). Sporne nie są ukrywane, tylko oznaczone własną
 * wagą, żeby Operator wiedział, czyja to ocena.
 *
 *   • `bieg-stanal`  (pewna)  — `loop.stopped = true`: koordynator sam ogłosił,
 *     że bieg naprawczy stanął, i podał powód. Sprawdzany pierwszy, bo jest
 *     najbardziej szczegółowy: niesie okno koordynatora i przyczynę.
 *   • `wstrzymany`   (pewna)  — `status = paused`.
 *   • `zatrzymany`   (pewna)  — `status = stopped`.
 *   • `usterka`      (sporna) — `status = failed`; nakładka budzi, bo proces
 *     stoi, a nikt inny go nie ruszy.
 *   • `bez-ruchu`    (sporna) — `status = pending` dłużej niż próg. Pozycja,
 *     która nigdy nie ruszyła, jest nieodróżnialna od pozycji zapomnianej.
 *
 * Źródło sygnału jest częścią odpowiedzi: odczyt nadrabiający i zdarzenie na
 * żywo mają różną świeżość. Pole `loop` w `progress.changed` jest martwe —
 * kontrakt je obiecuje, `telemetria_proces.go` nigdy go nie wypełnia. Stan
 * biegu naprawczego dojeżdża wyłącznie zdarzeniem `window.state.changed`,
 * którego treść jest odpisem `WindowStateGetResponse`, i stamtąd ta reguła
 * bierze `loop`.
 */

/** Sygnał, z którego rozpoznanie wzięło obserwację. */
export const ZrodloSygnalu = {
  /** `monitor.status` — odczyt nadrabiający; działa bez argumentów. */
  Monitor: 'monitor.status',
  /** `progress.changed` — zdarzenie rozgłaszane do każdego połączenia. */
  Postep: 'progress.changed',
  /** `window.state.changed` — jedyny nośnik stanu biegu naprawczego. */
  StanOkna: 'window.state.changed',
} as const;
export type ZrodloSygnalu = (typeof ZrodloSygnalu)[keyof typeof ZrodloSygnalu];

/** Powód, dla którego nakładka uznała, że proces czeka na Operatora. */
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

/** Czyja to ocena: stanu rdzenia czy nakładki. */
export const WagaDecyzji = {
  /** Proces stoi — stan rdzenia mówi to wprost. */
  Pewna: 'pewna',
  /** To ocena nakładki, a nie stan orzeczony przez rdzeń. */
  Sporna: 'sporna',
} as const;
export type WagaDecyzji = (typeof WagaDecyzji)[keyof typeof WagaDecyzji];

/**
 * Próg braku ruchu dla stanu `pending`.
 *
 * Wartość pochodzi z opracowania: „Próg czasu oczekiwania zadania w kolejce —
 * 15 minut. Przekroczenie tworzy sugestię klasy «stan kolejki zadań»"
 * (`docs/funkcje-globalne/always-on-display.md`, rozdz. 3.4). Pozycja `pending`
 * to właśnie zadanie oczekujące w kolejce, więc reguła nakładki bierze próg
 * stamtąd, zamiast stanowić własny.
 *
 * Kontrakt progu nie niesie; `LoopState` ma własny `threshold`, ale liczy
 * obiegi, nie czas, i dotyczy wyłącznie biegu naprawczego. Próg jest argumentem
 * reguły, więc jego zmiana dotyka jednej stałej, a nie kodu.
 */
export const PROG_BEZ_RUCHU_MS = PROGI_AOD.oczekiwanieWKolejceMs;

/** Przedrostek klucza wpisu, którego sygnał nie niósł identyfikatora procesu. */
export const PRZEDROSTEK_KLUCZA_OKNA = 'okno:';

/** Jedna obserwacja procesu sprowadzona do wspólnego kształtu. */
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

/** Zdarzenie rozpoznane jako wymagające decyzji Operatora. */
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
  /** Rodzaj sugestii wg katalogu opracowania (rozdz. 4). */
  rodzaj: RodzajSugestii;
  /** Waga ujawnienia wg rozdz. 4.5 — jak głośno sugestia ma wejść. */
  wagaUjawnienia: WagaUjawnienia;
  /**
   * Czy to punkt decyzyjny wstrzymujący proces.
   *
   * Wyjątek wagi krytycznej z rozdz. 3.5: taka sugestia ujawnia się mimo
   * wyciszenia — plakietką, bez dymka.
   */
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
 * Rozstrzyga, czy obserwacja jest zdarzeniem wymagającym decyzji.
 *
 * Zwraca `null` dla procesu, który biegnie albo skończył — i to jest większość
 * ruchu. Nakładka, która budzi przy każdym zdarzeniu, nie budzi przy żadnym.
 *
 * @param teraz chwila odczytu w milisekundach epoki — podawana z zewnątrz,
 *   żeby reguła była sprawdzalna bez zegara.
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

/** Waga przypisana każdemu powodowi. */
const WAGI: Readonly<Record<PowodDecyzji, WagaDecyzji>> = {
  [PowodDecyzji.BiegStanal]: WagaDecyzji.Pewna,
  [PowodDecyzji.Wstrzymany]: WagaDecyzji.Pewna,
  [PowodDecyzji.Zatrzymany]: WagaDecyzji.Pewna,
  [PowodDecyzji.Usterka]: WagaDecyzji.Sporna,
  [PowodDecyzji.BezRuchu]: WagaDecyzji.Sporna,
};

/**
 * Rodzaj sugestii przypisany każdemu powodowi — katalog rozdz. 4 opracowania.
 *
 * Cztery powody mówiące, że proces stoi, są wskazaniem problemu (rozdz. 4.3:
 * „zadanie w stanie błędu, kolejka zatrzymana, pętla przerwana"). Powód
 * `bez-ruchu` jest kolejnym krokiem (rozdz. 3.2, klasa „stan kolejki zadań":
 * zadanie oczekujące dłużej niż próg czasu) — nic się nie zepsuło, czeka na
 * ruszenie.
 */
const RODZAJE: Readonly<Record<PowodDecyzji, RodzajSugestii>> = {
  [PowodDecyzji.BiegStanal]: RodzajSugestii.Problem,
  [PowodDecyzji.Wstrzymany]: RodzajSugestii.Problem,
  [PowodDecyzji.Zatrzymany]: RodzajSugestii.Problem,
  [PowodDecyzji.Usterka]: RodzajSugestii.Problem,
  [PowodDecyzji.BezRuchu]: RodzajSugestii.KolejnyKrok,
};

/**
 * Waga ujawnienia każdego powodu — rozdz. 4.5 opracowania.
 *
 * Wysoka: „punkt decyzyjny wstrzymujący proces, kolejka zatrzymana". Średnia:
 * „przekroczony czas oczekiwania zadania" — i to jest dokładnie `bez-ruchu`.
 */
const WAGI_UJAWNIENIA: Readonly<Record<PowodDecyzji, WagaUjawnienia>> = {
  [PowodDecyzji.BiegStanal]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.Wstrzymany]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.Zatrzymany]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.Usterka]: WagaUjawnienia.Wysoka,
  [PowodDecyzji.BezRuchu]: WagaUjawnienia.Srednia,
};

/**
 * Powody będące punktem decyzyjnym wstrzymującym proces — wyjątek wagi
 * krytycznej z rozdz. 3.5, ujawniany mimo wyciszenia.
 *
 * Bieg naprawczy, który stanął, i proces wstrzymany zatrzymują pracę i czekają
 * wprost na rozstrzygnięcie człowieka. Zatrzymanie, usterka i brak ruchu
 * przez wyciszenie przeczekają.
 */
const KRYTYCZNE: ReadonlySet<PowodDecyzji> = new Set([
  PowodDecyzji.BiegStanal,
  PowodDecyzji.Wstrzymany,
]);

/** Nazwa powodu widziana przez Operatora. */
export const NAZWY_POWODOW: Readonly<Record<PowodDecyzji, string>> = {
  [PowodDecyzji.BiegStanal]: 'bieg naprawczy stanął',
  [PowodDecyzji.Wstrzymany]: 'proces wstrzymany',
  [PowodDecyzji.Zatrzymany]: 'proces zatrzymany',
  [PowodDecyzji.Usterka]: 'proces stanął usterką',
  [PowodDecyzji.BezRuchu]: 'proces czeka bez ruchu',
};

/** Nazwa sygnału widziana przez Operatora — „skąd wiadomo". */
export const NAZWY_ZRODEL: Readonly<Record<ZrodloSygnalu, string>> = {
  [ZrodloSygnalu.Monitor]: 'odczyt nadrabiający monitor.status',
  [ZrodloSygnalu.Postep]: 'zdarzenie progress.changed',
  [ZrodloSygnalu.StanOkna]: 'zdarzenie window.state.changed',
};

/** Pierwszy pasujący powód albo `null`, gdy proces nie stoi. */
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

/** Składa zdanie „co czeka" — jedno zdanie, bez skrótów rodem z dziennika. */
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

/** Sprowadza wiersz `monitor.status` do obserwacji. */
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
 * Sprowadza treść `progress.changed` do obserwacji.
 *
 * Zdarzenie NIE NIESIE chwili zmiany, więc chwilę podaje wywołujący — jest nią
 * moment odebrania zdarzenia. Nie niesie też sesji ani nazwy: telemetria zna
 * proces i okno, a resztę dokłada odczyt nadrabiający.
 *
 * Pola `loop` tu nie czytamy, choć kontrakt je obiecuje: `telemetria_proces.go`
 * nigdy go nie wypełnia. Czytanie martwego pola udawałoby drogę, której nie ma —
 * bieg naprawczy przychodzi zdarzeniem `window.state.changed`.
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
 * Sprowadza treść `window.state.changed` do obserwacji.
 *
 * Treść pola `state` jest odpisem `WindowStateGetResponse` złożonym przez
 * `core/stan_sesji_nadzor.go` — stąd stan procesu okna i JEDYNY żywy nośnik
 * `LoopState`. Kontrakt opisuje `state` jako `unknown`, więc pole nieznanego
 * kształtu daje `null`, a nie wyjątek.
 *
 * Klucz wpisu jest tu oknem, nie procesem: `WindowStateGetResponse` nie niesie
 * identyfikatora procesu i nakładka nie ma go skąd wziąć.
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

/** To, co nakładce potrzebne z odpisu stanu okna. */
interface StanOknaZeZdarzenia {
  stan: ProgressStatus;
  idSesji?: string;
  nazwa?: string;
  bieg?: LoopState;
}

/** Czyta odpis stanu okna z treści nieznanego kształtu; `null` przy każdej niezgodności. */
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

/** Czy wartość jest jedną z sześciu wartości `ProgressStatus`. */
function czyStanProcesu(wartosc: unknown): wartosc is ProgressStatus {
  return (
    typeof wartosc === 'string' &&
    (Object.values(ProgressStatus) as string[]).includes(wartosc)
  );
}

/**
 * Czy wartość jest stanem biegu naprawczego.
 *
 * Sprawdzamy pola, na których stoi reguła (`coordinatorWindowId`, `stopped`),
 * a nie komplet `LoopState` — odpis uboższy o pole nieużywane jest nadal
 * odpisem prawdziwym, a odrzucenie go zgubiłoby powód `bieg-stanal`.
 */
function czyStanBiegu(wartosc: unknown): wartosc is LoopState {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const zapis = wartosc as Record<string, unknown>;
  return typeof zapis['coordinatorWindowId'] === 'string' && typeof zapis['stopped'] === 'boolean';
}
