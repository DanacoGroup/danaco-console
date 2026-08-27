/* Moduł składa okno Centrum dowodzenia z kafli, stref oraz nawigacji startowej, tworząc główny widok programu. */
/* ══════════════════════════════════════════════════════════════════════════
   SKRYPT LOKALNY — Centrum dowodzenia v2. Rozszerza zachowania powłoki
   (menu, motyw, TRYBY) obsługiwane przez wspolne.js / prototyp.js / stanowisko.js.
   ══════════════════════════════════════════════════════════════════════════ */
(function () {
  'use strict';
  function q(s, k) { return (k || document).querySelector(s); }
  function qq(s, k) { return Array.prototype.slice.call((k || document).querySelectorAll(s)); }
  function oglos(t) { if (window.dnOglos) { window.dnOglos(t); } }
  function toast(tyt, tre, rodz, ms) {
    if (window.dnToast) { window.dnToast(tyt, tre, rodz || 'informacja', ms || 3200); }
  }

  /* ── WIERSZ PIERWSZEGO WEJŚCIA ─────────────────────────────────────── */
  var wiersz = q('#cd-start');
  var zamknijStart = q('#cd-start-zamknij');
  function ustawWskazowke(widoczna) {
    if (wiersz) { wiersz.hidden = !widoczna; }
    qq('[data-cd-wskazowka]').forEach(function (b) {
      b.setAttribute('aria-checked', widoczna ? 'true' : 'false');
    });
  }
  if (wiersz && zamknijStart) {
    zamknijStart.addEventListener('click', function () {
      ustawWskazowke(false);
      toast('Wskazówka ukryta', 'Wskazówka wraca z menu widoku w prawym rogu paska okna.', 'informacja');
      oglos('Wskazówka startowa ukryta.');
    });
  }
  qq('[data-cd-wskazowka]').forEach(function (b) {
    b.addEventListener('click', function () {
      ustawWskazowke(b.getAttribute('aria-checked') !== 'true');
    });
  });

  /* Stan samouczka niosą dwie kontrolki: przycisk pasa narzędzi (`aria-pressed`)
     i pozycja menu widoku (`aria-checked`). Przełącznik panelu ustawia tylko
     pierwszą z nich, więc drugą dostrajamy tu — inaczej ptaszek w menu kłamie. */
  qq('[data-przelacz-samouczek]').forEach(function (b) {
    b.addEventListener('click', function () {
      window.setTimeout(function () {
        var panel = q('#panel-samouczek');
        if (!panel) { return; }
        qq('[data-przelacz-samouczek][role="menuitemcheckbox"]').forEach(function (m) {
          m.setAttribute('aria-checked', panel.hidden ? 'false' : 'true');
        });
      }, 0);
    });
  });

  /* ── MENU WIDOKU · prawy róg paska okna ────────────────────────────────── */
  qq('[data-cd-odswiez]').forEach(function (b) {
    b.addEventListener('click', function () {
      toast('Widok odświeżony', 'Wykaz środowisk, komponentów i ustawień wczytany na nowo.', 'informacja');
      oglos('Widok odświeżony.');
    });
  });

  /* ── SAMOUCZEK · panel po prawej ─────────────────────────────────────────
     Przełączanie panelu obsługuje `zasoby/okna/zakladki-paneli.js` — wspólnie
     dla wyzwalacza w pasie narzędzi i dla przycisku „Otwórz samouczek" w wierszu
     pierwszego wejścia (oba noszą `data-przelacz-samouczek`). Tu nie dublujemy
     nasłuchu: dwa nasłuchy na tym samym kliknięciu zamykały panel w tej samej
     chwili, w której go otwierały. */

  /* ── PANEL SESJI · nowa sesja ──────────────────────────────────────────── */
  qq('.dn-obszar-panel--boczny [aria-label="Nowa sesja"]').forEach(function (b) {
    b.addEventListener('click', function () {
      toast('Nowa sesja', 'Nowe okno robocze na karcie Centrum dowodzenia. Sesja powstaje przy pierwszym poleceniu.', 'informacja', 3600);
    });
  });

  /* ── PASMO KART · nowa karta ───────────────────────────────────────────── */
  qq('.dn-obszar-panel--glowny .dn-karty-dodaj').forEach(function (b) {
    b.addEventListener('click', function () {
      toast('Nowa karta', 'Karta otwiera moduł wybrany w Centrum dowodzenia i zakłada grupę zadania.', 'informacja', 3600);
    });
  });

  /* ── STREFA 1 · wejście do środowiska ────────────────────────────────── */
  qq('.cd-wejdz').forEach(function (b) {
    b.addEventListener('click', function (e) {
      e.stopPropagation();
      var srod = b.getAttribute('data-wejdz');
      toast('Wejście do środowiska', 'Otwarcie przedsionka środowiska ' + srod + '. Moduł wiodący wybierasz kaflem w przedsionku.', 'sukces');
      oglos('Otwieram środowisko ' + srod + '.');
    });
  });
  qq('.dn-karta-srodowiska').forEach(function (k) {
    k.addEventListener('click', function (e) {
      /* Menu operacji leży na karcie, ale nie jest wejściem do środowiska —
         bez tego warunku każde otwarcie menu otwierało też środowisko. */
      if (e.target.closest('.cd-karta-menu')) { return; }
      var przy = q('.cd-wejdz', k);
      if (przy) { przy.click(); }
    });
  });

  /* ── STREFA 1 · menu operacji środowiska (opracowanie rozdz. 3.2) ─────── */
  var opisyKarty = {
    'nowa-karta': 'Otwarcie środowiska %s w nowej karcie sesji — dotychczasowe karty zostają nietknięte.',
    'przypnij': 'Środowisko %s przypięte: jego karta stoi pierwsza w strefie środowisk i w szynie nawigacji.',
    'wyczysc': 'Karty sesji środowiska %s zamknięte. Zadania serwerowe biegną dalej jako sesje w tle.'
  };
  qq('[data-karta-akcja]').forEach(function (b) {
    b.addEventListener('click', function (e) {
      e.stopPropagation();
      var akcja = b.getAttribute('data-karta-akcja');
      var srod = b.getAttribute('data-srodowisko');
      toast(srod, (opisyKarty[akcja] || 'Operacja %s').replace('%s', srod), akcja === 'wyczysc' ? 'ostrzezenie' : 'informacja', 3600);
      oglos(srod + ' — ' + akcja + '.');
    });
  });

  /* ── STREFA 1 · nawigacja klawiaturą po kartach ──────────────────────── */
  var karty = qq('.dn-karta-srodowiska');
  karty.forEach(function (k, i) {
    k.setAttribute('tabindex', i === 0 ? '0' : '-1');
    k.addEventListener('keydown', function (e) {
      var kier = e.key === 'ArrowRight' ? 1 : e.key === 'ArrowLeft' ? -1 : 0;
      if (kier) {
        e.preventDefault();
        var cel = karty[(i + kier + karty.length) % karty.length];
        karty.forEach(function (x) { x.setAttribute('tabindex', '-1'); });
        cel.setAttribute('tabindex', '0');
        cel.focus();
      } else if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        var przy = q('.cd-wejdz', k);
        if (przy) { przy.click(); }
      }
    });
  });

  /* ── STREFA 2 · komponenty aplikacji ─────────────────────────────────── */
  var opisyKafli = {
    Automations: 'Okno konfiguracji Automations: warunek uruchomienia, kroki i wykonawca automatyki.',
    Agents: 'Okno konfiguracji Agents: rola agenta, zakres samodzielności i granice działania.',
    Workspace: 'Okno konfiguracji projektu: nazwa, katalog roboczy, materiały wejściowe i zespół.',
    Assistant: 'Okno konfiguracji Assistant: ton wypowiedzi, zakres pamięci i stopień samodzielności.'
  };
  qq('.cd-kafel').forEach(function (b) {
    b.addEventListener('click', function () {
      var k = b.getAttribute('data-komponent');
      toast(k, opisyKafli[k] || ('Tworzenie komponentu: ' + k), 'informacja', 3600);
      oglos(opisyKafli[k] || k);
    });
  });

  var dodaj = q('#cd-dodaj-komponent');
  if (dodaj) {
    dodaj.addEventListener('click', function () {
      toast('Nowy komponent', 'Wybór rodzaju składnika, a po nim okno konfiguracji nowego komponentu własnego.', 'informacja', 3600);
      oglos('Zakładam komponent własny.');
    });
  }

  /* ── STREFA 3 · komponenty własne Operatora ──────────────────────────── */
  qq('[data-otworz-komponent]').forEach(function (b) {
    b.addEventListener('click', function () {
      var nazwa = b.getAttribute('data-otworz-komponent');
      toast(nazwa, 'Otwarcie zapisanego komponentu w jego oknie konfiguracji — bez zakładania nowego.', 'informacja');
      oglos('Otwieram ' + nazwa + '.');
    });
  });
  var opisyOperacji = {
    duplikuj: 'Kopia komponentu powstaje z tą samą definicją i nową nazwą.',
    eksport: 'Eksport zapisuje definicję komponentu do pliku przenośnego między instalacjami.',
    import: 'Import wczytuje definicję z pliku i dokłada ją do zasobów własnych konta.',
    usun: 'Usunięcie zdejmuje komponent z zasobów własnych. Sesje, w których był użyty, zostają.'
  };
  qq('[data-operacja]').forEach(function (b) {
    b.addEventListener('click', function () {
      var op = b.getAttribute('data-operacja');
      toast('Operacje komponentu', opisyOperacji[op] || op, op === 'usun' ? 'ostrzezenie' : 'informacja', 3600);
    });
  });

  /* ── STREFA 3 · listwa ustawień ──────────────────────────────────────── */
  var konf = q('#cd-konfiguracja');
  if (konf) {
    konf.addEventListener('click', function () {
      toast('Okno konfiguracji', 'Pełny zakres ustawień platformy, środowisk, sesji i modeli — trzynaście zakresów.', 'informacja');
    });
  }

  var opisySzybkie = {
    wyglad: 'Motyw, gęstość układu i wielkość pisma — zakres wyglądu okna konfiguracji.',
    jezyk: 'Język interfejsu platformy. Zmiana obejmuje wszystkie okna i funkcje globalne.',
    gestosc: 'Gęstość układu: zwarta dla pracy z danymi, swobodna dla czytania długich treści.',
    diagnostyka: 'Dzienniki procesów, podgląd stanu sesji i diagnostyka połączeń modeli.'
  };
  qq('[data-szybkie]').forEach(function (b) {
    b.addEventListener('click', function () {
      var k = b.getAttribute('data-szybkie');
      toast('Ustawienia szybkie', opisySzybkie[k] || k, 'informacja', 3600);
    });
  });

  var mobile = q('#cd-mobile');
  if (mobile) {
    mobile.addEventListener('click', function () {
      var wl = mobile.getAttribute('aria-pressed') === 'true';
      mobile.setAttribute('aria-pressed', wl ? 'false' : 'true');
      toast('Mobile', wl ? 'Funkcja globalna Mobile wyłączona.' : 'Funkcja globalna Mobile aktywna — dostęp mobilny i monitoring procesów. Nie tworzy nowej przestrzeni roboczej.', wl ? 'informacja' : 'sukces');
    });
  }

  /* Stan włączenia niesie sam przycisk (`aria-pressed`) — napis obok niego
     powtarzał tę samą informację słowem i był jedyną treścią w listwie, która
     nie była nazwą czynności. */
  var aod = q('#cd-aod-przelacz');
  if (aod) {
    aod.addEventListener('click', function () {
      var wl = aod.getAttribute('aria-pressed') === 'true';
      aod.setAttribute('aria-pressed', wl ? 'false' : 'true');
      if (!wl) {
        toast('Always On Display', 'Globalny agent towarzyszący przywołany: proaktywne doradztwo i pomoc kontekstowa nad całą platformą.', 'sukces', 3600);
      }
      oglos(wl ? 'Always On Display wyłączony.' : 'Always On Display włączony.');
    });
  }
})();
