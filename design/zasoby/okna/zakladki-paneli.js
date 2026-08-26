/* Zakładki paneli obszaru roboczego oraz zwijanie panelu samouczka.
   Karty okna głównego należą do `karty-okna.js` — ten plik ich nie dotyka:
   kilka kart może wskazywać ten sam panel, a wtedy kolejność pętli decydowałaby
   o tym, czy panel karty bieżącej zostaje odsłonięty. */
(function () {
  'use strict';
  document.addEventListener('click', function (e) {
    var z = e.target.closest ? e.target.closest('.dn-obszar-zakladka[aria-controls], .dn-karta-widoku[aria-controls]') : null;
    if (z && z.closest('.dn-obszar-panel--glowny .dn-karty-pasmo')) { return; }
    if (z) {
      var pasmo = z.closest('.dn-obszar-zakladki, .dn-karty-pasmo');
      Array.prototype.forEach.call(pasmo.querySelectorAll('.dn-obszar-zakladka, .dn-karta-widoku'), function (b) {
        var wybrana = b === z;
        b.setAttribute('aria-selected', wybrana ? 'true' : 'false');
        var panel = b.getAttribute('aria-controls') && document.getElementById(b.getAttribute('aria-controls'));
        if (panel) panel.hidden = !wybrana;
      });
      return;
    }
    var p = e.target.closest ? e.target.closest('[data-przelacz-samouczek], [data-zwin-samouczek]') : null;
    if (!p) return;
    var panel = document.getElementById('panel-samouczek');
    var przel = document.querySelector('[data-przelacz-samouczek]');
    if (!panel) return;
    panel.hidden = !panel.hidden;
    if (przel) przel.setAttribute('aria-pressed', panel.hidden ? 'false' : 'true');
  });
})();
