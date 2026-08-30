/**
 * Wiązanie wnętrza okna modułu Agents z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * agent.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  AgentVisibility,
  Command,
  EventType,
  MemoryLevel,
  ProviderTransport,
  type Agent,
  type AgentConnector,
  type AgentPlugin,
  type IsolationTechnicalSwitch,
  type Module,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  sklonuj,
  wpiszTekst,
  wpiszTekstAlboZdejmij,
  zdejmijAlboPostaw,
} from './okno-modulu.ts';

/** Panele prototypu, którym rodzina agent.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-plan', 'panel-zadania', 'panel-pliki', 'panel-artefakty'];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  pozycja: HTMLElement | null;
  karta: HTMLElement | null;
  plakietka: HTMLElement | null;
  opcja: HTMLElement | null;
  wybor: HTMLElement | null;
  rozszerzenie: HTMLElement | null;
  umiejetnosc: HTMLElement | null;
  izolacja: HTMLElement | null;
  model: HTMLElement | null;
  transport: HTMLElement | null;
  wersja: HTMLElement | null;
}

/**
 * Węzły niosące jedną wartość, zdejmowane na czas jej braku. Wskazania stoją
 * tu, a nie w wyszukaniu przy każdym odczycie: węzeł zdjęty jest poza
 * dokumentem, więc wyszukanie już by go nie znalazło i wartość, która dojdzie,
 * nie miałaby czego postawić.
 */
interface Wezly {
  belka: Element | null;
  podsumowanie: Element[];
}

/** Ekspert wczytany do edytora i rejestr modułów rdzenia; oba pochodzą z odpowiedzi rdzenia. */
interface Stan {
  wybrany: Agent | undefined;
  moduly: Module[];
}

/**
 * Wiąże wnętrze okna modułu Agents. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazAgents(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const wezly: Wezly = {
    belka: cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'),
    podsumowanie: [...cialo.querySelectorAll('.ab-sum-w')],
  };
  const stan: Stan = { wybrany: undefined, moduly: [] };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt biblioteki ekspertów; wołane przy wejściu i po każdej zmianie eksperta. */
  const odswiez = (): void => {
    void wypelnijBiblioteke(kanal, idOkna, cialo, wzory, wezly, stan);
  };

  zwiazWidoki(cialo);
  zwiazGrupyZakresu(cialo);
  zwiazWybor(kanal, idOkna, cialo, wzory, wezly, stan);
  zwiazZapis(kanal, cialo, stan, odswiez);
  kanal.naZdarzenie(EventType.AgentChanged, () => {
    odswiez();
  });

  odswiez();
}

/** Zdejmuje wzory pozycji, kart i wierszy, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const pozycja = sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .pt-pozycja'));
  // Znak pracy i wskazanie pozycji bieżącej nadaje dopiero wybór Operatora.
  pozycja?.querySelector('.pt-tetno')?.remove();
  pozycja?.removeAttribute('aria-current');

  const karta = sklonuj(cialo.querySelector('#widok-biblioteka .ab-karta'));
  // Powielenie definicji nie jest komendą kontraktu, a edycja idzie kliknięciem
  // w samą kartę, więc pas czynności karty schodzi.
  karta?.querySelector('.ab-karta-akcje')?.remove();

  const rozszerzenie = sklonuj(cialo.querySelector('#panel-connectors .pk-wiersz'));
  const umiejetnosc = sklonuj(cialo.querySelector('#panel-skills .pk-wiersz'));
  // Umiejętność niesie sam identyfikator przypisania: rejestru umiejętności
  // wraz z ich nazwą i stanem kontrakt nie ma, więc podpis wiersza schodzi.
  umiejetnosc?.querySelector('small')?.remove();

  const izolacja = sklonuj(cialo.querySelector('#grupa-izolacja .pk-mrz'));
  // Znak zakresu izolacji stoi w prototypie osobny dla każdego z ośmiu
  // wierszy; jeden wzór dałby wszystkim ten sam, więc znak schodzi.
  izolacja?.querySelector('svg')?.remove();

  const model = sklonuj(cialo.querySelector('#panel-model .pk-model'));
  // Podpis modelu opisuje jego przeznaczenie; rejestr kanałów niesie nazwę
  // i identyfikator modelu, nie jego przeznaczenie.
  model?.querySelector('small')?.remove();

  const wersja = sklonuj(cialo.querySelector('.ab-historia .dn-wykaz-modulu-poz'));
  wersja?.querySelector('.dn-kropka')?.remove();

  const grupy = [...cialo.querySelectorAll('#tab-tozsamosc .ab-grupa')];
  return {
    pozycja,
    karta,
    plakietka: sklonuj(cialo.querySelector('#widok-biblioteka .ab-karta-plak .dn-plakietka')),
    opcja: sklonuj(grupy[0]?.querySelector('.ab-opcja') ?? null),
    wybor: sklonuj(grupy[1]?.querySelector('.ab-opcja') ?? null),
    rozszerzenie,
    umiejetnosc,
    izolacja,
    model,
    transport: sklonuj(cialo.querySelector('#panel-model .pk-pig button')),
    wersja,
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Historia rozmowy, przypięty kontekst
 * i monitor wiszą na rodzinach spoza agent.*, a podpisy grupowania, znak
 * zapisu i polityka wypisana słowem nie mają pola w całym kontrakcie.
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
  zdejmijPolaBezPokrycia(cialo);
}

/** Zdejmuje znaki pasa czynności; ani projekt, ani ścieżka definicji, ani drzewo robocze, ani przyrost wersji nie mają pola przy ekspercie. */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  for (const chip of cialo.querySelectorAll('.sta-kontekst-akcji > .sta-chip')) chip.remove();
}

/**
 * Zdejmuje węzły, dla których rodzina agent.* nie ma pola: znak czynności
 * eksperta, czas zapisu, tabelę zakresu rozszerzeń wraz z jej zakładką,
 * politykę wypisaną słowem i liczniki zakresów.
 */
function zdejmijPolaBezPokrycia(cialo: HTMLElement): void {
  cialo.querySelector('.ab-edy-naglowek .dn-plakietka')?.remove();
  cialo.querySelector('.ab-zapis')?.remove();
  // Zakres uprawnień wypisany jednym słowem nie ma pola: polityka efektywna
  // niesie wykaz zgód i przełączników, a nie ich streszczenie.
  [...cialo.querySelectorAll('.ab-sum-w')][4]?.remove();
  cialo.querySelector('.ab-historia button')?.remove();
  // Zakres rozszerzenia rozpisany na odczyt, zapis i wywołanie nie ma pola:
  // uprawnienie eksperta niesie jedną zgodę na całą grupę zakresu.
  cialo.querySelector('#grupa-rozsz')?.remove();
  cialo.querySelector('[data-grupa="grupa-rozsz"]')?.remove();
  cialo.querySelector('[data-reset-upr]')?.remove();
  for (const licznik of cialo.querySelectorAll('#grupa-izolacja p, #panel-permissions .pk-kod')) {
    licznik.remove();
  }
  const etykiety = [...cialo.querySelectorAll('#panel-permissions .pk-etyk')];
  etykiety[etykiety.length - 1]?.remove();
  // Wybór widoku startowego, liczby kolumn i autozapisu jest nastawą okna,
  // której rejestr nastaw nie zna.
  cialo.querySelector('#konfig-builder')?.remove();
  cialo.querySelector('[data-konfig-okno="konfig-builder"]')?.remove();
  zdejmijDaneDostepowe(cialo);
}

/**
 * Zdejmuje z konfiguracji modelu opis drogi wywołania, dane dostępowe kanału
 * i próbę połączenia. Rejestr kanałów oddaje odwołanie do poświadczenia, nigdy
 * jego treść, a próba połączenia sięga dostawcy poza maszyną.
 */
function zdejmijDaneDostepowe(cialo: HTMLElement): void {
  const panel = cialo.querySelector('#panel-model');
  if (panel === null) return;
  panel.querySelector('#kanal-opis')?.remove();
  const etykiety = [...panel.querySelectorAll('.pk-etyk')];
  etykiety[2]?.remove();
  panel.querySelector('.pk-pole')?.remove();
  panel.querySelector('[data-testuj]')?.closest('.dn-wykaz-modulu-poz')?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina agent.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, wybór modelu i nakładu rozmowy,
 * urządzenia wejścia dźwięku oraz konektory i wtyczki menu dodawania.
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
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
  // Rejestr umiejętności stoi poza rodziną agent.*, więc instalacji nie ma
  // czym wywołać.
  cialo.querySelector('#panel-skills button.dn-btn')?.remove();
  cialo.querySelector('#panel-connectors > .sta-okno-tresc > button.dn-btn')?.remove();
}

/** Wiąże przełączniki widoku i zakładek edytora; oba stoją w znaczniku, a przełączania nie niesie żaden skrypt wydania. */
function zwiazWidoki(cialo: HTMLElement): void {
  const przelacz = (tory: Element, wybrany: HTMLElement, cel: string, klasa: string): void => {
    for (const inny of tory.querySelectorAll('[role="tab"]')) {
      inny.setAttribute('aria-selected', String(inny === wybrany));
    }
    for (const widok of cialo.querySelectorAll(`.${klasa}`)) {
      widok.toggleAttribute('hidden', widok.id !== cel);
    }
  };
  for (const przycisk of cialo.querySelectorAll<HTMLElement>('[data-widok]')) {
    const tory = przycisk.parentElement;
    if (tory === null) continue;
    przycisk.addEventListener('click', () => {
      przelacz(tory, przycisk, przycisk.dataset.widok ?? '', 'ab-widok');
    });
  }
  for (const przycisk of cialo.querySelectorAll<HTMLElement>('[data-tab]')) {
    const tory = przycisk.parentElement;
    if (tory === null) continue;
    przycisk.addEventListener('click', () => {
      przelacz(tory, przycisk, przycisk.dataset.tab ?? '', 'ab-tab');
    });
  }
}

/** Wiąże wybór grupy zakresu i stawia na wierzchu pierwszą grupę, która została po zdjęciu grupy bez pokrycia. */
function zwiazGrupyZakresu(cialo: HTMLElement): void {
  const zakladki = [...cialo.querySelectorAll<HTMLElement>('[data-grupa]')];
  const pokaz = (wybrana: HTMLElement): void => {
    for (const zakladka of zakladki) {
      zakladka.setAttribute('aria-selected', String(zakladka === wybrana));
    }
    for (const grupa of cialo.querySelectorAll('.pk-grz')) {
      grupa.toggleAttribute('hidden', grupa.id !== wybrana.dataset.grupa);
    }
  };
  for (const zakladka of zakladki) {
    zakladka.addEventListener('click', () => {
      pokaz(zakladka);
    });
  }
  const pierwsza = zakladki[0];
  if (pierwsza !== undefined) pokaz(pierwsza);
}

/** Wiąże wybór eksperta pozycją szyny i kartą biblioteki; oba wskazania wczytują tę samą definicję do edytora. */
function zwiazWybor(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
  stan: Stan,
): void {
  cialo.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wskazanie = cel.closest<HTMLElement>('[data-ekspert]');
    const identyfikator = wskazanie?.dataset.ekspert;
    if (identyfikator === undefined) return;
    void pokazEksperta(kanal, idOkna, cialo, wzory, wezly, stan, identyfikator);
  });
}

/** Wiąże zapis definicji eksperta: tożsamość idzie komendą zmiany, moduły zastosowania osobną komendą, bo kontrakt trzyma je rozdzielnie. */
function zwiazZapis(kanal: Kanal, cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  for (const przycisk of cialo.querySelectorAll('[data-zapisz]')) {
    przycisk.addEventListener('click', () => {
      void zapiszEksperta(kanal, cialo, stan, odswiez);
    });
  }
}

/** Odkłada w rdzeniu tożsamość eksperta wraz z zasięgiem, pamięcią i modułami zastosowania. */
async function zapiszEksperta(
  kanal: Kanal,
  cialo: HTMLElement,
  stan: Stan,
  odswiez: () => void,
): Promise<void> {
  const ekspert = stan.wybrany;
  if (ekspert === undefined) return;
  const nazwa = wartoscPola(cialo, '#ab-nazwa');
  if (nazwa === '') return;
  const wynik = await wywolaj(kanal, Command.AgentUpdate, {
    agentId: ekspert.id,
    name: nazwa,
    description: wartoscPola(cialo, '#ab-opis'),
    systemPrompt: wartoscPola(cialo, '#ab-instr'),
    visibility: zaznaczonyZasieg(cialo),
    memoryLevels: zaznaczonePoziomy(cialo),
  });
  if (!wynik.udany) return;
  await wywolaj(kanal, Command.AgentModulesSet, {
    agentId: ekspert.id,
    moduleCodes: zaznaczoneModuly(cialo),
  });
  odswiez();
}

/** Treść pola formularza edytora; brak pola daje łańcuch pusty, bo komenda zmiany traktuje puste jako brak zmiany. */
function wartoscPola(cialo: HTMLElement, wybor: string): string {
  const pole = cialo.querySelector<HTMLInputElement | HTMLTextAreaElement>(wybor);
  return pole === null ? '' : pole.value;
}

/** Zasięg widoczności zaznaczony w edytorze; brak zaznaczenia bierze stan wyjściowy platformy. */
function zaznaczonyZasieg(cialo: HTMLElement): AgentVisibility {
  const wybrany = cialo.querySelector<HTMLElement>('[data-zasieg] input:checked');
  const wartosc = wybrany?.closest<HTMLElement>('[data-zasieg]')?.dataset.zasieg ?? '';
  return wartosc === AgentVisibility.Project ? AgentVisibility.Project : AgentVisibility.Global;
}

/** Poziomy pamięci zaznaczone w edytorze; zbiór pusty jest zapisem pamięci wyłączonej, nie brakiem odpowiedzi. */
function zaznaczonePoziomy(cialo: HTMLElement): MemoryLevel[] {
  const poziomy: MemoryLevel[] = [];
  for (const opcja of cialo.querySelectorAll<HTMLElement>('[data-poziom]')) {
    const pole = opcja.querySelector('input');
    if (pole?.checked !== true) continue;
    const wartosc = Object.values(MemoryLevel).find((kandydat) => kandydat === opcja.dataset.poziom);
    if (wartosc !== undefined) poziomy.push(wartosc);
  }
  return poziomy;
}

/** Kody modułów zaznaczone w edytorze; zbiór pusty znaczy brak ograniczenia, czyli dostępność we wszystkich modułach. */
function zaznaczoneModuly(cialo: HTMLElement): string[] {
  const kody: string[] = [];
  const opcje = [...cialo.querySelectorAll<HTMLElement>('#tab-tozsamosc [data-modul-kod]')];
  const wszystkie = opcje.length;
  for (const opcja of opcje) {
    if (opcja.querySelector('input')?.checked === true) kody.push(opcja.dataset.modulKod ?? '');
  }
  return kody.length === wszystkie ? [] : kody;
}

/** Pyta rdzeń o bibliotekę ekspertów i rejestr modułów, po czym wypełnia nimi szynę, kafle biblioteki i edytor. */
async function wypelnijBiblioteke(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
  stan: Stan,
): Promise<void> {
  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  stan.moduly = moduly.udany ? (moduly.wynik?.modules ?? []) : [];

  const eksperci = await wywolaj(kanal, Command.AgentList, {});
  const wykaz = eksperci.udany ? (eksperci.wynik?.agents ?? []) : [];
  const przypisania = await wywolaj(kanal, Command.AgentAssignmentList, {});
  const liczba = new Map<string, number>();
  for (const przypisanie of przypisania.udany ? (przypisania.wynik?.assignments ?? []) : []) {
    liczba.set(przypisanie.agentId, (liczba.get(przypisanie.agentId) ?? 0) + 1);
  }

  wypelnijSzyne(cialo, wzory, wykaz);
  wypelnijKarty(cialo, wzory, wykaz, liczba);
  const biezacy = stan.wybrany;
  const wskazany = wykaz.find((ekspert) => ekspert.id === biezacy?.id) ?? wykaz[0];
  await pokazEksperta(kanal, idOkna, cialo, wzory, wezly, stan, wskazany?.id ?? '');
}

/** Stawia w szynie po jednej pozycji na eksperta biblioteki; podpisy grupowania schodzą, bo grupy kontrakt przy ekspercie nie niesie. */
function wypelnijSzyne(cialo: HTMLElement, wzory: Wzory, eksperci: Agent[]): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null || wzory.pozycja === null) return;
  lista.replaceChildren();
  for (const ekspert of eksperci) {
    const pozycja = wzory.pozycja.cloneNode(true) as HTMLElement;
    pozycja.dataset.ekspert = ekspert.id;
    const tytul = pozycja.querySelector('.pt-pozycja-tytul');
    if (tytul !== null) tytul.textContent = ekspert.name;
    lista.appendChild(pozycja);
  }
}

/** Stawia w bibliotece po jednej karcie na eksperta; wersja i zasięg wchodzą w plakietki, a liczba przypisań w podpis karty. */
function wypelnijKarty(
  cialo: HTMLElement,
  wzory: Wzory,
  eksperci: Agent[],
  przypisania: Map<string, number>,
): void {
  const siatka = cialo.querySelector<HTMLElement>('#widok-biblioteka .ab-siatka');
  if (siatka === null || wzory.karta === null) return;
  siatka.replaceChildren();
  for (const ekspert of eksperci) {
    const karta = wzory.karta.cloneNode(true) as HTMLElement;
    karta.dataset.ekspert = ekspert.id;
    const nazwa = karta.querySelector('.ab-karta-nazwa');
    if (nazwa !== null) nazwa.textContent = ekspert.name;
    const opis = karta.querySelector('.ab-karta-opis');
    if (opis !== null) {
      if (ekspert.description === undefined) opis.remove();
      else opis.textContent = ekspert.description;
    }
    const podpisy = [...karta.querySelectorAll('.ab-karta-meta')];
    opiszKanalKarty(podpisy[0] ?? null, ekspert);
    opiszPrzypisania(podpisy[1] ?? null, przypisania.get(ekspert.id) ?? 0);
    opiszPlakietki(karta, wzory, ekspert);
    siatka.appendChild(karta);
  }
}

/** Nanosi na podpis karty model bazowy i drogę wywołania kanału; brak obu zdejmuje podpis, bo pusty niczego nie mówi. */
function opiszKanalKarty(podpis: Element | null, ekspert: Agent): void {
  if (podpis === null) return;
  const czesci = [ekspert.model ?? '', ekspert.transport ?? ''].filter((czesc) => czesc !== '');
  if (czesci.length === 0) podpis.remove();
  else podpis.textContent = czesci.join(' · ');
}

/** Nanosi na podpis karty liczbę przypisań eksperta; podmieniana jest sama liczba, podpis prototypu zostaje. */
function opiszPrzypisania(podpis: Element | null, liczba: number): void {
  if (podpis === null) return;
  podpis.textContent = (podpis.textContent ?? '').replace(/^\s*\d+/u, String(liczba));
}

/** Stawia na karcie plakietkę wersji i plakietkę zasięgu; stan czynności eksperta nie ma pola nazwanego słowem, więc trzecia plakietka schodzi. */
function opiszPlakietki(karta: HTMLElement, wzory: Wzory, ekspert: Agent): void {
  const pas = karta.querySelector<HTMLElement>('.ab-karta-plak');
  if (pas === null || wzory.plakietka === null) return;
  pas.replaceChildren();
  if (ekspert.version !== undefined) {
    const wersja = wzory.plakietka.cloneNode(true) as HTMLElement;
    wersja.textContent = (wersja.textContent ?? '').replace(/\d+/u, String(ekspert.version));
    pas.appendChild(wersja);
  }
  const zasieg = wzory.plakietka.cloneNode(true) as HTMLElement;
  zasieg.textContent = ekspert.visibility;
  pas.appendChild(zasieg);
}

/** Wczytuje definicję eksperta do edytora wraz z wersjami, rozszerzeniami, umiejętnościami i polityką; brak eksperta czyści edytor. */
async function pokazEksperta(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
  stan: Stan,
  identyfikator: string,
): Promise<void> {
  const eksperci = await wywolaj(kanal, Command.AgentList, {});
  const wykaz = eksperci.udany ? (eksperci.wynik?.agents ?? []) : [];
  const ekspert = wykaz.find((kandydat) => kandydat.id === identyfikator);
  stan.wybrany = ekspert;
  for (const pozycja of cialo.querySelectorAll<HTMLElement>('.dn-szyna-modulu-lista .pt-pozycja')) {
    pozycja.toggleAttribute('aria-current', pozycja.dataset.ekspert === identyfikator);
  }
  wpiszTekstAlboZdejmij(wezly.belka, ekspert?.name ?? '');
  wypelnijTozsamosc(cialo, wzory, ekspert, stan.moduly);
  wypelnijPodsumowanie(wezly, ekspert);
  wypelnijUmiejetnosci(cialo, wzory, ekspert);
  await wypelnijKanaly(kanal, cialo, wzory, ekspert);
  await wypelnijWersje(kanal, cialo, wzory, ekspert);
  await wypelnijRozszerzenia(kanal, cialo, wzory, ekspert);
  await wypelnijPolityke(kanal, idOkna, cialo, wzory, ekspert, stan.moduly);
}

/** Wypełnia zakładkę tożsamości: nazwę, opis, instrukcje, moduły zastosowania, zasięg widoczności i poziomy pamięci. */
function wypelnijTozsamosc(
  cialo: HTMLElement,
  wzory: Wzory,
  ekspert: Agent | undefined,
  moduly: Module[],
): void {
  wpiszWartosc(cialo, '#ab-nazwa', ekspert?.name ?? '');
  wpiszWartosc(cialo, '.ab-nazwa-inline', ekspert?.name ?? '');
  wpiszWartosc(cialo, '#ab-opis', ekspert?.description ?? '');
  wpiszWartosc(cialo, '#ab-instr', ekspert?.systemPrompt ?? '');
  const grupy = [...cialo.querySelectorAll<HTMLElement>('#tab-tozsamosc .ab-grupa')];
  wypelnijModuly(grupy[0] ?? null, wzory, moduly, ekspert?.moduleCodes ?? []);
  wypelnijZasieg(grupy[1] ?? null, wzory, ekspert?.visibility);
  wypelnijPamiec(grupy[2] ?? null, wzory, ekspert?.memoryLevels ?? []);
  const tytul = cialo.querySelector('.ab-historia .ab-sum-tytul span');
  if (tytul !== null) tytul.textContent = ekspert?.name ?? '';
}

/**
 * Wpisuje wartość w pole formularza edytora. Treść wyjściowa pola idzie razem
 * z wartością: prototypowy tekst obszaru i atrybut wartości wpisu zostają
 * w dokumencie po samym podstawieniu wartości i nadal są treścią przykładową.
 */
function wpiszWartosc(cialo: HTMLElement, wybor: string, wartosc: string): void {
  const pole = cialo.querySelector<HTMLInputElement | HTMLTextAreaElement>(wybor);
  if (pole === null) return;
  if (pole instanceof HTMLTextAreaElement) pole.textContent = wartosc;
  else pole.setAttribute('value', wartosc);
  pole.value = wartosc;
}

/** Stawia po jednej opcji na moduł rejestru; zbiór pusty przy ekspercie znaczy brak ograniczenia, więc zaznacza wszystkie. */
function wypelnijModuly(
  grupa: HTMLElement | null,
  wzory: Wzory,
  moduly: Module[],
  kody: string[],
): void {
  if (grupa === null || wzory.opcja === null) return;
  grupa.replaceChildren();
  for (const modul of moduly) {
    const opcja = wzory.opcja.cloneNode(true) as HTMLElement;
    opcja.dataset.modulKod = modul.code;
    const pole = opcja.querySelector('input');
    if (pole !== null) pole.checked = kody.length === 0 || kody.includes(modul.code);
    if (wpiszTekst(opcja, modul.name)) grupa.appendChild(opcja);
  }
}

/** Stawia po jednej opcji na wartość zasięgu widoczności; zaznaczenie bierze się z definicji eksperta. */
function wypelnijZasieg(
  grupa: HTMLElement | null,
  wzory: Wzory,
  zasieg: AgentVisibility | undefined,
): void {
  if (grupa === null || wzory.wybor === null) return;
  grupa.replaceChildren();
  for (const wartosc of Object.values(AgentVisibility)) {
    const opcja = wzory.wybor.cloneNode(true) as HTMLElement;
    opcja.dataset.zasieg = wartosc;
    const pole = opcja.querySelector('input');
    if (pole !== null) pole.checked = wartosc === zasieg;
    if (wpiszTekst(opcja, wartosc)) grupa.appendChild(opcja);
  }
}

/**
 * Stawia po jednej opcji na poziom pamięci. Piątej opcji „wyłączona" nie ma:
 * wyłączenie pamięci jest zbiorem pustym poziomów, a nie osobną wartością.
 */
function wypelnijPamiec(
  grupa: HTMLElement | null,
  wzory: Wzory,
  poziomy: MemoryLevel[],
): void {
  if (grupa === null || wzory.opcja === null) return;
  grupa.replaceChildren();
  for (const wartosc of Object.values(MemoryLevel)) {
    const opcja = wzory.opcja.cloneNode(true) as HTMLElement;
    opcja.dataset.poziom = wartosc;
    const pole = opcja.querySelector('input');
    if (pole !== null) pole.checked = poziomy.includes(wartosc);
    if (wpiszTekst(opcja, wartosc)) grupa.appendChild(opcja);
  }
}

/** Wypełnia podsumowanie definicji wartościami z tożsamości eksperta; wiersz bez wartości schodzi do czasu, gdy rdzeń ją poda. */
function wypelnijPodsumowanie(wezly: Wezly, ekspert: Agent | undefined): void {
  const rozszerzenia = (ekspert?.connectorIds?.length ?? 0) + (ekspert?.pluginIds?.length ?? 0);
  const wartosci =
    ekspert === undefined
      ? []
      : [
          ekspert.model ?? '',
          ekspert.transport ?? '',
          String(ekspert.skillIds?.length ?? 0),
          String(rozszerzenia),
          (ekspert.memoryLevels ?? []).join(' · '),
        ];
  for (const [numer, wiersz] of wezly.podsumowanie.entries()) {
    const wartosc = wartosci[numer] ?? '';
    const pole = wiersz.querySelector('b');
    if (pole !== null && wartosc !== '') pole.textContent = wartosc;
    zdejmijAlboPostaw(wiersz, wartosc !== '');
  }
}

/** Wypełnia wykaz umiejętności przypisaniami eksperta; kontrakt niesie przy ekspercie sam identyfikator umiejętności, więc on staje w wierszu. */
function wypelnijUmiejetnosci(
  cialo: HTMLElement,
  wzory: Wzory,
  ekspert: Agent | undefined,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-skills .sta-okno-tresc');
  if (panel === null || wzory.umiejetnosc === null) return;
  for (const stojacy of panel.querySelectorAll('.pk-wiersz')) stojacy.remove();
  const umiejetnosci = ekspert?.skillIds ?? [];
  opiszZnacznik(cialo, '#panel-skills', umiejetnosci.length);
  for (const umiejetnosc of umiejetnosci) {
    const wiersz = wzory.umiejetnosc.cloneNode(true) as HTMLElement;
    const nazwa = wiersz.querySelector('.nazwa');
    if (nazwa !== null) nazwa.textContent = umiejetnosc;
    wiersz.querySelector('[role="switch"]')?.setAttribute('aria-checked', 'true');
    panel.appendChild(wiersz);
  }
}

/** Wypełnia wykaz rozszerzeń konektorami i wtyczkami eksperta; rodzaj konektora wchodzi w plakietkę, a wtyczka plakietki rodzaju nie ma. */
async function wypelnijRozszerzenia(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  ekspert: Agent | undefined,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-connectors .sta-okno-tresc');
  if (panel === null || wzory.rozszerzenie === null) return;
  for (const stojacy of panel.querySelectorAll('.pk-wiersz')) stojacy.remove();
  if (ekspert === undefined) {
    opiszZnacznik(cialo, '#panel-connectors', 0);
    return;
  }
  const konektory = await wywolaj(kanal, Command.AgentConnectorList, { agentId: ekspert.id });
  const wtyczki = await wywolaj(kanal, Command.AgentPluginList, { agentId: ekspert.id });
  const wykazKonektorow = konektory.udany ? (konektory.wynik?.connectors ?? []) : [];
  const wykazWtyczek = wtyczki.udany ? (wtyczki.wynik?.plugins ?? []) : [];
  opiszZnacznik(cialo, '#panel-connectors', wykazKonektorow.length + wykazWtyczek.length);
  const baner = panel.querySelector('.pk-baner');
  for (const konektor of wykazKonektorow) {
    panel.insertBefore(wierszRozszerzenia(wzory, konektor, undefined), baner);
  }
  for (const wtyczka of wykazWtyczek) {
    panel.insertBefore(wierszRozszerzenia(wzory, undefined, wtyczka), baner);
  }
}

/** Buduje wiersz rozszerzenia z konektora albo z wtyczki; punkt dostępu i źródło wchodzą w podpis wiersza. */
function wierszRozszerzenia(
  wzory: Wzory,
  konektor: AgentConnector | undefined,
  wtyczka: AgentPlugin | undefined,
): HTMLElement {
  const wiersz = (wzory.rozszerzenie as HTMLElement).cloneNode(true) as HTMLElement;
  const nazwa = wiersz.querySelector('.nazwa');
  const rodzaj = wiersz.querySelector('.dn-plakietka');
  if (konektor !== undefined && rodzaj !== null) rodzaj.textContent = konektor.kind;
  else rodzaj?.remove();
  if (nazwa !== null) wpiszTekst(nazwa, konektor?.name ?? wtyczka?.name ?? '');
  const podpis = wiersz.querySelector('small');
  const opis = konektor?.accessPointId ?? wtyczka?.source ?? '';
  if (podpis !== null) {
    if (opis === '') podpis.remove();
    else podpis.textContent = opis;
  }
  const przelacznik = wiersz.querySelector('[role="switch"]');
  const czynne = konektor?.enabled ?? wtyczka?.enabled ?? false;
  przelacznik?.setAttribute('aria-checked', String(czynne));
  return wiersz;
}

/** Wypełnia historię wersji tożsamości eksperta; etykieta wersji wchodzi w wyróżnienie, opis zmiany w tekst, a czas w podpis. */
async function wypelnijWersje(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  ekspert: Agent | undefined,
): Promise<void> {
  const wykaz = cialo.querySelector<HTMLElement>('.ab-historia .dn-wykaz-modulu');
  if (wykaz === null || wzory.wersja === null) return;
  wykaz.replaceChildren();
  if (ekspert === undefined) return;
  const wynik = await wywolaj(kanal, Command.AgentVersionList, { agentId: ekspert.id });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const wersja of wynik.wynik.versions) {
    const wiersz = wzory.wersja.cloneNode(true) as HTMLElement;
    const etykieta = wiersz.querySelector('b');
    if (etykieta !== null) {
      if (wersja.label === undefined) etykieta.remove();
      else etykieta.textContent = wersja.label;
    }
    const meta = wiersz.querySelector('.dn-meta');
    if (meta !== null) meta.textContent = data(wersja.createdAt);
    wpiszTekst(wiersz, ` ${wersja.summary ?? ''} `);
    wykaz.appendChild(wiersz);
  }
}

/** Wypełnia wybór modelu bazowego rejestrem kanałów rdzenia oraz wybór drogi wywołania wartościami kontraktu. */
async function wypelnijKanaly(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  ekspert: Agent | undefined,
): Promise<void> {
  const modele = cialo.querySelector<HTMLElement>('#panel-model .pk-modele');
  const drogi = cialo.querySelector<HTMLElement>('#panel-model .pk-pig');
  if (modele === null || drogi === null || wzory.model === null || wzory.transport === null) return;

  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  modele.replaceChildren();
  for (const wpis of wynik.udany ? (wynik.wynik?.channels ?? []) : []) {
    const przycisk = wzory.model.cloneNode(true) as HTMLElement;
    const nazwa = przycisk.querySelector('b');
    if (nazwa !== null) nazwa.textContent = wpis.model ?? wpis.name;
    przycisk.setAttribute('aria-pressed', String(wpis.id === ekspert?.channelId));
    przycisk.addEventListener('click', () => {
      if (ekspert === undefined) return;
      void wywolaj(kanal, Command.AgentModelSet, {
        agentId: ekspert.id,
        channelId: wpis.id,
        model: wpis.model,
      });
    });
    modele.appendChild(przycisk);
  }

  drogi.replaceChildren();
  for (const droga of Object.values(ProviderTransport)) {
    const przycisk = wzory.transport.cloneNode(true) as HTMLElement;
    przycisk.textContent = droga;
    przycisk.setAttribute('aria-pressed', String(droga === ekspert?.transport));
    drogi.appendChild(przycisk);
  }
}

/** Wypełnia centrum zakresu: moduły dostępne ekspertowi, osiem zakresów izolacji technicznej i dostępność sieci podagentów. */
async function wypelnijPolityke(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  ekspert: Agent | undefined,
  moduly: Module[],
): Promise<void> {
  const grupa = cialo.querySelector<HTMLElement>('#grupa-moduly .ab-grupa');
  const macierz = cialo.querySelector<HTMLElement>('#grupa-izolacja .pk-macierz');
  macierz?.replaceChildren();
  if (ekspert === undefined) {
    if (grupa !== null) grupa.replaceChildren();
    return;
  }
  const wynik = await wywolaj(kanal, Command.AgentPolicyGet, {
    agentId: ekspert.id,
    windowId: idOkna,
  });
  const polityka = wynik.udany ? wynik.wynik?.policy : undefined;
  if (polityka === undefined) return;
  wypelnijModuly(grupa, wzory, moduly, polityka.moduleCodes);
  wypelnijIzolacje(macierz, wzory, polityka.technicalSwitches);
  const przelacznik = cialo.querySelector('#grupa-mtai [role="switch"]');
  przelacznik?.setAttribute('aria-checked', String(polityka.subagentEnabled));
  const granica = cialo.querySelector<HTMLInputElement>('#grupa-mtai .pk-num');
  if (granica !== null) granica.value = String(polityka.subagentLimit);
}

/** Stawia po jednym wierszu na zakres izolacji technicznej; nazwa zakresu pochodzi z kontraktu, a objaśnienie z odpowiedzi rdzenia. */
function wypelnijIzolacje(
  macierz: HTMLElement | null,
  wzory: Wzory,
  przelaczniki: IsolationTechnicalSwitch[],
): void {
  if (macierz === null || wzory.izolacja === null) return;
  for (const przelacznik of przelaczniki) {
    const wiersz = wzory.izolacja.cloneNode(true) as HTMLElement;
    const nazwa = wiersz.querySelector('.rozc');
    if (nazwa !== null) nazwa.textContent = przelacznik.scope;
    const kontrolka = wiersz.querySelector('[role="switch"]');
    kontrolka?.setAttribute('aria-checked', String(przelacznik.isolated));
    kontrolka?.setAttribute('aria-label', przelacznik.scope);
    if (przelacznik.explanation !== undefined) wiersz.title = przelacznik.explanation;
    macierz.appendChild(wiersz);
  }
}

/** Nanosi liczbę pozycji na znak panelu; podmieniana jest sama liczba, podpis prototypu zostaje. */
function opiszZnacznik(cialo: HTMLElement, panel: string, liczba: number): void {
  const znacznik = cialo.querySelector<HTMLElement>(`${panel} .sta-okno-znacznik`);
  if (znacznik === null) return;
  znacznik.textContent = (znacznik.textContent ?? '').replace(/^\s*\d+/u, String(liczba));
}

/** Data zapisu w postaci rocznej; kontrakt niesie czas w milisekundach epoki, a nie jego zapis słowny. */
function data(znacznik: number): string {
  return new Date(znacznik).toISOString().slice(0, 10);
}
