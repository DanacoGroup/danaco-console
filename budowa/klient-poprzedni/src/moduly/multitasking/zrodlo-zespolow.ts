import {
  Command,
  EventType,
  type Agent,
  type AgentListRequest,
  type Team,
  type TeamChangedEvent,
  type TeamDuplicateRequest,
  type TeamListRequest,
  type TeamLoadRequest,
  type TeamSaveRequest,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Źródło sekcji zespołów panelu orkiestracji — jedyne miejsce modułu, które zna
 * nazwy rodziny `team.*` oraz komendy `agent.list`. Zespół i ekspert stoją
 * w jednym źródle, ponieważ zespół kontraktu jest wyłącznie wykazem
 * identyfikatorów ekspertów.
 */
export interface ZrodloZespolow {
  /** `team.list` — zespoły zapisane przez Operatora, w kolejności nazwy. */
  zespoly(zadanie: TeamListRequest): Promise<Wynik<Team[]>>;
  /** `team.save` — zapis zespołu; brak identyfikatora zakłada nowy. */
  zapisz(zadanie: TeamSaveRequest): Promise<Wynik<Team>>;
  /** `team.load` — jeden zespół wraz ze składem. */
  wczytaj(zadanie: TeamLoadRequest): Promise<Wynik<Team>>;
  /** `team.duplicate` — kopia zespołu pod nową nazwą. */
  powiel(zadanie: TeamDuplicateRequest): Promise<Wynik<Team>>;
  /** `agent.list` — eksperci, z których składa się zespół. */
  eksperci(zadanie: AgentListRequest): Promise<Wynik<Agent[]>>;
  /** Subskrypcja `team.changed` — zmiana zespołu poza tym widokiem. */
  naZespol(sluchacz: (tresc: TeamChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloZespolow(kanal: Kanal): ZrodloZespolow {
  return {
    async zespoly(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamList, zadanie),
        Command.TeamList,
        (tresc) => czyTablica(tresc.teams),
      );
      return przenies(wynik, (tresc) => tresc.teams);
    },

    async zapisz(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamSave, zadanie),
        Command.TeamSave,
        (tresc) => czyObiekt(tresc.team),
      );
      return przenies(wynik, (tresc) => tresc.team);
    },

    async wczytaj(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamLoad, zadanie),
        Command.TeamLoad,
        (tresc) => czyObiekt(tresc.team),
      );
      return przenies(wynik, (tresc) => tresc.team);
    },

    async powiel(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamDuplicate, zadanie),
        Command.TeamDuplicate,
        (tresc) => czyObiekt(tresc.team),
      );
      return przenies(wynik, (tresc) => tresc.team);
    },

    async eksperci(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentList, zadanie),
        Command.AgentList,
        (tresc) => czyTablica(tresc.agents),
      );
      return przenies(wynik, (tresc) => tresc.agents);
    },

    naZespol: (sluchacz) => kanal.naZdarzenie(EventType.TeamChanged, sluchacz),
  };
}
