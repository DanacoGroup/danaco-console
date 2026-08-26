/* ============================================================================
   SKŁADNIK — ODSYŁACZ POMOCY
   ----------------------------------------------------------------------------
   Czynność, nie adres — otwiera objaśnienie w oknie, nie prowadzi na zewnątrz.
   Dlatego niesie go przycisk, a nie odnośnik: klawiatura i czytnik ekranu mają
   dostać przycisk. Podkreślenia w spoczynku nie ma, bo odsyłacz stoi samodzielnie,
   poza tokiem zdania.

   Właściwości:
     klucz      klucz katalogu — treść odsyłacza
     dane       atrybuty `data-*`
     naOtwarcie wywołanie zwrotne
   ============================================================================ */
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
