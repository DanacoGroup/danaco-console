import { NavigationKind, type Environment, type Module } from '../../../shared/contract';
import type { NawigacjaPlatformy } from '../protokol/nawigacja-platformy';
import {
  srodowiskoNiedostepne,
  srodowiskoZKontraktu,
  type KluczSrodowiska,
  type Srodowisko,
} from './srodowiska';

/**
 * Źródło wykazu bocznej nawigacji — jedyna prawda powłoki o środowiskach
 * i modułach.
 *
 * Jedna odpowiedzialność: zamiana odpowiedzi rdzenia na wykaz nawigacji.
 * Plik nie buduje ani jednego elementu i nie zna nazwy ani jednej komendy —
 * bierze gotowe opakowania warstwy protokołu.
 *
 * `environment.enter` oddaje środowisko razem z jego modułami, każdy
 * z katalogiem okien operacyjnych, więc jedno wywołanie wystarcza na cały
 * wykaz. `module.list` zawężony do środowiska jest drugim podejściem na
 * wypadek, gdy wejście oddało wykaz pusty mimo nawigacji modułowej.
 *
 * Nieudane wywołanie daje wykaz pusty wraz z treścią odmowy, a nie zaszytą
 * listę modułów — kopia katalogu w kliencie byłaby drugim źródłem prawdy
 * o tym, co rdzeń niesie.
 */
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

/** Środowisko rozpoznane to takie, któremu rdzeń nadał kod — pustka nim nie jest. */
function czyRozpoznane(dane: Environment | undefined): dane is Environment {
  return dane !== undefined && typeof dane.code === 'string' && dane.code.length > 0;
}

/** Panel orkiestracji nie ma modułów z założenia — pusty wykaz nie jest brakiem. */
function czyPanel(dane: Environment): boolean {
  return dane.navigationKind !== NavigationKind.Modules;
}
