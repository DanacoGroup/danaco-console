// Wiązanie karty modułu Diagnostics z rdzeniem: wykaz błędów i rekomendacji. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  ExecutionEnv,
  PermissionMode,
  WindowRole,
  WindowStatus,
  type DiagnosticError,
  type DiagnosticRecommendation,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { kartyOkna, oknaRobocze, przypiszOknoKomunikacji } from './okna-robocze.ts';
import { opiszNaglowek, zdejmijSterowanieWspolne, zdejmijTrescWspolna } from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';
import { zapewnijSesje } from './sesja-biezaca.ts';

const KOD_MODULU = 'diagnostics';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazDiagnostyke(
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
  zwolnijDiagnostyke(idKarty);

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
    await Promise.all([wypelnijBledy(kanal, korzen), wypelnijRekomendacje(kanal, korzen)]);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOkno(kanal, idKarty, idOkna);
    if (idOkna === '') return;
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-bledy .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void odswiez();
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Diagnostyka');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijDiagnostyke(idKarty: string): void {
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
  niegotowe(panel(korzen, 'panel-bledy'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-rekomendacje'), 'Wykaz czeka na odpowiedź rdzenia.');
  niegotowe(panel(korzen, 'panel-logi'), 'Rdzeń nie podaje dziennika dla tego okna.');
  niegotowe(panel(korzen, 'panel-centrum'), 'Centrum diagnostyki czeka na pierwszy błąd.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan naprawy powstaje z rekomendacji.');
  niegotowe(panel(korzen, 'panel-artefakty'), 'Rdzeń nie podaje wytworów diagnostyki.');
  niegotowe(panel(korzen, 'panel-kolejka'), 'Rdzeń nie podaje zadań w tle dla tego okna.');
}

async function zapewnijOkno(kanal: Kanal, idKarty: string, stojace: string): Promise<string> {
  if (stojace !== '') return stojace;
  const idSesji = await zapewnijSesje(kanal, idKarty, 'Diagnostics');
  if (idSesji === '') return '';
  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  const idModulu = moduly.wynik?.modules.find((m) => m.code === KOD_MODULU)?.id ?? '';
  if (idModulu === '') return '';
  const wolne = await wskazOknoWolne(kanal, idSesji, idModulu);
  if (wolne !== '') {
    przypiszOknoKomunikacji(idKarty, wolne);
    return wolne;
  }
  const kanaly = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const idKanalu = kanaly.wynik?.channels[0]?.id ?? '';
  if (idKanalu === '') return '';
  const okno = await wywolaj(kanal, Command.WindowCreate, {
    sessionId: idSesji,
    moduleId: idModulu,
    modelChannelId: idKanalu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  const powstale = okno.wynik?.window.id ?? '';
  if (powstale !== '') przypiszOknoKomunikacji(idKarty, powstale);
  return powstale;
}

async function wskazOknoWolne(kanal: Kanal, idSesji: string, idModulu: string): Promise<string> {
  const odpowiedz = await wywolaj(kanal, Command.WindowList, {
    sessionId: idSesji,
    status: WindowStatus.Open,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) return '';
  const zajete = new Set<string>();
  for (const okno of oknaRobocze()) {
    for (const karta of kartyOkna(okno)) {
      if (karta.idOknaKomunikacji !== '') zajete.add(karta.idOknaKomunikacji);
    }
  }
  return odpowiedz.wynik.windows.find(
    (okno) => okno.moduleId === idModulu && !zajete.has(okno.id),
  )?.id ?? '';
}

async function wypelnijBledy(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = panel(korzen, 'panel-bledy');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.DiagnosticsErrorList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz błędów nie doszedł.');
    return;
  }
  const bledy: DiagnosticError[] = odpowiedz.wynik.errors;
  if (bledy.length === 0) {
    niegotowe(cialo, 'Rdzeń nie odnotował żadnego błędu.');
    return;
  }
  cialo.replaceChildren(...bledy.map((b) => pozycja(cialo, b.message, b.source ?? '')));
}

async function wypelnijRekomendacje(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = panel(korzen, 'panel-rekomendacje');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.DiagnosticsRecommendationList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz rekomendacji nie doszedł.');
    return;
  }
  const wykaz: DiagnosticRecommendation[] = odpowiedz.wynik.recommendations;
  if (wykaz.length === 0) {
    niegotowe(cialo, 'Rdzeń nie ma rekomendacji — nie ma z czego ich wyprowadzić.');
    return;
  }
  cialo.replaceChildren(...wykaz.map((r) => pozycja(cialo, r.title, r.detail ?? '')));
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
