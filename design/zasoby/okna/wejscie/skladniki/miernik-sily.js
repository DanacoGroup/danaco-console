/* ============================================================================
   SKŁADNIK — MIERNIK SIŁY HASŁA

   Cztery odcinki toru, zdanie o wyniku i cztery warunki. Znak warunku różni
   się KSZTAŁTEM, nie samą barwą: niespełniony niesie puste kółko, spełniony
   ptaszka. Ptaszek w każdym stanie czytał się jako „zrobione", a barwa jako
   jedyna różnica łamie WCAG 1.4.1.

   Miernik wiąże się z polem przez `data-sila-dla`, nie przez sąsiedztwo
   w drzewie — sąsiedztwo bywa różne w różnych oknach i cicho się rozjeżdża.

   Właściwości:
     dla   identyfikator pola hasła
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

var WARUNKI = ['dlugosc', 'wielkosc', 'cyfra', 'znak'];

W.skladniki.miernikSily = function (N, w) {
  var odcinki = [0, 1, 2, 3].map(function () { return N.el('span', { klasa: 'au-sila-odc' }); });

  var warunki = WARUNKI.map(function (k) {
    var kolko = N.zeZnacznika(W.ikony.kolko);
    kolko.setAttribute('class', 'ik-brak');
    kolko.setAttribute('aria-hidden', 'true');
    var ptaszek = N.zeZnacznika(W.ikony.ptaszek);
    ptaszek.setAttribute('class', 'ik-jest');
    ptaszek.setAttribute('aria-hidden', 'true');
    return N.el('span', {
      klasa: 'au-sila-warunek', dane: { warunek: k, spelniony: 'nie' }
    }, [kolko, ptaszek, N.tekst('dostep.sila.warunki.' + k)]);
  });

  return N.el('div', {
    klasa: 'au-sila', dane: { stopien: '0', 'sila-dla': w.dla }
  }, [
    N.el('div', { klasa: 'au-sila-tor', 'aria-hidden': 'true' }, odcinki),
    N.el('span', {
      klasa: 'au-sila-opis', id: w.dla + '-sila-opis',
      'aria-live': 'polite', tekst: N.tekst('dostep.sila.puste'),
      dane: { 'sila-opis': true }
    }),
    N.el('div', { klasa: 'au-sila-warunki' }, warunki)
  ]);
};
})();
