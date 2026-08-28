import { NavigationKind, type Environment, type Module } from '../../../shared/contract';
import type { NawigacjaPlatformy } from '../protokol/nawigacja-platformy';
import {
  srodowiskoNiedostepne,
  srodowiskoZKontraktu,
  type KluczSrodowiska,
  type Srodowisko,
} from './srodowiska';

/** Źródło wykazu bocznej nawigacji — jedyna prawda powłoki o środowiskach i modułach tego samego klienta. */
export type ZrodloNawigacji = (klucz: KluczSrodowiska) => Promise<Srodowisko>;

export function utworzZrodloNawigacji(
  platforma: NawigacjaPlatformy,
  idKlienta: string,
): ZrodloNawigacji {
  /** Moduły środowiska drugim podejściem, gdy wejście oddało wykaz pusty. */
  async function doczytajModuly(kod: string): Promise<readonly Module[]> {
    const wynik = await platforma.wykazModulow({ environmentId: kod });
    return wynik.udany ? (wynik.wynik?.modules ?? []) : [];
  }

  return async (klucz) => {
    const wynik = await platforma.wejdzDoSrodowiska({
      environmentId: klucz,
      clientId: idKlienta,
    });
    if (!wynik.udany) {
      return srodowiskoNiedostepne(klucz, wynik.blad?.message ?? 'Rdzeń odmówił wykazu modułów.');
    }

    const dane = wynik.wynik?.environment;
    if (!czyRozpoznane(dane)) {
      return srodowiskoNiedostepne(klucz, `Rdzeń nie zna środowiska o kodzie „${klucz}".`);
    }

    let moduly: readonly Module[] = wynik.wynik?.modules ?? [];
    if (moduly.length === 0 && !czyPanel(dane)) moduly = await doczytajModuly(dane.code);
    return srodowiskoZKontraktu(dane, moduly);
  };
}

/** Środowisko rozpoznane to takie, któremu rdzeń nadał kod — pustka takim środowiskiem nigdy tu nie jest. */
function czyRozpoznane(dane: Environment | undefined): dane is Environment {
  return dane !== undefined && typeof dane.code === 'string' && dane.code.length > 0;
}

/** Panel orkiestracji nie ma modułów z założenia — pusty wykaz w tym wypadku wcale nie jest ich brakiem. */
function czyPanel(dane: Environment): boolean {
  return dane.navigationKind !== NavigationKind.Modules;
}
