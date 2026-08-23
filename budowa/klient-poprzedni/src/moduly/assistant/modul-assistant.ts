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
 * Moduł Assistant — pięć okien operacyjnych osadzonych w jednym układzie.
 *
 * Układ wynika z roli okna. Voice Console jest oknem wiodącym i punktem wejścia
 * modułu, więc stoi w pasie pierwszym na całą szerokość. Actions Monitor
 * i Activity Feed monitorują ten sam przebieg — zlecenie kończy się w monitorze
 * i wchodzi do dziennika — więc stoją w pasie drugim obok siebie. Memory &
 * Context Manager i Command & Tools Hub zarządzają zasobem, a nie przebiegiem:
 * pamięcią, kontekstami, narzędziami i rutynami. Stoją w pasie trzecim, bo
 * praca w nich poprzedza polecenie albo je przeżywa, a nie towarzyszy mu.
 *
 * Chat Window i Execution Loop Window są oknami wspólnymi platformy i leżą poza
 * tym katalogiem: pas komunikacji montuje scena sesji, ta sama we wszystkich
 * modułach.
 *
 * Cały moduł ma jeden stan: zlecenie założone w Voice Console pojawia się
 * w monitorze bez drugiego odczytu, a wybór zlecenia w monitorze zawęża
 * dziennik. Dwa równoległe stany dałyby dwie prawdy o tym samym zleceniu.
 *
 * Źródła wywołań są trzy, bo trzy są rodziny komend, po które moduł sięga:
 * `zrodlo-assistant.ts` prowadzi obszar `assistant.*` wraz z jego zdarzeniem,
 * `zrodlo-pamieci.ts` — `memory.*` i `knowledge.*`, `zrodlo-narzedzi.ts` —
 * `tools.*`, `session.tool.*` i `automation.*`. Jedno źródło o trzech
 * rodzinach urosłoby do pliku, w którym odmowa pamięci i odmowa harmonogramu
 * leżą obok siebie bez żadnego związku.
 *
 * Wszystkie trzy komendy obszaru `assistant.*` mają uchwyt w rdzeniu; wpina je
 * `zarejestrujAsystenta`
 * (`server/internal/core/adapter_modul_asystent_uchwyty.go`), wołane
 * z `kompozycja.go`, a port `Asystent` wypełnia `montaz_porty.go`. Rdzeń
 * przyjmuje `assistant.voice.command`, prowadzi zlecenie przez stany
 * (w kolejce → w toku → wykonane/błąd) i ogłasza każdą zmianę zdarzeniem
 * `assistant.action.changed`. Actions Monitor odświeża się na to zdarzenie bez
 * odpytywania, więc stan zlecenia zmienia się w oknie w chwili, w której
 * zmienia się w rdzeniu. Odmowę merytoryczną okna pokazują jako swój stan błędu
 * wraz z powodem.
 *
 * Tam, gdzie kontrakt nie niesie drogi — przesył dźwięku, odsłuch syntezy
 * i nagrania, reguły retencji pamięci, nazwane konteksty, zakres uprawnień
 * narzędzia, schowek i skróty tekstowe — okno zgłasza brak wprost
 * (`braki-kontraktu.ts`).
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
  // Kafel akcji wypełnia pole polecenia Voice Console, tak samo jak kafel
  // siatki w samym oknie wiodącym: wysyłkę rozstrzyga Operator, a droga do
  // rdzenia zostaje jedna.
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
      // Okno modułu ustalamy przed odczytami: `assistant.voice.command` wymaga
      // pola windowId, a odczyty zawężają się do okna, jeżeli rdzeń je zna.
      await stan.ustalOkno(idSesji);
      // Odczyty idą równolegle: każdy dotyczy innej komendy, a żaden nie
      // warunkuje drugiego. Odmowa jednego zostaje w jego oknie.
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
