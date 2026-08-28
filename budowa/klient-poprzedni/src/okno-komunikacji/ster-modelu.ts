import type { Agent, Channel } from '../../../shared/contract';
import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import {
  PRZEDROSTEK_AGENTA,
  wartoscBiezaca,
  zlecenieWyboru,
} from '../sterowanie/model-glowny';
import type { RejestrAgentow } from '../sterowanie/rejestr-agentow';
import type { RejestrKanalow } from '../sterowanie/rejestr-kanalow';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

// Model to ster paska zlecenia — dwie grupy, jeden wybór, tak samo jak w kolumnie sterowania.

/** Nazwa zmiany w komunikacie do rdzenia — ta sama, którą wysyła lista modeli w kolumnie sterowania oknem. */
const NAZWA = 'Model';

/** Etykieta uchwytu pokazywana, gdy okno nie ma jeszcze przypisanego ani kanału modelu, ani eksperta wcale. */
export const BRAK_MODELU = 'Model niewskazany';

/** Zależności steru modelu — wąskie i wstrzykiwane, obejmujące migawkę stanu oraz wysyłkę zmiany do rdzenia. */
export interface ZaleznosciSteruModelu {
  /** Kanał modelu i ekspert nałożony na okno, ze stanu potwierdzonego. */
  migawka(): { modelChannelId: string; agentId?: string };
  /** Wysyła zmianę pól okna; odrzucenie niesie zdanie odmowy rdzenia. */
  zastosuj(nazwa: string, zmiana: { modelChannelId?: string; agentId: string }): Promise<void>;
}

export function utworzSterModelu(
  zaleznosci: ZaleznosciSteruModelu,
  rejestrKanalow: RejestrKanalow,
  rejestrAgentow: RejestrAgentow,
): SterPaska {
  const ster = utworzSterNastawy({
    nastawa: NAZWA,
    ikona: 'agent',
    wykonaj: (klucz) =>
      zaleznosci.zastosuj(NAZWA, zlecenieWyboru(klucz, rejestrAgentow)),
    odswiez: () => odswiez(),
  });

  function odswiez(): void {
    const stan = zaleznosci.migawka();
    const wybrana = wartoscBiezaca(stan.agentId, stan.modelChannelId);
    ster.ustaw(etykieta(rejestrKanalow, rejestrAgentow, wybrana), drzewo(wybrana));
  }

  function drzewo(wybrana: string): PozycjaMenu[] {
    const pozycje: PozycjaMenu[] = [];
    const kanaly = rejestrKanalow.kanaly();
    const agenci = rejestrAgentow.agenci();

    // Grupa bez pozycji nie powstaje — nagłówek nad pustką byłby myleniem.
    if (kanaly.length > 0) {
      pozycje.push({
        rodzaj: 'grupa',
        nazwa: 'Modele',
        dzieci: kanaly.map((kanal) => ({
          rodzaj: 'wybor',
          klucz: kanal.id,
          nazwa: kanal.name,
          opis: opisKanalu(kanal),
          wybrany: kanal.id === wybrana,
        })),
      });
    }

    if (agenci.length > 0) {
      pozycje.push({
        rodzaj: 'grupa',
        nazwa: 'Moi agenci',
        dzieci: agenci.map((agent) => ({
          rodzaj: 'wybor',
          klucz: `${PRZEDROSTEK_AGENTA}${agent.id}`,
          nazwa: imieAgenta(agent),
          opis: opisAgenta(agent),
          wybrany: `${PRZEDROSTEK_AGENTA}${agent.id}` === wybrana,
        })),
      });
    }

    return pozycje;
  }

  rejestrKanalow.naZmiane(odswiez);
  rejestrAgentow.naZmiane(odswiez);
  rejestrAgentow.odswiez();
  odswiez();

  return { element: ster.element, odswiez };
}

/**
 * Wartość na uchwyt niesie nazwę własną eksperta albo kanału; wartość spoza obu rejestrów nie znika, tylko pokazuje surowy identyfikator.
 */
function etykieta(
  rejestrKanalow: RejestrKanalow,
  rejestrAgentow: RejestrAgentow,
  wybrana: string,
): string {
  if (wybrana === '') return BRAK_MODELU;
  if (wybrana.startsWith(PRZEDROSTEK_AGENTA)) {
    const kod = wybrana.slice(PRZEDROSTEK_AGENTA.length);
    const agent = rejestrAgentow.agenci().find((pozycja) => pozycja.id === kod);
    return agent === undefined ? kod : imieAgenta(agent);
  }
  const kanal = rejestrKanalow.kanaly().find((pozycja) => pozycja.id === wybrana);
  return kanal?.name ?? wybrana;
}

/** Imię własne eksperta ze znakiem graficznym; puste pole nazwy wyświetlanej zastępowane jest nazwą podstawową. */
function imieAgenta(agent: Agent): string {
  const znak = agent.favicon !== undefined && agent.favicon.length > 0 ? `${agent.favicon} ` : '';
  const imie =
    agent.displayName !== undefined && agent.displayName.length > 0
      ? agent.displayName
      : agent.name;
  return `${znak}${imie}`;
}

/** Zdanie przy kanale opisujące rodzaj drogi, użyty model oraz bieżący stan wpisu w rejestrze rdzenia platformy. */
function opisKanalu(kanal: Channel): string {
  const model = kanal.model !== undefined && kanal.model.length > 0 ? ` · ${kanal.model}` : '';
  const czynny = kanal.enabled ? '' : ' · kanał oznaczony jako nieczynny';
  return `Model surowy przez kanał ${kanal.kind}${model}${czynny}.`;
}

/**
 * Zdanie przy ekspercie: jego własny opis, a gdy go nie ma — zdanie mówiące,
 * na czym ekspert stoi. Bez zdania zapasowego pozycja byłaby w menu
 * nierozpoznawalna.
 */
function opisAgenta(agent: Agent): string {
  if (agent.description !== undefined && agent.description.length > 0) return agent.description;
  const model = agent.model !== undefined && agent.model.length > 0 ? agent.model : 'kanale bazowym';
  return `Ekspert nałożony na model ${model}.`;
}
