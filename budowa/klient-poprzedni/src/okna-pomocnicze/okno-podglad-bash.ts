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
 * Podgląd w tle (bash) — okno pomocnicze modułów Developer i Diagnostics.
 * Pokazuje wyjście powłok biegnących w tle, bez otwierania karty terminala
 * i bez wpisywania poleceń; poleceń nie wykonuje, od tego jest moduł Terminal.
 *
 * Jedno wywołanie `terminal.output.stream` robi dwie rzeczy: zapisuje to okno
 * na zbiorcze wyjście wszystkich otwartych kart terminala i oddaje ogon
 * historii. Nowe wiersze dochodzą potem zdarzeniem `stream.chunk` z `windowId`
 * tego okna.
 *
 * Trzy stany dostają trzy różne zdania:
 *
 *  1. Odmowa rdzenia — okno pokazuje stan błędu wraz z treścią odmowy, nigdy
 *     pustą listę: „nie udało się zapytać" znaczy co innego niż „nic nie ma".
 *  2. Ogon pusty — odczyt się udał, a wierszy nie ma. Dziennik wyjścia rdzenia
 *     jest pierścieniem w pamięci żyjącym jeden bieg rdzenia, więc po jego
 *     restarcie pusty ogon jest prawdą o rdzeniu, a nie ukrytą stratą.
 *  3. Zapis niedoszły (`subscribed: false`) — historia wróciła, ale okno nie
 *     jest zapisane na żywo, bo żądanie poszło bez `windowId`.
 *
 * Widok trzyma ostatnie `POJEMNOSC_WIDOKU` wierszy, bo okno pomocnicze nie jest
 * drugą konsolą; ucięcie jest oznaczone, nie zamaskowane.
 */
export interface OknoPodgladuBash {
  /** Sekcja osadzana w pasie okien pomocniczych modułu. */
  element: HTMLElement;
  /** Ponawia zapis na wyjście i odczyt ogona. */
  odswiez(): void;
  /** Zamyka subskrypcję `stream.chunk` założoną przez to okno. */
  zamknij(): void;
}

/** Zależności okna; `okno` puste znaczy „rdzeń nie dał temu modułowi okna". */
export interface OpcjePodgladuBash {
  kanal: Kanal;
  /** Okno wykonania modułu — odbiorca zbiorczego wyjścia. */
  okno: string;
  /** Nazwa modułu w etykiecie dostępności — `Developer`, `Diagnostics`. */
  modul: string;
  /** Przedrostek klas modułu: `mdev` dla Developera, `dg` dla Diagnostics. */
  przedrostek: string;
}

/** Ile ostatnich wierszy trzyma widok okna. */
const POJEMNOSC_WIDOKU = 500;

/** Ile wierszy historii okno prosi przy zapisie na wyjście. */
const OGON_DOMYSLNY = 200;

/** Jeden wiersz podglądu — wspólna postać dla historii i dla żywego strumienia. */
interface WierszPodgladu {
  tekst: string;
  /** Wiersz z wyjścia diagnostycznego (`stderr` / fragment rodzaju `error`). */
  bledny: boolean;
  /** Skąd wiersz przyszedł — karta terminala albo proces. */
  zrodlo: string;
}

/** Stan zmienny okna, trzymany jednym obiektem — funkcje obsługi stoją poza wytwórnią. */
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

  // Oba przyciski są czynne zawsze; odmowa rdzenia ląduje w stanie błędu okna,
  // a nie w wyszarzeniu przycisku.
  const odswiezPrzycisk = przyciskAkcji('Odśwież ogon', 'dn-btn dn-btn--atrament');
  const wyczyscPrzycisk = przyciskAkcji('Wyczyść widok');
  odswiezPrzycisk.addEventListener('click', odswiez);
  wyczyscPrzycisk.addEventListener('click', wyczysc);
  rama.akcje.append(odswiezPrzycisk, wyczyscPrzycisk);

  ogon.setAttribute('aria-label', 'Liczba wierszy ogona pobieranych z rdzenia');
  rama.narzedzia.append(ogon);
  rama.cialo.append(tresc.element);

  // Subskrypcja `stream.chunk` filtrowana po oknie TEGO modułu. Rdzeń wysyła
  // obserwatorowi wiersze z `windowId` okna obserwującego, więc filtr jest
  // jedynym, co oddziela wyjście naszych kart od cudzych okien.
  const odsubskrybuj = zrodlo.naFragmentWyjscia((fragment) =>
    przyjmijFragment(kontekst, opcje.okno, fragment, rysuj),
  );

  // Pierwszego odczytu okno nie robi samo, tak jak okna operacyjne Developera
  // i Diagnostics czekające na `odswiez()` złożenia. Inaczej montaż wołałby
  // `terminal.output.stream` dwa razy pod rząd.
  rysuj();

  return { element: rama.element, odswiez, zamknij: odsubskrybuj };
}

/**
 * Żądanie komendy. `windowId` idzie tylko wtedy, gdy moduł okno zna: rdzeń
 * odmawia zapisu na okno, którego rejestr nie ma, a żądanie bez okna jest
 * kontraktem dopuszczone i znaczy sam odczyt ogona.
 */
function zadaniePodgladu(okno: string, ogonTekst: string): TerminalOutputStreamRequest {
  const zadanie: TerminalOutputStreamRequest = { tail: ogonZTekstu(ogonTekst) };
  if (okno.trim() !== '') zadanie.windowId = okno;
  return zadanie;
}

/** Liczba wierszy ogona z pola formularza; wartość nieczytelna wraca do domyślnej. */
function ogonZTekstu(tekst: string): number {
  const liczba = Number.parseInt(tekst, 10);
  if (!Number.isFinite(liczba) || liczba < 0) return OGON_DOMYSLNY;
  return liczba;
}

/** Odpowiedź `terminal.output.stream` — poza wytwórnią, bierze kontekst wprost. */
function przyjmijOdpowiedz(
  kontekst: KontekstPodgladu,
  tresc: StanTresci,
  wynik: Wynik<TerminalOutputStreamResponse>,
  rysuj: () => void,
): void {
  if (!wynik.udany || wynik.wynik === undefined) {
    // Odmowa nie może wyglądać jak brak danych: okno pokazuje, że rdzeń
    // odmówił, i z jakim kodem — na przykład `internal_error`, gdy rdzeń nie ma
    // nadajnika wyjścia i zbiorczego strumienia nie prowadzi w ogóle.
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

/** Zdarzenie `stream.chunk` przeznaczone dla okna tego modułu — poza wytwórnią. */
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

/** Historia z rdzenia w postaci wiersza widoku. */
function wierszZHistorii(wiersz: TerminalOutputLine): WierszPodgladu {
  return {
    tekst: wiersz.text,
    bledny: wiersz.channel === TerminalOutputChannel.Stderr,
    zrodlo: `karta ${wiersz.terminalSessionId}`,
  };
}

/**
 * Opis źródła wiersza żywego. Fragment niesie identyfikator procesu, nie karty,
 * bo kontrakt `stream.chunk` karty nie zna; okno nazywa więc to, co dostało,
 * zamiast podstawiać kartę, której rdzeń w tym zdarzeniu nie podał.
 */
function opisZrodlaProcesu(idProcesu: string): string {
  return idProcesu === '' ? 'proces nieznany' : `proces ${idProcesu}`;
}

/** Przycina widok do pojemności i zapamiętuje, że przycinał. */
function przytnij(kontekst: KontekstPodgladu): void {
  const nadmiar = kontekst.wiersze.length - POJEMNOSC_WIDOKU;
  if (nadmiar <= 0) return;
  kontekst.wiersze.splice(0, nadmiar);
  kontekst.uciety = true;
}

/**
 * Zdanie stanu pustego. Pustka bywa poprawna, więc zdanie nazywa tę właściwą:
 * przed pierwszym odczytem, po odczycie bez wierszy, po odczycie bez zapisu
 * na żywo.
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
 * Zdanie o braku zapisu na żywo — jedno miejsce, bo pada w trzech stanach okna.
 *
 * Powód jest jeden: moduł nie ma okna nadanego przez rdzeń, więc żądanie poszło
 * bez `windowId`, a rdzeń odpowiedział `subscribed: false`. To nie odmowa i nie
 * awaria, ale przemilczane kazałoby czytać okno jako żywe i zamarłe zarazem.
 */
const ZDANIE_BEZ_ZAPISU =
  'Okno NIE jest zapisane na żywo: moduł nie ma okna nadanego przez rdzeń, więc żądanie poszło ' +
  'bez wskazania okna i podgląd pokazuje wyłącznie historię. Nowe wiersze przyjdą dopiero po ' +
  'naciśnięciu „Odśwież ogon".';

/** Potwierdzenie czynności — rozdziela zapis na żywo od samego odczytu historii. */
function zdaniePotwierdzenia(odpowiedz: TerminalOutputStreamResponse): string {
  const ile = `Ogon historii: ${odpowiedz.lines.length} wierszy.`;
  return odpowiedz.subscribed
    ? `${ile} Okno zapisane na zbiorcze wyjście kart — nowe wiersze dochodzą na żywo.`
    : `${ile} ${ZDANIE_BEZ_ZAPISU}`;
}

/**
 * Widok wierszy: znacznik ucięcia, ostrzeżenie o braku zapisu, treść.
 *
 * Klasy noszą przedrostek `dnp-` obszaru okien pomocniczych, a nie przedrostek
 * modułu: ten sam widok stoi w Developerze i w Diagnostics, więc jego wygląd
 * ma jedno miejsce (`pomocnicze.css`), a nie dwa arkusze do rozjechania się.
 * Przedrostek modułu zostaje przy stanach treści, bo tam pokrycie już jest.
 */
function rysujWiersze(kontekst: KontekstPodgladu): DocumentFragment {
  const widok = document.createDocumentFragment();
  if (kontekst.uciety) {
    widok.append(uwaga(`Widok pokazuje ostatnie ${POJEMNOSC_WIDOKU} wierszy — starsze z niego wypadły.`));
  }
  // Uwaga o braku zapisu dopiero po odczycie: przed nim okno nie wie, czy jest
  // zapisane, a zdanie o niezapisaniu byłoby orzeczeniem bez podstawy.
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

/** Uwaga nad treścią — nie zastępuje wierszy, stoi obok nich. */
function uwaga(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis dnp-uwaga';
  element.setAttribute('role', 'note');
  element.textContent = zdanie;
  return element;
}
