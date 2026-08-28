/* Pasek postępu z licznikiem miary pracy w toku wymaga głowy nazywającej, czego miara dotyczy, bo dwie liczby stojące obok siebie inaczej czyta się jako jedną z rozjazdem wyglądającym na usterkę. */
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
    // Stan nieokreślony zdejmuje aria-valuenow, bo rolę ogłasza progressbar.
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
