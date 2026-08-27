/* Zgoda na trwałą sesję jest polem wyboru z opisem skutku, który mówi, co ta zgoda daje i kiedy operator nie powinien jej zaznaczać podczas logowania. */
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
