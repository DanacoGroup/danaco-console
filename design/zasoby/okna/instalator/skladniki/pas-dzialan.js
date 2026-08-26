/* ============================================================================
   SKŁADNIK — STOPKA Z PRZYCISKAMI
   ----------------------------------------------------------------------------
   Pas czynności biegnący pod obiema kolumnami okna. Skład pasa wyprowadza się
   z kroku i odsłony, nie ustawia po kawałku — każdy krok ma własny komplet
   czynności i to on rozstrzyga, których nie ma.

   Czynność niedostępna z powodu niespełnionego warunku nosi `aria-disabled`,
   nie `disabled`: musi przyjąć kliknięcie, żeby móc odpowiedzieć, czego brakuje.

   Właściwości:
     poboczna   { klucz, dane } — czynność po lewej stronie pary
     glowna     { klucz, dane } — czynność domyślna
     dodatkowa  { klucz, dane } — czynność trzecia (opcjonalna, np. kopiowanie)
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

function przycisk(opis, klasa) {
  if (!opis) return null;
  return el('button', {
    klasa: klasa, type: 'button',
    hidden: !!opis.ukryty,
    dane: opis.dane || null,
    tekst: tekst(opis.klucz)
  });
}

K.skladniki.pasDzialan = function (w) {
  return el('div', { klasa: 'dn-kreator-pas' }, [
    el('div', { klasa: 'dn-kreator-pas-tresc' }, [
      el('span', { klasa: 'dn-kreator-pas-odstep', 'aria-hidden': 'true' }),
      przycisk(w.dodatkowa, 'dn-btn dn-btn--zarys'),
      przycisk(w.poboczna, 'dn-btn dn-btn--zarys'),
      przycisk(w.glowna, 'dn-btn dn-btn--atrament')
    ])
  ]);
};
})();
