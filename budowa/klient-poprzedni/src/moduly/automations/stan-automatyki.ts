import type { AutomationWorkflow } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Automatyka bieżąca modułu — jedna prawda dla pięciu okien.
 *
 * Wszystkie okna pracują nad tą samą automatyką: Workflow Builder buduje jej
 * kroki, Orchestrator układa między nimi zależności, Scheduler nadaje jej
 * cykliczność, Queue Manager uruchamia jej kolejkę, a Execution Monitor
 * pokazuje jej przebiegi. Gdyby każde okno trzymało własne wskazanie, zmiana
 * automatyki w jednym rozjechałaby cztery pozostałe.
 *
 * Kolejka bieżąca stoi obok automatyki, bo Queue Manager musi wiedzieć, którą
 * kolejkę posuwa, a Execution Monitor — po której przyszedł stan przebiegu.
 * Kolejka jest bytem silnika kolejek, nie definicji, więc nie należy do
 * `AutomationWorkflow` i nie da się jej z niego odczytać.
 */
export interface StanAutomatyki {
  /** Identyfikator automatyki bieżącej; pusty znaczy „nie wskazano”. */
  automatyka(): string;
  /** Definicja automatyki bieżącej; pusta do pierwszego odczytu. */
  definicja(): AutomationWorkflow | null;
  /** Przestawia moduł na inną automatykę i powiadamia okna. */
  ustawAutomatyke(idAutomatyki: string, definicja?: AutomationWorkflow): void;
  /** Zapisuje definicję po zmianie i powiadamia okna zależne. */
  ustawDefinicje(definicja: AutomationWorkflow): void;
  /** Kolejka wykonująca automatykę; pusta, dopóki żadnej nie założono. */
  kolejka(): string;
  /** Zapisuje kolejkę bieżącą i powiadamia okna. */
  ustawKolejke(idKolejki: string): void;
  /** Subskrypcja zmiany wskazania albo definicji. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzeń rdzenia. */
  zamknij(): void;
}

/** Zależności stanu: automatyka i kolejka otwierane od razu. */
export interface OpcjeStanu {
  automatyka?: string;
  kolejka?: string;
}

export function utworzStanAutomatyki(
  zrodlo: ZrodloAutomations,
  opcje: OpcjeStanu = {},
): StanAutomatyki {
  const sluchacze = new Set<() => void>();
  let idAutomatyki = opcje.automatyka ?? '';
  let idKolejki = opcje.kolejka ?? '';
  let definicja: AutomationWorkflow | null = null;

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  // Stan przebiegu przychodzi także z pracy innego okna albo innego urządzenia
  // tego konta. Dotyczy naszej automatyki — więc definicja w oknie jest już
  // nieaktualna i okna mają się odczytać ponownie.
  const odsubskrybujPrzebieg = zrodlo.naStanPrzebiegu((tresc) => {
    if (tresc.execution.workflowId !== idAutomatyki) return;
    powiadom();
  });

  const odsubskrybujKolejke = zrodlo.naZmianeKolejki((tresc) => {
    if (tresc.queue.id !== idKolejki) return;
    powiadom();
  });

  return {
    automatyka: () => idAutomatyki,

    definicja: () => definicja,

    ustawAutomatyke(nowa, nowaDefinicja) {
      const przyciety = nowa.trim();
      if (przyciety === idAutomatyki && nowaDefinicja === undefined) return;
      idAutomatyki = przyciety;
      definicja = nowaDefinicja ?? null;
      powiadom();
    },

    ustawDefinicje(nowa) {
      definicja = nowa;
      idAutomatyki = nowa.id;
      powiadom();
    },

    kolejka: () => idKolejki,

    ustawKolejke(nowa) {
      const przyciety = nowa.trim();
      if (przyciety === idKolejki) return;
      idKolejki = przyciety;
      powiadom();
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybujPrzebieg();
      odsubskrybujKolejke();
    },
  };
}
