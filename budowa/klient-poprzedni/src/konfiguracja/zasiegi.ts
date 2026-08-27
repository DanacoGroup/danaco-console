/**
 * Słownik dwóch prostopadłych wymiarów konfiguracji: poziomu zasięgu oraz osi.
 * Plik ustawia oba wymiary w porządku pierwszeństwa, nazywa je po polsku wraz
 * z bytem, który wskazują, i odczytuje jedno oraz drugie z wpisu konfiguracji.
 */
import { ConfigAxis, ConfigScope, type ConfigEntry } from '../../../shared/contract';

/**
 * Osiem poziomów zasięgu w kolejności od najwęższego do najszerszego. Kolejność
 * jest zarazem porządkiem pierwszeństwa: okno wygrywa z każdym innym poziomem,
 * a globalny przegrywa z każdym.
 */
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

/**
 * Trzy osie w kolejności od najwęższej do najszerszej. Konto wygrywa z modelem,
 * a model z platformą, więc wartość nadana kontu obowiązuje mimo wartości
 * nadanej modelowi albo całej platformie.
 */
export const OSIE_OD_NAJWEZSZEJ: readonly ConfigAxis[] = [
  ConfigAxis.Account,
  ConfigAxis.Model,
  ConfigAxis.Platform,
];

/**
 * Nazwa poziomu zasięgu pokazywana Operatorowi. Wykaz obejmuje wszystkie poziomy
 * kontraktu, także zasięg aplikacji, który nie wchodzi do porządku pierwszeństwa,
 * lecz bywa nazwany w oknie konfiguracji.
 */
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

/**
 * Nazwa osi pokazywana Operatorowi. Zdanie po polsku stoi wyłącznie tutaj, żeby
 * wszystkie okna konfiguracji nazywały tę samą oś tak samo.
 */
export const NAZWY_OSI: Readonly<Record<ConfigAxis, string>> = {
  [ConfigAxis.Platform]: 'platforma',
  [ConfigAxis.Model]: 'model',
  [ConfigAxis.Account]: 'konto',
};

/**
 * Nazwa bytu, którego identyfikator podaje się przy poziomie zasięgu. Napis pusty
 * znaczy poziom bez bytu, czyli poziom globalny; pozostałe poziomy nazywają wprost,
 * czego oczekują od Operatora w polu identyfikatora.
 */
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

/**
 * Nazwa bytu osi. Oś platformy bytu nie ma i stoi z napisem pustym, natomiast oś
 * modelu oraz oś konta wymagają identyfikatora, który nazywa ten wykaz.
 */
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

/**
 * Pierwszeństwo osi liczone tak samo jak pierwszeństwo poziomu: im mniejsza liczba,
 * tym węziej i tym mocniej. Oś nierozpoznana trafia za platformę, zamiast wywrócić
 * porównanie dwóch wpisów.
 */
export function pierwszenstwoOsi(os: string): number {
  const miejsce = (OSIE_OD_NAJWEZSZEJ as readonly string[]).indexOf(os);
  return miejsce === -1 ? OSIE_OD_NAJWEZSZEJ.length : miejsce;
}

/**
 * Oś wpisu konfiguracji. Pole nieobecne we wpisie znaczy oś `platform`, ponieważ
 * kontrakt pomija oś najszerszą, a porównanie wpisów potrzebuje jej nazwanej wprost.
 */
export function osWpisu(wpis: ConfigEntry): ConfigAxis {
  return wpis.axis ?? ConfigAxis.Platform;
}

/**
 * Byt poziomu zapisany we wpisie. Pole nieobecne daje napis pusty, więc wpis
 * poziomu globalnego i wpis z bytem porównuje się tą samą drogą, bez rozgałęzienia
 * po obecności pola.
 */
export function bytZasieguWpisu(wpis: ConfigEntry): string {
  return wpis.scopeId ?? '';
}

/**
 * Byt osi zapisany we wpisie. Pole nieobecne daje napis pusty, więc wpis osi
 * platformy i wpis osi z identyfikatorem porównuje się tą samą drogą, bez
 * rozgałęzienia po obecności pola.
 */
export function bytOsiWpisu(wpis: ConfigEntry): string {
  return wpis.axisId ?? '';
}

/**
 * Nazwa poziomu gotowa do wydruku. Poziom spoza kontraktu pokazuje własny kod
 * zamiast pustego miejsca, żeby widoczne było to, co przysłał rdzeń.
 */
export function nazwaZasiegu(zasieg: string): string {
  return NAZWY_ZASIEGOW[zasieg as ConfigScope] ?? zasieg;
}

/**
 * Nazwa osi gotowa do wydruku. Oś spoza kontraktu pokazuje własny kod zamiast
 * pustego miejsca, żeby widoczne było to, co przysłał rdzeń.
 */
export function nazwaOsi(os: string): string {
  return NAZWY_OSI[os as ConfigAxis] ?? os;
}

/**
 * Czy poziom wymaga wskazania bytu. Globalny jako jedyny nie wymaga, ponieważ
 * obejmuje całość i nie ma czego wskazywać; pozostałe poziomy bez bytu byłyby
 * wpisem bez adresu.
 */
export function zasiegWymagaBytu(zasieg: ConfigScope): boolean {
  return zasieg !== ConfigScope.Global;
}

/**
 * Czy oś wymaga wskazania bytu. Platforma jako jedyna nie wymaga, ponieważ obejmuje
 * całość; oś modelu i oś konta bez identyfikatora byłyby wpisem bez adresu.
 */
export function osWymagaBytu(os: ConfigAxis): boolean {
  return os !== ConfigAxis.Platform;
}
