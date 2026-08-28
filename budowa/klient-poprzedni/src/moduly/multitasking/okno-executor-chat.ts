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
 * Executor Chat jest oknem wiodącym wykonawcy: wykonuje, nie zarządza, ma
 * polecenie, przerwanie tury, panel podagentów i wgląd we własny wynik;
 * przynależność do koordynatora i tryb współpracy pochodzą z kontraktu.
 */
export interface OknoWykonawcy {
  element: HTMLElement;
  /** Przerysowuje okno po zmianie stanu wspólnego. */
  odswiez(): void;
  // Odpina nasłuchy okna, wołane przez zamknięcie modułu, żeby nasłuch nie przeżywał zdjętej sceny.
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

/** Co w ramie obu wykonawców jest wspólne — różni je wyłącznie kod okna oraz jego tytuł widoczny w nagłówku. */
const RAMA_WYKONAWCY = {
  rola: 'wiodące',
  przeznaczenie: 'Wykonanie zadania z kolejki, uruchomienie podagentów i zwrot wyniku.',
  ogniskowalne: true,
} as const;

export function utworzOknoWykonawcy(opcje: OpcjeWykonawcy): OknoWykonawcy {
  const { stan, numer } = opcje;

  // Kod okna stoi wprost, osobno na każde wystąpienie, bo kodu składanego z numeru nie da się wyszukać.
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
    // Żywy wykaz powołanych bierze się z obsady tego wykonawcy, odczytywanej w chwili pytania.
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

  // Zdanie potwierdzenia mówi o tym, co oddał rdzeń; o otwarciu tury pisze tylko przy stanie streamingu.
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

  // Żywy wykaz podagentów odczytuje się przy zmianie okna wykonawcy, a nie przy każdym powiadomieniu.
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
    // Wyślij polecenie pozostaje czynne bez okna wykonawcy: czynność mówi o przesłance pełnym zdaniem.
    const brakOkna = okno === null;
    wyslij.title = brakOkna ? `Executor ${numer} nie ma jeszcze okna — naciśnięcie powie to wprost.` : '';
    if (brakOkna) wyslij.dataset['przeslanka'] = 'brak okna wykonawcy';
    else delete wyslij.dataset['przeslanka'];
    // Rachunek plakietki stoi w osobnym pliku razem ze znacznikiem koordynatora, jednym sprawdzianem.
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

/** Kontrolki własne okna wykonawcy wraz z miejscem na opis więzi i stan relacji z jego koordynatorem pracy. */
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
 * Składa kontrolki, pasek akcji, narzędzia i ciało okna wykonawcy; czysta
 * konstrukcja, bo panel podagentów bierze gotowy, a nasłuchy przycisków
 * zakłada wytwórnia.
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

  // Zatrzymanie czynne zawsze, także gdy okno nie prowadzi tury: brak zatrzymania jest informacją.
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

