/* Ekran szósty kreatora instalacji przedstawia wynik zakończonej instalacji programu oraz czynności dostępne po pierwszym uruchomieniu.

   Znak wyniku stoi przy tytule, bo to tytuł orzeka o wyniku instalacji; blok
   „przy pierwszym uruchomieniu" zapowiada przyszłe czynności i żadnego znaku
   nie bierze. Dwa wyniki — gotowe i z ostrzeżeniami — niesie `data-wynik`. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst, S = K.skladniki;
K.ekrany = K.ekrany || {};

K.ekrany[6] = function () {
  return el('section', {
    klasa: 'dn-kreator-ekran',
    dane: { ekran: 6, wynik: 'gotowe' },
    'aria-label': tekst('szyna.kroki')[5]
  }, [
    S.naglowekEkranu({ nadtytul: 'krok6.nadtytul', tytul: 'krok6.tytulGotowe', podtytul: 'krok6.podtytulGotowe', znak: 'gotowe' }),
    S.baner({ rodzaj: 'ostrzezenie', ikona: 'ostrzezenie', glowa: 'krok6.ostrzezenie.glowa', tresc: 'krok6.ostrzezenie.tresc', dane: { ostrzezenia: true }, ukryty: true }),
    el('div', {}, [
      S.naglowekBloku({ klucz: 'krok6.pierwszeUruchomienie.naglowek' }),
      S.tekstCiagly({ klucze: ['krok6.pierwszeUruchomienie.tresc'] })
    ]),
    S.blokDanych({
      naglowek: 'krok6.szczegoly.naglowek',
      wiersze: [
        { etykieta: tekst('krok6.szczegoly.lokalizacja.etykieta'), wartosc: tekst('krok6.szczegoly.lokalizacja.wartosc') },
        { etykieta: tekst('krok6.szczegoly.wersja.etykieta'), wartosc: tekst('krok6.szczegoly.wersja.wartosc') },
        { etykieta: tekst('krok6.szczegoly.procesor.etykieta'), wartosc: tekst('krok3.x64.nazwa'), dane: { 'podsumowanie-procesor': true } }
      ]
    }),
    S.poleWyboru({ etykieta: 'krok6.przewodnik.etykieta', opis: 'krok6.przewodnik.opis', id: 'opis-przewodnik', zaznaczone: true }),
    S.frazaNawigacyjna({ klucz: 'krok6.odinstalowanie' })
  ]);
};
})();
