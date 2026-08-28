import { KnownModuleIds, type Environment } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/** Cztery środowiska produktu jako pozycje strefy pierwszej strony głównej niosą wyłącznie treść kart, a ich kody pochodzą z kontraktu, żeby zmiana w kontrakcie przerywała kompilację klienta. */
export type KodSrodowiska = (typeof KnownModuleIds)[number];

export interface PozycjaSrodowiska {
  /** Kod z kontraktu — trafia do zdarzenia wyboru. */
  kod: KodSrodowiska;
  /** Nazwa kanoniczna widoczna na karcie. */
  nazwa: string;
  /** Motto środowiska niesie zasięg pracy karty: liczbę modułów bocznej nawigacji i charakter środowiska. */
  motto: string;
  /** Jednozdaniowy opis trybu pracy widoczny na karcie środowiska. */
  opis: string;
  /** Godło środowiska z zestawu ikon. */
  godlo: NazwaIkony;
}

/** Opis pozycji bez kodu; kod dokłada wykaz z kontraktu, żeby jeden słownik nie powielał dwóch źródeł prawdy. */
type TrescSrodowiska = Omit<PozycjaSrodowiska, 'kod'>;

/** Rekord wymusza komplet: brak jednego środowiska albo kod spoza kontraktu jest błędem kompilacji, nie brakiem karty. */
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

/** Treść kart widoczna, dopóki rdzeń nie odpowie; wartość początkowa, nie źródło prawdy, bo środowiska mieszkają w bazie. */
export const POZYCJE_SRODOWISK: readonly PozycjaSrodowiska[] = KnownModuleIds.map(
  (kod) => ({ kod, ...TRESCI[kod] }),
);

/** Rozstrzyga, czy podany kod środowiska należy do zamkniętego wykazu kontraktu, zamiast wpaść w kod obcy. */
export function kodZnany(kod: string): kod is KodSrodowiska {
  return (KnownModuleIds as readonly string[]).includes(kod);
}

/** Buduje pozycję karty ze środowiska rdzenia: nazwa, motto i opis idą z rdzenia, a godło pochodzi z tablicy miejscowej. */
export function pozycjaZeSrodowiska(srodowisko: Environment): PozycjaSrodowiska {
  const tresc = kodZnany(srodowisko.code) ? TRESCI[srodowisko.code] : undefined;
  // Kontrakt oznacza te pola jako nieobowiązkowe; pusty napis i brak pola znaczą tu to samo.
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
