/**
 * Wiązanie wnętrza okna modułu Browser z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * browser.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla których
 * kontrakt nie ma pola, zdejmuje. Stronę renderuje rdzeń, więc płótno okna
 * bierze tekst migawki, nie osadzoną przeglądarkę.
 */
import {
  BrowserTabState,
  Command,
  EventType,
  type BrowserNote,
  type BrowserSnapshot,
  type BrowserSource,
  type BrowserTab,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  godzina,
  opiszWezel,
  sklonuj,
  uzgodnijPrzelacznikiPaneli,
  wpiszTekst,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';

/** Panele prototypu, którym rodzina browser.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-artefakty', 'panel-zadania', 'panel-pliki'];

/** Wzory wierszy i pozycji zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone. */
interface Wzory {
  grupa: HTMLElement | null;
  pozycja: HTMLElement | null;
  karta: HTMLElement | null;
  zrodlo: HTMLElement | null;
  gwiazdka: HTMLElement | null;
  notatka: HTMLElement | null;
  notatkaPrzypieta: HTMLElement | null;
}

/** Karta bieżąca okna i adres jej strony; obie wartości pochodzą z odpowiedzi rdzenia. */
interface Stan {
  karta: string;
  adres: string;
  zrodla: Map<string, BrowserSource>;
}

/**
 * Wiąże wnętrze okna modułu Browser. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazBrowser(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { karta: '', adres: '', zrodla: new Map() };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt kart, strony, źródeł i notatek; wołane przy wejściu i po każdej czynności przeglądania. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, stan);
  };

  zwiazPasekAdresu(kanal, idOkna, cialo, stan, odswiez);
  zwiazKarty(kanal, idOkna, cialo, stan, odswiez);
  zwiazZrodla(kanal, idOkna, cialo, wzory, stan, odswiez);
  zwiazNotatki(kanal, idOkna, cialo, stan);
  zwiazSzyne(kanal, idOkna, cialo, odswiez);
  zwiazZdarzeniaStrony(kanal, idOkna, cialo, stan);

  odswiez();
}

/** Zdejmuje wzory kart, źródeł i notatek; wzór traci oznaczenia, których kontrakt bytowi nie nadaje. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const pozycja = sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .pt-pozycja'));
  pozycja?.querySelector('.pt-tetno')?.remove();
  pozycja?.removeAttribute('aria-current');
  const karta = sklonuj(cialo.querySelector('.brw-zakladki-tory .brw-tab'));
  karta?.querySelector('.pt-tetno')?.remove();
  karta?.setAttribute('aria-selected', 'false');
  const zrodla = [...cialo.querySelectorAll('#panel-sources .brw-zrodlo-poz')];
  const notatki = [...cialo.querySelectorAll('#panel-notes .brw-notatka')];
  const zrodlo = sklonuj(zrodla[1] ?? null);
  // Podgląd migawki schodzi ze wzoru: okno nie ma miejsca, w którym pokazałoby
  // zrzut inny niż strona bieżąca płótna.
  zdejmijPrzycisk(zrodlo, 'Podgląd');
  const notatka = sklonuj(notatki[1] ?? null);
  const notatkaPrzypieta = sklonuj(notatki[0] ?? null);
  // Zmiana notatki bierze jej nową treść, a panel pola do wpisania nie niesie.
  zdejmijPrzycisk(notatka, 'Edytuj');
  zdejmijPrzycisk(notatkaPrzypieta, 'Edytuj');
  return {
    grupa: sklonuj(cialo.querySelector('.dn-szyna-modulu-lista .dn-etyk-mono')),
    pozycja,
    karta,
    // Wzorem źródła jest pozycja zwykła; gwiazdka wypełniona wchodzi osobno,
    // bo oznacza źródło kluczowe, a to jest pole odpowiedzi.
    zrodlo,
    gwiazdka: sklonuj(zrodla[0]?.querySelector('.brw-zrodlo-tyt svg') ?? null),
    notatka,
    notatkaPrzypieta,
  };
}

/** Zdejmuje ze wzoru przycisk rozpoznany po początku podpisu; wzór bez węzła kończy funkcję bez działania. */
function zdejmijPrzycisk(wzor: HTMLElement | null, podpis: string): void {
  if (wzor === null) return;
  for (const przycisk of wzor.querySelectorAll('.sta-chip-rzad .dn-btn')) {
    if (przycisk.textContent?.startsWith(podpis) === true) przycisk.remove();
  }
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów wisi na rodzinach spoza
 * browser.*, a zaznaczony fragment strony, projekt i karta sesji nie mają pola
 * w rodzinie okna przeglądarki.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  zdejmijTrescWspolna(cialo);
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  oznaczChipyKontekstu(cialo);
  oznaczChipyCzynnosci(cialo);
  cialo.querySelector('#panel-browser .dn-plakietka--informacja')?.remove();
  cialo.querySelector('#panel-browser .brw-obecnosc')?.remove();
  cialo.querySelector('#panel-browser .brw-strona')?.replaceChildren();
  cialo.querySelector('.brw-zakladki-tory')?.replaceChildren();
}

/** Zostawia w pasie kontekstu wyłącznie znak strony i oznacza go do wypełnienia; zaznaczony fragment nie ma pola w migawce. */
function oznaczChipyKontekstu(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kom-kontekst > .sta-zrodlo')];
  chipy[1]?.remove();
  const strona = chipy[0];
  if (strona !== undefined) strona.dataset.pole = 'strona';
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania, projekt i karta sesji
 * schodzą, a znak liczby źródeł zostaje oznaczony do wypełnienia wykazem
 * źródeł okna.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  for (const numer of [0, 1, 2, 4]) chipy[numer]?.remove();
  const zrodla = chipy[3];
  if (zrodla !== undefined) zrodla.dataset.pole = 'zrodla';
  // Przekazanie zebranego materiału do modułu Research nie ma komendy
  // w rodzinie browser.*: wykaz źródeł jest bytem okna przeglądarki.
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina browser.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, ruch wstecz i naprzód, emulację
 * urządzenia bez wskazania go, tryb czytnika, podział ekranu, wywóz wykazu
 * i zmianę notatki bez pola jej treści.
 */
function zdejmijSterowanieBezPokrycia(cialo: HTMLElement): void {
  uzgodnijPrzelacznikiPaneli(cialo, PANELE_BEZ_POKRYCIA);
  zdejmijSterowanieWspolne(cialo);
  const naglowek = [...cialo.querySelectorAll('#panel-browser .sta-okno-akcje > .dn-btn-ikona')];
  for (const [numer, przycisk] of naglowek.entries()) {
    if (numer < 3) przycisk.remove();
  }
  cialo.querySelector('#panel-browser .sta-okno-akcje .sta-menu')?.remove();
  // Historia przeglądania nie ma komendy: rdzeń zna przejście pod adres,
  // ruchu wstecz i naprzód — nie.
  const pasek = [...cialo.querySelectorAll('.brw-pasek > .dn-btn-ikona')];
  pasek[0]?.remove();
  pasek[1]?.remove();
  cialo.querySelector('.brw-stopka')?.remove();
  // Założenie rozmowy jest czynnością karty sesji, nie okna przeglądarki.
  cialo.querySelector('.dn-szyna-modulu-glowa .dn-btn')?.remove();
  cialo.querySelector('#panel-notes .dn-btn--atrament')?.remove();
  cialo.querySelector('#panel-notes .dn-zakladki')?.remove();
  cialo.querySelector('#panel-notes .sta-okno-tresc > .dn-btn')?.remove();
  const stopkaZrodel = [...cialo.querySelectorAll('#panel-sources .sta-okno-tresc > .sta-chip-rzad')];
  stopkaZrodel[stopkaZrodel.length - 1]?.remove();
}

/** Pyta rdzeń o karty, stronę bieżącą, źródła i notatki okna, po czym wypełnia nimi pasek adresu, płótno i panele. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const karty = await wywolaj(kanal, Command.BrowserTabList, { windowId: idOkna });
  wypelnijKarty(cialo, wzory, stan, karty.udany ? (karty.wynik?.tabs ?? []) : []);
  const migawka = await wywolaj(kanal, Command.BrowserSnapshotGet, { windowId: idOkna });
  opiszStrone(cialo, stan, migawka.udany ? migawka.wynik?.snapshot : undefined);
  await wypelnijZrodla(kanal, idOkna, cialo, wzory, stan, '');
  await wypelnijNotatki(kanal, idOkna, cialo, wzory, stan);
  await wypelnijSzyne(kanal, idOkna, cialo, wzory, stan);
}

/** Stawia w pasie kart po jednej karcie na kartę przeglądania; karta czynna dostaje wskazanie, a jej identyfikator wchodzi w stan okna. */
function wypelnijKarty(
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  karty: BrowserTab[],
): void {
  const tory = cialo.querySelector<HTMLElement>('.brw-zakladki-tory');
  if (tory === null || wzory.karta === null) return;
  tory.replaceChildren();
  for (const karta of karty) {
    const pozycja = wzory.karta.cloneNode(true) as HTMLElement;
    pozycja.dataset.karta = karta.id;
    const podpis = karta.title ?? karta.url ?? '';
    if (podpis === '' || !wpiszTekst(pozycja, podpis)) continue;
    if (karta.state === BrowserTabState.Active) {
      pozycja.setAttribute('aria-selected', 'true');
      stan.karta = karta.id;
    }
    tory.appendChild(pozycja);
  }
}

/** Nanosi migawkę strony na pasek adresu, znak kontekstu i płótno okna; brak migawki zostawia płótno puste i zdejmuje znak. */
function opiszStrone(
  cialo: HTMLElement,
  stan: Stan,
  migawka: BrowserSnapshot | undefined,
): void {
  stan.adres = migawka?.url ?? '';
  const adres = cialo.querySelector<HTMLInputElement>('.brw-adres input');
  if (adres !== null) {
    /* Adres przykładowy stoi w znaczniku atrybutem, więc sama własność pola go
       nie zdejmuje — atrybut idzie razem z wartością. */
    adres.value = stan.adres;
    adres.setAttribute('value', stan.adres);
  }
  opiszWezel(cialo.querySelector('.sta-kom-kontekst [data-pole="strona"]'), stan.adres);
  const plotno = cialo.querySelector<HTMLElement>('#panel-browser .brw-strona');
  if (plotno !== null) plotno.textContent = migawka?.text ?? '';
}

/** Wiąże pasek adresu: przejście pod adres wpisany, odświeżenie strony bieżącej i odłożenie jej do zakładek. */
function zwiazPasekAdresu(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  stan: Stan,
  odswiez: () => void,
): void {
  const adres = cialo.querySelector<HTMLInputElement>('.brw-adres input');
  adres?.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter') return;
    zdarzenie.preventDefault();
    void przejdz(kanal, idOkna, adres.value.trim(), odswiez);
  });
  const przyciski = [...cialo.querySelectorAll('.brw-pasek > .dn-btn-ikona')];
  przyciski[0]?.addEventListener('click', () => {
    void przejdz(kanal, idOkna, stan.adres, odswiez);
  });
  przyciski[1]?.addEventListener('click', () => {
    if (stan.adres === '') return;
    void wywolaj(kanal, Command.BrowserBookmarkAdd, { windowId: idOkna, url: stan.adres });
  });
}

/** Przechodzi pod wskazany adres i odświeża stanowisko; adres pusty wstrzymuje czynność, bo rdzeń nie ma dokąd pójść. */
async function przejdz(
  kanal: Kanal,
  idOkna: string,
  adres: string,
  odswiez: () => void,
): Promise<void> {
  if (adres === '') return;
  const wynik = await wywolaj(kanal, Command.BrowserNavigate, { windowId: idOkna, url: adres });
  if (wynik.udany) odswiez();
}

/** Wiąże pas kart: wybór karty czyni ją czynną, a przycisk pasa otwiera kartę nową. */
function zwiazKarty(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  stan: Stan,
  odswiez: () => void,
): void {
  cialo.querySelector('.brw-zakladki-tory')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.brw-tab')?.dataset.karta;
    if (karta === undefined || karta === stan.karta) return;
    void wywolaj(kanal, Command.BrowserTabUpdate, { tabId: karta, active: true }).then(odswiez);
  });
  cialo.querySelector('.brw-zakladki > .dn-btn-ikona')?.addEventListener('click', () => {
    void wywolaj(kanal, Command.BrowserTabOpen, { windowId: idOkna }).then(odswiez);
  });
}

/** Stawia w panelu źródeł po jednej pozycji na źródło okna i nanosi ich liczbę na znacznik panelu oraz znak pasa czynności. */
async function wypelnijZrodla(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  fraza: string,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-sources .sta-okno-tresc');
  if (panel === null) return;
  for (const stojace of panel.querySelectorAll('.brw-zrodlo-poz')) stojace.remove();
  const wynik = await wywolaj(kanal, Command.BrowserSourceList, {
    windowId: idOkna,
    query: fraza === '' ? undefined : fraza,
  });
  const zrodla = wynik.udany ? (wynik.wynik?.sources ?? []) : [];
  stan.zrodla.clear();
  for (const zrodlo of zrodla) stan.zrodla.set(zrodlo.id, zrodlo);
  const znacznik = cialo.querySelector<HTMLElement>('#panel-sources .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = String(zrodla.length);
  opiszZnak(cialo, 'zrodla', String(zrodla.length));
  if (wzory.zrodlo === null) return;
  for (const zrodlo of zrodla) {
    const pozycja = postawZrodlo(wzory, zrodlo);
    if (pozycja !== null) panel.appendChild(pozycja);
  }
}

/** Składa pozycję źródła z jego tytułu, godziny zebrania i adresu; źródło kluczowe dostaje gwiazdkę wypełnioną. */
function postawZrodlo(wzory: Wzory, zrodlo: BrowserSource): HTMLElement | null {
  if (wzory.zrodlo === null) return null;
  const pozycja = wzory.zrodlo.cloneNode(true) as HTMLElement;
  pozycja.dataset.zrodlo = zrodlo.id;
  const tytul = pozycja.querySelector<HTMLElement>('.brw-zrodlo-tyt');
  if (tytul === null) return null;
  if (!wpiszTekst(tytul, `${zrodlo.title ?? zrodlo.url} `)) return null;
  if (zrodlo.key === true && wzory.gwiazdka !== null) {
    tytul.querySelector('svg')?.replaceWith(wzory.gwiazdka.cloneNode(true));
  }
  const czas = tytul.querySelector('.dn-meta');
  if (czas !== null) czas.textContent = godzina(zrodlo.createdAt);
  const adres = pozycja.querySelector('.brw-zrodlo-meta');
  if (adres !== null) adres.textContent = zrodlo.url;
  return pozycja;
}

/** Nanosi wartość na znak pasa czynności rozpoznany po oznaczeniu. */
function opiszZnak(cialo: HTMLElement, pole: string, wartosc: string): void {
  opiszWezel(cialo.querySelector(`.sta-kontekst-akcji [data-pole="${pole}"]`), wartosc);
}

/**
 * Wiąże panel źródeł: pole zawężania szuka po tytule i adresie, przycisk
 * nagłówka odkłada stronę bieżącą jako źródło, a pozycja wykazu otwiera swój
 * adres.
 */
function zwiazZrodla(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-sources .sta-okno-tresc');
  if (panel === null) return;
  const pole = panel.querySelector<HTMLInputElement>('.dn-szukaj input');
  pole?.addEventListener('input', () => {
    void wypelnijZrodla(kanal, idOkna, cialo, wzory, stan, pole.value.trim());
  });
  cialo.querySelector('#panel-sources .sta-okno-akcje .dn-btn')?.addEventListener('click', () => {
    if (stan.adres === '') return;
    void wywolaj(kanal, Command.BrowserSourceAdd, {
      windowId: idOkna,
      url: stan.adres,
    }).then(odswiez);
  });
  panel.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('.sta-chip-rzad') === null) return;
    const zrodlo = cel.closest<HTMLElement>('.brw-zrodlo-poz')?.dataset.zrodlo;
    const adres = zrodlo === undefined ? undefined : stan.zrodla.get(zrodlo)?.url;
    if (adres === undefined) return;
    void przejdz(kanal, idOkna, adres, odswiez);
  });
}

/** Stawia w panelu notatek po jednej notatce na notatkę okna; notatka przypięta bierze wzór przypięcia. */
async function wypelnijNotatki(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-notes .sta-okno-tresc');
  if (panel === null) return;
  panel.replaceChildren();
  const wynik = await wywolaj(kanal, Command.BrowserNoteList, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const notatka of wynik.wynik.notes) {
    const wpis = postawNotatke(wzory, stan, notatka);
    if (wpis !== null) panel.appendChild(wpis);
  }
}

/**
 * Składa notatkę z jej klasyfikacji, źródła i cytatu. Notatka bez cytatu
 * pokazuje swoją treść — pole cytatu jest w kontrakcie osobne i puste zostaje
 * puste tylko wtedy, gdy treści też nie ma.
 */
function postawNotatke(wzory: Wzory, stan: Stan, notatka: BrowserNote): HTMLElement | null {
  const wzor = notatka.pinned === true ? wzory.notatkaPrzypieta : wzory.notatka;
  if (wzor === null) return null;
  const wpis = wzor.cloneNode(true) as HTMLElement;
  wpis.dataset.zrodlo = notatka.sourceId ?? '';
  const zrodlo = notatka.sourceId === undefined ? undefined : stan.zrodla.get(notatka.sourceId);
  const nazwaZrodla = zrodlo?.title ?? zrodlo?.url ?? '';
  const podpis = notatka.classification ?? '';
  const tytul = wpis.querySelector<HTMLElement>('.brw-zrodlo-tyt');
  if (tytul !== null) {
    tytul.querySelector('.dn-meta')?.remove();
    const opis = nazwaZrodla === '' ? podpis : `${podpis} · źródło: ${nazwaZrodla}`;
    if (!wpiszTekst(tytul, opis === '' ? ' ' : opis)) return null;
  }
  const cytat = wpis.querySelector('.brw-cytat');
  if (cytat !== null) cytat.textContent = notatka.quote ?? notatka.content;
  return wpis;
}

/**
 * Wiąże panel notatek: pozycja wykazu otwiera źródło notatki. Zmiana notatki
 * schodzi — komenda bierze nową treść, a panel pola do jej wpisania nie
 * niesie.
 */
function zwiazNotatki(kanal: Kanal, idOkna: string, cialo: HTMLElement, stan: Stan): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-notes .sta-okno-tresc');
  if (panel === null) return;
  panel.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('.sta-chip-rzad') === null) return;
    const zrodlo = cel.closest<HTMLElement>('.brw-notatka')?.dataset.zrodlo ?? '';
    const adres = stan.zrodla.get(zrodlo)?.url;
    if (adres === undefined) return;
    void przejdz(kanal, idOkna, adres, () => undefined);
  });
}

/** Stawia w szynie zestawy tematyczne źródeł wraz z ich składem; źródło poza zestawem stoi w panelu źródeł, nie w szynie. */
async function wypelnijSzyne(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null || wzory.grupa === null || wzory.pozycja === null) return;
  lista.replaceChildren();
  const wynik = await wywolaj(kanal, Command.BrowserSourceGroupList, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const zestaw of wynik.wynik.groups) {
    const grupa = wzory.grupa.cloneNode(true) as HTMLElement;
    grupa.textContent = zestaw.name;
    lista.appendChild(grupa);
    for (const kod of zestaw.sourceIds ?? []) {
      const zrodlo = stan.zrodla.get(kod);
      if (zrodlo === undefined) continue;
      const pozycja = wzory.pozycja.cloneNode(true) as HTMLElement;
      pozycja.dataset.adres = zrodlo.url;
      const tytul = pozycja.querySelector('.pt-pozycja-tytul');
      if (tytul === null) continue;
      tytul.textContent = zrodlo.title ?? zrodlo.url;
      lista.appendChild(pozycja);
    }
  }
}

/** Wiąże wybór pozycji szyny z przejściem pod adres źródła tego zestawu. */
function zwiazSzyne(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null) return;
  lista.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest<HTMLElement>('.pt-pozycja');
    const adres = pozycja?.dataset.adres;
    if (pozycja === null || pozycja === undefined || adres === undefined) return;
    for (const inna of lista.querySelectorAll('.pt-pozycja')) inna.removeAttribute('aria-current');
    pozycja.setAttribute('aria-current', 'true');
    void przejdz(kanal, idOkna, adres, odswiez);
  });
}

/** Nasłuchuje zmian strony i kart okna: płótno bierze migawkę po zmianie, a pas kart odczytuje się na nowo. */
function zwiazZdarzeniaStrony(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  stan: Stan,
): void {
  kanal.naZdarzenie(EventType.BrowserPageChanged, (tresc) => {
    if (tresc.windowId !== idOkna) return;
    opiszStrone(cialo, stan, tresc.snapshot);
  });
  kanal.naZdarzenie(EventType.BrowserTabChanged, (tresc) => {
    if (tresc.windowId !== idOkna || tresc.tab === undefined) return;
    const karta = cialo.querySelector<HTMLElement>(`.brw-tab[data-karta="${tresc.tab.id}"]`);
    if (karta === null) return;
    opiszWezel(karta, tresc.tab.title ?? tresc.tab.url ?? '');
    karta.setAttribute('aria-selected', String(tresc.tab.state === BrowserTabState.Active));
  });
}
