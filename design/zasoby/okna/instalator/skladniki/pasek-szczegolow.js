/* ============================================================================
   SKŁADNIK — PASEK SZCZEGÓŁÓW ETAPU
   ----------------------------------------------------------------------------
   Wiersz pod paskiem postępu: numer etapu, licznik tego, co etap właśnie
   przenosi, i pozostały czas. Licznik należy do etapu oznaczonego jako w toku —
   etap, który nie liczy nic, zostawia miejsce puste, bo cudza miara pod nazwą
   bieżącego etapu jest sprzecznością danych.

   Właściwości:
     dane       atrybuty `data-*` na pasku
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.pasekSzczegolow = function (w) {
  return el('div', { klasa: 'dn-zestawienie', dane: w.dane || { szczegoly: true } }, [
    el('span', { dane: { 'licznik-etapu': true } }),
    el('span', { klasa: 'dn-meta', dane: { pozostalo: true } })
  ]);
};
})();
