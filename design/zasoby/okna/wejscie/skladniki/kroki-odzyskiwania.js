/* Tor trzech kroków odzyskiwania prowadzi przez adres, potwierdzenie i nowe hasło, a stan kroku bieżącego oraz zrobionego niosą dwie dane elementu, data-biezacy i data-zrobiony.
   Znak kroku zrobionego rysuje arkusz przez `::after` — numer ustępuje
   ptaszkowi, więc stan nie stoi na samej barwie (norma WCAG 1.4.1). Wstawianie
   tu rysunku dublowałoby ten sam znak.

   Strzałki między krokami są rysunkiem, nie treścią — czytnik ekranu je pomija.

   Właściwości:
     biezacy   numer kroku bieżącego (1..3)
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.krokiOdzyskiwania = function (N, w) {
  var nazwy = N.tekst('dostep.odzyskiwanie.kroki');
  var dzieci = [];
  nazwy.forEach(function (nazwa, i) {
    var nr = i + 1;
    var dane = { biezacy: nr === w.biezacy ? 'tak' : 'nie' };
    if (nr < w.biezacy) dane.zrobiony = 'tak';
    dzieci.push(N.el('span', { klasa: 'au-krok', dane: dane }, [
      N.el('span', { klasa: 'au-krok-nr', tekst: String(nr) }),
      nazwa
    ]));
    if (nr < nazwy.length) {
      dzieci.push(N.el('span', { klasa: 'au-krok-strzalka', 'aria-hidden': 'true', tekst: '→' }));
    }
  });
  return N.el('div', { klasa: 'au-kroki' }, dzieci);
};
})();
