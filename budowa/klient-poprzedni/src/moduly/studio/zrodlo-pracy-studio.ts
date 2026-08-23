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
 * Drogi do rdzenia potrzebne oknu pracy z dokumentem, a nieznane oknom
 * poprzednim.
 *
 * `zrodlo-studio.ts` niesie sześć komend obszaru, na których stały Studio Editor
 * i Diff/Grep Panel: otwarcie, zapis, operacja, porównanie, wykaz wersji,
 * przywrócenie. Jedno okno pracy potrzebuje ponad to sześciu dalszych rodzin,
 * które w kontrakcie są, a w kliencie nie miały wołającego:
 *
 *   — `studio.tracking.*`  — zmiany śledzone: wykaz, decyzja, przełącznik.
 *     To one niosą wynik modelu w miejscu, w którym stoi treść, więc bez nich
 *     „zmiany autora model przyjmowane po kolei" nie ma czym działać;
 *   — `studio.proposal.decide` — decyzja o propozycji po stronie RDZENIA,
 *     z fragmentami wskazanymi wybiórczo (`hunkIndexes`);
 *   — `studio.export.profile.*` — nastawy strony, jedyne miejsce w kontrakcie,
 *     w którym kartka, marginesy, nagłówek i stopka mają zapis trwały;
 *   — `studio.preview.render` — paginacja i typografia liczone przez rdzeń;
 *   — `studio.template.*` — galeria szablonów i zakładanie dokumentu z szablonu;
 *   — `studio.search.semantic` i `studio.diff.source` — pomiary panelu
 *     Redaktora: bliskość znaczeniowa i rozbieżność wobec materiału wejściowego.
 *
 * Źródło nie ma stanu i nie buduje elementu — sprawdza kształt odpowiedzi
 * i oddaje ją oknu.
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
  /**
   * Pobiera TREŚĆ zasobu magazynu rdzenia (`design.asset.content.get`).
   *
   * Rodzina komendy nazywa się `design.*`, ale komenda nie należy do modułu
   * Design: kontrakt mówi wprost, że dotyczy KAŻDEGO zasobu magazynu, bo magazyn
   * jest jeden dla rodzin `design.*`, `document.*`, `media.*` i `archive.*`.
   * Podgląd wydania sięga nią po obrazy stron wyrysowane przez rdzeń — pole `uri`
   * zasobu jest ścieżką w systemie plików rdzenia, więc przeglądarka nie wczyta
   * spod niego niczego.
   *
   * Druga komenda „pobierz zasób Studia" byłaby drugą drogą do tego samego
   * magazynu, więc jej tu nie ma i nie ma jej w kontrakcie.
   */
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
          // Wykaz pusty znaczy „w całości" i pola wtedy nie ma: kontrakt czyta
          // brak `hunkIndexes` jako całość, a wykaz pusty byłby wskazaniem
          // zera fragmentów.
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
        // `contentBase64` jest w kontrakcie nieobowiązkowe (postać odsyłania go
        // nie niesie), więc sprawdzenie kształtu pyta o pola OBOWIĄZKOWE. Brak
        // treści przy postaci „content" rozstrzyga wołający — i mówi o tym wprost,
        // zamiast rysować pustą kartkę.
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
          // Zakres jedzie wyłącznie wtedy, gdy jest: komentarz bez zaznaczenia
          // dotyczy dokumentu, a zakres 0–0 byłby kotwicą przy pierwszym znaku.
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
