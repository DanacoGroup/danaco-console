import type { DiagnosticAnalysis } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Analiza bieżąca i zakres czasu modułu — jedna prawda dla czterech okien.
 *
 * Recommendations Panel wyświetla rekomendacje tej analizy, którą uruchomiło
 * Diagnostics Center: kontrakt wiąże je polem `analysisId`
 * (`DiagnosticRecommendation.analysisId`, `DiagnosticsRecommendationListRequest.analysisId`).
 * Gdyby panel trzymał własne wskazanie, pokazywałby rekomendacje analizy
 * poprzedniej obok wyniku nowej — a Operator nie miałby z czego rozpoznać, że
 * patrzy na dwie różne migawki.
 *
 * Zakres czasu stoi obok analizy, bo trzy komendy przyjmują `fromTime`
 * i `toTime` niezależnie: dziennik, wykaz błędów i sama analiza. Rozjechany
 * zakres dałby Errors Panel z błędami jednej doby, dziennik z innej i analizę
 * z trzeciej — zestawienie wewnętrznie sprzeczne, po którym nie da się orzec
 * przyczyny. Zakres pusty (oba końce nieustawione) znaczy „bez zawężenia”
 * i jest stanem poprawnym, nie brakiem.
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

/** Zakres czasu w milisekundach epoki; końce nieustawione znaczą „bez granicy”. */
export interface ZakresCzasu {
  od?: number;
  do?: number;
}

/** Zależności stanu: analiza i zakres otwierane od razu. */
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

  // Analiza zmienia się także pracą innego okna albo innego urządzenia tego
  // konta — rdzeń dopisuje do niej rekomendacje po zakończeniu przebiegu.
  // Dotyczy analizy bieżącej, więc migawka w oknach jest nieaktualna i okna
  // mają się odczytać ponownie. Migawkę podmieniamy od razu, bo zdarzenie
  // niesie ją w całości i drugie pytanie rdzenia byłoby zbędne.
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
