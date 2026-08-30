/**
 * Wiązanie wnętrza okna modułu Research z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * research.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  Command,
  EventType,
  type ResearchFinding,
  type ResearchReport,
  type ResearchSource,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  sklonuj,
  wpiszAlboZdejmij,
  wpiszTekst,
  wpiszTekstAlboZdejmij,
  zdejmijAlboPostaw,
} from './okno-modulu.ts';

/**
 * Panele prototypu, którym rodzina research.* nie oddaje ani jednego pola;
 * schodzą razem z przełącznikami menu sesji. Panel eksportu schodzi z nimi:
 * niesie wybór formatu docelowego, a kontrakt zna wyłącznie ślady eksportów
 * już wykonanych, więc wykaz formatów pokazywałby wybór, którego rdzeń nie zna.
 */
const PANELE_BEZ_POKRYCIA = [
  'panel-plan',
  'panel-zadania',
  'panel-pliki',
  'panel-kolejka',
  'panel-subagenci',
  'panel-przegladarka',
  'panel-export',
];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  etap: HTMLElement | null;
  zrodlo: HTMLElement | null;
  ustalenie: HTMLElement | null;
  sekcja: HTMLElement | null;
}

/**
 * Węzły niosące jedną wartość, zdejmowane na czas jej braku. Wskazania stoją
 * tu, a nie w wyszukaniu przy każdym odczycie: węzeł zdjęty jest poza
 * dokumentem, więc wyszukanie już by go nie znalazło i wartość, która dojdzie,
 * nie miałaby czego postawić.
 */
interface Wezly {
  tytul: Element | null;
  chip: Element | null;
  belka: Element | null;
  sugestia: Element | null;
}

/**
 * Wiąże wnętrze okna modułu Research. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazResearch(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);
  const wezly = zbierzWezly(cialo);

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt przestrzeni, źródeł, ustaleń i raportu; wołane przy wejściu i po każdym zdarzeniu modułu. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, wezly);
  };

  zwiazNawigacjePaneli(cialo);
  zwiazSzukanieZrodel(kanal, idOkna, cialo, wzory);
  zwiazZdarzeniaBadania(kanal, idOkna, odswiez);

  odswiez();
}

/** Zdejmuje wzory wierszy i pozycji, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const plakietki = [...cialo.querySelectorAll<HTMLElement>('.rs-etapy .dn-plakietka')];
  // Etap przestrzeni jest samą nazwą: stopnia zaawansowania kontrakt przy nim
  // nie niesie, więc wzorem jest plakietka bez odmiany barwnej.
  const etap = sklonuj(plakietki.find((wzor) => wzor.className.trim() === 'dn-plakietka') ?? null);

  const zrodlo = sklonuj(cialo.querySelector('#panel-sources .dn-wykaz-modulu-poz'));
  // Numer cytowania i barwa oceny nie mają pola: źródło niesie ocenę słowem,
  // a nie kropką, i nie niesie żadnego numeru w wykazie.
  zrodlo?.querySelector('.dn-plakietka')?.remove();
  zrodlo?.querySelector('.dn-kropka')?.remove();

  const ustalenia = [...cialo.querySelectorAll<HTMLElement>('#panel-findings .dn-wykaz-modulu-poz')];
  const ustalenie = sklonuj(ustalenia[1] ?? ustalenia[0] ?? null);
  // Kolejny numer ustalenia i odsyłacze do źródeł numerem nie mają pola;
  // zostaje miejsce na stan ustalenia.
  ustalenie?.querySelector('.dn-plakietka')?.remove();

  const sekcja = sklonuj(cialo.querySelector('#panel-report .dn-wykaz-modulu-poz'));
  // Znak postępu sekcji nie ma pola: raport niesie tytuł, treść i kolejność.
  sekcja?.querySelector('.dn-kropka')?.remove();
  sekcja?.querySelector('.dn-meta')?.remove();

  return { etap, zrodlo, ustalenie, sekcja };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz badań szyny, historia rozmowy
 * i przypięte źródła wiszą na rodzinach spoza research.*, a miara postępu
 * raportu, mapa powiązań i widoki badania nie mają pola w całym kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  for (const chip of cialo.querySelectorAll('.sta-kom-kontekst > .sta-zrodlo')) chip.remove();
  oznaczChipyCzynnosci(cialo);
  zdejmijWidokiBadania(cialo);
  zdejmijMiaryBezPokrycia(cialo);
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania i miara różnicy schodzą,
 * migawka badania schodzi razem z nimi, a znaki zakresu, źródeł i ustaleń
 * zostają oznaczone do wypełnienia.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  chipy[0]?.remove();
  const pola = ['zakres', 'zrodla', 'ustalenia'];
  for (const [numer, pole] of pola.entries()) {
    const chip = chipy[numer + 1];
    if (chip !== undefined) chip.dataset.pole = pole;
  }
  cialo.querySelector('.sta-kontekst-akcji .sta-chip--diff')?.remove();
  // Migawka badania nie jest czynnością kontraktu: wersjonowanie zna wyłącznie
  // raport, a nie całe badanie.
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
}

/**
 * Zdejmuje przełącznik widoków badania i płótno mapy powiązań. Kontrakt niesie
 * graf dowodów jako wykaz węzłów i krawędzi, a nie jako oś czasu, tablicę ani
 * rysunek, więc żaden z czterech widoków nie ma czym stanąć.
 */
function zdejmijWidokiBadania(cialo: HTMLElement): void {
  cialo.querySelector('#panel-workspace .rs-zakladki')?.remove();
  cialo.querySelector('#panel-workspace .rs-graf')?.remove();
  // Zakres badania zmienia się komendą przestrzeni, ale okna edycji zakresu
  // prototyp nie niesie, więc przycisk odsyłałby donikąd.
  cialo.querySelector('#panel-workspace .sta-okno-tresc button')?.remove();
}

/** Zdejmuje miary, dla których kontrakt nie ma pola: stopień skompletowania raportu, styl cytowania i podpisy sortowania. */
function zdejmijMiaryBezPokrycia(cialo: HTMLElement): void {
  const liczniki = [...cialo.querySelectorAll('#panel-workspace .rs-licznik')];
  liczniki[2]?.remove();
  cialo.querySelector('#panel-report .sta-okno-znacznik')?.remove();
  // Styl cytowania bieżącego badania nie ma pola: rejestr stylów jest wykazem,
  // a nie wskazaniem stylu wybranego.
  cialo.querySelector('#panel-sources .sta-chip')?.remove();
  for (const meta of cialo.querySelectorAll('#panel-sources .dn-wykaz-modulu > .dn-meta')) {
    meta.remove();
  }
  for (const meta of cialo.querySelectorAll('#panel-findings .dn-wykaz-modulu > .dn-meta')) {
    meta.remove();
  }
}

/**
 * Zdejmuje sterowanie, którego rodzina research.* nie obsługuje: panele bez
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

/** Zostawia w menu dodawania sam tytuł i trzy pozycje plików; konektory i wtyczki należą do rodzin spoza research.*. */
function zdejmijKonektory(cialo: HTMLElement): void {
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Pyta rdzeń o przestrzeń badania, źródła, ustalenia, raport i luki, po czym wypełnia nimi belki, znaki i wykazy paneli. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
): Promise<void> {
  const przestrzen = await wywolaj(kanal, Command.ResearchWorkspaceGet, {});
  const zrodla = await wywolaj(kanal, Command.ResearchSourceList, { windowId: idOkna });
  const ustalenia = await wywolaj(kanal, Command.ResearchFindingList, { windowId: idOkna });
  const raport = await wywolaj(kanal, Command.ResearchReportGet, { windowId: idOkna });

  const zakres = przestrzen.udany ? (przestrzen.wynik?.scope ?? '') : '';
  opiszZakres(wezly, zakres);
  wypelnijEtapy(cialo, wzory, przestrzen.udany ? (przestrzen.wynik?.stages ?? []) : []);

  const wykazZrodel = zrodla.udany ? (zrodla.wynik?.sources ?? []) : [];
  const wykazUstalen = ustalenia.udany ? (ustalenia.wynik?.findings ?? []) : [];
  opiszMiare(cialo, 'zrodla', '#panel-sources', zrodla.udany ? (zrodla.wynik?.total ?? 0) : 0, 0);
  opiszMiare(
    cialo,
    'ustalenia',
    '#panel-findings',
    ustalenia.udany ? (ustalenia.wynik?.total ?? 0) : 0,
    1,
  );
  wypelnijZrodla(cialo, wzory, wykazZrodel);
  wypelnijUstalenia(cialo, wzory, wykazUstalen);
  wypelnijRaport(cialo, wzory, raport.udany ? raport.wynik?.report : undefined);
  await opiszLuke(kanal, idOkna, wezly);
}

/** Nanosi zakres badania na tytuł okna wiodącego, na znak belki rozmowy i na znak pasa czynności; zakres pusty zdejmuje każdy z nich. */
function opiszZakres(wezly: Wezly, zakres: string): void {
  wpiszAlboZdejmij(wezly.tytul, zakres);
  wpiszTekstAlboZdejmij(wezly.chip, zakres);
  wpiszTekstAlboZdejmij(wezly.belka, zakres);
}

/** Zbiera węzły jednowartościowe po zdjęciu treści przykładowej; oznaczenia pól nadaje dopiero to zdjęcie. */
function zbierzWezly(cialo: HTMLElement): Wezly {
  return {
    tytul: cialo.querySelector('#panel-workspace .sta-okno-tresc b'),
    chip: cialo.querySelector('.sta-kontekst-akcji [data-pole="zakres"]'),
    belka: cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'),
    sugestia: cialo.querySelector('#panel-workspace .rs-sugestia'),
  };
}

/** Stawia w pasie etapów po jednej plakietce na etap przestrzeni; kolejność bierze się z odpowiedzi rdzenia. */
function wypelnijEtapy(cialo: HTMLElement, wzory: Wzory, etapy: string[]): void {
  const pas = cialo.querySelector<HTMLElement>('#panel-workspace .rs-etapy');
  if (pas === null) return;
  pas.replaceChildren();
  if (wzory.etap === null) return;
  for (const etap of etapy) {
    const plakietka = wzory.etap.cloneNode(true) as HTMLElement;
    plakietka.textContent = etap;
    pas.appendChild(plakietka);
  }
}

/** Nanosi liczbę pozycji na licznik okna wiodącego, na znak panelu i na znak pasa czynności; podmieniana jest sama liczba, podpis zostaje. */
function opiszMiare(
  cialo: HTMLElement,
  pole: string,
  panel: string,
  liczba: number,
  numerLicznika: number,
): void {
  const licznik = [...cialo.querySelectorAll<HTMLElement>('#panel-workspace .rs-licznik b')][
    numerLicznika
  ];
  if (licznik !== undefined) licznik.textContent = String(liczba);
  const znacznik = cialo.querySelector<HTMLElement>(`${panel} .sta-okno-znacznik`);
  if (znacznik !== null) znacznik.textContent = String(liczba);
  const chip = cialo.querySelector<HTMLElement>(`.sta-kontekst-akcji [data-pole="${pole}"]`);
  if (chip === null) return;
  const tekst = (chip.textContent ?? '').replace(/^\s*\d+/u, String(liczba));
  if (!wpiszTekst(chip, tekst)) chip.remove();
}

/** Wypełnia wykaz źródeł pozycjami odpowiedzi; ocena wiarygodności wchodzi w podpis pozycji, bo tylko ona ma pole. */
function wypelnijZrodla(cialo: HTMLElement, wzory: Wzory, zrodla: ResearchSource[]): void {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-sources .dn-wykaz-modulu');
  if (wykaz === null || wzory.zrodlo === null) return;
  for (const stojaca of wykaz.querySelectorAll('.dn-wykaz-modulu-poz')) stojaca.remove();
  const kotwica = koniecWykazu(wykaz);
  for (const zrodlo of zrodla) {
    const wiersz = wzory.zrodlo.cloneNode(true) as HTMLElement;
    wiersz.dataset.zrodlo = zrodlo.id;
    const meta = wiersz.querySelector('.dn-meta');
    if (meta !== null) meta.textContent = zrodlo.credibility;
    if (wpiszTekst(wiersz, ` ${zrodlo.title} `)) wykaz.insertBefore(wiersz, kotwica);
  }
}

/** Wypełnia wykaz ustaleń treścią ustaleń; stan ustalenia wchodzi w podpis pozycji, a odsyłacze do źródeł schodzą razem z numeracją. */
function wypelnijUstalenia(cialo: HTMLElement, wzory: Wzory, ustalenia: ResearchFinding[]): void {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-findings .dn-wykaz-modulu');
  if (wykaz === null || wzory.ustalenie === null) return;
  for (const stojaca of wykaz.querySelectorAll('.dn-wykaz-modulu-poz')) stojaca.remove();
  const kotwica = koniecWykazu(wykaz);
  for (const ustalenie of ustalenia) {
    const wiersz = wzory.ustalenie.cloneNode(true) as HTMLElement;
    wiersz.dataset.ustalenie = ustalenie.id;
    const meta = wiersz.querySelector('.dn-meta');
    if (meta !== null) meta.textContent = ustalenie.status;
    if (wpiszTekst(wiersz, ` ${ustalenie.content} `)) wykaz.insertBefore(wiersz, kotwica);
  }
}

/** Wypełnia strukturę raportu tytułami jego sekcji; raport niezłożony zostawia wykaz pusty, bo sekcji nie ma skąd wziąć. */
function wypelnijRaport(
  cialo: HTMLElement,
  wzory: Wzory,
  raport: ResearchReport | undefined,
): void {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-report .dn-wykaz-modulu');
  if (wykaz === null || wzory.sekcja === null) return;
  for (const stojaca of wykaz.querySelectorAll('.dn-wykaz-modulu-poz')) stojaca.remove();
  const kotwica = koniecWykazu(wykaz);
  for (const sekcja of raport?.sections ?? []) {
    const wiersz = wzory.sekcja.cloneNode(true) as HTMLElement;
    wiersz.dataset.sekcja = sekcja.id;
    if (wpiszTekst(wiersz, ` ${sekcja.title} `)) wykaz.insertBefore(wiersz, kotwica);
  }
}

/** Przycisk zamykający wykaz zostaje na końcu panelu; pozycje wchodzą przed nim, żeby czynność nie stała w środku wykazu. */
function koniecWykazu(wykaz: HTMLElement): Element | null {
  const ostatni = wykaz.lastElementChild;
  return ostatni !== null && ostatni.tagName === 'BUTTON' ? ostatni : null;
}

/**
 * Nanosi na pas sugestii pierwszą lukę wskazaną przez rdzeń. Znak pracy
 * i przycisk zbadania schodzą: praca w tej chwili nie trwa, a wyszukanie
 * źródeł sięga sieci, której okno badania w tym wydaniu nie otwiera.
 */
async function opiszLuke(kanal: Kanal, idOkna: string, wezly: Wezly): Promise<void> {
  const pas = wezly.sugestia;
  if (pas === null) return;
  pas.querySelector('.pt-tetno')?.remove();
  pas.querySelector('button')?.remove();
  const wynik = await wywolaj(kanal, Command.ResearchGapFind, { windowId: idOkna });
  const luka = wynik.udany ? wynik.wynik?.gaps[0] : undefined;
  const opis = pas.querySelector('span');
  if (opis !== null && luka !== undefined) opis.textContent = luka.summary;
  zdejmijAlboPostaw(pas, luka !== undefined && opis !== null);
}

/** Wiąże odsyłacze okna wiodącego z panelami, do których prowadzą; odsyłacz do panelu zdjętego schodzi razem z nim. */
function zwiazNawigacjePaneli(cialo: HTMLElement): void {
  for (const przycisk of cialo.querySelectorAll<HTMLElement>('[data-panel-cel]')) {
    const panel = cialo.querySelector(`#${przycisk.dataset.panelCel ?? ''}`);
    if (panel === null) {
      przycisk.remove();
      continue;
    }
    przycisk.addEventListener('click', () => {
      panel.scrollIntoView({ block: 'nearest' });
    });
  }
}

/** Wiąże pole zawężania wykazu źródeł z komendą wykazu; fraza pusta przywraca wykaz pełny. */
function zwiazSzukanieZrodel(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): void {
  const pole = cialo.querySelector<HTMLInputElement>('#panel-sources input[type="search"]');
  if (pole === null) return;
  pole.addEventListener('input', () => {
    const fraza = pole.value.trim();
    void wywolaj(kanal, Command.ResearchSourceList, {
      windowId: idOkna,
      query: fraza === '' ? undefined : fraza,
    }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) return;
      wypelnijZrodla(cialo, wzory, wynik.wynik.sources);
    });
  });
}

/** Nasłuchuje zmian badania: nowe źródło, ustalenie albo raport odświeżają całe stanowisko. */
function zwiazZdarzeniaBadania(kanal: Kanal, idOkna: string, odswiez: () => void): void {
  kanal.naZdarzenie(EventType.ResearchSourceChanged, (tresc) => {
    if (tresc.source.windowId === idOkna) odswiez();
  });
  kanal.naZdarzenie(EventType.ResearchFindingChanged, (tresc) => {
    if (tresc.finding.windowId === idOkna) odswiez();
  });
  kanal.naZdarzenie(EventType.ResearchReportChanged, (tresc) => {
    if (tresc.report.windowId === idOkna) odswiez();
  });
}
