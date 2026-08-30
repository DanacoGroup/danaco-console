/**
 * Wiązanie wnętrza okna modułu Developer z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * developer.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  BuildStatus,
  Command,
  ContextualOpKind,
  EventType,
  GitActionKind,
  GitDiffLineKind,
  GitFileState,
  TreeNodeKind,
  type DeveloperBuild,
  type DeveloperGitStatus,
  type DeveloperTreeNode,
  type GitStatusEntry,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina developer.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = [
  'panel-plan',
  'panel-zadania',
  'panel-terminal',
  'panel-kolejka',
  'panel-subagenci',
  'panel-artefakty',
];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  katalog: HTMLElement | null;
  plik: HTMLElement | null;
  plakietka: HTMLElement | null;
  zmiana: HTMLElement | null;
  wierszDodany: HTMLElement | null;
  wierszUsuniety: HTMLElement | null;
  karta: HTMLElement | null;
}

/** Katalog roboczy okna i ścieżka pliku wczytanego do edytora; obie wartości pochodzą z odpowiedzi rdzenia. */
interface Stan {
  katalog: string;
  plik: string;
}

/**
 * Wiąże wnętrze okna modułu Developer. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazDeveloper(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { katalog: '', plik: '' };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt drzewa i stanu repozytorium; wołane przy wejściu i po każdej czynności repozytorium. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, stan);
  };

  zwiazDrzewo(kanal, idOkna, cialo, wzory, stan);
  zwiazEdytor(kanal, idOkna, cialo, wzory, stan);
  zwiazCzynnosciRepozytorium(kanal, idOkna, cialo, odswiez);
  zwiazZdarzeniaBudowania(kanal, idOkna, cialo);

  odswiez();
  void wypelnijBudowanie(kanal, idOkna, cialo);
}

/** Zdejmuje wzory wierszy i pozycji, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const wpisy = [...cialo.querySelectorAll<HTMLElement>('#panel-drzewo .dv-drzewo .wpis')];
  const katalog = oczyscWpis(sklonuj(wpisy.find((wpis) => nazwaWpisu(wpis).endsWith('/')) ?? null));
  const plik = oczyscWpis(sklonuj(wpisy.find((wpis) => !nazwaWpisu(wpis).endsWith('/')) ?? null));
  const zmiana = sklonuj(cialo.querySelector('#panel-git .dn-wykaz-modulu-poz'));
  // Metryka wierszy dodanych i usuniętych na plik nie ma pola ani w stanie
  // repozytorium, ani w zatwierdzeniu, ani we fragmencie różnicy.
  zmiana?.querySelector('.dn-meta')?.remove();
  const karta = sklonuj(cialo.querySelector('#panel-edytor .dv-karta'));
  karta?.querySelector('.pt-tetno')?.remove();
  return {
    katalog,
    plik,
    plakietka: sklonuj(cialo.querySelector('#panel-drzewo .dv-drzewo .dn-plakietka')),
    zmiana,
    wierszDodany: sklonuj(cialo.querySelector('#panel-git .pt-wiersz-dodany')),
    wierszUsuniety: sklonuj(cialo.querySelector('#panel-git .pt-wiersz-usuniety')),
    karta,
  };
}

/** Nazwa wpisu drzewa z treści przykładowej; ukośnik na końcu odróżnia wzór katalogu od wzoru pliku. */
function nazwaWpisu(wpis: HTMLElement): string {
  return (wpis.textContent ?? '').trim();
}

/** Zdejmuje ze wzoru wpisu drzewa oznaczenie stanu i wskazanie pozycji bieżącej; oba nadaje dopiero odpowiedź rdzenia. */
function oczyscWpis(wzor: HTMLElement | null): HTMLElement | null {
  if (wzor === null) return null;
  wzor.querySelector('.dn-plakietka')?.remove();
  wzor.removeAttribute('aria-current');
  wzor.classList.remove('wciecie');
  return wzor;
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów, historia i wpisy nagłówka
 * wiszą na rodzinach spoza developer.*, więc wiązanie ich nie wypełnia; miara
 * różnicy, drzewo robocze i przypięty plik nie mają pola w całym kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  zdejmijPoleNaglowka(cialo, 'Model');
  zdejmijPoleNaglowka(cialo, 'Wysiłek');
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  oznaczChipyKontekstu(cialo);
  oznaczChipyCzynnosci(cialo);
}

/** Zdejmuje pole nagłówka rozpoznane po podpisie; rodzina developer.* nie oddaje ani modelu kanału, ani nakładu rozumowania. */
function zdejmijPoleNaglowka(cialo: HTMLElement, podpis: string): void {
  const pola = [...cialo.querySelectorAll('.sta-kom-naglowek .sta-kom-pole')];
  pola.find((pole) => pole.textContent?.startsWith(podpis) === true)?.remove();
}

/** Zostawia w pasie kontekstu wyłącznie znak gałęzi i oznacza go do wypełnienia; przypięty plik nie ma pola w kontrakcie. */
function oznaczChipyKontekstu(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kom-kontekst > .sta-zrodlo')];
  chipy[0]?.remove();
  const galaz = chipy[1];
  if (galaz !== undefined) galaz.dataset.pole = 'galaz';
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania i znak drzewa roboczego
 * schodzą, miara różnicy schodzi, a znaki katalogu i gałęzi zostają oznaczone
 * do wypełnienia stanem repozytorium.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  chipy[0]?.remove();
  chipy[3]?.remove();
  const katalog = chipy[1];
  const galaz = chipy[2];
  if (katalog !== undefined) katalog.dataset.pole = 'katalog';
  if (galaz !== undefined) galaz.dataset.pole = 'galaz';
  cialo.querySelector('.sta-kontekst-akcji .sta-chip--diff')?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina developer.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, wybór modelu i nakładu, urządzenia
 * wejścia dźwięku oraz konektory i wtyczki menu dodawania.
 */
function zdejmijSterowanieBezPokrycia(cialo: HTMLElement): void {
  for (const nazwa of PANELE_BEZ_POKRYCIA) {
    cialo.querySelector(`#${nazwa}`)?.remove();
    cialo.querySelector(`[data-panel-toggle="${nazwa}"]`)?.remove();
  }
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    przelacznik.setAttribute('aria-checked', String(panel?.hasAttribute('hidden') === false));
  }
  cialo.querySelector('#pop-mik')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-model')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-wysilek')?.closest('.sta-nrz')?.remove();
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
  zdejmijKonektory(cialo);
}

/** Zostawia w menu dodawania sam tytuł i trzy pozycje plików; konektory i wtyczki należą do rodzin spoza developer.*. */
function zdejmijKonektory(cialo: HTMLElement): void {
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Pyta rdzeń o drzewo projektu i stan repozytorium, po czym wypełnia nimi belki, znaki, drzewo i panel repozytorium. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const drzewo = await wywolaj(kanal, Command.DeveloperTreeGet, { windowId: idOkna });
  const repozytorium = await wywolaj(kanal, Command.DeveloperGitStatus, { windowId: idOkna });
  const wykaz = drzewo.udany ? drzewo.wynik : undefined;
  const stanRepo = repozytorium.udany ? repozytorium.wynik?.status : undefined;
  stan.katalog = wykaz?.root ?? '';
  opiszKatalog(cialo, stan.katalog);
  opiszGalaz(cialo, stanRepo?.branch ?? '');
  wypelnijDrzewo(cialo, wzory, wykaz?.nodes ?? [], stanRepo);
  wypelnijZmiany(cialo, wzory, stanRepo);
  await wypelnijRoznice(kanal, idOkna, cialo, wzory, []);
}

/** Nanosi katalog roboczy okna na plakietkę belki, znak pasa czynności i znacznik panelu drzewa; pustka zdejmuje każdy z nich. */
function opiszKatalog(cialo: HTMLElement, katalog: string): void {
  const miejsca = [
    cialo.querySelector<HTMLElement>('#okno-czat-1 .sta-okno-belka > .sta-chip'),
    cialo.querySelector<HTMLElement>('.sta-kontekst-akcji [data-pole="katalog"]'),
    cialo.querySelector<HTMLElement>('#panel-drzewo .sta-okno-znacznik'),
  ];
  for (const miejsce of miejsca) {
    if (miejsce === null) continue;
    if (katalog === '' || !wpiszTekst(miejsce, katalog)) miejsce.remove();
  }
}

/** Nanosi gałąź bieżącą na znaki kontekstu i czynności oraz na znacznik panelu repozytorium; odłączona głowa zdejmuje każdy z nich. */
function opiszGalaz(cialo: HTMLElement, galaz: string): void {
  const miejsca = [
    cialo.querySelector<HTMLElement>('.sta-kom-kontekst [data-pole="galaz"]'),
    cialo.querySelector<HTMLElement>('.sta-kontekst-akcji [data-pole="galaz"]'),
    cialo.querySelector<HTMLElement>('#panel-git .sta-okno-znacznik'),
  ];
  for (const miejsce of miejsca) {
    if (miejsce === null) continue;
    if (galaz === '' || !wpiszTekst(miejsce, galaz)) miejsce.remove();
  }
}

/** Stawia w drzewie po jednym wpisie na węzeł odpowiedzi; rodzaj węzła rozstrzyga wzór, a ścieżka nadrzędna wcięcie. */
function wypelnijDrzewo(
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: DeveloperTreeNode[],
  stanRepo: DeveloperGitStatus | undefined,
): void {
  const drzewo = cialo.querySelector<HTMLElement>('#panel-drzewo .dv-drzewo');
  if (drzewo === null) return;
  drzewo.replaceChildren();
  const stany = stanyPlikow(stanRepo);
  for (const wezel of wezly) {
    const katalog = wezel.kind === TreeNodeKind.Directory;
    const wzor = katalog ? wzory.katalog : wzory.plik;
    if (wzor === null) continue;
    const wpis = wzor.cloneNode(true) as HTMLElement;
    wpis.classList.toggle('wciecie', (wezel.parentPath ?? '') !== '');
    if (!katalog) wpis.dataset.sciezka = wezel.path;
    wpis.dataset.nazwa = wezel.name;
    if (!wpiszTekst(wpis, wezel.name)) continue;
    oznaczStanPliku(wpis, wzory.plakietka, stany.get(wezel.path));
    drzewo.appendChild(wpis);
  }
}

/** Stany plików repozytorium w spisie po ścieżce; pozycje o stanie niezmienionym odpadają, bo nie mają czego oznaczyć. */
function stanyPlikow(stanRepo: DeveloperGitStatus | undefined): Map<string, GitStatusEntry> {
  const spis = new Map<string, GitStatusEntry>();
  for (const pozycja of stanRepo?.entries ?? []) {
    if (stanPozycji(pozycja) !== '') spis.set(pozycja.path, pozycja);
  }
  return spis;
}

/** Stan pozycji repozytorium: katalog roboczy przed indeksem, bo Operator widzi w drzewie zmianę jeszcze nieprzygotowaną. */
function stanPozycji(pozycja: GitStatusEntry): string {
  if (pozycja.worktree !== GitFileState.Unmodified) return pozycja.worktree;
  if (pozycja.index !== GitFileState.Unmodified) return pozycja.index;
  return '';
}

/** Dokłada do wpisu drzewa oznaczenie stanu pliku; nazwa stanu pochodzi z rejestru kontraktu, nie ze skrótu tego pliku. */
function oznaczStanPliku(
  wpis: HTMLElement,
  wzor: HTMLElement | null,
  pozycja: GitStatusEntry | undefined,
): void {
  if (wzor === null || pozycja === undefined) return;
  const plakietka = wzor.cloneNode(true) as HTMLElement;
  plakietka.textContent = stanPozycji(pozycja);
  wpis.appendChild(plakietka);
}

/** Wypełnia wykaz zmian panelu repozytorium pozycjami stanu i nanosi ich liczbę na nagłówek wykazu. */
function wypelnijZmiany(
  cialo: HTMLElement,
  wzory: Wzory,
  stanRepo: DeveloperGitStatus | undefined,
): void {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-git .dn-wykaz-modulu');
  if (wykaz === null) return;
  for (const stojaca of wykaz.querySelectorAll('.dn-wykaz-modulu-poz')) stojaca.remove();
  const pozycje = (stanRepo?.entries ?? []).filter((pozycja) => stanPozycji(pozycja) !== '');
  const naglowek = wykaz.querySelector('.pt-etykieta');
  if (naglowek !== null) {
    naglowek.textContent = (naglowek.textContent ?? '').replace(/\(\d+\)/u, `(${pozycje.length})`);
  }
  if (wzory.zmiana === null) return;
  for (const pozycja of pozycje) {
    const wiersz = wzory.zmiana.cloneNode(true) as HTMLElement;
    wiersz.dataset.sciezka = pozycja.path;
    if (wpiszTekst(wiersz, ` ${pozycja.path} `)) wykaz.appendChild(wiersz);
  }
}

/**
 * Wypełnia obszar różnicy wierszami odpowiedzi rdzenia. Wiersze kontekstu
 * odpadają, bo wzoru dla nich znacznik nie niesie, więc różnica idzie bez
 * otoczenia zmiany.
 */
async function wypelnijRoznice(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  sciezki: string[],
): Promise<void> {
  const numery = cialo.querySelector<HTMLElement>('#panel-git .pt-edytor-numery');
  const kod = cialo.querySelector<HTMLElement>('#panel-git .pt-edytor-kod');
  if (numery === null || kod === null) return;
  const wynik = await wywolaj(kanal, Command.DeveloperGitDiff, {
    windowId: idOkna,
    paths: sciezki.length === 0 ? undefined : sciezki,
    contextLines: 0,
  });
  numery.textContent = '';
  kod.replaceChildren();
  if (!wynik.udany || wynik.wynik === undefined) return;
  const numeracja: string[] = [];
  for (const fragment of wynik.wynik.hunks) {
    for (const linia of fragment.lines) {
      const wzor = linia.kind === GitDiffLineKind.Added ? wzory.wierszDodany : wzory.wierszUsuniety;
      if (linia.kind === GitDiffLineKind.Context || wzor === null) continue;
      const wiersz = wzor.cloneNode(true) as HTMLElement;
      wiersz.textContent = `${linia.text}\n`;
      kod.appendChild(wiersz);
      numeracja.push(String(linia.newLine ?? linia.oldLine ?? ''));
    }
  }
  numery.textContent = numeracja.join('\n');
}

/** Wiąże wybór pozycji drzewa z wczytaniem pliku do edytora; katalog wskazania nie niesie, bo komenda otwarcia bierze plik. */
function zwiazDrzewo(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): void {
  const drzewo = cialo.querySelector<HTMLElement>('#panel-drzewo .dv-drzewo');
  if (drzewo === null) return;
  drzewo.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wpis = cel.closest<HTMLElement>('.wpis');
    const sciezka = wpis?.dataset.sciezka;
    if (wpis === null || wpis === undefined || sciezka === undefined) return;
    for (const inny of drzewo.querySelectorAll('.wpis')) inny.removeAttribute('aria-current');
    wpis.setAttribute('aria-current', 'true');
    void otworzPlik(kanal, idOkna, cialo, wzory, stan, sciezka, wpis.dataset.nazwa ?? sciezka);
  });
}

/**
 * Wiąże edytor: karta pliku wczytuje go ponownie, a Ctrl+S odkłada treść
 * w repozytorium. Obszar kodu prototypu stoi zamknięty na edycję, bo niósł
 * plik przykładowy; z treścią rdzenia ma być polem pracy.
 */
function zwiazEdytor(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): void {
  const kod = cialo.querySelector<HTMLElement>('#panel-edytor .pt-edytor-kod');
  const karty = cialo.querySelector<HTMLElement>('#panel-edytor .dv-karty');
  karty?.replaceChildren();
  if (kod === null) return;
  kod.addEventListener('input', () => {
    ponumeruj(cialo, kod.textContent ?? '');
  });
  kod.addEventListener('keydown', (zdarzenie) => {
    if (!zdarzenie.ctrlKey || zdarzenie.key.toLowerCase() !== 's') return;
    zdarzenie.preventDefault();
    if (stan.plik === '') return;
    void wywolaj(kanal, Command.DeveloperFileSave, {
      windowId: idOkna,
      path: stan.plik,
      content: kod.textContent ?? '',
      createVersion: true,
    });
  });
  karty?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.dv-karta');
    const sciezka = karta?.dataset.sciezka;
    if (karta === null || karta === undefined || sciezka === undefined) return;
    void otworzPlik(kanal, idOkna, cialo, wzory, stan, sciezka, karta.dataset.nazwa ?? sciezka);
  });
}

/** Wczytuje plik repozytorium do edytora wraz z jego kartą, numeracją wierszy i nazwą języka. */
async function otworzPlik(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  sciezka: string,
  nazwa: string,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperFileOpen, { windowId: idOkna, path: sciezka });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const plik = wynik.wynik.file;
  stan.plik = plik.path;
  const kod = cialo.querySelector<HTMLElement>('#panel-edytor .pt-edytor-kod');
  if (kod !== null) {
    kod.textContent = plik.content ?? '';
    kod.setAttribute('contenteditable', 'true');
    ponumeruj(cialo, kod.textContent);
  }
  opiszJezyk(cialo, plik.language ?? '');
  postawKarte(cialo, wzory, plik.path, nazwa);
}

/** Nanosi na pas numeracji tyle liczb, ile wierszy niesie treść pliku. */
function ponumeruj(cialo: HTMLElement, tresc: string): void {
  const numery = cialo.querySelector<HTMLElement>('#panel-edytor .pt-edytor-numery');
  if (numery === null) return;
  const numeracja: string[] = [];
  const ile = tresc === '' ? 0 : tresc.split('\n').length;
  for (let numer = 1; numer <= ile; numer += 1) numeracja.push(String(numer));
  numery.textContent = numeracja.join('\n');
}

/**
 * Nanosi na pas stanu edytora nazwę języka pliku. Kodowanie, koniec wiersza,
 * miejsce kursora i znak edycji przez model schodzą — plik repozytorium niesie
 * ścieżkę, treść, język, rozmiar, wersję i czas zmiany, nic ponadto.
 */
function opiszJezyk(cialo: HTMLElement, jezyk: string): void {
  const pas = cialo.querySelector<HTMLElement>('#panel-edytor .dv-stan-linia');
  if (pas === null) return;
  pas.querySelector('.dn-plakietka')?.remove();
  const pola = [...pas.querySelectorAll('span')];
  for (const [numer, pole] of pola.entries()) {
    if (numer === 0) pole.textContent = jezyk;
    else pole.remove();
  }
}

/** Stawia kartę wczytanego pliku i czyni ją kartą wybraną; plik już otwarty dostaje wybór, nie drugą kartę. */
function postawKarte(cialo: HTMLElement, wzory: Wzory, sciezka: string, nazwa: string): void {
  const karty = cialo.querySelector<HTMLElement>('#panel-edytor .dv-karty');
  if (karty === null || wzory.karta === null) return;
  const stojace = [...karty.querySelectorAll<HTMLElement>('.dv-karta')];
  let karta = stojace.find((kandydat) => kandydat.dataset.sciezka === sciezka);
  if (karta === undefined) {
    karta = wzory.karta.cloneNode(true) as HTMLElement;
    karta.dataset.sciezka = sciezka;
    karta.dataset.nazwa = nazwa;
    if (!wpiszTekst(karta, nazwa)) return;
    karty.appendChild(karta);
  }
  for (const inna of karty.querySelectorAll('.dv-karta')) inna.setAttribute('aria-selected', 'false');
  karta.setAttribute('aria-selected', 'true');
}

/**
 * Wiąże czynności repozytorium: wybór pozycji zawęża różnicę do jednego pliku,
 * przycisk opisu woła operację kontekstową modelu, a zatwierdzenie i wysłanie
 * idą komendą czynności repozytorium.
 */
function zwiazCzynnosciRepozytorium(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-git');
  if (panel === null) return;
  const opis = panel.querySelector<HTMLElement>('.sta-chip');
  // Opis przykładowy schodzi przed pierwszym wywołaniem: pole puste jest
  // uczciwe, pole z cudzym opisem zmiany — nie.
  if (opis !== null) opis.textContent = '';

  panel.querySelector('.dn-wykaz-modulu')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const sciezka = cel.closest<HTMLElement>('.dn-wykaz-modulu-poz')?.dataset.sciezka;
    if (sciezka === undefined) return;
    void wypelnijRoznice(kanal, idOkna, cialo, zdejmijWzory(cialo), [sciezka]);
  });

  panel.querySelector('label.pt-etykieta button')?.addEventListener('click', () => {
    void opiszZmiane(kanal, idOkna, cialo, opis);
  });
  panel
    .querySelector('.dn-wykaz-modulu-poz button.dn-btn--atrament')
    ?.addEventListener('click', () => {
      void zatwierdz(kanal, idOkna, opis?.textContent ?? '', odswiez);
    });
  panel.querySelector('.dn-wykaz-modulu-poz button.dn-btn--zarys')?.addEventListener('click', () => {
    void wykonajCzynnosc(kanal, idOkna, GitActionKind.Push, undefined, odswiez);
  });
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.addEventListener('click', () => {
    void zatwierdz(kanal, idOkna, opis?.textContent ?? '', odswiez);
  });
}

/** Woła operację kontekstową modelu nad plikami zmienionymi i wpisuje jej wynik w pole opisu zmiany. */
async function opiszZmiane(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  opis: HTMLElement | null,
): Promise<void> {
  if (opis === null) return;
  const sciezki = [...cialo.querySelectorAll<HTMLElement>('#panel-git .dn-wykaz-modulu-poz')]
    .map((pozycja) => pozycja.dataset.sciezka ?? '')
    .filter((sciezka) => sciezka !== '');
  const wynik = await wywolaj(kanal, Command.DeveloperContextualOp, {
    windowId: idOkna,
    operation: ContextualOpKind.Generate,
    contextPaths: sciezki.length === 0 ? undefined : sciezki,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  opis.textContent = wynik.wynik.result;
}

/** Zatwierdza zmiany opisem z pola panelu; opis pusty wstrzymuje czynność, bo zatwierdzenie bez opisu nie mówi, co niesie. */
async function zatwierdz(
  kanal: Kanal,
  idOkna: string,
  opis: string,
  odswiez: () => void,
): Promise<void> {
  const wiadomosc = opis.trim();
  if (wiadomosc === '') return;
  await wykonajCzynnosc(kanal, idOkna, GitActionKind.Commit, wiadomosc, odswiez);
}

/** Wykonuje czynność repozytorium i odświeża stanowisko, gdy rdzeń ją przyjął. */
async function wykonajCzynnosc(
  kanal: Kanal,
  idOkna: string,
  czynnosc: GitActionKind,
  wiadomosc: string | undefined,
  odswiez: () => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperGitAction, {
    windowId: idOkna,
    action: czynnosc,
    message: wiadomosc,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  odswiez();
}

/** Wczytuje ostatni przebieg budowania okna wraz z jego logiem i licznikiem testów; brak przebiegu zdejmuje monitor i konsolę. */
async function wypelnijBudowanie(kanal: Kanal, idOkna: string, cialo: HTMLElement): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperBuildList, { windowId: idOkna, limit: 1 });
  const przebieg = wynik.udany ? wynik.wynik?.builds[0] : undefined;
  if (przebieg === undefined) {
    cialo.querySelector('.sta-kom-monitor')?.remove();
    cialo.querySelector('#panel-build .sta-okno-znacznik')?.remove();
    opiszLog(cialo, '');
    return;
  }
  opiszPrzebieg(cialo, przebieg);
  await wypelnijLog(kanal, cialo, przebieg.id);
  await opiszLicznikTestow(kanal, cialo, przebieg);
}

/** Nanosi zadanie przebiegu na znacznik panelu budowania i na monitor okna komunikacji. */
function opiszPrzebieg(cialo: HTMLElement, przebieg: DeveloperBuild): void {
  const znacznik = cialo.querySelector<HTMLElement>('#panel-build .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = przebieg.task;
  opiszMonitor(cialo, przebieg, '');
}

/** Nanosi na monitor zadanie przebiegu wraz z licznikiem testów; znak pracy zostaje wyłącznie przy przebiegu trwającym. */
function opiszMonitor(cialo: HTMLElement, przebieg: DeveloperBuild, licznik: string): void {
  const monitor = cialo.querySelector<HTMLElement>('.sta-kom-monitor');
  if (monitor === null) return;
  if (przebieg.status !== BuildStatus.Running) monitor.querySelector('.pt-tetno')?.remove();
  const opis = licznik === '' ? przebieg.task : `${przebieg.task} · ${licznik}`;
  if (!wpiszTekst(monitor, ` ${opis}`)) monitor.remove();
}

/** Rozwija log przebiegu i wpisuje jego wiersze w konsolę panelu budowania. */
async function wypelnijLog(kanal: Kanal, cialo: HTMLElement, idPrzebiegu: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperBuildLogGet, { buildId: idPrzebiegu });
  opiszLog(cialo, wynik.udany && wynik.wynik !== undefined ? wynik.wynik.lines.join('\n') : '');
}

/**
 * Wpisuje treść w konsolę panelu budowania. Znak zachęty powłoki i podział
 * wierszy na wagi schodzą razem z treścią przykładową — log rdzenia jest
 * wykazem łańcuchów, a wagę niesie wyłącznie zgłoszenie budowania.
 */
function opiszLog(cialo: HTMLElement, tresc: string): void {
  const konsola = cialo.querySelector<HTMLElement>('#panel-build .pt-konsola');
  if (konsola === null) return;
  konsola.textContent = tresc;
}

/** Dokłada do monitora licznik testów przebiegu; wynik zestawu bierze się po identyfikatorze przebiegu, nie po oknie. */
async function opiszLicznikTestow(
  kanal: Kanal,
  cialo: HTMLElement,
  przebieg: DeveloperBuild,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperTestResultGet, { buildId: przebieg.id });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const zestaw = wynik.wynik;
  const razem = zestaw.passed + zestaw.failed + zestaw.skipped;
  opiszMonitor(cialo, przebieg, `${zestaw.passed}/${razem}`);
}

/** Nasłuchuje zmian budowania okna: monitor bierze stan przebiegu, a konsola dopisuje wiersz logu przyrostem. */
function zwiazZdarzeniaBudowania(kanal: Kanal, idOkna: string, cialo: HTMLElement): void {
  kanal.naZdarzenie(EventType.DeveloperBuildChanged, (tresc) => {
    if (tresc.build.windowId !== idOkna) return;
    opiszPrzebieg(cialo, tresc.build);
    const konsola = cialo.querySelector<HTMLElement>('#panel-build .pt-konsola');
    if (konsola === null || tresc.logLine === undefined) return;
    konsola.textContent = `${konsola.textContent ?? ''}${tresc.logLine}\n`;
  });
}

/** Wpisuje wartość w pierwszy niepusty węzeł tekstowy, zostawiając ikonę i przyciski znacznika; fałsz znaczy węzeł bez miejsca na tekst. */
function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
