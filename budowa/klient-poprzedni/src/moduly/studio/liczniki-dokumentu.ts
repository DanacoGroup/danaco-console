/** Liczniki dokumentu, rachunek na treści bez elementów strony i bez rdzenia: wejściem jest napis, wyjściem liczby słów, znaków, zdań i akapitów. */

/** Tempo czytania tekstu ciągłego przyjęte do szacunku czasu potrzebnego na przeczytanie treści, w słowach na minutę. */
const SLOW_NA_MINUTE = 200;

/** Komplet liczb opisujących treść dokumentu: znaki, znaki bez odstępów, słowa, zdania, akapity i szacowany czas czytania. */
export interface LicznikiDokumentu {
  znaki: number;
  /** Znaki bez odstępów — miara objętości używana w rozliczeniach redakcyjnych. */
  znakiBezOdstepow: number;
  slowa: number;
  zdania: number;
  akapity: number;
  /** Szacowany czas czytania w pełnych minutach; najmniej 1 dla treści niepustej. */
  minutyCzytania: number;
}

/** Liczy komplet miar liczbowych opisujących treść dokumentu na podstawie samego napisu, bez odwołania do rdzenia. */
export function policzTresc(tresc: string): LicznikiDokumentu {
  const slowa = tresc.split(/\s+/u).filter((slowo) => slowo !== '').length;
  return {
    znaki: tresc.length,
    znakiBezOdstepow: tresc.replace(/\s/gu, '').length,
    slowa,
    zdania: policzZdania(tresc),
    akapity: policzAkapity(tresc),
    minutyCzytania: slowa === 0 ? 0 : Math.max(1, Math.ceil(slowa / SLOW_NA_MINUTE)),
  };
}

/** Zdania rozpoznawane po znaku kończącym zdanie, a nie po każdej napotkanej kropce, licząc przybliżoną liczbę zdań treści. */
function policzZdania(tresc: string): number {
  const dopasowania = tresc.match(/[.!?…]+(?=\s+\p{Lu}|\s*$)/gu);
  if (dopasowania !== null) return dopasowania.length;
  // Treść bez znaku kończącego jest jednym zdaniem, o ile w ogóle coś niesie.
  return tresc.trim() === '' ? 0 : 1;
}

/** Akapity rozdzielone co najmniej jednym pustym wierszem; treść bez pustych wierszy liczy się jako jeden akapit. */
function policzAkapity(tresc: string): number {
  return tresc
    .split(/\n\s*\n/u)
    .filter((akapit) => akapit.trim() !== '').length;
}

/** Zdanie licznikowe paska statusu, jedyne miejsce składania liczb dokumentu w czytelny napis dla operatora. */
export function opiszLiczniki(liczniki: LicznikiDokumentu): string {
  if (liczniki.znaki === 0) return 'dokument pusty';
  return (
    `słów ${liczniki.slowa} · znaków ${liczniki.znaki} (bez odstępów ${liczniki.znakiBezOdstepow}) · ` +
    `zdań ${liczniki.zdania} · akapitów ${liczniki.akapity} · czytanie ~${liczniki.minutyCzytania} min`
  );
}
