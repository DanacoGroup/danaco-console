import type { Kanal } from '../protokol/kanal';
import type { OpisKolumnyAod } from './kolumna-aod';
import { utworzWarstweAod, type WarstwaAod } from './warstwa-aod';

/**
 * Warstwa Always On Display — punkt zbiorczy katalogu.
 *
 * Reszta aplikacji zna stąd dwie czynności:
 *   • `zaczepAod(kanal, opis?)` — buduje warstwę do osadzenia w powłoce; to
 *     ona wnosi pływający awatar widoczny bez interakcji (warstwa 1);
 *   • `otworzPowierzchnieAod(kanal, opis?)` — otwiera powierzchnię interakcji
 *     jako rozszerzenie boczne; wołana z listwy ustawień strefy 3.
 *
 * Warstwa jest JEDNA na klienta i żyje między otwarciami powierzchni: kolejka
 * decyzji dosypuje się z `progress.changed` i `window.state.changed` także
 * wtedy, gdy kolumna stoi zwinięta, więc plakietka awatara jest prawdziwa bez
 * jednego kliknięcia.
 *
 * Kanał podajemy przy pierwszym zaczepieniu. Wywołanie z innym kanałem — po
 * ponownym połączeniu z rdzeniem — buduje warstwę na nowo, żeby przyszłe
 * komendy `aod.*` nie szły przez transport, którego już nie ma.
 */

/** Warstwa zbudowana przy pierwszym wywołaniu; `null` przed nim. */
let warstwa: WarstwaAod | null = null;

/** Kanał, na którym zbudowano warstwę — podstawa rozpoznania zmiany połączenia. */
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
 * Otwiera powierzchnię interakcji jako rozszerzenie boczne.
 *
 * Powierzchnia nie czeka na rdzeń: kolumna pojawia się od razu, a odczyt
 * `aod.status.get`, `aod.context.get` i `aod.suggestion` dojeżdża do niej
 * odpowiedzią. Czynności — `aod.observe.attach`, `aod.observe.detach`,
 * `aod.chat.send` i `aod.voice.command` — jadą z sekcji na żądanie Operatora.
 *
 * Warstwa niezaczepiona w powłoce osadza się przy tym wywołaniu w korzeniu
 * dokumentu: pozycja „Always On Display" listwy ustawień ma otwierać
 * powierzchnię naprawdę, a nie odmawiać z powodu kolejności montażu.
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
// PODAWANY, więc część druga zlecenia — wyciszenie jako byt rdzenia — przełoży
// zapis, nie ruszając ani jednego wołacza.
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
