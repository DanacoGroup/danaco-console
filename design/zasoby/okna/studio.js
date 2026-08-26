/* ══════════════════════════════════════════════════════════════════════════
   DANACO CONSOLE — STUDIO: ZACHOWANIA OKNA
   --------------------------------------------------------------------------
   Kompetencja: karty okna roboczego tego jednego okna — przełączanie myszą
   i klawiaturą wg wzorca ARIA „tabs”, zamykanie kart pomocniczych, wykaz kart
   w menu „⋮ → Karty okna roboczego” oraz przełącznik trybów w pasie stanowiska.
   Wygląd: `okna/studio.css`. Zachowania powłoki (menu, motyw, okna modalne)
   pochodzą z `wspolne.js` / `prototyp.js` / `rama.js` / `stanowisko.js`.
   Zasięg: wyłącznie `05-okna/moduly/studio.html`.
   ══════════════════════════════════════════════════════════════════════════ */
/* KARTY OKNA ROBOCZEGO — jedna karta robocza i po jednej karcie na każde okno
   pomocnicze. Wzorzec ARIA „tabs" z wędrującym tabindexem: Tab wchodzi w pasmo
   raz, strzałki przechodzą między kartami, Home/End skacze na koniec, Enter
   i Spacja otwierają kartę wskazaną strzałką. Stan wybrany niesie barwa karty
   (bryła zrośnięta ze wstążką okna), nie sam atrybut ARIA. */
(function(){
  var pasmo = document.querySelector('.st-karty');
  if (!pasmo) { return; }
  function karty(){ return Array.prototype.slice.call(pasmo.querySelectorAll('[role="tab"]')).filter(function(k){ return !k.hidden; }); }
  function panel(k){ var id = k.getAttribute('aria-controls'); return id ? document.getElementById(id) : null; }

  function otworz(karta, ustawFokus){
    karty().forEach(function(k){
      var wybrana = (k === karta);
      k.setAttribute('aria-selected', wybrana ? 'true' : 'false');
      k.setAttribute('tabindex', wybrana ? '0' : '-1');
      var p = panel(k); if (p) { p.hidden = !wybrana; }
    });
    if (ustawFokus) { karta.focus(); }
    if (window.dnOglos) { window.dnOglos('Karta ' + karta.textContent.trim()); }
  }

  function zamknij(karta){
    var lista = karty();
    if (lista.length < 2 || karta.classList.contains('st-karta--robocza')) { return; }
    var i = lista.indexOf(karta);
    var nastepna = lista[i + 1] || lista[i - 1];
    var byla = (karta.getAttribute('aria-selected') === 'true');
    var p = panel(karta); if (p) { p.hidden = true; }
    /* Karta odłożona nie może zostać z wybranym stanem: wykaz kart musi mieć
       dokładnie jedną kartę bieżącą, także po zamknięciu tej, która nią była. */
    karta.setAttribute('aria-selected', 'false');
    karta.setAttribute('tabindex', '-1');
    karta.hidden = true;
    if (byla) { otworz(nastepna, true); }
    if (window.dnToast) { window.dnToast('Karta zamknięta', karta.textContent.trim(), 'informacja'); }
  }

  pasmo.addEventListener('click', function(e){
    var zam = e.target.closest('[data-karta-zamknij]');
    if (zam) { zamknij(zam.closest('[role="tab"]')); return; }
    var karta = e.target.closest('[role="tab"]');
    if (karta) { otworz(karta, true); }
  });

  pasmo.addEventListener('keydown', function(e){
    var karta = e.target.closest('[role="tab"]');
    if (!karta) { return; }
    var lista = karty(), i = lista.indexOf(karta), cel = null;
    if (e.key === 'ArrowRight') { cel = lista[(i + 1) % lista.length]; }
    else if (e.key === 'ArrowLeft') { cel = lista[(i - 1 + lista.length) % lista.length]; }
    else if (e.key === 'Home') { cel = lista[0]; }
    else if (e.key === 'End') { cel = lista[lista.length - 1]; }
    else if (e.key === 'Enter' || e.key === ' ' || e.key === 'Spacebar') { e.preventDefault(); otworz(karta, true); return; }
    else if (e.key === 'Delete') { e.preventDefault(); zamknij(karta); return; }
    if (cel) { e.preventDefault(); otworz(cel, true); }
  });

  /* Zamknięcie okna pomocniczego z jego własnej belki zamyka także kartę. */
  document.addEventListener('click', function(e){
    var zam = e.target.closest('[data-karta-zamknij]');
    if (!zam || pasmo.contains(zam)) { return; }
    var okno = zam.closest('.sta-okno');
    if (!okno) { return; }
    var karta = pasmo.querySelector('[aria-controls="' + okno.id + '"]');
    if (karta) { zamknij(karta); }
  });

  /* Wykaz kart w menu okna komunikacji i w pasmie prowadzi do tej samej karty. */
  document.addEventListener('click', function(e){
    var poz = e.target.closest('[data-karta-przelacz]');
    if (!poz) { return; }
    var karta = pasmo.querySelector('[data-karta="' + poz.getAttribute('data-karta-przelacz') + '"]');
    if (!karta) { return; }
    if (karta.hidden) { karta.hidden = false; }
    otworz(karta, false);
  });

  document.querySelectorAll('[data-srod]').forEach(function(b){
    b.addEventListener('click', function(){
      document.querySelectorAll('[data-srod]').forEach(function(x){ x.setAttribute('aria-pressed', x === b ? 'true' : 'false'); });
      if (window.dnToast) { window.dnToast('Tryb: ' + b.textContent, 'Obszar polecenia i panele dostrojone do trybu', 'informacja'); }
    });
  });
})();
