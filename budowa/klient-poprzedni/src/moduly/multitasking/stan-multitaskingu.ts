import {
  MessageRole,
  MessageStatus,
  type LoopState,
  type Message,
  type Queue,
  type Window,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { czyWykonawca, zlozObsade, OBSADA_PUSTA, type Obsada } from './obsada-rol';
import { utworzStrumienWykonawcy, type StrumienWykonawcy } from './strumien-wykonawcy';
import { TrybWspolpracy, type TrybWspolpracy as Tryb } from './tryby-wspolpracy';
import { Wcielenie, type Wcielenie as WcielenieRoli } from './wcielenia-analizy';
import type { ZrodloBiegu } from './zrodlo-biegu';

/**
 * Stan wspólny czterem oknom ról — jedna prawda zamiast czterech.
 *
 * Koordynator steruje kolejką, wykonawcy w niej pracują, a analityk czyta ich
 * wyniki. Własna obsada w każdym oknie sprawiłaby, że przepięcie wykonawcy pod
 * innego koordynatora zmieniłoby jedno okno i zostawiło trzy z obrazem
 * nieaktualnym.
 *
 * Okna żyją ze zdarzeń: `stream.chunk` niesie strumień wykonawcy do
 * koordynatora, `message.changed` domknięty wynik do analityka, `queue.changed`
 * stan etapu. Odczyty `window.list`, `window.state.get` i `message.list` służą
 * pierwszemu wypełnieniu i wznowieniu po rozłączeniu.
 *
 * Wewnątrz stoją trzy rejestry — obsada ról, tury wykonawców i kolejki etapów —
 * każdy z własnymi polami; na zewnątrz okna widzą je wyłącznie przez wspólny
 * `StanMultitaskingu`, bo koordynator musi widzieć ten sam przebieg co
 * wykonawca, którym steruje.
 */
export interface StanMultitaskingu {
  /** Sesja, w której stoją okna ról. */
  sesja(): string;
  /** Obsada czterech ról odczytana z okien sesji. */
  obsada(): Obsada;
  /** Zapisuje okna sesji i przelicza obsadę. */
  ustawOkna(okna: readonly Window[]): void;
  /** Okno wskazane jako analityk; utrwalane na poziomie sesji. */
  ustawAnalityka(idOkna: string): void;
  /** Strumień wykonawców widziany przez koordynatora. */
  strumien(): StrumienWykonawcy;
  /** Stan biegu naprawczego koordynatora; pusty przed pierwszym odczytem. */
  bieg(): LoopState | null;
  /** Zapisuje stan biegu z `window.state.get`. */
  ustawBieg(stan: LoopState | null): void;
  /** Kolejki etapów założone dla tej obsady. */
  kolejki(): readonly Queue[];
  /** Wstawia albo podmienia kolejkę etapu. */
  zapiszKolejke(kolejka: Queue): void;
  /** Kolejka wskazana w sterowaniu; pusta, gdy nie ma ani jednej. */
  kolejkaBiezaca(): Queue | null;
  /** Przestawia kolejkę bieżącą. */
  ustawKolejke(idKolejki: string): void;
  /** Tryb współpracy wykonawców. */
  tryb(): Tryb;
  /** Przestawia tryb współpracy. */
  ustawTryb(tryb: Tryb): void;
  /** Wcielenie roli analityka. */
  wcielenie(): WcielenieRoli;
  /** Przestawia wcielenie analityka. */
  ustawWcielenie(wcielenie: WcielenieRoli): void;
  /** Okno wykonawcy, które dostało zlecenie poprzednio. */
  poprzedniWykonawca(): string;
  /** Zapamiętuje adresata ostatniego przekazania. */
  zapamietajPrzekazanie(idOkna: string): void;
  /** Ostatnia domknięta odpowiedź wykonawcy; pusta przed pierwszą turą. */
  wynikWykonawcy(idOkna: string): Message | null;
  /** Zapisuje wiadomości okna po odczycie `message.list`. */
  ustawWiadomosci(idOkna: string, wiadomosci: readonly Message[]): void;
  /** Czy w oknie trwa tura — po niej gaśnie znacznik pracy. */
  czyTura(idOkna: string): boolean;
  /** Subskrypcja zmiany stanu modułu. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzeń rdzenia. */
  zamknij(): void;
}

/** Zależności stanu. */
export interface OpcjeStanu {
  /** Sesja, w której stoją okna ról. */
  sesja: string;
  /** Okno wskazane jako analityk; puste, dopóki Operator nie wskaże. */
  analityk?: string;
}

export function utworzStanMultitaskingu(
  zrodlo: ZrodloBiegu,
  opcje: OpcjeStanu,
): StanMultitaskingu {
  const sluchacze = new Set<() => void>();
  const strumien = utworzStrumienWykonawcy();
  let bieg: LoopState | null = null;
  let tryb: Tryb = TrybWspolpracy.Niezalezna;
  let wcielenie: WcielenieRoli = Wcielenie.Validator;

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  const role = utworzRejestrObsady(opcje.analityk ?? '', powiadom);
  const tury = utworzRejestrTurWykonawcow(powiadom);
  const etapy = utworzRejestrKolejekEtapow(powiadom);

  const odsubskrybujFragment = zrodlo.naFragment((tresc) => {
    // Strumień jest wspólny całej platformie. Do widoku koordynatora wchodzi
    // wyłącznie fragment okna należącego do jego obsady — inaczej pokazywałby
    // cudzą rozmowę jako pracę swojego wykonawcy.
    if (!czyWykonawca(role.obsada(), tresc.windowId)) return;
    tury.otworz(tresc.windowId);
    if (strumien.dopisz(tresc, Date.now())) powiadom();
  });

  const odsubskrybujWiadomosc = zrodlo.naWiadomosc((tresc) => {
    const idOkna = tresc.message.windowId;
    if (!czyWykonawca(role.obsada(), idOkna) && idOkna !== role.analityk()) return;
    tury.przyjmij(tresc.message);
  });

  const odsubskrybujKolejke = zrodlo.naKolejke((tresc) => {
    etapy.przyjmij(tresc.queue, opcje.sesja);
  });

  return {
    sesja: () => opcje.sesja,
    obsada: () => role.obsada(),
    ustawOkna: (nowe) => role.ustawOkna(nowe),
    ustawAnalityka: (idOkna) => role.ustawAnalityka(idOkna),
    strumien: () => strumien,
    bieg: () => bieg,

    ustawBieg(stan) {
      bieg = stan;
      powiadom();
    },

    kolejki: () => etapy.wszystkie(),
    zapiszKolejke: (kolejka) => etapy.zapisz(kolejka),
    kolejkaBiezaca: () => etapy.biezaca(),
    ustawKolejke: (idKolejki) => etapy.ustawBiezaca(idKolejki),
    tryb: () => tryb,

    ustawTryb(nowy) {
      tryb = nowy;
      powiadom();
    },

    wcielenie: () => wcielenie,

    ustawWcielenie(nowe) {
      wcielenie = nowe;
      powiadom();
    },

    poprzedniWykonawca: () => tury.poprzedni(),
    zapamietajPrzekazanie: (idOkna) => tury.zapamietaj(idOkna),
    wynikWykonawcy: (idOkna) => tury.wynik(idOkna),
    ustawWiadomosci: (idOkna, wiadomosci) => tury.ustawWiadomosci(idOkna, wiadomosci),
    czyTura: (idOkna) => tury.trwa(idOkna),

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybujFragment();
      odsubskrybujWiadomosc();
      odsubskrybujKolejke();
    },
  };
}

/** Rejestr obsady: okna sesji, wskazanie analityka i przeliczona obsada ról. */
interface RejestrObsady {
  obsada(): Obsada;
  /** Okno wskazane jako analityk; puste, dopóki Operator nie wskaże. */
  analityk(): string;
  ustawOkna(okna: readonly Window[]): void;
  ustawAnalityka(idOkna: string): void;
}

/**
 * Rejestr obsady ról.
 *
 * Domyka się na trzech własnych polach, a o zmianie melduje przekazanym
 * wywołaniem zwrotnym — o oknach ról i o rdzeniu nie wie nic.
 */
function utworzRejestrObsady(analitykPoczatkowy: string, powiadom: () => void): RejestrObsady {
  let okna: readonly Window[] = [];
  let obsada: Obsada = OBSADA_PUSTA;
  let analityk = analitykPoczatkowy;

  function przelicz(): void {
    obsada = zlozObsade(okna, analityk);
    powiadom();
  }

  return {
    obsada: () => obsada,
    analityk: () => analityk,

    ustawOkna(nowe) {
      okna = [...nowe];
      przelicz();
    },

    ustawAnalityka(idOkna) {
      analityk = idOkna;
      przelicz();
    },
  };
}

/** Rejestr tur i wyników wykonawców wraz z adresatem ostatniego przekazania. */
interface RejestrTurWykonawcow {
  /** Otwiera turę okna; sam nie rozgłasza — robi to nadawca fragmentu. */
  otworz(idOkna: string): void;
  trwa(idOkna: string): boolean;
  wynik(idOkna: string): Message | null;
  /** Zdarzenie `message.changed`; cudza rola i cudze okno są tu bez znaczenia. */
  przyjmij(wiadomosc: Message): void;
  ustawWiadomosci(idOkna: string, wiadomosci: readonly Message[]): void;
  poprzedni(): string;
  zapamietaj(idOkna: string): void;
}

/**
 * Rejestr tur i wyników wykonawców.
 *
 * Przynależność okna do obsady rozstrzyga nadawca zdarzenia przed wywołaniem,
 * więc rejestr trzyma wyłącznie własne trzy pola i obsady nie zna.
 */
function utworzRejestrTurWykonawcow(powiadom: () => void): RejestrTurWykonawcow {
  const wyniki = new Map<string, Message>();
  const tury = new Set<string>();
  let poprzedni = '';

  return {
    otworz(idOkna) {
      tury.add(idOkna);
    },

    trwa: (idOkna) => tury.has(idOkna),
    wynik: (idOkna) => wyniki.get(idOkna) ?? null,
    poprzedni: () => poprzedni,

    przyjmij(wiadomosc) {
      if (wiadomosc.role !== MessageRole.Assistant) return;
      if (wiadomosc.status !== MessageStatus.Streaming) tury.delete(wiadomosc.windowId);
      wyniki.set(wiadomosc.windowId, wiadomosc);
      powiadom();
    },

    ustawWiadomosci(idOkna, wiadomosci) {
      const ostatnia = [...wiadomosci].reverse().find((w) => w.role === MessageRole.Assistant);
      if (ostatnia === undefined) wyniki.delete(idOkna);
      else wyniki.set(idOkna, ostatnia);
      powiadom();
    },

    zapamietaj(idOkna) {
      poprzedni = idOkna;
      tury.add(idOkna);
      powiadom();
    },
  };
}

/** Rejestr kolejek etapów wraz ze wskazaniem kolejki bieżącej. */
interface RejestrKolejekEtapow {
  wszystkie(): readonly Queue[];
  zapisz(kolejka: Queue): void;
  biezaca(): Queue | null;
  ustawBiezaca(idKolejki: string): void;
  /** Zdarzenie `queue.changed`; kolejka spoza sesji i spoza rejestru przechodzi bokiem. */
  przyjmij(kolejka: Queue, sesja: string): void;
}

/**
 * Rejestr kolejek etapów.
 *
 * Poza mapą kolejek i wskazaniem bieżącej nie dotyka niczego; sesję do odsiania
 * cudzych kolejek bierze parametrem, a o zmianie melduje przekazanym wywołaniem
 * zwrotnym.
 */
function utworzRejestrKolejekEtapow(powiadom: () => void): RejestrKolejekEtapow {
  const kolejki = new Map<string, Queue>();
  let wskazana = '';

  function wstaw(kolejka: Queue): void {
    kolejki.set(kolejka.id, kolejka);
    if (wskazana === '') wskazana = kolejka.id;
    powiadom();
  }

  return {
    wszystkie: () => [...kolejki.values()].sort((a, b) => a.createdAt - b.createdAt),
    zapisz: (kolejka) => wstaw(kolejka),
    biezaca: () => kolejki.get(wskazana) ?? null,

    ustawBiezaca(idKolejki) {
      if (wskazana === idKolejki) return;
      wskazana = idKolejki;
      powiadom();
    },

    przyjmij(kolejka, sesja) {
      if (!kolejki.has(kolejka.id) && kolejka.sessionId !== sesja) return;
      wstaw(kolejka);
    },
  };
}
