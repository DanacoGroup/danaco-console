/* Ekran pierwszy kreatora instalacji przedstawia wymagania sprzętowe i systemowe niezbędne do rozpoczęcia wdrożenia programu.

   Otwarcie kreatora. Nagłówek nazywa czynność, nie wita. Podmiotem zdań jest
   instalacja albo ten komputer, nie program mówiący o sobie w trzeciej osobie.
   Wymóg połączenia stoi wyłącznie w banerze — w tabliczce byłby tą samą rzeczą
   powiedzianą dwa razy. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst, S = K.skladniki;
K.ekrany = K.ekrany || {};

K.ekrany[1] = function () {
  return el('section', {
    klasa: 'dn-kreator-ekran',
    dane: { ekran: 1, widoczny: 'tak' },
    'aria-label': tekst('szyna.kroki')[0]
  }, [
    S.naglowekEkranu({ nadtytul: 'krok1.nadtytul', tytul: 'krok1.tytul', podtytul: 'krok1.podtytul' }),
    S.tekstCiagly({ klucze: ['krok1.akapit1', 'krok1.akapit2'] }),
    S.blokDanych({
      naglowek: 'krok1.wymagania.naglowek',
      wiersze: [
        { etykieta: tekst('krok1.wymagania.wersja.etykieta'), wartosc: tekst('krok1.wymagania.wersja.wartosc'), dane: { dane: true } },
        { etykieta: tekst('krok1.wymagania.system.etykieta'), wartosc: tekst('krok1.wymagania.system.wartosc') },
        { etykieta: tekst('krok1.wymagania.procesor.etykieta'), wartosc: tekst('krok1.wymagania.procesor.wartosc') },
        { etykieta: tekst('krok1.wymagania.miejsce.etykieta'), wartosc: tekst('krok1.wymagania.miejsce.wartosc') },
        { etykieta: tekst('krok1.wymagania.uprawnienia.etykieta'), wartosc: tekst('krok1.wymagania.uprawnienia.wartosc') }
      ]
    }),
    S.baner({ rodzaj: 'ostrzezenie', ikona: 'wifi', glowa: 'krok1.polaczenie.glowa', tresc: 'krok1.polaczenie.tresc' }),
    S.frazaNawigacyjna({ klucz: 'krok1.nawigacja' })
  ]);
};
})();
