/* Składnik frazy nawigacyjnej przedstawia zdanie zamykające krok instalacji, zmieniające treść i barwę po nieudanej próbie przejścia dalej.

   Zdanie zamykające krok — wzorzec Windows Installer. Ma dwa stany: w spoczynku
   mówi, co zrobić; po nieudanej próbie przejścia dalej mówi, czego brakuje,
   i bierze barwę błędu. Stan niesie treść i barwa naraz, nigdy sama barwa.

   Właściwości:
     klucz      klucz katalogu — treść w spoczynku
     dane       atrybuty `data-*` */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.frazaNawigacyjna = function (w) {
  return el('p', {
    klasa: 'dn-kreator-nawigacja',
    'aria-live': 'polite',
    dane: w.dane || null,
    tekst: tekst(w.klucz)
  });
};
})();
