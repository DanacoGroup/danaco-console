/* ============================================================================
   KREATOR INSTALACJI — mapa stanów

   Stany warunkowe ekranów nie są rozgałęzieniami w kodzie ekranu: ekran rysuje
   się raz, a stan jest DANĄ, którą się na nim ustawia. Ta mapa jest jedynym
   miejscem, gdzie zapisano, co dany stan zmienia — jaki tytuł, jaki podtytuł,
   które bloki widać.

   Wartości nie zawierają tekstu: stoją tu klucze katalogu treści. Zmiana słowa
   wchodzi w `tresci.json`, nie tutaj.

   ── EKRAN 3 · wykrycie procesora ────────────────────────────────────────────
      x64   wykryto Intel/AMD          podtytuł „x64",   zaznaczenie i plakietka na x64
      arm   wykryto ARM                podtytuł „arm",   zaznaczenie i plakietka na arm
      brak  nie rozpoznano procesora   podtytuł „brak",  nic nie zaznaczone, bez plakietek,
                                       wskazówka mówi, gdzie sprawdzić, czynność główna wyłączona

   ── EKRAN 5 · przebieg zapisu ───────────────────────────────────────────────
      przebieg     zapis trwa           pasek z miarą, wykaz etapów, pasek szczegółów
      wycofywanie  cofanie zmian        pasek bez miary, bez wykazu, czynności wyłączone
      blad         zapis nie powiódł się   bez paska i wykazu, blok błędu, trzy czynności

   ── EKRAN 6 · wynik instalacji ──────────────────────────────────────────────
      gotowe       ukończona            znak sukcesu przy tytule
      ostrzezenia  ukończona z uwagami  znak ostrzeżenia, baner z listą uwag

   ── STANY POZAEKRANOWE ──────────────────────────────────────────────────────
      krok 2 bez zgody     czynność główna wygaszona (`aria-disabled`), fraza
                           nawigacyjna zmienia treść i barwę po próbie przejścia
      wybór niezgodny      ostrzeżenie pod kartami, pytanie przy próbie przejścia
      anulowanie           pytanie o skutek; potwierdzenie uruchamia wycofywanie
   ============================================================================ */
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
