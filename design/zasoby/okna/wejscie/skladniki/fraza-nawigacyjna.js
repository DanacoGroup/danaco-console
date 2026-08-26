/* ============================================================================
   SKŁADNIK — FRAZA NAWIGACYJNA

   Zdanie nad pasem działań: mówi, co będzie dalej, albo stawia pytanie
   z wplecioną czynnością („Nie masz konta? Utwórz konto Operatora"). Stoi
   w treści, nie w pasie — w pasie zostają same czynności, inaczej zdanie,
   czynność poboczna i główna nie mieszczą się w jednym wierszu.

   Właściwości:
     klucz     klucz katalogu — całe zdanie
     pytanie   klucz katalogu — część przed odsyłaczem
     czynnosc  klucz katalogu — napis odsyłacza
     cel       nazwa widoku, do którego odsyłacz prowadzi
     grupa     grupa widoku
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.frazaNawigacyjna = function (N, w) {
  if (w.klucz) return N.el('p', { klasa: 'we-fraza' }, [N.el('span', { tekst: N.tekst(w.klucz) })]);
  return N.el('p', { klasa: 'we-fraza' }, [
    w.pytanie && N.el('span', { klasa: 'we-fraza-pytanie', tekst: N.tekst(w.pytanie) }),
    w.pytanie && ' ',
    N.el('a', {
      klasa: 'au-link', href: '#s-' + w.cel, tekst: N.tekst(w.czynnosc),
      dane: { idz: w.cel, 'idz-grupa': w.grupa || 'stan' }
    })
  ]);
};
})();
