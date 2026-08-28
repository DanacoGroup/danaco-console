/**
 * Warstwa Always On Display — punkt zbiorczy katalogu. Wystawia dwie czynności:
 * zaczepienie warstwy do osadzenia w powłoce oraz otwarcie powierzchni interakcji
 * jako rozszerzenia bocznego. Warstwa jest jedna na klienta i żyje między
 * otwarciami.
 */
import type { Kanal } from '../protokol/kanal';
import type { OpisKolumnyAod } from './kolumna-aod';
import { utworzWarstweAod, type WarstwaAod } from './warstwa-aod';

/**
 * Warstwa zbudowana przy pierwszym wywołaniu; `null` przed nim. Trzymana
 * w module, ponieważ jedna warstwa obsługuje całego klienta i żyje między
 * otwarciami powierzchni.
 */
let warstwa: WarstwaAod | null = null;

/**
 * Kanał, na którym zbudowano warstwę — podstawa rozpoznania zmiany połączenia.
 * Wywołanie z innym kanałem rozłącza warstwę dotychczasową i buduje ją na nowo.
 */
let osadzonyKanal: Kanal | null = null;

/**
 * Buduje albo zwraca warstwę Always On Display.
 *
 * @param opis nastawy wykraczające poza kanał — dziś wyłącznie tożsamość
 *   klienta z powitania, bez której ster przejęcia nie przestawia ogniska.
 */
export function zaczepAod(kanal: Kanal, opis: OpisKolumnyAod = {}): WarstwaAod {
  if (warstwa !== null && osadzonyKanal !== kanal) {
    warstwa.rozlacz();
    warstwa = null;
  }

  if (warstwa === null) {
    warstwa = utworzWarstweAod(kanal, opis);
    osadzonyKanal = kanal;
  }

  return warstwa;
}

/**
 * Otwiera powierzchnię interakcji jako rozszerzenie boczne. Powierzchnia nie
 * czeka na rdzeń: kolumna pojawia się od razu, a odczyty dojeżdżają do niej
 * odpowiedzią. Warstwa niezaczepiona w powłoce osadza się przy tym wywołaniu
 * w korzeniu dokumentu.
 */
export function otworzPowierzchnieAod(kanal: Kanal, opis: OpisKolumnyAod = {}): WarstwaAod {
  const zaczepiona = zaczepAod(kanal, opis);
  if (!zaczepiona.element.isConnected) document.body.append(zaczepiona.element);
  zaczepiona.otworzPowierzchnie();
  return zaczepiona;
}

export type { KolumnaAod, OpisKolumnyAod } from './kolumna-aod';
export type { WarstwaAod } from './warstwa-aod';
export type { TozsamoscDlaOgniska } from './cztery-stery';
export { TrybObecnosci, type StanObecnosci } from './tryb-obecnosci';
// Wyciszenie nakładki: wykaz wyciszeń czynnych i magazyn jego stanu. Magazyn jest
// podawany z zewnątrz, więc przeniesienie zapisu do rdzenia nie rusza żadnego
// wołacza.
export {
  KlasaZdarzen,
  NAZWY_KLAS,
  ZakresKontekstu,
  zdanieWyciszenia,
  type MagazynWyciszen,
  type StanWyciszen,
  type WyciszenieCzynne,
} from './wyciszenie-aod';
export { BRAKI_WYCISZENIA, brakiCzynne, zdanieGranicyWyciszenia } from './wyciszenie-braki-kontraktu';
