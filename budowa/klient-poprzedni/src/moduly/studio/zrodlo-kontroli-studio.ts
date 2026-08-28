import {
  Command,
  type ConfigScope,
  type StudioActionKind,
  type StudioActionState,
  type StudioAgentConflictPolicy,
  type StudioAgentsClaimResponse,
  type StudioAgentsReleaseResponse,
  type StudioAgentsSettingsSetResponse,
  type StudioAuthor,
  type StudioAutosaveGetResponse,
  type StudioAutosaveRunResponse,
  type StudioAutosaveSetResponse,
  type StudioBackupCreateResponse,
  type StudioBackupListResponse,
  type StudioBackupReason,
  type StudioBackupRestoreResponse,
  type StudioClipboardCopyResponse,
  type StudioClipboardPasteResponse,
  type StudioDiffFormCompareResponse,
  type StudioDiffHunkApplyResponse,
  type StudioJournalListResponse,
  type StudioJournalRedoResponse,
  type StudioJournalRevertResponse,
  type StudioLockAddResponse,
  type StudioLockListResponse,
  type StudioLockRemoveResponse,
  type StudioLockScope,
  type StudioMarkupAddResponse,
  type StudioMarkupDecideResponse,
  type StudioMarkupKind,
  type StudioMarkupListResponse,
  type StudioMarkupRemoveResponse,
  type StudioMarkupState,
  type StudioMarkupTypeDeleteResponse,
  type StudioMarkupTypeListResponse,
  type StudioMarkupTypeSaveResponse,
  type StudioModelChangesListResponse,
  type StudioModelChangesNavigateResponse,
  type StudioModelChangesRevertResponse,
  type StudioPasteMode,
  type StudioVersionRestoreInitialResponse,
  type StudioVersionSeries,
  type StudioVersionSeriesListResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import {
  czyLiczba,
  czyLogiczna,
  czyObiekt,
  czyTablica,
  sprawdzKsztalt,
} from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/** Kontrola pracy i bezpieczeństwo dokumentu: dziennik, zmiany modelu, autozapis, blokady i wykonawcy. */

/** Wskazanie wykonawcy czynności w żądaniu do rdzenia; pole puste znaczy, że czynność wykonał Operator. */
export interface WykonawcaCzynnosci {
  autor?: StudioAuthor;
  idWykonawcy?: string;
  nazwaWykonawcy?: string;
  idPodagenta?: string;
}

/** Zawężenie wykazu dziennika czynności dokumentu wedle autora, rodzaju, stanu i granicy liczby wpisów. */
export interface ZawezenieDziennika {
  autor?: StudioAuthor;
  rodzaj?: StudioActionKind;
  stan?: StudioActionState;
  granica?: number;
}

/** Nastawy autozapisu dokumentu przestawiane jednym wywołaniem: częstość, zdarzenia zapisu i liczba kopii. */
export interface NastawyAutozapisu {
  czynny: boolean;
  odstepSekund?: number;
  przyOdejsciu?: boolean;
  przyZamknieciu?: boolean;
  przyPrzelaczeniu?: boolean;
  ileKopiiZachowac?: number;
  poIluGodzinachWygasa?: number;
}

/** Nastawy pracy kilku wykonawców naraz nad tym samym dokumentem, wraz z polityką rozstrzygania sporów. */
export interface NastawyWykonawcow {
  zasieg?: ConfigScope;
  idBytuZasiegu?: string;
  petlaCzynna?: boolean;
  wieluCzynnych?: boolean;
  ilu?: number;
  przySpieciu?: StudioAgentConflictPolicy;
  granicaObiegow?: number;
  progBezPostepu?: number;
  wymagajZajecia?: boolean;
}

export interface ZrodloKontroliStudio {
  /* ── Dziennik czynności ─────────────────────────────────────────────────── */

  /** Odwracalny dziennik czynności dokumentu wraz z zależnościami. */
  kontrolaDziennik(
    idDokumentu: string,
    zawezenie: ZawezenieDziennika,
  ): Promise<Wynik<StudioJournalListResponse>>;
  // Cofa wskazane czynności pojedynczo i nie po kolei; nakłada różnicę drzew, nie migawkę stanu.
  kontrolaCofnijCzynnosci(
    idDokumentu: string,
    kody: readonly string[],
    zalozWersje: boolean,
  ): Promise<Wynik<StudioJournalRevertResponse>>;
  /** Ponawia cofnięte czynności — obejmuje treść i postać. */
  kontrolaPonowCzynnosci(
    idDokumentu: string,
    kody: readonly string[],
  ): Promise<Wynik<StudioJournalRedoResponse>>;

  /* ── Zmiany modelu ─────────────────────────────────────────────────────── */

  /** Wszystkie zmiany modelu wraz z licznikiem do przełącznika podświetlenia. */
  kontrolaZmianyModelu(
    idDokumentu: string,
    zRozstrzygnietymi: boolean,
    tylkoPostac: boolean,
  ): Promise<Wynik<StudioModelChangesListResponse>>;
  /** Wskazuje następną albo poprzednią zmianę modelu wraz z miejscem przewinięcia. */
  kontrolaSkokPoZmianachModelu(
    idDokumentu: string,
    odMiejsca: number,
    wPrzod: boolean,
    tylkoPostac: boolean,
  ): Promise<Wynik<StudioModelChangesNavigateResponse>>;
  // Cofa zmiany modelu z zachowaniem zmian Operatora naniesionych w tym czasie; kopia jedzie zawsze.
  kontrolaCofnijZmianyModelu(
    idDokumentu: string,
    kody: readonly string[],
    wszystkie: boolean,
  ): Promise<Wynik<StudioModelChangesRevertResponse>>;

  /* ── Różnica postaci i przeniesienie fragmentu ──────────────────────────── */

  /** Porównuje POSTAĆ dwóch wersji — zmiana kroju nie milczy. */
  kontrolaRoznicaPostaci(
    idDokumentu: string,
    odniesienie: string,
    porownywana: string,
    obszar: string,
  ): Promise<Wynik<StudioDiffFormCompareResponse>>;
  /** Przenosi pojedynczy fragment ze wskazanej wersji do stanu bieżącego. */
  kontrolaPrzeniesFragment(
    idDokumentu: string,
    idWersjiZrodlowej: string,
    wskazanie: { numerFragmentu?: number; poczatek?: number; koniec?: number },
    zPostacia: boolean,
    wykonawca: WykonawcaCzynnosci,
  ): Promise<Wynik<StudioDiffHunkApplyResponse>>;

  /* ── Autozapis, kopie, szeregi wersji ──────────────────────────────────── */

  /** Nastawy autozapisu wraz z czasem ostatniego zapisu i jego niepowodzeniem. */
  kontrolaAutozapisNastawy(
    idDokumentu: string,
    idOkna: string,
  ): Promise<Wynik<StudioAutosaveGetResponse>>;
  /** Przestawia autozapis — odstęp, zdarzenia okna, wygasanie kopii. */
  kontrolaAutozapisUstaw(
    idDokumentu: string,
    idOkna: string,
    nastawy: NastawyAutozapisu,
  ): Promise<Wynik<StudioAutosaveSetResponse>>;
  /** Wykonuje zapis samoczynny w OSOBNYM szeregu wersji. */
  kontrolaAutozapisWykonaj(
    idDokumentu: string,
    tresc: string,
    postac: unknown,
    powod: StudioBackupReason,
  ): Promise<Wynik<StudioAutosaveRunResponse>>;
  /** Zakłada kopię zapasową niezależnie od historii wersji. */
  kontrolaKopiaZaloz(
    idDokumentu: string,
    powod: StudioBackupReason,
    tresc: string,
    postac: unknown,
  ): Promise<Wynik<StudioBackupCreateResponse>>;
  /** Wykaz kopii wraz z liczbą kopii niosących zmiany niezapisane. */
  kontrolaKopieWykaz(
    idDokumentu: string,
    idOkna: string,
    tylkoNiezapisane: boolean,
  ): Promise<Wynik<StudioBackupListResponse>>;
  /** Przywraca kopię — na miejsce albo DO NOWEGO DOKUMENTU. */
  kontrolaKopiaPrzywroc(
    idKopii: string,
    doNowegoDokumentu: boolean,
    idOkna: string,
    tytul: string,
  ): Promise<Wynik<StudioBackupRestoreResponse>>;
  /** Wersje w rozbiciu na szereg Operatora i szereg autozapisu. */
  kontrolaSzeregiWersji(
    idDokumentu: string,
    szereg: StudioVersionSeries | null,
    tylkoKluczowe: boolean,
    granica: number,
  ): Promise<Wynik<StudioVersionSeriesListResponse>>;
  /** Wraca dokumentem do wersji założycielskiej jednym poleceniem. */
  kontrolaPowrotDoZalozycielskiej(
    idDokumentu: string,
    zalozWersje: boolean,
  ): Promise<Wynik<StudioVersionRestoreInitialResponse>>;

  /* ── Znakowanie fragmentu w rdzeniu ────────────────────────────────────── */

  /** Znakuje fragment — wyróżnienie, znacznik własny albo propozycja z brzmieniem. */
  kontrolaZnakowanieDodaj(
    idDokumentu: string,
    rodzaj: StudioMarkupKind,
    zakres: { poczatek: number; koniec: number },
    tresc: { barwa?: string; rodzajZnacznika?: string; uzasadnienie?: string; brzmienie?: string },
    wykonawca: WykonawcaCzynnosci,
  ): Promise<Wynik<StudioMarkupAddResponse>>;
  /** Wykaz znakowań dokumentu jako spis do przejścia. */
  kontrolaZnakowaniaWykaz(
    idDokumentu: string,
    zawezenie: {
      rodzaj?: StudioMarkupKind;
      autor?: StudioAuthor;
      stan?: StudioMarkupState;
      rodzajZnacznika?: string;
      zKomentarzami?: boolean;
      zeZmianami?: boolean;
    },
  ): Promise<Wynik<StudioMarkupListResponse>>;
  /** Zdejmuje znakowanie fragmentu. */
  kontrolaZnakowanieZdejmij(
    idDokumentu: string,
    idZnakowania: string,
    wykonawca: WykonawcaCzynnosci,
  ): Promise<Wynik<StudioMarkupRemoveResponse>>;
  /** Rozstrzyga propozycje z marginesu — przyjęcie wnosi brzmienie do treści. */
  kontrolaZnakowanieRozstrzygnij(
    idDokumentu: string,
    kody: readonly string[],
    przyjmij: boolean,
    poprawioneBrzmienie: string,
  ): Promise<Wynik<StudioMarkupDecideResponse>>;
  /** Rodzaje znaczników własnych wraz z barwą i liczbą użyć. */
  kontrolaRodzajeZnacznikow(idDokumentu: string): Promise<Wynik<StudioMarkupTypeListResponse>>;
  /** Zakłada albo zmienia rodzaj znacznika własnego. */
  kontrolaRodzajZnacznikaZapisz(
    nazwa: string,
    nazwaWidoczna: string,
    barwa: string,
  ): Promise<Wynik<StudioMarkupTypeSaveResponse>>;
  /** Usuwa rodzaj znacznika własnego; fabrycznego rdzeń nie usuwa. */
  kontrolaRodzajZnacznikaUsun(nazwa: string): Promise<Wynik<StudioMarkupTypeDeleteResponse>>;

  /* ── Blokady fragmentów ────────────────────────────────────────────────── */

  /** Zakłada blokadę na fragmencie wraz z nazwą, powodem i zasięgiem. */
  kontrolaBlokadaZaloz(
    idDokumentu: string,
    zakres: { poczatek: number; koniec: number },
    nazwa: string,
    powod: string,
    zasieg: StudioLockScope,
  ): Promise<Wynik<StudioLockAddResponse>>;
  /** Wykaz blokad dokumentu. */
  kontrolaBlokadyWykaz(idDokumentu: string): Promise<Wynik<StudioLockListResponse>>;
  /** Zdejmuje blokadę; zdejmuje ją WYŁĄCZNIE Operator. */
  kontrolaBlokadaZdejmij(
    idDokumentu: string,
    idBlokady: string,
  ): Promise<Wynik<StudioLockRemoveResponse>>;

  /* ── Zajęcia fragmentów i nastawy wykonawców ───────────────────────────── */

  /** Zajmuje fragment dla wykonawcy; zajęty przez kogo innego wraca odmową. */
  kontrolaZajmijFragment(
    idDokumentu: string,
    zakres: { poczatek: number; koniec: number },
    wykonawca: WykonawcaCzynnosci,
    ttlSekund: number,
  ): Promise<Wynik<StudioAgentsClaimResponse>>;
  /** Zwalnia zajęcie fragmentu — wskazane albo wszystkie wykonawcy. */
  kontrolaZwolnijFragment(
    idDokumentu: string,
    idZajecia: string,
    wykonawca: WykonawcaCzynnosci,
    wszystkie: boolean,
  ): Promise<Wynik<StudioAgentsReleaseResponse>>;
  /** Przestawia nastawy pętli wykonawczej i pracy kilku wykonawców naraz. */
  kontrolaNastawyWykonawcow(
    nastawy: NastawyWykonawcow,
  ): Promise<Wynik<StudioAgentsSettingsSetResponse>>;

  /* ── Schowek dokumentu ─────────────────────────────────────────────────── */

  /** Odkłada fragment dokumentu do schowka wraz z postacią; wycięcie opcjonalne. */
  kontrolaSchowekOdlozFragment(
    idDokumentu: string,
    zakres: { poczatek: number; koniec: number },
    wytnij: boolean,
    samaPostac: boolean,
    wykonawca: WykonawcaCzynnosci,
  ): Promise<Wynik<StudioClipboardCopyResponse>>;
  /** Wkleja wpis schowka — z postacią albo jako czysty tekst, wedle wyboru. */
  kontrolaSchowekWklejFragment(
    idDokumentu: string,
    miejsce: number,
    wskazanie: { idWpisu?: string; tresc?: string },
    tryb: StudioPasteMode,
    zamieniany: { poczatek?: number; koniec?: number },
    wykonawca: WykonawcaCzynnosci,
  ): Promise<Wynik<StudioClipboardPasteResponse>>;
}

/** Pola wykonawcy przekształcone do kształtu żądania wysyłanego do rdzenia; nieznane pola nie jadą wcale. */
function poleWykonawcy(wykonawca: WykonawcaCzynnosci): Record<string, unknown> {
  return {
    ...(wykonawca.autor === undefined ? {} : { author: wykonawca.autor }),
    ...(wykonawca.idWykonawcy === undefined || wykonawca.idWykonawcy === ''
      ? {}
      : { agentId: wykonawca.idWykonawcy }),
    ...(wykonawca.nazwaWykonawcy === undefined || wykonawca.nazwaWykonawcy === ''
      ? {}
      : { agentName: wykonawca.nazwaWykonawcy }),
    ...(wykonawca.idPodagenta === undefined || wykonawca.idPodagenta === ''
      ? {}
      : { subagentId: wykonawca.idPodagenta }),
  };
}

/** Pole liczbowe dołączane do żądania rdzenia, gdy podane i skończone; wartość ujemna nie jedzie wcale. */
function poleLiczby(nazwa: string, wartosc: number | undefined): Record<string, number> {
  if (wartosc === undefined || !Number.isFinite(wartosc) || wartosc < 0) return {};
  return { [nazwa]: wartosc };
}

export function utworzZrodloKontroliStudio(kanal: Kanal): ZrodloKontroliStudio {
  return {
    async kontrolaDziennik(idDokumentu, zawezenie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioJournalList, {
          documentId: idDokumentu,
          ...(zawezenie.autor === undefined ? {} : { author: zawezenie.autor }),
          ...(zawezenie.rodzaj === undefined ? {} : { kind: zawezenie.rodzaj }),
          ...(zawezenie.stan === undefined ? {} : { state: zawezenie.stan }),
          ...poleLiczby('limit', zawezenie.granica),
        }),
        Command.StudioJournalList,
        (tresc) => czyTablica(tresc.actions) && czyLiczba(tresc.total),
      );
    },

    async kontrolaCofnijCzynnosci(idDokumentu, kody, zalozWersje) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioJournalRevert, {
          documentId: idDokumentu,
          actionIds: [...kody],
          createVersion: zalozWersje,
        }),
        Command.StudioJournalRevert,
        // Bilans jest polem obowiązkowym: to on niesie zależność, przez którą cofnięcie nie zaszło.
        (tresc) => czyTablica(tresc.reverted) && czyObiekt(tresc.balance),
      );
    },

    async kontrolaPonowCzynnosci(idDokumentu, kody) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioJournalRedo, {
          documentId: idDokumentu,
          actionIds: [...kody],
        }),
        Command.StudioJournalRedo,
        (tresc) => czyTablica(tresc.redone) && czyObiekt(tresc.balance),
      );
    },

    async kontrolaZmianyModelu(idDokumentu, zRozstrzygnietymi, tylkoPostac) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioModelChangesList, {
          documentId: idDokumentu,
          includeDecided: zRozstrzygnietymi,
          onlyFormChanges: tylkoPostac,
        }),
        Command.StudioModelChangesList,
        (tresc) => czyObiekt(tresc.summary),
      );
    },

    async kontrolaSkokPoZmianachModelu(idDokumentu, odMiejsca, wPrzod, tylkoPostac) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioModelChangesNavigate, {
          documentId: idDokumentu,
          ...poleLiczby('fromOffset', odMiejsca),
          // Kierunek jedzie napisem wprost z kontraktu: pole direction nie ma osobnego wykazu wartości.
          direction: wPrzod ? 'next' : 'previous',
          onlyFormChanges: tylkoPostac,
        }),
        Command.StudioModelChangesNavigate,
        // Brak zmiany w odpowiedzi znaczy koniec wykazu, więc sprawdzana jest wyłącznie liczba wpisów.
        (tresc) => czyLiczba(tresc.total),
      );
    },

    async kontrolaCofnijZmianyModelu(idDokumentu, kody, wszystkie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioModelChangesRevert, {
          documentId: idDokumentu,
          ...(kody.length === 0 ? {} : { changeIds: [...kody] }),
          all: wszystkie,
          // Kopia zapasowa przed czynnością nieodwracalną jedzie zawsze, niezależnie od wyboru okna.
          createBackup: true,
          createVersion: true,
        }),
        Command.StudioModelChangesRevert,
        (tresc) =>
          czyLiczba(tresc.revertedCount) &&
          czyLiczba(tresc.keptOperatorChanges) &&
          czyObiekt(tresc.balance),
      );
    },

    async kontrolaRoznicaPostaci(idDokumentu, odniesienie, porownywana, obszar) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDiffFormCompare, {
          documentId: idDokumentu,
          // Brak odniesienia czyta rdzeń jako wersję założycielską; pola puste nie jadą, by nie zgadywać wersji.
          ...(odniesienie === '' ? {} : { baseVersionId: odniesienie }),
          ...(porownywana === '' ? {} : { targetVersionId: porownywana }),
          ...(obszar === '' ? {} : { area: obszar }),
        }),
        Command.StudioDiffFormCompare,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async kontrolaPrzeniesFragment(
      idDokumentu,
      idWersjiZrodlowej,
      wskazanie,
      zPostacia,
      wykonawca,
    ) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDiffHunkApply, {
          documentId: idDokumentu,
          sourceVersionId: idWersjiZrodlowej,
          ...poleLiczby('hunkIndex', wskazanie.numerFragmentu),
          ...poleLiczby('rangeStart', wskazanie.poczatek),
          ...poleLiczby('rangeEnd', wskazanie.koniec),
          includeForm: zPostacia,
          ...poleWykonawcy(wykonawca),
        }),
        Command.StudioDiffHunkApply,
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.balance),
      );
    },

    async kontrolaAutozapisNastawy(idDokumentu, idOkna) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAutosaveGet, {
          ...(idDokumentu === '' ? {} : { documentId: idDokumentu }),
          ...(idOkna === '' ? {} : { windowId: idOkna }),
        }),
        Command.StudioAutosaveGet,
        (tresc) => czyObiekt(tresc.settings),
      );
    },

    async kontrolaAutozapisUstaw(idDokumentu, idOkna, nastawy) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAutosaveSet, {
          ...(idDokumentu === '' ? {} : { documentId: idDokumentu }),
          ...(idOkna === '' ? {} : { windowId: idOkna }),
          enabled: nastawy.czynny,
          ...poleLiczby('intervalSeconds', nastawy.odstepSekund),
          ...(nastawy.przyOdejsciu === undefined ? {} : { onBlur: nastawy.przyOdejsciu }),
          ...(nastawy.przyZamknieciu === undefined ? {} : { onClose: nastawy.przyZamknieciu }),
          ...(nastawy.przyPrzelaczeniu === undefined
            ? {}
            : { onSwitch: nastawy.przyPrzelaczeniu }),
          ...poleLiczby('backupRetentionCount', nastawy.ileKopiiZachowac),
          ...poleLiczby('backupRetentionHours', nastawy.poIluGodzinachWygasa),
        }),
        Command.StudioAutosaveSet,
        (tresc) => czyObiekt(tresc.settings),
      );
    },

    async kontrolaAutozapisWykonaj(idDokumentu, tresc, postac, powod) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAutosaveRun, {
          documentId: idDokumentu,
          content: tresc,
          ...(postac === undefined || postac === null ? {} : { form: postac }),
          trigger: powod,
        }),
        Command.StudioAutosaveRun,
        // Odpowiedź udana z polem saved false opisuje zapis nieudany; sprawdzany jest sam fakt obecności pola.
        (odpowiedz) => czyLogiczna(odpowiedz.saved),
      );
    },

    async kontrolaKopiaZaloz(idDokumentu, powod, tresc, postac) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioBackupCreate, {
          documentId: idDokumentu,
          reason: powod,
          ...(tresc === '' ? {} : { content: tresc }),
          ...(postac === undefined || postac === null ? {} : { form: postac }),
        }),
        Command.StudioBackupCreate,
        (odpowiedz) => czyObiekt(odpowiedz.backup),
      );
    },

    async kontrolaKopieWykaz(idDokumentu, idOkna, tylkoNiezapisane) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioBackupList, {
          ...(idDokumentu === '' ? {} : { documentId: idDokumentu }),
          ...(idOkna === '' ? {} : { windowId: idOkna }),
          unsavedOnly: tylkoNiezapisane,
        }),
        Command.StudioBackupList,
        (tresc) => czyTablica(tresc.backups) && czyLiczba(tresc.unsavedCount),
      );
    },

    async kontrolaKopiaPrzywroc(idKopii, doNowegoDokumentu, idOkna, tytul) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioBackupRestore, {
          backupId: idKopii,
          asNewDocument: doNowegoDokumentu,
          ...(idOkna === '' ? {} : { windowId: idOkna }),
          ...(tytul === '' ? {} : { title: tytul }),
        }),
        Command.StudioBackupRestore,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async kontrolaSzeregiWersji(idDokumentu, szereg, tylkoKluczowe, granica) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioVersionSeriesList, {
          documentId: idDokumentu,
          ...(szereg === null ? {} : { series: szereg }),
          milestonesOnly: tylkoKluczowe,
          ...poleLiczby('limit', granica),
        }),
        Command.StudioVersionSeriesList,
        (tresc) =>
          czyTablica(tresc.versions) &&
          czyLiczba(tresc.operatorCount) &&
          czyLiczba(tresc.autosaveCount),
      );
    },

    async kontrolaPowrotDoZalozycielskiej(idDokumentu, zalozWersje) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioVersionRestoreInitial, {
          documentId: idDokumentu,
          createVersion: zalozWersje,
        }),
        Command.StudioVersionRestoreInitial,
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.version),
      );
    },

    async kontrolaZnakowanieDodaj(idDokumentu, rodzaj, zakres, tresc, wykonawca) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupAdd, {
          documentId: idDokumentu,
          kind: rodzaj,
          rangeStart: zakres.poczatek,
          rangeEnd: zakres.koniec,
          ...(tresc.barwa === undefined || tresc.barwa === '' ? {} : { color: tresc.barwa }),
          ...(tresc.rodzajZnacznika === undefined || tresc.rodzajZnacznika === ''
            ? {}
            : { markType: tresc.rodzajZnacznika }),
          ...(tresc.uzasadnienie === undefined || tresc.uzasadnienie === ''
            ? {}
            : { body: tresc.uzasadnienie }),
          ...(tresc.brzmienie === undefined || tresc.brzmienie === ''
            ? {}
            : { suggestedText: tresc.brzmienie }),
          ...poleWykonawcy(wykonawca),
        }),
        Command.StudioMarkupAdd,
        (odpowiedz) => czyObiekt(odpowiedz.markup),
      );
    },

    async kontrolaZnakowaniaWykaz(idDokumentu, zawezenie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupList, {
          documentId: idDokumentu,
          ...(zawezenie.rodzaj === undefined ? {} : { kind: zawezenie.rodzaj }),
          ...(zawezenie.autor === undefined ? {} : { author: zawezenie.autor }),
          ...(zawezenie.stan === undefined ? {} : { state: zawezenie.stan }),
          ...(zawezenie.rodzajZnacznika === undefined || zawezenie.rodzajZnacznika === ''
            ? {}
            : { markType: zawezenie.rodzajZnacznika }),
          ...(zawezenie.zKomentarzami === undefined
            ? {}
            : { includeComments: zawezenie.zKomentarzami }),
          ...(zawezenie.zeZmianami === undefined
            ? {}
            : { includeTrackedChanges: zawezenie.zeZmianami }),
        }),
        Command.StudioMarkupList,
        (tresc) => czyTablica(tresc.markups) && czyLiczba(tresc.total),
      );
    },

    async kontrolaZnakowanieZdejmij(idDokumentu, idZnakowania, wykonawca) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupRemove, {
          documentId: idDokumentu,
          markupId: idZnakowania,
          ...poleWykonawcy(wykonawca),
        }),
        Command.StudioMarkupRemove,
        (tresc) => czyLogiczna(tresc.removed),
      );
    },

    async kontrolaZnakowanieRozstrzygnij(idDokumentu, kody, przyjmij, poprawioneBrzmienie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupDecide, {
          documentId: idDokumentu,
          markupIds: [...kody],
          accept: przyjmij,
          ...(poprawioneBrzmienie === '' ? {} : { editedText: poprawioneBrzmienie }),
          createVersion: true,
        }),
        Command.StudioMarkupDecide,
        (tresc) => czyLiczba(tresc.decided) && czyObiekt(tresc.document),
      );
    },

    async kontrolaRodzajeZnacznikow(idDokumentu) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupTypeList, {
          ...(idDokumentu === '' ? {} : { documentId: idDokumentu }),
        }),
        Command.StudioMarkupTypeList,
        (tresc) => czyTablica(tresc.markupTypes),
      );
    },

    async kontrolaRodzajZnacznikaZapisz(nazwa, nazwaWidoczna, barwa) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupTypeSave, {
          name: nazwa,
          ...(nazwaWidoczna === '' ? {} : { label: nazwaWidoczna }),
          ...(barwa === '' ? {} : { color: barwa }),
        }),
        Command.StudioMarkupTypeSave,
        (tresc) => czyObiekt(tresc.markupType),
      );
    },

    async kontrolaRodzajZnacznikaUsun(nazwa) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioMarkupTypeDelete, { name: nazwa }),
        Command.StudioMarkupTypeDelete,
        (tresc) => czyLogiczna(tresc.deleted),
      );
    },

    async kontrolaBlokadaZaloz(idDokumentu, zakres, nazwa, powod, zasieg) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioLockAdd, {
          documentId: idDokumentu,
          rangeStart: zakres.poczatek,
          rangeEnd: zakres.koniec,
          name: nazwa,
          ...(powod === '' ? {} : { reason: powod }),
          scope: zasieg,
        }),
        Command.StudioLockAdd,
        (tresc) => czyObiekt(tresc.lock) && czyTablica(tresc.locks),
      );
    },

    async kontrolaBlokadyWykaz(idDokumentu) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioLockList, { documentId: idDokumentu }),
        Command.StudioLockList,
        (tresc) => czyTablica(tresc.locks),
      );
    },

    async kontrolaBlokadaZdejmij(idDokumentu, idBlokady) {
      return sprawdzKsztalt(
        // Pola wykonawcy tu nie jadą: blokadę zdejmuje wyłącznie Operator, inny autor obszedłby blokadę.
        await wywolajUczciwie(kanal, Command.StudioLockRemove, {
          documentId: idDokumentu,
          lockId: idBlokady,
        }),
        Command.StudioLockRemove,
        (tresc) => czyLogiczna(tresc.removed) && czyTablica(tresc.locks),
      );
    },

    async kontrolaZajmijFragment(idDokumentu, zakres, wykonawca, ttlSekund) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAgentsClaim, {
          documentId: idDokumentu,
          rangeStart: zakres.poczatek,
          rangeEnd: zakres.koniec,
          ...poleWykonawcy(wykonawca),
          ...poleLiczby('ttlSeconds', ttlSekund),
        }),
        Command.StudioAgentsClaim,
        // Odmowa zajęcia jest odpowiedzią udaną: pole claimed false niesie wykonawcę trzymającego fragment.
        (tresc) => czyLogiczna(tresc.claimed),
      );
    },

    async kontrolaZwolnijFragment(idDokumentu, idZajecia, wykonawca, wszystkie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAgentsRelease, {
          documentId: idDokumentu,
          ...(idZajecia === '' ? {} : { slotId: idZajecia }),
          ...(wykonawca.idWykonawcy === undefined || wykonawca.idWykonawcy === ''
            ? {}
            : { agentId: wykonawca.idWykonawcy }),
          ...(wykonawca.idPodagenta === undefined || wykonawca.idPodagenta === ''
            ? {}
            : { subagentId: wykonawca.idPodagenta }),
          all: wszystkie,
        }),
        Command.StudioAgentsRelease,
        (tresc) => czyLiczba(tresc.released) && czyTablica(tresc.slots),
      );
    },

    async kontrolaNastawyWykonawcow(nastawy) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAgentsSettingsSet, {
          ...(nastawy.zasieg === undefined ? {} : { scope: nastawy.zasieg }),
          ...(nastawy.idBytuZasiegu === undefined || nastawy.idBytuZasiegu === ''
            ? {}
            : { scopeId: nastawy.idBytuZasiegu }),
          ...(nastawy.petlaCzynna === undefined
            ? {}
            : { executionLoopEnabled: nastawy.petlaCzynna }),
          ...(nastawy.wieluCzynnych === undefined
            ? {}
            : { multiAgentEnabled: nastawy.wieluCzynnych }),
          ...poleLiczby('maxConcurrentAgents', nastawy.ilu),
          ...(nastawy.przySpieciu === undefined ? {} : { conflictPolicy: nastawy.przySpieciu }),
          ...poleLiczby('loopMaxIterations', nastawy.granicaObiegow),
          ...poleLiczby('loopNoProgressThreshold', nastawy.progBezPostepu),
          ...(nastawy.wymagajZajecia === undefined
            ? {}
            : { requireFragmentClaim: nastawy.wymagajZajecia }),
        }),
        Command.StudioAgentsSettingsSet,
        (tresc) => czyObiekt(tresc.settings),
      );
    },

    async kontrolaSchowekOdlozFragment(idDokumentu, zakres, wytnij, samaPostac, wykonawca) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioClipboardCopy, {
          documentId: idDokumentu,
          rangeStart: zakres.poczatek,
          rangeEnd: zakres.koniec,
          cut: wytnij,
          formatOnly: samaPostac,
          ...poleWykonawcy(wykonawca),
        }),
        Command.StudioClipboardCopy,
        (tresc) => typeof tresc.clipboardEntryId === 'string',
      );
    },

    async kontrolaSchowekWklejFragment(
      idDokumentu,
      miejsce,
      wskazanie,
      tryb,
      zamieniany,
      wykonawca,
    ) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioClipboardPaste, {
          documentId: idDokumentu,
          offset: miejsce,
          ...(wskazanie.idWpisu === undefined || wskazanie.idWpisu === ''
            ? {}
            : { clipboardEntryId: wskazanie.idWpisu }),
          ...(wskazanie.tresc === undefined || wskazanie.tresc === ''
            ? {}
            : { text: wskazanie.tresc }),
          mode: tryb,
          ...poleLiczby('replaceRangeStart', zamieniany.poczatek),
          ...poleLiczby('replaceRangeEnd', zamieniany.koniec),
          ...poleWykonawcy(wykonawca),
        }),
        Command.StudioClipboardPaste,
        (tresc) => czyObiekt(tresc.balance),
      );
    },
  };
}
