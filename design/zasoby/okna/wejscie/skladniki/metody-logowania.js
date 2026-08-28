/* Metody logowania inne niż główna tworzą wykaz trzech dróg pobocznych, z których każda niesie własny stan dostępności zamiast znikać, gdy jest wyłączona.
   Pusta lista dróg pobocznych nie mówi, że są jakieś, a lista z wygaszonymi
   pozycjami mówi, co można włączyć w ustawieniach.

   Składnik zwraca WYKAZ węzłów, nie jeden węzeł: etykieta i siatka metod są
   rodzeństwem w kolumnie panelu. Opakowanie ich w pudełko wprowadziłoby
   dodatkowy poziom, przez który odstęp kolumny liczyłby się raz zamiast dwa.

   Właściwości:
     dostepne   wykaz nazw metod aktywnych
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

var METODY = [
  { klucz: 'pin', ikona: 'klawiatura' },
  { klucz: 'klucz', ikona: 'tarcza' },
  { klucz: 'email', ikona: 'koperta' }
];

W.skladniki.metodyLogowania = function (N, w) {
  var dostepne = w && w.dostepne ? w.dostepne : ['email'];
  return [
    N.el('p', { klasa: 'au-etykieta', tekst: N.tekst('dostep.logowanie.metody.naglowek') }),
    N.el('div', { klasa: 'au-metody' }, METODY.map(function (m) {
      var znak = N.zeZnacznika(W.ikony[m.ikona]);
      znak.setAttribute('aria-hidden', 'true');
      var czynna = dostepne.indexOf(m.klucz) !== -1;
      return N.el('button', { klasa: 'dn-kafel dn-kafel--wybor', type: 'button' }, [
        znak,
        N.el('span', { tekst: N.tekst('dostep.logowanie.metody.' + m.klucz) }),
        N.el('span', {
          klasa: 'dn-kafel-stan',
          tekst: N.tekst('dostep.logowanie.metody.' + (czynna ? 'aktywna' : 'nieaktywna'))
        })
      ]);
    }))
  ];
};
})();
