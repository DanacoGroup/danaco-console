import { MessageStatus, type Message, type Window } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { poleTresci, przyciskAkcji as przycisk, wybor } from '../../modele/kontrolki-formularza';
import type { WykazKomendRdzenia } from './braki-kontraktu';
import { WIERSZE_POLA, wierszOpisu } from './kontrolki';
import { utworzPanelPodagentow } from './panel-subagent-network';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { powiesStanRelacji, znacznikWykonawcy } from './stany-relacji';
import { NAZWY_TRYBOW, trybZWartosci, zapiszTrybNaKoordynatorze } from './tryby-wspolpracy';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloBiegu } from './zrodlo-biegu';
import type { ZrodloOkien } from './zrodlo-okien';
import type { ZrodloPodagentow } from './zrodlo-podagentow';

/**
 * Executor Chat — okno wiodące wykonawcy; na scenie stoją dwa wystąpienia.
 *
 * Wykonawca wykonuje, nie zarządza: okno nie ma sterowania kolejką ani planu
 * etapów. Ma polecenie, przerwanie tury, panel Subagent Network i wgląd we
 * własny wynik.
 *
 * Przynależność do koordynatora pochodzi z kontraktu. `coordinatorWindowId`
 * okna wykonawczego jest jedynym oznaczeniem więzi, więc okno wypisuje je
 * wprost — inaczej wykonawca własnej pary jest nie do odróżnienia od cudzego.
 *
 * Tryb współpracy jest wspólny obu wykonawcom. Rozstrzyga, kto dostaje zlecenie
 * przy przekazaniu, więc zapisuje się na oknie koordynatora, a nie osobno
 * w każdym wykonawcy — dwa zapisy tej samej rzeczy rozjeżdżałyby się przy
 * pierwszej zmianie.
 */
export interface OknoWykonawcy {
  element: HTMLElement;
  /** Przerysowuje okno po zmianie stanu wspólnego. */
  odswiez(): void;
  /**
   * Odpina nasłuchy okna. Woła je zamknięcie modułu — bez tego nasłuch
   * `subagent.changed` przeżywałby scenę i odświeżał widok zdjęty z ekranu.
   */
  rozlacz(): void;
}

export interface OpcjeWykonawcy {
  okna: ZrodloOkien;
  bieg: ZrodloBiegu;
  stan: StanMultitaskingu;
  /** Źródło `subagent.list`/`subagent.result.collect` — żywy wykaz panelu. */
  podagenci: ZrodloPodagentow;
  /** Wykaz komend rdzenia — stąd bierze się powód nieczynnych kontrolek. */
  komendy: WykazKomendRdzenia;
  /** Numer wykonawcy na scenie: 1 albo 2. */
  numer: number;
}

/** Co w ramie obu wykonawców jest wspólne — różni je wyłącznie kod i tytuł. */
const RAMA_WYKONAWCY = {
  rola: 'wiodące',
  przeznaczenie: 'Wykonanie zadania z kolejki, uruchomienie podagentów i zwrot wyniku.',
  ogniskowalne: true,
} as const;

export function utworzOknoWykonawcy(opcje: OpcjeWykonawcy): OknoWykonawcy {
  const { stan, numer } = opcje;

  // Kod okna stoi wprost, osobno na każde wystąpienie. Kodu składanego z numeru
  // nie da się ani wyszukać w drzewie, ani zestawić z katalogiem rdzenia,
  // a numer spoza zakresu dawałby kod pusty — okno zbudowane bez wiersza
  // katalogu, o czym nikt by nie zameldował.
  const rama =
    numer === 1
      ? utworzRameOkna({
          kod: 'executor-chat-1',
          tytul: 'Executor Chat (Executor 1)',
          ...RAMA_WYKONAWCY,
        })
      : utworzRameOkna({
          kod: 'executor-chat-2',
          tytul: 'Executor Chat (Executor 2)',
          ...RAMA_WYKONAWCY,
        });

  const tresci = utworzStanTresci();

  const podagenci = utworzPanelPodagentow({
    zrodlo: opcje.okna,
    komendy: opcje.komendy,
    sesja: () => stan.sesja(),
    potwierdz: (zdanie, udane) => tresci.potwierdzenie(zdanie, udane),
    // Żywy wykaz powołanych: okno bierze się z obsady TEGO wykonawcy,
    // odczytywanej w chwili pytania — obsada bywa późniejsza niż panel.
    zywi: {
      podagenci: opcje.podagenci,
      okno: () => moje()?.id ?? '',
    },
  });

  const { polecenie, wyslij, zatrzymaj, trybPracy, opis, powodStanu } =
    zlozPowierzchnieWykonawcy(rama, tresci.element, podagenci.element, opcje.komendy);

  wyslij.addEventListener('click', () => {
    void wyslijPolecenie();
  });
  zatrzymaj.addEventListener('click', () => {
    void zatrzymajTureWykonawcy(opcje.bieg, moje(), tresci);
  });
  trybPracy.addEventListener('change', () => {
    void zapiszTrybNaKoordynatorze(opcje.okna, stan, trybZWartosci(trybPracy.value)).then(
      ({ zdanie, udane }) => tresci.potwierdzenie(zdanie, udane),
    );
  });

  /** Okno tego wykonawcy z obsady; puste, dopóki nie powstało. */
  function moje(): Window | null {
    return stan.obsada().wykonawcy[numer - 1] ?? null;
  }

  /**
   * Wydanie polecenia wykonawcy.
   *
   * Zdanie potwierdzenia mówi o tym, co oddał rdzeń: `message.send` zwraca
   * zapisaną wiadomość Operatora i nic nie orzeka o turze modelu. O otwarciu
   * tury okno pisze tylko wtedy, gdy stan wiadomości jest strumieniowy.
   */
  async function wyslijPolecenie(): Promise<void> {
    const okno = moje();
    if (okno === null) {
      tresci.blad(`Executor ${numer} nie istnieje — załóż okno wykonawcy w obsadzie.`);
      return;
    }
    const tresc = polecenie.value.trim();
    if (tresc === '') {
      tresci.potwierdzenie('Polecenie puste — tura się nie zaczyna.', false);
      return;
    }
    const wynik = await opcje.bieg.wyslij({ windowId: okno.id, content: tresc });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.blad('Rdzeń odmówił przyjęcia polecenia.', wynik.blad);
      return;
    }
    const zapisana = wynik.wynik;
    polecenie.value = '';
    stan.zapamietajPrzekazanie(okno.id);
    tresci.potwierdzenie(zdanieOPoleceniu(zapisana, okno.id), zapisana.windowId === okno.id);
  }

  /**
   * Identyfikator okna, przy którym panel podagentów odczytano ostatnio.
   *
   * Żywy wykaz podagentów odczytuje się przy zmianie okna wykonawcy, a nie przy
   * każdym powiadomieniu: `stan.naZmiane` woła się także na fragment strumienia,
   * więc odczyt bezwarunkowy zasypałby rdzeń wywołaniami `subagent.list`
   * w tempie strumienia.
   */
  let oknoPanelu = '';

  function odswiez(): void {
    const okno = moje();
    trybPracy.value = stan.tryb();
    const idOkna = okno?.id ?? '';
    if (idOkna !== oknoPanelu) {
      oknoPanelu = idOkna;
      podagenci.odswiezZywych();
    }
    opis.replaceChildren(
      ...wieziWykonawcy(okno, okno === null ? 0 : stan.strumien().liczbaWywolan(okno.id)),
    );
    // „Wyślij polecenie" pozostaje czynne także bez okna wykonawcy: czynność
    // zaczyna się od sprawdzenia tej samej przesłanki i mówi o niej pełnym
    // zdaniem („Executor N nie istnieje — załóż okno wykonawcy w obsadzie"),
    // a wygaszenie odebrałoby wyłącznie powód.
    const brakOkna = okno === null;
    wyslij.title = brakOkna ? `Executor ${numer} nie ma jeszcze okna — naciśnięcie powie to wprost.` : '';
    if (brakOkna) wyslij.dataset['przeslanka'] = 'brak okna wykonawcy';
    else delete wyslij.dataset['przeslanka'];
    // Rachunek plakietki stoi w `stany-relacji.ts` razem ze znacznikiem
    // koordynatora — oba są jednym zagadnieniem („czyja jest teraz kolej")
    // i jednym sprawdzianem.
    powiesStanRelacji(
      rama,
      znacznikWykonawcy(
        okno,
        okno !== null && stan.czyTura(okno.id),
        okno === null ? null : stan.wynikWykonawcy(okno.id),
      ),
      powodStanu,
    );
  }

  stan.naZmiane(odswiez);
  odswiez();

  return {
    element: rama.element,
    odswiez() {
      odswiez();
      podagenci.odswiez();
    },
    rozlacz() {
      podagenci.rozlacz();
    },
  };
}

/** Kontrolki własne okna wykonawcy wraz z miejscem na opis więzi. */
interface PowierzchniaWykonawcy {
  polecenie: HTMLTextAreaElement;
  wyslij: HTMLButtonElement;
  zatrzymaj: HTMLButtonElement;
  trybPracy: HTMLSelectElement;
  opis: HTMLElement;
  /** Pełne zdanie o stanie relacji; plakietka nagłówka mieści dwa słowa. */
  powodStanu: HTMLElement;
}

/**
 * Składa kontrolki, pasek akcji, narzędzia i ciało okna wykonawcy.
 *
 * Czysta konstrukcja: nie domyka się ani na stanie wspólnym, ani na źródłach
 * rdzenia — panel podagentów bierze gotowy, a nasłuchy przycisków zakłada
 * wytwórnia, więc fragment wyszedł bez przenoszenia zależności.
 */
function zlozPowierzchnieWykonawcy(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  panelPodagentow: HTMLElement,
  komendy: WykazKomendRdzenia,
): PowierzchniaWykonawcy {
  const polecenie = poleTresci(
    'Polecenie wykonawcy',
    WIERSZE_POLA,
    'Treść polecenia dla tego wykonawcy…',
  );
  const wyslij = przycisk('Wyślij polecenie', 'dn-btn dn-btn--atrament');

  // Zatrzymanie czynne zawsze, także gdy okno nie prowadzi tury:
  // odpowiedź `stopped=false` jest informacją, nie usterką.
  const zatrzymaj = przycisk('Zatrzymaj turę', 'dn-btn dn-btn--niebezpieczny');
  const trybPracy = wybor('Tryb współpracy wykonawców', NAZWY_TRYBOW);

  const opis = document.createElement('div');
  opis.className = 'dm-wiez';

  const powodStanu = document.createElement('p');
  powodStanu.className = 'dm-wiez';
  powodStanu.dataset['rola'] = 'powod-stanu-relacji';

  rama.akcje.append(wyslij, zatrzymaj);
  rama.narzedzia.append(trybPracy, komendy.przyciskNieczynny('Nadaj rolę', 'role.assign', 'role.update'));
  rama.cialo.append(opis, powodStanu, polecenie, panelPodagentow, stanTresci, komendy.wykazBrakow([
    'subagent.spawn',
    'subagent.list',
    'subagent.result.collect',
    'role.update',
  ]));

  return { polecenie, wyslij, zatrzymaj, trybPracy, opis, powodStanu };
}

/**
 * Wiersze więzi wykonawcy: okno, rola z kontraktu, koordynator, licznik wywołań.
 *
 * Czysta funkcja danych: liczbę wywołań dostaje policzoną, więc nie sięga ani do
 * strumienia, ani do stanu wspólnego.
 */
function wieziWykonawcy(okno: Window | null, wywolan: number): HTMLElement[] {
  return [
    wierszOpisu('Okno wykonawcy', okno?.id ?? 'nie założone'),
    wierszOpisu('Rola z kontraktu', okno?.windowRole ?? 'brak'),
    wierszOpisu('Koordynator', okno?.coordinatorWindowId ?? 'nie przypięty'),
    wierszOpisu('Wywołań narzędzi', String(wywolan)),
  ];
}

/**
 * Zdanie o przyjętym poleceniu — złożone z wiadomości, którą oddał rdzeń.
 *
 * Czysta funkcja danych: wiadomość i oczekiwane okno bierze wprost, więc nie zna
 * ani stanu wspólnego, ani okna wykonawcy.
 */
function zdanieOPoleceniu(wiadomosc: Message, oczekiwaneOkno: string): string {
  if (wiadomosc.windowId !== oczekiwaneOkno) {
    return `Rdzeń zapisał wiadomość ${wiadomosc.id} w oknie ${wiadomosc.windowId}, a polecenie szło do ${oczekiwaneOkno} — czynność trafiła gdzie indziej.`;
  }
  const stanTury =
    wiadomosc.status === MessageStatus.Streaming
      ? 'tura wykonawcy otwarta'
      : 'o otwarciu tury rdzeń jeszcze nic nie powiedział — czekaj na strumień';
  return `Rdzeń zapisał polecenie jako wiadomość ${wiadomosc.id} (stan ${wiadomosc.status}); ${stanTury}.`;
}

/**
 * Przerwanie tury wykonawcy.
 *
 * Wyjęte, bo domyka się wyłącznie na źródle biegu, oknie i stanie treści — a te
 * trzy bierze parametrami; z wnętrza okna nie potrzebuje niczego.
 */
async function zatrzymajTureWykonawcy(
  bieg: ZrodloBiegu,
  okno: Window | null,
  tresci: StanTresci,
): Promise<void> {
  if (okno === null) {
    tresci.potwierdzenie('Nie ma okna, w którym dałoby się przerwać turę.', false);
    return;
  }
  const wynik = await bieg.zatrzymaj({ windowId: okno.id });
  if (!wynik.udany) {
    tresci.blad('Rdzeń odmówił zatrzymania tury.', wynik.blad);
    return;
  }
  tresci.potwierdzenie(
    wynik.wynik?.stopped === true ? 'Tura przerwana.' : 'Żadna tura nie biegła.',
    true,
  );
}

