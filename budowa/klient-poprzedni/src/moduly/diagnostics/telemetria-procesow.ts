import {
  ProgressStatus,
  type MonitorStatus,
  type MonitorStatusRequest,
  type MonitorSubscribeRequest,
  type ProgressChangedEvent,
} from '../../../../shared/contract';
import { elementIkony } from '../../ikony/ikony';
import {
  pobierzPlik,
  przyciskAkcji,
  przyciskBezKomendy,
  pozycjaWykazu,
  wykaz,
} from '../../modele/kontrolki-formularza-braki';
import { wariantPlakietki, znakStanuPostepu } from '../../motyw/znaczenia-stanow';
import { powodBezKomendy } from './braki-kontraktu';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { cialoNarzedzia, objasnienieNarzedzia, pasekNarzedzia } from './zakladki-narzedzi';
import type { WynikDiagnostyki } from './zrodlo-diagnostics';
import type { ZrodloObserwowalnosci } from './zrodlo-obserwowalnosci';

/**
 * Telemetria procesów — jeden odczyt, dwie perspektywy: zakładki Metrics & Performance
 * i Health & Uptime czytają jeden materiał (`MonitorStatus`) z jednego odczytu, żeby
 * Operator nie widział dwóch migawek tej samej telemetrii naraz.
 */
export interface TelemetriaProcesow {
  /** Ciało zakładki Metrics & Performance. */
  metryki: HTMLElement;
  /** Ciało zakładki Health & Uptime. */
  kondycja: HTMLElement;
  /** Ponowny odczyt stanu procesów z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch `progress.changed`. */
  zamknij(): void;
}

export function utworzTelemetrieProcesow(
  zrodlo: ZrodloObserwowalnosci,
  idOkna: string,
): TelemetriaProcesow {
  const metryki = utworzStanTresci();
  const kondycja = utworzStanTresci();
  /** Procesy z ostatniego udanego odczytu, w kolejności oddanej przez rdzeń. */
  let procesy: readonly MonitorStatus[] = [];
  /** Zmiany doniesione zdarzeniem po wysłaniu ostatniego żądania, po procesie. */
  const zmiany = new Map<string, ProgressChangedEvent>();
  /** Los ostatniego odczytu — rozstrzyga, czy zdarzenie ma co przerysować. */
  let odczyt: 'w-toku' | 'udany' | 'nieudany' = 'w-toku';
  /** Los zapisu obserwacji; `undefined`, dopóki obserwacji nie zakładano. */
  let obserwacjaZapisana: boolean | undefined;

  /** Pobiera stan procesów — obserwację zakłada tylko pierwszy odczyt, dalej idzie monitor.status. */
  async function pobierz(): Promise<WynikDiagnostyki<MonitorStatus[]>> {
    if (obserwacjaZapisana !== undefined) {
      return zrodlo.stanProcesow(zadanieStanu());
    }
    const wynik = await zrodlo.obserwujProcesy(zadanieObserwacji(idOkna));
    if (!wynik.udany || wynik.wynik === undefined) {
      return {
        udany: false,
        ...(wynik.powod === undefined ? {} : { powod: wynik.powod }),
        ...(wynik.blad === undefined ? {} : { blad: wynik.blad }),
      };
    }
    obserwacjaZapisana = wynik.wynik.subscribed;
    return { udany: true, wynik: wynik.wynik.statuses };
  }

  function odswiez(): void {
    // Zmiany sprzed żądania zawiera już migawka; zmiany po żądaniu zostają, bo migawka może ich nie objąć.
    zmiany.clear();
    odczyt = 'w-toku';
    metryki.ladowanie('Odczyt telemetrii procesów…');
    kondycja.ladowanie('Odczyt telemetrii procesów…');
    void pobierz().then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Odmowa nie zostaje pustym wykazem — pusty wykaz nazwałby to brakiem procesów, co nie jest prawdą.
        procesy = [];
        odczyt = 'nieudany';
        const zdanie = zdanieNiepowodzenia('odczytu telemetrii procesów', wynik.powod);
        metryki.blad(zdanie, wynik.blad);
        kondycja.blad(zdanie, wynik.blad);
        return;
      }
      procesy = wynik.wynik;
      odczyt = 'udany';
      rysuj();
      const zdanie = zdanieOOdczycie(procesy.length, obserwacjaZapisana, idOkna);
      metryki.potwierdzenie(zdanie, true);
      kondycja.potwierdzenie(zdanie, true);
    });
  }

  function rysuj(): void {
    rysujMetryki(procesy, zmiany, metryki);
    rysujKondycje(procesy, zmiany, kondycja);
  }

  function eksportuj(): void {
    if (procesy.length === 0) {
      metryki.potwierdzenie(
        'Nie ma czego wyeksportować — ostatni odczyt nie oddał ani jednego procesu.',
        false,
      );
      return;
    }
    pobierzPlik('telemetria-procesow.json', JSON.stringify(procesy, null, 2), 'application/json');
    metryki.potwierdzenie(
      `Telemetria ${procesy.length} proces(ów) pobrana jako plik JSON — stan z ostatniego odczytu, ` +
        'bez zmian doniesionych zdarzeniem.',
      true,
    );
  }

  const odsubskrybuj = zrodlo.naPostep((tresc) => {
    // Zdarzenie o procesie spoza odczytu też ma znaczenie — wiersz powstaje wtedy z samego zdarzenia.
    zmiany.set(tresc.processId, tresc);
    // Przerysowanie ma sens wyłącznie nad udanym odczytem — po odmowie zdarzenia zasłoniłyby treść odmowy.
    if (odczyt !== 'udany') return;
    rysuj();
  });

  return {
    metryki: zakladkaMetryk(metryki.element, eksportuj),
    kondycja: zakladkaKondycji(kondycja.element),
    odswiez,
    zamknij: odsubskrybuj,
  };
}

/**
 * Żądanie obserwacji: windowId jest tu oknem odbierającym telemetrię, nie sitem
 * procesów — sitem są processIds i sessionId, a pusty wykaz procesów oznacza komplet.
 */
function zadanieObserwacji(idOkna: string): MonitorSubscribeRequest {
  return idOkna === '' ? {} : { windowId: idOkna };
}

/**
 * Żądanie odczytu stanu, bez ani jednego pola: windowId znaczy tu sito procesów, nie
 * okno obserwatora jak w monitor.subscribe, więc podanie go zawęziłoby odczyt.
 */
function zadanieStanu(): MonitorStatusRequest {
  return {};
}

/** Zdanie o odczycie — liczba odczytanych procesów rejestru rdzenia oraz los samego zapisu obserwacji na oknie. */
function zdanieOOdczycie(ile: number, zapisana: boolean | undefined, idOkna: string): string {
  const podstawa = `Odczytano telemetrię ${ile} proces(ów) rejestru rdzenia.`;
  if (zapisana === true) return `${podstawa} Obserwacja zapisana na oknie ${idOkna}.`;
  return (
    `${podstawa} Obserwacji rdzeń NIE zapisał (subscribed: false)` +
    `${idOkna === '' ? ' — moduł zmontowano bez okna komunikacji' : ''}. ` +
    'Wykaz i tak zmienia się na żywo: zdarzenie progress.changed dochodzi do wszystkich ' +
    'połączeń konta, a zapis mówi rdzeniowi wyłącznie, które okno których procesów pilnuje.'
  );
}

/** Zakładka Metrics & Performance: pasek czynności nad miejscem treści, z przyciskiem eksportu telemetrii procesów. */
function zakladkaMetryk(stanTresci: HTMLElement, eksportuj: () => void): HTMLElement {
  const eksport = przyciskAkcji('Eksportuj telemetrię procesów');
  eksport.addEventListener('click', eksportuj);

  return cialoNarzedzia(
    objasnienieNarzedzia(
      'Materiałem zakładki jest telemetria procesów rdzenia (monitor.status i monitor.subscribe): ' +
        'stan, etap, stopień ukończenia i licznik obiegów naprawczych. Szeregów czasowych, ' +
        'percentyli ani profilu kontrakt nie przechowuje, więc wykresu trendu nie ma z czego zbudować.',
    ),
    pasekNarzedzia(
      eksport,
      przyciskBezKomendy(
        'Percentyle p50/p95/p99',
        powodBezKomendy(
          'Percentyl liczy się z szeregu czasowego, a rodzina monitor oddaje wyłącznie stan bieżący procesu.',
          'monitor',
        ),
      ),
      przyciskBezKomendy(
        'Profil pprof (wykres płomieniowy)',
        powodBezKomendy('Zrzut profilu procesora i pamięci wymagałby komendy profilowania.', 'monitor'),
      ),
      przyciskBezKomendy(
        'Definiowalne zapytanie metryk',
        powodBezKomendy('Własny agregat nad metrykami wymagałby komendy zapytania o metryki.', 'monitor'),
      ),
    ),
    stanTresci,
  );
}

/** Zakładka Health & Uptime: zestawienie stanów procesów wraz z jasną granicą tego, co kontrakt naprawdę wie. */
function zakladkaKondycji(stanTresci: HTMLElement): HTMLElement {
  return cialoNarzedzia(
    objasnienieNarzedzia(
      'Zakładka mówi o kondycji procesów sesji — jedynym wymiarze kondycji, jaki kontrakt niesie. ' +
        'O usługach zewnętrznych, punktach końcowych i modelach nie orzeka, bo nie ma czym.',
    ),
    pasekNarzedzia(
      przyciskBezKomendy(
        'Zdefiniuj sondę zdrowia',
        powodBezKomendy('Okresowe sprawdzanie kondycji komponentu wymagałoby komendy sondy.', 'monitor'),
      ),
      przyciskBezKomendy(
        'Monitoring syntetyczny',
        powodBezKomendy('Odtwarzalny scenariusz sondy wymagałby komendy zlecającej jego wykonanie.', 'monitor'),
      ),
      przyciskBezKomendy(
        'Budżet błędów (SLO)',
        powodBezKomendy(
          'Cel dostępności i zużycie budżetu wymagają serii pomiarów oraz komendy jej odczytu.',
          'monitor',
        ),
      ),
      przyciskBezKomendy(
        'Mapa zależności komponentów',
        powodBezKomendy('Graf zależności usług wymagałby komendy oddającej powiązania komponentów.', 'monitor'),
      ),
    ),
    stanTresci,
  );
}

/** Wykaz procesów — perspektywa Metrics & Performance, przebieg proces po procesie, w kolejności odczytu rdzenia. */
function rysujMetryki(
  procesy: readonly MonitorStatus[],
  zmiany: ReadonlyMap<string, ProgressChangedEvent>,
  tresc: StanTresci,
): void {
  const znane = new Set(procesy.map((proces) => proces.processId));
  const doniesione = [...zmiany.values()].filter((zmiana) => !znane.has(zmiana.processId));
  if (procesy.length === 0 && doniesione.length === 0) {
    // Pustka jest stanem poprawnym: instalacja bez procesu w biegu nie jest
    // usterką telemetrii.
    tresc.pusto('Rejestr telemetrii rdzenia nie prowadzi w tej chwili żadnego procesu.');
    return;
  }
  const lista = wykaz('Procesy telemetrii rdzenia', 'dg-wykaz');
  for (const proces of procesy) lista.append(wierszProcesu(proces, zmiany.get(proces.processId)));
  for (const zmiana of doniesione) lista.append(wierszZeZdarzenia(zmiana));
  tresc.tresc().append(lista);
}

/** Jeden wiersz procesu; zmiana doniesiona zdarzeniem progress.changed bierze pierwszeństwo przed danymi z odczytu. */
function wierszProcesu(proces: MonitorStatus, zmiana: ProgressChangedEvent | undefined): HTMLElement {
  const stan = zmiana?.status ?? proces.status;
  const pozycja = pozycjaWykazu(proces.label ?? proces.processId, opisProcesu(proces, zmiana), 'dg');
  pozycja.element.dataset['stanProcesu'] = stan;

  const czas = document.createElement('span');
  czas.className = 'dg-pozycja__czas';
  czas.textContent = new Date(proces.updatedAt).toLocaleString('pl-PL');
  pozycja.element.append(czas);

  pozycja.akcje.append(plakietkaStanu(stan));
  return pozycja.element;
}

/**
 * Wiersz procesu znanego wyłącznie ze zdarzenia. Czasu zmiany nie ma czym
 * wypełnić — `progress.changed` go nie niesie — więc wiersz mówi to zamiast
 * wpisywać zegar klienta w miejsce czasu rdzenia.
 */
function wierszZeZdarzenia(zmiana: ProgressChangedEvent): HTMLElement {
  const pozycja = pozycjaWykazu(
    zmiana.processId,
    `${opisPostepu(zmiana)} — proces zgłoszony zdarzeniem po ostatnim odczycie; ` +
      'czasu zmiany zdarzenie nie niesie, odśwież telemetrię, aby odczytać stan pełny.',
    'dg',
  );
  pozycja.element.dataset['stanProcesu'] = zmiana.status;
  pozycja.akcje.append(plakietkaStanu(zmiana.status));
  return pozycja.element;
}

/** Zdanie opisu procesu: etap, ukończenie, obiegi naprawcze, okno i sesja, złożone w jedną linię tekstu. */
function opisProcesu(proces: MonitorStatus, zmiana: ProgressChangedEvent | undefined): string {
  const czesci: string[] = [`proces ${proces.processId}`];
  if (proces.stage !== undefined && proces.stage !== '') czesci.push(`etap ${proces.stage}`);
  if (proces.stageIndex !== undefined && proces.stageCount !== undefined) {
    czesci.push(`krok ${proces.stageIndex}/${proces.stageCount}`);
  }
  if (proces.completion !== undefined) czesci.push(`ukończenie ${proces.completion}%`);
  if (proces.cycle !== undefined) czesci.push(`obiegi naprawcze: ${proces.cycle}`);
  if (proces.windowId !== undefined) czesci.push(`okno ${proces.windowId}`);
  if (proces.sessionId !== undefined) czesci.push(`sesja ${proces.sessionId}`);
  if (zmiana !== undefined) czesci.push(`na żywo: ${opisPostepu(zmiana)} (czas zmiany nieznany)`);
  return czesci.join(' · ');
}

/** Zdanie o treści zdarzenia postępu — wyłącznie te liczby, które rdzeń faktycznie podał w tym zdarzeniu. */
function opisPostepu(zmiana: ProgressChangedEvent): string {
  const etap =
    zmiana.stepLabel === undefined || zmiana.stepLabel === '' ? '' : `${zmiana.stepLabel}, `;
  // totalSteps równe zeru znaczy liczbę etapów nieznaną, nie zero etapów.
  const kroki =
    zmiana.totalSteps === 0
      ? `krok ${zmiana.currentStep} z nieznanej liczby`
      : `krok ${zmiana.currentStep}/${zmiana.totalSteps}`;
  return `${etap}${kroki}, ukończenie ${zmiana.percent}%`;
}

/** Zestawienie stanów — perspektywa Health & Uptime: liczba procesów w każdym możliwym stanie postępu pracy. */
function rysujKondycje(
  procesy: readonly MonitorStatus[],
  zmiany: ReadonlyMap<string, ProgressChangedEvent>,
  tresc: StanTresci,
): void {
  if (procesy.length === 0 && zmiany.size === 0) {
    tresc.pusto(
      'Rejestr telemetrii rdzenia nie prowadzi w tej chwili żadnego procesu — ' +
        'nie ma z czego orzec o kondycji procesów sesji.',
    );
    return;
  }
  const liczniki = policzStany(procesy, zmiany);
  const lista = wykaz('Kondycja procesów sesji', 'dg-wykaz');
  for (const stan of KOLEJNOSC_STANOW) {
    const znak = znakStanuPostepu(stan);
    const pozycja = pozycjaWykazu(znak.etykieta, `procesów: ${liczniki.get(stan) ?? 0}`, 'dg');
    pozycja.element.dataset['stanProcesu'] = stan;
    pozycja.akcje.append(plakietkaStanu(stan));
    lista.append(pozycja.element);
  }

  const uwaga = document.createElement('p');
  uwaga.className = 'dg-stan';
  uwaga.textContent =
    'Zestawienie policzone w oknie z pól status oddanych przez rdzeń — gotowego agregatu ' +
    'kondycji rdzeń nie oddaje. Procent dostępności, historia incydentów i budżet błędów ' +
    'wymagają serii pomiarów, której kontrakt nie przechowuje.';
  tresc.tresc().append(lista, uwaga);
}

/** Kolejność stanów w zestawieniu Health & Uptime: od pracy trwającej do wszystkich stanów zakończenia procesu. */
const KOLEJNOSC_STANOW: readonly ProgressStatus[] = [
  ProgressStatus.Running,
  ProgressStatus.Pending,
  ProgressStatus.Paused,
  ProgressStatus.Stopped,
  ProgressStatus.Done,
  ProgressStatus.Failed,
];

/** Liczba procesów w każdym stanie; zmiana doniesiona zdarzeniem liczy się nad wynikiem ostatniego odczytu. */
function policzStany(
  procesy: readonly MonitorStatus[],
  zmiany: ReadonlyMap<string, ProgressChangedEvent>,
): Map<ProgressStatus, number> {
  const stany = new Map<string, ProgressStatus>();
  for (const proces of procesy) stany.set(proces.processId, proces.status);
  for (const [idProcesu, zmiana] of zmiany) stany.set(idProcesu, zmiana.status);

  const liczniki = new Map<ProgressStatus, number>();
  for (const stan of stany.values()) liczniki.set(stan, (liczniki.get(stan) ?? 0) + 1);
  return liczniki;
}

/**
 * Plakietka stanu procesu — ikona i etykieta obok barwy, wzięte z katalogu motywu, żeby
 * stan procesu miał tę samą nazwę i znak wszędzie w aplikacji.
 */
function plakietkaStanu(stan: ProgressStatus): HTMLElement {
  const znak = znakStanuPostepu(stan);
  const wariant = wariantPlakietki(znak.rodzina);
  const element = document.createElement('span');
  element.className = wariant === '' ? 'dn-plakietka' : `dn-plakietka ${wariant}`;
  element.append(elementIkony(znak.ikona, { rozmiar: 14 }), document.createTextNode(znak.etykieta));
  return element;
}
