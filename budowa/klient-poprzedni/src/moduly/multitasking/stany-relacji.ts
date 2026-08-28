import { QueueStatus, type LoopState, type Message, type Queue, type Window } from '../../../../shared/contract';

// Stany relacji koordynator–wykonawca stoją w jednym miejscu jako czyste funkcje danych.

/** Waga znacznika ramy okna informacyjnego; podzbiór pełnej wagi znacznika, bez wartości oznaczającej błąd. */
export type WagaStanu = 'neutralna' | 'sukces' | 'ostrzezenie';

/** Jeden stan relacji: napis na plakietce, jego waga, kod stały oraz pełne zdanie, które go dokładnie wyjaśnia. */
export interface StanRelacji {
  /** Napis plakietki — krótki, bo rama daje mu jeden wiersz. */
  readonly tekst: string;
  readonly waga: WagaStanu;
  // Kod stanu dla znacznika ramy — po nim odczyt automatyczny rozpoznaje stan bez czytania napisu.
  readonly kod: string;
  /** Pełne zdanie: skąd ten stan się wziął. Idzie w `aria-description` ramy. */
  readonly powod: string;
}

/** Przesłanki znacznika koordynatora — wszystkie odczytane wprost, żadna z nich nie jest liczona tutaj. */
export interface PrzeslankiKoordynatora {
  /** Bieg z `window.state.get`; pusty, dopóki odczyt nie wrócił. */
  bieg: LoopState | null;
  /** Kolejka etapu bieżącego; pusta, gdy plan nie ma ani jednego etapu. */
  kolejka: Queue | null;
  /** Czy któryś wykonawca tej obsady prowadzi turę. */
  wykonawcaWTurze: boolean;
}

/**
 * Znacznik okna koordynatora.
 *
 * Czysta funkcja danych: bieg, kolejkę i tury bierze odczytane, więc nie zna ani
 * stanu wspólnego, ani rdzenia.
 */
export function znacznikKoordynatora(przeslanki: PrzeslankiKoordynatora): StanRelacji {
  const { bieg, kolejka, wykonawcaWTurze } = przeslanki;

  if (bieg === null) {
    return {
      tekst: 'bieg nie rozpoczęty',
      waga: 'neutralna',
      kod: 'bieg-nierozpoczety',
      powod:
        'Rdzeń nie oddał jeszcze stanu biegu tego okna (window.state.get bez pola loop) — to nie znaczy „zero obiegów", tylko „pomiaru nie ma".',
    };
  }

  if (bieg.stopped) {
    return {
      tekst: 'bieg zatrzymany',
      waga: 'ostrzezenie',
      kod: 'bieg-zatrzymany',
      powod: `Rdzeń zatrzymał bieg po ${bieg.loops} obiegach. Powód: ${bieg.stopReason ?? 'rdzeń go nie podał'}. Bez postępu ${bieg.loopsWithoutProgress} z ${bieg.threshold}.`,
    };
  }

  if (kolejka !== null && kolejka.status === QueueStatus.Paused) {
    return {
      tekst: 'kolejka wstrzymana',
      waga: 'ostrzezenie',
      kod: 'kolejka-wstrzymana',
      powod: `Etap „${kolejka.name ?? kolejka.id}" stoi w stanie ${kolejka.status} — wznawia go „Wznów" (queue.action resume). Bieg koordynatora nie jest zatrzymany: obiegów ${bieg.loops}.`,
    };
  }

  if (wykonawcaWTurze) {
    return {
      tekst: 'wykonawca pracuje',
      waga: 'neutralna',
      kod: 'wykonawca-pracuje',
      powod: `Co najmniej jeden wykonawca tej obsady prowadzi turę — koordynator czeka na jej zamknięcie. Obiegów dotąd: ${bieg.loops}.`,
    };
  }

  if (bieg.loops > 0) {
    const poKim =
      bieg.lastExecutorWindowId === undefined
        ? 'rdzeń nie podał, czyja tura go wybudziła'
        : `ostatnia tura: okno ${bieg.lastExecutorWindowId}`;
    const zPowodu =
      bieg.lastTurnReason === undefined ? '' : `, powód zamknięcia tury: ${bieg.lastTurnReason}`;
    return {
      tekst: `koordynator wybudzony · obiegów ${bieg.loops}`,
      waga: 'sukces',
      kod: 'koordynator-wybudzony',
      powod: `Pętla rdzenia (session.Petla) zliczyła ${bieg.loops} obiegów i żadna tura wykonawcy nie biegnie — pałeczkę trzyma koordynator (${poKim}${zPowodu}). Licznik prowadzi rdzeń, nie ten widok.`,
    };
  }

  return {
    tekst: 'bieg bez obiegu',
    waga: 'neutralna',
    kod: 'bieg-bez-obiegu',
    powod:
      'Rdzeń oddał stan biegu, ale licznik obiegów stoi na zerze — żadna tura wykonawcy jeszcze się nie domknęła, więc koordynatora nie było czym wybudzić.',
  };
}

/**
 * Znacznik okna wykonawcy.
 *
 * Czysta funkcja danych: okno, znacznik tury i ostatni wynik przychodzą
 * odczytane, więc funkcja nie zna ani stanu wspólnego, ani rdzenia.
 */
export function znacznikWykonawcy(
  okno: Window | null,
  wTurze: boolean,
  wynik: Message | null,
): StanRelacji {
  if (okno === null) {
    return {
      tekst: 'brak okna',
      waga: 'ostrzezenie',
      kod: 'brak-okna',
      powod:
        'Obsada nie ma jeszcze okna dla tego wykonawcy — załóż je przyciskiem „Załóż wykonawcę" w panelu obsady.',
    };
  }
  if (wTurze) {
    return {
      tekst: 'tura w biegu',
      waga: 'neutralna',
      kod: 'tura-w-biegu',
      powod: `Okno ${okno.id} prowadzi turę — fragmenty idą strumieniem, wynik jeszcze się nie domknął.`,
    };
  }
  // Waga idzie za wynikiem, nie za istnieniem okna, bo barwa sukcesu przy braku wyniku mówiłaby inaczej.
  if (wynik === null) {
    return {
      tekst: 'bez wyniku',
      waga: 'neutralna',
      kod: 'bez-wyniku',
      powod: `Okno ${okno.id} stoi, ale rdzeń nie oddał jeszcze ani jednej domkniętej odpowiedzi tego wykonawcy.`,
    };
  }
  return {
    tekst: 'wynik gotowy',
    waga: 'sukces',
    kod: 'wynik-gotowy',
    powod: `Ostatnia domknięta odpowiedź okna ${okno.id}: wiadomość ${wynik.id} w stanie ${wynik.status}.`,
  };
}

/**
 * Wiesza stan relacji na ramie okna: plakietkę, kod stały do odczytu
 * automatycznego i pełne zdanie w opisie dostępności, a gdy wołający poda
 * miejsce — także w widocznym wierszu okna.
 */
export function powiesStanRelacji(
  rama: { element: HTMLElement; ustawZnacznik(tekst: string, waga?: WagaStanu): void },
  stan: StanRelacji,
  wierszPowodu?: HTMLElement,
): void {
  rama.ustawZnacznik(stan.tekst, stan.waga);
  rama.element.dataset['stanRelacji'] = stan.kod;
  rama.element.setAttribute('aria-description', stan.powod);
  if (wierszPowodu !== undefined) wierszPowodu.textContent = stan.powod;
}
