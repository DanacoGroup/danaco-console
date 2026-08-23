/**
 * Komplet sterowania okna komunikacji — interfejs katalogu.
 *
 * Punkt wejścia klienta montuje komplet dla każdego okna z osobna:
 *
 *   const rejestrKanalow = utworzRejestrKanalow(kanal);
 *   rejestrKanalow.odswiez();
 *   const sterowanie = utworzPanelSterowania({ kanal, okno, rejestrKanalow });
 *   kontener.append(sterowanie.element);
 *
 * Rejestr kanałów powstaje raz na klienta — jest katalogiem wyboru, nie
 * ustawieniem okna. Panel powstaje raz na okno i nie ma z innym panelem
 * żadnej wspólnej zmiennej.
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
