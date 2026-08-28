import {
  Command,
  EventType,
  type RoundtableConsensus,
  type RoundtableConsensusGetRequest,
  type RoundtableDebateChangedEvent,
  type RoundtableDebateStartRequest,
  type RoundtableDebateStartResponse,
  type RoundtableModelAddRequest,
  type RoundtableModeratorDirectRequest,
  type RoundtableModeratorDirectResponse,
  type RoundtableParticipant,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Moduł opisuje interfejs klienta do obszaru roundtable: cztery komendy debaty wieloosobowej
 * i jedno zdarzenie jej przyrostu, każda czynność zwraca Wynik zamiast samej treści.
 */
export interface ZrodloRoundtable {
  /** `roundtable.model.add` — dodanie uczestnika debaty pod własną tożsamością. */
  dodajModel(zadanie: RoundtableModelAddRequest): Promise<Wynik<RoundtableParticipant>>;
  /** `roundtable.debate.start` — pytanie kierowane jednocześnie do wszystkich. */
  uruchomDebate(
    zadanie: RoundtableDebateStartRequest,
  ): Promise<Wynik<RoundtableDebateStartResponse>>;
  /** `roundtable.moderator.direct` — czynność moderatora nad turą i uczestnikami. */
  moderuj(
    zadanie: RoundtableModeratorDirectRequest,
  ): Promise<Wynik<RoundtableModeratorDirectResponse>>;
  /** `roundtable.consensus.get` — stanowisko końcowe tury albo całej debaty. */
  stanowisko(zadanie: RoundtableConsensusGetRequest): Promise<Wynik<RoundtableConsensus>>;
  /** Subskrypcja `roundtable.debate.changed` — Debate Panel na żywo. */
  naZmianeDebaty(sluchacz: (tresc: RoundtableDebateChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloRoundtable(kanal: Kanal): ZrodloRoundtable {
  return {
    async dodajModel(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModelAdd, zadanie),
        Command.RoundtableModelAdd,
        (tresc) => czyObiekt(tresc.participant),
      );
      return przenies(wynik, (tresc) => tresc.participant);
    },

    async uruchomDebate(zadanie) {
      // Sprawdzenie kształtu obejmuje oba pola koperty, bo okno rysuje turę i wykaz uczestników razem.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDebateStart, zadanie),
        Command.RoundtableDebateStart,
        (tresc) => czyObiekt(tresc.turn) && czyTablica(tresc.participantIds),
      );
    },

    async moderuj(zadanie) {
      // Pole participants jest w kontrakcie nieobowiązkowe, sprawdzana jest więc wyłącznie tura.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableModeratorDirect, zadanie),
        Command.RoundtableModeratorDirect,
        (tresc) => czyObiekt(tresc.turn),
      );
    },

    async stanowisko(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableConsensusGet, zadanie),
        Command.RoundtableConsensusGet,
        (tresc) => czyObiekt(tresc.consensus),
      );
      return przenies(wynik, (tresc) => tresc.consensus);
    },

    naZmianeDebaty(sluchacz) {
      return kanal.naZdarzenie(EventType.RoundtableDebateChanged, (tresc) => sluchacz(tresc));
    },
  };
}
