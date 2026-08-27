/* Pole ścieżki z przyciskiem zmiany pokazuje katalog wskazany do instalacji jako kontrolkę tylko do odczytu, którą można zaznaczyć i skopiować klawiaturą, wraz z czynnością zmiany. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.poleSciezki = function (w) {
  var przycisk = el('button', {
    klasa: 'dn-btn dn-btn--zarys', type: 'button',
    disabled: w.stan === 'nieczynne',
    tekst: tekst(w.zmien)
  });
  if (w.naZmiane) przycisk.addEventListener('click', w.naZmiane);
  return el('div', { klasa: 'dn-pole' }, [
    el('label', { klasa: 'dn-pole-etykieta', for: w.id, tekst: tekst(w.etykieta) }),
    el('span', { klasa: 'dn-pole-zestaw' }, [
      el('input', {
        klasa: 'dn-pole-kontrolka', id: w.id, type: 'text', readonly: true,
        disabled: w.stan === 'nieczynne', value: tekst(w.wartosc)
      }),
      przycisk
    ]),
    el('span', { klasa: 'dn-pole-opis', tekst: tekst(w.opis) })
  ]);
};
})();
