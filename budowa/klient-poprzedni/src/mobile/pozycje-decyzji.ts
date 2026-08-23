import {
  LoopStopReason,
  ProgressStatus,
  QueueStatus,
  WindowRole,
  type MonitorStatus,
  type Queue,
  type Window,
  type WindowRoleAssignment,
} from '../../../shared/contract';
import type { BiegKoordynatora, ObrazInterwencji } from './zrodlo-interwencji';

/**
 * Skład wykazu „co czeka na twoją decyzję” — czysty rachunek na obrazie
 * interwencji, bez dokumentu i bez wywołań.
 *
 * Pozycja powstaje wyłącznie z tego, co rdzeń umie udowodnić odpowiedzią
 * komendy. Cztery dowody, cztery rodzaje pozycji:
 *
 *   `queue.list`       → kolejka w stanie `paused`      — praca stoi na kolejce
 *   `window.state.get` → `loop.stopped = true`          — pętla stoi, z powodem
 *   `monitor.status`   → proces `failed`                — krok padł
 *   `monitor.status`   → proces `paused`                — krok wstrzymany
 *
 * Kolejek eskalacji — przepływów wstrzymanych z pytaniem do człowieka — rdzeń
 * nie wystawia: nie ma na nie ani komendy odczytu, ani zdarzenia. Pulpit mówi
 * o tym wprost (`mission-control/pas-decyzji.ts`), a warstwa mobilna nie
 * zamalowuje braku pozycjami zmyślonymi. Gdy rdzeń wykaz eskalacji wystawi,
 * wejdzie on portem `port-kolejki-decyzji.ts`.
 *
 * Kontekst decyzji jest częścią pozycji: Operator otwiera telefon na minutę
 * i musi wiedzieć, na czym praca stoi, zanim cokolwiek naciśnie. Każde zdanie
 * kontekstu niesie nazwę komendy, z której przyszło, więc da się sprawdzić
 * jego źródło.
 */

/** Rodzaj pozycji — zarazem dowód, z którego powstała. */
export const RodzajPozycji = {
  KolejkaWstrzymana: 'kolejka-wstrzymana',
  PetlaZatrzymana: 'petla-zatrzymana',
  ProcesBledny: 'proces-bledny',
  ProcesWstrzymany: 'proces-wstrzymany',
  /** Pozycja wniesiona portem kolejki decyzji — nie powstaje w tym pliku. */
  Eskalacja: 'eskalacja',
} as const;
export type RodzajPozycji = (typeof RodzajPozycji)[keyof typeof RodzajPozycji];

/** Jedna rzecz czekająca na decyzję Operatora wraz z kontekstem i uchwytami. */
export interface PozycjaDecyzji {
  /** Klucz stały między odczytami — rodzaj i identyfikator bytu. */
  id: string;
  rodzaj: RodzajPozycji;
  /** Zdanie tytułowe karty; jedna linia, bez skrótów. */
  naglowek: string;
  /** Zdania kontekstu; każde kończy się nazwą komendy, z której pochodzi. */
  kontekst: readonly string[];
  /** Kolejka, na której stoi praca; obecność otwiera drogi kolejkowe. */
  queueId?: string;
  /** Okno, którego dotyczy pozycja. */
  windowId?: string;
  /** Okno koordynatora — adresat nastaw i przejęcia sterowania. */
  koordynatorWindowId?: string;
  sessionId?: string;
  processId?: string;
  /** Chwila ostatniej zmiany stanu w milisekundach epoki — „stoi od”. */
  stoiOd?: number;
  /** Nazwa komendy, która pozycję udowodniła. */
  zrodlo: string;
}

/** Nagłówek ekranu — „co się w ogóle dzieje”, zanim spojrzysz na wykaz. */
export interface ObrazSkrocony {
  /** Zdania stanu; każde z nazwą komendy albo z treścią odmowy. */
  zdania: readonly string[];
  /** Ile odczytów obrazu odmówił rdzeń — wprost, bez chowania. */
  odmowy: number;
}

const ETYKIETY_STANU: Readonly<Record<string, string>> = {
  [ProgressStatus.Pending]: 'oczekuje',
  [ProgressStatus.Running]: 'w biegu',
  [ProgressStatus.Paused]: 'wstrzymany',
  [ProgressStatus.Stopped]: 'zatrzymany',
  [ProgressStatus.Done]: 'zakończony',
  [ProgressStatus.Failed]: 'zakończony błędem',
};

const POWODY_ZATRZYMANIA: Readonly<Record<string, string>> = {
  [LoopStopReason.NoProgress]: 'wykonawca powtarza się bez postępu',
  [LoopStopReason.Manual]: 'zatrzymał ją Operator',
  [LoopStopReason.Failure]: 'obiegu nie udało się rozpocząć',
};

/** Zdanie „stoi od…”; bez znanego czasu — bez zdania, zamiast zmyślonego. */
export function zdanieCzasu(stoiOd: number | undefined, teraz: number): string | null {
  if (stoiOd === undefined || !Number.isFinite(stoiOd) || stoiOd <= 0) return null;
  const minut = Math.floor((teraz - stoiOd) / 60_000);
  if (minut < 0) return null;
  if (minut === 0) return 'stoi krócej niż minutę';
  if (minut < 60) return `stoi od ${minut} min`;
  const godzin = Math.floor(minut / 60);
  return `stoi od ${godzin} godz. ${minut % 60} min`;
}

/**
 * Składa wykaz pozycji z obrazu interwencji.
 *
 * Kolejność: kolejki, pętle, procesy błędne, procesy wstrzymane — od tego,
 * co blokuje najwięcej pracy, do tego, co blokuje jeden krok.
 */
export function zlozPozycjeDecyzji(obraz: ObrazInterwencji, teraz: number): PozycjaDecyzji[] {
  const okna = new Map((obraz.okna.wynik ?? []).map((o) => [o.id, o]));
  const procesy = obraz.procesy.wynik ?? [];
  const pozycje: PozycjaDecyzji[] = [];

  for (const kolejka of obraz.kolejki.wynik ?? []) {
    if (kolejka.status !== QueueStatus.Paused) continue;
    pozycje.push(pozycjaKolejki(kolejka, okna, obraz.role.wynik ?? [], procesy, teraz));
  }

  for (const bieg of obraz.biegi) {
    const pozycja = pozycjaPetli(bieg, okna, obraz.kolejki.wynik ?? [], teraz);
    if (pozycja !== null) pozycje.push(pozycja);
  }

  for (const proces of procesy) {
    if (proces.status === ProgressStatus.Failed) {
      pozycje.push(pozycjaProcesu(proces, RodzajPozycji.ProcesBledny, okna, teraz));
    }
  }
  for (const proces of procesy) {
    if (proces.status === ProgressStatus.Paused) {
      pozycje.push(pozycjaProcesu(proces, RodzajPozycji.ProcesWstrzymany, okna, teraz));
    }
  }

  return pozycje;
}

/** Zdania nagłówka ekranu wraz z liczbą odmów — stan czytany, nie zgadywany. */
export function skrocObraz(obraz: ObrazInterwencji): ObrazSkrocony {
  const zdania: string[] = [];
  let odmowy = 0;

  if (obraz.procesy.udany) {
    const wBiegu = (obraz.procesy.wynik ?? []).filter(
      (p) => p.status === ProgressStatus.Running,
    ).length;
    zdania.push(
      `Procesy w biegu: ${wBiegu} z ${(obraz.procesy.wynik ?? []).length} (monitor.status).`,
    );
  } else {
    odmowy += 1;
    zdania.push('Telemetrii procesów rdzeń dziś odmówił (monitor.status).');
  }

  if (obraz.kolejki.udany) {
    const kolejki = obraz.kolejki.wynik ?? [];
    const wstrzymane = kolejki.filter((k) => k.status === QueueStatus.Paused).length;
    zdania.push(`Kolejki: ${kolejki.length}, w tym wstrzymanych ${wstrzymane} (queue.list).`);
  } else {
    odmowy += 1;
    zdania.push('Wykazu kolejek rdzeń dziś odmówił (queue.list).');
  }

  if (obraz.okna.udany) {
    zdania.push(`Okna komunikacji: ${(obraz.okna.wynik ?? []).length} (window.list).`);
  } else {
    odmowy += 1;
    zdania.push('Wykazu okien rdzeń dziś odmówił (window.list).');
  }

  const odmowyBiegow = obraz.biegi.filter((b) => b.odmowa !== undefined).length;
  if (odmowyBiegow > 0) {
    odmowy += odmowyBiegow;
    zdania.push(`Stanu ${odmowyBiegow} okna koordynatora rdzeń nie oddał (window.state.get).`);
  }

  return { zdania, odmowy };
}

function pozycjaKolejki(
  kolejka: Queue,
  okna: ReadonlyMap<string, Window>,
  role: readonly WindowRoleAssignment[],
  procesy: readonly MonitorStatus[],
  teraz: number,
): PozycjaDecyzji {
  const windowId = kolejka.windowIds?.[0];
  const koordynator = wskazKoordynatora(windowId, okna, role);
  const proces = procesy.find((p) => p.windowId !== undefined && p.windowId === windowId);

  const kontekst: string[] = [
    `Kolejka jest wstrzymana; obiegów naprawczych: ${kolejka.cycle ?? 0} (queue.list).`,
  ];
  const czas = zdanieCzasu(kolejka.updatedAt, teraz);
  if (czas !== null) kontekst.push(`${czas} (queue.list).`);
  if (proces !== undefined) kontekst.push(zdanieEtapu(proces));
  if (koordynator !== undefined) {
    kontekst.push(`Prowadzi ją okno koordynatora ${koordynator} (role.list / window.list).`);
  }

  return {
    id: `${RodzajPozycji.KolejkaWstrzymana}:${kolejka.id}`,
    rodzaj: RodzajPozycji.KolejkaWstrzymana,
    naglowek: `Kolejka „${kolejka.name ?? kolejka.id}” stoi wstrzymana`,
    kontekst,
    queueId: kolejka.id,
    windowId,
    koordynatorWindowId: koordynator,
    sessionId: kolejka.sessionId,
    processId: proces?.processId,
    stoiOd: kolejka.updatedAt,
    zrodlo: 'queue.list',
  };
}

function pozycjaPetli(
  bieg: BiegKoordynatora,
  okna: ReadonlyMap<string, Window>,
  kolejki: readonly Queue[],
  teraz: number,
): PozycjaDecyzji | null {
  const petla = bieg.loop;
  if (petla === undefined || !petla.stopped) return null;

  const okno = okna.get(bieg.windowId);
  const kolejka = kolejki.find((k) => (k.windowIds ?? []).includes(bieg.windowId));
  const powod = petla.stopReason === undefined ? null : POWODY_ZATRZYMANIA[petla.stopReason];

  const kontekst: string[] = [
    powod === null
      ? 'Bieg naprawczy stoi; powodu rdzeń nie podał (window.state.get).'
      : `Bieg naprawczy stoi — ${powod} (window.state.get).`,
    `Obiegów: ${petla.loops}, bez postępu ${petla.loopsWithoutProgress} przy progu ${petla.threshold} (window.state.get).`,
  ];
  const czas = zdanieCzasu(petla.updatedAt, teraz);
  if (czas !== null) kontekst.push(`${czas} (window.state.get).`);
  if (petla.lastTurnReason !== undefined) {
    kontekst.push(`Ostatnia tura wykonawcy skończyła się: ${petla.lastTurnReason} (window.state.get).`);
  }

  return {
    id: `${RodzajPozycji.PetlaZatrzymana}:${bieg.windowId}`,
    rodzaj: RodzajPozycji.PetlaZatrzymana,
    naglowek: `Pętla koordynatora „${okno?.title ?? bieg.windowId}” zatrzymana`,
    kontekst,
    queueId: kolejka?.id,
    windowId: bieg.windowId,
    koordynatorWindowId: petla.coordinatorWindowId,
    sessionId: okno?.sessionId ?? kolejka?.sessionId,
    stoiOd: petla.updatedAt,
    zrodlo: 'window.state.get',
  };
}

function pozycjaProcesu(
  proces: MonitorStatus,
  rodzaj: RodzajPozycji,
  okna: ReadonlyMap<string, Window>,
  teraz: number,
): PozycjaDecyzji {
  const okno = proces.windowId === undefined ? undefined : okna.get(proces.windowId);
  const kontekst: string[] = [zdanieEtapu(proces)];
  const czas = zdanieCzasu(proces.updatedAt, teraz);
  if (czas !== null) kontekst.push(`${czas} (monitor.status).`);
  if (okno !== undefined) {
    kontekst.push(`Okno ${okno.title ?? okno.id}, tryb uprawnień ${okno.permissionMode} (window.list).`);
  }

  return {
    id: `${rodzaj}:${proces.processId}`,
    rodzaj,
    naglowek: `Proces „${proces.label ?? proces.processId}” — ${ETYKIETY_STANU[proces.status] ?? proces.status}`,
    kontekst,
    windowId: proces.windowId,
    koordynatorWindowId: okno?.coordinatorWindowId ?? okno?.id,
    sessionId: proces.sessionId ?? okno?.sessionId,
    processId: proces.processId,
    stoiOd: proces.updatedAt,
    zrodlo: 'monitor.status',
  };
}

/** Zdanie „na czym proces stoi” — etap, numer etapu, procent, obieg. */
function zdanieEtapu(proces: MonitorStatus): string {
  const czesci: string[] = [];
  if (proces.stage !== undefined) czesci.push(`etap „${proces.stage}”`);
  if (proces.stageIndex !== undefined) {
    czesci.push(
      proces.stageCount === undefined
        ? `krok ${proces.stageIndex}`
        : `krok ${proces.stageIndex} z ${proces.stageCount}`,
    );
  }
  if (proces.completion !== undefined) czesci.push(`${proces.completion}% ukończenia`);
  if (proces.cycle !== undefined) czesci.push(`obieg ${proces.cycle}`);
  if (czesci.length === 0) return 'Etapu procesu rdzeń nie podał (monitor.status).';
  return `${czesci.join(', ')} (monitor.status).`;
}

/** Okno koordynatora dla okna pozycji — z `window.list`, a gdy brak, z `role.list`. */
function wskazKoordynatora(
  windowId: string | undefined,
  okna: ReadonlyMap<string, Window>,
  role: readonly WindowRoleAssignment[],
): string | undefined {
  if (windowId === undefined) return undefined;
  const okno = okna.get(windowId);
  if (okno !== undefined) {
    if (okno.windowRole === WindowRole.Coordinator) return okno.id;
    if (okno.coordinatorWindowId !== undefined) return okno.coordinatorWindowId;
  }
  const nadanie = role.find((r) => r.windowId === windowId);
  if (nadanie === undefined) return undefined;
  if (nadanie.role === WindowRole.Coordinator) return nadanie.windowId;
  return nadanie.coordinatorWindowId;
}
