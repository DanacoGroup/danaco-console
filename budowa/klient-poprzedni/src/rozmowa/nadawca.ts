import { MessageRole, WindowRole } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/**
 * Dziewięć rodzajów nadawcy wpisu w historii okna komunikacji, rozróżnianych ikoną, barwą
 * medalionu i etykietą słowną, nigdy barwą tła.
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
 * Klasa semantyczna nadawcy dzieli dziewięć rodzajów na trzy grupy: człowieka,
 * inteligencję i system, każdą z własną barwą kreski i medalionu.
 */
export type KlasaNadawcy = 'czlowiek' | 'inteligencja' | 'system';

/**
 * Znaki rozpoznawcze nadawcy: etykieta słowna wersalikami, przypisana ikona oraz jego
 * klasa semantyczna.
 */
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

/**
 * Znaki rozpoznawcze wskazanego rodzaju nadawcy: jego etykieta słowna, ikona oraz
 * przypisana klasa semantyczna.
 */
export function znakiNadawcy(rodzaj: RodzajNadawcy): ZnakiNadawcy {
  return ZNAKI[rodzaj];
}

/**
 * Rozpoznaje nadawcę wpisu na podstawie roli wiadomości z kontraktu oraz roli okna w pętli
 * koordynator-wykonawca.
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

/**
 * Wypowiedź modelu w oknie o wskazanej roli, zapisywana jako wpis automatyzacji w historii
 * rozmowy okna.
 */
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
