import { liczba, obiekt, tekst } from './odczyt-fragmentu';

/** Interfejs opisuje metadane konta użytego przez kanał, niosące wyłącznie odwołania do konta i profilu, nigdy treść poświadczenia. */
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

/** Funkcja odczytuje metadane konta z ładunku fragmentu, zwracając wartość pustą, gdy ładunek jest nieczytelny. */
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

/** Funkcja składa opis konta w jednej linii tekstu, wyświetlany w podtytule bloku przejrzystości kanału komunikacji. */
export function opisKonta(k: MetadaneKonta): string {
  const czesci = [k.konto, k.profil].filter((czesc) => czesc.length > 0);
  if (k.powod.length > 0) czesci.push(k.powod);
  return czesci.join(' · ');
}
