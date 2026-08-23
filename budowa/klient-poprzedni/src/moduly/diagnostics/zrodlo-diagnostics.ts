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
 * Moduł Diagnostics widziany przez klienta — cztery komendy obszaru
 * `diagnostics.*` oraz zdarzenie zmiany analizy.
 *
 * Każda czynność oddaje `Wynik`, nie samą treść, i nigdzie nie podstawia pustej
 * kolekcji za nieudany odczyt. Port `Diagnostyka` jest w rdzeniu odbiorcą odmów
 * wykonania komend (`server/internal/core/kompozycja.go`, pole `Diagnostyka`),
 * a Errors Panel pokazuje te odmowy; źródło, które połknęłoby odmowę odczytu
 * i oddało pusty wykaz, kłamałoby o własnym niepowodzeniu w oknie poświęconym
 * niepowodzeniom cudzym.
 *
 * Asymetria kontraktu jest zamierzona: `windowId` występuje wyłącznie
 * w `diagnostics.analyze.run` i jest tam nieobowiązkowy — analiza zapisuje,
 * z którego okna ją uruchomiono. Pozostałe trzy komendy czytają dziennik, błędy
 * i rekomendacje całej instalacji, więc pojęcia okna nie mają.
 *
 * Niepowodzenie niesie swoją drogę. `Wynik` warstwy protokołu nie odróżnia odmowy
 * rdzenia od odpowiedzi, której klient nie zrozumiał, ani od odpowiedzi bez treści
 * — a to trzy różne zdania dla Operatora i tylko jedno z nich brzmi „rdzeń odmówił"
 * (`niepowodzenie-odczytu.ts`). Rozstrzygnięcie zapada tutaj, bo tylko tutaj widać
 * surową odpowiedź przed sprawdzianem kształtu.
 */
/** Wynik komendy modułu wraz z drogą, którą przyszło niepowodzenie. */
export type WynikDiagnostyki<T> = Wynik<T> & { powod?: PowodNiepowodzenia };

export interface ZrodloDiagnostics {
  /** `diagnostics.analyze.run` — uruchomienie analizy zagregowanego stanu. */
  uruchomAnalize(zadanie: DiagnosticsAnalyzeRunRequest): Promise<WynikDiagnostyki<DiagnosticAnalysis>>;
  /**
   * `diagnostics.log.query` — przeszukanie dziennika. Oddaje wpisy wraz
   * z licznikiem i znacznikiem przycięcia, bo Logs Viewer musi rozpoznać
   * wynik przycięty granicą bufora i powiedzieć to Operatorowi.
   */
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
 * Rozstrzyga odpowiedź na komendę, zapamiętując, gdzie się potknęła.
 *
 * Kolejność sprawdzeń jest istotna: dopóki nie wiadomo, czy `surowy` był odmową,
 * nie wolno wołać `sprawdzKsztalt` i wziąć jego kodu za kod rdzenia.
 *
 * Wyjście poza plik ma jednego odbiorcę: `zrodlo-obserwowalnosci.ts` sięga po
 * rodziny `monitor.*`, a rozróżnienie odmowy od odpowiedzi nieczytelnej ma
 * w Observability Tools tę samą wagę co w Errors Panelu. Druga kopia tego
 * rozstrzygnięcia dałaby dwa zdania o jednej ciszy rdzenia.
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
