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
 * Debata bieżąca modułu — jeden stan wspólny dla sześciu okien: Model Panels,
 * Debate Panel, Argument Map & Analysis, Voting & Evaluation Center, Moderator
 * Panel i Consensus Panel widzą ten sam skład i tę samą turę.
 *
 * Uczestnik jest kluczowany po `participantId`, nie po `channelId`: kontrakt przy
 * `roundtable.model.add` dopuszcza dwa wystąpienia tego samego kanału pod odrębnymi
 * tożsamościami, a wykaz po kanale nadpisałby jedno z nich drugim.
 *
 * Wykaz kanałów pochodzi z rejestru rdzenia (`sterowanie/rejestr-kanalow`), tego
 * samego, którym jedzie okno rozmowy — moduł nie prowadzi drugiej listy modeli.
 */
export interface StanDebaty {
  /** Okno debaty; wymagane przez każdą komendę obszaru, którą moduł wywołuje. */
  okno(): string;
  /** Przestawia moduł na inne okno debaty i czyści byty poprzedniej. */
  ustawOkno(idOkna: string): void;

  /**
   * Uczestnicy debaty; wspólni dla Model Panels i Moderator Panel.
   *
   * Kolejność jest kolejnością głosu, gdy rdzeń ją podał (`order` u KAŻDEGO
   * uczestnika), a w przeciwnym razie kolejnością dodania — patrz `poKolejnosci`.
   */
  uczestnicy(): RoundtableParticipant[];
  /** Uczestnik o podanej tożsamości; `null`, gdy nieznany. */
  uczestnik(idUczestnika: string): RoundtableParticipant | null;
  /** Dokłada albo podmienia uczestnika — klucz po `participantId`. */
  dodajUczestnika(uczestnik: RoundtableParticipant): void;
  /**
   * Podmienia skład po czynności moderatora oddającej `participants`.
   * Wykaz pusty jest pomijany — uzasadnienie przy ciele czynności.
   */
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

/** Zależności stanu: okno debaty i tura otwierane od razu. */
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
  // Mapa zachowuje kolejność wstawiania, więc kolejność dodania uczestników
  // jest kolejnością wykazu w oknie bez osobnego sortowania.
  const skladDebaty = new Map<string, RoundtableParticipant>();
  let idOkna = opcje.okno ?? '';
  let idTury = opcje.tura ?? '';
  let definicja: RoundtableTurn | null = null;
  let wypowiedziTury: RoundtableStatement[] = [];

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  // Przyrost debaty przychodzi także z pracy innego okna albo innego urządzenia
  // tego konta. Tura z innego okna debaty nas nie dotyczy i nie rusza widoku.
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
      // Uczestnicy, tura i wypowiedzi należą do okna debaty. Zostawione po
      // poprzednim oknie pokazywałyby cudzą debatę pod nowym identyfikatorem.
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

    /**
     * Wykaz pusty od rdzenia nie kasuje składu — pusta tablica znaczy „bez zmian
     * w składzie", nie „skład jest pusty".
     *
     * `participants` w odpowiedzi `roundtable.moderator.direct` jest polem
     * nieobowiązkowym; rdzeń dokłada je tylko przy zmianie składu albo kolejności
     * głosu. Odczyt składu jest w kontrakcie (`roundtable.model.list`), ale żadne
     * okno modułu go jeszcze nie wywołuje, więc wykaz raz wymazany zostaje
     * w tym stanie nieodzyskiwalny.
     */
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
 * Wchłonięcie przyrostu wypowiedzi — jedna wypowiedź to jeden wpis.
 *
 * Rdzeń rozgłasza tę samą wypowiedź dwukrotnie: jako `created` w chwili otwarcia
 * głosu, gdy treść jest jeszcze pusta, i jako `updated` po domknięciu strumienia
 * modelu. Klucz jest `id` wypowiedzi, a nie jej treść ani mówca: przyrost
 * o znanym identyfikatorze podmienia wpis w miejscu (treść rośnie, kolejność
 * zostaje), przyrost nieznany dokleja się na koniec.
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
 * Kolejność głosu wykazu uczestników.
 *
 * `RoundtableParticipant.order` jest w kontrakcie polem nieobowiązkowym,
 * a `speakingOrder` przy `roundtable.moderator.direct` obejmuje cały stół —
 * kolejność głosu jest własnością składu, nie pojedynczego uczestnika. Wykaz
 * z `order` tylko u części uczestników znaczy, że rdzeń kolejności nie ustalił.
 *
 * Dlatego sortujemy wyłącznie wtedy, gdy `order` ma każdy uczestnik; inaczej
 * zostaje kolejność dodania (Mapa zachowuje kolejność wstawiania). Decyzja stoi
 * w stanie wspólnym, żeby wszystkie okna widziały tę samą kolejność.
 */
function poKolejnosci(lista: RoundtableParticipant[]): RoundtableParticipant[] {
  const kazdyZKolejnoscia = lista.every((uczestnik) => uczestnik.order !== undefined);
  if (!kazdyZKolejnoscia) return lista;
  return [...lista].sort((pierwszy, drugi) => (pierwszy.order ?? 0) - (drugi.order ?? 0));
}

/**
 * Nazwa uczestnika widziana przez Operatora.
 *
 * Tożsamość jest ważniejsza od kanału, bo dwóch uczestników potrafi jechać tym
 * samym kanałem i sama nazwa modelu ich nie odróżnia. Gdy tożsamości nie nadano,
 * zostaje opis kanału i identyfikator uczestnika, żeby dwa wiersze wykazu nie
 * wyglądały identycznie.
 */
export function nazwaUczestnika(uczestnik: RoundtableParticipant, opisKanalu: string): string {
  const kanal = opisKanalu === '' ? uczestnik.channelId : opisKanalu;
  const tozsamosc = uczestnik.personaName;
  if (tozsamosc !== undefined && tozsamosc.length > 0) return `${tozsamosc} · ${kanal}`;
  return `${kanal} · ${uczestnik.id}`;
}
