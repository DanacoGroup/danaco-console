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
 * Czterdzieści jeden komend obszaru `apps.*` i dwa jego zdarzenia — cały
 * kontrakt modułu.
 *
 * Obszar jest dwukierunkowy: obok trzech komend zapisu (architektura, plik
 * warsztatu, wdrożenie) stoją trzy komendy odczytu — `apps.deployment.list`,
 * `apps.architecture.get`, `apps.workspace.list`. Bez nich moduł po
 * odświeżeniu okna przeglądarki zaczynałby od zera, bo wypełniałyby go tylko
 * zdarzenia bieżącej sesji gniazda, a historia wdrożeń trzymana w bazie
 * znikałaby Operatorowi z oczu. Każdy odczyt oddaje ten sam kształt, którym
 * odpowiada jego komenda zapisu — druga struktura o tym samym bycie byłaby
 * drugą prawdą.
 *
 * `apps.workspace.changed` jest drogą rozgłoszenia dla `apps.workspace.update`:
 * bez niej plik zapisany w jednym oknie nie docierałby do drugiego, dopóki
 * Operator nie odczytałby warsztatu ręcznie. Wchodzi tą samą bramą co
 * `apps.build.changed` — subskrypcją źródła, nie własnym gniazdem.
 *
 * Źródło nie ma własnego stanu: jest warstwą wywołań i sprawdzianu kształtu
 * odpowiedzi. Stan produktu mieszka w `stan-produktu.ts`, żeby pięć okien
 * patrzyło na jeden zbiór, a nie na pięć kopii.
 *
 * Komendy zapisu mają uchwyt w rdzeniu — wpina je `zarejestrujAplikacje`
 * (`adapter_modul_aplikacje_uchwyty.go`), więc wywołanie wraca zwykłą
 * odpowiedzią albo odmową merytoryczną, nie kopertą `apps.unknown`.
 *
 * Droga przez `odmowa-rdzenia.ts` zabezpiecza ścieżkę fail-open: gdyby uchwyt
 * zniknął albo rdzeń nie rozpoznał którejś komendy, obietnica wywołania ma się
 * czym rozstrzygnąć zamiast wisieć bez końca, a okno nazywa odmowę zamiast ją
 * ukryć — koperta zdarzenia odmowy sama korelacji nie rozstrzyga.
 *
 * Pola tożsamości idą do rdzenia tak, jak je wpisano — bez `trim()` po drodze.
 * `trim()` przeglądarki i `strings.TrimSpace` rdzenia nie są tym samym
 * przycięciem i rozjeżdżają się na dwóch znakach: JavaScript zdejmuje U+FEFF
 * (ZWNBSP), którego Go zostawia, a Go zdejmuje U+0085 (NEL), którego JavaScript
 * zostawia. Ma to znaczenie na ścieżce warsztatu, bo klucz
 * `(okno, warstwa, ścieżka)` czyni ją tożsamością pliku: przycięcie po stronie
 * klienta zapisałoby plik pod ścieżką inną niż wpisana albo odrzuciłoby
 * ścieżkę, którą rdzeń przyjmuje. Rdzeń przycina te pola po swojemu
 * (`adapter_modul_aplikacje_wdrozenie.go`: `strings.TrimSpace(z.Path)`
 * rozstrzyga o pustce, a zapisywana jest wartość surowa) i to on jest tu
 * jedyną władzą. O tym, czy pole jest puste, rozstrzyga więc pustka dosłowna;
 * okno zestawia potem wpisane z oddanym i ogłasza różnicę.
 */
export interface ZlecenieArchitektury {
  idOkna: string;
  idArchitektury: string;
  nazwa: string;
  szablon: AppArchitectureTemplate;
  komponenty: readonly AppComponent[];
}

/** Zlecenie zapisu pliku warsztatu — jedna komenda obsługuje obie warstwy. */
export interface ZlecenieWarsztatu {
  idOkna: string;
  warstwa: AppWorkspaceLayer;
  sciezka: string;
  tresc: string;
  idKomponentu: string;
}

/** Zlecenie wdrożenia; cofnięcie to to samo wywołanie z odnośnikiem wersji. */
export interface ZlecenieWdrozenia {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  strategia: AppDeployStrategy;
  wersja: string;
  notatki: string;
  cofnijDo: string;
}

/**
 * Zawężenie odczytu wdrożeń — pola żądania `apps.deployment.list`.
 *
 * Puste środowisko i zerowa granica znaczą „bez zawężenia": pola są w kontrakcie
 * opcjonalne, a granica pominięta bierze granicę rdzenia. Klient nie podstawia
 * tu własnych wartości domyślnych, bo podstawiona granica byłaby cudzą decyzją
 * przebraną za kontrakt.
 */
export interface ZapytanieWdrozen {
  idOkna: string;
  srodowisko: AppDeployEnvironment | '';
  granica: number;
}

/** Metadane produktu — pola żądania `apps.product.save`. */
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

/** Zlecenie zapisu kamienia milowego. */
export interface ZlecenieKamienia {
  idOkna: string;
  idKamienia: string;
  nazwa: string;
  termin: number;
  stan: AppMilestoneStatus | '';
  idEtapow: readonly string[];
}

/** Zlecenie zapisu notatki projektowej kanwy. */
export interface ZlecenieAdnotacji {
  idOkna: string;
  idAdnotacji: string;
  idKomponentu: string;
  zaleznoscZ: string;
  zaleznoscDo: string;
  tresc: string;
}

/** Zlecenie zapytania próbnego do punktu końcowego. */
export interface ZlecenieProby {
  idOkna: string;
  metoda: AppEndpointMethod;
  sciezka: string;
  naglowki: unknown;
  tresc: string;
  srodowisko: AppDeployEnvironment | '';
}

/** Zlecenie zapisu zmiennej środowiskowej — wartość jawna ALBO sekret. */
export interface ZlecenieZmiennej {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  nazwa: string;
  wartosc: string;
  odwolanieSekretu: string;
}

/** Zlecenie nadania domeny środowiska. */
export interface ZlecenieDomeny {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  domena: string;
  wpisyDns: unknown;
}

/** Zlecenie nastawy skalowania; wartość ujemna znaczy „bez zmiany". */
export interface ZlecenieSkalowania {
  idOkna: string;
  srodowisko: AppDeployEnvironment;
  instancje: number;
  minimum: number;
  maksimum: number;
  reguly: unknown;
}

/** Zlecenie publikacji pakietu do rejestru organizacji. */
export interface ZleceniePublikacji {
  idOkna: string;
  idPakietu: string;
  notatki: string;
  widocznosc: AppPackageVisibility | '';
}

/**
 * Cała reszta obszaru — trzydzieści pięć komend dobudowanych do sześciu, od
 * których moduł zaczynał.
 *
 * Jedno źródło, nie trzydzieści pięć: każde wywołanie idzie tą samą drogą
 * (`utworzWywolanieApps`), więc odmowa `apps.unknown` rozstrzyga obietnicę
 * wszędzie tak samo, a okno nazywa brak zamiast wisieć.
 *
 * Sprawdzian kształtu odpowiedzi jest przy każdej komendzie osobny i pyta
 * o pole, którego okno naprawdę używa. Sprawdzanie „czy cokolwiek wróciło"
 * przepuściłoby odpowiedź o kształcie innym niż kontraktowy, a okno wywróciłoby
 * się dopiero przy rysowaniu.
 *
 * Pustka NIE jest tu uszkodzonym kształtem. Wykaz pusty (`total: 0`) i pole
 * opcjonalne bez wartości (`product`, `theme`, `schema`) są odpowiedziami
 * prawdziwymi i znaczą „w tym oknie tego jeszcze nie ma" — dlatego przy nich
 * sprawdzian pyta o tablicę albo przepuszcza wszystko, zamiast żądać obiektu.
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
        // Sprawdzamy tablicę, a nie jej długość: wykaz pusty jest odpowiedzią
        // poprawną i znaczy „rdzeń nie zna wdrożeń tego okna". Pomylenie pustki
        // z uszkodzonym kształtem odebrałoby oknu jedyny stan, w którym wolno mu
        // powiedzieć „nic tu jeszcze nie ma".
        (tresc) => czyTablica(tresc.deployments),
      );
    },

    async architektura(idOkna) {
      const zadanie: AppsArchitectureGetRequest = { windowId: idOkna };
      return sprawdzKsztalt(
        await wywolaj(Command.AppsArchitectureGet, zadanie),
        Command.AppsArchitectureGet,
        // Pole `architecture` jest w kontrakcie opcjonalne — jego brak znaczy
        // „okno nie ma jeszcze żadnej architektury" i jest odpowiedzią udaną.
        // Żądanie kształtu obiektu odrzuciłoby tę odpowiedź jako uszkodzoną
        // i okno ogłosiłoby odmowę tam, gdzie rdzeń rzetelnie powiedział „pusto".
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
        // Pole `product` jest opcjonalne: okno bez zapisanego produktu dostaje
        // odpowiedź udaną i pustą. Żądanie obiektu ogłaszałoby odmowę tam,
        // gdzie rdzeń rzetelnie powiedział „jeszcze nic".
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
      // Pusty łańcuch JEST wartością: zdejmuje przypisanie wykonawcy. Wysyłamy
      // pole zawsze, gdy nie jest `null` — patrz komentarz przy ZlecenieEtapu.
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
        // Układ bez zastrzeżeń oddaje pustą tablicę i to jest wynik najlepszy
        // z możliwych — sprawdzamy tablicę, nie jej długość.
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
        // Odwołanie puste znaczyłoby wytwór bez bajtów — dokładnie ten wzorzec
        // szkody, którego moduł ma nie powtórzyć.
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
        // Usługa, która nie odpowiada, wraca wynikiem o kodzie 0 i powodem —
        // to odpowiedź udana, bo konstruktor zapytań ma pokazać także awarię.
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
      // Kontrakt wyklucza te dwa pola wzajemnie, więc klient nie wysyła obu —
      // wysłanie obu przeniosłoby rozstrzygnięcie na rdzeń i wróciłoby odmową.
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
          // Do rdzenia idzie ODWOŁANIE do klucza, nigdy jego treść: klucz
          // wydawcy leży w warstwie sekretów i przez kontrakt nie przechodzi.
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
