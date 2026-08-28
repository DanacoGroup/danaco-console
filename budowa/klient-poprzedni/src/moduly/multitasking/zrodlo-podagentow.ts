import {
  Command,
  EventType,
  type SubagentChangedEvent,
  type Subagent,
  type SubagentListRequest,
  type SubagentResultCollectRequest,
  type SubagentResultCollectResponse,
  type SubagentSpawnRequest,
  type SubagentStopRequest,
  type SubagentStopResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Źródło podagentów udostępnia klientowi komendy obszaru `subagent.*` wraz ze
 * zdarzeniem zmiany stanu. Literały nazw komend stoją wyłącznie tutaj, więc
 * przemianowanie komendy w kontrakcie zatrzymuje kompilację w jednym pliku.
 */
export interface ZrodloPodagentow {
  /** `subagent.spawn` — powołanie podagentów; liczbę powołanych podaje rdzeń. */
  powolaj(zadanie: SubagentSpawnRequest): Promise<Wynik<Subagent[]>>;
  /** `subagent.list` — podagenci okna wykonawcy albo całej karty sesji. */
  wykaz(zadanie: SubagentListRequest): Promise<Wynik<Subagent[]>>;
  /** `subagent.result.collect` — zebranie wyników wraz z polem `complete`. */
  zbierz(zadanie: SubagentResultCollectRequest): Promise<Wynik<SubagentResultCollectResponse>>;
  /** `subagent.changed` — rdzeń rozgłasza każdą zmianę stanu podagenta. */
  naPodagenta(sluchacz: (tresc: SubagentChangedEvent) => void): Odsubskrybuj;
  /** `subagent.stop` — zatrzymanie podagentów; odpowiedź dzieli dwa stany. */
  zatrzymaj(zadanie: SubagentStopRequest): Promise<Wynik<SubagentStopResponse>>;
}

export function utworzZrodloPodagentow(kanal: Kanal): ZrodloPodagentow {
  return {
    async powolaj(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.SubagentSpawn, zadanie),
        Command.SubagentSpawn,
        (tresc) => czyTablica(tresc.subagents),
      );
      return przenies(wynik, (tresc) => tresc.subagents);
    },

    async wykaz(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.SubagentList, zadanie),
        Command.SubagentList,
        (tresc) => czyTablica(tresc.subagents),
      );
      return przenies(wynik, (tresc) => tresc.subagents);
    },

    async zatrzymaj(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SubagentStop, zadanie),
        Command.SubagentStop,
        (tresc) =>
          czyTablica(tresc.stopped) && czyTablica(tresc.notRunning) && czyTablica(tresc.subagents),
      );
    },

    async zbierz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SubagentResultCollect, zadanie),
        Command.SubagentResultCollect,
        (tresc) => czyTablica(tresc.subagents) && typeof tresc.complete === 'boolean',
      );
    },

    naPodagenta: (sluchacz) => kanal.naZdarzenie(EventType.SubagentChanged, sluchacz),
  };
}
