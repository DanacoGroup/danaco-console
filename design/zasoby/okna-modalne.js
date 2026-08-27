/* Plik nadaje oknu nakładkowemu przeciąganie, zmianę rozmiaru, zwinięcie do doku i rozłożenie na cały obszar przeglądarki. ============================================================================
   DANACO CONSOLE — OKNA NAKŁADKOWE (zachowanie)
   ----------------------------------------------------------------------------
   Nadaje każdemu `<dialog class="dn-modal">` zachowanie okna: przeciąganie za
   nagłówek, zmianę rozmiaru uchwytami krawędzi i narożników, zwinięcie do
   nagłówka zadokowanego przy dolnej krawędzi oraz rozłożenie na cały obszar
   przeglądarki. Powód: nakładka bywa otwarta długo i zasłania dane, które
   operator musi w tym samym czasie odczytać.

   Uzbrojenie jest leniwe — dialog dostaje oprzyrządowanie dopiero przy
   otwarciu, bo dopiero wtedy ma zmierzalny prostokąt, z którego wyliczamy
   pozycję bezwzględną.

   Stan mieszka w atrybutach data-* na elemencie dialogu, więc jest widoczny
   w drzewie dokumentu i możliwy do odczytania przez inne skrypty prototypu.

   Style: okna-modalne.css. Bez zależności zewnętrznych.
   ============================================================================ */
(function () {
  'use strict';

  var D = document;

  /* ── 1 · STAŁE GEOMETRII ──────────────────────────────────────────────── */

  var MIN_SZER = 420;      /* poniżej tego nagłówek modala łamie się na dwa wiersze */
  var MIN_WYS = 260;
  var MARGINES_MAX = 16;   /* odstęp okna zmaksymalizowanego od krawędzi */
  var WIDOCZNE = 80;       /* tyle nagłówka zawsze zostaje w polu widzenia */
  var KROK = 16;           /* skok klawiaturowy — cztery jednostki siatki 4 px */
  var DOK_ODSTEP = 8;      /* przerwa między zadokowanymi nagłówkami */
  var DOK_SZER = 260;      /* szerokość okna zwiniętego do nagłówka */

  var STRONY = ['n', 's', 'e', 'w', 'ne', 'nw', 'se', 'sw'];

  function wszystkie(s, k) { return Array.prototype.slice.call((k || D).querySelectorAll(s)); }
  function ogranicz(v, min, max) { return v < min ? min : (v > max ? max : v); }

  /* ── 2 · IKONY ────────────────────────────────────────────────────────── */

  /* Rysowane w kodzie, bo prototypy nie ładują zestawu ikon dla samej ramy
     okna. Kreska, kwadrat i kwadrat w kwadracie — trzy stany geometrii okna. */
  function svg(sciezki) {
    return '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" ' +
      'stroke="currentColor" stroke-width="1.75" stroke-linecap="round" ' +
      'stroke-linejoin="round" aria-hidden="true" focusable="false">' + sciezki + '</svg>';
  }
  var IKONA_MIN = svg('<path d="M5 12h14"/>');
  var IKONA_MAX = svg('<rect width="14" height="14" x="2" y="8" rx="2" ry="2"/><path d="M20 16c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2h-10c-1.1 0-2 .9-2 2"/>');
  var IKONA_PRZYWROC = svg('<path d="m14 10 7-7"/><path d="M20 10h-6V4"/><path d="m3 21 7-7"/><path d="M4 14h6v6"/>');

  /* ── 3 · ODCZYT I ZAPIS STANU ─────────────────────────────────────────── */

  function stan(dlg) { return dlg.getAttribute('data-okno-stan') || 'normalny'; }

  /* Przejście na współrzędne bezwzględne. Dialog modalny domyślnie centruje
     się marginesem auto; gdybyśmy nadali left/top bez przepisania bieżącego
     prostokąta, okno skoczyłoby przy pierwszym chwycie. */
  function uwolnij(dlg) {
    if (dlg.getAttribute('data-okno-pozycja') === 'wolna') return;
    var r = dlg.getBoundingClientRect();
    dlg.setAttribute('data-okno-pozycja', 'wolna');
    ustaw(dlg, r.left, r.top, r.width, r.height);
  }

  function ustaw(dlg, x, y, w, h) {
    dlg.style.left = Math.round(x) + 'px';
    dlg.style.top = Math.round(y) + 'px';
    if (w != null) dlg.style.width = Math.round(w) + 'px';
    if (h != null) dlg.style.height = Math.round(h) + 'px';
    dlg.setAttribute('data-okno-x', Math.round(x));
    dlg.setAttribute('data-okno-y', Math.round(y));
    if (w != null) dlg.setAttribute('data-okno-szer', Math.round(w));
    if (h != null) dlg.setAttribute('data-okno-wys', Math.round(h));
  }

  /* Zapamiętanie prostokąta sprzed skoku do stanu specjalnego — bez tego
     przywrócenie musiałoby zgadywać rozmiar wyjściowy. */
  function zapamietaj(dlg) {
    var r = dlg.getBoundingClientRect();
    dlg.setAttribute('data-okno-zapis', [r.left, r.top, r.width, r.height].map(Math.round).join(' '));
  }
  function przywroc(dlg) {
    var z = (dlg.getAttribute('data-okno-zapis') || '').split(' ');
    if (z.length !== 4) return false;
    ustaw(dlg, parseFloat(z[0]), parseFloat(z[1]), parseFloat(z[2]), parseFloat(z[3]));
    return true;
  }

  /* Trzymanie nagłówka w polu widzenia. Okno wolno wysunąć poza ekran, ale
     nigdy tak, by zniknął pas, za który można je przyciągnąć z powrotem. */
  function utrzymaj(dlg) {
    if (dlg.getAttribute('data-okno-pozycja') !== 'wolna') return;
    var r = dlg.getBoundingClientRect();
    var x = ogranicz(r.left, WIDOCZNE - r.width, window.innerWidth - WIDOCZNE);
    var y = ogranicz(r.top, 0, window.innerHeight - Math.min(WIDOCZNE, r.height));
    ustaw(dlg, x, y, null, null);
  }

  /* ── 4 · DOKOWANIE OKIEN ZWINIĘTYCH ───────────────────────────────────── */

  /* Zwinięte nagłówki układają się w rząd przy dolnej krawędzi, od prawej.
     Miejsce liczone jest z listy okien już zwiniętych, żeby drugie okno nie
     wylądowało na pierwszym. */
  function przelicz_dok() {
    var zwiniete = wszystkie('.dn-modal[data-okno-stan="zminimalizowane"]');
    var y = window.innerHeight - MARGINES_MAX;
    zwiniete.forEach(function (dlg, i) {
      var wys = dlg.getBoundingClientRect().height;
      var x = window.innerWidth - MARGINES_MAX - (i + 1) * DOK_SZER - i * DOK_ODSTEP;
      ustaw(dlg, Math.max(MARGINES_MAX, x), y - wys, DOK_SZER, null);
      dlg.style.height = 'auto';
    });
  }

  /* ── 5 · PRZEŁĄCZANIE STANÓW ──────────────────────────────────────────── */

  function odswiez_przyciski(dlg) {
    var s = stan(dlg);
    var bMin = dlg.querySelector('.dn-modal-btn--min');
    var bMax = dlg.querySelector('.dn-modal-btn--max');
    if (bMin) {
      bMin.setAttribute('aria-pressed', s === 'zminimalizowane' ? 'true' : 'false');
      bMin.setAttribute('aria-label', s === 'zminimalizowane' ? 'Przywróć okno z paska' : 'Zwiń okno do nagłówka');
    }
    if (bMax) {
      bMax.setAttribute('aria-pressed', s === 'zmaksymalizowane' ? 'true' : 'false');
      bMax.setAttribute('aria-label', s === 'zmaksymalizowane' ? 'Przywróć rozmiar okna' : 'Rozłóż okno na cały obszar');
      bMax.innerHTML = s === 'zmaksymalizowane' ? IKONA_PRZYWROC : IKONA_MAX;
    }
  }

  function do_normalnego(dlg) {
    dlg.setAttribute('data-okno-stan', 'normalny');
    dlg.style.height = '';
    if (!przywroc(dlg)) utrzymaj(dlg);
    odswiez_przyciski(dlg);
    przelicz_dok();
  }

  function minimalizuj(dlg) {
    if (stan(dlg) === 'zminimalizowane') { do_normalnego(dlg); return; }
    uwolnij(dlg);
    if (stan(dlg) === 'normalny') zapamietaj(dlg);
    dlg.setAttribute('data-okno-stan', 'zminimalizowane');
    odswiez_przyciski(dlg);
    przelicz_dok();
  }

  function maksymalizuj(dlg) {
    if (stan(dlg) === 'zmaksymalizowane') { do_normalnego(dlg); return; }
    uwolnij(dlg);
    if (stan(dlg) === 'normalny') zapamietaj(dlg);
    dlg.setAttribute('data-okno-stan', 'zmaksymalizowane');
    dlg.style.height = '';
    ustaw(dlg, MARGINES_MAX, MARGINES_MAX,
      window.innerWidth - 2 * MARGINES_MAX, window.innerHeight - 2 * MARGINES_MAX);
    odswiez_przyciski(dlg);
    przelicz_dok();
  }

  /* ── 6 · UZBROJENIE DIALOGU ───────────────────────────────────────────── */

  function przycisk(klasa, ikona, etykieta) {
    var b = D.createElement('button');
    b.type = 'button';
    b.className = 'dn-modal-btn ' + klasa;
    b.innerHTML = ikona;
    b.setAttribute('aria-label', etykieta);
    b.setAttribute('aria-pressed', 'false');
    return b;
  }

  function uzbroj(dlg) {
    if (dlg.getAttribute('data-okno') === '1') return;
    dlg.setAttribute('data-okno', '1');
    dlg.setAttribute('data-okno-stan', 'normalny');

    var nag = dlg.querySelector('.dn-modal-naglowek');
    if (nag) {
      /* Nagłówek wchodzi w kolejność fokusu, bo obsługuje strzałki. */
      if (!nag.hasAttribute('tabindex')) nag.setAttribute('tabindex', '0');
      nag.setAttribute('aria-label', 'Nagłówek okna — przeciąganie i zmiana rozmiaru strzałkami');

      var zam = nag.querySelector('.dn-modal-zamknij');
      var bMin = przycisk('dn-modal-btn--min', IKONA_MIN, 'Zwiń okno do nagłówka');
      var bMax = przycisk('dn-modal-btn--max', IKONA_MAX, 'Rozłóż okno na cały obszar');
      if (zam) { nag.insertBefore(bMin, zam); nag.insertBefore(bMax, zam); }
      else { nag.appendChild(bMin); nag.appendChild(bMax); }
    }

    STRONY.forEach(function (s) {
      var u = D.createElement('div');
      u.className = 'dn-modal-uchwyt dn-modal-uchwyt--' + s;
      u.setAttribute('data-okno-uchwyt', s);
      u.setAttribute('aria-hidden', 'true');
      dlg.appendChild(u);
    });

    odswiez_przyciski(dlg);
  }

  /* Dialog jest uzbrajany przy otwarciu — showModal() jest wywoływane z rozmaitych
     miejsc prototypu, więc podpinamy się pod samą metodę zamiast tropić
     wszystkie wywołania. */
  var showModal = window.HTMLDialogElement && window.HTMLDialogElement.prototype.showModal;
  if (showModal) {
    window.HTMLDialogElement.prototype.showModal = function () {
      var wynik = showModal.apply(this, arguments);
      if (this.classList.contains('dn-modal')) uzbroj(this);
      return wynik;
    };
  }
  D.addEventListener('DOMContentLoaded', function () {
    wszystkie('.dn-modal[open]').forEach(uzbroj);
  });

  /* ── 6a · ZAMKNIĘCIE OKNA — obsługa własna, niezależna od reszty ──────── */

  /* Przycisk zamknięcia obsługuje także prototyp.js. Powiela się tu obsługę
     świadomie: gdy którykolwiek wcześniejszy skrypt prototypu przerwie się
     błędem, jego nasłuch nie powstanie i okna przestaną się zamykać, choć menu
     (obsługiwane w innym pliku) nadal będą działać — objaw myli, bo wygląda na
     usterkę samego okna. Dwa niezależne nasłuchy znoszą tę zależność.
     Wywołanie jest bezpieczne przy powtórzeniu: close() na zamkniętym oknie
     nie robi nic. */
  D.addEventListener('click', function (e) {
    var b = e.target.closest
      ? e.target.closest('[data-zamknij-modal], .dn-modal-zamknij')
      : null;
    if (!b) return;
    var d = b.closest('dialog');
    if (!d) return;
    if (typeof d.close === 'function') d.close();
    else d.removeAttribute('open');
  });

  /* Klawisz Escape zamyka okno także wtedy, gdy nie jest ono modalne (open bez
     showModal): natywne zachowanie dialogu obejmuje wyłącznie tryb modalny. */
  D.addEventListener('keydown', function (e) {
    if (e.key !== 'Escape') return;
    var d = D.querySelector('.dn-modal[open]');
    if (d && !d.matches(':modal') && typeof d.close === 'function') d.close();
  });

  /* ── 7 · PRZECIĄGANIE I ZMIANA ROZMIARU ───────────────────────────────── */

  var gest = null;   /* { dlg, tryb, strona, x0, y0, r } */

  D.addEventListener('pointerdown', function (e) {
    if (e.button !== 0) return;
    var dlg = e.target.closest ? e.target.closest('.dn-modal[data-okno="1"]') : null;
    if (!dlg) return;

    var uchwyt = e.target.closest('[data-okno-uchwyt]');
    var nag = e.target.closest('.dn-modal-naglowek');
    /* Przyciski i pola w nagłówku zachowują własne działanie — uchwytem jest
       tylko pusta przestrzeń belki. */
    if (nag && e.target.closest('button, a, input, select, textarea')) return;
    if (!uchwyt && !nag) return;
    if (uchwyt && stan(dlg) !== 'normalny') return;

    uwolnij(dlg);
    var r = dlg.getBoundingClientRect();
    gest = {
      dlg: dlg,
      tryb: uchwyt ? 'rozmiar' : 'ruch',
      strona: uchwyt ? uchwyt.getAttribute('data-okno-uchwyt') : '',
      x0: e.clientX, y0: e.clientY,
      r: { x: r.left, y: r.top, w: r.width, h: r.height },
      ruszony: false
    };
    dlg.setAttribute('data-okno-chwyt', '1');
    if (e.target.setPointerCapture) e.target.setPointerCapture(e.pointerId);
  });

  D.addEventListener('pointermove', function (e) {
    if (!gest) return;
    var g = gest, r = g.r;
    var dx = e.clientX - g.x0, dy = e.clientY - g.y0;
    if (Math.abs(dx) > 2 || Math.abs(dy) > 2) g.ruszony = true;

    if (g.tryb === 'ruch') {
      /* Okno zwinięte po chwycie wraca do rozmiaru — przeciąganie doku nie
         miałoby sensu, skoro dok sam układa nagłówki w rząd. */
      if (stan(g.dlg) !== 'normalny' && g.ruszony) {
        g.dlg.setAttribute('data-okno-stan', 'normalny');
        g.dlg.style.height = '';
        odswiez_przyciski(g.dlg);
        przelicz_dok();
      }
      var x = ogranicz(r.x + dx, WIDOCZNE - r.w, window.innerWidth - WIDOCZNE);
      var y = ogranicz(r.y + dy, 0, window.innerHeight - Math.min(WIDOCZNE, r.h));
      ustaw(g.dlg, x, y, null, null);
      return;
    }

    var x1 = r.x, y1 = r.y, w1 = r.w, h1 = r.h, s = g.strona;
    if (s.indexOf('e') > -1) w1 = r.w + dx;
    if (s.indexOf('s') > -1) h1 = r.h + dy;
    if (s.indexOf('w') > -1) { w1 = r.w - dx; x1 = r.x + dx; }
    if (s.indexOf('n') > -1) { h1 = r.h - dy; y1 = r.y + dy; }
    /* Krawędź zachodnia i północna ciągną róg przeciwny w miejscu: przy
       dojściu do minimum korygujemy pozycję, inaczej okno by się przesuwało. */
    w1 = ogranicz(w1, MIN_SZER, window.innerWidth);
    h1 = ogranicz(h1, MIN_WYS, window.innerHeight);
    if (s.indexOf('w') > -1) x1 = r.x + r.w - w1;
    if (s.indexOf('n') > -1) y1 = r.y + r.h - h1;
    ustaw(g.dlg, x1, y1, w1, h1);
  });

  function koniec_gestu() {
    if (!gest) return;
    gest.dlg.removeAttribute('data-okno-chwyt');
    if (gest.tryb === 'ruch') utrzymaj(gest.dlg);
    gest = null;
  }
  D.addEventListener('pointerup', koniec_gestu);
  D.addEventListener('pointercancel', koniec_gestu);

  /* ── 8 · PRZYCISKI I DWUKLIK NAGŁÓWKA ─────────────────────────────────── */

  D.addEventListener('click', function (e) {
    var b = e.target.closest ? e.target.closest('.dn-modal-btn') : null;
    if (!b) return;
    var dlg = b.closest('.dn-modal[data-okno="1"]');
    if (!dlg) return;
    if (b.classList.contains('dn-modal-btn--min')) minimalizuj(dlg);
    else if (b.classList.contains('dn-modal-btn--max')) maksymalizuj(dlg);
  });

  /* Naciśnięcie zwiniętego nagłówka przywraca okno — to najkrótsza droga
     powrotu z doku i odpowiednik kliknięcia pozycji na pasku zadań. */
  D.addEventListener('click', function (e) {
    var nag = e.target.closest ? e.target.closest('.dn-modal-naglowek') : null;
    if (!nag || e.target.closest('button, a, input, select, textarea')) return;
    var dlg = nag.closest('.dn-modal[data-okno="1"]');
    if (dlg && stan(dlg) === 'zminimalizowane') do_normalnego(dlg);
  });

  D.addEventListener('dblclick', function (e) {
    var nag = e.target.closest ? e.target.closest('.dn-modal-naglowek') : null;
    if (!nag || e.target.closest('button, a, input, select, textarea')) return;
    var dlg = nag.closest('.dn-modal[data-okno="1"]');
    if (dlg && stan(dlg) !== 'zminimalizowane') maksymalizuj(dlg);
  });

  /* ── 9 · KLAWIATURA ───────────────────────────────────────────────────── */

  var OSIE = { ArrowLeft: [-1, 0], ArrowRight: [1, 0], ArrowUp: [0, -1], ArrowDown: [0, 1] };

  D.addEventListener('keydown', function (e) {
    var os = OSIE[e.key];
    if (!os) return;
    var nag = e.target.closest ? e.target.closest('.dn-modal-naglowek') : null;
    if (!nag || e.target !== nag) return;
    var dlg = nag.closest('.dn-modal[data-okno="1"]');
    if (!dlg) return;

    e.preventDefault();
    uwolnij(dlg);
    if (stan(dlg) !== 'normalny') do_normalnego(dlg);
    var r = dlg.getBoundingClientRect();

    if (e.shiftKey) {
      var w = ogranicz(r.width + os[0] * KROK, MIN_SZER, window.innerWidth);
      var h = ogranicz(r.height + os[1] * KROK, MIN_WYS, window.innerHeight);
      ustaw(dlg, r.left, r.top, w, h);
    } else {
      ustaw(dlg,
        ogranicz(r.left + os[0] * KROK, WIDOCZNE - r.width, window.innerWidth - WIDOCZNE),
        ogranicz(r.top + os[1] * KROK, 0, window.innerHeight - Math.min(WIDOCZNE, r.height)),
        null, null);
    }
  });

  /* ── 10 · ZMIANA ROZMIARU PRZEGLĄDARKI ────────────────────────────────── */

  /* Po zwężeniu okna przeglądarki nakładka mogłaby zostać poza ekranem —
     przy każdej zmianie przywracamy warunek widoczności nagłówka. */
  window.addEventListener('resize', function () {
    wszystkie('.dn-modal[data-okno="1"]').forEach(function (dlg) {
      var s = stan(dlg);
      if (s === 'zmaksymalizowane') {
        ustaw(dlg, MARGINES_MAX, MARGINES_MAX,
          window.innerWidth - 2 * MARGINES_MAX, window.innerHeight - 2 * MARGINES_MAX);
      } else if (s === 'normalny') {
        utrzymaj(dlg);
      }
    });
    przelicz_dok();
  });

  /* ── 11 · ZAMKNIĘCIE ──────────────────────────────────────────────────── */

  /* Zamknięte okno znika z doku, więc pozostałe zwinięte trzeba przesunąć. */
  D.addEventListener('close', function (e) {
    var dlg = e.target;
    if (!dlg.classList || !dlg.classList.contains('dn-modal')) return;
    if (stan(dlg) === 'zminimalizowane') {
      dlg.setAttribute('data-okno-stan', 'normalny');
      dlg.style.height = '';
      przywroc(dlg);
      odswiez_przyciski(dlg);
    }
    przelicz_dok();
  }, true);
})();
