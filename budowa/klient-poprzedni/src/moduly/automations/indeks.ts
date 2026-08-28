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
 * Moduł Automations składa pięć okien operacyjnych wokół jednej automatyki, jako komponent
 * własny strefy drugiej strony głównej, bez okna modułowego.
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

/** Zależności złożenia: automatyka otwierana od razu, kolejka wykonująca ją i okno Execution Monitora do obserwacji telemetrii. */
export interface OpcjeAutomations {
  /** Automatyka otwierana od razu; pusta zostawia wskazanie Operatorowi. */
  automatyka?: string;
  /** Kolejka wykonująca automatykę; pusta do pierwszego założenia. */
  kolejka?: string;
  /** Identyfikator okna Execution Monitora; pusty bierze kod z rejestru okien operacyjnych. */
  oknoMonitora?: string;
}

export function zamontujAutomations(
  gospodarz: HTMLElement,
  kanal: Kanal,
  opcje: OpcjeAutomations = {},
): ZamontowaneAutomations {
  const zrodlo = utworzZrodloAutomations(kanal);
  const stan = utworzStanAutomatyki(zrodlo, opcje);

  // Pokrycie bierze rozstrzygnięcie z wykazu komend rdzenia, jedno na moduł, dla wszystkich okien.
  const pokrycie = utworzPokrycieKomend(kanal);

  const obszar = document.createElement('div');
  obszar.className = 'da-modul';
  obszar.dataset['modul'] = KOD_MODULU;

  /** Skok do okna o danym kodzie jest jedyną drogą przywołania w module; bez ogniska tylko dla myszy. */
  function pokaz(kodOkna: string): void {
    nawigacja.wskaz(kodOkna);
    const okno = obszar.querySelector<HTMLElement>(`[data-okno="${kodOkna}"]`);
    okno?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    okno?.focus({ preventScroll: true });
  }

  // Nawigacja niczego nie otwiera ani nie chowa: okna zostają, menu skraca drogę do wybranego.
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

  // Wykaz gotowych pętli stoi przed kreatorami; cztery panele dopełniające osadzają się wewnątrz okien.
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
  // Wykaz komend rdzenia idzie raz na połączenie; kontrolki pokrycia przerysowują się po nim same.
  void pokrycie.odczytaj();

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      // Menu zdejmuje własny nasłuch z dokumentu — bez tego moduł reagowałby po zamknięciu.
      nawigacja.zwin();
      wykazPetli.zamknij();
      orchestrator.zamknij();
      monitor.zamknij();
      pokrycie.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych, po którym nawigacja modułu skacze wprost do tego okna. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/** Samoopisujący się moduł dla rejestru powłoki: kod siedzi w module, nie w mapie po stronie samej powłoki. */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) => widokZMontazu(zamontujAutomations, kanal),
};
