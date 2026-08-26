/* ============================================================================
   DANACO CONSOLE — RAMA OKNA APLIKACJI (zachowanie)
   ----------------------------------------------------------------------------
   Obsługuje kontrolki belki tytułowej i pasa narzędzi. Zgodnie z zasadą zero
   blokad każda kontrolka odpowiada — czynnością albo komunikatem mówiącym,
   co się stanie i gdzie tę czynność wykonać.

   Rozwijane menu korzystają z mechaniki data-menu obsłużonej w prototyp.js.
   ============================================================================ */
(function () {
  'use strict';
  var D = document;
  function wszystkie(s, k) { return Array.prototype.slice.call((k || D).querySelectorAll(s)); }
  function toast(t, o, r) { return (window.dnToast ? window.dnToast(t, o, r) : null); }

  /* ── 1 · STEROWANIE OKNEM SYSTEMOWYM ──────────────────────────────────── */

  var OKNO = {
    'Minimalizuj okno': ['Minimalizacja okna', 'W powłoce natywnej okno schodzi do paska zadań systemu.'],
    'Maksymalizuj okno': ['Maksymalizacja okna', 'W powłoce natywnej okno wypełnia ekran; ponowne naciśnięcie przywraca rozmiar.'],
    'Zamknij okno': ['Zamknięcie okna', 'Sesje pracujące w tle trwają dalej po stronie serwera i wracają przy kolejnym uruchomieniu.']
  };
  D.addEventListener('click', function (e) {
    var b = e.target.closest('.dn-belka-btn');
    if (!b) return;
    var etykieta = b.getAttribute('aria-label');

    /* Zamknięcie okna musi zamykać okno. Dotąd przycisk pokazywał wyłącznie
       dymek o tym, co zamknięcie oznacza w powłoce natywnej — okno zostawało
       otwarte, więc kontrolka wyglądała na zepsutą. W prototypie zamknięciem
       jest powrót tam, skąd okno otwarto; okno otwarte skryptem zamyka się
       naprawdę. */
    if (etykieta === 'Zamknij okno') {
      if (window.opener && !window.opener.closed) { window.close(); return; }
      if (D.referrer && history.length > 1) { history.back(); return; }
      var dom = oknoZestawu('przeplyw/centrum-dowodzenia.html');
      location.href = dom || '../przeplyw/centrum-dowodzenia.html';
      return;
    }

    var opis = OKNO[etykieta];
    if (opis) toast(opis[0], opis[1], 'informacja');
  });

  /* ── 2 · ZWIJANIE PANELU SESJI ────────────────────────────────────────── */

  /* Nazewnictwo: „szyna nawigacji" to pionowy pas środowisk i modułów przy
     krawędzi okna; panel sesji to kolumna sesji wewnątrz przestrzeni roboczej.
     Kontrolka pasa edycji zwija tę drugą. Selektor obejmuje też
     `.dn-obszar-panel--boczny`: w oknach zbudowanych na obszarze trójkolumnowym
     panel sesji nosi tę nazwę i bez niej kontrolka nie robiła nic. */
  function panelSesji() {
    return D.querySelector('.dn-obszar-panel--boczny, .pd-szyna, .szyna, .sta-cialo > .dn-boczna, .sta-cialo > .boczna');
  }
  D.addEventListener('click', function (e) {
    var b = e.target.closest('[data-zwin-szyne]');
    if (!b) return;
    var s = panelSesji();
    if (!s) { toast('Brak panelu sesji', 'Ten widok nie ma kolumny sesji do zwinięcia.', 'informacja'); return; }
    var zwinieta = b.getAttribute('aria-pressed') === 'true';
    s.hidden = !zwinieta;
    b.setAttribute('aria-pressed', zwinieta ? 'false' : 'true');
    b.setAttribute('data-etykietka', zwinieta ? 'Zwiń panel' : 'Rozwiń panel');
    b.setAttribute('aria-label', zwinieta ? 'Zwiń panel' : 'Rozwiń panel');
  });

  /* ── 3 · NAWIGACJA HISTORII, WIDOK, ODŚWIEŻENIE ───────────────────────── */

  var CZYNNOSCI = {
    'Wstecz': ['Wstecz', 'Powrót do widoku poprzedniego w tej karcie sesji.'],
    'Do przodu': ['Do przodu', 'Przejście do widoku, z którego nastąpił powrót.'],
    'Centrum dowodzenia': ['Centrum dowodzenia', 'Powrót na stronę główną. Karty sesji trwają w tle.'],
    'Odśwież widok': ['Widok odświeżony', 'Dane okien roboczych pobrane ponownie ze źródła.'],
    'Kolejka zadań': ['Kolejka zadań', 'Zadania czekające, w toku i zakończone — wraz z ich kolejnością.'],
    'Magistrala kontekstu': ['Magistrala kontekstu', 'Podgląd tego, co moduły przekazują sobie w bieżącej sesji.'],
    'Powiadomienia': ['Powiadomienia', 'Zdarzenia platformy, zakończone zadania i prośby o zgodę.'],
    'Izolacja kontekstu': ['Izolacja kontekstu', 'Zasięg dostępu sesji do plików, sieci i pamięci projektów.'],
    'Tryb skupienia': ['Tryb skupienia', 'Okno wiodące zajmuje całą przestrzeń; pozostałe ustępują.'],
    'Pełny ekran': ['Pełny ekran', 'Okno aplikacji wypełnia ekran; powrót klawiszem F11.'],
    'Dodaj skrót': ['Dodanie skrótu', 'Wskaż komponent, który ma stanąć w szybkim wyborze szyny nawigacji.']
  };
  D.addEventListener('click', function (e) {
    var b = e.target.closest('.dn-nrz-btn[aria-label]');
    if (!b || b.hasAttribute('data-menu') || b.hasAttribute('data-zwin-szyne') || b.hasAttribute('data-przelacz-motyw')) return;
    var c = CZYNNOSCI[b.getAttribute('aria-label')];
    if (c) toast(c[0], c[1], 'informacja');
  });

  /* ── 3a · SZYNA NAWIGACJI — rozwijanie środowiska ─────────────────────── */

  /* Rozwinięte jest najwyżej jedno środowisko naraz: wybór kolejnego zwija
     poprzednie. Moduły wchodzą pod środowiskiem, przesuwając następne w dół. */
  D.addEventListener('click', function (e) {
    var b = e.target.closest('.dn-szyna-poz--srodowisko');
    if (!b) return;
    var otwarte = b.getAttribute('aria-expanded') === 'true';
    wszystkie('.dn-szyna-poz--srodowisko').forEach(function (x) {
      var grupa = D.getElementById(x.getAttribute('aria-controls'));
      var ma = (x === b) && !otwarte;
      x.setAttribute('aria-expanded', ma ? 'true' : 'false');
      /* Zwinięcie środowiska gasi jego podświetlenie: znacznik bieżącego zapala
         wejście w moduł, więc po zwinięciu listy nie ma czego wskazywać. Bez
         tego ramka sygnałowa zostawała na zwiniętej pozycji. */
      if (!ma) x.removeAttribute('data-biezace');
      if (grupa) grupa.hidden = !ma;
    });
    if (!otwarte) {
      var etyk = b.querySelector('.dn-szyna-etyk b');
      toast('Środowisko: ' + (etyk ? etyk.textContent : ''), 'Wybierz moduł z rozwiniętej listy.', 'informacja');
    }
  });

  /* Pozycje menu z data-nawiguj otwierają wskazane okno zestawu (instrukcja,
     instalator, wersja i licencja). Trasa jest względna wobec katalogu 05-okna. */
  D.addEventListener('click', function (e) {
    var n = e.target.closest('[data-nawiguj]');
    if (!n) return;
    var cel = oknoZestawu(n.getAttribute('data-nawiguj'));
    if (cel) { toast(n.textContent.trim(), 'Otwieram okno…', 'informacja'); setTimeout(function () { location.href = cel; }, 300); }
  });

  /* Strefa pracy szyny: wybór środowiska dla nowej sesji i nowego projektu
     oraz otwarcie pełnej historii sesji.

     W prototypie przejścia prowadzą do okien istniejących w zestawie: nowa
     sesja otwiera przedsionek wskazanego środowiska, nowy projekt — okno
     zakładania projektu, historia sesji — okno rejestru sesji konta. */
  function oknoZestawu(sciezka) {
    var i = location.pathname.indexOf('/05-okna/');
    if (i < 0) return null;
    return location.pathname.slice(0, i + 9) + sciezka;
  }
  function przejdz(sciezka) {
    var cel = oknoZestawu(sciezka);
    if (cel) setTimeout(function () { location.href = cel; }, 400);
  }
  D.addEventListener('click', function (e) {
    var s = e.target.closest('[data-nowa-sesja-srodowisko]');
    if (s) {
      toast('Nowa sesja — ' + s.textContent.split('  ')[0].trim(),
        'Otwiera się przedsionek wskazanego środowiska. Moduł wiodący wybierasz kaflem w przedsionku.', 'sukces');
      przejdz('srodowiska/' + s.getAttribute('data-nowa-sesja-srodowisko') + '-przedsionek.html');
      return;
    }
    var pr = e.target.closest('[data-nowy-projekt-srodowisko]');
    if (pr) {
      toast('Nowy projekt — ' + pr.textContent.split('  ')[0].trim(),
        'Otwiera się okno Workspace wskazanego środowiska w trybie zakładania projektu.', 'sukces');
      przejdz('platformowe/nowy-projekt.html');
      return;
    }
    if (e.target.closest('[data-otwarz-historie]')) {
      toast('Historia sesji',
        'Rejestr wszystkich sesji konta w trzech zbiorach: czynne, zakończone i archiwalne.', 'informacja');
      przejdz('platformowe/historia-sesji.html');
    }
  });

  /* Szybki wybór — pozycje komponentów i dodawanie kolejnych. */
  D.addEventListener('click', function (e) {
    var s = e.target.closest('[data-skrot-komponent]');
    if (s) {
      toast(s.getAttribute('data-skrot-komponent'),
        'Komponent otwiera się w bieżącym środowisku jako kolejna karta sesji.', 'informacja');
      return;
    }
    if (e.target.closest('[data-dodaj-skrot]')) {
      toast('Dodanie skrótu',
        'Wskaż komponent, który ma stanąć w szybkim wyborze szyny nawigacji. Kolejność zmienia się przeciągnięciem.',
        'informacja');
    }
  });

  /* Wskazanie modułu bieżącego — jedno na całą szynę. Podświetlenie sygnałowe
     należy wyłącznie modułowi, w którym Operator faktycznie pracuje: naciśnięcie
     modułu wchodzi w niego i zapala go na sygnałowo, a jego środowisko zostaje
     rozwinięte i oznaczone jako bieżące. */
  function oznaczModul(m) {
    wszystkie('.dn-szyna-poz--modul').forEach(function (x) {
      if (x === m) x.setAttribute('aria-current', 'true'); else x.removeAttribute('aria-current');
    });
    var grupa = m ? m.closest('.dn-szyna-moduly') : null;
    wszystkie('.dn-szyna-poz--srodowisko').forEach(function (s) {
      var jego = grupa && D.getElementById(s.getAttribute('aria-controls')) === grupa;
      s.setAttribute('aria-expanded', jego ? 'true' : 'false');
      s.setAttribute('data-biezace', jego ? 'true' : 'false');
      var g = D.getElementById(s.getAttribute('aria-controls'));
      if (g) g.hidden = !jego;
    });
  }
  D.addEventListener('click', function (e) {
    var m = e.target.closest('.dn-szyna-poz--modul');
    if (m) { oznaczModul(m); return; }
    /* Wyjście z modułu — powrót na stronę główną albo do przedsionka: żaden
       moduł nie jest już bieżący, więc środowiska zwijają się, a podświetlenie
       gaśnie. Rozwinięte pozostaje tylko to środowisko, które ma moduł bieżący,
       a takiego po wyjściu nie ma. */
    var wyjscie = e.target.closest('[aria-label="Centrum dowodzenia"], [data-wyjscie-modulu]');
    if (wyjscie) {
      wszystkie('.dn-szyna-poz--modul').forEach(function (x) { x.removeAttribute('aria-current'); });
      wszystkie('.dn-szyna-poz--srodowisko').forEach(function (s) {
        s.setAttribute('aria-expanded', 'false');
        s.removeAttribute('data-biezace');
        var g = D.getElementById(s.getAttribute('aria-controls'));
        if (g) g.hidden = true;
      });
    }
  });

  /* ── 3b · POŁOŻENIE ETYKIETY POZYCJI SZYNY ────────────────────────────── */

  /* Strefy szyny przewijają się, a kontener przewijany przycina wszystko, co
     z niego wystaje. Etykieta ma więc położenie względem okna, nadawane w
     chwili najechania: dzięki temu wysuwa się poza szynę niezależnie od
     przewinięcia i niezależnie od wysokości okna. */
  /* Etykieta stoi na wysokości pozycji, ale nie wolno jej wyjść poza okno:
     pozycja pierwsza dotyka górnej krawędzi, ostatnia — dolnej, a panel jest
     od nich wyższy. Dociskamy więc jego środek do zakresu, w którym cały panel
     mieści się w oknie. Wysokość panelu trzeba zmierzyć przy odsłoniętym
     układzie: `transform` nie zmienia prostokąta układu, więc `offsetHeight`
     jest wiarygodne także przy kryciu 0. */
  var MARGINES_ETYK = 8;

  function ustawEtykiete(poz) {
    var e = poz.querySelector('.dn-szyna-etyk');
    if (!e) return;
    var r = poz.getBoundingClientRect();
    var wys = e.offsetHeight || 0;
    var polowa = wys / 2;
    var gora = MARGINES_ETYK + polowa;
    var dol = window.innerHeight - MARGINES_ETYK - polowa;
    var srodek = r.top + r.height / 2;
    /* Okno niższe od etykiety: wtedy zakres jest pusty i środek idzie na
       środek okna — panel zostaje w całości widoczny mimo obcięcia treści. */
    if (dol < gora) { srodek = window.innerHeight / 2; }
    else if (srodek < gora) { srodek = gora; }
    else if (srodek > dol) { srodek = dol; }
    e.style.left = (r.right + MARGINES_ETYK) + 'px';
    e.style.top = srodek + 'px';
  }
  ['mouseenter', 'focusin'].forEach(function (zdarzenie) {
    D.addEventListener(zdarzenie, function (e) {
      var poz = e.target.closest ? e.target.closest('.dn-szyna-poz') : null;
      if (poz) ustawEtykiete(poz);
    }, true);
  });

  /* ── 3b' · SYGNAŁ PRZEWIJANIA STREFY ŚRODOWISK ────────────────────────── */

  /* Przewija się wyłącznie strefa środowisk (rama.css). Krawędź, za którą jest
     dalszy ciąg wykazu, ma gasnąć — inaczej pozycja urywa się w pół ikony i
     nic nie mówi Operatorowi, że wykaz biegnie dalej. */
  var strefy = wszystkie('.dn-szyna-nawigacji-lista, .dn-szyna-skroty');

  /* Strefa przewijana ma się kończyć na całej pozycji. Wysokość z zginania
     wypada zwykle w połowie ikony, a wtedy reszta pola czyta się jak szpara
     między strefami. Docinamy więc wysokość w dół do wielokrotności rytmu
     (pozycja + przerwa strefy); reszta wraca do wspólnego nadmiaru, który leży
     pod ostatnią strefą, przed stopką. */
  function rytmStrefy(strefa) {
    strefa.style.removeProperty('--szyna-strefa-wys');
    var poz = strefa.querySelector('.dn-szyna-poz');
    if (!poz) return;
    var przerwa = parseFloat(window.getComputedStyle(strefa).rowGap) || 0;
    var krok = poz.getBoundingClientRect().height + przerwa;
    if (krok <= 0) return;
    var wolne = strefa.clientHeight;
    if (strefa.scrollHeight <= wolne + 1) return;      /* mieści się w całości */
    var ile = Math.max(1, Math.floor((wolne + przerwa) / krok));
    strefa.style.setProperty('--szyna-strefa-wys', (ile * krok - przerwa) + 'px');
  }

  function znakPrzewijania(strefa) {
    var zapas = strefa.scrollHeight - strefa.clientHeight;
    if (zapas <= 1) { strefa.removeAttribute('data-przewijanie'); return; }
    var stan = [];
    if (strefa.scrollTop > 1) stan.push('gora');
    if (strefa.scrollTop < zapas - 1) stan.push('dol');
    strefa.setAttribute('data-przewijanie', stan.join(' '));
  }

  /* Docięcie zmienia podział miejsca, więc obie strefy przechodzą przez rytm
     po kolei: najpierw środowiska (ustępują pierwsze), potem szybki wybór. */
  function odswiezStrefy() {
    strefy.forEach(rytmStrefy);
    strefy.forEach(znakPrzewijania);
  }

  strefy.forEach(function (strefa) {
    strefa.addEventListener('scroll', function () { znakPrzewijania(strefa); }, { passive: true });
  });
  if (strefy.length) {
    window.addEventListener('resize', odswiezStrefy);
    D.addEventListener('click', function (e) {
      if (e.target.closest && e.target.closest('.dn-szyna-poz--srodowisko')) {
        window.setTimeout(odswiezStrefy, 0);
      }
    }, true);
    odswiezStrefy();
  }

  /* ── 3c · FILTRY WYSZUKIWANIA ─────────────────────────────────────────── */

  /* Ikona filtrów stoi wewnątrz pola wyszukiwania i rozwija listę zakresu:
     środowiska oraz rodzaj wyniku. Ikona niesie znak zawężenia, żeby Operator
     widział, że wynik jest ograniczony, zanim zdziwi się brakiem trafienia. */
  function zakresFiltru() { return D.getElementById('szukaj-zakres'); }

  function odswiezZnakFiltru() {
    var lista = zakresFiltru();
    var ikona = D.querySelector('[data-szukaj-filtr]');
    if (!lista || !ikona) return;
    var pola = wszystkie('input[type="checkbox"]', lista);
    var zawezony = pola.some(function (p) { return !p.checked; });
    ikona.setAttribute('data-zawezony', zawezony ? 'tak' : 'nie');
  }

  function zwinFiltr() {
    var lista = zakresFiltru();
    var ikona = D.querySelector('[data-szukaj-filtr]');
    if (!lista || !ikona) return;
    lista.hidden = true;
    ikona.setAttribute('aria-expanded', 'false');
  }

  D.addEventListener('click', function (e) {
    var wyzw = e.target.closest('[data-szukaj-filtr]');
    var lista = zakresFiltru();
    if (!lista) return;

    if (wyzw) {
      var otwarte = wyzw.getAttribute('aria-expanded') === 'true';
      lista.hidden = otwarte;
      wyzw.setAttribute('aria-expanded', otwarte ? 'false' : 'true');
      return;
    }

    if (e.target.closest('[data-zakres-zamknij]')) { zwinFiltr(); return; }

    if (e.target.closest('[data-zakres-wyczysc]')) {
      wszystkie('input[type="checkbox"]', lista).forEach(function (p) { p.checked = true; });
      odswiezZnakFiltru();
      toast('Filtry wyczyszczone', 'Wyszukiwanie obejmuje wszystkie środowiska i wszystkie rodzaje wyników.', 'sukces');
      return;
    }

    /* Kliknięcie poza listą zwija ją — tak jak każde menu ramy. */
    if (!lista.hidden && !e.target.closest('#szukaj-zakres')) zwinFiltr();
  });

  /* „Wszystkie środowiska" prowadzi pozostałe pozycje środowisk. */
  D.addEventListener('change', function (e) {
    var lista = zakresFiltru();
    if (!lista || !e.target.closest('#szukaj-zakres')) return;

    if (e.target.hasAttribute('data-zakres-wszystkie')) {
      wszystkie('[data-zakres-srod]', lista).forEach(function (p) { p.checked = e.target.checked; });
    } else if (e.target.hasAttribute('data-zakres-srod')) {
      var srod = wszystkie('[data-zakres-srod]', lista);
      var wszystkieSrod = lista.querySelector('[data-zakres-wszystkie]');
      if (wszystkieSrod) wszystkieSrod.checked = srod.every(function (p) { return p.checked; });
    }
    odswiezZnakFiltru();
  });

  /* ── 4 · PRZYPISANIE SKRÓTU KLAWISZOWEGO ──────────────────────────────── */

  var MODYFIKATORY = { Control: 'Ctrl', Shift: 'Shift', Alt: 'Alt', Meta: 'Meta' };

  function nasluchujSkrotu(pole) {
    pole.addEventListener('focus', function () {
      pole.dataset.poprzedni = pole.value;
      pole.value = 'naciśnij kombinację…';
    });
    pole.addEventListener('blur', function () {
      if (pole.value === 'naciśnij kombinację…') pole.value = pole.dataset.poprzedni || '';
    });
    pole.addEventListener('keydown', function (e) {
      e.preventDefault();
      if (e.key === 'Escape') { pole.value = pole.dataset.poprzedni || ''; pole.blur(); return; }
      if (MODYFIKATORY[e.key]) return;
      var czlony = [];
      if (e.ctrlKey) czlony.push('Ctrl');
      if (e.shiftKey) czlony.push('Shift');
      if (e.altKey) czlony.push('Alt');
      if (e.metaKey) czlony.push('Meta');
      czlony.push(e.key.length === 1 ? e.key.toUpperCase() : e.key);
      var nowy = czlony.join(' + ');
      var zajety = wszystkie('[data-skrot]').filter(function (p) { return p !== pole && p.value === nowy; })[0];
      if (zajety) {
        toast('Kombinacja zajęta', nowy + ' — przypisana do czynności „' +
          zajety.closest('.dn-menu-skrot').querySelector('.dn-menu-skrot-etykieta').textContent + '".', 'ostrzezenie');
        return;
      }
      pole.value = nowy;
      pole.dataset.poprzedni = nowy;
      toast('Skrót przypisany', nowy, 'sukces');
      pole.blur();
    });
  }

  /* ── 5 · POZYCJE MENU BEZ WŁASNEJ CZYNNOŚCI W PROTOTYPIE ──────────────── */

  D.addEventListener('click', function (e) {
    var p = e.target.closest('.dn-nrz-menu .sta-menu-poz');
    if (!p) return;
    var rola = p.getAttribute('role');
    if (rola === 'menuitemcheckbox') {
      var wl = p.getAttribute('aria-checked') === 'true';
      p.setAttribute('aria-checked', wl ? 'false' : 'true');
      return;
    }
    if (rola === 'menuitemradio') {
      var grupa = p.parentNode.querySelectorAll('[role="menuitemradio"]');
      Array.prototype.forEach.call(grupa, function (g) { g.setAttribute('aria-checked', g === p ? 'true' : 'false'); });
      return;
    }
    var nazwa = p.textContent.replace(/\s+/g, ' ').trim();
    toast(nazwa, 'Czynność opisana w dokumentacji platformy; w prototypie pokazana jest wyłącznie ścieżka wywołania.', 'informacja');
  });

  /* ── 5a · OKNO „DOSTOSUJ WSTĄŻKĘ” ─────────────────────── */

  /* Przenoszenie składników między trzema kolumnami — poza paskami, wstążka
     pozioma, szyna pionowa — działa naprawdę: prototyp ma pokazywać mechanikę
     dostosowania, nie jej atrapę. */

  var NAZWY_LIST = { dostepne: 'Poza paskami', wstazka: 'Wstążka pozioma', szyna: 'Szyna pionowa' };

  function listaDost(rodzaj) { return D.querySelector('[data-dost-lista="' + rodzaj + '"]'); }
  function wskazana(lista) { return lista ? lista.querySelector('[aria-selected="true"]') : null; }

  function wskaz(poz) {
    var lista = poz.closest('.dn-dost-lista');
    wszystkie('.dn-dost-poz', lista).forEach(function (x) {
      x.setAttribute('aria-selected', x === poz ? 'true' : 'false');
    });
  }

  D.addEventListener('click', function (e) {
    var poz = e.target.closest('.dn-dost-poz');
    if (poz) { wskaz(poz); return; }

    var przenies = e.target.closest('[data-dost-przenies]');
    if (przenies) {
      var zK = przenies.getAttribute('data-dost-z');
      var doK = przenies.getAttribute('data-dost-do');
      var zrodlo = listaDost(zK), cel = listaDost(doK);
      if (!zrodlo || !cel) return;
      var p = wskazana(zrodlo);
      if (!p) {
        toast('Nie wskazano składnika', 'Wskaż pozycję w kolumnie „' + NAZWY_LIST[zK] + '”.', 'informacja');
        return;
      }
      cel.appendChild(p);
      wskaz(p);
      toast(p.querySelector('.dn-dost-poz-nazwa').textContent,
        'Składnik przeniesiony do kolumny „' + NAZWY_LIST[doK] + '”.', 'sukces');
      return;
    }

    var porzadek = e.target.closest('[data-dost-gora], [data-dost-dol]');
    if (porzadek) {
      var lista = listaDost(porzadek.getAttribute('data-dost-lista'));
      var w = wskazana(lista);
      var wGore = porzadek.hasAttribute('data-dost-gora');
      if (!w) {
        toast('Nie wskazano składnika', 'Kolejność ustawiasz na pozycji wskazanej w tej kolumnie.', 'informacja');
        return;
      }
      var sasiad = wGore ? w.previousElementSibling : w.nextElementSibling;
      if (!sasiad) {
        toast('Pozycja skrajna', wGore ? 'Składnik stoi już na początku.' : 'Składnik stoi już na końcu.', 'informacja');
        return;
      }
      if (wGore) lista.insertBefore(w, sasiad); else lista.insertBefore(sasiad, w);
      return;
    }

    var inne = e.target.closest('[data-dost-domyslne], [data-dost-zastosuj]');
    if (!inne) return;

    if (inne.hasAttribute('data-dost-domyslne')) {
      var okno = inne.closest('dialog');
      if (okno && okno.dataset.uklad) {
        okno.querySelector('.dn-dost-przenoszenie').innerHTML = okno.dataset.uklad;
        wszystkie('input[type="radio"], input[type="checkbox"]', okno).forEach(function (pole) {
          pole.checked = pole.defaultChecked;
        });
      }
      toast('Układ domyślny przywrócony',
        'Skład obu pasów i ustawienia kart sesji wróciły do wartości fabrycznych.', 'sukces');
      return;
    }

    var naWstazce = (listaDost('wstazka') || { children: [] }).children.length;
    var naSzynie = (listaDost('szyna') || { children: [] }).children.length;
    var d = inne.closest('dialog');
    if (d && d.close) d.close();
    toast('Układ zapisany',
      'Na wstążce poziomej: ' + naWstazce + ' składników; na szynie pionowej: ' + naSzynie +
      '. Ustawienie należy do konta.', 'sukces');
  });

  /* Klawiatura: pozycję listy wskazuje się także spacją i enterem. */
  D.addEventListener('keydown', function (e) {
    var poz = e.target.closest ? e.target.closest('.dn-dost-poz') : null;
    if (!poz) return;
    if (e.key === ' ' || e.key === 'Enter') { e.preventDefault(); wskaz(poz); }
  });

  /* ── 6 · URUCHOMIENIE ─────────────────────────────────────────────────── */

  function odswiezStan() {
    var pole = D.querySelector('[data-stan-motyw]');
    if (pole) pole.textContent = D.documentElement.dataset.theme === 'dark' ? 'ciemny' : 'jasny';
  }

  function start() {
    wszystkie('[data-skrot]').forEach(nasluchujSkrotu);
    odswiezZnakFiltru();
    /* Układ fabryczny zapamiętany przy starcie — służy przywróceniu domyślnych. */
    var oknoDost = D.getElementById('dostosuj-wstazke');
    if (oknoDost) {
      var przen = oknoDost.querySelector('.dn-dost-przenoszenie');
      if (przen) oknoDost.dataset.uklad = przen.innerHTML;
    }
    odswiezStan();
    D.addEventListener('click', function (e) {
      if (e.target.closest('[data-przelacz-motyw]')) setTimeout(odswiezStan, 0);
    });
    /* Ctrl/Cmd + K ustawia kursor w polu wyszukiwania pasa narzędzi. */
    D.addEventListener('keydown', function (e) {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
        var pole = D.querySelector('.dn-narzedzia-szukaj input');
        if (pole) { e.preventDefault(); pole.focus(); pole.select(); }
      }
    });
  }
  if (D.readyState === 'loading') D.addEventListener('DOMContentLoaded', start); else start();
})();
