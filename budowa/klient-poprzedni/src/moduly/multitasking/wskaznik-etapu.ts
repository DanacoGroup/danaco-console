/**
 * Wiersz etapu przepływu: nazwa etapu i wskaźnik kropkowy postępu. Etap nie
 * pochodzi z podagenta — niesie go telemetria `monitor.status`, którą panel wiąże
 * z przepływem po polu `windowId`.
 */

import type { MonitorStatus } from '../../../../shared/contract';
import { kropkiEtapu } from './format-zadan';


/**
 * Telemetria okien odczytana jednym `monitor.status` dla całej karty sesji. Jeden
 * odczyt na kartę, a nie na przepływ, bo monitor oddaje stan wszystkich procesów
 * naraz, a pytanie osobno dla każdego okna mnożyłoby wywołania bez nowej treści.
 */
export interface TelemetriaOkien {
  /** Stan procesu okna; pusty, gdy monitor takiego procesu nie oddał. */
  stan(idOkna: string): MonitorStatus | null;
  /** Powód milczenia telemetrii; pusty, gdy odczyt się udał. */
  odmowa(): string;
}

/**
 * Telemetria złożona z odpowiedzi `monitor.status`.
 *
 * Jeden proces na okno: gdy monitor odda dla okna kilka procesów, wygrywa
 * najświeższy po `updatedAt` — starszy wpis pokazywałby etap sprzed zmiany.
 */
export function zlozTelemetrie(
  statusy: readonly MonitorStatus[],
  odmowa: string,
): TelemetriaOkien {
  const wedlugOkna = new Map<string, MonitorStatus>();
  for (const status of statusy) {
    const idOkna = status.windowId ?? '';
    if (idOkna === '') continue;
    const poprzedni = wedlugOkna.get(idOkna);
    if (poprzedni === undefined || poprzedni.updatedAt <= status.updatedAt) {
      wedlugOkna.set(idOkna, status);
    }
  }
  return {
    stan: (idOkna) => wedlugOkna.get(idOkna) ?? null,
    odmowa: () => odmowa,
  };
}

/**
 * Wiersz etapu dla jednego przepływu. Wiersz powstaje także wtedy, gdy telemetria
 * milczy: wypisuje wówczas przyczynę ciszy, zamiast rysować wskaźnik, którego nikt
 * nie odczytał.
 */
export function wierszEtapu(idOkna: string, telemetria: TelemetriaOkien): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dm-etap';

  const status = telemetria.stan(idOkna);
  if (status === null) {
    element.dataset['pokrycie'] = 'brak';
    element.textContent =
      telemetria.odmowa() === ''
        ? 'Etap nieznany — monitor nie oddał procesu tego okna.'
        : `Etap nieznany — ${telemetria.odmowa()}.`;
    return element;
  }

  const nazwa = document.createElement('span');
  nazwa.className = 'dm-etap__nazwa';
  nazwa.textContent = status.stage ?? 'etap bez nazwy w telemetrii';
  element.append(nazwa);

  const kropki = wskaznikKropkowy(status);
  if (kropki !== null) element.append(kropki);
  return element;
}

/**
 * Wskaźnik kropkowy albo nic.
 *
 * Bez `stageCount` nie ma skali, a wskaźnik bez skali nic nie mówi. Stopień
 * ukończenia w procentach idzie wtedy tekstem, gdy telemetria go niesie — ta
 * sama wiedza o innej ziarnistości.
 */
function wskaznikKropkowy(status: MonitorStatus): HTMLElement | null {
  const skala = status.stageCount ?? 0;
  if (skala <= 0) {
    if (status.completion === undefined) return null;
    const procent = document.createElement('span');
    procent.className = 'dm-etap__procent';
    procent.textContent = `${status.completion}%`;
    return procent;
  }

  const numer = status.stageIndex ?? 0;
  const kropki = document.createElement('span');
  kropki.className = 'dm-etap__kropki';
  kropki.textContent = kropkiEtapu(numer, skala);
  kropki.setAttribute('aria-label', `etap ${numer} z ${skala}`);
  kropki.title = `Etap ${numer} z ${skala} wg telemetrii monitor.status.`;
  return kropki;
}
