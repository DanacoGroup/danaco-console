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
 * Podagenci widziani przez klienta — jedno z trzech źródeł modułu, obok
 * `zrodlo-biegu.ts` (bieg pracy) i `zrodlo-okien.ts` (obsada ról).
 *
 * Okno nie zna nazw komend. Panel „Zadania w tle" woła `wykaz()` i `zbierz()`,
 * a nie `Command.SubagentList` — literał komendy stoi wyłącznie tutaj, więc gdy
 * kontrakt przemianuje komendę, kompilacja pęka w jednym pliku, a nie w każdym
 * widoku, który ją wołał.
 *
 * Źródło stoi osobno od `zrodlo-biegu.ts`, bo bieg pracy porusza oknami ról —
 * tura, przekazanie, kolejka etapu, telemetria — a podagent jest bytem pod
 * oknem wykonawcy, o którym kontrakt nie rozstrzyga, czym jest: oknem, procesem
 * czy pozycją kolejki.
 *
 * Telemetria zostaje w biegu: panel bierze `monitor.status` z `ZrodloBiegu`,
 * które tę komendę już niesie — drugie wywołanie tej samej komendy w drugim
 * źródle byłoby kopią bez powodu.
 */
export interface ZrodloPodagentow {
  /**
   * `subagent.spawn` — powołanie podagentów pod oknem wykonawcy.
   *
   * Oddaje podagentów tak, jak założył je rdzeń, a nie tak, jak prosił
   * formularz: liczba powołanych bierze się z odpowiedzi, bo rdzeń ma prawo
   * powołać ich mniej (górna granica piętnastu należy do niego, nie do widoku).
   */
  powolaj(zadanie: SubagentSpawnRequest): Promise<Wynik<Subagent[]>>;
  /**
   * `subagent.list` — podagenci okna wykonawcy albo całej karty sesji.
   *
   * Zawężenie po stanie robi rdzeń, nie widok: żądanie niesie pole `status`,
   * więc sito po stronie klienta byłoby drugim, rozjeżdżającym się sitem.
   */
  wykaz(zadanie: SubagentListRequest): Promise<Wynik<Subagent[]>>;
  /**
   * `subagent.result.collect` — zebranie wyników pracy podagentów.
   *
   * Oddaje komplet podagentów wraz z polem `complete`, które mówi, czy wszyscy
   * objęci zbieraniem domknęli pracę. Panel tego pola nie zgaduje ze stanów:
   * bierze je z odpowiedzi.
   */
  zbierz(zadanie: SubagentResultCollectRequest): Promise<Wynik<SubagentResultCollectResponse>>;
  /**
   * `subagent.changed` — podagent zmienił stan.
   *
   * Panel ma słuchać, nie odpytywać: rdzeń rozgłasza to zdarzenie przy
   * powołaniu, wejściu w bieg, zakończeniu i zatrzymaniu, więc żywy wykaz
   * nadąża za wykonawcą pracującym w tle bez ręcznego `subagent.list`.
   */
  naPodagenta(sluchacz: (tresc: SubagentChangedEvent) => void): Odsubskrybuj;
  /**
   * `subagent.stop` — zatrzymanie wskazanych podagentów okna wykonawcy.
   *
   * Zatrzymanie idzie tą komendą wprost, a nie obejściem przez `queue.action`
   * na kolejce o nazwie podagenta: podagent powołany bez wpiętej kolejki nie
   * miałby wtedy czego zatrzymać.
   *
   * Odpowiedź rozróżnia `stopped` od `notRunning`: podagent już zakończony nie
   * jest błędem, tylko innym stanem, więc panel mówi o nim osobno.
   */
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
