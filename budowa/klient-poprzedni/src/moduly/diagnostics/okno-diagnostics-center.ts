import type { DiagnosticAnalysis, DiagnosticsAnalyzeRunRequest } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  przyciskAkcji,
  przyciskBezKomendy,
  wiersz,
} from '../../modele/kontrolki-formularza-braki';
import { powodBezKomendy } from './braki-kontraktu';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import type { StanDiagnostyki, ZakresCzasu } from './stan-diagnostyki';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Diagnostics Center — okno wiodące modułu Diagnostics.
 *
 * Przegląd zagregowanego stanu systemu i uruchomienie analizy: punkt wejścia
 * agregujący Logs Viewer oraz Errors Panel. „Uruchom analizę” jest jedynym
 * przyciskiem sprawczym widoku.
 *
 * Po udanym biegu okno zapisuje identyfikator analizy do stanu
 * (`stan.ustawAnalize`) — Recommendations Panel pyta o rekomendacje tej analizy
 * (`analysisId`). Zakres czasu idzie do `stan.ustawZakres`, bo Logs Viewer
 * i Errors Panel jadą tym samym zakresem.
 *
 * Druga subskrypcja `diagnostics.analysis.changed` tu nie stoi:
 * `stan-diagnostyki.ts` sama nasłuchuje zdarzenia i filtruje po analizie
 * bieżącej — oknu wystarcza `stan.naZmiane(...)`.
 */
export interface OknoDiagnosticsCenter {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane` założony przy montażu okna. */
  zamknij(): void;
}

export function utworzOknoDiagnosticsCenter(
  zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
  idOkna: string,
): OknoDiagnosticsCenter {
  const rama = utworzRameOkna({
    tytul: 'Diagnostics Center',
    rola: 'wiodące',
    kod: 'diagnostics-center',
    przeznaczenie:
      'Przegląd zagregowanego stanu systemu i uruchomienie analizy — punkt wejścia agregujący dane z Logs Viewer i Errors Panel.',
    przedrostek: 'dg',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieCentrum(rama, tresc.element);

  function rysuj(): void {
    rysujAnalize(stan.migawka(), tresc);
  }

  function uruchom(): void {
    tresc.ladowanie('Analiza w toku…');
    const zadanie = zadanieAnalizy(idOkna, stan, powierzchnia);
    void zrodlo.uruchomAnalize(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(zdanieNiepowodzenia('uruchomienia analizy', wynik.powod), wynik.blad);
        return;
      }
      stan.ustawAnalize(wynik.wynik.id, wynik.wynik);
      const ocena = zdanieOAnalizie(zadanie, wynik.wynik);
      tresc.potwierdzenie(ocena.zdanie, ocena.udane);
    });
  }

  function eksportuj(): void {
    const migawka = stan.migawka();
    if (migawka === null) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — żadna analiza jeszcze nie biegła.', false);
      return;
    }
    pobierzPlik(`${migawka.id}-raport-diagnostyczny.json`, JSON.stringify(migawka, null, 2), 'application/json');
    tresc.potwierdzenie('Raport diagnostyczny pobrany jako plik JSON.', true);
  }

  /** Wpisuje zakres w pola okna i podaje go modułowi — jedna droga dla obu. */
  function ustawZakres(zakres: ZakresCzasu): void {
    powierzchnia.od.value = zakres.od === undefined ? '' : String(zakres.od);
    powierzchnia.do.value = zakres.do === undefined ? '' : String(zakres.do);
    stan.ustawZakres(zakres);
  }

  powierzchnia.uruchomPrzycisk.addEventListener('click', uruchom);
  powierzchnia.eksportPrzycisk.addEventListener('click', eksportuj);
  powierzchnia.zakresPrzycisk.addEventListener('click', () => {
    stan.ustawZakres(zakresZKontrolek(powierzchnia));
  });
  for (const skrot of powierzchnia.skroty) {
    skrot.przycisk.addEventListener('click', () => {
      // Jedna chwila odczytu na oba końce: dwa wywołania zegara dałyby zakres
      // o końcu późniejszym niż początek o czas własnego wyliczenia.
      const teraz = Date.now();
      ustawZakres(skrot.okno === undefined ? {} : { od: teraz - skrot.okno, do: teraz });
    });
  }

  const odsubskrybuj = stan.naZmiane(rysuj);
  rysuj();

  return { element: rama.element, odswiez: rysuj, zamknij: odsubskrybuj };
}

/** Kontrolki okna: zakres czasu, porównanie z migawką i pasek akcji. */
interface PowierzchniaCentrum {
  od: HTMLInputElement;
  do: HTMLInputElement;
  zakresPrzycisk: HTMLButtonElement;
  skroty: readonly SkrotZakresu[];
  compareId: HTMLInputElement;
  errorIds: HTMLInputElement;
  uruchomPrzycisk: HTMLButtonElement;
  eksportPrzycisk: HTMLButtonElement;
}

/** Skrót zakresu: przycisk wraz z długością okna; brak okna znaczy „bez granicy”. */
interface SkrotZakresu {
  przycisk: HTMLButtonElement;
  okno?: number;
}

/**
 * Skróty zakresu czasu wymienione w opracowaniu okna: ostatnia godzina, dzień,
 * tydzień oraz zakres własny. Pozycja bez okna znosi zawężenie — zakres pusty
 * jest stanem poprawnym modułu, nie brakiem.
 *
 * Skrót nie liczy niczego, czego nie widać: wyliczoną chwilę wpisuje w pola
 * początku i końca, więc Operator czyta z okna dokładnie te liczby, które idą
 * do rdzenia.
 */
const OKNA_SKROTOW: ReadonlyArray<{ etykieta: string; okno?: number }> = [
  { etykieta: 'Ostatnia godzina', okno: 60 * 60 * 1000 },
  { etykieta: 'Ostatni dzień', okno: 24 * 60 * 60 * 1000 },
  { etykieta: 'Ostatni tydzień', okno: 7 * 24 * 60 * 60 * 1000 },
  { etykieta: 'Bez zawężenia' },
];

/**
 * Składa narzędzia kontekstowe i pasek akcji ramy; ciało dostaje stan treści.
 *
 * Fragment jest czystą konstrukcją — nie domyka się na stanie modułu.
 */
function zlozPowierzchnieCentrum(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaCentrum {
  const od = poleLiczbowe('Początek zakresu (ms epoki)', 'domyślnie bez granicy');
  const doPole = poleLiczbowe('Koniec zakresu (ms epoki)', 'domyślnie bez granicy');
  const zakresPrzycisk = przyciskAkcji('Zastosuj zakres', 'dn-btn dn-btn--zarys');
  const zakres = document.createElement('div');
  zakres.className = 'dg-zakres';
  zakres.append(
    wiersz('Początek zakresu', od, { klasa: 'dg-wiersz' }),
    wiersz('Koniec zakresu', doPole, { klasa: 'dg-wiersz' }),
  );

  const skroty = OKNA_SKROTOW.map((pozycja) => ({
    przycisk: przyciskAkcji(pozycja.etykieta, 'dn-btn dn-btn--zarys dn-btn--sm'),
    ...(pozycja.okno === undefined ? {} : { okno: pozycja.okno }),
  }));
  const pasSkrotow = document.createElement('div');
  pasSkrotow.className = 'dg-skroty';
  pasSkrotow.setAttribute('aria-label', 'Skróty zakresu czasu');
  pasSkrotow.append(...skroty.map((skrot) => skrot.przycisk));

  const compareId = pole('Porównaj z migawką (id analizy)', 'opcjonalnie — id analizy wcześniejszej');
  const errorIds = pole('Błędy objęte analizą', 'opcjonalnie — id błędów po przecinku');

  rama.narzedzia.append(
    zakres,
    pasSkrotow,
    zakresPrzycisk,
    // Objaśnienia mówią, co robi okno, nie co zrobi rdzeń: rdzeń tych dwóch pól
    // dziś nie honoruje. O skutku orzeka wyłącznie zdanie potwierdzenia, bo ono
    // jedno czyta odpowiedź.
    wiersz('Porównanie z migawką wcześniejszą', compareId, {
      klasa: 'dg-wiersz',
      objasnienie: 'Pole idzie do rdzenia jako compareAnalysisId. Czy porównanie się odbyło, orzeka zdanie potwierdzenia — po polu comparedAnalysisId w odpowiedzi, nie po tym, że pole wypełniono.',
    }),
    wiersz('Błędy do analizy', errorIds, {
      klasa: 'dg-wiersz',
      objasnienie: 'Pole idzie do rdzenia jako errorIds. Które błędy analiza naprawdę objęła, orzeka zdanie potwierdzenia — po polu errorIds w odpowiedzi.',
    }),
  );

  const uruchomPrzycisk = przyciskAkcji('Uruchom analizę ✨');
  const eksportPrzycisk = przyciskAkcji('Eksportuj raport diagnostyczny', 'dn-btn dn-btn--zarys');
  rama.akcje.append(
    uruchomPrzycisk,
    eksportPrzycisk,
    przyciskBezKomendy(
      'Odnośniki szybkiego przejścia',
      powodBezKomendy('Przeskok do innego okna modułu wymagałby komendy nawigacji; to zadanie powłoki, nie modułu Diagnostics.'),
    ),
  );

  rama.cialo.append(stanTresci);
  return { od, do: doPole, zakresPrzycisk, skroty, compareId, errorIds, uruchomPrzycisk, eksportPrzycisk };
}

/** Zakres czasu odczytany z pól okna; pole niepoprawne znaczy „bez granicy”. */
function zakresZKontrolek(powierzchnia: PowierzchniaCentrum): ZakresCzasu {
  const zakres: ZakresCzasu = {};
  const od = Number.parseInt(powierzchnia.od.value, 10);
  if (Number.isInteger(od)) zakres.od = od;
  const konieca = Number.parseInt(powierzchnia.do.value, 10);
  if (Number.isInteger(konieca)) zakres.do = konieca;
  return zakres;
}

/**
 * Zadanie uruchomienia analizy z zakresu wspólnego modułu i pól okna.
 * `windowId` idzie tylko wtedy, gdy okno je zna — to jedyna komenda modułu,
 * która w ogóle je niesie, i jest tam opcjonalne.
 */
function zadanieAnalizy(
  idOkna: string,
  stan: StanDiagnostyki,
  powierzchnia: PowierzchniaCentrum,
): DiagnosticsAnalyzeRunRequest {
  const zakres = stan.zakres();
  const zadanie: DiagnosticsAnalyzeRunRequest = {};
  if (idOkna !== '') zadanie.windowId = idOkna;
  if (zakres.od !== undefined) zadanie.fromTime = zakres.od;
  if (zakres.do !== undefined) zadanie.toTime = zakres.do;
  const compareId = powierzchnia.compareId.value.trim();
  if (compareId !== '') zadanie.compareAnalysisId = compareId;
  const errorIds = powierzchnia.errorIds.value
    .split(',')
    .map((id) => id.trim())
    .filter((id) => id !== '');
  if (errorIds.length > 0) zadanie.errorIds = errorIds;
  return zadanie;
}

/**
 * Zdanie potwierdzenia biegu analizy — zestawienie żądania z odpowiedzią.
 *
 * Rdzeń nie honoruje dziś dwóch pól żądania: `UruchomAnalize`
 * (`core/adapter_modul_diagnostics_analiza.go`) dobiera błędy wyłącznie zakresem
 * czasu (`FiltrBledow{Od, Do}`), a `PorownanaKod` migawki nie ustawia nigdy,
 * więc `comparedAnalysisId` nie pojawi się w odpowiedzi nawet wtedy, gdy
 * porównywana analiza istnieje. Potwierdzenie zawsze udane robiłoby z tych
 * dwóch pól bez skutku pola pozornie działające.
 *
 * Zdanie mówi więc liczbami z odpowiedzi, a każde pominięte pole żądania
 * wychodzi na wierzch tonem nieudanym. Gdy rdzeń zacznie te pola honorować,
 * zastrzeżenia znikną same — nic o zachowaniu rdzenia nie jest tu wpisane
 * na sztywno.
 */
function zdanieOAnalizie(
  zadanie: DiagnosticsAnalyzeRunRequest,
  analiza: DiagnosticAnalysis,
): { zdanie: string; udane: boolean } {
  const objete = analiza.errorIds ?? [];
  const zastrzezenia: string[] = [];

  if (zadanie.compareAnalysisId !== undefined && analiza.comparedAnalysisId === undefined) {
    zastrzezenia.push(
      `porównania z migawką ${zadanie.compareAnalysisId} NIE BYŁO — rdzeń oddał analizę bez pola comparedAnalysisId`,
    );
  }
  const wskazane = zadanie.errorIds ?? [];
  const pominiete = wskazane.filter((id) => !objete.includes(id));
  if (pominiete.length > 0) {
    zastrzezenia.push(
      `z ${wskazane.length} wskazanych błędów analiza nie objęła ${pominiete.length} — ` +
        'rdzeń dobrał materiał zakresem czasu, nie wykazem z pola',
    );
  }

  const podstawa =
    `Analiza ${analiza.id} powstała: błędów objętych ${objete.length}, ` +
    `rekomendacji ${analiza.recommendationIds?.length ?? 0}.`;
  return zastrzezenia.length === 0
    ? { zdanie: podstawa, udane: true }
    : { zdanie: `${podstawa} Uwaga: ${zastrzezenia.join('; ')}.`, udane: false };
}

/** Rysuje migawkę analizy bieżącej albo stan pustki — brak biegu jest poprawny. */
function rysujAnalize(migawka: DiagnosticAnalysis | null, tresc: StanTresci): void {
  if (migawka === null) {
    tresc.pusto('Analiza nie została jeszcze uruchomiona. Ustaw zakres (opcjonalnie) i uruchom analizę.');
    return;
  }
  tresc.tresc().append(widokAnalizy(migawka));
}

/** Karta analizy — wszystkie pola nośne `DiagnosticAnalysis` widoczne wprost. */
function widokAnalizy(migawka: DiagnosticAnalysis): HTMLElement {
  const karta = document.createElement('div');
  karta.className = 'dg-analiza';

  const wynik = document.createElement('p');
  wynik.className = 'dg-analiza__wynik';
  wynik.textContent = migawka.summary === undefined || migawka.summary === ''
    ? 'Analiza nie oddała podsumowania.'
    : migawka.summary;

  const znacznik = document.createElement('p');
  znacznik.className = 'dg-analiza__znacznik';
  znacznik.textContent = opisMigawki(migawka);

  karta.append(wynik, znacznik);
  return karta;
}

/** Zdanie znacznika: identyfikator, czas, okno, zakres, liczności i migawka porównawcza. */
function opisMigawki(migawka: DiagnosticAnalysis): string {
  const czesci = [
    `id ${migawka.id}`,
    `utworzono ${new Date(migawka.createdAt).toLocaleString('pl-PL')}`,
  ];
  if (migawka.windowId !== undefined) czesci.push(`okno ${migawka.windowId}`);
  if (migawka.fromTime !== undefined || migawka.toTime !== undefined) {
    czesci.push(`zakres ${migawka.fromTime ?? '…'}–${migawka.toTime ?? '…'}`);
  }
  czesci.push(`błędy objęte: ${migawka.errorIds?.length ?? 0}`);
  czesci.push(`rekomendacje: ${migawka.recommendationIds?.length ?? 0}`);
  if (migawka.comparedAnalysisId !== undefined) {
    czesci.push(`porównano z ${migawka.comparedAnalysisId}`);
  }
  return czesci.join(' · ');
}
