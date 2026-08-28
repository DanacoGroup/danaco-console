/* Składnik banera informacyjnego przedstawia komunikat wymogu, ostrzeżenia, błędu albo wyniku, niezwiązany z żadną pojedynczą kontrolką ekranu.

   Komunikat, który nie należy do żadnej kontrolki: wymóg, ostrzeżenie, wynik.
   Głowa niesie komunikat, treść pod nią go uzasadnia — rozdziela je stopień
   i waga, a barwa zostaje barwą rodzaju komunikatu.

   Znak stoi w osi pionowej komunikatu, nie przy pierwszym wierszu.

   Właściwości:
     rodzaj     'ostrzezenie' | 'blad' | 'info' | 'sukces'
     ikona      nazwa z zestawu ikon
     glowa      klucz katalogu — głowa komunikatu (opcjonalna)
     tresc      klucz katalogu — uzasadnienie
     dane       atrybuty `data-*`
     ukryty     true — baner czeka na swój stan */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.baner = function (w) {
  var znak = N.zeZnacznika(K.ikony[w.ikona]);
  var tresc = [];
  if (w.glowa) tresc.push(el('b', { klasa: 'dn-alert-tytul', tekst: tekst(w.glowa) }));
  tresc.push(document.createTextNode(tekst(w.tresc)));
  return el('p', {
    klasa: 'dn-alert dn-alert--' + w.rodzaj,
    hidden: !!w.ukryty,
    dane: w.dane || null
  }, [
    el('span', { klasa: 'dn-alert-znak', 'aria-hidden': 'true' }, [znak]),
    el('span', { klasa: 'dn-alert-tresc' }, tresc)
  ]);
};
})();
