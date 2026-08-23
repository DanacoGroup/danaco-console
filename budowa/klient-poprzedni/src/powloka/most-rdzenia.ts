import { invoke } from '@tauri-apps/api/core';

import { czyPowlokaNatywna } from './powloka-natywna';

/**
 * Most do wiedzy powłoki natywnej o rdzeniu — konsument poleceń
 * `adres_rdzenia` i `stan_rdzenia` (`desktop/src-tauri/src/polecenia.rs`).
 *
 * Powłoka stawia proces rdzenia (`desktop/src-tauri/src/rdzen/`), zna jego port
 * ze zmiennej `DANACO_PORT`, jego identyfikator procesu i ścieżkę dziennika.
 * Interfejs nie zna żadnej z tych rzeczy: adres gniazda wylicza z lokalizacji
 * dokumentu (`polaczenie/adres-rdzenia.ts`), a powodu ciszy nie zna wcale.
 *
 * Stąd dwa polecenia:
 *
 *   `adres_rdzenia` — adres HTTP rdzenia lokalnego. Potrzebny, gdy okno dostało
 *   interfejs z pakietu osadzonego w powłoce: strona ma wtedy pochodzenie
 *   `tauri.localhost`, a wyliczenie z lokalizacji daje adres gniazda
 *   `ws://tauri.localhost:…`, pod którym nie nasłuchuje nikt (rozbieżność
 *   opisuje `desktop/src-tauri/src/zrodlo_interfejsu.rs`).
 *
 *   `stan_rdzenia` — opis rdzenia w tle wraz z przebiegiem uruchomienia. Bez
 *   niego wskaźnik łączności umie powiedzieć wyłącznie „Rozłączony", bez powodu
 *   i bez wskazania dziennika.
 *
 * Poza powłoką natywną oraz przy niepowodzeniu polecenia odpowiedzią jest
 * `null`. Żadna ścieżka nie rzuca wyjątkiem i nie odrzuca obietnicy — brak
 * odpowiedzi powłoki niczego nie wstrzymuje, bo interfejs ma własną drogę
 * ustalenia adresu i własny stan łączności.
 *
 * Polecenia zmieniającego stan rdzenia tu nie ma: zatrzymanie rdzenia jest
 * czynnością z zasobnika (`desktop/src-tauri/src/rdzen/uchwyt.rs`), a most nie
 * tworzy drugiej drogi sterowania platformą.
 */

/** Nazwa polecenia powłoki; odpowiednik `polecenia::adres_rdzenia`. */
const POLECENIE_ADRESU = 'adres_rdzenia';

/** Nazwa polecenia powłoki; odpowiednik `polecenia::stan_rdzenia`. */
const POLECENIE_STANU = 'stan_rdzenia';

/**
 * Stan rdzenia widziany przez powłokę.
 *
 * Nazwy pól są przepisane z `OpisRdzenia` (`desktop/src-tauri/src/rdzen/uchwyt.rs`)
 * znak w znak, bo serde oddaje je bez przemianowania. Zmiana nazwy po stronie
 * powłoki rozspaja most po cichu — dlatego stoją tu dosłownie.
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
 * Sprawdzenie kształtu odpowiedzi.
 *
 * Most nie ufa kształtowi z drugiej strony granicy procesu bardziej niż
 * kształtowi z sieci: odpowiedź niepasująca do umowy jest traktowana jak brak
 * odpowiedzi, a nie wpuszczana do widoku jako `undefined` w środku zdania.
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
