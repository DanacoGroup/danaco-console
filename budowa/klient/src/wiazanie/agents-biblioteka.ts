// Biblioteka ekspertów: siatka kart z rdzenia, filtry wykazu oraz archiwum.
// Karta niesie definicję eksperta, jego wersję, zasięg i liczbę przypisań.

import {
  AgentVisibility,
  Command,
  type Agent,
  type AgentAssignment,
} from '../../../shared/contract.ts';
import { oglos } from './ogloszenie.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  BEZ_POKRYCIA,
  nazwaEksperta,
  pustka,
  tekst,
  wzorem,
  zdejmijPustke,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

type StanWykazu = 'aktywni' | 'wszyscy' | 'archiwum';

const STANY: readonly StanWykazu[] = ['aktywni', 'wszyscy', 'archiwum'];
const PODPISY_STANU: Record<StanWykazu, string> = {
  aktywni: 'Stan: aktywni',
  wszyscy: 'Stan: wszyscy',
  archiwum: 'Stan: archiwum',
};

type Zasieg = '' | typeof AgentVisibility.Global | typeof AgentVisibility.Project;

const ZASIEGI: readonly Zasieg[] = ['', AgentVisibility.Global, AgentVisibility.Project];
const PODPISY_ZASIEGU: Record<Zasieg, string> = {
  '': 'Widoczność: dowolna',
  global: 'Widoczność: globalna',
  project: 'Widoczność: projektowa',
};

export interface Biblioteka {
  odswiez: () => void;
}

export function zwiazBiblioteke(
  stan: StanowiskoEkspertow,
  otworz: (agent: Agent) => void,
): Biblioteka | null {
  const widok = stan.korzen.querySelector('#widok-biblioteka');
  const siatkaMoze = widok?.querySelector('.ab-siatka') ?? null;
  if (widok === null || siatkaMoze === null) return null;
  const siatka = siatkaMoze;

  const wzor = wzorem(siatka, '.ab-karta');
  const szukajka = widok.querySelector('.ab-filtry input');
  const znaki = [...widok.querySelectorAll<HTMLElement>('.ab-filtry .sta-chip')];
  const [znakSrodowiska, znakZasiegu, znakStanu] = znaki;

  let stanWykazu: StanWykazu = 'aktywni';
  let zasieg: Zasieg = '';
  let szukane = '';
  let przypisania = new Map<string, AgentAssignment[]>();

  // Rdzeń nie zna filtru środowiska dla eksperta: wykaz zawęża projekt, nie środowisko.
  if (znakSrodowiska !== undefined) {
    znakSrodowiska.setAttribute('aria-disabled', 'true');
    znakSrodowiska.title = 'Rdzeń nie zawęża wykazu ekspertów środowiskiem.';
  }
  tekst(znakZasiegu ?? null, PODPISY_ZASIEGU[zasieg]);
  tekst(znakStanu ?? null, PODPISY_STANU[stanWykazu]);

  async function pobierzPrzypisania(): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.AgentAssignmentList, {});
    przypisania = new Map<string, AgentAssignment[]>();
    if (!wynik.udany || wynik.wynik === undefined) return;
    for (const przypisanie of wynik.wynik.assignments) {
      const zebrane = przypisania.get(przypisanie.agentId) ?? [];
      zebrane.push(przypisanie);
      przypisania.set(przypisanie.agentId, zebrane);
    }
  }

  async function pobierzEkspertow(): Promise<Agent[]> {
    if (stanWykazu === 'archiwum') {
      const wynik = await wywolaj(stan.kanal, Command.AgentArchiveList, { limit: 200 });
      return wynik.udany && wynik.wynik !== undefined ? wynik.wynik.agents : [];
    }
    const wynik = await wywolaj(stan.kanal, Command.AgentList, {
      query: szukane === '' ? undefined : szukane,
      enabledOnly: stanWykazu === 'aktywni',
      limit: 200,
    });
    return wynik.udany && wynik.wynik !== undefined ? wynik.wynik.agents : [];
  }

  function wypelnijKarte(karta: HTMLElement, agent: Agent): void {
    karta.dataset.ekspert = agent.id;
    tekst(karta.querySelector('.ab-karta-nazwa'), nazwaEksperta(agent));
    tekst(karta.querySelector('.ab-karta-opis'), agent.description ?? '');
    const opisy = [...karta.querySelectorAll<HTMLElement>('.ab-karta-meta')];
    const kanal = agent.transport ?? '';
    tekst(opisy[0] ?? null, [agent.model ?? '', kanal].filter((czesc) => czesc !== '').join(' · '));
    const zebrane = przypisania.get(agent.id) ?? [];
    tekst(opisy[1] ?? null, zebrane.length === 0 ? 'bez przypisań' : `${zebrane.length} przypisań`);
    const plakietki = karta.querySelector('.ab-karta-plak');
    if (plakietki !== null) {
      const wzorPlakietki = wzorem(plakietki, '.dn-plakietka');
      const opisyPlakietek = [
        agent.version === undefined ? '' : `w.${agent.version}`,
        agent.visibility === AgentVisibility.Global ? 'globalny' : 'projektowy',
        agent.enabled ? 'aktywny' : 'wyłączony',
      ].filter((opis) => opis !== '');
      for (const opis of opisyPlakietek) {
        if (wzorPlakietki === null) break;
        const plakietka = wzorPlakietki.cloneNode(true) as HTMLElement;
        plakietka.textContent = opis;
        plakietki.appendChild(plakietka);
      }
    }
    const czynnosci = [...karta.querySelectorAll<HTMLElement>('.ab-karta-akcje > *')];
    const archiwum = stanWykazu === 'archiwum';
    tekst(czynnosci[0] ?? null, archiwum ? 'Przywróć' : 'Edytuj');
    tekst(czynnosci[1] ?? null, archiwum ? 'Usuń' : 'Archiwizuj');
    czynnosci[0]?.setAttribute('data-czynnosc', archiwum ? 'przywroc' : 'edytuj');
    czynnosci[1]?.setAttribute('data-czynnosc', archiwum ? 'usun' : 'archiwizuj');
  }

  async function odswiez(): Promise<void> {
    if (wzor === null) return;
    await pobierzPrzypisania();
    const eksperci = await pobierzEkspertow();
    for (const karta of [...siatka.querySelectorAll('.ab-karta')]) karta.remove();
    zdejmijPustke(siatka);
    const widoczni = eksperci.filter((agent) => {
      if (zasieg !== '' && agent.visibility !== zasieg) return false;
      if (stanWykazu !== 'archiwum' || szukane === '') return true;
      const opis = `${agent.name} ${agent.description ?? ''}`.toLowerCase();
      return opis.includes(szukane.toLowerCase());
    });
    if (widoczni.length === 0) {
      pustka(siatka, 'Rdzeń nie oddał żadnego eksperta dla tego zawężenia.');
      return;
    }
    for (const agent of widoczni) {
      const karta = wzor.cloneNode(true) as HTMLElement;
      wypelnijKarte(karta, agent);
      siatka.appendChild(karta);
    }
  }

  const przeladuj = (): void => {
    void odswiez();
  };

  async function otworzWskazanego(id: string): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.AgentList, { limit: 200 });
    if (!wynik.udany || wynik.wynik === undefined) return;
    const znaleziony = wynik.wynik.agents.find((agent) => agent.id === id);
    if (znaleziony !== undefined) otworz(znaleziony);
  }

  async function wykonajCzynnosc(czynnosc: string, id: string): Promise<void> {
    if (czynnosc === 'edytuj') {
      await otworzWskazanego(id);
      return;
    }
    if (czynnosc === 'archiwizuj') {
      const wynik = await wywolaj(stan.kanal, Command.AgentArchive, { agentId: id });
      if (!wynik.udany) oglos('Archiwizacja', wynik.blad?.message ?? BEZ_POKRYCIA, 'blad');
      przeladuj();
      return;
    }
    if (czynnosc === 'przywroc') {
      const wynik = await wywolaj(stan.kanal, Command.AgentRestore, { agentId: id });
      if (!wynik.udany) oglos('Przywrócenie', wynik.blad?.message ?? BEZ_POKRYCIA, 'blad');
      przeladuj();
      return;
    }
    if (czynnosc === 'usun') {
      const wynik = await wywolaj(stan.kanal, Command.AgentDelete, { agentId: id });
      if (!wynik.udany) oglos('Usunięcie', wynik.blad?.message ?? BEZ_POKRYCIA, 'blad');
      przeladuj();
    }
  }

  siatka.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const kartaMozeMoze = cel.closest<HTMLElement>('.ab-karta');
    if (kartaMozeMoze === null) return;
    const kartaMoze = kartaMozeMoze;
    const karta = kartaMoze;
    const id = karta.dataset.ekspert ?? '';
    if (id === '') return;
    const czynnosc = cel.closest<HTMLElement>('[data-czynnosc]')?.dataset.czynnosc;
    void wykonajCzynnosc(czynnosc ?? (stanWykazu === 'archiwum' ? '' : 'edytuj'), id);
  }, stan.przy);

  szukajka?.addEventListener('input', (zdarzenie) => {
    const zrodlo = zdarzenie.target;
    szukane = zrodlo instanceof HTMLInputElement ? zrodlo.value.trim() : '';
    przeladuj();
  }, stan.przy);

  znakZasiegu?.addEventListener('click', () => {
    const kolejny = ZASIEGI[(ZASIEGI.indexOf(zasieg) + 1) % ZASIEGI.length];
    zasieg = kolejny ?? '';
    tekst(znakZasiegu, PODPISY_ZASIEGU[zasieg]);
    przeladuj();
  }, stan.przy);

  znakStanu?.addEventListener('click', () => {
    const kolejny = STANY[(STANY.indexOf(stanWykazu) + 1) % STANY.length];
    stanWykazu = kolejny ?? 'aktywni';
    tekst(znakStanu, PODPISY_STANU[stanWykazu]);
    przeladuj();
  }, stan.przy);

  return { odswiez: przeladuj };
}
