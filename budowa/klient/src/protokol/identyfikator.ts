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

/** Człon losowy: UUID środowiska, a przy jego braku czas i liczba pseudolosowa. */
function losowyCzlon(): string {
  const kryptografia = globalThis.crypto;
  if (kryptografia && typeof kryptografia.randomUUID === 'function') {
    return kryptografia.randomUUID();
  }
  return `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 10)}`;
}
