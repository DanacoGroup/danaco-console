/* ============================================================================
   SKŁADNIK — NAGŁÓWEK BLOKU
   ----------------------------------------------------------------------------
   Nazwa grupy składników wewnątrz kroku. Stoi niżej od podtytułu kroku:
   nazywa jeden blok, a nie cały krok. Poniżej pewnego stopnia rangi rozróżnia
   już nie rozmiar, tylko waga i kontrast.

   Właściwości:
     klucz      klucz katalogu — treść nagłówka
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.naglowekBloku = function (w) {
  return el('h3', { klasa: 'dn-kreator-podtytul', tekst: tekst(w.klucz) });
};
})();
