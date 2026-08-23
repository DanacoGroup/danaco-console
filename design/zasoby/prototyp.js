/* ============================================================================
   DANACO CONSOLE — WARSTWA PROTOTYPU (zachowanie) · v2.0
   ----------------------------------------------------------------------------
   Deklaratywna warstwa interakcji wspólna wszystkim klikalnym prototypom okien.
   Sterowanie wyłącznie atrybutami data-* — bez pisania JS w każdym pliku.

   ZASADA ZERO BLOKAD: żadna kontrolka nie jest wyłączana. Zamiast blokady —
   komunikat (toast) albo opis obok.

   Wymaga wcześniejszego wczytania wspolne.js (przełącznik motywu).
   ============================================================================ */

(function () {
  'use strict';

  var D = document;

  /* ── narzędzia ────────────────────────────────────────────────────────── */

  function wszystkie(sel, korzen) { return Array.prototype.slice.call((korzen || D).querySelectorAll(sel)); }
  function pierwszy(sel, korzen) { return (korzen || D).querySelector(sel); }

  /* ── 1 · POWIADOMIENIA (toast) ────────────────────────────────────────── */

  function pojemnikToastow() {
    var p = pierwszy('.dn-toasty');
    if (!p) {
      p = D.createElement('div');
      p.className = 'dn-toasty';
      p.setAttribute('role', 'status');
      p.setAttribute('aria-live', 'polite');
      D.body.appendChild(p);
    }
    return p;
  }

  var IKONY_TOASTU = {
    sukces: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m20 6-11 11-5-5"/></svg>',
    ostrzezenie: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m10.3 3.7-8 14A2 2 0 0 0 4 20.6h16a2 2 0 0 0 1.7-2.9l-8-14a2 2 0 0 0-3.4 0Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>',
    blad: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>',
    informacja: '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>'
  };

  function toast(tytul, tresc, rodzaj, czas) {
    rodzaj = rodzaj || 'informacja';
    var t = D.createElement('div');
    t.className = 'dn-toast dn-toast--' + rodzaj;
    t.innerHTML =
      '<span class="dn-toast-ikona" aria-hidden="true">' + (IKONY_TOASTU[rodzaj] || IKONY_TOASTU.informacja) + '</span>' +
      '<div><div class="dn-toast-tytul"></div>' + (tresc ? '<div class="dn-toast-tresc"></div>' : '') + '</div>';
    pierwszy('.dn-toast-tytul', t).textContent = tytul;
    if (tresc) pierwszy('.dn-toast-tresc', t).textContent = tresc;
    pojemnikToastow().appendChild(t);
    setTimeout(function () {
      t.style.transition = 'opacity .22s, transform .22s';
      t.style.opacity = '0';
      t.style.transform = 'translateY(-4px)';
      setTimeout(function () { if (t.parentNode) t.parentNode.removeChild(t); }, 240);
    }, czas || 3200);
    return t;
  }
  window.dnToast = toast;

  /* ── 2 · PRZEŁĄCZANIE WIDOKÓW ─────────────────────────────────────────── */
  /* <button data-przelacz="nazwa-widoku" data-grupa="g1">  →  [data-widok="nazwa-widoku"][data-grupa-widoku="g1"] */

  function przelaczWidok(nazwa, grupa) {
    var selGrupa = grupa ? '[data-grupa-widoku="' + grupa + '"]' : '[data-widok]:not([data-grupa-widoku])';
    wszystkie(selGrupa).forEach(function (w) {
      var aktywny = w.getAttribute('data-widok') === nazwa;
      w.setAttribute('data-widok-aktywny', aktywny ? 'tak' : 'nie');
      if (aktywny) { w.classList.remove('pt-wejscie'); void w.offsetWidth; w.classList.add('pt-wejscie'); }
    });
    var selPrzelacznik = grupa ? '[data-przelacz][data-grupa="' + grupa + '"]' : '[data-przelacz]:not([data-grupa])';
    wszystkie(selPrzelacznik).forEach(function (p) {
      var aktywny = p.getAttribute('data-przelacz') === nazwa;
      p.setAttribute('aria-selected', aktywny ? 'true' : 'false');
      if (p.hasAttribute('aria-current') || p.classList.contains('dn-boczna-pozycja')) {
        if (aktywny) p.setAttribute('aria-current', 'page'); else p.removeAttribute('aria-current');
      }
      p.classList.toggle('dn-zakladka--wybrana', aktywny && p.classList.contains('dn-zakladka'));
    });
  }
  window.dnPrzelaczWidok = przelaczWidok;

  /* ── 3 · ZAKŁADKI (rola tablist) ──────────────────────────────────────── */

  function obsluzZakladki(pojemnik) {
    var zakladki = wszystkie('[role="tab"]', pojemnik);
    zakladki.forEach(function (z, i) {
      z.setAttribute('tabindex', z.getAttribute('aria-selected') === 'true' ? '0' : '-1');
      z.addEventListener('keydown', function (e) {
        var kierunek = e.key === 'ArrowRight' ? 1 : e.key === 'ArrowLeft' ? -1 : e.key === 'Home' ? -999 : e.key === 'End' ? 999 : 0;
        if (!kierunek) return;
        e.preventDefault();
        var cel = kierunek === -999 ? 0 : kierunek === 999 ? zakladki.length - 1
          : (i + kierunek + zakladki.length) % zakladki.length;
        zakladki[cel].click();
        zakladki[cel].focus();
      });
    });
  }

  /* ── 4 · WYBÓR POZYCJI NA LIŚCIE ──────────────────────────────────────── */
  /* [data-lista-wyboru] > [data-pozycja]  — wybór jednokrotny */

  function obsluzListeWyboru(lista) {
    lista.addEventListener('click', function (e) {
      var poz = e.target.closest('[data-pozycja]');
      if (!poz || !lista.contains(poz)) return;
      wszystkie('[data-pozycja]', lista).forEach(function (p) { p.setAttribute('aria-selected', p === poz ? 'true' : 'false'); });
      var cel = poz.getAttribute('data-pokaz');
      if (cel) przelaczWidok(cel, lista.getAttribute('data-grupa') || null);
      var zdarz = new CustomEvent('dn:wybor', { bubbles: true, detail: { wartosc: poz.getAttribute('data-pozycja') } });
      poz.dispatchEvent(zdarz);
    });
  }

  /* ── 5 · MODALE (dialog) ──────────────────────────────────────────────── */

  D.addEventListener('click', function (e) {
    var otw = e.target.closest('[data-otworz-modal]');
    if (otw) {
      var m = D.getElementById(otw.getAttribute('data-otworz-modal'));
      if (m && m.showModal) { m.showModal(); }
      return;
    }
    var zam = e.target.closest('[data-zamknij-modal]');
    if (zam) {
      var d = zam.closest('dialog');
      if (d && d.close) d.close();
    }
  });
  /* zamknięcie kliknięciem w podłoże */
  D.addEventListener('click', function (e) {
    if (e.target.tagName === 'DIALOG' && e.target.hasAttribute('data-zamknij-podlozem')) e.target.close();
  });

  /* ── 6 · MENU PODRĘCZNE ───────────────────────────────────────────────── */

  D.addEventListener('click', function (e) {
    var przycisk = e.target.closest('[data-menu]');
    var otwarte = wszystkie('[data-menu-tresc][data-otwarte="tak"]');
    otwarte.forEach(function (m) {
      if (!przycisk || m.id !== przycisk.getAttribute('data-menu')) m.setAttribute('data-otwarte', 'nie');
    });
    if (!przycisk) return;
    var m2 = D.getElementById(przycisk.getAttribute('data-menu'));
    if (!m2) return;
    var teraz = m2.getAttribute('data-otwarte') === 'tak';
    m2.setAttribute('data-otwarte', teraz ? 'nie' : 'tak');
    przycisk.setAttribute('aria-expanded', teraz ? 'false' : 'true');
  });

  /* ── 7 · KOPIOWANIE DO SCHOWKA ────────────────────────────────────────── */

  D.addEventListener('click', function (e) {
    var el = e.target.closest('[data-kopiuj]');
    if (!el) return;
    var tekst = el.getAttribute('data-kopiuj');
    if (tekst === 'self') tekst = el.textContent.trim();
    function udalo() { toast('Skopiowano', tekst.length > 60 ? tekst.slice(0, 57) + '…' : tekst, 'sukces', 2000); }
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(tekst).then(udalo, zapasoweKopiowanie);
    } else { zapasoweKopiowanie(); }
    function zapasoweKopiowanie() {
      var ta = D.createElement('textarea');
      ta.value = tekst; ta.style.position = 'fixed'; ta.style.opacity = '0';
      D.body.appendChild(ta); ta.select();
      try { D.execCommand('copy'); udalo(); } catch (err) { toast('Nie udało się skopiować', null, 'ostrzezenie'); }
      D.body.removeChild(ta);
    }
  });

  /* ── 8 · ZERO BLOKAD — komunikat zamiast bramy ────────────────────────── */
  /* <button data-komunikat="treść" data-komunikat-rodzaj="ostrzezenie"> */

  D.addEventListener('click', function (e) {
    var el = e.target.closest('[data-komunikat]');
    if (!el) return;
    toast(
      el.getAttribute('data-komunikat-tytul') || 'Informacja',
      el.getAttribute('data-komunikat'),
      el.getAttribute('data-komunikat-rodzaj') || 'informacja'
    );
  });

  /* ── 9 · SYMULACJA PRACY (stan → stan) ────────────────────────────────── */
  /* <button data-symuluj="#cel" data-symuluj-czas="1400" data-symuluj-stan="gotowe"> */

  D.addEventListener('click', function (e) {
    var el = e.target.closest('[data-symuluj]');
    if (!el) return;
    var cel = pierwszy(el.getAttribute('data-symuluj'));
    if (!cel) return;
    var czas = parseInt(el.getAttribute('data-symuluj-czas') || '1400', 10);
    var stanKoncowy = el.getAttribute('data-symuluj-stan') || 'gotowe';
    cel.setAttribute('data-stan', 'pracuje');
    var etykieta = el.querySelector('[data-etykieta]');
    var pierwotna = etykieta ? etykieta.textContent : null;
    if (etykieta) etykieta.textContent = el.getAttribute('data-symuluj-etykieta') || 'Pracuję…';
    setTimeout(function () {
      cel.setAttribute('data-stan', stanKoncowy);
      if (etykieta && pierwotna) etykieta.textContent = pierwotna;
      var kom = el.getAttribute('data-symuluj-komunikat');
      if (kom) toast(kom, null, 'sukces');
    }, czas);
  });

  /* ── 10 · POSTĘP ANIMOWANY ────────────────────────────────────────────── */
  /* <div class="dn-postep" data-postep-do="72" data-postep-czas="1600"> */

  function animujPostep(el) {
    var docelowy = parseFloat(el.getAttribute('data-postep-do') || '0');
    var czas = parseInt(el.getAttribute('data-postep-czas') || '1600', 10);
    var pasek = pierwszy('.dn-postep-wartosc', el);
    var licznik = pierwszy('[data-postep-licznik]', el);
    if (!pasek) return;
    if (matchMedia('(prefers-reduced-motion: reduce)').matches) {
      pasek.style.width = docelowy + '%';
      if (licznik) licznik.textContent = Math.round(docelowy) + '%';
      return;
    }
    var start = null;
    function krok(t) {
      if (start === null) start = t;
      var p = Math.min((t - start) / czas, 1);
      var wart = docelowy * (1 - Math.pow(1 - p, 3));
      pasek.style.width = wart + '%';
      if (licznik) licznik.textContent = Math.round(wart) + '%';
      if (p < 1) requestAnimationFrame(krok);
    }
    requestAnimationFrame(krok);
  }

  /* ── 11 · TABELA — sortowanie i filtrowanie ───────────────────────────── */

  function obsluzTabele(tabela) {
    wszystkie('th[data-sortuj]', tabela).forEach(function (th, kolumna) {
      th.setAttribute('tabindex', '0');
      th.setAttribute('role', 'columnheader');
      function sortuj() {
        var rosnaco = th.getAttribute('aria-sort') !== 'ascending';
        var tbody = pierwszy('tbody', tabela);
        var wiersze = wszystkie('tr', tbody);
        var idx = Array.prototype.indexOf.call(th.parentNode.children, th);
        var typ = th.getAttribute('data-sortuj');
        wiersze.sort(function (a, b) {
          var wa = a.children[idx] ? a.children[idx].textContent.trim() : '';
          var wb = b.children[idx] ? b.children[idx].textContent.trim() : '';
          if (typ === 'liczba') {
            var la = parseFloat(wa.replace(/[^\d.,-]/g, '').replace(',', '.')) || 0;
            var lb = parseFloat(wb.replace(/[^\d.,-]/g, '').replace(',', '.')) || 0;
            return rosnaco ? la - lb : lb - la;
          }
          return rosnaco ? wa.localeCompare(wb, 'pl') : wb.localeCompare(wa, 'pl');
        });
        wiersze.forEach(function (w) { tbody.appendChild(w); });
        wszystkie('th[data-sortuj]', tabela).forEach(function (x) { x.removeAttribute('aria-sort'); });
        th.setAttribute('aria-sort', rosnaco ? 'ascending' : 'descending');
      }
      th.addEventListener('click', sortuj);
      th.addEventListener('keydown', function (e) { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); sortuj(); } });
    });
  }

  function obsluzFiltr(pole) {
    var cel = pierwszy(pole.getAttribute('data-filtruj'));
    if (!cel) return;
    pole.addEventListener('input', function () {
      var fraza = pole.value.trim().toLowerCase();
      var elementy = wszystkie('[data-filtr-pozycja]', cel);
      var widoczne = 0;
      elementy.forEach(function (el) {
        var tekst = (el.getAttribute('data-filtr-tekst') || el.textContent).toLowerCase();
        var pasuje = !fraza || tekst.indexOf(fraza) !== -1;
        el.style.display = pasuje ? '' : 'none';
        if (pasuje) widoczne++;
      });
      var licznik = pierwszy('[data-filtr-licznik="' + (pole.id || '') + '"]');
      if (licznik) licznik.textContent = widoczne + ' / ' + elementy.length;
      var pusty = pierwszy('[data-filtr-pusty]', cel.parentNode || cel);
      if (pusty) pusty.style.display = widoczne ? 'none' : '';
    });
  }

  /* ── 12 · PASEK PROTOTYPU ─────────────────────────────────────────────── */

  function zbudujPasekPrototypu() {
    var meta = pierwszy('[data-prototyp]');
    if (!meta) return;
    var etykieta = meta.getAttribute('data-prototyp') || '';
    var zrodlo = meta.getAttribute('data-prototyp-zrodlo') || '';
    var okno = meta.getAttribute('data-prototyp-okno') || '';

    var pasek = D.createElement('div');
    pasek.className = 'pt-pasek-prototypu';
    pasek.setAttribute('data-widoczny', 'nie');
    pasek.innerHTML =
      '<b></b><span class="pt-okno-nazwa"></span>' +
      '<span class="pt-pasek-prawa"><span class="pt-zrodlo"></span>' +
      '<button class="dn-btn dn-btn--duch dn-btn--sm" type="button" data-zamknij-pasek>Ukryj</button></span>';
    pierwszy('b', pasek).textContent = etykieta;
    pierwszy('.pt-okno-nazwa', pasek).textContent = okno;
    pierwszy('.pt-zrodlo', pasek).textContent = zrodlo;

    var uchwyt = D.createElement('button');
    uchwyt.type = 'button';
    uchwyt.className = 'pt-uchwyt-prototypu';
    uchwyt.setAttribute('aria-expanded', 'false');
    uchwyt.innerHTML =
      '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>' +
      '<span>Prototyp</span>';

    function przelacz() {
      var widoczny = pasek.getAttribute('data-widoczny') === 'tak';
      pasek.setAttribute('data-widoczny', widoczny ? 'nie' : 'tak');
      uchwyt.setAttribute('aria-expanded', widoczny ? 'false' : 'true');
      uchwyt.style.bottom = widoczny ? '' : 'calc(var(--dn-od-3) + 28px)';
    }
    uchwyt.addEventListener('click', przelacz);
    pasek.addEventListener('click', function (e) { if (e.target.closest('[data-zamknij-pasek]')) przelacz(); });

    D.body.appendChild(pasek);
    D.body.appendChild(uchwyt);
  }

  /* ── 13 · SPIS TREŚCI Z PODŚWIETLANIEM ────────────────────────────────── */

  function obsluzSpis() {
    var spis = pierwszy('[data-spis]');
    if (!spis || !('IntersectionObserver' in window)) return;
    var odnosniki = wszystkie('a[href^="#"]', spis);
    var mapa = {};
    var cele = [];
    odnosniki.forEach(function (a) {
      var id = a.getAttribute('href').slice(1);
      var el = D.getElementById(id);
      if (el) { mapa[id] = a; cele.push(el); }
    });
    var obserwator = new IntersectionObserver(function (wpisy) {
      wpisy.forEach(function (w) {
        var a = mapa[w.target.id];
        if (!a) return;
        if (w.isIntersecting) {
          odnosniki.forEach(function (x) { x.removeAttribute('aria-current'); });
          a.setAttribute('aria-current', 'true');
        }
      });
    }, { rootMargin: '-10% 0px -75% 0px', threshold: 0 });
    cele.forEach(function (c) { obserwator.observe(c); });
  }

  /* ── 14 · SKRÓTY KLAWISZOWE ───────────────────────────────────────────── */

  D.addEventListener('keydown', function (e) {
    /* Esc zamyka menu podręczne */
    if (e.key === 'Escape') {
      wszystkie('[data-menu-tresc][data-otwarte="tak"]').forEach(function (m) { m.setAttribute('data-otwarte', 'nie'); });
    }
    /* Ctrl/Cmd + K — pole wyszukiwania paska górnego */
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
      var szukaj = pierwszy('.dn-pasek-szukaj input, .dn-szukaj input, [data-szukaj]');
      if (szukaj) { e.preventDefault(); szukaj.focus(); szukaj.select && szukaj.select(); }
    }
  });

  /* ── 15 · URUCHOMIENIE ────────────────────────────────────────────────── */

  function uruchom() {
    /* przełączniki widoków */
    D.addEventListener('click', function (e) {
      var p = e.target.closest('[data-przelacz]');
      if (!p) return;
      przelaczWidok(p.getAttribute('data-przelacz'), p.getAttribute('data-grupa') || null);
    });

    wszystkie('[role="tablist"]').forEach(obsluzZakladki);
    wszystkie('[data-lista-wyboru]').forEach(obsluzListeWyboru);
    wszystkie('table[data-tabela]').forEach(obsluzTabele);
    wszystkie('[data-filtruj]').forEach(obsluzFiltr);
    wszystkie('.dn-postep[data-postep-do]').forEach(animujPostep);

    /* domyślny widok każdej grupy */
    var grupy = {};
    wszystkie('[data-widok]').forEach(function (w) {
      var g = w.getAttribute('data-grupa-widoku') || '__';
      if (!grupy[g]) grupy[g] = [];
      grupy[g].push(w);
    });
    Object.keys(grupy).forEach(function (g) {
      var lista = grupy[g];
      var maAktywny = lista.some(function (w) { return w.getAttribute('data-widok-aktywny') === 'tak'; });
      if (!maAktywny && lista.length) lista[0].setAttribute('data-widok-aktywny', 'tak');
    });

    zbudujPasekPrototypu();
    obsluzSpis();
  }

  if (D.readyState === 'loading') D.addEventListener('DOMContentLoaded', uruchom);
  else uruchom();
})();
