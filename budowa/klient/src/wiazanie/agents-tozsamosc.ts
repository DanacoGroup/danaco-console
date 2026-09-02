// Zakładka tożsamości eksperta: nazwa, przeznaczenie, instrukcje systemowe,
// moduły zastosowania, zasięg widoczności i poziomy pamięci.

import {
  AgentVisibility,
  Command,
  IdentityLayer,
  IdentityMode,
  MemoryLevel,
  type Agent,
  type IdentityCategory,
} from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  nazwaEksperta,
  obszar,
  pole,
  poleZPodpisem,
  wypelnijGrupe,
  wzorem,
  zaznaczoneKlucze,
  godzinaZapisu,
  tekst,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

const PODPISY_PAMIECI: Record<string, string> = {
  global: 'Globalna',
  project: 'Projekt',
  session: 'Sesja',
  environment: 'Środowisko',
};

const PODPISY_ZASIEGU: Record<string, string> = {
  global: 'Globalny',
  project: 'Projektowy',
};

export interface Tozsamosc {
  pokaz: (agent: Agent | null) => void;
}

export function zwiazTozsamosc(stan: StanowiskoEkspertow): Tozsamosc | null {
  const zakladkaMozeMoze = stan.korzen.querySelector('#tab-tozsamosc');
  if (zakladkaMozeMoze === null) return null;
  const zakladkaMoze = zakladkaMozeMoze;
  const zakladka = zakladkaMoze;

  const nazwa = pole(zakladka, '#ab-nazwa');
  const opis = obszar(zakladka, '#ab-opis');
  const instrukcje = obszar(zakladka, '#ab-instr');
  const grupaModulow = poleZPodpisem(zakladka, 'Moduły')?.querySelector('.ab-grupa') ?? null;
  const grupaZasiegu = poleZPodpisem(zakladka, 'Zasięg')?.querySelector('.ab-grupa') ?? null;
  const grupaPamieci = poleZPodpisem(zakladka, 'Konfiguracja pamięci')?.querySelector('.ab-grupa')
    ?? null;
  const nazwaInline = pole(stan.korzen, '.ab-nazwa-inline');
  const znacznikStanu = stan.korzen.querySelector('.ab-edy-naglowek .dn-plakietka');
  const znacznikZapisu = stan.korzen.querySelector('.ab-zapis');
  const stopka = stan.korzen.querySelector('.ab-stopka');
  const zapisz = stopka?.querySelector('[data-zapisz]') ?? null;
  const anuluj = stopka?.querySelector('button:not([data-zapisz])') ?? null;

  // Adnotacja prototypu opisuje przykładowy dobór modułów, więc schodzi z dokumentu.
  zakladka.querySelector('.ab-adnotacja')?.remove();
  wzorem(znacznikZapisu, '.dn-kropka');

  let kodyModulow: string[] = [];
  let kategoria: IdentityCategory | null = null;

  // Warstwę, do której trafiają instrukcje, wskazuje rejestr kategorii tożsamości;
  // bez niego klient sam rozstrzygałby o kształcie konstytucji eksperta.
  async function wczytajKategorie(): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.IdentityCategoryList, {});
    if (!wynik.udany || wynik.wynik === undefined) return;
    const czynne = wynik.wynik.categories.filter((wpis) => wpis.enabled);
    kategoria = czynne.find((wpis) => wpis.layer === IdentityLayer.Expertise) ?? czynne[0] ?? null;
    const podpis = zakladka.querySelector('label[for="ab-instr"]');
    if (kategoria !== null) tekst(podpis, `Instrukcje systemowe — ${kategoria.name}`);
  }

  async function wczytajModuly(): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.ModuleList, {});
    if (!wynik.udany || wynik.wynik === undefined) return;
    kodyModulow = wynik.wynik.modules.map((modul) => modul.code);
    const agent = stan.ekspert();
    wypelnijGrupe(
      grupaModulow,
      wynik.wynik.modules.map((modul) => ({
        klucz: modul.code,
        podpis: modul.name,
        zaznaczona: agent?.moduleCodes.includes(modul.code) === true,
      })),
    );
  }

  function pokaz(agent: Agent | null): void {
    if (nazwa !== null) nazwa.value = agent?.name ?? '';
    if (nazwaInline !== null) nazwaInline.value = agent === null ? '' : nazwaEksperta(agent);
    if (opis !== null) opis.value = agent?.description ?? '';
    if (instrukcje !== null) instrukcje.value = agent?.systemPrompt ?? '';
    tekst(znacznikStanu, agent === null ? '' : agent.enabled ? 'AKTYWNY' : 'WYŁĄCZONY');
    tekst(znacznikZapisu, agent === null ? '' : `zapisano o ${godzinaZapisu(agent.updatedAt)}`);
    wypelnijGrupe(
      grupaZasiegu,
      Object.values(AgentVisibility).map((wartosc) => ({
        klucz: wartosc,
        podpis: PODPISY_ZASIEGU[wartosc] ?? wartosc,
        zaznaczona: agent?.visibility === wartosc,
      })),
    );
    wypelnijGrupe(
      grupaPamieci,
      Object.values(MemoryLevel).map((wartosc) => ({
        klucz: wartosc,
        podpis: PODPISY_PAMIECI[wartosc] ?? wartosc,
        zaznaczona: agent?.memoryLevels.includes(wartosc) === true,
      })),
    );
    wypelnijGrupe(
      grupaModulow,
      kodyModulow.map((kod) => ({
        klucz: kod,
        podpis: kod,
        zaznaczona: agent?.moduleCodes.includes(kod) === true,
      })),
    );
  }

  // Tryb DOŁĄCZ dokłada instrukcje do konstytucji, więc ich treść idzie warstwą
  // wskazaną kategorią; tryb ZASTĄP oddaje ją wprost polem `systemPrompt`.
  async function zapiszWarstwe(agent: Agent, tresc: string): Promise<void> {
    const warstwa = kategoria?.layer ?? IdentityLayer.Expertise;
    if ((agent.mode ?? kategoria?.defaultMode) !== IdentityMode.DOLACZ) return;
    if (tresc === '') {
      await wywolaj(stan.kanal, Command.AgentLayerRemove, { agentId: agent.id, layer: warstwa });
      return;
    }
    await wywolaj(stan.kanal, Command.AgentLayerSet, {
      agentId: agent.id,
      layer: warstwa,
      content: tresc,
      enabled: true,
    });
  }

  async function zapiszEksperta(): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const zasiegi = zaznaczoneKlucze(grupaZasiegu);
    const wybranyZasieg = zasiegi[0];
    const wynik = await wywolaj(stan.kanal, Command.AgentUpdate, {
      agentId: agent.id,
      name: nazwa?.value.trim() ?? agent.name,
      description: opis?.value.trim() ?? agent.description,
      systemPrompt: instrukcje?.value ?? agent.systemPrompt,
      displayName: nazwaInline?.value.trim() ?? agent.displayName,
      memoryLevels: zaznaczoneKlucze(grupaPamieci) as Agent['memoryLevels'],
      visibility: (wybranyZasieg ?? agent.visibility) as Agent['visibility'],
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Zapis eksperta', wynik.blad?.message ?? 'Rdzeń odrzucił zapis.', 'blad');
      return;
    }
    const moduly = await wywolaj(stan.kanal, Command.AgentModulesSet, {
      agentId: agent.id,
      moduleCodes: zaznaczoneKlucze(grupaModulow),
    });
    await zapiszWarstwe(agent, instrukcje?.value.trim() ?? '');
    stan.postaw(moduly.udany && moduly.wynik !== undefined ? moduly.wynik.agent : wynik.wynik.agent);
    stan.odswiez();
  }

  async function przestawCzynnosc(): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentUpdate, {
      agentId: agent.id,
      enabled: !agent.enabled,
    });
    if (wynik.udany && wynik.wynik !== undefined) {
      stan.postaw(wynik.wynik.agent);
      stan.odswiez();
    }
  }

  zapisz?.addEventListener('click', () => {
    void zapiszEksperta();
  }, stan.przy);

  anuluj?.addEventListener('click', () => {
    pokaz(stan.ekspert());
  }, stan.przy);

  znacznikStanu?.addEventListener('click', () => {
    void przestawCzynnosc();
  }, stan.przy);

  void wczytajModuly();
  void wczytajKategorie();
  return { pokaz };
}
