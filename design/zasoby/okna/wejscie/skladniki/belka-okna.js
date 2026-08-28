/* ============================================================================
   SKŁADNIK — BELKA OKNA SYSTEMOWEGO

   Okno przed uwierzytelnieniem nie ma ramy aplikacji: ani szyny nawigacji, ani
   wstążki, ani paska stanu. Ma wyłącznie belkę — znak, tytuł i trzy kontrolki
   okna. Znak w belce dorównuje wielkością kontrolkom, więc niesie kropkę
   w barwie sygnału: przygaszony monochromat czytał się jak brakująca ikona.

   Właściwości:
     tytul   klucz katalogu — tytuł okna w belce
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.belkaOkna = function (N, w) {
  var znak = N.zeZnacznika(W.ikony.godlo);
  znak.setAttribute('class', 'dn-okno-wejsciowe-belka-znak');
  znak.setAttribute('aria-hidden', 'true');

  /* Zamkniecie niesie modyfikator: czerwien na wskazaniu jest jego wlasnym
     stanem, nie stanem sterowania w ogole. */
  function kontrolka(ikona, etykieta) {
    var i = N.zeZnacznika(W.ikony[ikona]);
    i.setAttribute('aria-hidden', 'true');
    var klasa = 'dn-okno-wejsciowe-belka-btn';
    if (ikona === 'zamknij') { klasa += ' dn-okno-wejsciowe-belka-btn--zamknij'; }
    return N.el('button', { klasa: klasa, type: 'button', 'aria-label': N.tekst(etykieta) }, [i]);
  }

  return N.el('div', { klasa: 'dn-okno-wejsciowe-belka' }, [
    znak,
    N.el('p', { klasa: 'dn-okno-wejsciowe-belka-tytul dn-okno-wejsciowe-belka-uchwyt', tekst: N.tekst(w.tytul) }),
    N.el('div', { klasa: 'we-belka-sterowanie' }, [
      kontrolka('zwin', 'okno.zwin'),
      kontrolka('rozwin', 'okno.rozwin'),
      kontrolka('zamknij', 'okno.zamknij')
    ])
  ]);
};
})();
