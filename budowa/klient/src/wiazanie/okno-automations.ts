/**
 * Wiązanie wnętrza okna modułu Automations z rdzeniem. Znacznik należy do
 * Właściciela — ten plik wypełnia stojące węzły odpowiedziami rodziny
 * automation.*, powiela wzory zdjęte z treści przykładowej, a węzły, dla
 * których kontrakt nie ma pola, zdejmuje.
 */
import {
  AutomationExecutionStatus,
  Command,
  EventType,
  type AutomationExecution,
  type AutomationWorkflow,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/**
 * Panele prototypu, którym rodzina automation.* nie oddaje ani jednego pola;
 * schodzą razem z przełącznikami menu sesji. Scheduler i Queue Manager schodzą
 * z powodu odczytu: harmonogram i kolejka mają komendę zapisu, lecz ani jednej
 * komendy, którą stan zapisany wróciłby do okna.
 */
const PANELE_BEZ_POKRYCIA = [
  'panel-scheduler',
  'panel-queue',
  'panel-zadania',
  'panel-plan',
  'panel-artefakty',
  'panel-terminal',
];

/** Klasa kropki przebiegu; nazwa klasy pochodzi ze znacznika Właściciela, a przypisanie z rejestru stanów kontraktu. */
const KROPKA_PRZEBIEGU: Record<string, string> = {
  [AutomationExecutionStatus.Succeeded]: 'dn-kropka--sukces',
  [AutomationExecutionStatus.Failed]: 'dn-kropka--blad',
  [AutomationExecutionStatus.Stopped]: 'dn-kropka--ostrzezenie',
  [AutomationExecutionStatus.Running]: 'dn-kropka--sygnal',
  [AutomationExecutionStatus.Paused]: 'dn-kropka--ostrzezenie',
  [AutomationExecutionStatus.Pending]: 'dn-kropka--neutralna',
};

/** Klasa plakietki przebiegu; ta sama zasada co przy kropce, z rejestru stanów kontraktu. */
const PLAKIETKA_PRZEBIEGU: Record<string, string> = {
  [AutomationExecutionStatus.Succeeded]: 'dn-plakietka dn-plakietka--sukces',
  [AutomationExecutionStatus.Failed]: 'dn-plakietka dn-plakietka--blad',
  [AutomationExecutionStatus.Stopped]: 'dn-plakietka dn-plakietka--ostrzezenie',
  [AutomationExecutionStatus.Running]: 'dn-plakietka dn-plakietka--sygnal',
  [AutomationExecutionStatus.Paused]: 'dn-plakietka dn-plakietka--ostrzezenie',
  [AutomationExecutionStatus.Pending]: 'dn-plakietka',
};

/** Wzory zdjęte z treści przykładowej, zanim wykazy zostaną wyczyszczone; klon zachowuje ikonę i klasy nadane przez bibliotekę. */
interface Wzory {
  pozycja: HTMLElement | null;
  przebieg: HTMLElement | null;
}

/** Automatyka wiodąca okna; jej kod rozstrzyga, czyje przebiegi i zależności stoją w panelach. */
interface Stan {
  automatyka: string;
}

/**
 * Znak pasa czynności zdjęty z okna wraz z miejscem, w które wraca. Znak bez
 * wartości schodzi, bo znak pusty niczego nie mówi; wartość podana później —
 * automatyka zapisana już po otwarciu okna — stawia ten sam znak z powrotem.
 */
interface Znak {
  wezel: HTMLElement;
  rodzic: Element;
  miejsce: number;
  wzorTekstu: string;
}

/** Znaki wypełniane odpowiedzią rdzenia: nazwa i wersja automatyki w pasie czynności oraz ostrzeżenie walidacji w Workflow Builderze. */
interface Znaki {
  automatyka: Znak | null;
  wersja: Znak | null;
  ostrzezenie: Znak | null;
}

/**
 * Wiąże wnętrze okna modułu Automations. Brak węzłów wnętrza kończy funkcję bez
 * żadnego działania — okna w ramie nie ma, więc nie ma czego wypełniać.
 */
export function zwiazAutomations(kanal: Kanal, idOkna: string): void {
  const cialo = document.querySelector<HTMLElement>('.sta-cialo');
  if (cialo === null) return;

  const wzory = zdejmijWzory(cialo);
  zdejmijTrescPrzykladowa(cialo);
  const znaki = zdejmijZnakiCzynnosci(cialo);
  zdejmijSterowanieBezPokrycia(cialo);

  const stan: Stan = { automatyka: '' };

  /* Odwołanie stoi wartością, nie deklaracją: deklaracja wchodzi na wierzch
     zasięgu, więc kompilator bierze dla węzła wnętrza typ sprzed rozpoznania
     pustki. */
  /** Powtarza odczyt automatyk i ich przebiegów; wołane przy wejściu i po każdej zmianie ogłoszonej przez rdzeń. */
  const odswiez = (): void => {
    void wypelnijStanowisko(kanal, idOkna, cialo, wzory, znaki, stan);
  };

  zwiazSzyne(cialo, stan, odswiez);
  zwiazCzynnosci(kanal, cialo, znaki, stan, odswiez);
  zwiazZdarzenia(kanal, cialo, wzory, stan, odswiez);

  odswiez();
}

/** Zdejmuje wzory pozycji i przebiegów, zanim wykazy zostaną wyczyszczone; wzór traci oznaczenia bez pokrycia w kontrakcie. */
function zdejmijWzory(cialo: HTMLElement): Wzory {
  const pozycja = sklonuj(cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista .pt-pozycja'));
  // Znak pracy nad automatyką nie ma pola: definicja niesie nazwę, kroki
  // i czynność, nic o toczącym się przebiegu.
  pozycja?.querySelector('.pt-tetno')?.remove();
  pozycja?.removeAttribute('aria-current');

  return {
    pozycja,
    przebieg: sklonuj(cialo.querySelector<HTMLElement>('#panel-monitor .au-przebieg')),
  };
}

/**
 * Zdejmuje treść przykładową wnętrza. Wykaz rozmów i historia wiszą na
 * rodzinach spoza automation.*; przypięty krok, znak pracy, monitor symulacji
 * i miara różnicy nie mają pola w kontrakcie.
 */
function zdejmijTrescPrzykladowa(cialo: HTMLElement): void {
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  cialo.querySelector('.sta-kom-kontekst')?.replaceChildren();
  zdejmijKanwy(cialo);
  zdejmijMiaryMonitora(cialo);
  oznaczChipyCzynnosci(cialo);
}

/**
 * Zdejmuje z Workflow Buildera i Orchestratora kanwy grafu wraz z pasami
 * kreślenia. Położenie węzła zna wyłącznie komenda `automation.step.layout.set`
 * i sama je oddaje; żadna komenda odczytu go nie zwraca, więc graf narysowany
 * w oknie miałby układ zmyślony, nie wczytany.
 */
function zdejmijKanwy(cialo: HTMLElement): void {
  cialo.querySelector('#panel-builder .au-kanwa')?.remove();
  cialo.querySelector('#panel-orch .sta-okno-tresc > svg')?.remove();
  const pasy = [...cialo.querySelectorAll('#panel-builder .au-narzedzia')];
  pasy.at(-1)?.remove();
  for (const chip of pasy[0]?.querySelectorAll('.sta-chip') ?? []) chip.remove();
  // Zapis definicji bierze komplet kroków, których okno nie ma skąd zebrać po
  // zdjęciu kanwy; zostaje sam bieg próbny.
  pasy[0]?.querySelectorAll('button')[1]?.remove();
}

/**
 * Zdejmuje z Execution Monitor miary bez pokrycia: wskaźnik powodzenia liczy
 * się z okna czasu, którego kontrakt nie zwraca, histogram nie ma pola wcale,
 * a przerwania przebiegu rodzina automation.* nie niesie żadną komendą.
 */
function zdejmijMiaryMonitora(cialo: HTMLElement): void {
  cialo.querySelector('#panel-monitor .dn-wykaz-modulu-poz')?.remove();
  for (const histogram of cialo.querySelectorAll('#panel-monitor .au-hist')) histogram.remove();
  const przyciski = [...cialo.querySelectorAll<HTMLElement>('#panel-monitor .au-narzedzia button')];
  przyciski[1]?.remove();
  cialo.querySelector('#panel-monitor .au-narzedzia .sta-chip')?.remove();
  if (przyciski[0] !== undefined) przyciski[0].dataset.czynnosc = 'ponow';
}

/**
 * Porządkuje pas czynności: znak zasięgu wykonania i miara różnicy schodzą,
 * a przycisk sugestii schodzi razem z nimi, bo rodzina automation.* nie niesie
 * ani zasięgu wykonania, ani miary zmian definicji.
 */
function oznaczChipyCzynnosci(cialo: HTMLElement): void {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  chipy[0]?.remove();
  cialo.querySelector('.sta-kontekst-akcji .sta-chip--diff')?.remove();
  cialo.querySelector('.sta-kontekst-akcji button.dn-btn')?.remove();
}

/**
 * Zdejmuje znaki wypełniane odpowiedzią rdzenia wraz z miejscem, w które
 * wracają. Ostrzeżenie walidacji schodzi od razu: mówi o biegu próbnym, którego
 * przy otwarciu okna jeszcze nie było.
 */
function zdejmijZnakiCzynnosci(cialo: HTMLElement): Znaki {
  const chipy = [...cialo.querySelectorAll<HTMLElement>('.sta-kontekst-akcji > .sta-chip')];
  const ostrzezenie = zdejmijZnak(
    cialo.querySelector<HTMLElement>('#panel-builder .au-narzedzia .dn-plakietka'),
  );
  opiszZnak(ostrzezenie, '', false);
  return {
    automatyka: zdejmijZnak(chipy[1] ?? null),
    wersja: zdejmijZnak(chipy[2] ?? null),
    ostrzezenie,
  };
}

/**
 * Zdejmuje sterowanie, którego rodzina automation.* nie obsługuje: panele bez
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
  // Grupowanie równoległe nie ma komendy; walidacja i ścieżka krytyczna zostają.
  const chipyGrafu = [
    ...cialo.querySelectorAll<HTMLElement>('#panel-orch .au-narzedzia .sta-chip'),
  ];
  chipyGrafu[1]?.remove();
  if (chipyGrafu[0] !== undefined) chipyGrafu[0].dataset.czynnosc = 'waliduj';
  if (chipyGrafu[2] !== undefined) chipyGrafu[2].dataset.pole = 'sciezka';
  zdejmijKonektory(cialo);
}

/** Zostawia w menu dodawania sam tytuł i trzy pozycje materiału; konektory i wtyczki należą do rodzin spoza automation.*. */
function zdejmijKonektory(cialo: HTMLElement): void {
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Pyta rdzeń o automatyki Operatora, po czym wypełnia nimi szynę, belki i panele automatyki wiodącej. */
async function wypelnijStanowisko(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  znaki: Znaki,
  stan: Stan,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AutomationWorkflowList, {});
  const automatyki = wynik.udany ? (wynik.wynik?.workflows ?? []) : [];
  wypelnijSzyne(cialo, wzory, automatyki, stan);
  const wiodaca = automatyki.find((pozycja) => pozycja.id === stan.automatyka) ?? automatyki[0];
  stan.automatyka = wiodaca?.id ?? '';
  opiszAutomatyke(cialo, znaki, wiodaca);
  await wypelnijPrzebiegi(kanal, idOkna, cialo, wzory, stan.automatyka);
  await wypelnijZaleznosci(kanal, cialo, wiodaca);
}

/** Stawia w szynie po jednej pozycji na automatykę; nagłówki grup schodzą, bo automatyk kontrakt nie grupuje po module docelowym. */
function wypelnijSzyne(
  cialo: HTMLElement,
  wzory: Wzory,
  automatyki: AutomationWorkflow[],
  stan: Stan,
): void {
  const lista = cialo.querySelector<HTMLElement>('.dn-szyna-modulu-lista');
  if (lista === null) return;
  lista.replaceChildren();
  if (wzory.pozycja === null) return;
  for (const automatyka of automatyki) {
    const pozycja = wzory.pozycja.cloneNode(true) as HTMLElement;
    pozycja.dataset.automatyka = automatyka.id;
    if (automatyka.id === stan.automatyka) pozycja.setAttribute('aria-current', 'true');
    const tytul = pozycja.querySelector('.pt-pozycja-tytul');
    if (tytul !== null) tytul.textContent = automatyka.name;
    lista.appendChild(pozycja);
  }
}

/** Nanosi nazwę i wersję automatyki na tytuły paneli, tytuł okna komunikacji i znaki pasa czynności. */
function opiszAutomatyke(
  cialo: HTMLElement,
  znaki: Znaki,
  automatyka: AutomationWorkflow | undefined,
): void {
  const nazwa = automatyka?.name ?? '';
  opiszTytul(cialo.querySelector('#panel-builder .sta-okno-tytul b'), nazwa);
  opiszTytul(cialo.querySelector('#okno-czat-1 .sta-okno-tytul b'), nazwa);
  opiszZnak(znaki.automatyka, nazwa, false);
  const znacznik = cialo.querySelector('#panel-builder .sta-okno-znacznik');
  const wersja = automatyka?.version;
  if (znacznik !== null) {
    // Znak niezapisanych zmian nie ma pola: definicja niesie numer wersji,
    // a nie stan edycji, której okno po zdjęciu kanwy nie prowadzi.
    if (wersja === undefined) znacznik.remove();
    else znacznik.textContent = (znacznik.textContent ?? '')
      .replace(/\s*·.*$/u, '').replace(/\d+/u, String(wersja));
  }
  opiszZnak(znaki.wersja, wersja === undefined ? '' : String(wersja), true);
}

/** Zapamiętuje znak wraz z miejscem w pasie czynności; znak, którego prototyp nie niesie, daje pustkę. */
function zdejmijZnak(wezel: HTMLElement | null): Znak | null {
  const rodzic = wezel?.parentElement ?? null;
  if (wezel === null || rodzic === null) return null;
  return {
    wezel,
    rodzic,
    miejsce: [...rodzic.children].indexOf(wezel),
    wzorTekstu: wezel.textContent ?? '',
  };
}

/**
 * Wpisuje wartość w znak pasa czynności i stawia go w jego miejscu, gdy stał
 * zdjęty. Znak liczbowy zachowuje podpis wzoru i podmienia w nim samą liczbę;
 * znak nazwy bierze wartość w całości.
 */
function opiszZnak(znak: Znak | null, wartosc: string, liczbowy: boolean): void {
  if (znak === null) return;
  if (wartosc === '') {
    znak.wezel.remove();
    return;
  }
  znak.wezel.textContent = liczbowy ? znak.wzorTekstu.replace(/\d+/u, wartosc) : wartosc;
  if (znak.wezel.isConnected) return;
  znak.rodzic.insertBefore(znak.wezel, znak.rodzic.children[znak.miejsce] ?? null);
}

/** Wypełnia Execution Monitor przebiegami automatyki wiodącej i zakłada na nie obserwację telemetrii. */
async function wypelnijPrzebiegi(
  kanal: Kanal,
  idOkna: string,
  cialo: HTMLElement,
  wzory: Wzory,
  idAutomatyki: string,
): Promise<void> {
  const os = cialo.querySelector<HTMLElement>('#panel-monitor .au-osczasu');
  if (os === null) return;
  os.replaceChildren();
  if (idAutomatyki === '') return;
  const wynik = await wywolaj(kanal, Command.AutomationExecutionSubscribe, {
    workflowId: idAutomatyki,
    windowId: idOkna,
    limit: 20,
  });
  if (wzory.przebieg === null) return;
  for (const przebieg of wynik.udany ? (wynik.wynik?.executions ?? []) : []) {
    os.appendChild(zbudujPrzebieg(wzory.przebieg, przebieg));
  }
}

/** Składa wiersz przebiegu: znak stanu, czas rozpoczęcia, czas trwania i nazwa stanu z rejestru kontraktu. */
function zbudujPrzebieg(wzor: HTMLElement, przebieg: AutomationExecution): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.przebieg = przebieg.id;
  const kropka = wiersz.querySelector<HTMLElement>('.dn-kropka');
  if (kropka !== null) {
    kropka.className = `dn-kropka ${KROPKA_PRZEBIEGU[przebieg.status] ?? 'dn-kropka--neutralna'}`;
  }
  const pola = [...wiersz.querySelectorAll<HTMLElement>('span')];
  const poczatek = pola.find((pole) => pole.className === '');
  if (poczatek !== undefined) poczatek.textContent = chwila(przebieg.startedAt);
  const czas = wiersz.querySelector('.czas');
  if (czas !== null) {
    if (przebieg.finishedAt === undefined) czas.remove();
    else czas.textContent = trwanie(przebieg.finishedAt - przebieg.startedAt);
  }
  const plakietka = wiersz.querySelector<HTMLElement>('.dn-plakietka');
  if (plakietka !== null) {
    plakietka.className = PLAKIETKA_PRZEBIEGU[przebieg.status] ?? 'dn-plakietka';
    plakietka.textContent = przebieg.status;
  }
  return wiersz;
}

/**
 * Sprawdza układ zależności automatyki i nanosi jego ocenę na Orchestratora.
 * Żądanie bez zależności jest samym sprawdzeniem układu zastanego, więc odczyt
 * niczego w rdzeniu nie zmienia.
 */
async function wypelnijZaleznosci(
  kanal: Kanal,
  cialo: HTMLElement,
  automatyka: AutomationWorkflow | undefined,
): Promise<void> {
  const opis = cialo.querySelector('#panel-orch .dn-meta');
  const sciezka = cialo.querySelector<HTMLElement>('#panel-orch [data-pole="sciezka"]');
  if (automatyka === undefined) {
    opis?.remove();
    sciezka?.remove();
    return;
  }
  const wynik = await wywolaj(kanal, Command.AutomationOrchestratorDefine, {
    workflowId: automatyka.id,
  });
  const ocena = wynik.udany ? wynik.wynik : undefined;
  const zastrzezenia = ocena?.issues ?? [];
  if (opis !== null) {
    if (zastrzezenia.length === 0) opis.remove();
    else opis.textContent = zastrzezenia.join(' · ');
  }
  if (sciezka === null) return;
  const nazwy = new Map((automatyka.steps ?? []).map((krok) => [krok.id, krok.name ?? krok.id]));
  const kroki = (ocena?.criticalPathStepIds ?? []).map((krok) => nazwy.get(krok) ?? krok);
  if (kroki.length === 0) sciezka.remove();
  else sciezka.textContent = kroki.join(' → ');
}

/** Wiąże wybór pozycji szyny z podmianą automatyki wiodącej okna. */
function zwiazSzyne(cialo: HTMLElement, stan: Stan, odswiez: () => void): void {
  cialo.querySelector('.dn-szyna-modulu-lista')?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const wybrana = cel.closest<HTMLElement>('.pt-pozycja')?.dataset.automatyka;
    if (wybrana === undefined || wybrana === stan.automatyka) return;
    stan.automatyka = wybrana;
    odswiez();
  });
}

/**
 * Wiąże czynności okna: bieg próbny automatyki, ponowienie ostatniego przebiegu
 * i sprawdzenie układu zależności.
 */
function zwiazCzynnosci(
  kanal: Kanal,
  cialo: HTMLElement,
  znaki: Znaki,
  stan: Stan,
  odswiez: () => void,
): void {
  cialo.querySelector('#panel-builder .au-narzedzia button')?.addEventListener('click', () => {
    if (stan.automatyka === '') return;
    void wywolaj(kanal, Command.AutomationWorkflowSimulate, { workflowId: stan.automatyka })
      .then((wynik) => {
        // Ostrzeżenie Workflow Buildera niesie zastrzeżenia biegu próbnego:
        // to one mówią o kroku bez parametru, a nie ocena układu zależności.
        opiszZnak(znaki.ostrzezenie, (wynik.wynik?.issues ?? []).join(' · '), false);
        odswiez();
      });
  });
  cialo.querySelector('#panel-orch [data-czynnosc="waliduj"]')?.addEventListener('click', () => {
    odswiez();
  });
  cialo.querySelector('#panel-monitor [data-czynnosc="ponow"]')?.addEventListener('click', () => {
    const przebieg = cialo.querySelector<HTMLElement>('#panel-monitor .au-przebieg')?.dataset
      .przebieg;
    if (przebieg === undefined) return;
    void wywolaj(kanal, Command.AutomationExecutionReplay, { executionId: przebieg })
      .then(odswiez);
  });
}

/** Nasłuchuje zmian rdzenia: zmiana definicji odświeża stanowisko, a zmiana stanu przebiegu — samą oś czasu monitora. */
function zwiazZdarzenia(
  kanal: Kanal,
  cialo: HTMLElement,
  wzory: Wzory,
  stan: Stan,
  odswiez: () => void,
): void {
  kanal.naZdarzenie(EventType.AutomationWorkflowChanged, () => {
    odswiez();
  });
  kanal.naZdarzenie(EventType.AutomationExecutionStatus, (tresc) => {
    if (tresc.execution.workflowId !== stan.automatyka || wzory.przebieg === null) return;
    const os = cialo.querySelector<HTMLElement>('#panel-monitor .au-osczasu');
    if (os === null) return;
    const wiersz = zbudujPrzebieg(wzory.przebieg, tresc.execution);
    const stojacy = os.querySelector<HTMLElement>(
      `.au-przebieg[data-przebieg="${tresc.execution.id}"]`,
    );
    if (stojacy === null) os.prepend(wiersz);
    else stojacy.replaceWith(wiersz);
  });
}

/** Zostawia w tytule człon stały prototypu i dokłada do niego wartość rdzenia; wartość pusta zostawia sam człon stały. */
function opiszTytul(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  const staly = (wezel.textContent ?? '').split(' — ')[0] ?? '';
  wezel.textContent = wartosc === '' ? staly : `${staly} — ${wartosc}`;
}

/** Chwila w zapisie daty i zegara ściennego; czas kontraktu idzie w milisekundach epoki. */
function chwila(czas: number): string {
  const data = new Date(czas);
  const zegar = data.toLocaleTimeString('pl-PL', { hour12: false });
  return `${data.toLocaleDateString('sv-SE')} ${zegar}`;
}

/** Czas trwania przebiegu w minutach i sekundach; zapis idzie skrótem prototypu. */
function trwanie(milisekundy: number): string {
  const sekundy = Math.max(0, Math.round(milisekundy / 1000));
  return `${Math.floor(sekundy / 60)}m ${String(sekundy % 60).padStart(2, '0')}s`;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj<T extends Element>(wezel: T | null): T | null {
  return wezel === null ? null : (wezel.cloneNode(true) as T);
}
