import { ChunkKind } from '../../../shared/contract';

/**
 * Rodzaje fragmentu strumienia odpowiedzi, pochodzące wprost z kontraktu i tylko nazwane
 * językiem dziedziny.
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
  // Wersja ostateczna zastępuje fragmenty tekstowe treści, niesiona kopertą domykającą
  // strumień.
  WersjaOstateczna: ChunkKind.Final,
} as const;

export type RodzajFragmentu = ChunkKind;
