/**
 * Wiązanie wnętrza okna modułu Assistant z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * assistant.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  AssistantActionControl,
  AssistantActionStatus,
  Command,
  EventType,
  type AssistantAction,
  type AssistantActivityEntry,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Panele prototypu, którym rodzina assistant.* nie oddaje ani jednego pola; schodzą razem z przełącznikami menu sesji. */
const PANELE_BEZ_POKRYCIA = ['panel-zadania', 'panel-artefakty', 'panel-plan', 'panel-pliki'];

/** Stany zlecenia, które wolno wstrzymać albo odwołać; pozostałe nie mają czego przyjąć od Actions Monitor. */
const STANY_STEROWALNE: string[] = [
  AssistantActionStatus.Queued,
  AssistantActionStatus.Running,
  AssistantActionStatus.Paused,
];

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  karta: HTMLElement | null;
  krok: HTMLElement | null;
  ponow: HTMLElement | null;
  wpisDziennika: HTMLElement | null;
}

/**
 * Wiąże wnętrze okna modułu Assistant. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazAssistant(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt zleceń i dziennika; wołane przy wejściu, po poleceniu i po każdej zmianie ogłoszonej przez rdzeń. */
  const odswiez = (): void => {
    void wypelnijZlecenia(kanal, idOkna, cialo, wzory);
    void wypelnijDziennik(kanal, idOkna, cialo, wzory, '');
  };

  zwiazSterowanieZleceniem(kanal, idOkna, cialo, odswiez);
  zwiazSzukanieWDzienniku(kanal, idOkna, cialo, wzory);
  zwiazPolecenie(kanal, idOkna, cialo, odswiez);
  kanal.naZdarzenie(EventType.AssistantActionChanged, (tresc) => {
    if (tresc.action.windowId !== idOkna) return;
    void wypelnijZlecenia(kanal, idOkna, cialo, wzory);
  });

  odswiez();
}

/** Zdejmuje wzory kart i wpisów, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const karta = sklonuj(cialo.querySelector<HTMLElement>('#panel-actions .am-karta'));
  // Wykaz kroków zlecenia jest treścią przykładową: zlecenie niesie numer
  // etapu i ich liczbę, nie zapis każdego kroku z osobna.
  for (const krok of [...(karta?.querySelectorAll('.am-krok') ?? [])]) krok.remove();

  return {
    karta,
    krok: sklonuj(cialo.querySelector<HTMLElement>('#panel-actions .am-karta--blad .am-krok')),
    ponow: sklonuj(
      cialo.querySelector<HTMLElement>('#panel-actions .am-karta--blad .sta-chip-rzad button'),
    ),
    wpisDziennika: zdejmijWzorWpisu(cialo),
  };
}

/** Zdejmuje wzór wpisu dziennika wraz z jego oczyszczeniem; wynik działania nie ma pola w kontrakcie wpisu. */
function zdejmijWzorWpisu(cialo: HTMLElement): HTMLElement | null {
  const wpis = sklonuj(cialo.querySelector<HTMLElement>('#panel-activity .af-poz'));
  wpis?.querySelector('.wynik')?.remove();
  return wpis;
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów i historia wiszą na
 * rodzinach spoza assistant.*; profil, język, pewność rozpoznania, znak
 * nasłuchu, monitor syntezy i pas czynności nie mają pola w kontrakcie.
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
  cialo.querySelector('.sta-kontekst-akcji')?.remove();
  zdejmijOpisyKonsoli(cialo);
  zdejmijPasyDziennika(cialo);
}

/**
 * Zdejmuje z Voice Console podpisy bez pokrycia: profil, język i pewność
 * rozpoznania nie wracają z rdzenia żadną komendą, transkrypcja stoi w oknie
 * przykładowa, a szybkie akcje nie mają za sobą ani jednej komendy.
 */
function zdejmijOpisyKonsoli(cialo: HTMLElement): void {
  cialo.querySelector('#panel-voice .vc-nag')?.remove();
  cialo.querySelector('#panel-voice .vc-transkrypcja')?.remove();
  cialo.querySelector('#panel-voice .vc-akcje')?.remove();
  cialo.querySelector('#panel-voice .sta-okno-belka > .dn-plakietka')?.remove();
}

/**
 * Zdejmuje z Activity Feed nagłówki dni oraz czynności bez komendy: odtworzenia
 * przebiegu, odsłuchu nagrania i wywiedzenia dziennika rodzina assistant.* nie
 * niesie.
 */
function zdejmijPasyDziennika(cialo: HTMLElement): void {
  for (const grupa of cialo.querySelectorAll('#panel-activity .af-grupa')) grupa.remove();
  for (const rzad of cialo.querySelectorAll('#panel-activity .sta-chip-rzad')) rzad.remove();
  for (const przycisk of cialo.querySelectorAll('#panel-activity > .sta-okno-tresc > .dn-btn')) {
    przycisk.remove();
  }
}

/**
 * Zdejmuje sterowanie, którego rodzina assistant.* nie obsługuje: panele bez
 * pokrycia wraz z ich przełącznikami, wybór modelu i nakładu, urządzenia
 * wejścia dźwięku, nastawy głosu Voice Console, filtr stanu zleceń oraz
 * konektory i wtyczki menu dodawania.
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
  // Głos, tempo i fraza wybudzająca to nastawa profilu spoza rodziny
  // assistant.*; stan zlecenia nie ma w niej pola zawężającego wykaz.
  cialo.querySelector('#menu-vc')?.closest('.sta-menu')?.remove();
  cialo.querySelector('#menu-am-filtr')?.closest('.sta-menu')?.remove();
  zdejmijKonektory(cialo);
}

/** Zostawia w menu dodawania sam tytuł i trzy pozycje materiału; konektory i wtyczki należą do rodzin spoza assistant.*. */
function zdejmijKonektory(cialo: HTMLElement): void {
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Wypełnia Actions Monitor zleceniami okna; brak zleceń zostawia panel pusty, bo pusty wykaz jest prawdą o oknie bez pracy. */
async function wypelnijZlecenia(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): Promise<void> {
  const wykaz = cialo.querySelector<HTMLElement>('#panel-actions .dn-wykaz-modulu');
  if (wykaz === null) return;
  const wynik = await wywolaj(kanal, Command.AssistantActionStatus, { windowId: idOkna });
  const zlecenia = wynik.udany ? (wynik.wynik?.actions ?? []) : [];
  wykaz.replaceChildren();
  const wzor = wzory.karta;
  if (wzor === null) return;
  for (const zlecenie of zlecenia) wykaz.appendChild(zbudujKarte(wzor, wzory, zlecenie));
}

/** Składa kartę zlecenia: nazwa, etap, postęp, wynik i sterowanie właściwe dla stanu. */
function zbudujKarte(wzor: HTMLElement, wzory: Wzory, zlecenie: AssistantAction): HTMLElement {
  const karta = wzor.cloneNode(true) as HTMLElement;
  karta.dataset.zlecenie = zlecenie.id;
  karta.classList.toggle('am-karta--blad', zlecenie.status === AssistantActionStatus.Failed);
  opiszTytulZlecenia(karta, zlecenie);
  opiszPostep(karta, zlecenie);
  opiszWynik(karta, wzory, zlecenie);
  opiszSterowanie(karta, wzory, zlecenie);
  return karta;
}

/** Nanosi na kartę nazwę zlecenia i jego etap; znak pracy zostaje wyłącznie przy zleceniu trwającym, a etap bez liczby kroków schodzi. */
function opiszTytulZlecenia(karta: HTMLElement, zlecenie: AssistantAction): void {
  const tytul = karta.querySelector<HTMLElement>('.am-tyt');
  if (tytul === null) return;
  if (zlecenie.status !== AssistantActionStatus.Running) tytul.querySelector('.pt-tetno')?.remove();
  if (!wpiszTekst(tytul, ` ${zlecenie.title ?? zlecenie.id} `)) return;
  const etap = tytul.querySelector('.dn-meta');
  if (etap === null) return;
  if (zlecenie.currentStep === undefined || zlecenie.totalSteps === undefined) etap.remove();
  else etap.textContent = (etap.textContent ?? '')
    .replace(/\d+\s*\/\s*\d+/u, `${zlecenie.currentStep}/${zlecenie.totalSteps}`);
}

/** Nanosi na pasek postępu udział etapu bieżącego w liczbie etapów; zlecenie bez liczby etapów zdejmuje pasek. */
function opiszPostep(karta: HTMLElement, zlecenie: AssistantAction): void {
  const pasek = karta.querySelector<HTMLElement>('.dn-postep');
  if (pasek === null) return;
  const etapy = zlecenie.totalSteps ?? 0;
  const wartosc = pasek.querySelector<HTMLElement>('.dn-postep-wartosc');
  if (etapy === 0 || zlecenie.currentStep === undefined || wartosc === null) {
    pasek.remove();
    return;
  }
  const udzial = Math.max(0, Math.min(100, (zlecenie.currentStep / etapy) * 100));
  pasek.setAttribute('data-postep-do', String(udzial));
  wartosc.style.width = `${udzial}%`;
}

/** Dokłada do karty wynik zlecenia; zlecenie bez wyniku nie dostaje wiersza, bo wiersz pusty niczego nie mówi. */
function opiszWynik(karta: HTMLElement, wzory: Wzory, zlecenie: AssistantAction): void {
  const tresc = zlecenie.result ?? zlecenie.confirmationReason ?? '';
  if (tresc === '' || wzory.krok === null) return;
  const krok = wzory.krok.cloneNode(true) as HTMLElement;
  krok.textContent = tresc;
  const rzad = karta.querySelector('.sta-chip-rzad');
  if (rzad === null) karta.appendChild(krok);
  else rzad.before(krok);
}

/**
 * Dobiera sterowanie karty do stanu zlecenia: zlecenie w toku wolno wstrzymać
 * i odwołać, zlecenie zakończone błędem — ponowić, a zlecenie domknięte traci
 * cały rząd, bo rejestr sterowania nie ma dla niego czynności.
 */
function opiszSterowanie(karta: HTMLElement, wzory: Wzory, zlecenie: AssistantAction): void {
  const rzad = karta.querySelector<HTMLElement>('.sta-chip-rzad');
  if (rzad === null) return;
  const przyciski = [...rzad.querySelectorAll<HTMLElement>('button')];
  if (zlecenie.status === AssistantActionStatus.Failed) {
    if (wzory.ponow === null) {
      rzad.remove();
      return;
    }
    const ponow = wzory.ponow.cloneNode(true) as HTMLElement;
    ponow.dataset.sterowanie = AssistantActionControl.Retry;
    rzad.replaceChildren(ponow);
    return;
  }
  if (!STANY_STEROWALNE.includes(zlecenie.status)) {
    rzad.remove();
    return;
  }
  const wstrzymaj = przyciski[0];
  const odwolaj = przyciski[1];
  if (wstrzymaj !== undefined) wstrzymaj.dataset.sterowanie = AssistantActionControl.Pause;
  if (odwolaj !== undefined) odwolaj.dataset.sterowanie = AssistantActionControl.Cancel;
}

/** Wiąże przyciski kart Actions Monitor ze sterowaniem zleceniem; po przyjęciu czynności wykaz czyta się na nowo. */
function zwiazSterowanieZleceniem(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const panel = cialo.querySelector<HTMLElement>('#panel-actions');
  panel?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-sterowanie]');
    const zlecenie = przycisk?.closest<HTMLElement>('.am-karta')?.dataset.zlecenie;
    const sterowanie = przycisk?.dataset.sterowanie;
    if (zlecenie === undefined || sterowanie === undefined) return;
    void wywolaj(kanal, Command.AssistantActionStatus, {
      windowId: idOkna,
      actionId: zlecenie,
      control: sterowanie as AssistantActionControl,
    }).then(odswiez);
  });
}

/** Wypełnia Activity Feed dziennikiem działań okna; zapytanie puste bierze dziennik bez zawężenia. */
async function wypelnijDziennik(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  zapytanie: string,
): Promise<void> {
  const panel = cialo.querySelector<HTMLElement>('#panel-activity .dn-wykaz-modulu');
  if (panel === null) return;
  const wynik = await wywolaj(kanal, Command.AssistantActivityList, {
    windowId: idOkna,
    query: zapytanie === '' ? undefined : zapytanie,
    limit: 100,
  });
  for (const stojacy of panel.querySelectorAll('.af-poz')) stojacy.remove();
  if (wzory.wpisDziennika === null) return;
  for (const wpis of wynik.udany ? (wynik.wynik?.entries ?? []) : []) {
    panel.appendChild(zbudujWpis(wzory.wpisDziennika, wpis));
  }
}

/** Składa wiersz dziennika: godzina wpisu i jego treść; wpis wyróżniony przez Operatora niesie to oznaczenie dalej. */
function zbudujWpis(wzor: HTMLElement, wpis: AssistantActivityEntry): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.wpis = wpis.id;
  if (wpis.important === true) wiersz.setAttribute('aria-current', 'true');
  const czas = wiersz.querySelector('.czas');
  if (czas !== null) czas.textContent = godzina(wpis.createdAt);
  wpiszTekst(wiersz, ` ${wpis.content}`);
  return wiersz;
}

/** Wiąże pole szukania Activity Feed z zawężeniem dziennika; rodzina assistant.* zawęża go treścią wpisu. */
function zwiazSzukanieWDzienniku(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
): void {
  const pole = cialo.querySelector<HTMLInputElement>('#panel-activity .dn-szukaj input');
  pole?.addEventListener('change', () => {
    void wypelnijDziennik(kanal, idOkna, cialo, wzory, pole.value.trim());
  });
}

/**
 * Wiąże pole polecenia z komendą asystenta. Polecenie idzie transkrypcją, bo
 * nagrania okno nie zbiera — nagranie należy do rodziny mowy, nie do
 * assistant.*.
 */
function zwiazPolecenie(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  odswiez: () => void,
): void {
  const pole = cialo.querySelector<HTMLTextAreaElement>('.sta-prompt-obszar');
  if (pole === null) return;
  const wyslij = (): void => {
    const tresc = pole.value.trim();
    if (tresc === '') return;
    pole.value = '';
    void wywolaj(kanal, Command.AssistantVoiceCommand, {
      windowId: idOkna,
      transcript: tresc,
    }).then(odswiez);
  };
  cialo.querySelector('.sta-prompt')?.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    wyslij();
  });
  for (const przycisk of cialo.querySelectorAll('.sta-prompt-wyslij')) {
    przycisk.addEventListener('click', wyslij);
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
