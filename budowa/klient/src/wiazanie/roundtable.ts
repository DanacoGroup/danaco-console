// Wiązanie karty modułu Roundtable z rdzeniem: okno debaty, wykaz modeli i argumentów. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type RoundtableArgumentNode,
  type RoundtableParticipant,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  opiszNaglowek,
  zapewnijOknoModulu,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

const KOD_MODULU = 'roundtable';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazDebate(
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
  zwolnijDebate(idKarty);

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
    await Promise.all([wypelnijModele(kanal, korzen, idOkna), wypelnijArgumenty(kanal, korzen, idOkna)]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Roundtable', idOkna);
    if (idOkna === '') {
      niegotowe(panel(korzen, 'panel-debate'), 'Rdzeń nie dał okna debaty dla tej karty.');
      return;
    }
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-debate .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void odswiez();
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Debata');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijDebate(idKarty: string): void {
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
  niegotowe(panel(korzen, 'panel-debate'), 'Wykaz czeka na odpowiedź rdzenia.');
  /* Siatka mówców niesie w prototypie całą debatę wpisaną wprost — wypowiedzi,
     miary czasu i znaków. Uczestników poda rdzeń, więc panele schodzą puste. */
  for (const mowca of korzen.querySelectorAll('#model-panel-a, #model-panel-b')) mowca.remove();
  niegotowe(panel(korzen, 'panel-moderator'), 'Wykaz modeli czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-consensus'), 'Zgody nie ma czym policzyć — debata jeszcze nie ruszyła.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan debaty czeka na pierwszą turę.');
  niegotowe(panel(korzen, 'panel-artefakty'), 'Rdzeń nie podaje wytworów tej debaty.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań debaty dla tego okna.');
}

async function wypelnijModele(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-moderator');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.RoundtableModelList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz modeli nie doszedł.');
    return;
  }
  const modele: RoundtableParticipant[] = odpowiedz.wynik.participants;
  if (modele.length === 0) {
    niegotowe(cialo, 'Żaden model nie stanął jeszcze do tej debaty.');
    return;
  }
  cialo.replaceChildren(...modele.map((m) => pozycja(cialo, m.personaName ?? m.id, m.channelId)));
}

async function wypelnijArgumenty(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-debate');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.RoundtableArgumentList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz argumentów nie doszedł.');
    return;
  }
  const argumenty: RoundtableArgumentNode[] = odpowiedz.wynik.graph.nodes;
  if (argumenty.length === 0) {
    niegotowe(cialo, 'Debata nie ma jeszcze ani jednej wypowiedzi.');
    return;
  }
  cialo.replaceChildren(...argumenty.map((a) => pozycja(cialo, a.text, a.speechAct)));
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
