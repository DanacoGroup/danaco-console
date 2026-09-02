/**
 * Notatki projektu w panelu „Artefakty"; wywołania z poprzedniego klienta
 * (`wiki-projektu.ts`): `workspace.note.list`, `.note.tree.get`, `.note.get`,
 * `.note.backlink.list` i `workspace.comment.list`.
 */

import {
  Command,
  WorkspaceEntityKind,
  type WorkspaceNote,
  type WorkspaceNoteNode,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { data } from './okno-modulu.ts';
import { odmien } from './workspace-projekt.ts';
import {
  cialoPanelu,
  kopiaWzoru,
  niegotowe,
  przyjmij,
  wpisz,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

export interface StanNotatek {
  odswiez: () => Promise<void>;
}

export function zwiazNotatki(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  przy: AddEventListenerOptions,
): StanNotatek {
  const cialo = cialoPanelu(wezly.artefakty);
  const wzor = kopiaWzoru(cialo?.querySelector<HTMLElement>('.dn-wykaz-modulu-poz') ?? null);
  for (const stary of cialo?.querySelectorAll('.dn-wykaz-modulu-poz') ?? []) stary.remove();

  const odswiez = async (): Promise<void> => {
    const projekt = idProjektu();
    if (cialo === null || projekt === '') return;
    const notatki = await pobierzNotatki(kanal, projekt);
    const glebokosci = await pobierzGlebokosci(kanal, projekt);
    cialo.replaceChildren();
    if (notatki.length === 0) {
      niegotowe(cialo, 'Projekt nie ma jeszcze notatek.');
      return;
    }
    for (const notatka of notatki) {
      const wiersz = wierszNotatki(wzor, notatka, glebokosci.get(notatka.id) ?? 0);
      if (wiersz !== null) cialo.appendChild(wiersz);
    }
  };

  cialo?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wierszMoze = cel.closest<HTMLElement>('[data-notatka]');
    if (wierszMoze === null) return;
    const wiersz = wierszMoze;
    void opiszPowiazania(kanal, idProjektu(), wiersz);
  }, przy);

  return { odswiez };
}

async function pobierzNotatki(kanal: Kanal, idProjektu: string): Promise<WorkspaceNote[]> {
  const wynik = await wywolaj(kanal, Command.WorkspaceNoteList, {
    projectId: idProjektu,
    limit: 100,
  });
  return przyjmij('Notatki projektu', wynik)?.notes ?? [];
}

/** Wcięcie wiersza bierze się z drzewa; wykaz płaski nie niesie zagnieżdżenia. */
async function pobierzGlebokosci(
  kanal: Kanal,
  idProjektu: string,
): Promise<Map<string, number>> {
  const glebokosci = new Map<string, number>();
  const wynik = await wywolaj(kanal, Command.WorkspaceNoteTreeGet, { projectId: idProjektu });
  if (!wynik.udany || wynik.wynik === undefined) return glebokosci;
  const wezly = new Map<string, WorkspaceNoteNode>();
  for (const wezel of wynik.wynik.nodes) wezly.set(wezel.noteId, wezel);
  for (const wezel of wynik.wynik.nodes) {
    let glebokosc = 0;
    let rodzic = wezel.parentNoteId;
    while (rodzic !== undefined && glebokosc < 12) {
      glebokosc += 1;
      rodzic = wezly.get(rodzic)?.parentNoteId;
    }
    glebokosci.set(wezel.noteId, glebokosc);
  }
  return glebokosci;
}

function wierszNotatki(
  wzor: HTMLElement | null,
  notatka: WorkspaceNote,
  glebokosc: number,
): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
  wiersz.replaceChildren();
  wiersz.append(`${notatka.title} `);
  if (meta !== null) {
    meta.textContent = data(notatka.updatedAt);
    wiersz.appendChild(meta);
  }
  wiersz.style.paddingInlineStart = `${glebokosc * 12}px`;
  wiersz.dataset.notatka = notatka.id;
  return wiersz;
}

async function opiszPowiazania(
  kanal: Kanal,
  idProjektu: string,
  wiersz: HTMLElement,
): Promise<void> {
  const idNotatki = wiersz.dataset.notatka ?? '';
  if (idProjektu === '' || idNotatki === '') return;
  const [notatka, odnosniki, komentarze] = await Promise.all([
    wywolaj(kanal, Command.WorkspaceNoteGet, { noteId: idNotatki }),
    wywolaj(kanal, Command.WorkspaceNoteBacklinkList, { noteId: idNotatki, limit: 50 }),
    wywolaj(kanal, Command.WorkspaceCommentList, {
      projectId: idProjektu,
      targetKind: WorkspaceEntityKind.Note,
      targetId: idNotatki,
      limit: 50,
    }),
  ]);
  const tresc = notatka.wynik?.note.content ?? '';
  if (tresc !== '') wiersz.setAttribute('title', tresc);
  const czesci: string[] = [];
  if (odnosniki.udany && odnosniki.wynik !== undefined) {
    czesci.push(odmien(odnosniki.wynik.total, ['odnośnik', 'odnośniki', 'odnośników']));
  }
  if (komentarze.udany && komentarze.wynik !== undefined) {
    czesci.push(odmien(komentarze.wynik.total, ['komentarz', 'komentarze', 'komentarzy']));
  }
  if (czesci.length > 0) wpisz(wiersz.querySelector('.dn-meta'), czesci.join(' · '));
}
