import { ConfigAxis, ConfigScope, type ConfigEntry } from '../../../shared/contract';

/**
 * Słownik dwóch prostopadłych wymiarów konfiguracji.
 *
 * Poziom zasięgu mówi, jak wąsko obowiązuje wartość: okno jest najwęższe
 * i wygrywa, globalny najszerszy i przegrywa z każdym innym.
 * Oś mówi, dla czego wartość obowiązuje: dla platformy, dla wskazanego
 * modelu albo dla wskazanego konta. Oś pominięta znaczy `platform`.
 */

/** Osiem poziomów zasięgu w kolejności od najwęższego. Okno wygrywa. */
export const ZASIEGI_OD_NAJWEZSZEGO: readonly ConfigScope[] = [
  ConfigScope.Window,
  ConfigScope.Role,
  ConfigScope.Session,
  ConfigScope.Project,
  ConfigScope.ModulePair,
  ConfigScope.Module,
  ConfigScope.Environment,
  ConfigScope.Global,
];

/** Trzy osie w kolejności od najwęższej. Konto wygrywa z modelem, model z platformą. */
export const OSIE_OD_NAJWEZSZEJ: readonly ConfigAxis[] = [
  ConfigAxis.Account,
  ConfigAxis.Model,
  ConfigAxis.Platform,
];

/** Nazwa poziomu zasięgu pokazywana Operatorowi. */
export const NAZWY_ZASIEGOW: Readonly<Record<ConfigScope, string>> = {
  // Zasięg aplikacji obejmuje całą aplikację — nie pojedyncze okno, sesję ani środowisko.
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]: 'globalny',
  [ConfigScope.Environment]: 'środowisko',
  [ConfigScope.Module]: 'moduł',
  [ConfigScope.ModulePair]: 'para modułów',
  [ConfigScope.Project]: 'projekt',
  [ConfigScope.Session]: 'karta sesji',
  [ConfigScope.Role]: 'rola',
  [ConfigScope.Window]: 'okno komunikacji',
};

/** Nazwa osi pokazywana Operatorowi. */
export const NAZWY_OSI: Readonly<Record<ConfigAxis, string>> = {
  [ConfigAxis.Platform]: 'platforma',
  [ConfigAxis.Model]: 'model',
  [ConfigAxis.Account]: 'konto',
};

/** Nazwa bytu, którego identyfikator podaje się przy poziomie zasięgu. */
export const BYTY_ZASIEGOW: Readonly<Record<ConfigScope, string>> = {
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]: '',
  [ConfigScope.Environment]: 'kod środowiska',
  [ConfigScope.Module]: 'kod modułu',
  [ConfigScope.ModulePair]: 'para kodów modułów',
  [ConfigScope.Project]: 'identyfikator projektu',
  [ConfigScope.Session]: 'identyfikator sesji',
  [ConfigScope.Role]: 'nazwa roli',
  [ConfigScope.Window]: 'identyfikator okna',
};

/** Nazwa bytu osi; oś platformy bytu nie ma. */
export const BYTY_OSI: Readonly<Record<ConfigAxis, string>> = {
  [ConfigAxis.Platform]: '',
  [ConfigAxis.Model]: 'identyfikator modelu',
  [ConfigAxis.Account]: 'identyfikator konta',
};

/**
 * Pierwszeństwo poziomu: im mniejsza liczba, tym węziej i tym mocniej.
 * Poziom nierozpoznany trafia za globalny, zamiast wywrócić porównanie.
 */
export function pierwszenstwoZasiegu(zasieg: string): number {
  const miejsce = (ZASIEGI_OD_NAJWEZSZEGO as readonly string[]).indexOf(zasieg);
  return miejsce === -1 ? ZASIEGI_OD_NAJWEZSZEGO.length : miejsce;
}

/** Pierwszeństwo osi liczone tak samo jak pierwszeństwo poziomu. */
export function pierwszenstwoOsi(os: string): number {
  const miejsce = (OSIE_OD_NAJWEZSZEJ as readonly string[]).indexOf(os);
  return miejsce === -1 ? OSIE_OD_NAJWEZSZEJ.length : miejsce;
}

/** Oś wpisu konfiguracji; pole puste znaczy `platform`. */
export function osWpisu(wpis: ConfigEntry): ConfigAxis {
  return wpis.axis ?? ConfigAxis.Platform;
}

/** Byt poziomu zapisany we wpisie; pole puste znaczy napis pusty. */
export function bytZasieguWpisu(wpis: ConfigEntry): string {
  return wpis.scopeId ?? '';
}

/** Byt osi zapisany we wpisie; pole puste znaczy napis pusty. */
export function bytOsiWpisu(wpis: ConfigEntry): string {
  return wpis.axisId ?? '';
}

/** Nazwa poziomu gotowa do wydruku; poziom spoza kontraktu pokazuje własny kod. */
export function nazwaZasiegu(zasieg: string): string {
  return NAZWY_ZASIEGOW[zasieg as ConfigScope] ?? zasieg;
}

/** Nazwa osi gotowa do wydruku; oś spoza kontraktu pokazuje własny kod. */
export function nazwaOsi(os: string): string {
  return NAZWY_OSI[os as ConfigAxis] ?? os;
}

/** Czy poziom wymaga wskazania bytu. Globalny jako jedyny nie wymaga. */
export function zasiegWymagaBytu(zasieg: ConfigScope): boolean {
  return zasieg !== ConfigScope.Global;
}

/** Czy oś wymaga wskazania bytu. Platforma jako jedyna nie wymaga. */
export function osWymagaBytu(os: ConfigAxis): boolean {
  return os !== ConfigAxis.Platform;
}
