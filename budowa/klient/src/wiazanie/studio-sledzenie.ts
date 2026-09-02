// Śledzenie zmian dokumentu Studia w panelu różnic: wykaz zmian do decyzji
// wraz z ich przyjęciem albo odrzuceniem. Wzory wierszy podaje panel.

import {
  Command,
  StudioChangeDecision,
  StudioChangeKind,
  type StudioTrackedChange,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { nazwijOdmowe } from './studio-dokument.ts';

export interface WzoryZmian {
  usuniety: HTMLElement | null;
  dodany: HTMLElement | null;
  formatowanie: HTMLElement | null;
  decyzja: HTMLElement | null;
  nota: HTMLElement | null;
}

// Bez włączonego śledzenia zapis nadpisuje treść i wykaz nie ma z czego powstać.
export async function przestawSledzenie(
  kanal: Kanal,
  idDokumentu: string,
  wlaczone: boolean,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.StudioTrackingSet, {
    documentId: idDokumentu,
    enabled: wlaczone,
  });
  if (wynik.udany) return;
  oglos('Studio', nazwijOdmowe('Śledzenie zmian', wynik.blad), 'ostrzezenie');
}

export async function wypelnijZmianySledzone(
  kanal: Kanal,
  cel: HTMLElement,
  wzory: WzoryZmian,
  idDokumentu: string,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.StudioTrackingList, { documentId: idDokumentu });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Studio', nazwijOdmowe('Wykaz zmian śledzonych', wynik.blad), 'ostrzezenie');
    return;
  }
  const zmiany = wynik.wynik.changes;
  for (const zmiana of zmiany) {
    const wiersz = wierszZmiany(wzory, zmiana);
    if (wiersz === null) continue;
    cel.appendChild(wiersz);
    if (zmiana.decision !== StudioChangeDecision.Oczekuje) continue;
    const decyzja = wierszDecyzji(kanal, wzory, idDokumentu, zmiana.id, odswiez);
    if (decyzja !== null) cel.appendChild(decyzja);
  }
  const nota = opiszNote(wzory.nota, zmiany);
  if (nota !== null) cel.appendChild(nota);
}

function wierszZmiany(wzory: WzoryZmian, zmiana: StudioTrackedChange): HTMLElement | null {
  if (zmiana.kind === StudioChangeKind.Wstawienie) {
    return wypelnijWiersz(wzory.dodany, '.dn-diff-dod', zmiana.after);
  }
  if (zmiana.kind === StudioChangeKind.Usuniecie) {
    return wypelnijWiersz(wzory.usuniety, '.dn-diff-usu', zmiana.before);
  }
  return wypelnijWiersz(wzory.formatowanie, '.dn-diff-fmt', zmiana.after ?? zmiana.before);
}

function wypelnijWiersz(
  wzor: HTMLElement | null,
  selektor: string,
  tresc: string | undefined,
): HTMLElement | null {
  if (wzor === null || tresc === undefined) return null;
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const miejsce = wiersz.querySelector(selektor);
  if (miejsce === null) return null;
  miejsce.textContent = tresc;
  return wiersz;
}

function wierszDecyzji(
  kanal: Kanal,
  wzory: WzoryZmian,
  idDokumentu: string,
  idZmiany: string,
  odswiez: () => void,
): HTMLElement | null {
  if (wzory.decyzja === null) return null;
  const wiersz = wzory.decyzja.cloneNode(true) as HTMLElement;
  const przyjmij = wiersz.querySelector('.dn-btn--zarys');
  const odrzuc = wiersz.querySelector('.dn-btn--duch');
  if (przyjmij === null || odrzuc === null) return null;
  przyjmij.addEventListener('click', () => {
    void rozstrzygnij(kanal, idDokumentu, idZmiany, true, odswiez);
  });
  odrzuc.addEventListener('click', () => {
    void rozstrzygnij(kanal, idDokumentu, idZmiany, false, odswiez);
  });
  return wiersz;
}

async function rozstrzygnij(
  kanal: Kanal,
  idDokumentu: string,
  idZmiany: string,
  przyjmij: boolean,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.StudioTrackingDecide, {
    documentId: idDokumentu,
    changeIds: [idZmiany],
    accept: przyjmij,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Studio', nazwijOdmowe('Decyzja o zmianie', wynik.blad), 'ostrzezenie');
    return;
  }
  odswiez();
}

function opiszNote(wzor: HTMLElement | null, zmiany: StudioTrackedChange[]): HTMLElement | null {
  if (wzor === null) return null;
  const nota = wzor.cloneNode(true) as HTMLElement;
  const czekajace = zmiany.filter(
    (zmiana) => zmiana.decision === StudioChangeDecision.Oczekuje,
  ).length;
  nota.textContent = `Zmian: ${zmiany.length} · czeka na decyzję: ${czekajace}`;
  return nota;
}
