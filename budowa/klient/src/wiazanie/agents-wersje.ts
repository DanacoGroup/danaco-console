// Historia wersji eksperta: wykaz utrwalonych tożsamości, podgląd jednej z nich,
// przywrócenie oraz założenie nowego eksperta z bieżącej definicji.

import { Command, type Agent } from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  BEZ_EKSPERTA,
  godzinaZapisu,
  nazwaEksperta,
  pustka,
  tekst,
  wzorem,
  zdejmijPustke,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

export interface PanelWersji {
  pokaz: (agent: Agent | null) => void;
}

export function zwiazWersje(
  stan: StanowiskoEkspertow,
  otworz: (agent: Agent) => void,
): PanelWersji | null {
  const historia = stan.korzen.querySelector('.ab-historia');
  const listaMoze = historia?.querySelector('.dn-wykaz-modulu') ?? null;
  if (historia === null || listaMoze === null) return null;
  const lista = listaMoze;

  const tytul = historia.querySelector('.ab-sum-tytul');
  const przydomek = tytul?.querySelector('span') ?? null;
  const wiecej = historia.querySelector('button');
  const wzor = wzorem(lista, '.dn-wykaz-modulu-poz');

  function pokaz(agent: Agent | null): void {
    tekst(przydomek, agent === null ? '' : `— ${nazwaEksperta(agent)}`);
    for (const stary of [...lista.querySelectorAll('.dn-wykaz-modulu-poz')]) stary.remove();
    zdejmijPustke(lista);
    if (agent === null) {
      pustka(lista, BEZ_EKSPERTA);
      return;
    }
    void wczytaj(agent);
  }

  async function wczytaj(agent: Agent): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.AgentVersionList, {
      agentId: agent.id,
      limit: 50,
    });
    for (const stary of [...lista.querySelectorAll('.dn-wykaz-modulu-poz')]) stary.remove();
    zdejmijPustke(lista);
    const wersje = wynik.udany ? wynik.wynik?.versions ?? [] : [];
    if (wersje.length === 0 || wzor === null) {
      pustka(lista, 'Rdzeń nie oddał żadnej utrwalonej wersji tego eksperta.');
      return;
    }
    for (const wersja of wersje) {
      const wiersz = wzor.cloneNode(true) as HTMLElement;
      wiersz.dataset.wersja = wersja.id;
      tekst(wiersz.querySelector('b'), wersja.label ?? wersja.id);
      tekst(wiersz.querySelector('.dn-meta'), godzinaZapisu(wersja.createdAt));
      const podpis = [...wiersz.childNodes].find(
        (dziecko) => dziecko.nodeType === Node.TEXT_NODE && (dziecko.nodeValue ?? '').trim() !== '',
      );
      if (podpis !== undefined) podpis.nodeValue = ` ${wersja.summary ?? ''} `;
      lista.appendChild(wiersz);
    }
  }

  async function podejrzyj(idWersji: string): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentVersionGet, {
      agentId: agent.id,
      versionId: idWersji,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Historia wersji', wynik.blad?.message ?? 'Rdzeń nie oddał wersji.', 'blad');
      return;
    }
    const wersja = wynik.wynik.version;
    const autor = wersja.author ?? 'bez autora';
    oglos(wersja.label ?? wersja.id, `${autor} · ${wersja.checksum ?? 'bez sumy kontrolnej'}`);
  }

  async function przywroc(idWersji: string): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentVersionRestore, {
      agentId: agent.id,
      versionId: idWersji,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Historia wersji', wynik.blad?.message ?? 'Rdzeń odrzucił przywrócenie.', 'blad');
      return;
    }
    stan.postaw(wynik.wynik.agent);
    stan.odswiez();
  }

  // Kontrakt nie zna komendy powielenia eksperta: nowy powstaje założeniem
  // definicji z pól, które `agent.create` przyjmuje.
  async function powiel(): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentCreate, {
      name: `${agent.name} (kopia)`,
      description: agent.description,
      systemPrompt: agent.systemPrompt,
      channelId: agent.channelId,
      model: agent.model,
      displayName: agent.displayName,
      favicon: agent.favicon,
      mode: agent.mode,
      memoryLevels: agent.memoryLevels,
      visibility: agent.visibility,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Nowy ekspert', wynik.blad?.message ?? 'Rdzeń odrzucił założenie.', 'blad');
      return;
    }
    otworz(wynik.wynik.agent);
    stan.odswiez();
  }

  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-wersja]');
    if (wiersz?.dataset.wersja !== undefined) void podejrzyj(wiersz.dataset.wersja);
  }, stan.przy);

  lista.addEventListener('contextmenu', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-wersja]');
    if (wiersz?.dataset.wersja === undefined) return;
    zdarzenie.preventDefault();
    void przywroc(wiersz.dataset.wersja);
  }, stan.przy);

  wiecej?.addEventListener('click', () => {
    void powiel();
  }, stan.przy);

  return { pokaz };
}
