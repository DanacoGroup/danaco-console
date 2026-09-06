/**
 * Przedsionek środowiska CodeStudio: repozytoria konta, ostatnie przebiegi
 * budowania i miary kafli modułów pobrane z rdzenia. Komendy: `session.list`,
 * `project.list`, `module.list`, `window.list`, `developer.build.list`
 * oraz `terminal.session.list`.
 */

import {
  BuildStatus,
  Command,
  type DeveloperBuild,
  type Session,
  type SessionPresence,
  type WorkspaceProject,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zdanieOdmowy } from './wejscie-odmowa.ts';

const NAGLOWEK = 'CodeStudio';
const KOD_SRODOWISKA = 'codestudio';
const KOD_MODULU_BUDOWANIA = 'developer';
const KOD_MODULU_POWLOKI = 'terminal';
const LICZBA_PRZEBIEGOW = 6;

const STANY_BUDOWANIA: Readonly<Record<BuildStatus, string>> = {
  [BuildStatus.Running]: 'trwa',
  [BuildStatus.Succeeded]: 'powodzenie',
  [BuildStatus.Failed]: 'błąd',
  [BuildStatus.Stopped]: 'przerwane',
};

interface PrzebiegZSesja {
  budowanie: DeveloperBuild;
  idSesji: string;
}

let sesje: Session[] = [];
let odpisy: SessionPresence[] = [];
let projekty: WorkspaceProject[] = [];
let przebiegi: PrzebiegZSesja[] = [];
let powloki = 0;
let wypelniony = false;
/* Wzorce kafli giną przy pierwszej wymianie zawartości siatki, więc obie
   siatki oddają swój wzór do pamięci wiązania przed pierwszym wypełnieniem. */
const wzorceSiatek = new Map<string, HTMLElement>();

/**
 * Wiąże przedsionek CodeStudio: dane rdzenia wchodzą przy każdym wejściu
 * w środowisko, a wybór repozytorium albo przebiegu wraca do sesji, w której
 * praca stoi. Kafle modułów otwiera Centrum, więc ich kliknięcie idzie dalej.
 */
export function zwiazPrzedsionekCodeStudio(kanal: Kanal): void {
  sledzWidok(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const repozytorium = cel.closest<HTMLElement>('[data-repozytorium]');
    if (repozytorium !== null) {
      zdarzenie.stopPropagation();
      wskazSesjeProjektu(repozytorium.dataset.repozytorium ?? '');
      return;
    }
    const przebieg = cel.closest<HTMLElement>('[data-przebieg]');
    if (przebieg !== null) {
      zdarzenie.stopPropagation();
      wskazSesjeWierszem(przebieg.dataset.sesjaPrzebiegu ?? '');
    }
  }, true);
}

/* Widok przedsionka powstaje w płótnie dopiero przy wejściu w środowisko,
   a kafle modułów wypełnia osobne wiązanie, które zdejmuje miarę kafla.
   Obserwator wraca więc po każdej przebudowie i dopisuje miarę z pamięci. */
function sledzWidok(kanal: Kanal): void {
  const obserwator = new MutationObserver(() => {
    const plotno = plotnoCodeStudio();
    if (plotno === null) {
      wypelniony = false;
      return;
    }
    if (wypelniony) {
      opiszKafleModulow(plotno);
      return;
    }
    wypelniony = true;
    void odswiez(kanal, plotno);
  });
  obserwator.observe(document.body, { childList: true, subtree: true });
}

function plotnoCodeStudio(): HTMLElement | null {
  const widok = document.querySelector<HTMLElement>('.cd-tresc--przedsionek');
  if (widok === null || widok.hidden) return null;
  const listwa = widok.querySelector<HTMLElement>('[data-nowy-projekt-srodowisko]');
  const kod = listwa?.dataset.nowyProjektSrodowisko ?? '';
  return kod === KOD_SRODOWISKA ? widok : null;
}

async function odswiez(kanal: Kanal, plotno: HTMLElement): Promise<void> {
  zapamietajWzorce(plotno);
  await pobierzSesje(kanal);
  await Promise.all([
    pobierzProjekty(kanal),
    pobierzPrzebiegi(kanal),
    pobierzPowloki(kanal),
  ]);
  opiszKafleModulow(plotno);
  wypelnijRepozytoria(plotno);
  wypelnijPrzebiegi(plotno);
}

function zapamietajWzorce(plotno: HTMLElement): void {
  for (const nazwa of ['repozytoria', 'przebiegi']) {
    if (wzorceSiatek.has(nazwa)) continue;
    const kafel = plotno.querySelector<HTMLElement>(`[data-${nazwa}] .pd-kafel`);
    if (kafel !== null) wzorceSiatek.set(nazwa, kafel.cloneNode(true) as HTMLElement);
  }
}

/* Sesje CodeStudio biorą się z wykazu konta, bo przedsionek pyta o nie sam:
   wejście w środowisko prowadzi wspólne wiązanie przedsionków i nie oddaje
   swojej odpowiedzi na zewnątrz. */
async function pobierzSesje(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.SessionList, { includePresence: true });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'ostrzezenie');
    sesje = [];
    odpisy = [];
    return;
  }
  sesje = wynik.wynik.sessions.filter((sesja) => sesja.environmentCode === KOD_SRODOWISKA);
  const znane = new Set(sesje.map((sesja) => sesja.id));
  odpisy = (wynik.wynik.presence ?? []).filter((odpis) => znane.has(odpis.sessionId));
}

async function pobierzProjekty(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ProjectList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'ostrzezenie');
    projekty = [];
    return;
  }
  projekty = wynik.wynik.projects;
}

/* Dziennik budowania kontrakt prowadzi po oknie modułu Developer, nie po
   środowisku, więc przebiegi zbiera się z okien tego modułu i po nich wraca
   do sesji, w której przebieg powstał. */
async function pobierzPrzebiegi(kanal: Kanal): Promise<void> {
  przebiegi = [];
  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  if (!moduly.udany || moduly.wynik === undefined) return;
  const idModulu = moduly.wynik.modules
    .find((modul) => modul.code.toLowerCase() === KOD_MODULU_BUDOWANIA)?.id ?? '';
  if (idModulu === '') return;
  const okna = await wywolaj(kanal, Command.WindowList, {});
  if (!okna.udany || okna.wynik === undefined) return;
  const znane = new Set(sesje.map((sesja) => sesja.id));
  const budowania: PrzebiegZSesja[] = [];
  for (const okno of okna.wynik.windows) {
    if (okno.moduleId !== idModulu || !znane.has(okno.sessionId)) continue;
    const wykaz = await wywolaj(kanal, Command.DeveloperBuildList, {
      windowId: okno.id,
      limit: LICZBA_PRZEBIEGOW,
    });
    if (!wykaz.udany || wykaz.wynik === undefined) continue;
    for (const budowanie of wykaz.wynik.builds) {
      budowania.push({ budowanie, idSesji: okno.sessionId });
    }
  }
  budowania.sort((jeden, drugi) => drugi.budowanie.startedAt - jeden.budowanie.startedAt);
  przebiegi = budowania.slice(0, LICZBA_PRZEBIEGOW);
}

async function pobierzPowloki(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalSessionList, { includeExited: false });
  powloki = wynik.udany && wynik.wynik !== undefined ? wynik.wynik.sessions.length : 0;
}

/* Miarę kafla niesie odpis sesji: to on wiąże sesję z modułem, w którym
   Operator ją prowadzi. Kafel bez ani jednej sesji zostaje bez miary,
   bo liczba zero mówi mniej niż jej brak. */
function opiszKafleModulow(plotno: HTMLElement): void {
  const liczby = new Map<string, number>();
  for (const odpis of odpisy) {
    const kod = (odpis.moduleCode ?? '').toLowerCase();
    if (kod === '') continue;
    liczby.set(kod, (liczby.get(kod) ?? 0) + 1);
  }
  for (const kafel of plotno.querySelectorAll<HTMLElement>('.pd-siatka .pd-kafel[data-modul]')) {
    const kod = (kafel.dataset.modul ?? '').toLowerCase();
    ustawMiareKafla(kafel, opisMiaryModulu(kod, liczby.get(kod) ?? 0));
  }
}

function opisMiaryModulu(kod: string, ile: number): string {
  const sesjami = ile === 0 ? '' : `sesji: ${ile}`;
  if (kod !== KOD_MODULU_POWLOKI || powloki === 0) return sesjami;
  const powlokami = `powłok: ${powloki}`;
  return sesjami === '' ? powlokami : `${sesjami} · ${powlokami}`;
}

/* Miara stoi w kaflu prototypu, ale wiązanie kafli modułów zdejmuje ją przy
   każdej przebudowie siatki, więc pusty kafel dostaje ją tu z powrotem. */
function ustawMiareKafla(kafel: HTMLElement, tresc: string): void {
  const stojaca = kafel.querySelector<HTMLElement>('.pd-kafel-meta');
  if (tresc === '') {
    stojaca?.remove();
    return;
  }
  if (stojaca !== null) {
    wpisz(stojaca, tresc);
    return;
  }
  const miara = kafel.ownerDocument.createElement('span');
  miara.className = 'pd-kafel-meta';
  miara.textContent = tresc;
  kafel.appendChild(miara);
}

/* Repozytorium przedsionka to projekt konta, w którym stoi praca CodeStudio:
   kontrakt nie wiąże projektu ze środowiskiem inaczej niż przez sesje. */
function wypelnijRepozytoria(plotno: HTMLElement): void {
  const siatka = plotno.querySelector<HTMLElement>('[data-repozytoria]');
  const wzor = wzorceSiatek.get('repozytoria');
  if (siatka === null || wzor === undefined) return;
  const liczby = new Map<string, number>();
  for (const sesja of sesje) {
    const projekt = sesja.projectId ?? '';
    if (projekt !== '') liczby.set(projekt, (liczby.get(projekt) ?? 0) + 1);
  }
  const widoczne = projekty.filter((projekt) => liczby.has(projekt.id));
  ustawStrefe(plotno, 'strefaRepozytoria', widoczne.length > 0);
  siatka.replaceChildren(...widoczne.map(
    (projekt) => zbudujKafelRepozytorium(wzor, projekt, liczby.get(projekt.id) ?? 0),
  ));
}

function zbudujKafelRepozytorium(
  wzor: HTMLElement,
  projekt: WorkspaceProject,
  ile: number,
): HTMLElement {
  const kafel = wzor.cloneNode(true) as HTMLElement;
  kafel.dataset.repozytorium = projekt.id;
  const nazwa = kafel.querySelector<HTMLElement>('.pd-kafel-nazwa');
  if (nazwa !== null) wpisz(nazwa, projekt.name);
  const opis = kafel.querySelector<HTMLElement>('.pd-kafel-opis');
  if (opis !== null) {
    if (projekt.description === undefined || projekt.description === '') opis.remove();
    else wpisz(opis, projekt.description);
  }
  ustawMiareKafla(kafel, `sesji: ${ile}`);
  return kafel;
}

function wypelnijPrzebiegi(plotno: HTMLElement): void {
  const siatka = plotno.querySelector<HTMLElement>('[data-przebiegi]');
  const wzor = wzorceSiatek.get('przebiegi');
  if (siatka === null || wzor === undefined) return;
  ustawStrefe(plotno, 'strefaPrzebiegi', przebiegi.length > 0);
  siatka.replaceChildren(...przebiegi.map((pozycja) => zbudujKafelPrzebiegu(wzor, pozycja)));
}

function zbudujKafelPrzebiegu(wzor: HTMLElement, pozycja: PrzebiegZSesja): HTMLElement {
  const kafel = wzor.cloneNode(true) as HTMLElement;
  kafel.dataset.przebieg = pozycja.budowanie.id;
  kafel.dataset.sesjaPrzebiegu = pozycja.idSesji;
  const nazwa = kafel.querySelector<HTMLElement>('.pd-kafel-nazwa');
  if (nazwa !== null) wpisz(nazwa, pozycja.budowanie.task);
  const opis = kafel.querySelector<HTMLElement>('.pd-kafel-opis');
  if (opis !== null) wpisz(opis, nazwaSesji(pozycja.idSesji));
  ustawMiareKafla(kafel, STANY_BUDOWANIA[pozycja.budowanie.status]);
  return kafel;
}

function nazwaSesji(idSesji: string): string {
  const sesja = sesje.find((pozycja) => pozycja.id === idSesji);
  return sesja?.title ?? idSesji;
}

/* Strefa bez ani jednego wiersza schodzi z widoku: pusty nagłówek obiecywałby
   dane, których rdzeń nie ma. */
function ustawStrefe(plotno: HTMLElement, cecha: string, widoczna: boolean): void {
  const nazwaCechy = 'data-' + cecha.replace(/[A-Z]/g, (znak) => '-' + znak.toLowerCase());
  const strefa = plotno.querySelector<HTMLElement>(`[${nazwaCechy}]`);
  if (strefa !== null && strefa.hidden === widoczna) strefa.hidden = !widoczna;
}

function wskazSesjeProjektu(idProjektu: string): void {
  if (idProjektu === '') return;
  const nalezace = sesje
    .filter((sesja) => (sesja.projectId ?? '') === idProjektu)
    .sort((jedna, druga) => druga.updatedAt - jedna.updatedAt);
  const najnowsza = nalezace[0];
  if (najnowsza === undefined) {
    oglos(NAGLOWEK, 'Ten projekt nie ma sesji w środowisku CodeStudio.', 'ostrzezenie');
    return;
  }
  wskazSesjeWierszem(najnowsza.id);
}

/* Otwarcie sesji prowadzi szyna przedsionka wspólna dla środowisk: wybór
   repozytorium naciska jej wiersz, żeby wznowienie i wiązanie szły jedną
   drogą. Wiersz spoza zakresu szyny nie stoi, więc wybór nazywa powód. */
function wskazSesjeWierszem(idSesji: string): void {
  if (idSesji === '') return;
  const wiersz = document.querySelector<HTMLElement>(
    `.cd-tresc--przedsionek .pd-sesja[data-id-sesji="${idSesji}"]`,
  );
  if (wiersz === null) {
    oglos(NAGLOWEK, 'Sesja stoi poza zakresem szyny — zmień zakres na wszystkie.', 'ostrzezenie');
    return;
  }
  wiersz.click();
}

/* Zapis przez porównanie: obserwator widoku odpowiada na każdą zmianę drzewa,
   a wpis tej samej treści zawracałby go w kółko. */
function wpisz(wezel: HTMLElement, tresc: string): void {
  if (wezel.textContent !== tresc) wezel.textContent = tresc;
}

