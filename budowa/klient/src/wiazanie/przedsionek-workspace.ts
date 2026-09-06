/**
 * Przedsionek środowiska WorkSpace: strefa projektów konta wypełniana rdzeniem,
 * wejście w projekt oraz założenie projektu. Komendy: `project.list`,
 * `project.create`, `session.list`, `session.create` i `workspace.task.list`.
 */

import {
  Command,
  SessionStatus,
  WorkspaceProjectStatus,
  type Session,
  type WorkspaceProject,
} from '../../../shared/contract.ts';
import { miaraSesji } from '../model/miary.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zdanieOdmowy } from './wejscie-odmowa.ts';

const NAGLOWEK = 'Projekty środowiska';

const KOD_SRODOWISKA = 'workspace';

/** Górna granica wykazu sesji; przedsionek pokazuje dorobek, nie archiwum konta. */
const LIMIT_SESJI = 400;

let sesjeSrodowiska: Session[] = [];
let wypelniony = false;

/* Wzór kafla ginie przy pierwszej wymianie zawartości siatki, więc strefa
   zapamiętuje go, zanim po raz pierwszy wypisze projekty rdzenia. */
let wzorzecKafla: HTMLElement | null = null;

/* Widok przedsionka powstaje w płótnie dopiero przy wejściu w środowisko, więc
   wiązanie idzie za jego przebudową obserwatorem, nie odczytem znacznika. */
export function zwiazPrzedsionekWorkspace(kanal: Kanal): void {
  sledzWidok(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || strefaProjektow() === null) return;
    if (cel.closest('[contenteditable]') !== null) return;
    if (cel.closest('[data-nowy-projekt-srodowisko]') !== null) {
      zdarzenie.stopPropagation();
      zdarzenie.preventDefault();
      void zalozProjekt(kanal);
      return;
    }
    const kafel = cel.closest<HTMLElement>('.pd-kafel[data-projekt]');
    /* Kafel z przypisaną sesją przepuszcza kliknięcie do Centrum: to ono
       prowadzi okna robocze i wie, czy okno tej sesji już stoi. */
    if (kafel === null || kafel.dataset.idSesji !== undefined) return;
    zdarzenie.stopPropagation();
    zdarzenie.preventDefault();
    void wejdzWProjekt(kanal, kafel);
  }, true);
}

function sledzWidok(kanal: Kanal): void {
  const obserwator = new MutationObserver(() => {
    const strefa = strefaProjektow();
    if (strefa === null) {
      wypelniony = false;
      return;
    }
    if (wypelniony) return;
    wypelniony = true;
    void wczytaj(kanal, strefa);
  });
  obserwator.observe(document.body, { childList: true, subtree: true });
}

/** Strefa projektów stojącego przedsionka WorkSpace; poza nim brak strefy. */
function strefaProjektow(): HTMLElement | null {
  const widok = document.querySelector<HTMLElement>('.cd-tresc--przedsionek');
  if (widok === null || widok.hidden) return null;
  return widok.querySelector<HTMLElement>(
    `[data-projekty-srodowiska="${KOD_SRODOWISKA}"]`,
  );
}

async function wczytaj(kanal: Kanal, strefa: HTMLElement): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.ProjectList, { includeArchived: false });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(wykaz.blad), 'ostrzezenie');
    /* Bez odpowiedzi rdzenia w siatce stoją kafle wzorcowe prototypu.
       Zostawione podawałyby wymyślone projekty za dorobek Operatora. */
    wypelnijStrefe(strefa, []);
    return;
  }
  sesjeSrodowiska = await pobierzSesje(kanal);
  const projekty = wykaz.wynik.projects.filter(
    (projekt) => projekt.status !== WorkspaceProjectStatus.Archived,
  );
  wypelnijStrefe(strefa, projekty);
  await opiszZadania(kanal, strefa, projekty);
}

/* Sesje bierze się z wykazu konta i zawęża kodem środowiska, bo kafel projektu
   ma otwierać pracę tego środowiska, a nie pierwszą sesję projektu z brzegu. */
async function pobierzSesje(kanal: Kanal): Promise<Session[]> {
  const wynik = await wywolaj(kanal, Command.SessionList, { limit: LIMIT_SESJI });
  const sesje = wynik.wynik?.sessions ?? [];
  return sesje.filter(
    (sesja) => sesja.environmentCode === KOD_SRODOWISKA
      && sesja.status !== SessionStatus.Archived,
  );
}

function wypelnijStrefe(strefa: HTMLElement, projekty: WorkspaceProject[]): void {
  const siatka = strefa.querySelector<HTMLElement>('[data-projekty-lista]');
  if (siatka === null) return;
  const wzor = siatka.querySelector<HTMLElement>('.pd-kafel');
  if (wzor !== null) wzorzecKafla = wzor.cloneNode(true) as HTMLElement;
  const wzorzec = wzorzecKafla;
  if (wzorzec === null) return;
  siatka.replaceChildren(...projekty.map((projekt) => zbudujKafel(wzorzec, projekt)));
  if (projekty.length === 0) siatka.append(zdaniePustej(siatka));
}

function zbudujKafel(wzorzec: HTMLElement, projekt: WorkspaceProject): HTMLElement {
  const kafel = wzorzec.cloneNode(true) as HTMLElement;
  kafel.dataset.projekt = projekt.id;
  delete kafel.dataset.idSesji;
  const nazwa = kafel.querySelector('.pd-kafel-nazwa');
  if (nazwa !== null) nazwa.textContent = projekt.name;
  const opis = kafel.querySelector('.pd-kafel-opis');
  if (opis !== null) {
    /* Opisu pustego nie zastępuje się zdaniem ułożonym po stronie widoku:
       kafel miałby wtedy treść, której rdzeń o projekcie nie zna. */
    if (projekt.description === undefined || projekt.description === '') opis.remove();
    else opis.textContent = projekt.description;
  }
  const miara = kafel.querySelector('.pd-kafel-meta');
  if (miara !== null) miara.textContent = miaraSesji(sesjeProjektu(projekt.id).length);
  return kafel;
}

function zdaniePustej(siatka: HTMLElement): HTMLElement {
  const zdanie = siatka.ownerDocument.createElement('p');
  zdanie.className = 'dn-pusty';
  zdanie.textContent = 'Nie masz jeszcze projektów. Załóż pierwszy przyciskiem „Nowy projekt”.';
  return zdanie;
}

/* Liczba zadań dochodzi osobnym biegiem po wypisaniu kafli: wykaz zadań idzie
   po jednym wywołaniu na projekt i czekanie na całość trzymałoby pustą siatkę. */
async function opiszZadania(
  kanal: Kanal,
  strefa: HTMLElement,
  projekty: WorkspaceProject[],
): Promise<void> {
  const miary = await Promise.all(projekty.map(
    (projekt) => wywolaj(kanal, Command.WorkspaceTaskList, {
      projectId: projekt.id,
      includeDone: true,
      limit: 1,
    }),
  ));
  for (const [pozycja, projekt] of projekty.entries()) {
    const wynik = miary[pozycja];
    if (wynik === undefined || !wynik.udany || wynik.wynik === undefined) continue;
    const miara = strefa.querySelector(
      `.pd-kafel[data-projekt="${projekt.id}"] .pd-kafel-meta`,
    );
    if (miara === null) continue;
    const sesje = miaraSesji(sesjeProjektu(projekt.id).length);
    miara.textContent = `${miaraZadan(wynik.wynik.total)} · ${sesje}`;
  }
}

/** Liczba zadań wraz z odmianą rzeczownika; miary sesji podaje model miar. */
function miaraZadan(ile: number): string {
  if (ile === 1) return '1 zadanie';
  const reszta = ile % 10;
  const setka = ile % 100;
  const mnoga = reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14);
  return `${ile} ${mnoga ? 'zadania' : 'zadań'}`;
}

function sesjeProjektu(idProjektu: string): Session[] {
  return sesjeSrodowiska.filter((sesja) => sesja.projectId === idProjektu);
}

/* Wejście w projekt prowadzi do jego sesji w tym środowisku; projekt bez sesji
   dostaje ją przy wejściu, bo inaczej kafel nie miałby dokąd prowadzić. */
async function wejdzWProjekt(kanal: Kanal, kafel: HTMLElement): Promise<void> {
  const idProjektu = kafel.dataset.projekt ?? '';
  if (idProjektu === '') return;
  const ostatnia = najswiezszaSesja(idProjektu);
  if (ostatnia !== undefined) {
    otworzSesjeKafla(kafel, ostatnia.id);
    return;
  }
  const nazwa = kafel.querySelector('.pd-kafel-nazwa')?.textContent ?? '';
  const zalozona = await wywolaj(kanal, Command.SessionCreate, {
    title: nazwa,
    projectId: idProjektu,
    environmentCode: KOD_SRODOWISKA,
  });
  if (!zalozona.udany || zalozona.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(zalozona.blad), 'ostrzezenie');
    return;
  }
  sesjeSrodowiska = [...sesjeSrodowiska, zalozona.wynik.session];
  otworzSesjeKafla(kafel, zalozona.wynik.session.id);
}

function najswiezszaSesja(idProjektu: string): Session | undefined {
  return sesjeProjektu(idProjektu).reduce<Session | undefined>(
    (najlepsza, sesja) => najlepsza === undefined || sesja.updatedAt > najlepsza.updatedAt
      ? sesja
      : najlepsza,
    undefined,
  );
}

/* Otwarcie sesji prowadzi Centrum po znaczniku sesji, więc kafel bierze ten
   znacznik i sam się naciska — dublowanie otwierania rozjechałoby okna robocze. */
function otworzSesjeKafla(kafel: HTMLElement, idSesji: string): void {
  kafel.dataset.idSesji = idSesji;
  kafel.click();
}

async function zalozProjekt(kanal: Kanal): Promise<void> {
  const strefa = strefaProjektow();
  const siatka = strefa?.querySelector<HTMLElement>('[data-projekty-lista]') ?? null;
  const wzorzec = wzorzecKafla;
  if (strefa === null || siatka === null || wzorzec === null) return;
  const nazwa = await zapytajONazwe(siatka, wzorzec);
  if (nazwa === null || nazwa === '') return;
  const wynik = await wywolaj(kanal, Command.ProjectCreate, { name: nazwa });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'ostrzezenie');
    return;
  }
  wypelniony = false;
  await wczytaj(kanal, strefa);
  wypelniony = true;
}

/* Nazwę projektu Operator wpisuje wprost w kaflu roboczym: przedsionek nie ma
   własnego okna dialogowego, a kafel stoi już w miejscu przyszłego projektu. */
function zapytajONazwe(siatka: HTMLElement, wzorzec: HTMLElement): Promise<string | null> {
  const roboczy = wzorzec.cloneNode(true) as HTMLElement;
  delete roboczy.dataset.projekt;
  roboczy.querySelector('.pd-kafel-opis')?.remove();
  roboczy.querySelector('.pd-kafel-meta')?.remove();
  const pole = roboczy.querySelector<HTMLElement>('.pd-kafel-nazwa');
  if (pole === null) return Promise.resolve(null);
  pole.textContent = '';
  siatka.prepend(roboczy);
  return new Promise((rozstrzygnij) => {
    let domkniete = false;
    const domknij = (wpis: string | null): void => {
      if (domkniete) return;
      domkniete = true;
      roboczy.remove();
      rozstrzygnij(wpis);
    };
    pole.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') {
        zdarzenie.preventDefault();
        domknij((pole.textContent ?? '').trim());
        return;
      }
      if (zdarzenie.key === 'Escape') {
        zdarzenie.preventDefault();
        domknij(null);
      }
    });
    pole.addEventListener('blur', () => domknij(null));
    pole.setAttribute('contenteditable', 'plaintext-only');
    pole.focus();
  });
}
