(function () {
  'use strict';
  document.addEventListener('click', function (e) {
    var b = e.target.closest('.cd-wejdz');
    if (!b) { return; }
    var srod = (b.getAttribute('data-wejdz') || '').toLowerCase();
    if (!srod) { return; }
    window.setTimeout(function () {
      window.location.href = '../srodowiska/' + srod + '-przedsionek.html';
    }, 220);
  }, true);
})();
