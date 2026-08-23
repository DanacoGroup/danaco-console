import { Command } from '../../../../shared/contract';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import type { StrazOdmow } from './straz-odmow';
import type { ZrodloBiblioteki } from './zrodlo-biblioteki';

/**
 * Trzy wywołania historii dokumentu: odczyt wersji, dołożenie kolejnej
 * i powrót do wcześniejszej.
 *
 * Wydzielone ze `zrodlo-biblioteki.ts` wzdłuż odpowiedzialności: tutaj historia
 * dokumentu, tam plik, wyszukiwanie i kolekcje. Kształt pozostaje jeden —
 * `ZrodloBiblioteki` — żeby okna widziały jedno źródło, a nie dwa. Komendy stoją
 * razem, bo są jednym łańcuchem: bez dołożenia wersji odczyt historii zwraca
 * jeden wpis, a przywrócenie nie ma dokąd wracać.
 */
export type WywolaniaWersji = Pick<ZrodloBiblioteki, 'wersje' | 'dolozWersje' | 'przywroc'>;

export function utworzWywolaniaWersji(straz: StrazOdmow): WywolaniaWersji {
  return {
    async wersje(idPliku) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryVersionList, { fileId: idPliku }),
        Command.LibraryVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
    },

    async dolozWersje(zadanie) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryVersionAdd, zadanie),
        Command.LibraryVersionAdd,
        (tresc) => czyObiekt(tresc.version) && czyObiekt(tresc.file),
      );
    },

    async przywroc(idPliku, idWersji) {
      return sprawdzKsztalt(
        await straz.wywolaj(Command.LibraryVersionRestore, {
          fileId: idPliku,
          versionId: idWersji,
        }),
        Command.LibraryVersionRestore,
        (tresc) => czyObiekt(tresc.file),
      );
    },
  };
}
