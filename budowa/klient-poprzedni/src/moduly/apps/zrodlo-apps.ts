import {
  Command,
  EventType,
  type AppAnnotation,
  type AppArchitectureVersion,
  type AppArtifact,
  type AppDeploymentHealth,
  type AppEndpoint,
  type AppEndpointMethod,
  type AppEndpointProbeResult,
  type AppEnvironment,
  type AppEnvironmentVariable,
  type AppExportFormat,
  type AppMilestone,
  type AppMilestoneStatus,
  type AppPackage,
  type AppPackageFormat,
  type AppPackageManifest,
  type AppPackageVisibility,
  type AppPreviewStatus,
  type AppProduct,
  type AppProductLink,
  type AppProductPlatform,
  type AppRoute,
  type AppSchema,
  type AppStage,
  type AppStageStatus,
  type AppTimelineEntry,
  type AppValidationIssue,
  type Extension,
  type ExtensionSignature,
  type AppArchitecture,
  type AppComponent,
  type AppArchitectureTemplate,
  type AppDeployEnvironment,
  type AppDeployStrategy,
  type AppWorkspaceLayer,
  type AppsArchitectureDefineRequest,
  type AppsArchitectureAnnotationSaveRequest,
  type AppsArchitectureGetRequest,
  type AppsArchitectureVersionListRequest,
  type AppsArtifactListRequest,
  type AppsDeploymentDomainSetRequest,
  type AppsDeploymentHealthGetRequest,
  type AppsDeploymentLogReadRequest,
  type AppsDeploymentScaleSetRequest,
  type AppsEndpointListRequest,
  type AppsEndpointProbeRequest,
  type AppsEnvironmentVariableSetRequest,
  type AppsMilestoneListRequest,
  type AppsMilestoneSaveRequest,
  type AppsPackageBuildRequest,
  type AppsPackageManifestSaveRequest,
  type AppsPackagePublishRequest,
  type AppsPreviewStartRequest,
  type AppsProductSaveRequest,
  type AppsSchemaGetRequest,
  type AppsServiceLogReadRequest,
  type AppsStageListRequest,
  type AppsStageSaveRequest,
  type AppsTimelineListRequest,
  type AppsArchitectureGetResponse,
  type AppsBuildChangedEvent,
  type AppsDeploymentListRequest,
  type AppsDeploymentListResponse,
  type AppsDeploymentRunRequest,
  type AppsWorkspaceChangedEvent,
  type AppsWorkspaceListRequest,
  type AppsWorkspaceListResponse,
  type AppsWorkspaceUpdateRequest,
  type AppDeployment,
  type DeveloperFile,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

/**
 * Czterdzieści jeden komend obszaru apps.* i dwa jego zdarzenia stanowią cały
 * kontrakt modułu Apps; komendom zapisu odpowiadają komendy odczytu tego
 * samego kształtu.
 */
export interface ZlecenieArchitektury {
  idOkna: string;
  idArchitektury: string;
  nazwa: string;
  szablon: AppArchitectureTemplate;
  komponenty: readonly AppComponent[];
}

/** Zlecenie zapisu pliku warsztatu jednej z dwóch warstw — jedna komenda kontraktu obsługuje obie warstwy interfejsu i zaplecza. */
export interface ZlecenieWarsztatu {
  idOkna: string;
  warstwa: AppWorkspaceLayer;
  sciezka: string;
  tresc: string;
  idKomponentu: string;
}

/** Zlecenie wdrożenia produktu na środowisko; cofnięcie do wcześniejszej wersji to to samo wywołanie z podanym odnośnikiem wersji. */
export interface ZlecenieWdrozenia {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  strategia: AppDeployStrategy;
  wersja: string;
  notatki: string;
  cofnijDo: string;
}

/** Zawężenie odczytu wdrożeń — pola żądania apps.deployment.list; puste środowisko i zerowa granica znaczą brak zawężenia. */
export interface ZapytanieWdrozen {
  idOkna: string;
  srodowisko: AppDeployEnvironment | '';
  granica: number;
}

/** Metadane produktu modułu Apps — pola żądania apps.product.save: nazwa, opis, platformy docelowe i repozytorium. */
export interface ZlecenieProduktu {
  idOkna: string;
  nazwa: string;
  opis: string;
  platformy: readonly AppProductPlatform[];
  repozytorium: string;
}

/**
 * Zlecenie zapisu etapu.
 *
 * `wykonawca` rozróżnia trzy stany, bo tak każe kontrakt: `null` znaczy „nie
 * ruszaj przypisania", pusty łańcuch — „zdejmij je", a wartość — „przypisz".
 * Jedno pole o dwóch znaczeniach zgubiłoby zdjęcie przypisania.
 */
export interface ZlecenieEtapu {
  idOkna: string;
  idEtapu: string;
  nazwa: string;
  kolejnosc: number;
  stan: AppStageStatus | '';
  wykonawca: string | null;
}

/** Zlecenie zapisu kamienia milowego produktu wraz z terminem, stanem realizacji i powiązanymi etapami budowy trackera. */
export interface ZlecenieKamienia {
  idOkna: string;
  idKamienia: string;
  nazwa: string;
  termin: number;
  stan: AppMilestoneStatus | '';
  idEtapow: readonly string[];
}

/** Zlecenie zapisu notatki projektowej kanwy architektury wraz z jej powiązaniami zależności między komponentami. */
export interface ZlecenieAdnotacji {
  idOkna: string;
  idAdnotacji: string;
  idKomponentu: string;
  zaleznoscZ: string;
  zaleznoscDo: string;
  tresc: string;
}

/** Zlecenie zapytania próbnego do punktu końcowego, wykonywanego naprawdę wobec wskazanego środowiska wdrożeniowego. */
export interface ZlecenieProby {
  idOkna: string;
  metoda: AppEndpointMethod;
  sciezka: string;
  naglowki: unknown;
  tresc: string;
  srodowisko: AppDeployEnvironment | '';
}

/** Zlecenie zapisu zmiennej środowiskowej — wartość jawna albo odwołanie do sekretu, nigdy oba naraz w jednym żądaniu. */
export interface ZlecenieZmiennej {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  nazwa: string;
  wartosc: string;
  odwolanieSekretu: string;
}

/** Zlecenie nadania domeny środowiska wraz z wymaganymi wpisami DNS potwierdzającymi własność podanej domeny. */
export interface ZlecenieDomeny {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  domena: string;
  wpisyDns: unknown;
}

/** Zlecenie nastawy skalowania usługi wdrożonej na środowisku; wartość ujemna znaczy brak zmiany bieżącej nastawy. */
export interface ZlecenieSkalowania {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  instancje: number;
  minimum: number;
  maksimum: number;
  reguly: unknown;
}

/** Zlecenie publikacji pakietu rozszerzenia do prywatnego rejestru organizacji wraz z notatką towarzyszącą wydaniu. */
export interface ZleceniePublikacji {
  idOkna: string;
  idPakietu: string;
  notatki: string;
  widocznosc: AppPackageVisibility | '';
}

/**
 * Cała reszta obszaru apps.* — trzydzieści pięć komend dobudowanych do
 * sześciu początkowych, wywoływanych jednym źródłem przez wspólną drogę
 * wywołania.
 */
export interface ZrodloApps {
  zapiszArchitekture(z: ZlecenieArchitektury): Promise<Wynik<{ architecture: AppArchitecture }>>;
  zapiszPlik(
    z: ZlecenieWarsztatu,
  ): Promise<Wynik<{ layer: AppWorkspaceLayer; file: DeveloperFile }>>;
  uruchomWdrozenie(z: ZlecenieWdrozenia): Promise<Wynik<{ deployment: AppDeployment }>>;
  /** Odczyt wdrożeń okna — odpowiednik odczytu dla `apps.deployment.run`. */
  wdrozenia(z: ZapytanieWdrozen): Promise<Wynik<AppsDeploymentListResponse>>;
  /** Odczyt architektury okna — odpowiednik odczytu dla `apps.architecture.define`. */
  architektura(idOkna: string): Promise<Wynik<AppsArchitectureGetResponse>>;
  /** Odczyt plików warsztatu — odpowiednik odczytu dla `apps.workspace.update`. */
  plikiWarsztatu(
    idOkna: string,
    warstwa: AppWorkspaceLayer | '',
  ): Promise<Wynik<AppsWorkspaceListResponse>>;
  // ── Product Builder ───────────────────────────────────────────────────

  /** `apps.product.get` — metadane produktu okna; brak produktu to pustka. */
  produkt(idOkna: string): Promise<Wynik<{ product?: AppProduct }>>;
  /** `apps.product.save` — nazwa, opis, platformy docelowe i repozytorium. */
  zapiszProdukt(z: ZlecenieProduktu): Promise<Wynik<{ product: AppProduct }>>;
  /** `apps.product.link.list` — stan powiązań z innymi modułami. */
  powiazaniaProduktu(idOkna: string): Promise<Wynik<{ links: AppProductLink[] }>>;
  /** `apps.stage.list` — etapy budowy, opcjonalnie zawężone stanem. */
  etapy(
    idOkna: string,
    stan: AppStageStatus | '',
  ): Promise<Wynik<{ stages: AppStage[]; total: number }>>;
  /** `apps.stage.save` — założenie albo zmiana etapu trackera. */
  zapiszEtap(z: ZlecenieEtapu): Promise<Wynik<{ stage: AppStage }>>;
  /** `apps.milestone.list` — kamienie milowe, opcjonalnie zawężone stanem. */
  kamienieMilowe(
    idOkna: string,
    stan: AppMilestoneStatus | '',
  ): Promise<Wynik<{ milestones: AppMilestone[]; total: number }>>;
  /** `apps.milestone.save` — założenie albo zmiana kamienia milowego. */
  zapiszKamien(z: ZlecenieKamienia): Promise<Wynik<{ milestone: AppMilestone }>>;
  /** `apps.milestone.delete` — usunięcie kamienia milowego. */
  usunKamien(idOkna: string, idKamienia: string): Promise<Wynik<{ deleted: boolean }>>;
  /** `apps.timeline.list` — chronologia projektu z pięciu źródeł rdzenia. */
  osCzasu(
    idOkna: string,
    granica: number,
  ): Promise<Wynik<{ entries: AppTimelineEntry[]; total: number }>>;

  // ── Architecture Designer ─────────────────────────────────────────────

  /** `apps.architecture.validate` — zastrzeżenia układu; żadne nie blokuje. */
  sprawdzArchitekture(
    idOkna: string,
  ): Promise<Wynik<{ issues: AppValidationIssue[]; validatedAt: number }>>;
  /** `apps.architecture.version.list` — historia wersji układu. */
  wersjeArchitektury(
    idOkna: string,
    granica: number,
  ): Promise<Wynik<{ versions: AppArchitectureVersion[]; total: number }>>;
  /** `apps.architecture.annotation.save` — notatka przypięta do kanwy. */
  zapiszAdnotacje(z: ZlecenieAdnotacji): Promise<Wynik<{ annotation: AppAnnotation }>>;
  /** `apps.architecture.export` — plik diagramu w magazynie rdzenia. */
  wyeksportujArchitekture(
    idOkna: string,
    format: AppExportFormat,
  ): Promise<Wynik<{ artifactRef: string; sizeBytes: number }>>;

  // ── Frontend i Backend Workspace ──────────────────────────────────────

  /** `apps.preview.start` — podniesienie serwera podglądu warstwy. */
  uruchomPodglad(
    idOkna: string,
    warstwa: AppWorkspaceLayer | '',
  ): Promise<Wynik<{ previewUrl: string; status: AppPreviewStatus; startedAt: number }>>;
  /** `apps.preview.stop` — zatrzymanie serwera podglądu okna. */
  zatrzymajPodglad(idOkna: string): Promise<Wynik<{ stopped: boolean }>>;
  /** `apps.route.list` — mapa routingu odczytana z plików warstwy interfejsu. */
  trasy(idOkna: string): Promise<Wynik<{ routes: AppRoute[]; total: number }>>;
  /** `apps.theme.get` — motyw produktu; brak motywu to pustka. */
  motyw(idOkna: string): Promise<Wynik<{ theme?: unknown }>>;
  /** `apps.theme.set` — zapis motywu produktu. */
  ustawMotyw(idOkna: string, motyw: unknown): Promise<Wynik<{ theme: unknown }>>;
  /** `apps.endpoint.list` — punkty końcowe z kontraktów API komponentów. */
  punktyKoncowe(
    idOkna: string,
    idKomponentu: string,
  ): Promise<Wynik<{ endpoints: AppEndpoint[]; total: number }>>;
  /** `apps.endpoint.probe` — prawdziwe zapytanie testowe do punktu końcowego. */
  zapytajPunkt(z: ZlecenieProby): Promise<Wynik<{ result: AppEndpointProbeResult }>>;
  /** `apps.schema.get` — schemat bazy produktu; brak schematu to pustka. */
  schemat(idOkna: string, idKomponentu: string): Promise<Wynik<{ schema?: AppSchema }>>;

  // ── Deployment Panel ──────────────────────────────────────────────────

  /** `apps.environment.list` — środowiska wdrożeniowe produktu. */
  srodowiska(
    idOkna: string,
  ): Promise<Wynik<{ environments: AppEnvironment[]; total: number }>>;
  /** `apps.environment.variable.list` — zmienne jednego środowiska. */
  zmienneSrodowiska(
    idOkna: string,
    srodowisko: AppDeployEnvironment,
  ): Promise<Wynik<{ variables: AppEnvironmentVariable[]; total: number }>>;
  /** `apps.environment.variable.set` — wartość jawna ALBO odwołanie do sekretu. */
  ustawZmienna(z: ZlecenieZmiennej): Promise<Wynik<{ variable: AppEnvironmentVariable }>>;
  /** `apps.deployment.domain.set` — domena środowiska wraz z wpisami DNS. */
  ustawDomene(z: ZlecenieDomeny): Promise<Wynik<{ domain: string; verified: boolean }>>;
  /** `apps.deployment.scale.set` — nastawa skalowania usługi. */
  ustawSkalowanie(
    z: ZlecenieSkalowania,
  ): Promise<Wynik<{ applied: boolean; effectiveInstances?: number }>>;
  /** `apps.deployment.health.get` — zmierzony stan wdrożonego produktu. */
  kondycja(
    idOkna: string,
    srodowisko: AppDeployEnvironment | '',
  ): Promise<Wynik<{ health: AppDeploymentHealth }>>;

  // ── Dzienniki i artefakty ─────────────────────────────────────────────

  /** `apps.service.log.read` — dziennik usług okna. */
  dziennikUslugi(
    idOkna: string,
    idKomponentu: string,
    granica: number,
  ): Promise<Wynik<{ lines: string[]; total: number; streaming: boolean }>>;
  /** `apps.deployment.log.read` — dziennik jednego przebiegu wdrożenia. */
  dziennikWdrozenia(
    idOkna: string,
    idWdrozenia: string,
    granica: number,
  ): Promise<Wynik<{ lines: string[]; total: number; streaming: boolean }>>;
  /** `apps.artifact.list` — artefakty budowania wraz z ich rozmiarem i sumą. */
  artefakty(
    idOkna: string,
    idWdrozenia: string,
  ): Promise<Wynik<{ artifacts: AppArtifact[]; total: number }>>;

  // ── Publisher Panel ───────────────────────────────────────────────────

  /** `apps.package.build` — złożenie archiwum pakietu z artefaktu. */
  zbudujPakiet(
    idOkna: string,
    odwolanieArtefaktu: string,
    format: AppPackageFormat | '',
  ): Promise<Wynik<{ package: AppPackage }>>;
  /** `apps.package.manifest.save` — tożsamość, wersja, narzędzia, uprawnienia. */
  zapiszManifest(
    idOkna: string,
    idPakietu: string,
    manifest: AppPackageManifest,
  ): Promise<Wynik<{ package: AppPackage }>>;
  /** `apps.package.validate` — raport zgodności z kontraktem rozszerzenia. */
  sprawdzPakiet(
    idOkna: string,
    idPakietu: string,
  ): Promise<Wynik<{ issues: AppValidationIssue[]; validatedAt: number }>>;
  /** `apps.package.sign` — podpis Ed25519 kluczem wskazanym odwołaniem. */
  podpiszPakiet(
    idOkna: string,
    idPakietu: string,
    odwolanieKlucza: string,
  ): Promise<Wynik<{ signature: ExtensionSignature }>>;
  /** `apps.package.publish` — pozycja w prywatnym rejestrze organizacji. */
  opublikujPakiet(
    z: ZleceniePublikacji,
  ): Promise<Wynik<{ extension: Extension; publishedAt: number }>>;

  /** Subskrypcja `apps.build.changed` — zdarzenia etapów i wdrożeń. */
  naZmianeBudowy(sluchacz: (tresc: AppsBuildChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja `apps.workspace.changed` — zdarzenia pliku warsztatu. */
  naZmianeWarsztatu(sluchacz: (tresc: AppsWorkspaceChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloApps(kanal: Kanal): ZrodloApps {
  const wywolaj: WywolanieApps = utworzWywolanieApps(kanal);

  return {
    async zapiszArchitekture(z) {
      const zadanie: AppsArchitectureDefineRequest = {
        windowId: z.idOkna,
        template: z.szablon,
        components: [...z.komponenty],
      };
      if (z.idArchitektury !== '') zadanie.architectureId = z.idArchitektury;
      if (z.nazwa !== '') zadanie.name = z.nazwa;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureDefine, zadanie),
        Command.AppsArchitectureDefine,
        (tresc) => czyObiekt(tresc.architecture),
      );
    },

    async zapiszPlik(z) {
      const zadanie: AppsWorkspaceUpdateRequest = {
        windowId: z.idOkna,
        layer: z.warstwa,
        path: z.sciezka,
        content: z.tresc,
      };
      if (z.idKomponentu !== '') zadanie.componentId = z.idKomponentu;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsWorkspaceUpdate, zadanie),
        Command.AppsWorkspaceUpdate,
        (tresc) => czyObiekt(tresc.file),
      );
    },

    async uruchomWdrozenie(z) {
      const zadanie: AppsDeploymentRunRequest = {
        windowId: z.idOkna,
        environment: z.srodowisko,
        strategy: z.strategia,
      };
      if (z.wersja !== '') zadanie.version = z.wersja;
      if (z.notatki !== '') zadanie.releaseNotes = z.notatki;
      if (z.cofnijDo !== '') zadanie.rollbackToDeploymentId = z.cofnijDo;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsDeploymentRun, zadanie),
        Command.AppsDeploymentRun,
        (tresc) => czyObiekt(tresc.deployment),
      );
    },

    async wdrozenia(z) {
      const zadanie: AppsDeploymentListRequest = { windowId: z.idOkna };
      if (z.srodowisko !== '') zadanie.environment = z.srodowisko;
      if (z.granica > 0) zadanie.limit = z.granica;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsDeploymentList, zadanie),
        Command.AppsDeploymentList,
        // Sprawdza się tablicę, nie jej długość: wykaz pusty to odpowiedź poprawna, nie kształt uszkodzony.
        (tresc) => czyTablica(tresc.deployments),
      );
    },

    async architektura(idOkna) {
      const zadanie: AppsArchitectureGetRequest = { windowId: idOkna };
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureGet, zadanie),
        Command.AppsArchitectureGet,
        // Pole architecture jest opcjonalne; jego brak to odpowiedź udana, nie kształt uszkodzony.
        () => true,
      );
    },

    async plikiWarsztatu(idOkna, warstwa) {
      const zadanie: AppsWorkspaceListRequest = { windowId: idOkna };
      if (warstwa !== '') zadanie.layer = warstwa;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsWorkspaceList, zadanie),
        Command.AppsWorkspaceList,
        (tresc) => czyTablica(tresc.files),
      );
    },

    // ── Product Builder ─────────────────────────────────────────────────

    async produkt(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsProductGet, { windowId: idOkna }),
        Command.AppsProductGet,
        // Pole product jest opcjonalne: brak zapisanego produktu daje odpowiedź udaną i pustą.
        () => true,
      );
    },

    async zapiszProdukt(z) {
      const zadanie: AppsProductSaveRequest = { windowId: z.idOkna, name: z.nazwa };
      if (z.opis !== '') zadanie.description = z.opis;
      if (z.platformy.length > 0) zadanie.platforms = [...z.platformy];
      if (z.repozytorium !== '') zadanie.repositoryUrl = z.repozytorium;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsProductSave, zadanie),
        Command.AppsProductSave,
        (tresc) => czyObiekt(tresc.product),
      );
    },

    async powiazaniaProduktu(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsProductLinkList, { windowId: idOkna }),
        Command.AppsProductLinkList,
        (tresc) => czyTablica(tresc.links),
      );
    },

    async etapy(idOkna, stan) {
      const zadanie: AppsStageListRequest = { windowId: idOkna };
      if (stan !== '') zadanie.status = stan;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsStageList, zadanie),
        Command.AppsStageList,
        (tresc) => czyTablica(tresc.stages),
      );
    },

    async zapiszEtap(z) {
      const zadanie: AppsStageSaveRequest = { windowId: z.idOkna };
      if (z.idEtapu !== '') zadanie.stageId = z.idEtapu;
      if (z.nazwa !== '') zadanie.name = z.nazwa;
      if (z.kolejnosc > 0) zadanie.order = z.kolejnosc;
      if (z.stan !== '') zadanie.status = z.stan;
      // Pusty łańcuch jest wartością: zdejmuje wykonawcę; wysyła się pole zawsze, gdy nie jest null.
      if (z.wykonawca !== null) zadanie.ownerAgentId = z.wykonawca;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsStageSave, zadanie),
        Command.AppsStageSave,
        (tresc) => czyObiekt(tresc.stage),
      );
    },

    async kamienieMilowe(idOkna, stan) {
      const zadanie: AppsMilestoneListRequest = { windowId: idOkna };
      if (stan !== '') zadanie.status = stan;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsMilestoneList, zadanie),
        Command.AppsMilestoneList,
        (tresc) => czyTablica(tresc.milestones),
      );
    },

    async zapiszKamien(z) {
      const zadanie: AppsMilestoneSaveRequest = { windowId: z.idOkna, name: z.nazwa };
      if (z.idKamienia !== '') zadanie.milestoneId = z.idKamienia;
      if (z.termin > 0) zadanie.dueAt = z.termin;
      if (z.stan !== '') zadanie.status = z.stan;
      if (z.idEtapow.length > 0) zadanie.stageIds = [...z.idEtapow];
      return sprawdzKsztalt(
        await wywolaj(Command.AppsMilestoneSave, zadanie),
        Command.AppsMilestoneSave,
        (tresc) => czyObiekt(tresc.milestone),
      );
    },

    async usunKamien(idOkna, idKamienia) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsMilestoneDelete, {
          windowId: idOkna,
          milestoneId: idKamienia,
        }),
        Command.AppsMilestoneDelete,
        (tresc) => typeof tresc.deleted === 'boolean',
      );
    },

    async osCzasu(idOkna, granica) {
      const zadanie: AppsTimelineListRequest = { windowId: idOkna };
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsTimelineList, zadanie),
        Command.AppsTimelineList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    // ── Architecture Designer ───────────────────────────────────────────

    async sprawdzArchitekture(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureValidate, { windowId: idOkna }),
        Command.AppsArchitectureValidate,
        // Układ bez zastrzeżeń oddaje pustą tablicę — wynik najlepszy; sprawdza się tablicę, nie jej długość.
        (tresc) => czyTablica(tresc.issues),
      );
    },

    async wersjeArchitektury(idOkna, granica) {
      const zadanie: AppsArchitectureVersionListRequest = { windowId: idOkna };
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureVersionList, zadanie),
        Command.AppsArchitectureVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
    },

    async zapiszAdnotacje(z) {
      const zadanie: AppsArchitectureAnnotationSaveRequest = {
        windowId: z.idOkna,
        text: z.tresc,
      };
      if (z.idAdnotacji !== '') zadanie.annotationId = z.idAdnotacji;
      if (z.idKomponentu !== '') zadanie.componentId = z.idKomponentu;
      if (z.zaleznoscZ !== '') zadanie.dependencyFrom = z.zaleznoscZ;
      if (z.zaleznoscDo !== '') zadanie.dependencyTo = z.zaleznoscDo;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureAnnotationSave, zadanie),
        Command.AppsArchitectureAnnotationSave,
        (tresc) => czyObiekt(tresc.annotation),
      );
    },

    async wyeksportujArchitekture(idOkna, format) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureExport, { windowId: idOkna, format }),
        Command.AppsArchitectureExport,
        // Odwołanie puste znaczyłoby wytwór bez bajtów — wzorzec szkody, którego moduł ma nie powtórzyć.
        (tresc) => typeof tresc.artifactRef === 'string' && tresc.artifactRef !== '',
      );
    },

    // ── Frontend i Backend Workspace ────────────────────────────────────

    async uruchomPodglad(idOkna, warstwa) {
      const zadanie: AppsPreviewStartRequest = { windowId: idOkna };
      if (warstwa !== '') zadanie.layer = warstwa;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPreviewStart, zadanie),
        Command.AppsPreviewStart,
        (tresc) => typeof tresc.previewUrl === 'string' && tresc.previewUrl !== '',
      );
    },

    async zatrzymajPodglad(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPreviewStop, { windowId: idOkna }),
        Command.AppsPreviewStop,
        (tresc) => typeof tresc.stopped === 'boolean',
      );
    },

    async trasy(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsRouteList, { windowId: idOkna }),
        Command.AppsRouteList,
        (tresc) => czyTablica(tresc.routes),
      );
    },

    async motyw(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsThemeGet, { windowId: idOkna }),
        Command.AppsThemeGet,
        () => true,
      );
    },

    async ustawMotyw(idOkna, motyw) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsThemeSet, { windowId: idOkna, theme: motyw }),
        Command.AppsThemeSet,
        () => true,
      );
    },

    async punktyKoncowe(idOkna, idKomponentu) {
      const zadanie: AppsEndpointListRequest = { windowId: idOkna };
      if (idKomponentu !== '') zadanie.componentId = idKomponentu;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsEndpointList, zadanie),
        Command.AppsEndpointList,
        (tresc) => czyTablica(tresc.endpoints),
      );
    },

    async zapytajPunkt(z) {
      const zadanie: AppsEndpointProbeRequest = {
        windowId: z.idOkna,
        method: z.metoda,
        path: z.sciezka,
      };
      if (z.naglowki !== undefined) zadanie.headers = z.naglowki;
      if (z.tresc !== '') zadanie.body = z.tresc;
      if (z.srodowisko !== '') zadanie.environment = z.srodowisko;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsEndpointProbe, zadanie),
        Command.AppsEndpointProbe,
        // Usługa, która nie odpowiada, wraca kodem 0 i powodem — to odpowiedź udana, pokazująca awarię.
        (tresc) => czyObiekt(tresc.result),
      );
    },

    async schemat(idOkna, idKomponentu) {
      const zadanie: AppsSchemaGetRequest = { windowId: idOkna };
      if (idKomponentu !== '') zadanie.componentId = idKomponentu;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsSchemaGet, zadanie),
        Command.AppsSchemaGet,
        () => true,
      );
    },

    // ── Deployment Panel ────────────────────────────────────────────────

    async srodowiska(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsEnvironmentList, { windowId: idOkna }),
        Command.AppsEnvironmentList,
        (tresc) => czyTablica(tresc.environments),
      );
    },

    async zmienneSrodowiska(idOkna, srodowisko) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsEnvironmentVariableList, {
          windowId: idOkna,
          environment: srodowisko,
        }),
        Command.AppsEnvironmentVariableList,
        (tresc) => czyTablica(tresc.variables),
      );
    },

    async ustawZmienna(z) {
      const zadanie: AppsEnvironmentVariableSetRequest = {
        windowId: z.idOkna,
        environment: z.srodowisko,
        name: z.nazwa,
      };
      // Kontrakt wyklucza te dwa pola wzajemnie; wysłanie obu przeniosłoby rozstrzygnięcie na rdzeń.
      if (z.odwolanieSekretu !== '') zadanie.secretRef = z.odwolanieSekretu;
      else zadanie.value = z.wartosc;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsEnvironmentVariableSet, zadanie),
        Command.AppsEnvironmentVariableSet,
        (tresc) => czyObiekt(tresc.variable),
      );
    },

    async ustawDomene(z) {
      const zadanie: AppsDeploymentDomainSetRequest = {
        windowId: z.idOkna,
        environment: z.srodowisko,
        domain: z.domena,
      };
      if (z.wpisyDns !== undefined) zadanie.dnsRecords = z.wpisyDns;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsDeploymentDomainSet, zadanie),
        Command.AppsDeploymentDomainSet,
        (tresc) => typeof tresc.verified === 'boolean',
      );
    },

    async ustawSkalowanie(z) {
      const zadanie: AppsDeploymentScaleSetRequest = {
        windowId: z.idOkna,
        environment: z.srodowisko,
      };
      if (z.instancje >= 0) zadanie.instances = z.instancje;
      if (z.minimum >= 0) zadanie.minInstances = z.minimum;
      if (z.maksimum >= 0) zadanie.maxInstances = z.maksimum;
      if (z.reguly !== undefined) zadanie.rules = z.reguly;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsDeploymentScaleSet, zadanie),
        Command.AppsDeploymentScaleSet,
        (tresc) => typeof tresc.applied === 'boolean',
      );
    },

    async kondycja(idOkna, srodowisko) {
      const zadanie: AppsDeploymentHealthGetRequest = { windowId: idOkna };
      if (srodowisko !== '') zadanie.environment = srodowisko;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsDeploymentHealthGet, zadanie),
        Command.AppsDeploymentHealthGet,
        (tresc) => czyObiekt(tresc.health),
      );
    },

    // ── Dzienniki i artefakty ───────────────────────────────────────────

    async dziennikUslugi(idOkna, idKomponentu, granica) {
      const zadanie: AppsServiceLogReadRequest = { windowId: idOkna };
      if (idKomponentu !== '') zadanie.componentId = idKomponentu;
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsServiceLogRead, zadanie),
        Command.AppsServiceLogRead,
        (tresc) => czyTablica(tresc.lines),
      );
    },

    async dziennikWdrozenia(idOkna, idWdrozenia, granica) {
      const zadanie: AppsDeploymentLogReadRequest = {
        windowId: idOkna,
        deploymentId: idWdrozenia,
      };
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsDeploymentLogRead, zadanie),
        Command.AppsDeploymentLogRead,
        (tresc) => czyTablica(tresc.lines),
      );
    },

    async artefakty(idOkna, idWdrozenia) {
      const zadanie: AppsArtifactListRequest = { windowId: idOkna };
      if (idWdrozenia !== '') zadanie.deploymentId = idWdrozenia;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArtifactList, zadanie),
        Command.AppsArtifactList,
        (tresc) => czyTablica(tresc.artifacts),
      );
    },

    // ── Publisher Panel ─────────────────────────────────────────────────

    async zbudujPakiet(idOkna, odwolanieArtefaktu, format) {
      const zadanie: AppsPackageBuildRequest = { windowId: idOkna };
      if (odwolanieArtefaktu !== '') zadanie.artifactRef = odwolanieArtefaktu;
      if (format !== '') zadanie.format = format;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPackageBuild, zadanie),
        Command.AppsPackageBuild,
        (tresc) => czyObiekt(tresc.package),
      );
    },

    async zapiszManifest(idOkna, idPakietu, manifest) {
      const zadanie: AppsPackageManifestSaveRequest = { windowId: idOkna, manifest };
      if (idPakietu !== '') zadanie.packageId = idPakietu;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPackageManifestSave, zadanie),
        Command.AppsPackageManifestSave,
        (tresc) => czyObiekt(tresc.package),
      );
    },

    async sprawdzPakiet(idOkna, idPakietu) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPackageValidate, {
          windowId: idOkna,
          packageId: idPakietu,
        }),
        Command.AppsPackageValidate,
        (tresc) => czyTablica(tresc.issues),
      );
    },

    async podpiszPakiet(idOkna, idPakietu, odwolanieKlucza) {
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPackageSign, {
          windowId: idOkna,
          packageId: idPakietu,
          // Do rdzenia idzie odwołanie do klucza, nie treść; klucz leży w warstwie sekretów.
          signingKeyRef: odwolanieKlucza,
        }),
        Command.AppsPackageSign,
        (tresc) => czyObiekt(tresc.signature),
      );
    },

    async opublikujPakiet(z) {
      const zadanie: AppsPackagePublishRequest = {
        windowId: z.idOkna,
        packageId: z.idPakietu,
      };
      if (z.notatki !== '') zadanie.releaseNotes = z.notatki;
      if (z.widocznosc !== '') zadanie.visibility = z.widocznosc;
      return sprawdzKsztalt(
        await wywolaj(Command.AppsPackagePublish, zadanie),
        Command.AppsPackagePublish,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    naZmianeBudowy(sluchacz) {
      return kanal.naZdarzenie(EventType.AppsBuildChanged, (tresc) => sluchacz(tresc));
    },

    naZmianeWarsztatu(sluchacz) {
      return kanal.naZdarzenie(EventType.AppsWorkspaceChanged, (tresc) => sluchacz(tresc));
    },
  };
}
