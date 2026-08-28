import { Command, type Agent } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/** Wykaz ekspertów operatora czytany z rdzenia — druga sekcja wyboru modelu, osobna od rejestru kanałów. */
export interface RejestrAgentow {
  /** Eksperci znani w tej chwili; pusty wykaz, dopóki rdzeń nie odpowie. */
  agenci(): Agent[];
  // Czy rdzeń odpowiedział na wykaz ekspertów choć raz; pusty wykaz bywa dwuznaczny.
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
      // Tylko eksperci czynni: wyłączony ekspert nie wchodzi do wyboru modelu.
      kanal.wyslij(Command.AgentList, { enabledOnly: true }, (wynik) => {
        // Odpowiedź odnotowuje się także po odmowie, żeby wskaźnik odczytu zgasł.
        odpowiedziano = true;
        const odebrane = wynik.wynik?.agents;
        wykaz = wynik.udany && odebrane !== undefined ? odebrane : wykaz;
        zmiany.oglos(wykaz);
      });
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Nazwa eksperta pokazywana operatorowi: imię własne wyprzedza nazwę techniczną, favikon idzie przed imieniem. */
export function nazwaAgenta(agent: Agent): string {
  const znak = agent.favicon !== undefined && agent.favicon.length > 0 ? `${agent.favicon} ` : '';
  const imie = agent.displayName !== undefined && agent.displayName.length > 0
    ? agent.displayName
    : agent.name;
  const model = agent.model !== undefined && agent.model.length > 0 ? ` · ${agent.model}` : '';
  return `${znak}${imie}${model}`;
}
