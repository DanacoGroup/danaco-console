// Wiązanie karty modułu Automations z rdzeniem: wykaz automatyk. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type AutomationAlertRule,
  type AutomationAuditEntry,
  type AutomationSecretRef,
  type AutomationTemplate,
  type AutomationWorkflow,
  type Queue,
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
import { zwiazCzynnosciAutomatyk } from './automations-czynnosci.ts';

const KOD_MODULU = 'automations';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazAutomatyzacje(
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
  zwolnijAutomatyzacje(idKarty);

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
  /* Czynności wierszy: nasłuch stoi na korzeniu okna, więc przetrwa każde
     odświeżenie wykazu automatyk. */
  zwiazCzynnosciAutomatyk(kanal, korzen, () => {
    void wypelnijAutomatyki(kanal, korzen);
  });

  let idOkna = idOknaStojacego;
  const odswiez = async (): Promise<void> => {
    if (idOkna === '') return;
    await Promise.all([      wypelnijAutomatyki(kanal, korzen),
      wykaz(kanal, korzen, 'panel-builder', Command.AutomationTemplateList,
        'Rdzeń nie ma wzorców automatyki.',
        (o) => (o.templates as AutomationTemplate[]).map((s) => [s.name, s.description ?? ''] as const)),
      wykaz(kanal, korzen, 'panel-monitor', Command.AutomationAlertRuleList,
        'Żadna reguła powiadamiania nie została ustawiona.',
        (o) => (o.rules as AutomationAlertRule[]).map((r) => [r.trigger, r.condition ?? ''] as const)),
      wykaz(kanal, korzen, 'panel-scheduler', Command.AutomationSecretList,
        'Rdzeń nie ma poświadczeń automatyki.',
        (o) => (o.secrets as AutomationSecretRef[]).map((s) => [s.name, s.scope ?? ''] as const)),
      wykaz(kanal, korzen, 'panel-queue', Command.QueueList,
        'Żadna kolejka zadań nie stoi.',
        (o) => (o.queues as Queue[]).map((k) => [k.name ?? k.id, k.status] as const)),
      wykaz(kanal, korzen, 'panel-artefakty', Command.AutomationAuditList,
        'Rdzeń nie odnotował czynności automatyk.',
        (o) => (o.entries as AutomationAuditEntry[]).map((w) => [w.action, w.actorId ?? ''] as const))]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Automations', idOkna);
    if (idOkna === '') return;
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-orch .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void odswiez();
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Automations');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijAutomatyzacje(idKarty: string): void {
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
  niegotowe(panel(korzen, 'panel-orch'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-builder'), 'Budowniczy czeka na wskazanie automatyki.');
  niegotowe(panel(korzen, 'panel-scheduler'), 'Rdzeń nie podaje harmonogramu dla tego okna.');
  niegotowe(panel(korzen, 'panel-monitor'), 'Monitor czeka na pierwszy przebieg.');
  niegotowe(panel(korzen, 'panel-queue'), 'Rdzeń nie podaje kolejki dla tego okna.');
  niegotowe(panel(korzen, 'panel-terminal'), 'Terminal automatyki nie jest jeszcze wystawiony przez rdzeń.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan automatyki czeka na pierwszy krok.');
  niegotowe(panel(korzen, 'panel-artefakty'), 'Rdzeń nie podaje wytworów tego okna.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań dla tego okna.');
}

async function wypelnijAutomatyki(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = panel(korzen, 'panel-orch');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.AutomationWorkflowList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz automatyk nie doszedł.');
    return;
  }
  const wykaz: AutomationWorkflow[] = odpowiedz.wynik.workflows;
  if (wykaz.length === 0) {
    niegotowe(cialo, 'Żadna automatyka nie została jeszcze złożona.');
    return;
  }
  cialo.replaceChildren(...wykaz.map((w) => wierszAutomatyki(cialo, w)));
}

/* Panel wykazu bez własnego kształtu: komenda bez pól wymaganych. */
async function wykaz(
  kanal: Kanal,
  korzen: Element,
  panelId: string,
  komenda: Parameters<typeof wywolaj>[1],
  pusty: string,
  mapuj: (wynik: Record<string, unknown>) => readonly (readonly [string, string])[],
): Promise<void> {
  const cialo = cialoPanelu(korzen, panelId);
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, komenda as never, {} as never);
  const wynik = odpowiedz.wynik as Record<string, unknown> | undefined;
  wykazPanelu(cialo, odpowiedz.udany, odpowiedz.blad?.message,
    wynik === undefined ? undefined : [...mapuj(wynik)], pusty, (x) => x);
}

/* Wiersz automatyki niesie jej kod i czynności; sam wykaz zostaje wykazem. */
function wierszAutomatyki(cialo: Element, automatyka: AutomationWorkflow): HTMLElement {
  const wiersz = pozycja(cialo, automatyka.name, automatyka.description ?? '');
  wiersz.dataset.automatyka = automatyka.id;
  const czynnosci: readonly (readonly [string, string])[] = [
    ['uruchom', 'Uruchom'],
    ['zatrzymaj', 'Zatrzymaj'],
    [automatyka.enabled === true ? 'wylacz' : 'wlacz',
      automatyka.enabled === true ? 'Wyłącz' : 'Włącz'],
  ];
  for (const [kod, etykieta] of czynnosci) {
    const przycisk = cialo.ownerDocument.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
    przycisk.dataset.czynnoscAutomatyki = kod;
    przycisk.textContent = etykieta;
    wiersz.append(przycisk);
  }
  return wiersz;
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
