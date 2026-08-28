/**
 * Rozstrzygnięcie wartości ustawienia wraz z jej pochodzeniem. Okno
 * konfiguracji pokazuje nie tylko, ile wynosi wartość, ale i z którego poziomu
 * zasięgu oraz z której osi przyszła.
 */
import {
  ConfigAxis,
  ConfigScope,
  type ConfigEntry,
  type SettingDefinition,
} from '../../../shared/contract';
import type { AdresUstawienia } from './adres-ustawienia';
import {
  bytOsiWpisu,
  bytZasieguWpisu,
  osWpisu,
  pierwszenstwoOsi,
  pierwszenstwoZasiegu,
} from './zasiegi';

/**
 * Punkt widzenia to adres, czyli miejsce, względem którego liczone jest
 * dziedziczenie wartości. Wskazuje go pasek u góry okna konfiguracji, ponieważ
 * klient pełnej ścieżki bytów nie zna.
 */
export type PunktWidzenia = AdresUstawienia;

/**
 * Skąd pochodzi wartość obowiązująca: sama wartość, wpis będący jej źródłem,
 * pełny łańcuch zapisów klucza oraz wskaźnik mówiący, czy obowiązuje wartość
 * domyślna katalogu ustawień.
 */
export interface Rozstrzygniecie {
  /** Wartość obowiązująca w punkcie widzenia. */
  wartosc: unknown;
  /** Wpis, z którego wartość pochodzi; brak znaczy wartość domyślną. */
  zrodlo: ConfigEntry | null;
  /** Wszystkie zapisy klucza, od najwęższego; podstawa podglądu dziedziczenia. */
  lancuch: readonly ConfigEntry[];
  /** Czy w punkcie widzenia obowiązuje wartość domyślna katalogu. */
  domyslna: boolean;
}

/**
 * Wpisy dotyczące jednego klucza, uporządkowane od zapisu najwęższego do
 * najszerszego. Porządek jest podstawą zarówno rozstrzygnięcia wartości, jak
 * i podglądu dziedziczenia w oknie konfiguracji.
 */
export function wpisyKlucza(
  wpisy: readonly ConfigEntry[],
  klucz: string,
): ConfigEntry[] {
  return wpisy.filter((wpis) => wpis.key === klucz).sort(porownaj);
}

/**
 * Wartość obowiązująca w punkcie widzenia wraz z jej pochodzeniem.
 *
 * Brak zapisu nie jest błędem ani blokadą: obowiązuje wartość domyślna
 * z katalogu, a wskaźnik mówi to wprost.
 */
export function rozstrzygnij(
  definicja: SettingDefinition,
  wpisy: readonly ConfigEntry[],
  punkt: PunktWidzenia,
): Rozstrzygniecie {
  const lancuch = wpisyKlucza(wpisy, definicja.key);
  const zrodlo = lancuch.find((wpis) => wpisObowiazuje(wpis, punkt)) ?? null;

  return {
    wartosc: zrodlo === null ? definicja.defaultValue : zrodlo.value,
    zrodlo,
    lancuch,
    domyslna: zrodlo === null,
  };
}

/**
 * Czy wpis jest tym, z którego pochodzi wartość obowiązująca. Podgląd
 * dziedziczenia zaznacza nim jeden wiersz łańcucha, żeby Operator widział, który
 * zapis przeważył nad pozostałymi.
 */
export function czyZrodlo(
  rozstrzygniecie: Rozstrzygniecie,
  wpis: ConfigEntry,
): boolean {
  return rozstrzygniecie.zrodlo === wpis;
}

/**
 * Czy wpis wchodzi do łańcucha punktu widzenia. Wchodzą wpisy globalne oraz
 * zapisane dokładnie na wskazanym poziomie i dla wskazanego bytu; zapis dla
 * innego bytu dotyczy kogoś innego.
 */
function wpisObowiazuje(wpis: ConfigEntry, punkt: PunktWidzenia): boolean {
  return zasiegObowiazuje(wpis, punkt) && osObowiazuje(wpis, punkt);
}

function zasiegObowiazuje(wpis: ConfigEntry, punkt: PunktWidzenia): boolean {
  if (wpis.scope === ConfigScope.Global) return true;
  return wpis.scope === punkt.zasieg && bytZasieguWpisu(wpis) === punkt.bytZasiegu;
}

function osObowiazuje(wpis: ConfigEntry, punkt: PunktWidzenia): boolean {
  const os = osWpisu(wpis);
  if (os === ConfigAxis.Platform) return true;
  return os === punkt.os && bytOsiWpisu(wpis) === punkt.bytOsi;
}

/**
 * Porządek wpisów w łańcuchu: najpierw najwęższy poziom zasięgu, a przy równym
 * poziomie najwęższa oś. Dzięki temu wpis stojący na początku łańcucha jest
 * zawsze tym, z którego pochodzi wartość obowiązująca.
 */
function porownaj(pierwszy: ConfigEntry, drugi: ConfigEntry): number {
  const poziom =
    pierwszenstwoZasiegu(pierwszy.scope) - pierwszenstwoZasiegu(drugi.scope);
  if (poziom !== 0) return poziom;

  const os = pierwszenstwoOsi(osWpisu(pierwszy)) - pierwszenstwoOsi(osWpisu(drugi));
  if (os !== 0) return os;

  return bytZasieguWpisu(pierwszy).localeCompare(bytZasieguWpisu(drugi), 'pl-PL');
}
