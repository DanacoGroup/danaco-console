/* PASEK STANU — miary
   Udział toru wchodzi żetonem `--dn-stan-udzial` z atrybutu `data-udzial`,
   więc znacznik nie nosi wstawek stylu. */
(function () {
  'use strict';

  function odswiez() {
    Array.prototype.forEach.call(document.querySelectorAll('.dn-stan-tor[data-udzial]'), function (tor) {
      var udzial = parseFloat(tor.getAttribute('data-udzial'));
      if (isNaN(udzial)) { return; }
      tor.style.setProperty('--dn-stan-udzial', Math.max(0, Math.min(100, udzial)) + '%');
      tor.setAttribute('data-poziom', udzial >= 80 ? 'wysoki' : 'zwykly');
    });
  }

  window.dnPasekStanu = { odswiez: odswiez };
  document.addEventListener('DOMContentLoaded', odswiez);
  odswiez();
})();
