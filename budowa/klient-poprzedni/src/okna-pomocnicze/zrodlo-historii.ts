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
 * Źródło historii rozmowy okna opiera się na trzech komendach kontraktu i jednej subskrypcji, przy czym nazwy komend padają wyłącznie w tym pliku, więc zmiana nazwy w kontrakcie przerywa kompilację właśnie tutaj.
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
      // Sprawdzian obejmuje oba pola: wiersze oraz zdanie o tym, czy jest ich więcej niż widać.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HistoryLoad, zadanie),
        Command.HistoryLoad,
        (tresc) => czyTablica(tresc.entries) && czyLiczba(tresc.total),
      );
    },

    async usun(zadanie) {
      // Liczba usuniętych jest treścią właściwą, nie ozdobą — bez niej to nie jest potwierdzenie.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HistoryDelete, zadanie),
        Command.HistoryDelete,
        (tresc) => czyLiczba(tresc.deleted),
      );
    },

    async ustawZasade(zadanie) {
      // Rdzeń oddaje zasadę obowiązującą po zmianie, nie echo żądania; progi zostają niesprawdzone.
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
