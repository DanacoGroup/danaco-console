import {
  Command,
  type StudioActionBalance,
  type StudioDocumentCopyRequest,
  type StudioDocumentCopyResponse,
  type StudioDocumentCreateRequest,
  type StudioDocumentCreateResponse,
  type StudioDocumentExportBatchRequest,
  type StudioDocumentExportBatchResponse,
  type StudioDocumentExportFormatRequest,
  type StudioDocumentExportFormatResponse,
  type StudioDocumentFormatSetRequest,
  type StudioDocumentFormatSetResponse,
  type StudioDocumentImageImportRequest,
  type StudioDocumentImageImportResponse,
  type StudioDocumentImportFileRequest,
  type StudioDocumentImportFileResponse,
  type StudioDocumentImportPdfRequest,
  type StudioDocumentImportPdfResponse,
  type StudioDocumentSaveAsRequest,
  type StudioDocumentSaveAsResponse,
  type StudioExportResult,
  type StudioImportBalance,
  type StudioIngestCorrectionSetRequest,
  type StudioIngestCorrectionSetResponse,
  type StudioIngestDeviceListResponse,
  type StudioIngestItemAcceptRequest,
  type StudioIngestItemAcceptResponse,
  type StudioIngestQueueAddRequest,
  type StudioIngestQueueAddResponse,
  type StudioIngestQueueListRequest,
  type StudioIngestQueueListResponse,
  type StudioIngestRecognizeRequest,
  type StudioIngestRecognizeResponse,
  type StudioInsertFromLibraryRequest,
  type StudioInsertFromLibraryResponse,
  type StudioInsertFromWebRequest,
  type StudioInsertFromWebResponse,
  type StudioSkippedItem,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/** Interfejs ZrodloWstawienStudio zestawia siedemnaście komend studio wnoszących treść do dokumentu, wydających go oraz prowadzących kolejkę cyfryzacji rdzenia. */
export interface ZrodloWstawienStudio {
  /* ── Wejście do edytora ────────────────────────────────────────────────── */

  /** Zakłada dokument pusty — nową stronę gotową do pisania. */
  zalozDokument(
    zadanie: StudioDocumentCreateRequest,
  ): Promise<Wynik<StudioDocumentCreateResponse>>;
  /** Wnosi plik Operatora wprost do edytora, z zachowaną postacią. */
  wniesPlik(
    zadanie: StudioDocumentImportFileRequest,
  ): Promise<Wynik<StudioDocumentImportFileResponse>>;
  /** Zamienia PDF na dokument edytowalny wprost w edytorze. */
  wniesPdf(
    zadanie: StudioDocumentImportPdfRequest,
  ): Promise<Wynik<StudioDocumentImportPdfResponse>>;
  /** Wnosi obraz w miejsce kursora — z pliku, magazynu, Designu, bazy zdjęciowej. */
  wniesObraz(
    zadanie: StudioDocumentImageImportRequest,
  ): Promise<Wynik<StudioDocumentImageImportResponse>>;

  /* ── Zapis, kopia, wydanie ─────────────────────────────────────────────── */

  /** Zapisuje dokument pod nową nazwą albo do wskazanego pliku, wraz z postacią. */
  zapiszJako(
    zadanie: StudioDocumentSaveAsRequest,
  ): Promise<Wynik<StudioDocumentSaveAsResponse>>;
  /** Zakłada kopię dokumentu — osobny dokument, nie drugie odwołanie. */
  skopiuj(zadanie: StudioDocumentCopyRequest): Promise<Wynik<StudioDocumentCopyResponse>>;
  /** Wydaje dokument do formatu docelowego wraz z wykazem cech pominiętych. */
  wydaj(
    zadanie: StudioDocumentExportFormatRequest,
  ): Promise<Wynik<StudioDocumentExportFormatResponse>>;
  /** Wydaje wiele dokumentów naraz; odmowa jednego nie wstrzymuje pozostałych. */
  wydajWsad(
    zadanie: StudioDocumentExportBatchRequest,
  ): Promise<Wynik<StudioDocumentExportBatchResponse>>;
  /** Przestawia format dokumentu Studia. */
  przestawFormat(
    zadanie: StudioDocumentFormatSetRequest,
  ): Promise<Wynik<StudioDocumentFormatSetResponse>>;

  /* ── Wniesienie ze źródła ──────────────────────────────────────────────── */

  /** Wnosi plik, wzór albo obraz z Biblioteki wprost do dokumentu. */
  wniesZBiblioteki(
    zadanie: StudioInsertFromLibraryRequest,
  ): Promise<Wynik<StudioInsertFromLibraryResponse>>;
  /** Wnosi fragment albo obraz ze strony sieci wprost do dokumentu. */
  wniesZeStrony(
    zadanie: StudioInsertFromWebRequest,
  ): Promise<Wynik<StudioInsertFromWebResponse>>;

  /* ── Cyfryzacja ────────────────────────────────────────────────────────── */

  /** Dokłada materiał do kolejki wczytywania po stronie rdzenia. */
  dolozDoKolejki(
    zadanie: StudioIngestQueueAddRequest,
  ): Promise<Wynik<StudioIngestQueueAddResponse>>;
  /** Oddaje kolejkę wczytywania okna wraz ze stanem każdej pozycji. */
  kolejka(zadanie: StudioIngestQueueListRequest): Promise<Wynik<StudioIngestQueueListResponse>>;
  /** Rozpoznaje tekst pozycji z pełnym sterowaniem silnikiem, językami i progiem. */
  rozpoznaj(zadanie: StudioIngestRecognizeRequest): Promise<Wynik<StudioIngestRecognizeResponse>>;
  /** Poprawia rozpoznane słowo PRZED przyjęciem wyniku do edytora. */
  popraw(
    zadanie: StudioIngestCorrectionSetRequest,
  ): Promise<Wynik<StudioIngestCorrectionSetResponse>>;
  /** Przyjmuje wynik kolejki jako dokument roboczy wraz z pierwszą wersją. */
  przyjmij(zadanie: StudioIngestItemAcceptRequest): Promise<Wynik<StudioIngestItemAcceptResponse>>;
  /** Oddaje urządzenia wejściowe widziane przez rdzeń: skanery i kamery. */
  urzadzenia(): Promise<Wynik<StudioIngestDeviceListResponse>>;
}

export function utworzZrodloWstawienStudio(kanal: Kanal): ZrodloWstawienStudio {
  return {
    async zalozDokument(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentCreate, zadanie),
        Command.StudioDocumentCreate,
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.form),
      );
    },

    async wniesPlik(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentImportFile, zadanie),
        Command.StudioDocumentImportFile,
        // Bilans wniesienia jest polem obowiązkowym odpowiedzi, nie polem ozdobnym.
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.balance),
      );
    },

    async wniesPdf(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentImportPdf, zadanie),
        Command.StudioDocumentImportPdf,
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.balance),
      );
    },

    async wniesObraz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentImageImport, zadanie),
        Command.StudioDocumentImageImport,
        (tresc) => czyObiekt(tresc.object) && czyObiekt(tresc.balance),
      );
    },

    async zapiszJako(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentSaveAs, zadanie),
        Command.StudioDocumentSaveAs,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async skopiuj(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentCopy, zadanie),
        Command.StudioDocumentCopy,
        (tresc) => czyObiekt(tresc.document) && czyLiczba(tresc.versionsCopied),
      );
    },

    async wydaj(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentExportFormat, zadanie),
        Command.StudioDocumentExportFormat,
        (tresc) => czyObiekt(tresc.result),
      );
    },

    async wydajWsad(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentExportBatch, zadanie),
        Command.StudioDocumentExportBatch,
        (tresc) =>
          czyTablica(tresc.results) && czyLiczba(tresc.succeeded) && czyLiczba(tresc.failed),
      );
    },

    async przestawFormat(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentFormatSet, zadanie),
        Command.StudioDocumentFormatSet,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async wniesZBiblioteki(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioInsertFromLibrary, zadanie),
        Command.StudioInsertFromLibrary,
        // Pochodzenie fragmentu jest polem obowiązkowym, nie polem opcjonalnym odpowiedzi rdzenia.
        (tresc) => czyObiekt(tresc.provenance) && czyObiekt(tresc.balance),
      );
    },

    async wniesZeStrony(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioInsertFromWeb, zadanie),
        Command.StudioInsertFromWeb,
        (tresc) => czyObiekt(tresc.provenance) && czyObiekt(tresc.balance),
      );
    },

    async dolozDoKolejki(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioIngestQueueAdd, zadanie),
        Command.StudioIngestQueueAdd,
        (tresc) => czyTablica(tresc.items),
      );
    },

    async kolejka(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioIngestQueueList, zadanie),
        Command.StudioIngestQueueList,
        (tresc) => czyTablica(tresc.items),
      );
    },

    async rozpoznaj(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioIngestRecognize, zadanie),
        Command.StudioIngestRecognize,
        // Słowa i układ rozpoznania są nieobowiązkowe, więc ich brak jest wynikiem poprawnym.
        (tresc) => czyObiekt(tresc.item),
      );
    },

    async popraw(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioIngestCorrectionSet, zadanie),
        Command.StudioIngestCorrectionSet,
        (tresc) => czyObiekt(tresc.item),
      );
    },

    async przyjmij(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioIngestItemAccept, zadanie),
        Command.StudioIngestItemAccept,
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.version),
      );
    },

    async urzadzenia() {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioIngestDeviceList, {}),
        Command.StudioIngestDeviceList,
        (tresc) => czyTablica(tresc.devices),
      );
    },
  };
}

/* ── Bilans widoczny — trzy zdania, jedno miejsce ─────────────────────────── */

/** Zdanie o bilansie czynności podaje liczbę zmienionych i pominiętych miejsc wraz z nazwą blokady, która pominięcie spowodowała. */
export function wstawieniaOpiszBilans(bilans: StudioActionBalance): string {
  const czesci: string[] = [
    bilans.applied === 0
      ? 'Czynność nie zmieniła ani jednego miejsca'
      : `Zmienionych miejsc: ${bilans.applied}`,
  ];
  if (bilans.skippedCount > 0) czesci.push(`pominiętych: ${bilans.skippedCount}`);
  if (bilans.deferredCount !== undefined && bilans.deferredCount > 0) {
    czesci.push(`odłożonych, żeby nie nadpisać cudzej pracy: ${bilans.deferredCount}`);
  }
  if (bilans.conflicts !== undefined && bilans.conflicts.length > 0) {
    czesci.push(`spięć o ten sam fragment: ${bilans.conflicts.length}`);
  }
  const zdania = [`${czesci.join(', ')}.`];
  if (bilans.note !== undefined && bilans.note !== '') zdania.push(bilans.note);
  const pominiecia = wstawieniaOpiszPominiecia(bilans.skipped);
  if (pominiecia !== '') zdania.push(pominiecia);
  return zdania.join(' ');
}

/** Zdanie o bilansie wniesienia podaje, ile stron, akapitów, tabel, obrazów, stylów i sekcji rdzeń odzyskał z wniesionego materiału. */
export function wstawieniaOpiszWniesienie(bilans: StudioImportBalance): string {
  const czesci: string[] = [`Format wniesiony: ${bilans.format}`];
  if (bilans.encoding !== undefined && bilans.encoding !== '') {
    czesci.push(`zapis znaków rozpoznany jako ${bilans.encoding}`);
  }
  if (bilans.pages !== undefined) czesci.push(`stron ${bilans.pages}`);
  if (bilans.pagesWithText !== undefined) czesci.push(`z warstwą tekstową ${bilans.pagesWithText}`);
  if (bilans.pagesWithoutText !== undefined && bilans.pagesWithoutText > 0) {
    czesci.push(`BEZ warstwy tekstowej ${bilans.pagesWithoutText}`);
  }
  if (bilans.paragraphsRecovered !== undefined) {
    czesci.push(`akapitów odtworzonych ${bilans.paragraphsRecovered}`);
  }
  if (bilans.tablesRecognized !== undefined) {
    czesci.push(`tabel rozpoznanych ${bilans.tablesRecognized}`);
  }
  if (bilans.tablesMissed !== undefined && bilans.tablesMissed > 0) {
    czesci.push(`układów tabelarycznych NIEROZPOZNANYCH ${bilans.tablesMissed}`);
  }
  if (bilans.imagesEmbedded !== undefined) {
    czesci.push(`obrazów osadzonych ${bilans.imagesEmbedded}`);
  }
  if (bilans.imagesSkipped !== undefined && bilans.imagesSkipped > 0) {
    czesci.push(`obrazów POMINIĘTYCH ${bilans.imagesSkipped}`);
  }
  if (bilans.stylesRecovered !== undefined) {
    czesci.push(`stylów nazwanych przejętych ${bilans.stylesRecovered}`);
  }
  if (bilans.sectionsRecovered !== undefined) {
    czesci.push(`sekcji przejętych ${bilans.sectionsRecovered}`);
  }

  const zdania = [`${czesci.join(' · ')}.`];
  if (bilans.needsTextRecognition === true) {
    zdania.push(
      'Ten materiał NIE MA warstwy tekstowej — odzyskania nie ma z czego zrobić i rdzeń kieruje ' +
        'go na rozpoznanie tekstu. Otwórz narzędziownię cyfryzacji, rozpoznaj pozycję i przyjmij ' +
        'wynik; wniesienie udające konwersję oddałoby pustą kartkę.',
    );
  }
  if (bilans.note !== undefined && bilans.note !== '') zdania.push(bilans.note);
  const pominiecia = wstawieniaOpiszPominiecia(bilans.skipped);
  if (pominiecia !== '') zdania.push(pominiecia);
  return zdania.join(' ');
}

/** Zdanie o wyniku wydania dokumentu podaje ścieżkę albo zasób zapisu wraz z wykazem cech, których docelowy format nie niesie. */
export function wstawieniaOpiszWydanie(wynik: StudioExportResult): string {
  const czesci: string[] = [`Wydanie do formatu ${wynik.format}`];
  if (wynik.path !== undefined && wynik.path !== '') czesci.push(`plik ${wynik.path}`);
  if (wynik.assetId !== undefined && wynik.assetId !== '') {
    czesci.push(`zasób magazynu rdzenia ${wynik.assetId}`);
  }
  if (wynik.bytes !== undefined) czesci.push(`${wynik.bytes} bajtów`);

  const zdania = [`${czesci.join(' · ')}.`];
  const pominiete = wynik.droppedFeatures ?? [];
  if (pominiete.length === 0) {
    zdania.push('Rdzeń nie zgłosił ani jednej cechy dokumentu, której ten format nie niesie.');
  } else {
    zdania.push(
      `Ten format NIE NIESIE ${pominiete.length} cech dokumentu i one z wydania odpadły: ` +
        `${pominiete.map(opiszPominiecie).join('; ')}.`,
    );
  }
  if (wynik.note !== undefined && wynik.note !== '') zdania.push(wynik.note);
  return zdania.join(' ');
}

/** Wykaz pominięć jednym zdaniem łączy powód, szczegół i blokadę każdej pominiętej pozycji, a pusty wykaz nie daje zdania wcale. */
export function wstawieniaOpiszPominiecia(
  pominiecia: readonly StudioSkippedItem[] | undefined,
): string {
  const wykaz = pominiecia ?? [];
  if (wykaz.length === 0) return '';
  return `Pominięte: ${wykaz.map(opiszPominiecie).join('; ')}.`;
}

/** Jedno pominięcie niesie powód, czego dotyczy, oraz nazwę albo identyfikator blokady, jeśli to ona je zatrzymała. */
function opiszPominiecie(pozycja: StudioSkippedItem): string {
  const czesci: string[] = [pozycja.reason];
  if (pozycja.detail !== undefined && pozycja.detail !== '') czesci.push(pozycja.detail);
  if (pozycja.lockName !== undefined && pozycja.lockName !== '') {
    czesci.push(`zatrzymała to blokada „${pozycja.lockName}"`);
  } else if (pozycja.lockId !== undefined && pozycja.lockId !== '') {
    czesci.push(`zatrzymała to blokada ${pozycja.lockId}`);
  }
  if (pozycja.rangeStart !== undefined && pozycja.rangeEnd !== undefined) {
    czesci.push(`znaki ${pozycja.rangeStart}–${pozycja.rangeEnd}`);
  }
  return czesci.join(', ');
}
