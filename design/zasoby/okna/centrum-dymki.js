(function () {
  var ODSTEP = 4;
  var MARGINES = 8;
  var ostatni = null;
  var plotno = null;

  function miara(el) {
    var s = window.getComputedStyle(el, '::after');
    var wyscielka = parseFloat(s.paddingLeft) + parseFloat(s.paddingRight)
    + parseFloat(s.borderLeftWidth) + parseFloat(s.borderRightWidth);
    var wysokosc = parseFloat(s.height);
    return {
      krój: s.fontWeight + ' ' + s.fontSize + ' ' + s.fontFamily,
      wyscielka: isFinite(wyscielka) ? wyscielka : 20,
      wysokosc: isFinite(wysokosc) && wysokosc > 0 ? wysokosc : 26,
    };
  }

  function szerokosc(m, tekst) {
    if (!plotno) { plotno = document.createElement('canvas').getContext('2d'); }
    plotno.font = m.krój;
    return plotno.measureText(tekst).width + m.wyscielka;
  }

  var lewaZapamietana = null;
  function strefaLewa() {
    if (lewaZapamietana !== null) { return lewaZapamietana; }
    var szyna = document.querySelector('.dn-szyna-nawigacji');
    var r = szyna ? szyna.getBoundingClientRect() : null;
    lewaZapamietana = (r && r.width) ? Math.max(MARGINES, r.right + MARGINES) : MARGINES;
    return lewaZapamietana;
  }

  function zdejmij(el) {
    if (!el) { return; }
    el.style.removeProperty('--cd-dymek-x');
    el.style.removeProperty('--cd-dymek-y');
    el.removeAttribute('data-dymek');
  }

  function ustaw(el) {
    var tekst = el.getAttribute('data-etykietka');
    if (!tekst) { return; }
    var r = el.getBoundingClientRect();
    if (!r.width) { return; }

    if (ostatni && ostatni !== el) { zdejmij(ostatni); }
    ostatni = el;

    var m = miara(el);
    var polowa = szerokosc(m, tekst) / 2;
    var x = r.left + r.width / 2;
    var lewaStrefa = strefaLewa();
    var lewaGranica = lewaStrefa + polowa;
    var prawaGranica = window.innerWidth - MARGINES - polowa;
    if (prawaGranica > lewaGranica) {
      x = Math.min(Math.max(x, lewaGranica), prawaGranica);
    } else {
      x = window.innerWidth / 2;
    }

    var wysokosc = m.wysokosc;
    var y = r.bottom + ODSTEP;
    if (y + wysokosc > window.innerHeight - MARGINES) {
      y = r.top - ODSTEP - wysokosc;
    }

    el.style.setProperty('--cd-dymek-x', Math.round(x) + 'px');
    el.style.setProperty('--cd-dymek-y', Math.round(y) + 'px');
    el.setAttribute('data-dymek', 'tak');
  }

  function zWezla(cel) {
    return cel && cel.closest ? cel.closest('[data-etykietka]') : null;
  }

  ['pointerover', 'focusin'].forEach(function (n) {
    document.addEventListener(n, function (e) {
      var el = zWezla(e.target);
      if (el) { ustaw(el); } else if (ostatni) { zdejmij(ostatni); ostatni = null; }
    }, true);
  });
  ['pointerout', 'focusout'].forEach(function (n) {
    document.addEventListener(n, function (e) {
      var el = zWezla(e.target);
      if (el && el === ostatni) { zdejmij(el); ostatni = null; }
    }, true);
  });

  document.addEventListener('pointermove', function (e) {
    var el = zWezla(e.target);
    if (el && el.getAttribute('data-dymek') !== 'tak') { ustaw(el); }
  }, true);
  window.addEventListener('scroll', function () {
    if (ostatni) { ustaw(ostatni); }
  }, true);
  window.addEventListener('resize', function () {
    lewaZapamietana = null;
    if (ostatni) { ustaw(ostatni); }
  });
})();
