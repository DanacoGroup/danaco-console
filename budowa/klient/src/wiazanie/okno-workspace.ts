/**
 * Wiązanie wnętrza okna modułu Workspace z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * workspace.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje. Każda odpowiedź rodziny bierze
 * projekt, więc wnętrze wypełnia się dopiero po jego rozpoznaniu.
 */
import {
  Command,
  EventType,
  WorkspaceProjectStatus,
  type LibraryFile,
  type WorkspaceActivityEntry,
  type WorkspaceDashboard,
  type WorkspaceMemoryEntry,
  type WorkspaceTask,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  data,
  godzina,
  opiszWezel,
  sklonuj,
  uzgodnijPrzelacznikiPaneli,
  wpiszTekst,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';

/** Panele prototypu, którym rodzina workspace.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-zadania', 'panel-plan', 'panel-artefakty', 'panel-terminal'];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  grupa: HTMLElement | null;
  pozycja: HTMLElement | null;
  zadanie: HTMLElement | null;
  wpisOsi: HTMLElement | null;
  wpisPamieci: HTMLElement | null;
  plik: HTMLElement | null;
  ekspert: HTMLElement | null;
}

/** Projekt, którego stan wnętrze pokazuje. Pustka znaczy sesję bez projektu — wtedy rodzina workspace.* nie ma o co pytać. */
interface Stan {
  projekt: string;
}

/**
 * Wiąże wnętrze okna modułu Workspace. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazWorkspace(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { projekt: '' };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt stanu projektu; wołane przy wejściu, po wyborze projektu w szynie i po jego zmianie w rdzeniu. */
  const odswiez = (): void => {
    void wypelnijProjekt(kanal, cialo, wzory, stan);
  };

  zwiazSzyne(cialo, stan, odswiez);
  zwiazSzukaniePamieci(kanal, cialo, wzory, stan);
  zwiazCzynnosciProjektu(kanal, cialo, stan, odswiez);
  zwiazZdarzeniaProjektu(kanal, cialo, stan);

  void wypelnijSzyne(kanal, idOkna, cialo, wzory, stan, odswiez);
}

/** Zdejmuje wzory wierszy i pozycji, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const pozycja = sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .pt-pozycja'));
  // Tętno wskazuje sesję pracującą; karta sesji stanu pracy nie niesie.
  pozycja?.querySelector('.pt-tetno')?.remove();
  pozycja?.removeAttribute('aria-current');
  const zadanie = sklonuj(cialo.querySelector('.wk-zadania .wk-zadanie'));
  // Barwa kropki dzieli zadania na trzy stany, a kontrakt zna ich sześć.
  zadanie?.querySelector('.dn-kropka')?.remove();
  const plik = sklonuj(cialo.querySelector('#panel-library .dn-wykaz-modulu-poz'));
  const ekspert = sklonuj(cialo.querySelector('#panel-agents .dn-wykaz-modulu-poz'));
  // Inicjały eksperta i plakietka domyślności nie mają pola w przypisaniu.
  ekspert?.querySelector('.rt-awatar')?.remove();
  ekspert?.querySelector('.dn-plakietka')?.remove();
  ekspert?.querySelector('.dn-meta')?.remove();
  return {
    grupa: sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .dn-etyk-mono')),
    pozycja,
    zadanie,
    wpisOsi: sklonuj(cialo.querySelector('.wk-os .wk-os-poz')),
    wpisPamieci: sklonuj(cialo.querySelector('#panel-context .rt-kluczowy')),
    plik,
    ekspert,
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów i znaki kontekstu wiszą na
 * rodzinach spoza workspace.*, a profil izolacji, miara różnicy i zasięg
 * wykonania nie mają pola w całym kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  zdejmijTrescWspolna(cialo);
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znak of cialo.querySelectorAll('.sta-kom-kontekst > .sta-zrodlo')) znak.remove();
  oznaczChipyCzynnosci(cialo);
  cialo.querySelector('#panel-instructions .dn-plakietka--ostrzezenie')?.remove();
  cialo.querySelector('#panel-instructions .rt-mp-dol')?.remove();
  cialo.querySelector('#panel-context .dn-alert')?.remove();
  // Zajętość pamięci ma w pulpicie parę pól, których rdzeń nie wypełnia.
  cialo.querySelector('#panel-context .dn-wykaz-modulu-poz')?.remove();
  cialo.querySelector('#panel-context .dn-postep')?.remove();
  cialo.querySelector('#panel-library .rt-mp-dol')?.remove();
  // Warstwa instrukcji i miara tokenów: zasięg jest polem żądania, nie
  // odpowiedzi, a licznika tokenów kontrakt nie prowadzi nigdzie.
  const wiersze = [...cialo.querySelectorAll('#panel-instructions .dn-wykaz-modulu-poz')];
  wiersze[1]?.remove();
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania, profil izolacji i miara
 * różnicy schodzą, a znaki projektu i liczby plików zostają oznaczone do
 * wypełnienia pulpitem projektu.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  chipy[0]?.remove();
  chipy[2]?.remove();
  const projekt = chipy[1];
  const pliki = chipy[3];
  if (projekt !== undefined) projekt.dataset.pole = 'projekt';
  if (pliki !== undefined) pliki.dataset.pole = 'pliki';
  cialo.querySelector('.sta-kontekst-akcji .sta-chip--diff')?.remove();
  // Odłożenie wyniku rozmowy do biblioteki projektu nie ma komendy: rodzina
  // workspace.* bibliotekę czyta, wpisu do niej nie zakłada.
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina workspace.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, założenie projektu i eksperta, wgranie
 * pliku oraz czynności projektu poza zmianą jego stanu.
 */
function zdejmijSterowanieBezPokrycia(cialo: HTMLElement): void {
  uzgodnijPrzelacznikiPaneli(cialo, PANELE_BEZ_POKRYCIA);
  zdejmijSterowanieWspolne(cialo);
  // Założenie projektu bierze jego nazwę, a szyna pola nazwy nie niesie.
  cialo.querySelector('.dn-szyna-modulu-glowa .dn-btn')?.remove();
  cialo.querySelector('#panel-context .sta-okno-akcje .dn-btn-ikona')?.remove();
  cialo.querySelector('#panel-library .sta-okno-akcje .dn-btn-ikona')?.remove();
  cialo.querySelector('#panel-agents .sta-okno-akcje .dn-btn-ikona')?.remove();
  cialo.querySelector('#panel-agents .dn-btn--zarys')?.remove();
  // Widok tablicy ma komendę własną, ale wykaz zadań pulpitu jest listą:
  // rysowanie kolumn w liście byłoby komponowaniem okna.
  cialo.querySelector('.wk-dash .dn-zakladki')?.remove();
  cialo.querySelector('#panel-dashboard .sta-okno-akcje .dn-btn-ikona:last-child')?.remove();
  // W menu projektu zostaje archiwizacja: zmianę stanu projektu komenda zna,
  // powielenia, wywozu i usunięcia — nie.
  const menu = [...cialo.querySelectorAll<HTMLElement>('#menu-proj .sta-menu-poz')];
  for (const [numer, pozycja] of menu.entries()) {
    if (numer !== 2) pozycja.remove();
  }
}

/** Wczytuje projekt okna i wykaz projektów kont do szyny; sesja bez projektu zostawia wnętrze puste, bo rodzina workspace.* nie ma o co pytać. */
async function wypelnijSzyne(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): Promise<void> {
  const okno = await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna });
  const wejscie = okno.udany && okno.wynik !== undefined
    ? await wywolaj(kanal, Command.WorkspaceEnter, {
      sessionId: okno.wynik.window.sessionId,
      moduleId: okno.wynik.window.moduleId,
      windowId: idOkna,
    })
    : undefined;
  stan.projekt = wejscie?.udany === true ? (wejscie.wynik?.session.projectId ?? '') : '';

  const sesje = await wywolaj(kanal, Command.SessionList, {});
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista !== null && sesje.udany && sesje.wynik !== undefined) {
    const wykaz = new Map<string, string[]>();
    for (const sesja of sesje.wynik.sessions) {
      const projekt = sesja.projectId ?? '';
      if (projekt === '') continue;
      wykaz.set(projekt, [...(wykaz.get(projekt) ?? []), sesja.title ?? '']);
    }
    for (const [projekt, tytuly] of wykaz) {
      await postawGrupe(kanal, lista, wzory, projekt, tytuly, stan);
    }
  }
  odswiez();
}

/** Stawia w szynie grupę jednego projektu wraz z jego kartami sesji; nazwę grupy niesie pulpit projektu, nie karta sesji. */
async function postawGrupe(
  kanal: Kanal,
  lista: HTMLElement,
  wzory: Wzory,
  projekt: string,
  tytuly: string[],
  stan: Stan,
): Promise<void> {
  const pulpit = await wywolaj(kanal, Command.WorkspaceDashboardGet, { projectId: projekt });
  const nazwa = pulpit.udany ? (pulpit.wynik?.dashboard.project.name ?? '') : '';
  if (nazwa === '' || wzory.grupa === null) return;
  const grupa = wzory.grupa.cloneNode(true) as HTMLElement;
  grupa.textContent = nazwa;
  lista.appendChild(grupa);
  if (wzory.pozycja === null) return;
  for (const tytul of tytuly) {
    if (tytul === '') continue;
    const pozycja = wzory.pozycja.cloneNode(true) as HTMLElement;
    pozycja.dataset.projekt = projekt;
    if (projekt === stan.projekt) pozycja.setAttribute('aria-current', 'true');
    const tytulPozycji = pozycja.querySelector('.pt-pozycja-tytul');
    if (tytulPozycji === null) continue;
    tytulPozycji.textContent = tytul;
    lista.appendChild(pozycja);
  }
}

/** Wiąże wybór pozycji szyny z przestawieniem wnętrza na projekt tej karty sesji. */
function zwiazSzyne(cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null) return;
  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('.pt-pozycja');
    const projekt = pozycja?.dataset.projekt;
    if (pozycja === null || pozycja === undefined || projekt === undefined) return;
    for (const inna of lista.querySelectorAll('.pt-pozycja')) inna.removeAttribute('aria-current');
    pozycja.setAttribute('aria-current', 'true');
    stan.projekt = projekt;
    odswiez();
  });
}

/** Pyta rdzeń o stan projektu i wypełnia nim pulpit, instrukcje, pamięć, bibliotekę i wykaz ekspertów. */
async function wypelnijProjekt(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const pulpit = stan.projekt === ''
    ? undefined
    : await wywolaj(kanal, Command.WorkspaceDashboardGet, { projectId: stan.projekt });
  const zestawienie = pulpit?.udany === true ? pulpit.wynik?.dashboard : undefined;
  opiszPulpit(cialo, zestawienie);
  wypelnijEkspertow(cialo, wzory, zestawienie?.assignedAgentIds ?? []);
  await wypelnijZadania(kanal, cialo, wzory, stan, zestawienie);
  await wypelnijOsCzasu(kanal, cialo, wzory, stan);
  await wypelnijInstrukcje(kanal, cialo, stan);
  await wypelnijPamiec(kanal, cialo, wzory, stan, '');
  await wypelnijBiblioteke(kanal, cialo, wzory, stan);
}

/**
 * Nanosi projekt na nagłówek pulpitu, znaki pasa czynności i kafle
 * nawigacyjne. Sesja poza projektem zostawia je puste — projekt wskazany
 * w szynie wchodzi w te same węzły, więc pustka ich nie zdejmuje.
 */
function opiszPulpit(cialo: HTMLElement, zestawienie: WorkspaceDashboard | undefined): void {
  const nazwa = zestawienie?.project.name ?? '';
  const naglowek = cialo.querySelector<HTMLElement>('.wk-naglowek h2');
  if (naglowek !== null) naglowek.textContent = nazwa;
  const opis = cialo.querySelector<HTMLElement>('.wk-naglowek p');
  if (opis !== null) opis.textContent = zestawienie?.project.description ?? '';
  const stanProjektu = cialo.querySelector<HTMLElement>('#panel-dashboard .dn-plakietka--sukces');
  opiszWezel(stanProjektu, zestawienie === undefined ? '' : ` ${zestawienie.project.status}`);
  opiszWezel(cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'), nazwa);
  opiszZnak(cialo, 'projekt', nazwa);
  opiszZnak(cialo, 'pliki', liczba(zestawienie?.libraryFileCount));
  opiszKafle(cialo, zestawienie);
}

/** Nanosi wartość na znak pasa czynności rozpoznany po oznaczeniu. */
function opiszZnak(cialo: HTMLElement, pole: string, wartosc: string): void {
  opiszWezel(cialo.querySelector(`.sta-kontekst-akcji [data-pole="${pole}"]`), wartosc);
}

/** Nanosi liczniki pulpitu na kafle nawigacyjne; kafel bez licznika zostaje samą nazwą panelu, bo nazwa panelu wartością nie jest. */
function opiszKafle(cialo: HTMLElement, zestawienie: WorkspaceDashboard | undefined): void {
  const liczniki = [
    liczba(zestawienie?.instructionSetCount),
    liczba(zestawienie?.memoryEntryCount),
    liczba(zestawienie?.libraryFileCount),
    zestawienie === undefined ? '' : String(zestawienie.assignedAgentIds?.length ?? 0),
  ];
  for (const [numer, kafel] of [...cialo.querySelectorAll('.wk-kafle .wk-kafel')].entries()) {
    const licznik = kafel.querySelector('span');
    if (licznik !== null) licznik.textContent = liczniki[numer] ?? '';
  }
}

/** Liczba pulpitu w zapisie znacznika; pole nieoddane przez rdzeń daje pustkę, która zdejmuje węzeł. */
function liczba(wartosc: number | undefined): string {
  return wartosc === undefined ? '' : String(wartosc);
}

/** Stawia w wykazie zadań po jednym wierszu na zadanie projektu i nanosi postęp policzony z liczników pulpitu. */
async function wypelnijZadania(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  zestawienie: WorkspaceDashboard | undefined,
): Promise<void> {
  opiszPostep(cialo, zestawienie);
  const wykaz = cialo.querySelector<HTMLElement>('.wk-zadania');
  if (wykaz === null || wzory.zadanie === null) return;
  wykaz.replaceChildren();
  if (stan.projekt === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceTaskList, {
    projectId: stan.projekt,
    includeDone: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const zadanie of wynik.wynik.tasks) {
    const wiersz = wzory.zadanie.cloneNode(true) as HTMLElement;
    if (!wpiszTekst(wiersz, zadanie.title)) continue;
    opiszMetaZadania(wiersz, zadanie);
    wykaz.appendChild(wiersz);
  }
}

/** Nanosi na wiersz zadania jego stan i wykonawcę; wykonawca wchodzi identyfikatorem, bo przypisanie nazwiska nie niesie. */
function opiszMetaZadania(wiersz: HTMLElement, zadanie: WorkspaceTask): void {
  const meta = wiersz.querySelector('.dn-meta');
  if (meta === null) return;
  const wykonawca = zadanie.assigneeId ?? '';
  meta.textContent = wykonawca === '' ? zadanie.status : `${zadanie.status} · ${wykonawca}`;
}

/**
 * Nanosi postęp projektu policzony z liczby zadań i zadań ukończonych. Brak
 * któregokolwiek licznika zostawia miarę pustą, a pasek na zerze — postęp
 * wpisany bez liczników byłby postępem wymyślonym.
 */
function opiszPostep(cialo: HTMLElement, zestawienie: WorkspaceDashboard | undefined): void {
  const razem = zestawienie?.taskCount ?? 0;
  const gotowe = zestawienie?.taskDoneCount;
  const znany = zestawienie !== undefined && gotowe !== undefined && razem !== 0;
  const procent = znany ? Math.round(((gotowe ?? 0) / razem) * 100) : 0;
  const miara = cialo.querySelector<HTMLElement>('.wk-dash .dn-wykaz-modulu-poz .dn-meta');
  if (miara !== null) miara.textContent = znany ? `Postęp ${procent}%` : '';
  const pasek = cialo.querySelector<HTMLElement>('.wk-dash .dn-postep');
  if (pasek === null) return;
  pasek.dataset.postepDo = String(procent);
  const wartosc = pasek.querySelector<HTMLElement>('.dn-postep-wartosc');
  if (wartosc !== null) wartosc.style.width = `${procent}%`;
}

/** Stawia na osi czasu po jednym wpisie na zdarzenie projektu; godzina bierze się ze znacznika zdarzenia. */
async function wypelnijOsCzasu(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const os = cialo.querySelector<HTMLElement>('.wk-os');
  if (os === null || wzory.wpisOsi === null) return;
  os.replaceChildren();
  if (stan.projekt === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceActivityList, { projectId: stan.projekt });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const zdarzenie of wynik.wynik.entries) {
    const wpis = postawWpisOsi(wzory.wpisOsi, zdarzenie);
    if (wpis !== null) os.appendChild(wpis);
  }
}

/** Składa wpis osi czasu z godziny i streszczenia zdarzenia; wpis bez miejsca na tekst nie wchodzi. */
function postawWpisOsi(wzor: HTMLElement, zdarzenie: WorkspaceActivityEntry): HTMLElement | null {
  const wpis = wzor.cloneNode(true) as HTMLElement;
  const czas = wpis.querySelector('time');
  if (czas !== null) czas.textContent = godzina(zdarzenie.occurredAt);
  return wpiszTekst(wpis, ` ${zdarzenie.summary}`) ? wpis : null;
}

/** Wpisuje w panel instrukcji treść wersji obowiązującej wraz z jej miarą; brak wersji zostawia panel pusty. */
async function wypelnijInstrukcje(kanal: Kanal, cialo: HTMLElement, stan: Stan): Promise<void> {
  const konsola = cialo.querySelector<HTMLElement>('#panel-instructions .pt-konsola');
  const miara = cialo.querySelector<HTMLElement>('#panel-instructions .dn-meta');
  const wynik = stan.projekt === ''
    ? undefined
    : await wywolaj(kanal, Command.WorkspaceInstructionsVersionList, {
      projectId: stan.projekt,
      limit: 1,
    });
  const wersja = wynik?.udany === true ? wynik.wynik?.versions[0] : undefined;
  if (konsola !== null) konsola.textContent = wersja?.content ?? '';
  if (miara !== null) {
    miara.textContent = wersja === undefined ? '' : `${wersja.content.length} znaków`;
  }
}

/** Stawia w panelu pamięci po jednym wpisie na wpis pamięci projektu i nanosi ich liczbę na plakietkę panelu. */
async function wypelnijPamiec(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  fraza: string,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-context .sta-okno-tresc');
  if (panel === null) return;
  for (const stojacy of panel.querySelectorAll('.rt-kluczowy')) stojacy.remove();
  const wynik = stan.projekt === ''
    ? undefined
    : await wywolaj(kanal, Command.WorkspaceContextGet, {
      projectId: stan.projekt,
      query: fraza === '' ? undefined : fraza,
    });
  const wpisy = wynik?.udany === true ? (wynik.wynik?.entries ?? []) : [];
  const plakietka = cialo.querySelector<HTMLElement>('#panel-context .sta-okno-belka .dn-plakietka');
  if (plakietka !== null) plakietka.textContent = String(wpisy.length);
  if (wzory.wpisPamieci === null) return;
  for (const wpis of wpisy) panel.appendChild(postawWpisPamieci(wzory.wpisPamieci, wpis));
}

/** Składa wpis pamięci z jego treści oraz podpisu złożonego z etykiet, daty założenia i pochodzenia. */
function postawWpisPamieci(wzor: HTMLElement, wpis: WorkspaceMemoryEntry): HTMLElement {
  const wpisWezel = wzor.cloneNode(true) as HTMLElement;
  wpiszTekst(wpisWezel, wpis.content);
  const podpis = wpisWezel.querySelector('small');
  if (podpis !== null) {
    const etykiety = (wpis.tags ?? []).join(' · ');
    const opis = `${data(wpis.createdAt)} · ${wpis.origin}`;
    podpis.textContent = etykiety === '' ? opis : `${etykiety} · ${opis}`;
  }
  return wpisWezel;
}

/** Wiąże pole zawężania panelu pamięci z wyszukiwaniem wpisów po frazie. */
function zwiazSzukaniePamieci(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): void {
  const pole = cialo.querySelector<HTMLInputElement>('#panel-context .dn-szukaj input');
  if (pole === null) return;
  pole.addEventListener('input', () => {
    if (stan.projekt === '') return;
    void wypelnijPamiec(kanal, cialo, wzory, stan, pole.value.trim());
  });
}

/** Stawia w bibliotece projektu po jednym wierszu na plik i nanosi ich liczbę na plakietkę panelu. */
async function wypelnijBiblioteke(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-library .sta-okno-tresc');
  if (panel === null) return;
  panel.replaceChildren();
  const wynik = stan.projekt === ''
    ? undefined
    : await wywolaj(kanal, Command.WorkspaceLibraryList, { projectId: stan.projekt });
  const pliki = wynik?.udany === true ? (wynik.wynik?.files ?? []) : [];
  const plakietka = cialo.querySelector<HTMLElement>('#panel-library .sta-okno-belka .dn-plakietka');
  if (plakietka !== null) plakietka.textContent = String(pliki.length);
  if (wzory.plik === null) return;
  for (const plik of pliki) {
    const wiersz = postawWierszPliku(wzory.plik, plik);
    if (wiersz !== null) panel.appendChild(wiersz);
  }
}

/** Składa wiersz pliku projektu z jego nazwy i etykiet; plik bez etykiet traci podpis, bo podpis pusty niczego nie mówi. */
function postawWierszPliku(wzor: HTMLElement, plik: LibraryFile): HTMLElement | null {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  if (!wpiszTekst(wiersz, `${plik.name} `)) return null;
  const meta = wiersz.querySelector('.dn-meta');
  const etykiety = (plik.tags ?? []).join(', ');
  if (meta !== null) {
    if (etykiety === '') meta.remove();
    else meta.textContent = etykiety;
  }
  return wiersz;
}

/** Stawia w menedżerze ekspertów po jednym wierszu na przypisanie i nanosi ich liczbę na plakietkę panelu. */
function wypelnijEkspertow(cialo: HTMLElement, wzory: Wzory, eksperci: string[]): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-agents .sta-okno-tresc');
  if (panel === null) return;
  panel.replaceChildren();
  const plakietka = cialo.querySelector<HTMLElement>('#panel-agents .sta-okno-belka .dn-plakietka');
  if (plakietka !== null) plakietka.textContent = String(eksperci.length);
  if (wzory.ekspert === null) return;
  for (const ekspert of eksperci) {
    const wiersz = wzory.ekspert.cloneNode(true) as HTMLElement;
    if (wpiszTekst(wiersz, ekspert)) panel.appendChild(wiersz);
  }
}

/** Wiąże archiwizację projektu z menu pulpitu; czynność bez projektu nie ma na czym stanąć. */
function zwiazCzynnosciProjektu(
  kanal: Kanal,
  cialo: HTMLElement,
  stan: Stan,
  odswiez: () => void,
): void {
  cialo.querySelector('#menu-proj .sta-menu-poz')?.addEventListener('click', () => {
    if (stan.projekt === '') return;
    void wywolaj(kanal, Command.WorkspaceProjectStatusSet, {
      projectId: stan.projekt,
      status: WorkspaceProjectStatus.Archived,
    }).then(odswiez);
  });
}

/** Nasłuchuje zmian projektu: nagłówek pulpitu bierze nazwę, opis i stan projektu po zmianie. */
function zwiazZdarzeniaProjektu(kanal: Kanal, cialo: HTMLElement, stan: Stan): void {
  kanal.naZdarzenie(EventType.WorkspaceProjectChanged, (tresc) => {
    if (tresc.project.id !== stan.projekt) return;
    const naglowek = cialo.querySelector<HTMLElement>('.wk-naglowek h2');
    if (naglowek !== null) naglowek.textContent = tresc.project.name;
    opiszWezel(
      cialo.querySelector('#panel-dashboard .dn-plakietka--sukces'),
      ` ${tresc.project.status}`,
    );
  });
}
