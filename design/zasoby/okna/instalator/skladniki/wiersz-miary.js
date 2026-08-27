/* Wiersz miary przedstawia rachunek w jednym wierszu, gdzie liczby idą pismem maszynowym jako dane, a nie tekst interfejsu, co jest jedynym uzasadnieniem mieszania krojów. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.wierszMiary = function (w) {
  var dane = tekst(w.dane), wzor = tekst(w.wzor);
  var wyroznione = w.wytluszcz || [];
  var czesci = wzor.split(/(\{\w+\})/);
  return el('p', { klasa: 'dn-meta' }, czesci.map(function (c) {
    var m = c.match(/^\{(\w+)\}$/);
    if (!m) return document.createTextNode(c);
    var v = dane[m[1]] !== undefined ? dane[m[1]] : c;
    return wyroznione.indexOf(m[1]) >= 0 ? el('b', { tekst: v }) : document.createTextNode(v);
  }));
};
})();
