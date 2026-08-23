/**
 * Liczniki dokumentu — rachunek na treści, bez DOM i bez rdzenia.
 *
 * Opracowanie stawia liczniki w warstwie pierwszej modułu (pasek statusu:
 * „słów: 1 284"), a rdzeń żadnej komendy liczącej nie niesie i nie musi:
 * treść stoi w buforze edytora, więc liczenie jej po stronie klienta nie jest
 * obejściem braku, tylko właściwym miejscem tej czynności. Wywołanie rdzenia po
 * liczbę słów byłoby przesyłaniem dokumentu po odpowiedź, którą klient ma
 * natychmiast.
 *
 * Plik nie zna okna: wejściem jest napis, wyjściem liczby. Dzięki temu rachunek
 * sprawdza się bez stawiania widoku.
 *
 * Czas czytania liczony jest tempem **200 słów na minutę** — wartością przyjętą
 * w typografii użytkowej dla tekstu ciągłego. Stała jest nazwana, bo bez nazwy
 * byłaby liczbą magiczną, a przy zmianie wymagania trzeba by jej szukać.
 */

/** Tempo czytania tekstu ciągłego przyjęte do szacunku czasu, w słowach na minutę. */
const SLOW_NA_MINUTE = 200;

/** Komplet liczb opisujących treść dokumentu. */
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

/** Liczy komplet miar dla treści dokumentu. */
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

/**
 * Zdania rozpoznawane po znaku kończącym, a nie po każdej kropce.
 *
 * Kropka rozdzielająca skrót („np.", „itd.") kończy zdanie tylko wtedy, gdy po
 * niej idzie odstęp i wielka litera albo koniec treści. Rachunek jest
 * przybliżony i takim ma pozostać: pełna segmentacja zdań wymaga słownika
 * skrótów, którego moduł nie ma i którego dla licznika w pasku statusu nie
 * warto zakładać.
 */
function policzZdania(tresc: string): number {
  const dopasowania = tresc.match(/[.!?…]+(?=\s+\p{Lu}|\s*$)/gu);
  if (dopasowania !== null) return dopasowania.length;
  // Treść bez znaku kończącego jest jednym zdaniem, o ile w ogóle coś niesie.
  return tresc.trim() === '' ? 0 : 1;
}

/** Akapity rozdzielone pustym wierszem; treść bez pustych wierszy to jeden akapit. */
function policzAkapity(tresc: string): number {
  return tresc
    .split(/\n\s*\n/u)
    .filter((akapit) => akapit.trim() !== '').length;
}

/** Zdanie licznikowe paska statusu — jedno miejsce składania tych liczb w napis. */
export function opiszLiczniki(liczniki: LicznikiDokumentu): string {
  if (liczniki.znaki === 0) return 'dokument pusty';
  return (
    `słów ${liczniki.slowa} · znaków ${liczniki.znaki} (bez odstępów ${liczniki.znakiBezOdstepow}) · ` +
    `zdań ${liczniki.zdania} · akapitów ${liczniki.akapity} · czytanie ~${liczniki.minutyCzytania} min`
  );
}
