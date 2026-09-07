// Wiązanie karty modułu Apps z rdzeniem: wykaz wdrożeń aplikacji. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type AppDeployment,
  type AppEndpoint,
  type AppMilestone,
  type AppRoute,
  type AppArchitectureVersion,
  type AppArtifact,
  type AppEnvironment,
  type AppProductLink,
  type AppStage,
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
import { zwiazCzynnosciAplikacji } from './apps-czynnosci.ts';

const KOD_MODULU = 'apps';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazAplikacje(
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
  zwolnijAplikacje(idKarty);

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
    await Promise.all([      wypelnijWdrozenia(kanal, korzen, idOkna),
      wykaz(kanal, korzen, idOkna, 'panel-plan', Command.AppsStageList,
        'Rdzeń nie ma etapów tej aplikacji.',
        (o) => (o.stages as AppStage[]).map((e) => [e.name, ''] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-zadania', Command.AppsMilestoneList,
        'Rdzeń nie ma kamieni milowych tej aplikacji.',
        (o) => (o.milestones as AppMilestone[]).map((k) => [k.name, ''] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-frontend', Command.AppsRouteList,
        'Rdzeń nie ma dróg widoku tej aplikacji.',
        (o) => (o.routes as AppRoute[]).map((d) => [d.path, d.viewName ?? ''] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-architektura', Command.AppsArchitectureVersionList,
        'Rdzeń nie ma wersji architektury tej aplikacji.',
        (o) => (o.versions as AppArchitectureVersion[]).map((w) => [String(w.version), String(w.componentCount ?? '')] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-pliki', Command.AppsArtifactList,
        'Rdzeń nie ma wytworów tej aplikacji.',
        (o) => (o.artifacts as AppArtifact[]).map((a) => [a.kind, a.id] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-terminal', Command.AppsEnvironmentList,
        'Rdzeń nie ma środowisk tej aplikacji.',
        (o) => (o.environments as AppEnvironment[]).map((s) => [s.code, ''] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-builder', Command.AppsProductLinkList,
        'Rdzeń nie ma powiązań z modułami.',
        (o) => (o.links as AppProductLink[]).map((l) => [l.moduleCode, l.detail ?? ''] as const)),
      wykaz(kanal, korzen, idOkna, 'panel-backend', Command.AppsEndpointList,
        'Rdzeń nie ma końcówek tej aplikacji.',
        (o) => (o.endpoints as AppEndpoint[]).map((k) => [k.path, k.method] as const))]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Apps', idOkna);
    if (idOkna === '') return;
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-deployment .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void odswiez();
  }, przy);

  zwiazCzynnosciAplikacji(kanal, korzen, () => idOkna, () => {
    void odswiez();
  }, (zdejmij) => {
    odlaczenia.push(zdejmij);
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Apps');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijAplikacje(idKarty: string): void {
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
  niegotowe(panel(korzen, 'panel-deployment'), 'Wykaz czeka na odpowiedź rdzenia.');
  /* Tytuły paneli niosą nazwę aplikacji wymyśloną na pokaz; rdzeń poda swoją. */
  for (const tytul of korzen.querySelectorAll('.sta-okno-tytul b')) {
    tytul.textContent = (tytul.textContent ?? '').split(' — ')[0] ?? '';
  }
  niegotowe(panel(korzen, 'panel-builder'), 'Budowniczy czeka na wskazanie aplikacji.');
  niegotowe(panel(korzen, 'panel-architektura'), 'Rdzeń nie podaje architektury tej aplikacji.');
  niegotowe(panel(korzen, 'panel-frontend'), 'Rdzeń nie podaje warstwy widoku tej aplikacji.');
  niegotowe(panel(korzen, 'panel-backend'), 'Rdzeń nie podaje warstwy rdzenia tej aplikacji.');
  niegotowe(panel(korzen, 'panel-terminal'), 'Terminal aplikacji nie jest jeszcze wystawiony przez rdzeń.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan budowy czeka na pierwsze zadanie.');
  niegotowe(panel(korzen, 'panel-pliki'), 'Rdzeń nie podaje plików tej aplikacji.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań dla tego okna.');
}

async function wypelnijWdrozenia(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-deployment');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.AppsDeploymentList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz wdrożeń nie doszedł.');
    return;
  }
  const wdrozenia: AppDeployment[] = odpowiedz.wynik.deployments;
  if (wdrozenia.length === 0) {
    niegotowe(cialo, 'Żadne wdrożenie nie ruszyło z tego okna.');
    return;
  }
  cialo.replaceChildren(...wdrozenia.map((w) => pozycja(cialo, w.version ?? w.id, w.status)));
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
