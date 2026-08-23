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
 * Dwie komendy obszaru dokumentów, z których korzysta Format Studio — cudzy
 * obszar, którego moduł nie prowadzi.
 *
 * Obszar `translate` nie ma ani jednej komendy przyjmującej plik: tekst
 * źródłowy wchodzi do niego napisem. Wydobycie tekstu z dokumentu i rozpoznanie
 * pisma ze skanu są w kontrakcie, ale w obszarze `document`, i to one są
 * jedynym wejściem Format Studio od strony pliku. Moduł ich nie kopiuje i nie
 * przepisuje — woła je wprost, tak jak woła `window.list` i `channel.list`.
 *
 * Wywołanie idzie zwykłą drogą protokołu, nie ścieżką odmowy z
 * `wywolanie-translate.ts`: ta rozpoznaje kopertę `translate.unknown`, właściwą
 * wyłącznie obszarowi tego modułu. Odmowa obszaru dokumentów przychodzi
 * kopertą ze statusem i korelacja rozpoznaje ją bez pomocy.
 *
 * Ścieżka pliku jest ścieżką po stronie rdzenia. Klient dysku nie czyta i nie
 * zapisuje — podaje wskazanie i oddaje Operatorowi to, co rdzeń odpowiedział.
 */
export interface ZrodloDokumentuTranslate {
  /** Wydobywa tekst z dokumentu albo obrazu; `rozpoznajPismo` wymusza rozpoznanie ze skanu. */
  wydobadzTekst(
    sciezka: string,
    rozpoznajPismo: boolean,
  ): Promise<Wynik<DocumentTextExtractResponse>>;
  /** Zamienia dokument między formatami rdzenia. */
  zamienFormat(
    sciezka: string,
    formatDocelowy: string,
  ): Promise<Wynik<DocumentConvertResponse>>;
}

export function utworzZrodloDokumentuTranslate(kanal: Kanal): ZrodloDokumentuTranslate {
  return {
    async wydobadzTekst(sciezka, rozpoznajPismo) {
      const zadanie: DocumentTextExtractRequest = { sourcePath: sciezka };
      // Pole wymuszenia wchodzi wyłącznie wtedy, gdy Operator je zaznaczył:
      // brak pola znaczy „weź warstwę tekstową, gdy dokument ją ma", a to jest
      // zachowanie domyślne kontraktu, nie wybór okna.
      if (rozpoznajPismo) zadanie.forceOcr = true;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DocumentTextExtract, zadanie),
        Command.DocumentTextExtract,
        (tresc) => czyTekst(tresc.text) && czyLogiczna(tresc.usedOcr),
      );
    },

    async zamienFormat(sciezka, formatDocelowy) {
      const zadanie: DocumentConvertRequest = {
        sourcePath: sciezka,
        toFormat: formatDocelowy,
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DocumentConvert, zadanie),
        Command.DocumentConvert,
        (tresc) => czyObiekt(tresc.asset),
      );
    },
  };
}
