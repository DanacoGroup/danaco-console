import { Command, type SessionDeleteResponse } from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { przenies } from '../protokol/wynik-czastkowy';
import { wywolaj } from '../protokol/wywolanie';

// Trwałe usunięcie sesji — jedyna droga utraty danych sesji w produkcie; widoku ten plik nie zna.

/** Odpowiedź rdzenia rozłożona na dwa rozłączne wykazy: sesje usunięte oraz sesje całkiem nieznalezione. */
export interface RozliczenieUsuniecia {
  /** Sesje, których zapis rdzeń skasował. */
  usuniete: readonly string[];
  /** Wskazania bez odpowiednika w historii; nie są błędem. */
  nieznalezione: readonly string[];
  /** Licznik rdzenia — brany z odpowiedzi, nie liczony z długości wykazu. */
  liczba: number;
}

/** Nazwa czynności dla opisu odmowy — jedna wspólna nazwa dla wszystkich odmów usunięcia sesji rdzenia. */
export const CZYNNOSC_USUNIECIA = 'Trwałe usunięcie sesji';

/** `session.delete` — przenosi wskazane sesje do kosza tego rdzenia, wskazań może być naraz bardzo wiele. */
export function zadajUsuniecieSesji(
  kanal: Kanal,
  idSesji: readonly string[],
  potwierdzone: boolean,
): Promise<Wynik<RozliczenieUsuniecia>> {
  return wywolaj(kanal, Command.SessionDelete, {
    sessionIds: [...idSesji],
    confirm: potwierdzone,
  })
    .then((wynik) =>
      sprawdzKsztalt(
        wynik,
        Command.SessionDelete,
        (tresc) => czyTablica(tresc.deletedIds) && czyLiczba(tresc.deletedCount),
      ),
    )
    .then((wynik) => przenies(wynik, rozlicz));
}

/** Rozkład odpowiedzi rdzenia na dwa rozłączne wykazy sesji usuniętych oraz sesji nieznalezionych wcale. */
function rozlicz(odpowiedz: SessionDeleteResponse): RozliczenieUsuniecia {
  return {
    usuniete: odpowiedz.deletedIds,
    nieznalezione: odpowiedz.missingIds ?? [],
    liczba: odpowiedz.deletedCount,
  };
}
