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
 * Moduł Roundtable widziany przez klienta — cztery komendy obszaru
 * `roundtable.*` i jedno zdarzenie przyrostu debaty.
 *
 * Cztery, choć obszar niesie ich dziś czterdzieści sześć: pozostałych to źródło
 * jeszcze nie wywołuje, bo ich obsługi nie zbudowano. Wykaz tego, czego brakuje
 * i którą komendą to wejdzie, prowadzi `katalog-funkcji.ts`; okna nazywają
 * brakującą obsługę przy każdej pozycji z osobna.
 *
 * Pole `windowId` jest wymagane w każdej z czterech czynności, bo debata jest
 * bytem okna, nie sesji: uczestnicy, tury i stanowisko należą do jednego okna
 * debaty i rdzeń bez tego identyfikatora nie wie, o którą debatę pytamy. Wymóg
 * stoi w kontrakcie, więc żadne pole żądania nie ma tu wartości domyślnej —
 * typy żądań kontraktu wnosi się w całości.
 *
 * Każda czynność oddaje `Wynik`, nie samą treść. Debata z wieloma modelami
 * odmawia z powodów zwyczajnych — kanał uczestnika bywa nieczynny, tura bywa
 * już zamknięta — a okno musi odróżnić „debata jeszcze nic nie powiedziała" od
 * „nie udało się zapytać". Dlatego nie ma tu ani jednego `?? []`.
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
      // Kształt sprawdzamy po obu polach koperty: tura mówi, co się otworzyło,
      // a wykaz uczestników — kto w tej turze odpowiada. Okno rysuje jedno
      // obok drugiego, więc brak któregokolwiek jest odpowiedzią bezużyteczną.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.RoundtableDebateStart, zadanie),
        Command.RoundtableDebateStart,
        (tresc) => czyObiekt(tresc.turn) && czyTablica(tresc.participantIds),
      );
    },

    async moderuj(zadanie) {
      // `participants` jest w kontrakcie polem nieobowiązkowym — rdzeń dokłada
      // je tylko wtedy, gdy czynność zmieniła skład albo kolejność głosu.
      // Sprawdzana jest więc wyłącznie tura, obecna zawsze.
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
