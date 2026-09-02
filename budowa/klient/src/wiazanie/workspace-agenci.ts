/**
 * Panel ekspertów projektu, wywołania z poprzedniego klienta
 * (`zarzadca-agentow.ts`): `workspace.dashboard.get`, `agent.list`,
 * `workspace.agent.assign` i `workspace.agent.unassign`.
 */

import { Command, type Agent, type WorkspaceDashboard } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { odmien } from './workspace-projekt.ts';
import {
  cialoPanelu,
  kopiaWzoru,
  niegotowe,
  plakietkaPanelu,
  przyjmij,
  usun,
  wpisz,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

export interface StanAgentow {
  opisz: (pulpit: WorkspaceDashboard) => void;
}

export function zwiazAgentow(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  poZmianie: () => Promise<void>,
  przy: AddEventListenerOptions,
): StanAgentow {
  const panel = wezly.agenci;
  const cialo = cialoPanelu(panel);
  const wzor = kopiaWzoru(cialo?.querySelector<HTMLElement>('.dn-wykaz-modulu-poz') ?? null);

  // Założenie nowego eksperta prowadzi inny moduł; to okno tylko przypisuje istniejącego.
  usun(cialo?.querySelector('.dn-btn'));
  for (const stary of cialo?.querySelectorAll('.dn-wykaz-modulu-poz') ?? []) stary.remove();

  cialo?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const kandydat = cel.closest<HTMLElement>('[data-kandydat]');
    if (kandydat !== null) {
      void przypisz(kanal, idProjektu(), kandydat.dataset.kandydat ?? '', poZmianie);
      return;
    }
    const wierszMoze = cel.closest<HTMLElement>('[data-agent]');
    if (wierszMoze === null) return;
    const wiersz = wierszMoze;
    void odlacz(kanal, idProjektu(), wiersz.dataset.agent ?? '', poZmianie);
  }, przy);

  const dodaj = panel?.querySelector<HTMLElement>(
    '.sta-okno-belka .dn-btn-ikona:not([data-okno-zamknij])',
  );
  dodaj?.addEventListener('click', () => {
    void pokazKandydatow(kanal, idProjektu(), cialo, wzor);
  }, przy);

  return {
    opisz: (pulpit) => {
      if (cialo === null) return;
      const agenci = pulpit.assignedAgentIds ?? [];
      wpisz(plakietkaPanelu(panel), odmien(agenci.length, ['agent', 'agenci', 'agentów']));
      cialo.replaceChildren();
      if (agenci.length === 0) {
        niegotowe(cialo, 'Do tego projektu nie przypisano jeszcze eksperta.');
        return;
      }
      for (const agent of agenci) {
        const wiersz = wierszAgenta(wzor, agent, agent, 'Odłącz', 'agent');
        if (wiersz !== null) cialo.appendChild(wiersz);
      }
    },
  };
}

/** Wykaz kandydatów zawężony do projektu — takim zapytaniem szedł poprzedni klient. */
async function pokazKandydatow(
  kanal: Kanal,
  idProjektu: string,
  cialo: HTMLElement | null,
  wzor: HTMLElement | null,
): Promise<void> {
  if (cialo === null || idProjektu === '') return;
  const wynik = await wywolaj(kanal, Command.AgentList, {
    projectId: idProjektu,
    enabledOnly: true,
    limit: 50,
  });
  const agenci = przyjmij('Eksperci', wynik)?.agents ?? [];
  cialo.replaceChildren();
  if (agenci.length === 0) {
    niegotowe(cialo, 'Rdzeń nie ma eksperta, którego można przypisać.');
    return;
  }
  for (const agent of agenci) {
    const wiersz = wierszKandydata(wzor, agent);
    if (wiersz !== null) cialo.appendChild(wiersz);
  }
}

function wierszKandydata(wzor: HTMLElement | null, agent: Agent): HTMLElement | null {
  return wierszAgenta(wzor, agent.id, agent.displayName ?? agent.name, 'Przypisz', 'kandydat');
}

function wierszAgenta(
  wzor: HTMLElement | null,
  identyfikator: string,
  nazwa: string,
  czynnosc: string,
  rodzaj: 'agent' | 'kandydat',
): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const awatar = wiersz.querySelector<HTMLElement>('.rt-awatar');
  const meta =
    wiersz.querySelector<HTMLElement>('.dn-meta') ??
    wiersz.querySelector<HTMLElement>('.dn-plakietka');
  wiersz.replaceChildren();
  if (awatar !== null) {
    awatar.textContent = nazwa.slice(0, 2);
    wiersz.appendChild(awatar);
  }
  wiersz.append(` ${nazwa} `);
  if (meta !== null) {
    meta.className = 'dn-meta';
    meta.textContent = czynnosc;
    wiersz.appendChild(meta);
  }
  if (rodzaj === 'agent') wiersz.dataset.agent = identyfikator;
  else wiersz.dataset.kandydat = identyfikator;
  return wiersz;
}

async function przypisz(
  kanal: Kanal,
  idProjektu: string,
  idAgenta: string,
  poZmianie: () => Promise<void>,
): Promise<void> {
  if (idProjektu === '' || idAgenta === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceAgentAssign, {
    projectId: idProjektu,
    agentId: idAgenta,
  });
  if (przyjmij('Eksperci projektu', wynik) === null) return;
  oglos('Eksperci projektu', `Przypisano ${idAgenta}.`);
  await poZmianie();
}

async function odlacz(
  kanal: Kanal,
  idProjektu: string,
  idAgenta: string,
  poZmianie: () => Promise<void>,
): Promise<void> {
  if (idProjektu === '' || idAgenta === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceAgentUnassign, {
    projectId: idProjektu,
    agentId: idAgenta,
  });
  if (przyjmij('Eksperci projektu', wynik) === null) return;
  oglos('Eksperci projektu', `Odłączono ${idAgenta}.`);
  await poZmianie();
}
