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
 * Build Output i Run & Debug tworzą jedno okno monitora modułu Developer
 * o dwóch zakładkach: Build Output prowadzi budowanie, testy i log,
 * a Run & Debug — konfiguracje uruchomień i debugger krokowy.
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

  // Trzy części kolumny monitora: bieżący przebieg, sterowanie zatrzymanym programem i historia.
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

  // Zdanie o niewidocznej historii przychodzi z rdzenia, nie z napisu na stałe w kliencie.
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
    // Log czyszczony jest teraz, przed odpowiedzią, żeby zdarzenie trafiło do logu już pustego.
    kontekst.logi.length = 0;
    kontekst.oczekujeWlasnegoStartu = true;
    kontekst.ostatnieZadanie = zadanie;
    void zrodlo.budowanie(zadanie).then((wynik) => obsluzOdpowiedzUruchomienia(kontekst, tresc, wynik, rysuj));
  }

  // Ponowne uruchomienie komendą developer.build.run z zadaniem, które okno już zna.
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

  // Przycisk zatrzymania nie ma warunku; odmowa rdzenia ląduje w stanie błędu, nie w wyszarzeniu.
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

  // Zawężenie przelicza wyłącznie widok — przerysowanie jest całą jego obsługą.
  function przeliczZawezenie(): void {
    if (kontekst.przebieg !== null) rysuj();
  }
  powierzchnia.szukaj.addEventListener('input', przeliczZawezenie);
  powierzchnia.waga.addEventListener('change', przeliczZawezenie);

  rysuj();

  // Druga subskrypcja tego zdarzenia co stan-developer.ts — celowo, bierze pełną treść zdarzenia.
  const odsubskrybuj = zrodlo.naZmianeBudowania((zdarzenie) =>
    obsluzZdarzenieBudowania(kontekst, stan.okno(), zdarzenie, rysuj));

  return {
    element: rama.element,
    odswiez: rysuj,
    zamknij: () => {
      // Panel debugowania zamyka się razem z oknem, bo sesja debugowania jest stanem żywym.
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
  // Czynne od wysłania run do związania id, bo rdzeń nie gwarantuje kolejności odpowiedzi.
  oczekujeWlasnegoStartu: boolean;
  uciety: boolean;
  logi: string[];
  /** Zdanie stanu pustego złożone z rejestru komend rdzenia (`katalog-komend.ts`). */
  zdanieHistorii: string;
  // Ostatnie żądanie uruchomienia wysłane przez to okno — podstawa wiernego ponowienia zadania.
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

/** Powtórzenie wraz ze zdaniem o tym, czego rdzeń o powtarzanym przebiegu nie mówi wprost Operatorowi tego okna. */
interface Powtorzenie {
  zadanie: DeveloperBuildRunRequest;
  /** Pusty, gdy powtórzenie jest wierne; inaczej mówi, czego zabrakło. */
  dopisek: string;
}

/**
 * Składa żądanie ponowienia wiernie z żądania wysłanego przez to okno, a gdy
 * takiego nie było, z migawki przebiegu, którą okno zna, i wtedy nie udaje
 * wierności powtórzenia.
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

/** Odpowiedź developer.build.run w wariancie uruchomienia budowania — funkcja poza wytwórnią, bierze kontekst okna wprost. */
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
  // Zdanie budowane z przebiegu scalonego, nie z odpowiedzi — zdarzenie bywa szybsze niż odpowiedź.
  const potwierdzenie = zdanieUruchomienia(kontekst.przebieg);
  tresc.potwierdzenie(potwierdzenie.zdanie + dopisek, potwierdzenie.udane);
  rysuj();
}

/** Odpowiedź developer.build.run w wariancie przerwania budowania — funkcja poza wytwórnią, bierze kontekst okna wprost. */
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
  // Czy było co przerywać, wie tylko okno — rdzeń odpowiada tą samą migawką dla obu przypadków.
  const poprzedni = kontekst.przebieg;
  const bylZakonczony =
    poprzedni !== null && poprzedni.id === wynik.wynik.id && poprzedni.finishedAt !== undefined;
  kontekst.przebieg = scalPrzebieg(poprzedni, wynik.wynik);
  const potwierdzenie = zdaniePrzerwania(kontekst.przebieg, bylZakonczony);
  tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
  rysuj();
}

/** Zdarzenie developer.build.changed filtrowane po oknie modułu — funkcja stoi poza wytwórnią tego okna monitora. */
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
 * Łączy migawkę przebiegu ze świeżo przyszłą; ten sam przebieg nie cofa się
 * z zakończonego do trwającego, bo migawka z polem finishedAt wygrywa
 * z migawką bez niego.
 */
function scalPrzebieg(biezacy: DeveloperBuild | null, przychodzacy: DeveloperBuild): DeveloperBuild {
  if (biezacy === null || biezacy.id !== przychodzacy.id) return przychodzacy;
  if (biezacy.finishedAt !== undefined && przychodzacy.finishedAt === undefined) return biezacy;
  return przychodzacy;
}

/** Eksport logu zebranego przez to okno — czysta czynność, znacznik ucięcia zapisany wprost w treści pliku. */
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
 * Składa pasek akcji: uruchomienie, przerwanie bez warunku, ponowienie
 * i eksport logu zebranego w oknie, oraz trzy pozycje bez komendy w kontrakcie.
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
  // Wariant biblioteczny dla czynności przerywającej pracę, jak w Execution Monitorze tego produktu.
  const przerwijPrzycisk = przycisk('Przerwij budowanie', 'dn-btn dn-btn--niebezpieczny');
  // Bez warunku, jak przerwanie: brak czego powtarzać jest odpowiedzią okna w stanie błędu.
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

/** Kontrolki zakładki Build Output wraz z jej obszarem, osadzane w pasie zakładek tego okna monitora modułu Developer. */
interface PowierzchniaBudowania {
  /** Obszar zakładki osadzany w pasie zakładek okna. */
  obszar: HTMLElement;
  zadanie: HTMLInputElement;
  argumenty: HTMLInputElement;
  szukaj: HTMLInputElement;
  waga: HTMLSelectElement;
}

/** Składa kontrolki zakładki Build Output: pola zadania, zawężenia widoku logu oraz miejsce treści wynikowej. */
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

/** Zawężenie widoku przebiegu odczytane z pól zakładki Build Output bieżącego okna monitora modułu Developer. */
function zawezenieZPol(powierzchnia: PowierzchniaBudowania): ZawezeniePrzebiegu {
  return {
    szukaj: powierzchnia.szukaj.value,
    waga: czyWaga(powierzchnia.waga.value) ? powierzchnia.waga.value : '',
  };
}

/** Rozstrzyga, czy wartość wybrana w polu wagi zgłoszenia jest wagą znaną kontraktowi tego produktu programistycznego. */
function czyWaga(wartosc: string): wartosc is ProblemSeverity {
  return (Object.values(ProblemSeverity) as readonly string[]).includes(wartosc);
}

/** Zadanie uruchomienia budowania; zadanie o treści pustej znaczy „nie wiem, co uruchomić” zgłoszone Operatorowi. */
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
 * Zadanie przerwania bierze pole zadania Operatora, a bez niego zadanie
 * przebiegu znanego oknu; przycisk działa nawet, gdy oba są puste.
 */
function zadanieZatrzymania(
  idOkna: string,
  zadanieTekst: string,
  przebieg: DeveloperBuild | null,
): DeveloperBuildRunRequest {
  const zadanie = zadanieTekst.trim() !== '' ? zadanieTekst.trim() : (przebieg?.task ?? '');
  return { windowId: idOkna, task: zadanie, stop: true };
}
