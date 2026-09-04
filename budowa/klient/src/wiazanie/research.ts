// Wiązanie karty modułu Research z rdzeniem: okno badania, wykaz źródeł
// i wykaz ustaleń. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  ResearchSourceKind,
  type ResearchExcerpt,
  type ResearchExportRecord,
  type ResearchFinding,
  type ResearchMonitor,
  type ResearchReportTemplate,
  type ResearchSource,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  opiszNaglowek,
  zapewnijOknoModulu,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

const KOD_MODULU = 'research';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazBadania(
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
  zwolnijBadania(idKarty);

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
    await Promise.all([
      wypelnijZrodla(kanal, korzen, idOkna),
      wypelnijUstalenia(kanal, korzen, idOkna),
      wypelnijWyciagi(kanal, korzen, idOkna),
      wypelnijObserwacje(kanal, korzen, idOkna),
      wypelnijSzablony(kanal, korzen),
      wypelnijWydania(kanal, korzen, idOkna),
    ]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Research', idOkna);
    if (idOkna === '') {
      niegotowe(panel(korzen, 'panel-sources'), 'Rdzeń nie dał okna badania dla tej karty.');
      niegotowe(panel(korzen, 'panel-findings'), 'Rdzeń nie dał okna badania dla tej karty.');
      return;
    }
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-sources .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void wniesZrodlo(kanal, korzen, idOkna).then(odswiez);
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Badania');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijBadania(idKarty: string): void {
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
  for (const identyfikator of ['panel-sources', 'panel-findings']) {
    niegotowe(panel(korzen, identyfikator), 'Wykaz czeka na odpowiedź rdzenia.');
  }
  for (const zeton of korzen.querySelectorAll('.rs-licznik, .rs-sugestia, .rs-etapy, .rs-graf')) {
    zeton.remove();
  }
  // Nazwa badania, żetony miar i struktura raportu opisują cudzą pracę.
  for (const zeton of korzen.querySelectorAll('.sta-kom-stan ~ * .sta-chip')) zeton.remove();
  const nazwa = korzen.querySelector('.rs-nawig b, #panel-workspace .sta-okno-tresc b');
  if (nazwa !== null) nazwa.textContent = '';
  niegotowe(panel(korzen, 'panel-report'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-pliki'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-przegladarka'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan badania czeka na pierwsze źródło.');
  niegotowe(panel(korzen, 'panel-kolejka'), 'Rdzeń nie podaje zadań w tle dla tego okna.');
}

async function wypelnijZrodla(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-sources');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchSourceList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz źródeł nie doszedł.');
    return;
  }
  const zrodla: ResearchSource[] = odpowiedz.wynik.sources;
  if (zrodla.length === 0) {
    niegotowe(cialo, 'Żadne źródło nie zostało jeszcze wniesione.');
    return;
  }
  cialo.replaceChildren(...zrodla.map((zrodlo) => pozycja(cialo, zrodlo.title, zrodlo.url ?? zrodlo.kind)));
}

async function wypelnijUstalenia(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-findings');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchFindingList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz ustaleń nie doszedł.');
    return;
  }
  const ustalenia: ResearchFinding[] = odpowiedz.wynik.findings;
  if (ustalenia.length === 0) {
    niegotowe(cialo, 'Żadne ustalenie nie zostało jeszcze zapisane.');
    return;
  }
  cialo.replaceChildren(...ustalenia.map((u) => pozycja(cialo, u.content, u.status ?? '')));
}

async function wypelnijWyciagi(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-pliki');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchExcerptList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz wyciągów nie doszedł.');
    return;
  }
  const wyciagi: ResearchExcerpt[] = odpowiedz.wynik.excerpts;
  if (wyciagi.length === 0) {
    niegotowe(cialo, 'Z żadnego źródła nie wyjęto jeszcze cytatu.');
    return;
  }
  cialo.replaceChildren(...wyciagi.map((w) => pozycja(cialo, w.quote, w.sourceTitle ?? '')));
}

async function wypelnijObserwacje(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-przegladarka');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchMonitorList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz obserwacji nie doszedł.');
    return;
  }
  const obserwacje: ResearchMonitor[] = odpowiedz.wynik.monitors;
  if (obserwacje.length === 0) {
    niegotowe(cialo, 'Żadne źródło nie jest obserwowane.');
    return;
  }
  cialo.replaceChildren(...obserwacje.map((o) => pozycja(cialo, o.query ?? o.url ?? o.id, o.kind)));
}

async function wypelnijSzablony(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = panel(korzen, 'panel-report');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchReportTemplateList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz szablonów raportu nie doszedł.');
    return;
  }
  const szablony: ResearchReportTemplate[] = odpowiedz.wynik.templates;
  if (szablony.length === 0) {
    niegotowe(cialo, 'Rdzeń nie ma szablonu raportu.');
    return;
  }
  cialo.replaceChildren(...szablony.map((s) => pozycja(cialo, s.name, String(s.sectionTitles.length) + ' sekcji')));
}

async function wypelnijWydania(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cialo = panel(korzen, 'panel-export');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchExportList, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz wydań badania nie doszedł.');
    return;
  }
  const wydania: ResearchExportRecord[] = odpowiedz.wynik.exports;
  if (wydania.length === 0) {
    niegotowe(cialo, 'Żadnego raportu jeszcze nie wydano.');
    return;
  }
  cialo.replaceChildren(...wydania.map((w) => pozycja(cialo, w.format, w.target ?? '')));
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
async function wniesZrodlo(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  if (idOkna === '') {
    oglos('Badania', 'Okno badania jeszcze nie powstało, więc nie ma do czego wnieść źródła.');
    return;
  }
  const tytul = globalThis.prompt?.('Tytuł źródła') ?? '';
  if (tytul.trim() === '') return;
  const odpowiedz = await wywolaj(kanal, Command.ResearchSourceAdd, {
    windowId: idOkna,
    title: tytul.trim(),
    kind: ResearchSourceKind.Document,
  });
  if (!odpowiedz.udany) {
    oglos('Badania', odpowiedz.blad?.message ?? 'Rdzeń odmówił wniesienia źródła.', 'blad');
    return;
  }
  void korzen;
}
