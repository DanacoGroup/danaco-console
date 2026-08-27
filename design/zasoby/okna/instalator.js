/* ============================================================================
   INSTALATOR — złożenie okna

   Mechanika kreatora sześciu kroków: przełączanie ekranów, wykrycie procesora
   i jego trzy wyniki, przebieg zapisu z odsłonami wycofywania i błędu, pas
   działań wyprowadzany z kroku oraz okna potwierdzeń.

   Wygląd nie należy do tego pliku ani do znacznika okna — stoi w bibliotece.
   Tutaj są wyłącznie: stan kreatora, treści komunikatów zależne od stanu
   i reguły, które z czynności są w danym kroku dostępne.
   ============================================================================ */
/* Przełączanie kroków instalatora — jeden ekran naraz.
   Okno powstaje z komponentów asynchronicznie (katalog treści wczytuje się
   z pliku), więc przepływ rusza dopiero na zdarzenie `kreator-gotowy`. */
(function () {
  'use strict';
  document.addEventListener('kreator-gotowy', uruchom, { once: true });

function uruchom() {
  var K = window.DanacoKreator, tekst = K.narzedzia.tekst, podstaw = K.narzedzia.podstaw;
  var ILE = 6, biezacy = 1;
  /* Bryła stoi w komplecie od pierwszego kroku — to znak produktu, nie miara
     postępu; miarę niesie pasek na kroku 5. Z krokiem zmienia się wyłącznie
     to, czy bryła pracuje. Sam rysunek jest składnikiem biblioteki. */
  function bryla() {
    var pole = document.querySelector('[data-bryla]');
    return pole ? pole.bryla : null;
  }
  function osadzWarstwy(n) {
    var b = bryla();
    if (!b) return;
    /* Bryła prowadzi własną pętlę — okno nie podaje jej miary, bo bryła jest
       znakiem produktu, a nie wskaźnikiem postępu; miarę niesie pasek kroku 5.
       Z kroku bierze jedynie to, czy praca trwa, i błysk na zakończenie. */
    b.pracuje(n === 5);
    if (n === ILE) b.domknij();
  }

  /* Odsłony pokazywane adresem — prototyp musi dać się obejrzeć w stanach,
     do których przepływ sam nie doprowadzi: `?procesor=arm|brak`,
     `?stan=blad|wycofywanie`, `?wynik=ostrzezenia`. */
  var ADRES = new URLSearchParams(location.search);

  /* ——— Krok 5: instalowanie ————————————————————————————————————————————
     Prototyp odgrywa zapis, bo statyczny zrzut nie pokazuje tego, co w tym
     kroku rozstrzyga: krok nie ma wyboru ani drogi wstecz, jedyną czynnością
     operatora jest przerwanie, a na krok 6 instalator przechodzi sam.
     Udziały etapów są miarą czasu, nie treścią — nazwy, liczniki i pozostały
     czas stoją przy etapach w dokumencie, nie tutaj. */
  var PRZEBIEG_CZAS = 9000, WYCOFYWANIE_CZAS = 3000;
  var UDZIALY = [0.08, 0.46, 0.34, 0.12];
  var ZNAK_GOTOWE = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" '
    + 'stroke-width="2" stroke-linecap="round" stroke-linejoin="round">'
    + '<path d="m5 13 4 4L19 7"/></svg>';
  var SLOWA = (function () {
    var t = tekst('krok5.stany');
    return { gotowe: t.gotowe, 'w-toku': t.wToku, oczekuje: t.oczekuje };
  })();
  var ODSLONY = (function () {
    var m = {};
    Object.keys(K.stany.odslona).forEach(function (nazwa) {
      var o = K.stany.odslona[nazwa];
      m[nazwa] = { tytul: tekst(o.tytul), wstep: tekst(o.podtytul) };
    });
    return m;
  })();
  var przebieg = null, zegar = null;

  function ekran5() { return document.querySelector('.dn-kreator-ekran[data-ekran="5"]'); }
  function odslona() { var e = ekran5(); return (e && e.dataset.odslona) || 'przebieg'; }
  function etapy() {
    var e = ekran5();
    return e ? Array.prototype.slice.call(e.querySelectorAll('[data-etapy] > .dn-krok')) : [];
  }
  /* Stan etapu niesie znak i słowo naraz: ✓ gotowe, wskaźnik pracy w toku,
     pusty pierścień oczekuje. Sama barwa stanu nie niesie. */
  function ustawEtap(li, stan) {
    if (li.dataset.etap === stan) return;
    li.dataset.etap = stan;
    li.classList.toggle('dn-krok--poprawny', stan === 'gotowe');
    li.classList.toggle('dn-krok--pracuje', stan === 'w-toku');
    var slowo = li.querySelector('[data-etap-slowo]');
    if (slowo) slowo.textContent = SLOWA[stan];
    var znak = li.querySelector('[data-etap-znak]');
    if (!znak) return;
    if (stan === 'gotowe') znak.innerHTML = ZNAK_GOTOWE;
    else if (stan === 'w-toku') znak.innerHTML = '<span class="dn-spinner"></span>';
    else znak.textContent = '';
  }
  /* Pasek szczegółów należy do etapu oznaczonego jako w toku i liczy to, co ten
     etap właśnie przenosi. Etap, który nie liczy nic, zostawia pasek pusty —
     cudza miara pod nazwą bieżącego etapu jest sprzecznością danych. */
  function szczegolyEtapu(li, udzial) {
    var czasownik = li.dataset.licznikCzasownik;
    if (!czasownik) return { licznik: '', czas: '' };
    var cel = Number(li.dataset.licznikCel || 0);
    var jedn = li.dataset.licznikJednostka ? ' ' + li.dataset.licznikJednostka : '';
    var rzecz = li.dataset.licznikRzecz ? ' ' + li.dataset.licznikRzecz : '';
    return {
      licznik: czasownik + ' <b>' + Math.round(udzial * cel) + jedn + '</b> z <b>'
        + cel + jedn + '</b>' + rzecz,
      czas: li.dataset.pozostalo || ''
    };
  }
  function ustawMiare(proc) {
    var e = ekran5(); if (!e) return;
    var tor = e.querySelector('[data-postep-tor]');
    var pasek = tor && tor.querySelector('.dn-postep-wartosc');
    if (pasek) pasek.style.width = proc + '%';
    if (tor) tor.setAttribute('aria-valuenow', String(proc));
    var miara = e.querySelector('[data-etap-miara]');
    if (miara) miara.textContent = proc + '%';
  }
  /* Wycofywanie nie ma miary — tor idzie bez wartości, a rola `progressbar`
     traci `aria-valuenow` i to ona ogłasza stan nieokreślony. */
  function torBezMiary(wl) {
    var e = ekran5(); if (!e) return;
    var tor = e.querySelector('[data-postep-tor]');
    var miara = e.querySelector('[data-etap-miara]');
    if (tor) {
      tor.classList.toggle('dn-postep-tor--nieokreslony', wl);
      if (wl) tor.removeAttribute('aria-valuenow');
      var pasek = tor.querySelector('.dn-postep-wartosc');
      if (pasek && wl) pasek.style.width = '';
    }
    if (miara) miara.hidden = wl;
  }
  function rysujPrzebieg(p) {
    var lista = etapy();
    if (!lista.length) return;
    var suma = 0, idx = lista.length - 1, udzial = 1;
    for (var i = 0; i < lista.length; i++) {
      var w = UDZIALY[i] || 0;
      if (p < suma + w) { idx = i; udzial = w ? (p - suma) / w : 1; break; }
      suma += w;
    }
    lista.forEach(function (li, i) {
      ustawEtap(li, p >= 1 || i < idx ? 'gotowe' : (i === idx ? 'w-toku' : 'oczekuje'));
    });
    ustawMiare(Math.round(p * 100));
    var e = ekran5();
    /* Miara nad torem jest miarą CAŁOŚCI i tak jest podpisana w dokumencie.
       Numer etapu i licznik etapowy stoją razem w pasku szczegółów, więc dwie
       różne miary nigdy nie stoją obok siebie bez wyjaśnienia, czego dotyczą. */
    var sz = szczegolyEtapu(lista[idx], udzial);
    var numer = podstaw(tekst('krok5.numerEtapu'), { numer: idx + 1, ile: lista.length });
    var licznik = e.querySelector('[data-licznik-etapu]');
    var czas = e.querySelector('[data-szczegoly] [data-pozostalo]');
    if (licznik) licznik.innerHTML = sz.licznik ? numer + ' · ' + sz.licznik : numer;
    if (czas) { czas.textContent = sz.czas; czas.hidden = !sz.czas; }
  }
  function ustawOdslone(nazwa) {
    var e = ekran5(); if (!e) return;
    e.dataset.odslona = nazwa;
    var tytul = e.querySelector('[data-tytul]');
    var wstep = e.querySelector('[data-wstep]');
    if (tytul) tytul.textContent = ODSLONY[nazwa].tytul;
    if (wstep) wstep.textContent = ODSLONY[nazwa].wstep;
    /* Ekran błędu mówi wyłącznie o błędzie: miara i wykaz etapów opisują zapis,
       którego już nie ma — zmiany zostały cofnięte. */
    var postep = e.querySelector('[data-blok="postep"]');
    var kolejka = e.querySelector('[data-blok="kolejka"]');
    var szczegoly = e.querySelector('[data-szczegoly]');
    if (postep) postep.hidden = nazwa === 'blad';
    if (kolejka) kolejka.hidden = nazwa !== 'przebieg';
    if (szczegoly) szczegoly.hidden = nazwa !== 'przebieg';
    e.querySelectorAll('[data-blok="blad"]').forEach(function (b) { b.hidden = nazwa !== 'blad'; });
    torBezMiary(nazwa === 'wycofywanie');
    /* Bryła pracuje wyłącznie w czasie zapisu — wycofywanie i błąd ją wyciszają. */
    var b = bryla();
    if (b) b.pracuje(nazwa === 'przebieg');
    if (nazwa === 'wycofywanie') {
      var n = e.querySelector('[data-etap-nazwa]');
      if (n) n.textContent = ODSLONY.wycofywanie.tytul;
    }
    ustawPasDzialan();
  }
  function zatrzymajPrzebieg() {
    if (przebieg) { cancelAnimationFrame(przebieg); przebieg = null; }
    if (zegar) { clearTimeout(zegar); zegar = null; }
  }
  /* Jedna pętla dla obu przebiegów: zapisu z miarą i wycofywania bez miary.
     Przy wygaszonym ruchu obie kończą się od razu — animacja jest tu jedynym
     nośnikiem czasu, a stan końcowy niesie treść, nie ruch. */
  function odegraj(czasTrwania, naKlatce, naKoncu) {
    zatrzymajPrzebieg();
    if (matchMedia('(prefers-reduced-motion: reduce)').matches) {
      if (naKlatce) naKlatce(1);
      zegar = setTimeout(naKoncu, 600);
      return;
    }
    var start = null;
    (function klatka(t) {
      if (start === null) start = t;
      var p = Math.min((t - start) / czasTrwania, 1);
      if (naKlatce) naKlatce(p);
      if (p < 1) { przebieg = requestAnimationFrame(klatka); return; }
      przebieg = null;
      naKoncu();
    })(performance.now());
  }
  function uruchomZapis() {
    ustawOdslone('przebieg');
    odegraj(PRZEBIEG_CZAS, rysujPrzebieg, function () { pokaz(6); });
  }
  function uruchomWycofywanie() {
    ustawOdslone('wycofywanie');
    odegraj(WYCOFYWANIE_CZAS, null, wyjdz);
  }
  function pokazBlad5() {
    zatrzymajPrzebieg();
    ustawOdslone('blad');
  }


  function zgodaPrzyjeta() {
    var z = document.querySelector('[data-zgoda-licencja]');
    return !z || z.checked;
  }
  /* Pas działań jest wyprowadzany z kroku i odsłony, nie ustawiany po kawałku:
     każdy krok ma własny komplet czynności i to on rozstrzyga, których nie ma. */
  function ustawPasDzialan() {
    var wstecz = document.querySelector('[data-krok-wstecz]');
    var dalej = document.querySelector('[data-krok-dalej]');
    var kopiuj = document.querySelector('[data-blad-kopiuj]');
    if (!wstecz || !dalej) return;
    var o = odslona();
    if (kopiuj) kopiuj.hidden = !(biezacy === 5 && o === 'blad');
    wstecz.hidden = false; wstecz.setAttribute('aria-disabled', 'false');
    dalej.hidden = false;
    dalej.setAttribute('aria-disabled', 'false');

    if (biezacy === 5) {
      if (o === 'blad') {
        wstecz.textContent = tekst('dzialania.zamknij');
        dalej.textContent = tekst('dzialania.sprobujPonownie');
        return;
      }
      /* Zapis w toku nie ma czego pominąć ani dokąd się cofnąć — pas niesie
         jedną czynność, przerwanie. W czasie wycofywania nie ma i jej. */
      wstecz.textContent = tekst('dzialania.anuluj');
      wstecz.setAttribute('aria-disabled', o === 'wycofywanie' ? 'true' : 'false');
      dalej.hidden = true;
      return;
    }
    if (biezacy === ILE) {
      /* Po instalacji nie ma dokąd wracać, a wyjście bez uruchamiania jest
         drogą poprawną — stąd „Zakończ" obok czynności domyślnej. */
      wstecz.textContent = tekst('dzialania.zakoncz');
      dalej.textContent = tekst('dzialania.uruchom');
      return;
    }
    wstecz.textContent = tekst(biezacy === 1 ? 'dzialania.anuluj' : 'dzialania.wstecz');
    dalej.textContent = tekst(biezacy === 4 ? 'dzialania.instaluj' : 'dzialania.dalej');
    /* Warunek niespełniony niesie wyłączona czynność główna — komunikat błędu
       obok niej byłby powtórzeniem tego, co widać, a przed pierwszym kliknięciem
       byłby zarzutem postawionym, zanim cokolwiek się wydarzyło. */
    /* Krok 2 nie wyłącza czynności naprawdę: przycisk musi przyjąć kliknięcie,
       żeby móc odpowiedzieć, czego brakuje. `aria-disabled` niesie stan wygaszony
       i ogłasza go czytnikowi, a klikalność zostaje. Krok 3 bez wybranej wersji
       nie ma czego wyjaśniać — pytanie jeszcze nie padło. */
    dalej.setAttribute('aria-disabled',
      ((biezacy === 2 && !zgodaPrzyjeta()) || (biezacy === 3 && !wybranaWersja())) ? 'true' : 'false');
  }
  var NAWIGACJA = {
    spoczynek: tekst(K.stany.nawigacjaKroku2.spoczynek),
    brakZgody: tekst(K.stany.nawigacjaKroku2.brakZgody)
  };
  function pokazBlad(wl) {
    var f = document.querySelector('.dn-kreator-ekran[data-ekran="2"] [data-nawigacja]');
    var z = document.querySelector('[data-zgoda-licencja]');
    if (f) {
      f.textContent = wl ? NAWIGACJA.brakZgody : NAWIGACJA.spoczynek;
      f.classList.toggle('dn-kreator-nawigacja--blad', wl);
    }
    if (z) {
      z.setAttribute('aria-invalid', wl ? 'true' : 'false');
      if (wl) z.focus();
    }
  }
  document.addEventListener('change', function (e) {
    if (e.target.matches('[data-zgoda-licencja]')) {
      if (e.target.checked) pokazBlad(false);
      ustawPasDzialan();
    }
    /* Odejście od wykrytej wersji jest wyborem operatora, nie usterką — okno
       nazywa skutek, a nie odbiera drogi. Pytanie pada dopiero przy „Dalej". */
    if (e.target.matches('[name="wariant"]')) { wersjaMimoTo = false; odswiezNiezgodnosc(); }
  });

  /* Zestawienie szersze od kolumny przewija się w poziomie, więc musi dać się
     osiągnąć z klawiatury. Ogniskowalne są wyłącznie te, które faktycznie się
     przewijają — inaczej dokument zbiera kilkanaście pustych przystanków. */
  function oznaczZestawienia() {
    document.querySelectorAll('.dn-zestawienie-pole').forEach(function (z) {
      var przewija = z.scrollWidth > z.clientWidth;
      if (przewija === z.hasAttribute('tabindex')) return;
      if (przewija) {
        z.setAttribute('tabindex', '0');
        z.setAttribute('role', 'group');
        z.setAttribute('aria-label', 'Zestawienie — przewijane w poziomie');
      } else {
        z.removeAttribute('tabindex');
        z.removeAttribute('role');
        z.removeAttribute('aria-label');
      }
    });
  }

  /* Treść warunków licencji nie jest pisana w tym pliku — wstawia ją skrypt
     `zbuduj-licencje.py` z dokumentu źródłowego, między znaczniki w kroku 2. */
  (function licencja() {
    var pole = document.querySelector('[data-licencja]');
    if (!pole) return;
    /* Odsyłacz spisu przewija samo pole, nie stronę — okno instalatora nie ma
       własnego przewijania, a adres nie ma powodu zbierać kotwic. */
    pole.addEventListener('click', function (e) {
      var a = e.target.closest('.dn-spis-tresci a[href^="#"]');
      if (!a) return;
      e.preventDefault();
      var cel = pole.querySelector('[id="' + a.getAttribute('href').slice(1) + '"]');
      if (cel) cel.scrollIntoView({ block: 'start', behavior: 'smooth' });
    });
    /* Krok bywa pokazywany różnymi drogami, a ukryty nie ma szerokości.
       Obserwator rozmiaru łapie moment, w którym zestawienie ją dostaje. */
    if (window.ResizeObserver) {
      var obserwator = new ResizeObserver(oznaczZestawienia);
      pole.querySelectorAll('.dn-zestawienie-pole').forEach(function (z) { obserwator.observe(z); });
    } else {
      oznaczZestawienia();
    }
  })();

  function pokaz(n) {
    biezacy = Math.min(ILE, Math.max(1, n));
    document.querySelectorAll('.dn-kreator-ekran').forEach(function (e) {
      if (Number(e.dataset.ekran) === biezacy) e.dataset.widoczny = 'tak';
      else delete e.dataset.widoczny;
    });
    document.querySelectorAll('.dn-kreator-krok').forEach(function (k) {
      var nr = Number(k.dataset.krok);
      k.dataset.stan = nr < biezacy ? 'zrobiony' : (nr === biezacy ? 'biezacy' : 'oczekuje');
    });
    if (biezacy !== 2) pokazBlad(false);
    osadzWarstwy(biezacy);
    if (biezacy === 5) {
      var stan = ADRES.get('stan');
      if (stan === 'blad') pokazBlad5();
      else if (stan === 'wycofywanie') ustawOdslone('wycofywanie');
      else uruchomZapis();
    } else zatrzymajPrzebieg();
    if (biezacy === ILE) ustawKrok6();
    ustawPasDzialan();
    if (n === 2) oznaczZestawienia();
    /* Fokus wchodzi na pierwszą kontrolkę kroku, nie zostaje na szynie — inaczej
       operator po przejściu dalej nie wie, gdzie stoi klawiatura. */
    var ekran = document.querySelector('.dn-kreator-ekran[data-widoczny="tak"]');
    if (ekran) {
      /* Fokus trafia na kontrolkę wskazaną przez krok, a w jej braku na pierwszą
         kontrolkę formularza. Pole przewijane też jest ogniskowalne, więc nie
         może wygrać wyboru samą kolejnością w kodzie. */
      var pierwsza = ekran.querySelector('[data-fokus]')
        || ekran.querySelector('input:not([type=hidden]), button');
      if (pierwsza) pierwsza.focus({ preventScroll: true });
    }
  }
  var pyt = document.querySelector('[data-potwierdzenie]');
  /* Pytanie zależy od tego, co instalator zdążył zrobić. Przed zapisem nie ma
     czego wycofywać i okno po prostu się zamyka; w trakcie zapisu wycofywanie
     jest osobną, widoczną czynnością, więc potwierdzenie ją uruchamia, a nie
     kończy okno. Po zapisie i po błędzie nic nie stoi na szali — pytania nie ma. */
  var TRESCI = {
    przed: [tekst('dialogi.zamkniecie.tytul'), tekst('dialogi.zamkniecie.tresc'),
            tekst('dialogi.zamkniecie.zostan'), tekst('dialogi.zamkniecie.wyjdz')],
    wtrakcie: [tekst('dialogi.anulowanie.tytul'), tekst('dialogi.anulowanie.tresc'),
               tekst('dialogi.anulowanie.zostan'), tekst('dialogi.anulowanie.wyjdz')]
  };
  var poPotwierdzeniu = wyjdz;
  function zapisWToku() { return biezacy === 5 && odslona() === 'przebieg'; }
  function zapytaj(otworz) {
    if (!pyt) return;
    var t = TRESCI[zapisWToku() ? 'wtrakcie' : 'przed'];
    poPotwierdzeniu = zapisWToku() ? uruchomWycofywanie : wyjdz;
    pyt.querySelector('[data-pot-tytul]').textContent = t[0];
    pyt.querySelector('[data-pot-tresc]').textContent = t[1];
    pyt.querySelector('[data-pot-zostan]').textContent = t[2];
    pyt.querySelector('[data-pot-wyjdz]').textContent = t[3];
    pyt.hidden = !otworz;
    /* Wyjście domyślne prowadzi do dalszej pracy, nie do zamknięcia. */
    if (otworz) { var b = pyt.querySelector('[data-potwierdzenie-nie]'); if (b) b.focus(); }
  }
  function zamknijOkno() {
    /* W czasie wycofywania nie ma czynności — także tej. */
    if (biezacy === 5 && odslona() === 'wycofywanie') return;
    if (biezacy <= 4 || zapisWToku()) zapytaj(true);
    else wyjdz();
  }
  function wyjdz() {
    if (window.opener && !window.opener.closed) { window.close(); return; }
    if (document.referrer && history.length > 1) { history.back(); return; }
    location.href = '../przeplyw/centrum-dowodzenia.html';
  }

  /* ——— Okna pomocnicze kroku 3 ——————————————————————————————————————— */
  function modal(nazwa) { return document.querySelector('[data-modal="' + nazwa + '"]'); }
  function pokazModal(nazwa, otworz, fokus) {
    var m = modal(nazwa);
    if (!m) return;
    m.hidden = !otworz;
    if (otworz && fokus) { var b = m.querySelector(fokus); if (b) b.focus(); }
  }
  function pokazNiezgodnosc(otworz) {
    var m = modal('niezgodnosc');
    if (!m) return;
    if (otworz) {
      m.querySelector('[data-nz-tresc]').textContent = podstaw(tekst('dialogi.niezgodnosc.tresc'),
        { wykryty: nazwaWersji(wykryto), wybrany: nazwaWersji(wybranaWersja()) });
    }
    pokazModal('niezgodnosc', otworz, '[data-nz-popraw]');
  }
  function zamknijModale() {
    if (pyt) pyt.hidden = true;
    ['niezgodnosc', 'pomoc-procesor'].forEach(function (n) { pokazModal(n, false); });
  }

  document.addEventListener('keydown', function (e) {
    if (e.key !== 'Escape') return;
    if (pyt && !pyt.hidden) { zapytaj(false); return; }
    zamknijModale();
  });

  document.addEventListener('click', function (e) {
    /* Krzyżyk w belce pyta tym samym pytaniem co czynność poboczna. */
    if (e.target.closest('.dn-kreator-belka-btn--zamknij')) { zamknijOkno(); return; }
    if (e.target.closest('[data-potwierdzenie-nie]')) { zapytaj(false); return; }
    if (e.target.closest('[data-potwierdzenie-tak]')) {
      var czynnosc = poPotwierdzeniu;
      zapytaj(false);
      czynnosc();
      return;
    }
    if (e.target === pyt) { zapytaj(false); return; }

    if (e.target.closest('[data-pomoc-procesor]')) { pokazModal('pomoc-procesor', true, '[data-pp-zamknij]'); return; }
    if (e.target.closest('[data-pp-zamknij]') || e.target === modal('pomoc-procesor')) {
      pokazModal('pomoc-procesor', false); return;
    }
    if (e.target.closest('[data-nz-popraw]')) {
      var zgodny = document.querySelector('[name="wariant"][value="' + wykryto + '"]');
      if (zgodny) zgodny.checked = true;
      wersjaMimoTo = false;
      pokazNiezgodnosc(false);
      odswiezNiezgodnosc();
      return;
    }
    if (e.target.closest('[data-nz-mimo-to]')) {
      wersjaMimoTo = true;
      pokazNiezgodnosc(false);
      pokaz(4);
      return;
    }
    if (e.target === modal('niezgodnosc')) { pokazNiezgodnosc(false); return; }

    if (e.target.closest('[data-blad-kopiuj]')) {
      var szcz = document.querySelector('[data-blad-szczegoly]');
      if (szcz && navigator.clipboard) navigator.clipboard.writeText(szcz.textContent);
      return;
    }

    var wstecz = e.target.closest('[data-krok-wstecz]');
    if (wstecz) {
      if (wstecz.getAttribute('aria-disabled') === 'true') return;
      if (biezacy === 1 || zapisWToku()) { zamknijOkno(); return; }
      if (biezacy === 5 || biezacy === ILE) { wyjdz(); return; }
      pokaz(biezacy - 1);
      return;
    }
    var dalej = e.target.closest('[data-krok-dalej]');
    if (dalej) {
      if (biezacy === 2 && !zgodaPrzyjeta()) { pokazBlad(true); return; }
      if (biezacy === 3 && !wybranaWersja()) return;
      if (biezacy === 3 && !zgodnaWersja() && !wersjaMimoTo) { pokazNiezgodnosc(true); return; }
      if (biezacy === 5) { if (odslona() === 'blad') uruchomZapis(); return; }
      if (biezacy === ILE) { wyjdz(); return; }
      pokaz(biezacy + 1);
      return;
    }
    var k = e.target.closest('.dn-kreator-krok');
    if (k) pokaz(Number(k.dataset.krok));
  });

  /* ——— Krok 3: wersja programu ————————————————————————————————————————
     Wykrycie jest przesłanką, nie rozstrzygnięciem. Gdy procesora nie udało się
     rozpoznać, okno nie zgaduje: nie zaznacza nic i nie pokazuje plakietek. */
  var WYKRYCIA = (function () {
    var m = {};
    Object.keys(K.stany.wykrycie).forEach(function (nazwa) {
      var w = K.stany.wykrycie[nazwa];
      m[nazwa] = { wstep: tekst(w.podtytul), wskazowka: tekst(w.wskazowka) };
    });
    return m;
  })();
  var wykryto = { arm: 'arm', brak: 'brak' }[ADRES.get('procesor')] || 'x64';
  var wersjaMimoTo = false;

  function ekran3() { return document.querySelector('.dn-kreator-ekran[data-ekran="3"]'); }
  function wybranaWersja() {
    var r = document.querySelector('[name="wariant"]:checked');
    return r ? r.value : null;
  }
  function nazwaWersji(w) { return tekst(w === 'arm' ? 'krok3.arm.nazwa' : 'krok3.x64.nazwa'); }
  function zgodnaWersja() { return wykryto === 'brak' || wybranaWersja() === wykryto; }
  function odswiezNiezgodnosc() {
    var e = ekran3(); if (!e) return;
    var a = e.querySelector('[data-niezgodnosc]');
    if (a) a.hidden = !wybranaWersja() || zgodnaWersja();
    ustawPasDzialan();
  }
  function ustawKrok3() {
    var e = ekran3(); if (!e) return;
    e.dataset.wykryto = wykryto;
    e.querySelector('[data-wstep]').textContent = WYKRYCIA[wykryto].wstep;
    e.querySelector('[data-wskazowka]').textContent = WYKRYCIA[wykryto].wskazowka;
    e.querySelectorAll('[name="wariant"]').forEach(function (r) {
      var moja = r.value === wykryto;
      r.checked = wykryto !== 'brak' && moja;
      r.closest('label').querySelectorAll('[data-plakietka]').forEach(function (p) {
        p.hidden = wykryto === 'brak' || !moja;
      });
    });
    odswiezNiezgodnosc();
  }

  /* ——— Krok 6: gotowe ————————————————————————————————————————————————— */
  function ustawKrok6() {
    var e = document.querySelector('.dn-kreator-ekran[data-ekran="' + ILE + '"]');
    if (!e) return;
    var zOstrzezeniami = ADRES.get('wynik') === 'ostrzezenia';
    e.dataset.wynik = zOstrzezeniami ? 'ostrzezenia' : 'gotowe';
    var w = K.stany.wynik[zOstrzezeniami ? 'ostrzezenia' : 'gotowe'];
    e.querySelector('[data-tytul]').textContent = tekst(w.tytul);
    e.querySelector('[data-wstep]').textContent = tekst(w.podtytul);
    /* `hidden` jako właściwość stoi na HTMLElement — svg jej nie ma i cicho
       przyjmuje przypisanie bez zmiany atrybutu. Znak przełącza się atrybutem. */
    e.querySelector('[data-znak="gotowe"]').toggleAttribute('hidden', zOstrzezeniami);
    e.querySelector('[data-znak="ostrzezenia"]').toggleAttribute('hidden', !zOstrzezeniami);
    e.querySelector('[data-ostrzezenia]').hidden = !zOstrzezeniami;
    /* Wiersz procesora zamyka pętlę z krokiem 3 — pokazuje wybór, nie wykrycie. */
    var proc = e.querySelector('[data-podsumowanie-procesor]');
    if (proc) proc.textContent = nazwaWersji(wybranaWersja() || wykryto);
  }
  ustawKrok3();
  /* Bryła zakłada się na `DOMContentLoaded`, czyli po pierwszym `pokaz`.
     Jej skrypt stoi wyżej, więc gdy dojdzie do tego wpisu, pole już istnieje. */
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () { osadzWarstwy(biezacy); });
  }
  pokaz(1);
}
})();
