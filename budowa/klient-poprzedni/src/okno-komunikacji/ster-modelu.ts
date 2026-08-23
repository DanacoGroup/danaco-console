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

/**
 * Model — ster paska zlecenia. Dwie grupy, jeden wybór, tak samo jak
 * w kolumnie sterowania: modele surowe w grupie „Modele", eksperci
 * w „Moi agenci". Grupa, a nie gałąź, bo obie rodziny mieszczą się na jednym
 * poziomie, a gałąź schowałaby wykaz o jeden ruch dalej.
 *
 * Znaczenie wyboru jest pisane raz: przedrostek eksperta, wartość zaznaczona
 * i przekład wyboru na treść `window.update` pochodzą z
 * `sterowanie/model-glowny.ts`. Gdyby pasek składał to zlecenie po swojemu,
 * wybór eksperta w pasku i w kolumnie znaczyłby dwie różne rzeczy.
 *
 * Etykieta uchwytu jest krótka: w kolumnie sterowania wiersz niesie
 * `nazwaKanalu` („nazwa (rodzaj) · model · nieczynny"), a w pasku stoi sama
 * nazwa własna i reszta schodzi do opisu pozycji, bo pasek ma zostać jednym
 * rzędem. Oba napisy składane są z tych samych pól `Channel`/`Agent`.
 *
 * Kanał nieczynny zostaje na wykazie i pozostaje wybieralny — wyszarzenie
 * byłoby blokadą, a o stanie kanału mówi opis pozycji.
 */

/** Nazwa zmiany w komunikacie — ta sama, którą wysyła lista w kolumnie. */
const NAZWA = 'Model';

/** Etykieta uchwytu, gdy okno nie ma jeszcze ani kanału, ani eksperta. */
export const BRAK_MODELU = 'Model niewskazany';

/** Zależności steru — wąskie i wstrzykiwane. */
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
 * Wartość na uchwyt: nazwa własna eksperta albo kanału.
 *
 * Wartość spoza obu rejestrów nie znika — uchwyt pokazuje surowy identyfikator.
 * Okno pracuje wtedy na kanale, którego wykaz jeszcze nie zna (rdzeń nie
 * odpowiedział) albo już nie zna; napis „Model niewskazany" byłby w obu
 * przypadkach nieprawdą.
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

/** Imię własne eksperta ze znakiem graficznym; puste `displayName` = `name`. */
function imieAgenta(agent: Agent): string {
  const znak = agent.favicon !== undefined && agent.favicon.length > 0 ? `${agent.favicon} ` : '';
  const imie =
    agent.displayName !== undefined && agent.displayName.length > 0
      ? agent.displayName
      : agent.name;
  return `${znak}${imie}`;
}

/** Zdanie przy kanale: rodzaj drogi, model i stan wpisu rejestru. */
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
