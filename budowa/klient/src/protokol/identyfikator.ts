/**
 * Identyfikatory nadawane po stronie klienta.
 *
 * Kontrakt wymaga identyfikatora w każdej kopercie: odpowiedź i fragmenty
 * strumienia powtarzają identyfikator żądania, dzięki czemu klient wiąże
 * wynik z wywołaniem.
 */
export function nowyIdentyfikator(przedrostek: string): string {
  return `${przedrostek}-${losowyCzlon()}`;
}

/** Człon losowy identyfikatora: UUID środowiska, a przy jego braku znacznik czasu i liczba pseudolosowa łączone w jeden ciąg. */
function losowyCzlon(): string {
  const kryptografia = globalThis.crypto;
  if (kryptografia && typeof kryptografia.randomUUID === 'function') {
    return kryptografia.randomUUID();
  }
  return `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 10)}`;
}
