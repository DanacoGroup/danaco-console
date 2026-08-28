/**
 * Nazwy własne sekcji dostępów: wartość wyliczenia kontraktu w zdaniu po
 * polsku. Plik nazywa rodzaj punktu, jego grupę oraz ikonę, stan wraz z klasą
 * stylu i tryb dostępu, opisuje chwilę odczytu i składa adres punktu.
 */
import {
  AccessMode,
  AccessPointKind,
  AccessPointStatus,
} from '../../../shared/contract';

/**
 * Rodzaj punktu dostępu w jednym słowie, gotowym do wstawienia w kartę punktu.
 * Wartość spoza wyliczenia wraca jako własny napis, żeby widoczne było to, co
 * przysłał rdzeń, zamiast pustego miejsca.
 */
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

/**
 * Nagłówek grupy wykazu punktów — nazwa rodzaju w liczbie mnogiej, ponieważ
 * grupa zbiera wszystkie punkty tego samego rodzaju. Rodzaj nieznany klientowi
 * podaje własną wartość kontraktu.
 */
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

/**
 * Ikona grupy dobrana do rodzaju punktu; nazwy pochodzą z zestawu `ikony/`.
 * Most protokołu MCP dostaje znak sieci, katalog lokalny znak katalogu, więc
 * rodzaj grupy czyta się z wykazu bez sięgania do treści wiersza.
 */
export function ikonaRodzaju(rodzaj: AccessPointKind): 'siec' | 'folder' {
  return rodzaj === AccessPointKind.McpBridge ? 'siec' : 'folder';
}

/**
 * Wynik ostatniego sprawdzenia punktu w jednym słowie: odpowiada, nie odpowiada
 * albo niesprawdzony. Stan nierozpoznany wraca własną wartością kontraktu,
 * ponieważ brak wiedzy o stanie nie jest tym samym co stan zły.
 */
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

/**
 * Tryb nadania w jednym słowie: odczyt albo zapis. Tryb spoza wyliczenia wraca
 * własną wartością, żeby nadanie nieznane klientowi było widoczne w karcie,
 * a nie milcząco pominięte.
 */
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

/**
 * Tryby w kolejności rosnącego uprawnienia; ta sama kolejność jest porządkiem
 * przełącznika w karcie nadania, więc ruch w prawo zawsze znaczy uprawnienie
 * szersze, a nie węższe.
 */
export const TRYBY: readonly AccessMode[] = [AccessMode.Read, AccessMode.Write];

/**
 * Rodzaje punktów w kolejności grup wykazu. Wykaz idzie tą kolejnością zamiast
 * kolejnością nadejścia z rdzenia, żeby grupy stały zawsze w tym samym miejscu.
 */
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

/**
 * Adres punktu widoczny w karcie: dla mostu jest nim punkt końcowy, host albo
 * nazwa mostu, a dla katalogu lokalnego identyfikator urządzenia. Pierwsza
 * wartość niepusta wygrywa, ponieważ kontrakt nie wypełnia wszystkich pól.
 */
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

/**
 * Pierwsza wartość niepusta z podanych; komplet pusty daje zdanie zastępcze
 * zamiast napisu pustego, żeby karta punktu nie pokazywała pola bez treści.
 */
function pierwszyNiepusty(wartosci: readonly (string | undefined)[]): string {
  for (const wartosc of wartosci) {
    if (typeof wartosc === 'string' && wartosc !== '') return wartosc;
  }
  return 'adres nieokreślony';
}
