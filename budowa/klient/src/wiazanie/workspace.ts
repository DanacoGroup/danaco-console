// Wiązanie karty WorkSpace z rdzeniem: wejście w moduł, projekt bieżący, panele.
import { Command, EventType, type WorkspaceProject } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  uzgodnijPrzelacznikiPaneli,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { sesjaBiezaca } from './sesja-biezaca.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';
import { zwiazAgentow } from './workspace-agenci.ts';
import { zwiazBiblioteke } from './workspace-biblioteka.ts';
import { zwiazInstrukcje } from './workspace-instrukcje.ts';
import { zwiazNotatki } from './workspace-notatki.ts';
import { zwiazPamiec } from './workspace-pamiec.ts';
import {
  opiszKafle,
  opiszOsCzasu,
  opiszProjekt,
  pobierzAktywnosc,
  pobierzPulpit,
  zwiazMenuProjektu,
} from './workspace-projekt.ts';
import { opiszPlan, zwiazZadania } from './workspace-zadania.ts';
import { cialoPanelu, niegotowe, zbierzWezly, type WezlyWorkspace } from './workspace-wezly.ts';

const KOD_MODULU = 'workspace';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

// Każda karta ma własne okno i uchwyty; zamknięcie zwalnia wyłącznie jej wiązanie.
const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazWorkspace(
  kanal: Kanal,
  _nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  const wezlyMoze = zbierzWezly(korzen);
  if (wezlyMoze === null) return false;
  const wezly = wezlyMoze;
  if (wezly === null) return false;
  zwolnijWorkspace(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  // Nasłuchy karty schodzą razem z nią, zdjęte sterownikiem przerwania.
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  let idProjektu = '';
  const projekt = (): string => idProjektu;

  zdejmijTrescWspolna(korzen as HTMLElement);
  zdejmijSterowanieWspolne(korzen as HTMLElement);
  zdejmijPanelBezPokrycia(wezly);
  const zadania = zwiazZadania(wezly, kanal, projekt, przy);
  const instrukcje = zwiazInstrukcje(wezly, kanal, projekt, przy);
  const pamiec = zwiazPamiec(wezly, kanal, projekt, przy);
  const biblioteka = zwiazBiblioteke(wezly, kanal, projekt, przy);
  const notatki = zwiazNotatki(wezly, kanal, projekt, przy);
  const agenci = zwiazAgentow(wezly, kanal, projekt, () => odswiezPulpit(), przy);

  async function odswiezPulpit(): Promise<void> {
    if (idProjektu === '') return;
    const pulpit = await pobierzPulpit(kanal, idProjektu);
    if (pulpit === null) return;
    opiszProjekt(wezly, pulpit.project);
    opiszKafle(wezly, pulpit);
    pamiec.opiszZajetosc(pulpit);
    agenci.opisz(pulpit);
  }

  const odswiezWszystko = async (): Promise<void> => {
    if (idProjektu === '') return;
    await odswiezPulpit();
    opiszOsCzasu(wezly, await pobierzAktywnosc(kanal, idProjektu));
    await Promise.all([
      zadania.odswiez(),
      pamiec.odswiez(),
      biblioteka.odswiez(),
      notatki.odswiez(),
      instrukcje.wczytaj(),
      opiszPlan(wezly, kanal, idProjektu),
    ]);
    uzgodnijPrzelacznikiPaneli(korzen as HTMLElement, []);
  };

  zwiazMenuProjektu(wezly, kanal, projekt, (zmieniony: WorkspaceProject) => {
    opiszProjekt(wezly, zmieniony);
  }, przy);

  odlaczenia.push(
    zglosUchwyt(EventType.WorkspaceProjectChanged, (tresc) => {
      if (idProjektu === '' || tresc.project.id !== idProjektu) return;
      opiszProjekt(wezly, tresc.project);
      void odswiezWszystko();
    }),
  );

  void (async (): Promise<void> => {
    idProjektu = await ustalProjekt(kanal, idOknaStojacego);
    if (idProjektu === '') {
      niegotowe(
        wezly.pulpit.querySelector('.wk-naglowek'),
        'Rdzeń nie wskazał projektu dla tej sesji.',
      );
      return;
    }
    await odswiezWszystko();
  })();

  const katalog = zwiazKatalogModulu(kanal, idOknaStojacego, korzen, 'workspace', 'Warsztat');
  if (katalog !== null) odlaczenia.push(katalog);

  return true;
}

export function zwolnijWorkspace(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

/** Sesja bez projektu każe sięgnąć po wykaz konta: komendy modułu żądają identyfikatora. */
async function ustalProjekt(kanal: Kanal, idOknaStojacego: string): Promise<string> {
  const idSesji = sesjaBiezaca();
  if (idSesji !== '') {
    const wejscie = await wywolaj(kanal, Command.WorkspaceEnter, {
      sessionId: idSesji,
      moduleId: KOD_MODULU,
      ...(idOknaStojacego === '' ? {} : { windowId: idOknaStojacego }),
    });
    const zProjektu = wejscie.wynik?.session.projectId ?? '';
    if (zProjektu !== '') return zProjektu;
  }
  const wykaz = await wywolaj(kanal, Command.ProjectList, {});
  return wykaz.wynik?.projects[0]?.id ?? '';
}

/** Zadania w tle i terminal nie mają odpowiednika w rodzinie `workspace.*`. */
function zdejmijPanelBezPokrycia(wezly: WezlyWorkspace): void {
  niegotowe(cialoPanelu(wezly.wTle), 'Rdzeń nie podaje zadań w tle dla tego modułu.');
  niegotowe(cialoPanelu(wezly.terminal), 'Terminal projektu nie jest jeszcze wystawiony przez rdzeń.');
}
