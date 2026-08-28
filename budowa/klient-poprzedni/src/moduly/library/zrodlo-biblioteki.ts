import {
  Command,
  EventType,
  type LibraryFile,
  type LibraryFileChangedEvent,
  type LibraryFileListRequest,
  type LibraryFileUploadRequest,
  type LibraryPreview,
  type LibraryVersion,
  type LibraryVersionAddRequest,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import type { StrazOdmow } from './straz-odmow';
import { utworzWywolaniaWersji } from './zrodlo-wersji';
import { utworzZarzadRepozytorium, type ZarzadRepozytorium } from './zrodlo-zarzadu';

/**
 * Interfejs udostępnia komendy obszaru biblioteki widziane przez okna modułu:
 * źródło niczego nie pamięta, jest warstwą wywołania i sprawdzianu kształtu
 * odpowiedzi, a odmowa rdzenia wraca jako wynik z błędem.
 */
export interface ZrodloBiblioteki extends ZarzadRepozytorium {
  wykaz(zadanie: LibraryFileListRequest): Promise<Wynik<{ files: LibraryFile[]; total?: number }>>;
  szukaj(fraza: string, granica: number): Promise<Wynik<{ files: LibraryFile[]; total: number }>>;
  wgraj(zadanie: LibraryFileUploadRequest): Promise<Wynik<{ file: LibraryFile }>>;
  podglad(idPliku: string, strona: number): Promise<Wynik<{ preview: LibraryPreview }>>;
  wersje(idPliku: string): Promise<Wynik<{ versions: LibraryVersion[] }>>;
  /** Dokłada wersję istniejącemu plikowi — jedyna droga wzrostu jego historii. */
  dolozWersje(
    zadanie: LibraryVersionAddRequest,
  ): Promise<Wynik<{ version: LibraryVersion; file: LibraryFile }>>;
  przywroc(idPliku: string, idWersji: string): Promise<Wynik<{ file: LibraryFile }>>;
  nadajEtykiety(
    idPliku: string,
    etykiety: readonly string[],
    kolekcje?: readonly string[],
  ): Promise<Wynik<{ file: LibraryFile }>>;
  zalozKolekcje(
    nazwa: string,
    opis: string,
  ): Promise<Wynik<{ collectionId: string; name: string }>>;
  przypiszDoKolekcji(
    idKolekcji: string,
    idPlikow: readonly string[],
  ): Promise<Wynik<{ collectionId: string; assignedCount: number }>>;
  /** Subskrypcja `library.file.changed` — także artefaktów z innych modułów. */
  naZmianePliku(sluchacz: (tresc: LibraryFileChangedEvent) => void): Odsubskrybuj;
}

/** Górna granica podglądu tekstowego — okno pokazuje fragment, nie cały plik, chroniąc wydajność wyświetlania. */
const ZNAKI_PODGLADU = 4000;

export function utworzZrodloBiblioteki(kanal: Kanal, straz: StrazOdmow): ZrodloBiblioteki {
  return {
    // Historia dokumentu ma własny plik: trzy komendy wersjonowania są jednym łańcuchem.
    ...utworzWywolaniaWersji(straz),
    // Zarząd repozytorium ma własny plik: jeden kształt, trzy pliki wedle odpowiedzialności.
    ...utworzZarzadRepozytorium(straz),

    async wykaz(zadanie) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileList, zadanie),
        Command.LibraryFileList,
        (tresc) => czyTablica(tresc.files),
      );
    },

    async szukaj(fraza, granica) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileSearch, { query: fraza, limit: granica }),
        Command.LibraryFileSearch,
        (tresc) => czyTablica(tresc.files),
      );
    },

    async wgraj(zadanie) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFileUpload, zadanie),
        Command.LibraryFileUpload,
        (tresc) => czyObiekt(tresc.file),
      );
    },

    async podglad(idPliku, strona) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryFilePreview, {
          fileId: idPliku,
          page: strona,
          maxChars: ZNAKI_PODGLADU,
        }),
        Command.LibraryFilePreview,
        (tresc) => czyObiekt(tresc.preview),
      );
    },

    async nadajEtykiety(idPliku, etykiety, kolekcje) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryTagSet, {
          fileId: idPliku,
          tags: [...etykiety],
          ...(kolekcje === undefined ? {} : { collectionIds: [...kolekcje] }),
        }),
        Command.LibraryTagSet,
        (tresc) => czyObiekt(tresc.file),
      );
    },

    async zalozKolekcje(nazwa, opis) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryCollectionCreate, {
          name: nazwa,
          ...(opis === '' ? {} : { description: opis }),
        }),
        Command.LibraryCollectionCreate,
        // Nazwa wchodzi do sprawdzianu, bo okno buduje z niej zdanie potwierdzające zapis kolekcji.
        (tresc) => typeof tresc.collectionId === 'string' && typeof tresc.name === 'string',
      );
    },

    async przypiszDoKolekcji(idKolekcji, idPlikow) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryCollectionAssign, {
          collectionId: idKolekcji,
          fileIds: [...idPlikow],
        }),
        Command.LibraryCollectionAssign,
        (tresc) => typeof tresc.assignedCount === 'number',
      );
    },

    naZmianePliku(sluchacz) {
      return kanal.naZdarzenie(EventType.LibraryFileChanged, (tresc) => sluchacz(tresc));
    },
  };
}
