// Wiązanie karty modułu Assistant z rdzeniem: dziennik działań asystenta. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type AssistantAction,
  type AssistantActivityEntry,
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
import { zwiazGlosAsystenta } from './assistant-glos.ts';
import { zwiazOknoAsystenta } from './assistant-okno.ts';
import { zwiazPoczteAsystenta } from './assistant-poczta.ts';
import { zwiazSchowekAsystenta } from './assistant-schowek.ts';
import { oglos } from './ogloszenie.ts';

const KOD_MODULU = 'assistant';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazAsystenta(
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
  zwolnijAsystenta(idKarty);

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
    await Promise.all([      wypelnijDziennik(kanal, korzen, idOkna),
      wypelnijZlecenia(kanal, korzen, idOkna)]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Assistant', idOkna);
    if (idOkna === '') return;
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('[data-wpis-dziennika]');
    if (wiersz !== null && cel.closest('[data-wpis-wyroznij]') !== null) {
      zdarzenie.stopPropagation();
      void przelaczWyroznienie(kanal, korzen, wiersz, idOkna);
      return;
    }
    if (cel.closest('[data-polecenie-wyslij]') !== null) {
      zdarzenie.stopPropagation();
      void wyslijPolecenie(kanal, korzen, idOkna);
      return;
    }
    if (cel.closest('#panel-activity .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void odswiez();
  }, przy);

  zwiazGlosAsystenta(kanal, korzen, () => idOkna, (zdejmij) => {
    odlaczenia.push(zdejmij);
  }, przy);

  zwiazOknoAsystenta(kanal, korzen, () => idOkna, (zdejmij) => {
    odlaczenia.push(zdejmij);
  }, przy);

  zwiazSchowekAsystenta(kanal, korzen, () => idOkna, (zdejmij) => {
    odlaczenia.push(zdejmij);
  }, przy);

  zwiazPoczteAsystenta(kanal, korzen, (zdejmij) => {
    odlaczenia.push(zdejmij);
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Assistant');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijAsystenta(idKarty: string): void {
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
  niegotowe(panel(korzen, 'panel-activity'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-actions'), 'Rdzeń nie podaje zleceń dla tego okna.');
  postawPolecenie(panel(korzen, 'panel-voice'));
  niegotowe(panel(korzen, 'panel-plan'), 'Plan pracy czeka na pierwsze zlecenie.');
  niegotowe(panel(korzen, 'panel-artefakty'), 'Rdzeń nie podaje wytworów tego okna.');
  niegotowe(panel(korzen, 'panel-pliki'), 'Rdzeń nie podaje plików tego okna.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań dla tego okna.');
}

async function wypelnijDziennik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-activity');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.AssistantActivityList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Dziennik działań nie doszedł.');
    return;
  }
  const wpisy: AssistantActivityEntry[] = odpowiedz.wynik.entries;
  if (wpisy.length === 0) {
    niegotowe(cialo, 'Asystent nie odnotował jeszcze żadnego działania.');
    return;
  }
  cialo.replaceChildren(...wpisy.map((w) => wpisDziennika(cialo, w)));
}

/* Wpis dziennika niesie własne wyróżnienie: `assistant.activity.flag` odwraca
   oznaczenie, a wykaz wczytuje się na nowo, bo rdzeń rozstrzyga o kolejności. */
function wpisDziennika(cialo: Element, wpis: AssistantActivityEntry): HTMLElement {
  const wiersz = pozycja(cialo, wpis.content, wpis.kind);
  wiersz.dataset.wpisDziennika = wpis.id;
  const przelacz = cialo.ownerDocument.createElement('button');
  przelacz.type = 'button';
  przelacz.className = 'dn-btn-ikona';
  przelacz.dataset.wpisWyroznij = wpis.important === true ? 'zdejmij' : 'nadaj';
  przelacz.setAttribute(
    'aria-label',
    wpis.important === true ? 'Zdejmij wyróżnienie wpisu' : 'Wyróżnij wpis dziennika',
  );
  przelacz.textContent = wpis.important === true ? '★' : '☆';
  wiersz.appendChild(przelacz);
  return wiersz;
}

/* Nagrywania okno nie prowadzi, ale komenda przyjmuje samą transkrypcję, więc
   polecenie da się podyktować pisemnie: pole i przycisk zamiast zdania o braku. */
function postawPolecenie(cialo: Element | null): void {
  if (cialo === null) return;
  const pole = cialo.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Treść polecenia';
  pole.setAttribute('aria-label', 'Treść polecenia dla asystenta');
  pole.dataset.polecenieTresc = '';
  const przycisk = cialo.ownerDocument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przycisk.dataset.polecenieWyslij = '';
  przycisk.textContent = 'Wyślij polecenie';
  const zdanie = cialo.ownerDocument.createElement('div');
  zdanie.className = 'dn-meta';
  zdanie.textContent = 'Nagrywanie głosu nie wchodzi do tego wydania; polecenie idzie pisemnie.';
  cialo.replaceChildren(pole, przycisk, zdanie);
}

async function wyslijPolecenie(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
): Promise<void> {
  const pole = korzen.querySelector<HTMLInputElement>('[data-polecenie-tresc]');
  const tresc = (pole?.value ?? '').trim();
  if (idOkna === '' || tresc === '') return;
  const wynik = await wywolaj(kanal, Command.AssistantVoiceCommand, {
    windowId: idOkna,
    transcript: tresc,
  });
  if (!wynik.udany) {
    oglos('Asystent', wynik.blad?.message ?? 'Rdzeń odmówił przyjęcia polecenia.', 'ostrzezenie');
    return;
  }
  if (pole !== null) pole.value = '';
  oglos('Asystent', `Zlecenie założone: ${wynik.wynik?.action.title ?? 'bez nazwy'}.`);
  await wypelnijZlecenia(kanal, korzen, idOkna);
}

async function przelaczWyroznienie(
  kanal: Kanal,
  korzen: Element,
  wiersz: HTMLElement,
  idOkna: string,
): Promise<void> {
  const idWpisu = wiersz.dataset.wpisDziennika ?? '';
  const przycisk = wiersz.querySelector<HTMLElement>('[data-wpis-wyroznij]');
  if (idWpisu === '' || przycisk === null) return;
  const wynik = await wywolaj(kanal, Command.AssistantActivityFlag, {
    entryId: idWpisu,
    important: przycisk.dataset.wpisWyroznij === 'nadaj',
  });
  if (!wynik.udany) {
    oglos('Dziennik asystenta',
      wynik.blad?.message ?? 'Rdzeń odmówił zmiany wyróżnienia wpisu.', 'ostrzezenie');
    return;
  }
  await wypelnijDziennik(kanal, korzen, idOkna);
}

/* Zlecenia asystenta: te w biegu i te czekające na zgodę Operatora idą jednym
   wykazem, bo w oknie stoją obok siebie. */
async function wypelnijZlecenia(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = cialoPanelu(korzen, 'panel-actions');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.AssistantActionStatus, { windowId: idOkna });
  const wynik = odpowiedz.wynik;
  const zlecenia: AssistantAction[] = wynik === undefined
    ? [] : [...(wynik.actions ?? []), ...(wynik.awaiting ?? [])];
  wykazPanelu(cialo, odpowiedz.udany, odpowiedz.blad?.message,
    wynik === undefined ? undefined : zlecenia,
    'Asystent nie ma zleceń ani czekających, ani w biegu.',
    (z) => [z.title ?? z.id, z.status ?? ''] as const);
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
