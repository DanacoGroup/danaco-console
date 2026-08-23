import './pomocnicze.css';

import type { HistoryEntry, HistoryLoadRequest, HistoryLoadResponse } from '../../../shared/contract';
import { utworzRameOkna } from '../komponenty/rama-okna';
import { utworzStanTresci, type StanTresci } from '../komponenty/stan-tresci';
import { przyciskAkcji, poleLiczbowe } from '../modele/kontrolki-formularza-braki';
import type { Wynik } from '../protokol/kanal';
import { utworzNastawaRetencji } from './nastawa-retencji';
import type { OpcjePanelu, PanelPomocniczy } from './panel-pomocniczy';
import { otworzCzyszczenieHistorii } from './potwierdzenie-czyszczenia';
import { utworzWykazHistorii } from './wykaz-historii';
import { utworzZrodloHistorii } from './zrodlo-historii';

/**
 * Historia rozmowy — okno pomocnicze stojące obok rozmowy każdego gospodarza.
 *
 * Panel jest drogą Operatora do `history.load`, `history.delete`
 * i `retention.set`: rdzeń trzyma wypowiedzi okna trwale i umie je oddać od
 * najnowszej, skasować oraz przyciąć zasadą przechowywania.
 *
 * To nie jest drugi widok rozmowy. Kolumna rozmowy pokazuje turę bieżącą,
 * rosnącą na żywo strumieniem fragmentów; panel pokazuje zapis trwały tego
 * samego okna od pozycji najnowszej wstecz, stronicowany kursorem czasu —
 * czyli to, czego rozmowa nie pokazuje, bo nie sięga poza swoją sesję. Panel
 * nie nadaje wiadomości i nie zna komendy nadania.
 *
 * Zamiast pytania „czy na pewno" przed czynnością stoi odpowiedź po niej:
 * rdzeń oddaje liczbę usuniętych pozycji, a panel mówi ją wprost, żeby nic nie
 * znikało po cichu. Wyjątkiem jest wyczyszczenie całej historii okna — ten sam
 * wyjątek, który ma kasowanie sesji, o tym samym zasięgu, czyli czynność
 * nieodwracalna obejmująca całość. Usunięcie pozycji wskazanych potwierdzenia
 * nie ma: Operator wskazał je własną ręką.
 *
 * Wszystkie przyciski są czynne zawsze; brak wskazanych pozycji nie
 * wyszarza przycisku, tylko wraca zdaniem, co się stało.
 *
 * Gospodarz bywa bez okna nadanego przez rdzeń; `history.load` wymaga wtedy
 * czegoś, czego nie ma, i panel mówi to wprost zamiast pokazać pustą listę.
 * Nastawa retencji zostaje czynna, bo zakres `global` okna nie potrzebuje.
 */
export interface OknoHistoriiRozmowy extends PanelPomocniczy {
  element: HTMLElement;
  odswiez(): void;
  zamknij(): void;
}

/** Ile pozycji panel prosi jedną stroną wykazu. */
const STRONA_DOMYSLNA = 50;

/** Stan zmienny panelu — jednym obiektem, bo obsługa stoi poza wytwórnią. */
interface KontekstHistorii {
  pozycje: HistoryEntry[];
  /** Ile pozycji rdzeń naliczył przy ostatnim odczycie — może być więcej niż widać. */
  razem: number;
  /** Czy odczyt doszedł do skutku choć raz. */
  odczytany: boolean;
}

export function utworzOknoHistoriiRozmowy(opcje: OpcjePanelu): OknoHistoriiRozmowy {
  const zrodlo = utworzZrodloHistorii(opcje.kanal);
  const rama = utworzRameOkna({
    tytul: 'Historia rozmowy',
    rola: 'pomocnicze',
    kod: 'historia-rozmowy',
    przeznaczenie:
      'Trwały zapis wypowiedzi tego okna od najnowszej wstecz — z usuwaniem pozycji ' +
      'i z zasadą przechowywania, którą rdzeń egzekwuje.',
    modul: opcje.modul,
    przedrostek: opcje.przedrostek,
  });
  const tresc = utworzStanTresci(opcje.przedrostek);
  const wykaz = utworzWykazHistorii();
  const kontekst: KontekstHistorii = { pozycje: [], razem: 0, odczytany: false };

  const strona = poleLiczbowe('Ile pozycji na stronę wykazu', String(STRONA_DOMYSLNA));
  strona.value = String(STRONA_DOMYSLNA);

  const nastawa = utworzNastawaRetencji({
    okno: opcje.okno,
    naZapis: (zadanie) => {
      tresc.potwierdzenie('Zapis zasady przechowywania w toku…', true);
      void zrodlo.ustawZasade(zadanie).then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.potwierdzenie(`Rdzeń odmówił zapisu zasady. ${powod(wynik)}`, false);
          return;
        }
        nastawa.pokazZasade(wynik.wynik.policy);
        // Rdzeń egzekwuje zasadę zaraz po zapisie, więc wykaz sprzed zapisu
        // jest już nieprawdziwy — bez ponownego odczytu panel pokazywałby
        // pozycje, których w bazie nie ma.
        odswiez();
      });
    },
  });

  function rysuj(): void {
    if (kontekst.pozycje.length === 0) {
      tresc.pusto(zdaniePustego(kontekst, opcje.okno));
      return;
    }
    wykaz.rysuj(kontekst.pozycje, nastawa.pozaZasada(kontekst.pozycje));
    tresc.tresc().append(wykaz.element);
  }

  /** Odczyt od początku — pierwsza strona wykazu, kursor zerowany. */
  function odswiez(): void {
    if (opcje.okno.trim() === '') {
      tresc.pusto(ZDANIE_BEZ_OKNA);
      return;
    }
    tresc.ladowanie('Odczyt historii rozmowy tego okna…');
    void zrodlo
      .wczytaj(zadanieOdczytu(opcje.okno, strona.value, null))
      .then((wynik) => przyjmij(kontekst, tresc, wynik, false, () => { odczytanoSesje(); rysuj(); }));
  }

  /** Dociągnięcie strony starszej — kursorem czasu, wprost z pozycji najstarszej. */
  function starsze(): void {
    const najstarsza = kontekst.pozycje[kontekst.pozycje.length - 1];
    if (najstarsza === undefined) {
      odswiez();
      return;
    }
    tresc.ladowanie('Odczyt pozycji starszych…');
    void zrodlo
      .wczytaj(zadanieOdczytu(opcje.okno, strona.value, najstarsza.createdAt))
      .then((wynik) => przyjmij(kontekst, tresc, wynik, true, () => { odczytanoSesje(); rysuj(); }));
  }

  /**
   * Sesja odczytana z pozycji — jedyne miejsce, w którym panel ją widzi.
   * Bez niej nastawa retencji nie ma czym wypełnić zakresu `session`.
   */
  function odczytanoSesje(): void {
    const zSesja = kontekst.pozycje.find((pozycja) => (pozycja.sessionId ?? '') !== '');
    nastawa.ustawSesje(zSesja?.sessionId ?? '');
  }

  function usunWskazane(): void {
    const wskazane = wykaz.zaznaczone();
    if (wskazane.length === 0) {
      // Czynność się wykonała i odpowiedziała. Gdyby żądanie poszło bez
      // pozycji, rdzeń wyczyściłby całą historię okna — a to jest druga
      // czynność, nie ta sama (`handlers_historia.go`).
      tresc.potwierdzenie(
        'Nie wskazano ani jednej pozycji, więc nic nie usunięto. Do wyczyszczenia całej ' +
          'historii tego okna służy osobna czynność obok.',
        false,
      );
      return;
    }
    kasuj({ windowId: opcje.okno, entryIds: wskazane }, `wskazane pozycje (${wskazane.length})`);
  }

  /**
   * Czyszczenie całej historii okna — jedyna czynność panelu z potwierdzeniem.
   *
   * Wyjątek jest wąski i taki sam jak przy kasowaniu sesji: obejmuje czynność
   * nieodwracalną sięgającą całości, a nie każde kasowanie. „Usuń wskazane"
   * zostaje bez potwierdzenia, bo wskazanie pozycji jest zgodą.
   *
   * Modal woła rdzeń sam (dostaje wysyłkę), bo to on ma pokazać odpowiedź —
   * odmowa i liczba usuniętych padają tam, gdzie Operator patrzy w chwili
   * czynności. Panel powtarza je u siebie i przeładowuje wykaz.
   */
  function wyczysc(): void {
    if (opcje.okno.trim() === '') {
      tresc.potwierdzenie(ZDANIE_BEZ_OKNA, false);
      return;
    }
    void otworzCzyszczenieHistorii(
      {
        okno: opcje.okno,
        wWidoku: kontekst.pozycje.length,
        razem: kontekst.razem,
        odczytany: kontekst.odczytany,
      },
      () => zrodlo.usun({ windowId: opcje.okno }),
    ).then((odpowiedz) => {
      // `null` znaczy „Operator zostawił historię" albo „rdzeń odmówił" —
      // w obu razach modal już powiedział, co się stało, a wykaz jest
      // nietknięty i nie ma czego odświeżać.
      if (odpowiedz === null) return;
      tresc.potwierdzenie(
        `Usunięto ${odpowiedz.deleted} pozycji historii (całą historię tego okna).`,
        true,
      );
      odswiez();
    });
  }

  /** Wspólna droga obu kasowań — jedna komenda, jedno zdanie po niej. */
  function kasuj(zadanie: { windowId: string; entryIds?: string[] }, co: string): void {
    if (opcje.okno.trim() === '') {
      tresc.potwierdzenie(ZDANIE_BEZ_OKNA, false);
      return;
    }
    tresc.potwierdzenie(`Usuwanie: ${co}…`, true);
    void zrodlo.usun(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.potwierdzenie(`Rdzeń odmówił usunięcia. ${powod(wynik)}`, false);
        return;
      }
      // Liczba idzie do Operatora także wtedy, gdy wynosi zero: „usunięto 0"
      // to prawda o czynności, a milczenie na tym miejscu byłoby zniknięciem
      // po cichu.
      tresc.potwierdzenie(`Usunięto ${wynik.wynik.deleted} pozycji historii (${co}).`, true);
      odswiez();
    });
  }

  const odswiezPrzycisk = przyciskAkcji('Odśwież', 'dn-btn dn-btn--atrament');
  const starszePrzycisk = przyciskAkcji('Wczytaj starsze');
  const usunPrzycisk = przyciskAkcji('Usuń wskazane');
  const wyczyscPrzycisk = przyciskAkcji('Wyczyść historię okna', 'dn-btn dn-btn--niebezpieczny');
  odswiezPrzycisk.addEventListener('click', odswiez);
  starszePrzycisk.addEventListener('click', starsze);
  usunPrzycisk.addEventListener('click', usunWskazane);
  wyczyscPrzycisk.addEventListener('click', wyczysc);
  rama.akcje.append(odswiezPrzycisk, starszePrzycisk, usunPrzycisk, wyczyscPrzycisk);

  strona.setAttribute('aria-label', 'Ile pozycji historii pobierać jedną stroną');
  rama.narzedzia.append(strona);
  rama.cialo.append(tresc.element, nastawa.element);

  // Przerysowanie po każdej zmianie progów: zapowiedź „co zasada zrobi" ma iść
  // za ręką Operatora, a nie za zapisem. Po zapisie byłaby już nie zapowiedzią,
  // tylko sprawozdaniem ze straty.
  const odsubskrybujProgi = nastawa.naZmiane(rysuj);
  // Zaznaczenie zmienia wyłącznie podpis czynności, a nie wykaz — przerysowanie
  // wykazu przy każdym kliknięciu odznaczałoby to, co Operator właśnie wskazał.
  const odsubskrybujWybor = wykaz.naZmianeZaznaczenia(() => {
    const ile = wykaz.zaznaczone().length;
    usunPrzycisk.textContent = ile === 0 ? 'Usuń wskazane' : `Usuń wskazane (${ile})`;
  });

  // Historia zmienia się także cudzą ręką: inne okno kasuje pozycje, zasada
  // zakresu globalnego tnie wiele okien naraz, przemiatanie przy starcie
  // rdzenia zabiera stare wypowiedzi. Bez tej subskrypcji panel pokazywałby
  // wykaz, który przestał być prawdziwy — i nic by o tym nie powiedział.
  const odsubskrybujZmiane = zrodlo.naZmianeHistorii((zdarzenie) => {
    if (opcje.okno.trim() === '' || zdarzenie.windowId !== opcje.okno) return;
    odswiez();
  });

  // Pierwszego odczytu panel nie robi sam — robi go gospodarz przez `odswiez()`,
  // tak jak przy podglądzie w tle. Inaczej otwarcie panelu wołałoby
  // `history.load` dwa razy pod rząd.
  rysuj();

  return {
    element: rama.element,
    odswiez,
    zamknij: () => {
      odsubskrybujProgi();
      odsubskrybujWybor();
      odsubskrybujZmiane();
    },
  };
}

/**
 * Żądanie odczytu. `before` idzie wyłącznie przy dociąganiu strony starszej:
 * kontrakt nazywa je kursorem czasu, a podanie go przy odczycie od początku
 * odcięłoby pozycje najnowsze — czyli te, po które Operator sięga najczęściej.
 */
function zadanieOdczytu(okno: string, stronaTekst: string, kursor: number | null): HistoryLoadRequest {
  const zadanie: HistoryLoadRequest = { windowId: okno, limit: stronaZTekstu(stronaTekst) };
  if (kursor !== null) zadanie.before = kursor;
  return zadanie;
}

/** Wielkość strony z pola Operatora; wartość nieczytelna wraca do domyślnej. */
function stronaZTekstu(tekst: string): number {
  const liczba = Number.parseInt(tekst, 10);
  if (!Number.isFinite(liczba) || liczba <= 0) return STRONA_DOMYSLNA;
  return liczba;
}

/**
 * Odpowiedź `history.load` — poza wytwórnią, bierze kontekst wprost.
 *
 * @param dolaczaj `true` przy dociąganiu strony starszej: pozycje dochodzą na
 *   koniec wykazu zamiast go zastąpić. Zastąpienie przy dociąganiu skróciłoby
 *   widok do jednej strony i wyglądało jak utrata pozycji.
 */
function przyjmij(
  kontekst: KontekstHistorii,
  tresc: StanTresci,
  wynik: Wynik<HistoryLoadResponse>,
  dolaczaj: boolean,
  rysuj: () => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    // Odmowa nie może wyglądać jak pusta historia: „nie udało się zapytać" to
    // nie to samo co „nic nie zapisano".
    tresc.blad('Rdzeń odmówił odczytu historii rozmowy tego okna.', wynik.blad);
    return;
  }
  kontekst.odczytany = true;
  kontekst.razem = wynik.wynik.total;
  kontekst.pozycje = dolaczaj
    ? [...kontekst.pozycje, ...wynik.wynik.entries]
    : [...wynik.wynik.entries];
  rysuj();
  tresc.potwierdzenie(zdanieOdczytu(kontekst, wynik.wynik.entries.length, dolaczaj), true);
}

/** Zdanie po odczycie — mówi, ile widać i ile jeszcze zostało w rdzeniu. */
function zdanieOdczytu(kontekst: KontekstHistorii, doszlo: number, dolaczaj: boolean): string {
  const wstep = dolaczaj
    ? `Doszło ${doszlo} pozycji starszych.`
    : `Wykaz odczytany: ${doszlo} pozycji.`;
  if (doszlo === 0 && dolaczaj) {
    return `${wstep} Starszych pozycji rdzeń już nie ma — to jest początek historii tego okna.`;
  }
  return `${wstep} W widoku stoi ${kontekst.pozycje.length}, rdzeń naliczył ${kontekst.razem}.`;
}

/**
 * Zdanie stanu pustego. Pustka bywa poprawna, ale Operator ma wiedzieć, która
 * pustka go spotkała: brak okna, brak odczytu czy pusty zapis.
 */
function zdaniePustego(kontekst: KontekstHistorii, okno: string): string {
  if (okno.trim() === '') return ZDANIE_BEZ_OKNA;
  if (!kontekst.odczytany) return 'Odczyt historii jeszcze nie wrócił z rdzenia.';
  return (
    'To okno nie ma ani jednej zapisanej wypowiedzi. Pustka bywa też skutkiem ZASADY ' +
    'PRZECHOWYWANIA: rdzeń stosuje ją przy każdym odczycie, więc pozycje spoza progu ' +
    'nie wchodzą do wykazu wcale.'
  );
}

/** Brak okna nadanego przez rdzeń — powód wypisany, bo pustka sama go nie niesie. */
const ZDANIE_BEZ_OKNA =
  'Gospodarz nie ma okna nadanego przez rdzeń, a historia jest własnością OKNA — nie ma więc ' +
  'czego odczytać ani czego usunąć. Zasadę przechowywania w zakresie całej instalacji nadal ' +
  'da się stąd ustawić: ona okna nie potrzebuje.';

/** Treść odmowy rdzenia w jednym zdaniu — kod błędu wprost, bez tłumaczenia. */
function powod(wynik: Wynik<unknown>): string {
  return wynik.blad === undefined ? '' : `Powód: ${wynik.blad.message} (${wynik.blad.code}).`;
}
