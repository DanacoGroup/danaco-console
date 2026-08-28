import { invoke } from '@tauri-apps/api/core';

import { czyPowlokaNatywna } from './powloka-natywna';

/** Nazwa polecenia powłoki wywoływanego przy odczycie adresu rdzenia lokalnego; odpowiednik polecenia adres_rdzenia. */
const POLECENIE_ADRESU = 'adres_rdzenia';

/** Nazwa polecenia powłoki wywoływanego przy odczycie stanu rdzenia postawionego w tle; odpowiednik polecenia stan_rdzenia. */
const POLECENIE_STANU = 'stan_rdzenia';

/**
 * Stan rdzenia widziany przez powłokę. Nazwy pól są przepisane znak w znak z opisu rdzenia po
 * stronie powłoki, bo serde oddaje je bez przemianowania, a zmiana nazwy po tamtej stronie
 * rozspaja most po cichu.
 */
export interface StanRdzenia {
  /** Czy rdzeń odpowiada na porcie w chwili zapytania. */
  pracuje: boolean;
  /** Identyfikator procesu, gdy rdzeń postawiła powłoka. */
  pid: number | null;
  /** Czy proces rdzenia postawiła ta powłoka, czy zastała pracujący. */
  postawiony_przez_powloke: boolean;
  /** Adres HTTP rdzenia lokalnego. */
  adres: string;
  /** Ścieżka dziennika rdzenia prowadzonego przez powłokę. */
  dziennik: string;
  /** Zdanie opisujące przebieg uruchomienia — także przy niepowodzeniu. */
  opis: string;
}

/**
 * Pyta powłokę o adres HTTP rdzenia lokalnego.
 *
 * Zwraca `null` poza powłoką natywną i przy niepowodzeniu polecenia — w obu
 * przypadkach wywołujący zostaje przy własnym rozstrzygnięciu adresu.
 */
export async function adresRdzeniaZPowloki(): Promise<string | null> {
  if (!czyPowlokaNatywna()) return null;
  try {
    const adres = await invoke<string>(POLECENIE_ADRESU);
    return typeof adres === 'string' && adres !== '' ? adres : null;
  } catch (blad) {
    console.warn('[powłoka] adres rdzenia nieustalony', blad);
    return null;
  }
}

/**
 * Pyta powłokę o stan rdzenia postawionego w tle.
 *
 * Zwraca `null` poza powłoką natywną i przy niepowodzeniu polecenia. Odpowiedź
 * jest wyłącznie informacją — niczego nie wyłącza i nie blokuje.
 */
export async function stanRdzeniaZPowloki(): Promise<StanRdzenia | null> {
  if (!czyPowlokaNatywna()) return null;
  try {
    const stan = await invoke<StanRdzenia>(POLECENIE_STANU);
    return czyOpisRdzenia(stan) ? stan : null;
  } catch (blad) {
    console.warn('[powłoka] stan rdzenia nieustalony', blad);
    return null;
  }
}

/**
 * Sprawdzenie kształtu odpowiedzi. Most nie ufa kształtowi z drugiej strony granicy procesu:
 * odpowiedź niepasująca do umowy jest traktowana jak brak odpowiedzi, a nie wpuszczana do widoku.
 */
function czyOpisRdzenia(wartosc: unknown): wartosc is StanRdzenia {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const opis = wartosc as Partial<StanRdzenia>;
  return (
    typeof opis.pracuje === 'boolean' &&
    typeof opis.adres === 'string' &&
    typeof opis.dziennik === 'string' &&
    typeof opis.opis === 'string'
  );
}
