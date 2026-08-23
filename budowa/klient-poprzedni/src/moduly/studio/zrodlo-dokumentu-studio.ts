import {
  Command,
  type DocumentConvertRequest,
  type DocumentConvertResponse,
  type DocumentTextExtractRequest,
  type DocumentTextExtractResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Dwie komendy obszaru dokumentów, z których korzysta moduł Studio — cudzy
 * obszar, którego moduł nie prowadzi.
 *
 * Obszar `studio` nie ma komendy przyjmującej plik ani wydającej dokument
 * w formacie wyjściowym: `studio.document.open` ze wskazaniem ścieżki albo pliku
 * Library zakłada dokument PUSTY, a `studio.document.save` przyjmuje sam napis.
 * Cyfryzacja materiału i zamiana formatu są w kontrakcie, ale w obszarze
 * `document` — i to one są jedynym wejściem modułu od strony pliku oraz jedynym
 * wyjściem do formatu binarnego. Moduł ich nie kopiuje: woła je wprost, tak samo
 * jak woła `window.list` i `action.list`.
 *
 * Wywołanie idzie zwykłą drogą protokołu, nie osłoną `wywolajUczciwie`
 * z `odmowa-rdzenia.ts`: ta rozpoznaje koperty `studio.unknown` i `window.unknown`,
 * a obszar `document` własnego zdarzenia nierozpoznanej nie ma — kontrakt kieruje
 * go na zdarzenie połączenia, więc odmowa przychodzi kopertą ze statusem
 * i korelacja rozpoznaje ją bez pomocy.
 *
 * Ścieżka pliku jest ścieżką po stronie rdzenia. Klient dysku nie czyta i nie
 * zapisuje — podaje wskazanie i oddaje Operatorowi to, co rdzeń odpowiedział.
 */

/** Wskazanie materiału wejściowego — jedna z dwóch dróg do treści pliku. */
export interface WskazanieMaterialu {
  /** Ścieżka pliku widziana przez rdzeń. */
  sciezka: string;
  /** Zasób magazynu rdzenia (Design albo Library). */
  idZasobu: string;
  /** Język rozpoznawania pisma; pusty zostawia rdzeniowi wartość domyślną. */
  jezyk: string;
  /** Pierwsza strona zakresu; 0 znaczy „bez ograniczenia". */
  stronaOd: number;
  /** Ostatnia strona zakresu; 0 znaczy „bez ograniczenia". */
  stronaDo: number;
  /** Wymuszenie rozpoznania pisma mimo obecnej warstwy tekstowej. */
  wymusRozpoznanie: boolean;
}

export interface ZrodloDokumentuStudio {
  /** Wydobywa tekst z dokumentu albo obrazu wraz z rozpoznaniem pisma. */
  wydobadzTekst(
    wskazanie: WskazanieMaterialu,
  ): Promise<Wynik<DocumentTextExtractResponse>>;
  /**
   * Zamienia treść dokumentu na format docelowy i odkłada wynik w magazynie.
   *
   * Treść jedzie polem `content`, a nie ścieżką: dokument Studia mieszka
   * w rdzeniu pod `studio.document.save`, a nie na dysku Operatora, więc
   * wskazanie ścieżki oddawałoby do zamiany plik, którego treść mogła się
   * już rozejść z treścią zaakceptowaną w module.
   */
  zamienFormat(
    idOkna: string,
    tresc: string,
    formatZrodlowy: string,
    formatDocelowy: string,
  ): Promise<Wynik<DocumentConvertResponse>>;
}

export function utworzZrodloDokumentuStudio(kanal: Kanal): ZrodloDokumentuStudio {
  return {
    async wydobadzTekst(wskazanie) {
      const zadanie: DocumentTextExtractRequest = {};
      // Zasób magazynu ma pierwszeństwo przed ścieżką: jego treść stoi już pod
      // sumą kontrolną, więc rdzeń nie wciąga jej po raz drugi.
      if (wskazanie.idZasobu !== '') zadanie.assetId = wskazanie.idZasobu;
      else zadanie.sourcePath = wskazanie.sciezka;
      if (wskazanie.jezyk !== '') zadanie.language = wskazanie.jezyk;
      if (wskazanie.stronaOd > 0) zadanie.pageFrom = wskazanie.stronaOd;
      if (wskazanie.stronaDo > 0) zadanie.pageTo = wskazanie.stronaDo;
      // Pole wymuszenia wchodzi wyłącznie po wskazaniu Operatora: brak pola
      // znaczy „weź warstwę tekstową, gdy dokument ją ma", czyli zachowanie
      // domyślne kontraktu, a nie wybór okna.
      if (wskazanie.wymusRozpoznanie) zadanie.forceOcr = true;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DocumentTextExtract, zadanie),
        Command.DocumentTextExtract,
        (tresc) => czyTekst(tresc.text) && czyLogiczna(tresc.usedOcr),
      );
    },

    async zamienFormat(idOkna, tresc, formatZrodlowy, formatDocelowy) {
      const zadanie: DocumentConvertRequest = {
        content: tresc,
        toFormat: formatDocelowy,
      };
      // Okno wchodzi do żądania, bo od niego zależy, czy zasób wyniku pojawi się
      // w wykazie okna, czy zostanie w magazynie bez wykazu. Brak okna czynności
      // nie wstrzymuje — kontrakt mówi to wprost przy polu.
      if (idOkna !== '') zadanie.windowId = idOkna;
      // Format źródłowy podaje się jawnie, bo rozpoznanie z pliku nie ma tu na
      // czym pracować: treść jedzie napisem, bez nazwy i bez rozszerzenia.
      if (formatZrodlowy !== '') zadanie.fromFormat = formatZrodlowy;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DocumentConvert, zadanie),
        Command.DocumentConvert,
        (odpowiedz) => czyObiekt(odpowiedz.asset),
      );
    },
  };
}
