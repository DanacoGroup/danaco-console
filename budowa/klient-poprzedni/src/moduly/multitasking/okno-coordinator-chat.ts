import { type Window, type WindowStateGetResponse } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { poleTresci, przyciskAkcji, wybor } from '../../modele/kontrolki-formularza';
import type { WykazKomendRdzenia } from './braki-kontraktu';
import { WIERSZE_POLA } from './kontrolki';
import { KOD_KOORDYNATORA } from './obsada-rol';
import { utworzPlanEtapow } from './plan-etapow';
import { utworzSterowanie } from './sterowanie-koordynatora';
import { utworzStanTresci } from './stany-okna';
import { powiesStanRelacji, znacznikKoordynatora } from './stany-relacji';
import { NAZWY_TRYBOW, trybZWartosci, zapiszTrybNaKoordynatorze } from './tryby-wspolpracy';
import { utworzWidokStrumienia } from './widok-strumienia';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloBiegu } from './zrodlo-biegu';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Coordinator Chat jest oknem zarządcy pętli koordynator–wykonawca: ma plan
 * etapów, kreator promptu, przyciski sterowania i podgląd strumienia, a pętla
 * żyje w rdzeniu, więc klient tylko odczytuje jej stan.
 */
export interface OknoKoordynatora {
  element: HTMLElement;
  /** Odczytuje stan biegu z rdzenia i przerysowuje okno. */
  odswiez(): void;
}

export interface OpcjeKoordynatora {
  okna: ZrodloOkien;
  bieg: ZrodloBiegu;
  stan: StanMultitaskingu;
  /** Wykaz komend rdzenia — stąd bierze się powód nieczynnych kontrolek. */
  komendy: WykazKomendRdzenia;
}

export function utworzOknoKoordynatora(opcje: OpcjeKoordynatora): OknoKoordynatora {
  const { okna, stan } = opcje;

  const rama = utworzRameOkna({
    tytul: 'Coordinator Chat',
    rola: 'zarządca',
    kod: KOD_KOORDYNATORA,
    przeznaczenie:
      'Planowanie etapów, podział pracy, sterowanie kolejką i przekazanie zlecenia wykonawcom.',
    ogniskowalne: true,
  });

  const tresci = utworzStanTresci();
  const { kreator, trybPracy } = kontrolkiKoordynatora();

  // Miejsce na pełne zdanie o stanie relacji; plakietka nagłówka mieści tylko dwa słowa.
  const stanRelacji = document.createElement('p');
  stanRelacji.className = 'dm-wiez';
  stanRelacji.dataset['rola'] = 'powod-stanu-relacji';

  trybPracy.addEventListener('change', () => {
    void zapiszTrybNaKoordynatorze(okna, stan, trybZWartosci(trybPracy.value)).then(
      ({ zdanie, udane }) => tresci.potwierdzenie(zdanie, udane),
    );
  });

  const sterowanie = utworzSterowanie({
    zrodlo: opcje.bieg,
    stan,
    trescZlecenia: () => kreator.value,
    zwaliduj: () => {
      void zwaliduj();
    },
    potwierdz: (zdanie, udane) => tresci.potwierdzenie(zdanie, udane),
  });

  const plan = utworzPlanEtapow({
    zrodlo: opcje.bieg,
    stan,
    komendy: opcje.komendy,
    potwierdz: (zdanie, udane) => tresci.potwierdzenie(zdanie, udane),
  });

  const podglad = utworzWidokStrumienia(stan);

  // Subskrypcja monitora zapisuje okno koordynatora na telemetrię postępu procesów jego sceny.
  const subskrybuj = przyciskAkcji('Subskrybuj monitor', 'dn-btn dn-btn--zarys');
  subskrybuj.addEventListener('click', () => {
    void subskrybujMonitor();
  });

  zmontujOknoKoordynatora(rama, opcje.komendy, {
    kreator,
    trybPracy,
    subskrybuj,
    stanRelacji,
    sterowanie: sterowanie.element,
    plan: plan.element,
    podglad: podglad.element,
    stanTresci: tresci.element,
  });

  // Bez wskazanego koordynatora nie ma czego zapisać — okno mówi to wprost, zamiast udawać obserwację.
  async function subskrybujMonitor(): Promise<void> {
    const koordynator = stan.obsada().koordynator;
    if (koordynator === null) {
      tresci.blad('Obsada nie ma koordynatora — nie ma okna, które zapisać na telemetrię.');
      return;
    }
    tresci.ladowanie('Zakładanie obserwacji monitora…');
    const wynik = await opcje.bieg.subskrybujMonitor({
      windowId: koordynator.id,
      sessionId: stan.sesja(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.blad('Rdzeń odmówił subskrypcji monitora.', wynik.blad);
      return;
    }
    const { subscribed, statuses } = wynik.wynik;
    // `subscribed` bierze się z odpowiedzi, nie z faktu wysłania — rdzeń mógłby
    // obserwacji nie założyć.
    tresci.potwierdzenie(
      subscribed
        ? `Rdzeń założył obserwację monitora na oknie koordynatora ${koordynator.id}; procesów w stanie bieżącym: ${statuses.length}.`
        : `Rdzeń NIE założył obserwacji (subscribed=false) — oddał sam stan bieżący ${statuses.length} procesów.`,
      subscribed,
    );
  }

  /** Walidacja przebiegu — powód jej zakresu stoi przy `raportPrzebieguKoordynatora`. */
  async function zwaliduj(): Promise<void> {
    const obsada = stan.obsada();
    if (obsada.koordynator === null) {
      tresci.blad('Nie ma koordynatora, którego bieg dałoby się sprawdzić.');
      return;
    }
    tresci.ladowanie('Sprawdzanie przebiegu…');
    const raport = await raportPrzebieguKoordynatora(
      [obsada.koordynator, ...obsada.wykonawcy],
      (zadanie) => okna.stanOkna(zadanie),
    );
    tresci.tresc().append(raport);
  }

  function odswiez(): void {
    trybPracy.value = stan.tryb();
    sterowanie.odswiez();
    plan.odswiez();
    podglad.odswiez();
    // Rachunek znacznika stoi w osobnym pliku, bo jest czystą funkcją danych o biegu i kolejce etapu.
    powiesStanRelacji(
      rama,
      znacznikKoordynatora({
        bieg: stan.bieg(),
        kolejka: stan.kolejkaBiezaca(),
        wykonawcaWTurze: stan.obsada().wykonawcy.some((okno) => stan.czyTura(okno.id)),
      }),
      stanRelacji,
    );
  }

  stan.naZmiane(odswiez);
  odswiez();
  void odczytajBiegKoordynatora(okna, stan);

  return {
    element: rama.element,
    odswiez() {
      odswiez();
      void odczytajBiegKoordynatora(okna, stan);
    },
  };
}

/** Kontrolki własne okna koordynatora; panele ról stoją obok nich, poza tą powierzchnią samych kontrolek. */
interface PowierzchniaKoordynatora {
  kreator: HTMLTextAreaElement;
  trybPracy: HTMLSelectElement;
}

/**
 * Dwie kontrolki własne zarządcy.
 *
 * Czysta konstrukcja: ani kreator promptu, ani selektor trybu nie domykają się
 * na stanie okna, więc wyszły bez przenoszenia jakiejkolwiek zależności.
 */
function kontrolkiKoordynatora(): PowierzchniaKoordynatora {
  return {
    kreator: poleTresci(
      'Kreator promptu dla wykonawców',
      WIERSZE_POLA,
      'Treść zlecenia przekazywanego wykonawcom…',
    ),
    trybPracy: wybor('Tryb współpracy wykonawców', NAZWY_TRYBOW),
  };
}

/**
 * Montaż paska akcji, narzędzi i ciała okna zarządcy.
 *
 * Czysta konstrukcja: same wstawienia gotowych elementów w gotową ramę — funkcja
 * nie zna ani stanu wspólnego, ani źródeł rdzenia, więc dała się wyjąć wprost.
 */
function zmontujOknoKoordynatora(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  komendy: WykazKomendRdzenia,
  czesci: {
    kreator: HTMLElement;
    trybPracy: HTMLElement;
    subskrybuj: HTMLElement;
    /** Pełne zdanie o stanie relacji koordynator–wykonawca. */
    stanRelacji: HTMLElement;
    sterowanie: HTMLElement;
    plan: HTMLElement;
    podglad: HTMLElement;
    stanTresci: HTMLElement;
  },
): void {
  rama.akcje.append(czesci.sterowanie);
  rama.narzedzia.append(czesci.trybPracy, czesci.subskrybuj);
  // Podgląd strumienia stoi w ciele okna na stałe; raport walidacji ma go uzupełniać, nie wypychać.
  rama.cialo.append(
    // Stan relacji idzie pierwszy — mówi, czyja jest teraz kolej.
    czesci.stanRelacji,
    czesci.kreator,
    czesci.plan,
    czesci.podglad,
    czesci.stanTresci,
    // Dwie komendy, których rdzeń dla tego okna nie udostępnia: przypisanie roli i zależność etapów planu.
    komendy.wykazBrakow(['orchestration.dependency.set', 'role.assign']),
  );
}

/**
 * Raport walidacji przebiegu jest jedyną walidacją, którą kontrakt pokrywa:
 * odczyt stanu koordynatora i wykonawców mówi, czy tura trwa, ile wiadomości
 * zapisano i czy bieg stoi.
 */
async function raportPrzebieguKoordynatora(
  okna: readonly Window[],
  odczytStanu: ZrodloOkien['stanOkna'],
): Promise<HTMLElement> {
  const wiersze: string[] = [];
  for (const okno of okna) {
    const wynik = await odczytStanu({ windowId: okno.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      wiersze.push(`${okno.id}: odczyt odmówiony (kod ${wynik.blad?.code ?? 'brak'})`);
      continue;
    }
    wiersze.push(opisStanu(okno.id, wynik.wynik));
  }
  const raport = document.createElement('pre');
  raport.className = 'dm-walidacja';
  raport.textContent = wiersze.join('\n');
  return raport;
}

/**
 * Odczyt stanu koordynatora wraz z licznikiem obiegów; zasila znacznik okna.
 *
 * Wyjęte, bo poza źródłem okien i stanem wspólnym — obydwoma z parametru — nie
 * dotyka niczego z wnętrza okna.
 */
async function odczytajBiegKoordynatora(
  okna: ZrodloOkien,
  stan: StanMultitaskingu,
): Promise<void> {
  const koordynator = stan.obsada().koordynator;
  if (koordynator === null) return;
  const wynik = await okna.stanOkna({ windowId: koordynator.id });
  if (!wynik.udany || wynik.wynik === undefined) return;
  stan.ustawBieg(wynik.wynik.loop ?? null);
}

/** Jeden wiersz raportu walidacji przebiegu, złożony z pól odczytanego stanu okna koordynatora albo wykonawcy. */
function opisStanu(idOkna: string, stan: WindowStateGetResponse): string {
  const bieg =
    stan.loop === undefined
      ? 'bez biegu'
      : `obiegów ${stan.loop.loops}, bez postępu ${stan.loop.loopsWithoutProgress}/${stan.loop.threshold}${stan.loop.stopped ? `, ZATRZYMANY: ${stan.loop.stopReason ?? 'bez powodu'}` : ''}`;
  return `${idOkna}: rola ${stan.window.windowRole}, proces ${stan.processStatus}, wiadomości ${stan.messageCount}, tura ${stan.streaming ? 'trwa' : 'stoi'}, ${bieg}`;
}
