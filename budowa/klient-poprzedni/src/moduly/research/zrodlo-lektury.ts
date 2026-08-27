import {
  Command,
  type LibraryFilePreviewRequest,
  type LibraryPreview,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzNasluchOdmow, type NasluchOdmow, type OdpowiedzBadania } from './nasluch-odmow';

/**
 * Treść źródła wczytywana do Reading View przez podgląd zasobu repozytorium, jedyną dziś
 * przejezdną drogę komendy modułu odczytu treści.
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
        // Nazwa nieznanego typu przechodzi dalej, by okno ją wypisało — ten sam zabieg co w źródle research.
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
