import { ChunkKind } from '../../../shared/contract';

/**
 * Rodzaje fragmentu strumienia odpowiedzi.
 *
 * Wszystkie rodzaje pochodzą wprost z kontraktu — moduł ich nie przepisuje,
 * wyłącznie nadaje im nazwy dziedziny. Zmiana nazwy w `shared/contract.json`
 * przerywa kompilację tego pliku.
 *
 * Prowenancja wywołania i metadane konta kanału są diagnostyczne i przelotowe:
 * strumień je niesie, baza ich nie zapisuje.
 */
export const RodzajFragmentu = {
  /** Tekst odpowiedzi — narasta w treści wpisu. */
  Tekst: ChunkKind.Text,
  /** Tok rozumowania — podgląd pracy modelu na żywo, zwijany. */
  Rozumowanie: ChunkKind.Thinking,
  /** Wywołanie narzędzia przez model. */
  WywolanieNarzedzia: ChunkKind.ToolUse,
  /** Wynik narzędzia zwrócony modelowi. */
  WynikNarzedzia: ChunkKind.ToolResult,
  /** Treść obrazowa. */
  Obraz: ChunkKind.Image,
  /** Treść dźwiękowa. */
  Dzwiek: ChunkKind.Audio,
  /** Błąd w trakcie strumienia; kończy wywołanie, nie sesję. */
  Blad: ChunkKind.Error,
  /** Prowenancja wywołania: argv i prompt systemowy. */
  Prowenancja: ChunkKind.Provenance,
  /** Metadane konta użytego przez kanał. */
  Konto: ChunkKind.Account,
  /**
   * Wersja ostateczna odpowiedzi — zastępuje treść złożoną z fragmentów
   * tekstowych, nie dokłada się do niej. Niesie ją koperta domykająca strumień.
   */
  WersjaOstateczna: ChunkKind.Final,
} as const;

export type RodzajFragmentu = ChunkKind;
