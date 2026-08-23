import { MessageRole, WindowRole } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/**
 * Dziewięć rodzajów nadawcy wpisu w historii okna komunikacji.
 *
 * Rodzaje rozróżniają się trzema nośnikami naraz i ani jeden z nich nie jest
 * barwą tła:
 *   1. ikona w medalionie pierwszej kolumny wpisu,
 *   2. barwa medalionu i kreski krawędzi — wyłącznie w rozdzielczości klasy
 *      semantycznej, nie rodzaju (patrz `KlasaNadawcy` niżej),
 *   3. etykieta słowna wersalikami.
 * Barwa sama nie może być jedynym nośnikiem znaczenia, więc tło wpisu pozostaje
 * jedno dla wszystkich dziewięciu rodzajów.
 */
export const RodzajNadawcy = {
  /** Operator przy klawiaturze. */
  Uzytkownik: 'uzytkownik',
  /** Model w oknie samodzielnym. */
  Model: 'model',
  /** Agent zbudowany przez Operatora. */
  Agent: 'agent',
  /** Koordynator pętli — planuje, nie wytwarza. */
  Koordynator: 'koordynator',
  /** Wykonawca pętli — zamknięcie tury wybudza koordynatora. */
  Wykonawca: 'wykonawca',
  /** Walidator — ocenia wynik kroku kolejki. */
  Walidator: 'walidator',
  /** Automatyzacja — proces, harmonogram, kolejka. */
  Automatyzacja: 'automatyzacja',
  /** Always On Display — pływający doradca całej platformy. */
  Aod: 'aod',
  /** Wynik narzędzia zwrócony modelowi. */
  Narzedzie: 'narzedzie',
} as const;

export type RodzajNadawcy = (typeof RodzajNadawcy)[keyof typeof RodzajNadawcy];

/**
 * Klasa semantyczna nadawcy — trzy, nie dziewięć.
 *
 * Arkusz `budowa/client/src/komponenty/wpis.css` zna dokładnie trzy
 * modyfikatory `.dn-wpis--czlowiek` / `--inteligencja` / `--system` i wiąże
 * z nimi barwę kreski oraz barwę medalionu:
 *
 *   człowiek      atrament   `--dn-atrament`             (Operator)
 *   inteligencja  sygnał     `--dn-sygnal-wypelnienie`   (model, agent,
 *                                                        koordynator,
 *                                                        wykonawca, walidator)
 *   system        neutralna  `--dn-obrys-mocny`          (automatyzacja,
 *                                                        Always On Display,
 *                                                        wynik narzędzia)
 *
 * Rodzaj wewnątrz klasy różnicuje ikona i etykieta, nigdy kolor.
 */
export type KlasaNadawcy = 'czlowiek' | 'inteligencja' | 'system';

/** Znaki rozpoznawcze nadawcy: etykieta, ikona i klasa semantyczna. */
export interface ZnakiNadawcy {
  etykieta: string;
  ikona: NazwaIkony;
  klasa: KlasaNadawcy;
}

const ZNAKI: Readonly<Record<RodzajNadawcy, ZnakiNadawcy>> = {
  uzytkownik: { etykieta: 'Operator', ikona: 'uzytkownik', klasa: 'czlowiek' },
  model: { etykieta: 'Model', ikona: 'gwiazdka', klasa: 'inteligencja' },
  agent: { etykieta: 'Agent', ikona: 'tarcza', klasa: 'inteligencja' },
  koordynator: { etykieta: 'Koordynator', ikona: 'waga', klasa: 'inteligencja' },
  wykonawca: { etykieta: 'Wykonawca', ikona: 'uruchom', klasa: 'inteligencja' },
  walidator: { etykieta: 'Walidator', ikona: 'ptaszek-kolo', klasa: 'inteligencja' },
  automatyzacja: { etykieta: 'Automatyzacja', ikona: 'odswiez', klasa: 'system' },
  aod: { etykieta: 'Always On Display', ikona: 'dzwonek', klasa: 'system' },
  narzedzie: { etykieta: 'Wynik narzędzia', ikona: 'kod', klasa: 'system' },
};

/** Znaki rozpoznawcze wskazanego rodzaju nadawcy. */
export function znakiNadawcy(rodzaj: RodzajNadawcy): ZnakiNadawcy {
  return ZNAKI[rodzaj];
}

/**
 * Rozpoznaje nadawcę wpisu z roli wiadomości i roli okna.
 *
 * Rola wiadomości pochodzi z kontraktu (`MessageRole`) i zna cztery wartości;
 * rola okna (`WindowRole`) rozdziela wypowiedź modelu na koordynatora
 * i wykonawcę pętli. Rola systemowa oznacza wypowiedź warstwy
 * automatycznej platformy, a nie żadnego z ośmiu pozostałych nadawców.
 */
export function rozpoznajNadawce(
  rola: MessageRole,
  rolaOkna: WindowRole | null = null,
): RodzajNadawcy {
  switch (rola) {
    case MessageRole.User:
      return RodzajNadawcy.Uzytkownik;
    case MessageRole.Tool:
      return RodzajNadawcy.Narzedzie;
    case MessageRole.System:
      return RodzajNadawcy.Automatyzacja;
    case MessageRole.Assistant:
      return nadawcaModelu(rolaOkna);
  }
}

/** Wypowiedź modelu w oknie o wskazanej roli. */
function nadawcaModelu(rolaOkna: WindowRole | null): RodzajNadawcy {
  switch (rolaOkna) {
    case WindowRole.Coordinator:
      return RodzajNadawcy.Koordynator;
    case WindowRole.Executor:
      return RodzajNadawcy.Wykonawca;
    default:
      return RodzajNadawcy.Model;
  }
}
