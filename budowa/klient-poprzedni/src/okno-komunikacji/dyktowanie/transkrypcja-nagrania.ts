import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Kanal } from '../../protokol/kanal';
import { wywolaj } from '../../protokol/wywolanie';
import { OdmowaDostarczenia, type Dostarczenie } from './dostarczenie-nagrania';
import { wynikNieprzetworzony, zlozWynik, type WynikDyktowania } from './wynik-dyktowania';

/**
 * Transkrypcja nagrania dostarcza plik na maszynę silnika, wywołuje `speech.transcribe` i składa wynik jako jeden z trzech stanów bez logiki widoku.
 */
export interface Transkrypcja {
  /** Nagranie → tekst; `jezyk` pusty znaczy rozpoznanie automatyczne. */
  wykonaj(nagranie: Blob, rodzajTresci: string, jezyk?: string): Promise<WynikDyktowania>;
}

export function utworzTranskrypcje(kanal: Kanal, dostarczenie: Dostarczenie): Transkrypcja {
  return {
    async wykonaj(nagranie: Blob, rodzajTresci: string, jezyk?: string): Promise<WynikDyktowania> {
      // Dostarczenie zamienia bajty nagrania na `audioRef` przez moduł `dostarczenie-nagrania.ts`.
      let audioRef: string;
      try {
        audioRef = await dostarczenie.dostarcz(nagranie, rodzajTresci);
      } catch (powod) {
        return wynikNieprzetworzony(zdanieOdmowyDostarczenia(powod));
      }

      // Puste pole `language` zostaje nieustawione — kontrakt czyta brak pola jako rozpoznanie automatyczne.
      const zadany = (jezyk ?? '').trim();
      const odpowiedz = await wywolaj(kanal, Command.SpeechTranscribe, {
        audioRef,
        ...(zadany === '' ? {} : { language: zadany }),
      });

      if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
        return wynikNieprzetworzony(
          opisOdmowyBledu('Transkrypcja nagrania (speech.transcribe)', odpowiedz.blad),
        );
      }

      // Rozdział trzech stanów wyniku należy do `wynik-dyktowania.ts` — tu cisza nie zamienia się w awarię.
      return zlozWynik(odpowiedz.wynik);
    },
  };
}

/**
 * Zdanie odmowy dostarczenia przenosi powód bez straty, a każdy inny wyjątek zostaje nazwany jako zatrzymanie na etapie dostarczenia nagrania.
 */
function zdanieOdmowyDostarczenia(powod: unknown): string {
  if (powod instanceof OdmowaDostarczenia) return powod.zdanie;
  const tresc = powod instanceof Error ? powod.message : String(powod);
  return `Dostarczenie nagrania do silnika zatrzymało się przed wysłaniem: ${tresc}. Rdzeń nagrania nie otrzymał, więc transkrypcja nie została nawet rozpoczęta.`;
}
