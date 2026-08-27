/** Plik stanowi interfejs katalogu warstwy rozmowy okna komunikacji: powłoka sięga po rozmowę wyłącznie stąd, a pozostałe pliki katalogu pozostają wewnętrzne. */

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
