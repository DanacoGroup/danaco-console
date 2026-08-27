/* Pionowy wykaz kroków kreatora w kolumnie bocznej pokazuje, gdzie stoi przepływ i co go czeka, a pozycja niesie stan atrybutem data-stan, ponieważ stan jest daną, a nie odmianą składnika. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst;
K.skladniki = K.skladniki || {};

K.skladniki.nawigacjaKrokow = function (w) {
  var nazwy = tekst(w.klucz);
  var pozycje = nazwy.map(function (nazwa, i) {
    var nr = i + 1;
    var stan = nr < w.biezacy ? 'zrobiony' : (nr === w.biezacy ? 'biezacy' : 'oczekuje');
    return el('span', {
      klasa: 'dn-kreator-krok',
      dane: { krok: nr, stan: stan }
    }, [el('span', { klasa: 'dn-kreator-krok-nr', tekst: String(nr) }), nazwa]);
  });
  var nawigacja = el('nav', { klasa: 'dn-kreator-kroki', 'aria-label': tekst('szyna.etykieta') },
    [el('span', { klasa: 'dn-kreator-marka' }, [
        N.zeZnacznika(K.ikony.godloDuze.replace('<svg', '<svg class="dn-kreator-marka-znak" aria-hidden="true"')),
        el('span', { klasa: 'dn-kreator-marka-nazwa', tekst: tekst('szyna.marka') })
      ])].concat(pozycje));
  if (w.naKrok) {
    nawigacja.addEventListener('click', function (e) {
      var p = e.target.closest('.dn-kreator-krok');
      if (p) w.naKrok(Number(p.dataset.krok));
    });
  }
  return nawigacja;
};
})();
