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
 * Wskazanie materiału wejściowego — jedna z dwóch dróg do treści pliku obszaru dokumentów,
 * z którego korzysta moduł Studio.
 */
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
  /** Zamienia treść na format docelowy, wynik zapisuje w magazynie; jedzie polem `content`. */
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
      // Zasób magazynu ma pierwszeństwo przed ścieżką: jego treść stoi już pod sumą kontrolną.
      if (wskazanie.idZasobu !== '') zadanie.assetId = wskazanie.idZasobu;
      else zadanie.sourcePath = wskazanie.sciezka;
      if (wskazanie.jezyk !== '') zadanie.language = wskazanie.jezyk;
      if (wskazanie.stronaOd > 0) zadanie.pageFrom = wskazanie.stronaOd;
      if (wskazanie.stronaDo > 0) zadanie.pageTo = wskazanie.stronaDo;
      // Pole wymuszenia wchodzi tylko po wskazaniu Operatora; brak znaczy zachowanie domyślne.
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
      // Okno wchodzi do żądania, bo od niego zależy widoczność wyniku w wykazie, nie w magazynie.
      if (idOkna !== '') zadanie.windowId = idOkna;
      // Format źródłowy podaje się jawnie: treść jedzie napisem, bez nazwy i rozszerzenia pliku.
      if (formatZrodlowy !== '') zadanie.fromFormat = formatZrodlowy;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DocumentConvert, zadanie),
        Command.DocumentConvert,
        (odpowiedz) => czyObiekt(odpowiedz.asset),
      );
    },
  };
}
