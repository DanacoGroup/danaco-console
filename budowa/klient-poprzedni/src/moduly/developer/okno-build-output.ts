import {
  ProblemSeverity,
  type DeveloperBuild,
  type DeveloperBuildChangedEvent,
  type DeveloperBuildRunRequest,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza-braki';
import type { Wynik } from '../../protokol/kanal';
import { powodBezKomendy } from './braki-kontraktu';
import {
  KATALOG_W_ODCZYCIE,
  odczytajKatalogKomend,
  zdanieHistoriiBudowania,
} from './katalog-komend';
import { utworzPanelHistoriiBudowan } from './panel-historii-budowan';
import { utworzPanelRunDebug } from './panel-run-debug';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { rysujPrzebieg, type ZawezeniePrzebiegu } from './widok-przebiegu';
import { zdaniePrzerwania, zdanieUruchomienia } from './zdania-odpowiedzi';
import { rysujZaleznosci, zaleznosci } from './zaleznosci-zewnetrzne';
import { cialoZakladki, utworzZakladkiOkna } from './zakladki-okna';
import type { ZrodloDeveloper } from './zrodlo-developer';

/**
 * Build Output i Run & Debug — okno monitora modułu Developer.
 *
 * Opracowanie modułu opisuje monitor jako JEDNO okno o dwóch częściach
 * (rozdz. 2, pozycja 6; rozdz. 3.6): Build Output prowadzi budowanie, testy
 * i log, Run & Debug prowadzi konfiguracje uruchomień oraz debugger krokowy,
 * a obie części zajmują tę samą kolumnę i przełącza je pas zakładek w jej
 * nagłówku. Stąd jedna rama, dwie zakładki i jeden kod okna `build-output`
 * z katalogu rdzenia — Run & Debug nie jest osobnym oknem i osobnego wiersza
 * katalogu nie dostaje.
 *
 * Aktualizacja na żywo idzie ze zdarzenia. Okno subskrybuje
 * `developer.build.changed` bezpośrednio u źródła, mimo że `stan-developer.ts`
 * subskrybuje to samo zdarzenie i budzi okna przez `stan.naZmiane(...)`: stan
 * tylko rozgłasza „coś się zmieniło" po filtrze `windowId`, nie niesie
 * `logLine` ani przebiegu z treści zdarzenia, a kontrakt nie ma komendy, którą
 * dałoby się dogonić stan inaczej. Ta subskrypcja bierze więc co innego niż
 * subskrypcja stanu i nie jest powieleniem tej samej pracy.
 *
 * Zgłoszenie przebiegu prowadzi do pliku: zgłoszenia niosą `path` i `line`,
 * a każde ze ścieżką ma przejście „Otwórz w edytorze” — woła `stan.wskazPlik`,
 * a Code Editor otwiera plik sam.
 *
 * Log jest ucięty przy otwarciu okna w trakcie przebiegu. Rośnie wyłącznie
 * z `logLine` kolejnych zdarzeń, więc okno otwarte po starcie przebiegu (albo
 * przebiegu uruchomionego przez inne okno tego konta) widzi tylko ogon.
 * Ucięcie jest oznaczone wprost w treści okna, nie zamaskowane.
 *
 * Szukanie w logu i zawężanie zgłoszeń wagą są czynnościami wyłącznie
 * klienckimi nad materiałem już zebranym — kontrakt nie ma komendy, którą
 * dałoby się dopytać rdzeń o wiersze pominięte, i widok mówi to wprost.
 * Waga zgłoszenia pochodzi z pola `severity` oddanego przez rdzeń, nie
 * z rozpoznawania treści wiersza logu.
 *
 * Zdanie stanu pustego o tym, czym rozporządza rdzeń, składa `katalog-komend.ts`
 * z rejestru komend wziętego z `connection.hello` — z tego, co rdzeń
 * zarejestrował, a nie z napisu na stałe.
 */
export interface OknoBudowania {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka subskrypcję `developer.build.changed` założoną przez to okno. */
  zamknij(): void;
}

export function utworzOknoBudowania(
  zrodlo: ZrodloDeveloper,
  warsztat: ZrodloWarsztatu,
  stan: StanDevelopera,
): OknoBudowania {
  const rama = utworzRameOkna({
    tytul: 'Build Output i Run & Debug',
    rola: 'monitor',
    kod: 'build-output',
    przeznaczenie:
      'Uruchomienie i przerwanie budowania; log na żywo, wynik i zgłoszenia przebiegu. ' +
      'Druga zakładka prowadzi konfiguracje uruchomień i debugger krokowy.',
    modul: 'Developer',
    przedrostek: 'mdev',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieBudowania(tresc.element);
  const kontekst = utworzKontekstBudowania();

  // Trzy części kolumny monitora, bo trzy różne pytania: co się teraz dzieje
  // (Build Output), jak sterować zatrzymanym programem (Run & Debug) i co
  // zostało po przebiegach zakończonych (Historia i pomiary).
  const runDebug = utworzPanelRunDebug(warsztat, stan);
  const historia = utworzPanelHistoriiBudowan(warsztat, stan);
  const zakladki = utworzZakladkiOkna('Części kolumny monitora', [
    { kod: 'build-output', nazwa: 'Build Output', element: powierzchnia.obszar },
    { kod: 'run-debug', nazwa: 'Run & Debug', element: runDebug.element },
    { kod: 'historia', nazwa: 'Historia i pomiary', element: historia.element },
  ]);
  rama.pasek.append(zakladki.pasek);
  rama.cialo.append(zakladki.obszary);

  function rysuj(): void {
    if (kontekst.przebieg === null) {
      tresc.pusto(kontekst.zdanieHistorii);
      return;
    }
    tresc.tresc().append(
      rysujPrzebieg(kontekst.przebieg, {
        logi: kontekst.logi,
        uciety: kontekst.uciety,
        zawezenie: zawezenieZPol(powierzchnia),
        naWskazanie: (sciezka) => stan.wskazPlik(sciezka),
      }),
    );
  }

  // Zdanie o niewidocznej historii przychodzi Z RDZENIA, a nie z napisu — patrz
  // `katalog-komend.ts`. Odczyt idzie raz, a jego wynik odświeża stan pusty.
  void odczytajKatalogKomend(zrodlo).then((katalog) => {
    kontekst.zdanieHistorii = zdanieHistoriiBudowania(katalog);
    if (kontekst.przebieg === null) rysuj();
  });

  function uruchom(): void {
    const zadanie = zadanieUruchomienia(stan.okno(), powierzchnia.zadanie.value, powierzchnia.argumenty.value);
    if (zadanie === null) {
      tresc.blad('Podaj zadanie budowania (np. „test” albo „build”) przed uruchomieniem.');
      return;
    }
    tresc.ladowanie('Uruchamianie budowania…');
    // Log czyścimy TERAZ, przed odpowiedzią — nie w `.then()`. Zdarzenie tego
    // przebiegu mogące przyjść wcześniej niż odpowiedź trafia więc do logu już
    // pustego, zamiast zostać skasowane późniejszym czyszczeniem po odpowiedzi.
    kontekst.logi.length = 0;
    kontekst.oczekujeWlasnegoStartu = true;
    kontekst.ostatnieZadanie = zadanie;
    void zrodlo.budowanie(zadanie).then((wynik) => obsluzOdpowiedzUruchomienia(kontekst, tresc, wynik, rysuj));
  }

  /**
   * Ponowne uruchomienie przez komendę `developer.build.run` z zadaniem, które
   * okno już zna.
   *
   * Powtórzenie znaczy „to samo żądanie”, a nie „to, co teraz stoi w polach” —
   * od tego jest przycisk uruchomienia. Gdy okno nie wysłało jeszcze niczego
   * (przebieg zaczęło inne okno tego konta), powtarza samo zadanie z migawki
   * rdzenia i mówi, że parametrów nie zna, bo `DeveloperBuild` ich nie niesie.
   */
  function ponow(): void {
    const powtorzenie = zadaniePonowienia(stan.okno(), kontekst);
    if (powtorzenie === null) {
      tresc.blad('Nie ma czego powtórzyć — to okno nie zna jeszcze żadnego przebiegu budowania.');
      return;
    }
    tresc.ladowanie(`Ponowne uruchomienie zadania „${powtorzenie.zadanie.task}”…`);
    kontekst.logi.length = 0;
    kontekst.oczekujeWlasnegoStartu = true;
    kontekst.ostatnieZadanie = powtorzenie.zadanie;
    void zrodlo
      .budowanie(powtorzenie.zadanie)
      .then((wynik) =>
        obsluzOdpowiedzUruchomienia(kontekst, tresc, wynik, rysuj, powtorzenie.dopisek),
      );
  }

  // Przycisk zatrzymania nie ma warunku. Czynny bez odczytu i po zakończeniu
  // przebiegu; odmowa rdzenia (np. „nic nie trwa”) ląduje w stanie błędu okna,
  // nie w wyszarzeniu przycisku.
  function przerwij(): void {
    const zadanie = zadanieZatrzymania(stan.okno(), powierzchnia.zadanie.value, kontekst.przebieg);
    tresc.ladowanie('Przerywanie budowania…');
    void zrodlo.budowanie(zadanie).then((wynik) => obsluzOdpowiedzPrzerwania(kontekst, tresc, wynik, rysuj));
  }

  zlozAkcjeBudowania(rama.akcje, {
    uruchom,
    przerwij,
    ponow,
    eksportujLog: () => eksportujLogBudowania(tresc, kontekst.przebieg, kontekst.logi, kontekst.uciety),
  });

  // Zawężenie przelicza wyłącznie widok — nic nie jedzie do rdzenia, więc
  // przerysowanie jest całą jego obsługą. Warunek na stan „treść” chroni
  // komunikat ładowania i odmowy przed startem: zawężenie nie ma prawa
  // zamienić błędu odczytu w pusty log.
  function przeliczZawezenie(): void {
    if (kontekst.przebieg !== null) rysuj();
  }
  powierzchnia.szukaj.addEventListener('input', przeliczZawezenie);
  powierzchnia.waga.addEventListener('change', przeliczZawezenie);

  rysuj();

  // Druga subskrypcja tego samego zdarzenia co `stan-developer.ts` — celowo:
  // bierze treść zdarzenia (przebieg + logLine), której `stan.naZmiane()` nie
  // przekazuje. Filtr po `windowId` chroni przed wpisami cudzego okna modułu.
  const odsubskrybuj = zrodlo.naZmianeBudowania((zdarzenie) =>
    obsluzZdarzenieBudowania(kontekst, stan.okno(), zdarzenie, rysuj));

  return {
    element: rama.element,
    odswiez: rysuj,
    zamknij: () => {
      // Panel debugowania zamyka się razem z oknem: sesja debugowania jest
      // stanem żywym i uchwyt do niej nie ma prawa przeżyć zejścia ze sceny.
      runDebug.zamknij();
      odsubskrybuj();
    },
  };
}

/**
 * Stan zmienny okna, trzymany jednym obiektem, żeby funkcje obsługi odpowiedzi
 * i zdarzenia mogły stać poza wytwórnią okna i dzielić tę samą prawdę bez
 * domykania się na jej zmiennych.
 */
interface KontekstBudowania {
  przebieg: DeveloperBuild | null;
  idZnanyOdPoczatku: string | null;
  /**
   * Czynne od chwili wysłania `run` do chwili związania id — rdzeń NIE
   * gwarantuje, że odpowiedź przyjdzie przed pierwszym zdarzeniem tego
   * samego przebiegu, więc id wiąże cokolwiek przyjdzie pierwsze.
   */
  oczekujeWlasnegoStartu: boolean;
  uciety: boolean;
  logi: string[];
  /** Zdanie stanu pustego złożone z rejestru komend rdzenia (`katalog-komend.ts`). */
  zdanieHistorii: string;
  /**
   * Ostatnie żądanie uruchomienia WYSŁANE PRZEZ TO OKNO — podstawa ponowienia.
   *
   * Rdzeń parametrów przebiegu nie oddaje (`DeveloperBuild` niesie `task`, nie
   * `arguments`), więc wierne powtórzenie zna wyłącznie okno, które je wysłało.
   */
  ostatnieZadanie: DeveloperBuildRunRequest | null;
}

function utworzKontekstBudowania(): KontekstBudowania {
  return {
    przebieg: null,
    idZnanyOdPoczatku: null,
    oczekujeWlasnegoStartu: false,
    uciety: false,
    logi: [],
    zdanieHistorii: KATALOG_W_ODCZYCIE,
    ostatnieZadanie: null,
  };
}

/** Powtórzenie wraz ze zdaniem o tym, czego rdzeń o powtarzanym przebiegu nie mówi. */
interface Powtorzenie {
  zadanie: DeveloperBuildRunRequest;
  /** Pusty, gdy powtórzenie jest wierne; inaczej mówi, czego zabrakło. */
  dopisek: string;
}

/**
 * Składa żądanie ponowienia — wiernie z żądania wysłanego przez to okno, a gdy
 * takiego nie było, z migawki przebiegu, którą okno zna, i wtedy NIE UDAJE
 * wierności.
 *
 * `null` znaczy „nie ma czego powtórzyć” i jest odróżnione od powtórzenia
 * przybliżonego: pierwsze to odmowa, drugie to czynność z zastrzeżeniem.
 */
function zadaniePonowienia(idOkna: string, kontekst: KontekstBudowania): Powtorzenie | null {
  if (kontekst.ostatnieZadanie !== null) {
    return { zadanie: { ...kontekst.ostatnieZadanie, windowId: idOkna }, dopisek: '' };
  }
  if (kontekst.przebieg === null) return null;
  return {
    zadanie: { windowId: idOkna, task: kontekst.przebieg.task },
    dopisek:
      ' Powtórzenie jest przybliżone: przebieg uruchomiło inne okno tego konta, a rdzeń nie oddaje ' +
      'parametrów zadania — poszło samo zadanie, bez nich.',
  };
}

/** Wiąże pierwszą wiadomość (odpowiedź albo zdarzenie) własnego startu z jego id — cokolwiek przyjdzie pierwsze. */
function powiazZWlasnymStartem(kontekst: KontekstBudowania, id: string): void {
  if (!kontekst.oczekujeWlasnegoStartu) return;
  kontekst.idZnanyOdPoczatku = id;
  kontekst.uciety = false;
  kontekst.oczekujeWlasnegoStartu = false;
}

/** Odpowiedź `developer.build.run` (uruchomienie) — poza wytwórnią, bierze kontekst wprost. */
function obsluzOdpowiedzUruchomienia(
  kontekst: KontekstBudowania,
  tresc: StanTresci,
  wynik: Wynik<DeveloperBuild>,
  rysuj: () => void,
  dopisek = '',
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    kontekst.oczekujeWlasnegoStartu = false;
    tresc.blad('Rdzeń odmówił uruchomienia budowania.', wynik.blad);
    return;
  }
  powiazZWlasnymStartem(kontekst, wynik.wynik.id);
  kontekst.przebieg = scalPrzebieg(kontekst.przebieg, wynik.wynik);
  // Zdanie z przebiegu SCALONEGO, nie z samej odpowiedzi: zdarzenie potrafi
  // wyprzedzić odpowiedź i donieść stan końcowy, a wtedy „uruchomione” byłoby
  // zdaniem o czymś, co już się skończyło.
  const potwierdzenie = zdanieUruchomienia(kontekst.przebieg);
  tresc.potwierdzenie(potwierdzenie.zdanie + dopisek, potwierdzenie.udane);
  rysuj();
}

/** Odpowiedź `developer.build.run` (przerwanie) — poza wytwórnią, bierze kontekst wprost. */
function obsluzOdpowiedzPrzerwania(
  kontekst: KontekstBudowania,
  tresc: StanTresci,
  wynik: Wynik<DeveloperBuild>,
  rysuj: () => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    tresc.blad('Rdzeń odmówił przerwania budowania.', wynik.blad);
    return;
  }
  // Czy przerywać BYŁO CZEGO — rdzeń tego nie powie, bo na przebieg biegnący
  // i na dawno domknięty odpowiada tą samą migawką. Wie to wyłącznie okno,
  // z przebiegu, który trzymało przed naciśnięciem.
  const poprzedni = kontekst.przebieg;
  const bylZakonczony =
    poprzedni !== null && poprzedni.id === wynik.wynik.id && poprzedni.finishedAt !== undefined;
  kontekst.przebieg = scalPrzebieg(poprzedni, wynik.wynik);
  const potwierdzenie = zdaniePrzerwania(kontekst.przebieg, bylZakonczony);
  tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
  rysuj();
}

/** Zdarzenie `developer.build.changed` filtrowane po oknie modułu — poza wytwórnią. */
function obsluzZdarzenieBudowania(
  kontekst: KontekstBudowania,
  idOkna: string,
  zdarzenie: DeveloperBuildChangedEvent,
  rysuj: () => void,
): void {
  if (zdarzenie.build.windowId !== idOkna) return;
  powiazZWlasnymStartem(kontekst, zdarzenie.build.id);
  kontekst.uciety = kontekst.idZnanyOdPoczatku === null || kontekst.idZnanyOdPoczatku !== zdarzenie.build.id;
  kontekst.przebieg = scalPrzebieg(kontekst.przebieg, zdarzenie.build);
  if (zdarzenie.logLine !== undefined && zdarzenie.logLine !== '') kontekst.logi.push(zdarzenie.logLine);
  rysuj();
}

/**
 * Łączy migawkę przebiegu ze świeżo przyszłą — odpowiedź `run` bywa starszą
 * migawką niż zdarzenie, które zdążyło już donieść stan końcowy (wyścig).
 * Ten sam przebieg (to samo `id`) nie cofa się z zakończonego do trwającego:
 * migawka z `finishedAt` wygrywa z migawką bez niego.
 */
function scalPrzebieg(biezacy: DeveloperBuild | null, przychodzacy: DeveloperBuild): DeveloperBuild {
  if (biezacy === null || biezacy.id !== przychodzacy.id) return przychodzacy;
  if (biezacy.finishedAt !== undefined && przychodzacy.finishedAt === undefined) return biezacy;
  return przychodzacy;
}

/** Eksport logu zebranego przez okno — czysta czynność, znacznik ucięcia w treści pliku. */
function eksportujLogBudowania(
  tresc: StanTresci,
  przebieg: DeveloperBuild | null,
  logi: readonly string[],
  uciety: boolean,
): void {
  if (logi.length === 0) {
    tresc.potwierdzenie('Nie ma czego wyeksportować — log jest jeszcze pusty.', false);
    return;
  }
  const naglowek = uciety
    ? '# Log budowania (ucięty — okno nie widziało początku przebiegu)\n\n'
    : '# Log budowania\n\n';
  pobierzPlik(`${przebieg?.id ?? 'budowanie'}-log.md`, naglowek + logi.join('\n'), 'text/markdown');
  tresc.potwierdzenie('Log zebrany w tym oknie pobrany jako plik.', true);
}

/**
 * Składa pasek akcji: uruchomienie, przerwanie (bez warunku), ponowienie
 * i eksport logu zebranego w oknie, oraz trzy pozycje bez komendy w kontrakcie.
 *
 * Wydzielone z wytwórni okna dla progu długości funkcji. Przyciski podpinają się
 * tutaj, bo wytwórnia nie ma po nich żadnej innej potrzeby niż podpięcie —
 * oddawanie ich na zewnątrz wyłącznie po to, by je zaraz podpiąć, byłoby drogą
 * bez odbiorcy.
 */
function zlozAkcjeBudowania(
  gospodarz: HTMLElement,
  obsluga: {
    uruchom: () => void;
    przerwij: () => void;
    ponow: () => void;
    eksportujLog: () => void;
  },
): void {
  const uruchomPrzycisk = przycisk('Uruchom budowanie', 'dn-btn dn-btn--atrament');
  // Wariant biblioteczny dla czynności przerywającej pracę — jak w Execution
  // Monitorze; wariant `--ostrzezenie` nie istnieje i jedno użycie go nie uzasadnia.
  const przerwijPrzycisk = przycisk('Przerwij budowanie', 'dn-btn dn-btn--niebezpieczny');
  // Bez warunku, tak jak przerwanie: brak czego powtarzać jest
  // odpowiedzią okna w stanie błędu, nie wyszarzeniem przycisku.
  const ponowPrzycisk = przycisk('Uruchom ponownie');
  const eksportujPrzycisk = przycisk('Eksportuj log zebrany');

  uruchomPrzycisk.addEventListener('click', obsluga.uruchom);
  przerwijPrzycisk.addEventListener('click', obsluga.przerwij);
  ponowPrzycisk.addEventListener('click', obsluga.ponow);
  eksportujPrzycisk.addEventListener('click', obsluga.eksportujLog);

  gospodarz.append(
    uruchomPrzycisk,
    przerwijPrzycisk,
    ponowPrzycisk,
    eksportujPrzycisk,
    przyciskBezKomendy(
      'Pobierz pełny log',
      powodBezKomendy(
        'DeveloperBuild.logRef odsyła do pełnego logu, a rozwinięcie tego odnośnika wymagałoby ' +
          'komendy developer.build.log.get — log tego okna to wyłącznie to, co przyszło ' +
          'zdarzeniami. Dziennik przebiegów rdzeń już prowadzi, więc brakuje samej komendy.',
      ),
    ),
    przyciskBezKomendy(
      'Wykaz przebiegów',
      powodBezKomendy(
        'Historia przebiegów wymagałaby komendy developer.build.list; rdzeń prowadzi dziennik ' +
          'przebiegów, ale kontrakt nie ma czym go odczytać.',
      ),
    ),
    przyciskBezKomendy(
      'Wynik testów i pokrycie',
      powodBezKomendy(
        'Rozróżnienie przeszedł/nie przeszedł/pominięty oraz pokrycie z rozbiciem na pliki ' +
          'wymagałyby komend developer.test.result.get i developer.coverage.get; ' +
          'DeveloperBuild niesie zgłoszenia, nie wynik zestawu testów.',
      ),
    ),
    przyciskBezKomendy(
      'Otwórz w Diagnostics Center',
      'Odmowa uruchomienia trafia do Diagnostics automatycznie w rdzeniu, ale przejście między oknami modułów wymaga nawigacji powłoki, do której to okno nie ma dostępu.',
    ),
  );
}

/** Kontrolki zakładki Build Output wraz z jej obszarem. */
interface PowierzchniaBudowania {
  /** Obszar zakładki osadzany w pasie zakładek okna. */
  obszar: HTMLElement;
  zadanie: HTMLInputElement;
  argumenty: HTMLInputElement;
  szukaj: HTMLInputElement;
  waga: HTMLSelectElement;
}

/** Składa kontrolki zakładki Build Output: pola zadania, zawężenia i miejsce treści. */
function zlozPowierzchnieBudowania(stanTresci: HTMLElement): PowierzchniaBudowania {
  const zadanie = pole('Zadanie budowania', 'np. test albo build');
  const argumenty = pole('Parametry (rozdzielone spacją)', 'opcjonalne');
  const szukaj = pole('Szukaj w logu', 'fraza zawężająca wiersze logu');
  const waga = wybor('Waga zgłoszeń', [
    ['', 'wszystkie wagi'],
    [ProblemSeverity.Error, 'wyłącznie błędy'],
    [ProblemSeverity.Warning, 'wyłącznie ostrzeżenia'],
    [ProblemSeverity.Info, 'wyłącznie informacje'],
  ]);

  const obszar = cialoZakladki(
    wiersz('Zadanie', zadanie, {
      klasa: 'mdev-wiersz',
      objasnienie:
        'Wymagane przy uruchomieniu. Pierwsze słowo jest programem uruchamianym na maszynie ' +
        'rdzenia, reszta jego wiodącymi parametrami.',
    }),
    wiersz('Parametry', argumenty, {
      klasa: 'mdev-wiersz',
      objasnienie: 'Puste pole nie wysyła argumentów budowania.',
    }),
    wiersz('Szukaj w logu', szukaj, {
      klasa: 'mdev-wiersz',
      objasnienie:
        'Zawężenie klienckie nad wierszami zebranymi w tym oknie — do rdzenia nic nie jedzie.',
    }),
    wiersz('Waga zgłoszeń', waga, {
      klasa: 'mdev-wiersz',
      objasnienie: 'Waga pochodzi z pola zgłoszenia oddanego przez rdzeń, nie z treści wiersza logu.',
    }),
    rysujZaleznosci(zaleznosci(['zadanie-budowania'])),
    stanTresci,
  );
  return { obszar, zadanie, argumenty, szukaj, waga };
}

/** Zawężenie widoku odczytane z pól okna. */
function zawezenieZPol(powierzchnia: PowierzchniaBudowania): ZawezeniePrzebiegu {
  return {
    szukaj: powierzchnia.szukaj.value,
    waga: czyWaga(powierzchnia.waga.value) ? powierzchnia.waga.value : '',
  };
}

/** Rozstrzyga, czy wartość pola wyboru jest wagą kontraktu. */
function czyWaga(wartosc: string): wartosc is ProblemSeverity {
  return (Object.values(ProblemSeverity) as readonly string[]).includes(wartosc);
}

/** Zadanie uruchomienia budowania; zadanie puste znaczy „nie wiem co uruchomić”. */
function zadanieUruchomienia(
  idOkna: string,
  zadanieTekst: string,
  argumentyTekst: string,
): DeveloperBuildRunRequest | null {
  const zadanie = zadanieTekst.trim();
  if (zadanie === '') return null;
  const args = argumentyTekst.trim();
  const wynik: DeveloperBuildRunRequest = { windowId: idOkna, task: zadanie };
  if (args !== '') wynik.arguments = args.split(/\s+/);
  return wynik;
}

/**
 * Zadanie przerwania — `task` jest wymagane kontraktem nawet dla przerwania,
 * więc bierze pole zadania Operatora, a bez niego zadanie przebiegu znanego
 * oknu. Przycisk działa nawet, gdy oba są puste; odmowa rdzenia
 * zostaje wtedy widoczna w stanie błędu.
 */
function zadanieZatrzymania(
  idOkna: string,
  zadanieTekst: string,
  przebieg: DeveloperBuild | null,
): DeveloperBuildRunRequest {
  const zadanie = zadanieTekst.trim() !== '' ? zadanieTekst.trim() : (przebieg?.task ?? '');
  return { windowId: idOkna, task: zadanie, stop: true };
}
