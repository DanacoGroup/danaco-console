/* Składnik bloku danych przedstawia wiersze złożone z nazwy i odczytu w jednolitym układzie, właściwym tabliczce znamionowej instalacji.

   Tabliczka znamionowa: wiersze „nazwa — odczyt" w jednym rytmie. Nazwa
   i odczyt idą tym samym pismem i w tej samej barwie; rozdziela je wyłącznie
   położenie, bo są jednym wierszem jednej tabliczki, nie dwiema rangami.

   Właściwości:
     naglowek   klucz katalogu — nagłówek bloku (opcjonalny)
     wiersze    tablica { etykieta, wartosc, dane } — łańcuchy albo klucze
     dane       atrybuty `data-*` na wierszu, do podmiany odczytu w czasie pracy */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.blokDanych = function (w) {
  var wykaz = el('div', { klasa: 'dn-wykaz-cichy' }, w.wiersze.map(function (r) {
    return el('div', { klasa: 'dn-zestawienie dn-zestawienie--cichy' }, [
      el('span', { tekst: r.etykieta }),
      el('b', { tekst: r.wartosc, dane: r.dane || null })
    ]);
  }));
  if (!w.naglowek) return wykaz;
  return el('div', {}, [K.skladniki.naglowekBloku({ klucz: w.naglowek }), wykaz]);
};
})();
