/* ============================================================================
   SKŁADNIK — POLE ŚCIEŻKI Z PRZYCISKIEM ZMIANY
   ----------------------------------------------------------------------------
   Katalog wskazany do instalacji: etykieta, ścieżka tylko do odczytu
   i czynność zmiany. Ścieżka jest kontrolką, nie tekstem — musi dać się
   zaznaczyć i skopiować klawiaturą.

   Właściwości:
     etykieta   klucz katalogu — nazwa pola
     wartosc    klucz katalogu — ścieżka
     opis       klucz katalogu — zdanie o zawartości katalogu
     id         identyfikator kontrolki, wiązany etykietą
     zmien      klucz katalogu — napis przycisku zmiany
     naZmiane   wywołanie zwrotne przycisku
     stan       'nieczynne' | null
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.poleSciezki = function (w) {
  var przycisk = el('button', {
    klasa: 'dn-btn dn-btn--zarys', type: 'button',
    disabled: w.stan === 'nieczynne',
    tekst: tekst(w.zmien)
  });
  if (w.naZmiane) przycisk.addEventListener('click', w.naZmiane);
  return el('div', { klasa: 'dn-pole' }, [
    el('label', { klasa: 'dn-pole-etykieta', for: w.id, tekst: tekst(w.etykieta) }),
    el('span', { klasa: 'dn-pole-zestaw' }, [
      el('input', {
        klasa: 'dn-pole-kontrolka', id: w.id, type: 'text', readonly: true,
        disabled: w.stan === 'nieczynne', value: tekst(w.wartosc)
      }),
      przycisk
    ]),
    el('span', { klasa: 'dn-pole-opis', tekst: tekst(w.opis) })
  ]);
};
})();
