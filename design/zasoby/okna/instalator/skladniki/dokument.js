/* Składnik dokumentu osadzonego wyświetla w całości akt prawny albo instrukcję w polu przewijanym o stałej wysokości, pobierając treść z rejestru.

   Akt prawny albo instrukcja wyświetlana w całości w polu przewijanym o stałej
   wysokości. Treść NIE stoi w znaczniku — składnik bierze ją z rejestru treści
   `window.DanacoTresci`, który wypełniają wpięte pliki treści. Dokument ma więc
   jedno źródło i da się go wymienić bez dotykania okna, a okno otwiera się
   wprost z dysku, bo nic się nie dociąga.

   Pole jest ogniskowalne klawiaturą, bo samo się przewija: bez tego czytający
   klawiaturą nie ma jak do niego dojść.

   Właściwości:
     zrodlo     nazwa pozycji w rejestrze `window.DanacoTresci`
     etykieta   klucz katalogu — nazwa obszaru dla czytnika ekranu
     rosnace    true — pole bierze całą wolną wysokość kroku */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.dokument = function (w) {
  var pole = el('div', {
    klasa: 'dn-kreator-dokument dn-tekst-ciagly dn-tekst-ciagly--osadzony',
    tabindex: '0', role: 'region',
    'aria-label': tekst(w.etykieta),
    dane: { licencja: true }
  });
  var html = (window.DanacoTresci || {})[w.zrodlo];
  /* Brak treści jest usterką widoczną, nie pustym polem: dokument, którego nie
     widać, a którego akceptacji się wymaga, byłby zgodą pozorną. */
  if (html) pole.innerHTML = html;
  else pole.textContent = '⟨brak treści w rejestrze: ' + w.zrodlo + '⟩';
  return el('div', { klasa: 'dn-pole' + (w.rosnace ? ' dn-pole--rosnace' : '') }, [pole]);
};
})();
