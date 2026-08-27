/**
 * Punkt zbiorczy pływającego Asystenta. Na zewnątrz wychodzi jedno wywołanie
 * `zaczepAsystenta`, natomiast favikon, dymek, stan rozmowy i dostępność głosu
 * pozostają wnętrzem tej warstwy, niewidocznym z aplikacji.
 */

export { brakujaceOgniwa, czyGlosDziala, OGNIWA_GLOSU, type OgniwoGlosu } from './dostepnosc-mowy';
export {
  czyAsystent,
  nazwijSprawce,
  rozpoznajSprawce,
  zdanieOSprawcy,
  type KopertaZeSprawca,
  type ZnanySprawca,
} from './sprawca-zdarzenia';
export { zaczepAsystenta, type AsystentPlywajacy } from './zaczep-asystenta';
