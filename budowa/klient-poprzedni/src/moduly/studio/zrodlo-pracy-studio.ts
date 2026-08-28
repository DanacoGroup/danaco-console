import {
  AssetContentDisposition,
  Command,
  type StudioComment,
  type StudioDiffSourceResponse,
  type StudioDiffVisualResponse,
  type StudioDocument,
  type StudioExportProfile,
  type StudioPageSetup,
  type StudioProposalDecideResponse,
  type StudioSearchSemanticResponse,
  type StudioTemplate,
  type StudioTrackedChange,
  type StudioTrackingDecideResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Drogi do rdzenia potrzebne oknu pracy z dokumentem, a nieznane oknom poprzednim; źródło nie
 * ma stanu i nie buduje elementu, sprawdza kształt odpowiedzi.
 */
export interface ZrodloPracyStudio {
  /** Zmiany śledzone dokumentu; `tylkoOczekujace` zawęża do nierozstrzygniętych. */
  zmiany(
    idDokumentu: string,
    tylkoOczekujace: boolean,
  ): Promise<Wynik<{ changes: StudioTrackedChange[] }>>;
  /** Przyjmuje albo odrzuca wskazane zmiany śledzone. */
  rozstrzygnijZmiany(
    idDokumentu: string,
    kody: readonly string[],
    przyjmij: boolean,
  ): Promise<Wynik<StudioTrackingDecideResponse>>;
  /** Włącza albo wyłącza śledzenie zmian dokumentu. */
  ustawSledzenie(idDokumentu: string, czynne: boolean): Promise<Wynik<{ enabled: boolean }>>;
  /** Rozstrzyga propozycję w rdzeniu, w całości albo wskazanymi fragmentami. */
  rozstrzygnijPropozycje(
    idDokumentu: string,
    idPropozycji: string,
    przyjmij: boolean,
    fragmenty: readonly number[],
  ): Promise<Wynik<StudioProposalDecideResponse>>;
  /** Profile wydania — nastawy strony zapisane w rdzeniu. */
  profile(): Promise<Wynik<{ profiles: StudioExportProfile[] }>>;
  /** Zapisuje nastawy strony jako nazwany profil wydania. */
  zapiszProfil(
    nazwa: string,
    format: string,
    nastawy: StudioPageSetup,
  ): Promise<Wynik<{ profile: StudioExportProfile }>>;
  /** Zleca rdzeniowi render stron podglądu wraz z paginacją. */
  render(
    idDokumentu: string,
    format: string,
    idOkna: string,
    idProfilu: string,
  ): Promise<Wynik<{ pageAssetIds: string[]; pages: number }>>;
  /** Pobiera treść zasobu magazynu rdzenia; rodzina komendy `design.*` dotyczy każdego zasobu. */
  trescZasobu(
    idZasobu: string,
    granicaBajtow: number,
  ): Promise<Wynik<{ assetId: string; contentBase64?: string; mediaType: string; sizeBytes: number }>>;
  /** Szablony dokumentu — zawartość galerii. */
  szablony(): Promise<Wynik<{ templates: StudioTemplate[] }>>;
  /** Zakłada dokument z szablonu wraz z wypełnieniem jego pól. */
  zastosujSzablon(
    idOkna: string,
    idSzablonu: string,
    wartosci: Record<string, string>,
    tytul: string,
  ): Promise<Wynik<{ document: StudioDocument }>>;
  /** Komentarze redakcyjne dokumentu wraz z wątkami. */
  komentarze(
    idDokumentu: string,
    zRozwiazanymi: boolean,
  ): Promise<Wynik<{ comments: StudioComment[] }>>;
  /** Zakłada komentarz przypięty do fragmentu albo odpowiedź w wątku. */
  dodajKomentarz(
    idDokumentu: string,
    tresc: string,
    zakres: { poczatek: number; koniec: number } | null,
    idWatku: string,
  ): Promise<Wynik<{ comment: StudioComment }>>;
  /** Oznacza wątek jako rozwiązany albo zdejmuje to oznaczenie. */
  rozstrzygnijKomentarz(
    idKomentarza: string,
    rozwiazany: boolean,
  ): Promise<Wynik<{ comment: StudioComment }>>;
  /** Różnica wyglądu dwóch wersji — obszary i nakładki stron. */
  roznicaWygladu(
    idDokumentu: string,
    odniesienie: string,
    porownywana: string,
    idOkna: string,
  ): Promise<Wynik<StudioDiffVisualResponse>>;
  /** Wyszukiwanie znaczeniowe — podstawa pomiaru podobieństw w panelu Redaktora. */
  wyszukajZnaczeniowo(
    idDokumentu: string,
    zapytanie: string,
  ): Promise<Wynik<StudioSearchSemanticResponse>>;
  /** Zestawienie dokumentu z materiałem wejściowym. */
  zestawZeZrodlem(
    idDokumentu: string,
    wskazanie: { plikLibrary?: string; zasob?: string },
  ): Promise<Wynik<StudioDiffSourceResponse>>;
}

export function utworzZrodloPracyStudio(kanal: Kanal): ZrodloPracyStudio {
  return {
    async zmiany(idDokumentu, tylkoOczekujace) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioTrackingList, {
          documentId: idDokumentu,
          pendingOnly: tylkoOczekujace,
        }),
        Command.StudioTrackingList,
        (tresc) => czyTablica(tresc.changes),
      );
    },

    async rozstrzygnijZmiany(idDokumentu, kody, przyjmij) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioTrackingDecide, {
          documentId: idDokumentu,
          changeIds: [...kody],
          accept: przyjmij,
          createVersion: true,
        }),
        Command.StudioTrackingDecide,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async ustawSledzenie(idDokumentu, czynne) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioTrackingSet, {
          documentId: idDokumentu,
          enabled: czynne,
        }),
        Command.StudioTrackingSet,
        (tresc) => typeof tresc.enabled === 'boolean',
      );
    },

    async rozstrzygnijPropozycje(idDokumentu, idPropozycji, przyjmij, fragmenty) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioProposalDecide, {
          documentId: idDokumentu,
          proposalId: idPropozycji,
          accept: przyjmij,
          // Wykaz pusty znaczy w całości; brak hunkIndexes to całość, nie zero fragmentów.
          ...(fragmenty.length === 0 ? {} : { hunkIndexes: [...fragmenty] }),
          createVersion: true,
        }),
        Command.StudioProposalDecide,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async profile() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioExportProfileList, {}),
        Command.StudioExportProfileList,
        (tresc) => czyTablica(tresc.profiles),
      );
    },

    async zapiszProfil(nazwa, format, nastawy) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioExportProfileSave, {
          name: nazwa,
          format,
          pageSetup: nastawy,
        }),
        Command.StudioExportProfileSave,
        (tresc) => czyObiekt(tresc.profile),
      );
    },

    async render(idDokumentu, format, idOkna, idProfilu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioPreviewRender, {
          documentId: idDokumentu,
          format,
          ...(idOkna === '' ? {} : { windowId: idOkna }),
          ...(idProfilu === '' ? {} : { profileId: idProfilu }),
        }),
        Command.StudioPreviewRender,
        (tresc) => czyTablica(tresc.pageAssetIds) && typeof tresc.pages === 'number',
      );
    },

    async trescZasobu(idZasobu, granicaBajtow) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetContentGet, {
          assetId: idZasobu,
          disposition: AssetContentDisposition.Inline,
          ...(granicaBajtow > 0 ? { maxBytes: granicaBajtow } : {}),
        }),
        Command.DesignAssetContentGet,
        // Pole contentBase64 jest nieobowiązkowe; brak treści rozstrzyga wołający.
        (tresc) => typeof tresc.mediaType === 'string' && typeof tresc.sizeBytes === 'number',
      );
    },

    async szablony() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioTemplateList, {}),
        Command.StudioTemplateList,
        (tresc) => czyTablica(tresc.templates),
      );
    },

    async zastosujSzablon(idOkna, idSzablonu, wartosci, tytul) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioTemplateApply, {
          windowId: idOkna,
          templateId: idSzablonu,
          values: wartosci,
          ...(tytul === '' ? {} : { title: tytul }),
        }),
        Command.StudioTemplateApply,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async komentarze(idDokumentu, zRozwiazanymi) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioCommentList, {
          documentId: idDokumentu,
          includeResolved: zRozwiazanymi,
        }),
        Command.StudioCommentList,
        (tresc) => czyTablica(tresc.comments),
      );
    },

    async dodajKomentarz(idDokumentu, tresc, zakres, idWatku) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioCommentAdd, {
          documentId: idDokumentu,
          body: tresc,
          // Zakres jedzie tylko, gdy jest: brak zaznaczenia dotyczy dokumentu, nie zera.
          ...(zakres === null
            ? {}
            : { selectionStart: zakres.poczatek, selectionEnd: zakres.koniec }),
          ...(idWatku === '' ? {} : { parentCommentId: idWatku }),
        }),
        Command.StudioCommentAdd,
        (odpowiedz) => czyObiekt(odpowiedz.comment),
      );
    },

    async rozstrzygnijKomentarz(idKomentarza, rozwiazany) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioCommentResolve, {
          commentId: idKomentarza,
          resolved: rozwiazany,
        }),
        Command.StudioCommentResolve,
        (tresc) => czyObiekt(tresc.comment),
      );
    },

    async roznicaWygladu(idDokumentu, odniesienie, porownywana, idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioDiffVisual, {
          documentId: idDokumentu,
          baseVersionId: odniesienie,
          targetVersionId: porownywana,
          ...(idOkna === '' ? {} : { windowId: idOkna }),
        }),
        Command.StudioDiffVisual,
        (tresc) => czyTablica(tresc.regions),
      );
    },

    async wyszukajZnaczeniowo(idDokumentu, zapytanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioSearchSemantic, {
          documentId: idDokumentu,
          query: zapytanie,
        }),
        Command.StudioSearchSemantic,
        (tresc) => czyTablica(tresc.matches),
      );
    },

    async zestawZeZrodlem(idDokumentu, wskazanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioDiffSource, {
          documentId: idDokumentu,
          ...(wskazanie.plikLibrary === undefined
            ? {}
            : { sourceLibraryFileId: wskazanie.plikLibrary }),
          ...(wskazanie.zasob === undefined ? {} : { sourceAssetId: wskazanie.zasob }),
        }),
        Command.StudioDiffSource,
        (tresc) => czyTablica(tresc.hunks),
      );
    },
  };
}
