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
 * Historia rozmowy to okno pomocnicze stojące obok rozmowy każdego gospodarza, dające Operatorowi drogę do odczytu, usuwania i przycinania trwale zapisanych wypowiedzi zasadą przechowywania.
 */
export interface OknoHistoriiRozmowy extends PanelPomocniczy {
  element: HTMLElement;
  odswiez(): void;
  zamknij(): void;
}

/** Ile pozycji panel prosi jedną stroną wykazu przy każdym kolejnym odczycie historii rozmowy z rdzenia. */
const STRONA_DOMYSLNA = 50;

/** Stan zmienny panelu historii trzymany jednym obiektem, bo obsługa zdarzeń stoi poza wytwórnią panelu. */
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
        // Rdzeń egzekwuje zasadę zaraz po zapisie, więc panel odczytuje wykaz zamiast pokazać stare pozycje.
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

  /** Sesja odczytana z pozycji jest jedynym miejscem, gdzie panel ją widzi; wypełnia nią zakres retencji. */
  function odczytanoSesje(): void {
    const zSesja = kontekst.pozycje.find((pozycja) => (pozycja.sessionId ?? '') !== '');
    nastawa.ustawSesje(zSesja?.sessionId ?? '');
  }

  function usunWskazane(): void {
    const wskazane = wykaz.zaznaczone();
    if (wskazane.length === 0) {
      // Żądanie bez pozycji czyściłoby całą historię okna — to inna czynność niż usunięcie wskazanych.
      tresc.potwierdzenie(
        'Nie wskazano ani jednej pozycji, więc nic nie usunięto. Do wyczyszczenia całej ' +
          'historii tego okna służy osobna czynność obok.',
        false,
      );
      return;
    }
    kasuj({ windowId: opcje.okno, entryIds: wskazane }, `wskazane pozycje (${wskazane.length})`);
  }

  /** Czyszczenie całej historii okna to jedyna czynność panelu z potwierdzeniem, bo jest nieodwracalna. */
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
      // Wynik null znaczy, że Operator zostawił historię albo rdzeń odmówił; modal już to pokazał.
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
      // Liczba usuniętych idzie do Operatora nawet gdy wynosi zero — milczenie byłoby zniknięciem po cichu.
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

  // Przerysowanie idzie za ręką Operatora, nie za zapisem — inaczej byłoby już sprawozdaniem ze straty.
  const odsubskrybujProgi = nastawa.naZmiane(rysuj);
  // Zaznaczenie zmienia tylko podpis czynności; przerysowanie wykazu odznaczałoby wskazanie Operatora.
  const odsubskrybujWybor = wykaz.naZmianeZaznaczenia(() => {
    const ile = wykaz.zaznaczone().length;
    usunPrzycisk.textContent = ile === 0 ? 'Usuń wskazane' : `Usuń wskazane (${ile})`;
  });

  // Historia zmienia się też cudzą ręką, więc bez subskrypcji panel pokazywałby nieprawdziwy już wykaz.
  const odsubskrybujZmiane = zrodlo.naZmianeHistorii((zdarzenie) => {
    if (opcje.okno.trim() === '' || zdarzenie.windowId !== opcje.okno) return;
    odswiez();
  });

  // Pierwszego odczytu panel nie robi sam — robi go gospodarz, inaczej odczyt poszedłby dwukrotnie.
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

/** Wielkość strony żądania z pola formularza Operatora; wartość nieczytelna wraca do domyślnej liczby pozycji strony. */
function stronaZTekstu(tekst: string): number {
  const liczba = Number.parseInt(tekst, 10);
  if (!Number.isFinite(liczba) || liczba <= 0) return STRONA_DOMYSLNA;
  return liczba;
}

/**
 * Odpowiedź żądania odczytu historii poza wytwórnią bierze kontekst wprost, a flaga dociągania rozstrzyga, czy pozycje zastępują wykaz, czy dochodzą na jego koniec.
 */
function przyjmij(
  kontekst: KontekstHistorii,
  tresc: StanTresci,
  wynik: Wynik<HistoryLoadResponse>,
  dolaczaj: boolean,
  rysuj: () => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    // Odmowa nie może wyglądać jak pusta historia — brak odpowiedzi to nie to samo co brak pozycji.
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

/** Zdanie pokazywane po odczycie historii — mówi, ile pozycji widać i ile jeszcze zostało w samym rdzeniu. */
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

/** Brak okna nadanego przez rdzeń — powód niedostępności wypisany wprost, bo sama pustka go nie niesie. */
const ZDANIE_BEZ_OKNA =
  'Gospodarz nie ma okna nadanego przez rdzeń, a historia jest własnością OKNA — nie ma więc ' +
  'czego odczytać ani czego usunąć. Zasadę przechowywania w zakresie całej instalacji nadal ' +
  'da się stąd ustawić: ona okna nie potrzebuje.';

/** Treść odmowy rdzenia podana jednym zdaniem wprost Operatorowi — sam kod błędu, bez dodatkowego tłumaczenia. */
function powod(wynik: Wynik<unknown>): string {
  return wynik.blad === undefined ? '' : `Powód: ${wynik.blad.message} (${wynik.blad.code}).`;
}
