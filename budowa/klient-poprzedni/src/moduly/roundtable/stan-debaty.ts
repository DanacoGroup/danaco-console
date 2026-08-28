import type {
  Channel,
  RoundtableParticipant,
  RoundtableStatement,
  RoundtableTurn,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { nazwaKanalu, type RejestrKanalow } from '../../sterowanie/rejestr-kanalow';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Debata bieżąca modułu — jeden stan wspólny dla sześciu okien, kluczowany po tożsamości
 * uczestnika, nie po kanale, z wykazem kanałów z rejestru rdzenia.
 */
export interface StanDebaty {
  /** Okno debaty; wymagane przez każdą komendę obszaru, którą moduł wywołuje. */
  okno(): string;
  /** Przestawia moduł na inne okno debaty i czyści byty poprzedniej. */
  ustawOkno(idOkna: string): void;

  // Uczestnicy debaty; kolejność jest kolejnością głosu, gdy rdzeń ją podał, inaczej dodania.
  uczestnicy(): RoundtableParticipant[];
  /** Uczestnik o podanej tożsamości; `null`, gdy nieznany. */
  uczestnik(idUczestnika: string): RoundtableParticipant | null;
  /** Dokłada albo podmienia uczestnika — klucz po `participantId`. */
  dodajUczestnika(uczestnik: RoundtableParticipant): void;
  // Podmienia skład po czynności moderatora; wykaz pusty jest pomijany, znaczy brak zmiany.
  ustawUczestnikow(lista: RoundtableParticipant[]): void;
  /** Ile razy ten kanał wystąpił w debacie — dwie tożsamości to dwa wystąpienia. */
  wystapieniaKanalu(idKanalu: string): number;

  /** Tura bieżąca; pusta znaczy „żadnej nie uruchomiono". */
  tura(): string;
  /** Definicja tury bieżącej; `null` do pierwszego odczytu. */
  definicjaTury(): RoundtableTurn | null;
  /** Zapisuje turę bieżącą i powiadamia okna. */
  ustawTure(tura: RoundtableTurn): void;

  /** Wypowiedzi tury bieżącej w kolejności przyjścia. */
  wypowiedzi(): RoundtableStatement[];

  /** Kanały modelu znane rejestrowi rdzenia — katalog wyboru Model Panels. */
  kanaly(): Channel[];
  /** Czy rdzeń odpowiedział na `channel.list` choć raz (pustka a ładowanie). */
  kanalyOdczytane(): boolean;
  /** Zamawia wykaz kanałów z rdzenia. */
  odswiezKanaly(): void;
  /** Nazwa kanału uczestnika pokazywana Operatorowi; pusta, gdy kanał nieznany. */
  opisKanalu(idKanalu: string): string;

  /** Subskrypcja zmiany składu, tury, wypowiedzi albo wykazu kanałów. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzeń rdzenia. */
  zamknij(): void;
}

/** Zależności stanu debaty: okno debaty i tura otwierane od razu, obie pola opcjonalne o wartości domyślnej pustej. */
export interface OpcjeStanuDebaty {
  okno?: string;
  tura?: string;
}

export function utworzStanDebaty(
  zrodlo: ZrodloRoundtable,
  rejestr: RejestrKanalow,
  opcje: OpcjeStanuDebaty = {},
): StanDebaty {
  const sluchacze = new Set<() => void>();
  // Mapa zachowuje kolejność wstawiania, więc dodanie uczestników ustala wykaz bez osobnego sortowania.
  const skladDebaty = new Map<string, RoundtableParticipant>();
  let idOkna = opcje.okno ?? '';
  let idTury = opcje.tura ?? '';
  let definicja: RoundtableTurn | null = null;
  let wypowiedziTury: RoundtableStatement[] = [];

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  // Przyrost debaty przychodzi też z pracy innego okna; tura z innego okna debaty nas nie dotyczy.
  const odsubskrybujDebate = zrodlo.naZmianeDebaty((tresc) => {
    if (idOkna !== '' && tresc.turn.windowId !== idOkna) return;
    const zmianaTury = tresc.turn.id !== idTury;
    if (zmianaTury) wypowiedziTury = [];
    idTury = tresc.turn.id;
    definicja = tresc.turn;
    if (tresc.statement !== undefined) wypowiedziTury = wchlonWypowiedz(wypowiedziTury, tresc.statement);
    powiadom();
  });

  const odsubskrybujKanaly = rejestr.naZmiane(() => powiadom());

  return {
    okno: () => idOkna,

    ustawOkno(nowe) {
      const przyciety = nowe.trim();
      if (przyciety === idOkna) return;
      // Uczestnicy, tura i wypowiedzi należą do okna debaty; zostawione pokazywałyby cudzą debatę.
      idOkna = przyciety;
      skladDebaty.clear();
      idTury = '';
      definicja = null;
      wypowiedziTury = [];
      powiadom();
    },

    uczestnicy: () => poKolejnosci([...skladDebaty.values()]),

    uczestnik: (idUczestnika) => skladDebaty.get(idUczestnika) ?? null,

    dodajUczestnika(uczestnik) {
      skladDebaty.set(uczestnik.id, uczestnik);
      powiadom();
    },

    // Wykaz pusty od rdzenia nie kasuje składu — pusta tablica znaczy brak zmiany, nie skład pusty.
    ustawUczestnikow(lista) {
      if (lista.length === 0) return;
      skladDebaty.clear();
      for (const uczestnik of lista) skladDebaty.set(uczestnik.id, uczestnik);
      powiadom();
    },

    wystapieniaKanalu(idKanalu) {
      let ile = 0;
      for (const uczestnik of skladDebaty.values()) {
        if (uczestnik.channelId === idKanalu) ile += 1;
      }
      return ile;
    },

    tura: () => idTury,

    definicjaTury: () => definicja,

    ustawTure(tura) {
      if (tura.id !== idTury) wypowiedziTury = [];
      idTury = tura.id;
      definicja = tura;
      powiadom();
    },

    wypowiedzi: () => [...wypowiedziTury],

    kanaly: () => rejestr.kanaly(),

    kanalyOdczytane: () => rejestr.odpowiedzOtrzymana(),

    odswiezKanaly: () => rejestr.odswiez(),

    opisKanalu(idKanalu) {
      const znaleziony = rejestr.kanaly().find((kanal) => kanal.id === idKanalu);
      return znaleziony === undefined ? '' : nazwaKanalu(znaleziony);
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybujDebate();
      odsubskrybujKanaly();
    },
  };
}

/**
 * Wchłonięcie przyrostu wypowiedzi — jedna wypowiedź to jeden wpis, kluczowany identyfikatorem,
 * nie treścią ani mówcą.
 */
function wchlonWypowiedz(
  wykaz: readonly RoundtableStatement[],
  przyrost: RoundtableStatement,
): RoundtableStatement[] {
  const miejsce = wykaz.findIndex((wypowiedz) => wypowiedz.id === przyrost.id);
  if (miejsce === -1) return [...wykaz, przyrost];
  const nowy = [...wykaz];
  nowy[miejsce] = przyrost;
  return nowy;
}

/**
 * Kolejność głosu wykazu uczestników: sortujemy wyłącznie, gdy kolejność ma każdy uczestnik,
 * inaczej zostaje kolejność dodania.
 */
function poKolejnosci(lista: RoundtableParticipant[]): RoundtableParticipant[] {
  const kazdyZKolejnoscia = lista.every((uczestnik) => uczestnik.order !== undefined);
  if (!kazdyZKolejnoscia) return lista;
  return [...lista].sort((pierwszy, drugi) => (pierwszy.order ?? 0) - (drugi.order ?? 0));
}

/**
 * Nazwa uczestnika widziana przez Operatora, ważniejsza tożsamością niż kanałem, bo dwóch
 * uczestników potrafi jechać tym samym kanałem.
 */
export function nazwaUczestnika(uczestnik: RoundtableParticipant, opisKanalu: string): string {
  const kanal = opisKanalu === '' ? uczestnik.channelId : opisKanalu;
  const tozsamosc = uczestnik.personaName;
  if (tozsamosc !== undefined && tozsamosc.length > 0) return `${tozsamosc} · ${kanal}`;
  return `${kanal} · ${uczestnik.id}`;
}
