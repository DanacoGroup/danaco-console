/* Wiersz pod paskiem postępu pokazuje numer etapu, licznik tego, co etap właśnie przenosi, i pozostały czas, przy czym licznik należy wyłącznie do etapu oznaczonego jako w toku. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.pasekSzczegolow = function (w) {
  return el('div', { klasa: 'dn-zestawienie', dane: w.dane || { szczegoly: true } }, [
    el('span', { dane: { 'licznik-etapu': true } }),
    el('span', { klasa: 'dn-meta', dane: { pozostalo: true } })
  ]);
};
})();
