/**
 * Widok sterowania okna — interfejs katalogu.
 *
 * Katalog nie buduje ani jednej kontrolki. Buduje oprawę, w której kontrolki
 * z `sterowanie/` stają się widoczne: kolumnę obok sceny okna, nagłówek
 * z rolą okna, podsumowanie ośmiu wartości oraz dwa uchwyty rozwijania —
 * w panelu i na pasku górnym.
 *
 * Punkt wejścia klienta montuje widok dla okna potwierdzonego przez rdzeń:
 *
 *   const rejestrKanalow = utworzRejestrKanalow(rdzen.kanal);
 *   rejestrKanalow.odswiez();
 *
 *   rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
 *     zamontujWidokSterowania({
 *       kanal: rdzen.kanal,
 *       okno,
 *       rejestrKanalow,
 *       panel: korzen.panel,
 *       akcjePaska: korzen.akcje,
 *     });
 *   });
 *
 * Widok powstaje raz na okno; rejestr kanałów raz na klienta, bo jest
 * katalogiem wyboru, nie ustawieniem okna.
 */
export {
  zamontujWidokSterowania,
  type MiejscaWidoku,
  type WidokSterowania,
  type ZaleznosciWidoku,
} from './montaz-widoku';

export { utworzObserwatorUstawien, type ObserwatorUstawien } from './obserwator-ustawien';

export { pozycjePodsumowania, type PozycjaPodsumowania } from './nazwy-ustawien';
