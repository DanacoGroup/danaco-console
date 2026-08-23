import { Command, type SessionDeleteResponse } from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { przenies } from '../protokol/wynik-czastkowy';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Trwałe usunięcie sesji — jedyna droga utraty danych sesji w produkcie.
 *
 * Jedna odpowiedzialność: wysłanie `session.delete` i rozłożenie odpowiedzi
 * rdzenia na dwa rozłączne wykazy. Widoku ten plik nie zna.
 *
 * Usunięcie nie jest zamknięciem ani archiwizacją. Rdzeń rozdziela trzy
 * czynności i tak samo rozdziela je powłoka (`adapter_sesje_usuwanie.go`):
 *  - `session.close` (pas kart, `wpiecie-kart-sesji.ts`) zmienia sam stan
 *    sesji; wiadomości, okna i katalog roboczy zostają nietknięte,
 *  - `session.archive` (`protokol/archiwum-sesji.ts`) wyprowadza sesję
 *    z historii bieżącej i pozwala ją przywrócić — zapis zostaje w całości,
 *  - `session.delete` wyprowadza zapis do kosza rdzenia
 *    (`migracja_097_kosz_sesji.sql`); po terminie kosza rdzeń kasuje go
 *    trwale — a w oknie terminu zapis wraca przez `session.restore`.
 * Dlatego czynność nazywa się „Usuń trwale", nigdy „Zamknij" — nazwa wzięta
 * od sąsiada byłaby kłamstwem o skutku.
 *
 * Potwierdzenie jest warunkiem kontraktu, nie ozdobą widoku. `confirm`
 * przechodzi tędy dokładnie takie, jakie podał wywołujący — powłoka nie
 * dopisuje go sama i nie uprzedza odmowy rdzenia własnym sprawdzeniem.
 * Żądanie bez potwierdzenia wraca z rdzenia jako `validation_failed`
 * i tak ma je zobaczyć Operator.
 *
 * Dwa wykazy znaczą dwie różne rzeczy. `deletedIds` to sesje, których zapis
 * zniknął; `missingIds` to wskazania bez odpowiednika w historii, które nie są
 * błędem — widok nie ma prawa ani ich przemilczeć, ani wliczyć do usuniętych,
 * bo milczenie o nich znaczyłoby „usunąłem", gdy nie było czego usuwać.
 */

/** Odpowiedź rdzenia rozłożona na dwa rozłączne wykazy. */
export interface RozliczenieUsuniecia {
  /** Sesje, których zapis rdzeń skasował. */
  usuniete: readonly string[];
  /** Wskazania bez odpowiednika w historii; nie są błędem. */
  nieznalezione: readonly string[];
  /** Licznik rdzenia — brany z odpowiedzi, nie liczony z długości wykazu. */
  liczba: number;
}

/** Nazwa czynności dla `opisOdmowy` — jedna dla wszystkich odmów usunięcia. */
export const CZYNNOSC_USUNIECIA = 'Trwałe usunięcie sesji';

/**
 * `session.delete` — przenosi wskazane sesje do kosza rdzenia.
 *
 * Wskazań może być wiele, bo „usuń tę" i „usuń zaznaczone" to w historii ten
 * sam gest. Kształt sprawdzamy na `deletedIds` i `deletedCount`: `missingIds`
 * jest w kontrakcie polem opcjonalnym, więc jego brak jest odpowiedzią
 * poprawną, a nie odpowiedzią o złym kształcie.
 */
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

/**
 * Rozkład odpowiedzi rdzenia.
 *
 * Dopełnienie `missingIds` do wykazu pustego jest tu wolne od zarzutu, który
 * ciąży na `?? []` w źródłach: tędy przechodzą wyłącznie odpowiedzi udane
 * (niepowodzenie zatrzymuje `przenies` wraz z powodem), a brak pola
 * opcjonalnego w odpowiedzi udanej znaczy dokładnie „żadne wskazanie nie było
 * bez odpowiednika". Odmowy ta gałąź nie widzi i nie ma czego przesłonić.
 */
function rozlicz(odpowiedz: SessionDeleteResponse): RozliczenieUsuniecia {
  return {
    usuniete: odpowiedz.deletedIds,
    nieznalezione: odpowiedz.missingIds ?? [],
    liczba: odpowiedz.deletedCount,
  };
}
