/* ============================================================================
   EKRAN 5 — INSTALACJA
   ----------------------------------------------------------------------------
   Krok bez wyboru: jedyną czynnością operatora jest przerwanie, a na krok 6
   instalator przechodzi sam. Miara nad torem jest miarą CAŁOŚCI i tak jest
   podpisana; numer etapu i licznik etapowy stoją razem w pasku szczegółów.

   Trzy odsłony — przebieg, wycofywanie, błąd — niesie `data-odslona`.
   ============================================================================ */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst, S = K.skladniki;
K.ekrany = K.ekrany || {};

K.ekrany[5] = function () {
  return el('section', {
    klasa: 'dn-kreator-ekran',
    dane: { ekran: 5, odslona: 'przebieg' },
    'aria-label': tekst('szyna.kroki')[4]
  }, [
    S.naglowekEkranu({ nadtytul: 'krok5.nadtytul', tytul: 'krok5.przebieg.tytul', podtytul: 'krok5.przebieg.podtytul' }),
    S.pasekPostepu({ etykieta: 'krok5.postep.etykieta', opisPaska: 'krok5.postep.opisPaska', wartosc: 0 }),
    S.listaEtapow({ naglowek: 'krok5.etapyNaglowek', etapy: 'krok5.etapy', stany: 'krok5.stany',
      poczatkowe: ['gotowe', 'w-toku', 'oczekuje', 'oczekuje'] }),
    S.pasekSzczegolow({ dane: { szczegoly: true } })
  ].concat(S.blokBledu({ szczegoly: 'krok5.blad.szczegoly', rada: 'krok5.blad.rada', dane: { blok: 'blad' }, ukryty: true })));
};
})();
