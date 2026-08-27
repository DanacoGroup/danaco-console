/* Ekran czwarty kreatora instalacji pozwala wskazać katalog docelowy programu oraz wybrać skróty tworzone na pulpicie i w menu startowym.

   Dwa katalogi i dwa skróty. Opisy pól mówią to, czego z etykiety nie widać —
   nie tłumaczą, czym jest skrót na pulpicie. Dane operatora idą do `Roaming`,
   nie do `Local`: projekty i ustawienia są pracą użytkownika, nie pamięcią
   podręczną, więc mają wchodzić do kopii profilu. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst, S = K.skladniki;
K.ekrany = K.ekrany || {};

K.ekrany[4] = function () {
  return el('section', {
    klasa: 'dn-kreator-ekran',
    dane: { ekran: 4 },
    'aria-label': tekst('szyna.kroki')[3]
  }, [
    S.naglowekEkranu({ nadtytul: 'krok4.nadtytul', tytul: 'krok4.tytul', podtytul: 'krok4.podtytul' }),
    S.poleSciezki({
      etykieta: 'krok4.katalogProgramu.etykieta', wartosc: 'krok4.katalogProgramu.wartosc',
      opis: 'krok4.katalogProgramu.opis', id: 'katalog-programu', zmien: 'krok4.zmien'
    }),
    S.poleSciezki({
      etykieta: 'krok4.katalogDanych.etykieta', wartosc: 'krok4.katalogDanych.wartosc',
      opis: 'krok4.katalogDanych.opis', id: 'katalog-danych', zmien: 'krok4.zmien'
    }),
    S.wierszMiary({ wzor: 'krok4.miejsce', dane: 'krok4.miejsceDane', wytluszcz: ['wymagane', 'dostepne'] }),
    el('div', {}, [
      S.naglowekBloku({ klucz: 'krok4.skroty.naglowek' }),
      S.poleWyboru({ etykieta: 'krok4.skroty.pulpit.etykieta', opis: 'krok4.skroty.pulpit.opis', id: 'opis-pulpit', zaznaczone: true })
    ]),
    S.poleWyboru({ etykieta: 'krok4.skroty.start.etykieta', opis: 'krok4.skroty.start.opis', id: 'opis-start', zaznaczone: true })
  ]);
};
})();
