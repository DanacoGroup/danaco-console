/* Mapa stanów kreatora instalacji przechowuje wyłącznie klucze katalogu treści przypisane do stanów warunkowych ekranów, ponieważ stan jest daną ustawianą na ekranie, a nie rozgałęzieniem w jego kodzie. */
(function () {
'use strict';
var K = window.DanacoKreator;

K.stany = {
  wykrycie: {
    x64:  { podtytul: 'krok3.podtytulX64',  wskazowka: 'krok3.wskazowkaWykryta', zaznacz: 'x64' },
    arm:  { podtytul: 'krok3.podtytulArm',  wskazowka: 'krok3.wskazowkaWykryta', zaznacz: 'arm' },
    brak: { podtytul: 'krok3.podtytulBrak', wskazowka: 'krok3.wskazowkaBrak',    zaznacz: null }
  },
  odslona: {
    przebieg:    { tytul: 'krok5.przebieg.tytul',    podtytul: 'krok5.przebieg.podtytul',
                   pokaz: ['postep', 'kolejka', 'szczegoly'] },
    wycofywanie: { tytul: 'krok5.wycofywanie.tytul', podtytul: 'krok5.wycofywanie.podtytul',
                   pokaz: ['postep'], miara: 'nieokreslona' },
    blad:        { tytul: 'krok5.blad.tytul',        podtytul: 'krok5.blad.podtytul',
                   pokaz: ['blad'] }
  },
  wynik: {
    gotowe:      { tytul: 'krok6.tytulGotowe',      podtytul: 'krok6.podtytulGotowe',      znak: 'gotowe' },
    ostrzezenia: { tytul: 'krok6.tytulOstrzezenia', podtytul: 'krok6.podtytulOstrzezenia', znak: 'ostrzezenia' }
  },
  nawigacjaKroku2: { spoczynek: 'krok2.nawigacja', brakZgody: 'krok2.nawigacjaBrakZgody' }
};
})();
