/**
 * Punkt zbiorczy pływającego Asystenta.
 *
 * Na zewnątrz wychodzi jedno wywołanie — `zaczepAsystenta(kanal, posuniecia?,
 * idKlienta?)`. Reszta (favikon, dymek, stan rozmowy, dostępność głosu) jest
 * wnętrzem tej warstwy i nie ma po co być widoczna z aplikacji.
 *
 * Wyjątkiem są `czyGlosDziala` i `brakujaceOgniwa`: stan kanału głosowego
 * mówi o produkcie, nie o widoku, i przyda się każdemu, kto będzie budował
 * głos naprawdę.
 *
 * Drugim wyjątkiem jest `sprawca-zdarzenia.ts`: odczyt pól
 * `actor`/`actorClientId` z koperty przydaje się poza tą warstwą —
 * `aplikacja/zrodlo-posuniec.ts` rozpoznaje dziś sprawcę odciskiem okna, mając
 * prawdę w kopercie. Wychodzi tu po to, żeby nikt nie napisał go drugi raz
 * i inaczej.
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
