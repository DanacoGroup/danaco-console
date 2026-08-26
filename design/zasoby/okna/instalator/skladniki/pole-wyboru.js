/* ============================================================================
   SKŁADNIK — POLE WYBORU Z OPISEM
   ----------------------------------------------------------------------------
   Pole wyboru i zdanie mówiące, co się stanie po jego zaznaczeniu. Opis nie
   tłumaczy, czym jest pole wyboru — mówi o skutku. Kontrolka i jej opis są
   jedną grupą: dzieli je odstęp mniejszy niż do następnej grupy.

   Właściwości:
     etykieta   klucz katalogu — treść przy kontrolce
     opis       klucz katalogu — zdanie o skutku (opcjonalne)
     id         identyfikator opisu, wiązany przez `aria-describedby`
     zaznaczone true | false
     dane       atrybuty `data-*` na kontrolce
     stan       'nieczynne' | null
   ============================================================================ */
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
