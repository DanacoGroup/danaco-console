(function () {
  'use strict';

  /* Wejscie do srodowiska z karty centrum.
     Srodowisko, ktore ma wlasna karte w oknie roboczym, otwiera sie na miejscu:
     okno lewe i okno prawe zostaja, zmienia sie tylko widok w oknie roboczym.
     Srodowisko bez karty dziala jak dotad — przechodzi do wlasnego okna. */

  document.addEventListener('click', function (e) {
    var b = e.target.closest('.cd-wejdz');
    if (!b) { return; }
    var srod = b.getAttribute('data-wejdz') || '';
    if (!srod) { return; }

    var pasmo = document.querySelector('.dn-obszar-panel--glowny .dn-karty-pasmo');
    var karta = pasmo && pasmo.querySelector('.dn-karta-widoku[data-srodowisko="' + srod + '"]');
    if (karta) {
      e.preventDefault();
      karta.click();
      return;
    }

    window.setTimeout(function () {
      window.location.href = '../srodowiska/' + srod.toLowerCase() + '-przedsionek.html';
    }, 220);
  }, true);
})();
