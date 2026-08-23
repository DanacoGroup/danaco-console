/**
 * Okna pomocnicze modułów inżynierskich — jedno wejście obszaru.
 *
 * Obszar niesie okna towarzyszące rozmowie w modułach Developer i Diagnostics
 * (terminal, przeglądarka, artefakty, pliki, podgląd w tle bash, pliki
 * środowiska) oraz spis tych, których jeszcze nie ma, wraz z powodem braku.
 *
 * Obszar jest wspólny, nie modułowy: stoi poza `moduly/`, bo pas składa się
 * identycznie w obu modułach, a Terminal ma dojść w module Apps. Kopia
 * w każdym module byłaby kolejnym miejscem do rozejścia się treści.
 */
export { budzetRozmowy, type BudzetRozmowy } from './budzet-rozmowy';
export { type OpcjePanelu, type PanelPomocniczy } from './panel-pomocniczy';
export { paneleNieotwieralne, paneleOtwieralne, type PozycjaOtwieralna } from './panele-otwieralne';
export { KODY_Z_WYTWORNIA, wytworniaPanelu, type WytworniaPanelu } from './wytwornia-paneli';
export { utworzOknoHistoriiRozmowy, type OknoHistoriiRozmowy } from './okno-historii-rozmowy';
export { utworzOknoPodgladuBash, type OknoPodgladuBash, type OpcjePodgladuBash } from './okno-podglad-bash';
export { utworzZrodloHistorii, type ZrodloHistorii } from './zrodlo-historii';
export { utworzPasPomocniczych, type OpcjePasa, type PasPomocniczych } from './pas-pomocniczych';
export {
  oknaPomocnicze,
  type OpisPomocniczego,
  type StanPomocniczego,
} from './rejestr-pomocniczych';
export { utworzZrodloPodgladuBash, type ZrodloPodgladuBash } from './zrodlo-podgladu-bash';
