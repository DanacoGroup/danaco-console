import type { TerminalProcess, TerminalSession } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import { utworzBuforWyjscia, type BuforWyjscia } from './bufor-wyjscia';
import type { ZrodloTerminala } from './zrodlo-terminala';

/**
 * Stan modułu Terminal jest jedną prawdą dla trzech okien operacyjnych: kart, procesów i wyjścia, budzonych zdarzeniami, a nie odpytywaniem.
 */
export interface StanTerminala {
  /** Okno komunikacji, w którym pracuje moduł. */
  okno(): string;
  /** Karty powłok w kolejności otwarcia. */
  karty(): readonly TerminalSession[];
  /** Karta bieżąca; pusta, gdy nie ma ani jednej. */
  kartaBiezaca(): TerminalSession | null;
  /** Przestawia ognisko na wskazaną kartę; fałsz znaczy „widok takiej karty nie ma”. */
  ustawKarte(idKarty: string): boolean;
  /** Dokłada kartę i czyni ją bieżącą. */
  dodajKarte(karta: TerminalSession): void;
  /** Zamyka kartę w widoku klienta; fałsz znaczy „widok takiej karty nie ma”. */
  zamknijKarte(idKarty: string): boolean;
  /** Zmienia nazwę karty w widoku klienta. */
  przemianujKarte(idKarty: string, nazwa: string): void;
  /** Przypina albo odpina kartę. */
  przypnijKarte(idKarty: string, przypieta: boolean): void;
  /** Czy karta jest przypięta. */
  czyPrzypieta(idKarty: string): boolean;
  /** Procesy znane oknu, od najnowszego. */
  procesy(): readonly TerminalProcess[];
  /** Zapisuje komplet procesów po odczycie z rdzenia. */
  ustawProcesy(procesy: readonly TerminalProcess[]): void;
  /** Wstawia albo podmienia jeden proces. */
  zapiszProces(proces: TerminalProcess): void;
  /** Ostatnie polecenie karty — nośnik „uruchom ponownie”. */
  ostatniePolecenie(idKarty: string): string;
  /** Zapamiętuje polecenie wysłane do karty. */
  zapamietajPolecenie(idKarty: string, polecenie: string): void;
  /** Bufor wyjścia wspólny dla wszystkich kart. */
  bufor(): BuforWyjscia;
  /** Subskrypcja zmiany stanu modułu. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzeń rdzenia. */
  zamknij(): void;
}

/**
 * Zależności stanu modułu Terminal: rejestr otwartych procesów, bufor wyjścia i źródło zdarzeń terminala.
 */
export interface OpcjeStanu {
  /** Okno komunikacji modułu — bez niego nie ma czego otworzyć. */
  okno: string;
}

export function utworzStanTerminala(zrodlo: ZrodloTerminala, opcje: OpcjeStanu): StanTerminala {
  const sluchacze = new Set<() => void>();
  const procesy = new Map<string, TerminalProcess>();
  const bufor = utworzBuforWyjscia();

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  const rejestr = utworzRejestrKart(powiadom);

  function kartaProcesu(idProcesu: string): string {
    return procesy.get(idProcesu)?.sessionId ?? '';
  }

  const odsubskrybujProces = zrodlo.naZmianeProcesu((tresc) => {
    procesy.set(tresc.process.id, tresc.process);
    if (tresc.process.status !== 'running') bufor.domknij(tresc.process.id);
    powiadom();
  });

  const odsubskrybujWyjscie = zrodlo.naFragmentWyjscia((tresc, koperta) => {
    // Strumień jest wspólny dla platformy; do konsoli wchodzi wyłącznie fragment znanego procesu.
    if (tresc.windowId !== opcje.okno || !procesy.has(tresc.messageId)) return;
    const dopisane = bufor.dopisz(
      tresc.messageId,
      kartaProcesu(tresc.messageId),
      tresc.kind,
      tresc.text ?? '',
      koperta.timestamp,
    );
    if (koperta.done === true) bufor.domknij(tresc.messageId);
    if (dopisane > 0 || koperta.done === true) powiadom();
  });

  return {
    okno: () => opcje.okno,

    karty: () => rejestr.karty(),
    kartaBiezaca: () => rejestr.kartaBiezaca(),
    ustawKarte: (idKarty) => rejestr.ustawKarte(idKarty),
    dodajKarte: (karta) => rejestr.dodajKarte(karta),
    zamknijKarte: (idKarty) => rejestr.zamknijKarte(idKarty),
    przemianujKarte: (idKarty, nazwa) => rejestr.przemianujKarte(idKarty, nazwa),
    przypnijKarte: (idKarty, przypieta) => rejestr.przypnijKarte(idKarty, przypieta),
    czyPrzypieta: (idKarty) => rejestr.czyPrzypieta(idKarty),
    ostatniePolecenie: (idKarty) => rejestr.ostatniePolecenie(idKarty),
    zapamietajPolecenie: (idKarty, polecenie) => rejestr.zapamietajPolecenie(idKarty, polecenie),

    procesy: () => [...procesy.values()].sort((a, b) => b.startedAt - a.startedAt),

    ustawProcesy(nowe) {
      procesy.clear();
      for (const proces of nowe) procesy.set(proces.id, proces);
      powiadom();
    },

    zapiszProces(proces) {
      procesy.set(proces.id, proces);
      powiadom();
    },

    bufor: () => bufor,

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybujProces();
      odsubskrybujWyjscie();
    },
  };
}

/**
 * Karty okna terminala: wykaz otwartych kart, ognisko, przypięcia i ostatnio wykonane polecenia każdej z nich.
 */
interface RejestrKart {
  karty(): readonly TerminalSession[];
  kartaBiezaca(): TerminalSession | null;
  ustawKarte(idKarty: string): boolean;
  dodajKarte(karta: TerminalSession): void;
  zamknijKarte(idKarty: string): boolean;
  przemianujKarte(idKarty: string, nazwa: string): void;
  przypnijKarte(idKarty: string, przypieta: boolean): void;
  czyPrzypieta(idKarty: string): boolean;
  ostatniePolecenie(idKarty: string): string;
  zapamietajPolecenie(idKarty: string, polecenie: string): void;
}

function utworzRejestrKart(powiadom: () => void): RejestrKart {
  const karty: TerminalSession[] = [];
  const przypiete = new Set<string>();
  const polecenia = new Map<string, string>();
  let biezaca = '';

  return {
    karty: () => karty,

    kartaBiezaca: () => karty.find((karta) => karta.id === biezaca) ?? null,

    // Ognisko wolno przestawić wyłącznie na kartę znaną temu widokowi, nie na kartę z innej sesji.
    ustawKarte(idKarty) {
      if (biezaca === idKarty) return true;
      if (!karty.some((karta) => karta.id === idKarty)) return false;
      biezaca = idKarty;
      powiadom();
      return true;
    },

    dodajKarte(karta) {
      karty.push(karta);
      biezaca = karta.id;
      powiadom();
    },

    zamknijKarte(idKarty) {
      const miejsce = karty.findIndex((karta) => karta.id === idKarty);
      if (miejsce < 0) return false;
      karty.splice(miejsce, 1);
      przypiete.delete(idKarty);
      polecenia.delete(idKarty);
      if (biezaca === idKarty) biezaca = karty[karty.length - 1]?.id ?? '';
      powiadom();
      return true;
    },

    przemianujKarte(idKarty, nazwa) {
      const karta = karty.find((pozycja) => pozycja.id === idKarty);
      if (karta === undefined) return;
      karta.title = nazwa;
      powiadom();
    },

    przypnijKarte(idKarty, przypieta) {
      if (przypieta) przypiete.add(idKarty);
      else przypiete.delete(idKarty);
      powiadom();
    },

    czyPrzypieta: (idKarty) => przypiete.has(idKarty),

    ostatniePolecenie: (idKarty) => polecenia.get(idKarty) ?? '',

    zapamietajPolecenie(idKarty, polecenie) {
      polecenia.set(idKarty, polecenie);
    },
  };
}
