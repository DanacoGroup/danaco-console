import {
  Command,
  EventType,
  type DesignAnnotationListRequest,
  type DesignAnnotationListResponse,
  type DesignAnnotationSetRequest,
  type DesignAnnotationSetResponse,
  type DesignAssetChangedEvent,
  type DesignAssetContentGetRequest,
  type DesignAssetContentGetResponse,
  type DesignAssetExportBatchRequest,
  type DesignAssetExportBatchResponse,
  type DesignAssetExportRequest,
  type DesignAssetExportResponse,
  type DesignAssetGenerateRequest,
  type DesignAssetGenerateResponse,
  type DesignAssetKind,
  type DesignAssetListRequest,
  type DesignAssetListResponse,
  type DesignAssetFavoriteSetRequest,
  type DesignAssetFavoriteSetResponse,
  type DesignAssetRemoveRequest,
  type DesignAssetRemoveResponse,
  type DesignAssetTagSetRequest,
  type DesignAssetTagSetResponse,
  type DesignAssetUploadRequest,
  type DesignAssetUploadResponse,
  type DesignBoardExportRequest,
  type DesignBoardExportResponse,
  type DesignBoardLayer,
  type DesignBoardListRequest,
  type DesignBoardListResponse,
  type DesignBoardPresenceEvent,
  type DesignBoardRegion,
  type DesignBoardUpdateRequest,
  type DesignBoardUpdateResponse,
  type DesignBoardVersionListRequest,
  type DesignBoardVersionListResponse,
  type DesignBoardVersionRestoreRequest,
  type DesignBoardVersionRestoreResponse,
  type DesignBoardVersionSaveRequest,
  type DesignBoardVersionSaveResponse,
  type DesignCollectionAssignRequest,
  type DesignCollectionAssignResponse,
  type DesignCollectionCreateRequest,
  type DesignCollectionCreateResponse,
  type DesignCollectionListRequest,
  type DesignCollectionListResponse,
  type DesignPresenceReportRequest,
  type DesignPresenceReportResponse,
  type DesignPrompt,
  type DesignPromptHistoryListRequest,
  type DesignPromptHistoryListResponse,
  type DesignPromptTemplateListRequest,
  type DesignPromptTemplateListResponse,
  type DesignPromptTemplateSaveRequest,
  type DesignPromptTemplateSaveResponse,
  type DesignStyleguidePublishRequest,
  type DesignStyleguidePublishResponse,
  type DesignToken,
  type DesignTokenTarget,
  type DesignTokensetExportRequest,
  type DesignTokensetExportResponse,
  type DesignTokensetImportRequest,
  type DesignTokensetImportResponse,
  type DesignTokensetListRequest,
  type DesignTokensetListResponse,
  type DesignTokensetSaveRequest,
  type DesignTokensetSaveResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { przytnijKod } from './przyciecie-pol';

/**
 * Osiem komend obszaru `design.*` widzianych przez okna modułu.
 *
 * Źródło to warstwa wywołań i sprawdzianu kształtu odpowiedzi. Nie ma stanu
 * i niczego nie pamięta — zbiór zasobów mieszka w `stan-designu.ts`, żeby trzy
 * okna patrzyły na jeden zbiór.
 *
 * Zasób wchodzi do modułu dwiema drogami, które się nie zastępują:
 *   `design.asset.generate` — zasób z pracy modelu; wymaga kanału obrazowego
 *                             (adapter „obrazy"), a jego brak jest odmową
 *                             nazywającą brak, nigdy obrazem zastępczym;
 *   `design.asset.upload`   — zasób z pliku wskazanego przez Operatora; działa
 *                             bez żadnego kanału modelu i wnosi treść już
 *                             istniejącą.
 *
 * Przy wgraniu Operator wskazuje plik, klient czyta go i oddaje treść w polu
 * `contentBase64`, a rdzeń zapisuje ją w swoim magazynie pod sumą kontrolną.
 * Od tej chwili zasób nie zależy już od pliku na dysku, który Operator może
 * nadpisać albo skasować.
 *
 * Treść idzie bajtami, nie ścieżką, mimo że kontrakt zna oba pola. `sourcePath`
 * każe rdzeniowi otworzyć plik pod wskazaną ścieżką, a rdzeń biegnie na innej
 * maszynie niż przeglądarka Operatora: ścieżka z jego pulpitu nie znaczy tam
 * nic albo znaczy coś zupełnie innego. Klient czyta plik sam i wysyła to, co
 * przeczytał (`wczytanie-pliku.ts`).
 *
 * Odmowa wraca polem `blad`, nie wyjątkiem. Rdzeń odmawia czynności na zasobie,
 * którego nie zna; okno pokazuje to zdanie i zostaje czynne.
 */

/** Warunki zawężające odczyt zasobów — pola żądania `design.asset.list`. */
export interface ZapytanieZasobow {
  idOkna: string;
  rodzaj: DesignAssetKind | '';
  etykiety: readonly string[];
  tylkoUlubione: boolean;
  granica: number;
}

/** Zlecenie zapisu kompozycji — pola żądania `design.board.update`. */
export interface ZlecenieKompozycji {
  idOkna: string;
  idKompozycji: string;
  nazwa: string;
  warstwy: readonly DesignBoardLayer[];
}

/**
 * Zlecenie generowania — pola żądania `design.asset.generate`.
 *
 * Wskazanie kanału jedzie polem `channelId`: rdzeń sprawdza po nim rodzaj
 * kanału i odmawia wskazania tekstowego. Puste `idKanalu` nie jest brakiem —
 * kontrakt opisuje pominięcie pola jako „weź pierwszy czynny kanał obrazowy
 * konta" i tak właśnie zachowuje się okno, gdy Operator kanału nie wskazał.
 */
export interface ZlecenieGenerowania {
  idOkna: string;
  prompt: DesignPrompt;
  idReferencji: string;
  /** Kanał obrazowy zlecenia; pusty zostawia wybór rdzeniowi (pole pominięte). */
  idKanalu: string;
  /** Rodzaj zasobu, który ma powstać; pusty zostawia rdzeniowi rodzaj obrazu. */
  rodzaj: DesignAssetKind | '';
}

/**
 * Zlecenie nadania etykiet — pola żądania `design.asset.tag.set`.
 *
 * `etykiety` to zestaw docelowy, nie dokładka: zestaw pusty znaczy „zdejmij
 * wszystkie", a nie „nic nie rób".
 */
export interface ZlecenieEtykiet {
  idZasobu: string;
  etykiety: readonly string[];
}

/**
 * Zlecenie wgrania zasobu — pola żądania `design.asset.upload`.
 *
 * Treść jest tu zawsze bajtami pliku w zapisie base64. Pola opisowe (nazwa,
 * format, wymiary) idą tylko wtedy, gdy klient je zmierzył; podstawienie
 * wartości domyślnej byłoby wpisaniem rdzeniowi metadanych, których nikt nie
 * sprawdził.
 */
export interface ZlecenieWgrania {
  idOkna: string;
  nazwa: string;
  rodzaj: DesignAssetKind;
  trescBase64: string;
  format: string;
  szerokosc: number;
  wysokosc: number;
  etykiety: readonly string[];
}

/**
 * Zlecenie wydania zasobu — pola żądania `design.asset.export`.
 *
 * Skala zerowa znaczy „naturalna", nie „zerowa": pole jest w kontrakcie
 * opcjonalne, a zero wysłane do rdzenia byłoby żądaniem obrazu o zerowym boku
 * i wróciłoby odmową. Tak samo jakość — zero nie jest stopniem kompresji.
 */
export interface ZlecenieWydania {
  idZasobu: string;
  format: string;
  skala: number;
  jakosc: number;
}

/** Zlecenie wydania partii — pola żądania `design.asset.export.batch`. */
export interface ZleceniePartii {
  idZasobow: readonly string[];
  format: string;
  skale: readonly number[];
  jakosc: number;
}

/** Zlecenie przypisania do kolekcji — pola żądania `design.collection.assign`. */
export interface ZlecenieKolekcji {
  idKolekcji: string;
  idZasobow: readonly string[];
  zdejmij: boolean;
}

/** Zlecenie zapisu szablonu promptu — pola żądania `design.prompt.template.save`. */
export interface ZlecenieSzablonu {
  idOkna: string;
  nazwa: string;
  prompt: DesignPrompt;
  /** Szablon nadpisywany; pusty zakłada nowy. */
  idSzablonu: string;
}

/** Zlecenie zapisu wersji kompozycji — pola żądania `design.board.version.save`. */
export interface ZlecenieWersji {
  idKompozycji: string;
  nazwa: string;
  uzasadnienie: string;
}

/** Zlecenie wyrysu kompozycji — pola żądania `design.board.export`. */
export interface ZlecenieWyrysu {
  idKompozycji: string;
  format: string;
  skala: number;
  /** Obszar wyrysu; `null` bierze całość. */
  obszar: DesignBoardRegion | null;
}

/** Zlecenie adnotacji — pola żądania `design.annotation.set`. */
export interface ZlecenieAdnotacji {
  idKompozycji: string;
  tresc: string;
  idWarstwy: string;
  idAdnotacji: string;
  idNadrzednej: string;
  /** `null` zostawia stan zamknięcia bez zmiany (pole pominięte). */
  zamknieta: boolean | null;
}

/** Zlecenie zgłoszenia obecności — pola żądania `design.presence.report`. */
export interface ZlecenieObecnosci {
  idKompozycji: string;
  x: number;
  y: number;
  zaznaczone: readonly string[];
  odchodzi: boolean;
}

/** Zlecenie zapisu zestawu żetonów — pola żądania `design.tokenset.save`. */
export interface ZlecenieZestawu {
  idOkna: string;
  nazwa: string;
  zetony: readonly DesignToken[];
  /** Zestaw nadpisywany; pusty zakłada nowy. */
  idZestawu: string;
  motyw: string;
}

/** Zlecenie wczytania zestawu — pola żądania `design.tokenset.import`. */
export interface ZlecenieWczytaniaZetonow {
  idOkna: string;
  nazwa: string;
  trescBase64: string;
  /** Postać zapisu; pusta zostawia rozpoznanie rdzeniowi. */
  postac: string;
}

export interface ZrodloDesignu {
  zasoby(zapytanie: ZapytanieZasobow): Promise<Wynik<DesignAssetListResponse>>;
  generuj(zlecenie: ZlecenieGenerowania): Promise<Wynik<DesignAssetGenerateResponse>>;
  zapiszKompozycje(zlecenie: ZlecenieKompozycji): Promise<Wynik<DesignBoardUpdateResponse>>;
  ustawEtykiety(zlecenie: ZlecenieEtykiet): Promise<Wynik<DesignAssetTagSetResponse>>;
  /** Wnosi do modułu zasób z pliku wskazanego przez Operatora. */
  wgrajZasob(zlecenie: ZlecenieWgrania): Promise<Wynik<DesignAssetUploadResponse>>;
  /** Usuwa zasób z Assets Panel; wykonuje się bez pytania o potwierdzenie. */
  usunZasob(idZasobu: string): Promise<Wynik<DesignAssetRemoveResponse>>;
  /** Ustawia albo zdejmuje oznaczenie „ulubiony". */
  ustawUlubiony(idZasobu: string, ulubiony: boolean): Promise<Wynik<DesignAssetFavoriteSetResponse>>;
  /** Odczyt kompozycji okna wraz z warstwami — powrót zapisanej planszy. */
  kompozycje(idOkna: string): Promise<Wynik<DesignBoardListResponse>>;
  /** Oddaje treść zasobu z magazynu rdzenia — jedyna droga do jego bajtów. */
  trescZasobu(idZasobu: string, granicaBajtow: number): Promise<Wynik<DesignAssetContentGetResponse>>;
  /** Wydaje zasób jako plik w wskazanym formacie i skali. */
  wydajZasob(zlecenie: ZlecenieWydania): Promise<Wynik<DesignAssetExportResponse>>;
  /** Wydaje partię zasobów w komplecie skal; odmowa jednego nie wstrzymuje reszty. */
  wydajPartie(zlecenie: ZleceniePartii): Promise<Wynik<DesignAssetExportBatchResponse>>;

  /** Zakłada kolekcję zasobów. */
  zalozKolekcje(idOkna: string, nazwa: string, opis: string): Promise<Wynik<DesignCollectionCreateResponse>>;
  /** Dokłada zasoby do kolekcji albo je z niej zdejmuje. */
  przypiszDoKolekcji(zlecenie: ZlecenieKolekcji): Promise<Wynik<DesignCollectionAssignResponse>>;
  /** Odczyt kolekcji okna; wskazanie zasobu zawęża do kolekcji, które go mają. */
  kolekcje(idOkna: string, idZasobu: string): Promise<Wynik<DesignCollectionListResponse>>;

  /** Utrwala prompt jako szablon do wielokrotnego użycia. */
  zapiszSzablonPromptu(zlecenie: ZlecenieSzablonu): Promise<Wynik<DesignPromptTemplateSaveResponse>>;
  /** Odczyt szablonów promptów okna. */
  szablonyPromptu(idOkna: string): Promise<Wynik<DesignPromptTemplateListResponse>>;
  /** Odczyt promptów wydanych w oknie wraz z zasobami, które z nich powstały. */
  historiaPromptow(idOkna: string, granica: number): Promise<Wynik<DesignPromptHistoryListResponse>>;

  /** Utrwala bieżący układ kompozycji jako nazwaną wersję. */
  zapiszWersje(zlecenie: ZlecenieWersji): Promise<Wynik<DesignBoardVersionSaveResponse>>;
  /** Odczyt wersji kompozycji — bez warstw. */
  wersje(idKompozycji: string, granica: number): Promise<Wynik<DesignBoardVersionListResponse>>;
  /** Przywraca układ z wersji; odkłada przy tym wersję sprzed przywrócenia. */
  przywrocWersje(idWersji: string): Promise<Wynik<DesignBoardVersionRestoreResponse>>;
  /** Wyrysowuje kompozycję albo jej obszar do jednego pliku. */
  wyrysujKompozycje(zlecenie: ZlecenieWyrysu): Promise<Wynik<DesignBoardExportResponse>>;

  /** Zakłada albo zmienia adnotację przypiętą do warstwy lub kompozycji. */
  ustawAdnotacje(zlecenie: ZlecenieAdnotacji): Promise<Wynik<DesignAnnotationSetResponse>>;
  /** Odczyt adnotacji kompozycji wraz z wątkami. */
  adnotacje(idKompozycji: string, tylkoOtwarte: boolean): Promise<Wynik<DesignAnnotationListResponse>>;
  /** Zgłasza obecność i położenie kursora na kompozycji; zgłoszenie jest ulotne. */
  zglosObecnosc(zlecenie: ZlecenieObecnosci): Promise<Wynik<DesignPresenceReportResponse>>;

  /** Utrwala zestaw żetonów systemu projektowego. */
  zapiszZestawZetonow(zlecenie: ZlecenieZestawu): Promise<Wynik<DesignTokensetSaveResponse>>;
  /** Odczyt zestawów żetonów okna wraz z ich żetonami. */
  zestawyZetonow(idOkna: string, idZestawu: string): Promise<Wynik<DesignTokensetListResponse>>;
  /** Wydaje zestaw żetonów w postaci przyjmowanej przez kod. */
  wydajZetony(
    idZestawu: string,
    postac: DesignTokenTarget,
    idModuluDocelowego: string,
  ): Promise<Wynik<DesignTokensetExportResponse>>;
  /** Wczytuje zestaw żetonów z zapisu zewnętrznego. */
  wczytajZetony(zlecenie: ZlecenieWczytaniaZetonow): Promise<Wynik<DesignTokensetImportResponse>>;
  /** Wydaje przewodnik systemu projektowego do modułu docelowego. */
  wydajPrzewodnik(
    idZestawu: string,
    idModuluDocelowego: string,
    idKolekcji: string,
  ): Promise<Wynik<DesignStyleguidePublishResponse>>;

  /** Subskrypcja `design.asset.changed` — zdarzenia zasobu. */
  naZmianeZasobu(sluchacz: (tresc: DesignAssetChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja `design.board.presence` — kursorów współpracy na kompozycji. */
  naObecnosc(sluchacz: (tresc: DesignBoardPresenceEvent) => void): Odsubskrybuj;
}

export function utworzZrodloDesignu(kanal: Kanal): ZrodloDesignu {
  return {
    async zasoby(zapytanie) {
      const zadanie: DesignAssetListRequest = { favoriteOnly: zapytanie.tylkoUlubione };
      if (zapytanie.idOkna !== '') zadanie.windowId = zapytanie.idOkna;
      if (zapytanie.rodzaj !== '') zadanie.kind = zapytanie.rodzaj;
      if (zapytanie.etykiety.length > 0) zadanie.tags = [...zapytanie.etykiety];
      if (zapytanie.granica > 0) zadanie.limit = zapytanie.granica;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetList, zadanie),
        Command.DesignAssetList,
        (tresc) => czyTablica(tresc.assets),
      );
    },

    async generuj(zlecenie) {
      const zadanie: DesignAssetGenerateRequest = {
        windowId: zlecenie.idOkna,
        prompt: zlecenie.prompt,
      };
      if (zlecenie.idReferencji !== '') zadanie.referenceAssetId = zlecenie.idReferencji;
      // Pola opcjonalne idą wyłącznie wskazane. Pustka wpisana w `channelId`
      // nie jest tym samym co pole pominięte: rdzeń przycina wskazanie i puste
      // traktuje jak brak (`kanalObrazowyZadania`), ale wysyłanie pustego napisu
      // twierdziłoby, że Operator coś wskazał. Pominięcie pola mówi prawdę —
      // wyboru nie było, więc rdzeń bierze pierwszy czynny kanał obrazowy.
      if (zlecenie.idKanalu !== '') zadanie.channelId = zlecenie.idKanalu;
      if (zlecenie.rodzaj !== '') zadanie.kind = zlecenie.rodzaj;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetGenerate, zadanie),
        Command.DesignAssetGenerate,
        (tresc) => czyTablica(tresc.assets),
      );
    },

    async zapiszKompozycje(zlecenie) {
      const zadanie: DesignBoardUpdateRequest = { windowId: zlecenie.idOkna };
      if (zlecenie.idKompozycji !== '') zadanie.boardId = zlecenie.idKompozycji;
      // Rdzeń nazwy nie przycina, więc przycięcie jest tutaj jedyne i to ono
      // musi zgadzać się z porównaniem skutku zapisu.
      const nazwa = przytnijKod(zlecenie.nazwa);
      if (nazwa !== '') zadanie.name = nazwa;
      zadanie.layers = [...zlecenie.warstwy];
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignBoardUpdate, zadanie),
        Command.DesignBoardUpdate,
        (tresc) => czyObiekt(tresc.board),
      );
    },

    async ustawEtykiety(zlecenie) {
      // Zestaw idzie do rdzenia zawsze, także pusty — inaczej Operator nie
      // miałby czym zdjąć ostatniej etykiety, a pominięcie pola byłoby
      // milczącym nieposłuszeństwem wobec kontraktu (`tags` jest wymagane).
      const zadanie: DesignAssetTagSetRequest = {
        assetId: zlecenie.idZasobu,
        tags: [...zlecenie.etykiety],
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetTagSet, zadanie),
        Command.DesignAssetTagSet,
        (tresc) => czyObiekt(tresc.asset),
      );
    },

    async wgrajZasob(zlecenie) {
      const zadanie: DesignAssetUploadRequest = {
        windowId: zlecenie.idOkna,
        kind: zlecenie.rodzaj,
        contentBase64: zlecenie.trescBase64,
      };
      // Nazwa przycinana tak samo jak nazwa kompozycji (`przyciecie-pol.ts`),
      // żeby nazwa nadana przy wgraniu trafiała we własne zawężenie wykazu.
      const nazwa = przytnijKod(zlecenie.nazwa);
      if (nazwa !== '') zadanie.name = nazwa;
      if (zlecenie.format !== '') zadanie.format = zlecenie.format;
      // Wymiary idą wyłącznie zmierzone. Zero jest tu znakiem „nie zmierzono",
      // a nie szerokością — wpisane do żądania byłoby metadaną zmyśloną.
      if (zlecenie.szerokosc > 0) zadanie.width = zlecenie.szerokosc;
      if (zlecenie.wysokosc > 0) zadanie.height = zlecenie.wysokosc;
      if (zlecenie.etykiety.length > 0) zadanie.tags = [...zlecenie.etykiety];
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetUpload, zadanie),
        Command.DesignAssetUpload,
        (tresc) => czyObiekt(tresc.asset),
      );
    },

    async usunZasob(idZasobu) {
      const zadanie: DesignAssetRemoveRequest = { assetId: idZasobu };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetRemove, zadanie),
        Command.DesignAssetRemove,
        // Pole `removed` jest logiczne i niesie treść: `false` znaczy „zasobu
        // nie było", co jest odpowiedzią udaną i innym zdaniem niż odmowa.
        // Sprawdzamy więc typ, a nie prawdziwość.
        (tresc) => typeof tresc.removed === 'boolean',
      );
    },

    async ustawUlubiony(idZasobu, ulubiony) {
      const zadanie: DesignAssetFavoriteSetRequest = { assetId: idZasobu, favorite: ulubiony };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetFavoriteSet, zadanie),
        Command.DesignAssetFavoriteSet,
        (tresc) => czyObiekt(tresc.asset),
      );
    },

    async kompozycje(idOkna) {
      const zadanie: DesignBoardListRequest = { windowId: idOkna };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignBoardList, zadanie),
        Command.DesignBoardList,
        (tresc) => czyTablica(tresc.boards),
      );
    },

    async trescZasobu(idZasobu, granicaBajtow) {
      const zadanie: DesignAssetContentGetRequest = { assetId: idZasobu };
      // Granica idzie wyłącznie wskazana. Zero wysłane do rdzenia byłoby
      // żądaniem treści o zerowej wielkości, a nie brakiem granicy.
      if (granicaBajtow > 0) zadanie.maxBytes = granicaBajtow;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetContentGet, zadanie),
        Command.DesignAssetContentGet,
        // Sprawdzamy sumę i miarę, nie samą treść: przy postaci odsyłania
        // `contentBase64` jest puste zgodnie z kontraktem, a odpowiedź jest
        // wtedy pełna.
        (tresc) => typeof tresc.checksum === 'string' && typeof tresc.sizeBytes === 'number',
      );
    },

    async wydajZasob(zlecenie) {
      const zadanie: DesignAssetExportRequest = {
        assetId: zlecenie.idZasobu,
        format: zlecenie.format,
      };
      if (zlecenie.skala > 0) zadanie.scale = zlecenie.skala;
      if (zlecenie.jakosc > 0) zadanie.quality = zlecenie.jakosc;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetExport, zadanie),
        Command.DesignAssetExport,
        // Wydanie bez treści jest kopertą udaną i pustą — sprawdzamy bajty,
        // a nie samo powodzenie.
        (tresc) => typeof tresc.contentBase64 === 'string' && tresc.contentBase64.length > 0,
      );
    },

    async wydajPartie(zlecenie) {
      const zadanie: DesignAssetExportBatchRequest = {
        assetIds: [...zlecenie.idZasobow],
        format: zlecenie.format,
      };
      if (zlecenie.skale.length > 0) zadanie.scales = [...zlecenie.skale];
      if (zlecenie.jakosc > 0) zadanie.quality = zlecenie.jakosc;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAssetExportBatch, zadanie),
        Command.DesignAssetExportBatch,
        (tresc) => czyTablica(tresc.files) && typeof tresc.exported === 'number',
      );
    },

    async zalozKolekcje(idOkna, nazwa, opis) {
      const zadanie: DesignCollectionCreateRequest = {
        windowId: idOkna,
        name: przytnijKod(nazwa),
      };
      const przyciety = przytnijKod(opis);
      if (przyciety !== '') zadanie.description = przyciety;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignCollectionCreate, zadanie),
        Command.DesignCollectionCreate,
        (tresc) => czyObiekt(tresc.collection),
      );
    },

    async przypiszDoKolekcji(zlecenie) {
      const zadanie: DesignCollectionAssignRequest = {
        collectionId: zlecenie.idKolekcji,
        assetIds: [...zlecenie.idZasobow],
      };
      // Zdjęcie idzie polem wskazanym; dokładka jest brakiem pola, bo tak
      // opisuje ją kontrakt („brak znaczy dołóż").
      if (zlecenie.zdejmij) zadanie.remove = true;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignCollectionAssign, zadanie),
        Command.DesignCollectionAssign,
        (tresc) => czyObiekt(tresc.collection) && typeof tresc.changed === 'number',
      );
    },

    async kolekcje(idOkna, idZasobu) {
      const zadanie: DesignCollectionListRequest = { windowId: idOkna };
      if (idZasobu !== '') zadanie.assetId = idZasobu;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignCollectionList, zadanie),
        Command.DesignCollectionList,
        (tresc) => czyTablica(tresc.collections),
      );
    },

    async zapiszSzablonPromptu(zlecenie) {
      const zadanie: DesignPromptTemplateSaveRequest = {
        windowId: zlecenie.idOkna,
        name: przytnijKod(zlecenie.nazwa),
        prompt: zlecenie.prompt,
      };
      if (zlecenie.idSzablonu !== '') zadanie.templateId = zlecenie.idSzablonu;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignPromptTemplateSave, zadanie),
        Command.DesignPromptTemplateSave,
        (tresc) => czyObiekt(tresc.template),
      );
    },

    async szablonyPromptu(idOkna) {
      const zadanie: DesignPromptTemplateListRequest = { windowId: idOkna };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignPromptTemplateList, zadanie),
        Command.DesignPromptTemplateList,
        (tresc) => czyTablica(tresc.templates),
      );
    },

    async historiaPromptow(idOkna, granica) {
      const zadanie: DesignPromptHistoryListRequest = { windowId: idOkna };
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignPromptHistoryList, zadanie),
        Command.DesignPromptHistoryList,
        (tresc) => czyTablica(tresc.prompts),
      );
    },

    async zapiszWersje(zlecenie) {
      const zadanie: DesignBoardVersionSaveRequest = { boardId: zlecenie.idKompozycji };
      const nazwa = przytnijKod(zlecenie.nazwa);
      if (nazwa !== '') zadanie.name = nazwa;
      const uzasadnienie = przytnijKod(zlecenie.uzasadnienie);
      if (uzasadnienie !== '') zadanie.note = uzasadnienie;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignBoardVersionSave, zadanie),
        Command.DesignBoardVersionSave,
        (tresc) => czyObiekt(tresc.version),
      );
    },

    async wersje(idKompozycji, granica) {
      const zadanie: DesignBoardVersionListRequest = { boardId: idKompozycji };
      if (granica > 0) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignBoardVersionList, zadanie),
        Command.DesignBoardVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
    },

    async przywrocWersje(idWersji) {
      const zadanie: DesignBoardVersionRestoreRequest = { versionId: idWersji };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignBoardVersionRestore, zadanie),
        Command.DesignBoardVersionRestore,
        (tresc) => czyObiekt(tresc.board),
      );
    },

    async wyrysujKompozycje(zlecenie) {
      const zadanie: DesignBoardExportRequest = {
        boardId: zlecenie.idKompozycji,
        format: zlecenie.format,
      };
      if (zlecenie.skala > 0) zadanie.scale = zlecenie.skala;
      if (zlecenie.obszar !== null) zadanie.region = zlecenie.obszar;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignBoardExport, zadanie),
        Command.DesignBoardExport,
        (tresc) => typeof tresc.contentBase64 === 'string' && tresc.contentBase64.length > 0,
      );
    },

    async ustawAdnotacje(zlecenie) {
      const zadanie: DesignAnnotationSetRequest = {
        boardId: zlecenie.idKompozycji,
        text: zlecenie.tresc,
      };
      if (zlecenie.idWarstwy !== '') zadanie.layerId = zlecenie.idWarstwy;
      if (zlecenie.idAdnotacji !== '') zadanie.annotationId = zlecenie.idAdnotacji;
      if (zlecenie.idNadrzednej !== '') zadanie.parentId = zlecenie.idNadrzednej;
      // Stan zamknięcia idzie wyłącznie wskazany: `false` wysłane przy zapisie
      // treści otwierałoby wątek, którego nikt nie kazał otwierać.
      if (zlecenie.zamknieta !== null) zadanie.resolved = zlecenie.zamknieta;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAnnotationSet, zadanie),
        Command.DesignAnnotationSet,
        (tresc) => czyObiekt(tresc.annotation),
      );
    },

    async adnotacje(idKompozycji, tylkoOtwarte) {
      const zadanie: DesignAnnotationListRequest = { boardId: idKompozycji };
      if (tylkoOtwarte) zadanie.openOnly = true;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignAnnotationList, zadanie),
        Command.DesignAnnotationList,
        (tresc) => czyTablica(tresc.annotations),
      );
    },

    async zglosObecnosc(zlecenie) {
      const zadanie: DesignPresenceReportRequest = { boardId: zlecenie.idKompozycji };
      // Położenie idzie tylko zmierzone. Zero jest tu prawdziwym położeniem
      // (lewy górny róg kanwy), więc rozstrzyga skończoność liczby, nie jej
      // wartość — inaczej kursor w rogu znikałby z widoku pozostałych.
      if (Number.isFinite(zlecenie.x)) zadanie.x = zlecenie.x;
      if (Number.isFinite(zlecenie.y)) zadanie.y = zlecenie.y;
      if (zlecenie.zaznaczone.length > 0) zadanie.selectedLayerIds = [...zlecenie.zaznaczone];
      if (zlecenie.odchodzi) zadanie.leaving = true;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignPresenceReport, zadanie),
        Command.DesignPresenceReport,
        (tresc) => czyTablica(tresc.participants),
      );
    },

    async zapiszZestawZetonow(zlecenie) {
      const zadanie: DesignTokensetSaveRequest = {
        windowId: zlecenie.idOkna,
        name: przytnijKod(zlecenie.nazwa),
        tokens: [...zlecenie.zetony],
      };
      if (zlecenie.idZestawu !== '') zadanie.tokenSetId = zlecenie.idZestawu;
      if (zlecenie.motyw !== '') zadanie.theme = zlecenie.motyw;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignTokensetSave, zadanie),
        Command.DesignTokensetSave,
        (tresc) => czyObiekt(tresc.tokenSet),
      );
    },

    async zestawyZetonow(idOkna, idZestawu) {
      const zadanie: DesignTokensetListRequest = { windowId: idOkna };
      if (idZestawu !== '') zadanie.tokenSetId = idZestawu;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignTokensetList, zadanie),
        Command.DesignTokensetList,
        (tresc) => czyTablica(tresc.tokenSets),
      );
    },

    async wydajZetony(idZestawu, postac, idModuluDocelowego) {
      const zadanie: DesignTokensetExportRequest = { tokenSetId: idZestawu, target: postac };
      if (idModuluDocelowego !== '') zadanie.targetModuleId = idModuluDocelowego;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignTokensetExport, zadanie),
        Command.DesignTokensetExport,
        (tresc) => typeof tresc.content === 'string' && tresc.content.length > 0,
      );
    },

    async wczytajZetony(zlecenie) {
      const zadanie: DesignTokensetImportRequest = {
        windowId: zlecenie.idOkna,
        name: przytnijKod(zlecenie.nazwa),
        contentBase64: zlecenie.trescBase64,
      };
      if (zlecenie.postac !== '') zadanie.format = zlecenie.postac;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignTokensetImport, zadanie),
        Command.DesignTokensetImport,
        (tresc) => czyObiekt(tresc.tokenSet),
      );
    },

    async wydajPrzewodnik(idZestawu, idModuluDocelowego, idKolekcji) {
      const zadanie: DesignStyleguidePublishRequest = {
        tokenSetId: idZestawu,
        targetModuleId: idModuluDocelowego,
      };
      if (idKolekcji !== '') zadanie.collectionId = idKolekcji;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DesignStyleguidePublish, zadanie),
        Command.DesignStyleguidePublish,
        // Zasób przewodnika jest tu treścią odpowiedzi: `published: true` bez
        // identyfikatora zasobu byłoby obietnicą wydania, którego nie da się
        // odnaleźć.
        (tresc) => typeof tresc.assetId === 'string' && tresc.assetId.length > 0,
      );
    },

    naZmianeZasobu(sluchacz) {
      return kanal.naZdarzenie(EventType.DesignAssetChanged, (tresc) => sluchacz(tresc));
    },

    naObecnosc(sluchacz) {
      return kanal.naZdarzenie(EventType.DesignBoardPresence, (tresc) => sluchacz(tresc));
    },
  };
}
