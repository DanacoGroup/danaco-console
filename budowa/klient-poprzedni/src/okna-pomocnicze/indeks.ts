/**
 * Plik zbiera i udostępnia elementy obszaru okien pomocniczych, wspólnego dla modułów Developer i Diagnostics, obejmującego terminal, przeglądarkę, artefakty, pliki oraz podgląd bash.
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
