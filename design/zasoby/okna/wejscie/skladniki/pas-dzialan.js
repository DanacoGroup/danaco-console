/* Pas działań domyka okno na całej jego szerokości, także pod kolumną tożsamości, a czynności zbiera przy prawej krawędzi okna zgodnie z resztą tej rodziny okien.
   Pas urwany na granicy kolumny czyta się jak niedokończony rysunek, a gdy
   zostaje jedna czynność, nie zawisa samotnie po lewej stronie.

   Właściwości:
     widok      nazwa widoku, z którym pas się przełącza
     grupa      grupa widoku
     aktywny    true — pas widoczny na starcie
     czynnosci  wykaz: { klucz, klasa, cel, grupa, dane, komunikat }
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.pasDzialan = function (N, w) {
  var czynnosci = (w.czynnosci || []).map(function (c) {
    var dzieci = [];
    if (c.ikona) {
      var i = N.zeZnacznika(W.ikony[c.ikona]);
      i.setAttribute('aria-hidden', 'true');
      dzieci.push(i);
    }
    dzieci.push(N.tekst(c.klucz));
    var dane = {};
    if (c.cel) { dane.idz = c.cel; dane['idz-grupa'] = c.grupa || 'etap'; }
    if (c.dane) Object.keys(c.dane).forEach(function (k) { dane[k] = c.dane[k]; });
    if (c.komunikat) {
      dane.komunikat = N.tekst('komunikaty.' + c.komunikat + '.tresc');
      dane['komunikat-tytul'] = N.tekst('komunikaty.' + c.komunikat + '.tytul');
      dane['komunikat-rodzaj'] = 'informacja';
    }
    return N.el('button', {
      klasa: c.klasa || 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button', id: c.id || null, dane: dane
    }, dzieci);
  });

  return N.el('div', {
    klasa: 'we-pas',
    dane: { widok: w.widok, 'grupa-widoku': w.grupa, 'widok-aktywny': w.aktywny ? 'tak' : 'nie' }
  }, czynnosci);
};
})();
