/* ============================================================================
   SKŁADNIK — KOD POTWIERDZAJĄCY

   Sześć pól po jednym znaku, odliczanie ważności i czynność poboczna. Pola
   zaczynają PUSTE: okno pokazuje odsłonę przed wpisaniem, a nie udaje, że
   użytkownik zdążył już coś wpisać.

   Każdy zestaw rządzi się sam — mechanika wiąże pola grupami, więc kursor nie
   przeskakuje między oknami.

   Składnik zwraca WYKAZ węzłów: zestaw pól i stopka z odliczaniem są
   rodzeństwem w kolumnie panelu.

   Właściwości:
     odliczanie  napis odliczania, na przykład '09:12'
     czynnosc    'wklej' | 'ponow' — czynność w stopce zestawu. Wklejenie jest
                 wygodą, więc niesie przycisk; ponowienie wysyłki jest wyjściem
                 z sytuacji bez wyjścia, więc niesie odsyłacz — waga czynności
                 rozstrzyga o postaci kontrolki, nie odwrotnie
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.poleKodu = function (N, w) {
  var pola = [1, 2, 3, 4, 5, 6].map(function (nr) {
    return N.el('input', {
      klasa: 'au-kod-pole', type: 'text', inputmode: 'numeric', maxlength: '1',
      autocomplete: nr === 1 ? 'one-time-code' : null,
      'aria-label': N.podstaw(N.tekst('dostep.kod.znak'), { numer: nr })
    });
  });

  return [
    N.el('div', {
      klasa: 'au-kod', role: 'group', 'aria-label': N.tekst('dostep.kod.obszar')
    }, pola),
    N.el('div', { klasa: 'au-kod-stopka' }, [
      N.el('span', { klasa: 'au-odliczanie' }, [
        N.tekst('dostep.kod.odliczanie') + ' ', N.el('b', { tekst: w.odliczanie })
      ]),
      w.czynnosc === 'ponow'
        ? N.el('button', {
            klasa: 'au-link', type: 'button', tekst: N.tekst('dostep.kod.ponow'),
            dane: {
              komunikat: N.tekst('komunikaty.kodPonowiony.tresc'),
              'komunikat-tytul': N.tekst('komunikaty.kodPonowiony.tytul'),
              'komunikat-rodzaj': 'sukces'
            }
          })
        : N.el('button', {
            klasa: 'dn-btn dn-btn--duch dn-btn--sm', type: 'button',
            tekst: N.tekst('dostep.kod.wklej'), dane: { 'wklej-kod': true }
          })
    ])
  ];
};
})();
