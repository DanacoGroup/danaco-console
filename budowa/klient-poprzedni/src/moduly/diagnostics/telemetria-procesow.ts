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
 * Telemetria procesów — jeden odczyt, dwie perspektywy.
 *
 * Zakładki Metrics & Performance oraz Health & Uptime opisują w opracowaniu
 * dwa różne pytania („jak szybko idzie praca” i „czy platforma stoi”), lecz
 * kontrakt niesie na oba jeden materiał: `MonitorStatus`. Gdyby każda zakładka
 * czytała rdzeń osobno, w jednej chwili pokazywałyby dwie migawki tej samej
 * telemetrii i Operator nie miałby jak orzec, która jest bieżąca. Odczyt jest
 * więc jeden, a zakładki różnią się wyłącznie tym, co z niego czytają:
 * Metrics & Performance przebieg procesu po procesie, Health & Uptime
 * zestawienie stanów.
 *
 * Zakres czasu wspólny modułowi tego odczytu nie dotyczy: `monitor.status`
 * i `monitor.subscribe` nie mają w kontrakcie pól `fromTime` ani `toTime`, więc
 * telemetria jest zawsze stanem bieżącym rejestru procesów. Zawężanie jej
 * zakresem po stronie okna byłoby zawężaniem pozorowanym.
 *
 * Obserwacja zakłada się raz, odczyt powtarza się: pierwsze pytanie idzie
 * `monitor.subscribe` (zapisuje okno na telemetrię), każde następne
 * `monitor.status` (czyta i niczego nie zapisuje). Pole `windowId` znaczy
 * w tych dwóch komendach co innego — w pierwszej jest oknem obserwatora,
 * w drugiej sitem procesów — więc do odczytu nie idzie wcale.
 *
 * Na żywo idzie `progress.changed`. Zdarzenie niesie stan, etap i stopień
 * ukończenia, lecz NIE niesie czasu zmiany — dlatego wiersz odświeżony
 * zdarzeniem pokazuje liczby ze zdarzenia, a czas nadal ten z odczytu, wraz
 * ze zdaniem, skąd się bierze. Podstawienie czasu klienta w miejsce czasu
 * rdzenia byłoby wpisaniem do telemetrii wartości, której rdzeń nie orzekł.
 *
 * Eksport telemetrii stoi przy danych, a nie w pasku akcji okna: potwierdzenie
 * eksportu wypisuje się w tej samej zakładce, w której widać eksportowany
 * materiał. Potwierdzenie na zakładce zasłoniętej nie byłoby potwierdzeniem.
 *
 * Wytwórnia sama nie czyta — pierwszy odczyt zleca złożenie modułu wywołaniem
 * `odswiez()`, tak samo jak w Recommendations Panelu. Odczyt w wytwórni obok
 * odczytu ze złożenia dałby dwa żądania na jedno zmontowanie okna.
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

  /**
   * Pobiera stan procesów. Obserwację zakłada się RAZ — przy pierwszym
   * odczycie — a każdy następny idzie `monitor.status`, bo powtarzanie zapisu
   * przy każdym odświeżeniu byłoby zapisywaniem tego samego okna na tę samą
   * telemetrię bez powodu.
   */
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
    // Zmiany sprzed żądania zawiera już migawka, którą rdzeń właśnie zbuduje;
    // zmiany doniesione PO wysłaniu żądania zostają, bo migawka może ich nie
    // objąć, a wyczyszczenie mapy dopiero przy odpowiedzi by je zgubiło.
    zmiany.clear();
    odczyt = 'w-toku';
    metryki.ladowanie('Odczyt telemetrii procesów…');
    kondycja.ladowanie('Odczyt telemetrii procesów…');
    void pobierz().then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Odmowa nie zostaje pustym wykazem: rdzeń bez wpiętego rejestru
        // telemetrii odmawia głośno, a wykaz pusty nazwałby to „brakiem
        // procesów”. Wykaz poprzedni też odchodzi — po odmowie nie jest już
        // bieżący, a eksport oddałby go jako gdyby był.
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
    // Zdarzenie o procesie spoza ostatniego odczytu też ma znaczenie: rejestr
    // rdzenia zna proces, którego okno jeszcze nie widziało. Wiersz powstaje
    // wtedy z samego zdarzenia i mówi to o sobie wprost.
    zmiany.set(tresc.processId, tresc);
    // Przerysowanie ma sens wyłącznie nad udanym odczytem. W trakcie odczytu
    // nie ma czego przerysowywać, a po odmowie wykaz zbudowany ze zdarzeń
    // zastąpiłby na ekranie treść odmowy — czyli ukryłby ją.
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
 * Żądanie obserwacji. `windowId` jest tu oknem ODBIERAJĄCYM telemetrię, nie
 * sitem procesów; sitem są `processIds` i `sessionId`. `processIds` nie idzie
 * wcale — pusty wykaz znaczy w kontrakcie komplet procesów, a to jest zakres
 * właściwy oknu, które patrzy na kondycję całej instalacji, nie jednego
 * przebiegu.
 */
function zadanieObserwacji(idOkna: string): MonitorSubscribeRequest {
  return idOkna === '' ? {} : { windowId: idOkna };
}

/**
 * Żądanie odczytu stanu — bez ani jednego pola.
 *
 * `windowId` znaczy w `monitor.status` co innego niż w `monitor.subscribe`:
 * tam jest oknem obserwatora, tutaj SITEM procesów. Podanie go zawęziłoby
 * odświeżenie do procesów okna Diagnostics, czyli do innego zbioru niż ten,
 * który oddał odczyt pierwszy — wykaz kurczyłby się po każdym odświeżeniu bez
 * żadnej zmiany w rdzeniu.
 */
function zadanieStanu(): MonitorStatusRequest {
  return {};
}

/** Zdanie o odczycie — liczba procesów oraz los samego zapisu obserwacji. */
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

/** Zakładka Metrics & Performance: pasek czynności nad miejscem treści. */
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

/** Zakładka Health & Uptime: zestawienie stanów wraz z granicą tego, co wiadomo. */
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

/** Wykaz procesów — perspektywa Metrics & Performance. */
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

/** Jeden wiersz procesu; zmiana ze zdarzenia bierze pierwszeństwo nad odczytem. */
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

/** Zdanie opisu procesu: etap, ukończenie, obiegi, okno i sesja. */
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

/** Zdanie o treści zdarzenia postępu — wyłącznie liczby, które rdzeń podał. */
function opisPostepu(zmiana: ProgressChangedEvent): string {
  const etap =
    zmiana.stepLabel === undefined || zmiana.stepLabel === '' ? '' : `${zmiana.stepLabel}, `;
  // `totalSteps` równe zeru znaczy w kontrakcie „liczba etapów nieznana”,
  // a nie „zero etapów” — ułamek z zerem w mianowniku byłby wtedy zmyśleniem.
  const kroki =
    zmiana.totalSteps === 0
      ? `krok ${zmiana.currentStep} z nieznanej liczby`
      : `krok ${zmiana.currentStep}/${zmiana.totalSteps}`;
  return `${etap}${kroki}, ukończenie ${zmiana.percent}%`;
}

/** Zestawienie stanów — perspektywa Health & Uptime. */
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

/** Kolejność stanów w zestawieniu: od pracy trwającej do zakończeń. */
const KOLEJNOSC_STANOW: readonly ProgressStatus[] = [
  ProgressStatus.Running,
  ProgressStatus.Pending,
  ProgressStatus.Paused,
  ProgressStatus.Stopped,
  ProgressStatus.Done,
  ProgressStatus.Failed,
];

/** Liczba procesów w każdym stanie; zmiana ze zdarzenia liczy się nad odczytem. */
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
 * Plakietka stanu procesu — ikona i etykieta obok barwy.
 *
 * Znak bierze się z katalogu motywu (`motyw/znaczenia-stanow.ts`), a nie
 * z własnej mapy okna: stan procesu jest tym samym pojęciem tutaj, na karcie
 * sesji i w monitorze pracy ciągłej, więc ma jedną nazwę i jeden znak.
 */
function plakietkaStanu(stan: ProgressStatus): HTMLElement {
  const znak = znakStanuPostepu(stan);
  const wariant = wariantPlakietki(znak.rodzina);
  const element = document.createElement('span');
  element.className = wariant === '' ? 'dn-plakietka' : `dn-plakietka ${wariant}`;
  element.append(elementIkony(znak.ikona, { rozmiar: 14 }), document.createTextNode(znak.etykieta));
  return element;
}
