// Panel różnic okna Studia: porównanie treści i postaci wersji dokumentu oraz
// wykaz zmian śledzonych, wedle trybu wskazanego na pasku panelu.

import {
  ChangeKind,
  Command,
  DiffHunkKind,
  EventType,
  type StudioDiffHunk,
  type StudioFormDiffEntry,
  type StudioVersion,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { przestawSledzenie, wypelnijZmianySledzone } from './studio-sledzenie.ts';

const TRYB_TRESCI = 0;
const TRYB_POSTACI = 1;
const TRYB_SLEDZENIA = 2;
const NAZWY_TRYBU = ['różnica treści', 'różnica postaci', 'zmiany śledzone'];

interface WezlyRoznic {
  tresc: HTMLElement;
  wierszWersji: HTMLElement;
  wersjaOdniesienia: HTMLElement;
  wersjaPorownywana: HTMLElement;
  tryb: HTMLElement | null;
  poleWzorca: HTMLElement;
  wzorzec: HTMLInputElement;
}

interface WzoryWierszy {
  usuniety: HTMLElement | null;
  dodany: HTMLElement | null;
  formatowanie: HTMLElement | null;
  decyzja: HTMLElement | null;
  nota: HTMLElement | null;
}

interface WyborPorownania {
  idDokumentu: string;
  wersje: StudioVersion[];
  odniesienie: number;
  porownywana: number;
  tryb: number;
}

interface Panel {
  wezly: WezlyRoznic;
  wzory: WzoryWierszy;
  wybor: WyborPorownania;
}

interface Zwiazany {
  kanal: Kanal;
  idOkna: string;
  panel: Panel;
  odswiez: () => void;
}

// Panel opróżniony wzoru już nie oddaje, więc zdjęcie stoi raz dla wszystkich kart.
let wzoryPanelu: WzoryWierszy | null = null;

// Porównanie ze wstążki ma trafić w panel własnej karty, nie karty sąsiedniej.
const ZWIAZANE = new Map<string, Zwiazany>();

export function zdejmijTrescPrzykladowaRoznic(korzen: ParentNode): boolean {
  return przygotujPanel(korzen) !== null;
}

function przygotujPanel(korzen: ParentNode): { wezly: WezlyRoznic; wzory: WzoryWierszy } | null {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return null;
  wzoryPanelu ??= zdejmijWzory(wezly);
  zdejmijTrescPrzykladowa(wezly);
  return { wezly, wzory: wzoryPanelu };
}

export async function porownajWersjeDokumentu(idKarty: string): Promise<boolean> {
  const biezacy = ZWIAZANE.get(idKarty);
  if (biezacy === undefined) return false;
  if (biezacy.panel.wybor.idDokumentu === '') {
    const otwarty = await wywolaj(biezacy.kanal, Command.StudioDocumentOpen, {
      windowId: biezacy.idOkna,
    });
    if (!otwarty.udany || otwarty.wynik === undefined) return false;
    biezacy.panel.wybor.idDokumentu = otwarty.wynik.document.id;
  }
  await wczytaj(biezacy.kanal, biezacy.panel, biezacy.odswiez);
  return true;
}

// Dokument bierze się ze zdarzenia zmiany: komendy pytającej o dokument okna nie ma.
export function zwiazRoznice(
  kanal: Kanal,
  idOkna: string,
  idKarty: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const przygotowany = przygotujPanel(korzen);
  if (przygotowany === null) return null;
  const wezly: WezlyRoznic = przygotowany.wezly;

  const panel: Panel = {
    wezly,
    wzory: przygotowany.wzory,
    wybor: { idDokumentu: '', wersje: [], odniesienie: 0, porownywana: 0, tryb: TRYB_TRESCI },
  };
  opiszTryb(panel);

  function odswiez(): void {
    if (panel.wybor.idDokumentu === '') return;
    void porownaj(kanal, panel, odswiez);
  }

  wezly.wersjaOdniesienia.addEventListener('click', () => {
    panel.wybor.odniesienie = nastepneMiejsce(panel.wybor.odniesienie, panel.wybor.wersje.length);
    opiszPasek(panel);
    odswiez();
  });

  wezly.wersjaPorownywana.addEventListener('click', () => {
    panel.wybor.porownywana = nastepneMiejsce(panel.wybor.porownywana, panel.wybor.wersje.length);
    opiszPasek(panel);
    odswiez();
  });

  wezly.wzorzec.addEventListener('change', () => {
    odswiez();
  });

  wezly.tryb?.addEventListener('click', () => {
    const poprzedni = panel.wybor.tryb;
    panel.wybor.tryb = (poprzedni + 1) % NAZWY_TRYBU.length;
    opiszTryb(panel);
    if (panel.wybor.idDokumentu !== '' && poprzedni !== panel.wybor.tryb) {
      const wchodzi = panel.wybor.tryb === TRYB_SLEDZENIA;
      if (wchodzi || poprzedni === TRYB_SLEDZENIA) {
        void przestawSledzenie(kanal, panel.wybor.idDokumentu, wchodzi);
      }
    }
    odswiez();
  });

  const odlacz = kanal.naZdarzenie(EventType.StudioDocumentChanged, (tresc) => {
    if (tresc.document.windowId !== idOkna) return;
    if (tresc.change === ChangeKind.Deleted) {
      panel.wybor.idDokumentu = '';
      panel.wybor.wersje = [];
      zdejmijPozycje(wezly);
      return;
    }
    panel.wybor.idDokumentu = tresc.document.id;
    void wczytaj(kanal, panel, odswiez);
  });

  ZWIAZANE.set(idKarty, { kanal, idOkna, panel, odswiez });
  return () => {
    odlacz();
    if (ZWIAZANE.get(idKarty)?.panel === panel) ZWIAZANE.delete(idKarty);
  };
}

async function wczytaj(kanal: Kanal, panel: Panel, odswiez: () => void): Promise<void> {
  const wynik = await wywolaj(kanal, Command.StudioRepositoryList, {
    documentId: panel.wybor.idDokumentu,
  });
  if (wynik.udany && wynik.wynik !== undefined) {
    // Wersja bez nazwy nie wchodzi do wykazu paska: przełącznik nie ma jej czym nazwać.
    panel.wybor.wersje = wynik.wynik.versions.filter((wersja) => (wersja.label ?? '') !== '');
    panel.wybor.porownywana = 0;
    panel.wybor.odniesienie = Math.min(1, panel.wybor.wersje.length - 1);
  }
  opiszPasek(panel);
  await porownaj(kanal, panel, odswiez);
}

async function porownaj(kanal: Kanal, panel: Panel, odswiez: () => void): Promise<void> {
  if (panel.wybor.tryb === TRYB_SLEDZENIA) {
    zdejmijPozycje(panel.wezly);
    await wypelnijZmianySledzone(
      kanal, panel.wezly.tresc, panel.wzory, panel.wybor.idDokumentu, odswiez,
    );
    return;
  }
  if (panel.wybor.tryb === TRYB_POSTACI) {
    await porownajPostac(kanal, panel);
    return;
  }
  const wzorzec = panel.wezly.wzorzec.value.trim();
  const odniesienie = idWersji(panel.wybor, panel.wybor.odniesienie);
  const porownywana = idWersji(panel.wybor, panel.wybor.porownywana);
  const wynik = await wywolaj(kanal, Command.StudioDiffCompare, {
    documentId: panel.wybor.idDokumentu,
    baseVersionId: odniesienie === '' ? undefined : odniesienie,
    targetVersionId: porownywana === '' ? undefined : porownywana,
    pattern: wzorzec === '' ? undefined : wzorzec,
    // Pole wzorca niesie w znaczniku podpis wyrażenia regularnego, więc wzorzec idzie do rdzenia jako wyrażenie.
    regex: wzorzec === '' ? undefined : true,
  });
  zdejmijPozycje(panel.wezly);
  if (!wynik.udany || wynik.wynik === undefined) return;
  /* Trafienia wzorca odpowiedź niesie osobnym polem, ale panel nie ma dla nich
     ani jednego węzła, więc wykaz wypełniają wyłącznie fragmenty różnicy. */
  wypelnijFragmenty(kanal, panel, wynik.wynik.hunks ?? [], odswiez);
}

// Postać zestawiona osobno: zmiana kroju czy wcięcia milczy w różnicy treści.
async function porownajPostac(kanal: Kanal, panel: Panel): Promise<void> {
  const odniesienie = idWersji(panel.wybor, panel.wybor.odniesienie);
  const porownywana = idWersji(panel.wybor, panel.wybor.porownywana);
  const wynik = await wywolaj(kanal, Command.StudioDiffFormCompare, {
    documentId: panel.wybor.idDokumentu,
    baseVersionId: odniesienie === '' ? undefined : odniesienie,
    targetVersionId: porownywana === '' ? undefined : porownywana,
  });
  zdejmijPozycje(panel.wezly);
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const pozycja of wynik.wynik.entries) {
    const wiersz = wierszPostaci(panel.wzory, pozycja);
    if (wiersz !== null) panel.wezly.tresc.appendChild(wiersz);
  }
  const nota = panel.wzory.nota;
  if (nota === null) return;
  const opis = nota.cloneNode(true) as HTMLElement;
  const rachunek = wynik.wynik;
  opis.textContent = `Cech doszło: ${rachunek.added} · odpadło: ${rachunek.removed}`
    + ` · zmienionych: ${rachunek.changed}`;
  panel.wezly.tresc.appendChild(opis);
}

function wierszPostaci(wzory: WzoryWierszy, pozycja: StudioFormDiffEntry): HTMLElement | null {
  return wypelnijWiersz(wzory.formatowanie, '.dn-diff-fmt', `${pozycja.area}: ${pozycja.detail}`);
}

function wypelnijFragmenty(
  kanal: Kanal,
  panel: Panel,
  fragmenty: StudioDiffHunk[],
  odswiez: () => void,
): void {
  for (const fragment of fragmenty) {
    const wiersze = wierszeFragmentu(panel.wzory, fragment);
    if (wiersze.length === 0) continue;
    for (const wiersz of wiersze) panel.wezly.tresc.appendChild(wiersz);
    const decyzja = wierszDecyzji(kanal, panel, fragment.index, odswiez);
    if (decyzja !== null) panel.wezly.tresc.appendChild(decyzja);
  }
}

// Fragment podany dla kontekstu odpada, bo wzoru dla niego znacznik nie niesie.
function wierszeFragmentu(wzory: WzoryWierszy, fragment: StudioDiffHunk): HTMLElement[] {
  const wiersze: HTMLElement[] = [];
  if (fragment.kind === DiffHunkKind.Removed || fragment.kind === DiffHunkKind.Changed) {
    const ubytek = wypelnijWiersz(wzory.usuniety, '.dn-diff-usu', fragment.before);
    if (ubytek !== null) wiersze.push(ubytek);
  }
  if (fragment.kind === DiffHunkKind.Added || fragment.kind === DiffHunkKind.Changed) {
    const przyrost = wypelnijWiersz(wzory.dodany, '.dn-diff-dod', fragment.after);
    if (przyrost !== null) wiersze.push(przyrost);
  }
  return wiersze;
}

function wypelnijWiersz(wzor: HTMLElement | null, selektor: string, tresc?: string): HTMLElement | null {
  if (wzor === null || tresc === undefined) return null;
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const miejsce = wiersz.querySelector(selektor);
  if (miejsce === null) return null;
  miejsce.textContent = tresc;
  return wiersz;
}

// Odrzucenie fragmentu schodzi z wiersza: komendę ma wyłącznie przeniesienie.
function wierszDecyzji(
  kanal: Kanal,
  panel: Panel,
  numer: number,
  odswiez: () => void,
): HTMLElement | null {
  if (panel.wzory.decyzja === null) return null;
  if (idWersji(panel.wybor, panel.wybor.porownywana) === '') return null;
  const wiersz = panel.wzory.decyzja.cloneNode(true) as HTMLElement;
  wiersz.querySelector('.dn-btn--duch')?.remove();
  const przyjmij = wiersz.querySelector('button');
  if (przyjmij === null) return null;
  przyjmij.addEventListener('click', () => {
    void przyjmijFragment(kanal, panel, numer, odswiez);
  });
  return wiersz;
}

async function przyjmijFragment(
  kanal: Kanal,
  panel: Panel,
  numer: number,
  odswiez: () => void,
): Promise<void> {
  const zrodlo = idWersji(panel.wybor, panel.wybor.porownywana);
  if (zrodlo === '') return;
  const wynik = await wywolaj(kanal, Command.StudioDiffHunkApply, {
    documentId: panel.wybor.idDokumentu,
    sourceVersionId: zrodlo,
    hunkIndex: numer,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  odswiez();
}

function opiszPasek(panel: Panel): void {
  const wersje = panel.wybor.wersje;
  if (wersje.length === 0) {
    zdejmijWyborWersji(panel.wezly);
    return;
  }
  panel.wezly.wersjaOdniesienia.textContent = wersje[panel.wybor.odniesienie]?.label ?? '';
  panel.wezly.wersjaPorownywana.textContent = wersje[panel.wybor.porownywana]?.label ?? '';
}

// Przełącznik trybu zostaje: wykaz zmian śledzonych stoi bez nazwanych wersji.
function zdejmijWyborWersji(wezly: WezlyRoznic): void {
  for (const wezel of [...wezly.wierszWersji.childNodes]) {
    if (wezel instanceof Element && wezel.classList.contains('dn-meta')) continue;
    wezel.remove();
  }
}

function opiszTryb(panel: Panel): void {
  if (panel.wezly.tryb === null) return;
  panel.wezly.tryb.textContent = `${NAZWY_TRYBU[panel.wybor.tryb] ?? ''} ▾`;
}

function idWersji(wybor: WyborPorownania, miejsce: number): string {
  return wybor.wersje[miejsce]?.id ?? '';
}

// Wybór krąży po wykazie, bo rozwijanego spisu wersji znacznik nie niesie.
function nastepneMiejsce(miejsce: number, ile: number): number {
  return ile === 0 ? 0 : (miejsce + 1) % ile;
}

function zbierzWezly(korzen: ParentNode): WezlyRoznic | null {
  const tresc = korzen.querySelector('#panel-diff .sta-okno-tresc');
  const wierszWersji = tresc?.querySelector(':scope > .st-panel-wiersz');
  const poleWzorca = tresc?.querySelector('.dn-szukaj');
  const wzorzec = poleWzorca?.querySelector('input[type="search"]');
  const przelaczniki = [...(wierszWersji?.querySelectorAll(':scope > button') ?? [])];
  const odniesienie = przelaczniki[0];
  const porownywana = przelaczniki[1];
  const tryb = wierszWersji?.querySelector('.dn-meta .dn-btn');
  if (
    !(tresc instanceof HTMLElement) ||
    !(wierszWersji instanceof HTMLElement) ||
    !(poleWzorca instanceof HTMLElement) ||
    !(wzorzec instanceof HTMLInputElement) ||
    !(odniesienie instanceof HTMLElement) ||
    !(porownywana instanceof HTMLElement)
  ) {
    return null;
  }
  return {
    tresc,
    wierszWersji,
    wersjaOdniesienia: odniesienie,
    wersjaPorownywana: porownywana,
    tryb: tryb instanceof HTMLElement ? tryb : null,
    poleWzorca,
    wzorzec,
  };
}

function zdejmijWzory(wezly: WezlyRoznic): WzoryWierszy {
  const wiersze = [...wezly.tresc.querySelectorAll(':scope > .st-panel-wiersz')];
  const decyzja = sklonuj(wiersze[1] ?? null);
  if (decyzja !== null) oczyscDecyzje(decyzja);
  return {
    usuniety: sklonuj(wezly.tresc.querySelector('.dn-diff-usu')?.closest('div') ?? null),
    dodany: oczyscDodany(sklonuj(wezly.tresc.querySelector('.dn-diff-dod')?.closest('div') ?? null)),
    formatowanie: sklonuj(wezly.tresc.querySelector('.dn-diff-fmt')?.closest('div') ?? null),
    decyzja,
    nota: sklonuj(wezly.tresc.querySelector('.dn-nota')),
  };
}

// Adnotacja schodzi: jej założenie wymaga treści, której znacznik nie przyjmie.
function oczyscDecyzje(wiersz: HTMLElement): void {
  wiersz.querySelector('[aria-label="Adnotacja"]')?.remove();
}

// Kontekst rdzeń oddaje osobnym fragmentem, nie doklejką do zdania wzorcowego.
function oczyscDodany(wzor: HTMLElement | null): HTMLElement | null {
  if (wzor === null) return null;
  for (const wezel of [...wzor.childNodes]) {
    if (wezel.nodeType === Node.TEXT_NODE) wezel.remove();
  }
  return wzor;
}

// Pasek wersji z przełącznikiem trybu i pole wzorca zostają; wiersze schodzą.
function zdejmijTrescPrzykladowa(wezly: WezlyRoznic): void {
  for (const dziecko of [...wezly.tresc.children]) {
    if (dziecko === wezly.wierszWersji || dziecko === wezly.poleWzorca) continue;
    dziecko.remove();
  }
}

function zdejmijPozycje(wezly: WezlyRoznic): void {
  let pozycja = wezly.poleWzorca.nextElementSibling;
  while (pozycja !== null) {
    const nastepna = pozycja.nextElementSibling;
    pozycja.remove();
    pozycja = nastepna;
  }
}

function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
