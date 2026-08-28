/* Tekst ciągły składa blok akapitów kroku, przy czym akapity jednej myśli stoją w rytmie akapitu, a nie sekcji, aby zdania jednego wprowadzenia nie rozchodziły się na osobne bloki. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.tekstCiagly = function (w) {
  return el('div', { klasa: 'dn-tekst-ciagly' }, w.klucze.map(function (k, i) {
    return el('p', { tekst: tekst(k), dane: (i === 0 && w.dane) ? w.dane : null });
  }));
};
})();
