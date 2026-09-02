// Panel Skills Manager: umiejętności wpięte w definicję eksperta.
// Rdzeń nie ma komendy wykazu rejestru, więc wiersz niesie sam identyfikator.

import { Command, type Agent } from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  BEZ_EKSPERTA,
  pole,
  przelacz,
  pustka,
  tekst,
  wzorem,
  zdejmijPustke,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

export interface PanelUmiejetnosci {
  pokaz: (agent: Agent | null) => void;
}

export function zwiazUmiejetnosci(stan: StanowiskoEkspertow): PanelUmiejetnosci | null {
  const panel = stan.korzen.querySelector('#panel-skills');
  const trescMoze = panel?.querySelector('.sta-okno-tresc') ?? null;
  if (panel === null || trescMoze === null) return null;
  const tresc = trescMoze;

  const wzor = wzorem(tresc, '.pk-wiersz');
  const szukajka = pole(panel, '.dn-szukaj input');
  const dodaj = panel.querySelector('.sta-okno-tresc > button');
  const znacznik = panel.querySelector('.sta-okno-znacznik');

  function pokaz(agent: Agent | null): void {
    for (const stary of [...tresc.querySelectorAll('.pk-wiersz')]) stary.remove();
    zdejmijPustke(tresc);
    const identyfikatory = agent?.skillIds ?? [];
    tekst(znacznik, agent === null ? '' : `${identyfikatory.length} wpiętych`);
    if (agent === null) {
      pustka(tresc, BEZ_EKSPERTA);
      return;
    }
    if (identyfikatory.length === 0) {
      pustka(tresc, 'Ekspert nie ma wpiętej żadnej umiejętności.');
      return;
    }
    if (wzor === null) return;
    for (const identyfikator of identyfikatory) {
      const wiersz = wzor.cloneNode(true) as HTMLElement;
      wiersz.dataset.umiejetnosc = identyfikator;
      tekst(wiersz.querySelector('.nazwa'), identyfikator);
      tekst(wiersz.querySelector('.tresc small'), 'wpięta w definicję eksperta');
      const przelacznik = wiersz.querySelector('.dn-przelacznik');
      przelacz(przelacznik, true);
      przelacznik?.setAttribute('aria-label', identyfikator);
      tresc.insertBefore(wiersz, dodaj);
    }
  }

  async function odepnij(identyfikator: string): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentSkillRemove, {
      agentId: agent.id,
      skillId: identyfikator,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Umiejętności', wynik.blad?.message ?? 'Rdzeń odrzucił odpięcie.', 'blad');
      return;
    }
    stan.postaw(wynik.wynik.agent);
  }

  async function wepnij(): Promise<void> {
    const agent = stan.ekspert();
    const identyfikator = szukajka?.value.trim() ?? '';
    if (agent === null || identyfikator === '') return;
    const wynik = await wywolaj(stan.kanal, Command.AgentSkillAdd, {
      agentId: agent.id,
      skillId: identyfikator,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Umiejętności', wynik.blad?.message ?? 'Rdzeń odrzucił wpięcie.', 'blad');
      return;
    }
    if (szukajka !== null) szukajka.value = '';
    stan.postaw(wynik.wynik.agent);
  }

  tresc.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('.pk-wiersz');
    if (wiersz?.dataset.umiejetnosc === undefined) return;
    void odepnij(wiersz.dataset.umiejetnosc);
  }, stan.przy);

  dodaj?.addEventListener('click', () => {
    void wepnij();
  }, stan.przy);

  szukajka?.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') void wepnij();
  }, stan.przy);

  return { pokaz };
}
