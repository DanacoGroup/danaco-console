import {
  Command,
  KnowledgeScope,
  type KnowledgeSearchRequest,
  type LibraryFileSearchRequest,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzNasluchOdmow, type NasluchOdmow, type OdpowiedzBadania } from './nasluch-odmow';
import { zHitowWiedzy, zPlikowRepozytorium, type WynikOdkrycia } from './wynik-odkrycia';

/**
 * Wyszukiwanie zasilające Discovery Panel: wyszukiwanie po znaczeniu w wiedzy Operatora oraz
 * pełnotekstowe w zasobach repozytorium, dwa tryby z czterech kontraktu.
 */
export interface ZrodloOdkrywania {
  /** Wyszukiwanie po znaczeniu w wiedzy Operatora i w przestrzeni okna badania. */
  poZnaczeniu(
    zapytanie: string,
    idOkna: string,
    limit: number,
  ): Promise<OdpowiedzBadania<WynikOdkrycia[]>>;
  /** Wyszukiwanie pełnotekstowe w zasobach repozytorium Library. */
  poTresci(zapytanie: string, limit: number): Promise<OdpowiedzBadania<WynikOdkrycia[]>>;
  rozlacz(): void;
}

export function utworzZrodloOdkrywania(kanal: Kanal): ZrodloOdkrywania {
  const nasluch: NasluchOdmow = utworzNasluchOdmow(kanal);

  return {
    async poZnaczeniu(zapytanie, idOkna, limit) {
      // Zakres all obejmuje bibliotekę, rozmowy i pliki przestrzeni; windowId wskazuje czyją przeszukać.
      const zadanie: KnowledgeSearchRequest = {
        query: zapytanie,
        scope: KnowledgeScope.All,
        limit,
      };
      if (idOkna !== '') zadanie.windowId = idOkna;
      const wynik = sprawdz(
        await nasluch.wyslij(Command.KnowledgeSearch, zadanie),
        Command.KnowledgeSearch,
        (tresc) => czyTablica(tresc.results),
      );
      if (!wynik.udany || wynik.wynik === undefined) return przenies(wynik);
      return { udany: true, wynik: zHitowWiedzy(wynik.wynik.results) };
    },

    async poTresci(zapytanie, limit) {
      const zadanie: LibraryFileSearchRequest = { query: zapytanie, limit };
      const wynik = sprawdz(
        await nasluch.wyslij(Command.LibraryFileSearch, zadanie),
        Command.LibraryFileSearch,
        (tresc) => czyTablica(tresc.files),
      );
      if (!wynik.udany || wynik.wynik === undefined) return przenies(wynik);
      return { udany: true, wynik: zPlikowRepozytorium(wynik.wynik.files) };
    },

    rozlacz: () => nasluch.rozlacz(),
  };
}

/**
 * Sprawdzian kształtu zachowujący nazwę nieznanego typu — ten sam zabieg co
 * w `zrodlo-research.ts`. Bez niego odmowa nieznanej komendy dochodziłaby do
 * okna jako zwykły błąd i okno nie miałoby czego wypisać.
 */
function sprawdz<T>(
  odpowiedz: OdpowiedzBadania<T>,
  komenda: string,
  sprawdzian: (tresc: T) => boolean,
): OdpowiedzBadania<T> {
  const wynik = sprawdzKsztalt(odpowiedz, komenda, sprawdzian);
  if (wynik.udany || odpowiedz.nieznanyTyp === undefined) return wynik;
  return { ...wynik, nieznanyTyp: odpowiedz.nieznanyTyp };
}

/** Przeniesienie niepowodzenia z jednego wyniku na wynik o innej treści, bez gubienia powodu odmowy ani nazwy nieznanego typu. */
function przenies<T>(odpowiedz: OdpowiedzBadania<T>): OdpowiedzBadania<WynikOdkrycia[]> {
  return {
    udany: false,
    ...(odpowiedz.blad === undefined ? {} : { blad: odpowiedz.blad }),
    ...(odpowiedz.nieznanyTyp === undefined ? {} : { nieznanyTyp: odpowiedz.nieznanyTyp }),
  };
}
