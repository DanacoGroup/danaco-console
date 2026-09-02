// Panel Connectors Manager: wtyczki i konektory eksperta wraz z ich rodzajem,
// punktem dostępu i stanem włączenia.

import {
  AccessPointKind,
  AgentConnectorKind,
  Command,
  type AccessPoint,
  type Agent,
  type AgentConnector,
  type AgentPlugin,
  type Extension,
} from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  BEZ_EKSPERTA,
  czyWlaczony,
  pole,
  przelacz,
  pustka,
  tekst,
  wzorem,
  zdejmijPustke,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

const FILTRY: readonly string[] = ['wszystkie', 'plugin', 'api', 'mcp'];

export interface PanelKonektorow {
  pokaz: (agent: Agent | null) => void;
}

export function zwiazKonektory(stan: StanowiskoEkspertow): PanelKonektorow | null {
  const panel = stan.korzen.querySelector('#panel-connectors');
  const trescMoze = panel?.querySelector('.sta-okno-tresc') ?? null;
  if (panel === null || trescMoze === null) return null;
  const tresc = trescMoze;

  const wzor = wzorem(tresc, '.pk-wiersz');
  const szukajka = pole(panel, '.dn-szukaj input');
  const dodaj = panel.querySelector('.sta-okno-tresc > button');
  const baner = panel.querySelector('.pk-baner');
  const znacznik = panel.querySelector('.sta-okno-znacznik');
  const przyciskiFiltru = [...panel.querySelectorAll<HTMLElement>('[data-rodzaj]')];

  for (const [numer, przycisk] of przyciskiFiltru.entries()) {
    przycisk.dataset.filtr = FILTRY[numer] ?? 'wszystkie';
  }

  let filtr = 'wszystkie';
  let konektory: AgentConnector[] = [];
  let wtyczki: AgentPlugin[] = [];
  let katalog: Extension[] = [];
  let mosty: AccessPoint[] = [];

  function wypelnijWiersz(
    identyfikator: string,
    rodzaj: string,
    nazwa: string,
    opis: string,
    wlaczony: boolean,
    wtyczka: boolean,
  ): void {
    if (wzor === null) return;
    const wiersz = wzor.cloneNode(true) as HTMLElement;
    wiersz.dataset.rozszerzenie = identyfikator;
    wiersz.dataset.rodzajWpisu = wtyczka ? 'plugin' : rodzaj;
    tekst(wiersz.querySelector('.nazwa'), nazwa);
    tekst(wiersz.querySelector('.tresc small'), opis);
    const przelacznik = wiersz.querySelector('.dn-przelacznik');
    przelacz(przelacznik, wlaczony);
    przelacznik?.setAttribute('aria-label', nazwa);
    // Wtyczka nie ma w kontrakcie komendy zmiany stanu — jest tylko dopięcie i zdjęcie.
    if (wtyczka) przelacznik?.setAttribute('aria-disabled', 'true');
    tresc.insertBefore(wiersz, baner ?? dodaj);
  }

  function pokaz(agent: Agent | null): void {
    for (const stary of [...tresc.querySelectorAll('.pk-wiersz')]) stary.remove();
    zdejmijPustke(tresc);
    for (const przycisk of przyciskiFiltru) {
      przycisk.setAttribute('aria-pressed', String(przycisk.dataset.filtr === filtr));
    }
    if (agent === null) {
      tekst(znacznik, '');
      tekst(baner, BEZ_EKSPERTA);
      return;
    }
    const widoczneKonektory = konektory.filter(
      (konektor) => filtr === 'wszystkie' || konektor.kind === filtr,
    );
    const widoczneWtyczki = filtr === 'wszystkie' || filtr === 'plugin' ? wtyczki : [];
    tekst(znacznik, `${konektory.length + wtyczki.length} podłączonych`);
    for (const konektor of widoczneKonektory) {
      const punkt = konektor.accessPointId ?? 'bez punktu dostępu';
      wypelnijWiersz(konektor.id, konektor.kind, konektor.name, `${konektor.kind} · ${punkt}`,
        konektor.enabled, false);
    }
    for (const wtyczka of widoczneWtyczki) {
      const zrodlo = wtyczka.source ?? 'bez źródła';
      const wersja = wtyczka.version === undefined ? '' : ` · ${wtyczka.version}`;
      wypelnijWiersz(wtyczka.id, 'plugin', wtyczka.name, `${zrodlo}${wersja}`, wtyczka.enabled, true);
    }
    for (const pozycja of katalog) {
      const stanPozycji = pozycja.installed ? 'zainstalowane' : 'dostępne w katalogu';
      const punkt = pozycja.accessPointId ?? 'bez punktu dostępu';
      const wiersz = wzor === null ? null : (wzor.cloneNode(true) as HTMLElement);
      if (wiersz === null) break;
      wiersz.dataset.rozszerzenie = pozycja.id;
      wiersz.dataset.rodzajWpisu = 'katalog';
      tekst(wiersz.querySelector('.nazwa'), pozycja.name);
      tekst(wiersz.querySelector('.tresc small'), `${pozycja.kind} · ${stanPozycji} · ${punkt}`);
      przelacz(wiersz.querySelector('.dn-przelacznik'), pozycja.installed && pozycja.enabled);
      tresc.insertBefore(wiersz, baner ?? dodaj);
    }
    if (widoczneKonektory.length + widoczneWtyczki.length + katalog.length === 0) {
      pustka(tresc, 'Ekspert nie ma podłączonego rozszerzenia tego rodzaju.');
    }
    tekst(baner, `Punkty dostępu z rdzenia: ${mosty.length}. Menu podręczne zdejmuje pozycję.`);
  }

  async function wczytaj(): Promise<void> {
    const agent = stan.ekspert();
    if (agent === null) {
      konektory = [];
      wtyczki = [];
      pokaz(null);
      return;
    }
    const wykazKonektorow = await wywolaj(stan.kanal, Command.AgentConnectorList, {
      agentId: agent.id,
    });
    const wykazWtyczek = await wywolaj(stan.kanal, Command.AgentPluginList, { agentId: agent.id });
    const wykazKatalogu = await wywolaj(stan.kanal, Command.ExtensionList, { agentId: agent.id });
    const wykazMostow = await wywolaj(stan.kanal, Command.AccessPointList, {
      kind: AccessPointKind.McpBridge,
      enabledOnly: true,
    });
    konektory = wykazKonektorow.udany ? wykazKonektorow.wynik?.connectors ?? [] : [];
    wtyczki = wykazWtyczek.udany ? wykazWtyczek.wynik?.plugins ?? [] : [];
    katalog = wykazKatalogu.udany ? wykazKatalogu.wynik?.extensions ?? [] : [];
    mosty = wykazMostow.udany ? wykazMostow.wynik?.points ?? [] : [];
    pokaz(agent);
  }

  async function przestawKatalog(wiersz: HTMLElement): Promise<void> {
    const identyfikator = wiersz.dataset.rozszerzenie ?? '';
    const pozycja = katalog.find((wpis) => wpis.id === identyfikator);
    if (pozycja === undefined) return;
    const wynik = pozycja.installed
      ? await wywolaj(stan.kanal, Command.ExtensionToggle, {
        extensionId: pozycja.id,
        enabled: !pozycja.enabled,
      })
      : await wywolaj(stan.kanal, Command.ExtensionInstall, {
        code: pozycja.code,
        kind: pozycja.kind,
      });
    if (!wynik.udany) {
      oglos('Katalog rozszerzeń', wynik.blad?.message ?? 'Rdzeń odrzucił zmianę.', 'blad');
      return;
    }
    await wczytaj();
  }

  // Punkt dostępu wskazuje most MCP: rdzeń oddaje ich wykaz, a konfiguracja
  // pozycji katalogu jest jedynym miejscem, w którym wskazanie się zapisuje.
  async function wskazMost(wiersz: HTMLElement): Promise<void> {
    const identyfikator = wiersz.dataset.rozszerzenie ?? '';
    const most = mosty[0];
    if (identyfikator === '' || most === undefined) return;
    const wynik = await wywolaj(stan.kanal, Command.ExtensionConfigure, {
      extensionId: identyfikator,
      config: { accessPointId: most.id },
    });
    if (!wynik.udany) {
      oglos('Katalog rozszerzeń', wynik.blad?.message ?? 'Rdzeń odrzucił konfigurację.', 'blad');
      return;
    }
    await wczytaj();
  }

  async function odinstaluj(wiersz: HTMLElement): Promise<void> {
    const identyfikator = wiersz.dataset.rozszerzenie ?? '';
    if (identyfikator === '') return;
    const wynik = await wywolaj(stan.kanal, Command.ExtensionUninstall, {
      extensionId: identyfikator,
    });
    if (!wynik.udany) {
      oglos('Katalog rozszerzeń', wynik.blad?.message ?? 'Rdzeń odrzucił odinstalowanie.', 'blad');
      return;
    }
    await wczytaj();
  }

  async function przestaw(wiersz: HTMLElement): Promise<void> {
    const agent = stan.ekspert();
    const identyfikator = wiersz.dataset.rozszerzenie ?? '';
    if (agent === null || identyfikator === '') return;
    if (wiersz.dataset.rodzajWpisu === 'plugin') return;
    const wlaczony = czyWlaczony(wiersz.querySelector('.dn-przelacznik'));
    const wynik = await wywolaj(stan.kanal, Command.AgentConnectorConfigure, {
      agentId: agent.id,
      connectorId: identyfikator,
      enabled: !wlaczony,
    });
    if (!wynik.udany) {
      oglos('Konektory', wynik.blad?.message ?? 'Rdzeń odrzucił zmianę stanu.', 'blad');
      return;
    }
    await wczytaj();
  }

  async function zdejmij(wiersz: HTMLElement): Promise<void> {
    const agent = stan.ekspert();
    const identyfikator = wiersz.dataset.rozszerzenie ?? '';
    if (agent === null || identyfikator === '') return;
    const wynik = wiersz.dataset.rodzajWpisu === 'plugin'
      ? await wywolaj(stan.kanal, Command.AgentPluginRemove, {
        agentId: agent.id,
        pluginId: identyfikator,
      })
      : await wywolaj(stan.kanal, Command.AgentConnectorRemove, {
        agentId: agent.id,
        connectorId: identyfikator,
      });
    if (!wynik.udany) {
      oglos('Rozszerzenia', wynik.blad?.message ?? 'Rdzeń odrzucił zdjęcie.', 'blad');
      return;
    }
    if (wynik.wynik !== undefined && 'agent' in wynik.wynik) stan.postaw(wynik.wynik.agent);
    await wczytaj();
  }

  async function podlacz(): Promise<void> {
    const agent = stan.ekspert();
    const nazwa = szukajka?.value.trim() ?? '';
    if (agent === null || nazwa === '') return;
    const wynik = filtr === 'plugin'
      ? await wywolaj(stan.kanal, Command.AgentPluginAdd, { agentId: agent.id, name: nazwa })
      : await wywolaj(stan.kanal, Command.AgentConnectorAdd, {
        agentId: agent.id,
        name: nazwa,
        kind: (filtr === 'mcp' ? AgentConnectorKind.Mcp : AgentConnectorKind.Api),
        accessPointId: filtr === 'mcp' ? mosty[0]?.id : undefined,
      });
    if (!wynik.udany) {
      oglos('Rozszerzenia', wynik.blad?.message ?? 'Rdzeń odrzucił podłączenie.', 'blad');
      return;
    }
    if (szukajka !== null) szukajka.value = '';
    await wczytaj();
  }

  tresc.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-rodzaj]');
    if (przycisk !== null) {
      filtr = przycisk.dataset.filtr ?? 'wszystkie';
      pokaz(stan.ekspert());
      return;
    }
    const wierszMozeMoze = cel.closest<HTMLElement>('.pk-wiersz');
    if (wierszMozeMoze === null) return;
    const wierszMoze = wierszMozeMoze;
    const wiersz = wierszMoze;
    if (wiersz.dataset.rodzajWpisu === 'katalog') void przestawKatalog(wiersz);
    else void przestaw(wiersz);
  }, stan.przy);

  tresc.addEventListener('dblclick', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('.pk-wiersz[data-rodzaj-wpisu="katalog"]');
    if (wiersz !== null) void wskazMost(wiersz);
  }, stan.przy);

  tresc.addEventListener('contextmenu', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wierszMozeMoze = cel.closest<HTMLElement>('.pk-wiersz');
    if (wierszMozeMoze === null) return;
    const wierszMoze = wierszMozeMoze;
    const wiersz = wierszMoze;
    zdarzenie.preventDefault();
    if (wiersz.dataset.rodzajWpisu === 'katalog') void odinstaluj(wiersz);
    else void zdejmij(wiersz);
  }, stan.przy);

  dodaj?.addEventListener('click', () => {
    void podlacz();
  }, stan.przy);

  szukajka?.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') void podlacz();
  }, stan.przy);

  return {
    pokaz: () => {
      void wczytaj();
    },
  };
}
