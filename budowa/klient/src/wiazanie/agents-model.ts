// Panel Model Configuration: kanały bazowe z rdzenia, droga wywołania
// i sprawdzenie dostępności kanału wskazanego ekspertowi.

import {
  Command,
  ProviderTransport,
  type Agent,
  type Channel,
} from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { tekst, wzorem, type StanowiskoEkspertow } from './agents-wspolne.ts';

const PODPISY_DROGI: Record<string, string> = {
  api: 'API',
  cli: 'CLI',
  sdk: 'SDK',
  ssh: 'SSH',
  local: 'LOKALNIE',
};

export interface PanelModelu {
  pokaz: (agent: Agent | null) => void;
}

export function zwiazModel(stan: StanowiskoEkspertow): PanelModelu | null {
  const panelMozeMoze = stan.korzen.querySelector('#panel-model');
  if (panelMozeMoze === null) return null;
  const panelMoze = panelMozeMoze;
  const panel = panelMoze;

  const modele = panel.querySelector('.pk-modele');
  const pigulka = panel.querySelector('.pk-pig');
  const opisKanalu = panel.querySelector('#kanal-opis');
  const kluczPola = panel.querySelector('#klucz-pole');
  const maskuj = panel.querySelector('[data-maskuj]');
  const testuj = panel.querySelector('[data-testuj]');
  const wynikTestu = panel.querySelector('#test-wynik');

  const wzorModelu = wzorem(modele, '.pk-model');
  const wzorDrogi = wzorem(pigulka, '[data-kanal]');
  const adnotacja = panel.querySelector('.ab-adnotacja');
  wzorem(adnotacja, 'svg');

  // Rdzeń oddaje sam stan poświadczenia — ustawione czy brak i kiedy zmienione,
  // nigdy treść klucza; pole jawne z prototypu nie ma więc czego pokazać.
  tekst(kluczPola, '');
  maskuj?.setAttribute('aria-disabled', 'true');
  tekst(wynikTestu, '');

  let kanaly: Channel[] = [];

  function wypelnijDrogi(agent: Agent | null): void {
    if (pigulka === null || wzorDrogi === null) return;
    for (const stary of [...pigulka.querySelectorAll('[data-kanal]')]) stary.remove();
    for (const droga of Object.values(ProviderTransport)) {
      const przycisk = wzorDrogi.cloneNode(true) as HTMLElement;
      przycisk.dataset.droga = droga;
      przycisk.textContent = PODPISY_DROGI[droga] ?? droga;
      przycisk.setAttribute('aria-pressed', String(agent?.transport === droga));
      pigulka.appendChild(przycisk);
    }
  }

  function wypelnijModele(agent: Agent | null): void {
    if (modele === null || wzorModelu === null) return;
    for (const stary of [...modele.querySelectorAll('.pk-model')]) stary.remove();
    for (const kanal of kanaly) {
      const przycisk = wzorModelu.cloneNode(true) as HTMLElement;
      przycisk.dataset.kanal = kanal.id;
      tekst(przycisk.querySelector('b'), kanal.model ?? kanal.name);
      tekst(przycisk.querySelector('small'), kanal.kind);
      przycisk.setAttribute('aria-pressed', String(agent?.channelId === kanal.id));
      modele.appendChild(przycisk);
    }
  }

  function opiszKanal(agent: Agent | null): void {
    const wskazany = kanaly.find((kanal) => kanal.id === agent?.channelId);
    if (wskazany === undefined) {
      tekst(opisKanalu, 'Ekspert nie ma wskazanego kanału bazowego.');
      return;
    }
    const stanKanalu = wskazany.enabled ? 'włączony' : 'wyłączony';
    tekst(opisKanalu, `${wskazany.name} · ${wskazany.kind} · ${stanKanalu}`);
  }

  async function opiszPoswiadczenie(agent: Agent | null): Promise<void> {
    if (agent?.channelId === undefined) {
      tekst(kluczPola, 'Ekspert nie ma wskazanego kanału bazowego.');
      return;
    }
    const wynik = await wywolaj(stan.kanal, Command.ChannelCredentialStatus, {
      channelId: agent.channelId,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tekst(kluczPola, wynik.blad?.message ?? 'Rdzeń nie oddał stanu poświadczenia.');
      return;
    }
    const poswiadczenie = wynik.wynik.status;
    const rodzaj = poswiadczenie.kind === undefined ? '' : ` · ${poswiadczenie.kind}`;
    const zarzadca = poswiadczenie.managedBy === undefined ? '' : ` · ${poswiadczenie.managedBy}`;
    const zmiana = poswiadczenie.updatedAt === undefined
      ? ''
      : ` · zmienione ${new Date(poswiadczenie.updatedAt).toLocaleString('pl-PL')}`;
    tekst(kluczPola, `${poswiadczenie.present ? 'ustawione' : 'brak'}${rodzaj}${zarzadca}${zmiana}`);
  }

  async function opiszMozliwosci(agent: Agent | null): Promise<void> {
    if (adnotacja === null) return;
    if (agent?.channelId === undefined) {
      tekst(adnotacja, '');
      return;
    }
    const wynik = await wywolaj(stan.kanal, Command.ConfigCapabilitiesGet, {
      channelId: agent.channelId,
      transport: agent.transport,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tekst(adnotacja, wynik.blad?.message ?? 'Rdzeń nie oddał możliwości kanału.');
      return;
    }
    const pola = wynik.wynik.capabilities.fields;
    tekst(adnotacja, `Kanał przyjmuje ${pola.length} pól konfiguracji wywołania.`);
  }

  function pokaz(agent: Agent | null): void {
    wypelnijModele(agent);
    wypelnijDrogi(agent);
    opiszKanal(agent);
    tekst(wynikTestu, '');
    void opiszPoswiadczenie(agent);
    void opiszMozliwosci(agent);
  }

  async function wczytajKanaly(): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.ChannelList, { enabledOnly: false });
    if (!wynik.udany || wynik.wynik === undefined) return;
    kanaly = wynik.wynik.channels;
    pokaz(stan.ekspert());
  }

  async function przypisz(idKanalu: string, droga?: string): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wskazany = kanaly.find((kanal) => kanal.id === idKanalu);
    const wynik = await wywolaj(stan.kanal, Command.AgentModelSet, {
      agentId: agent.id,
      channelId: idKanalu,
      model: wskazany?.model,
      transport: (droga ?? agent.transport) as Agent['transport'],
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Model bazowy', wynik.blad?.message ?? 'Rdzeń odrzucił wskazanie kanału.', 'blad');
      return;
    }
    stan.postaw(wynik.wynik.agent);
  }

  modele?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('.pk-model');
    if (przycisk?.dataset.kanal === undefined) return;
    void przypisz(przycisk.dataset.kanal);
  }, stan.przy);

  pigulka?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-kanal]');
    const droga = przycisk?.dataset.droga;
    const agent = stan.ekspert();
    if (droga === undefined || agent?.channelId === undefined) return;
    void przypisz(agent.channelId, droga);
  }, stan.przy);

  testuj?.addEventListener('click', () => {
    const agent = stan.ekspert();
    if (agent?.channelId === undefined) {
      tekst(wynikTestu, 'Ekspert nie ma wskazanego kanału bazowego.');
      return;
    }
    void (async () => {
      const wynik = await wywolaj(stan.kanal, Command.ChannelCheck, { channelId: agent.channelId ?? '' });
      if (!wynik.udany || wynik.wynik === undefined) {
        tekst(wynikTestu, wynik.blad?.message ?? 'Rdzeń nie oddał wyniku sprawdzenia.');
        return;
      }
      const opoznienie = wynik.wynik.latencyMs === undefined ? '' : ` · ${wynik.wynik.latencyMs} ms`;
      const szczegol = wynik.wynik.detail === undefined ? '' : ` · ${wynik.wynik.detail}`;
      tekst(wynikTestu, `${wynik.wynik.reachable ? 'osiągalny' : 'nieosiągalny'}${opoznienie}${szczegol}`);
    })();
  }, stan.przy);

  void wczytajKanaly();
  return { pokaz };
}
