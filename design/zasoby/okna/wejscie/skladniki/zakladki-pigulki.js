/* ============================================================================
   SKŁADNIK — ZAKŁADKI LOGOWANIA I REJESTRACJI

   Dwie drogi równorzędne. Stan wybrany niesie PODŚWIETLENIE, nie wypełnienie:
   pigułka w pełnym błękicie czytała się jak przycisk działania, a to kontrolka
   wyboru — wskazuje, gdzie stoisz, a nie co wykonasz.

   Właściwości:
     wybrana   'logowanie' | 'rejestracja'
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

var CELE = ['logowanie', 'rejestracja'];

W.skladniki.zakladkiPigulki = function (N, w) {
  return N.el('div', {
    klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'tablist',
    'aria-label': N.tekst('dostep.zakladki').join(' albo ')
  }, N.tekst('dostep.zakladki').map(function (napis, i) {
    var cel = CELE[i];
    return N.el('button', {
      klasa: 'dn-zakladka', type: 'button', role: 'tab',
      'aria-selected': cel === w.wybrana ? 'true' : 'false',
      'aria-controls': 's-' + cel, tekst: napis,
      dane: { przelacz: cel, grupa: 'stan' }
    });
  }));
};
})();
