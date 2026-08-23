import {
  Command,
  EventType,
  type HistoryChangedEvent,
  type HistoryDeleteRequest,
  type HistoryDeleteResponse,
  type HistoryLoadRequest,
  type HistoryLoadResponse,
  type RetentionSetRequest,
  type RetentionSetResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Źródło historii rozmowy okna — trzy komendy kontraktu i jedna subskrypcja.
 *
 * Nazwy komend padają tylko tutaj: widok panelu nie zna ani jednej stałej
 * `Command.*`, dostaje cztery funkcje i tyle. Zmiana nazwy w kontrakcie
 * przerywa kompilację w tym jednym pliku.
 *
 * `history.changed` jest zdarzeniem, nie czwartą komendą. Rdzeń rozgłasza je po
 * każdym skasowaniu, także po tym z zasady przechowywania nadanej w innym oknie
 * i po przemiataniu przy starcie rdzenia. Bez tej subskrypcji panel pokazywałby
 * pozycje, których w bazie już nie ma. Pole `entry` jest niewymagane: po
 * czyszczeniu zbiorczym pozycji „po zmianie" nie ma żadnej i panel przeładowuje
 * wtedy wykaz w całości, zamiast doszywać wiersz.
 *
 * Źródło nie trzyma stanu, nie zna okna i nie rozstrzyga, czy zasada ma sens —
 * to należy do panelu i do rdzenia. Zdarzeń po oknie nie filtruje, bo filtr
 * wymaga znajomości okna gospodarza, a ta jest w panelu.
 */
export interface ZrodloHistorii {
  /** `history.load` — wykaz od pozycji najnowszej, stronicowany kursorem czasu. */
  wczytaj(zadanie: HistoryLoadRequest): Promise<Wynik<HistoryLoadResponse>>;
  /** `history.delete` — usuwa wskazane pozycje, a bez wskazania całą historię okna. */
  usun(zadanie: HistoryDeleteRequest): Promise<Wynik<HistoryDeleteResponse>>;
  /** `retention.set` — zapisuje zasadę przechowywania dla zakresu. */
  ustawZasade(zadanie: RetentionSetRequest): Promise<Wynik<RetentionSetResponse>>;
  /** Subskrypcja `history.changed` — zmiany cudze i zmiany zasady. */
  naZmianeHistorii(sluchacz: (tresc: HistoryChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloHistorii(kanal: Kanal): ZrodloHistorii {
  return {
    async wczytaj(zadanie) {
      // Sprawdzian obejmuje oba pola, bo panel rozstrzyga z nich dwie różne
      // rzeczy: `entries` daje wiersze, a `total` daje zdanie „jest ich więcej
      // niż widzisz". Brak `total` przy obecnych wierszach kazałby zgadywać,
      // czy jest co dociągać.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HistoryLoad, zadanie),
        Command.HistoryLoad,
        (tresc) => czyTablica(tresc.entries) && czyLiczba(tresc.total),
      );
    },

    async usun(zadanie) {
      // `deleted` jest treścią właściwą, nie ozdobą potwierdzenia: panel ma
      // powiedzieć, ile pozycji ubyło, więc odpowiedź bez tej liczby nie jest
      // potwierdzeniem czynności.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HistoryDelete, zadanie),
        Command.HistoryDelete,
        (tresc) => czyLiczba(tresc.deleted),
      );
    },

    async ustawZasade(zadanie) {
      // Rdzeń oddaje zasadę w kształcie obowiązującym po zmianie, a nie echo
      // żądania, i to jego odpowiedź panel pokazuje. Sprawdzamy `scope`, bo po
      // nim panel nazywa zakres; progi zostają niesprawdzone, bo ich brak jest
      // wartością prawdziwą — zasadą bez ograniczenia.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RetentionSet, zadanie),
        Command.RetentionSet,
        (tresc) => czyObiekt(tresc.policy) && czyTekst(tresc.policy.scope),
      );
    },

    naZmianeHistorii(sluchacz) {
      return kanal.naZdarzenie(EventType.HistoryChanged, (tresc) => sluchacz(tresc));
    },
  };
}
