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
 * Rozstrzygnięcie wartości ustawienia wraz z jej pochodzeniem.
 *
 * Okno konfiguracji pokazuje nie tylko, ile wynosi wartość, ale i skąd pochodzi:
 * z którego poziomu zasięgu i z której osi przyszła oraz jakie inne zapisy tego
 * samego klucza istnieją w systemie.
 *
 * Rachunek jest w całości po stronie klienta i opiera się wyłącznie na tym,
 * co przyszło z rdzenia komendą `config.get`: rdzeń oddaje wpisy wraz
 * z poziomem, bytem poziomu, osią i bytem osi, więc dziedziczenie da się
 * odtworzyć bez drugiej komendy.
 *
 * Klient nie zna pełnej ścieżki bytów (środowisko → moduł → … → okno), więc jej
 * nie zgaduje: punkt widzenia wskazuje pasek u góry okna. Łańcuch takiego
 * punktu składa się z wpisów globalnych
 * oraz z wpisów zapisanych dokładnie na wskazanym poziomie i dla wskazanego
 * bytu; zapis na tym samym poziomie, lecz dla innego bytu, dotyczy kogoś
 * innego i do rachunku nie wchodzi.
 */

/** Punkt widzenia to adres — miejsce, względem którego liczymy dziedziczenie. */
export type PunktWidzenia = AdresUstawienia;

/** Skąd pochodzi wartość obowiązująca. */
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

/** Wpisy dotyczące jednego klucza, uporządkowane od najwęższego zapisu. */
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

/** Czy wpis jest tym, z którego pochodzi wartość obowiązująca. */
export function czyZrodlo(
  rozstrzygniecie: Rozstrzygniecie,
  wpis: ConfigEntry,
): boolean {
  return rozstrzygniecie.zrodlo === wpis;
}

/** Czy wpis wchodzi do łańcucha punktu widzenia. */
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

/** Porządek wpisów: najpierw najwęższy poziom, potem najwęższa oś. */
function porownaj(pierwszy: ConfigEntry, drugi: ConfigEntry): number {
  const poziom =
    pierwszenstwoZasiegu(pierwszy.scope) - pierwszenstwoZasiegu(drugi.scope);
  if (poziom !== 0) return poziom;

  const os = pierwszenstwoOsi(osWpisu(pierwszy)) - pierwszenstwoOsi(osWpisu(drugi));
  if (os !== 0) return os;

  return bytZasieguWpisu(pierwszy).localeCompare(bytZasieguWpisu(drugi), 'pl-PL');
}
