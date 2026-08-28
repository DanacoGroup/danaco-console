/* Składnik bloku błędu przedstawia kod i okoliczność niepowodzenia instalacji wraz ze zdaniem opisującym zalecane dalsze postępowanie.

   Kod i okoliczność niepowodzenia oraz zdanie o tym, co z tym zrobić.
   Szczegóły idą pismem maszynowym — to jedyne miejsce w oknie, gdzie treść jest
   naprawdę danymi: kod i ścieżka mają być przepisywalne bez przekłamania.

   Właściwości:
     szczegoly  klucz katalogu — kod i okoliczność
     rada       klucz katalogu — co zrobić
     dane       atrybuty `data-*` na obu częściach
     ukryty     true — blok czeka na swoją odsłonę */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.blokBledu = function (w) {
  return [
    el('pre', {
      klasa: 'dn-kod', hidden: !!w.ukryty,
      dane: Object.assign({ 'blad-szczegoly': true }, w.dane || {}),
      tekst: tekst(w.szczegoly)
    }),
    el('div', { klasa: 'dn-tekst-ciagly', hidden: !!w.ukryty, dane: w.dane || null },
      [el('p', { tekst: tekst(w.rada) })])
  ];
};
})();
