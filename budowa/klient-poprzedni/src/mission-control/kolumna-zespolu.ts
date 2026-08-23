import { WindowRole } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { ETYKIETA_BRAKU_ZRODLA, ETYKIETA_ROLI_OKNA } from './etykiety-pulpitu';
import type { AgentZespolu } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import { utworzStanPusty } from './stan-pusty';

/**
 * Kolumna zespołu agentów: kto z zespołu pracuje i nad czym.
 *
 * Źródłem są otwarte okna komunikacji z odczytu `window.list` — rola okna
 * w pętli jest rolą agenta, a zajęcie pochodzi z etapu telemetrii
 * `progress.changed`. Liczba podagentów nie ma źródła w kontrakcie, więc wiersz
 * mówi to wprost zamiast pokazać wartość zastępczą.
 */
export interface KolumnaZespolu {
  element: HTMLElement;
  odswiez(zespol: AgentZespolu[]): void;
}

/** Buduje kolumnę zespołu agentów. */
export function utworzKolumneZespolu(zespol: AgentZespolu[]): KolumnaZespolu {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--zespol',
    {
      tytul: 'Zespół agentów',
      dopisek: 'Kto pracuje, w jakiej roli i z iloma podagentami.',
      ikona: 'uzytkownik',
    },
    'mc-tytul-zespol',
  );

  const wykaz = document.createElement('ul');
  wykaz.className = 'mc-zespol';
  cialo.append(wykaz);

  const odswiez = (dane: AgentZespolu[]): void => {
    if (dane.length === 0) {
      wykaz.replaceChildren(
        utworzStanPusty(
          'Zespół nie pracuje',
          'Odczyt window.list nie zna żadnego otwartego okna komunikacji.',
        ),
      );
      return;
    }
    wykaz.replaceChildren(...dane.map(wiersz));
  };
  odswiez(zespol);

  return { element, odswiez };
}

/** Jeden agent: awatar z inicjałem, nazwa, rola, zajęcie i liczba podagentów. */
function wiersz(agent: AgentZespolu): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'mc-agent';
  element.dataset.rola = agent.rola;

  const awatar = document.createElement('span');
  awatar.className = agent.rola === WindowRole.Coordinator
    ? 'dn-awatar dn-awatar--sm dn-awatar--inteligencja mc-agent__awatar'
    : 'dn-awatar dn-awatar--sm mc-agent__awatar';
  awatar.setAttribute('aria-hidden', 'true');
  awatar.textContent = agent.nazwa.slice(0, 1).toUpperCase();

  const blok = document.createElement('div');
  blok.className = 'mc-agent__blok';

  const gora = document.createElement('div');
  gora.className = 'mc-agent__gora';

  const nazwa = document.createElement('span');
  nazwa.className = 'mc-agent__nazwa';
  nazwa.textContent = agent.nazwa;

  const rola = document.createElement('span');
  rola.className = 'dn-plakietka dn-plakietka--rola mc-agent__rola';
  rola.textContent = ETYKIETA_ROLI_OKNA[agent.rola];

  gora.append(nazwa, rola);

  const zajecie = document.createElement('span');
  zajecie.className = 'mc-agent__zajecie';
  // Stan pracy niesie ikonę i słowo, nigdy samą barwę.
  zajecie.append(
    elementIkony(agent.czynny ? 'uruchom' : 'zegar', { rozmiar: 16 }),
    document.createTextNode(
      agent.zajecie ?? (agent.czynny ? 'proces bez nazwy etapu' : 'bez procesu w telemetrii'),
    ),
  );

  const podagenci = document.createElement('span');
  podagenci.className = 'mc-agent__podagenci';
  podagenci.textContent =
    agent.podagenci === null
      ? `podagenci: ${ETYKIETA_BRAKU_ZRODLA}`
      : `${agent.podagenci} podagentów`;
  podagenci.title =
    agent.podagenci === null
      ? 'Kontrakt nie niesie liczby podagentów okna. Brakujący odczyt zgłoszony.'
      : 'Wykonawca może uruchomić do 15 podagentów';

  blok.append(gora, zajecie, podagenci);
  element.append(awatar, blok);
  return element;
}
