import {
  BuildStatus,
  ContainerStatus,
  DataEngine,
  DebugStatus,
  ProblemSeverity,
  ErrorCode,
  ScanKind,
  TestStatus,
} from '../../../../shared/contract';
import type { Wynik } from '../../protokol/kanal';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Źródło próbne warsztatu modułu Developer — rdzeń zastąpiony zapisem wywołań
 * i gotowymi odpowiedziami.
 *
 * Jedna atrapa dla całej rodziny zamiast atrapy przepisywanej w każdym
 * sprawdzianie: rozjazd z umową `ZrodloWarsztatu` przerywa wtedy kompilację,
 * zamiast rozjeżdżać sprawdziany po cichu. Plik służy wyłącznie sprawdzianom
 * `*.test.ts` tego modułu i nie jest importowany przez żadne okno.
 *
 * Atrapa zapisuje NAZWY wywołanych czynności w kolejności wywołania. To jest
 * jej główny sens: sprawdzian pyta, czy z okna naprawdę prowadzi droga do danej
 * komendy — przycisk, który nie woła niczego, jest atrapą po stronie interfejsu
 * i wygląda tak samo jak przycisk działający.
 */
export interface ZrodloWarsztatuProbne extends ZrodloWarsztatu {
  /** Nazwy wywołanych czynności w kolejności wywołania. */
  wywolania: string[];
  /** Ładunki wywołań, po nazwie czynności — do sprawdzenia treści żądania. */
  ladunki: Record<string, unknown[]>;
  /** Gdy prawda, każda czynność wraca odmową — sprawdzian stanu błędu okna. */
  odmawiaj: boolean;
  /** Gdy fałsz, wykaz kontenerów wraca z `engineAvailable: false`. */
  silnikKontenerow: boolean;
  /** Gdy fałsz, nawigacja i analiza wracają z brakiem programu serwera. */
  programyWarsztatu: boolean;
}

export function utworzZrodloWarsztatuProbne(): ZrodloWarsztatuProbne {
  // Stan atrapy stoi osobno od jej metod: metody są przypisywane niżej, a bez
  // tego rozdziału deklaracja musiałaby wymienić trzydzieści trzy czynności
  // w jednym literale i przestałaby być czytelna.
  const stan = {
    wywolania: [] as string[],
    ladunki: {} as Record<string, unknown[]>,
    odmawiaj: false,
    silnikKontenerow: true,
    programyWarsztatu: true,
  };

  /** oddaj zapisuje wywołanie i zwraca odpowiedź albo odmowę. */
  function oddaj<T>(nazwa: string, zadanie: unknown, wynik: T): Promise<Wynik<T>> {
    stan.wywolania.push(nazwa);
    (stan.ladunki[nazwa] ??= []).push(zadanie);
    if (stan.odmawiaj) {
      return Promise.resolve({
        udany: false,
        blad: {
          code: ErrorCode.InternalError,
          message: `atrapa odmawia czynności ${nazwa}`,
          retryable: false,
        },
      });
    }
    return Promise.resolve({ udany: true, wynik });
  }

  const zapis = stan as ZrodloWarsztatuProbne;
  zapis.nawigujDoSymbolu = (zadanie) =>
    oddaj('nawigujDoSymbolu', zadanie, {
      symbols: stan.programyWarsztatu
        ? [{ name: 'Suma', path: 'suma.go', line: 12 }]
        : [],
      serverAvailable: stan.programyWarsztatu,
    });

  zapis.formatuj = (zadanie) =>
    oddaj('formatuj', zadanie, {
      file: { path: zadanie.path, updatedAt: 1_755_000_000_000 },
      changed: true,
      formatter: 'goimports',
    });

  zapis.analiza = (zadanie) =>
    oddaj('analiza', zadanie, {
      diagnostics: stan.programyWarsztatu
        ? [
            {
              path: 'suma.go',
              line: 12,
              severity: ProblemSeverity.Warning,
              message: 'zmienna bez użycia',
            },
          ]
        : [],
      linterAvailable: stan.programyWarsztatu,
    });

  zapis.refaktoryzuj = (zadanie) =>
    oddaj('refaktoryzuj', zadanie, {
      edits: [{ path: zadanie.path, startLine: 10, endLine: 12, newText: '' }],
      applied: false,
      changedPaths: [zadanie.path],
    });

  zapis.wersjePliku = (zadanie) =>
    oddaj('wersjePliku', zadanie, {
      versions: [{ id: 'fver-1', path: zadanie.path, createdAt: 1_755_000_000_000 }],
    });

  zapis.przywrocWersje = (zadanie) =>
    oddaj('przywrocWersje', zadanie, {
      path: 'suma.go',
      content: 'wersja pierwsza',
      updatedAt: 1_755_000_000_000,
    });

  zapis.wykazBudowan = (zadanie) =>
    oddaj('wykazBudowan', zadanie, {
      builds: [
        {
          id: 'build-1',
          windowId: zadanie.windowId,
          task: 'go test ./...',
          status: BuildStatus.Succeeded,
          startedAt: 1_755_000_000_000,
        },
      ],
    });

  zapis.logBudowania = (zadanie) =>
    oddaj('logBudowania', zadanie, { lines: ['krok pierwszy', 'krok drugi'], truncated: true });

  zapis.wynikTestow = (zadanie) =>
    oddaj('wynikTestow', zadanie, {
      results: [{ name: 'TestSuma', status: TestStatus.Passed }],
      passed: 1,
      failed: 0,
      skipped: 0,
    });

  zapis.pokrycie = (zadanie) =>
    oddaj('pokrycie', zadanie, {
      files: [{ path: 'suma.go', statements: 5, covered: 3, percent: 60 }],
      percent: 60,
    });

  zapis.rozpocznijDebugowanie = (zadanie) =>
    oddaj('rozpocznijDebugowanie', zadanie, {
      id: 'dbg-1',
      windowId: zadanie.windowId,
      adapter: 'dlv',
      status: DebugStatus.Running,
      startedAt: 1_755_000_000_000,
    });

  zapis.sterujDebugowaniem = (zadanie) =>
    oddaj('sterujDebugowaniem', zadanie, {
      id: zadanie.sessionId,
      windowId: 'okno-1',
      adapter: 'dlv',
      status: DebugStatus.Stopped,
      startedAt: 1_755_000_000_000,
    });

  zapis.punktPrzerwania = (zadanie) =>
    oddaj('punktPrzerwania', zadanie, {
      breakpoints:
        zadanie.remove === true
          ? []
          : [
              {
                id: 'bpt-1',
                path: zadanie.path,
                line: zadanie.line,
                kind: 'line',
                verified: true,
              },
            ],
    });

  zapis.zakresDebugowania = (zadanie) =>
    oddaj('zakresDebugowania', zadanie, {
      frames: [{ id: '1', name: 'main.main', path: 'main.go', line: 10 }],
      scopes: [{ name: 'Locals', variablesRef: '2' }],
      variables: [{ name: 'suma', value: '10' }],
    });

  zapis.obliczWyrazenie = (zadanie) =>
    oddaj('obliczWyrazenie', zadanie, { value: '10', type: 'int' });

  zapis.zapytanieApi = (zadanie) =>
    oddaj('zapytanieApi', zadanie, {
      status: 201,
      statusText: '201 Created',
      body: '{"stan":"przyjete"}',
      durationMs: 12,
    });

  zapis.zapiszKolekcje = (zadanie) =>
    oddaj('zapiszKolekcje', zadanie, {
      id: 'apic-1',
      windowId: zadanie.windowId,
      name: zadanie.name,
      requests: [],
      updatedAt: 1_755_000_000_000,
    });

  zapis.kolekcje = (zadanie) =>
    oddaj('kolekcje', zadanie, {
      collections: [
        {
          id: 'apic-1',
          windowId: zadanie.windowId,
          name: 'Usługa magazynu',
          requests: [],
          updatedAt: 1_755_000_000_000,
        },
      ],
    });

  zapis.importujOpenapi = (zadanie) =>
    oddaj('importujOpenapi', zadanie, {
      collection: {
        id: 'apic-2',
        windowId: zadanie.windowId,
        name: 'Usługa magazynu',
        requests: [],
        updatedAt: 1_755_000_000_000,
      },
      requestCount: 3,
    });

  zapis.ustawPolaczenie = (zadanie) =>
    oddaj('ustawPolaczenie', zadanie, {
      id: 'conn-1',
      name: zadanie.name,
      engine: zadanie.engine,
      database: zadanie.database,
      readOnly: true,
    });

  zapis.polaczenia = (zadanie) =>
    oddaj('polaczenia', zadanie, {
      connections: [
        {
          id: 'conn-1',
          name: 'magazyn',
          engine: DataEngine.Sqlite,
          database: 'magazyn.sqlite',
          readOnly: true,
        },
      ],
    });

  zapis.schemat = (zadanie) =>
    oddaj('schemat', zadanie, {
      nodes: [{ path: 'towar', name: 'towar', kind: 'table' }],
    });

  zapis.zapytanieDanych = (zadanie) =>
    oddaj('zapytanieDanych', zadanie, {
      columns: ['nazwa'],
      rows: [['młotek']],
      rowCount: 1,
      durationMs: 3,
      truncated: false,
    });

  zapis.migracje = (zadanie) =>
    oddaj('migracje', zadanie, { applied: ['001_klient.sql'], pending: [] });

  zapis.kontenery = (zadanie) =>
    oddaj('kontenery', zadanie, {
      containers: stan.silnikKontenerow
        ? [{ id: 'kon-1', name: 'usluga', status: ContainerStatus.Running }]
        : [],
      engineAvailable: stan.silnikKontenerow,
    });

  zapis.czynnoscKontenera = (zadanie) =>
    oddaj('czynnoscKontenera', zadanie, {
      container: { id: zadanie.containerId, name: 'usluga', status: ContainerStatus.Running },
      output: 'wiersz logu',
    });

  zapis.budujObraz = (zadanie) => oddaj('budujObraz', zadanie, { imageId: 'sha256:abc' });

  zapis.stosUslug = (zadanie) =>
    oddaj('stosUslug', zadanie, {
      services: [{ id: 'kon-2', name: 'baza', status: ContainerStatus.Running }],
      output: 'baza: uruchomiona',
    });

  zapis.zaleznosci = (zadanie) =>
    oddaj('zaleznosci', zadanie, {
      dependencies: [{ name: 'github.com/coder/websocket', version: 'v1.8.15', direct: true }],
      manifest: 'go.mod',
    });

  zapis.skan = (zadanie) =>
    oddaj('skan', zadanie, {
      id: 'scan-1',
      windowId: zadanie.windowId,
      kinds: zadanie.kinds,
      status: BuildStatus.Succeeded,
      findingCount: 1,
      startedAt: 1_755_000_000_000,
    });

  zapis.znaleziska = (zadanie) =>
    oddaj('znaleziska', zadanie, {
      findings: [
        {
          id: 'find-1',
          kind: ScanKind.Secrets,
          severity: ProblemSeverity.Error,
          title: 'klucz dostępowy w pliku konfiguracja.go',
          path: 'konfiguracja.go',
          line: 4,
        },
      ],
    });

  zapis.operacjaKontekstowa = (zadanie) =>
    oddaj('operacjaKontekstowa', zadanie, { result: 'Kod liczy sumę czterech pierwszych liczb.' });

  zapis.warsztat = (zadanie) =>
    oddaj('warsztat', zadanie, {
      programs: [{ program: 'gofmt', present: true, path: '/usr/local/go/bin/gofmt' }],
    });

  return zapis;
}
