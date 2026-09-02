// Panel Permissions Center: zakresy uprawnień, dostęp do modułów, izolacja
// techniczna, Subagent Network oraz polityka efektywna eksperta.

import {
  AgentPermissionGroup,
  Command,
  type Agent,
  type AgentPolicy,
  type IsolationTechnicalSwitch,
  type ToolScope,
} from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  BEZ_EKSPERTA,
  czyWlaczony,
  pole,
  przelacz,
  tekst,
  wypelnijGrupe,
  wzorem,
  zaznaczoneKlucze,
  type StanowiskoEkspertow,
} from './agents-wspolne.ts';

const PODPISY_GRUP: Record<string, string> = {
  files: 'Pliki i dokumenty',
  network: 'Dostęp sieciowy',
  processes: 'Procesy i wykonanie',
  integrations: 'Integracje i rozszerzenia',
  modules: 'Moduły i zasoby',
};

export interface PanelUprawnien {
  pokaz: (agent: Agent | null) => void;
}

export function zwiazUprawnienia(stan: StanowiskoEkspertow): PanelUprawnien | null {
  const panelMozeMoze = stan.korzen.querySelector('#panel-permissions');
  if (panelMozeMoze === null) return null;
  const panelMoze = panelMozeMoze;
  const panel = panelMoze;

  const grupy = [...panel.querySelectorAll<HTMLElement>('.pk-grupa[data-grupa]')];
  const tabela = panel.querySelector('#grupa-rozsz tbody');
  const grupaModulow = panel.querySelector('#grupa-moduly .ab-grupa');
  const macierz = panel.querySelector('#grupa-izolacja .pk-macierz');
  const przelacznikSieci = panel.querySelector('#grupa-mtai .pk-przelacznik, #grupa-mtai .dn-przelacznik');
  const liczbaPodagentow = pole(panel, '#grupa-mtai .pk-num');
  const kod = panel.querySelector('.pk-kod');
  const banerWiodacy = panel.querySelector('.pk-baner');
  const podsumowanieRozszerzen = panel.querySelector('#grupa-rozsz p');
  const podsumowanieIzolacji = panel.querySelector('#grupa-izolacja p');
  const podsumowanieModulow = panel.querySelector('#grupa-moduly p');

  const wzorWiersza = wzorem(tabela, 'tr');
  const wzorZakresu = wzorem(macierz, '.pk-mrz');

  let kodyModulow: string[] = [];
  let przelaczniki: IsolationTechnicalSwitch[] = [];
  let zakresyNarzedzi: ToolScope[] = [];

  // Nagłówki prototypu dzielą wiersz na odczyt, zapis i akcję; rdzeń trzyma
  // dostępność, wymóg potwierdzenia i granicę wywołań, więc podpisy idą za nim.
  const naglowki = [...panel.querySelectorAll('#grupa-rozsz thead th')];
  for (const [numer, podpis] of ['Zakres', 'dostępny', 'potwierdzenie', 'limit'].entries()) {
    tekst(naglowki[numer] ?? null, podpis);
  }

  async function wczytajModuly(): Promise<void> {
    const wynik = await wywolaj(stan.kanal, Command.ModuleList, {});
    if (!wynik.udany || wynik.wynik === undefined) return;
    kodyModulow = wynik.wynik.modules.map((modul) => modul.code);
  }

  function wypelnijZakresy(agent: Agent | null): void {
    if (tabela === null || wzorWiersza === null) return;
    for (const stary of [...tabela.querySelectorAll('tr')]) stary.remove();
    if (agent === null) return;
    let zawezone = 0;
    for (const grupa of Object.values(AgentPermissionGroup)) {
      const wpis = agent.permissions?.find((uprawnienie) => uprawnienie.group === grupa);
      const przyznane = wpis?.granted !== false;
      if (!przyznane) zawezone += 1;
      const wiersz = wzorWiersza.cloneNode(true) as HTMLElement;
      wiersz.dataset.grupa = grupa;
      tekst(wiersz.querySelector('td'), PODPISY_GRUP[grupa] ?? grupa);
      const zaznaczenia = [...wiersz.querySelectorAll<HTMLInputElement>('input')];
      const [odczyt, ...pozostale] = zaznaczenia;
      if (odczyt !== undefined) {
        odczyt.checked = przyznane;
        odczyt.setAttribute('aria-label', `${PODPISY_GRUP[grupa] ?? grupa} — przyznane`);
      }
      // Rdzeń trzyma jedną zgodę na grupę, nie rozdziela odczytu od zapisu i akcji.
      for (const zaznaczenie of pozostale) {
        zaznaczenie.checked = false;
        zaznaczenie.disabled = true;
        zaznaczenie.title = 'Rdzeń trzyma jedną zgodę na grupę zakresu.';
      }
      tabela.appendChild(wiersz);
    }
    for (const zakres of zakresyNarzedzi) {
      const wiersz = wzorWiersza.cloneNode(true) as HTMLElement;
      wiersz.dataset.narzedzie = zakres.toolName;
      wiersz.dataset.profil = zakres.profileId;
      tekst(wiersz.querySelector('td'), zakres.shortName);
      const zaznaczenia = [...wiersz.querySelectorAll<HTMLInputElement>('input')];
      const stany = [zakres.enabled, zakres.confirmRequired, zakres.callLimit > 0];
      for (const [numer, zaznaczenie] of zaznaczenia.entries()) {
        zaznaczenie.checked = stany[numer] === true;
        zaznaczenie.disabled = false;
        zaznaczenie.setAttribute('aria-label', `${zakres.shortName} — kolumna ${numer + 1}`);
      }
      tabela.appendChild(wiersz);
    }
    tekst(
      podsumowanieRozszerzen,
      `${zawezone} zawężonych grup zakresu, ${zakresyNarzedzi.length} zakresów narzędzi.`,
    );
  }

  function wypelnijIzolacje(): void {
    if (macierz === null || wzorZakresu === null) return;
    for (const stary of [...macierz.querySelectorAll('.pk-mrz')]) stary.remove();
    let wlaczone = 0;
    for (const zakres of przelaczniki) {
      if (zakres.isolated) wlaczone += 1;
      const wiersz = wzorZakresu.cloneNode(true) as HTMLElement;
      wiersz.dataset.zakres = zakres.scope;
      tekst(wiersz.querySelector('.rozc'), zakres.explanation ?? zakres.scope);
      const przelacznik = wiersz.querySelector('.dn-przelacznik');
      przelacz(przelacznik, zakres.isolated);
      przelacznik?.setAttribute('aria-label', zakres.explanation ?? zakres.scope);
      macierz.appendChild(wiersz);
    }
    tekst(
      podsumowanieIzolacji,
      `${wlaczone}/${przelaczniki.length} zakresów izolacji włączonych.`,
    );
  }

  function wypelnijSiec(agent: Agent | null): void {
    przelacz(przelacznikSieci, (agent?.subagentLimit ?? 0) > 0);
    if (liczbaPodagentow !== null) {
      liczbaPodagentow.value = String(agent?.subagentLimit ?? 0);
    }
  }

  function wypiszPolityke(polityka: AgentPolicy | null): void {
    if (polityka === null) {
      tekst(kod, '');
      return;
    }
    const zawezone = polityka.permissions.filter((wpis) => !wpis.granted).length;
    const izolowane = polityka.technicalSwitches.filter((wpis) => wpis.isolated).length;
    const nadpisanie = polityka.overriddenBy === undefined
      ? ''
      : `\nnadpisana zakresem: ${polityka.overriddenBy}`;
    tekst(
      kod,
      [
        `zakresy uprawnień: ${zawezone} zawężonych z ${polityka.permissions.length}`,
        `moduły i zasoby: ${polityka.moduleCodes.length} przyznanych`,
        `izolacja techniczna: ${izolowane}/${polityka.technicalSwitches.length} zakresów`,
        `Subagent Network: ${polityka.subagentEnabled ? 'tak' : 'nie'}, do ${polityka.subagentLimit}`,
      ].join('\n') + nadpisanie,
    );
  }

  async function wczytaj(): Promise<void> {
    const agent = stan.ekspert();
    if (agent === null) {
      przelaczniki = [];
      wypelnijZakresy(null);
      wypelnijIzolacje();
      wypelnijSiec(null);
      wypiszPolityke(null);
      tekst(banerWiodacy, BEZ_EKSPERTA);
      return;
    }
    tekst(banerWiodacy, 'Zakres zawężony tutaj obowiązuje eksperta we wszystkich modułach.');
    const izolacja = await wywolaj(stan.kanal, Command.AgentIsolationGet, { agentId: agent.id });
    przelaczniki = izolacja.udany ? izolacja.wynik?.switches ?? [] : [];
    const narzedzia = await wywolaj(stan.kanal, Command.ToolsScopeList, { profileId: agent.id });
    zakresyNarzedzi = narzedzia.udany ? narzedzia.wynik?.scopes ?? [] : [];
    wypelnijZakresy(agent);
    wypelnijIzolacje();
    wypelnijSiec(agent);
    wypelnijGrupe(
      grupaModulow,
      kodyModulow.map((kodModulu) => ({
        klucz: kodModulu,
        podpis: kodModulu,
        zaznaczona: agent.moduleCodes.includes(kodModulu),
      })),
    );
    tekst(podsumowanieModulow, `${agent.moduleCodes.length} z ${kodyModulow.length} modułów.`);
    const polityka = await wywolaj(stan.kanal, Command.AgentPolicyGet, {
      agentId: agent.id,
      windowId: stan.idOkna === '' ? undefined : stan.idOkna,
    });
    wypiszPolityke(polityka.udany ? polityka.wynik?.policy ?? null : null);
  }

  async function przestawZakres(wiersz: HTMLElement, przyznane: boolean): Promise<void> {
    const agent = stan.ekspert();
    const grupa = wiersz.dataset.grupa ?? '';
    if (agent === null || grupa === '') return;
    const wynik = await wywolaj(stan.kanal, Command.AgentPermissionSet, {
      agentId: agent.id,
      group: grupa as (typeof AgentPermissionGroup)[keyof typeof AgentPermissionGroup],
      granted: przyznane,
    });
    if (!wynik.udany) oglos('Uprawnienia', wynik.blad?.message ?? 'Rdzeń odrzucił zmianę.', 'blad');
    await wczytaj();
  }

  async function zdejmijZakres(wiersz: HTMLElement): Promise<void> {
    const agent = stan.ekspert();
    const grupa = wiersz.dataset.grupa ?? '';
    if (agent === null || grupa === '') return;
    await wywolaj(stan.kanal, Command.AgentPermissionRemove, {
      agentId: agent.id,
      group: grupa as (typeof AgentPermissionGroup)[keyof typeof AgentPermissionGroup],
    });
    await wczytaj();
  }

  async function przestawNarzedzie(
    wiersz: HTMLElement,
    kolumna: number,
    wlaczone: boolean,
  ): Promise<void> {
    const nazwa = wiersz.dataset.narzedzie ?? '';
    const profil = wiersz.dataset.profil ?? '';
    const zakres = zakresyNarzedzi.find((wpis) => wpis.toolName === nazwa);
    if (nazwa === '' || profil === '' || zakres === undefined) return;
    const wynik = await wywolaj(stan.kanal, Command.ToolsScopeSet, {
      profileId: profil,
      toolName: nazwa,
      enabled: kolumna === 0 ? wlaczone : undefined,
      confirmRequired: kolumna === 1 ? wlaczone : undefined,
      callLimit: kolumna === 2 ? (wlaczone ? Math.max(zakres.callLimit, 1) : 0) : undefined,
    });
    if (!wynik.udany) {
      oglos('Zakresy narzędzi', wynik.blad?.message ?? 'Rdzeń odrzucił zmianę.', 'blad');
    }
    await wczytaj();
  }

  async function przestawIzolacje(wiersz: HTMLElement): Promise<void> {
    const agent = stan.ekspert();
    const zakres = wiersz.dataset.zakres ?? '';
    if (agent === null || zakres === '') return;
    const zmienione = przelaczniki.map((wpis) => (
      wpis.scope === zakres ? { ...wpis, isolated: !wpis.isolated } : wpis
    ));
    const wynik = await wywolaj(stan.kanal, Command.AgentIsolationSet, {
      agentId: agent.id,
      switches: zmienione,
    });
    if (!wynik.udany) oglos('Izolacja', wynik.blad?.message ?? 'Rdzeń odrzucił zmianę.', 'blad');
    await wczytaj();
  }

  async function przestawSiec(wlaczona: boolean, granica: number): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentSubagentSet, {
      agentId: agent.id,
      enabled: wlaczona,
      limit: granica,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Subagent Network', wynik.blad?.message ?? 'Rdzeń odrzucił zmianę.', 'blad');
      return;
    }
    stan.postaw(wynik.wynik.agent);
  }

  async function zapiszModuly(): Promise<void> {
    const agentMozeMoze = stan.ekspert();
    if (agentMozeMoze === null) return;
    const agentMoze = agentMozeMoze;
    const agent = agentMoze;
    const wynik = await wywolaj(stan.kanal, Command.AgentModulesSet, {
      agentId: agent.id,
      moduleCodes: zaznaczoneKlucze(grupaModulow),
    });
    if (wynik.udany && wynik.wynik !== undefined) stan.postaw(wynik.wynik.agent);
  }

  panel.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przyciskGrupy = cel.closest<HTMLElement>('.pk-grupa[data-grupa]');
    if (przyciskGrupy !== null) {
      for (const przycisk of grupy) {
        const wybrana = przycisk === przyciskGrupy;
        przycisk.setAttribute('aria-selected', String(wybrana));
        const zawartosc = panel.querySelector(`#${przycisk.dataset.grupa ?? ''}`);
        if (wybrana) zawartosc?.removeAttribute('hidden');
        else zawartosc?.setAttribute('hidden', '');
      }
      return;
    }
    const zakres = cel.closest<HTMLElement>('#grupa-izolacja .pk-mrz');
    if (zakres !== null) {
      void przestawIzolacje(zakres);
      return;
    }
    if (cel.closest('#grupa-mtai .pk-mrz') !== null) {
      const wlaczona = !czyWlaczony(przelacznikSieci);
      void przestawSiec(wlaczona, wlaczona ? Number(liczbaPodagentow?.value ?? '1') || 1 : 0);
      return;
    }
    if (cel.closest('#grupa-moduly') !== null && cel instanceof HTMLInputElement) {
      void zapiszModuly();
      return;
    }
    if (!(cel instanceof HTMLInputElement) || cel.disabled) return;
    const wierszMozeMoze = cel.closest<HTMLElement>('#grupa-rozsz tr');
    if (wierszMozeMoze === null) return;
    const wierszMoze = wierszMozeMoze;
    const wiersz = wierszMoze;
    if (wiersz.dataset.grupa !== undefined) {
      void przestawZakres(wiersz, cel.checked);
      return;
    }
    const kolumna = [...wiersz.querySelectorAll('input')].indexOf(cel);
    if (wiersz.dataset.narzedzie !== undefined) {
      void przestawNarzedzie(wiersz, kolumna, cel.checked);
    }
  }, stan.przy);

  panel.addEventListener('contextmenu', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wierszMozeMoze = cel.closest<HTMLElement>('#grupa-rozsz tr[data-grupa]');
    if (wierszMozeMoze === null) return;
    const wierszMoze = wierszMozeMoze;
    const wiersz = wierszMoze;
    zdarzenie.preventDefault();
    void zdejmijZakres(wiersz);
  }, stan.przy);

  liczbaPodagentow?.addEventListener('change', () => {
    const granica = Number(liczbaPodagentow.value) || 0;
    void przestawSiec(granica > 0, granica);
  }, stan.przy);

  void wczytajModuly();
  return {
    pokaz: () => {
      void wczytaj();
    },
  };
}
