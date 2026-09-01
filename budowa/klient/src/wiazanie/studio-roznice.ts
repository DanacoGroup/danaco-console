/**
 * Wiązanie panelu różnic okna Studia z rdzeniem. Znacznik niesie biblioteka
 * Właściciela — ten plik nic nie buduje: woła porównanie wersji dokumentu,
 * wypełnia stojące węzły fragmentami odpowiedzi i powiela wzory wierszy zdjęte
 * z treści przykładowej.
 */

import {
  ChangeKind,
  Command,
  DiffHunkKind,
  EventType,
  type StudioDiffHunk,
  type StudioVersion,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Węzły panelu różnic, na których wiązanie pracuje. Brak któregokolwiek znaczy, że panel nie stoi w dokumencie. */
interface WezlyRoznic {
  tresc: HTMLElement;
  wierszWersji: HTMLElement;
  wersjaOdniesienia: HTMLElement;
  wersjaPorownywana: HTMLElement;
  poleWzorca: HTMLElement;
  wzorzec: HTMLInputElement;
}

/** Wzory wierszy zdjęte z treści przykładowej; klon zachowuje układ i klasy nadane przez bibliotekę. */
interface WzoryWierszy {
  usuniety: HTMLElement | null;
  dodany: HTMLElement | null;
  decyzja: HTMLElement | null;
}

/** Wybór paska: dokument porównywany, wykaz nazwanych wersji oraz miejsca obu przełączników w tym wykazie. */
interface WyborPorownania {
  idDokumentu: string;
  wersje: StudioVersion[];
  odniesienie: number;
  porownywana: number;
}

/** Panel po zebraniu węzłów: znacznik, wzory wierszy i wybór pary wersji, na których stoi porównanie. */
interface Panel {
  wezly: WezlyRoznic;
  wzory: WzoryWierszy;
  wybor: WyborPorownania;
}

/** Panel po ostatnim wiązaniu; po nim idzie porównanie wywołane spoza panelu — ze wstążki okna roboczego. */
interface Zwiazany {
  kanal: Kanal;
  idOkna: string;
  panel: Panel;
  odswiez: () => void;
}

/**
 * Wzory zdjęte przy pierwszym montażu okna. Powłoka wstawia wnętrze okna na
 * nowo przy każdym wejściu, a drugie zdjęcie zastałoby panel już opróżniony,
 * więc wzoru nie da się z niego wziąć po raz drugi.
 */
let wzoryPanelu: WzoryWierszy | null = null;

/** Panel wiązania bieżącego; pustka znaczy okno bez założonego stanowiska. */
let zwiazany: Zwiazany | null = null;

/**
 * Zdejmuje treść przykładową panelu różnic, zabierając z niej wzory wierszy.
 * Woła się przy montażu okna, przed powstaniem stanowiska: różnica z prototypu
 * opisuje cudzy dokument. Zwraca prawdę, gdy panel stał w dokumencie.
 */
export function zdejmijTrescPrzykladowaRoznic(): boolean {
  return przygotujPanel() !== null;
}

/** Zbiera węzły panelu i opróżnia je z treści przykładowej; pustka znaczy panel poza dokumentem. */
function przygotujPanel(): { wezly: WezlyRoznic; wzory: WzoryWierszy } | null {
  const wezly = zbierzWezly();
  if (wezly === null) return null;
  /* Znacznik inny niż związany znaczy panel wstawiony ponownie: wiązanie
     poprzednie trzyma węzły odczepione od dokumentu i porównanie wypełniłoby
     ekran, którego nie ma. */
  if (zwiazany !== null && zwiazany.panel.wezly.tresc !== wezly.tresc) zwiazany = null;
  wzoryPanelu ??= zdejmijWzory(wezly);
  zdejmijTrescPrzykladowa(wezly);
  return { wezly, wzory: wzoryPanelu };
}

/**
 * Wczytuje wersje dokumentu otwartego w oknie i zestawia parę wskazaną na
 * pasku panelu. Fałsz znaczy panel bez wiązania albo okno bez dokumentu —
 * wołający ma wtedy czym odmówić Operatorowi zamiast milczeć.
 */
export async function porownajWersjeDokumentu(): Promise<boolean> {
  const biezacy = zwiazany;
  if (biezacy === null) return false;
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

/**
 * Wiąże panel różnic okna z rdzeniem. Dokument bierze się ze zdarzenia zmiany
 * dokumentu wskazanego okna, bo komendy pytającej o dokument otwarty w oknie
 * kontrakt nie ma. Zwraca prawdę, gdy znacznik panelu stał.
 */
export function zwiazRoznice(kanal: Kanal, idOkna: string): boolean {
  const przygotowany = przygotujPanel();
  if (przygotowany === null) return false;
  const wezly: WezlyRoznic = przygotowany.wezly;

  const panel: Panel = {
    wezly,
    wzory: przygotowany.wzory,
    wybor: { idDokumentu: '', wersje: [], odniesienie: 0, porownywana: 0 },
  };

  /** Powtarza porównanie wybranej pary; bez znanego dokumentu nie ma o co pytać. */
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

  kanal.naZdarzenie(EventType.StudioDocumentChanged, (tresc) => {
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

  zwiazany = { kanal, idOkna, panel, odswiez };
  return true;
}

/** Wczytuje nazwane wersje dokumentu z repozytorium sesji i porównuje parę wskazaną na pasku. */
async function wczytaj(kanal: Kanal, panel: Panel, odswiez: () => void): Promise<void> {
  const wynik = await wywolaj(kanal, Command.StudioRepositoryList, {
    documentId: panel.wybor.idDokumentu,
  });
  if (wynik.udany && wynik.wynik !== undefined) {
    /* Nazwa wersji jest w kontrakcie nieobowiązkowa, a przełącznik nie ma czym
       nazwać wersji bez nazwy, więc taka wersja do wykazu paska nie wchodzi. */
    panel.wybor.wersje = wynik.wynik.versions.filter((wersja) => (wersja.label ?? '') !== '');
    panel.wybor.porownywana = 0;
    panel.wybor.odniesienie = Math.min(1, panel.wybor.wersje.length - 1);
  }
  opiszPasek(panel);
  await porownaj(kanal, panel, odswiez);
}

/** Woła porównanie pary wersji i wypełnia panel fragmentami odpowiedzi; odmowa rdzenia zostawia panel pusty. */
async function porownaj(kanal: Kanal, panel: Panel, odswiez: () => void): Promise<void> {
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

/** Wstawia fragmenty odpowiedzi rdzenia, powielając wzory wierszy; za każdym fragmentem staje jego wiersz decyzji. */
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

/**
 * Klony wierszy niosących treść fragmentu. Fragment zmieniony pokazuje obie
 * swoje treści we wzorach ubytku i przyrostu, a fragment podany dla kontekstu
 * odpada, bo wzoru dla niego znacznik nie niesie.
 */
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

/** Klon wzoru wiersza z treścią fragmentu; brak treści w odpowiedzi zdejmuje cały wiersz, bo nie ma czego pokazać. */
function wypelnijWiersz(wzor: HTMLElement | null, selektor: string, tresc?: string): HTMLElement | null {
  if (wzor === null || tresc === undefined) return null;
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const miejsce = wiersz.querySelector(selektor);
  if (miejsce === null) return null;
  miejsce.textContent = tresc;
  return wiersz;
}

/** Klon wiersza decyzji związany z numerem fragmentu; bez wybranej wersji porównywanej wiersz odpada, bo komenda przeniesienia jej wymaga. */
function wierszDecyzji(
  kanal: Kanal,
  panel: Panel,
  numer: number,
  odswiez: () => void,
): HTMLElement | null {
  if (panel.wzory.decyzja === null) return null;
  if (idWersji(panel.wybor, panel.wybor.porownywana) === '') return null;
  const wiersz = panel.wzory.decyzja.cloneNode(true) as HTMLElement;
  const przyjmij = wiersz.querySelector('button');
  if (przyjmij === null) return null;
  przyjmij.addEventListener('click', () => {
    void przyjmijFragment(kanal, panel, numer, odswiez);
  });
  return wiersz;
}

/** Przenosi fragment z wersji porównywanej do stanu bieżącego i powtarza porównanie na dokumencie po zmianie. */
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

/** Nadaje przełącznikom nazwy wybranych wersji; brak nazwanych wersji zdejmuje pasek, bo nie ma czego na nim wskazać. */
function opiszPasek(panel: Panel): void {
  const wersje = panel.wybor.wersje;
  if (wersje.length === 0) {
    panel.wezly.wierszWersji.remove();
    return;
  }
  panel.wezly.wersjaOdniesienia.textContent = wersje[panel.wybor.odniesienie]?.label ?? '';
  panel.wezly.wersjaPorownywana.textContent = wersje[panel.wybor.porownywana]?.label ?? '';
}

/** Identyfikator wersji spod wskazanego miejsca wykazu; pustka znaczy wykaz nazwanych wersji bez tego miejsca. */
function idWersji(wybor: WyborPorownania, miejsce: number): string {
  return wybor.wersje[miejsce]?.id ?? '';
}

/** Kolejne miejsce w wykazie wersji; wybór krąży po wykazie, bo rozwijanego spisu wersji znacznik nie niesie. */
function nastepneMiejsce(miejsce: number, ile: number): number {
  return ile === 0 ? 0 : (miejsce + 1) % ile;
}

/** Wskazuje węzły panelu różnic; pustka znaczy, że panel nie stoi w dokumencie. */
function zbierzWezly(): WezlyRoznic | null {
  const tresc = document.querySelector('#panel-diff .sta-okno-tresc');
  const wierszWersji = tresc?.querySelector(':scope > .st-panel-wiersz');
  const poleWzorca = tresc?.querySelector('.dn-szukaj');
  const wzorzec = poleWzorca?.querySelector('input[type="search"]');
  const przelaczniki = [...(wierszWersji?.querySelectorAll(':scope > button') ?? [])];
  const odniesienie = przelaczniki[0];
  const porownywana = przelaczniki[1];
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
    poleWzorca,
    wzorzec,
  };
}

/** Zdejmuje wzory wierszy, zanim treść przykładowa zniknie; wzór decyzji traci przyciski bez pokrycia w rdzeniu. */
function zdejmijWzory(wezly: WezlyRoznic): WzoryWierszy {
  const wiersze = [...wezly.tresc.querySelectorAll(':scope > .st-panel-wiersz')];
  const decyzja = sklonuj(wiersze[1] ?? null);
  if (decyzja !== null) oczyscDecyzje(decyzja);
  return {
    usuniety: sklonuj(wezly.tresc.querySelector('.dn-diff-usu')?.closest('div') ?? null),
    dodany: oczyscDodany(sklonuj(wezly.tresc.querySelector('.dn-diff-dod')?.closest('div') ?? null)),
    decyzja,
  };
}

/**
 * Zdejmuje ze wzoru decyzji przyciski bez pokrycia: odrzucenie fragmentu ma
 * komendę wyłącznie przy porównaniu z propozycją zmiany, a założenie adnotacji
 * wymaga treści, której znacznik nie ma gdzie przyjąć.
 */
function oczyscDecyzje(wiersz: HTMLElement): void {
  wiersz.querySelector('.dn-btn--duch')?.remove();
  wiersz.querySelector('[aria-label="Adnotacja"]')?.remove();
}

/** Zdejmuje ze wzoru przyrostu goły tekst dopisany do zdania wzorcowego: kontekst rdzeń oddaje osobnym fragmentem, nie doklejką. */
function oczyscDodany(wzor: HTMLElement | null): HTMLElement | null {
  if (wzor === null) return null;
  for (const wezel of [...wzor.childNodes]) {
    if (wezel.nodeType === Node.TEXT_NODE) wezel.remove();
  }
  return wzor;
}

/**
 * Zdejmuje treść przykładową panelu. Przełącznik stanu scalenia odpada, bo pola
 * trybu widoku ani stanu scalenia pary wersji nie ma żadna komenda różnicy;
 * wzór różnicy postaci i statystyka słów odpadają razem z wierszami, bo rdzeń
 * ani rodzaju postaci, ani licznika słów w różnicy nie oddaje.
 */
function zdejmijTrescPrzykladowa(wezly: WezlyRoznic): void {
  wezly.wierszWersji.querySelector('.dn-meta .dn-btn')?.remove();
  for (const dziecko of [...wezly.tresc.children]) {
    if (dziecko === wezly.wierszWersji || dziecko === wezly.poleWzorca) continue;
    dziecko.remove();
  }
}

/** Zdejmuje pozycje poprzedniego porównania; pasek wersji i pole wzorca zostają, bo należą do znacznika, nie do odpowiedzi. */
function zdejmijPozycje(wezly: WezlyRoznic): void {
  let pozycja = wezly.poleWzorca.nextElementSibling;
  while (pozycja !== null) {
    const nastepna = pozycja.nextElementSibling;
    pozycja.remove();
    pozycja = nastepna;
  }
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
