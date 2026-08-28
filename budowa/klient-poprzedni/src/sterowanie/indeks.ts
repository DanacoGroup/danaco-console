/**
 * Komplet sterowania okna komunikacji — interfejs katalogu montowany osobno dla każdego
 * okna klienta osobno.
 */
export {
  utworzPanelSterowania,
  type PanelSterowania,
  type ZaleznosciPanelu,
} from './panel-sterowania';

export { utworzRejestrKanalow, type RejestrKanalow } from './rejestr-kanalow';

export {
  KluczUstawieniaOkna,
  ZASIEG_OKNA,
  type UstawieniaOkna,
} from './klucze-ustawien';
