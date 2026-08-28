import {
  Command,
  type LibraryAuditEntry,
  type LibraryCollection,
  type LibraryDiff,
  type LibraryDuplicateGroup,
  type LibraryFile,
  type LibraryFixityResult,
  type LibraryMetadata,
  type LibraryNormalizeResult,
  type LibraryPreservationKind,
  type LibraryPreservationResult,
  type LibraryRetentionDue,
  type LibraryRetentionPolicy,
  type LibraryRule,
  type LibrarySuggestion,
  type LibraryTag,
  type LibraryDuplicateKind,
  type LibraryFieldDefinition,
  type LibraryPackageKind,
  type LibraryShare,
  type LibraryStats,
  type LibraryThesaurusRelation,
  type LibraryWebhook,
  type LibraryShareScope,
} from '../../../../shared/contract';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import type { Wynik } from '../../protokol/kanal';
import type { StrazOdmow } from './straz-odmow';

/**
 * Zarząd repozytorium udostępnia oknom czynności zarządcze biblioteki i sprawdza
 * kształt każdej odpowiedzi rdzenia, by nieprzysłane pole nie trafiło do interfejsu.
 */
export interface ZarzadRepozytorium {
  /** Opis zasobu; `zTechnicznymi` sięga po metadane osadzone w bajtach. */
  opis(idPliku: string, zTechnicznymi: boolean): Promise<Wynik<{ metadata: LibraryMetadata }>>;
  zapiszOpis(
    idPliku: string,
    opis: LibraryMetadata,
    podmien: boolean,
  ): Promise<Wynik<{ metadata: LibraryMetadata; file: LibraryFile }>>;

  slownikEtykiet(fraza?: string): Promise<Wynik<{ tags: LibraryTag[]; total: number }>>;
  zmienEtykiete(
    nazwa: string,
    nowaNazwa?: string,
    barwa?: string,
  ): Promise<Wynik<{ tag: LibraryTag; affectedFiles: number }>>;
  polaczEtykiety(
    zrodlowe: readonly string[],
    docelowa: string,
  ): Promise<Wynik<{ tag: LibraryTag; affectedFiles: number; mergedNames: string[] }>>;
  usunEtykiete(
    nazwa: string,
    potwierdzone: boolean,
  ): Promise<Wynik<{ removed: boolean; affectedFiles: number }>>;

  wykazKolekcji(): Promise<Wynik<{ collections: LibraryCollection[]; total: number }>>;
  ustawRelacje(
    zrodlo: string,
    cel: string,
    relacja: LibraryThesaurusRelation,
    zdejmij: boolean,
  ): Promise<Wynik<{ relationSet: boolean }>>;
  wywiezTezaurus(
    format: string,
  ): Promise<Wynik<{ content: string; format: string; conceptCount: number; relationCount: number }>>;

  wykazRegul(): Promise<Wynik<{ rules: LibraryRule[]; total: number }>>;
  zapiszRegule(
    regula: LibraryRule,
    przeliczTeraz: boolean,
  ): Promise<Wynik<{ rule: LibraryRule; matchedFiles?: number }>>;
  usunRegule(
    idReguly: string,
    zdejmijPliki: boolean,
  ): Promise<Wynik<{ removed: boolean; detachedFiles: number }>>;

  skanujDuplikaty(
    rodzaje: readonly LibraryDuplicateKind[],
  ): Promise<Wynik<{ groups: LibraryDuplicateGroup[]; total: number; skippedWithoutChecksum: number }>>;
  polaczDuplikaty(
    docelowy: string,
    zrodlowe: readonly string[],
  ): Promise<Wynik<{ file: LibraryFile; mergedCount: number; carriedVersions: number }>>;
  sprawdzIntegralnosc(): Promise<
    Wynik<{
      results: LibraryFixityResult[];
      checkedCount: number;
      mismatchedCount: number;
      missingCount: number;
    }>
  >;
  normalizujNazwy(
    idPlikow: readonly string[],
    probny: boolean,
  ): Promise<Wynik<{ results: LibraryNormalizeResult[]; changedCount: number }>>;
  pulpitStanu(): Promise<Wynik<{ stats: LibraryStats }>>;
  dziennikAudytu(
    idPliku?: string,
  ): Promise<Wynik<{ entries: LibraryAuditEntry[]; total: number }>>;

  przenies(
    idPlikow: readonly string[],
    sciezka: string,
  ): Promise<Wynik<{ movedCount: number; files: LibraryFile[] }>>;
  zarchiwizuj(
    idPlikow: readonly string[],
    powod: string,
  ): Promise<Wynik<{ archivedCount: number; files: LibraryFile[] }>>;
  przywrocZArchiwum(
    idPlikow: readonly string[],
  ): Promise<Wynik<{ restoredCount: number; files: LibraryFile[] }>>;
  usunTrwale(
    idPlikow: readonly string[],
  ): Promise<Wynik<{ deletedCount: number; deletedVersions: number }>>;

  utrwal(
    idPlikow: readonly string[],
    postac: LibraryPreservationKind,
  ): Promise<Wynik<{ results: LibraryPreservationResult[]; validCount: number; invalidCount: number }>>;
  wywiezPaczke(
    idPlikow: readonly string[],
    postac: LibraryPackageKind,
    format: string,
  ): Promise<Wynik<{ fileId: string; entries: number; sizeBytes: number; manifestChecksum: string }>>;
  wykazRetencji(
    wDniach?: number,
  ): Promise<Wynik<{ policies: LibraryRetentionPolicy[]; due?: LibraryRetentionDue[] }>>;
  zapiszRetencje(
    polityka: LibraryRetentionPolicy,
    usun: boolean,
  ): Promise<Wynik<{ policy: LibraryRetentionPolicy; affectedFiles: number }>>;

  wystawUdostepnienie(
    zakres: LibraryShareScope,
    idCelu: string,
    sekundy?: number,
  ): Promise<Wynik<{ share: LibraryShare; url: string }>>;
  wykazUdostepnien(): Promise<Wynik<{ shares: LibraryShare[]; total: number }>>;
  odwolajUdostepnienie(idUdostepnienia: string): Promise<Wynik<{ revoked: boolean }>>;

  wykazNasluchow(): Promise<Wynik<{ webhooks: LibraryWebhook[]; total: number }>>;
  zapiszNasluch(nasluch: LibraryWebhook): Promise<Wynik<{ webhook: LibraryWebhook }>>;
  usunNasluch(idNasluchu: string): Promise<Wynik<{ removed: boolean }>>;

  klasyfikuj(
    idPlikow: readonly string[],
    zapisz: boolean,
  ): Promise<Wynik<{ suggestions: LibrarySuggestion[]; processedCount: number; appliedCount: number }>>;
  wykazSugestii(
    idPliku?: string,
  ): Promise<Wynik<{ suggestions: LibrarySuggestion[]; total: number }>>;
  rozstrzygnijSugestie(
    idSugestii: readonly string[],
    przyjmij: boolean,
  ): Promise<Wynik<{ appliedCount: number; rejectedCount: number; files?: LibraryFile[] }>>;

  porownaj(
    lewy: string,
    prawy?: string,
    lewaWersja?: string,
    prawaWersja?: string,
  ): Promise<Wynik<{ diff: LibraryDiff }>>;

  /** Definicje pól niestandardowych schematu metadanych. */
  schematMetadanych(): Promise<Wynik<{ fields: LibraryFieldDefinition[] }>>;
  zapiszPoleSchematu(
    pole: LibraryFieldDefinition,
    usun: boolean,
  ): Promise<Wynik<{ field: LibraryFieldDefinition; affectedFiles: number }>>;
}

export function utworzZarzadRepozytorium(straz: StrazOdmow): ZarzadRepozytorium {
  return {
    async opis(idPliku, zTechnicznymi) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryMetadataGet, {
          fileId: idPliku,
          includeTechnical: zTechnicznymi,
        }),
        Command.LibraryMetadataGet,
        (tresc) => czyObiekt(tresc.metadata),
      );
    },

    async zapiszOpis(idPliku, opis, podmien) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryMetadataSet, {
          fileId: idPliku,
          metadata: opis,
          replace: podmien,
        }),
        Command.LibraryMetadataSet,
        (tresc) => czyObiekt(tresc.metadata) && czyObiekt(tresc.file),
      );
    },

    async schematMetadanych() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibrarySchemaGet, {}),
        Command.LibrarySchemaGet,
        (tresc) => czyTablica(tresc.fields),
      );
    },

    async zapiszPoleSchematu(pole, usun) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibrarySchemaSet, { field: pole, remove: usun }),
        Command.LibrarySchemaSet,
        (tresc) => czyObiekt(tresc.field) && typeof tresc.affectedFiles === 'number',
      );
    },

    async slownikEtykiet(fraza) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryTagList, fraza === undefined ? {} : { query: fraza }),
        Command.LibraryTagList,
        (tresc) => czyTablica(tresc.tags) && typeof tresc.total === 'number',
      );
    },

    async zmienEtykiete(nazwa, nowaNazwa, barwa) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryTagUpdate, {
          name: nazwa,
          ...(nowaNazwa === undefined ? {} : { newName: nowaNazwa }),
          ...(barwa === undefined ? {} : { color: barwa }),
        }),
        Command.LibraryTagUpdate,
        (tresc) => czyObiekt(tresc.tag) && typeof tresc.affectedFiles === 'number',
      );
    },

    async polaczEtykiety(zrodlowe, docelowa) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryTagMerge, {
          sourceNames: [...zrodlowe],
          targetName: docelowa,
        }),
        Command.LibraryTagMerge,
        (tresc) => czyObiekt(tresc.tag) && czyTablica(tresc.mergedNames),
      );
    },

    async usunEtykiete(nazwa, potwierdzone) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryTagRemove, { name: nazwa, confirm: potwierdzone }),
        Command.LibraryTagRemove,
        (tresc) => typeof tresc.removed === 'boolean' && typeof tresc.affectedFiles === 'number',
      );
    },

    async wykazKolekcji() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryCollectionList, {}),
        Command.LibraryCollectionList,
        (tresc) => czyTablica(tresc.collections) && typeof tresc.total === 'number',
      );
    },

    async ustawRelacje(zrodlo, cel, relacja, zdejmij) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryThesaurusRelate, {
          sourceName: zrodlo,
          targetName: cel,
          relation: relacja,
          remove: zdejmij,
        }),
        Command.LibraryThesaurusRelate,
        (tresc) => typeof tresc.relationSet === 'boolean',
      );
    },

    async wywiezTezaurus(format) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryThesaurusExport, { format, includeCollections: true }),
        Command.LibraryThesaurusExport,
        (tresc) =>
          typeof tresc.content === 'string' &&
          typeof tresc.conceptCount === 'number' &&
          typeof tresc.relationCount === 'number',
      );
    },

    async wykazRegul() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryRuleList, {}),
        Command.LibraryRuleList,
        (tresc) => czyTablica(tresc.rules) && typeof tresc.total === 'number',
      );
    },

    async zapiszRegule(regula, przeliczTeraz) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryRuleSet, { rule: regula, runNow: przeliczTeraz }),
        Command.LibraryRuleSet,
        (tresc) => czyObiekt(tresc.rule),
      );
    },

    async usunRegule(idReguly, zdejmijPliki) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryRuleRemove, {
          ruleId: idReguly,
          detachFiles: zdejmijPliki,
        }),
        Command.LibraryRuleRemove,
        (tresc) => typeof tresc.removed === 'boolean' && typeof tresc.detachedFiles === 'number',
      );
    },

    async skanujDuplikaty(rodzaje) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryDuplicateScan, { kinds: [...rodzaje] }),
        Command.LibraryDuplicateScan,
        (tresc) =>
          czyTablica(tresc.groups) &&
          typeof tresc.total === 'number' &&
          typeof tresc.skippedWithoutChecksum === 'number',
      );
    },

    async polaczDuplikaty(docelowy, zrodlowe) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryDuplicateMerge, {
          targetFileId: docelowy,
          sourceFileIds: [...zrodlowe],
          keepVersions: true,
        }),
        Command.LibraryDuplicateMerge,
        (tresc) => czyObiekt(tresc.file) && typeof tresc.mergedCount === 'number',
      );
    },

    async sprawdzIntegralnosc() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFixityCheck, {}),
        Command.LibraryFixityCheck,
        (tresc) =>
          czyTablica(tresc.results) &&
          typeof tresc.checkedCount === 'number' &&
          typeof tresc.mismatchedCount === 'number' &&
          typeof tresc.missingCount === 'number',
      );
    },

    async normalizujNazwy(idPlikow, probny) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryNameNormalize, {
          ...(idPlikow.length === 0 ? {} : { fileIds: [...idPlikow] }),
          dryRun: probny,
          transliterate: true,
        }),
        Command.LibraryNameNormalize,
        (tresc) => czyTablica(tresc.results) && typeof tresc.changedCount === 'number',
      );
    },

    async pulpitStanu() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryStatsGet, {}),
        Command.LibraryStatsGet,
        (tresc) => czyObiekt(tresc.stats),
      );
    },

    async dziennikAudytu(idPliku) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryAuditList, {
          ...(idPliku === undefined ? {} : { fileId: idPliku }),
          limit: 50,
        }),
        Command.LibraryAuditList,
        (tresc) => czyTablica(tresc.entries) && typeof tresc.total === 'number',
      );
    },

    async przenies(idPlikow, sciezka) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileMove, { fileIds: [...idPlikow], path: sciezka }),
        Command.LibraryFileMove,
        (tresc) => typeof tresc.movedCount === 'number' && czyTablica(tresc.files),
      );
    },

    async zarchiwizuj(idPlikow, powod) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileArchive, {
          fileIds: [...idPlikow],
          ...(powod === '' ? {} : { reason: powod }),
        }),
        Command.LibraryFileArchive,
        (tresc) => typeof tresc.archivedCount === 'number' && czyTablica(tresc.files),
      );
    },

    async przywrocZArchiwum(idPlikow) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileRestore, { fileIds: [...idPlikow] }),
        Command.LibraryFileRestore,
        (tresc) => typeof tresc.restoredCount === 'number' && czyTablica(tresc.files),
      );
    },

    async usunTrwale(idPlikow) {
      // Potwierdzenie jest zawsze true: okno pyta Operatora osobno, powtórne pytanie byłoby pozorne.
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileDelete, { fileIds: [...idPlikow], confirm: true }),
        Command.LibraryFileDelete,
        (tresc) =>
          typeof tresc.deletedCount === 'number' && typeof tresc.deletedVersions === 'number',
      );
    },

    async utrwal(idPlikow, postac) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryPreservationRun, {
          fileIds: [...idPlikow],
          kind: postac,
        }),
        Command.LibraryPreservationRun,
        (tresc) =>
          czyTablica(tresc.results) &&
          typeof tresc.validCount === 'number' &&
          typeof tresc.invalidCount === 'number',
      );
    },

    async wywiezPaczke(idPlikow, postac, format) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryPackageExport, {
          kind: postac,
          ...(idPlikow.length === 0 ? {} : { fileIds: [...idPlikow] }),
          format,
        }),
        Command.LibraryPackageExport,
        (tresc) =>
          typeof tresc.fileId === 'string' &&
          typeof tresc.entries === 'number' &&
          typeof tresc.sizeBytes === 'number' &&
          typeof tresc.manifestChecksum === 'string',
      );
    },

    async wykazRetencji(wDniach) {
      return sprawdzKsztalt(
        await straz.wywolaj(
          Command.LibraryRetentionList,
          wDniach === undefined ? {} : { dueWithinDays: wDniach },
        ),
        Command.LibraryRetentionList,
        (tresc) => czyTablica(tresc.policies),
      );
    },

    async zapiszRetencje(polityka, usun) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryRetentionSet, { policy: polityka, remove: usun }),
        Command.LibraryRetentionSet,
        (tresc) => czyObiekt(tresc.policy) && typeof tresc.affectedFiles === 'number',
      );
    },

    async wystawUdostepnienie(zakres, idCelu, sekundy) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryShareCreate, {
          scope: zakres,
          targetId: idCelu,
          ...(sekundy === undefined ? {} : { expiresInSeconds: sekundy }),
        }),
        Command.LibraryShareCreate,
        (tresc) => czyObiekt(tresc.share) && typeof tresc.url === 'string',
      );
    },

    async wykazUdostepnien() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryShareList, {}),
        Command.LibraryShareList,
        (tresc) => czyTablica(tresc.shares) && typeof tresc.total === 'number',
      );
    },

    async odwolajUdostepnienie(idUdostepnienia) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryShareRevoke, { shareId: idUdostepnienia }),
        Command.LibraryShareRevoke,
        (tresc) => typeof tresc.revoked === 'boolean',
      );
    },

    async wykazNasluchow() {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryWebhookList, {}),
        Command.LibraryWebhookList,
        (tresc) => czyTablica(tresc.webhooks) && typeof tresc.total === 'number',
      );
    },

    async zapiszNasluch(nasluch) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryWebhookSet, { webhook: nasluch }),
        Command.LibraryWebhookSet,
        (tresc) => czyObiekt(tresc.webhook),
      );
    },

    async usunNasluch(idNasluchu) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryWebhookRemove, { webhookId: idNasluchu }),
        Command.LibraryWebhookRemove,
        (tresc) => typeof tresc.removed === 'boolean',
      );
    },

    async klasyfikuj(idPlikow, zapisz) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryClassifyRun, {
          ...(idPlikow.length === 0 ? {} : { fileIds: [...idPlikow] }),
          apply: zapisz,
        }),
        Command.LibraryClassifyRun,
        (tresc) =>
          czyTablica(tresc.suggestions) &&
          typeof tresc.processedCount === 'number' &&
          typeof tresc.appliedCount === 'number',
      );
    },

    async wykazSugestii(idPliku) {
      return sprawdzKsztalt(
        await straz.wywolaj(
          Command.LibrarySuggestionList,
          idPliku === undefined ? {} : { fileId: idPliku },
        ),
        Command.LibrarySuggestionList,
        (tresc) => czyTablica(tresc.suggestions) && typeof tresc.total === 'number',
      );
    },

    async rozstrzygnijSugestie(idSugestii, przyjmij) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibrarySuggestionApply, {
          suggestionIds: [...idSugestii],
          accept: przyjmij,
        }),
        Command.LibrarySuggestionApply,
        (tresc) =>
          typeof tresc.appliedCount === 'number' && typeof tresc.rejectedCount === 'number',
      );
    },

    async porownaj(lewy, prawy, lewaWersja, prawaWersja) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryDiffCompare, {
          leftFileId: lewy,
          ...(prawy === undefined ? {} : { rightFileId: prawy }),
          ...(lewaWersja === undefined ? {} : { leftVersionId: lewaWersja }),
          ...(prawaWersja === undefined ? {} : { rightVersionId: prawaWersja }),
          contextLines: 0,
        }),
        Command.LibraryDiffCompare,
        (tresc) => czyObiekt(tresc.diff),
      );
    },
  };
}
