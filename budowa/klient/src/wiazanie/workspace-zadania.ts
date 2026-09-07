/**
 * Zadania projektu WorkSpace: lista i tablica na pulpicie, przestawianie stanu
 * oraz harmonogram w panelu „Plan". Komendy: `workspace.task.list`,
 * `workspace.board.get`, `workspace.task.update`, `workspace.task.move`
 * i `workspace.schedule.get`.
 */

import {
  Command,
  WorkspaceAssigneeKind,
  WorkspaceCalendarSpan,
  WorkspaceTaskStatus,
  type WorkspaceBoard,
  type WorkspaceScheduleBar,
  type WorkspaceTask,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { data } from './okno-modulu.ts';
import { oglos } from './ogloszenie.ts';
import {
  cialoPanelu,
  kopiaWzoru,
  niegotowe,
  przyjmij,
  ustawPostep,
  wpisz,
  wzorWiersza,
  type WezlyWorkspace,
} from './workspace-wezly.ts';

const NASTEPNY_STAN: Record<string, string> = {
  [WorkspaceTaskStatus.Todo]: WorkspaceTaskStatus.InProgress,
  [WorkspaceTaskStatus.InProgress]: WorkspaceTaskStatus.InReview,
  [WorkspaceTaskStatus.InReview]: WorkspaceTaskStatus.Done,
  [WorkspaceTaskStatus.Done]: WorkspaceTaskStatus.Todo,
  [WorkspaceTaskStatus.Blocked]: WorkspaceTaskStatus.InProgress,
  [WorkspaceTaskStatus.Cancelled]: WorkspaceTaskStatus.Todo,
};

export interface StanZadan {
  odswiez: () => Promise<void>;
}

export function zwiazZadania(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  przy: AddEventListenerOptions,
): StanZadan {
  const lista = wezly.pulpit.querySelector<HTMLElement>('.wk-zadania');
  const wzor = wzorWiersza(lista, '.wk-zadanie');
  const zakladki = [...wezly.pulpit.querySelectorAll<HTMLElement>('.dn-zakladki .dn-zakladka')];
  let tablica = false;

  const odswiez = async (): Promise<void> => {
    const projekt = idProjektu();
    if (projekt === '' || lista === null || wzor === null) return;
    if (tablica) {
      const plansza = await pobierzTablice(kanal, projekt);
      if (plansza !== null) wypelnijTablice(lista, wzor, plansza);
      return;
    }
    const zadania = await pobierzZadania(kanal, projekt);
    wypelnijListe(lista, wzor, zadania);
    opiszPostep(wezly, zadania);
  };

  for (const zakladka of zakladki) {
    zakladka.addEventListener(
      'click',
      () => {
        tablica = (zakladka.textContent ?? '').trim() === 'Tablica';
        for (const inna of zakladki) inna.setAttribute('aria-selected', String(inna === zakladka));
        void odswiez();
      },
      przy,
    );
  }

  lista?.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wierszMoze = cel.closest<HTMLElement>('.wk-zadanie[data-zadanie]');
      if (wierszMoze === null) return;
      const wiersz = wierszMoze;
      const idZadania = wiersz.dataset.zadanie ?? '';
      if (cel.closest('[data-zadanie-usun]') !== null) {
        zdarzenie.stopPropagation();
        void usunZadanie(kanal, wiersz, idZadania, odswiez);
        return;
      }
      if (cel.closest('[data-zadanie-zaleznosc]') !== null) {
        zdarzenie.stopPropagation();
        void przestawZaleznosc(kanal, lista, idProjektu(), idZadania, odswiez);
        return;
      }
      const stan = wiersz.dataset.stan ?? WorkspaceTaskStatus.Todo;
      const kolumna = wiersz.dataset.nastepnaKolumna ?? '';
      void (kolumna === ''
        ? przestawStan(kanal, idZadania, stan, odswiez)
        : przeniesDoKolumny(kanal, idZadania, kolumna, odswiez));
    },
    przy,
  );

  wezly.pulpit.querySelector('[data-zadanie-nowe]')?.addEventListener(
    'click',
    () => {
      void zalozZadanie(kanal, lista, wzor, idProjektu(), odswiez);
    },
    przy,
  );

  return { odswiez };
}

/* Okna zakładania zadania wydanie nie niesie, więc tytuł wchodzi wprost w
   wierszu wstawionym na czoło listy — tą samą drogą, którą idzie nazwa
   komponentu w Centrum. Pusty wpis znaczy odwołanie. */
async function zalozZadanie(
  kanal: Kanal,
  lista: HTMLElement | null,
  wzor: HTMLElement | null,
  idProjektu: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  if (lista === null || wzor === null || idProjektu === '') return;
  const wiersz = kopiaWzoru(wzor);
  if (wiersz === null) return;
  for (const zbedne of wiersz.querySelectorAll(
    '[data-zadanie-usun], [data-zadanie-zaleznosc], .dn-meta',
  )) zbedne.remove();
  const pole = wiersz.querySelector<HTMLElement>('.wk-zadanie-tytul') ?? wiersz;
  lista.prepend(wiersz);
  const tytul = await zapytajWWezle(pole, '');
  wiersz.remove();
  if (tytul === null || tytul === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceTaskCreate, { projectId: idProjektu, title: tytul });
  if (przyjmij('Nowe zadanie', wynik) !== null) await odswiez();
}

/*
przestawZaleznosc wiąże zadanie z poprzednikiem albo znosi wiązanie stojące.

Poprzednikiem jest zadanie stojące w wykazie nad wskazanym: okno nie prowadzi
wyboru zadania, a kolejność wykazu jest jedynym wskazaniem, które Operator widzi.
Zależność stojąca schodzi, więc ten sam przycisk działa w obie strony.
*/
async function przestawZaleznosc(
  kanal: Kanal,
  lista: HTMLElement | null,
  idProjektu: string,
  idZadania: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  if (lista === null || idProjektu === '' || idZadania === '') return;
  const wiersze = [...lista.querySelectorAll<HTMLElement>('.wk-zadanie[data-zadanie]')];
  const numer = wiersze.findIndex((w) => w.dataset.zadanie === idZadania);
  const poprzednik = wiersze[numer - 1]?.dataset.zadanie ?? '';
  if (poprzednik === '') {
    oglos('Zadania', 'Pierwsze zadanie wykazu nie ma poprzednika.', 'ostrzezenie');
    return;
  }
  const stojaca = wiersze[numer]?.dataset.zaleznosc ?? '';
  /* Obie drogi oddają inny kształt wyniku, więc każda przyjmowana jest osobno;
     wspólna zmienna zlewałaby dwa typy odpowiedzi w jeden. */
  if (stojaca === '') {
    const zalozona = await wywolaj(kanal, Command.WorkspaceTaskDependencySet, {
      projectId: idProjektu,
      predecessorTaskId: poprzednik,
      successorTaskId: idZadania,
    });
    if (przyjmij('Zależność zadań', zalozona) !== null) await odswiez();
    return;
  }
  const zniesiona = await wywolaj(kanal, Command.WorkspaceTaskDependencyRemove, {
    dependencyId: stojaca,
  });
  if (przyjmij('Zależność zadań', zniesiona) !== null) await odswiez();
}

/* Usunięcie zadania jest nieodwracalne, więc pierwsze naciśnięcie uzbraja
   przycisk, a dopiero drugie wysyła komendę. */
async function usunZadanie(
  kanal: Kanal,
  wiersz: HTMLElement,
  idZadania: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  if (idZadania === '') return;
  const przycisk = wiersz.querySelector<HTMLElement>('[data-zadanie-usun]');
  if (przycisk !== null && przycisk.dataset.uzbrojone !== 'tak') {
    przycisk.dataset.uzbrojone = 'tak';
    przycisk.setAttribute('aria-label', 'Naciśnij ponownie, aby usunąć zadanie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.WorkspaceTaskDelete, { taskId: idZadania });
  if (przyjmij('Usunięcie zadania', wynik) !== null) await odswiez();
}

/* Wpis w węźle listy: węzeł staje się polem na czas pisania i wraca do swojej
   treści po zatwierdzeniu albo odwołaniu. */
function zapytajWWezle(wezel: HTMLElement, wartosc: string): Promise<string | null> {
  const przed = wezel.textContent ?? '';
  return new Promise((rozstrzygnij) => {
    let domkniete = false;
    const domknij = (wpis: string | null): void => {
      if (domkniete) return;
      domkniete = true;
      wezel.removeAttribute('contenteditable');
      wezel.textContent = przed;
      rozstrzygnij(wpis);
    };
    wezel.setAttribute('contenteditable', 'plaintext-only');
    wezel.textContent = wartosc;
    wezel.focus();
    wezel.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') {
        zdarzenie.preventDefault();
        domknij((wezel.textContent ?? '').trim());
      }
      if (zdarzenie.key === 'Escape') domknij(null);
    });
    wezel.addEventListener('blur', () => domknij(null));
  });
}

function wypelnijListe(lista: HTMLElement, wzor: HTMLElement, zadania: WorkspaceTask[]): void {
  lista.replaceChildren();
  for (const zadanie of zadania) {
    const wiersz = zbudujWiersz(wzor, zadanie, podpisWykonawcy(zadanie));
    if (wiersz !== null) lista.appendChild(wiersz);
  }
}

/** Tablica wchodzi w ten sam węzeł co lista; kolumnę niesie podpis wiersza. */
function wypelnijTablice(lista: HTMLElement, wzor: HTMLElement, plansza: WorkspaceBoard): void {
  lista.replaceChildren();
  const kolumny = [...plansza.columns].sort((jedna, druga) => jedna.order - druga.order);
  for (const [numer, kolumna] of kolumny.entries()) {
    const nastepna = kolumny[numer + 1]?.id ?? kolumny[0]?.id ?? '';
    for (const zadanie of plansza.tasks.filter((cel) => cel.boardColumnId === kolumna.id)) {
      const wiersz = zbudujWiersz(wzor, zadanie, kolumna.name);
      if (wiersz === null) continue;
      wiersz.dataset.nastepnaKolumna = nastepna;
      lista.appendChild(wiersz);
    }
  }
}

function zbudujWiersz(
  wzor: HTMLElement,
  zadanie: WorkspaceTask,
  podpis: string,
): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const znak = wiersz.querySelector<HTMLElement>('.dn-kropka, .pt-tetno');
  const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
  wiersz.replaceChildren();
  if (znak !== null) {
    znak.className = klasaZnaku(zadanie.status);
    wiersz.appendChild(znak);
  }
  wiersz.append(` ${zadanie.title} `);
  if (meta !== null) {
    meta.textContent = podpis;
    wiersz.appendChild(meta);
  }
  wiersz.dataset.zadanie = zadanie.id;
  wiersz.dataset.stan = zadanie.status;
  return wiersz;
}

function klasaZnaku(stan: string): string {
  if (stan === WorkspaceTaskStatus.Done) return 'dn-kropka dn-kropka--sukces';
  if (stan === WorkspaceTaskStatus.InProgress) return 'pt-tetno';
  return 'dn-kropka dn-kropka--neutralna';
}

function podpisWykonawcy(zadanie: WorkspaceTask): string {
  if (zadanie.assigneeKind === WorkspaceAssigneeKind.Operator) return 'Operator';
  if (zadanie.assigneeId !== undefined && zadanie.assigneeId !== '') return zadanie.assigneeId;
  return '';
}

function opiszPostep(wezly: WezlyWorkspace, zadania: WorkspaceTask[]): void {
  const zrobione = zadania.filter((cel) => cel.status === WorkspaceTaskStatus.Done).length;
  const udzial = zadania.length === 0 ? 0 : (zrobione / zadania.length) * 100;
  ustawPostep(wezly.pulpit.querySelector('.dn-postep'), udzial);
  wpisz(
    wezly.pulpit.querySelector('.dn-wykaz-modulu-poz .dn-meta'),
    zadania.length === 0 ? 'Brak zadań' : `Postęp ${Math.round(udzial)}%`,
  );
}

async function pobierzZadania(kanal: Kanal, idProjektu: string): Promise<WorkspaceTask[]> {
  const wynik = await wywolaj(kanal, Command.WorkspaceTaskList, {
    projectId: idProjektu,
    includeDone: true,
    limit: 100,
  });
  return przyjmij('Zadania projektu', wynik)?.tasks ?? [];
}

async function pobierzTablice(kanal: Kanal, idProjektu: string): Promise<WorkspaceBoard | null> {
  const wynik = await wywolaj(kanal, Command.WorkspaceBoardGet, {
    projectId: idProjektu,
    includeDone: true,
  });
  return przyjmij('Tablica zadań', wynik)?.board ?? null;
}

async function przestawStan(
  kanal: Kanal,
  idZadania: string,
  stan: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  if (idZadania === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceTaskUpdate, {
    taskId: idZadania,
    status: (NASTEPNY_STAN[stan] ?? WorkspaceTaskStatus.InProgress) as WorkspaceTaskStatus,
  });
  if (przyjmij('Zadanie', wynik) !== null) await odswiez();
}

async function przeniesDoKolumny(
  kanal: Kanal,
  idZadania: string,
  idKolumny: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  if (idZadania === '' || idKolumny === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceTaskMove, {
    taskId: idZadania,
    boardColumnId: idKolumny,
  });
  if (przyjmij('Przeniesienie zadania', wynik) !== null) await odswiez();
}

/** Panel „Plan" niesie harmonogram: paski zadań w kolejności rozpoczęcia. */
export async function opiszPlan(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: string,
): Promise<void> {
  const cialo = cialoPanelu(wezly.plan);
  if (cialo === null || idProjektu === '') return;
  const naglowek = cialo.querySelector<HTMLElement>('.pt-etykieta');
  const wzor = wzorWiersza(cialo, '.dn-wykaz-modulu-poz');
  const wynik = await wywolaj(kanal, Command.WorkspaceScheduleGet, { projectId: idProjektu });
  const odpowiedz = przyjmij('Harmonogram projektu', wynik);
  cialo.replaceChildren();
  if (naglowek !== null) {
    naglowek.textContent = 'Harmonogram projektu';
    cialo.appendChild(naglowek);
  }
  const paski = odpowiedz?.bars ?? [];
  if (wzor === null || paski.length === 0) {
    niegotowe(cialo, 'Rdzeń nie podał jeszcze pasków harmonogramu dla tego projektu.');
    if (naglowek !== null) cialo.prepend(naglowek);
    return;
  }
  for (const pasek of [...paski].sort((jeden, drugi) => jeden.startAt - drugi.startAt)) {
    const wiersz = wierszPaska(wzor, pasek);
    if (wiersz !== null) cialo.appendChild(wiersz);
  }
  await dopiszKalendarz(cialo, wzor, kanal, idProjektu);
}

/* Harmonogram pokazuje paski zadań, a kalendarz terminy i kamienie milowe —
   to dwa różne wykazy rdzenia, więc miesiąc bieżący dochodzi pod paskami
   zamiast je zastępować. Kotwicą jest chwila otwarcia panelu. */
async function dopiszKalendarz(
  cialo: HTMLElement,
  wzor: HTMLElement | null,
  kanal: Kanal,
  idProjektu: string,
): Promise<void> {
  if (wzor === null) return;
  const wynik = await wywolaj(kanal, Command.WorkspaceCalendarGet, {
    projectId: idProjektu,
    span: WorkspaceCalendarSpan.Month,
    anchorAt: Date.now(),
  });
  const pozycje = przyjmij('Kalendarz projektu', wynik)?.entries ?? [];
  if (pozycje.length === 0) return;
  const naglowek = cialo.ownerDocument.createElement('div');
  naglowek.className = 'pt-etykieta';
  naglowek.textContent = 'Kalendarz miesiąca';
  cialo.appendChild(naglowek);
  for (const pozycja of pozycje) {
    const wiersz = kopiaWzoru(wzor);
    if (wiersz === null) continue;
    const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
    wiersz.replaceChildren();
    wiersz.append(`${pozycja.milestone === true ? '◆ ' : ''}${pozycja.title} `);
    if (meta !== null) {
      meta.textContent = data(pozycja.startAt);
      wiersz.appendChild(meta);
    }
    cialo.appendChild(wiersz);
  }
}

function wierszPaska(wzor: HTMLElement, pasek: WorkspaceScheduleBar): HTMLElement | null {
  const wierszMoze = kopiaWzoru(wzor);
  if (wierszMoze === null) return null;
  const wiersz = wierszMoze;
  const znak = wiersz.querySelector<HTMLElement>('.dn-kropka, .pt-tetno');
  const meta = wiersz.querySelector<HTMLElement>('.dn-meta');
  wiersz.replaceChildren();
  if (znak !== null) {
    znak.className = (pasek.progressPercent ?? 0) >= 100 ? 'dn-kropka dn-kropka--sukces' : 'pt-tetno';
    wiersz.appendChild(znak);
  }
  wiersz.append(` ${pasek.title} `);
  if (meta !== null) {
    meta.textContent = `${data(pasek.startAt)} – ${data(pasek.endAt)}`;
    wiersz.appendChild(meta);
  }
  wiersz.dataset.zadanie = pasek.taskId;
  return wiersz;
}
