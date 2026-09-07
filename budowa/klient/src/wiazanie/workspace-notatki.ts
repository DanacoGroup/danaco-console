/**
 * Notatki projektu w panelu „Artefakty": wykaz, powiązania, założenie strony
 * i usunięcie stojącej. Wywołania: `workspace.note.list`, `.note.tree.get`,
 * `.note.get`, `.note.backlink.list`, `.note.save`, `.note.delete`
 * i `workspace.comment.list`.
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
import { oglos } from './ogloszenie.ts';
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
  const pole = zalozPasDzialan(cialo);

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

  wezly.artefakty?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-notatka-nowa]') !== null) {
      zdarzenie.stopPropagation();
      void zalozNotatke(kanal, pole, idProjektu(), odswiez);
      return;
    }
    const wierszMoze = cel.closest<HTMLElement>('[data-notatka]');
    if (wierszMoze === null) return;
    const wiersz = wierszMoze;
    if (cel.closest('[data-notatka-usun]') !== null) {
      zdarzenie.stopPropagation();
      void usunNotatke(kanal, wiersz, odswiez);
      return;
    }
    if (cel.closest('[data-notatka-komentarz]') !== null) {
      zdarzenie.stopPropagation();
      void dolozKomentarz(kanal, idProjektu(), wiersz);
      return;
    }
    if (cel.closest('[data-komentarz-zdejmij]') !== null) {
      zdarzenie.stopPropagation();
      void zdejmijKomentarz(kanal, idProjektu(), wiersz);
      return;
    }
    void opiszPowiazania(kanal, idProjektu(), wiersz);
  }, przy);

  return { odswiez };
}

/* Pas działań stoi nad wykazem, bo ciało panelu jest wymieniane przy każdym
   odświeżeniu i wszystko w nim postawione znikłoby wraz z wykazem. */
function zalozPasDzialan(cialo: HTMLElement | null): HTMLInputElement | null {
  if (cialo === null || cialo.parentElement === null) return null;
  const pas = cialo.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const pole = cialo.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Tytuł nowej strony';
  pole.setAttribute('aria-label', 'Tytuł nowej strony');
  const przycisk = cialo.ownerDocument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przycisk.dataset.notatkaNowa = '';
  przycisk.textContent = 'Nowa strona';
  pas.append(pole, przycisk);
  cialo.parentElement.insertBefore(pas, cialo);
  return pole;
}

/* Nowa strona wchodzi z pustą treścią: `workspace.note.save` bez `noteId`
   zakłada notatkę, a treść dopisuje się w oknie strony. */
async function zalozNotatke(
  kanal: Kanal,
  pole: HTMLInputElement | null,
  idProjektu: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  const tytul = (pole?.value ?? '').trim();
  if (idProjektu === '' || tytul === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceNoteSave, {
    projectId: idProjektu,
    title: tytul,
    content: '',
  });
  if (przyjmij('Nowa strona', wynik) === null) return;
  if (pole !== null) pole.value = '';
  await odswiez();
}

/* Usunięcie jest nieodwracalne, więc pierwsze naciśnięcie uzbraja przycisk.
   Bez `withChildren` strony podrzędne przechodzą pod stronę nadrzędną
   usuwanej, zamiast zniknąć razem z nią. */
async function usunNotatke(
  kanal: Kanal,
  wiersz: HTMLElement,
  odswiez: () => Promise<void>,
): Promise<void> {
  const idNotatki = wiersz.dataset.notatka ?? '';
  if (idNotatki === '') return;
  const przycisk = wiersz.querySelector<HTMLElement>('[data-notatka-usun]');
  if (przycisk !== null && przycisk.dataset.uzbrojone !== 'tak') {
    przycisk.dataset.uzbrojone = 'tak';
    przycisk.setAttribute('aria-label', 'Naciśnij ponownie, aby usunąć stronę');
    return;
  }
  const wynik = await wywolaj(kanal, Command.WorkspaceNoteDelete, { noteId: idNotatki });
  if (przyjmij('Usunięcie strony', wynik) !== null) await odswiez();
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
  const zdejmij = wiersz.ownerDocument.createElement('button');
  zdejmij.type = 'button';
  zdejmij.className = 'dn-btn-ikona';
  zdejmij.dataset.komentarzZdejmij = '';
  zdejmij.setAttribute('aria-label', 'Zdejmij ostatni komentarz strony');
  zdejmij.textContent = '⌫';
  wiersz.appendChild(zdejmij);
  const skomentuj = wiersz.ownerDocument.createElement('button');
  skomentuj.type = 'button';
  skomentuj.className = 'dn-btn-ikona';
  skomentuj.dataset.notatkaKomentarz = '';
  skomentuj.setAttribute('aria-label', 'Dołóż komentarz do strony');
  skomentuj.textContent = '💬';
  wiersz.appendChild(skomentuj);
  const usun = wiersz.ownerDocument.createElement('button');
  usun.type = 'button';
  usun.className = 'dn-btn-ikona';
  usun.dataset.notatkaUsun = '';
  usun.setAttribute('aria-label', 'Usuń stronę');
  usun.textContent = '×';
  wiersz.appendChild(usun);
  return wiersz;
}

/* Treść komentarza wchodzi w wierszu strony: okna komentarzy wydanie nie
   niesie, a licznik przy wierszu i tak pokazuje, ile ich jest. */
async function dolozKomentarz(
  kanal: Kanal,
  idProjektu: string,
  wiersz: HTMLElement,
): Promise<void> {
  const idNotatki = wiersz.dataset.notatka ?? '';
  if (idProjektu === '' || idNotatki === '') return;
  const tresc = await zapytajWWierszu(wiersz);
  if (tresc === null || tresc === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceCommentAdd, {
    projectId: idProjektu,
    targetKind: WorkspaceEntityKind.Note,
    targetId: idNotatki,
    content: tresc,
  });
  if (przyjmij('Komentarz', wynik) !== null) await opiszPowiazania(kanal, idProjektu, wiersz);
}

/*
zdejmijKomentarz usuwa komentarz ostatni, bo wykazu komentarzy okno nie niesie.

Usunięcie jest nieodwracalne i zabiera także odpowiedzi w wątku — taki jest
domyślny kształt komendy — więc pierwsze naciśnięcie uzbraja przycisk.
*/
async function zdejmijKomentarz(
  kanal: Kanal,
  idProjektu: string,
  wiersz: HTMLElement,
): Promise<void> {
  const idNotatki = wiersz.dataset.notatka ?? '';
  if (idProjektu === '' || idNotatki === '') return;
  const przycisk = wiersz.querySelector<HTMLElement>('[data-komentarz-zdejmij]');
  if (przycisk !== null && przycisk.dataset.uzbrojone !== 'tak') {
    przycisk.dataset.uzbrojone = 'tak';
    przycisk.setAttribute('aria-label', 'Naciśnij ponownie, aby zdjąć komentarz wraz z wątkiem');
    return;
  }
  const wykaz = await wywolaj(kanal, Command.WorkspaceCommentList, {
    projectId: idProjektu,
    targetKind: WorkspaceEntityKind.Note,
    targetId: idNotatki,
    limit: 50,
  });
  const komentarze = przyjmij('Komentarze strony', wykaz)?.comments ?? [];
  const ostatni = komentarze[komentarze.length - 1]?.id ?? '';
  if (ostatni === '') {
    oglos('Komentarze strony', 'Ta strona nie ma jeszcze komentarza do zdjęcia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.WorkspaceCommentDelete, { commentId: ostatni });
  if (przyjmij('Komentarze strony', wynik) !== null) {
    await opiszPowiazania(kanal, idProjektu, wiersz);
  }
}

/* Wiersz staje się polem na czas pisania i wraca do swojej treści; wpis
   zatwierdza Enter, odwołuje Escape albo odejście wskazania. */
function zapytajWWierszu(wiersz: HTMLElement): Promise<string | null> {
  const pole = wiersz.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Treść komentarza';
  pole.setAttribute('aria-label', 'Treść komentarza');
  wiersz.appendChild(pole);
  pole.focus();
  return new Promise((rozstrzygnij) => {
    let domkniete = false;
    const domknij = (wpis: string | null): void => {
      if (domkniete) return;
      domkniete = true;
      pole.remove();
      rozstrzygnij(wpis);
    };
    pole.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') domknij(pole.value.trim());
      if (zdarzenie.key === 'Escape') domknij(null);
    });
    pole.addEventListener('blur', () => domknij(null));
  });
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
