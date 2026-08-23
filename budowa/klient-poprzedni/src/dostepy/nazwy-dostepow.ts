import {
  AccessMode,
  AccessPointKind,
  AccessPointStatus,
} from '../../../shared/contract';

/**
 * Nazwy własne sekcji dostępów — jedyne miejsce, w którym wartość wyliczenia
 * kontraktu zamienia się w zdanie po polsku.
 *
 * Wartości pochodzą wyłącznie ze stałych `shared/contract`; ten plik
 * dokłada do nich warstwę językową i nic więcej. Rozsypanie tych zdań po
 * widokach dałoby dwie nazwy tego samego stanu.
 *
 * Wartość spoza wyliczenia nie jest błędem i niczego nie wygasza: wraca jako
 * własny napis, żeby widoczne było, co przysłał rdzeń, zamiast pustego miejsca.
 */

/** Rodzaj punktu dostępu w jednym słowie. */
export function nazwaRodzaju(rodzaj: AccessPointKind): string {
  switch (rodzaj) {
    case AccessPointKind.McpBridge:
      return 'Maszyna (most MCP)';
    case AccessPointKind.LocalDirectory:
      return 'Katalog lokalny';
    default:
      return String(rodzaj);
  }
}

/** Nagłówek grupy wykazu punktów. */
export function nazwaGrupy(rodzaj: AccessPointKind): string {
  switch (rodzaj) {
    case AccessPointKind.McpBridge:
      return 'Maszyny';
    case AccessPointKind.LocalDirectory:
      return 'Katalogi lokalne';
    default:
      return String(rodzaj);
  }
}

/** Ikona grupy; nazwy pochodzą z zestawu `ikony/`. */
export function ikonaRodzaju(rodzaj: AccessPointKind): 'siec' | 'folder' {
  return rodzaj === AccessPointKind.McpBridge ? 'siec' : 'folder';
}

/** Wynik ostatniego sprawdzenia punktu. */
export function nazwaStanu(stan: AccessPointStatus): string {
  switch (stan) {
    case AccessPointStatus.Reachable:
      return 'odpowiada';
    case AccessPointStatus.Unreachable:
      return 'nie odpowiada';
    case AccessPointStatus.Unknown:
      return 'niesprawdzony';
    default:
      return String(stan);
  }
}

/**
 * Odmiana plakietki stanu. Stan nierozpoznany dostaje plakietkę neutralną —
 * nie jest błędem, jest brakiem wiedzy.
 */
export function klasaStanu(stan: AccessPointStatus): string {
  switch (stan) {
    case AccessPointStatus.Reachable:
      return 'dn-plakietka dn-plakietka--sukces';
    case AccessPointStatus.Unreachable:
      return 'dn-plakietka dn-plakietka--blad';
    default:
      return 'dn-plakietka';
  }
}

/** Tryb nadania w jednym słowie. */
export function nazwaTrybu(tryb: AccessMode): string {
  switch (tryb) {
    case AccessMode.Read:
      return 'odczyt';
    case AccessMode.Write:
      return 'zapis';
    default:
      return String(tryb);
  }
}

/** Tryby w kolejności rosnącego uprawnienia — porządek przełącznika. */
export const TRYBY: readonly AccessMode[] = [AccessMode.Read, AccessMode.Write];

/** Rodzaje w kolejności grup wykazu. */
export const RODZAJE: readonly AccessPointKind[] = [
  AccessPointKind.McpBridge,
  AccessPointKind.LocalDirectory,
];

/**
 * Chwila w zapisie czytelnym. Brak chwili nie daje pustego miejsca, tylko
 * zdanie — Operator ma wiedzieć, że punktu nie sprawdzono.
 */
export function opisChwili(chwila: number | undefined): string {
  if (chwila === undefined || !Number.isFinite(chwila)) return 'nigdy';
  try {
    return new Date(chwila).toLocaleString('pl-PL');
  } catch {
    return String(chwila);
  }
}

/** Adres punktu widoczny w karcie: most albo urządzenie. */
export function adresPunktu(dane: {
  kind: AccessPointKind;
  host?: string;
  endpoint?: string;
  bridgeName?: string;
  deviceId?: string;
}): string {
  if (dane.kind === AccessPointKind.McpBridge) {
    return pierwszyNiepusty([dane.endpoint, dane.host, dane.bridgeName]);
  }
  return pierwszyNiepusty([dane.deviceId]);
}

/**
 * Nazwa punktu wyprowadzona ze ścieżki katalogu: ostatni człon albo cała
 * ścieżka, gdy członów nie ma. Operator wskazuje katalog oknem powłoki i nie
 * ma gdzie podać nazwy — ta musi powstać sama, i musi być czytelna.
 */
export function nazwaZeSciezki(sciezka: string): string {
  const czlony = sciezka.split(/[\\/]+/u).filter((czlon) => czlon !== '');
  const ostatni = czlony[czlony.length - 1];
  return ostatni === undefined || ostatni === '' ? sciezka : ostatni;
}

/** Pierwsza wartość niepusta; komplet pusty daje zdanie zastępcze. */
function pierwszyNiepusty(wartosci: readonly (string | undefined)[]): string {
  for (const wartosc of wartosci) {
    if (typeof wartosc === 'string' && wartosc !== '') return wartosc;
  }
  return 'adres nieokreślony';
}
