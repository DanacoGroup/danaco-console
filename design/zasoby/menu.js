/* Moduł menu rozwijanego oblicza położenie panelu względem wyzwalacza oraz obsługuje nawigację klawiaturą dla wszystkich menu pakietu.
   Panel menu jest przenoszony do warstwy okna i ustawiany względem swojego
   wyzwalacza. Jedno miejsce liczy położenie wszystkich menu w pakiecie, więc
   dodanie kolejnego menu nie wymaga własnych reguł: wystarczy wyzwalacz
   `data-menu="ID"` i panel `.sta-menu-tresc` o tym `ID`.

   Kierunek wysuwania: `dol` (domyślny) albo `prawo`. Kierunek `prawo` biorą
   menu wyzwalane z pasa przy krawędzi okna (szyna nawigacji), bo w dół nie ma
   tam miejsca. Wyzwalacz albo panel może go nadać wprost atrybutem
   `data-menu-kierunek`. Wyrównanie: `data-menu-wyrownanie` albo — zgodnie ze
   znacznikiem pakietu — `data-kotwica="prawo"`.

   Otwieranie i zamykanie zostaje w `prototyp.js` (mechanika `data-menu`);
   ten plik zmienia wyłącznie położenie panelu i obsługę klawiatury.

   Wymaga: `menu.css`. */
(function () {
  'use strict';

  var ODSTEP = 8;
  var miejsce = new WeakMap();

  function wyzwalacz(panel) { return document.querySelector('[data-menu="' + panel.id + '"]'); }

  function kierunek(panel, w) {
    var jawny = panel.getAttribute('data-menu-kierunek') || (w && w.getAttribute('data-menu-kierunek'));
    if (jawny) { return jawny; }
    return w && w.closest('.dn-szyna-nawigacji') ? 'prawo' : 'dol';
  }

  function doWarstwy(panel) {
    if (!miejsce.has(panel)) {
      miejsce.set(panel, { rodzic: panel.parentNode, nastepnik: panel.nextSibling });
    }
    if (panel.parentNode !== document.body) { document.body.appendChild(panel); }
  }

  function zWarstwy(panel) {
    panel.style.position = '';
    panel.style.left = '';
    panel.style.top = '';
    panel.style.right = '';
    panel.style.bottom = '';
    var m = miejsce.get(panel);
    if (m && panel.parentNode === document.body) { m.rodzic.insertBefore(panel, m.nastepnik); }
  }

  /* Obszar, w którym menu ma się zmieścić: okno albo panel, w którym stoi
     wyzwalacz. Bez tego panel rozwinięty przy prawej krawędzi okna wchodzi na
     okno sąsiednie. */
  function obszar(w) {
    var e = w.closest('.dn-obszar-panel, .dn-rama-prawa, .dn-belka, .dn-narzedzia-pas, .dn-szyna-nawigacji');
    if (!e) { return { left: ODSTEP, top: ODSTEP, right: window.innerWidth - ODSTEP, bottom: window.innerHeight - ODSTEP }; }
    var r = e.getBoundingClientRect();
    /* W poziomie panel zostaje w regionie wyzwalacza — inaczej wchodzi na okno
       sąsiednie. W pionie obowiązuje całe okno aplikacji: rozwinięcie ma prawo
       przykryć treść pod sobą, na tym polega menu. */
    return {
      left: Math.max(ODSTEP, r.left),
      top: ODSTEP,
      right: Math.min(window.innerWidth - ODSTEP, Math.max(r.right, r.left + 360)),
      bottom: window.innerHeight - ODSTEP
    };
  }

  function poziomo(kier, wyrownanie, t, m, box) {
    if (kier === 'prawo') {
      if (t.right + ODSTEP / 2 + m.width <= box.right) { return t.right + ODSTEP / 2; }
      if (t.left - ODSTEP / 2 - m.width >= box.left) { return t.left - ODSTEP / 2 - m.width; }
      return Math.max(box.left, Math.min(t.right + ODSTEP / 2, window.innerWidth - m.width - ODSTEP));
    }
    if (kier === 'lewo') { return Math.max(box.left, t.left - ODSTEP / 2 - m.width); }
    /* Kierunek pionowy: wyrównanie do lewej albo prawej krawędzi wyzwalacza.
       Domyślnie do lewej; przy braku miejsca panel wyrównuje się do prawej,
       zamiast wychodzić poza obszar. */
    if (wyrownanie === 'prawo') { return Math.max(box.left, t.right - m.width); }
    if (t.left + m.width <= box.right) { return t.left; }
    if (t.right - m.width >= box.left) { return t.right - m.width; }
    return Math.max(box.left, Math.min(t.left, window.innerWidth - m.width - ODSTEP));
  }

  function pionowo(kier, t, m, box) {
    if (kier === 'prawo' || kier === 'lewo') {
      if (t.top + m.height <= box.bottom) { return t.top; }
      return Math.max(box.top, box.bottom - m.height);
    }
    if (kier === 'gora') { return Math.max(box.top, t.top - ODSTEP / 2 - m.height); }
    if (t.bottom + ODSTEP / 2 + m.height <= box.bottom) { return t.bottom + ODSTEP / 2; }
    if (t.top - ODSTEP / 2 - m.height >= box.top) { return t.top - ODSTEP / 2 - m.height; }
    return Math.max(box.top, box.bottom - m.height);
  }

  function ustaw(panel) {
    var w = wyzwalacz(panel);
    if (!w) { return; }
    var t = w.getBoundingClientRect();
    var box = obszar(w);
    doWarstwy(panel);
    /* Kotwice z arkuszy muszą zniknąć: razem z ustawianym tu `left`/`top`
       rozciągałyby panel na całą szerokość albo wysokość okna. */
    panel.style.position = 'fixed';
    panel.style.right = 'auto';
    panel.style.bottom = 'auto';
    panel.style.left = '0';
    panel.style.top = '0';

    var m = panel.getBoundingClientRect();
    var kier = kierunek(panel, w);
    var wyrownanie = panel.getAttribute('data-menu-wyrownanie') || w.getAttribute('data-menu-wyrownanie')
      || (panel.getAttribute('data-kotwica') === 'prawo' ? 'prawo' : '');
    panel.style.left = Math.round(poziomo(kier, wyrownanie, t, m, box)) + 'px';
    panel.style.top = Math.round(pionowo(kier, t, m, box)) + 'px';
  }

  var obserwator = new MutationObserver(function (zmiany) {
    zmiany.forEach(function (z) {
      var panel = z.target;
      if (panel.getAttribute('data-otwarte') === 'tak') { ustaw(panel); } else { zWarstwy(panel); }
    });
  });

  Array.prototype.forEach.call(document.querySelectorAll('[data-menu-tresc]'), function (panel) {
    if (!panel.id) { return; }
    obserwator.observe(panel, { attributes: true, attributeFilter: ['data-otwarte'] });
  });

  /* Klawiatura: Esc zamyka, strzałki prowadzą po pozycjach. */
  document.addEventListener('keydown', function (e) {
    var panel = document.querySelector('[data-menu-tresc][data-otwarte="tak"]');
    if (!panel) { return; }
    if (e.key === 'Escape') {
      panel.setAttribute('data-otwarte', 'nie');
      var w = wyzwalacz(panel);
      if (w) { w.setAttribute('aria-expanded', 'false'); w.focus(); }
      return;
    }
    if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp') { return; }
    var pozycje = Array.prototype.filter.call(panel.querySelectorAll('.sta-menu-poz'), function (p) {
      return !p.disabled && p.offsetParent !== null;
    });
    if (!pozycje.length) { return; }
    e.preventDefault();
    var teraz = pozycje.indexOf(document.activeElement);
    var krok = e.key === 'ArrowDown' ? 1 : -1;
    pozycje[(teraz + krok + pozycje.length) % pozycje.length].focus();
  });

  window.addEventListener('resize', function () {
    Array.prototype.forEach.call(document.querySelectorAll('[data-menu-tresc][data-otwarte="tak"]'), ustaw);
  });
})();
