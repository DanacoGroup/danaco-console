import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Kanal } from '../../protokol/kanal';
import { wywolaj } from '../../protokol/wywolanie';
import { OdmowaDostarczenia, type Dostarczenie } from './dostarczenie-nagrania';
import { wynikNieprzetworzony, zlozWynik, type WynikDyktowania } from './wynik-dyktowania';

/**
 * Transkrypcja nagrania — dostarczenie, wywołanie silnika, złożenie wyniku.
 *
 * Ogniwo dostarcza nagranie na maszynę silnika, prosi o `speech.transcribe`
 * i składa wynik według trzech stanów. Nie ma tu logiki nagrywania (to sąsiedni
 * plik pakietu) ani logiki widoku — oddaje `WynikDyktowania` i kończy pracę.
 *
 * Zatrzymanie na dostarczeniu i odmowa transkrypcji dają ten sam stan
 * `nieprzetworzone`, bo skutek dla Operatora jest jeden: tekstu nie ma. Powód
 * jest jednak inny i inna jest naprawa — w pierwszym przypadku brakuje komendy
 * w kontrakcie, w drugim rdzeń miał nagranie i go nie przerobił. Oba zdania są
 * więc rozdzielone.
 *
 * Obietnica nie jest odrzucana: odmowa wraca jako wynik ze stanem
 * `nieprzetworzone`, nie jako wyjątek — widok nie zakłada `try`, a nieudane
 * dyktowanie nie wywraca okna.
 *
 * Dźwięk nie opuszcza maszyny Operatora: nagranie idzie wyłącznie przez
 * `dostarczenie`, a to ogniwo pilnuje tej granicy samo.
 */
export interface Transkrypcja {
  /** Nagranie → tekst; `jezyk` pusty znaczy rozpoznanie automatyczne. */
  wykonaj(nagranie: Blob, rodzajTresci: string, jezyk?: string): Promise<WynikDyktowania>;
}

export function utworzTranskrypcje(kanal: Kanal, dostarczenie: Dostarczenie): Transkrypcja {
  return {
    async wykonaj(nagranie: Blob, rodzajTresci: string, jezyk?: string): Promise<WynikDyktowania> {
      // Dostarczenie — jedyne miejsce, w którym bajty nagrania stają się
      // `audioRef`. Odmawia, dopóki kontrakt nie ma komendy przyjmującej dźwięk;
      // powód niesie `dostarczenie-nagrania.ts`.
      let audioRef: string;
      try {
        audioRef = await dostarczenie.dostarcz(nagranie, rodzajTresci);
      } catch (powod) {
        return wynikNieprzetworzony(zdanieOdmowyDostarczenia(powod));
      }

      // Silnik. `language` puste zostawiamy nieustawione, bo kontrakt
      // czyta brak pola jako „rozpoznaj automatycznie"; pusty łańcuch byłby
      // wskazaniem języka o pustej nazwie.
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

      // Złożenie. Rozdział trzech stanów należy do `wynik-dyktowania.ts`
      // i tam zostaje; tutaj nie ma ani jednego `if` o pustym tekście, żeby cisza
      // nie zamieniła się po drodze w awarię.
      return zlozWynik(odpowiedz.wynik);
    },
  };
}

/**
 * Zdanie odmowy dostarczenia — powód przenoszony bez straty.
 *
 * `OdmowaDostarczenia` niesie zdanie gotowe dla Operatora i idzie dalej co do
 * znaku. Wyjątek innego rodzaju (pomyłka w kodzie, brak `arrayBuffer` w danym
 * środowisku) też nie ginie: zostaje nazwany jako zatrzymanie na dostarczeniu,
 * żeby Operator wiedział, że rdzeń nagrania jeszcze nie widział.
 */
function zdanieOdmowyDostarczenia(powod: unknown): string {
  if (powod instanceof OdmowaDostarczenia) return powod.zdanie;
  const tresc = powod instanceof Error ? powod.message : String(powod);
  return `Dostarczenie nagrania do silnika zatrzymało się przed wysłaniem: ${tresc}. Rdzeń nagrania nie otrzymał, więc transkrypcja nie została nawet rozpoczęta.`;
}
