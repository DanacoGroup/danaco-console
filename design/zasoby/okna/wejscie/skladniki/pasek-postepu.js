/* Pasek postępu z głową przedstawia nazwę czynności trwającej, jej wartość liczbową oraz tor, którego rolę progressbar niesie sam tor, nie nazwa nad nim.
   Wartość stoi w danej na znaczniku, nie w atrybucie `style` — miara jest
   daną, a wygląd należy do arkusza.

   Właściwości:
     etykieta   klucz katalogu — co się dzieje
     opisPaska  klucz katalogu — nazwa toru dla czytnika ekranu
     wartosc    0..100
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.pasekPostepu = function (N, w) {
  return N.el('div', { klasa: 'pg-postep' }, [
    N.el('div', { klasa: 'pg-postep-wiersz' }, [
      N.el('span', { tekst: N.tekst(w.etykieta) }),
      N.el('b', { tekst: w.wartosc + '%' })
    ]),
    N.el('div', { klasa: 'dn-postep' }, [
      N.el('div', {
        klasa: 'dn-postep-tor', role: 'progressbar',
        'aria-valuenow': String(w.wartosc), 'aria-valuemin': '0', 'aria-valuemax': '100',
        'aria-label': N.tekst(w.opisPaska)
      }, [
        N.el('div', { klasa: 'dn-postep-wartosc', dane: { wartosc: w.wartosc } })
      ])
    ])
  ]);
};
})();
