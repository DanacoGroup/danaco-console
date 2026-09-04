// Wiązanie karty modułu Developer z rdzeniem: wykaz przebiegów budowania i gałęzi. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type ApiCollection,
  type ContainerInfo,
  type DataConnection,
  type DependencyNode,
  type DeveloperBuild,
  type GitBranch,
  type ScanFinding,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  cialoPanelu,
  opiszNaglowek,
  zapewnijOknoModulu,
  zdejmijSterowanieWspolne,
  wykazPanelu,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

const KOD_MODULU = 'developer';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazDevelopera(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijDevelopera(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  if (korzen instanceof HTMLElement) {
    zdejmijTrescWspolna(korzen);
    zdejmijSterowanieWspolne(korzen);
  }
  zdejmijTrescPrzykladowa(korzen);
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);

  let idOkna = idOknaStojacego;
  const odswiez = async (): Promise<void> => {
    if (idOkna === '') return;
    await Promise.all([      wypelnijBudowania(kanal, korzen, idOkna),
      wypelnijGalezie(kanal, korzen, idOkna),
      wykaz(kanal, korzen, idOkna, 'panel-drzewo', Command.DeveloperDependencyList,
        'Rdzeń nie widzi zależności w katalogu roboczym.',
        (o) => (o.dependencies as DependencyNode[]).map((d) => [d.name, d.version] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-terminal', Command.DeveloperContainerList,
        'Żaden pojemnik nie stoi dla tego okna.',
        (o) => (o.containers as ContainerInfo[]).map((c) => [c.name, c.status] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-artefakty', Command.DeveloperApiCollectionList,
        'Rdzeń nie ma kolekcji zapytań dla tego okna.',
        (o) => (o.collections as ApiCollection[]).map((k) => [k.name, ''] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-zadania', Command.DeveloperDataConnectionList,
        'Rdzeń nie ma połączeń z bazami dla tego okna.',
        (o) => (o.connections as DataConnection[]).map((p) => [p.name, p.engine] as const)),
      wykaz(kanal, korzen, '', 'panel-plan', Command.DeveloperScanResultList,
        'Rdzeń nie odnotował zgłoszeń przeglądu bezpieczeństwa.',
        (o) => (o.findings as ScanFinding[]).map((z) => [z.title, z.severity] as const))]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Developer', idOkna);
    if (idOkna === '') return;
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-build .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void odswiez();
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Developer');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijDevelopera(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

function panel(korzen: Element, identyfikator: string): Element | null {
  return korzen.querySelector(`#${identyfikator} .sta-okno-tresc`);
}

function niegotowe(cialo: Element | null, zdanie: string): void {
  if (cialo === null) return;
  const napis = cialo.ownerDocument.createElement('div');
  napis.className = 'dn-meta';
  napis.textContent = zdanie;
  cialo.replaceChildren(napis);
}

/* Znacznik niesie źródła, ustalenia i miary wpisane wprost. Schodzą przed
   pierwszym pytaniem rdzenia, żeby okno nie pokazywało cudzego badania. */
function zdejmijTrescPrzykladowa(korzen: Element): void {
  niegotowe(panel(korzen, 'panel-build'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-git'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-drzewo'), 'Rdzeń nie podaje drzewa plików dla tego okna.');
  niegotowe(panel(korzen, 'panel-edytor'), 'Edytor czeka na wskazanie pliku.');
  // Zakładki otwartych plików i miary kursora prototyp wpisuje wprost.
  for (const karta of korzen.querySelectorAll('.dv-karta')) karta.remove();
  niegotowe(panel(korzen, 'panel-terminal'), 'Terminal projektu nie jest jeszcze wystawiony przez rdzeń.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan pracy czeka na pierwsze zadanie.');
  niegotowe(panel(korzen, 'panel-artefakty'), 'Rdzeń nie podaje wytworów tego okna.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań dla tego okna.');
}

async function wypelnijBudowania(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-build');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.DeveloperBuildList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz budowań nie doszedł.');
    return;
  }
  const przebiegi: DeveloperBuild[] = odpowiedz.wynik.builds;
  if (przebiegi.length === 0) {
    niegotowe(cialo, 'Żadne budowanie nie ruszyło w tym oknie.');
    return;
  }
  cialo.replaceChildren(...przebiegi.map((b) => pozycja(cialo, b.task, b.status)));
}

async function wypelnijGalezie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-git');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.DeveloperGitBranchList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz gałęzi nie doszedł.');
    return;
  }
  const galezie: GitBranch[] = odpowiedz.wynik.branches;
  if (galezie.length === 0) {
    niegotowe(cialo, 'Rdzeń nie widzi repozytorium w katalogu roboczym tego okna.');
    return;
  }
  cialo.replaceChildren(...galezie.map((g) => pozycja(cialo, g.name, g.remote ?? '')));
}

/* Panel wykazu bez własnego kształtu; okno puste znaczy komendę bez pola okna. */
async function wykaz(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  panelId: string,
  komenda: Parameters<typeof wywolaj>[1],
  pusty: string,
  mapuj: (wynik: Record<string, unknown>) => readonly (readonly [string, string])[],
): Promise<void> {
  const cialo = cialoPanelu(korzen, panelId);
  if (cialo === null) return;
  const zadanie = idOkna === '' ? {} : { windowId: idOkna };
  const odpowiedz = await wywolaj(kanal, komenda as never, zadanie as never);
  const wynik = odpowiedz.wynik as Record<string, unknown> | undefined;
  wykazPanelu(cialo, odpowiedz.udany, odpowiedz.blad?.message,
    wynik === undefined ? undefined : [...mapuj(wynik)], pusty, (x) => x);
}

function pozycja(cialo: Element, tytul: string, podpis: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  const nazwa = cialo.ownerDocument.createElement('span');
  nazwa.textContent = tytul;
  wiersz.append(nazwa);
  if (podpis !== '') {
    const meta = cialo.ownerDocument.createElement('span');
    meta.className = 'dn-meta';
    meta.textContent = podpis;
    wiersz.append(meta);
  }
  return wiersz;
}

/* Okna zakładania źródła wydanie nie niesie: tytuł wchodzi zapytaniem w oknie,
   tak samo jak nazwa komponentu w Centrum. */
