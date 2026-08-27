/* Pole wyboru z opisem łączy kontrolkę wyboru ze zdaniem mówiącym o skutku jej zaznaczenia, nie o tym, czym jest pole wyboru, a oba elementy tworzą jedną grupę wizualną. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.poleWyboru = function (w) {
  var kontrolka = el('input', {
    klasa: 'dn-check', type: 'checkbox',
    checked: !!w.zaznaczone,
    disabled: w.stan === 'nieczynne',
    'aria-describedby': w.opis ? w.id : null,
    dane: w.dane || null
  });
  var dzieci = [el('label', { klasa: 'dn-wybor' }, [kontrolka, el('span', { tekst: tekst(w.etykieta) })])];
  if (w.opis) dzieci.push(el('p', { klasa: 'dn-pole-opis', id: w.id, tekst: tekst(w.opis) }));
  return el('div', { klasa: 'dn-pole' }, dzieci);
};
})();
