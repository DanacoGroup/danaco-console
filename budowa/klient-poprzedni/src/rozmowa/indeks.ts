/**
 * Warstwa rozmowy okna komunikacji — interfejs katalogu.
 *
 * Powłoka sięga po rozmowę wyłącznie stąd; pozostałe pliki katalogu są
 * wewnętrzne. Złożenie w punkcie wejścia sprowadza się do dwóch wierszy:
 *
 *   rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
 *     zamontujRozmowe(korzen.scena, rdzen.kanal, okno.id, {
 *       persona: opis.kanalModelu,
 *       rolaOkna: okno.windowRole,
 *     });
 *   });
 *
 * Nazwy komend i zdarzeń, których warstwa używa — `message.send`,
 * `message.stop`, `stream.chunk`, `message.changed`, `window.changed`,
 * `window.state.get`, `action.list` — pochodzą wyłącznie z `shared/contract`.
 *
 * Okno przestawia się na moduł samo: śledzi moduł okna w rdzeniu i przy jego
 * zmianie rekonfiguruje pasek narzędzi promptu, panel akcji i panel kontekstu,
 * zachowując wątek rozmowy. Powłoka może przestawić okno wprost —
 * `zamontujRozmowe(...).ustawModul(kod)` — gdy zna wynik `workspace.enter`.
 */

export { zamontujRozmowe, type OpcjeMontazu, type ZamontowanaRozmowa } from './montaz-rozmowy';

export { utworzRozmowe, type OpcjeRozmowy, type Rozmowa } from './rozmowa';

export {
  utworzWidokRozmowy,
  type OpcjeWidokuRozmowy,
  type WidokRozmowy,
} from './widok-rozmowy';

/**
 * Widok transkryptu — cztery tryby pokazywania wątku. Menu sesji powłoki bierze
 * stąd wykaz i klucz techniczny (`widok_zapisu`), a przestawia okno przez
 * `ZamontowanaRozmowa`.
 */
export {
  czyWpisWZapisie,
  rozpoznajWidokZapisu,
  warstwyZapisu,
  WIDOK_ZAPISU_DOMYSLNY,
  WIDOKI_ZAPISU,
  type OpisWidokuZapisu,
  type WarstwyZapisu,
  type WidokZapisu,
} from './widok-zapisu';

export { sledzModulOkna, type SledzenieModulu } from './modul-okna';

export { opisStanu, type StanRozmowy } from './stan-rozmowy';

export {
  czyUlotna,
  politykaModulu,
  politykaTrwala,
  ZDANIE_O_ZAPISIE_RDZENIA,
  type KontekstRoboczy,
  type PolitykaUlotnosci,
} from './ulotnosc';

export {
  RodzajNadawcy,
  rozpoznajNadawce,
  znakiNadawcy,
  type ZnakiNadawcy,
} from './nadawca';

export { RodzajFragmentu } from './rodzaje-fragmentow';

export {
  godzinaWpisu,
  type BladWpisu,
  type StanWpisu,
  type WpisRozmowy,
  type WywolanieNarzedzia,
} from './wpis-rozmowy';

export { wierszWywolania, type Prowenancja } from './prowenancja';

export { opisKonta, type MetadaneKonta } from './metadane-konta';

export { opisPodsumowania, type PodsumowanieTury } from './podsumowanie-tury';
