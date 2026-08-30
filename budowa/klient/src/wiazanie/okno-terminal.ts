/**
 * Wiązanie wnętrza okna modułu Terminal z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * terminal.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  Command,
  EventType,
  TerminalOutputChannel,
  TerminalProcessStatus,
  TerminalSessionStatus,
  type TerminalOutputLine,
  type TerminalProcess,
  type TerminalSession,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina terminal.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-plan', 'panel-pliki'];

/** Ile milisekund odczyt wyjścia czeka na domknięcie procesu; komenda uruchomienia kończy się w chwili startu, więc bez czekania wyjście byłoby puste. */
const CZAS_ODCZYTU_WYJSCIA = 3000;

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  karta: HTMLElement | null;
  wierszWyjscia: HTMLElement | null;
  proces: HTMLElement | null;
  tetno: HTMLElement | null;
  ponow: HTMLElement | null;
  sciaga: HTMLElement | null;
}

/** Karta terminala wybrana w oknie wraz ze spisem kart rdzenia; obie wartości pochodzą z odpowiedzi rdzenia. */
interface Stan {
  karta: string;
  karty: Map<string, TerminalSession>;
}

/**
 * Wiąże wnętrze okna modułu Terminal. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazTerminal(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { karta: '', karty: new Map() };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt kart, procesów i wyjścia; wołane przy wejściu i po każdej czynności zmieniającej stan rdzenia. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, stan);
  };

  zwiazKarty(kanal, cialo, stan, odswiez);
  zwiazWiersz(kanal, cialo, stan, odswiez);
  zwiazProcesy(kanal, cialo, stan, odswiez);
  zwiazZdarzeniaProcesow(kanal, idOkna, odswiez);

  odswiez();
  void wypelnijSciage(kanal, cialo, wzory);
}

/** Zdejmuje wzory kart, wierszy i pozycji, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const karta = sklonuj(cialo.querySelector('#panel-tabs .tt-karta'));
  karta?.removeAttribute('data-powloka');
  const proces = sklonuj(cialo.querySelector('#panel-process .pm-wiersz'));
  // Wykres obciążenia nie ma pola: proces niesie zużycie procesora i pamięci
  // wartością chwilową, nie przebiegiem.
  proces?.querySelector('.spark')?.remove();
  proces?.querySelector('.dn-btn-ikona')?.remove();
  const wiersze = [...cialo.querySelectorAll<HTMLElement>('#panel-process .pm-wiersz')];
  return {
    karta,
    wierszWyjscia: sklonuj(cialo.querySelector('#panel-output .oc-wiersz')),
    proces,
    tetno: sklonuj(cialo.querySelector('#panel-process .pt-tetno')),
    ponow: sklonuj(wiersze[2]?.querySelector('button.dn-btn--zarys') ?? null),
    sciaga: sklonuj(cialo.querySelector('#menu-sciaga .sta-menu-poz[data-sciaga]')),
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz sesji szyny, historia rozmowy
 * i znaki pasa czynności wiszą na rodzinach spoza terminal.*, więc wiązanie
 * ich nie wypełnia; gałąź, profil powłoki i miary widoku nie mają pola
 * w całym kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  zdejmijPoleNaglowka(cialo, 'Wysiłek');
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  // Odpięcie znaku kontekstu nie ma komendy w żadnej z rodzin tego okna.
  for (const odepnij of cialo.querySelectorAll('.sta-kom-kontekst .odepnij')) odepnij.remove();
  // Zasięg wykonania, katalog projektu, gałąź, drzewo robocze, miara różnicy
  // i zatwierdzenie zmian należą do rodzin spoza terminal.*.
  cialo.querySelector('.sta-kontekst-akcji')?.remove();
  zdejmijOzdobyKonsoli(cialo);
  zdejmijSterowanieWyjscia(cialo);
  zdejmijSterowanieProcesow(cialo);
}

/** Zdejmuje pole nagłówka rozpoznane po podpisie; rodzina terminal.* nie oddaje nakładu rozumowania kanału. */
function zdejmijPoleNaglowka(cialo: HTMLElement, podpis: string): void {
  const pola = [...cialo.querySelectorAll('.sta-kom-naglowek .sta-kom-pole')];
  pola.find((pole) => pole.textContent?.startsWith(podpis) === true)?.remove();
}

/**
 * Czyści konsolę karty do samego wiersza wpisu. Znak zachęty powłoki, wyjście
 * przykładowe, gałąź pasa stanu i wybór profilu schodzą: karta terminala
 * niesie powłokę, katalog roboczy, stan i kod wyjścia, nic ponadto.
 */
function zdejmijOzdobyKonsoli(cialo: HTMLElement): void {
  const konsola = cialo.querySelector<HTMLElement>('#panel-tabs .tt-konsola');
  const wpis = konsola?.querySelector('.tt-wiersz');
  if (konsola !== null && konsola !== undefined) {
    konsola.replaceChildren();
    if (wpis !== null && wpis !== undefined) konsola.appendChild(wpis);
  }
  for (const zacheta of cialo.querySelectorAll('#panel-tabs [data-zacheta]')) zacheta.remove();
  cialo.querySelector('#panel-tabs .tt-stan .dn-plakietka')?.remove();
  cialo.querySelector('#panel-tabs .tt-stan span[style]')?.remove();
  // Nowa karta powłoki nie staje: `terminal.session.open` żąda rodzaju powłoki,
  // a okno nie niesie węzła, którym Operator by go wskazał.
  cialo.querySelector('#panel-tabs .tt-karty .tt-ntool')?.remove();
  zdejmijNarzedziaKarty(cialo);
}

/** Zostawia na pasie narzędzi karty sam znak stanu; podział karty, stopień pisma i schemat kolorów nie mają pola w kontrakcie. */
function zdejmijNarzedziaKarty(cialo: HTMLElement): void {
  const pas = cialo.querySelector<HTMLElement>('#panel-tabs .tt-narzedzia');
  if (pas === null) return;
  const stan = pas.lastElementChild;
  pas.replaceChildren();
  if (stan !== null) pas.appendChild(stan);
}

/**
 * Zdejmuje sterowanie zbiorczego wyjścia. Widok, poziom, wyszukiwanie
 * wyrażeniem, zawijanie, znaczniki czasu, przewijanie i wydanie logu nie mają
 * pola: wiersz wyjścia niesie kartę, proces, strumień, treść i czas.
 */
function zdejmijSterowanieWyjscia(cialo: HTMLElement): void {
  const tresc = cialo.querySelector<HTMLElement>('#panel-output .sta-okno-tresc');
  if (tresc === null) return;
  tresc.firstElementChild?.remove();
  tresc.lastElementChild?.remove();
  cialo.querySelector('#panel-output .oc-akcje')?.remove();
}

/** Zdejmuje zawężanie i wydanie migawki monitora procesów; `terminal.process.list` zawęża stanem i inicjatorem, a okno nie niesie węzła wyboru. */
function zdejmijSterowanieProcesow(cialo: HTMLElement): void {
  const tresc = cialo.querySelector<HTMLElement>('#panel-process .sta-okno-tresc');
  if (tresc === null) return;
  tresc.firstElementChild?.remove();
  tresc.lastElementChild?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina terminal.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, ustawienia sesji, urządzenia wejścia
 * dźwięku, wybór modelu i nakładu oraz menu dodawania.
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
  for (const pozycja of cialo.querySelectorAll('#menu-sesja .sta-menu-poz[role="menuitem"]')) {
    pozycja.remove();
  }
  cialo.querySelector('#menu-sesja .sta-menu-sep')?.remove();
  cialo.querySelector('#konfig-czat')?.remove();
  cialo.querySelector('[data-konfig-okno="konfig-czat"]')?.remove();
  cialo.querySelector('#pop-mik')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-model')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-wysilek')?.closest('.sta-nrz')?.remove();
  // Dołączenie pliku logu żąda ścieżki, a polecenia ukośnikowe własnego
  // rejestru; ani jednego, ani drugiego okno nie niesie.
  cialo.querySelector('#pop-plus')?.closest('.sta-nrz')?.remove();
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
  cialo.querySelector('#menu-sciaga .sta-menu-sep')?.nextElementSibling?.remove();
  cialo.querySelector('#menu-sciaga .sta-menu-sep')?.remove();
}

/** Pyta rdzeń o karty powłoki, procesy rejestru i zbiorcze wyjście okna, po czym wypełnia nimi pasek kart, monitor i konsolę. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const karty = await wywolaj(kanal, Command.TerminalSessionList, { windowId: idOkna });
  const wykaz = karty.udany ? (karty.wynik?.sessions ?? []) : [];
  stan.karty = new Map(wykaz.map((karta) => [karta.id, karta]));
  if (!stan.karty.has(stan.karta)) stan.karta = wykaz[0]?.id ?? '';
  wypelnijKarty(cialo, wzory, wykaz, stan);
  opiszKarte(cialo, stan.karty.get(stan.karta));
  await wypelnijProcesy(kanal, idOkna, cialo, wzory, stan);
  await wypelnijWyjscie(kanal, idOkna, cialo, wzory, stan);
}

/** Stawia w pasku po jednej karcie na kartę powłoki rdzenia; stan karty rozstrzyga odmianę kropki, a wybór — oznaczenie karty bieżącej. */
function wypelnijKarty(
  cialo: HTMLElement,
  wzory: Wzory,
  karty: TerminalSession[],
  stan: Stan,
): void {
  const tory = cialo.querySelector<HTMLElement>('#panel-tabs .tt-karty-tory');
  if (tory === null) return;
  tory.replaceChildren();
  if (wzory.karta === null) return;
  for (const sesja of karty) {
    const karta = wzory.karta.cloneNode(true) as HTMLElement;
    karta.dataset.karta = sesja.id;
    karta.setAttribute('aria-selected', String(sesja.id === stan.karta));
    oznaczKropke(karta, odmianaKarty(sesja.status));
    if (!wpiszTekst(karta, sesja.title ?? sesja.shell)) continue;
    tory.appendChild(karta);
  }
}

/** Nazwa odmiany kropki dla stanu karty powłoki; nazwy pochodzą z biblioteki, rozróżnienie ze stanu karty w kontrakcie. */
function odmianaKarty(status: TerminalSessionStatus): string {
  return status === TerminalSessionStatus.Running ? 'dn-kropka--sygnal' : 'dn-kropka--neutralna';
}

/** Podmienia odmianę kropki stanu w klonie, zostawiając klasę podstawową nadaną przez bibliotekę. */
function oznaczKropke(wezel: HTMLElement, odmiana: string): void {
  const kropka = wezel.querySelector<HTMLElement>('.dn-kropka');
  if (kropka === null) return;
  for (const klasa of [...kropka.classList]) {
    if (klasa.startsWith('dn-kropka--')) kropka.classList.remove(klasa);
  }
  kropka.classList.add(odmiana);
}

/**
 * Opisuje kartą bieżącą belkę okna komunikacji, znaki kontekstu, pas stanu
 * karty i profil powłoki. Brak karty zdejmuje każdy z tych węzłów — karta
 * jest jedynym źródłem powłoki i katalogu roboczego.
 */
function opiszKarte(cialo: HTMLElement, sesja: TerminalSession | undefined): void {
  const nazwa = sesja?.title ?? sesja?.shell ?? '';
  const katalog = sesja?.workingDir ?? '';
  opiszWezel(cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip'), nazwa);
  const znaki = [...cialo.querySelectorAll<HTMLElement>('.sta-kom-kontekst > .sta-zrodlo')];
  opiszWezel(znaki[0] ?? null, nazwa);
  opiszWezel(znaki[1] ?? null, katalog);
  const sciezka = cialo.querySelector<HTMLElement>('#panel-tabs .tt-breadcrumb');
  if (sciezka !== null) sciezka.textContent = katalog;
  const stanKarty = cialo.querySelector<HTMLElement>('#panel-tabs .tt-narzedzia span');
  if (stanKarty !== null) stanKarty.textContent = sesja?.status ?? '';
  opiszProfil(cialo, sesja);
}

/** Wypełnia profil powłoki karty bieżącej; schemat kolorów schodzi, bo karta niesie powłokę i katalog roboczy, nie wygląd konsoli. */
function opiszProfil(cialo: HTMLElement, sesja: TerminalSession | undefined): void {
  const pola = [...cialo.querySelectorAll<HTMLElement>('#konfig-tabs .sta-konfig-pole')];
  const wartosci = [sesja?.shell ?? '', sesja?.workingDir ?? ''];
  for (const [numer, pole] of pola.entries()) {
    const wartosc = wartosci[numer];
    const miejsce = pole.querySelector('.pt-mono');
    if (wartosc === undefined || miejsce === null) pole.remove();
    else miejsce.textContent = wartosc;
  }
}

/** Wypełnia monitor procesów wierszami rejestru rdzenia i nanosi proces czynny na monitor okna komunikacji. */
async function wypelnijProcesy(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalProcessList, { windowId: idOkna });
  const procesy = wynik.udany ? (wynik.wynik?.processes ?? []) : [];
  const panel = cialo.querySelector<HTMLElement>('#panel-process .sta-okno-tresc');
  const wzor = wzory.proces;
  if (panel !== null && wzor !== null) {
    panel.replaceChildren();
    for (const proces of procesy) panel.appendChild(zbudujProces(wzor, wzory, proces, stan));
  }
  opiszMonitor(cialo, procesy);
}

/** Zwraca wiersz monitora opisany procesem rejestru; czynność wiersza rozstrzyga stan procesu, bo zakończonego nie ma czego kończyć. */
function zbudujProces(
  wzor: HTMLElement,
  wzory: Wzory,
  proces: TerminalProcess,
  stan: Stan,
): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.proces = proces.id;
  wiersz.dataset.polecenie = proces.command;
  const czynny = proces.status === TerminalProcessStatus.Running;
  if (czynny && wzory.tetno !== null) {
    wiersz.querySelector('.dn-kropka')?.replaceWith(wzory.tetno.cloneNode(true));
  } else {
    oznaczKropke(wiersz, odmianaProcesu(proces.status));
  }
  const polecenie = wiersz.querySelector('.pm-cmd');
  if (polecenie !== null) polecenie.textContent = proces.command;
  const plakietka = wiersz.querySelector('.dn-plakietka');
  if (plakietka !== null) plakietka.textContent = proces.initiator;
  opiszMiary(wiersz, miaryProcesu(proces, stan));
  opiszCzynnoscProcesu(wiersz, wzory, czynny);
  return wiersz;
}

/** Nazwa odmiany kropki dla stanu procesu; rozróżnienie bierze się z rejestru stanów kontraktu, nie ze skrótu tego pliku. */
function odmianaProcesu(status: TerminalProcessStatus): string {
  if (status === TerminalProcessStatus.Finished) return 'dn-kropka--sukces';
  if (status === TerminalProcessStatus.Failed) return 'dn-kropka--blad';
  return 'dn-kropka--neutralna';
}

/** Miary procesu w kolejności pasa: powłoka karty, obciążenie procesora, zajęta pamięć i kod wyjścia; wartość nieobecna miary nie stawia. */
function miaryProcesu(proces: TerminalProcess, stan: Stan): string[] {
  const miary: string[] = [];
  const karta = proces.sessionId === undefined ? undefined : stan.karty.get(proces.sessionId);
  if (karta !== undefined) miary.push(karta.shell);
  if (proces.cpuPercent !== undefined) miary.push(`${proces.cpuPercent}%`);
  if (proces.memoryBytes !== undefined) miary.push(`${Math.round(proces.memoryBytes / 1024 / 1024)} MB`);
  if (proces.exitCode !== undefined) miary.push(`kod wyjścia: ${proces.exitCode}`);
  return miary;
}

/** Wpisuje w pas miar procesu kolejne wartości odpowiedzi; pole bez wartości schodzi, bo pole puste niczego nie mówi. */
function opiszMiary(wiersz: HTMLElement, miary: string[]): void {
  const pas = wiersz.querySelector('.pm-meta');
  if (pas === null) return;
  const pola = [...pas.querySelectorAll('span')];
  for (const [numer, pole] of pola.entries()) {
    const miara = miary[numer];
    if (miara === undefined) pole.remove();
    else pole.textContent = miara;
  }
  if (miary.length === 0) pas.remove();
}

/** Stawia w wierszu monitora czynność właściwą stanowi: proces czynny da się zakończyć, zakończony — uruchomić ponownie. */
function opiszCzynnoscProcesu(wiersz: HTMLElement, wzory: Wzory, czynny: boolean): void {
  const pas = wiersz.lastElementChild;
  if (pas === null) return;
  if (czynny) return;
  pas.replaceChildren();
  if (wzory.ponow === null) pas.remove();
  else pas.appendChild(wzory.ponow.cloneNode(true));
}

/** Nanosi na monitor okna komunikacji polecenie procesu czynnego; brak procesu czynnego monitor zdejmuje. */
function opiszMonitor(cialo: HTMLElement, procesy: TerminalProcess[]): void {
  const monitor = cialo.querySelector<HTMLElement>('.sta-kom-monitor');
  if (monitor === null) return;
  const czynny = procesy.find((proces) => proces.status === TerminalProcessStatus.Running);
  if (czynny === undefined || !wpiszTekst(monitor, ` ${czynny.command}`)) monitor.remove();
}

/** Zapisuje okno na zbiorcze wyjście kart i wypełnia jego ogonem konsolę wyjścia oraz konsolę karty bieżącej. */
async function wypelnijWyjscie(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalOutputStream, { windowId: idOkna });
  const wiersze = wynik.udany ? (wynik.wynik?.lines ?? []) : [];
  const konsola = cialo.querySelector<HTMLElement>('#panel-output .sta-okno-tresc > div');
  const wzor = wzory.wierszWyjscia;
  if (konsola !== null && wzor !== null) {
    konsola.replaceChildren();
    for (const linia of wiersze) konsola.appendChild(zbudujWierszWyjscia(wzor, linia, stan));
  }
  const wlasne = wiersze
    .filter((linia) => linia.terminalSessionId === stan.karta)
    .map((linia) => linia.text);
  wpiszKonsole(cialo, wlasne.join('\n'));
}

/** Zwraca wiersz zbiorczego wyjścia opisany kartą, czasem wypisania i treścią; wiersz strumienia błędów dostaje oznaczenie biblioteki. */
function zbudujWierszWyjscia(
  wzor: HTMLElement,
  linia: TerminalOutputLine,
  stan: Stan,
): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const karta = stan.karty.get(linia.terminalSessionId);
  const zrodlo = wiersz.querySelector('.zrodlo');
  if (zrodlo !== null) zrodlo.textContent = karta?.title ?? karta?.shell ?? linia.terminalSessionId;
  const czas = wiersz.querySelector('.czas');
  if (czas !== null) czas.textContent = new Date(linia.at).toLocaleTimeString();
  const tresc = wiersz.querySelectorAll('span')[2];
  if (tresc !== undefined) tresc.textContent = linia.text;
  wiersz.classList.toggle('blad', linia.channel === TerminalOutputChannel.Stderr);
  return wiersz;
}

/** Wpisuje treść w konsolę karty bieżącej, zostawiając wiersz wpisu na końcu; konsola jest wykazem wierszy rdzenia, nie zapisem powłoki. */
function wpiszKonsole(cialo: HTMLElement, tresc: string): void {
  const konsola = cialo.querySelector<HTMLElement>('#panel-tabs .tt-konsola');
  if (konsola === null) return;
  const wpis = konsola.querySelector('.tt-wiersz');
  konsola.replaceChildren();
  if (tresc !== '') konsola.appendChild(document.createTextNode(`${tresc}\n`));
  if (wpis !== null) konsola.appendChild(wpis);
}

/** Dopisuje treść do konsoli karty bieżącej przed wierszem wpisu; przyrost idzie po odczycie wyjścia uruchomionego procesu. */
function dopiszKonsole(cialo: HTMLElement, tresc: string): void {
  const konsola = cialo.querySelector<HTMLElement>('#panel-tabs .tt-konsola');
  if (konsola === null || tresc === '') return;
  const wpis = konsola.querySelector('.tt-wiersz');
  const wezel = document.createTextNode(`${tresc}\n`);
  if (wpis === null) konsola.appendChild(wezel);
  else konsola.insertBefore(wezel, wpis);
}

/** Wypełnia szufladę ściągi pozycjami biblioteki skryptów; nagłówek szuflady bierze ich liczbę, a pozycja treść do wstawienia w wiersz. */
async function wypelnijSciage(kanal: Kanal, cialo: HTMLElement, wzory: Wzory): Promise<void> {
  const menu = cialo.querySelector<HTMLElement>('#menu-sciaga');
  if (menu === null) return;
  const wynik = await wywolaj(kanal, Command.TerminalScriptList, {});
  const pozycje = wynik.udany ? (wynik.wynik?.scripts ?? []) : [];
  const naglowek = menu.querySelector('.sta-menu-etyk');
  if (naglowek !== null) {
    naglowek.textContent = (naglowek.textContent ?? '').replace(/\(\d+\)/u, `(${pozycje.length})`);
  }
  for (const stojaca of menu.querySelectorAll('.sta-menu-poz')) stojaca.remove();
  if (wzory.sciaga === null) return;
  for (const skrypt of pozycje) {
    const pozycja = wzory.sciaga.cloneNode(true) as HTMLElement;
    pozycja.dataset.sciaga = skrypt.content;
    if (!wpiszTekst(pozycja, skrypt.name)) continue;
    menu.appendChild(pozycja);
  }
}

/** Wiąże pasek kart: wskazanie karty zmienia kartę bieżącą, a znak zamknięcia zamyka kartę powłoki w rdzeniu. */
function zwiazKarty(kanal: Kanal, cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  const tory = cialo.querySelector<HTMLElement>('#panel-tabs .tt-karty-tory');
  if (tory === null) return;
  tory.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.tt-karta')?.dataset.karta;
    if (karta === undefined) return;
    if (cel.classList.contains('x')) {
      void zamknijKarte(kanal, karta, odswiez);
      return;
    }
    stan.karta = karta;
    odswiez();
  });
}

/** Zamyka kartę powłoki i odświeża stanowisko, gdy rdzeń zamknięcie przyjął. */
async function zamknijKarte(kanal: Kanal, idKarty: string, odswiez: () => void): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalSessionClose, { sessionId: idKarty });
  if (!wynik.udany) return;
  odswiez();
}

/**
 * Wiąże wiersz poleceń karty bieżącej. Enter uruchamia polecenie w karcie,
 * a wyjście dochodzi osobnym odczytem: komenda uruchomienia kończy się
 * w chwili startu procesu, więc wyniku nieść nie może.
 */
function zwiazWiersz(kanal: Kanal, cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  const wpis = cialo.querySelector<HTMLInputElement>('#panel-tabs input[data-wpis]');
  if (wpis === null) return;
  wpis.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter') return;
    zdarzenie.preventDefault();
    const polecenie = wpis.value.trim();
    if (polecenie === '' || stan.karta === '') return;
    wpis.value = '';
    void uruchomPolecenie(kanal, cialo, stan.karta, polecenie, odswiez);
  });
  const menu = cialo.querySelector<HTMLElement>('#menu-sciaga');
  menu?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const tresc = cel.closest<HTMLElement>('.sta-menu-poz')?.dataset.sciaga;
    if (tresc === undefined) return;
    wpis.value = tresc;
    wpis.focus();
  });
}

/** Uruchamia polecenie w karcie, dopisuje jego wyjście do konsoli i odświeża monitor procesów. */
async function uruchomPolecenie(
  kanal: Kanal,
  cialo: HTMLElement,
  idKarty: string,
  polecenie: string,
  odswiez: () => void,
): Promise<void> {
  const start = await wywolaj(kanal, Command.TerminalCommandExec, {
    sessionId: idKarty,
    command: polecenie,
  });
  if (!start.udany || start.wynik === undefined) return;
  odswiez();
  const wyjscie = await wywolaj(kanal, Command.TerminalOutputRead, {
    processId: start.wynik.process.id,
    waitMs: CZAS_ODCZYTU_WYJSCIA,
  });
  if (!wyjscie.udany || wyjscie.wynik === undefined) return;
  dopiszKonsole(cialo, [wyjscie.wynik.stdout, wyjscie.wynik.stderr].filter((tekst) => tekst !== '').join('\n'));
  odswiez();
}

/** Wiąże czynności monitora: zakończenie procesu idzie sygnałem łagodnym, a powtórzenie uruchamia to samo polecenie w tej samej karcie. */
function zwiazProcesy(kanal: Kanal, cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-process .sta-okno-tresc');
  if (panel === null) return;
  panel.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wiersz = cel.closest<HTMLElement>('.pm-wiersz');
    const proces = wiersz?.dataset.proces;
    if (wiersz === null || wiersz === undefined || proces === undefined) return;
    if (cel.closest('button.dn-btn--niebezpieczny') !== null) {
      void zakonczProces(kanal, proces, odswiez);
      return;
    }
    if (cel.closest('button.dn-btn--zarys') === null) return;
    const polecenie = wiersz.dataset.polecenie ?? '';
    if (polecenie === '' || stan.karta === '') return;
    void uruchomPolecenie(kanal, cialo, stan.karta, polecenie, odswiez);
  });
}

/** Kończy proces rejestru sygnałem łagodnym i odświeża stanowisko, gdy rdzeń zakończenie przyjął. */
async function zakonczProces(kanal: Kanal, idProcesu: string, odswiez: () => void): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalProcessKill, { processId: idProcesu });
  if (!wynik.udany) return;
  odswiez();
}

/**
 * Nasłuchuje zmian rejestru procesów okna. Zmiana procesu przerysowuje całe
 * stanowisko, bo obszar nie ma zdarzenia ani dla kart powłoki, ani dla
 * zbiorczego wyjścia, a proces zmienia i jedno, i drugie.
 */
function zwiazZdarzeniaProcesow(kanal: Kanal, idOkna: string, odswiez: () => void): void {
  kanal.naZdarzenie(EventType.TerminalProcessChanged, (tresc) => {
    if (tresc.process.windowId !== idOkna) return;
    odswiez();
  });
}

/**
 * Wpisuje wartość w węzeł tekstowy znacznika. Wartość pusta zostawia węzeł
 * pustym, a nie zdejmuje go: wartość dochodzi dopiero z pracą podjętą w oknie,
 * więc znak na nią czeka. Węzeł raz opróżniony pozostaje celem kolejnych
 * wpisów, bo wybór pada na ostatni węzeł tekstowy, gdy niepustego już nie ma.
 */
function opiszWezel(wezel: HTMLElement | null, wartosc: string): void {
  if (wezel === null) return;
  const wezly = [...wezel.childNodes].filter((dziecko) => dziecko.nodeType === Node.TEXT_NODE);
  const cel = wezly.find((dziecko) => (dziecko.nodeValue ?? '').trim() !== '')
    ?? wezly[wezly.length - 1];
  if (cel !== undefined) cel.nodeValue = wartosc;
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
