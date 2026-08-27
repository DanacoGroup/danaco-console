import {
  Command,
  type Extension,
  type ExtensionAuthKind,
  type ExtensionBulkAction,
  type ExtensionCollection,
  type ExtensionDefinitionFormat,
  type ExtensionDetail,
  type ExtensionHealth,
  type ExtensionHealthStatus,
  type ExtensionHistoryEntry,
  type ExtensionKind,
  type ExtensionMapping,
  type ExtensionOrigin,
  type ExtensionPermission,
  type ExtensionRejection,
  type ExtensionScanFinding,
  type ExtensionSignature,
  type ExtensionToolEntry,
  type ExtensionToolKind,
  type ExtensionUpdate,
  type ExtensionUsage,
  type ExtensionWebhook,
  type ExtensionWebhookDirection,
  type McpTransport,
  type ProtocolFrame,
  type SecretRef,
  type ExtensionAdminBulkRequest,
  type ExtensionAuditListRequest,
  type ExtensionCollectionListRequest,
  type ExtensionCollectionSaveRequest,
  type ExtensionCredentialBindRequest,
  type ExtensionDefinitionImportRequest,
  type ExtensionHealthCheckRequest,
  type ExtensionHistoryListRequest,
  type ExtensionMappingSaveRequest,
  type ExtensionPackageUploadRequest,
  type ExtensionProtocolLogListRequest,
  type ExtensionRegistryListRequest,
  type ExtensionSearchRequest,
  type ExtensionSecretListRequest,
  type ExtensionSecretShareRequest,
  type ExtensionToolCallRequest,
  type ExtensionToolListRequest,
  type ExtensionTransportSetRequest,
  type ExtensionUpdateCheckRequest,
  type ExtensionUsageGetRequest,
  type ExtensionVersionPinRequest,
  type ExtensionWebhookListRequest,
  type ExtensionWebhookSaveRequest,
  type ExtensionAuditEntry,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

// Dobudowa obszaru extension — trzydzieści dwie komendy stojące obok pięciu z cyklu życia pozycji.

/** Zawężenie wyszukiwarki katalogu — pola żądania `extension.search`: fraza, rodzaj, pochodzenie, tylko zainstalowane, granica i odsunięcie. */
export interface ZapytanieKatalogu {
  fraza: string;
  rodzaj: ExtensionKind | '';
  pochodzenie: ExtensionOrigin | '';
  tylkoZainstalowane: boolean;
  granica: number;
  odsuniecie: number;
}

/** Zlecenie zapisu kolekcji kuratorskiej: identyfikator, nazwa, opis, oznaczenie barwne i wykaz pozycji kolekcji. */
export interface ZlecenieKolekcji {
  idKolekcji: string;
  nazwa: string;
  opis: string;
  oznaczenieBarwne: string;
  idPozycji: readonly string[];
}

/** Zlecenie ustawienia transportu integracji: identyfikator rozszerzenia, rodzaj transportu, adres, polecenie i żądanie sprawdzenia. */
export interface ZlecenieTransportu {
  idRozszerzenia: string;
  transport: McpTransport;
  adres: string;
  polecenie: string;
  sprawdz: boolean;
}

/** Zlecenie powiązania poświadczenia — do rdzenia idzie odwołanie, nie treść: rozszerzenie, sposób uwierzytelnienia, odwołanie i zakresy. */
export interface ZleceniePoswiadczenia {
  idRozszerzenia: string;
  sposob: ExtensionAuthKind;
  odwolanie: string;
  zakresy: readonly string[];
}

/** Zlecenie zapisu webhooka: identyfikator, rozszerzenie, kierunek, adres, wykaz zdarzeń, odwołanie sekretu i znacznik czynności. */
export interface ZlecenieWebhooka {
  idWebhooka: string;
  idRozszerzenia: string;
  kierunek: ExtensionWebhookDirection;
  adres: string;
  zdarzenia: readonly string[];
  odwolanieSekretu: string;
  czynny: boolean;
}

/** Zlecenie importu definicji API: postać zapisu, kod źródłowy, pochodzenie i treść definicji do rozpoznania. */
export interface ZlecenieImportu {
  format: ExtensionDefinitionFormat;
  kod: string;
  zrodlo: string;
  tresc: string;
}

export interface ZrodloDobudowyRozszerzen {
  /** `extension.search` — trafienia po polach, które pozycja niesie. */
  szukaj(
    z: ZapytanieKatalogu,
  ): Promise<Wynik<{ extensions: Extension[]; total: number; suggestions?: string[] }>>;
  /** `extension.detail.get` — pełna metryka pozycji. */
  szczegol(idRozszerzenia: string): Promise<Wynik<{ detail: ExtensionDetail }>>;
  /** `extension.collection.list` — kolekcje kuratorskie. */
  kolekcje(
    idKolekcji: string,
  ): Promise<Wynik<{ collections: ExtensionCollection[]; total: number }>>;
  /** `extension.collection.save` — nazwa, opis, barwa i skład kolekcji. */
  zapiszKolekcje(z: ZlecenieKolekcji): Promise<Wynik<{ collection: ExtensionCollection }>>;
  /** `extension.collection.apply` — grupowe włączenie albo wyłączenie zestawu. */
  zastosujKolekcje(
    idKolekcji: string,
    wlacz: boolean,
  ): Promise<Wynik<{ applied: Extension[]; rejected: ExtensionRejection[] }>>;
  /** `extension.registry.list` — prywatny rejestr organizacji. */
  rejestrOrganizacji(
    fraza: string,
    rodzaj: ExtensionKind | '',
  ): Promise<Wynik<{ extensions: Extension[]; total: number; registryUrl?: string }>>;
  /** `extension.update.check` — dostępne wydania wobec zainstalowanych. */
  aktualizacje(
    idRozszerzenia: string,
  ): Promise<Wynik<{ updates: ExtensionUpdate[]; checkedAt: number }>>;
  /** `extension.package.upload` — przesyłka paczki instalacji Personal. */
  przeslijPaczke(
    nazwaPliku: string,
    trescBase64: string,
    suma: string,
  ): Promise<Wynik<{ uploadRef: string; sizeBytes: number }>>;
  /** `extension.version.pin` — przypięcie wersji; pusta zdejmuje przypięcie. */
  przypnijWersje(
    idRozszerzenia: string,
    wersja: string,
  ): Promise<Wynik<{ extension: Extension; pinnedVersion?: string }>>;
  /** `extension.version.rollback` — powrót do wersji zarejestrowanej. */
  cofnijWersje(
    idRozszerzenia: string,
    wersjaDocelowa: string,
  ): Promise<Wynik<{ extension: Extension; rolledBackFromVersion: string }>>;
  /** `extension.bundle.install` — odtworzenie środowiska z manifestu zestawu. */
  zainstalujZestaw(
    manifest: unknown,
    wlacz: boolean,
  ): Promise<Wynik<{ installed: Extension[]; rejected: ExtensionRejection[] }>>;
  /** `extension.history.list` — dziennik cyklu życia pozycji. */
  historia(
    idRozszerzenia: string,
    granica: number,
  ): Promise<Wynik<{ entries: ExtensionHistoryEntry[]; total: number }>>;
  /** `extension.admin.bulk` — operacja zbiorcza rejestru. */
  operacjaZbiorcza(
    idPozycji: readonly string[],
    czynnosc: ExtensionBulkAction,
  ): Promise<Wynik<{ affected: Extension[]; rejected: ExtensionRejection[] }>>;

  /** `extension.transport.set` — transport rozmowy z serwerem integracji. */
  ustawTransport(
    z: ZlecenieTransportu,
  ): Promise<Wynik<{ extension: Extension; probeStatus?: ExtensionHealthStatus }>>;
  /** `extension.credential.bind` — powiązanie odwołania do poświadczenia. */
  powiazPoswiadczenie(
    z: ZleceniePoswiadczenia,
  ): Promise<Wynik<{ extension: Extension; authorizationUrl?: string }>>;
  /** `extension.tool.list` — wykaz narzędzi, zasobów i promptów integracji. */
  narzedzia(
    idRozszerzenia: string,
    rodzaj: ExtensionToolKind | '',
    odswiez: boolean,
  ): Promise<Wynik<{ entries: ExtensionToolEntry[]; discoveredAt: number; protocolVersion?: string }>>;
  /** `extension.tool.call` — próbne wywołanie narzędzia z inspektora. */
  wywolajNarzedzie(
    idRozszerzenia: string,
    narzedzie: string,
    argumenty: unknown,
  ): Promise<Wynik<{ ok: boolean; text?: string; durationMs: number; errorDetail?: string }>>;
  /** `extension.protocol.log.list` — strumień ramek JSON-RPC. */
  logProtokolu(
    idRozszerzenia: string,
    granica: number,
  ): Promise<Wynik<{ frames: ProtocolFrame[]; total: number }>>;
  /** `extension.sandbox.run` — przebieg próbny bez podłączania do eksperta. */
  piaskownica(
    idRozszerzenia: string,
    wejscie: unknown,
  ): Promise<Wynik<{ ok: boolean; logRef?: string; durationMs: number }>>;
  /** `extension.definition.import` — konektor zbudowany z opisu API. */
  zaimportujDefinicje(
    z: ZlecenieImportu,
  ): Promise<Wynik<{ extension: Extension; operations: ExtensionToolEntry[] }>>;
  /** `extension.webhook.list` — webhooki przychodzące i wychodzące. */
  webhooki(
    idRozszerzenia: string,
    kierunek: ExtensionWebhookDirection | '',
  ): Promise<Wynik<{ webhooks: ExtensionWebhook[]; total: number }>>;
  /** `extension.webhook.save` — rejestracja albo zmiana webhooka. */
  zapiszWebhook(z: ZlecenieWebhooka): Promise<Wynik<{ webhook: ExtensionWebhook }>>;
  /** `extension.mapping.save` — odwzorowanie i transformacja danych. */
  zapiszMapowanie(
    idRozszerzenia: string,
    idMapowania: string,
    nazwa: string,
    reguly: unknown,
  ): Promise<Wynik<{ mapping: ExtensionMapping; validationIssues?: string[] }>>;
  /** `extension.usage.get` — metryki użycia w oknie czasu. */
  uzycie(
    idRozszerzenia: string,
  ): Promise<Wynik<{ usage: ExtensionUsage[]; windowStart: number; windowEnd: number }>>;
  /** `extension.health.check` — zmierzony stan integracji. */
  kondycja(
    idRozszerzenia: string,
  ): Promise<Wynik<{ results: ExtensionHealth[]; checkedAt: number }>>;
  /** `extension.audit.list` — użycie widziane od strony eksperta. */
  audyt(
    idRozszerzenia: string,
    idEksperta: string,
    granica: number,
  ): Promise<Wynik<{ entries: ExtensionAuditEntry[]; total: number }>>;

  /** `extension.permission.list` — deklaracje, nadania i nadmiar. */
  uprawnienia(idRozszerzenia: string): Promise<
    Wynik<{
      declared?: ExtensionPermission[];
      granted?: ExtensionPermission[];
      excessive?: string[];
    }>
  >;
  /** `extension.permission.grant` — nadanie WYMIENIA komplet uprawnień. */
  nadajUprawnienia(
    idRozszerzenia: string,
    uprawnienia: readonly ExtensionPermission[],
    idEksperta: string,
  ): Promise<Wynik<{ granted: ExtensionPermission[] }>>;
  /** `extension.signature.verify` — weryfikacja podpisu liczona od nowa. */
  zweryfikujPodpis(
    idRozszerzenia: string,
  ): Promise<Wynik<{ signature: ExtensionSignature; checkedAt: number }>>;
  /** `extension.manifest.scan` — skaner manifestu; sygnał, nie brama. */
  skanujManifest(
    idRozszerzenia: string,
  ): Promise<Wynik<{ findings: ExtensionScanFinding[]; scannedAt: number }>>;
  /** `extension.secret.list` — rejestr referencji i przypomnienia o rotacji. */
  sekrety(
    idRozszerzenia: string,
    wCiaguDni: number,
  ): Promise<Wynik<{ secrets: SecretRef[]; total: number }>>;
  /** `extension.secret.share` — zakres współdzielenia referencji sekretu. */
  udostepnijSekret(
    odwolanie: string,
    idPozycji: readonly string[],
    idRol: readonly string[],
  ): Promise<Wynik<{ secret: SecretRef }>>;
}

export function utworzZrodloDobudowyRozszerzen(kanal: Kanal): ZrodloDobudowyRozszerzen {
  const wywolaj: WywolanieApps = utworzWywolanieApps(kanal);

  return {
    async szukaj(z) {
      const zadanie: ExtensionSearchRequest = { query: z.fraza };
      if (z.rodzaj !== '') zadanie.kind = z.rodzaj;
      if (z.pochodzenie !== '') zadanie.origin = z.pochodzenie;
      if (z.tylkoZainstalowane) zadanie.installedOnly = true;
      if (z.granica > 0) zadanie.limit = z.granica;
      if (z.odsuniecie > 0) zadanie.offset = z.odsuniecie;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionSearch, zadanie),
        Command.ExtensionSearch,
        (tresc) => czyTablica(tresc.extensions),
      );
    },

    async szczegol(idRozszerzenia) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionDetailGet, { extensionId: idRozszerzenia }),
        Command.ExtensionDetailGet,
        (tresc) => czyObiekt(tresc.detail),
      );
    },

    async kolekcje(idKolekcji) {
      const zadanie: ExtensionCollectionListRequest = {};
      if (idKolekcji !== '') zadanie.collectionId = idKolekcji;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionCollectionList, zadanie),
        Command.ExtensionCollectionList,
        (tresc) => czyTablica(tresc.collections),
      );
    },

    async zapiszKolekcje(z) {
      const zadanie: ExtensionCollectionSaveRequest = {
        name: z.nazwa,
        extensionIds: [...z.idPozycji],
      };
      if (z.idKolekcji !== '') zadanie.collectionId = z.idKolekcji;
      if (z.opis !== '') zadanie.description = z.opis;
      if (z.oznaczenieBarwne !== '') zadanie.colorTag = z.oznaczenieBarwne;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionCollectionSave, zadanie),
        Command.ExtensionCollectionSave,
        (tresc) => czyObiekt(tresc.collection),
      );
    },

    async zastosujKolekcje(idKolekcji, wlacz) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionCollectionApply, {
          collectionId: idKolekcji,
          enable: wlacz,
        }),
        Command.ExtensionCollectionApply,
        // Odrzucenia są częścią wyniku, nie odmową: zestaw robi, ile się da, i mówi, czego nie zrobił.
        (tresc) => czyTablica(tresc.applied) && czyTablica(tresc.rejected),
      );
    },

    async rejestrOrganizacji(fraza, rodzaj) {
      const zadanie: ExtensionRegistryListRequest = {};
      if (fraza !== '') zadanie.query = fraza;
      if (rodzaj !== '') zadanie.kind = rodzaj;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionRegistryList, zadanie),
        Command.ExtensionRegistryList,
        (tresc) => czyTablica(tresc.extensions),
      );
    },

    async aktualizacje(idRozszerzenia) {
      const zadanie: ExtensionUpdateCheckRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionUpdateCheck, zadanie),
        Command.ExtensionUpdateCheck,
        (tresc) => czyTablica(tresc.updates),
      );
    },

    async przeslijPaczke(nazwaPliku, trescBase64, suma) {
      const zadanie: ExtensionPackageUploadRequest = {
        fileName: nazwaPliku,
        contentBase64: trescBase64,
      };
      if (suma !== '') zadanie.checksumSha256 = suma;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionPackageUpload, zadanie),
        Command.ExtensionPackageUpload,
        // Odwołanie puste znaczyłoby przesyłkę bez bajtów.
        (tresc) => typeof tresc.uploadRef === 'string' && tresc.uploadRef !== '',
      );
    },

    async przypnijWersje(idRozszerzenia, wersja) {
      const zadanie: ExtensionVersionPinRequest = { extensionId: idRozszerzenia };
      // Pole pominięte zdejmuje przypięcie, więc pusty łańcuch nie jedzie na drut jako wartość.
      if (wersja !== '') zadanie.version = wersja;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionVersionPin, zadanie),
        Command.ExtensionVersionPin,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async cofnijWersje(idRozszerzenia, wersjaDocelowa) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionVersionRollback, {
          extensionId: idRozszerzenia,
          targetVersion: wersjaDocelowa,
        }),
        Command.ExtensionVersionRollback,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async zainstalujZestaw(manifest, wlacz) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionBundleInstall, { manifest, enable: wlacz }),
        Command.ExtensionBundleInstall,
        (tresc) => czyTablica(tresc.installed) && czyTablica(tresc.rejected),
      );
    },

    async historia(idRozszerzenia, granica) {
      const zadanie: ExtensionHistoryListRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionHistoryList, zadanie),
        Command.ExtensionHistoryList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async operacjaZbiorcza(idPozycji, czynnosc) {
      const zadanie: ExtensionAdminBulkRequest = {
        extensionIds: [...idPozycji],
        action: czynnosc,
      };
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionAdminBulk, zadanie),
        Command.ExtensionAdminBulk,
        (tresc) => czyTablica(tresc.affected) && czyTablica(tresc.rejected),
      );
    },

    async ustawTransport(z) {
      const zadanie: ExtensionTransportSetRequest = {
        extensionId: z.idRozszerzenia,
        transport: z.transport,
      };
      if (z.adres !== '') zadanie.endpoint = z.adres;
      if (z.polecenie !== '') zadanie.command = z.polecenie;
      if (z.sprawdz) zadanie.probe = true;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionTransportSet, zadanie),
        Command.ExtensionTransportSet,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async powiazPoswiadczenie(z) {
      const zadanie: ExtensionCredentialBindRequest = {
        extensionId: z.idRozszerzenia,
        authKind: z.sposob,
        // Do rdzenia idzie odwołanie, nigdy treść poświadczenia — hasło zostaje po stronie Operatora.
        credentialRef: z.odwolanie,
      };
      if (z.zakresy.length > 0) zadanie.scopes = [...z.zakresy];
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionCredentialBind, zadanie),
        Command.ExtensionCredentialBind,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async narzedzia(idRozszerzenia, rodzaj, odswiez) {
      const zadanie: ExtensionToolListRequest = { extensionId: idRozszerzenia };
      if (rodzaj !== '') zadanie.kind = rodzaj;
      if (odswiez) zadanie.refresh = true;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionToolList, zadanie),
        Command.ExtensionToolList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async wywolajNarzedzie(idRozszerzenia, narzedzie, argumenty) {
      const zadanie: ExtensionToolCallRequest = {
        extensionId: idRozszerzenia,
        toolName: narzedzie,
      };
      if (argumenty !== undefined) zadanie.arguments = argumenty;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionToolCall, zadanie),
        Command.ExtensionToolCall,
        // Niepowodzenie serwera integracji jest odpowiedzią udaną z polem ok:false, nie odmową komendy.
        (tresc) => typeof tresc.ok === 'boolean',
      );
    },

    async logProtokolu(idRozszerzenia, granica) {
      const zadanie: ExtensionProtocolLogListRequest = { extensionId: idRozszerzenia };
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionProtocolLogList, zadanie),
        Command.ExtensionProtocolLogList,
        (tresc) => czyTablica(tresc.frames),
      );
    },

    async piaskownica(idRozszerzenia, wejscie) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionSandboxRun, {
          extensionId: idRozszerzenia,
          input: wejscie,
        }),
        Command.ExtensionSandboxRun,
        (tresc) => typeof tresc.ok === 'boolean',
      );
    },

    async zaimportujDefinicje(z) {
      const zadanie: ExtensionDefinitionImportRequest = { format: z.format, code: z.kod };
      if (z.zrodlo !== '') zadanie.source = z.zrodlo;
      if (z.tresc !== '') zadanie.content = z.tresc;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionDefinitionImport, zadanie),
        Command.ExtensionDefinitionImport,
        (tresc) => czyObiekt(tresc.extension) && czyTablica(tresc.operations),
      );
    },

    async webhooki(idRozszerzenia, kierunek) {
      const zadanie: ExtensionWebhookListRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      if (kierunek !== '') zadanie.direction = kierunek;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionWebhookList, zadanie),
        Command.ExtensionWebhookList,
        (tresc) => czyTablica(tresc.webhooks),
      );
    },

    async zapiszWebhook(z) {
      const zadanie: ExtensionWebhookSaveRequest = {
        extensionId: z.idRozszerzenia,
        direction: z.kierunek,
      };
      if (z.idWebhooka !== '') zadanie.webhookId = z.idWebhooka;
      if (z.adres !== '') zadanie.url = z.adres;
      if (z.zdarzenia.length > 0) zadanie.eventTypes = [...z.zdarzenia];
      if (z.odwolanieSekretu !== '') zadanie.secretRef = z.odwolanieSekretu;
      zadanie.enabled = z.czynny;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionWebhookSave, zadanie),
        Command.ExtensionWebhookSave,
        (tresc) => czyObiekt(tresc.webhook),
      );
    },

    async zapiszMapowanie(idRozszerzenia, idMapowania, nazwa, reguly) {
      const zadanie: ExtensionMappingSaveRequest = {
        extensionId: idRozszerzenia,
        name: nazwa,
        rules: reguly,
      };
      if (idMapowania !== '') zadanie.mappingId = idMapowania;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionMappingSave, zadanie),
        Command.ExtensionMappingSave,
        (tresc) => czyObiekt(tresc.mapping),
      );
    },

    async uzycie(idRozszerzenia) {
      const zadanie: ExtensionUsageGetRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionUsageGet, zadanie),
        Command.ExtensionUsageGet,
        (tresc) => czyTablica(tresc.usage),
      );
    },

    async kondycja(idRozszerzenia) {
      const zadanie: ExtensionHealthCheckRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionHealthCheck, zadanie),
        Command.ExtensionHealthCheck,
        (tresc) => czyTablica(tresc.results),
      );
    },

    async audyt(idRozszerzenia, idEksperta, granica) {
      const zadanie: ExtensionAuditListRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      if (idEksperta !== '') zadanie.agentId = idEksperta;
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionAuditList, zadanie),
        Command.ExtensionAuditList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async uprawnienia(idRozszerzenia) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionPermissionList, { extensionId: idRozszerzenia }),
        Command.ExtensionPermissionList,
        (tresc) => czyTablica(tresc.declared) && czyTablica(tresc.granted),
      );
    },

    async nadajUprawnienia(idRozszerzenia, uprawnienia, idEksperta) {
      const zadanie = {
        extensionId: idRozszerzenia,
        permissions: [...uprawnienia],
        ...(idEksperta === '' ? {} : { agentId: idEksperta }),
      };
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionPermissionGrant, zadanie),
        Command.ExtensionPermissionGrant,
        (tresc) => czyTablica(tresc.granted),
      );
    },

    async zweryfikujPodpis(idRozszerzenia) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionSignatureVerify, { extensionId: idRozszerzenia }),
        Command.ExtensionSignatureVerify,
        (tresc) => czyObiekt(tresc.signature),
      );
    },

    async skanujManifest(idRozszerzenia) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionManifestScan, { extensionId: idRozszerzenia }),
        Command.ExtensionManifestScan,
        // Brak spostrzeżeń jest wynikiem najlepszym z możliwych — sprawdzamy
        // tablicę, nie jej długość.
        (tresc) => czyTablica(tresc.findings),
      );
    },

    async sekrety(idRozszerzenia, wCiaguDni) {
      const zadanie: ExtensionSecretListRequest = {};
      if (idRozszerzenia !== '') zadanie.extensionId = idRozszerzenia;
      if (wCiaguDni > 0) zadanie.expiringWithinDays = wCiaguDni;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionSecretList, zadanie),
        Command.ExtensionSecretList,
        (tresc) => czyTablica(tresc.secrets),
      );
    },

    async udostepnijSekret(odwolanie, idPozycji, idRol) {
      const zadanie: ExtensionSecretShareRequest = { secretRef: odwolanie };
      if (idPozycji.length > 0) zadanie.extensionIds = [...idPozycji];
      if (idRol.length > 0) zadanie.roleIds = [...idRol];
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionSecretShare, zadanie),
        Command.ExtensionSecretShare,
        (tresc) => czyObiekt(tresc.secret),
      );
    },
  };
}
