/* ============================================================================
   SKŁADNIK — OKNO DIALOGOWE
   ----------------------------------------------------------------------------
   Pytanie zadane w toku pracy: potwierdzenie zamknięcia, ostrzeżenie o wyborze,
   objaśnienie. Wyjście domyślne prowadzi do dalszej pracy, nie do jej przerwania
   — dlatego to ono, a nie czynność niszcząca, bierze fokus przy otwarciu.

   Rola `alertdialog` należy się pytaniu o skutek, `dialog` — objaśnieniu.

   Właściwości:
     nazwa      identyfikator okna (`data-modal`)
     dane       dodatkowe atrybuty `data-*` na tle okna
     rola       'alertdialog' | 'dialog'
     tytul      klucz katalogu
     tresc      klucz katalogu (pomijany, gdy podano `cialo`)
     cialo      gotowy węzeł treści (dla objaśnień dłuższych niż zdanie)
     czynnosci  tablica { klucz, klasa, dane } — od domyślnej do ostatecznej
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.oknoDialogowe = function (w) {
  var idT = 'dlg-' + w.nazwa + '-tytul';
  var idO = 'dlg-' + w.nazwa + '-tresc';
  var cialo = w.cialo || el('p', { klasa: 'dn-modal-lid', id: idO, dane: w.daneTresci || null, tekst: tekst(w.tresc) });
  var dane = { modal: w.nazwa };
  if (w.dane) Object.keys(w.dane).forEach(function (k) { dane[k] = w.dane[k]; });
  return el('div', { klasa: 'dn-modal-tlo', hidden: true, dane: dane }, [
    el('div', {
      klasa: 'dn-modal', role: w.rola, 'aria-modal': 'true',
      'aria-labelledby': idT, 'aria-describedby': w.cialo ? null : idO
    }, [
      el('div', { klasa: 'dn-modal-naglowek' },
        [el('h2', { klasa: 'dn-modal-tytul', id: idT, dane: w.daneTytulu || null, tekst: tekst(w.tytul) })]),
      el('div', { klasa: 'dn-modal-cialo' }, [cialo]),
      el('div', { klasa: 'dn-modal-stopka' }, w.czynnosci.map(function (c) {
        return el('button', { klasa: c.klasa, type: 'button', dane: c.dane || null, tekst: tekst(c.klucz) });
      }))
    ])
  ]);
};
})();
