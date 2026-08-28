/* Baner komunikatu przedstawia stan komunikatu — informację, ostrzeżenie, błąd albo potwierdzenie sukcesu — łącząc głowę nazywającą rzecz z treścią opisującą skutek albo zalecane działanie.
   Wygląd wnosi składnik biblioteki `.dn-alert` w wariancie ze wstęgą: barwa
   stanu obejmuje znak i głowę, a wstęga przy lewej krawędzi niesie stan
   kształtem, bo barwa jako jedyna różnica nie wystarcza (norma WCAG 1.4.1).

   Właściwości:
     rodzaj  'informacja' | 'ostrzezenie' | 'blad' | 'sukces'
     ikona   nazwa znaku z zestawu
     glowa   klucz katalogu
     tresc   klucz katalogu
     dane      dane do podstawienia w treść
     daneGlowy dane do podstawienia w głowę
     odliczanie        czas w sekundach — licznik w GŁOWIE
     odliczanieTresci  czas w sekundach — licznik w TREŚCI
     postacOdliczania  'sekundy' — sama liczba zamiast zapisu mm:ss
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

/* Głowa z licznikiem rozpada się na trzy części: to, co przed liczbą, sam
   licznik i to, co po niej. Inaczej mechanika musiałaby przepisywać całe
   zdanie co sekundę, a wtedy czytnik ekranu ogłaszałby je od nowa. */
function zLicznikiem(N, znacznik, zdanie, sekundy, postac) {
  var czesci = String(zdanie).split('{odliczanie}');
  return N.el(znacznik, {}, [
    czesci[0] || '',
    N.el('span', { dane: { odliczanie: sekundy, 'odliczanie-postac': postac || null } }),
    czesci[1] || ''
  ]);
}

function zbudujGlowe(N, w, glowa) {
  if (!w.odliczanie) return N.el('b', { tekst: glowa });
  return zLicznikiem(N, 'b', glowa, w.odliczanie, w.postacOdliczania);
}

W.skladniki.baner = function (N, w) {
  var znak = N.zeZnacznika(W.ikony[w.ikona || 'informacja']);
  znak.setAttribute('aria-hidden', 'true');
  var tresc = N.tekst(w.tresc);
  if (w.dane) tresc = N.podstaw(tresc, w.dane);
  var glowa = N.tekst(w.glowa);
  if (w.daneGlowy) glowa = N.podstaw(glowa, w.daneGlowy);
  var RODZAJE = { informacja: 'info', ostrzezenie: 'ostrzezenie', blad: 'blad', sukces: 'sukces' };
  return N.el('div', {
    klasa: 'dn-alert dn-alert--wstega dn-alert--' + (RODZAJE[w.rodzaj] || 'info'),
    id: w.id || null,
    role: w.rodzaj === 'blad' || w.rodzaj === 'ostrzezenie' ? 'alert' : 'status'
  }, [
    N.el('span', { klasa: 'dn-alert-znak', 'aria-hidden': 'true' }, [znak]),
    N.el('span', { klasa: 'dn-alert-tresc' }, [
      zbudujGlowe(N, w, glowa),
      w.odliczanieTresci
        ? zLicznikiem(N, 'span', tresc, w.odliczanieTresci, w.postacOdliczania)
        : tresc
    ])
  ]);
};
})();
