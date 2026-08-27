import type { AutomationWorkflow } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Automatyka bieżąca modułu — jedna prawda dla pięciu okien. Wszystkie okna
 * pracują nad tą samą automatyką, więc gdyby każde trzymało własne wskazanie,
 * zmiana automatyki w jednym rozjechałaby cztery pozostałe.
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

/**
 * Zależności stanu: automatyka i kolejka otwierane od razu. Obie są nieobowiązkowe,
 * bo moduł otwarty bez wskazania czeka na wybór, zamiast podstawiać pierwszą pozycję
 * z wykazu i pracować nad nieswoją automatyką.
 */
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

  // Stan przebiegu z innego okna tego konta unieważnia definicję w oknie.
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
