// Wiązanie karty modułu Design z rdzeniem: zakłada okno modułu, zdejmuje treść
// przykładową prototypu i rozdaje panele podgrupom wiązania.

import {
  Command,
  ExecutionEnv,
  PermissionMode,
  WindowRole,
  WindowStatus,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { oglos } from './ogloszenie.ts';
import {
  opiszNaglowek,
  uzgodnijPrzelacznikiPaneli,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { kartyOkna, oknaRobocze, przypiszOknoKomunikacji } from './okna-robocze.ts';
import { zapewnijSesje } from './sesja-biezaca.ts';
import {
  KOD_MODULU_DESIGN,
  nieGotowe,
  panel,
  poproszony,
  stanPusty,
  tresc,
  type Kontekst,
} from './design-wspolne.ts';
import { odswiezPliki, zapomnijPliki } from './design-wykaz.ts';
import { zapomnijPlansze, zwiazPlansze } from './design-plansza.ts';
import { zwiazZasoby } from './design-zasoby.ts';
import { zwiazPrompt } from './design-prompt.ts';
import { zwiazPodglad } from './design-podglad.ts';
import { zwiazKolekcje } from './design-kolekcje.ts';
import { zwiazZetony } from './design-zetony.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';
import { pokazWarsztaty } from './design-warsztaty.ts';
import { zejdzZKompozycji } from './design-adnotacje.ts';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
  kontekst: Kontekst;
}

const WIAZANIA = new Map<string, WiazanieKarty>();

// Okno stojące podaje wołający; pustka znaczy okno zakładane dla tej karty.
export function zwiazDesign(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (!(korzen instanceof HTMLElement)) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijDesign(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  // Nasłuchy karty schodzą razem z nią, zdjęte sterownikiem przerwania.
  const sterowanie = new AbortController();
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  const kontekst: Kontekst = {
    kanal,
    korzen,
    idKarty,
    przy: { signal: sterowanie.signal },
    stan: stanPusty(),
  };
  WIAZANIA.set(idKarty, { korzen, odlaczenia, kontekst });

  zdejmijTrescPrzykladowa(kontekst);
  zdejmijSterowanieWspolne(korzen);
  zdejmijTrescWspolna(korzen);
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
  uzgodnijPrzelacznikiPaneli(korzen, []);
  opiszSrodowisko(korzen, nazwaSrodowiska);

  zwiazKolekcje(kontekst);
  zwiazZetony(kontekst);
  pokazWarsztaty(kontekst);
  odswiezPliki(kontekst);

  void zapewnijOkno(kontekst, idOknaStojacego).then((gotowe) => {
    if (!gotowe) {
      oglos('Design', 'Rdzeń nie założył okna modułu — panele stoją puste.', 'ostrzezenie');
      return;
    }
    odlaczenia.push(...zwiazZasoby(kontekst));
    odlaczenia.push(...zwiazPlansze(kontekst));
    zwiazPrompt(kontekst);
    zwiazPodglad(kontekst);
  });

  /* Katalog stoi niezależnie od okna: bez kanału modelu rdzeń okna nie założy,
     a operacje rodziny „design” wskazania okna nie wymagają. */
  const katalog = zwiazKatalogModulu(kontekst.kanal, '', korzen, 'design', 'Design');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijDesign(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  // Obecność jest ulotna: rdzeń nie zdejmie jej sam, gdy karta schodzi.
  void zejdzZKompozycji(wiazanie.kontekst);
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
  zapomnijPlansze(idKarty);
  zapomnijPliki(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

function opiszSrodowisko(korzen: HTMLElement, nazwaSrodowiska: string): void {
  if (nazwaSrodowiska === '') return;
  for (const pole of korzen.querySelectorAll('.sta-kom-pole')) {
    if (pole.textContent?.startsWith('Środowisko') !== true) continue;
    const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
    if (dane !== null) dane.textContent = nazwaSrodowiska;
  }
}

// Prototyp niesie treść pokazową w każdym panelu; do czasu odpowiedzi rdzenia
// panel ma stać pusty, nie z cudzym przykładem.
function zdejmijTrescPrzykladowa(kontekst: Kontekst): void {
  for (const identyfikator of ['panel-board', 'panel-assets', 'panel-preview', 'panel-zadania', 'panel-plan', 'panel-pliki']) {
    const cialo = tresc(panel(kontekst.korzen, identyfikator));
    if (cialo === null) continue;
    if (identyfikator === 'panel-board') {
      nieGotowe(cialo.querySelector('.dg-kanwa'), 'Kompozycja czeka na odpowiedź rdzenia.');
      continue;
    }
    if (identyfikator === 'panel-assets') {
      nieGotowe(cialo.querySelector('.dg-siatka'), 'Wykaz zasobów czeka na odpowiedź rdzenia.');
      continue;
    }
    if (identyfikator === 'panel-preview') {
      nieGotowe(cialo.querySelector('.dg-podglad'), 'Podgląd czeka na wskazanie zasobu.');
      continue;
    }
    nieGotowe(cialo, 'Panel czeka na odpowiedź rdzenia.');
  }
}

async function zapewnijOkno(kontekst: Kontekst, idOknaStojacego: string): Promise<boolean> {
  if (idOknaStojacego !== '') {
    const stanowisko = await poproszony(kontekst.kanal, Command.WindowStateGet, {
      windowId: idOknaStojacego,
    });
    if (stanowisko === null) return false;
    kontekst.stan.idOkna = idOknaStojacego;
    kontekst.stan.idSesji = stanowisko.window.sessionId;
    kontekst.stan.idModulu = stanowisko.window.moduleId;
    return true;
  }

  const idSesji = await zapewnijSesje(kontekst.kanal, kontekst.idKarty, 'Design');
  if (idSesji === '') return false;
  kontekst.stan.idSesji = idSesji;

  const idModulu = await wskazModul(kontekst.kanal);
  if (idModulu === '') return false;
  kontekst.stan.idModulu = idModulu;

  // Wykaz okien idzie przed założeniem: zakładanie mnożyłoby okna rdzenia.
  const wolne = await wskazOknoWolne(kontekst.kanal, idSesji, idModulu);
  if (wolne !== '') {
    kontekst.stan.idOkna = wolne;
    przypiszOknoKomunikacji(kontekst.idKarty, wolne);
    return true;
  }

  const idKanalu = await wskazKanalModelu(kontekst.kanal);
  if (idKanalu === '') return false;
  const okno = await poproszony(kontekst.kanal, Command.WindowCreate, {
    sessionId: idSesji,
    moduleId: idModulu,
    modelChannelId: idKanalu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  if (okno === null) return false;
  kontekst.stan.idOkna = okno.window.id;
  przypiszOknoKomunikacji(kontekst.idKarty, okno.window.id);
  return true;
}

async function wskazModul(kanal: Kanal): Promise<string> {
  const odpowiedz = await poproszony(kanal, Command.ModuleList, {});
  const modul = odpowiedz?.modules.find((pozycja) => pozycja.code === KOD_MODULU_DESIGN);
  return modul?.id ?? '';
}

async function wskazOknoWolne(kanal: Kanal, idSesji: string, idModulu: string): Promise<string> {
  const odpowiedz = await poproszony(kanal, Command.WindowList, {
    sessionId: idSesji,
    status: WindowStatus.Open,
  });
  if (odpowiedz === null) return '';
  const zajete = new Set<string>();
  for (const okno of oknaRobocze()) {
    for (const karta of kartyOkna(okno)) {
      if (karta.idOknaKomunikacji !== '') zajete.add(karta.idOknaKomunikacji);
    }
  }
  return odpowiedz.windows.find(
    (okno) => okno.moduleId === idModulu && !zajete.has(okno.id),
  )?.id ?? '';
}

// Wyboru kanału znacznik prototypu nie niesie, a `window.create` go wymaga.
async function wskazKanalModelu(kanal: Kanal): Promise<string> {
  const odpowiedz = await poproszony(kanal, Command.ChannelList, { enabledOnly: true });
  return odpowiedz?.channels[0]?.id ?? '';
}
