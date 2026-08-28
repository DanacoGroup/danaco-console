/* Składnik wyświetla nagłówek ekranu kreatora złożony z nadtytułu, tytułu i opcjonalnego podtytułu, przy czym wariant ze znakiem wyniku dokłada przed tytułem ikonę powtarzającą jego treść. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

var ZNAKI = {
  gotowe: { ikona: 'wynikGotowe', klasa: 'dn-tytul-znak--sukces' },
  ostrzezenia: { ikona: 'wynikOstrzezenia', klasa: 'dn-tytul-znak--ostrzezenie' }
};

K.skladniki.naglowekEkranu = function (w) {
  var dzieci = [];
  if (w.nadtytul) dzieci.push(el('p', { klasa: 'dn-kreator-nadtytul', tekst: tekst(w.nadtytul) }));

  if (w.znak) {
    var znaki = Object.keys(ZNAKI).map(function (nazwa) {
      var z = ZNAKI[nazwa];
      var svg = N.zeZnacznika(K.ikony[z.ikona]);
      svg.setAttribute('class', 'dn-tytul-znak ' + z.klasa);
      svg.setAttribute('data-znak', nazwa);
      svg.setAttribute('aria-hidden', 'true');
      if (nazwa !== w.znak) svg.setAttribute('hidden', '');
      return svg;
    });
    dzieci.push(el('h2', { klasa: 'dn-kreator-tytul dn-tytul-ze-znakiem' },
      znaki.concat([el('span', { dane: { tytul: true }, tekst: tekst(w.tytul) })])));
  } else {
    dzieci.push(el('h2', { klasa: 'dn-kreator-tytul', dane: { tytul: true }, tekst: tekst(w.tytul) }));
  }

  if (w.podtytul) {
    dzieci.push(el('p', { klasa: 'dn-kreator-lid', dane: { wstep: true }, tekst: tekst(w.podtytul) }));
  }
  return el('div', {}, dzieci);
};
})();
