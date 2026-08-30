/**
 * Wiązanie wnętrza okna modułu Apps z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * apps.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  AppDeployStrategy,
  AppStageStatus,
  Command,
  EventType,
  type AppDeployEnvironment,
  type AppDeployment,
  type AppEndpoint,
  type AppMilestone,
  type AppProduct,
  type AppStage,
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
 * Panele prototypu, którym rodzina apps.* nie oddaje ani jednego pola;
 * schodzą razem z przełącznikami menu sesji. Warsztat frontendu schodzi
 * z nimi: jego węzłami są podgląd na żywo i wycinek kodu, a kontrakt niesie
 * wyłącznie polecenie postawienia serwera podglądu, bez odczytu jego stanu.
 */
const PANELE_BEZ_POKRYCIA = [
  'panel-frontend',
  'panel-plan',
  'panel-zadania',
  'panel-terminal',
  'panel-pliki',
];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  etap: HTMLElement | null;
  platforma: HTMLElement | null;
  kamien: HTMLElement | null;
  wydanie: HTMLElement | null;
  punkt: HTMLElement | null;
  wdrozenie: HTMLElement | null;
  srodowisko: HTMLElement | null;
  strategia: HTMLElement | null;
  powiazanie: HTMLElement | null;
}

/**
 * Węzły niosące jedną wartość, zdejmowane na czas jej braku. Wskazania stoją
 * tu, a nie w wyszukaniu przy każdym odczycie: węzeł zdjęty jest poza
 * dokumentem, więc wyszukanie już by go nie znalazło i wartość, która dojdzie,
 * nie miałaby czego postawić.
 */
interface Wezly {
  tytul: Element | null;
  belki: HTMLElement[];
  nastawy: Element[];
  repozytorium: Element | null;
  metryki: Element[];
  kondycja: Element | null;
  szablon: HTMLElement | null;
  wersja: HTMLElement | null;
  etap: Element | null;
}

/**
 * Wiąże wnętrze okna modułu Apps. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazApps(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);
  const wezly = zbierzWezly(cialo);

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt produktu, etapów, kamieni milowych i wdrożeń; wołane przy wejściu i po każdym zdarzeniu budowy. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, wezly);
  };

  zwiazNawigacjeKafli(cialo);
  zwiazWdrozenie(kanal, idOkna, cialo, odswiez);
  kanal.naZdarzenie(EventType.AppsBuildChanged, (tresc) => {
    if (tresc.stage.windowId === idOkna) odswiez();
  });

  odswiez();
}

/** Zdejmuje wzory pozycji i wierszy, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const etapy = [...cialo.querySelectorAll<HTMLElement>('.pb-tracker .pb-etap')];
  const etap = sklonuj(etapy[etapy.length - 1] ?? null);
  // Znak stanu etapu stoi w prototypie osobny dla każdego stopnia; etap niesie
  // stan wartością, a nie kropką, więc znak schodzi ze wzoru.
  etap?.querySelector('.dn-kropka')?.remove();
  etap?.querySelector('svg')?.remove();
  etap?.querySelector('.pt-tetno')?.remove();

  const listy = [...cialo.querySelectorAll<HTMLElement>('#panel-builder .pb-lista')];
  const kamien = sklonuj(listy[0]?.querySelector('.pb-wiersz') ?? null);
  kamien?.querySelector('.dn-kropka')?.remove();
  kamien?.querySelector('.pt-tetno')?.remove();

  const punkt = sklonuj(cialo.querySelector('#panel-backend .dn-wykaz-modulu-poz'));
  // Barwa metody i plakietka odpowiedzi są nadawane wartością punktu końcowego,
  // nie odmianą klasy; klasa odmiany schodzi, żeby nie barwiła cudzej wartości.
  const metoda = punkt?.querySelector('.metoda');
  if (metoda !== null && metoda !== undefined) metoda.className = 'metoda';
  const stanPunktu = punkt?.querySelector('.dn-plakietka');
  if (stanPunktu !== null && stanPunktu !== undefined) stanPunktu.className = 'dn-plakietka';

  const wdrozenie = sklonuj(cialo.querySelector('#panel-deployment .dn-tabela tr'));
  const stanWdrozenia = wdrozenie?.querySelector('.dn-plakietka');
  if (stanWdrozenia !== null && stanWdrozenia !== undefined) stanWdrozenia.className = 'dn-plakietka';
  wdrozenie?.querySelector('button')?.removeAttribute('data-komunikat');
  wdrozenie?.querySelector('button')?.removeAttribute('data-komunikat-tytul');

  return {
    etap,
    platforma: sklonuj(cialo.querySelector('.pb-naglowek .dn-plakietka')),
    kamien,
    wydanie: sklonuj(listy[2]?.querySelector('.pb-wiersz') ?? null),
    punkt,
    wdrozenie,
    srodowisko: sklonuj(cialo.querySelector('#panel-deployment option')),
    strategia: sklonuj(
      cialo.querySelector('#panel-deployment input[name="strategia"]')?.closest('label') ?? null,
    ),
    powiazanie: sklonuj(
      cialo.querySelector('#konfig-czat [role="switch"]')?.closest('.sta-konfig-pole') ?? null,
    ),
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz zleceń szyny, historia rozmowy
 * i przypięty kontekst wiszą na rodzinach spoza apps.*, a miary otwartych
 * zadań, błędów budowania i przypisania wykonawców nie mają pola w kontrakcie.
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
  for (const chip of cialo.querySelectorAll('.sta-kontekst-akcji > .sta-chip')) chip.remove();
  // Zatwierdzenie zmian warsztatu nie jest czynnością kontraktu: rodzina apps.*
  // zapisuje plik warsztatu, a nie zestaw zmian.
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
  zdejmijMiaryBezPokrycia(cialo);
  zdejmijNastawyBezPokrycia(cialo);
}

/** Zdejmuje miary i podpisy, dla których rodzina apps.* nie ma pola: otwarte zadania, błędy budowania, przypisanie wykonawców, wycinki kodu. */
function zdejmijMiaryBezPokrycia(cialo: HTMLElement): void {
  const metryki = [...cialo.querySelectorAll('#panel-builder .pb-metryka')];
  metryki[0]?.remove();
  metryki[1]?.remove();
  for (const podpis of cialo.querySelectorAll('#panel-builder .pb-kafel .op')) podpis.remove();
  // Wykaz wykonawców etapu wymagałby nazwy eksperta, a etap niesie wyłącznie
  // jego identyfikator; sam identyfikator nie nazywa wykonawcy.
  const listy = [...cialo.querySelectorAll('#panel-builder .pb-lista')];
  listy[1]?.parentElement?.remove();
  for (const kod of cialo.querySelectorAll('#panel-backend .kod-blok')) kod.remove();
  cialo.querySelector('#panel-backend .dn-plakietka--rola')?.remove();
  cialo.querySelector('#panel-backend')?.removeAttribute('data-stan');
  const plakietki = cialo.querySelector('#panel-backend .dn-plakietka--sygnal')?.parentElement;
  plakietki?.remove();
  // Stan strzałki między etapami jest rysunkiem prototypu, a etapy stawiane są
  // odpowiedzią rdzenia, więc strzałki nie mają się między czym stawiać.
  for (const strzalka of cialo.querySelectorAll('.pb-tracker .pb-strzalka')) strzalka.remove();
}

/**
 * Zdejmuje nastawy okna bez pokrycia: wykonawcę polecenia wraz z kanałem oraz
 * zastrzeżenia układu wypisane słowem. Walidacja układu oddaje wykaz zastrzeżeń,
 * a znacznik niesie na nie jedną plakietkę — jedno zastrzeżenie z wielu byłoby
 * wyborem wiązania, nie odpowiedzią rdzenia.
 */
function zdejmijNastawyBezPokrycia(cialo: HTMLElement): void {
  cialo.querySelector('#konfig-czat .sta-konfig-sekcja')?.remove();
  cialo.querySelector('#panel-architektura .ad-diagram')?.remove();
  cialo.querySelector('#panel-architektura .dn-plakietka--ostrzezenie')?.remove();
  cialo.querySelector('#panel-architektura .sta-okno-belka button.dn-btn')?.remove();
  cialo.querySelector('#panel-deployment [data-komunikat]')?.removeAttribute('data-komunikat');
  for (const przycisk of cialo.querySelectorAll('#panel-deployment [data-komunikat-tytul]')) {
    przycisk.removeAttribute('data-komunikat-tytul');
  }
}

/**
 * Zdejmuje sterowanie, którego rodzina apps.* nie obsługuje: panele bez
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
}

/** Pyta rdzeń o produkt, etapy, kamienie milowe, wdrożenia, architekturę i punkty końcowe, po czym wypełnia nimi wnętrze okna. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
): Promise<void> {
  const produkt = await wywolaj(kanal, Command.AppsProductGet, { windowId: idOkna });
  const etapy = await wywolaj(kanal, Command.AppsStageList, { windowId: idOkna });
  const kamienie = await wywolaj(kanal, Command.AppsMilestoneList, { windowId: idOkna });
  const wdrozenia = await wywolaj(kanal, Command.AppsDeploymentList, { windowId: idOkna });

  const wykazEtapow = etapy.udany ? (etapy.wynik?.stages ?? []) : [];
  const wykazWdrozen = wdrozenia.udany ? (wdrozenia.wynik?.deployments ?? []) : [];

  opiszProdukt(cialo, wzory, wezly, produkt.udany ? produkt.wynik?.product : undefined);
  wypelnijEtapy(cialo, wzory, wezly, wykazEtapow);
  wypelnijKamienie(cialo, wzory, kamienie.udany ? (kamienie.wynik?.milestones ?? []) : []);
  wypelnijDziennik(cialo, wzory, wykazWdrozen);
  wypelnijWdrozenia(cialo, wzory, wykazWdrozen);
  await opiszKondycje(kanal, idOkna, wezly, wykazWdrozen[0]);
  await opiszArchitekture(kanal, idOkna, wezly);
  await wypelnijPunktyKoncowe(kanal, idOkna, cialo, wzory);
  await wypelnijSrodowiska(kanal, idOkna, cialo, wzory);
  await wypelnijPowiazania(kanal, idOkna, cialo, wzory);
  wypelnijStrategie(cialo, wzory);
}

/** Nanosi metadane produktu na nagłówek okna wiodącego, na belki okien i na nastawę produktu; brak produktu zdejmuje każdy z tych węzłów. */
function opiszProdukt(
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
  produkt: AppProduct | undefined,
): void {
  const nazwa = produkt?.name ?? '';
  wpiszAlboZdejmij(wezly.tytul, nazwa);
  for (const belka of wezly.belki) {
    /* Tytuł okna niesie nazwę produktu za myślnikiem. Człon przed myślnikiem
       zapamiętany jest przy pierwszym odczycie: po podstawieniu nazwy nie
       dałoby się go już odczytać z samego tytułu. */
    belka.dataset.tytul ??= (belka.textContent ?? '').replace(/\s—\s.*$/u, '');
    belka.textContent = nazwa === '' ? belka.dataset.tytul : `${belka.dataset.tytul} — ${nazwa}`;
  }
  cialo.querySelector('#panel-builder .sta-okno-belka .dn-plakietka')?.remove();
  wypelnijPlatformy(cialo, wzory, produkt?.platforms ?? []);
  const repozytorium = produkt?.repositoryUrl ?? '';
  /* Nastawa produktu niesie dwa pola w tej kolejności: platformy docelowe
     i repozytorium źródłowe; trzecie miejsce repozytorium stoi w nagłówku. */
  const wartosci = [(produkt?.platforms ?? []).join(' · '), repozytorium];
  for (const [numer, miejsce] of wezly.nastawy.entries()) {
    wpiszAlboZdejmij(miejsce, wartosci[numer] ?? '');
  }
  wpiszAlboZdejmij(wezly.repozytorium, repozytorium);
}

/** Zbiera węzły jednowartościowe po zdjęciu treści przykładowej; część z nich zdejmuje dopiero brak wartości, więc wskazania muszą przetrwać odczyt. */
function zbierzWezly(cialo: HTMLElement): Wezly {
  const plakietki = [...cialo.querySelectorAll<HTMLElement>('#panel-architektura .dn-plakietka')];
  return {
    tytul: cialo.querySelector('#panel-builder .pb-naglowek h2'),
    /* Nazwę produktu niosą tylko te tytuły okien, które mają ją w prototypie
       za myślnikiem; pozostałe nazywają samo okno i produktu nie dotyczą. */
    belki: [...cialo.querySelectorAll<HTMLElement>('.sta-okno-belka .sta-okno-tytul b')].filter(
      (belka) => /\s—\s/u.test(belka.textContent ?? ''),
    ),
    nastawy: [...cialo.querySelectorAll('#konfig-builder .sta-konfig-pole .pt-mono')],
    repozytorium: cialo.querySelector('#panel-builder .pb-naglowek .pt-mono'),
    metryki: [...cialo.querySelectorAll('#panel-builder .pb-metryka')],
    kondycja: cialo.querySelector('#panel-deployment .dn-plakietka'),
    szablon: plakietki[0] ?? null,
    wersja: plakietki[1] ?? null,
    etap: cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'),
  };
}

/** Stawia w nagłówku produktu po jednej plakietce na platformę docelową; plakietki prototypu schodzą przed powieleniem wzoru. */
function wypelnijPlatformy(cialo: HTMLElement, wzory: Wzory, platformy: string[]): void {
  const naglowek = cialo.querySelector<HTMLElement>('#panel-builder .pb-naglowek');
  if (naglowek === null || wzory.platforma === null) return;
  for (const stojaca of naglowek.querySelectorAll('.dn-plakietka')) stojaca.remove();
  const kotwica = naglowek.querySelector('.pt-mono');
  for (const platforma of platformy) {
    const plakietka = wzory.platforma.cloneNode(true) as HTMLElement;
    plakietka.textContent = platforma;
    naglowek.insertBefore(plakietka, kotwica);
  }
}

/** Stawia w pasie etapów po jednym kroku na etap budowy; etap w toku dostaje wskazanie kroku bieżącego. */
function wypelnijEtapy(
  cialo: HTMLElement,
  wzory: Wzory,
  wezly: Wezly,
  etapy: AppStage[],
): void {
  const pas = cialo.querySelector<HTMLElement>('#panel-builder .pb-tracker');
  if (pas === null || wzory.etap === null) return;
  pas.replaceChildren();
  for (const etap of etapy) {
    const krok = wzory.etap.cloneNode(true) as HTMLElement;
    krok.textContent = etap.name;
    if (etap.status === AppStageStatus.Active) krok.setAttribute('aria-current', 'step');
    pas.appendChild(krok);
  }
  const czynny = etapy.find((etap) => etap.status === AppStageStatus.Active);
  wpiszTekstAlboZdejmij(wezly.etap, czynny?.name ?? '');
}

/** Wypełnia wykaz kamieni milowych; stan realizacji i termin wchodzą w podpis pozycji, bo tylko one mają pole. */
function wypelnijKamienie(cialo: HTMLElement, wzory: Wzory, kamienie: AppMilestone[]): void {
  const lista = [...cialo.querySelectorAll<HTMLElement>('#panel-builder .pb-lista')][0];
  if (lista === undefined || wzory.kamien === null) return;
  lista.replaceChildren();
  for (const kamien of kamienie) {
    const wiersz = wzory.kamien.cloneNode(true) as HTMLElement;
    const meta = wiersz.querySelector('.dn-meta');
    if (meta !== null) {
      meta.textContent =
        kamien.dueAt === undefined ? kamien.status : `${kamien.status} · ${data(kamien.dueAt)}`;
    }
    if (wpiszTekst(wiersz, ` ${kamien.name} `)) lista.appendChild(wiersz);
  }
}

/** Wypełnia dziennik wydań wdrożeniami okna; wersja i notatka wydania stoją w tekście pozycji, a czas rozpoczęcia w podpisie. */
function wypelnijDziennik(cialo: HTMLElement, wzory: Wzory, wdrozenia: AppDeployment[]): void {
  const listy = [...cialo.querySelectorAll<HTMLElement>('#panel-builder .pb-lista')];
  const lista = listy[listy.length - 1];
  if (lista === undefined || wzory.wydanie === null || lista === listy[0]) return;
  lista.replaceChildren();
  for (const wdrozenie of wdrozenia) {
    const wiersz = wzory.wydanie.cloneNode(true) as HTMLElement;
    const meta = wiersz.querySelector('.dn-meta');
    if (meta !== null) meta.textContent = data(wdrozenie.startedAt);
    const opis = [wdrozenie.version ?? '', wdrozenie.releaseNotes ?? ''].filter(
      (czesc) => czesc !== '',
    );
    if (opis.length > 0 && wpiszTekst(wiersz, ` ${opis.join(' — ')} `)) lista.appendChild(wiersz);
  }
}

/** Wypełnia historię wdrożeń wierszami tabeli; przycisk wiersza cofa produkt do tej wersji tą samą komendą, która wdraża. */
function wypelnijWdrozenia(cialo: HTMLElement, wzory: Wzory, wdrozenia: AppDeployment[]): void {
  const tabela = cialo.querySelector<HTMLElement>('#panel-deployment .dn-tabela tbody');
  if (tabela === null || wzory.wdrozenie === null) return;
  tabela.replaceChildren();
  for (const wdrozenie of wdrozenia) {
    const wiersz = wzory.wdrozenie.cloneNode(true) as HTMLElement;
    wiersz.dataset.wdrozenie = wdrozenie.id;
    wiersz.dataset.srodowisko = wdrozenie.environment;
    const komorki = [...wiersz.querySelectorAll('td')];
    const wersja = komorki[0];
    if (wersja !== undefined) wersja.textContent = wdrozenie.version ?? '';
    const czas = komorki[1];
    if (czas !== undefined) czas.textContent = data(wdrozenie.startedAt);
    const stan = wiersz.querySelector('.dn-plakietka');
    if (stan !== null) stan.textContent = wdrozenie.status;
    tabela.appendChild(wiersz);
  }
}

/** Nanosi na panel kondycji wersję ostatniego wdrożenia i dostępność produktu; brak którejkolwiek wartości zdejmuje jej metrykę. */
async function opiszKondycje(
  kanal: Kanal,
  idOkna: string,
  wezly: Wezly,
  ostatnie: AppDeployment | undefined,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsDeploymentHealthGet, { windowId: idOkna });
  const kondycja = wynik.udany ? wynik.wynik?.health : undefined;
  const wersja = ostatnie?.version ?? '';
  const dostepnosc = kondycja?.availabilityPercent ?? '';
  for (const [numer, metryka] of wezly.metryki.entries()) {
    const wartosc = [wersja, dostepnosc][numer] ?? '';
    const liczba = metryka.querySelector('.liczba');
    if (liczba !== null && wartosc !== '') liczba.textContent = wartosc;
    zdejmijAlboPostaw(metryka, wartosc !== '');
  }
  wezly.kondycja?.querySelector('.dn-kropka')?.remove();
  wpiszAlboZdejmij(wezly.kondycja, [wersja, dostepnosc].filter((czesc) => czesc !== '').join(' · '));
}

/** Nanosi na panel architektury szablon układu i numer wersji; brak architektury zdejmuje obie plakietki. */
async function opiszArchitekture(kanal: Kanal, idOkna: string, wezly: Wezly): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsArchitectureGet, { windowId: idOkna });
  const uklad = wynik.udany ? wynik.wynik?.architecture : undefined;
  const szablon = wezly.szablon;
  if (szablon !== null) {
    szablon.dataset.wzor ??= szablon.textContent ?? '';
    const podpis = szablon.dataset.wzor;
    wpiszAlboZdejmij(
      szablon,
      uklad === undefined ? '' : podpis.replace(/:.*$/u, `: ${uklad.template}`),
    );
  }
  const wersja = wezly.wersja;
  if (wersja === null) return;
  wersja.dataset.wzor ??= wersja.textContent ?? '';
  const podpis = wersja.dataset.wzor;
  wpiszAlboZdejmij(
    wersja,
    uklad?.version === undefined ? '' : podpis.replace(/\d+/u, String(uklad.version)),
  );
}

/** Wypełnia wykaz punktów końcowych warsztatu backendu; metoda i stan punktu stoją wartościami kontraktu, bez barwy nadanej z góry. */
async function wypelnijPunktyKoncowe(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): Promise<void> {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-backend .dn-wykaz-modulu');
  if (wykaz === null || wzory.punkt === null) return;
  for (const stojacy of wykaz.querySelectorAll('.dn-wykaz-modulu-poz')) stojacy.remove();
  const wynik = await wywolaj(kanal, Command.AppsEndpointList, { windowId: idOkna });
  for (const punkt of wynik.udany ? (wynik.wynik?.endpoints ?? []) : []) {
    wykaz.appendChild(wierszPunktu(wzory, punkt));
  }
}

/** Buduje wiersz punktu końcowego; brak stanu punktu zdejmuje plakietkę, bo pusta plakietka niczego nie mówi. */
function wierszPunktu(wzory: Wzory, punkt: AppEndpoint): HTMLElement {
  const wiersz = (wzory.punkt as HTMLElement).cloneNode(true) as HTMLElement;
  const metoda = wiersz.querySelector('.metoda');
  if (metoda !== null) metoda.textContent = punkt.method;
  const stan = wiersz.querySelector('.dn-plakietka');
  if (stan !== null) {
    if (punkt.status === undefined) stan.remove();
    else stan.textContent = punkt.status;
  }
  wpiszTekst(wiersz, ` ${punkt.path} `);
  return wiersz;
}

/** Wypełnia wybór środowiska wdrożenia środowiskami produktu; kod środowiska wchodzi w wartość, a nazwa w podpis pozycji. */
async function wypelnijSrodowiska(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): Promise<void> {
  const wybor = cialo.querySelector<HTMLSelectElement>('#panel-deployment select');
  if (wybor === null || wzory.srodowisko === null) return;
  const wynik = await wywolaj(kanal, Command.AppsEnvironmentList, { windowId: idOkna });
  wybor.replaceChildren();
  for (const srodowisko of wynik.udany ? (wynik.wynik?.environments ?? []) : []) {
    const pozycja = wzory.srodowisko.cloneNode(true) as HTMLOptionElement;
    pozycja.value = srodowisko.code;
    pozycja.textContent = srodowisko.name;
    wybor.appendChild(pozycja);
  }
}

/** Wypełnia wybór strategii wdrożenia wartościami kontraktu; żadna nie stoi zaznaczona, bo strategii domyślnej kontrakt nie wskazuje. */
function wypelnijStrategie(cialo: HTMLElement, wzory: Wzory): void {
  const wzor = wzory.strategia;
  if (wzor === null) return;
  const gniazdo = cialo.querySelector<HTMLElement>('#panel-deployment input[name="strategia"]')
    ?.parentElement?.parentElement;
  if (gniazdo === null || gniazdo === undefined) return;
  for (const stojaca of gniazdo.querySelectorAll('label.dn-wybor')) stojaca.remove();
  for (const strategia of Object.values(AppDeployStrategy)) {
    const opcja = wzor.cloneNode(true) as HTMLElement;
    opcja.dataset.strategia = strategia;
    const pole = opcja.querySelector('input');
    if (pole !== null) pole.checked = false;
    if (wpiszTekst(opcja, ` ${strategia}`)) gniazdo.appendChild(opcja);
  }
}

/** Wypełnia nastawę powiązań produktu stanem powiązań z innymi modułami; stan powiązania stoi przełącznikiem, a jego opis podpowiedzią. */
async function wypelnijPowiazania(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): Promise<void> {
  const sekcja = cialo.querySelector<HTMLElement>('#konfig-czat .sta-konfig-sekcja');
  if (sekcja === null || wzory.powiazanie === null) return;
  const wynik = await wywolaj(kanal, Command.AppsProductLinkList, { windowId: idOkna });
  for (const stojace of sekcja.querySelectorAll('.sta-konfig-pole')) stojace.remove();
  for (const powiazanie of wynik.udany ? (wynik.wynik?.links ?? []) : []) {
    const pole = wzory.powiazanie.cloneNode(true) as HTMLElement;
    const podpis = pole.querySelector('label');
    if (podpis !== null) podpis.textContent = powiazanie.moduleCode;
    const przelacznik = pole.querySelector('[role="switch"]');
    przelacznik?.setAttribute('aria-checked', String(powiazanie.enabled));
    przelacznik?.setAttribute('aria-label', powiazanie.moduleCode);
    if (powiazanie.detail !== undefined) pole.title = powiazanie.detail;
    sekcja.appendChild(pole);
  }
}

/** Wiąże kafle okna wiodącego z panelami, do których prowadzą; kafel panelu zdjętego schodzi razem z nim. */
function zwiazNawigacjeKafli(cialo: HTMLElement): void {
  for (const kafel of cialo.querySelectorAll<HTMLElement>('[data-skok]')) {
    const panel = cialo.querySelector(`#${kafel.dataset.skok ?? ''}`);
    if (panel === null) {
      kafel.remove();
      continue;
    }
    kafel.addEventListener('click', () => {
      panel.scrollIntoView({ block: 'nearest' });
    });
  }
}

/** Wiąże wdrożenie produktu: przycisk panelu wdraża do wybranego środowiska, a przycisk wiersza cofa produkt do wdrożenia tego wiersza. */
function zwiazWdrozenie(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-deployment');
  if (panel === null) return;
  panel.querySelector('button.dn-btn--atrament')?.addEventListener('click', () => {
    void wdroz(kanal, idOkna, cialo, undefined, odswiez);
  });
  panel.querySelector('.dn-tabela')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('tr');
    const identyfikator = wiersz?.dataset.wdrozenie;
    if (identyfikator === undefined) return;
    void wdroz(kanal, idOkna, cialo, identyfikator, odswiez);
  });
}

/** Uruchamia wdrożenie albo cofnięcie do wskazanego wdrożenia; środowisko i strategię bierze z pól panelu. */
async function wdroz(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  cofnijDo: string | undefined,
  odswiez: () => void,
): Promise<void> {
  const wybor = cialo.querySelector<HTMLSelectElement>('#panel-deployment select');
  const srodowisko = wybor?.value ?? '';
  if (srodowisko === '') return;
  const wynik = await wywolaj(kanal, Command.AppsDeploymentRun, {
    windowId: idOkna,
    environment: srodowisko as AppDeployEnvironment,
    strategy: wybranaStrategia(cialo),
    rollbackToDeploymentId: cofnijDo,
  });
  if (wynik.udany) odswiez();
}

/** Strategia zaznaczona w panelu wdrożeń; brak zaznaczenia zostawia wybór strategii rdzeniowi. */
function wybranaStrategia(cialo: HTMLElement): AppDeployStrategy | undefined {
  const zaznaczona = cialo.querySelector<HTMLElement>(
    '#panel-deployment input[name="strategia"]:checked',
  );
  const wartosc = zaznaczona?.closest<HTMLElement>('[data-strategia]')?.dataset.strategia ?? '';
  return Object.values(AppDeployStrategy).find((kandydat) => kandydat === wartosc);
}

/** Data zdarzenia w postaci rocznej; kontrakt niesie czas w milisekundach epoki, a nie jego zapis słowny. */
function data(znacznik: number): string {
  return new Date(znacznik).toISOString().slice(0, 10);
}
