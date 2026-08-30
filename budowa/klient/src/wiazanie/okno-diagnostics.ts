/**
 * Wiązanie wnętrza okna modułu Diagnostics z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * diagnostics.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  Command,
  DiagnosticPriority,
  EventType,
  LogLevel,
  type DiagnosticAnalysis,
  type DiagnosticError,
  type LogEntry,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina diagnostics.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = [
  'panel-plan',
  'panel-zadania',
  'panel-terminal',
  'panel-kolejka',
  'panel-subagenci',
  'panel-artefakty',
];

/** Klasa kropki wagi błędu; nazwa klasy pochodzi ze znacznika Właściciela, a przypisanie z rejestru priorytetów kontraktu. */
const KROPKA_WAGI: Record<string, string> = {
  [DiagnosticPriority.Critical]: 'dn-kropka--blad',
  [DiagnosticPriority.High]: 'dn-kropka--blad',
  [DiagnosticPriority.Medium]: 'dn-kropka--ostrzezenie',
  [DiagnosticPriority.Low]: 'dn-kropka--neutralna',
};

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  wierszDziennika: HTMLElement | null;
  wierszBledu: HTMLElement | null;
  miaraBledu: HTMLElement | null;
  wierszStanu: HTMLElement | null;
  wierszRekomendacji: HTMLElement | null;
  miaraRekomendacji: HTMLElement | null;
  karta: HTMLElement | null;
}

/**
 * Wiąże wnętrze okna modułu Diagnostics. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazDiagnostics(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  zwiazFiltrDziennika(kanal, cialo, wzory);
  kanal.naZdarzenie(EventType.DiagnosticsAnalysisChanged, (tresc) => {
    if (tresc.analysis.windowId !== idOkna) return;
    opiszAnalize(cialo, wzory, tresc.analysis);
    void wypelnijRekomendacje(kanal, cialo, wzory, tresc.analysis.id);
  });

  void wypelnijDziennik(kanal, cialo, wzory, '');
  void wypelnijBledy(kanal, cialo, wzory);
  void wypelnijCentrum(kanal, idOkna, cialo, wzory);
}

/** Zdejmuje wzory wierszy, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const wierszBledu = sklonuj(
    cialo.querySelector<HTMLElement>('#panel-bledy .dn-wykaz-modulu-poz'),
  );
  // Barwa licznika wystąpień jest w prototypie barwą ostrzeżenia; błąd niesie
  // wagę osobno, więc licznik idzie plakietką bez wagi.
  wierszBledu?.querySelector('.dn-plakietka')?.setAttribute('class', 'dn-plakietka dn-na-koniec');
  const wierszStanu = sklonuj(
    cialo.querySelector<HTMLElement>('#panel-centrum .dn-wykaz-modulu-poz'),
  );
  wierszStanu?.querySelector('.dn-plakietka')?.remove();
  wierszStanu?.querySelector('.dn-meta')?.remove();
  const wierszRekomendacji = sklonuj(
    cialo.querySelector<HTMLElement>('#panel-rekomendacje .dn-wykaz-modulu-poz'),
  );

  return {
    wierszDziennika: sklonuj(cialo.querySelector<HTMLElement>('#panel-logi .dg-log .wiersz')),
    wierszBledu,
    miaraBledu: sklonuj(cialo.querySelector<HTMLElement>('#panel-bledy .dn-meta')),
    wierszStanu,
    wierszRekomendacji,
    miaraRekomendacji: sklonuj(cialo.querySelector<HTMLElement>('#panel-rekomendacje .dn-meta')),
    karta: sklonuj(cialo.querySelector<HTMLElement>('#panel-centrum .dg-kpi-karta')),
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów i historia wiszą na
 * rodzinach spoza diagnostics.*; materiał sesji, przypięty plik dziennika,
 * znak pracy, monitor korelacji i miara różnicy nie mają pola w kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.replaceChildren();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  cialo.querySelector('.sta-kom-kontekst')?.replaceChildren();
  cialo.querySelector('#okno-czat-1 .sta-okno-belka > .sta-chip')?.remove();
  // Znacznik źródła dziennika: wpis niesie własne źródło, lecz wykaz idzie
  // z wielu źródeł naraz, więc jedna nazwa nad panelem mówiłaby nieprawdę.
  cialo.querySelector('#panel-logi .sta-okno-znacznik')?.remove();
  // Znacznik „do przeglądu” opisuje stan przeglądu, którego rekomendacja nie
  // niesie jako całość — stan ma pojedyncza rekomendacja, nie panel.
  cialo.querySelector('#panel-rekomendacje .sta-okno-znacznik')?.remove();
  zdejmijPasCzynnosci(cialo);
}

/**
 * Zdejmuje pas czynności okna komunikacji. Zasięg wykonania, materiał sesji,
 * plik dziennika, drzewo robocze i miara różnicy nie mają pola w rodzinie
 * diagnostics.*, a zastosowania poprawki nie niesie żadna jej komenda.
 */
function zdejmijPasCzynnosci(cialo: HTMLElement): void {
  cialo.querySelector('.sta-kontekst-akcji')?.remove();
}

/**
 * Zdejmuje sterowanie, którego rodzina diagnostics.* nie obsługuje: panele bez
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

/** Zostawia w menu dodawania sam tytuł i trzy pozycje materiału; konektory i wtyczki należą do rodzin spoza diagnostics.*. */
function zdejmijKonektory(cialo: HTMLElement): void {
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Wypełnia Logs Viewer wpisami dziennika rdzenia; poziom pusty bierze wykaz bez zawężenia. */
async function wypelnijDziennik(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  poziom: string,
): Promise<void> {
  const dziennik = cialo.querySelector<HTMLElement>('#panel-logi .dg-log');
  if (dziennik === null) return;
  const wynik = await wywolaj(kanal, Command.DiagnosticsLogQuery, {
    level: poziom === '' ? undefined : (poziom as LogLevel),
    limit: 200,
  });
  dziennik.replaceChildren();
  if (wzory.wierszDziennika === null) return;
  for (const wpis of wynik.udany ? (wynik.wynik?.entries ?? []) : []) {
    dziennik.appendChild(zbudujWierszDziennika(wzory.wierszDziennika, wpis));
  }
}

/** Składa wiersz dziennika: czas, poziom i treść wpisu; klasa poziomu bierze się z jego nazwy w kontrakcie. */
function zbudujWierszDziennika(wzor: HTMLElement, wpis: LogEntry): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const czas = wiersz.querySelector('.czas');
  if (czas !== null) czas.textContent = godzina(wpis.timestamp);
  const poziom = wiersz.querySelector<HTMLElement>('.poz');
  if (poziom !== null) {
    poziom.className = `poz poz-${skrotPoziomu(wpis.level)}`;
    poziom.textContent = wpis.level.toUpperCase();
  }
  const tresc = [...wiersz.querySelectorAll('span')].at(-1);
  if (tresc !== undefined) tresc.textContent = wpis.message;
  return wiersz;
}

/** Skrót poziomu w nazwie klasy prototypu; poziom spoza trzech nazwanych idzie klasą wpisu zwykłego. */
function skrotPoziomu(poziom: LogLevel): string {
  if (poziom === LogLevel.Error) return 'err';
  if (poziom === LogLevel.Warn) return 'warn';
  return 'info';
}

/** Wiąże filtr poziomu Logs Viewer z zawężeniem zapytania; pierwszy znak pasa nie zawęża, bo prototyp nazywa go „Wszystkie”. */
function zwiazFiltrDziennika(kanal: Kanal, cialo: HTMLElement, wzory: Wzory): void {
  const pas = cialo.querySelector<HTMLElement>('#panel-logi .dg-filtr');
  if (pas === null) return;
  const znaki = [...pas.querySelectorAll<HTMLElement>('.chip')];
  pas.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const znak = cel.closest<HTMLElement>('.chip');
    if (znak === null) return;
    for (const inny of znaki) inny.setAttribute('aria-pressed', String(inny === znak));
    const nazwa = (znak.textContent ?? '').trim().toLowerCase();
    const poziom = Object.values(LogLevel).find((kandydat) => kandydat === nazwa) ?? '';
    void wypelnijDziennik(kanal, cialo, wzory, poziom);
  });
}

/** Wypełnia Errors Panel błędami zgrupowanymi po odcisku i nanosi ich liczbę na znacznik panelu. */
async function wypelnijBledy(kanal: Kanal, cialo: HTMLElement, wzory: Wzory): Promise<void> {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-bledy .dn-wykaz-modulu');
  if (wykaz === null) return;
  const wynik = await wywolaj(kanal, Command.DiagnosticsErrorList, { limit: 50 });
  const bledy = wynik.udany ? (wynik.wynik?.errors ?? []) : [];
  wykaz.replaceChildren();
  /* Znacznik idzie samą liczbą: prototyp niesie ją z rzeczownikiem w mianowniku
     mnogim, a odmiana zależy od liczby, której znacznik z góry nie zna. */
  const znacznik = cialo.querySelector('#panel-bledy .sta-okno-znacznik');
  if (znacznik !== null) znacznik.textContent = String(wynik.wynik?.total ?? bledy.length);
  for (const blad of bledy) {
    if (wzory.wierszBledu !== null) wykaz.appendChild(zbudujWierszBledu(wzory.wierszBledu, blad));
    if (wzory.miaraBledu !== null) {
      const miara = wzory.miaraBledu.cloneNode(true) as HTMLElement;
      miara.textContent = opisBledu(blad);
      wykaz.appendChild(miara);
    }
  }
}

/** Składa wiersz błędu: kropka wagi, treść i licznik wystąpień; błąd bez licznika traci plakietkę. */
function zbudujWierszBledu(wzor: HTMLElement, blad: DiagnosticError): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const kropka = wiersz.querySelector<HTMLElement>('.dn-kropka');
  if (kropka !== null) {
    kropka.className = `dn-kropka ${KROPKA_WAGI[blad.priority] ?? 'dn-kropka--neutralna'}`;
  }
  const tresc = wiersz.querySelector('b');
  if (tresc !== null) tresc.textContent = blad.message;
  const plakietka = wiersz.querySelector('.dn-plakietka');
  if (plakietka === null) return wiersz;
  if (blad.occurrences === undefined) plakietka.remove();
  else plakietka.textContent = `×${blad.occurrences}`;
  return wiersz;
}

/** Miara pod wierszem błędu: odcisk, źródło i kod odmowy; człony bez wartości nie wchodzą do zapisu. */
function opisBledu(blad: DiagnosticError): string {
  const czlony = [`odcisk ${blad.fingerprint}`, blad.source ?? '', blad.errorCode ?? ''];
  return czlony.filter((czlon) => czlon !== '').join(' · ');
}

/** Uruchamia analizę stanu systemu i wypełnia nią Diagnostics Center wraz z Recommendations Panel. */
async function wypelnijCentrum(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DiagnosticsAnalyzeRun, { windowId: idOkna });
  const analiza = wynik.udany ? wynik.wynik?.analysis : undefined;
  if (analiza === undefined) {
    cialo.querySelector('#panel-centrum')?.remove();
    cialo.querySelector('[data-panel-toggle="panel-centrum"]')?.remove();
    cialo.querySelector('#panel-rekomendacje')?.remove();
    cialo.querySelector('[data-panel-toggle="panel-rekomendacje"]')?.remove();
    return;
  }
  opiszAnalize(cialo, wzory, analiza);
  await wypelnijRekomendacje(kanal, cialo, wzory, analiza.id);
}

/**
 * Nanosi analizę na Diagnostics Center: karta miary bierze liczbę błędów
 * objętych analizą, a wykaz obrazu stanu — jej wynik zbiorczy. Miara zdrowia
 * i oś korelacji czasowej schodzą: analiza niesie zakres czasu, błędy,
 * rekomendacje i wynik zbiorczy, nic ponadto.
 */
function opiszAnalize(cialo: HTMLElement, wzory: Wzory, analiza: DiagnosticAnalysis): void {
  const miary = cialo.querySelector<HTMLElement>('#panel-centrum .dg-kpi');
  if (miary !== null && wzory.karta !== null) {
    miary.replaceChildren();
    miary.appendChild(zbudujKarte(wzory.karta, (analiza.errorIds ?? []).length));
  }
  const sekcje = [...cialo.querySelectorAll('#panel-centrum .dg-sekcja')];
  const wykazy = [...cialo.querySelectorAll('#panel-centrum .dn-wykaz-modulu')];
  for (const [numer, wykaz] of wykazy.entries()) {
    if (numer === 0) opiszObrazStanu(wykaz, wzory, analiza.summary ?? '');
    else {
      wykaz.remove();
      sekcje[numer]?.remove();
    }
  }
  if ((analiza.summary ?? '') === '') sekcje[0]?.remove();
}

/**
 * Karta miary z liczbą błędów objętych analizą. Z podpisu schodzi zakres czasu
 * podany w nawiasie: analiza niesie własny zakres, a ten w prototypie jest
 * godziną przykładową, nie zakresem odczytanym.
 */
function zbudujKarte(wzor: HTMLElement, ile: number): HTMLElement {
  const karta = wzor.cloneNode(true) as HTMLElement;
  const podpis = karta.querySelector('.etk');
  if (podpis !== null) podpis.textContent = (podpis.textContent ?? '').replace(/\s*\(.*\)/u, '');
  const wartosc = karta.querySelector('.wart');
  if (wartosc !== null) wartosc.textContent = String(ile);
  return karta;
}

/** Wpisuje wynik zbiorczy analizy w wykaz obrazu stanu; wynik pusty zdejmuje wykaz, bo pusty wykaz niczego nie mówi. */
function opiszObrazStanu(wykaz: Element, wzory: Wzory, wynik: string): void {
  if (wynik === '' || wzory.wierszStanu === null) {
    wykaz.remove();
    return;
  }
  const wiersz = wzory.wierszStanu.cloneNode(true) as HTMLElement;
  // Kropka niesie wagę stanu, której wynik zbiorczy analizy nie stopniuje.
  wiersz.querySelector('.dn-kropka')?.remove();
  if (!wpiszTekst(wiersz, wynik)) return;
  wykaz.replaceChildren(wiersz);
}

/**
 * Wypełnia Recommendations Panel rekomendacjami analizy. Przyciski przeglądu
 * schodzą: rodzina diagnostics.* nie niesie ani otwarcia poprawki w edytorze,
 * ani zmiany stanu rekomendacji.
 */
async function wypelnijRekomendacje(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  idAnalizy: string,
): Promise<void> {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-rekomendacje .dn-wykaz-modulu');
  if (wykaz === null) return;
  const wynik = await wywolaj(kanal, Command.DiagnosticsRecommendationList, {
    analysisId: idAnalizy,
  });
  const rekomendacje = wynik.udany ? (wynik.wynik?.recommendations ?? []) : [];
  wykaz.replaceChildren();
  for (const rekomendacja of rekomendacje) {
    if (wzory.wierszRekomendacji !== null) {
      const wiersz = wzory.wierszRekomendacji.cloneNode(true) as HTMLElement;
      const kropka = wiersz.querySelector<HTMLElement>('.dn-kropka');
      if (kropka !== null) {
        kropka.className =
          `dn-kropka ${KROPKA_WAGI[rekomendacja.priority] ?? 'dn-kropka--neutralna'}`;
      }
      if (wpiszTekst(wiersz, rekomendacja.title)) wykaz.appendChild(wiersz);
    }
    const opis = rekomendacja.detail ?? rekomendacja.targetPath ?? '';
    if (opis === '' || wzory.miaraRekomendacji === null) continue;
    const miara = wzory.miaraRekomendacji.cloneNode(true) as HTMLElement;
    miara.textContent = opis;
    wykaz.appendChild(miara);
  }
}

/** Godzina wpisu w zapisie zegara ściennego; czas kontraktu idzie w milisekundach epoki. */
function godzina(czas: number): string {
  return new Date(czas).toLocaleTimeString('pl-PL', { hour12: false });
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
function sklonuj<T extends Element>(wezel: T | null): T | null {
  return wezel === null ? null : (wezel.cloneNode(true) as T);
}
