/**
 * Panel instrukcji projektu WorkSpace: treść warstwy, zapis, wykaz wersji
 * i przywrócenie wersji. Komendy przeniesione z poprzedniego klienta
 * (`moduly/workspace/panel-instrukcji.ts`): treść bieżącą oddaje `config.get`
 * pod kluczem `workspace.instrukcje`, zapis idzie `workspace.instructions.set`,
 * historię niosą `workspace.instructions.version.list` i `.restore`.
 */

import {
  Command,
  ConfigScope,
  type WorkspaceInstructionsVersion,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { data } from './okno-modulu.ts';
import { sesjaBiezaca } from './sesja-biezaca.ts';
import { odmien } from './workspace-projekt.ts';
import {
  cialoPanelu,
  kopiaWzoru,
  przyjmij,
  usun,
  wpisz,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

/** Klucz nastawy pod którym rdzeń trzyma instrukcje modułu — nazwa z poprzedniego klienta. */
const KLUCZ_INSTRUKCJI = 'workspace.instrukcje';

const WARSTWY: Record<string, ConfigScope> = {
  projekt: ConfigScope.Project,
  sesja: ConfigScope.Session,
};

export interface StanInstrukcji {
  wczytaj: () => Promise<void>;
}

export function zwiazInstrukcje(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  przy: AddEventListenerOptions,
): StanInstrukcji {
  const panel = wezly.instrukcje;
  const cialo = cialoPanelu(panel);
  const konsola = cialo?.querySelector<HTMLElement>('.pt-konsola') ?? null;
  const niezapisane = panel?.querySelector<HTMLElement>('.dn-plakietka--ostrzezenie') ?? null;
  const wiersze = [...(cialo?.querySelectorAll<HTMLElement>('.dn-wykaz-modulu-poz') ?? [])];
  const miara = wiersze[0]?.querySelector<HTMLElement>('.dn-meta') ?? null;
  const wzorWersji = kopiaWzoru(wiersze[0] ?? null);
  let warstwa: ConfigScope = ConfigScope.Project;

  // Prototyp niesie trzy przyciski; pokrycie w kontrakcie ma tylko wykaz wersji.
  const przyciski = [...(cialo?.querySelectorAll<HTMLButtonElement>('.rt-mp-dol button') ?? [])];
  const przyciskWersji =
    przyciski.find((cel) => (cel.textContent ?? '').startsWith('Wersje')) ?? null;
  for (const przycisk of przyciski) if (przycisk !== przyciskWersji) przycisk.remove();
  // Liczbę wersji podaje rdzeń po naciśnięciu; do tej chwili prototyp mówi „(4)".
  if (przyciskWersji !== null) przyciskWersji.textContent = 'Wersje ▾';

  const oznaczNiezapisane = (stan: boolean): void => {
    if (niezapisane === null) return;
    niezapisane.hidden = !stan;
    niezapisane.textContent = 'niezapisane ●';
  };

  const bytWarstwy = (): string =>
    warstwa === ConfigScope.Session ? sesjaBiezaca() : idProjektu();

  const wczytaj = async (): Promise<void> => {
    const tresc = await pobierzTresc(kanal, warstwa, bytWarstwy());
    if (konsola !== null) konsola.textContent = tresc;
    wpisz(miara, opisMiary(tresc));
    oznaczNiezapisane(false);
  };

  if (konsola !== null) {
    konsola.textContent = '';
    konsola.contentEditable = 'true';
    konsola.spellcheck = false;
    konsola.addEventListener('input', () => {
      oznaczNiezapisane(true);
      wpisz(miara, opisMiary(konsola.textContent ?? ''));
    }, przy);
    konsola.addEventListener('blur', () => {
      void zapisz(kanal, idProjektu(), warstwa, bytWarstwy(), konsola.textContent ?? '', oznaczNiezapisane);
    }, przy);
  }
  oznaczNiezapisane(false);
  wpisz(miara, opisMiary(''));

  zwiazWarstwy(wiersze[1] ?? null, (wybrana) => {
    warstwa = wybrana;
    void wczytaj();
  }, przy);

  przyciskWersji?.addEventListener('click', () => {
    void pokazWersje(kanal, idProjektu(), cialo, wzorWersji, przyciskWersji, konsola, przy);
  }, przy);

  return { wczytaj };
}

/** Żetony warstwy odpowiadają poziomom zasięgu z kontraktu: projekt i karta sesji. */
function zwiazWarstwy(
  wiersz: HTMLElement | null,
  wybierz: (warstwa: ConfigScope) => void,
  przy: AddEventListenerOptions,
): void {
  if (wiersz === null) return;
  const zetony = [...wiersz.querySelectorAll<HTMLElement>('.sta-chip')];
  for (const zeton of zetony) {
    const nazwa = (zeton.textContent ?? '').replace('●', '').trim();
    const poziom = WARSTWY[nazwa];
    if (poziom === undefined) {
      usun(zeton);
      continue;
    }
    zeton.textContent = nazwa;
    zeton.setAttribute('aria-pressed', String(poziom === ConfigScope.Project));
    zeton.addEventListener('click', () => {
      for (const inny of zetony) inny.setAttribute('aria-pressed', String(inny === zeton));
      wybierz(poziom);
    }, przy);
  }
}

async function pobierzTresc(
  kanal: Kanal,
  warstwa: ConfigScope,
  bytWarstwy: string,
): Promise<string> {
  if (bytWarstwy === '') return '';
  const wynik = await wywolaj(kanal, Command.ConfigGet, {
    key: KLUCZ_INSTRUKCJI,
    scope: warstwa,
    scopeId: bytWarstwy,
  });
  const wpis = wynik.wynik?.entries.find((pozycja) => pozycja.key === KLUCZ_INSTRUKCJI);
  return typeof wpis?.value === 'string' ? wpis.value : '';
}

function opisMiary(tresc: string): string {
  return odmien(tresc.length, ['znak', 'znaki', 'znaków']);
}

async function zapisz(
  kanal: Kanal,
  idProjektu: string,
  warstwa: ConfigScope,
  bytWarstwy: string,
  tresc: string,
  oznaczNiezapisane: (stan: boolean) => void,
): Promise<void> {
  if (idProjektu === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceInstructionsSet, {
    projectId: idProjektu,
    content: tresc,
    scope: warstwa,
    ...(bytWarstwy === '' ? {} : { scopeId: bytWarstwy }),
  });
  if (przyjmij('Instrukcje projektu', wynik) === null) return;
  oznaczNiezapisane(false);
}

async function pokazWersje(
  kanal: Kanal,
  idProjektu: string,
  cialo: HTMLElement | null,
  wzor: HTMLElement | null,
  przycisk: HTMLButtonElement,
  konsola: HTMLElement | null,
  przy: AddEventListenerOptions,
): Promise<void> {
  if (cialo === null || idProjektu === '') return;
  for (const stary of cialo.querySelectorAll('[data-wersja]')) stary.remove();
  const rozwiniete = przycisk.getAttribute('aria-expanded') === 'true';
  przycisk.setAttribute('aria-expanded', String(!rozwiniete));
  if (rozwiniete) return;
  const wynik = await wywolaj(kanal, Command.WorkspaceInstructionsVersionList, {
    projectId: idProjektu,
    limit: 20,
  });
  const odpowiedzMoze = przyjmij('Wersje instrukcji', wynik);
  if (odpowiedzMoze === null) return;
  const odpowiedz = odpowiedzMoze;
  przycisk.textContent = `Wersje (${odpowiedz.total}) ▾`;
  for (const wersja of odpowiedz.versions) {
    const wiersz = wierszWersji(wzor, wersja);
    if (wiersz === null) continue;
    wiersz.addEventListener('click', () => {
      void przywroc(kanal, idProjektu, wersja.id, konsola);
    }, przy);
    cialo.appendChild(wiersz);
  }
}

function wierszWersji(
  wzor: HTMLElement | null,
  wersja: WorkspaceInstructionsVersion,
): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
  wiersz.replaceChildren();
  wiersz.append(`${data(wersja.createdAt)} `);
  if (meta !== null) {
    meta.textContent = odmien(wersja.content.length, ['znak', 'znaki', 'znaków']);
    wiersz.appendChild(meta);
  }
  wiersz.dataset.wersja = wersja.id;
  return wiersz;
}

async function przywroc(
  kanal: Kanal,
  idProjektu: string,
  idWersji: string,
  konsola: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.WorkspaceInstructionsVersionRestore, {
    projectId: idProjektu,
    versionId: idWersji,
  });
  const odpowiedzMoze = przyjmij('Przywrócenie wersji', wynik);
  if (odpowiedzMoze === null) return;
  const odpowiedz = odpowiedzMoze;
  if (konsola !== null) konsola.textContent = odpowiedz.instructions.content;
  oglos('Instrukcje projektu', `Przywrócono wersję z ${data(odpowiedz.version.createdAt)}.`);
}
