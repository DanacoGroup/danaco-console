/* Pas czynności biegnie pod obiema kolumnami okna, a jego skład wyprowadza się z kroku i odsłony, ponieważ każdy krok ma własny komplet czynności i rozstrzyga, których w nim brakuje. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

function przycisk(opis, klasa) {
  if (!opis) return null;
  return el('button', {
    klasa: klasa, type: 'button',
    hidden: !!opis.ukryty,
    dane: opis.dane || null,
    tekst: tekst(opis.klucz)
  });
}

K.skladniki.pasDzialan = function (w) {
  return el('div', { klasa: 'dn-kreator-pas' }, [
    el('div', { klasa: 'dn-kreator-pas-tresc' }, [
      el('span', { klasa: 'dn-kreator-pas-odstep', 'aria-hidden': 'true' }),
      przycisk(w.dodatkowa, 'dn-btn dn-btn--zarys'),
      przycisk(w.poboczna, 'dn-btn dn-btn--zarys'),
      przycisk(w.glowna, 'dn-btn dn-btn--atrament')
    ])
  ]);
};
})();
