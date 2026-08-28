import { Command, ConfigScope } from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { OPERACJE_PASKA } from './kategorie-operacji';

/** Stała KLUCZ_UZYCIA jest kluczem nastawy liczników użycia operacji przybornika, zapisanym w rdzeniu na poziomie globalnym. */
export const KLUCZ_UZYCIA = 'studio_przybornik_uzycie_operacji';
/** Stała KLUCZ_PRZYPIETYCH jest kluczem nastawy czynności przypiętych przez Operatora, zapisanym w rdzeniu na poziomie globalnym. */
export const KLUCZ_PRZYPIETYCH = 'studio_przybornik_operacje_przypiete';
/** Stała KLUCZ_TRYBU jest kluczem nastawy trybu wykazu operacji przybornika, zapisanym w rdzeniu na poziomie globalnym. */
export const KLUCZ_TRYBU = 'studio_przybornik_tryb_operacji';

/** Stała TrybOperacji nazywa dwa równorzędne tryby wykazu operacji przybornika: narzędzia ukryte i stały panel; wybór należy do Operatora. */
export const TrybOperacji = {
  /** Narzędzia ukryte: pływak przy zaznaczeniu i menu pod uchwytem. */
  Ukryte: 'ukryte',
  /** Stały panel boczny zajmujący kolumnę. */
  Panel: 'panel',
} as const;
export type TrybOperacji = (typeof TrybOperacji)[keyof typeof TrybOperacji];

/** Interfejs UzycieOperacji niesie jedno użycie operacji: identyfikator akcji, ile razy uruchomiona i kiedy ostatnio. */
export interface UzycieOperacji {
  idAkcji: string;
  razy: number;
  ostatnio: number;
}

/** Interfejs NastawyPrzybornika niesie nastawy przybornika trzymane w rdzeniu: użycie operacji, czynności przypięte i tryb wykazu. */
export interface NastawyPrzybornika {
  uzycie: readonly UzycieOperacji[];
  przypiete: readonly string[];
  tryb: TrybOperacji;
}

/** Interfejs UzycieZrodlo opisuje drogę do nastaw przybornika: odczyt trzech nastaw i zapis jednej nastawy w rdzeniu. */
export interface UzycieZrodlo {
  /** Odczytuje trzy nastawy przybornika z poziomu globalnego. */
  przybornikNastawy(): Promise<Wynik<NastawyPrzybornika>>;
  /** Zapisuje jedną nastawę przybornika. */
  przybornikZapisz(klucz: string, wartosc: unknown): Promise<Wynik<{ zapisane: boolean }>>;
}

export function utworzUzycieZrodlo(kanal: Kanal): UzycieZrodlo {
  return {
    async przybornikNastawy() {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigGet, { scope: ConfigScope.Global }),
        Command.ConfigGet,
        (tresc) => czyTablica(tresc.entries),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
      }
      const wpisy = new Map(wynik.wynik.entries.map((wpis) => [wpis.key, wpis.value]));
      return {
        udany: true,
        wynik: {
          uzycie: przybornikOdczytajUzycie(wpisy.get(KLUCZ_UZYCIA)),
          przypiete: przybornikOdczytajPrzypiete(wpisy.get(KLUCZ_PRZYPIETYCH)),
          tryb: przybornikOdczytajTryb(wpisy.get(KLUCZ_TRYBU)),
        },
      };
    },

    async przybornikZapisz(klucz, wartosc) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigSet, {
          key: klucz,
          value: wartosc,
          scope: ConfigScope.Global,
        }),
        Command.ConfigSet,
        (tresc) => czyObiekt(tresc.entry),
      );
      if (!wynik.udany) return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
      return { udany: true, wynik: { zapisane: true } };
    },
  };
}

/**
 * Funkcja przybornikOdczytajUzycie odczytuje liczniki użycia operacji z wartości nastawy o nieznanym kształcie; wpis niezgodny jest pomijany, a nie naprawiany.
 */
export function przybornikOdczytajUzycie(wartosc: unknown): UzycieOperacji[] {
  const zrodlo = przybornikRozpakuj(wartosc);
  if (!Array.isArray(zrodlo)) return [];
  const wynik: UzycieOperacji[] = [];
  for (const wpis of zrodlo) {
    if (typeof wpis !== 'object' || wpis === null) continue;
    const pola = wpis as Record<string, unknown>;
    if (typeof pola['idAkcji'] !== 'string' || typeof pola['razy'] !== 'number') continue;
    wynik.push({
      idAkcji: pola['idAkcji'],
      razy: pola['razy'],
      ostatnio: typeof pola['ostatnio'] === 'number' ? pola['ostatnio'] : 0,
    });
  }
  return wynik;
}

/** Funkcja przybornikOdczytajPrzypiete odczytuje czynności przypięte z wartości nastawy; wpis nietekstowy schodzi z wykazu. */
export function przybornikOdczytajPrzypiete(wartosc: unknown): string[] {
  const zrodlo = przybornikRozpakuj(wartosc);
  if (!Array.isArray(zrodlo)) return [];
  return zrodlo.filter((wpis): wpis is string => typeof wpis === 'string');
}

/**
 * Funkcja przybornikOdczytajTryb odczytuje tryb wykazu operacji; brak nastawy znaczy narzędzia ukryte, bo stały panel jest trybem do wyboru, nie postacią domyślną.
 */
export function przybornikOdczytajTryb(wartosc: unknown): TrybOperacji {
  const zrodlo = przybornikRozpakuj(wartosc);
  return zrodlo === TrybOperacji.Panel ? TrybOperacji.Panel : TrybOperacji.Ukryte;
}

/**
 * Rozpakowuje wartość nastawy.
 *
 * Rdzeń zapisuje wartość jako tekst zapisu JSON i tak ją oddaje, gdy rodzaj
 * wartości nie jest prostym typem. Rozpakowanie stoi w jednym miejscu, żeby trzy
 * nastawy czytały ją tą samą drogą.
 */
function przybornikRozpakuj(wartosc: unknown): unknown {
  if (typeof wartosc !== 'string') return wartosc;
  try {
    return JSON.parse(wartosc);
  } catch {
    return wartosc;
  }
}

/** Interfejs UzyciePrzybornika niesie zbiór użycia operacji wraz z czynnościami na nim: liczeniem, przypinaniem i odczytem kolejności. */
export interface UzyciePrzybornika {
  /** Podnosi licznik operacji i oddaje nowy stan wykazu do zapisu. */
  policz(idAkcji: string): readonly UzycieOperacji[];
  /** Przypina operację albo zdejmuje przypięcie; oddaje nowy wykaz przypiętych. */
  przestawPrzypiecie(idAkcji: string): readonly string[];
  /** Czy operacja jest przypięta. */
  przypieta(idAkcji: string): boolean;
  /** Wstawia nastawy odczytane z rdzenia. */
  wchlon(nastawy: NastawyPrzybornika): void;
  /** Tryb wykazu operacji. */
  tryb(): TrybOperacji;
  ustawTryb(tryb: TrybOperacji): void;
  /** Kolejność pływaka: przypięte, potem najczęstsze, potem ostatnio użyte. */
  naWierzchu(ile: number): readonly string[];
  /** Zdanie o podstawie kolejności — Operator ma wiedzieć, skąd ona jest. */
  opiszKolejnosc(): string;
}

export function utworzUzyciePrzybornika(): UzyciePrzybornika {
  let uzycie: UzycieOperacji[] = [];
  let przypiete: string[] = [];
  let tryb: TrybOperacji = TrybOperacji.Ukryte;

  return {
    policz(idAkcji) {
      const wpis = uzycie.find((pozycja) => pozycja.idAkcji === idAkcji);
      if (wpis === undefined) uzycie.push({ idAkcji, razy: 1, ostatnio: Date.now() });
      else {
        wpis.razy += 1;
        wpis.ostatnio = Date.now();
      }
      return uzycie;
    },

    przestawPrzypiecie(idAkcji) {
      const gdzie = przypiete.indexOf(idAkcji);
      if (gdzie >= 0) przypiete.splice(gdzie, 1);
      else przypiete.push(idAkcji);
      return przypiete;
    },

    przypieta: (idAkcji) => przypiete.includes(idAkcji),

    wchlon(nastawy) {
      uzycie = [...nastawy.uzycie];
      przypiete = [...nastawy.przypiete];
      tryb = nastawy.tryb;
    },

    tryb: () => tryb,
    ustawTryb(nowy) {
      tryb = nowy;
    },

    naWierzchu: (ile) => przybornikNaWierzchu(uzycie, przypiete, ile),

    opiszKolejnosc() {
      if (uzycie.length === 0) {
        return (
          'Kolejność czynności na wierzchu jest jeszcze stanem początkowym: rdzeń nie ma ani ' +
          'jednego policzonego użycia. Po pierwszym uruchomieniu operacji kolejność liczy się ' +
          'z użycia — liczby użyć i czasu ostatniego — a nie z wykazu wpisanego w kod.'
        );
      }
      const razem = uzycie.reduce((suma, wpis) => suma + wpis.razy, 0);
      return (
        `Kolejność policzona z ${razem} uruchomień w ${uzycie.length} czynnościach; przypiętych ` +
        `${przypiete.length}. Przypięte stoją przed najczęstszymi, najczęstsze przed ostatnio ` +
        'używanymi. Nastawa leży w rdzeniu pod kluczem ' +
        `${KLUCZ_UZYCIA} i cofa się komendą config.reset.`
      );
    },
  };
}

/**
 * Funkcja przybornikNaWierzchu zwraca kolejność czynności na wierzchu pływaka: przypięte Operatora, potem najczęściej używane, a przy remisie użyte ostatnio.
 */
export function przybornikNaWierzchu(
  uzycie: readonly UzycieOperacji[],
  przypiete: readonly string[],
  ile: number,
): string[] {
  const kolejnosc: string[] = [...przypiete];
  const wedlugCzestosci = [...uzycie].sort((pierwsza, druga) =>
    druga.razy !== pierwsza.razy ? druga.razy - pierwsza.razy : druga.ostatnio - pierwsza.ostatnio,
  );
  for (const wpis of wedlugCzestosci) {
    if (!kolejnosc.includes(wpis.idAkcji)) kolejnosc.push(wpis.idAkcji);
  }
  for (const operacja of OPERACJE_PASKA) {
    if (!kolejnosc.includes(operacja.id)) kolejnosc.push(operacja.id);
  }
  return kolejnosc.slice(0, Math.max(1, ile));
}
