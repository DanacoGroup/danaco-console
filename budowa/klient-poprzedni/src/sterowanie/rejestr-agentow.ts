import { Command, type Agent } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/**
 * Wykaz ekspertów Operatora czytany z rdzenia — druga sekcja wyboru modelu.
 *
 * Moduł Agents komponuje ekspertów, a ich wynik ma dać się wybrać we wszystkich
 * oknach roboczych — bez tego wykazu portfolio jest martwe.
 *
 * Rejestr jest osobny, a nie doklejony do kanałów, bo agent kanałem nie jest:
 * kanał to droga do modelu, agent to tożsamość nałożona na tę drogę. Wrzucenie
 * go do `rejestr-kanalow.ts` zatarłoby tę różnicę już w typie, a przy setce
 * agentów surowe modele utonęłyby w wykazie. Dwa rejestry znaczą dwie sekcje
 * jednego menu: „Modele" i „Moi agenci".
 *
 * Kształt jest bliźniaczy wobec `rejestr-kanalow.ts` z zamysłem: ta sama umowa
 * (`kanaly`/`agenci`, `odpowiedzOtrzymana`, `odswiez`, `naZmiane`) pozwala
 * sterowaniu modelu obsłużyć oba wykazy jedną drogą, bez rozgałęzień na typ.
 *
 * Rejestr nie zna okna i niczego w nim nie zapisuje — jest katalogiem wyboru
 * wspólnym dla klienta, nie ustawieniem okna. Nie filtruje też ekspertów po
 * module ani po środowisku: agent nie jest własnością modułu.
 */
export interface RejestrAgentow {
  /** Eksperci znani w tej chwili; pusty wykaz, dopóki rdzeń nie odpowie. */
  agenci(): Agent[];
  /**
   * Czy rdzeń odpowiedział na `agent.list` choć raz.
   *
   * Pusty wykaz znaczy dwie różne rzeczy — „jeszcze nie wiem" i „Operator nie
   * ma ani jednego eksperta". Pierwsza to wskaźnik odczytu, druga to stan
   * poprawny, o którym nie ma po co mówić ani słowa.
   */
  odpowiedzOtrzymana(): boolean;
  /** Zamawia wykaz z rdzenia. */
  odswiez(): void;
  /** Subskrypcja zmian wykazu. */
  naZmiane(sluchacz: (agenci: Agent[]) => void): Odsubskrybuj;
}

export function utworzRejestrAgentow(kanal: Kanal): RejestrAgentow {
  const zmiany = utworzMagistrale<Agent[]>();
  let wykaz: Agent[] = [];
  let odpowiedziano = false;

  return {
    agenci: () => wykaz,

    odpowiedzOtrzymana: () => odpowiedziano,

    odswiez() {
      // Tylko eksperci czynni: ekspert wyłączony jest wyłączony przez Operatora
      // i wystawianie go do wyboru byłoby cofaniem jego decyzji.
      kanal.wyslij(Command.AgentList, { enabledOnly: true }, (wynik) => {
        // Odpowiedź odnotowuje się także po odmowie: pytanie zostało
        // rozstrzygnięte, więc wskaźnik odczytu ma zgasnąć zamiast wisieć.
        odpowiedziano = true;
        const odebrane = wynik.wynik?.agents;
        wykaz = wynik.udany && odebrane !== undefined ? odebrane : wykaz;
        zmiany.oglos(wykaz);
      });
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/**
 * Nazwa eksperta pokazywana Operatorowi.
 *
 * Imię własne wyprzedza nazwę techniczną: `name` jest nazwą, pod którą ekspert
 * zapisał się w rejestrze, a `displayName` tym, jak Operator go woła. Favikon
 * idzie przed imieniem, bo przy setce pozycji znak rozpoznaje się szybciej niż
 * napis.
 */
export function nazwaAgenta(agent: Agent): string {
  const znak = agent.favicon !== undefined && agent.favicon.length > 0 ? `${agent.favicon} ` : '';
  const imie = agent.displayName !== undefined && agent.displayName.length > 0
    ? agent.displayName
    : agent.name;
  const model = agent.model !== undefined && agent.model.length > 0 ? ` · ${agent.model}` : '';
  return `${znak}${imie}${model}`;
}
