/* Ekran drugi kreatora instalacji prezentuje treść umowy licencyjnej i wymaga jej zaakceptowania przed kontynuacją instalacji.

   Zgoda na dokument, którego nie da się przeczytać, byłaby zgodą pozorną —
   dlatego pole umowy bierze całą wolną wysokość kroku. Braku zgody nie tłumaczy
   osobny komunikat: mówi o nim ta sama fraza, która w spoczynku mówi, co zrobić. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst, S = K.skladniki;
K.ekrany = K.ekrany || {};

K.ekrany[2] = function () {
  return el('section', {
    klasa: 'dn-kreator-ekran',
    dane: { ekran: 2 },
    'aria-label': tekst('szyna.kroki')[1]
  }, [
    S.naglowekEkranu({ nadtytul: 'krok2.nadtytul', tytul: 'krok2.tytul', podtytul: 'krok2.podtytul' }),
    S.dokument({ zrodlo: 'licencja', etykieta: 'krok2.dokument.etykieta', rosnace: true }),
    S.poleWyboru({
      etykieta: 'krok2.zgoda.etykieta', opis: 'krok2.zgoda.opis', id: 'zgoda-opis',
      dane: { fokus: true, 'zgoda-licencja': true }
    }),
    S.frazaNawigacyjna({ klucz: 'krok2.nawigacja', dane: { nawigacja: true } })
  ]);
};
})();
