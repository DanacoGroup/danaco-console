/**
 * Wiązanie panelu planu okna Studia z rdzeniem. Znacznik należy do Właściciela —
 * ten plik nic nie buduje: pyta rdzeń o rozkład zlecenia, wypełnia nim stojące
 * węzły, a wiersze zadań powiela z wzorów zdjętych z treści przykładowej.
 */

import {
  Command,
  EventType,
  StudioTaskState,
  type StudioDocumentTask,
  type StudioTaskPlan,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Węzły panelu planu, na których wiązanie pracuje. Brak któregokolwiek znaczy, że panel nie stoi w dokumencie. */
interface WezlyPlanu {
  znacznik: HTMLElement;
  etykieta: HTMLElement;
  lista: HTMLElement;
}

/** Wzory wiersza zdjęte z treści przykładowej, po jednym na stan, któremu znacznik nadaje wygląd. */
interface WzoryWierszy {
  wykonane: HTMLElement | null;
  wBiegu: HTMLElement | null;
  czekajace: HTMLElement | null;
}

/**
 * Wiąże panel planu z rdzeniem. Zwraca prawdę, gdy węzły panelu stały i
 * wiązanie zostało założone; fałsz, gdy panelu w dokumencie nie ma.
 */
export function zwiazPlan(kanal: Kanal, idOkna: string): boolean {
  const znalezione = zbierzWezly();
  if (znalezione === null) return false;
  const wezly: WezlyPlanu = znalezione;

  const wzory = zdejmijWzoryWierszy(wezly.lista);
  const przedrostek = przedrostekZlecenia(wezly.etykieta);

  // Treść przykładowa schodzi, zanim padnie pierwsza odpowiedź rdzenia: panel
  // pusty jest uczciwy, panel z cudzym rozkładem nie.
  zdejmijWiersze(wezly.lista);
  wezly.etykieta.textContent = '';
  wezly.znacznik.textContent = '';

  /** Nanosi rozkład rdzenia na panel; węzły, których nie ma czym wypełnić, schodzą. */
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

  /* Rozkład nie ma własnego zdarzenia kanału. Zmiana dokumentu tego okna jest
     jedyną chwilą, o której rdzeń zawiadamia, więc po niej panel pyta o rozkład
     ponownie. */
  kanal.naZdarzenie(EventType.StudioDocumentChanged, (tresc) => {
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
  return true;
}

/** Pyta rdzeń o rozkład wskazanego dokumentu; odmowa rdzenia oddaje pustkę. */
async function wczytajRozklad(kanal: Kanal, idDokumentu: string): Promise<StudioTaskPlan | null> {
  const wynik = await wywolaj(kanal, Command.StudioPlanGet, { documentId: idDokumentu });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.plan;
}

/** Wskazuje dokument otwarty w oknie; rozkład idzie po dokumencie, a wiązanie dostaje identyfikator okna. */
async function wskazDokument(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  return wynik.wynik.document.id;
}

/** Wskazuje węzły panelu planu; pustka znaczy, że panel nie stoi w dokumencie. */
function zbierzWezly(): WezlyPlanu | null {
  const panel = document.querySelector('#panel-plan');
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

/** Zdejmuje z listy po jednym wierszu na stan; klon zachowuje znacznik stanu i glif nadane przez Właściciela. */
function zdejmijWzoryWierszy(lista: HTMLElement): WzoryWierszy {
  return {
    wykonane: sklonuj(wskazWiersz(lista, '.dn-kropka--sukces')),
    wBiegu: sklonuj(wskazWiersz(lista, '.pt-tetno')),
    czekajace: sklonuj(wskazWiersz(lista, '.dn-kropka--neutralna')),
  };
}

/** Pierwszy wiersz listy niosący wskazany znacznik stanu. */
function wskazWiersz(lista: HTMLElement, selektor: string): HTMLElement | null {
  for (const wiersz of lista.querySelectorAll('.st-panel-wiersz')) {
    if (wiersz instanceof HTMLElement && wiersz.querySelector(selektor) !== null) return wiersz;
  }
  return null;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

/** Zdejmuje wiersze zadań, zostawiając etykietę zlecenia, która stoi w tym samym pojemniku. */
function zdejmijWiersze(lista: HTMLElement): void {
  for (const wiersz of lista.querySelectorAll('.st-panel-wiersz')) wiersz.remove();
}

/** Stały początek etykiety wzięty ze znacznika; po nim wchodzi zlecenie Operatora zwrócone przez rdzeń. */
function przedrostekZlecenia(etykieta: HTMLElement): string {
  const tresc = etykieta.textContent ?? '';
  const koniec = tresc.indexOf(':');
  return koniec === -1 ? '' : tresc.slice(0, koniec + 1);
}

/** Zwraca klon wzoru właściwego stanowi zadania, opisany jego nazwą. */
function zbudujWiersz(wzory: WzoryWierszy, zadanie: StudioDocumentTask): HTMLElement | null {
  const pokryty = wzorDlaStanu(wzory, zadanie.state);
  const wzor = pokryty ?? wzory.czekajace;
  if (wzor === null) return null;
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  if (pokryty === null) {
    /* Stanom failed, blocked i skipped znacznik nie nadaje ani znacznika stanu,
       ani glifu, a wygląd dobrany poza znacznikiem mówiłby Operatorowi rzecz
       nierozstrzygniętą; zostaje sama nazwa zadania. */
    wiersz.querySelector('.dn-kropka')?.remove();
    wiersz.querySelector('.dn-meta')?.remove();
  }
  wpiszNazwe(wiersz, zadanie.title);
  return wiersz;
}

/** Wzór wiersza właściwy stanowi zadania; pustka znaczy stan, któremu znacznik wyglądu nie nadaje. */
function wzorDlaStanu(wzory: WzoryWierszy, stan: StudioTaskState): HTMLElement | null {
  if (stan === StudioTaskState.Done) return wzory.wykonane;
  if (stan === StudioTaskState.Running) return wzory.wBiegu;
  if (stan === StudioTaskState.Pending) return wzory.czekajace;
  return null;
}

/** Wpisuje nazwę zadania w goły węzeł tekstowy wiersza — miejsce, które w znaczniku niesie nazwę zadania. */
function wpiszNazwe(wiersz: HTMLElement, nazwa: string): void {
  for (const wezel of wiersz.childNodes) {
    if (wezel.nodeType !== Node.TEXT_NODE || (wezel.nodeValue ?? '').trim() === '') continue;
    wezel.nodeValue = ` ${nazwa} `;
    return;
  }
}

/** Ile zadań rozkładu jest domkniętych — licznik przy nazwie karty liczy się z zadań, bo kontrakt osobnego pola nie niesie. */
function policzWykonane(zadania: StudioDocumentTask[]): number {
  return zadania.filter((zadanie) => zadanie.state === StudioTaskState.Done).length;
}
