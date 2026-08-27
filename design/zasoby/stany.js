/* Stan zbiorczy porządkuje kolejność stanów pozycji i wylicza stan zbiorczy pozycji nadrzędnej na podstawie stanów jej bezpośrednich potomków.
   Arkusz `stany.css` rysuje stan, ten plik go wylicza.

   Kolejność: blad › reakcja › praca › zakonczone › brak.

   Pojemnik z atrybutem `data-stan-zbiorczy` dostaje stan najwyższy ze swoich
   potomków niosących `data-stan`; wynik trafia na wskazany znacznik:
     data-stan-zbiorczy="cel"   selektor znacznika w obrębie pojemnika
   Wywołanie ręczne: `dnStany.odswiez(pojemnik)`.
   ============================================================================ */
(function () {
  'use strict';

  var KOLEJNOSC = ['blad', 'reakcja', 'praca', 'zakonczone', 'brak'];

  function najwyzszy(stany) {
    for (var i = 0; i < KOLEJNOSC.length; i += 1) {
      if (stany.indexOf(KOLEJNOSC[i]) !== -1) { return KOLEJNOSC[i]; }
    }
    return 'brak';
  }

  function znacznikCelu(pojemnik) {
    var cel = pojemnik.getAttribute('data-stan-zbiorczy');
    if (!cel) { return pojemnik.querySelector('.dn-stan-znacznik, .dn-stan-liczba'); }
    try { return pojemnik.querySelector(cel); } catch (e) { return null; }
  }

  function zbierz(pojemnik) {
    var znacznik = znacznikCelu(pojemnik);
    if (!znacznik) { return; }
    var stany = Array.prototype.filter.call(
      pojemnik.querySelectorAll('[data-stan]:not([data-stan-zbiorczy])'),
      function (e) { return e !== znacznik; }
    ).map(function (e) { return e.getAttribute('data-stan'); })
      .filter(function (s) { return s && s !== 'brak'; });
    znacznik.setAttribute('data-stan', najwyzszy(stany));
  }

  function odswiez(zakres) {
    var korzen = zakres || document;
    Array.prototype.forEach.call(korzen.querySelectorAll('[data-stan-zbiorczy]'), zbierz);
  }

  window.dnStany = { KOLEJNOSC: KOLEJNOSC, najwyzszy: najwyzszy, odswiez: odswiez };

  document.addEventListener('DOMContentLoaded', function () { odswiez(); });
})();
