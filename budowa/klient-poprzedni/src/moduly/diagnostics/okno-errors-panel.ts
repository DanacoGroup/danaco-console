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
 * Errors Panel — okno pomocnicze modułu Diagnostics.
 *
 * Port `Diagnostyka` przyjmuje w rdzeniu odmowy wykonania komend; każda odmowa
 * staje się wierszem `DiagnosticError`, w którym `source` niesie nazwę
 * odrzuconej komendy. Stan pusty i stan odmowy odczytu są tu rozróżnione, bo
 * okno jest jedynym miejscem w produkcie, gdzie widać odrzucenie komendy.
 *
 * Kod odmowy nie leży w `message`, tylko osobno: w polu `errorCode` błędu,
 * a w dzisiejszym rdzeniu jeszcze w `context.errorCode`. Po nim rozpoznaje się
 * `*.unknown`, więc `opisBledu` wydobywa go wprost, nie tylko zrzuca `context`
 * jako JSON.
 *
 * Zakres czasu bierze się ze wspólnego stanu modułu (`stan.zakres()`) przez
 * `stan.naZmiane(...)`; okno nie prowadzi drugiego.
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

  // Subskrypcja jedzie przez zmienną, żeby `zapiszMigawkeBezPodwojnegoOdczytu`
  // mogła ją na chwilę wyłączyć: `stan.ustawMigawke` woła `powiadom()`, które
  // bez tego obudziłoby WŁASNY `odczytaj()` okna — miganie stanem ładowania
  // obok świeżo pokazanego potwierdzenia i zbędne drugie żądanie.
  let odsubskrybujZakres = stan.naZmiane(odczytaj);
  function zapiszMigawkeBezPodwojnegoOdczytu(migawka: Parameters<StanDiagnostyki['ustawMigawke']>[0]): void {
    odsubskrybujZakres();
    stan.ustawMigawke(migawka);
    odsubskrybujZakres = stan.naZmiane(odczytaj);
  }
  const przekazDoAnalizy = (): void =>
    przekazBledyDoAnalizy(zrodlo, stan, wykazBledow, tresc, zapiszMigawkeBezPodwojnegoOdczytu);

  podepnijAkcjeErrorsPanel(powierzchnia, { odczytaj, eksportuj, przekazDoAnalizy });

  // Wytwórnia okna nie czyta sama — pierwszy odczyt zleca złożenie modułu
  // wywołaniem `odswiez()`, tak jak w pozostałych oknach Diagnostics. Odczyt
  // w wytwórni obok odczytu ze złożenia dawał dwa żądania `error.list` na
  // jedno zmontowanie modułu.
  return { element: rama.element, odswiez: odczytaj, zamknij: () => odsubskrybujZakres() };
}

/**
 * Wykaz widoczny w oknie wraz z licznością, którą orzekł rdzeń.
 *
 * `wszystkich` trzyma się przy pozycjach, a nie w miejscu wywołania, bo
 * przerysowanie wykazu bez nowego odczytu (po przekazaniu błędów do analizy)
 * musi powtórzyć tę samą liczbę — inaczej ostrzeżenie o wykazie uciętym
 * znika po czynności, która niczego w wykazie nie zmieniła.
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

/** Eksport wykazu widocznego w oknie jako plik Markdown (bez pokrycia w kontrakcie, klienckie). */
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
 * Zdanie potwierdzenia przekazania — mówi, co objęła analiza, a nie co wysłało
 * okno.
 *
 * Rdzeń `errorIds` z żądania nie używa: `UruchomAnalize` dobiera błędy
 * wyłącznie zakresem czasu, więc liczba wzięta z żądania przeczyłaby temu, co
 * w tej samej chwili pokazuje Diagnostics Center. Wszystkie liczby zdania
 * biorą się z migawki oddanej przez rdzeń, a rozjazd między wykazem wysłanym
 * a objętym jest powiedziany wprost i tonem nieudanym.
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

/** Przyciski paska akcji Errors Panel. */
interface AkcjeErrorsPanel {
  odswiezPrzycisk: HTMLButtonElement;
  eksportPrzycisk: HTMLButtonElement;
  analizaPrzycisk: HTMLButtonElement;
}

/**
 * Składa pasek akcji okna. Kontrakt `diagnostics.error.list` jest WYŁĄCZNIE
 * odczytem — `DiagnosticError` niesie pola `status`, `priority`, `note`, ale
 * komendy zapisu nie ma. Panel może po tych polach FILTROWAĆ, nigdy ich
 * NADAWAĆ, więc trzy czynności z inwentarza są jawnie bez pokrycia.
 *
 * POWÓD KAŻDEJ NIECZYNNEJ POZYCJI SKŁADA SIĘ Z KONTRAKTU, a nie z napisu
 * (`braki-kontraktu.ts`): zdanie o braku komendy ma się zmienić samo w dniu,
 * w którym komenda się pojawi, bo napisu nikt wtedy nie zdejmie.
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

/** Filtry okna: stan, priorytet i górna granica wykazu. */
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

/** Kontrolki okna: filtry, pasek akcji i ciało ramy. */
interface PowierzchniaErrorsPanel extends AkcjeErrorsPanel, FiltryErrorsPanel {}

function zlozPowierzchnieErrorsPanel(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaErrorsPanel {
  // Rozwijanie z biblioteki, nie natywny `<select>`: `komponenty/menu-drzewo.ts`
  // przez obsadę `moduly/apps/wybor-z-menu.ts`.
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
      // NIE „puste pole zwraca wszystkie" — tego okno nie wie i wiedzieć nie
      // może. Rdzeń ma własną granicę domyślną, a jedyną prawdą o liczności
      // jest `total` z odpowiedzi, wypisany nad wykazem.
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
 * Rysuje wykaz błędów, stan pustki albo notatkę o wykazie uciętym.
 *
 * `total` większe od liczby oddanych błędów jest widoczne wprost — wykaz
 * ucięty bez ostrzeżenia jest kłamstwem tej samej rodziny co pusty wykaz
 * przy odmowie. Rozjazd jest osiągalny: rdzeń liczy `total` osobnym
 * zapytaniem, które bierze stan, priorytet i zakres czasu, a granicy nie
 * bierze wcale — wykaz ucina granica z pola „Granica" albo granica domyślna
 * rdzenia.
 */
function rysujBledy(bledy: readonly DiagnosticError[], total: number | undefined, tresc: StanTresci): void {
  const posortowane = [...bledy].sort((a, b) => b.lastSeenAt - a.lastSeenAt);
  if (posortowane.length === 0) {
    // Pustka jest stanem poprawnym — instalacja bez błędów w zakresie czasu
    // nie jest usterką, i to jest zdanie inne niż stan błędu odczytu.
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

/** Buduje wykaz błędów — czysta konstrukcja z danych. */
function listaBledow(bledy: readonly DiagnosticError[]): HTMLElement {
  const lista = wykaz('Błędy zgłoszone przez rdzeń', 'dg-wykaz');
  for (const blad of bledy) {
    // `source` jest nazwą odrzuconej komendy (np. library.collection.create),
    // nie ozdobnym opisem — dlatego stoi w tytule pozycji wprost.
    const pozycja = pozycjaWykazu(blad.source ?? '(źródło nieznane)', opisBledu(blad), 'dg');
    pozycja.element.dataset['priorytet'] = blad.priority;
    pozycja.element.dataset['stanBledu'] = blad.status;

    const czas = document.createElement('span');
    czas.className = 'dg-pozycja__czas';
    czas.textContent = `${new Date(blad.firstSeenAt).toLocaleString('pl-PL')} → ${new Date(blad.lastSeenAt).toLocaleString('pl-PL')}`;
    pozycja.element.append(czas);

    const licznik = document.createElement('span');
    licznik.className = 'dg-wpis__licznik';
    // Brak `occurrences` nie znaczy „jedna odmowa" — rdzeń nie podał
    // licznika. Podstawienie jedynki byłoby atrapą danych.
    licznik.textContent = opisLicznika(blad.occurrences);
    pozycja.akcje.append(licznik);

    lista.append(pozycja.element);
  }
  return lista;
}

/** Tekst licznika wystąpień — brak pola jest powiedziany jako brak, nie jako `×1`. */
function opisLicznika(occurrences: number | undefined): string {
  return occurrences === undefined ? 'liczba wystąpień nieznana' : `×${occurrences}`;
}

/**
 * Kod odmowy wiersza błędu — pole kontraktu przed kontekstem zapisu.
 *
 * Kontrakt trzyma kod w polu własnym błędu (`DiagnosticError.errorCode`):
 * odmowa komendy jest faktem o samym błędzie, nie o okolicznościach jego
 * zapisu. Dzisiejszy rdzeń tego pola jeszcze nie wypełnia i wkłada kod do
 * `context.errorCode` — dlatego czytane są oba miejsca, w tej kolejności.
 * Sam odczyt kontekstu nie jest tu obejściem: gdy rdzeń zacznie wypełniać pole
 * kontraktu, gałąź zapasowa przestanie być osiągalna sama z siebie, bez zmiany
 * w tym pliku.
 */
function kodOdmowy(blad: DiagnosticError): string | undefined {
  if (blad.errorCode !== undefined && blad.errorCode !== '') return blad.errorCode;
  const context = blad.context;
  if (typeof context !== 'object' || context === null) return undefined;
  const kod = (context as Record<string, unknown>)['errorCode'];
  return typeof kod === 'string' && kod !== '' ? kod : undefined;
}

/** Zdanie opisu pozycji: treść, kod odmowy, odcisk i kontekst zapisu. */
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

/** Raport błędów w Markdown — treść pliku eksportu. */
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
