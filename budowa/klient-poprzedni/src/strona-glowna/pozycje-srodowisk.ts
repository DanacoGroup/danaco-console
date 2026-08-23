import { KnownModuleIds, type Environment } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/**
 * Cztery środowiska produktu jako pozycje strefy pierwszej strony głównej.
 * Wyłącznie treść kart — bez elementów i bez stylu.
 *
 * Kody pochodzą z kontraktu: `KnownModuleIds` w `shared/contract.ts` odpowiada
 * kolumnie `srodowisko.kod`. Opisy trzymane są w rekordzie kluczowanym tym
 * typem, więc zmiana kodu w `shared/contract.json` przerywa kompilację klienta
 * zamiast dawać cichą lukę.
 */

/** Kod środowiska — typ wywiedziony z kontraktu, nie z literału. */
export type KodSrodowiska = (typeof KnownModuleIds)[number];

export interface PozycjaSrodowiska {
  /** Kod z kontraktu — trafia do zdarzenia wyboru. */
  kod: KodSrodowiska;
  /** Nazwa kanoniczna widoczna na karcie. */
  nazwa: string;
  /**
   * Motto środowiska — zasięg pracy karty: liczba modułów bocznej nawigacji
   * i charakter środowiska, jak na makiecie Centrum dowodzenia.
   */
  motto: string;
  /**
   * Jednozdaniowy opis trybu pracy — zdanie z tabeli opracowania „Kolejność
   * kart i opis trybu pracy na karcie" (rozdz. 3.2).
   */
  opis: string;
  /** Godło środowiska z zestawu ikon. */
  godlo: NazwaIkony;
}

/** Opis pozycji bez kodu — kod dokłada wykaz z kontraktu. */
type TrescSrodowiska = Omit<PozycjaSrodowiska, 'kod'>;

/**
 * Rekord wymusza komplet: brak jednego środowiska albo kod spoza kontraktu
 * jest błędem kompilacji, nie brakiem karty na ekranie.
 */
const TRESCI: Record<KodSrodowiska, TrescSrodowiska> = {
  talkin: {
    nazwa: 'TalkIn',
    motto: '9 modułów · rozmowa i wiedza',
    opis: 'Wiedza, komunikacja i praca z treścią.',
    godlo: 'srodowisko-talkin',
  },
  workspace: {
    nazwa: 'WorkSpace',
    motto: '9 modułów · produktywność',
    opis: 'Produktywność, organizacja i realizacja projektów.',
    godlo: 'srodowisko-workspace',
  },
  codestudio: {
    nazwa: 'CodeStudio',
    motto: '8 modułów · programowanie',
    opis: 'Programowanie — edytor, terminal, kontrola wersji.',
    godlo: 'srodowisko-codestudio',
  },
  multitaskingai: {
    nazwa: 'MultitaskingAI',
    motto: 'panel orkiestracji · 6 sekcji',
    opis: 'Orkiestracja autonomicznej pracy ciągłej.',
    godlo: 'srodowisko-multitaskingai',
  },
};

/**
 * Treść kart widoczna, dopóki rdzeń nie odpowie.
 *
 * Wartość początkowa, nie źródło prawdy: środowiska mieszkają w bazie i oddaje
 * je komenda `environment.list`. Stała zostaje po to, by przed pierwszą
 * odpowiedzią zamiast pustego ekranu stały cztery karty.
 */
export const POZYCJE_SRODOWISK: readonly PozycjaSrodowiska[] = KnownModuleIds.map(
  (kod) => ({ kod, ...TRESCI[kod] }),
);

/** Czy kod środowiska należy do wykazu kontraktu. */
export function kodZnany(kod: string): kod is KodSrodowiska {
  return (KnownModuleIds as readonly string[]).includes(kod);
}

/**
 * Środowisko rdzenia jako pozycja karty.
 *
 * Nazwa, motto i opis idą z rdzenia; godło pochodzi z tablicy miejscowej, bo
 * kolumna `srodowisko` go nie niesie. Środowisko o kodzie spoza kontraktu
 * zostaje na ekranie z godłem zastępczym — wykaz kontraktu jest informacyjny,
 * nie bramą.
 */
export function pozycjaZeSrodowiska(srodowisko: Environment): PozycjaSrodowiska {
  const tresc = kodZnany(srodowisko.code) ? TRESCI[srodowisko.code] : undefined;
  // Kontrakt oznacza te pola jako nieobowiązkowe. Pusty napis i brak pola
  // znaczą tu to samo, więc oba schodzą na treść z tablicy miejscowej.
  const nazwa = srodowisko.name ?? '';
  const opis = srodowisko.description ?? '';
  const motto = srodowisko.motto ?? '';
  return {
    kod: srodowisko.code as KodSrodowiska,
    nazwa: nazwa !== '' ? nazwa : (tresc?.nazwa ?? srodowisko.code),
    motto: motto !== '' ? motto : (tresc?.motto ?? ''),
    opis: opis !== '' ? opis : (tresc?.opis ?? ''),
    godlo: tresc?.godlo ?? 'globus',
  };
}
