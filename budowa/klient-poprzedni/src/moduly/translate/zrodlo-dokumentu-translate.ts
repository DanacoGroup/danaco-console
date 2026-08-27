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
 * Dwie komendy obszaru dokumentów wołane wprost przez moduł: wydobycie tekstu
 * z dokumentu albo skanu i zamiana formatu dokumentu po stronie rdzenia.
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
      // Pole wymuszenia wchodzi wyłącznie wtedy, gdy Operator je zaznaczył w oknie.
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
