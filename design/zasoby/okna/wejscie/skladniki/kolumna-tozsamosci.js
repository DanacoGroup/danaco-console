/* ============================================================================
   SKŁADNIK — KOLUMNA TOŻSAMOŚCI

   Lewa strefa okna: sygnet z pulsującą kropką, nazwa, motto i trzy zdania
   o tym, czym program jest. Kolumna nie jest planszą marki — jest strefą okna,
   więc powierzchnia różni się od panelu treści odcieniem, nie kontrastem.

   Nota wydawcy NIE należy do tego składnika: stoi w wierszu pasa działań,
   niżej niż kolumna, i jest jedna dla całego okna.

   Właściwości:
     odslona   'uruchamianie' | 'dostep' — który zestaw zdań
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

var ZNAKI_ZALET = ['srodowiska', 'modele', 'tarcza'];

W.skladniki.kolumnaTozsamosci = function (N, w) {
  var godlo = N.zeZnacznika(W.ikony.godlo);
  godlo.setAttribute('class', 'we-marka-godlo');
  godlo.setAttribute('role', 'img');
  godlo.setAttribute('aria-label', N.tekst('marka.nazwa'));

  var zalety = N.tekst('marka.' + w.odslona + '.zalety').map(function (z, i) {
    var znak = N.zeZnacznika(W.ikony[ZNAKI_ZALET[i] || 'tarcza']);
    znak.setAttribute('aria-hidden', 'true');
    return N.el('li', { klasa: 'we-marka-poz' }, [
      znak,
      N.el('span', {}, [N.el('b', { tekst: z.glowa }), z.tresc])
    ]);
  });

  return N.el('div', { klasa: 'we-marka' }, [
    N.el('div', {}, [
      godlo,
      N.el('p', { klasa: 'we-marka-nazwa', tekst: N.tekst('marka.nazwa') }),
      N.el('p', { klasa: 'we-marka-motto', tekst: N.tekst('marka.' + w.odslona + '.motto') })
    ]),
    N.el('ul', { klasa: 'we-marka-lista' }, zalety)
  ]);
};

/* Nota wydawcy — osobny węzeł, bo siada w innym wierszu siatki okna. */
W.skladniki.notaWydawcy = function (N) {
  return N.el('p', { klasa: 'we-marka-stopka' }, [
    N.tekst('marka.wydawca'), N.el('br'),
    N.tekst('marka.wydanie'), N.el('br'),
    N.tekst('marka.wsparcie')
  ]);
};
})();
