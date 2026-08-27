(function () {
  'use strict';
  var D = document;
  var obszar = D.querySelector('.dn-obszar');
  if (!obszar) { return; }
  var lewy = D.querySelector('.dn-obszar-panel--boczny');
  var prawy = D.querySelector('.dn-obszar-panel--samouczek');

  function zbudujUchwyt(panel, strona, zmienna, minSzer, maxSzer) {
    if (!panel) { return; }
    var u = D.createElement('span');
    u.className = 'dn-uchwyt';
    u.setAttribute('role', 'separator');
    u.setAttribute('aria-orientation', 'vertical');
    u.setAttribute('tabindex', '0');
    u.setAttribute('aria-label', strona === 'lewa'
    ? 'Szerokość okna sesji i projektów' : 'Szerokość okna samouczka');
    u.setAttribute('aria-valuemin', String(minSzer));
    u.setAttribute('aria-valuemax', String(maxSzer));
    u.setAttribute('aria-valuenow', String(Math.round(panel.getBoundingClientRect().width)));
    if (strona === 'lewa') { panel.after(u); } else { panel.before(u); }

    function ustaw(px) {
      var w = Math.max(minSzer, Math.min(maxSzer, Math.round(px)));
      D.documentElement.style.setProperty(zmienna, w + 'px');
      u.setAttribute('aria-valuenow', String(w));
      return w;
    }

    var ciagnie = false;
    u.addEventListener('pointerdown', function (e) {
      ciagnie = true;
      u.setPointerCapture(e.pointerId);
      u.setAttribute('data-ciagniety', 'tak');
      D.body.style.cursor = 'col-resize';
      D.body.style.userSelect = 'none';
      e.preventDefault();
    });
    u.addEventListener('pointermove', function (e) {
      if (!ciagnie) { return; }
      var b = panel.getBoundingClientRect();
      ustaw(strona === 'lewa' ? e.clientX - b.left : b.right - e.clientX);
    });
    function koniec(e) {
      if (!ciagnie) { return; }
      ciagnie = false;
      try { u.releasePointerCapture(e.pointerId); } catch (err) {  }
      u.removeAttribute('data-ciagniety');
      D.body.style.cursor = '';
      D.body.style.userSelect = '';
    }
    u.addEventListener('pointerup', koniec);
    u.addEventListener('pointercancel', koniec);

    u.addEventListener('keydown', function (e) {
      var krok = e.shiftKey ? 40 : 10;
      var teraz = panel.getBoundingClientRect().width;
      if (e.key === 'ArrowLeft') { ustaw(strona === 'lewa' ? teraz - krok : teraz + krok); e.preventDefault(); }
      if (e.key === 'ArrowRight') { ustaw(strona === 'lewa' ? teraz + krok : teraz - krok); e.preventDefault(); }
    });
  }

  zbudujUchwyt(lewy, 'lewa', '--cd-szer-boczny', 220, 560);
  zbudujUchwyt(prawy, 'prawa', '--cd-szer-samouczek', 260, 620);

  if (lewy) {
    D.addEventListener('click', function (e) {
      var b = e.target.closest('[data-zwin-szyne]');
      if (!b) { return; }
      e.stopPropagation();
      var zwinieta = lewy.getAttribute('data-zwiniety') === 'tak';
      if (zwinieta) {
        lewy.removeAttribute('data-zwiniety');
        lewy.removeAttribute('data-podglad');
      } else {
        lewy.setAttribute('data-zwiniety', 'tak');
      }
      b.setAttribute('aria-pressed', zwinieta ? 'false' : 'true');
      b.setAttribute('data-etykietka', zwinieta ? 'Zwiń panel' : 'Rozwiń panel');
      b.setAttribute('aria-label', zwinieta ? 'Zwiń panel' : 'Rozwiń panel');
    }, true);

    var krawedz = D.createElement('div');
    krawedz.className = 'cd-krawedz-podgladu';
    krawedz.setAttribute('aria-hidden', 'true');
    var szyna = D.querySelector('.dn-szyna, .dn-szyna-tresc');
    krawedz.style.left = szyna ? Math.round(szyna.getBoundingClientRect().right) + 'px' : '0';
    D.body.appendChild(krawedz);

    krawedz.addEventListener('pointerenter', function () {
      if (lewy.getAttribute('data-zwiniety') === 'tak') { lewy.setAttribute('data-podglad', 'tak'); }
    });
    lewy.addEventListener('pointerleave', function () { lewy.removeAttribute('data-podglad'); });
    krawedz.addEventListener('pointerleave', function (e) {
      if (!lewy.contains(e.relatedTarget)) { lewy.removeAttribute('data-podglad'); }
    });
  }
})();
