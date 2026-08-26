/* ============================================================================
   SKŁADNIK — TEKST CIĄGŁY
   ----------------------------------------------------------------------------
   Blok akapitów kroku. Akapity jednej myśli stoją w rytmie akapitu, nie
   w rytmie sekcji — inaczej dwa zdania tego samego wprowadzenia rozchodzą się
   na dwa osobne bloki.

   Właściwości:
     klucze     tablica kluczy katalogu — po jednym na akapit
     dane       atrybuty `data-*` na pierwszym akapicie (np. wskazówka kroku)
   ============================================================================ */
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
