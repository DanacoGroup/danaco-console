// Wiązanie karty modułu Apps z rdzeniem: wykaz wdrożeń aplikacji. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type AppDeployment,
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
    await Promise.all([wypelnijWdrozenia(kanal, korzen, idOkna)]);
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
