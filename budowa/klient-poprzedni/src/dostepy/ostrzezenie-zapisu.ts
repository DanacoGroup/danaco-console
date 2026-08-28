/**
 * Ostrzeżenie przy nadaniu trybu zapisu na maszynie chronionej. Zakazu ten
 * plik nie egzekwuje — nadanie pozostaje możliwe, ponieważ rozstrzyga o nim
 * Operator. Egzekwuje go jawne zdanie widoczne przy przełączniku trybu.
 */
import { AccessMode, type AccessPoint } from '../../../shared/contract';

/**
 * Przedrostek nazw mostów konsoli, odpowiednik stałej `core.PrefiksMostuKonsoli`
 * po stronie rdzenia. Zdjęcie go z nazwy mostu zostawia samą nazwę maszyny.
 */
const PRZEDROSTEK_MOSTU = 'mcp-danaco-pulpit-console-';

/**
 * Maszyny, na których nadanie trybu zapisu wymaga świadomej decyzji Operatora.
 * Wykaz stoi w kodzie, a nie w katalogu konfiguracji, więc nie daje się
 * wyłączyć zapisem w bazie.
 */
export const MASZYNY_CHRONIONE: readonly string[] = ['danaco-data'];

/**
 * Nazwa maszyny, do której odnosi się punkt dostępu. Kontrakt niesie ją
 * w polach `host` przy moście MCP, `bridgeName` przy moście nazwanym po
 * stronie klienta oraz `deviceId` przy katalogu lokalnym; wynikiem jest
 * pierwsze wypełnione, bez przedrostka.
 */
export function maszynaPunktu(punkt: AccessPoint): string {
  if (typeof punkt.host === 'string' && punkt.host !== '') return punkt.host;
  const most = punkt.bridgeName ?? '';
  if (most !== '') {
    return most.startsWith(PRZEDROSTEK_MOSTU)
      ? most.slice(PRZEDROSTEK_MOSTU.length)
      : most;
  }
  return punkt.deviceId ?? '';
}

/**
 * Czy punkt dostępu prowadzi na maszynę objętą ostrzeżeniem. Porównanie pomija
 * wielkość liter i uznaje zgodność także wtedy, gdy nazwa maszyny zawiera
 * nazwę chronioną jako część dłuższego napisu.
 */
export function czyMaszynaChroniona(punkt: AccessPoint): boolean {
  const maszyna = maszynaPunktu(punkt).toLowerCase();
  if (maszyna === '') return false;
  return MASZYNY_CHRONIONE.some(
    (chroniona) => maszyna === chroniona || maszyna.includes(chroniona),
  );
}

/**
 * Zdanie ostrzeżenia albo napis pusty, gdy ostrzeżenia nie ma.
 *
 * Napis pusty znaczy „brak ostrzeżenia" — ta sama umowa co w polu
 * `Kontrolka.ostrzezenie` okna konfiguracji.
 */
export function ostrzezenieZapisu(punkt: AccessPoint, tryb: AccessMode): string {
  if (tryb !== AccessMode.Write || !czyMaszynaChroniona(punkt)) return '';
  return zdanieOstrzezenia(maszynaPunktu(punkt));
}

/**
 * Treść ostrzeżenia dla wskazanej maszyny: nazwa maszyny, rola hosta
 * produkcyjnej platformy LEX, warunek wyraźnej potrzeby oraz zachęta do
 * pozostawienia trybu odczytu.
 */
export function zdanieOstrzezenia(maszyna: string): string {
  return (
    `Tryb zapisu na maszynie ${maszyna} — to host produkcyjnej platformy LEX. ` +
    'Specyfikacja mostu dopuszcza tu zapis wyłącznie przy wyraźnej ' +
    'potrzebie: model będzie mógł zmieniać zbiory, z których korzysta cała ' +
    'kancelaria. Zostaw odczyt, jeżeli zadanie nie wymaga zmiany danych.'
  );
}
