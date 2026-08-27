/* Belka okna systemowego zastępuje w oknie przed uwierzytelnieniem całą ramę aplikacji, niosąc wyłącznie znak, tytuł i trzy kontrolki okna, bez szyny nawigacji, wstążki i paska stanu.
   Znak w belce dorównuje wielkością kontrolkom, więc niesie kropkę w barwie
   sygnału — przygaszony monochromat czytał się jak brakująca ikona.

   Właściwości:
     tytul   klucz katalogu — tytuł okna w belce
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.belkaOkna = function (N, w) {
  var znak = N.zeZnacznika(W.ikony.godlo);
  znak.setAttribute('class', 'we-belka-znak');
  znak.setAttribute('aria-hidden', 'true');

  function kontrolka(ikona, etykieta) {
    var i = N.zeZnacznika(W.ikony[ikona]);
    i.setAttribute('aria-hidden', 'true');
    return N.el('button', { klasa: 'we-belka-btn', type: 'button', 'aria-label': N.tekst(etykieta) }, [i]);
  }

  return N.el('div', { klasa: 'we-belka' }, [
    znak,
    N.el('p', { klasa: 'we-belka-tytul', tekst: N.tekst(w.tytul) }),
    N.el('div', { klasa: 'we-belka-sterowanie' }, [
      kontrolka('zwin', 'okno.zwin'),
      kontrolka('rozwin', 'okno.rozwin'),
      kontrolka('zamknij', 'okno.zamknij')
    ])
  ]);
};
})();
