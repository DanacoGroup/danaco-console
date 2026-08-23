import {
  Command,
  type LibraryFilePreviewRequest,
  type LibraryPreview,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzNasluchOdmow, type NasluchOdmow, type OdpowiedzBadania } from './nasluch-odmow';

/**
 * Treść źródła wczytywana do Reading View.
 *
 * Drogi są dwie, ale przejezdna jest jedna. Kontrakt ma własną komendę modułu —
 * odczyt treści źródła niezależnie od tego, czy jest ono plikiem repozytorium —
 * a rdzeń nie ma dla niej jeszcze uchwytu. Przejezdny jest `library.file.preview`,
 * czyli podgląd zasobu repozytorium wraz z podziałem na strony i znacznikiem
 * skrócenia, i po nim to źródło sięga.
 *
 * Ta droga wystarcza dla źródeł związanych z dokumentem repozytorium (pole
 * `libraryFileId`) i tylko dla nich — źródło typu strona internetowa, notatka
 * albo zbiór danych, którego rdzeń nie trzyma jako pliku, treści nie ma dziś
 * skąd wziąć. Okno nazywa to Operatorowi zamiast pokazywać pusty czytnik;
 * dobudowa uchwytu własnej komendy tę granicę zdejmuje.
 *
 * Mechanizm podglądu jest wspólny z File Preview modułu Library — tak, jak
 * rozstrzyga opracowanie modułu (rozdz. 7.4: „podgląd współdzielony z Library
 * File Preview"). Research nie buduje drugiego czytnika obok tamtego.
 */
export interface ZrodloLektury {
  /** Podgląd strony dokumentu repozytorium wraz z jego treścią tekstową. */
  wczytajStrone(
    idPliku: string,
    strona: number,
    limitZnakow: number,
  ): Promise<OdpowiedzBadania<LibraryPreview>>;
  rozlacz(): void;
}

export function utworzZrodloLektury(kanal: Kanal): ZrodloLektury {
  const nasluch: NasluchOdmow = utworzNasluchOdmow(kanal);

  return {
    async wczytajStrone(idPliku, strona, limitZnakow) {
      const zadanie: LibraryFilePreviewRequest = {
        fileId: idPliku,
        page: strona,
        maxChars: limitZnakow,
      };
      const odpowiedz = await nasluch.wyslij(Command.LibraryFilePreview, zadanie);
      const wynik = sprawdzKsztalt(odpowiedz, Command.LibraryFilePreview, (tresc) =>
        czyObiekt(tresc.preview),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        // Nazwa nieznanego typu przechodzi dalej, żeby okno mogło ją wypisać
        // wprost — ten sam zabieg co w `zrodlo-research.ts`.
        return {
          udany: false,
          ...(wynik.blad === undefined ? {} : { blad: wynik.blad }),
          ...(odpowiedz.nieznanyTyp === undefined ? {} : { nieznanyTyp: odpowiedz.nieznanyTyp }),
        };
      }
      return { udany: true, wynik: wynik.wynik.preview };
    },

    rozlacz: () => nasluch.rozlacz(),
  };
}
