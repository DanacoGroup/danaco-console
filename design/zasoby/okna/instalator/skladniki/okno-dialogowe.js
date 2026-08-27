/* Okno dialogowe zadaje pytanie w toku pracy, takie jak potwierdzenie zamknięcia, ostrzeżenie o wyborze albo objaśnienie, a przy otwarciu fokus bierze wyjście domyślne prowadzące do dalszej pracy, nie czynność niszcząca. */
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
