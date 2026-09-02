// Tools Panel okna Studia: wykaz operacji kontekstowych, zakres działania
// i wywołanie operacji na dokumencie. Znacznik niesie biblioteka Właściciela.

import {
  ChangeKind,
  Command,
  EventType,
  StudioOperationScope,
  type StudioOperation,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const KLASA_WCISNIETY = 'dn-btn--zarys';
const KLASA_SPOCZYNKU = 'dn-btn--duch';

interface WezlyPanelu {
  lista: HTMLElement;
  zakladki: HTMLElement | null;
  zakresy: HTMLElement[];
  uruchom: HTMLElement | null;
}

interface WzoryPanelu {
  grupa: HTMLElement | null;
  wiersz: HTMLElement | null;
}

// Wzory zdejmuje pierwsze wiązanie: lista opróżniona wzoru już nie oddaje.
let wzory: WzoryPanelu | null = null;

// Karty Studia stoją obok siebie: przycisk wstążki trafia w panel własnej karty.
const ZWIAZANE = new Map<string, () => Promise<void>>();

// Woła się przy montażu karty: wykaz operacji z prototypu opisuje cudzy dokument.
export function zdejmijTrescPrzykladowaNarzedzi(korzen: ParentNode): boolean {
  return przygotujPanel(korzen) !== null;
}

function przygotujPanel(korzen: ParentNode): WezlyPanelu | null {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return null;
  wzory ??= zdejmijWzory(wezly.lista);
  zdejmijPozycje(wezly);
  return wezly;
}

// Fałsz znaczy kartę bez wiązania panelu — wołający ma wtedy czym odmówić.
export async function uruchomOperacjeDokumentu(idKarty: string): Promise<boolean> {
  const uruchom = ZWIAZANE.get(idKarty);
  if (uruchom === undefined) return false;
  await uruchom();
  return true;
}

// Dokument bierze się ze zdarzenia zmiany albo z otwarcia: komendy pytającej
// o dokument otwarty w oknie kontrakt nie ma.
export function zwiazNarzedzia(
  kanal: Kanal,
  idOkna: string,
  idKarty: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const wezly = przygotujPanel(korzen);
  if (wezly === null) return null;
  const kanwa = korzen.querySelector('.dn-kanwa');
  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });

  let wybrana = '';
  let zakres: StudioOperationScope = StudioOperationScope.Selection;
  let idDokumentu = '';

  const wskaz = (identyfikator: string): void => {
    wybrana = identyfikator;
    oznaczWybor(wezly.lista, wybrana);
  };

  const uruchom = async (): Promise<void> => {
    const dokument = idDokumentu === '' ? await otworzDokument(kanal, idOkna) : idDokumentu;
    idDokumentu = dokument;
    await wywolajOperacje(kanal, { idOkna, dokument, wybrana, zakres }, kanwa);
  };

  wezly.lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('.st-panel-wiersz[data-operacja]');
    if (wiersz !== null) wskaz(wiersz.dataset.operacja ?? '');
  }, przy);

  for (const [numer, przycisk] of wezly.zakresy.entries()) {
    przycisk.addEventListener('click', () => {
      zakres = numer === 0 ? StudioOperationScope.Selection : StudioOperationScope.Document;
      oznaczZakres(wezly.zakresy, numer);
    }, przy);
  }

  wezly.uruchom?.addEventListener('click', () => {
    void uruchom();
  }, przy);

  odlaczenia.push(
    zglosUchwyt(EventType.StudioDocumentChanged, (tresc) => {
      if (tresc.document.windowId !== idOkna) return;
      idDokumentu = tresc.change === ChangeKind.Deleted ? '' : tresc.document.id;
    }),
  );

  void odczytajOperacje(kanal).then((wykaz) => {
    wypelnij(wezly, wykaz);
    oznaczWybor(wezly.lista, wybrana);
  });
  oznaczZakres(wezly.zakresy, 0);
  ZWIAZANE.set(idKarty, uruchom);

  return () => {
    for (const odlacz of odlaczenia) odlacz();
    if (ZWIAZANE.get(idKarty) === uruchom) ZWIAZANE.delete(idKarty);
  };
}

interface Wywolanie {
  idOkna: string;
  dokument: string;
  wybrana: string;
  zakres: StudioOperationScope;
}

async function wywolajOperacje(
  kanal: Kanal,
  wywolanie: Wywolanie,
  kanwa: Element | null,
): Promise<void> {
  if (wywolanie.wybrana === '') {
    oglos('Studio', 'Wskaż operację na liście Tools Panel.');
    return;
  }
  if (wywolanie.dokument === '') {
    oglos('Studio', 'Okno nie prowadzi dokumentu — nie ma na czym wykonać operacji.');
    return;
  }
  const zaznaczenie = wywolanie.zakres === StudioOperationScope.Selection
    ? zaznaczenieKanwy(kanwa)
    : null;
  if (wywolanie.zakres === StudioOperationScope.Selection && zaznaczenie === null) {
    oglos('Studio', 'Zaznacz fragment dokumentu albo przestaw zakres na cały dokument.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.StudioContextualOp, {
    windowId: wywolanie.idOkna,
    documentId: wywolanie.dokument,
    actionId: wywolanie.wybrana,
    scope: wywolanie.zakres,
    ...(zaznaczenie === null
      ? {}
      : { selectionStart: zaznaczenie.poczatek, selectionEnd: zaznaczenie.koniec }),
  });
  if (!wynik.udany) {
    oglos('Studio', wynik.blad?.message ?? 'Rdzeń odmówił wykonania operacji.', 'ostrzezenie');
  }
}

function zaznaczenieKanwy(kanwa: Element | null): { poczatek: number; koniec: number } | null {
  if (kanwa === null) return null;
  const zaznaczenie = globalThis.getSelection();
  if (zaznaczenie === null || zaznaczenie.rangeCount === 0) return null;
  const zakres = zaznaczenie.getRangeAt(0);
  if (zakres.collapsed || !kanwa.contains(zakres.commonAncestorContainer)) return null;
  const przed = zakres.cloneRange();
  przed.selectNodeContents(kanwa);
  przed.setEnd(zakres.startContainer, zakres.startOffset);
  const poczatek = przed.toString().length;
  return { poczatek, koniec: poczatek + zakres.toString().length };
}

async function odczytajOperacje(kanal: Kanal): Promise<StudioOperation[]> {
  const wynik = await wywolaj(kanal, Command.StudioOperationList, {});
  if (!wynik.udany || wynik.wynik === undefined) return [];
  return wynik.wynik.operations;
}

async function otworzDokument(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  return wynik.wynik.document.id;
}

function wypelnij(wezly: WezlyPanelu, operacje: StudioOperation[]): void {
  zdejmijPozycje(wezly);
  let kategoria: string | null = null;
  for (const operacja of operacje) {
    if (operacja.category !== kategoria) {
      kategoria = operacja.category;
      wstaw(wezly, naglowek(kategoria));
    }
    wstaw(wezly, wiersz(operacja));
  }
}

function wstaw(wezly: WezlyPanelu, wezel: HTMLElement | null): void {
  if (wezel === null) return;
  wezly.lista.insertBefore(wezel, wezly.uruchom);
}

function naglowek(kategoria: string): HTMLElement | null {
  const wezel = sklonuj(wzory?.grupa ?? null);
  if (wezel === null || kategoria === '') return null;
  wezel.textContent = kategoria;
  return wezel;
}

function wiersz(operacja: StudioOperation): HTMLElement | null {
  const wezel = sklonuj(wzory?.wiersz ?? null);
  if (wezel === null) return null;
  wezel.dataset.operacja = operacja.id;
  const znak = wezel.querySelector('.dn-meta');
  wezel.replaceChildren(document.createTextNode(operacja.name));
  if (znak !== null) {
    znak.textContent = '';
    wezel.appendChild(znak);
  }
  return wezel;
}

// Znak wyboru stoi w węźle, który w prototypie niósł znak rozwinięcia.
function oznaczWybor(lista: HTMLElement, wybrana: string): void {
  for (const wezel of lista.querySelectorAll<HTMLElement>('.st-panel-wiersz[data-operacja]')) {
    const czynny = wezel.dataset.operacja === wybrana && wybrana !== '';
    wezel.setAttribute('aria-current', String(czynny));
    const znak = wezel.querySelector('.dn-meta');
    if (znak !== null) znak.textContent = czynny ? '✓' : '';
  }
}

function oznaczZakres(zakresy: HTMLElement[], czynny: number): void {
  for (const [numer, przycisk] of zakresy.entries()) {
    const wcisniety = numer === czynny;
    przycisk.setAttribute('aria-pressed', String(wcisniety));
    przycisk.classList.toggle(KLASA_WCISNIETY, wcisniety);
    przycisk.classList.toggle(KLASA_SPOCZYNKU, !wcisniety);
  }
}

// Nagłówkiem grupy jest drugi napis listy: pierwszy niesie miarę zaznaczenia,
// której kontrakt nie oddaje.
function zdejmijWzory(lista: HTMLElement): WzoryPanelu {
  const napisy = lista.querySelectorAll('.pt-etykieta');
  return {
    grupa: sklonuj(napisy[1] ?? null),
    wiersz: sklonuj(lista.querySelector('.st-panel-wiersz')),
  };
}

function zdejmijPozycje(wezly: WezlyPanelu): void {
  for (const wezel of wezly.lista.querySelectorAll('.pt-etykieta, .st-panel-wiersz')) {
    wezel.remove();
  }
}

function zbierzWezly(korzen: ParentNode): WezlyPanelu | null {
  const lista = korzen.querySelector('#panel-tools .sta-okno-tresc.st-panel-lista');
  if (!(lista instanceof HTMLElement)) return null;
  const zakladki = lista.querySelector('.dn-zakladki');
  const uruchom = lista.querySelector('.st-odsun-sekcja');
  return {
    lista,
    zakladki: zakladki instanceof HTMLElement ? zakladki : null,
    zakresy: [...(zakladki?.querySelectorAll<HTMLElement>('button') ?? [])],
    uruchom: uruchom instanceof HTMLElement ? uruchom : null,
  };
}

function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
