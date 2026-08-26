/* ============================================================================
   SKŁADNIK — BANER KOMUNIKATU

   Głowa nazywa rzecz, treść mówi, co z niej wynika albo co zrobić. Baner
   odcina się wstęgą przy lewej krawędzi, nie samym tłem — barwa jako jedyna
   różnica nie wystarcza (WCAG 1.4.1).

   Właściwości:
     rodzaj  'informacja' | 'ostrzezenie' | 'blad' | 'sukces'
     ikona   nazwa znaku z zestawu
     glowa   klucz katalogu
     tresc   klucz katalogu
     dane      dane do podstawienia w treść
     daneGlowy dane do podstawienia w głowę
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.baner = function (N, w) {
  var znak = N.zeZnacznika(W.ikony[w.ikona || 'informacja']);
  znak.setAttribute('aria-hidden', 'true');
  var tresc = N.tekst(w.tresc);
  if (w.dane) tresc = N.podstaw(tresc, w.dane);
  var glowa = N.tekst(w.glowa);
  if (w.daneGlowy) glowa = N.podstaw(glowa, w.daneGlowy);
  return N.el('div', {
    klasa: 'we-alarm we-alarm--' + (w.rodzaj || 'informacja'),
    id: w.id || null,
    role: w.rodzaj === 'blad' || w.rodzaj === 'ostrzezenie' ? 'alert' : 'status'
  }, [
    znak,
    N.el('span', {}, [N.el('b', { tekst: glowa }), tresc])
  ]);
};
})();
