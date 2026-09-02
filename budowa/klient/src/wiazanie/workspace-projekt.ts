/**
 * Pulpit projektu WorkSpace: nagłówek, kafle, stan i oś czasu z komend
 * `workspace.dashboard.get`, `.project.status.set` i `.activity.list`.
 */

import {
  Command,
  WorkspaceProjectStatus,
  type WorkspaceActivityEntry,
  type WorkspaceDashboard,
  type WorkspaceProject,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { data, godzina } from './okno-modulu.ts';
import {
  kopiaWzoru,
  przyjmij,
  wpisz,
  wzorWiersza,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

const STAN_NAZWA: Record<string, string> = {
  [WorkspaceProjectStatus.Active]: 'aktywny',
  [WorkspaceProjectStatus.Paused]: 'wstrzymany',
  [WorkspaceProjectStatus.Archived]: 'zarchiwizowany',
};

const STAN_KLASA: Record<string, string> = {
  [WorkspaceProjectStatus.Active]: 'dn-plakietka dn-plakietka--sukces',
  [WorkspaceProjectStatus.Paused]: 'dn-plakietka dn-plakietka--ostrzezenie',
  [WorkspaceProjectStatus.Archived]: 'dn-plakietka',
};

export function odmien(liczba: number, formy: [string, string, string]): string {
  const setki = liczba % 100;
  if (liczba === 1) return `${liczba} ${formy[0]}`;
  const dziesiatki = setki % 10;
  const mnoga = dziesiatki >= 2 && dziesiatki <= 4 && (setki < 12 || setki > 14);
  return `${liczba} ${mnoga ? formy[1] : formy[2]}`;
}

export function opiszProjekt(wezly: WezlyWorkspace, projekt: WorkspaceProject): void {
  wpisz(wezly.pulpit.querySelector('.wk-naglowek h2'), projekt.name);
  const czesci = [projekt.description ?? '', `założony ${data(projekt.createdAt)}`];
  if (projekt.ownerNote !== undefined && projekt.ownerNote !== '') czesci.push(projekt.ownerNote);
  wpisz(
    wezly.pulpit.querySelector('.wk-naglowek p'),
    czesci.filter((czesc) => czesc !== '').join(' · '),
  );
  opiszStan(wezly, projekt.status);
}

function opiszStan(wezly: WezlyWorkspace, stan: string): void {
  const plakietkaMoze = wezly.pulpit.querySelector<HTMLElement>('.sta-okno-belka .dn-plakietka');
  if (plakietkaMoze === null) return;
  const plakietka = plakietkaMoze;
  const kropka = plakietka.querySelector('.dn-kropka');
  plakietka.className = STAN_KLASA[stan] ?? 'dn-plakietka';
  plakietka.replaceChildren();
  if (kropka !== null && stan === WorkspaceProjectStatus.Active) plakietka.appendChild(kropka);
  plakietka.append(` ${STAN_NAZWA[stan] ?? stan}`);
}

export function opiszKafle(wezly: WezlyWorkspace, pulpit: WorkspaceDashboard): void {
  const podpisy: Record<string, string> = {
    'panel-instructions': licznik(pulpit.instructionSetCount, ['zestaw', 'zestawy', 'zestawów']),
    'panel-context': licznik(pulpit.memoryEntryCount, ['wpis', 'wpisy', 'wpisów']),
    'panel-library': licznik(pulpit.libraryFileCount, ['plik', 'pliki', 'plików']),
    'panel-agents': licznik(pulpit.assignedAgentIds?.length, ['agent', 'agenci', 'agentów']),
  };
  for (const kafel of wezly.pulpit.querySelectorAll<HTMLElement>('.wk-kafel')) {
    wpisz(kafel.querySelector('span'), podpisy[kafel.dataset.panelToggle ?? ''] ?? '');
  }
}

function licznik(wartosc: number | undefined, formy: [string, string, string]): string {
  return wartosc === undefined ? '' : odmien(wartosc, formy);
}

export function opiszOsCzasu(wezly: WezlyWorkspace, wpisy: WorkspaceActivityEntry[]): void {
  const lista = wezly.pulpit.querySelector<HTMLElement>('.wk-os');
  const wzor = wzorWiersza(lista, '.wk-os-poz');
  if (lista === null || wzor === null) return;
  for (const wpis of wpisy) {
    const wiersz = kopiaWzoru(wzor);
    if (wiersz === null) continue;
    const znacznik = wiersz.querySelector('time');
    wiersz.replaceChildren();
    if (znacznik !== null) {
      znacznik.textContent = godzina(wpis.occurredAt);
      znacznik.setAttribute('title', data(wpis.occurredAt));
      wiersz.appendChild(znacznik);
    }
    wiersz.append(` ${wpis.summary}`);
    wiersz.dataset.wpis = wpis.id;
    lista.appendChild(wiersz);
  }
}

export async function pobierzPulpit(
  kanal: Kanal,
  idProjektu: string,
): Promise<WorkspaceDashboard | null> {
  const wynik = await wywolaj(kanal, Command.WorkspaceDashboardGet, { projectId: idProjektu });
  return przyjmij('Pulpit projektu', wynik)?.dashboard ?? null;
}

export async function pobierzAktywnosc(
  kanal: Kanal,
  idProjektu: string,
): Promise<WorkspaceActivityEntry[]> {
  const wynik = await wywolaj(kanal, Command.WorkspaceActivityList, {
    projectId: idProjektu,
    limit: 40,
  });
  return przyjmij('Oś czasu projektu', wynik)?.entries ?? [];
}

/**
 * Menu projektu niesie w prototypie cztery pozycje; pokrycie w rodzinie
 * `workspace.*` ma wyłącznie archiwizacja, więc reszta schodzi z menu.
 */
export function zwiazMenuProjektu(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  poZmianie: (projekt: WorkspaceProject) => void,
  przy: AddEventListenerOptions,
): void {
  const menuMoze = wezly.pulpit.querySelector<HTMLElement>('#menu-proj');
  if (menuMoze === null) return;
  const menu = menuMoze;
  for (const pozycja of [...menu.querySelectorAll<HTMLElement>('.sta-menu-poz')]) {
    const podpis = pozycja.textContent ?? '';
    if (!podpis.startsWith('Zarchiwizuj')) {
      pozycja.remove();
      continue;
    }
    pozycja.addEventListener(
      'click',
      () => {
        void ustawStan(kanal, idProjektu(), WorkspaceProjectStatus.Archived, poZmianie);
      },
      przy,
    );
  }
}

async function ustawStan(
  kanal: Kanal,
  idProjektu: string,
  stan: string,
  poZmianie: (projekt: WorkspaceProject) => void,
): Promise<void> {
  if (idProjektu === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceProjectStatusSet, {
    projectId: idProjektu,
    status: stan as WorkspaceProjectStatus,
  });
  const odpowiedzMoze = przyjmij('Stan projektu', wynik);
  if (odpowiedzMoze === null) return;
  const odpowiedz = odpowiedzMoze;
  oglos('Projekt', `Stan projektu: ${STAN_NAZWA[odpowiedz.project.status] ?? stan}.`);
  poZmianie(odpowiedz.project);
}
