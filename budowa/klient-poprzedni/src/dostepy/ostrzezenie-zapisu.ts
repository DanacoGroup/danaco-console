import { AccessMode, type AccessPoint } from '../../../shared/contract';

/**
 * Ostrzeżenie przy nadaniu trybu zapisu na maszynie chronionej.
 *
 * Specyfikacja mostu `mcp-danaco-pulpit-console` zabrania nadawania zapisu na
 * `danaco-data` bez wyraźnej potrzeby: to host produkcyjnej platformy LEX,
 * a zapis modelu sięga tam zbiorów, z których korzysta cała kancelaria. Zakazu
 * nie egzekwuje ten plik — nadanie pozostaje możliwe, bo decyzja należy do
 * Operatora. Egzekwuje go jawne, widoczne zdanie przy przełączniku trybu, a nie
 * podpowiedź pod kursorem.
 *
 * Wykaz maszyn stoi tutaj, a nie w katalogu konfiguracji: to ostrzeżenie
 * bezpieczeństwa, więc nie może dać się wyłączyć zapisem w bazie. Rozszerzenie
 * wykazu to jeden wiersz w stałej poniżej.
 */

/** Przedrostek nazw mostów konsoli; odpowiednik `core.PrefiksMostuKonsoli`. */
const PRZEDROSTEK_MOSTU = 'mcp-danaco-pulpit-console-';

/** Maszyny, na których zapis wymaga świadomej decyzji Operatora. */
export const MASZYNY_CHRONIONE: readonly string[] = ['danaco-data'];

/**
 * Nazwa maszyny, do której odnosi się punkt dostępu.
 *
 * Kontrakt niesie ją w trzech polach zależnie od tego, jak punkt założono:
 * `host` przy moście MCP, `bridgeName` przy moście nazwanym po stronie
 * klienta, `deviceId` przy katalogu lokalnym. Bierzemy pierwsze wypełnione,
 * a z nazwy mostu zdejmujemy przedrostek konsoli — po nim zostaje sama nazwa
 * maszyny.
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

/** Czy punkt prowadzi na maszynę objętą ostrzeżeniem. */
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

/** Treść ostrzeżenia dla wskazanej maszyny. */
export function zdanieOstrzezenia(maszyna: string): string {
  return (
    `Tryb zapisu na maszynie ${maszyna} — to host produkcyjnej platformy LEX. ` +
    'Specyfikacja mostu dopuszcza tu zapis wyłącznie przy wyraźnej ' +
    'potrzebie: model będzie mógł zmieniać zbiory, z których korzysta cała ' +
    'kancelaria. Zostaw odczyt, jeżeli zadanie nie wymaga zmiany danych.'
  );
}
