import {
  Command,
  ConfigScope,
  type ConfigEntry,
  type SettingDefinition,
} from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { wywolaj } from '../protokol/wywolanie';
import { naZmianeKlucza, zapiszNastawe } from './zrodlo-nastaw';

/**
 * Źródło sekcji „Powiadomienia" — piętnaście nastaw czytanych dwoma pytaniami.
 *
 * Sekcja nie ma własnej rodziny kontraktu i mieć jej nie powinna. Model danych
 * mówi wprost (rozdz. 18.4, uwaga projektowa), że zakres klas zdarzeń i kanał
 * dostarczenia SĄ ustawieniami konfiguracyjnymi; katalog wnosi je migracją 377,
 * a droga do nich jest tą samą drogą, co do każdej innej nastawy platformy:
 * `settings.definition.list`, `config.get`, `config.set`, zdarzenie
 * `config.changed`.
 *
 * Dwa pytania, nie piętnaście. Katalog przychodzi jednym wywołaniem po
 * kategorii, wartości jednym wywołaniem po poziomie — `config.get` bez klucza
 * oddaje komplet zapisów zasięgu. Piętnaście par pytań na otwarcie sekcji byłoby
 * trzydziestoma kopertami po treść, którą rdzeń oddaje w dwóch.
 *
 * Poziom zapisu jest jeden — `globalny`. Katalog dopuszcza trzy (globalny,
 * środowisko, karta sesji), bo tak stoi w dwóch dokumentach dostawy, ale okno
 * Ustawień jest oknem poziomu aplikacji i nie ma w nim ani selektora
 * środowiska, ani karty sesji. Zawężenie środowiskiem należy do okna
 * Konfiguracji, gdzie łańcuch zasięgów jest widoczny i wybieralny.
 */

/** Kod klasy zdarzenia — wprost z `powiadomienie.klasa` modelu danych. */
export type KlasaZdarzenia =
  | 'zakonczenie'
  | 'decyzja'
  | 'blad'
  | 'wzmianka'
  | 'termin'
  | 'automatyka'
  | 'system';

/** Klasy w kolejności prezentacji z rozdz. 7.2 opracowania Ustawień. */
export const KLASY_ZDARZEN: readonly KlasaZdarzenia[] = [
  'zakonczenie',
  'decyzja',
  'blad',
  'wzmianka',
  'termin',
  'automatyka',
  'system',
];

/** Klucz przełącznika głównego. */
export const KLUCZ_GLOWNY = 'powiadomienia.wlaczone';

/** Klucz czynności klasy zdarzenia. */
export function kluczKlasy(klasa: KlasaZdarzenia): string {
  return `powiadomienia.klasa.${klasa}`;
}

/** Klucz kanałów dodatkowych klasy zdarzenia. */
export function kluczKanalow(klasa: KlasaZdarzenia): string {
  return `${kluczKlasy(klasa)}.kanaly`;
}

/** Stan jednej klasy zdarzenia po odczycie. */
export interface StanKlasy {
  klasa: KlasaZdarzenia;
  /** Nazwa i opis z katalogu rdzenia; `undefined` znaczy „katalog jej nie zna". */
  definicja?: SettingDefinition;
  /** Czy klasa jest czynna. */
  czynna: boolean;
  /** Kanały dodatkowe wybrane dla tej klasy; centrum stoi zawsze i nie jest tu wymienione. */
  kanaly: string[];
}

/** Stan całej sekcji po odczycie. */
export interface StanPowiadomien {
  /** Przełącznik główny. */
  wlaczone: boolean;
  klasy: StanKlasy[];
  /** Kanały dodatkowe dopuszczone przez katalog — opcje, nie wybór. */
  dostepneKanaly: { wartosc: string; etykieta: string }[];
  /** Zdanie odmowy rdzenia; puste znaczy „odczyt się udał". */
  odmowa: string;
}

export interface ZrodloPowiadomien {
  odczytaj(): Promise<StanPowiadomien>;
  /** Przestawia przełącznik główny albo czynność jednej klasy. */
  ustawLogiczna(klucz: string, wartosc: boolean): Promise<string>;
  /** Zapisuje komplet kanałów dodatkowych jednej klasy. */
  ustawKanaly(klasa: KlasaZdarzenia, kanaly: string[]): Promise<string>;
  /** Nasłuch zmian dowolnego z piętnastu kluczy — także z drugiego okna. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
}

/** Poziom zapisu sekcji; powód wyboru stoi w nagłówku pliku. */
const POZIOM = ConfigScope.Global;

/** Kategoria katalogu wniesiona migracją 377. */
const KATEGORIA = 'powiadomienia';

export function utworzZrodloPowiadomien(kanal: Kanal): ZrodloPowiadomien {
  /** Wszystkie klucze sekcji — do rozpoznania zdarzenia `config.changed`. */
  const klucze = new Set<string>([
    KLUCZ_GLOWNY,
    ...KLASY_ZDARZEN.map(kluczKlasy),
    ...KLASY_ZDARZEN.map(kluczKanalow),
  ]);

  async function odczytaj(): Promise<StanPowiadomien> {
    const pusty: StanPowiadomien = {
      wlaczone: true,
      klasy: [],
      dostepneKanaly: [],
      odmowa: '',
    };

    const katalog = await wywolaj(kanal, Command.SettingsDefinitionList, {
      categoryId: KATEGORIA,
      includeDisabled: false,
    });
    if (!katalog.udany || katalog.wynik === undefined) {
      return { ...pusty, odmowa: zdanieOdmowy(katalog.blad?.message) };
    }
    const definicje = new Map(katalog.wynik.definitions.map((p) => [p.key, p]));
    if (definicje.size === 0) {
      return {
        ...pusty,
        odmowa:
          'Katalog ustawień rdzenia nie zna kategorii „Powiadomienia" — rdzeń pracuje na ' +
          'bazie sprzed migracji 377 i nie ma czym sterować.',
      };
    }

    const wartosci = await wywolaj(kanal, Command.ConfigGet, { scope: POZIOM });
    if (!wartosci.udany || wartosci.wynik === undefined) {
      return { ...pusty, odmowa: zdanieOdmowy(wartosci.blad?.message) };
    }
    const zapisy = new Map<string, ConfigEntry>(
      wartosci.wynik.entries.map((wpis) => [wpis.key, wpis]),
    );

    return {
      wlaczone: logiczna(zapisy.get(KLUCZ_GLOWNY), definicje.get(KLUCZ_GLOWNY), true),
      klasy: KLASY_ZDARZEN.map((klasa) => {
        const definicja = definicje.get(kluczKlasy(klasa));
        return {
          klasa,
          ...(definicja === undefined ? {} : { definicja }),
          czynna: logiczna(zapisy.get(kluczKlasy(klasa)), definicja, true),
          kanaly: wykazTekstow(
            zapisy.get(kluczKanalow(klasa))?.value ??
              definicje.get(kluczKanalow(klasa))?.defaultValue,
          ),
        };
      }),
      dostepneKanaly: (definicje.get(kluczKanalow('zakonczenie'))?.options ?? []).map((opcja) => ({
        wartosc: opcja.value,
        etykieta: opcja.label,
      })),
      odmowa: '',
    };
  }

  return {
    odczytaj,

    async ustawLogiczna(klucz, wartosc) {
      const odpowiedz = await zapiszNastawe(kanal, klucz, wartosc, POZIOM);
      return odpowiedz.udany ? '' : zdanieOdmowy(odpowiedz.blad?.message);
    },

    async ustawKanaly(klasa, kanaly) {
      const odpowiedz = await zapiszNastawe(kanal, kluczKanalow(klasa), kanaly, POZIOM);
      return odpowiedz.udany ? '' : zdanieOdmowy(odpowiedz.blad?.message);
    },

    naZmiane(sluchacz) {
      // Jeden nasłuch na piętnaście kluczy, nie piętnaście nasłuchów: zdarzenie
      // niesie klucz, więc rozpoznanie jest sprawdzeniem przynależności.
      const odsubskrybuj = [...klucze].map((klucz) =>
        naZmianeKlucza(kanal, klucz, () => sluchacz()),
      );
      return () => odsubskrybuj.forEach((zdejmij) => zdejmij());
    },
  };
}

/**
 * Wartość logiczna nastawy: zapis, a przy jego braku wartość domyślna katalogu.
 *
 * Zapisu nie ma znaczy „obowiązuje domyślna", nie „fałsz" — inaczej sekcja
 * pokazywałaby wszystko wyłączone na świeżej bazie, choć katalog mówi
 * „aktywne".
 */
function logiczna(
  wpis: ConfigEntry | undefined,
  definicja: SettingDefinition | undefined,
  gdyBrak: boolean,
): boolean {
  const surowa = wpis?.value ?? definicja?.defaultValue;
  if (typeof surowa === 'boolean') return surowa;
  if (typeof surowa === 'string') return surowa === 'true';
  return gdyBrak;
}

/** Wykaz napisów z wartości kontraktu; kształt obcy daje wykaz pusty. */
function wykazTekstow(wartosc: unknown): string[] {
  if (Array.isArray(wartosc)) return wartosc.filter((p): p is string => typeof p === 'string');
  if (typeof wartosc === 'string' && wartosc.trim() !== '') {
    try {
      const odczytana: unknown = JSON.parse(wartosc);
      return Array.isArray(odczytana)
        ? odczytana.filter((p): p is string => typeof p === 'string')
        : [];
    } catch {
      return [];
    }
  }
  return [];
}

/** Zdanie odmowy rdzenia albo nazwanie milczenia; nigdy pustka. */
function zdanieOdmowy(wiadomosc: string | undefined): string {
  return wiadomosc === undefined || wiadomosc === ''
    ? 'Rdzeń odmówił bez podania powodu.'
    : wiadomosc;
}
