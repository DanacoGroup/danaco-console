import {
  Command,
  type BrowserNavigateResponse,
  type BrowserSnapshotGetResponse,
  type LibraryFileListResponse,
  type LibraryFilePreviewResponse,
  type StudioIngestUrlResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Przeglądarka i Biblioteka wewnątrz Studia: sześć komend wołanych z gotowych rodzin, bez rachunku przeglądania ani repozytorium plików. */
export interface OsadzenieZrodel {
  /** Pliki Biblioteki wedle frazy i etykiet. */
  osadzenieBiblioteka(fraza: string, granica: number): Promise<Wynik<LibraryFileListResponse>>;
  /** Podgląd pliku Biblioteki wraz z treścią tekstową. */
  osadzeniePodglad(
    idPliku: string,
    strona: number,
    maxZnakow: number,
  ): Promise<Wynik<LibraryFilePreviewResponse>>;
  /** Wciąga stronę oknem Studia i oczyszcza ją do postaci czytelnej. */
  osadzenieWciagnijStrone(
    idOkna: string,
    adres: string,
    zObrazami: boolean,
  ): Promise<Wynik<StudioIngestUrlResponse>>;
  /** Migawka strony otwartej w module Browser. */
  osadzenieMigawka(
    idOknaPrzegladarki: string,
    zeZrodlem: boolean,
  ): Promise<Wynik<BrowserSnapshotGetResponse>>;
  /** Otwiera adres w oknie modułu Browser i oddaje migawkę po przejściu. */
  osadzenieOtworz(
    idOknaPrzegladarki: string,
    adres: string,
  ): Promise<Wynik<BrowserNavigateResponse>>;
}

export function utworzOsadzenieZrodel(kanal: Kanal): OsadzenieZrodel {
  return {
    async osadzenieBiblioteka(fraza, granica) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.LibraryFileList, {
          ...(fraza === '' ? {} : { query: fraza }),
          limit: granica,
        }),
        Command.LibraryFileList,
        (tresc) => czyTablica(tresc.files),
      );
    },

    async osadzeniePodglad(idPliku, strona, maxZnakow) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.LibraryFilePreview, {
          fileId: idPliku,
          page: strona,
          maxChars: maxZnakow,
        }),
        Command.LibraryFilePreview,
        (tresc) => czyObiekt(tresc.preview),
      );
    },

    async osadzenieWciagnijStrone(idOkna, adres, zObrazami) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.StudioIngestUrl, {
          windowId: idOkna,
          url: adres,
          readability: true,
          includeImages: zObrazami,
        }),
        Command.StudioIngestUrl,
        (tresc) => czyObiekt(tresc.item),
      );
    },

    async osadzenieMigawka(idOknaPrzegladarki, zeZrodlem) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSnapshotGet, {
          windowId: idOknaPrzegladarki,
          includeHtml: zeZrodlem,
        }),
        Command.BrowserSnapshotGet,
        (tresc) => czyObiekt(tresc.snapshot),
      );
    },

    async osadzenieOtworz(idOknaPrzegladarki, adres) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNavigate, {
          windowId: idOknaPrzegladarki,
          url: adres,
        }),
        Command.BrowserNavigate,
        (tresc) => czyObiekt(tresc.snapshot),
      );
    },
  };
}
