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
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
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

/** Panel gotowy do wypełnienia: węzły znacznika oraz wzory i przedrostek zdjęte z treści przykładowej. */
interface PanelPlanu {
  wezly: WezlyPlanu;
  wzory: WzoryWierszy;
  przedrostek: string;
}

/**
 * Wzory i stały początek etykiety zdjęte przy pierwszym montażu karty. Każda
 * karta niesie ten sam znacznik, a panel opróżniony z treści przykładowej
 * wzoru już nie oddaje, więc zdjęcie stoi raz dla wszystkich kart.
 */
let zdjete: { wzory: WzoryWierszy; przedrostek: string } | null = null;

/**
 * Zdejmuje treść przykładową panelu planu, zabierając z niej wzory wierszy
 * i stały początek etykiety zlecenia. Woła się przy montażu okna, przed
 * powstaniem stanowiska: rozkład z prototypu jest cudzą pracą i Operator nie
 * ma prawa wziąć go za swoją. Zwraca prawdę, gdy panel stał w dokumencie.
 */
export function zdejmijTrescPrzykladowaPlanu(korzen: ParentNode): boolean {
  return przygotujPanel(korzen) !== null;
}

/** Zbiera węzły panelu i opróżnia je z treści przykładowej; pustka znaczy panel poza kartą. */
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

/**
 * Wiąże panel planu karty z rdzeniem; węzły idą od korzenia karty. Zwraca
 * odłączenie nasłuchu, a pustkę, gdy panelu w karcie nie ma.
 */
export function zwiazPlan(kanal: Kanal, idOkna: string, korzen: ParentNode): Odsubskrybuj | null {
  const przygotowany = przygotujPanel(korzen);
  if (przygotowany === null) return null;
  const wezly: WezlyPlanu = przygotowany.wezly;
  const wzory: WzoryWierszy = przygotowany.wzory;
  const przedrostek: string = przygotowany.przedrostek;

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

/** Wskazuje węzły panelu planu od korzenia karty; pustka znaczy, że panel nie stoi w karcie. */
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
