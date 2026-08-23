import { liczba, obiekt, tekst } from './odczyt-fragmentu';

/**
 * Metadane konta użytego przez kanał.
 *
 * Kanał nadaje fragment tego rodzaju także w chwili przełączenia konta
 * w trakcie tury, żeby rotacja była widoczna. Struktura niesie wyłącznie
 * odwołania: kod konta, nazwę profilu, odwołanie do danych dostępowych. Treść
 * poświadczenia nie trafia tu nigdy.
 *
 * Kształt odpowiada `models.MetadaneKonta` po stronie rdzenia.
 */
export interface MetadaneKonta {
  /** Kod konta w rejestrze platformy. */
  konto: string;
  /** Nazwa profilu poświadczeń kanału. */
  profil: string;
  /** Odwołanie do danych dostępowych — nigdy ich treść. */
  odwolanie: string;
  /** Dlaczego kanał używa właśnie tego konta. */
  powod: string;
  /** Ile kont pozostaje do przełączenia. */
  kolejne: number;
}

/** Odczytuje metadane konta z ładunku fragmentu; ładunek nieczytelny daje `null`. */
export function odczytajKonto(dane: unknown): MetadaneKonta | null {
  const zrodlo = obiekt(dane);
  if (zrodlo === null) return null;
  return {
    konto: tekst(zrodlo, 'account'),
    profil: tekst(zrodlo, 'profile'),
    odwolanie: tekst(zrodlo, 'credentialRef'),
    powod: tekst(zrodlo, 'reason'),
    kolejne: liczba(zrodlo, 'remaining'),
  };
}

/** Opis konta jedną linią — do podtytułu bloku przejrzystości. */
export function opisKonta(k: MetadaneKonta): string {
  const czesci = [k.konto, k.profil].filter((czesc) => czesc.length > 0);
  if (k.powod.length > 0) czesci.push(k.powod);
  return czesci.join(' · ');
}
