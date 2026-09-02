// Panel planu okna Studia: rozkład zlecenia dokumentowego wraz ze stanem
// każdego zadania i przestawieniem tego stanu wskazaniem wiersza.

import {
  Command,
  EventType,
  StudioTaskState,
  type StudioDocumentTask,
  type StudioTaskPlan,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

interface WezlyPlanu {
  znacznik: HTMLElement;
  etykieta: HTMLElement;
  lista: HTMLElement;
}

interface WzoryWierszy {
  wykonane: HTMLElement | null;
  wBiegu: HTMLElement | null;
  czekajace: HTMLElement | null;
}

interface PanelPlanu {
  wezly: WezlyPlanu;
  wzory: WzoryWierszy;
  przedrostek: string;
}

// Panel opróżniony wzoru nie oddaje, więc zdjęcie stoi raz dla wszystkich kart.
let zdjete: { wzory: WzoryWierszy; przedrostek: string } | null = null;

export function zdejmijTrescPrzykladowaPlanu(korzen: ParentNode): boolean {
  return przygotujPanel(korzen) !== null;
}

function przygotujPanel(korzen: ParentNode): PanelPlanu | null {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return null;
  zdjete ??= {
    wzory: zdejmijWzoryWierszy(wezly.lista),
    przedrostek: przedrostekZlecenia(wezly.etykieta),
  };
  zdejmijWiersze(wezly.lista);
  wezly.etykieta.textContent = '';
  wezly.znacznik.textContent = '';
  return { wezly, wzory: zdjete.wzory, przedrostek: zdjete.przedrostek };
}

export function zwiazPlan(kanal: Kanal, idOkna: string, korzen: ParentNode): Odsubskrybuj | null {
  const przygotowany = przygotujPanel(korzen);
  if (przygotowany === null) return null;
  const wezly: WezlyPlanu = przygotowany.wezly;
  const wzory: WzoryWierszy = przygotowany.wzory;
  const przedrostek: string = przygotowany.przedrostek;

  function pokaz(plan: StudioTaskPlan | null): void {
    zdejmijWiersze(wezly.lista);
    if (plan === null) {
      wezly.etykieta.remove();
      wezly.znacznik.remove();
      return;
    }
    wezly.etykieta.textContent =
      przedrostek === '' ? plan.order : `${przedrostek} ${plan.order}`;
    const zadania = plan.tasks;
    if (zadania === undefined) {
      wezly.znacznik.remove();
      return;
    }
    wezly.znacznik.textContent = `${policzWykonane(zadania)} / ${zadania.length}`;
    for (const zadanie of zadania) {
      const wiersz = zbudujWiersz(wzory, zadanie);
      if (wiersz !== null) wezly.lista.appendChild(wiersz);
    }
  }

  wezly.lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('.st-panel-wiersz[data-zadanie]');
    const idZadania = wiersz?.dataset.zadanie;
    if (idZadania === undefined) return;
    void przestawZadanie(kanal, idZadania, wiersz?.dataset.stan ?? '').then((plan) => {
      if (plan !== null) pokaz(plan);
    });
  });

  // Rozkład nie ma własnego zdarzenia; zmiana dokumentu okna jest jedyną chwilą,
  // o której rdzeń zawiadamia, więc po niej panel pyta o rozkład ponownie.
  const odlacz = kanal.naZdarzenie(EventType.StudioDocumentChanged, (tresc) => {
    if (tresc.document.windowId !== idOkna) return;
    void wczytajRozklad(kanal, tresc.document.id).then(pokaz);
  });

  void wskazDokument(kanal, idOkna).then((wskazany) => {
    if (wskazany === '') {
      pokaz(null);
      return;
    }
    void wczytajRozklad(kanal, wskazany).then(pokaz);
  });
  return odlacz;
}

// Wskazanie wiersza domyka zadanie albo je otwiera z powrotem; stanów failed,
// blocked i skipped panel nie nadaje, bo znacznik nie ma dla nich wyglądu.
async function przestawZadanie(
  kanal: Kanal,
  idZadania: string,
  stan: string,
): Promise<StudioTaskPlan | null> {
  const wynik = await wywolaj(kanal, Command.StudioPlanTaskUpdate, {
    taskId: idZadania,
    state: stan === StudioTaskState.Done ? StudioTaskState.Pending : StudioTaskState.Done,
  });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.plan;
}

async function wczytajRozklad(kanal: Kanal, idDokumentu: string): Promise<StudioTaskPlan | null> {
  const wynik = await wywolaj(kanal, Command.StudioPlanGet, { documentId: idDokumentu });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.plan;
}

async function wskazDokument(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  return wynik.wynik.document.id;
}

function zbierzWezly(korzen: ParentNode): WezlyPlanu | null {
  const panel = korzen.querySelector('#panel-plan');
  const znacznik = panel?.querySelector('.sta-okno-znacznik');
  const etykieta = panel?.querySelector('.pt-etykieta');
  const lista = panel?.querySelector('.sta-okno-tresc.st-panel-lista');
  if (
    !(znacznik instanceof HTMLElement) ||
    !(etykieta instanceof HTMLElement) ||
    !(lista instanceof HTMLElement)
  ) {
    return null;
  }
  return { znacznik, etykieta, lista };
}

function zdejmijWzoryWierszy(lista: HTMLElement): WzoryWierszy {
  return {
    wykonane: sklonuj(wskazWiersz(lista, '.dn-kropka--sukces')),
    wBiegu: sklonuj(wskazWiersz(lista, '.pt-tetno')),
    czekajace: sklonuj(wskazWiersz(lista, '.dn-kropka--neutralna')),
  };
}

function wskazWiersz(lista: HTMLElement, selektor: string): HTMLElement | null {
  for (const wiersz of lista.querySelectorAll('.st-panel-wiersz')) {
    if (wiersz instanceof HTMLElement && wiersz.querySelector(selektor) !== null) return wiersz;
  }
  return null;
}

function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

function zdejmijWiersze(lista: HTMLElement): void {
  for (const wiersz of lista.querySelectorAll('.st-panel-wiersz')) wiersz.remove();
}

function przedrostekZlecenia(etykieta: HTMLElement): string {
  const tresc = etykieta.textContent ?? '';
  const koniec = tresc.indexOf(':');
  return koniec === -1 ? '' : tresc.slice(0, koniec + 1);
}

function zbudujWiersz(wzory: WzoryWierszy, zadanie: StudioDocumentTask): HTMLElement | null {
  const pokryty = wzorDlaStanu(wzory, zadanie.state);
  const wzor = pokryty ?? wzory.czekajace;
  if (wzor === null) return null;
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.zadanie = zadanie.id;
  wiersz.dataset.stan = zadanie.state;
  if (pokryty === null) {
    // Stanom bez wyglądu w znaczniku zostaje sama nazwa zadania.
    wiersz.querySelector('.dn-kropka')?.remove();
    wiersz.querySelector('.dn-meta')?.remove();
  }
  wpiszNazwe(wiersz, zadanie.title);
  return wiersz;
}

function wzorDlaStanu(wzory: WzoryWierszy, stan: StudioTaskState): HTMLElement | null {
  if (stan === StudioTaskState.Done) return wzory.wykonane;
  if (stan === StudioTaskState.Running) return wzory.wBiegu;
  if (stan === StudioTaskState.Pending) return wzory.czekajace;
  return null;
}

function wpiszNazwe(wiersz: HTMLElement, nazwa: string): void {
  for (const wezel of wiersz.childNodes) {
    if (wezel.nodeType !== Node.TEXT_NODE || (wezel.nodeValue ?? '').trim() === '') continue;
    wezel.nodeValue = ` ${nazwa} `;
    return;
  }
}

// Licznik przy nazwie karty liczy się z zadań: kontrakt osobnego pola nie niesie.
function policzWykonane(zadania: StudioDocumentTask[]): number {
  return zadania.filter((zadanie) => zadanie.state === StudioTaskState.Done).length;
}
