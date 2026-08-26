/* ============================================================================
   SKŁADNIK — ZGODA NA TRWAŁĄ SESJĘ

   Pole wyboru z opisem skutku. Opis mówi, co ta zgoda daje i kiedy jej nie
   zaznaczać — samo „pozostań zalogowany" nie niesie ani jednego, ani drugiego.
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.poleSesji = function (N) {
  return N.el('div', { klasa: 'au-opcje' }, [
    N.el('label', { klasa: 'au-opcja' }, [
      N.el('input', { klasa: 'dn-check', type: 'checkbox' }),
      N.el('span', {}, [
        N.tekst('dostep.sesja.etykieta'),
        N.el('small', { tekst: N.tekst('dostep.sesja.opis') })
      ])
    ])
  ]);
};
})();
