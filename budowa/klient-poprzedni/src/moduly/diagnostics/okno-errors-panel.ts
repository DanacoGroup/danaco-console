import {
  DiagnosticErrorStatus,
  DiagnosticPriority,
  type DiagnosticAnalysis,
  type DiagnosticError,
  type DiagnosticsErrorListRequest,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  poleLiczbowe,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
  pozycjaWykazu,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza-braki';
import { utworzWyborZMenu, type WyborZMenu } from '../apps/wybor-z-menu';
import { powodBezKomendy } from './braki-kontraktu';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import type { StanDiagnostyki } from './stan-diagnostyki';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Errors Panel pokazuje odmowy wykonania komend zgłoszone przez rdzeń jako
 * wiersze `DiagnosticError`, w których pole `source` niesie nazwę odrzuconej
 * komendy, a zakres czasu wykazu pochodzi ze wspólnego stanu modułu.
 */
export interface OknoErrorsPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane` (zakres czasu wspólny modułu) założony przy montażu okna. */
  zamknij(): void;
}

export function utworzOknoErrorsPanel(
  zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
): OknoErrorsPanel {
  const rama = utworzRameOkna({
    tytul: 'Errors Panel',
    rola: 'pomocnicze',
    kod: 'errors-panel',
    przeznaczenie:
      'Przegląd błędu z kontekstem wystąpienia — źródło jest nazwą odrzuconej komendy, a licznik odróżnia jedną odmowę od trzystu.',
    przedrostek: 'dg',
  });
  const tresc = utworzStanTresci();
  const wykazBledow: WykazBledow = { pozycje: new Map<string, DiagnosticError>() };
  const powierzchnia = zlozPowierzchnieErrorsPanel(rama, tresc.element);

  const odczytaj = (): void => wczytajBledy(zrodlo, stan, powierzchnia, wykazBledow, tresc);
  const eksportuj = (): void => eksportujBledy(wykazBledow.pozycje, tresc);

  // Zapis migawki wyłącza subskrypcję na chwilę, żeby uniknąć drugiego odczytu.
  let odsubskrybujZakres = stan.naZmiane(odczytaj);
  function zapiszMigawkeBezPodwojnegoOdczytu(migawka: Parameters<StanDiagnostyki['ustawMigawke']>[0]): void {
    odsubskrybujZakres();
    stan.ustawMigawke(migawka);
    odsubskrybujZakres = stan.naZmiane(odczytaj);
  }
  const przekazDoAnalizy = (): void =>
    przekazBledyDoAnalizy(zrodlo, stan, wykazBledow, tresc, zapiszMigawkeBezPodwojnegoOdczytu);

  podepnijAkcjeErrorsPanel(powierzchnia, { odczytaj, eksportuj, przekazDoAnalizy });

  // Pierwszy odczyt zleca złożenie modułu wywołaniem `odswiez()`, nie wytwórnia.
  return { element: rama.element, odswiez: odczytaj, zamknij: () => odsubskrybujZakres() };
}

/**
 * Wykaz błędów widocznych w oknie wraz z licznością wszystkich pozycji
 * spełniających warunki, którą osobnym zapytaniem orzekł rdzeń.
 */
interface WykazBledow {
  pozycje: Map<string, DiagnosticError>;
  /** `total` z ostatniego UDANEGO odczytu; brak znaczy „rdzeń nie podał". */
  wszystkich?: number;
}

/**
 * Odczyt wykazu błędów. Odmowa NIE STAJE SIĘ pustym wykazem — okno mówi
 * wprost, co się stało, z kodem i treścią odmowy. Zdanie o powodzie
 * bierze się z drogi niepowodzenia, nie z założenia, że zawinił rdzeń.
 */
function wczytajBledy(
  zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
  filtry: FiltryErrorsPanel,
  wykaz: WykazBledow,
  tresc: StanTresci,
): void {
  tresc.ladowanie('Odczyt błędów…');
  void zrodlo.wykazBledow(zadanieBledow(stan, filtry)).then((wynik) => {
    wykaz.pozycje.clear();
    delete wykaz.wszystkich;
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(zdanieNiepowodzenia('odczytu błędów', wynik.powod), wynik.blad);
      return;
    }
    for (const blad of wynik.wynik.errors) wykaz.pozycje.set(blad.id, blad);
    if (wynik.wynik.total !== undefined) wykaz.wszystkich = wynik.wynik.total;
    rysujBledy([...wykaz.pozycje.values()], wykaz.wszystkich, tresc);
  });
}

/**
 * Eksportuje wykaz błędów widocznych w oknie jako plik Markdown; kontrakt
 * komendy eksportu nie niesie, więc plik składa i pobiera wyłącznie interfejs.
 */
function eksportujBledy(bledy: Map<string, DiagnosticError>, tresc: StanTresci): void {
  if (bledy.size === 0) {
    tresc.potwierdzenie('Nie ma czego wyeksportować — wykaz błędów jest pusty.', false);
    return;
  }
  pobierzPlik('errors-panel-raport.md', raportBledow([...bledy.values()]), 'text/markdown');
  tresc.potwierdzenie('Raport błędów pobrany jako plik Markdown.', true);
}

/**
 * Przekazuje błędy widoczne w oknie do analizy Diagnostics Center: prawdziwa
 * komenda `diagnostics.analyze.run` z `errorIds`, wynik ląduje w stanie
 * wspólnym modułu (`stan.ustawMigawke`), żeby Diagnostics Center go pokazał.
 */
function przekazBledyDoAnalizy(
  zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
  wykaz: WykazBledow,
  tresc: StanTresci,
  zapiszMigawke: (migawka: Parameters<StanDiagnostyki['ustawMigawke']>[0]) => void,
): void {
  if (wykaz.pozycje.size === 0) {
    tresc.potwierdzenie('Nie ma błędów widocznych w oknie — nie ma czego przekazać do analizy.', false);
    return;
  }
  const zakres = stan.zakres();
  const przekazane = [...wykaz.pozycje.keys()];
  tresc.ladowanie('Przekazuję widoczne błędy do analizy Diagnostics Center…');
  void zrodlo
    .uruchomAnalize({
      errorIds: przekazane,
      ...(zakres.od === undefined ? {} : { fromTime: zakres.od }),
      ...(zakres.do === undefined ? {} : { toTime: zakres.do }),
    })
    .then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(zdanieNiepowodzenia('uruchomienia analizy z tych błędów', wynik.powod), wynik.blad);
        return;
      }
      zapiszMigawke(wynik.wynik);
      rysujBledy([...wykaz.pozycje.values()], wykaz.wszystkich, tresc);
      const ocena = zdanieOPrzekazaniu(wynik.wynik, przekazane);
      tresc.potwierdzenie(ocena.zdanie, ocena.udane);
    });
}

/**
 * Zdanie potwierdzenia przekazania do analizy — mówi, co objęła analiza,
 * a nie co wysłało okno; wszystkie liczby zdania biorą się z migawki, którą
 * oddał rdzeń.
 */
function zdanieOPrzekazaniu(
  analiza: DiagnosticAnalysis,
  przekazane: readonly string[],
): { zdanie: string; udane: boolean } {
  const objete = analiza.errorIds;
  if (objete === undefined) {
    return {
      zdanie:
        `Analiza ${analiza.id} powstała, ale rdzeń nie podał, które błędy nią objął ` +
        `(przekazano ${przekazane.length}) — nie ma z czego orzec, czy wykaz z okna coś zmienił.`,
      udane: false,
    };
  }
  const pominiete = przekazane.filter((id) => !objete.includes(id));
  if (pominiete.length > 0 || objete.length !== przekazane.length) {
    return {
      zdanie:
        `Analiza ${analiza.id} objęła ${objete.length} błędów, choć z tego okna przekazano ${przekazane.length}` +
        `${pominiete.length === 0 ? '' : `, a ${pominiete.length} z przekazanych zostało poza nią`}. ` +
        'Rdzeń dobiera błędy zakresem czasu, nie wykazem z okna — Diagnostics Center pokazuje wynik faktyczny.',
      udane: false,
    };
  }
  return {
    zdanie:
      `Analiza ${analiza.id} objęła dokładnie te ${objete.length} błędów, które przekazano; ` +
      `rekomendacji: ${analiza.recommendationIds?.length ?? 0}. Diagnostics Center ją pokazuje.`,
    udane: true,
  };
}

/**
 * Przyciski paska akcji Errors Panel: odświeżenie wykazu, eksport do pliku
 * i przekazanie widocznych błędów do analizy Diagnostics Center.
 */
interface AkcjeErrorsPanel {
  odswiezPrzycisk: HTMLButtonElement;
  eksportPrzycisk: HTMLButtonElement;
  analizaPrzycisk: HTMLButtonElement;
}

/**
 * Składa pasek akcji okna. Kontrakt `diagnostics.error.list` jest wyłącznie
 * odczytem, więc panel może polami stanu, priorytetu i notatki błędu jedynie
 * filtrować, nigdy ich nadawać, a nieczynne pozycje niosą powód wprost.
 */
function zlozAkcjeErrorsPanel(gospodarz: HTMLElement): AkcjeErrorsPanel {
  const odswiezPrzycisk = przycisk('Odśwież błędy', 'dn-btn dn-btn--atrament');
  const eksportPrzycisk = przycisk('Eksportuj listę');
  const analizaPrzycisk = przycisk('Przekaż do Diagnostics Center', 'dn-btn dn-btn--atrament');

  gospodarz.append(
    odswiezPrzycisk,
    eksportPrzycisk,
    analizaPrzycisk,
    przyciskBezKomendy(
      'Analizuj',
      powodBezKomendy('Przestawienie stanu pojedynczego błędu na "analyzing" wymagałoby komendy zapisu.'),
    ),
    przyciskBezKomendy(
      'Rozwiąż',
      powodBezKomendy('Zamknięcie błędu wymagałoby komendy zapisu; error.list jest wyłącznie odczytem.'),
    ),
    przyciskBezKomendy(
      'Ignoruj',
      powodBezKomendy('Zignorowanie błędu wymagałoby komendy zapisu; error.list jest wyłącznie odczytem.'),
    ),
    przyciskBezKomendy(
      'Przypisz priorytet',
      powodBezKomendy('DiagnosticError niesie pole priority, lecz nadanie go wymagałoby komendy zapisu.'),
    ),
    przyciskBezKomendy(
      'Zapisz notatkę',
      powodBezKomendy('DiagnosticError niesie pole note, lecz zapis notatki Operatora wymagałby komendy zapisu.'),
    ),
    przyciskBezKomendy(
      'Przekaż do Developera',
      powodBezKomendy('Przekazanie błędu do modułu Developer wymagałoby komendy przyjmującej jego odniesienie.', 'developer'),
    ),
  );
  return { odswiezPrzycisk, eksportPrzycisk, analizaPrzycisk };
}

/**
 * Filtry okna: stan błędu, priorytet błędu i górna granica liczby pozycji
 * wykazu, którą przyjmuje żądanie odczytu skierowane do rdzenia.
 */
interface FiltryErrorsPanel {
  statusWybor: WyborZMenu;
  priorytetWybor: WyborZMenu;
  granica: HTMLInputElement;
}

const OPCJE_STANU = [
  { wartosc: '', etykieta: 'Każdy stan' },
  { wartosc: DiagnosticErrorStatus.New, etykieta: 'Nowy' },
  { wartosc: DiagnosticErrorStatus.Analyzing, etykieta: 'W analizie' },
  { wartosc: DiagnosticErrorStatus.Resolved, etykieta: 'Rozwiązany' },
  { wartosc: DiagnosticErrorStatus.Ignored, etykieta: 'Zignorowany' },
] as const;

const OPCJE_PRIORYTETU = [
  { wartosc: '', etykieta: 'Każdy priorytet' },
  { wartosc: DiagnosticPriority.Critical, etykieta: 'Krytyczny' },
  { wartosc: DiagnosticPriority.High, etykieta: 'Wysoki' },
  { wartosc: DiagnosticPriority.Medium, etykieta: 'Średni' },
  { wartosc: DiagnosticPriority.Low, etykieta: 'Niski' },
] as const;

/**
 * Kontrolki okna: filtry stanu i priorytetu, pasek akcji oraz ciało ramy,
 * do którego trafia wykaz błędów albo zdanie o stanie odczytu.
 */
interface PowierzchniaErrorsPanel extends AkcjeErrorsPanel, FiltryErrorsPanel {}

function zlozPowierzchnieErrorsPanel(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaErrorsPanel {
  // Rozwijanie z biblioteki `komponenty/menu-drzewo.ts`, nie natywny `<select>`.
  const statusWybor = utworzWyborZMenu('Stan błędu', OPCJE_STANU);
  const priorytetWybor = utworzWyborZMenu('Priorytet błędu', OPCJE_PRIORYTETU);
  const granica = poleLiczbowe('Górna granica liczby błędów', 'domyślnie wszystkie');
  const akcje = zlozAkcjeErrorsPanel(rama.akcje);

  rama.narzedzia.append(
    wiersz('Stan', statusWybor.element, {
      klasa: 'dg-wiersz',
      objasnienie: 'Filtruje po polu status.',
    }),
    wiersz('Priorytet', priorytetWybor.element, {
      klasa: 'dg-wiersz',
      objasnienie: 'Filtruje po polu priority.',
    }),
    wiersz('Granica', granica, {
      klasa: 'dg-wiersz',
      objasnienie: 'Puste pole nie narzuca granicy z okna; obowiązuje wtedy granica rdzenia. Ile błędów spełnia warunki, a ile rdzeń oddał, mówi zdanie nad wykazem.',
    }),
  );
  rama.cialo.append(stanTresci);
  return { statusWybor, priorytetWybor, granica, ...akcje };
}

function podepnijAkcjeErrorsPanel(
  powierzchnia: PowierzchniaErrorsPanel,
  obsluga: { odczytaj: () => void; eksportuj: () => void; przekazDoAnalizy: () => void },
): void {
  powierzchnia.odswiezPrzycisk.addEventListener('click', obsluga.odczytaj);
  powierzchnia.eksportPrzycisk.addEventListener('click', obsluga.eksportuj);
  powierzchnia.analizaPrzycisk.addEventListener('click', obsluga.przekazDoAnalizy);
  powierzchnia.statusWybor.naZmiane(obsluga.odczytaj);
  powierzchnia.priorytetWybor.naZmiane(obsluga.odczytaj);
  powierzchnia.granica.addEventListener('change', obsluga.odczytaj);
}

/**
 * Zadanie odczytu błędów — zakres czasu ZAWSZE ze wspólnego stanu modułu,
 * filtry stanu/priorytetu i granica z kontrolek lokalnych okna.
 */
function zadanieBledow(
  stan: StanDiagnostyki,
  filtry: FiltryErrorsPanel,
): DiagnosticsErrorListRequest {
  const zadanie: DiagnosticsErrorListRequest = {};
  const zakres = stan.zakres();
  if (zakres.od !== undefined) zadanie.fromTime = zakres.od;
  if (zakres.do !== undefined) zadanie.toTime = zakres.do;
  if (filtry.statusWybor.wartosc() !== '') {
    zadanie.status = filtry.statusWybor.wartosc() as DiagnosticErrorStatus;
  }
  if (filtry.priorytetWybor.wartosc() !== '') {
    zadanie.priority = filtry.priorytetWybor.wartosc() as DiagnosticPriority;
  }
  const limit = Number.parseInt(filtry.granica.value, 10);
  if (Number.isInteger(limit) && limit > 0) zadanie.limit = limit;
  return zadanie;
}

/**
 * Rysuje wykaz błędów, stan pustki albo notatkę o wykazie uciętym; rozjazd
 * między liczbą `total` a liczbą oddanych pozycji jest pokazany wprost.
 */
function rysujBledy(bledy: readonly DiagnosticError[], total: number | undefined, tresc: StanTresci): void {
  const posortowane = [...bledy].sort((a, b) => b.lastSeenAt - a.lastSeenAt);
  if (posortowane.length === 0) {
    // Pustka jest stanem poprawnym, odrębnym od stanu błędu odczytu.
    tresc.pusto('Żaden błąd nie spełnia warunków w wybranym zakresie i filtrach.');
    return;
  }
  const cialo = tresc.tresc();
  if (total !== undefined && total > posortowane.length) {
    const uwaga = document.createElement('p');
    uwaga.className = 'dg-stan';
    uwaga.textContent = `Rdzeń zgłasza ${total} błędów spełniających warunki — oddano ${posortowane.length}. Zawęź filtry albo podnieś granicę.`;
    cialo.append(uwaga);
  }
  cialo.append(listaBledow(posortowane));
}

/**
 * Buduje wykaz błędów jako czystą konstrukcję z danych, wraz z czasem
 * wystąpienia i licznikiem powtórzeń każdej pozycji.
 */
function listaBledow(bledy: readonly DiagnosticError[]): HTMLElement {
  const lista = wykaz('Błędy zgłoszone przez rdzeń', 'dg-wykaz');
  for (const blad of bledy) {
    // `source` jest nazwą odrzuconej komendy, nie opisem — stoi w tytule wprost.
    const pozycja = pozycjaWykazu(blad.source ?? '(źródło nieznane)', opisBledu(blad), 'dg');
    pozycja.element.dataset['priorytet'] = blad.priority;
    pozycja.element.dataset['stanBledu'] = blad.status;

    const czas = document.createElement('span');
    czas.className = 'dg-pozycja__czas';
    czas.textContent = `${new Date(blad.firstSeenAt).toLocaleString('pl-PL')} → ${new Date(blad.lastSeenAt).toLocaleString('pl-PL')}`;
    pozycja.element.append(czas);

    const licznik = document.createElement('span');
    licznik.className = 'dg-wpis__licznik';
    // Brak `occurrences` nie znaczy jedną odmowę — rdzeń nie podał licznika.
    licznik.textContent = opisLicznika(blad.occurrences);
    pozycja.akcje.append(licznik);

    lista.append(pozycja.element);
  }
  return lista;
}

/**
 * Tekst licznika wystąpień błędu — brak pola jest powiedziany jako brak
 * liczby, a nie zastąpiony domyślną wartością jednego wystąpienia.
 */
function opisLicznika(occurrences: number | undefined): string {
  return occurrences === undefined ? 'liczba wystąpień nieznana' : `×${occurrences}`;
}

/**
 * Kod odmowy wiersza błędu — czytany z pola kontraktu `errorCode`, a przy
 * jego braku z `context.errorCode`, gdzie kod wkłada dzisiejszy rdzeń.
 */
function kodOdmowy(blad: DiagnosticError): string | undefined {
  if (blad.errorCode !== undefined && blad.errorCode !== '') return blad.errorCode;
  const context = blad.context;
  if (typeof context !== 'object' || context === null) return undefined;
  const kod = (context as Record<string, unknown>)['errorCode'];
  return typeof kod === 'string' && kod !== '' ? kod : undefined;
}

/**
 * Zdanie opisu pozycji wykazu: treść błędu, kod odmowy, odcisk, notatka
 * i pozostały kontekst zapisu, w tej kolejności.
 */
function opisBledu(blad: DiagnosticError): string {
  const czesci = [blad.message];
  const kod = kodOdmowy(blad);
  if (kod !== undefined) czesci.push(`kod odmowy: ${kod}`);
  czesci.push(`odcisk ${blad.fingerprint}`);
  if (blad.note !== undefined && blad.note !== '') czesci.push(`notatka: ${blad.note}`);
  if (blad.commitId !== undefined) czesci.push(`zatwierdzenie ${blad.commitId}`);
  if (blad.deploymentId !== undefined) czesci.push(`wdrożenie ${blad.deploymentId}`);
  if (blad.context !== undefined) czesci.push(`kontekst: ${JSON.stringify(blad.context)}`);
  return czesci.join(' — ');
}

/**
 * Raport błędów w formacie Markdown — treść pliku, który pobiera eksport
 * wykazu błędów widocznego w oknie Errors Panel.
 */
function raportBledow(bledy: readonly DiagnosticError[]): string {
  const wiersze = ['# Raport błędów Errors Panel', ''];
  for (const blad of bledy) {
    wiersze.push(
      `## ${blad.source ?? '(źródło nieznane)'} — ${blad.fingerprint}`,
      `- Stan: ${blad.status}; priorytet: ${blad.priority}; wystąpienia: ${opisLicznika(blad.occurrences)}`,
      `- ${blad.message}`,
      `- Pierwsze wystąpienie: ${new Date(blad.firstSeenAt).toLocaleString('pl-PL')}`,
      `- Ostatnie wystąpienie: ${new Date(blad.lastSeenAt).toLocaleString('pl-PL')}`,
      blad.note === undefined || blad.note === '' ? '' : `- Notatka: ${blad.note}`,
      '',
    );
  }
  return wiersze.filter((wiersz) => wiersz !== '').join('\n');
}
