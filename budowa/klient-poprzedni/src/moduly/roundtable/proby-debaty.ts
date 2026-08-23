import {
  ChangeKind,
  ChunkKind,
  EventType,
  RoundtableFormat,
  RoundtableTurnStatus,
  type Channel,
  type Envelope,
  type RoundtableDebateChangedEvent,
  type RoundtableParticipant,
  type RoundtableStatement,
  type RoundtableTurn,
  type StreamChunkEvent,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import type { RejestrKanalow } from '../../sterowanie/rejestr-kanalow';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Zastawy debaty dla sprawdzianów modułu — jedno miejsce, w którym powstają
 * kształty kontraktu i podstawione zależności.
 *
 * Plik sam sprawdzianem nie jest i nie ma w nazwie `.test`: sprawdziany modułu
 * sięgają po te same kształty kontraktu (`RoundtableTurn`, `RoundtableStatement`,
 * `RoundtableParticipant`, `StreamChunkEvent`) i te same podstawienia
 * (`ZrodloRoundtable`, `RejestrKanalow`), a kopia zastawu w każdym z nich by się
 * rozjechała.
 *
 * Każda wytwórnia zwraca typ kontraktu w całości, więc dopisanie pola
 * obowiązkowego do `contract.ts` przerywa kompilację tutaj, a nie w każdym
 * sprawdzianie z osobna.
 */

/** Tura debaty; pola nadpisywalne, reszta wypełniona wartością sensowną. */
export function probnaTura(nadpisania: Partial<RoundtableTurn> = {}): RoundtableTurn {
  return {
    id: 'tura-1',
    windowId: 'okno-debaty',
    index: 1,
    format: RoundtableFormat.Free,
    status: RoundtableTurnStatus.Open,
    startedAt: 1_700_000_000_000,
    ...nadpisania,
  };
}

/** Wypowiedź uczestnika; treść pusta jest domyślna, bo rdzeń rozgłasza
 *  `created` właśnie z pustą treścią. */
export function probnaWypowiedz(
  nadpisania: Partial<RoundtableStatement> = {},
): RoundtableStatement {
  return {
    id: 'wypow-1',
    turnId: 'tura-1',
    participantId: 'uczest-1',
    content: '',
    createdAt: 1_700_000_000_000,
    ...nadpisania,
  };
}

/** Uczestnik debaty. */
export function probnyUczestnik(
  nadpisania: Partial<RoundtableParticipant> = {},
): RoundtableParticipant {
  return {
    id: 'uczest-1',
    windowId: 'okno-debaty',
    channelId: 'kan-1',
    ...nadpisania,
  };
}

/** Kanał modelu w rejestrze rdzenia. */
export function probnyKanal(nadpisania: Partial<Channel> = {}): Channel {
  return {
    id: 'kan-1',
    name: 'Kanał pierwszy',
    kind: 'model',
    enabled: true,
    ...nadpisania,
  } as Channel;
}

/** Źródło komend obszaru z ręcznie sterowanym zdarzeniem `debate.changed`. */
export interface ProbneZrodlo {
  zrodlo: ZrodloRoundtable;
  /** Rozgłasza przyrost debaty do wszystkich subskrybentów. */
  oglos(tresc: RoundtableDebateChangedEvent): void;
  /** Ilu słuchaczy zdarzenia żyje w tej chwili — mierzy zdjęcie subskrypcji. */
  sluchaczy(): number;
}

export function probneZrodlo(): ProbneZrodlo {
  const sluchacze = new Set<(tresc: RoundtableDebateChangedEvent) => void>();
  const odmowa = (): never => {
    throw new Error('Sprawdzian nie przewidział wywołania komendy przez ten byt.');
  };
  const zrodlo: ZrodloRoundtable = {
    dodajModel: odmowa,
    uruchomDebate: odmowa,
    moderuj: odmowa,
    stanowisko: odmowa,
    naZmianeDebaty(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
  return {
    zrodlo,
    oglos: (tresc) => {
      for (const sluchacz of [...sluchacze]) sluchacz(tresc);
    },
    sluchaczy: () => sluchacze.size,
  };
}

/** Przyrost debaty w kształcie, w jakim nadaje go rdzeń. */
export function przyrost(
  tura: RoundtableTurn,
  wypowiedz?: RoundtableStatement,
  rodzaj: ChangeKind = ChangeKind.Updated,
): RoundtableDebateChangedEvent {
  return wypowiedz === undefined
    ? { change: rodzaj, turn: tura }
    : { change: rodzaj, turn: tura, statement: wypowiedz };
}

/** Rejestr kanałów podstawiony — bez rdzenia, z ręcznie ustawianym wykazem. */
export interface ProbnyRejestr {
  rejestr: RejestrKanalow;
  /** Podstawia wykaz i powiadamia subskrybentów, tak jak zrobiłby to rdzeń. */
  ustaw(kanaly: Channel[]): void;
  /** Ile razy okno zamówiło `channel.list`. */
  zamowien(): number;
}

export function probnyRejestr(kanaly: Channel[] = []): ProbnyRejestr {
  const sluchacze = new Set<(lista: Channel[]) => void>();
  let wykaz = kanaly;
  let odpowiedziano = kanaly.length > 0;
  let zamowien = 0;
  const rejestr: RejestrKanalow = {
    kanaly: () => wykaz,
    odpowiedzOtrzymana: () => odpowiedziano,
    odswiez: () => {
      zamowien += 1;
    },
    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
  return {
    rejestr,
    ustaw(lista) {
      wykaz = lista;
      odpowiedziano = true;
      for (const sluchacz of [...sluchacze]) sluchacz(wykaz);
    },
    zamowien: () => zamowien,
  };
}

/** Kanał podstawiony — obsługuje wyłącznie subskrypcję zdarzeń i nadanie ramki. */
export interface ProbnyKanalZdarzen {
  kanal: Kanal;
  /** Nadaje fragment strumienia tak, jak zrobiłby to rdzeń. */
  fragment(tresc: StreamChunkEvent, ostatni?: boolean): void;
  /** Nadaje dowolne zdarzenie kontraktu. */
  zdarzenie(typ: EventType, tresc: unknown, koperta?: Partial<Envelope>): void;
  /** Ile subskrypcji zdarzeń żyje w tej chwili. */
  subskrypcji(): number;
}

/**
 * @param odpowiedzi treść oddawana na komendę o podanej nazwie. Komenda bez
 *   wpisu nie dostaje odpowiedzi w ogóle — to stan „rdzeń jeszcze nie
 *   odpowiedział", a nie odmowa, i okna mają go rozróżniać.
 */
export function probnyKanalZdarzen(
  odpowiedzi: Partial<Record<string, unknown>> = {},
): ProbnyKanalZdarzen {
  const sluchacze = new Map<EventType, Set<(tresc: never, koperta: Envelope) => void>>();

  function subskrybuj(typ: EventType, sluchacz: (tresc: never, koperta: Envelope) => void): Odsubskrybuj {
    const zbior = sluchacze.get(typ) ?? new Set();
    zbior.add(sluchacz);
    sluchacze.set(typ, zbior);
    return () => zbior.delete(sluchacz);
  }

  const kanal = {
    wyslij: (
      komenda: string,
      _zadanie: unknown,
      przyWyniku?: (wynik: { udany: boolean; wynik?: unknown }) => void,
    ) => {
      const wynik = odpowiedzi[komenda];
      if (przyWyniku !== undefined && wynik !== undefined) przyWyniku({ udany: true, wynik });
      return 'zad-1';
    },
    naZdarzenie: (typ: EventType, sluchacz: (tresc: never, koperta: Envelope) => void) =>
      subskrybuj(typ, sluchacz),
    naDowolny: () => () => undefined,
    sesja: () => ({ id: () => 'sesja-1' }),
    dziennikNieznanych: () => ({ zapisz: () => undefined, wykaz: () => [] }),
  } as unknown as Kanal;

  function zdarzenie(typ: EventType, tresc: unknown, koperta: Partial<Envelope> = {}): void {
    const pelna: Envelope = {
      type: typ,
      id: 'zad-1',
      timestamp: 1_700_000_000_000,
      ...koperta,
    };
    for (const sluchacz of [...(sluchacze.get(typ) ?? [])]) {
      (sluchacz as (tresc: unknown, koperta: Envelope) => void)(tresc, pelna);
    }
  }

  return {
    kanal,
    fragment: (tresc, ostatni = false) =>
      zdarzenie(EventType.StreamChunk, tresc, ostatni ? { done: true } : {}),
    zdarzenie,
    subskrypcji: () => [...sluchacze.values()].reduce((suma, zbior) => suma + zbior.size, 0),
  };
}

/** Fragment strumienia w kształcie kontraktu. */
export function probnyFragment(nadpisania: Partial<StreamChunkEvent> = {}): StreamChunkEvent {
  return {
    windowId: 'okno-debaty',
    messageId: 'wypow-1',
    kind: ChunkKind.Text,
    text: '',
    ...nadpisania,
  };
}
