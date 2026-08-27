import './pomocnicze.css';

import {
  ChunkKind,
  TerminalOutputChannel,
  type StreamChunkEvent,
  type TerminalOutputLine,
  type TerminalOutputStreamRequest,
  type TerminalOutputStreamResponse,
} from '../../../shared/contract';
import { utworzRameOkna } from '../komponenty/rama-okna';
import { utworzStanTresci, type StanTresci } from '../komponenty/stan-tresci';
import { przyciskAkcji, poleLiczbowe } from '../modele/kontrolki-formularza-braki';
import type { Kanal, Wynik } from '../protokol/kanal';
import { utworzZrodloPodgladuBash } from './zrodlo-podgladu-bash';

/**
 * Podgląd w tle pokazuje wyjście powłok biegnących w tle bez otwierania karty terminala i bez wykonywania poleceń, odróżniając stan błędu, ogon pusty i zapis niedoszły osobnymi zdaniami.
 */
export interface OknoPodgladuBash {
  /** Sekcja osadzana w pasie okien pomocniczych modułu. */
  element: HTMLElement;
  /** Ponawia zapis na wyjście i odczyt ogona. */
  odswiez(): void;
  /** Zamyka subskrypcję `stream.chunk` założoną przez to okno. */
  zamknij(): void;
}

/** Zależności okna; wartość okno pusta znaczy, że rdzeń nie dał temu modułowi żadnego własnego okna wykonania. */
export interface OpcjePodgladuBash {
  kanal: Kanal;
  /** Okno wykonania modułu — odbiorca zbiorczego wyjścia. */
  okno: string;
  /** Nazwa modułu w etykiecie dostępności — `Developer`, `Diagnostics`. */
  modul: string;
  /** Przedrostek klas modułu: `mdev` dla Developera, `dg` dla Diagnostics. */
  przedrostek: string;
}

/** Ile ostatnich wierszy trzyma widok okna podglądu, zanim starsze wiersze zostaną odrzucone z pamięci. */
const POJEMNOSC_WIDOKU = 500;

/** Ile wierszy historii wyjścia okno prosi przy każdym zapisie na zbiorcze wyjście kart terminala tego modułu. */
const OGON_DOMYSLNY = 200;

/** Jeden wiersz podglądu w tle — wspólna postać używana zarówno dla historii, jak i dla żywego strumienia. */
interface WierszPodgladu {
  tekst: string;
  /** Wiersz z wyjścia diagnostycznego (`stderr` / fragment rodzaju `error`). */
  bledny: boolean;
  /** Skąd wiersz przyszedł — karta terminala albo proces. */
  zrodlo: string;
}

/** Stan zmienny okna podglądu w tle trzymany jednym obiektem — funkcje obsługi zdarzeń stoją poza wytwórnią. */
interface KontekstPodgladu {
  wiersze: WierszPodgladu[];
  /** Czy widok wyrzucił wiersze starsze niż `POJEMNOSC_WIDOKU`. */
  uciety: boolean;
  /** Czy rdzeń potwierdził zapis okna na zbiorcze wyjście. */
  zapisane: boolean;
  /** Czy odczyt ogona doszedł do skutku choć raz. */
  odczytany: boolean;
}

export function utworzOknoPodgladuBash(opcje: OpcjePodgladuBash): OknoPodgladuBash {
  const zrodlo = utworzZrodloPodgladuBash(opcje.kanal);
  const rama = utworzRameOkna({
    tytul: 'Podgląd w tle (bash)',
    rola: 'pomocnicze',
    kod: 'podglad-bash',
    przeznaczenie:
      'Zbiorcze wyjście wszystkich otwartych kart powłoki — bez otwierania karty i bez wpisywania poleceń.',
    modul: opcje.modul,
    przedrostek: opcje.przedrostek,
  });
  const tresc = utworzStanTresci(opcje.przedrostek);
  const kontekst: KontekstPodgladu = { wiersze: [], uciety: false, zapisane: false, odczytany: false };

  const ogon = poleLiczbowe('Liczba wierszy ogona', String(OGON_DOMYSLNY));
  ogon.value = String(OGON_DOMYSLNY);

  function rysuj(): void {
    if (kontekst.wiersze.length === 0) {
      tresc.pusto(zdaniePustego(kontekst));
      return;
    }
    tresc.tresc().append(rysujWiersze(kontekst));
  }

  function odswiez(): void {
    tresc.ladowanie('Zapis na zbiorcze wyjście kart terminala…');
    void zrodlo
      .zapiszNaWyjscie(zadaniePodgladu(opcje.okno, ogon.value))
      .then((wynik) => przyjmijOdpowiedz(kontekst, tresc, wynik, rysuj));
  }

  function wyczysc(): void {
    kontekst.wiersze.length = 0;
    kontekst.uciety = false;
    rysuj();
    tresc.potwierdzenie('Widok wyczyszczony. Historia w rdzeniu została nietknięta.', true);
  }

  // Oba przyciski są czynne zawsze; odmowa rdzenia ląduje w stanie błędu okna, nie w wyszarzeniu.
  const odswiezPrzycisk = przyciskAkcji('Odśwież ogon', 'dn-btn dn-btn--atrament');
  const wyczyscPrzycisk = przyciskAkcji('Wyczyść widok');
  odswiezPrzycisk.addEventListener('click', odswiez);
  wyczyscPrzycisk.addEventListener('click', wyczysc);
  rama.akcje.append(odswiezPrzycisk, wyczyscPrzycisk);

  ogon.setAttribute('aria-label', 'Liczba wierszy ogona pobieranych z rdzenia');
  rama.narzedzia.append(ogon);
  rama.cialo.append(tresc.element);

  // Subskrypcja filtrowana po oknie modułu, bo rdzeń wysyła wiersze wszystkich obserwujących okien.
  const odsubskrybuj = zrodlo.naFragmentWyjscia((fragment) =>
    przyjmijFragment(kontekst, opcje.okno, fragment, rysuj),
  );

  // Pierwszego odczytu okno nie robi samo — robi go gospodarz, inaczej montaż wołałby zapis dwa razy.
  rysuj();

  return { element: rama.element, odswiez, zamknij: odsubskrybuj };
}

/**
 * Żądanie komendy zapisu; identyfikator okna idzie tylko wtedy, gdy moduł je zna, bo rdzeń odmawia zapisu na okno spoza rejestru.
 */
function zadaniePodgladu(okno: string, ogonTekst: string): TerminalOutputStreamRequest {
  const zadanie: TerminalOutputStreamRequest = { tail: ogonZTekstu(ogonTekst) };
  if (okno.trim() !== '') zadanie.windowId = okno;
  return zadanie;
}

/** Liczba wierszy ogona odczytana z pola formularza Operatora; wartość nieczytelna wraca do wartości domyślnej. */
function ogonZTekstu(tekst: string): number {
  const liczba = Number.parseInt(tekst, 10);
  if (!Number.isFinite(liczba) || liczba < 0) return OGON_DOMYSLNY;
  return liczba;
}

/** Odpowiedź zapisu na zbiorcze wyjście, obsługiwana poza wytwórnią panelu, bierze kontekst podglądu wprost. */
function przyjmijOdpowiedz(
  kontekst: KontekstPodgladu,
  tresc: StanTresci,
  wynik: Wynik<TerminalOutputStreamResponse>,
  rysuj: () => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    // Odmowa nie może wyglądać jak brak danych: okno pokazuje, że rdzeń odmówił, i z jakim kodem błędu.
    tresc.blad('Rdzeń odmówił zapisu na zbiorcze wyjście kart terminala.', wynik.blad);
    return;
  }
  kontekst.odczytany = true;
  kontekst.zapisane = wynik.wynik.subscribed;
  kontekst.wiersze = wynik.wynik.lines.map(wierszZHistorii);
  kontekst.uciety = false;
  przytnij(kontekst);
  rysuj();
  tresc.potwierdzenie(zdaniePotwierdzenia(wynik.wynik), wynik.wynik.subscribed);
}

/** Zdarzenie strumienia wyjścia przeznaczone dla okna tego modułu, obsługiwane osobno poza wytwórnią panelu. */
function przyjmijFragment(
  kontekst: KontekstPodgladu,
  okno: string,
  fragment: StreamChunkEvent,
  rysuj: () => void,
): void {
  if (okno.trim() === '' || fragment.windowId !== okno) return;
  const tekst = fragment.text;
  if (tekst === undefined || tekst === '') return;
  const bledny = fragment.kind === ChunkKind.Error;
  for (const wiersz of tekst.split('\n')) {
    if (wiersz === '') continue;
    kontekst.wiersze.push({ tekst: wiersz, bledny, zrodlo: opisZrodlaProcesu(fragment.messageId) });
  }
  przytnij(kontekst);
  rysuj();
}

/** Historia wyjścia oddana przez rdzeń, przepisana na postać jednego wiersza widoku podglądu pracującego w tle. */
function wierszZHistorii(wiersz: TerminalOutputLine): WierszPodgladu {
  return {
    tekst: wiersz.text,
    bledny: wiersz.channel === TerminalOutputChannel.Stderr,
    zrodlo: `karta ${wiersz.terminalSessionId}`,
  };
}

/**
 * Opis źródła wiersza żywego strumienia — fragment niesie identyfikator procesu, nie karty, bo kontrakt karty nie zna.
 */
function opisZrodlaProcesu(idProcesu: string): string {
  return idProcesu === '' ? 'proces nieznany' : `proces ${idProcesu}`;
}

/** Przycina widok wierszy do dopuszczalnej pojemności i zapamiętuje fakt, że przycinanie już nastąpiło. */
function przytnij(kontekst: KontekstPodgladu): void {
  const nadmiar = kontekst.wiersze.length - POJEMNOSC_WIDOKU;
  if (nadmiar <= 0) return;
  kontekst.wiersze.splice(0, nadmiar);
  kontekst.uciety = true;
}

/**
 * Zdanie stanu pustego nazywa właściwą pustkę: przed pierwszym odczytem, po odczycie bez wierszy albo bez zapisu na żywo.
 */
function zdaniePustego(kontekst: KontekstPodgladu): string {
  if (!kontekst.odczytany) return 'Odczyt zbiorczego wyjścia jeszcze nie wrócił z rdzenia.';
  const podstawa =
    'Żadna otwarta karta powłoki nic nie wypisała. Dziennik wyjścia rdzenia żyje jeden ' +
    'jego bieg — po restarcie rdzenia ogon jest pusty i to jest prawda o rdzeniu, nie strata.';
  return kontekst.zapisane
    ? `${podstawa} Okno jest zapisane na żywo: nowe wiersze pojawią się same.`
    : `${podstawa} ${ZDANIE_BEZ_ZAPISU}`;
}

/**
 * Zdanie o braku zapisu na żywo pada w jednym miejscu dla trzech stanów okna, bo moduł bez własnego okna nie może być zapisany na żywo mimo udanego odczytu.
 */
const ZDANIE_BEZ_ZAPISU =
  'Okno NIE jest zapisane na żywo: moduł nie ma okna nadanego przez rdzeń, więc żądanie poszło ' +
  'bez wskazania okna i podgląd pokazuje wyłącznie historię. Nowe wiersze przyjdą dopiero po ' +
  'naciśnięciu „Odśwież ogon".';

/** Potwierdzenie czynności zapisu na wyjście, wyraźnie rozdzielające zapis na żywo od samego odczytu historii. */
function zdaniePotwierdzenia(odpowiedz: TerminalOutputStreamResponse): string {
  const ile = `Ogon historii: ${odpowiedz.lines.length} wierszy.`;
  return odpowiedz.subscribed
    ? `${ile} Okno zapisane na zbiorcze wyjście kart — nowe wiersze dochodzą na żywo.`
    : `${ile} ${ZDANIE_BEZ_ZAPISU}`;
}

/**
 * Widok wierszy niesie znacznik ucięcia, ostrzeżenie o braku zapisu na żywo i samą treść, a klasy noszą przedrostek obszaru okien pomocniczych, nie przedrostek modułu.
 */
function rysujWiersze(kontekst: KontekstPodgladu): DocumentFragment {
  const widok = document.createDocumentFragment();
  if (kontekst.uciety) {
    widok.append(uwaga(`Widok pokazuje ostatnie ${POJEMNOSC_WIDOKU} wierszy — starsze z niego wypadły.`));
  }
  // Uwaga o braku zapisu pada dopiero po odczycie — przed nim zdanie o niezapisaniu byłoby bez podstawy.
  if (kontekst.odczytany && !kontekst.zapisane) widok.append(uwaga(ZDANIE_BEZ_ZAPISU));

  const pre = document.createElement('pre');
  pre.className = 'dnp-wyjscie';
  for (const wiersz of kontekst.wiersze) {
    const linia = document.createElement('span');
    linia.className = 'dnp-wyjscie__wiersz';
    linia.dataset['kanal'] = wiersz.bledny ? 'stderr' : 'stdout';
    linia.title = wiersz.zrodlo;
    linia.textContent = `${wiersz.tekst}\n`;
    pre.append(linia);
  }
  widok.append(pre);
  return widok;
}

/** Uwaga umieszczona nad treścią wierszy — nie zastępuje ich, tylko stoi obok, informując o stanie zapisu. */
function uwaga(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis dnp-uwaga';
  element.setAttribute('role', 'note');
  element.textContent = zdanie;
  return element;
}
