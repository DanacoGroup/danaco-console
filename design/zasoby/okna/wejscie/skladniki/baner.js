/* ============================================================================
   SKŁADNIK — BANER KOMUNIKATU

   Głowa nazywa rzecz, treść mówi, co z niej wynika albo co zrobić. Wygląd
   wnosi składnik biblioteki `.dn-alert` w wariancie ze wstęgą: barwa stanu
   obejmuje znak i głowę, a wstęga przy lewej krawędzi niesie stan kształtem —
   barwa jako jedyna różnica nie wystarcza (WCAG 1.4.1).

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
  var RODZAJE = { informacja: 'info', ostrzezenie: 'ostrzezenie', blad: 'blad', sukces: 'sukces' };
  return N.el('div', {
    klasa: 'dn-alert dn-alert--wstega dn-alert--' + (RODZAJE[w.rodzaj] || 'info'),
    id: w.id || null,
    role: w.rodzaj === 'blad' || w.rodzaj === 'ostrzezenie' ? 'alert' : 'status'
  }, [
    N.el('span', { klasa: 'dn-alert-znak', 'aria-hidden': 'true' }, [znak]),
    N.el('span', { klasa: 'dn-alert-tresc' }, [N.el('b', { tekst: glowa }), tresc])
  ]);
};
})();
