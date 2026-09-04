/**
 * Panel pamięci projektu WorkSpace: wpisy warstwy projektu i karty sesji,
 * wyszukiwanie, zajętość, dokładanie i usuwanie wpisu. Wywołania przeniesione
 * z poprzedniego klienta (`moduly/workspace/pamiec-kontekstu.ts`):
 * `workspace.context.get`, `workspace.context.set`, `memory.list`,
 * `memory.delete` oraz `workspace.search.project` przy zapytaniu bez trafień.
 */

import {
  Command,
  MemoryEntryOrigin,
  type WorkspaceDashboard,
  type WorkspaceMemoryEntry,
  type WorkspaceSearchHit,
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
  niegotowe,
  wpiszMiarePanelu,
  przyjmij,
  udzialProcentowy,
  ustawPostep,
  usun,
  wpisz,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

export interface StanPamieci {
  odswiez: () => Promise<void>;
  opiszZajetosc: (pulpit: WorkspaceDashboard) => void;
}

export function zwiazPamiec(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  przy: AddEventListenerOptions,
): StanPamieci {
  const panel = wezly.pamiec;
  const cialo = cialoPanelu(panel);
  const szukaj = cialo?.querySelector<HTMLInputElement>('input[type="search"]') ?? null;
  const zajetosc = cialo?.querySelector<HTMLElement>('.dn-wykaz-modulu-poz') ?? null;
  const postep = cialo?.querySelector<HTMLElement>('.dn-postep') ?? null;
  const wzorWpisu = kopiaWzoru(cialo?.querySelector<HTMLElement>('.rt-kluczowy') ?? null);
  const wzorPropozycji = kopiaWzoru(cialo?.querySelector<HTMLElement>('.dn-alert--info') ?? null);
  const gniazdo = document.createElement('div');
  gniazdo.className = 'dn-wykaz-modulu';

  for (const przykladowy of cialo?.querySelectorAll('.rt-kluczowy, .dn-alert--info') ?? []) {
    przykladowy.remove();
  }
  cialo?.appendChild(gniazdo);

  const odswiez = async (): Promise<void> => {
    const projekt = idProjektu();
    if (projekt === '') return;
    const pytanie = (szukaj?.value ?? '').trim();
    const wpisy = await pobierzWpisy(kanal, projekt, pytanie);
    wpiszMiarePanelu(panel, odmien(wpisy.length, ['wpis', 'wpisy', 'wpisów']));
    if (wpisy.length === 0 && pytanie !== '') {
      await pokazTrafienia(kanal, projekt, pytanie, gniazdo, wzorWpisu);
      return;
    }
    wypelnijWpisy(gniazdo, wzorWpisu, wzorPropozycji, wpisy, kanal, projekt, odswiez, przy);
  };

  let odliczanie = 0;
  szukaj?.addEventListener('input', () => {
    globalThis.clearTimeout(odliczanie);
    odliczanie = globalThis.setTimeout(() => void odswiez(), 250);
  }, przy);

  const dodaj = panel?.querySelector<HTMLElement>('.sta-okno-belka .dn-btn-ikona');
  if (dodaj !== null && dodaj !== undefined) {
    dodaj.setAttribute('title', 'Dodaj wpis z pola wyszukiwania');
    dodaj.addEventListener('click', () => {
      void dolozWpis(kanal, idProjektu(), szukaj, odswiez);
    }, przy);
  }

  return {
    odswiez,
    opiszZajetosc: (pulpit) => {
      const udzial = udzialProcentowy(pulpit.memoryUsedBytes, pulpit.memoryCapacityBytes);
      ustawPostep(postep, udzial);
      wpisz(zajetosc?.querySelector('.dn-meta'), udzial === 0 ? '—' : `${Math.round(udzial)}%`);
    },
  };
}

function wypelnijWpisy(
  gniazdo: HTMLElement,
  wzorWpisu: HTMLElement | null,
  wzorPropozycji: HTMLElement | null,
  wpisy: WorkspaceMemoryEntry[],
  kanal: Kanal,
  idProjektu: string,
  odswiez: () => Promise<void>,
  przy: AddEventListenerOptions,
): void {
  gniazdo.replaceChildren();
  if (wpisy.length === 0) {
    niegotowe(gniazdo, 'Pamięć tego projektu jest pusta.');
    return;
  }
  for (const wpis of wpisy) {
    const propozycja = wpis.origin === MemoryEntryOrigin.Model;
    const wezel = kopiaWzoru(propozycja ? wzorPropozycji : wzorWpisu);
    if (wezel === null) continue;
    wezel.dataset.wpis = wpis.id;
    if (propozycja) {
      zbudujPropozycje(wezel, wpis, kanal, idProjektu, odswiez, przy);
    } else {
      zbudujWpis(wezel, wpis);
    }
    gniazdo.appendChild(wezel);
  }
}

function zbudujWpis(wezel: HTMLElement, wpis: WorkspaceMemoryEntry): void {
  const stopka = wezel.querySelector<HTMLElement>('small');
  wezel.replaceChildren();
  wezel.append(`${wpis.pinned === true ? '📌 ' : ''}${wpis.content}`);
  if (stopka === null) return;
  const czesci = [...(wpis.tags ?? []).map((znacznik) => `#${znacznik}`), data(wpis.updatedAt)];
  stopka.textContent = czesci.join(' · ');
  wezel.append(document.createElement('br'), stopka);
}

function zbudujPropozycje(
  wezel: HTMLElement,
  wpis: WorkspaceMemoryEntry,
  kanal: Kanal,
  idProjektu: string,
  odswiez: () => Promise<void>,
  przy: AddEventListenerOptions,
): void {
  const przyciski = [...wezel.querySelectorAll<HTMLButtonElement>('button')];
  const akceptuj = przyciski[0] ?? null;
  const odrzuc = przyciski[1] ?? null;
  for (const przycisk of przyciski.slice(2)) usun(przycisk);
  wezel.replaceChildren();
  wezel.append(`🤖 ${wpis.content} `);
  if (akceptuj !== null) {
    akceptuj.textContent = 'Akceptuj';
    akceptuj.addEventListener('click', () => {
      void przyjmijPropozycje(kanal, idProjektu, wpis, odswiez);
    }, przy);
    wezel.appendChild(akceptuj);
  }
  if (odrzuc === null) return;
  odrzuc.textContent = 'Odrzuć';
  odrzuc.addEventListener('click', () => {
    void usunWpis(kanal, wpis.id, odswiez);
  }, przy);
  wezel.appendChild(odrzuc);
}

/**
 * Warstwa projektu idzie komendą kontekstu, warstwa karty sesji wykazem pamięci
 * — poprzedni klient rozdzielał je tak samo, bo rdzeń trzyma je pod dwoma
 * zasięgami i żadna z komend nie oddaje obu naraz.
 */
async function pobierzWpisy(
  kanal: Kanal,
  idProjektu: string,
  pytanie: string,
): Promise<WorkspaceMemoryEntry[]> {
  const wynik = await wywolaj(kanal, Command.WorkspaceContextGet, {
    projectId: idProjektu,
    limit: 60,
    ...(pytanie === '' ? {} : { query: pytanie }),
  });
  const wpisy = przyjmij('Pamięć projektu', wynik)?.entries ?? [];
  const idSesji = sesjaBiezaca();
  if (idSesji === '') return wpisy;
  const sesyjne = await wywolaj(kanal, Command.MemoryList, {
    sessionId: idSesji,
    projectId: idProjektu,
    limit: 60,
  });
  const znane = new Set(wpisy.map((wpis) => wpis.id));
  const dolozone = (sesyjne.wynik?.entries ?? []).filter((wpis) => !znane.has(wpis.id));
  return [...wpisy, ...dolozone];
}

async function usunWpis(
  kanal: Kanal,
  idWpisu: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.MemoryDelete, { entryId: idWpisu });
  if (przyjmij('Pamięć projektu', wynik) === null) return;
  await odswiez();
}

async function pokazTrafienia(
  kanal: Kanal,
  idProjektu: string,
  pytanie: string,
  gniazdo: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.WorkspaceSearchProject, {
    projectId: idProjektu,
    query: pytanie,
    limit: 30,
  });
  const trafienia = przyjmij('Szukanie w projekcie', wynik)?.hits ?? [];
  gniazdo.replaceChildren();
  if (trafienia.length === 0) {
    niegotowe(gniazdo, 'Rdzeń nie znalazł niczego na to zapytanie.');
    return;
  }
  for (const trafienie of trafienia) {
    const wezel = wierszTrafienia(wzor, trafienie);
    if (wezel !== null) gniazdo.appendChild(wezel);
  }
}

function wierszTrafienia(
  wzor: HTMLElement | null,
  trafienie: WorkspaceSearchHit,
): HTMLElement | null {
  const wezelMoze = kopiaWzoru(wzor);
  if (wezelMoze === null) return null;
  const wezel = wezelMoze;
  const stopka = wezel.querySelector<HTMLElement>('small');
  wezel.replaceChildren();
  wezel.append(trafienie.title);
  if (stopka !== null) {
    stopka.textContent = [trafienie.kind, trafienie.snippet ?? ''].filter((czesc) => czesc !== '').join(' · ');
    wezel.append(document.createElement('br'), stopka);
  }
  wezel.dataset.trafienie = trafienie.entityId;
  return wezel;
}

async function dolozWpis(
  kanal: Kanal,
  idProjektu: string,
  szukaj: HTMLInputElement | null,
  odswiez: () => Promise<void>,
): Promise<void> {
  const tresc = (szukaj?.value ?? '').trim();
  if (idProjektu === '' || tresc === '') {
    oglos('Pamięć projektu', 'Wpisz treść w polu wyszukiwania, aby założyć wpis.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.WorkspaceContextSet, {
    projectId: idProjektu,
    content: tresc,
    origin: MemoryEntryOrigin.Operator,
  });
  if (przyjmij('Pamięć projektu', wynik) === null) return;
  if (szukaj !== null) szukaj.value = '';
  await odswiez();
}

async function przyjmijPropozycje(
  kanal: Kanal,
  idProjektu: string,
  wpis: WorkspaceMemoryEntry,
  odswiez: () => Promise<void>,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.WorkspaceContextSet, {
    projectId: idProjektu,
    entryId: wpis.id,
    content: wpis.content,
    pinned: true,
    origin: MemoryEntryOrigin.Operator,
  });
  if (przyjmij('Pamięć projektu', wynik) === null) return;
  await odswiez();
}
