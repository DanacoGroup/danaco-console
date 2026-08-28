/* Składnik karty opcji wyboru przedstawia jedną z równorzędnych dróg instalacji, między którymi operator rozstrzyga raz na całą instalację.

   Jedna z dróg równorzędnych, między którymi trzeba rozstrzygnąć raz na całą
   instalację. Cała karta jest celem wskazania, a stan wybrania niesie wstęga
   przy krawędzi i zaznaczona kontrolka — nie sama barwa tła.

   Plakietka jest ukryta domyślnie: pokazuje ją dopiero ekran, gdy wie, która
   opcja jest zalecana.

   Właściwości:
     nazwa      klucz katalogu — nazwa opcji
     opis       klucz katalogu — zdanie objaśniające
     wartosc    wartość kontrolki radiowej
     grupa      nazwa grupy radiowej
     plakietka  klucz katalogu — napis plakietki (opcjonalna)
     zaznaczona true | false
     fokus      true — kontrolka bierze fokus przy wejściu w krok */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.kartaOpcji = function (w) {
  var nazwa = [tekst(w.nazwa)];
  if (w.plakietka) {
    nazwa.push(' ');            /* odstęp międzywyrazowy przed plakietką */
    nazwa.push(el('span', {
      klasa: 'dn-plakietka dn-plakietka--sygnal',
      dane: { plakietka: true }, hidden: true,
      tekst: tekst(w.plakietka)
    }));
  }
  return el('label', { klasa: 'dn-wybor dn-wybor--blokowy dn-wybor--pole' }, [
    el('input', {
      klasa: 'dn-radio', type: 'radio', name: w.grupa, value: w.wartosc,
      checked: !!w.zaznaczona, dane: w.fokus ? { fokus: true } : null
    }),
    el('span', {}, [
      el('span', { klasa: 'dn-wybor-nazwa' }, nazwa),
      el('span', { klasa: 'dn-wybor-opis', tekst: tekst(w.opis) })
    ])
  ]);
};
})();
