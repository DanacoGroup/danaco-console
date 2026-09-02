// Panele uniwersalne okna Agents: zespoły ekspertów, podagenci zlecani w tle
// oraz wywołanie konfiguracji utrwalone dla okna.

import {
  Command,
  type Agent,
  type Subagent,
  type Team,
} from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  BEZ_EKSPERTA,
  BEZ_POKRYCIA,
  nazwaEksperta,
  pustka,
  tekst,
  wzorem,
  zdejmijPustke,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

export interface PanelePomocnicze {
  pokaz: (agent: Agent | null) => void;
  odswiezPodagentow: () => void;
}

export function zwiazZespoly(stan: StanowiskoEkspertow): PanelePomocnicze | null {
  const planMoze = stan.korzen.querySelector('#panel-plan .sta-okno-tresc');
  const zadaniaMoze = stan.korzen.querySelector('#panel-zadania .sta-okno-tresc');
  if (planMoze === null || zadaniaMoze === null) return null;
  const plan = planMoze;
  const zadania = zadaniaMoze;
  const pliki = stan.korzen.querySelector('#panel-pliki .sta-okno-tresc');
  const artefakty = stan.korzen.querySelector('#panel-artefakty .sta-okno-tresc');

  const etykietaPlanu = plan.querySelector('.pt-etykieta');
  const wzorZespolu = wzorem(plan, '.dn-wykaz-modulu-poz');
  const wzorZadania = wzorem(zadania, '.dn-wykaz-modulu-poz');
  zadania.querySelector('.dn-postep')?.remove();
  zadania.querySelector('.dn-meta')?.remove();
  const wzorPliku = wzorem(pliki, '.dn-wykaz-modulu-poz');
  const wzorArtefaktu = wzorem(artefakty, '.dn-wykaz-modulu-poz');

  // Rdzeń nie wiąże plików roboczych z definicją eksperta żadną komendą kontraktu.
  if (wzorPliku !== null) pustka(pliki, BEZ_POKRYCIA);

  let zespoly: Team[] = [];

  async function wczytajZespoly(): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.TeamList, {});
    zespoly = wynik.udany ? wynik.wynik?.teams ?? [] : [];
    for (const stary of [...plan.querySelectorAll('.dn-wykaz-modulu-poz')]) stary.remove();
    zdejmijPustke(plan);
    tekst(etykietaPlanu, 'Zespoły ekspertów');
    if (zespoly.length === 0 || wzorZespolu === null) {
      pustka(plan, 'Rdzeń nie oddał żadnego zespołu.');
      return;
    }
    for (const zespol of zespoly) {
      const wiersz = wzorZespolu.cloneNode(true) as HTMLElement;
      wiersz.dataset.zespol = zespol.id;
      tekst(wiersz.querySelector('.dn-meta'), `${zespol.agentIds.length}`);
      const podpis = [...wiersz.childNodes].find(
        (dziecko) => dziecko.nodeType === Node.TEXT_NODE && (dziecko.nodeValue ?? '').trim() !== '',
      );
      if (podpis !== undefined) podpis.nodeValue = ` ${zespol.name} `;
      plan.appendChild(wiersz);
    }
  }

  async function wczytajPodagentow(): Promise<void> {
    for (const stary of [...zadania.querySelectorAll('.dn-wykaz-modulu-poz')]) stary.remove();
    zdejmijPustke(zadania);
    if (stan.idOkna === '') {
      pustka(zadania, 'Okno nie stoi jeszcze w rdzeniu — podagenci nie mają gdzie pracować.');
      return;
    }
    const wynik = await wywolaj(stan.kanal, Command.SubagentList, { windowId: stan.idOkna });
    const podagenci = wynik.udany ? wynik.wynik?.subagents ?? [] : [];
    if (podagenci.length === 0 || wzorZadania === null) {
      pustka(zadania, 'Rdzeń nie oddał żadnego podagenta dla tego okna.');
      return;
    }
    for (const podagent of podagenci) {
      wstawPodagenta(podagent);
    }
  }

  function wstawPodagenta(podagent: Subagent): void {
    if (wzorZadania === null) return;
    const wiersz = wzorZadania.cloneNode(true) as HTMLElement;
    wiersz.dataset.podagent = podagent.id;
    tekst(wiersz.querySelector('b'), podagent.name ?? podagent.id);
    const podpis = [...wiersz.childNodes].find(
      (dziecko) => dziecko.nodeType === Node.TEXT_NODE && (dziecko.nodeValue ?? '').trim() !== '',
    );
    if (podpis !== undefined) podpis.nodeValue = ` ${podagent.status} `;
    wiersz.title = podagent.result ?? podagent.task;
    zadania.appendChild(wiersz);
  }

  async function wczytajWywolanie(): Promise<void> {
    if (artefakty === null) return;
    for (const stary of [...artefakty.querySelectorAll('.dn-wykaz-modulu-poz')]) stary.remove();
    zdejmijPustke(artefakty);
    if (stan.idOkna === '' || wzorArtefaktu === null) {
      pustka(artefakty, 'Okno nie stoi jeszcze w rdzeniu — wywołanie nie jest utrwalone.');
      return;
    }
    const wynik = await wywolaj(stan.kanal, Command.ConfigExplainGet, { windowId: stan.idOkna });
    if (!wynik.udany || wynik.wynik === undefined) {
      pustka(artefakty, wynik.blad?.message ?? 'Rdzeń nie oddał wywołania dla tego okna.');
      return;
    }
    const wiersze = [
      wynik.wynik.argv.join(' '),
      `suma konstytucji: ${wynik.wynik.constitutionHash}`,
    ];
    for (const opis of wiersze) {
      const wiersz = wzorArtefaktu.cloneNode(true) as HTMLElement;
      wiersz.textContent = opis;
      artefakty.appendChild(wiersz);
    }
  }

  async function wczytajZespol(idZespolu: string): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.TeamLoad, { teamId: idZespolu });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Zespoły', wynik.blad?.message ?? 'Rdzeń nie oddał zespołu.', 'blad');
      return;
    }
    oglos(wynik.wynik.team.name, `${wynik.wynik.team.agentIds.length} ekspertów w zespole`);
  }

  async function powielZespol(idZespolu: string): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.TeamDuplicate, { teamId: idZespolu });
    if (!wynik.udany) oglos('Zespoły', wynik.blad?.message ?? 'Rdzeń odrzucił powielenie.', 'blad');
    await wczytajZespoly();
  }

  async function zapiszZespol(): Promise<void> {
    const agent = stan.ekspert();
    if (agent === null) {
      oglos('Zespoły', BEZ_EKSPERTA, 'ostrzezenie');
      return;
    }
    const wynik = await wywolaj(stan.kanal, Command.TeamSave, {
      name: `Zespół ${nazwaEksperta(agent)}`,
      description: agent.description,
      agentIds: [agent.id],
    });
    if (!wynik.udany) oglos('Zespoły', wynik.blad?.message ?? 'Rdzeń odrzucił zapis.', 'blad');
    await wczytajZespoly();
  }

  async function zlecPodagenta(): Promise<void> {
    const agent = stan.ekspert();
    if (agent === null || stan.idOkna === '') return;
    const wynik = await wywolaj(stan.kanal, Command.SubagentSpawn, {
      windowId: stan.idOkna,
      task: agent.description ?? nazwaEksperta(agent),
      name: nazwaEksperta(agent),
      count: 1,
      modelChannelId: agent.channelId,
    });
    if (!wynik.udany) oglos('Podagenci', wynik.blad?.message ?? 'Rdzeń odrzucił zlecenie.', 'blad');
    await wczytajPodagentow();
  }

  async function zbierzWynik(idPodagenta: string): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.SubagentResultCollect, {
      windowId: stan.idOkna,
      subagentIds: [idPodagenta],
      waitForAll: false,
    });
    if (!wynik.udany || wynik.wynik === undefined) return;
    const zebrany = wynik.wynik.subagents.find((podagent) => podagent.id === idPodagenta);
    if (zebrany !== undefined) oglos(zebrany.name ?? zebrany.id, zebrany.result ?? zebrany.status);
    await wczytajPodagentow();
  }

  async function zatrzymaj(idPodagenta: string): Promise<void> {
    await wywolaj(stan.kanal, Command.SubagentStop, {
      windowId: stan.idOkna,
      subagentIds: [idPodagenta],
    });
    await wczytajPodagentow();
  }

  plan.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-zespol]');
    if (wiersz?.dataset.zespol !== undefined) {
      void wczytajZespol(wiersz.dataset.zespol);
      return;
    }
    if (cel.closest('.pt-etykieta') !== null) void zapiszZespol();
  }, stan.przy);

  plan.addEventListener('contextmenu', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-zespol]');
    if (wiersz?.dataset.zespol === undefined) return;
    zdarzenie.preventDefault();
    void powielZespol(wiersz.dataset.zespol);
  }, stan.przy);

  zadania.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-podagent]');
    if (wiersz?.dataset.podagent === undefined) {
      void zlecPodagenta();
      return;
    }
    void zbierzWynik(wiersz.dataset.podagent);
  }, stan.przy);

  zadania.addEventListener('contextmenu', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-podagent]');
    if (wiersz?.dataset.podagent === undefined) return;
    zdarzenie.preventDefault();
    void zatrzymaj(wiersz.dataset.podagent);
  }, stan.przy);

  void wczytajZespoly();
  void wczytajPodagentow();
  void wczytajWywolanie();

  return {
    pokaz: () => {
      void wczytajZespoly();
    },
    odswiezPodagentow: () => {
      void wczytajPodagentow();
    },
  };
}
