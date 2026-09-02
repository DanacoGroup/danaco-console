// Wiązanie karty Terminala z rdzeniem: karty powłok, wyjście procesów, monitor
// procesów, ściąga poleceń i wykaz hostów.

import {
  Command,
  EventType,
  type TerminalHost,
  type TerminalProcess,
  type TerminalScript,
  type TerminalSession,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

interface WiazanieTerminala {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieTerminala>();

export function zwiazTerminal(
  kanal: Kanal,
  idOkna: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijTerminal(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  zdejmijTrescPrzykladowa(korzen);

  void wypelnijKarty(kanal, korzen, idOkna);
  void wypelnijProcesy(kanal, korzen, idOkna);
  void wypelnijSciage(kanal, korzen);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    const wstawka = cel.closest<HTMLElement>('[data-sciaga]');
    if (wstawka !== null) {
      wpiszPolecenie(korzen, wstawka.dataset.sciaga ?? '');
      return;
    }
    const zdjecie = cel.closest<HTMLElement>('[data-proces-zdejmij]');
    if (zdjecie !== null) {
      void zdejmijProces(kanal, korzen, idOkna, zdjecie.dataset.procesZdejmij ?? '');
      return;
    }
    const wstrzymanie = cel.closest<HTMLElement>('[data-proces-wstrzymaj]');
    if (wstrzymanie !== null) {
      void wstrzymajProces(kanal, korzen, idOkna, wstrzymanie.dataset.procesWstrzymaj ?? '');
      return;
    }
    const zamkniecie = cel.closest<HTMLElement>('[data-karta-zamknij]');
    if (zamkniecie !== null) {
      void zamknijKarte(kanal, korzen, idOkna, zamkniecie.dataset.kartaZamknij ?? '');
    }
  }, przy);

  const pole = korzen.querySelector<HTMLInputElement>('[data-polecenie]');
  pole?.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    void wykonajPolecenie(kanal, korzen, idOkna, pole.value);
  }, przy);

  odlaczenia.push(zglosUchwyt(EventType.TerminalProcessChanged, () => {
    void wypelnijProcesy(kanal, korzen, idOkna);
  }));

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, 'terminal', 'Terminal');
  if (katalog !== null) odlaczenia.push(katalog);

  return true;
}

export function zwolnijTerminal(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

/* Prototyp niesie karty, wiersze wyjścia, procesy i ściągę wpisane w znacznik.
   Wypełniacz zostaje w produkcie na zawsze, więc schodzi przy montażu. */
function zdejmijTrescPrzykladowa(korzen: Element): void {
  for (const wybor of ['.oc-wiersz', '.pm-wiersz', '[data-sciaga]', '.tt-karta']) {
    for (const wezel of korzen.querySelectorAll(wybor)) wezel.remove();
  }
}

async function wypelnijKarty(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gniazdo = korzen.querySelector('#panel-tabs .sta-okno-tresc');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalSessionList, { windowId: idOkna });
  if (!wynik.udany) {
    postawStanPusty(gniazdo, 'Wykaz kart terminala nie doszedł.');
    return;
  }
  const karty = (wynik.wynik as { sessions?: TerminalSession[] } | undefined)?.sessions ?? [];
  gniazdo.replaceChildren();
  if (karty.length === 0) {
    postawStanPusty(gniazdo, 'Żadna karta terminala nie stoi.');
    return;
  }
  korzen.setAttribute('data-karta-terminala', karty[0]?.id ?? '');
  for (const karta of karty) {
    const wiersz = document.createElement('div');
    wiersz.className = 'tt-karta';
    wiersz.dataset.karta = karta.id;
    const tytul = document.createElement('span');
    tytul.textContent = karta.title ?? karta.shell;
    const stan = document.createElement('span');
    stan.className = 'dn-plakietka';
    stan.textContent = karta.status;
    const zamknij = document.createElement('button');
    zamknij.className = 'dn-btn dn-btn--zarys dn-btn--sm';
    zamknij.type = 'button';
    zamknij.dataset.kartaZamknij = karta.id;
    zamknij.textContent = 'Zamknij';
    wiersz.append(tytul, stan, zamknij);
    gniazdo.append(wiersz);
  }
}

async function wypelnijProcesy(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gniazdo = korzen.querySelector('#panel-process .sta-okno-tresc');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalProcessList, { windowId: idOkna });
  if (!wynik.udany) {
    postawStanPusty(gniazdo, 'Wykaz procesów nie doszedł.');
    return;
  }
  const procesy = (wynik.wynik as { processes?: TerminalProcess[] } | undefined)?.processes ?? [];
  gniazdo.replaceChildren();
  if (procesy.length === 0) {
    postawStanPusty(gniazdo, 'Żaden proces nie biegnie.');
    return;
  }
  for (const proces of procesy) {
    const wiersz = document.createElement('div');
    wiersz.className = 'pm-wiersz';
    const glowny = document.createElement('div');
    glowny.className = 'pm-glowny';
    const polecenie = document.createElement('span');
    polecenie.className = 'pm-cmd';
    polecenie.textContent = proces.command;
    const zrodlo = document.createElement('span');
    zrodlo.className = 'dn-plakietka dn-na-koniec';
    zrodlo.textContent = proces.initiator;
    const wstrzymaj = document.createElement('button');
    wstrzymaj.className = 'dn-btn dn-btn--zarys dn-btn--sm';
    wstrzymaj.type = 'button';
    wstrzymaj.dataset.procesWstrzymaj = proces.id;
    wstrzymaj.textContent = 'Wstrzymaj';
    const zdejmij = document.createElement('button');
    zdejmij.className = 'dn-btn dn-btn--zarys dn-btn--sm';
    zdejmij.type = 'button';
    zdejmij.dataset.procesZdejmij = proces.id;
    zdejmij.textContent = 'Zakończ';
    glowny.append(polecenie, zrodlo, wstrzymaj, zdejmij);
    wiersz.append(glowny);
    gniazdo.append(wiersz);
  }
}

async function wypelnijSciage(kanal: Kanal, korzen: Element): Promise<void> {
  const gniazdo = korzen.querySelector('#menu-sciaga');
  if (gniazdo === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalScriptList, {});
  if (!wynik.udany) return;
  const skrypty = (wynik.wynik as { scripts?: TerminalScript[] } | undefined)?.scripts ?? [];
  for (const etykieta of gniazdo.querySelectorAll('.sta-menu-etyk')) {
    etykieta.textContent = `Ściąga poleceń (${skrypty.length})`;
  }
  for (const skrypt of skrypty) {
    const pozycja = document.createElement('button');
    pozycja.className = 'sta-menu-poz';
    pozycja.type = 'button';
    pozycja.setAttribute('role', 'menuitem');
    pozycja.dataset.sciaga = skrypt.content;
    pozycja.textContent = skrypt.alias ?? skrypt.name;
    gniazdo.append(pozycja);
  }
}

async function wykonajPolecenie(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  tresc: string,
): Promise<void> {
  const polecenie = tresc.trim();
  if (polecenie === '') return;
  const idKarty = korzen.getAttribute('data-karta-terminala') ?? '';
  if (idKarty === '') {
    oglos('Terminal', 'Żadna karta terminala nie stoi — polecenie nie ma gdzie ruszyć.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalCommandExec, {
    sessionId: idKarty,
    command: polecenie,
  });
  if (!wynik.udany) {
    oglos('Terminal', wynik.blad?.message ?? 'Polecenie nie ruszyło.', 'blad');
    return;
  }
  const pole = korzen.querySelector<HTMLInputElement>('[data-polecenie]');
  if (pole !== null) pole.value = '';
  void wypelnijProcesy(kanal, korzen, idOkna);
}

async function zdejmijProces(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  idProcesu: string,
): Promise<void> {
  if (idProcesu === '') return;
  const wynik = await wywolaj(kanal, Command.TerminalProcessKill, { processId: idProcesu });
  if (!wynik.udany) {
    oglos('Terminal', wynik.blad?.message ?? 'Proces nie został zakończony.', 'blad');
    return;
  }
  void wypelnijProcesy(kanal, korzen, idOkna);
}

async function wstrzymajProces(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  idProcesu: string,
): Promise<void> {
  if (idProcesu === '') return;
  const wynik = await wywolaj(kanal, Command.TerminalProcessSuspend, { processId: idProcesu });
  if (!wynik.udany) {
    oglos('Terminal', wynik.blad?.message ?? 'Proces nie został wstrzymany.', 'blad');
    return;
  }
  void wypelnijProcesy(kanal, korzen, idOkna);
}

async function zamknijKarte(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  idKarty: string,
): Promise<void> {
  if (idKarty === '') return;
  const wynik = await wywolaj(kanal, Command.TerminalSessionClose, { sessionId: idKarty });
  if (!wynik.udany) {
    oglos('Terminal', wynik.blad?.message ?? 'Karta nie została zamknięta.', 'blad');
    return;
  }
  void wypelnijKarty(kanal, korzen, idOkna);
}

function wpiszPolecenie(korzen: Element, tresc: string): void {
  const pole = korzen.querySelector<HTMLInputElement>('[data-polecenie]');
  if (pole === null || tresc === '') return;
  pole.value = tresc;
  pole.focus();
}

function postawStanPusty(gniazdo: Element, zdanie: string): void {
  gniazdo.replaceChildren();
  const wezel = document.createElement('p');
  wezel.className = 'dn-pusty';
  wezel.textContent = zdanie;
  gniazdo.append(wezel);
}

export type { TerminalHost };
