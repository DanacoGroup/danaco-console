import { Command, QueueAction, QueueStatus, type Queue } from '../../../../shared/contract';
import type { Wynik } from '../../protokol/kanal';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleLiczbowe,
  poleTresci,
  przyciskAkcji as przycisk,
  pozycjaWykazu,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import type { StanAutomatyki } from './stan-automatyki';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Queue Manager — okno zarządcy modułu Automations.
 *
 * Okno nie ma własnego silnika kolejek ani własnego cyklu życia zlecenia:
 * jeden silnik obsługuje pętlę sesyjną i MultitaskingAI, a okno posuwa kolejkę
 * komendami kontraktu i pokazuje to, co rdzeń o niej oddaje.
 *
 * Czego kontrakt nie oddaje, tego okno nie zmyśla. Byt `Queue` niesie stan,
 * licznik obiegów, okna, politykę, liczbę zleceń oczekujących i powiązaną
 * automatykę, lecz nie niesie samych zleceń — te oddaje osobna komenda wykazu
 * zleceń. Dopóki rdzeń nie ma jej uchwytu, przegląd pokazuje kolejkę jako
 * całość, a identyfikator zlecenia wpisuje Operator; okno mówi to wprost,
 * zamiast pokazywać pustą listę udającą przegląd.
 *
 * O tym, czy rdzeń ma uchwyt danej komendy, okno samo nie orzeka i nie
 * przepisuje stanu kontraktu do zdań: powód pozycji paska bierze się z wykazu
 * komend oddanego przez rdzeń (`pokrycie-komend.ts`), a wykaz działań silnika
 * z wyliczenia kontraktu. Oba przestają być prawdziwe same, bez edycji tego
 * pliku.
 *
 * Zdanie potwierdzenia powstaje z oddanej kolejki, nie z etykiety przycisku:
 * etykieta mówi o żądaniu, a rdzeń może oddać kolejkę w innym stanie, niż
 * żądanie zapowiadało.
 */
export interface OknoQueueManagera {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Nazwy działań silnika kolejek na ekranie.
 *
 * Zapis jest mapą zupełną po wyliczeniu, a nie wykazem przepisanym ręcznie:
 * gdy `QueueAction` urośnie, kompilacja zatrzyma się tutaj, zamiast zostawić
 * okno pokazujące mniej działań, niż silnik ma. Wykaz przepisany ręcznie
 * przestał raz odpowiadać kontraktowi i to jest zapora na powtórzenie.
 */
const NAZWY_DZIALAN: Readonly<Record<QueueAction, string>> = {
  [QueueAction.Start]: 'Uruchom',
  [QueueAction.Pause]: 'Wstrzymaj',
  [QueueAction.Resume]: 'Wznów',
  [QueueAction.Retry]: 'Powtórz (bieg naprawczy)',
  [QueueAction.Stop]: 'Zatrzymaj',
  [QueueAction.Clear]: 'Opróżnij',
  [QueueAction.Enqueue]: 'Wstaw zlecenie',
  [QueueAction.Dequeue]: 'Zdejmij zlecenie',
  [QueueAction.Delay]: 'Odłóż zlecenie',
  [QueueAction.Split]: 'Podziel zlecenie',
  [QueueAction.Merge]: 'Scal zlecenia',
  [QueueAction.Route]: 'Skieruj zlecenie',
  [QueueAction.Branch]: 'Rozgałęź zlecenie',
  [QueueAction.Condition]: 'Uwarunkuj zlecenie',
};

/**
 * Działania dotyczące pojedynczego zlecenia, nie kolejki jako całości.
 *
 * Rozdział wynika z kształtu żądania: `itemId` jest polem nieobowiązkowym, więc
 * kontrakt nie odmówi działania bez wskazania zlecenia, a rdzeń nie miałby
 * wtedy czego wstawić, zdjąć ani odłożyć. Okno pyta o zlecenie wcześniej,
 * zamiast wysyłać żądanie, które nie ma przedmiotu.
 */
const DZIALANIA_ZLECENIA: ReadonlySet<QueueAction> = new Set([
  QueueAction.Enqueue,
  QueueAction.Dequeue,
  QueueAction.Delay,
  QueueAction.Split,
  QueueAction.Merge,
  QueueAction.Route,
  QueueAction.Branch,
  QueueAction.Condition,
]);

/** Wszystkie działania wyliczenia w kolejności kontraktu. */
export const DZIALANIA: ReadonlyArray<[QueueAction, string]> = Object.values(QueueAction).map(
  (dzialanie) => [dzialanie, NAZWY_DZIALAN[dzialanie]],
);

/** Nazwy stanów kolejki na ekranie — mapa zupełna po wyliczeniu. */
const NAZWY_STANOW_KOLEJKI: Readonly<Record<QueueStatus, string>> = {
  [QueueStatus.Idle]: 'bezczynna',
  [QueueStatus.Running]: 'pracuje',
  [QueueStatus.Paused]: 'wstrzymana',
  [QueueStatus.Stopped]: 'zatrzymana',
  [QueueStatus.Done]: 'wyczerpana',
};

/** Stany kolejki w wykazie zawężenia; pusta wartość znaczy „wszystkie”. */
const STANY_KOLEJKI: ReadonlyArray<[string, string]> = [
  ['', 'wszystkie stany'],
  ...Object.values(QueueStatus).map((stan): [string, string] => [stan, NAZWY_STANOW_KOLEJKI[stan]]),
];

export function utworzOknoQueueManagera(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
  pokrycie: PokrycieKomend,
): OknoQueueManagera {
  const rama = utworzRameOkna({
    tytul: 'Queue Manager',
    rola: 'zarządca',
    przeznaczenie: 'Kolejka wykonująca automatykę na jednym silniku kolejek platformy.',
    przedrostek: 'da',
  });
  const tresc = utworzStanTresci();
  const kontrolki = zlozPowierzchnieKolejki(rama, tresc.element, wykonaj, pokrycie);
  const { idKolejki, sesja, zlecenia, idZlecenia, priorytet, kolejkaDocelowa } = kontrolki;

  /** Kolejka wskazana w polu, a gdy pole puste — kolejka bieżąca modułu. */
  function wskazanaKolejka(): string {
    return idKolejki.value.trim() === '' ? stan.kolejka() : idKolejki.value.trim();
  }

  function pokaz(kolejka: Queue): void {
    stan.ustawKolejke(kolejka.id);
    idKolejki.value = kolejka.id;
    tresc.tresc().append(wykazStanuKolejki(kolejka), zastrzezenieWykazuZlecen(pokrycie));
  }

  const odbierz = odbiorKolejki(tresc, pokaz);

  function wykonaj(dzialanie: QueueAction, czynnosc: string, zAutomatyka = false): void {
    const kolejka = wskazanaKolejka();
    if (kolejka === '') {
      tresc.pusto('Nie ma jeszcze kolejki. Załóż ją, zanim wykonasz na niej działanie.');
      return;
    }
    if (DZIALANIA_ZLECENIA.has(dzialanie) && idZlecenia.value.trim() === '') {
      tresc.potwierdzenie(
        `Działanie „${NAZWY_DZIALAN[dzialanie].toLowerCase()}” dotyczy pojedynczego zlecenia — ` +
          'wskaż je w polu „Zlecenie”. Żądanie bez zlecenia nie ma przedmiotu, więc okno go nie wysyła.',
        false,
      );
      return;
    }
    const zadanie = zadanieDzialaniaKolejki(kolejka, dzialanie, idZlecenia.value.trim(),
      zAutomatyka ? stan.automatyka() : '');
    tresc.ladowanie('Działanie na kolejce…');
    void zrodlo.dzialanieKolejki(zadanie).then((wynik) =>
      odbierz(wynik, 'Rdzeń odmówił działania na kolejce.', czynnosc));
  }

  function ulozZlecenie(zPriorytetem: boolean): void {
    if (idZlecenia.value.trim() === '') {
      tresc.potwierdzenie('Wskaż zlecenie — priorytet i skierowanie dotyczą pojedynczego zlecenia.', false);
      return;
    }
    const zadanie = zadanieUlozeniaZlecenia(wskazanaKolejka(), idZlecenia.value.trim(),
      zPriorytetem, priorytet.value, kolejkaDocelowa.value.trim());
    if (!zPriorytetem && zadanie.targetQueueId === '') {
      tresc.potwierdzenie('Wskaż kolejkę docelową skierowania.', false);
      return;
    }
    tresc.ladowanie('Układanie zlecenia…');
    void zrodlo.dzialanieKolejki(zadanie).then((wynik) =>
      odbierz(wynik, 'Rdzeń odmówił ułożenia zlecenia.',
        zPriorytetem
          ? `Rdzeń przyjął priorytet zlecenia ${idZlecenia.value.trim()}; priorytet nie jest osobnym ` +
            'działaniem silnika, więc kolejka szła przy tym działaniem „wstrzymaj”.'
          : `Rdzeń przyjął skierowanie zlecenia ${idZlecenia.value.trim()} do kolejki ${kolejkaDocelowa.value.trim()} ` +
            'działaniem „skieruj”.'));
  }

  function zalozKolejke(): void {
    const zadanie = zadanieZalozeniaKolejki(sesja.value, zlecenia.value, stan.automatyka());
    tresc.ladowanie('Zakładanie kolejki…');
    void zrodlo.zalozKolejke(zadanie).then((wynik) =>
      odbierz(wynik, 'Rdzeń nie założył kolejki.', 'Rdzeń założył kolejkę.'));
  }

  /**
   * Zasilenie kolejki krokami automatyki.
   *
   * Okno nie twierdzi, że zasilenie się odbyło. Rdzeń zasila kolejkę tylko
   * wtedy, gdy nie ma ona jeszcze ani jednej pozycji — powtórne uruchomienie
   * jest biegiem naprawczym po zleceniach zastanych, nie ich podwojeniem.
   * Odpowiedź niesie stan kolejki i liczbę zleceń oczekujących, lecz nie mówi,
   * które zlecenia w niej stoją — zdanie mówi więc to, co odpowiedź naprawdę
   * niesie, i nie orzeka, że kroki weszły.
   */
  function zasilKrokami(): void {
    if (stan.automatyka() === '') {
      tresc.potwierdzenie('Wskaż automatykę w Workflow Builderze — to jej kroki zasilają kolejkę.', false);
      return;
    }
    wykonaj(
      QueueAction.Start,
      'Rdzeń przyjął uruchomienie kolejki ze wskazaniem automatyki (zasilenie krokami wykonuje wyłącznie ' +
        'na kolejce bez zleceń; okno nie orzeka, które zlecenia weszły).',
      true,
    );
  }

  /**
   * Wykaz kolejek z rdzenia (`queue.list`).
   *
   * Zawężenie idzie w żądaniu, nie w rysunku: komenda przyjmuje sesję i stan
   * kolejki, więc rdzeń oddaje od razu to, o co Operator pyta, i nie ma czego
   * odsiewać po tej stronie. Zdanie mówi, ile kolejek rdzeń oddał i przy jakim
   * zawężeniu — pusty wykaz przy zawężeniu znaczy co innego niż pusty wykaz bez.
   */
  function pokazWykazKolejek(): void {
    const zadanie: Parameters<ZrodloAutomations['wykazKolejek']>[0] = {};
    if (sesja.value.trim() !== '') zadanie.sessionId = sesja.value.trim();
    if (kontrolki.stanKolejki.value !== '') {
      zadanie.status = kontrolki.stanKolejki.value as QueueStatus;
    }
    tresc.ladowanie('Odczyt wykazu kolejek…');
    void zrodlo.wykazKolejek(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał wykazu kolejek.', wynik.blad);
        return;
      }
      const kolejki = wynik.wynik;
      const zakres = opisZawezeniaKolejek(sesja.value.trim(), kontrolki.stanKolejki.value);
      if (kolejki.length === 0) {
        tresc.potwierdzenie(`Rdzeń oddał pusty wykaz kolejek${zakres}.`, true);
        return;
      }
      tresc.tresc().append(wykazKolejekRdzenia(kolejki));
      tresc.potwierdzenie(`Rdzeń oddał ${kolejki.length} kolejek${zakres}.`, true);
    });
  }

  /**
   * Powiązanie wskazanej kolejki z bieżącą automatyką (`queue.link`). Rdzeń
   * oddaje kolejkę po powiązaniu — zdanie i wykaz biorą się z niej, nie z żądania.
   */
  function powiazKolejke(): void {
    const kolejka = wskazanaKolejka();
    if (kolejka === '') {
      tresc.pusto('Nie ma kolejki do powiązania. Załóż ją albo wpisz jej identyfikator.');
      return;
    }
    if (stan.automatyka() === '') {
      tresc.potwierdzenie('Wskaż automatykę w Workflow Builderze — to z nią wiąże się kolejka.', false);
      return;
    }
    tresc.ladowanie('Powiązanie kolejki z automatyką…');
    void zrodlo
      .powiazKolejke({ queueId: kolejka, workflowId: stan.automatyka() })
      .then((wynik) =>
        odbierz(wynik, 'Rdzeń odmówił powiązania kolejki.',
          `Rdzeń powiązał kolejkę z automatyką ${stan.automatyka()}.`));
  }

  kontrolki.zaloz.addEventListener('click', zalozKolejke);
  kontrolki.zasil.addEventListener('click', zasilKrokami);
  kontrolki.priorytetPrzycisk.addEventListener('click', () => ulozZlecenie(true));
  kontrolki.skieruj.addEventListener('click', () => ulozZlecenie(false));
  kontrolki.wykazKolejek.addEventListener('click', pokazWykazKolejek);
  kontrolki.powiaz.addEventListener('click', powiazKolejke);

  stan.naZmiane(() => {
    if (stan.kolejka() !== '' && idKolejki.value === '') idKolejki.value = stan.kolejka();
  });

  function odswiez(): void {
    if (stan.kolejka() === '') {
      tresc.pusto('Nie ma jeszcze kolejki dla tej automatyki. Załóż ją albo zasil krokami automatyki.');
      return;
    }
    idKolejki.value = stan.kolejka();
    tresc.potwierdzenie('Kolejka bieżąca wskazana. Wykonaj działanie, aby odczytać jej stan.', true);
  }

  return { element: rama.element, odswiez };
}

/** Kontrolki okna Queue Managera. */
interface PowierzchniaKolejki {
  idKolejki: HTMLInputElement;
  sesja: HTMLInputElement;
  zlecenia: HTMLTextAreaElement;
  idZlecenia: HTMLInputElement;
  priorytet: HTMLInputElement;
  kolejkaDocelowa: HTMLInputElement;
  stanKolejki: HTMLSelectElement;
  zaloz: HTMLButtonElement;
  zasil: HTMLButtonElement;
  priorytetPrzycisk: HTMLButtonElement;
  skieruj: HTMLButtonElement;
  wykazKolejek: HTMLButtonElement;
  powiaz: HTMLButtonElement;
}

/** Zdanie o zawężeniu wykazu kolejek — dopowiedzenie do liczby oddanych pozycji. */
function opisZawezeniaKolejek(sesja: string, stanKolejki: string): string {
  const czesci: string[] = [];
  if (sesja !== '') czesci.push(`sesji ${sesja}`);
  if (stanKolejki !== '') czesci.push(`stanu ${stanKolejki}`);
  return czesci.length === 0 ? ' (bez zawężenia)' : ` dla ${czesci.join(' i ')}`;
}

/**
 * Składa kontrolki, pasek akcji i ciało okna.
 *
 * Czysta konstrukcja: fragment nie domyka się na stanie okna ani na źródle.
 * Przyciski działań biorą wywołanie zwrotne — zdanie potwierdzenia powstaje
 * przy przycisku, a praca zostaje w oknie. Po kolejności dokładania idą
 * sprawdziany widoku, więc nie zmienia się jej bez potrzeby.
 */
function zlozPowierzchnieKolejki(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  naDzialanie: (dzialanie: QueueAction, zdanie: string) => void,
  pokrycie: PokrycieKomend,
): PowierzchniaKolejki {
  const idKolejki = pole('Identyfikator kolejki', 'puste — załóż nową');
  const sesja = pole('Sesja kolejki', 'np. sesja-automatyki');
  const zlecenia = poleTresci('Zlecenia początkowe kolejki', 3);
  const idZlecenia = pole('Identyfikator zlecenia', 'wymagany przy priorytecie i skierowaniu');
  const priorytet = poleLiczbowe('Priorytet zlecenia', 'niższa liczba — wcześniej');
  const kolejkaDocelowa = pole('Kolejka docelowa', 'skierowanie zlecenia');

  const zaloz = przycisk('Załóż kolejkę', 'dn-btn dn-btn--atrament');
  const zasil = przycisk('Zasil krokami automatyki', 'dn-btn dn-btn--atrament');
  const priorytetPrzycisk = przycisk('Zastosuj priorytet');
  const skieruj = przycisk('Skieruj zlecenie');

  rama.akcje.append(zaloz, zasil, priorytetPrzycisk, skieruj);
  for (const [dzialanie, etykieta] of DZIALANIA) {
    const kontrolka = przycisk(etykieta);
    // Zdanie mówi, co okno wysłało; stan po działaniu dokłada `odbiorKolejki`
    // z odpowiedzi rdzenia. Etykieta przycisku nie orzeka o skutku.
    kontrolka.addEventListener('click', () =>
      naDzialanie(dzialanie, `Rdzeń przyjął działanie „${etykieta.toLowerCase()}”.`));
    rama.akcje.append(kontrolka);
  }
  // Wykaz kolejek i powiązanie kolejki wołają rdzeń wprost: `queue.list` oddaje
  // `queues`, a `queue.link` oddaje `queue` po powiązaniu.
  const wykazKolejek = przycisk('Wykaz kolejek');
  const powiaz = przycisk('Powiąż kolejkę z automatyką');

  rama.akcje.append(
    wykazKolejek,
    powiaz,
    // Osiem z jedenastu akcji silnika kolejek opracowania nie ma dziś
    // odpowiednika w wyliczeniu `QueueAction`; pozycje nazywają komendę,
    // która by je wykonała, zamiast opisywać brak zdaniem wpisanym w moduł.
    pokrycie.przycisk(
      'Dodaj zlecenie do kolejki',
      Command.QueueItemEnqueue,
      'dołożenie zlecenia do kolejki istniejącej, poza harmonogramem',
    ),
    pokrycie.przycisk(
      'Zdejmij zlecenie z kolejki',
      Command.QueueItemDequeue,
      'zdjęcie zlecenia z kolejki przed jego wykonaniem',
    ),
    pokrycie.przycisk(
      'Opóźnij zlecenie',
      Command.QueueItemDelay,
      'odłożenie wykonania zlecenia o wskazany czas',
    ),
    pokrycie.przycisk(
      'Podziel i scal zlecenia',
      Command.QueueItemSplit,
      'podział zlecenia na podzadania oraz scalenie kilku zleceń w jedno',
    ),
    pokrycie.przycisk(
      'Rozgałęź i uwarunkuj zlecenie',
      Command.QueueItemBranch,
      'rozejście zlecenia na tory równoległe oraz przetworzenie warunkowe',
    ),
    pokrycie.przycisk(
      'Wykaz zleceń kolejki',
      Command.QueueItemList,
      'przegląd zleceń oczekujących i przetwarzanych wraz z ich ładunkiem i próbami',
    ),
    pokrycie.przycisk(
      'Zasięg, współbieżność i przepustowość',
      Command.QueuePolicySet,
      'zasięg kolejki, limit zadań równoległych, tempo przetwarzania i polityka ponawiania',
    ),
    pokrycie.przycisk(
      'Kolejka zadań martwych',
      Command.QueueDeadList,
      'zlecenia trwale nieudane, do przeglądu ręcznego albo ponowienia zbiorczego',
    ),
  );

  const stanKolejki = wybor('Stan kolejki', STANY_KOLEJKI);
  rama.narzedzia.append(stanKolejki);

  rama.cialo.append(
    wiersz('Kolejka', idKolejki, {
      klasa: 'da-wiersz',
      objasnienie: 'Kolejka bieżąca modułu; wspólna dla Queue Managera i Execution Monitora.',
    }),
    wiersz('Sesja', sesja, {
      klasa: 'da-wiersz',
      objasnienie: 'Kolejka założona ze strony głównej może iść bez karty sesji.',
    }),
    wiersz('Zlecenia początkowe', zlecenia, {
      klasa: 'da-wiersz',
      objasnienie: 'Po jednym w wierszu; puste pole zakłada kolejkę bez zleceń.',
    }),
    wiersz('Zlecenie', idZlecenia, {
      klasa: 'da-wiersz',
      objasnienie:
        'Wymagany przy działaniach dotyczących pojedynczego zlecenia. Dopóki okno nie pokazuje ' +
        'wykazu zleceń, identyfikator wpisuje Operator.',
    }),
    wiersz('Priorytet', priorytet, { klasa: 'da-wiersz' }),
    wiersz('Kolejka docelowa', kolejkaDocelowa, { klasa: 'da-wiersz' }),
    stanTresci,
  );

  return {
    idKolejki, sesja, zlecenia, idZlecenia, priorytet, kolejkaDocelowa, stanKolejki,
    zaloz, zasil, priorytetPrzycisk, skieruj, wykazKolejek, powiaz,
  };
}

/**
 * Wykaz stanu kolejki — czysta konstrukcja z bytu `Queue`, bez dostępu do stanu
 * okna, więc wyszła z `pokaz` wprost.
 */
function wykazStanuKolejki(kolejka: Queue): HTMLElement {
  const lista = wykaz('Stan kolejki', 'da-wykaz');
  const wiersze: ReadonlyArray<[string, string]> = [
    ['Kolejka', `${kolejka.name ?? 'bez nazwy'} (${kolejka.id})`],
    ['Stan', kolejka.status],
    ['Sesja', kolejka.sessionId === '' ? 'bez powiązania z kartą sesji' : kolejka.sessionId],
    ['Okna obsługiwane', (kolejka.windowIds ?? []).join(', ') || 'brak'],
    ['Licznik obiegów naprawczych', String(kolejka.cycle ?? 0)],
    [
      'Zleceń oczekujących',
      kolejka.pendingCount === undefined ? 'rdzeń nie podał liczby' : String(kolejka.pendingCount),
    ],
    [
      'Automatyka kolejki',
      kolejka.workflowId === undefined || kolejka.workflowId === ''
        ? 'bez powiązania z automatyką'
        : kolejka.workflowId,
    ],
    ['Polityka kolejki', opisPolityki(kolejka.policy)],
    ['Ostatnia zmiana', new Date(kolejka.updatedAt).toLocaleString('pl-PL')],
  ];
  for (const [nazwa, wartosc] of wiersze) lista.append(pozycjaWykazu(nazwa, wartosc, 'da').element);
  return lista;
}

/**
 * Zdanie o polityce kolejki. Brak polityki nie jest brakiem danych: kolejka
 * bez zapisanej polityki pracuje na wartościach domyślnych silnika i tak to
 * nazywamy, zamiast pokazywać puste pole.
 */
function opisPolityki(polityka: Queue['policy']): string {
  if (polityka === undefined) return 'bez zapisanej polityki — wartości domyślne silnika';
  const czesci: string[] = [];
  if (polityka.scope !== undefined) czesci.push(`zasięg ${polityka.scope}`);
  if (polityka.maxConcurrent !== undefined) {
    czesci.push(
      polityka.maxConcurrent === 0
        ? 'bez limitu zadań równoległych'
        : `limit zadań równoległych ${polityka.maxConcurrent}`,
    );
  }
  if (polityka.ratePerMinute !== undefined && polityka.ratePerMinute > 0) {
    czesci.push(`tempo ${polityka.ratePerMinute} na minutę`);
  }
  if (polityka.maxAttempts !== undefined) czesci.push(`prób ${polityka.maxAttempts}`);
  if (polityka.backoff !== undefined) czesci.push(`wycofanie ${polityka.backoff}`);
  if (polityka.deadLetterEnabled === true) czesci.push('kolejka zadań martwych włączona');
  return czesci.length === 0 ? 'polityka zapisana, lecz bez ani jednej nastawy' : czesci.join(' · ');
}

/**
 * Wykaz kolejek oddany przez `queue.list` — po jednej pozycji na kolejkę wraz
 * z jej stanem i powiązaniami. Czysta konstrukcja z bytów `Queue`.
 */
function wykazKolejekRdzenia(kolejki: readonly Queue[]): HTMLElement {
  const lista = wykaz('Kolejki jednego silnika', 'da-wykaz');
  for (const kolejka of kolejki) {
    const okna = (kolejka.windowIds ?? []).join(', ');
    lista.append(
      pozycjaWykazu(
        `${kolejka.name ?? 'bez nazwy'} (${kolejka.id})`,
        `stan ${kolejka.status}` +
          (kolejka.sessionId === '' ? '' : `, sesja ${kolejka.sessionId}`) +
          (okna === '' ? '' : `, okna ${okna}`),
        'da',
      ).element,
    );
  }
  return lista;
}

/**
 * Zdanie o braku wykazu zleceń — treść stała, więc wyszła z `pokaz` jako czysta
 * konstrukcja. Zatajenie tego braku byłoby pustą listą udającą przegląd.
 */
function zastrzezenieWykazuZlecen(pokrycie: PokrycieKomend): HTMLElement {
  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'dn-pole-opis';
  zastrzezenie.textContent =
    `Okno pokazuje kolejkę jako całość. ${pokrycie.zdanie(
      Command.QueueItemList,
      'przegląd zleceń oczekujących i przetwarzanych',
    )}`;
  return zastrzezenie;
}

/**
 * Odbiór odpowiedzi kolejki wspólny dla czterech czynności okna: odmowa rdzenia
 * idzie w stan błędu wprost, a stan udany do rysunku. Fragment bierze wywołanie
 * zwrotne rysunku, więc nie musi znać ani pól okna, ani stanu modułu.
 *
 * Zdanie czynności dostaje dopowiedzenie z odpowiedzi: numer kolejki i jej stan
 * po działaniu. Bez tego potwierdzenie mówiłoby o żądaniu zamiast o skutku.
 */
function odbiorKolejki(
  tresc: StanTresci,
  pokaz: (kolejka: Queue) => void,
): (wynik: Wynik<Queue>, odmowa: string, zdanie: string) => void {
  return (wynik, odmowa, zdanie) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(odmowa, wynik.blad);
      return;
    }
    const kolejka = wynik.wynik;
    pokaz(kolejka);
    tresc.potwierdzenie(
      `${zdanie} Kolejka ${kolejka.id} po działaniu: ${kolejka.status}, ` +
        `licznik obiegów naprawczych ${kolejka.cycle ?? 0}.`,
      true,
    );
  };
}

/** Zadanie działania na kolejce — czysta konstrukcja z wartości pól. */
function zadanieDzialaniaKolejki(
  kolejka: string,
  dzialanie: QueueAction,
  zlecenie: string,
  automatyka: string,
): Parameters<ZrodloAutomations['dzialanieKolejki']>[0] {
  const zadanie: Parameters<ZrodloAutomations['dzialanieKolejki']>[0] = {
    queueId: kolejka,
    action: dzialanie,
  };
  if (zlecenie !== '') zadanie.itemId = zlecenie;
  if (automatyka !== '') zadanie.workflowId = automatyka;
  return zadanie;
}

/**
 * Zadanie ułożenia zlecenia — priorytet albo skierowanie.
 *
 * Skierowanie ma własne działanie silnika i tym działaniem idzie. Priorytet
 * własnego nie ma: jest polem żądania, nie działaniem, więc jedzie przy
 * wstrzymaniu — działaniu odwracalnym jednym naciśnięciem „Wznów”. Zdanie
 * potwierdzenia mówi o obu skutkach, nie tylko o tym, o który Operator prosił.
 */
function zadanieUlozeniaZlecenia(
  kolejka: string,
  zlecenie: string,
  zPriorytetem: boolean,
  priorytet: string,
  docelowa: string,
): Parameters<ZrodloAutomations['dzialanieKolejki']>[0] {
  const zadanie: Parameters<ZrodloAutomations['dzialanieKolejki']>[0] = {
    queueId: kolejka,
    action: zPriorytetem ? QueueAction.Pause : QueueAction.Route,
    itemId: zlecenie,
  };
  if (zPriorytetem) zadanie.priority = Number.parseInt(priorytet, 10) || 0;
  else zadanie.targetQueueId = docelowa;
  return zadanie;
}

/**
 * Zadanie założenia kolejki — czysta konstrukcja z wartości pól; zlecenia idą
 * po jednym w wierszu, a puste pole zakłada kolejkę bez zleceń.
 */
function zadanieZalozeniaKolejki(
  sesja: string,
  zlecenia: string,
  automatyka: string,
): Parameters<ZrodloAutomations['zalozKolejke']>[0] {
  const zadanie: Parameters<ZrodloAutomations['zalozKolejke']>[0] = {
    sessionId: sesja.trim() === '' ? 'automations' : sesja.trim(),
  };
  const wykazZlecen = zlecenia
    .split('\n')
    .map((zlecenie) => zlecenie.trim())
    .filter((zlecenie) => zlecenie !== '');
  if (wykazZlecen.length > 0) zadanie.items = wykazZlecen;
  if (automatyka !== '') zadanie.name = `Automatyka ${automatyka}`;
  return zadanie;
}
