import { SubagentStatus, type MonitorStatus } from '../../../../shared/contract';
import { przyciskAkcji, wybor } from '../../modele/kontrolki-formularza';
import { kartaPrzeplywu } from './karta-przeplywu';
import { utworzStanTresci } from './stany-okna';
import type { StanMultitaskingu } from './stan-multitaskingu';
import { zlozTelemetrie, type TelemetriaOkien } from './wskaznik-etapu';
import { zlozPrzeplywy, type Przeplyw } from './zadania-w-tle';
import type { ZrodloBiegu } from './zrodlo-biegu';
import type { ZrodloPodagentow } from './zrodlo-podagentow';
import { zwinieteZakonczone } from './zwiniete-zakonczone';
import './zadania-w-tle.css';

/**
 * Panel „Zadania w tle" — podagenci całej karty sesji w jednym miejscu.
 *
 * Każdy przepływ (okno wykonawcy wraz z jego zadaniem) dostaje kartę: nazwę
 * z odznaką stanu, wiersz podsumowania, etap z telemetrii i tabelę agentów.
 * Przepływy domknięte idą do zwiniętego licznika, żeby praca trwająca nie
 * schodziła poniżej krawędzi.
 *
 * Panel nie stawia pauzy ani usuwania podagenta, bo kontrakt nie niesie komendy,
 * która by je wykonała: rodzina `subagent.*` to `spawn`, `list` i
 * `result.collect`, a `queue.action` wymaga `queueId`, którego `Subagent` nie
 * niesie (ma sam `queueItemId`).
 */
export interface PanelZadanWTle {
  element: HTMLElement;
  /** Odczytuje podagentów sesji wraz z telemetrią okien. */
  odswiez(): void;
}

export interface OpcjeZadanWTle {
  /** Źródło `subagent.list` i `subagent.result.collect`. */
  podagenci: ZrodloPodagentow;
  /** Bieg pracy — stąd `monitor.status` z etapami procesów okien. */
  bieg: ZrodloBiegu;
  /** Stan modułu; z niego bierze się karta sesji. */
  stan: StanMultitaskingu;
}

/** Pozycje sita stanu; pusta wartość znaczy „bez zawężenia". */
const POZYCJE_SITA: ReadonlyArray<readonly [string, string]> = [
  ['', 'Wszystkie stany'],
  [SubagentStatus.Running, 'W toku'],
  [SubagentStatus.Pending, 'Oczekujące'],
  [SubagentStatus.Done, 'Zakończone'],
  [SubagentStatus.Failed, 'Błędne'],
  [SubagentStatus.Stopped, 'Zatrzymane'],
];

export function utworzPanelZadanWTle(opcje: OpcjeZadanWTle): PanelZadanWTle {
  const { podagenci, bieg, stan } = opcje;
  const tresci = utworzStanTresci();

  const licznik = document.createElement('span');
  licznik.className = 'dm-zadania__licznik';
  licznik.textContent = '0';

  const sito = wybor('Stan zadań w tle', POZYCJE_SITA);
  sito.classList.add('dm-zadania__sito');
  sito.addEventListener('change', () => {
    void odczytaj();
  });

  const odswiezanie = przyciskAkcji('Odśwież', 'dn-btn dn-btn--sm');
  odswiezanie.addEventListener('click', () => {
    void odczytaj();
  });

  const element = zlozPowierzchnie(licznik, sito, odswiezanie, tresci.element);

  /**
   * Odczyt podagentów karty sesji wraz z telemetrią etapów.
   *
   * Chwila odniesienia powstaje raz i wędruje w dół, żeby czasy trwania
   * wszystkich wierszy mierzyły się do tej samej milisekundy.
   */
  async function odczytaj(): Promise<void> {
    const sesja = stan.sesja();
    if (sesja === '') {
      tresci.blad(
        'Zadania w tle nieodczytane: panel nie zna karty sesji, więc nie ma o czyich podagentów zapytać. Wróć tu po wejściu w sesję z listy sesji.',
      );
      pokazLicznik(0);
      return;
    }

    tresci.ladowanie('Odczyt zadań w tle…');
    const wykaz = await podagenci.wykaz(zadanieWykazu(sesja, sito.value));
    if (!wykaz.udany || wykaz.wynik === undefined) {
      tresci.blad(
        'Rdzeń odmówił wykazu podagentów — panel nie pokazuje żadnego zadania w tle, choć zadania mogą biec. Powtórz przyciskiem „Odśwież"; gdy odmowa się utrzymuje, komendy subagent.list nie ma w tym rdzeniu i potrzebne jest jego rozszerzenie.',
        wykaz.blad,
      );
      pokazLicznik(0);
      return;
    }

    const teraz = Date.now();
    const telemetria = await odczytajTelemetrie(bieg, sesja);
    const przeplywy = zlozPrzeplywy(wykaz.wynik, teraz);
    pokazLicznik(przeplywy.filter((przeplyw) => !przeplyw.domkniety).length);

    if (przeplywy.length === 0) {
      tresci.pusto(zdaniePustki(sito.value));
      return;
    }
    narysuj(przeplywy, telemetria, teraz);
  }

  /** Rysuje karty przepływów czynnych i zwinięty licznik domkniętych. */
  function narysuj(
    przeplywy: readonly Przeplyw[],
    telemetria: TelemetriaOkien,
    teraz: number,
  ): void {
    const miejsce = tresci.tresc();
    const lista = document.createElement('ul');
    lista.className = 'dm-zadania__lista';

    const czynne = przeplywy.filter((przeplyw) => !przeplyw.domkniety);
    const domkniete = przeplywy.filter((przeplyw) => przeplyw.domkniety);

    for (const przeplyw of czynne) {
      lista.append(
        kartaPrzeplywu(przeplyw, telemetria, teraz, {
          zbierz: (wskazany) => {
            void zbierz(wskazany);
          },
        }),
      );
    }
    miejsce.append(lista);

    const zwiniete = zwinieteZakonczone(domkniete);
    if (zwiniete !== null) miejsce.append(zwiniete);
  }

  /** `subagent.result.collect` dla jednego przepływu; okno bierze się z niego. */
  async function zbierz(przeplyw: Przeplyw): Promise<void> {
    const wynik = await podagenci.zbierz({
      windowId: przeplyw.idOkna,
      subagentIds: przeplyw.podagenci.map((podagent) => podagent.id),
      waitForAll: false,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.potwierdzenie(
        `Rdzeń odmówił zebrania wyników przepływu „${przeplyw.nazwa}" — wyniki podagentów zostają nieściągnięte. Powód: ${wynik.blad?.message ?? 'rdzeń nie podał przyczyny'} (kod ${wynik.blad?.code ?? 'brak'}). Powtórz zbieranie; gdy odmowa się utrzymuje, rdzeń nie rejestruje subagent.result.collect.`,
        false,
      );
      return;
    }
    const zebrani = wynik.wynik.subagents.length;
    tresci.potwierdzenie(
      wynik.wynik.complete
        ? `Rdzeń oddał wyniki ${zebrani} podagentów przepływu „${przeplyw.nazwa}" — wszyscy objęci zbieraniem domknęli pracę.`
        : `Rdzeń oddał wyniki ${zebrani} podagentów przepływu „${przeplyw.nazwa}"; część nadal pracuje, więc zbiór jest niepełny.`,
      true,
    );
    await odczytaj();
  }

  /** Licznik przy tytule podaje liczbę przepływów czynnych. */
  function pokazLicznik(ile: number): void {
    licznik.textContent = String(ile);
    licznik.title = `Przepływy niedomknięte: ${ile}.`;
  }

  return {
    element,
    odswiez: () => {
      void odczytaj();
    },
  };
}

/** Żądanie `subagent.list`; sito puste nie dokłada pola do żądania. */
function zadanieWykazu(sesja: string, stanSita: string): {
  sessionId: string;
  status?: SubagentStatus;
} {
  return stanSita === ''
    ? { sessionId: sesja }
    : { sessionId: sesja, status: stanSita as SubagentStatus };
}

/**
 * Telemetria etapów albo jej brak — odmowa monitora nie przerywa panelu.
 *
 * Treścią panelu są podagenci z `subagent.list`; etap z `monitor.status` jest
 * dodatkiem, więc jego brak zubaża kartę, zamiast gasić cały wykaz.
 */
async function odczytajTelemetrie(bieg: ZrodloBiegu, sesja: string): Promise<TelemetriaOkien> {
  const wynik = await bieg.stanMonitora({ sessionId: sesja });
  if (!wynik.udany || wynik.wynik === undefined) {
    return zlozTelemetrie(
      [],
      `rdzeń odmówił telemetrii monitor.status (kod ${wynik.blad?.code ?? 'brak'})`,
    );
  }
  const statusy: readonly MonitorStatus[] = wynik.wynik.statuses;
  return zlozTelemetrie(statusy, '');
}

/** Zdanie stanu pustego; sito zawężające mówi o sobie wprost. */
function zdaniePustki(stanSita: string): string {
  const opis = POZYCJE_SITA.find(([wartosc]) => wartosc === stanSita)?.[1] ?? stanSita;
  return stanSita === ''
    ? 'Rdzeń nie oddał ani jednego podagenta tej sesji — nic nie biegnie w tle.'
    : `Rdzeń nie oddał ani jednego podagenta w stanie „${opis}". Zdejmij zawężenie, żeby zobaczyć pozostałe.`;
}

/** Powierzchnia panelu: pasek tytułu z licznikiem i sterowaniem oraz treść. */
function zlozPowierzchnie(
  licznik: HTMLElement,
  sito: HTMLElement,
  odswiezanie: HTMLElement,
  tresc: HTMLElement,
): HTMLElement {
  const tytul = document.createElement('h3');
  tytul.className = 'dm-zadania__tytul';
  tytul.textContent = 'Zadania w tle';

  const sterowanie = document.createElement('div');
  sterowanie.className = 'dm-zadania__sterowanie';
  sterowanie.append(sito, odswiezanie);

  const pasek = document.createElement('div');
  pasek.className = 'dm-zadania__pasek';
  pasek.append(tytul, licznik, sterowanie);

  const element = document.createElement('section');
  element.className = 'dm-zadania';
  element.setAttribute('aria-label', 'Zadania w tle');
  element.append(pasek, tresc);
  return element;
}
