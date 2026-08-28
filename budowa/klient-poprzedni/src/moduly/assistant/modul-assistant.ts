import './assistant.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzOknoActionsMonitor, type OknoActionsMonitor } from './okno-actions-monitor';
import { utworzOknoActivityFeed, type OknoActivityFeed } from './okno-activity-feed';
import {
  utworzOknoCommandToolsHub,
  type OknoCommandToolsHub,
} from './okno-command-tools-hub';
import {
  utworzOknoMemoryContextManager,
  type OknoMemoryContextManager,
} from './okno-memory-context-manager';
import { utworzOknoVoiceConsole, type OknoVoiceConsole } from './okno-voice-console';
import { utworzStanAssistant, type StanAssistant } from './stan-assistant';
import { utworzZrodloNarzedzi, type ZrodloNarzedzi } from './zrodlo-narzedzi';
import { utworzZrodloPamieci, type ZrodloPamieci } from './zrodlo-pamieci';

/**
 * Moduł Assistant składa pięć okien operacyjnych w jednym układzie: okno
 * wiodące poleceń głosowych, dwa okna monitorujące przebieg oraz dwa okna
 * zarządzające zasobem pamięci, kontekstów, narzędzi i rutyn.
 */
export interface ModulAssistant {
  /** Element osadzany w obszarze roboczym powłoki; moduł nie osadza go sam. */
  element: HTMLElement;
  /** Ustala okno modułu i czyta zlecenia, dziennik, pamięć oraz katalogi. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulAssistant(kanal: Kanal): ModulAssistant {
  const stan: StanAssistant = utworzStanAssistant(kanal);
  const pamiec: ZrodloPamieci = utworzZrodloPamieci(kanal);
  const narzedzia: ZrodloNarzedzi = utworzZrodloNarzedzi(kanal);

  const konsola: OknoVoiceConsole = utworzOknoVoiceConsole(stan);
  const monitor: OknoActionsMonitor = utworzOknoActionsMonitor(stan);
  const dziennik: OknoActivityFeed = utworzOknoActivityFeed(stan);
  const zarzadcaPamieci: OknoMemoryContextManager = utworzOknoMemoryContextManager(stan, pamiec);
  // Kafel akcji wypełnia pole polecenia; wysyłkę rozstrzyga Operator.
  const hub: OknoCommandToolsHub = utworzOknoCommandToolsHub(stan, narzedzia, (tresc) =>
    konsola.ustawPolecenie(tresc),
  );

  const pasMonitorow = document.createElement('div');
  pasMonitorow.className = 'ma-modul__pas ma-modul__pas--monitory';
  pasMonitorow.append(monitor.element, dziennik.element);

  const pasZarzadcow = document.createElement('div');
  pasZarzadcow.className = 'ma-modul__pas ma-modul__pas--zarzadcy';
  pasZarzadcow.append(zarzadcaPamieci.element, hub.element);

  const element = document.createElement('div');
  element.className = 'ma-modul';
  element.dataset['modul'] = 'assistant';
  element.setAttribute('aria-label', 'Moduł Assistant — okna operacyjne');
  element.append(konsola.element, pasMonitorow, pasZarzadcow);

  const odsubskrybuj = stan.obserwuj(() => {
    konsola.odswiez();
    monitor.odswiez();
    dziennik.odswiez();
  });

  return {
    element,

    async wczytaj(idSesji) {
      // Okno modułu ustala się przed odczytami, bo polecenie wymaga wskazania okna.
      await stan.ustalOkno(idSesji);
      // Odczyty idą równolegle; odmowa jednego zostaje w jego oknie.
      await Promise.all([
        stan.odswiezZlecenia(),
        stan.odswiezDziennik(),
        konsola.wczytaj(),
        zarzadcaPamieci.wczytaj(),
        hub.wczytaj(),
      ]);
    },

    rozlacz() {
      odsubskrybuj();
      stan.rozlacz();
    },
  };
}
