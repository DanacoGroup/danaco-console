import './automations.css';

import type { Kanal } from '../../protokol/kanal';
import { widokZMontazu, type OpisModulu } from '../rejestracja';
import { utworzPokrycieKomend } from '../pokrycie-komend';
import { KODY_OKIEN, KOD_MODULU, KOD_WYKAZU_PETLI } from './kody-okien';
import { utworzNawigacjeOkien } from './nawigacja-okien';
import { utworzOknoExecutionMonitora } from './okno-execution-monitor';
import { utworzOknoWykazuPetli } from './okno-wykaz-petli';
import { utworzOknoOrchestratora } from './okno-orchestrator';
import { utworzOknoQueueManagera } from './okno-queue-manager';
import { utworzOknoSchedulera } from './okno-scheduler';
import { utworzOknoWorkflowBuilder } from './okno-workflow-builder';
import { utworzPanelDozoru } from './panel-dozoru';
import { utworzPanelNadzoru } from './panel-nadzoru';
import { utworzPanelWersji } from './panel-wersji';
import { utworzPanelZlecen } from './panel-zlecen';
import { utworzStanAutomatyki, type StanAutomatyki } from './stan-automatyki';
import { utworzZrodloAutomations } from './zrodlo-automations';

/**
 * Moduł Automations — złożenie pięciu okien operacyjnych wokół jednej
 * automatyki.
 *
 * Moduł nie ma okna modułowego w żadnym środowisku: jest komponentem własnym
 * strefy 2 strony głównej, więc jego okna otwierają się stamtąd, a nie
 * z przestrzeni roboczej karty sesji. Dlatego złożenie nie przyjmuje ani
 * środowiska, ani karty sesji — przyjmuje automatykę i kolejkę.
 *
 * Układ wynika z ról okien. Workflow Builder jest kreatorem budującym
 * strukturę, więc stoi w obszarze głównym jako punkt wejścia; obok niego drugi
 * kreator — Orchestrator — który tę strukturę układa. Pod nimi pas wykonania:
 * dwaj zarządcy (Scheduler, Queue Manager) i monitor (Execution Monitor), bo
 * dopiero po zbudowaniu automatyki jest co uruchamiać i co obserwować. Okno
 * rozmowy modułu nie należy do tego złożenia: jest bytem sesji i składa je
 * warstwa rozmowy.
 *
 * Nad układem stoi pas przywołania. Menu z `nawigacja-okien.ts` skraca drogę do
 * okna do jednego wskazania i niczego nie chowa: siatka zostaje taka sama,
 * żadne okno nie znika i żadne nie dochodzi z urzędu.
 */
export interface ZamontowaneAutomations {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan automatyki wspólny pięciu oknom. */
  stan: StanAutomatyki;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu. */
  zamknij(): void;
}

/** Zależności złożenia. */
export interface OpcjeAutomations {
  /** Automatyka otwierana od razu; pusta zostawia wskazanie Operatorowi. */
  automatyka?: string;
  /** Kolejka wykonująca automatykę; pusta do pierwszego założenia. */
  kolejka?: string;
  /**
   * Identyfikator okna Execution Monitora — nośnik obserwacji telemetrii.
   *
   * Pusty bierze kod z rejestru okien operacyjnych (`kody-okien.ts`). Wartość
   * podana z zewnątrz jest identyfikatorem WYSTĄPIENIA okna (np. `okn_…`
   * z `window.create`), nie kodem rodzaju — rdzeń wpisuje ją do
   * `Queue.windowIds` i do telemetrii postępu, więc musi wskazywać okno, które
   * naprawdę istnieje.
   */
  oknoMonitora?: string;
}

export function zamontujAutomations(
  gospodarz: HTMLElement,
  kanal: Kanal,
  opcje: OpcjeAutomations = {},
): ZamontowaneAutomations {
  const zrodlo = utworzZrodloAutomations(kanal);
  const stan = utworzStanAutomatyki(zrodlo, opcje);

  // Pozycje, których kontrakt nie obsługuje, nie mówią o tym napisem wpisanym
  // w moduł: pokrycie bierze rozstrzygnięcie z wykazu komend rdzenia, więc
  // nazywa brak tam, gdzie jest, i przestaje go nazywać w dniu, w którym rdzeń
  // komendę dostanie. Jedno pokrycie na moduł, wspólne dla wszystkich okien.
  const pokrycie = utworzPokrycieKomend(kanal);

  const obszar = document.createElement('div');
  obszar.className = 'da-modul';
  obszar.dataset['modul'] = KOD_MODULU;

  /**
   * Skok do okna o danym kodzie — jedyna droga przywołania w tym module.
   *
   * Woła ją i nawigacja, i skok Workflow Buildera do Orchestratora, więc obie
   * kończą się tak samo: kafel wjeżdża na ekran i bierze ognisko. Bez ogniska
   * wędrówka klawiaturą wracałaby na początek układu i przywołanie byłoby drogą
   * wyłącznie dla myszy.
   *
   * Wskazanie w nawigacji przestawia się także przy skoku z okna, bo uchwyt ma
   * nieść okno bieżące, a nie ostatnie wybrane w menu.
   */
  function pokaz(kodOkna: string): void {
    nawigacja.wskaz(kodOkna);
    const okno = obszar.querySelector<HTMLElement>(`[data-okno="${kodOkna}"]`);
    okno?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    okno?.focus({ preventScroll: true });
  }

  // Nawigacja stoi na mechanizmie biblioteki (`komponenty/menu-drzewo.ts`)
  // i niczego nie otwiera ani nie chowa: okna zostają w układzie, a menu
  // jedynie skraca drogę do wybranego (`nawigacja-okien.ts`).
  const nawigacja = utworzNawigacjeOkien({ przywolaj: pokaz });

  const builder = utworzOknoWorkflowBuilder(zrodlo, stan, { pokaz }, pokrycie);
  const wykazPetli = utworzOknoWykazuPetli(zrodlo, stan);
  const orchestrator = utworzOknoOrchestratora(zrodlo, stan, pokrycie);
  const scheduler = utworzOknoSchedulera(zrodlo, stan, pokrycie);
  const kolejka = utworzOknoQueueManagera(zrodlo, stan, pokrycie);
  const monitor = utworzOknoExecutionMonitora(
    zrodlo,
    stan,
    opcje.oknoMonitora ?? KODY_OKIEN.executionMonitor,
    pokrycie,
  );

  // Wykaz gotowych pętli stoi przed kreatorami, bo najczęstsza praca Operatora
  // nie jest budowaniem: pętla zapisana wcześniej ma iść w ruch jednym
  // kliknięciem, a kreatory otwiera się dopiero wtedy, gdy trzeba ją zmienić.
  // Cztery panele dopełniające osadzają się WEWNĄTRZ okien, do których należy
  // ich praca: wersje i szablony w Workflow Builderze, okna wykonania i nadzór
  // w Schedulerze, zlecenia w Queue Managerze, log i skarbiec w Execution
  // Monitorze. Szósty kafel na siatce byłby szóstym oknem, którego dokument
  // projektowy modułu nie zna — okien operacyjnych jest pięć.
  builder.element.append(utworzPanelWersji(zrodlo, stan).element);
  scheduler.element.append(utworzPanelNadzoru(zrodlo, stan).element);
  kolejka.element.append(utworzPanelZlecen(zrodlo, stan).element);
  monitor.element.append(utworzPanelDozoru(zrodlo, stan).element);

  const pas = document.createElement('div');
  pas.className = 'da-modul__start';
  pas.append(oznacz(wykazPetli.element, KOD_WYKAZU_PETLI));

  const gorny = document.createElement('div');
  gorny.className = 'da-modul__gora';
  gorny.append(
    oznacz(builder.element, KODY_OKIEN.workflowBuilder),
    oznacz(orchestrator.element, KODY_OKIEN.orchestrator),
  );

  const dolny = document.createElement('div');
  dolny.className = 'da-modul__dol';
  dolny.append(
    oznacz(scheduler.element, KODY_OKIEN.scheduler),
    oznacz(kolejka.element, KODY_OKIEN.queueManager),
    oznacz(monitor.element, KODY_OKIEN.executionMonitor),
  );

  obszar.append(nawigacja.element, pas, gorny, dolny);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    wykazPetli.odswiez();
    builder.odswiez();
    orchestrator.odswiez();
    scheduler.odswiez();
    kolejka.odswiez();
    monitor.odswiez();
  }

  odswiez();
  // Wykaz komend rdzenia idzie raz na połączenie; wszystkie kontrolki pokrycia
  // czekają na tę jedną odpowiedź i przerysowują się po niej same.
  void pokrycie.odczytaj();

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      // Menu zdejmuje własny nasłuch `pointerdown` z dokumentu — bez tego
      // zamknięty moduł nadal reagowałby na kliknięcia poza sobą.
      nawigacja.zwin();
      wykazPetli.zamknij();
      orchestrator.zamknij();
      monitor.zamknij();
      pokrycie.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych — po nim skacze nawigacja. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod siedzi w module, nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) => widokZMontazu(zamontujAutomations, kanal),
};
