/* ============================================================================
   DANACO CONSOLE — PANEL SESJI I PROJEKTÓW
   ----------------------------------------------------------------------------
   Kompetencja: porządek i zakres płaskiego wykazu sesji, drzewo projektów,
   czynności wiersza, czynności panelu, zwinięcie i przypięcie panelu.

   Wymaga: `prototyp.js` (mechanika `data-menu`, `dnToast`, `dnOglos`),
   `stany.js` (stan zbiorczy pozycji), `panel-sesji.css`.

   Kontrakt znacznika:
     .dn-panel-wykaz[data-sesje-sort][data-sesje-widok]       wykaz sesji
     .dn-panel-poz[data-srodowisko][data-stan]                wymiary wiersza
     [data-sesje-sort="czynnosc|najstarsze|nazwa|nazwa-odwrotnie|srodowisko"]
     [data-sesje-widok="wszystkie|czynne|reakcja|zakonczone"]
     [data-projekty-sort="nazwa|czynnosc"]                    porządek drzewa
     [data-projekty-widok="wszystkie|praca|reakcja"]          zakres drzewa
     [data-panel-strona] · [data-panel-eksport]               czynności panelu
     [data-drzewo="rozwin|zwin"]                              stan gałęzi
     [data-poz-akcja="nazwa|przenies|archiwizuj|usun"]        czynność wiersza
     [data-panel-zwin] · [data-panel-przypnij]                stan panelu
   ============================================================================ */
(function () {
  'use strict';

  function q(s, k) { return (k || document).querySelector(s); }
  function qq(s, k) { return Array.prototype.slice.call((k || document).querySelectorAll(s)); }
  function toast(t, o, r, ms) { if (window.dnToast) { window.dnToast(t, o, r || 'informacja', ms || 3200); } }
  function oglos(t) { if (window.dnOglos) { window.dnOglos(t); } }

  var wykaz = q('#wykaz-sesji');
  var panel = q('.dn-obszar-panel--boczny');

  /* ── Porządek wykazu sesji ───────────────────────────────────────────────
     Kolejność wyjściowa znacznika jest porządkiem „ostatnio czynne"; pozostałe
     porządki liczą się z niej, więc powrót do niej nie wymaga przeładowania. */

  /* Zestaw stanów sesji jest zamknięty i liczy trzy wartości: praca (w toku),
     reakcja (czeka na Operatora), zakonczone. Stan `blad` nie występuje —
     sesja przerwana czeka na rozstrzygnięcie, czyli stoi w `reakcja`. */
  var STAN_WAGA = { reakcja: 0, praca: 1, zakonczone: 2 };
  var NAZWY_SORT = {
    czynnosc: 'ostatnio czynne',
    najstarsze: 'najstarsze',
    nazwa: 'nazwa A–Z',
    'nazwa-odwrotnie': 'nazwa Z–A',
    srodowisko: 'środowisko'
  };
  var NAZWY_WIDOK = {
    wszystkie: 'wszystkie sesje',
    czynne: 'w toku',
    reakcja: 'czekające na Operatora',
    zakonczone: 'zakończone'
  };
  var NAZWY_WIDOK_PROJEKTY = {
    wszystkie: 'wszystkie projekty',
    praca: 'z pracą w toku',
    reakcja: 'czekające na Operatora'
  };

  function obudowy(w) { return qq(':scope > .sta-menu', w); }
  function nazwaWiersza(o) {
    var n = q('.dn-obszar-pozycja-nazwa', o);
    return n ? n.textContent.trim() : '';
  }
  function stanWiersza(o) {
    var p = q('.dn-panel-poz', o);
    return p ? p.getAttribute('data-stan') || '' : '';
  }
  function srodWiersza(o) {
    var p = q('.dn-panel-poz', o);
    return p ? p.getAttribute('data-srodowisko') || '' : '';
  }

  var kolejnoscZnacznika = wykaz ? obudowy(wykaz) : [];

  function sortuj(tryb) {
    if (!wykaz) { return; }
    var lista = kolejnoscZnacznika.slice();
    if (tryb === 'najstarsze') { lista.reverse(); }
    if (tryb === 'nazwa' || tryb === 'nazwa-odwrotnie') {
      lista.sort(function (a, b) { return nazwaWiersza(a).localeCompare(nazwaWiersza(b), 'pl'); });
      if (tryb === 'nazwa-odwrotnie') { lista.reverse(); }
    }
    if (tryb === 'srodowisko') {
      lista.sort(function (a, b) {
        var r = srodWiersza(a).localeCompare(srodWiersza(b), 'pl');
        return r !== 0 ? r : STAN_WAGA[stanWiersza(a)] - STAN_WAGA[stanWiersza(b)];
      });
    }
    lista.forEach(function (o) { wykaz.appendChild(o); });
    wykaz.setAttribute('data-sesje-sort', tryb);
  }

  function ustawWybor(menu, atrybut, wartosc) {
    qq('[' + atrybut + ']', menu).forEach(function (b) {
      b.setAttribute('aria-checked', b.getAttribute(atrybut) === wartosc ? 'true' : 'false');
    });
  }

  qq('[data-sesje-sort]').filter(function (b) { return b.hasAttribute('role'); }).forEach(function (b) {
    b.addEventListener('click', function () {
      var tryb = b.getAttribute('data-sesje-sort');
      sortuj(tryb);
      ustawWybor(b.closest('[data-menu-tresc]'), 'data-sesje-sort', tryb);
      oglos('Wykaz sesji: ' + NAZWY_SORT[tryb] + '.');
    });
  });

  qq('[data-sesje-widok]').filter(function (b) { return b.hasAttribute('role'); }).forEach(function (b) {
    b.addEventListener('click', function () {
      var widok = b.getAttribute('data-sesje-widok');
      if (!wykaz) { return; }
      wykaz.setAttribute('data-sesje-widok', widok);
      ustawWybor(b.closest('[data-menu-tresc]'), 'data-sesje-widok', widok);
      var ile = obudowy(wykaz).filter(function (o) { return o.offsetParent !== null; }).length;
      oglos('Wykaz sesji: ' + NAZWY_WIDOK[widok] + ', pozycji ' + ile + '.');
    });
  });

  /* ── Drzewo projektów ───────────────────────────────────────────────────── */

  qq('[data-drzewo]').forEach(function (b) {
    b.addEventListener('click', function () {
      var rozwin = b.getAttribute('data-drzewo') === 'rozwin';
      qq('.dn-panel-galaz').forEach(function (g) { g.open = rozwin; });
      oglos(rozwin ? 'Drzewo projektów rozwinięte.' : 'Drzewo projektów zwinięte.');
    });
  });

  /* Projekt nie nosi własnego znacznika stanu — znacznik należy do sesji.
     Porządek „ostatnio czynne" liczy więc stan najpilniejszy spośród sesji
     gałęzi; gałąź bez sesji idzie na koniec. */
  function stanGalezi(galaz) {
    var wagi = qq('.dn-panel-poz', galaz).map(function (poz) {
      var w = STAN_WAGA[poz.getAttribute('data-stan')];
      return w === undefined ? 9 : w;
    });
    return wagi.length ? Math.min.apply(null, wagi) : 9;
  }

  qq('[data-projekty-widok]').filter(function (b) { return b.hasAttribute('role'); }).forEach(function (b) {
    b.addEventListener('click', function () {
      var widok = b.getAttribute('data-projekty-widok');
      var drzewo = q('.dn-panel-drzewo');
      if (!drzewo) { return; }
      drzewo.setAttribute('data-projekty-widok', widok);
      ustawWybor(b.closest('[data-menu-tresc]'), 'data-projekty-widok', widok);
      var ile = qq(':scope > .dn-panel-galaz', drzewo).filter(function (g) { return g.offsetParent !== null; }).length;
      oglos('Drzewo projektów: ' + NAZWY_WIDOK_PROJEKTY[widok] + ', gałęzi ' + ile + '.');
    });
  });

  qq('[data-projekty-sort]').filter(function (b) { return b.hasAttribute('role'); }).forEach(function (b) {
    b.addEventListener('click', function () {
      var tryb = b.getAttribute('data-projekty-sort');
      var drzewo = q('.dn-panel-drzewo');
      if (!drzewo) { return; }
      var galezie = qq(':scope > .dn-panel-galaz', drzewo);
      var luzne = galezie.filter(function (g) { return !g.hasAttribute('data-projekt'); });
      var projekty = galezie.filter(function (g) { return g.hasAttribute('data-projekt'); });
      if (tryb === 'nazwa') {
        projekty.sort(function (a, b2) {
          return q('.dn-panel-galaz-nazwa', a).textContent.localeCompare(
            q('.dn-panel-galaz-nazwa', b2).textContent, 'pl');
        });
      } else {
        projekty.sort(function (a, b2) { return stanGalezi(a) - stanGalezi(b2); });
      }
      projekty.concat(luzne).forEach(function (g) { drzewo.appendChild(g); });
      ustawWybor(b.closest('[data-menu-tresc]'), 'data-projekty-sort', tryb);
      oglos(tryb === 'nazwa' ? 'Projekty po nazwie.' : 'Projekty po ostatniej czynności.');
    });
  });

  qq('.dn-panel-galaz').forEach(function (galaz) {
    galaz.addEventListener('toggle', function () {
      if (window.dnStany) { window.dnStany.odswiez(galaz); }
    });
  });

  /* ── Czynności wiersza ──────────────────────────────────────────────────── */

  var OPISY = {
    nazwa: 'Nazwa do zmiany w miejscu — automat nadał ją z pierwszego polecenia sesji.',
    przenies: 'Przeniesienie do projektu: sesja wchodzi w drzewo projektu i zostaje w wykazie sesji.',
    archiwizuj: 'Archiwizacja zdejmuje sesję z wykazu czynnego; treść zostaje dostępna z historii.',
    usun: 'Usunięcie trwałe. Zadania serwerowe tej sesji zostają zatrzymane.'
  };

  qq('[data-poz-akcja]').forEach(function (b) {
    b.addEventListener('click', function () {
      var akcja = b.getAttribute('data-poz-akcja');
      var obudowa = b.closest('.sta-menu');
      var nazwa = obudowa ? nazwaWiersza(obudowa) : 'Sesja';
      toast(nazwa, OPISY[akcja] || akcja, akcja === 'usun' ? 'ostrzezenie' : 'informacja', 3600);
      oglos(nazwa + ' — ' + akcja + '.');
    });
  });

  /* ── Zwinięcie i przypięcie panelu ──────────────────────────────────────── */

  qq('[data-panel-zwin]').forEach(function (b) {
    b.addEventListener('click', function () {
      if (!panel) { return; }
      panel.hidden = !panel.hidden;
      b.setAttribute('aria-pressed', panel.hidden ? 'true' : 'false');
      oglos(panel.hidden ? 'Panel sesji zwinięty.' : 'Panel sesji rozwinięty.');
    });
  });

  qq('[data-panel-przypnij]').forEach(function (b) {
    b.addEventListener('click', function () {
      var wl = b.getAttribute('aria-pressed') === 'true';
      b.setAttribute('aria-pressed', wl ? 'false' : 'true');
      toast('Panel sesji', wl
        ? 'Przypięcie zdjęte: panel zwija się przy wejściu w moduł.'
        : 'Panel przypięty: zostaje otwarty niezależnie od modułu w oknie.', 'informacja');
    });
  });

  /* ── Pierwszeństwo nazwy w wierszu wykazu ────────────────────────────────
     Nazwa jest rzeczą, po której Operator poznaje wiersz, więc skraca się jako
     ostatnia: niedobór miejsca zdejmuje najpierw oznaczenie środowiska
     (`panel-sesji.css`, współczynnik kurczenia). Zostaje jeden przypadek, którego
     arkusz nie rozstrzygnie: gdy na oznaczenie zostaje mniej niż kilka znaków,
     wielokropek zostawia kikut w rodzaju „T…”. Wtedy oznaczenie ustępuje w
     całości — wiersz woli nie pokazać wymiaru porządkującego, niż pokazać jego
     resztkę. */
  var PROG_SROD = 4;   /* znaki — poniżej tego oznaczenie nic już nie mówi */

  function pierwszenstwoNazwy() {
    qq('.dn-panel-wiersz').forEach(function (wiersz) {
      var srod = q('.dn-panel-srod', wiersz);
      if (!srod) { return; }
      wiersz.removeAttribute('data-srod');
      if (!srod.offsetParent) { return; }
      var znak = srod.scrollWidth / Math.max(1, srod.textContent.length);
      if (srod.clientWidth < znak * PROG_SROD) { wiersz.setAttribute('data-srod', 'ukryte'); }
    });
  }

  window.addEventListener('resize', pierwszenstwoNazwy);
  document.addEventListener('click', function (e) {
    if (e.target.closest && e.target.closest('.dn-karta-widoku, .dn-panel-galaz > summary, [data-drzewo], [data-panel-strona]')) {
      window.setTimeout(pierwszenstwoNazwy, 0);
    }
  }, true);
  pierwszenstwoNazwy();

  /* ── Czynności panelu ───────────────────────────────────────────────────── */

  qq('[data-panel-strona]').forEach(function (b) {
    b.addEventListener('click', function () {
      var prawa = b.getAttribute('aria-checked') === 'true';
      b.setAttribute('aria-checked', prawa ? 'false' : 'true');
      if (panel) { panel.setAttribute('data-strona', prawa ? 'lewo' : 'prawo'); }
      toast('Panel sesji', prawa
        ? 'Panel wraca na lewą stronę okna.'
        : 'Panel przechodzi na prawą stronę okna.', 'informacja');
    });
  });

  qq('[data-panel-eksport]').forEach(function (b) {
    b.addEventListener('click', function () {
      toast('Eksport wykazu', 'Wykaz sesji i projektów zapisany do pliku przenośnego.', 'informacja', 3600);
    });
  });

  /* ── Wiersz wykazu: wejście w sesję ─────────────────────────────────────── */

  document.addEventListener('click', function (e) {
    var wiersz = e.target.closest && e.target.closest('.dn-obszar-pozycja');
    if (!wiersz || e.target.closest('.dn-obszar-pozycja-menu')) { return; }
    qq('.dn-obszar-pozycja[aria-current]').forEach(function (x) { x.removeAttribute('aria-current'); });
    wiersz.setAttribute('aria-current', 'true');
    var nazwa = (q('.dn-obszar-pozycja-nazwa', wiersz) || wiersz).textContent.trim();
    toast(nazwa, 'Sesja wraca w oknie roboczym na karcie Centrum dowodzenia; pozostałe karty stoją obok.', 'sukces', 3600);
    oglos('Otwieram sesję ' + nazwa + '.');
  });
})();
