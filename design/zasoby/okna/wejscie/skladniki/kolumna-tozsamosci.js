/* Kolumna tożsamości zajmuje lewą strefę okna: niesie sygnet z pulsującą kropką, nazwę, motto i trzy zdania o tym, czym program jest, w odsłonie zależnej od etapu wejścia.
   Kolumna nie jest planszą marki — jest strefą okna, więc powierzchnia różni
   się od panelu treści odcieniem, nie kontrastem.

   Nota wydawcy NIE należy do tego składnika: stoi w wierszu pasa działań,
   niżej niż kolumna, i jest jedna dla całego okna.

   W odsłonie przygotowania zamiast zdań stoi animacja powłok: na tym etapie
   użytkownik już nie wybiera programu, tylko czeka, aż się złoży.

   Właściwości:
     odslona   'uruchamianie' | 'dostep' | 'przygotowanie'
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

  var glowa = N.el('div', {}, [
    godlo,
    N.el('p', { klasa: 'we-marka-nazwa', tekst: N.tekst('marka.nazwa') }),
    N.el('p', { klasa: 'we-marka-motto', tekst: N.tekst('marka.' + w.odslona + '.motto') })
  ]);

  if (w.odslona === 'przygotowanie') {
    return N.el('div', { klasa: 'we-marka' }, [
      glowa,
      /* Pole puste — rysunek wnosi składnik biblioteki, nie znacznik okna. */
      N.el('div', {
        klasa: 'pg-bryla-pole dn-powloki', role: 'img',
        'aria-label': N.tekst('przygotowanie.obszarPowlok'),
        dane: { powloki: true }
      })
    ]);
  }

  var zalety = N.tekst('marka.' + w.odslona + '.zalety').map(function (z, i) {
    var znak = N.zeZnacznika(W.ikony[ZNAKI_ZALET[i] || 'tarcza']);
    znak.setAttribute('aria-hidden', 'true');
    return N.el('li', { klasa: 'we-marka-poz' }, [
      znak,
      N.el('span', {}, [N.el('b', { tekst: z.glowa }), z.tresc])
    ]);
  });

  return N.el('div', { klasa: 'we-marka' }, [glowa, N.el('ul', { klasa: 'we-marka-lista' }, zalety)]);
};

/* Nota w wierszu pasa działań — osobny węzeł, bo siada w innym wierszu siatki
   okna. Wydawca w oknach przed wejściem, stan przywracania w oknie
   przygotowania: w obu miejscach to rzecz, która ma stać przy dnie i nie
   należy do żadnej odsłony. */
W.skladniki.notaPasa = function (N, klucz) {
  return N.el('p', { klasa: 'we-marka-stopka', tekst: N.tekst(klucz) });
};

W.skladniki.notaWydawcy = function (N) {
  return N.el('p', { klasa: 'we-marka-stopka' }, [
    N.tekst('marka.wydawca'), N.el('br'),
    N.tekst('marka.wydanie'), N.el('br'),
    N.tekst('marka.wsparcie')
  ]);
};
})();
