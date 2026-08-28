import {
  Command,
  EventType,
  type DiagnosticAnalysis,
  type DiagnosticError,
  type DiagnosticRecommendation,
  type DiagnosticsAnalysisChangedEvent,
  type DiagnosticsAnalyzeRunRequest,
  type DiagnosticsErrorListRequest,
  type DiagnosticsLogQueryRequest,
  type DiagnosticsRecommendationListRequest,
  type LogEntry,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';
import type { PowodNiepowodzenia } from './niepowodzenie-odczytu';

/**
 * Wynik komendy modułu Diagnostics wraz z drogą, którą przyszło niepowodzenie:
 * żadna czynność nie podstawia pustej kolekcji za odczyt, który się nie udał.
 */
export type WynikDiagnostyki<T> = Wynik<T> & { powod?: PowodNiepowodzenia };

export interface ZrodloDiagnostics {
  /** `diagnostics.analyze.run` — uruchomienie analizy zagregowanego stanu. */
  uruchomAnalize(zadanie: DiagnosticsAnalyzeRunRequest): Promise<WynikDiagnostyki<DiagnosticAnalysis>>;
  /** Przeszukanie dziennika zwraca wpisy z licznikiem i znacznikiem przycięcia granicą bufora. */
  przeszukajDziennik(zadanie: DiagnosticsLogQueryRequest): Promise<
    WynikDiagnostyki<{ entries: LogEntry[]; total?: number; truncated?: boolean }>
  >;
  /** `diagnostics.error.list` — błędy zgrupowane po odcisku, wraz z liczbą. */
  wykazBledow(zadanie: DiagnosticsErrorListRequest): Promise<
    WynikDiagnostyki<{ errors: DiagnosticError[]; total?: number }>
  >;
  /** `diagnostics.recommendation.list` — rekomendacje powstałe z analizy. */
  wykazRekomendacji(
    zadanie: DiagnosticsRecommendationListRequest,
  ): Promise<WynikDiagnostyki<DiagnosticRecommendation[]>>;
  /** Subskrypcja `diagnostics.analysis.changed` — odświeża okna analizy. */
  naZmianeAnalizy(sluchacz: (tresc: DiagnosticsAnalysisChangedEvent) => void): Odsubskrybuj;
}

/**
 * Rozstrzyga odpowiedź na komendę, zapamiętując krok, na którym się potknęła:
 * odmowę rdzenia, nieczytelną odpowiedź albo brak treści.
 */
export function rozstrzygnij<Z, W>(
  surowy: Wynik<Z>,
  komenda: string,
  sprawdzian: (tresc: Z) => boolean,
  wybierz: (tresc: Z) => W,
): WynikDiagnostyki<W> {
  if (!surowy.udany) {
    return { udany: false, powod: 'odmowa-rdzenia', ...(surowy.blad === undefined ? {} : { blad: surowy.blad }) };
  }
  const sprawdzony = sprawdzKsztalt(surowy, komenda, sprawdzian);
  if (!sprawdzony.udany) {
    return {
      udany: false,
      powod: 'odpowiedz-nieczytelna',
      ...(sprawdzony.blad === undefined ? {} : { blad: sprawdzony.blad }),
    };
  }
  const wynik = przenies(sprawdzony, wybierz);
  return wynik.udany ? wynik : { ...wynik, powod: 'odpowiedz-bez-tresci' };
}

export function utworzZrodloDiagnostics(kanal: Kanal): ZrodloDiagnostics {
  return {
    async uruchomAnalize(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.DiagnosticsAnalyzeRun, zadanie),
        Command.DiagnosticsAnalyzeRun,
        (tresc) => czyObiekt(tresc.analysis),
        (tresc) => tresc.analysis,
      );
    },

    async przeszukajDziennik(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.DiagnosticsLogQuery, zadanie),
        Command.DiagnosticsLogQuery,
        (tresc) => czyTablica(tresc.entries),
        (tresc) => ({
          entries: tresc.entries,
          ...(tresc.total === undefined ? {} : { total: tresc.total }),
          ...(tresc.truncated === undefined ? {} : { truncated: tresc.truncated }),
        }),
      );
    },

    async wykazBledow(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.DiagnosticsErrorList, zadanie),
        Command.DiagnosticsErrorList,
        (tresc) => czyTablica(tresc.errors),
        (tresc) => ({
          errors: tresc.errors,
          ...(tresc.total === undefined ? {} : { total: tresc.total }),
        }),
      );
    },

    async wykazRekomendacji(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.DiagnosticsRecommendationList, zadanie),
        Command.DiagnosticsRecommendationList,
        (tresc) => czyTablica(tresc.recommendations),
        (tresc) => tresc.recommendations,
      );
    },

    naZmianeAnalizy(sluchacz) {
      return kanal.naZdarzenie(EventType.DiagnosticsAnalysisChanged, (tresc) => sluchacz(tresc));
    },
  };
}
