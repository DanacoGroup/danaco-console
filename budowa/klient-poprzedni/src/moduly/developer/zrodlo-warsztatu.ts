import {
  Command,
  type DeveloperApiCollectionListRequest,
  type DeveloperApiCollectionListResponse,
  type DeveloperApiCollectionSaveRequest,
  type DeveloperApiOpenapiImportRequest,
  type DeveloperApiOpenapiImportResponse,
  type DeveloperApiRequestRequest,
  type DeveloperBreakpointSetRequest,
  type DeveloperBreakpointSetResponse,
  type DeveloperBuildListRequest,
  type DeveloperBuildListResponse,
  type DeveloperBuildLogGetRequest,
  type DeveloperBuildLogGetResponse,
  type DeveloperComposeUpRequest,
  type DeveloperComposeUpResponse,
  type DeveloperContainerActionRequest,
  type DeveloperContainerActionResponse,
  type DeveloperContainerListRequest,
  type DeveloperContainerListResponse,
  type DeveloperContextualOpRequest,
  type DeveloperContextualOpResponse,
  type DeveloperCoverageGetRequest,
  type DeveloperCoverageGetResponse,
  type DeveloperDataConnectionListRequest,
  type DeveloperDataConnectionListResponse,
  type DeveloperDataConnectionSetRequest,
  type DeveloperDataMigrationRunRequest,
  type DeveloperDataMigrationRunResponse,
  type DeveloperDataQueryRunRequest,
  type DeveloperDataSchemaGetRequest,
  type DeveloperDataSchemaGetResponse,
  type DeveloperDebugEvaluateRequest,
  type DeveloperDebugEvaluateResponse,
  type DeveloperDebugScopeGetRequest,
  type DeveloperDebugScopeGetResponse,
  type DeveloperDebugSessionControlRequest,
  type DeveloperDebugSessionStartRequest,
  type DeveloperDependencyListRequest,
  type DeveloperDependencyListResponse,
  type DeveloperFile,
  type DeveloperFileVersionListRequest,
  type DeveloperFileVersionListResponse,
  type DeveloperFileVersionRestoreRequest,
  type DeveloperFormatRunRequest,
  type DeveloperFormatRunResponse,
  type DeveloperImageBuildRequest,
  type DeveloperImageBuildResponse,
  type DeveloperLintGetRequest,
  type DeveloperLintGetResponse,
  type DeveloperRefactorApplyRequest,
  type DeveloperRefactorApplyResponse,
  type DeveloperScanResultListRequest,
  type DeveloperScanResultListResponse,
  type DeveloperScanRunRequest,
  type DeveloperSymbolNavigateRequest,
  type DeveloperSymbolNavigateResponse,
  type DeveloperTestResultGetRequest,
  type DeveloperTestResultGetResponse,
  type DeveloperToolchainCheckRequest,
  type DeveloperToolchainCheckResponse,
  type ApiCollection,
  type ApiResponse,
  type DataConnection,
  type DataQueryResult,
  type DebugSession,
  type ScanRun,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Warsztat modułu Developer widziany przez klienta — trzydzieści trzy komendy
 * obszaru `developer.*` spoza pierwszej piątki okien.
 *
 * ── Dlaczego to jest osobne źródło ─────────────────────────────────────────
 * `ZrodloDeveloper` obsługuje okna wiodące: edytor, drzewo, repozytorium
 * i budowanie. Ten warsztat obsługuje rodziny, które w opracowaniu należą do
 * zakładek Dev Tools, panelu Run & Debug i paska operacji kontekstowych.
 * Rozdział jest po ODBIORCY, nie po wielkości pliku: Project Tree nie ma nic
 * wspólnego z kontenerami, a jedna umowa na wszystko kazałaby każdemu oknu
 * przyjmować zależność od czterdziestu metod, z których używa czterech.
 *
 * ── Każda czynność oddaje `Wynik`, nie samą treść ──────────────────────────
 * Okna mają obowiązkowy stan błędu, więc źródło nie połyka odmowy i nie zwraca
 * w jej miejsce pustego wykazu. Pusty wykaz kontenerów i odmowa odczytu to dwa
 * różne zdania, a Operator ma prawo wiedzieć, które obowiązuje — stąd brak
 * `?? []` w całym pliku.
 *
 * ── Sprawdzenie kształtu ───────────────────────────────────────────────────
 * Każda odpowiedź przechodzi przez `sprawdzKsztalt`. Rdzeń rozminięty
 * z kontraktem nie ma prawa dojść do okna jako `undefined` w środku rysowania:
 * to jest ta klasa usterki, która ujawnia się dopiero u Operatora.
 */
export interface ZrodloWarsztatu {
  // ── Warstwa językowa Code Editora ─────────────────────────────────────────
  /** `developer.symbol.navigate` — przejście do definicji, wystąpień, symboli. */
  nawigujDoSymbolu(
    zadanie: DeveloperSymbolNavigateRequest,
  ): Promise<Wynik<DeveloperSymbolNavigateResponse>>;
  /** `developer.format.run` — formatowanie pliku albo zaznaczenia. */
  formatuj(zadanie: DeveloperFormatRunRequest): Promise<Wynik<DeveloperFormatRunResponse>>;
  /** `developer.lint.get` — zgłoszenia analizy statycznej. */
  analiza(zadanie: DeveloperLintGetRequest): Promise<Wynik<DeveloperLintGetResponse>>;
  /** `developer.refactor.apply` — refaktoryzacja semantyczna. */
  refaktoryzuj(
    zadanie: DeveloperRefactorApplyRequest,
  ): Promise<Wynik<DeveloperRefactorApplyResponse>>;

  // ── Historia pliku edytora ────────────────────────────────────────────────
  /** `developer.file.version.list` — migawki pliku od najnowszej. */
  wersjePliku(
    zadanie: DeveloperFileVersionListRequest,
  ): Promise<Wynik<DeveloperFileVersionListResponse>>;
  /** `developer.file.version.restore` — powrót do migawki. */
  przywrocWersje(zadanie: DeveloperFileVersionRestoreRequest): Promise<Wynik<DeveloperFile>>;

  // ── Odczyt okna Build Output ──────────────────────────────────────────────
  /** `developer.build.list` — historia przebiegów budowania okna. */
  wykazBudowan(zadanie: DeveloperBuildListRequest): Promise<Wynik<DeveloperBuildListResponse>>;
  /** `developer.build.log.get` — log przebiegu wraz z informacją o przycięciu. */
  logBudowania(
    zadanie: DeveloperBuildLogGetRequest,
  ): Promise<Wynik<DeveloperBuildLogGetResponse>>;
  /** `developer.test.result.get` — wynik testów przebiegu. */
  wynikTestow(
    zadanie: DeveloperTestResultGetRequest,
  ): Promise<Wynik<DeveloperTestResultGetResponse>>;
  /** `developer.coverage.get` — pokrycie kodu przebiegu. */
  pokrycie(zadanie: DeveloperCoverageGetRequest): Promise<Wynik<DeveloperCoverageGetResponse>>;

  // ── Run & Debug ───────────────────────────────────────────────────────────
  /** `developer.debug.session.start` — rozpoczęcie sesji debugowania. */
  rozpocznijDebugowanie(
    zadanie: DeveloperDebugSessionStartRequest,
  ): Promise<Wynik<DebugSession>>;
  /** `developer.debug.session.control` — krok, wznowienie albo zatrzymanie. */
  sterujDebugowaniem(
    zadanie: DeveloperDebugSessionControlRequest,
  ): Promise<Wynik<DebugSession>>;
  /** `developer.breakpoint.set` — postawienie albo zdjęcie punktu przerwania. */
  punktPrzerwania(
    zadanie: DeveloperBreakpointSetRequest,
  ): Promise<Wynik<DeveloperBreakpointSetResponse>>;
  /** `developer.debug.scope.get` — stos wywołań, zakresy i zmienne. */
  zakresDebugowania(
    zadanie: DeveloperDebugScopeGetRequest,
  ): Promise<Wynik<DeveloperDebugScopeGetResponse>>;
  /** `developer.debug.evaluate` — wyrażenie w kontekście ramki. */
  obliczWyrazenie(
    zadanie: DeveloperDebugEvaluateRequest,
  ): Promise<Wynik<DeveloperDebugEvaluateResponse>>;

  // ── API Client ────────────────────────────────────────────────────────────
  /** `developer.api.request` — wykonanie zapytania HTTP w sieci serwera. */
  zapytanieApi(zadanie: DeveloperApiRequestRequest): Promise<Wynik<ApiResponse>>;
  /** `developer.api.collection.save` — zapis kolekcji zapytań. */
  zapiszKolekcje(zadanie: DeveloperApiCollectionSaveRequest): Promise<Wynik<ApiCollection>>;
  /** `developer.api.collection.list` — kolekcje zapytań okna. */
  kolekcje(
    zadanie: DeveloperApiCollectionListRequest,
  ): Promise<Wynik<DeveloperApiCollectionListResponse>>;
  /** `developer.api.openapi.import` — kolekcja wytworzona z kontraktu OpenAPI. */
  importujOpenapi(
    zadanie: DeveloperApiOpenapiImportRequest,
  ): Promise<Wynik<DeveloperApiOpenapiImportResponse>>;

  // ── Data Console ──────────────────────────────────────────────────────────
  /** `developer.data.connection.set` — opis połączenia bazodanowego. */
  ustawPolaczenie(zadanie: DeveloperDataConnectionSetRequest): Promise<Wynik<DataConnection>>;
  /** `developer.data.connection.list` — połączenia okna. */
  polaczenia(
    zadanie: DeveloperDataConnectionListRequest,
  ): Promise<Wynik<DeveloperDataConnectionListResponse>>;
  /** `developer.data.schema.get` — drzewo schematu połączenia. */
  schemat(zadanie: DeveloperDataSchemaGetRequest): Promise<Wynik<DeveloperDataSchemaGetResponse>>;
  /** `developer.data.query.run` — wykonanie polecenia SQL. */
  zapytanieDanych(zadanie: DeveloperDataQueryRunRequest): Promise<Wynik<DataQueryResult>>;
  /** `developer.data.migration.run` — migracje schematu bazy Operatora. */
  migracje(
    zadanie: DeveloperDataMigrationRunRequest,
  ): Promise<Wynik<DeveloperDataMigrationRunResponse>>;

  // ── Containers ────────────────────────────────────────────────────────────
  /** `developer.container.list` — kontenery i obrazy silnika. */
  kontenery(
    zadanie: DeveloperContainerListRequest,
  ): Promise<Wynik<DeveloperContainerListResponse>>;
  /** `developer.container.action` — czynność cyklu życia kontenera. */
  czynnoscKontenera(
    zadanie: DeveloperContainerActionRequest,
  ): Promise<Wynik<DeveloperContainerActionResponse>>;
  /** `developer.image.build` — budowanie obrazu z pliku Dockerfile. */
  budujObraz(zadanie: DeveloperImageBuildRequest): Promise<Wynik<DeveloperImageBuildResponse>>;
  /** `developer.compose.up` — podniesienie albo zatrzymanie stosu usług. */
  stosUslug(zadanie: DeveloperComposeUpRequest): Promise<Wynik<DeveloperComposeUpResponse>>;

  // ── Zależności, bezpieczeństwo, jakość ────────────────────────────────────
  /** `developer.dependency.list` — drzewo zależności z manifestu. */
  zaleznosci(
    zadanie: DeveloperDependencyListRequest,
  ): Promise<Wynik<DeveloperDependencyListResponse>>;
  /** `developer.scan.run` — przebieg skanowania repozytorium. */
  skan(zadanie: DeveloperScanRunRequest): Promise<Wynik<ScanRun>>;
  /** `developer.scan.result.list` — znaleziska skanów. */
  znaleziska(
    zadanie: DeveloperScanResultListRequest,
  ): Promise<Wynik<DeveloperScanResultListResponse>>;

  // ── Operacje kontekstowe i sonda warsztatu ────────────────────────────────
  /** `developer.contextual.op` — operacja paska pływającego Code Editora. */
  operacjaKontekstowa(
    zadanie: DeveloperContextualOpRequest,
  ): Promise<Wynik<DeveloperContextualOpResponse>>;
  /** `developer.toolchain.check` — obecność programów warsztatu na serwerze. */
  warsztat(
    zadanie: DeveloperToolchainCheckRequest,
  ): Promise<Wynik<DeveloperToolchainCheckResponse>>;
}

export function utworzZrodloWarsztatu(kanal: Kanal): ZrodloWarsztatu {
  return {
    async nawigujDoSymbolu(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperSymbolNavigate, zadanie),
        Command.DeveloperSymbolNavigate,
        (tresc) => czyTablica(tresc.symbols) && typeof tresc.serverAvailable === 'boolean',
      );
    },

    async formatuj(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperFormatRun, zadanie),
        Command.DeveloperFormatRun,
        (tresc) => czyObiekt(tresc.file) && czyTekst(tresc.file.path),
      );
    },

    async analiza(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperLintGet, zadanie),
        Command.DeveloperLintGet,
        (tresc) => czyTablica(tresc.diagnostics) && typeof tresc.linterAvailable === 'boolean',
      );
    },

    async refaktoryzuj(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperRefactorApply, zadanie),
        Command.DeveloperRefactorApply,
        (tresc) => czyTablica(tresc.edits) && typeof tresc.applied === 'boolean',
      );
    },

    async wersjePliku(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperFileVersionList, zadanie),
        Command.DeveloperFileVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
    },

    async przywrocWersje(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperFileVersionRestore, zadanie),
        Command.DeveloperFileVersionRestore,
        (tresc) => czyObiekt(tresc.file) && czyTekst(tresc.file.path),
      );
      return przenies(wynik, (tresc) => tresc.file);
    },

    async wykazBudowan(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperBuildList, zadanie),
        Command.DeveloperBuildList,
        (tresc) => czyTablica(tresc.builds),
      );
    },

    async logBudowania(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperBuildLogGet, zadanie),
        Command.DeveloperBuildLogGet,
        (tresc) => czyTablica(tresc.lines) && typeof tresc.truncated === 'boolean',
      );
    },

    async wynikTestow(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperTestResultGet, zadanie),
        Command.DeveloperTestResultGet,
        (tresc) =>
          czyTablica(tresc.results) &&
          czyLiczba(tresc.passed) &&
          czyLiczba(tresc.failed) &&
          czyLiczba(tresc.skipped),
      );
    },

    async pokrycie(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperCoverageGet, zadanie),
        Command.DeveloperCoverageGet,
        (tresc) => czyTablica(tresc.files) && czyLiczba(tresc.percent),
      );
    },

    async rozpocznijDebugowanie(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDebugSessionStart, zadanie),
        Command.DeveloperDebugSessionStart,
        (tresc) => czyObiekt(tresc.session) && czyTekst(tresc.session.id),
      );
      return przenies(wynik, (tresc) => tresc.session);
    },

    async sterujDebugowaniem(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDebugSessionControl, zadanie),
        Command.DeveloperDebugSessionControl,
        (tresc) => czyObiekt(tresc.session) && czyTekst(tresc.session.id),
      );
      return przenies(wynik, (tresc) => tresc.session);
    },

    async punktPrzerwania(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperBreakpointSet, zadanie),
        Command.DeveloperBreakpointSet,
        (tresc) => czyTablica(tresc.breakpoints),
      );
    },

    async zakresDebugowania(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDebugScopeGet, zadanie),
        Command.DeveloperDebugScopeGet,
        (tresc) => czyTablica(tresc.frames) && czyTablica(tresc.variables),
      );
    },

    async obliczWyrazenie(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDebugEvaluate, zadanie),
        Command.DeveloperDebugEvaluate,
        (tresc) => czyTekst(tresc.value),
      );
    },

    async zapytanieApi(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperApiRequest, zadanie),
        Command.DeveloperApiRequest,
        (tresc) => czyObiekt(tresc.response) && czyLiczba(tresc.response.status),
      );
      return przenies(wynik, (tresc) => tresc.response);
    },

    async zapiszKolekcje(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperApiCollectionSave, zadanie),
        Command.DeveloperApiCollectionSave,
        (tresc) => czyObiekt(tresc.collection) && czyTekst(tresc.collection.id),
      );
      return przenies(wynik, (tresc) => tresc.collection);
    },

    async kolekcje(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperApiCollectionList, zadanie),
        Command.DeveloperApiCollectionList,
        (tresc) => czyTablica(tresc.collections),
      );
    },

    async importujOpenapi(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperApiOpenapiImport, zadanie),
        Command.DeveloperApiOpenapiImport,
        (tresc) => czyObiekt(tresc.collection) && czyLiczba(tresc.requestCount),
      );
    },

    async ustawPolaczenie(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDataConnectionSet, zadanie),
        Command.DeveloperDataConnectionSet,
        (tresc) => czyObiekt(tresc.connection) && czyTekst(tresc.connection.id),
      );
      return przenies(wynik, (tresc) => tresc.connection);
    },

    async polaczenia(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDataConnectionList, zadanie),
        Command.DeveloperDataConnectionList,
        (tresc) => czyTablica(tresc.connections),
      );
    },

    async schemat(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDataSchemaGet, zadanie),
        Command.DeveloperDataSchemaGet,
        (tresc) => czyTablica(tresc.nodes),
      );
    },

    async zapytanieDanych(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDataQueryRun, zadanie),
        Command.DeveloperDataQueryRun,
        (tresc) => czyObiekt(tresc.result) && czyTablica(tresc.result.columns),
      );
      return przenies(wynik, (tresc) => tresc.result);
    },

    async migracje(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDataMigrationRun, zadanie),
        Command.DeveloperDataMigrationRun,
        (tresc) => czyTablica(tresc.applied) && czyTablica(tresc.pending),
      );
    },

    async kontenery(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperContainerList, zadanie),
        Command.DeveloperContainerList,
        (tresc) => czyTablica(tresc.containers) && typeof tresc.engineAvailable === 'boolean',
      );
    },

    async czynnoscKontenera(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperContainerAction, zadanie),
        Command.DeveloperContainerAction,
        (tresc) => czyObiekt(tresc.container) && czyTekst(tresc.container.id),
      );
    },

    async budujObraz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperImageBuild, zadanie),
        Command.DeveloperImageBuild,
        (tresc) => czyTekst(tresc.imageId),
      );
    },

    async stosUslug(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperComposeUp, zadanie),
        Command.DeveloperComposeUp,
        (tresc) => czyTablica(tresc.services),
      );
    },

    async zaleznosci(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperDependencyList, zadanie),
        Command.DeveloperDependencyList,
        (tresc) => czyTablica(tresc.dependencies) && czyTekst(tresc.manifest),
      );
    },

    async skan(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperScanRun, zadanie),
        Command.DeveloperScanRun,
        (tresc) => czyObiekt(tresc.scan) && czyTekst(tresc.scan.id),
      );
      return przenies(wynik, (tresc) => tresc.scan);
    },

    async znaleziska(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperScanResultList, zadanie),
        Command.DeveloperScanResultList,
        (tresc) => czyTablica(tresc.findings),
      );
    },

    async operacjaKontekstowa(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperContextualOp, zadanie),
        Command.DeveloperContextualOp,
        (tresc) => czyTekst(tresc.result),
      );
    },

    async warsztat(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperToolchainCheck, zadanie),
        Command.DeveloperToolchainCheck,
        (tresc) => czyTablica(tresc.programs),
      );
    },
  };
}
