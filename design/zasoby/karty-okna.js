/* ============================================================================
   DANACO CONSOLE — PASMO KART OKNA ROBOCZEGO
   ----------------------------------------------------------------------------
   Kompetencja: przełączanie kart, karta przypięta, zmiana nazwy karty,
   zamykanie, grupy kart, menu powłok, grot, „+" i ustawienia widoku.

   Wymaga: `prototyp.js` (mechanika `data-menu`, `dnToast`, `dnOglos`),
   `stany.js`, `karty-okna.css`.

   Kontrakt znacznika:
     .dn-karty-pasmo[data-grupowanie="tak|nie"][data-gestosc]  ustawienia widoku
     .dn-karty-lista[role="tablist"]                           wyłącznie zakładki
     .dn-karta-widoku[data-karta][data-karta-rodzaj]           karta pasma
       data-karta-rodzaj="centrum"                             karta główna
     [data-okno-akcja="przelacz|zamknij|nowe|incognito|galaz"] menu powłok
     [data-karta-przelacz="…"]                                 grot
     [data-nowa-karta-modul="…"]                                „+"
     [data-widok="grupowanie|rozgrupuj|gestosc|archiwalne"]     ustawienia widoku
   ============================================================================ */
(function () {
  'use strict';

  function q(s, k) { return (k || document).querySelector(s); }
  function qq(s, k) { return Array.prototype.slice.call((k || document).querySelectorAll(s)); }
  function toast(t, o, r, ms) { if (window.dnToast) { window.dnToast(t, o, r || 'informacja', ms || 3200); } }
  function oglos(t) { if (window.dnOglos) { window.dnOglos(t); } }

  var pasmo = q('.dn-obszar-panel--glowny .dn-karty-pasmo');
  if (!pasmo) { return; }

  var karty = qq('.dn-karta-widoku', pasmo);

  /* ── Przełączanie kart ──────────────────────────────────────────────────── */

  function panelKarty(k) {
    var id = k.getAttribute('aria-controls');
    return id ? document.getElementById(id) : null;
  }

  function aktywuj(karta) {
    /* Kilka kart może wskazywać ten sam panel (karty modułowe dzielą zasłonę
       treści), więc panele gasimy raz, a dopiero potem odsłaniamy panel karty
       bieżącej. */
    karty.forEach(function (k) {
      k.setAttribute('aria-selected', k === karta ? 'true' : 'false');
      var panel = panelKarty(k);
      if (panel) { panel.hidden = true; }
    });
    var biezacy = panelKarty(karta);
    if (biezacy) { biezacy.hidden = false; }
    document.dispatchEvent(new CustomEvent('dn-karta-zmiana', { detail: {
      karta: karta.getAttribute('data-karta'),
      rodzaj: karta.getAttribute('data-karta-rodzaj') || 'modul',
      modul: karta.getAttribute('data-modul') || ''
    } }));
  }

  karty.forEach(function (karta) {
    karta.addEventListener('click', function (e) {
      if (e.target.closest('.dn-karta-widoku-zamknij, .dn-karta-widoku-grot, .sta-menu')) { return; }
      aktywuj(karta);
    });
    karta.addEventListener('dblclick', function (e) {
      if (karta.getAttribute('data-karta-rodzaj') === 'centrum') { return; }
      var nazwa = q('.dn-karta-widoku-nazwa', karta);
      if (!nazwa || e.target.closest('.dn-karta-widoku-zamknij')) { return; }
      nazwa.setAttribute('contenteditable', 'true');
      nazwa.focus();
      document.getSelection().selectAllChildren(nazwa);
    });
  });

  qq('.dn-karta-widoku-nazwa', pasmo).forEach(function (nazwa) {
    nazwa.addEventListener('blur', function () {
      if (nazwa.getAttribute('contenteditable') !== 'true') { return; }
      nazwa.removeAttribute('contenteditable');
      toast('Nazwa karty', 'Nazwa zapisana: ' + nazwa.textContent.trim(), 'sukces');
    });
    nazwa.addEventListener('keydown', function (e) {
      if (e.key === 'Enter') { e.preventDefault(); nazwa.blur(); }
      if (e.key === 'Escape') { nazwa.blur(); }
    });
  });

  /* Zamknięcie karty. Przycisk zamknięcia leży wewnątrz pozycji `tab`, więc nie
     stoi w porządku tabulacji (`tabindex="-1"`) — pozycja o roli zakładki nie
     może mieć własnych celów tabulacji. Z klawiatury kartę zamyka Ctrl+W albo
     Delete na zaznaczonej karcie; myszą — przycisk. */
  function zamknij(karta) {
    if (!karta || karta.getAttribute('data-karta-rodzaj') === 'centrum') { return; }
    toast('Karta zamknięta', 'Zadania serwerowe tej karty biegną dalej jako sesja w tle.', 'informacja');
    oglos('Karta ' + ((q('.dn-karta-widoku-nazwa', karta) || {}).textContent || '') + ' zamknięta.');
  }

  qq('.dn-karta-widoku-zamknij', pasmo).forEach(function (b) {
    b.addEventListener('click', function (e) {
      e.stopPropagation();
      zamknij(b.closest('.dn-karta-widoku'));
    });
  });

  karty.forEach(function (karta) {
    karta.addEventListener('keydown', function (e) {
      if (e.key === 'Delete' || (e.key.toLowerCase() === 'w' && e.ctrlKey)) {
        e.preventDefault();
        zamknij(karta);
      }
    });
  });

  /* Grot karty grupy też stoi poza porządkiem tabulacji. Z klawiatury wykaz
     kart grupy otwiera Alt+Strzałka w dół na zaznaczonej karcie — tak samo jak
     rozwinięcie pola wyboru w oknach systemowych. */
  karty.forEach(function (karta) {
    karta.addEventListener('keydown', function (e) {
      if (e.altKey && e.key === 'ArrowDown') {
        var grot = q('.dn-karta-widoku-grot', karta);
        if (grot) { e.preventDefault(); grot.click(); }
      }
    });
  });

  /* ── Menu powłok: okna robocze bieżącej sesji ───────────────────────────── */

  var OPISY_OKNA = {
    przelacz: 'Przełączenie na okno robocze tej sesji; karty okna zostają tam, gdzie były.',
    zamknij: 'Zamknięcie okna roboczego. Sesja i jej zadania serwerowe zostają.',
    nowe: 'Nowe okno robocze bez izolacji: ta sama wiedza i te same pliki.',
    incognito: 'Nowe okno incognito: własna wiedza sesji, praca zapisywana we własnej gałęzi.',
    galaz: 'Nowe okno incognito na własnej gałęzi: izolacja wiedzy i plików rozłącznie.'
  };

  qq('[data-okno-akcja]').forEach(function (b) {
    b.addEventListener('click', function () {
      var akcja = b.getAttribute('data-okno-akcja');
      var okno = b.getAttribute('data-okno-nazwa') || 'Okno robocze';
      toast(okno, OPISY_OKNA[akcja] || akcja, akcja === 'zamknij' ? 'ostrzezenie' : 'informacja', 3600);
      oglos(okno + ' — ' + akcja + '.');
    });
  });

  /* ── Grot: przełączanie otwartych kart ─────────────────────────────────── */

  qq('[data-karta-przelacz]').forEach(function (b) {
    b.addEventListener('click', function () {
      var klucz = b.getAttribute('data-karta-przelacz');
      var karta = q('.dn-karta-widoku[data-karta="' + klucz + '"]', pasmo);
      if (karta) { aktywuj(karta); }
      oglos('Karta ' + klucz + '.');
    });
  });

  /* ── „+": moduł jako nowa karta zakładająca nową grupę ─────────────────── */

  /* ── Wejscie do srodowiska z karty centrum ─────────────────────────────
     Przycisk „Wejdz" na karcie srodowiska otwiera karte tego srodowiska
     w oknie roboczym; okno lewe i okno prawe zostaja na miejscu. */

  qq('[data-wejdz]').forEach(function (b) {
    b.addEventListener('click', function () {
      var nazwa = b.getAttribute('data-wejdz');
      var karta = q('.dn-karta-widoku[data-srodowisko="' + nazwa + '"]', pasmo);
      if (!karta) { return; }
      aktywuj(karta);
      oglos('Otwarto srodowisko ' + nazwa + ' w oknie roboczym.');
    });
  });

  qq('[data-nowa-karta-modul]').forEach(function (b) {
    b.addEventListener('click', function () {
      var modul = b.getAttribute('data-nowa-karta-modul');
      toast('Nowa karta', modul + ' otwiera się jako nowa karta i zakłada nową grupę zadania.', 'sukces', 3600);
      oglos('Nowa karta: ' + modul + '.');
    });
  });

  /* ── Ustawienia widoku pasma ───────────────────────────────────────────── */

  var OPISY_WIDOKU = {
    grupowanie: 'Karty jednego zadania stoją razem; okna pomocnicze wchodzą do grupy karty macierzystej.',
    rozgrupuj: 'Grupy rozwiązane: karty stoją płasko w kolejności otwarcia.',
    archiwalne: 'Karty zarchiwizowane ukryte w pasmie; wracają z historii sesji.'
  };

  qq('[data-widok]').forEach(function (b) {
    b.addEventListener('click', function () {
      var co = b.getAttribute('data-widok');
      if (co === 'grupowanie') {
        var wl = pasmo.getAttribute('data-grupowanie') !== 'nie';
        pasmo.setAttribute('data-grupowanie', wl ? 'nie' : 'tak');
        b.setAttribute('aria-checked', wl ? 'false' : 'true');
        toast('Grupowanie kart', wl ? OPISY_WIDOKU.rozgrupuj : OPISY_WIDOKU.grupowanie, 'informacja', 3600);
        return;
      }
      if (co === 'gestosc') {
        var kolejna = { zwarta: 'swobodna', swobodna: 'ciasna', ciasna: 'zwarta' };
        var teraz = pasmo.getAttribute('data-gestosc') || 'zwarta';
        var nowa = kolejna[teraz];
        pasmo.setAttribute('data-gestosc', nowa);
        var etykieta = q('[data-gestosc-stan]', b);
        if (etykieta) { etykieta.textContent = nowa; }
        toast('Gęstość kart', 'Szerokość karty: ' + nowa + '.', 'informacja');
        return;
      }
      if (co === 'archiwalne') {
        var ukryte = b.getAttribute('aria-checked') === 'true';
        b.setAttribute('aria-checked', ukryte ? 'false' : 'true');
        toast('Karty zarchiwizowane', OPISY_WIDOKU.archiwalne, 'informacja');
        return;
      }
      toast('Ustawienia widoku', OPISY_WIDOKU[co] || co, 'informacja');
    });
  });

  /* ── Ręczna korekta grup: przeciąganie karty ───────────────────────────── */

  var przenoszona = null;

  karty.forEach(function (karta) {
    if (karta.getAttribute('data-karta-rodzaj') === 'centrum') { return; }
    karta.setAttribute('draggable', 'true');
    karta.addEventListener('dragstart', function () {
      przenoszona = karta;
      karta.setAttribute('data-przenoszona', 'tak');
    });
    karta.addEventListener('dragend', function () {
      karta.removeAttribute('data-przenoszona');
      pasmo.removeAttribute('data-cel-upuszczenia');
      przenoszona = null;
    });
    karta.addEventListener('dragover', function (e) { e.preventDefault(); pasmo.setAttribute('data-cel-upuszczenia', 'tak'); });
    karta.addEventListener('drop', function (e) {
      e.preventDefault();
      if (!przenoszona || przenoszona === karta) { return; }
      karta.parentNode.insertBefore(przenoszona, karta);
      var grupa = karta.getAttribute('data-grupa') || 'nowe zadanie';
      przenoszona.setAttribute('data-grupa', grupa);
      toast('Karta przeniesiona', 'Karta należy teraz do grupy: ' + grupa + '.', 'sukces');
    });
  });
})();
