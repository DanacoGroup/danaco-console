import type { ResearchFinding, ResearchReport, ResearchSource } from '../../../../shared/contract';

/**
 * Pamięć jednego badania — dane, które pięć okien ogląda wspólnie.
 *
 * Jedna odpowiedzialność: przechowanie i ogłoszenie zmiany. Pamięć nie zna
 * kontraktu i nie woła rdzenia — dzięki temu odczyt (`odczyt-badania`) i wykazy
 * okien patrzą na ten sam zbiór, a nie na jego kopie.
 *
 * Fazy są trzy, nie dwie: „jeszcze nie pytałem", „pytam" i „rdzeń nic nie ma"
 * zostają rozróżnialne, bo zlanie ich w jedno kazałoby zgadywać, czy czekać,
 * czy działać.
 */
export type FazaBadania = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface PamiecBadania {
  idOkna(): string;
  faza(): FazaBadania;
  powod(): string;
  zakres(): string;
  etapy(): readonly string[];
  zrodla(): readonly ResearchSource[];
  ustalenia(): readonly ResearchFinding[];
  raport(): ResearchReport | null;
  /** Ustawia fazę odczytu wraz z jej powodem; pusty powód znaczy brak powodu. */
  ustawFaze(faza: FazaBadania, powod: string): void;
  /** Zapamiętuje okno badania wskazane przez rdzeń. */
  ustawOkno(idOkna: string): void;
  wchlonZakres(zakres: string, etapy: readonly string[]): void;
  wchlonZrodlo(zrodlo: ResearchSource): void;
  wchlonUstalenie(ustalenie: ResearchFinding): void;
  wchlonRaport(raport: ResearchReport | null): void;
  obserwuj(sluchacz: () => void): () => void;
  /** Zdejmuje wszystkich obserwatorów. */
  zapomnij(): void;
}

export function utworzPamiecBadania(): PamiecBadania {
  const sluchacze = new Set<() => void>();
  const dane = {
    okno: '',
    faza: 'spoczynek' as FazaBadania,
    powod: '',
    zakres: '',
    etapy: [] as readonly string[],
    zrodla: [] as readonly ResearchSource[],
    ustalenia: [] as readonly ResearchFinding[],
    raport: null as ResearchReport | null,
  };

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  return {
    idOkna: () => dane.okno,
    faza: () => dane.faza,
    powod: () => dane.powod,
    zakres: () => dane.zakres,
    etapy: () => dane.etapy,
    zrodla: () => dane.zrodla,
    ustalenia: () => dane.ustalenia,
    raport: () => dane.raport,

    ustawFaze(faza, powod) {
      dane.faza = faza;
      dane.powod = powod;
      oglos();
    },

    ustawOkno(idOkna) {
      dane.okno = idOkna;
    },

    wchlonZakres(zakres, etapy) {
      dane.zakres = zakres;
      dane.etapy = [...etapy];
      oglos();
    },

    wchlonZrodlo(zrodlo) {
      dane.zrodla = [...dane.zrodla.filter((wpis) => wpis.id !== zrodlo.id), zrodlo];
      oglos();
    },

    wchlonUstalenie(ustalenie) {
      const znane = dane.ustalenia.some((wpis) => wpis.id === ustalenie.id);
      dane.ustalenia = znane
        ? dane.ustalenia.map((wpis) => (wpis.id === ustalenie.id ? ustalenie : wpis))
        : [...dane.ustalenia, ustalenie];
      oglos();
    },

    wchlonRaport(raport) {
      dane.raport = raport;
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zapomnij: () => sluchacze.clear(),
  };
}
