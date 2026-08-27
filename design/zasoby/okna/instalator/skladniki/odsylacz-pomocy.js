/* Odsyłacz pomocy jest przyciskiem, nie odnośnikiem, ponieważ otwiera objaśnienie w oknie zamiast prowadzić na zewnątrz, a klawiatura i czytnik ekranu mają otrzymać rolę przycisku. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.odsylaczPomocy = function (w) {
  var b = el('button', {
    klasa: 'dn-link dn-link--cichy', type: 'button',
    'aria-haspopup': 'dialog',
    dane: w.dane || null,
    tekst: tekst(w.klucz)
  });
  if (w.naOtwarcie) b.addEventListener('click', w.naOtwarcie);
  return el('p', {}, [b]);
};
})();
