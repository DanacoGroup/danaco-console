/* ============================================================================
   DANACO CONSOLE — WARSTWA STANOWISKA (zachowanie) · v2.1
   ----------------------------------------------------------------------------
   Steruje elastycznym układem wielookiennym: liczba okien czatu (1/2/4),
   widoczność strefy roboczej, aktywne okno, popovery przybornika, panel
   konfiguracji per okno, prompt i symulacja wysyłki.

   Współpracuje z prototyp.js (toasty, motyw). Sterowanie atrybutami data-*.
   ============================================================================ */
(function () {
  'use strict';
  var D = document;
  function all(s, r) { return Array.prototype.slice.call((r || D).querySelectorAll(s)); }
  function toast(t, o, r) { return (window.dnToast ? window.dnToast(t, o, r) : null); }

  /* ── 1 · LICZBA OKIEN CZATU (1 / 2 / 4) ───────────────────────────────── */
  function ustawCzaty(n) {
    var strefa = D.querySelector('.sta-czaty');
    if (!strefa) return;
    var okna = all('.sta-okno', strefa);
    okna.forEach(function (o, i) { o.hidden = i >= n; });
    strefa.setAttribute('data-liczba', String(n));
    all('[data-czaty-ile]').forEach(function (b) {
      b.setAttribute('aria-pressed', b.getAttribute('data-czaty-ile') === String(n) ? 'true' : 'false');
    });
  }

  /* ── 2 · STREFA ROBOCZA (pokaż/ukryj, liczba) ─────────────────────────── */
  function przelaczRobocza() {
    var obszar = D.querySelector('.sta-obszar');
    if (!obszar) return;
    var ukryta = obszar.getAttribute('data-robocza') === 'ukryta';
    obszar.setAttribute('data-robocza', ukryta ? 'widoczna' : 'ukryta');
    all('[data-robocza-przelacz]').forEach(function (b) { b.setAttribute('aria-pressed', ukryta ? 'true' : 'false'); });
  }

  /* ── 3 · AKTYWNE OKNO (klik oznacza okno jako aktywne) ────────────────── */
  D.addEventListener('click', function (e) {
    var okno = e.target.closest('.sta-okno');
    if (!okno) return;
    all('.sta-okno').forEach(function (o) { o.setAttribute('data-aktywne', o === okno ? 'tak' : 'nie'); });
  });

  /* ── 4 · POPOVERY PRZYBORNIKA (przełączniki i rozwijane menu) ─────────── */
  function zamknijPopovery(oprocz) {
    all('.sta-popover').forEach(function (p) { if (p !== oprocz) { p.hidden = true;
      var b = p.closest('.sta-nrz'); if (b) { var pr = b.querySelector('[aria-expanded]'); if (pr) pr.setAttribute('aria-expanded', 'false'); } } });
  }
  D.addEventListener('click', function (e) {
    var przel = e.target.closest('[data-popover]');
    if (przel) {
      e.stopPropagation();
      var cel = D.getElementById(przel.getAttribute('data-popover'));
      if (!cel) return;
      var otw = !cel.hidden;
      zamknijPopovery(otw ? null : cel);
      cel.hidden = otw;
      przel.setAttribute('aria-expanded', otw ? 'false' : 'true');
      return;
    }
    var przelacznik = e.target.closest('[data-przelacznik]');
    if (przelacznik && !przelacznik.hasAttribute('data-popover')) {
      var on = przelacznik.getAttribute('aria-pressed') === 'true';
      przelacznik.setAttribute('aria-pressed', on ? 'false' : 'true');
      var ety = przelacznik.getAttribute('data-etykieta');
      if (ety) toast(ety + (on ? ' — wyłączone' : ' — włączone'), 'Ustawienie tego okna', on ? 'informacja' : 'sukces');
      return;
    }
    if (!e.target.closest('.sta-popover')) zamknijPopovery(null);
  });

  /* ── 5 · KONFIGURACJA PER OKNO (osobny komplet ustawień) ──────────────── */
  D.addEventListener('click', function (e) {
    var otw = e.target.closest('[data-konfig-okno]');
    if (otw) {
      var id = otw.getAttribute('data-konfig-okno');
      var panel = D.getElementById(id);
      if (panel) { panel.hidden = !panel.hidden; otw.setAttribute('aria-expanded', panel.hidden ? 'false' : 'true'); }
      return;
    }
    var zam = e.target.closest('[data-konfig-zamknij]');
    if (zam) { var p = zam.closest('.sta-konfig'); if (p) p.hidden = true; }
  });

  /* ── 6 · ZAMYKANIE / PRZYWRACANIE OKNA ────────────────────────────────── */
  D.addEventListener('click', function (e) {
    var zam = e.target.closest('[data-okno-zamknij]');
    if (!zam) return;
    var okno = zam.closest('.sta-okno');
    if (!okno) return;
    okno.hidden = true;
    var nazwa = okno.getAttribute('data-nazwa') || 'Okno';
    toast(nazwa + ' zamknięte', 'Przywróć z paska układu', 'informacja');
    var pas = D.querySelector('[data-przywroc-pas]');
    if (pas && okno.id) {
      var b = D.createElement('button');
      b.className = 'dn-btn dn-btn--zarys dn-btn--sm';
      b.textContent = '+ ' + nazwa;
      b.addEventListener('click', function () { okno.hidden = false; b.remove(); });
      pas.appendChild(b);
    }
  });

  /* ── 7 · PROMPT — autorozrost + wysyłka z symulacją odpowiedzi ─────────── */
  function podłączPrompt(form) {
    var pole = form.querySelector('.sta-prompt-obszar');
    var historia = form.closest('.sta-kom') && form.closest('.sta-kom').querySelector('.sta-kom-historia');
    if (pole) {
      pole.addEventListener('input', function () { pole.style.height = 'auto'; pole.style.height = Math.min(pole.scrollHeight, 140) + 'px'; });
      pole.addEventListener('keydown', function (e) { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); wyslij(); } });
    }
    form.addEventListener('submit', function (e) { e.preventDefault(); wyslij(); });
    var wyslijBtn = form.querySelector('.sta-prompt-wyslij');
    if (wyslijBtn) wyslijBtn.addEventListener('click', function (e) { e.preventDefault(); wyslij(); });

    function wpis(rodzaj, nadawca, tresc, medalion) {
      var el = D.createElement('article');
      el.className = 'sta-wpis sta-wpis--' + rodzaj;
      el.innerHTML = '<div class="sta-wpis-medalion" aria-hidden="true">' + medalion + '</div><div>' +
        '<div class="sta-wpis-tozsamosc"><span class="sta-wpis-nadawca"></span><span class="sta-wpis-godzina">teraz</span></div>' +
        '<div class="sta-wpis-tresc"></div></div>';
      el.querySelector('.sta-wpis-nadawca').textContent = nadawca;
      el.querySelector('.sta-wpis-tresc').textContent = tresc;
      return el;
    }
    var IK_CZLOWIEK = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>';
    var IK_AI = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M12 8V4H8"/><rect width="16" height="12" x="4" y="8" rx="2"/><path d="M2 14h2"/><path d="M20 14h2"/><path d="M15 13v2"/><path d="M9 13v2"/></svg>';

    function wyslij() {
      if (!pole || !historia) return;
      var tekst = pole.value.trim();
      if (!tekst) { toast('Pole polecenia jest puste', 'Wpisz polecenie dla modelu', 'ostrzezenie'); return; }
      historia.appendChild(wpis('czlowiek', 'Operator', tekst, IK_CZLOWIEK));
      pole.value = ''; pole.style.height = 'auto';
      historia.scrollTop = historia.scrollHeight;
      var w = wpis('inteligencja', 'Inteligencja', 'Przetwarzam polecenie…', IK_AI);
      w.querySelector('.sta-wpis-tresc').innerHTML = '<span class="pt-tetno" aria-hidden="true"></span> Przetwarzam polecenie…';
      historia.appendChild(w); historia.scrollTop = historia.scrollHeight;
      var okno = form.closest('.sta-okno'); if (okno) okno.setAttribute('data-stan', 'pracuje');
      setTimeout(function () {
        w.querySelector('.sta-wpis-tresc').textContent = 'Przyjęte. Wynik zależy od kontekstu tego okna — model, kanał i źródła są konfigurowane osobno dla tej sesji (ikona ustawień w belce).';
        if (okno) okno.removeAttribute('data-stan');
        historia.scrollTop = historia.scrollHeight;
      }, 1400);
    }
  }

  /* ── 8 · URUCHOMIENIE ─────────────────────────────────────────────────── */
  function start() {
    all('[data-czaty-ile]').forEach(function (b) {
      b.addEventListener('click', function () { ustawCzaty(parseInt(b.getAttribute('data-czaty-ile'), 10)); });
    });
    all('[data-robocza-przelacz]').forEach(function (b) { b.addEventListener('click', przelaczRobocza); });
    all('.sta-prompt').forEach(podłączPrompt);
    var pierwsze = D.querySelector('.sta-czaty .sta-okno');
    if (pierwsze) pierwsze.setAttribute('data-aktywne', 'tak');
    D.addEventListener('keydown', function (e) { if (e.key === 'Escape') zamknijPopovery(null); });
  }
  if (D.readyState === 'loading') D.addEventListener('DOMContentLoaded', start); else start();
})();
