import type { DiagnosticAnalysis } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Analiza bieżąca i zakres czasu modułu Diagnostics, stanowiące jedną prawdę
 * dla czterech okien. Zakres czasu stoi obok analizy, ponieważ dziennik, wykaz
 * błędów i sama analiza przyjmują końce zakresu niezależnie.
 */
export interface StanDiagnostyki {
  /** Identyfikator analizy bieżącej; pusty znaczy „nie uruchomiono”. */
  analiza(): string;
  /** Migawka analizy bieżącej; pusta do pierwszego uruchomienia. */
  migawka(): DiagnosticAnalysis | null;
  /** Przestawia moduł na inną analizę i powiadamia okna. */
  ustawAnalize(idAnalizy: string, migawka?: DiagnosticAnalysis): void;
  /** Zapisuje migawkę po uruchomieniu analizy i powiadamia okna zależne. */
  ustawMigawke(migawka: DiagnosticAnalysis): void;
  /** Zakres czasu wspólny dla dziennika, błędów i analizy. */
  zakres(): ZakresCzasu;
  /** Zapisuje zakres czasu i powiadamia okna; oba końce są nieobowiązkowe. */
  ustawZakres(zakres: ZakresCzasu): void;
  /** Subskrypcja zmiany analizy albo zakresu. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzeń rdzenia. */
  zamknij(): void;
}

/**
 * Zakres czasu podawany w milisekundach epoki; końce nieustawione znaczą brak
 * granicy. Zakres pusty jest stanem poprawnym i mówi, że zawężenia nie ma.
 */
export interface ZakresCzasu {
  od?: number;
  do?: number;
}

/**
 * Zależności stanu diagnostyki: identyfikator analizy oraz zakres czasu, które
 * moduł otwiera od razu przy montażu, zanim Operator wykona pierwszą czynność
 * w którymkolwiek z czterech okien.
 */
export interface OpcjeStanuDiagnostyki {
  analiza?: string;
  zakres?: ZakresCzasu;
}

export function utworzStanDiagnostyki(
  zrodlo: ZrodloDiagnostics,
  opcje: OpcjeStanuDiagnostyki = {},
): StanDiagnostyki {
  const sluchacze = new Set<() => void>();
  let idAnalizy = opcje.analiza ?? '';
  let migawka: DiagnosticAnalysis | null = null;
  let zakres: ZakresCzasu = opcje.zakres ?? {};

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  // Analiza zmienia się także pracą innego okna albo innego urządzenia konta.
  const odsubskrybujAnalize = zrodlo.naZmianeAnalizy((tresc) => {
    if (tresc.analysis.id !== idAnalizy) return;
    migawka = tresc.analysis;
    powiadom();
  });

  return {
    analiza: () => idAnalizy,

    migawka: () => migawka,

    ustawAnalize(nowa, nowaMigawka) {
      const przyciety = nowa.trim();
      if (przyciety === idAnalizy && nowaMigawka === undefined) return;
      idAnalizy = przyciety;
      migawka = nowaMigawka ?? null;
      powiadom();
    },

    ustawMigawke(nowa) {
      migawka = nowa;
      idAnalizy = nowa.id;
      powiadom();
    },

    zakres: () => ({ ...zakres }),

    ustawZakres(nowy) {
      if (nowy.od === zakres.od && nowy.do === zakres.do) return;
      zakres = {
        ...(nowy.od === undefined ? {} : { od: nowy.od }),
        ...(nowy.do === undefined ? {} : { do: nowy.do }),
      };
      powiadom();
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybujAnalize();
    },
  };
}
