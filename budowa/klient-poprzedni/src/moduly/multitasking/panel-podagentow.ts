import { SubagentStatus, type Subagent } from '../../../../shared/contract';
import { przyciskAkcji as przycisk, pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import { czasOdcinka, zapiszCzasOdcinka } from './format-zadan';
import { PRZEDROSTEK } from './kontrolki';
import type { ZrodloPodagentow } from './zrodlo-podagentow';

/**
 * Treść panelu Subagent Network — żywy wykaz podagentów okna wykonawcy.
 *
 * Dokłada się do wytwórni panelu (`panel-subagent-network.ts:
 * utworzPanelPodagentow`) jako sekcja domyślnie zamknięta (`<details>`), żeby
 * nie przykrywać definicji podagentów i kontrolek, które panel niesie wyżej.
 *
 * Zatrzymanie idzie komendą `subagent.stop` na wskazanym podagencie: tura okna
 * i podagenci niewskazani zostają nietknięci.
 *
 * Odmowa mówi, co się nie stało i czym to zmienić; wykaz pusty jest krótszą
 * listą, nie bramką.
 */
export interface TrescPodagentow {
  /** Odpina nasłuch `subagent.changed`; woła je zamknięcie okna wykonawcy. */
  rozlacz?(): void;
  element: HTMLElement;
  /** Ponawia odczyt wykazu żywych podagentów. */
  odswiez(): void;
}

export interface OpcjeTresciPodagentow {
  /** Źródło `subagent.list` i `subagent.result.collect`. */
  podagenci: ZrodloPodagentow;
  /** Okno tego wykonawcy; puste, dopóki obsada go nie założyła. */
  okno(): string;
  potwierdz(zdanie: string, udane: boolean): void;
}

/** Stan podagenta po polsku — słownik `SubagentStatus` kontraktu. */
const NAZWY_STANOW: Readonly<Record<SubagentStatus, string>> = {
  [SubagentStatus.Pending]: 'oczekuje',
  [SubagentStatus.Running]: 'w toku',
  [SubagentStatus.Done]: 'zakończony',
  [SubagentStatus.Failed]: 'błędny',
  [SubagentStatus.Stopped]: 'zatrzymany',
};

/** Górna długość wyniku pokazywanego w wierszu; całość niesie `title`. */
const SKROT_WYNIKU = 160;

export function utworzTrescPodagentow(opcje: OpcjeTresciPodagentow): TrescPodagentow {
  const lista = wykaz('Podagenci powołani pod tym wykonawcą', 'dm-wykaz');

  const licznik = document.createElement('span');
  licznik.className = 'dm-zadania__licznik';
  licznik.textContent = '0';

  const tytul = document.createElement('summary');
  tytul.className = 'dm-podagenci-zywi__tytul';
  tytul.append('Podagenci powołani (żywy wykaz) ', licznik);

  const odswiezenie = przycisk('Odśwież wykaz', 'dn-btn dn-btn--sm');
  odswiezenie.addEventListener('click', () => {
    void odswiez();
  });
  const zbieranie = przycisk('Zbierz wyniki', 'dn-btn dn-btn--sm');
  zbieranie.addEventListener('click', () => {
    void zbierz();
  });

  const pasek = document.createElement('div');
  pasek.className = 'dm-podagenci-zywi__pasek';
  pasek.append(odswiezenie, zbieranie);

  // Wykaz słucha, zamiast odpytywać: rdzeń rozgłasza `subagent.changed` przy
  // powołaniu, wejściu w bieg, zakończeniu i zatrzymaniu. Przycisk „Odśwież
  // wykaz" zostaje, bo zdarzenie mówi o zmianie, a nie o stanie zastanym,
  // i po ponownym nawiązaniu łącza wykaz trzeba odczytać raz od nowa.
  //
  // Sito po oknie stoi tutaj, bo zdarzenie idzie do wszystkich słuchaczy:
  // panel wykonawcy nie ma odświeżać się na cudzym podagencie.
  const odsubskrybuj = opcje.podagenci.naPodagenta((tresc) => {
    const okno = opcje.okno();
    if (okno === '' || tresc.subagent.windowId !== okno) return;
    void odswiez();
  });

  const oDrodze = document.createElement('p');
  oDrodze.className = 'dm-podagenci-zywi__droga';
  oDrodze.textContent =
    'Zatrzymanie: „Zatrzymaj" wykonuje subagent.stop na wskazanym podagencie tego okna. ' +
    'Tura okna i podagenci niewskazani zostają nietknięci; podagent już zakończony ' +
    'nie jest błędem — rdzeń zwraca go w wykazie nieczynnych.';

  // Sekcja jest domyślnie zamknięta: panel niesie wyżej definicje podagentów
  // i kontrolki wytwórni, a żywy wykaz otwiera Operator, gdy praca w tle biegnie.
  const element = document.createElement('details');
  element.className = 'dm-podagenci-zywi';
  element.append(tytul, pasek, lista, oDrodze);

  /** Odczyt wykazu żywych podagentów okna tego wykonawcy. */
  async function odswiez(): Promise<void> {
    const okno = opcje.okno();
    if (okno === '') {
      pokaz([], 'Executor nie ma jeszcze okna — załóż okno wykonawcy w obsadzie, a wykaz ożyje.');
      return;
    }
    const wynik = await opcje.podagenci.wykaz({ windowId: okno });
    if (!wynik.udany || wynik.wynik === undefined) {
      pokaz(
        [],
        `Rdzeń odmówił wykazu podagentów (kod ${wynik.blad?.code ?? 'brak'}: ${wynik.blad?.message ?? 'bez opisu'}). Powtórz „Odśwież wykaz".`,
      );
      return;
    }
    pokaz(
      wynik.wynik,
      'To okno nie powołało jeszcze ani jednego podagenta — powołuje ich model narzędziem subagent.spawn w trakcie tury albo Operator sekcją „Dodaj subagenta" pod tym wykazem.',
    );
  }

  /** `subagent.result.collect` — zbiór po drodze, bez blokowania panelu. */
  async function zbierz(): Promise<void> {
    const okno = opcje.okno();
    if (okno === '') {
      opcje.potwierdz('Executor nie ma jeszcze okna — nie ma czyich wyników zbierać.', false);
      return;
    }
    const wynik = await opcje.podagenci.zbierz({ windowId: okno, waitForAll: false });
    if (!wynik.udany || wynik.wynik === undefined) {
      opcje.potwierdz(
        `Rdzeń odmówił zebrania wyników (kod ${wynik.blad?.code ?? 'brak'}: ${wynik.blad?.message ?? 'bez opisu'}).`,
        false,
      );
      return;
    }
    opcje.potwierdz(
      wynik.wynik.complete
        ? `Zebrane wyniki ${wynik.wynik.subagents.length} podagentów — wszyscy domknęli pracę.`
        : `Zebrane wyniki ${wynik.wynik.subagents.length} podagentów — część nadal pracuje (zbiór po drodze; na komplet czeka dopiero waitForAll).`,
      true,
    );
    pokaz(wynik.wynik.subagents, 'Zbieranie nie zastało ani jednego podagenta.');
  }

  /**
   * Zatrzymanie wskazanego podagenta komendą `subagent.stop`.
   *
   * Zdanie powstaje z odpowiedzi rdzenia, a nie z żądania: rdzeń oddaje osobno
   * tych, których zatrzymał (`stopped`), i tych, którzy w chwili wywołania nie
   * pracowali (`notRunning`) — drugi przypadek nie jest błędem i panel go tak
   * nie nazywa.
   */
  async function zatrzymaj(podagent: Subagent): Promise<void> {
    const okno = opcje.okno();
    if (okno === '') {
      opcje.potwierdz('Executor nie ma jeszcze okna — nie ma czyich podagentów zatrzymywać.', false);
      return;
    }
    const wynik = await opcje.podagenci.zatrzymaj({ windowId: okno, subagentIds: [podagent.id] });
    if (!wynik.udany || wynik.wynik === undefined) {
      opcje.potwierdz(
        `Rdzeń odmówił zatrzymania podagenta ${podagent.id} (kod ${wynik.blad?.code ?? 'brak'}: ${wynik.blad?.message ?? 'bez opisu'}).`,
        false,
      );
      return;
    }
    const oddane = wynik.wynik;
    if (oddane.stopped.includes(podagent.id)) {
      opcje.potwierdz(`Rdzeń zatrzymał podagenta ${podagent.id}.`, true);
    } else if (oddane.notRunning.includes(podagent.id)) {
      opcje.potwierdz(
        `Podagent ${podagent.id} nie pracował w chwili wywołania — rdzeń zwrócił go jako nieczynnego, nie jako błąd.`,
        true,
      );
    } else {
      opcje.potwierdz(
        `Rdzeń przyjął wywołanie, ale nie wymienił podagenta ${podagent.id} ani wśród zatrzymanych, ani wśród nieczynnych — zatrzymanie niepotwierdzone.`,
        false,
      );
    }
    await odswiez();
  }

  /** Rysuje wykaz albo stan pusty; licznik mówi o pozycjach czynnych. */
  function pokaz(podagenci: readonly Subagent[], zdaniePustki: string): void {
    const teraz = Date.now();
    lista.replaceChildren();
    const czynni = podagenci.filter(
      (podagent) =>
        podagent.status === SubagentStatus.Pending || podagent.status === SubagentStatus.Running,
    ).length;
    licznik.textContent = String(czynni);
    licznik.title = `Podagenci czynni: ${czynni} z ${podagenci.length}.`;
    if (podagenci.length === 0) {
      const pusto = document.createElement('li');
      pusto.className = 'dm-pozycja dm-pozycja--pusta';
      pusto.textContent = zdaniePustki;
      lista.append(pusto);
      return;
    }
    for (const podagent of podagenci) {
      lista.append(wierszPodagenta(podagent, teraz, (wskazany) => void zatrzymaj(wskazany)));
    }
  }

  pokaz([], 'Wykaz nieodczytany — użyj „Odśwież wykaz".');

  return {
    element,
    odswiez: () => {
      void odswiez();
    },
    rozlacz: () => {
      odsubskrybuj();
    },
  };
}

/** Jeden wiersz wykazu: nazwa, stan z czasem i wynikiem, przycisk zatrzymania. */
function wierszPodagenta(
  podagent: Subagent,
  teraz: number,
  zatrzymaj: (podagent: Subagent) => void,
): HTMLElement {
  const { element, akcje } = pozycjaWykazu(
    podagent.name ?? podagent.id,
    opisPodagenta(podagent, teraz),
    PRZEDROSTEK,
  );
  element.dataset['podagent'] = podagent.id;
  if (podagent.result !== undefined && podagent.result !== '') {
    element.title = podagent.result;
  }
  const stop = przycisk('Zatrzymaj', 'dn-btn dn-btn--sm dn-btn--niebezpieczny');
  stop.addEventListener('click', () => {
    zatrzymaj(podagent);
  });
  akcje.append(stop);
  return element;
}

/** Zdanie opisu wiersza: stan · czas · wynik (skrócony). */
function opisPodagenta(podagent: Subagent, teraz: number): string {
  const czesci = [
    NAZWY_STANOW[podagent.status] ?? podagent.status,
    zapiszCzasOdcinka(czasOdcinka(podagent.startedAt, podagent.finishedAt, teraz)),
  ];
  if (podagent.result !== undefined && podagent.result !== '') {
    const skrot =
      podagent.result.length > SKROT_WYNIKU
        ? `${podagent.result.slice(0, SKROT_WYNIKU)}…`
        : podagent.result;
    czesci.push(`wynik: ${skrot}`);
  } else {
    czesci.push('wynik: jeszcze go nie ma');
  }
  return czesci.join(' · ');
}
