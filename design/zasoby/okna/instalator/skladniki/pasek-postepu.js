/* ============================================================================
   SKŁADNIK — PASEK POSTĘPU Z LICZNIKIEM
   ----------------------------------------------------------------------------
   Miara pracy w toku. Głowa nazywa, CZEGO miara dotyczy — bez tego dwie
   liczby stojące obok siebie czyta się jako jedną i rozjazd między nimi wygląda
   na usterkę rachunku.

   Odsłona bez miary (`nieokreslony`) zdejmuje `aria-valuenow`: to rola
   `progressbar`, a nie animacja, ogłasza stan nieokreślony.

   Właściwości:
     etykieta   klucz katalogu — nazwa miary
     opisPaska  klucz katalogu — etykieta dostępności paska
     wartosc    0..100 — miara początkowa
     stan       'okreslony' | 'nieokreslony'
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.pasekPostepu = function (w) {
  var nieokreslony = w.stan === 'nieokreslony';
  var tor = el('div', {
    klasa: 'dn-postep-tor' + (nieokreslony ? ' dn-postep-tor--nieokreslony' : ''),
    role: 'progressbar',
    'aria-label': tekst(w.opisPaska),
    'aria-valuenow': nieokreslony ? null : String(w.wartosc),
    'aria-valuemin': '0', 'aria-valuemax': '100',
    dane: { 'postep-tor': true }
  }, [el('div', { klasa: 'dn-postep-wartosc' })]);

  return el('div', { klasa: 'dn-postep dn-postep--z-glowa', dane: { blok: 'postep' } }, [
    el('div', { klasa: 'dn-postep-glowa' }, [
      el('span', { dane: { 'etap-nazwa': true }, tekst: tekst(w.etykieta) }),
      el('span', { klasa: 'dn-meta', dane: { 'etap-miara': true }, tekst: w.wartosc + '%' })
    ]),
    tor
  ]);
};
})();
