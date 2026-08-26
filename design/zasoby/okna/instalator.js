/* ============================================================================
   DANACO CONSOLE — INSTALATOR
   ----------------------------------------------------------------------------
   Złożenie jednego okna. Wygląd stoi w `wejscie.css` (§10); tutaj mieści się
   wyłącznie mechanizm:

   1 · przejście między sześcioma krokami bez przeładowania dokumentu —
       powitanie, warunki licencji, wariant, katalog i skróty, postęp,
       zakończenie,
   2 · wariant instalacji — rozmiar, etapy zapisu i zdanie zamykające wynikają
       z wyboru wariantu. Wariant rozstrzyga zakres instalacji w całości;
       osobnego wyboru zakresu instalator nie prowadzi,
   3 · przebieg instalacji — jedna liczba postępu (0…1) prowadzi pasek
       postępu, wykaz etapów ORAZ bryłę przestrzenną kroku postępu.
       Bryła nie ma własnego cyklu: bez postępu stoi nieruchomo,
   4 · bramka zgody licencyjnej wedle zasady zero blokad — przycisk pozostaje
       klikalny, niegotowość niesie stan pola i komunikat,
   5 · jawne przewijanie bloku dłuższego niż jego okno — warunki licencji.

   Ponadto (6) otwarcie okna: gotowa grafika przestrzenna sygnetu na całej
   powierzchni okna, po której zejściu zostaje kreator. Rzecz dotyczy
   otwarcia dokumentu, nie wejścia w krok — powrót na powitanie niczego
   nie odtwarza.

   ZASADA DOMYŚLNYCH WARTOŚCI. Kreator da się ukończyć samym „Dalej”:
   wariant szybki jest zaznaczony, katalog programu wskazany, skróty
   zaznaczone. Jedynym wyborem wymagającym świadomego działania jest zgoda
   licencyjna — stoi odznaczona i nie ma domyślnej wartości.
   ============================================================================ */
(function () {
  'use strict';

  var D = document;
  var okno = D.querySelector('.in-okno');
  if (!okno) return;

  var ILE_KROKOW      = 6;
  var KROK_ZGODY      = 2;
  var KROK_WARIANTU   = 3;
  var KROK_INSTALACJI = 5;

  /* ── 1 · WARIANTY INSTALACJI ──────────────────────────────────────────────
     Wariant niesie trzy rzeczy: nazwę i rozmiar pokazywane Operatorowi,
     zdanie zamykające oraz etapy zapisu. Etapy są czynnościami instalatora,
     nie wyborem: udział etapu w całości jest proporcjonalny do jego wagi,
     więc pasek postępu i bryła idą tym samym tempem, co treść wykazu. */

  var ROZPAKOWANIE = 'Rozpakowanie pakietu';
  var KLIENT       = 'Zapis plików klienta';
  var SKROTY       = 'Skróty i integracja';
  var REJESTRACJA  = 'Rejestracja w systemie';

  var WARIANTY = {
    szybka: {
      nazwa: 'Szybka (zalecana)',
      miara: '~80 MB',
      zdanie: 'Danaco Console jest zainstalowana na tym urządzeniu. Biblioteki i narzędzia pozostają na serwerze Danaco Console.',
      etapy: [
        { nazwa: ROZPAKOWANIE, waga: 24 },
        { nazwa: KLIENT,       waga: 38 },
        { nazwa: SKROTY,       waga: 4 },
        { nazwa: REJESTRACJA,  waga: 14 }
      ]
    },
    pelna: {
      nazwa: 'Pełna',
      miara: '~450 MB',
      zdanie: 'Danaco Console jest zainstalowana wraz z kompletem bibliotek i narzędzi. Praca pozostaje na tym urządzeniu.',
      etapy: [
        { nazwa: ROZPAKOWANIE,                waga: 40 },
        { nazwa: KLIENT,                      waga: 38 },
        { nazwa: 'Zapis bibliotek i narzędzi', waga: 340 },
        { nazwa: SKROTY,                      waga: 4 },
        { nazwa: REJESTRACJA,                 waga: 28 }
      ]
    }
  };

  var wariant = 'szybka';
  var biezacy = 1;

  function etapy() { return WARIANTY[wariant].etapy; }

  function wagaCalosci() {
    return etapy().reduce(function (a, e) { return a + e.waga; }, 0) || 1;
  }

  /* ── 2 · MIARY I ZDANIE ZAMYKAJĄCE ────────────────────────────────────── */

  var polaSumy     = Array.prototype.slice.call(D.querySelectorAll('[data-suma]'));
  var podsumowanie = D.querySelector('[data-podsumowanie]');
  var wynik        = D.querySelector('[data-wynik]');

  /* Podsumowanie pasa działań pojawia się dopiero od kroku, na którym wybór
     wariantu zapadł. Na powitaniu i na warunkach pas nie wyprzedza decyzji. */
  function odswiezPodsumowanie() {
    podsumowanie.textContent = biezacy < KROK_WARIANTU
      ? ''
      : WARIANTY[wariant].nazwa + ' · ' + WARIANTY[wariant].miara;
  }

  function odswiezWariant() {
    polaSumy.forEach(function (el) { el.textContent = WARIANTY[wariant].miara; });
    wynik.textContent = WARIANTY[wariant].zdanie;
    odswiezPodsumowanie();
    zbudujPlyty();
  }

  /* Katalog wspólny jest wyborem o skutku uprawnieniowym, więc przy otwarciu
     kroku pole stoi odznaczone, a program instaluje się w katalogu konta.
     Zaznaczenie przenosi katalog programu do miejsca wspólnego — i dopiero ono
     pociąga za sobą żądanie uprawnień administratora. */
  var SCIEZKA_KONTA   = 'C:\\Users\\Operator\\AppData\\Local\\Programs\\Danaco Console';
  var SCIEZKA_WSPOLNA = 'C:\\Program Files\\Danaco Console';

  var wszystkieKonta  = D.querySelector('[data-wszystkie-konta]');
  var sciezkaProgramu = D.querySelector('[data-sciezka-programu]');

  wszystkieKonta.addEventListener('change', function () {
    sciezkaProgramu.textContent = wszystkieKonta.checked ? SCIEZKA_WSPOLNA : SCIEZKA_KONTA;
  });

  D.querySelectorAll('[data-wariant]').forEach(function (r) {
    r.addEventListener('change', function () {
      if (!r.checked) return;
      wariant = r.value;
      odswiezWariant();
    });
  });

  /* ── 3 · BRYŁA PRZESTRZENNA ───────────────────────────────────────────────
     Podstawa bryły to urządzenie, każda płyta nad nią to jeden etap zapisu.
     Liczba płyt jest liczbą etapów wariantu, więc bryła zmienia się razem
     z wyborem wariantu. Otwarcie okna (§8) nie ma z tym nic wspólnego: jest
     gotowym obrazem, nie bryłą rysowaną regułami, i nie niesie stanu. */

  var bryla      = D.querySelector('[data-bryla]');
  var brylaLekka = D.querySelector('[data-bryla-lekka]');

  function plyty(pojemnik) {
    return Array.prototype.slice.call(pojemnik.querySelectorAll('.in-plyta:not(.in-plyta--podstawa)'));
  }

  function zbudujPlyty() {
    [bryla, brylaLekka].forEach(function (pojemnik) {
      plyty(pojemnik).forEach(function (p) { p.remove(); });
      etapy().forEach(function (e, i) {
        var p = D.createElement('div');
        p.className = 'in-plyta';
        p.setAttribute('data-plyta', String(Math.min(i + 1, 5)));
        p.setAttribute('data-stan', pojemnik === brylaLekka ? 'zakonczone' : 'oczekuje');
        var siatka = D.createElement('span');
        siatka.className = 'in-plyta-siatka';
        siatka.setAttribute('aria-hidden', 'true');
        p.appendChild(siatka);
        pojemnik.appendChild(p);
      });
    });
  }

  /* ── 4 · PRZEBIEG INSTALACJI ──────────────────────────────────────────────
     Jedna liczba (`postep`, 0…1) prowadzi wszystko: szerokość paska, wartość
     dla czytnika ekranu, stan wykazu etapów, stan płyt i obrót bryły. */

  var gniazdoFaz   = D.querySelector('[data-fazy]');
  var polaCzynnosc = D.querySelector('[data-czynnosc]');
  var poleProcent  = D.querySelector('[data-procent]');
  var tor          = D.querySelector('[data-tor]');

  var CZAS_PRZEBIEGU = 4400;   /* ms — cały przebieg instalacji w prototypie */
  var przebieg = null;
  var postep = 0;

  function zbudujFazy() {
    gniazdoFaz.replaceChildren();
    etapy().forEach(function (e) {
      var li = D.createElement('li');
      li.className = 'in-faza';
      li.setAttribute('data-stan', 'oczekuje');
      var nazwa = D.createElement('span');
      nazwa.textContent = e.nazwa;
      var stan = D.createElement('span');
      stan.className = 'in-faza-stan';
      stan.textContent = 'oczekuje';
      li.appendChild(nazwa);
      li.appendChild(stan);
      gniazdoFaz.appendChild(li);
    });
  }

  function ustawPostep(wartosc) {
    postep = Math.min(1, Math.max(0, wartosc));
    okno.style.setProperty('--in-postep', String(postep));

    var procent = Math.round(postep * 100);
    poleProcent.textContent = procent + '%';
    tor.setAttribute('aria-valuenow', String(procent));

    var lista = etapy();
    var calosc = wagaCalosci();
    var zrobione = postep * calosc;
    var narosle = 0;
    var biezacaNazwa = lista.length ? lista[lista.length - 1].nazwa : '';

    var wiersze = Array.prototype.slice.call(gniazdoFaz.children);
    var tafle = plyty(bryla);

    lista.forEach(function (e, i) {
      var poczatek = narosle;
      narosle += e.waga;
      var stan = zrobione >= narosle ? 'zakonczone' : (zrobione > poczatek ? 'w-toku' : 'oczekuje');
      if (stan === 'w-toku') biezacaNazwa = e.nazwa;

      var w = wiersze[i];
      if (w) {
        w.setAttribute('data-stan', stan);
        w.lastChild.textContent =
          stan === 'zakonczone' ? 'zapisany' : (stan === 'w-toku' ? 'w toku' : 'oczekuje');
      }
      var t = tafle[i];
      if (t) t.setAttribute('data-stan', stan);
    });

    polaCzynnosc.textContent = postep >= 1 ? 'Instalacja zakończona' : biezacaNazwa;
    if (biezacy === KROK_INSTALACJI) odswiezCzynnoscZapisu();
  }

  /* W trakcie zapisu czynnością główną jest przerwanie instalacji — tak jak
     w instalatorach systemowych, gdzie „Dalej” pojawia się dopiero po
     zakończeniu kopiowania. Przycisk nigdy nie zostaje wyłączony atrybutem:
     zmienia się jego znaczenie i etykieta. */
  function odswiezCzynnoscZapisu() {
    var skonczone = postep >= 1;
    btnDalej.textContent = skonczone ? 'Dalej' : 'Przerwij';
    btnDalej.className = skonczone ? 'dn-btn dn-btn--atrament' : 'dn-btn dn-btn--niebezpieczny';
  }

  function zatrzymajPrzebieg() {
    if (przebieg === null) return;
    cancelAnimationFrame(przebieg);
    przebieg = null;
  }

  function uruchomPrzebieg() {
    zatrzymajPrzebieg();
    zbudujFazy();
    zbudujPlyty();
    ustawPostep(0);

    var start = performance.now();
    function klatka(teraz) {
      var udzial = Math.min(1, (teraz - start) / CZAS_PRZEBIEGU);
      ustawPostep(udzial);
      if (udzial < 1) przebieg = requestAnimationFrame(klatka);
      else przebieg = null;
    }
    przebieg = requestAnimationFrame(klatka);
  }

  /* ── 5 · KROKI ────────────────────────────────────────────────────────── */

  var ETYKIETY_WSTECZ = ['Anuluj', 'Wstecz', 'Wstecz', 'Wstecz', 'Anuluj', 'Wstecz'];
  var ETYKIETY_DALEJ  = ['Dalej', 'Dalej', 'Dalej', 'Instaluj', 'Dalej', 'Zakończ'];

  var ekrany  = Array.prototype.slice.call(D.querySelectorAll('.in-ekran'));
  var pozycje = Array.prototype.slice.call(D.querySelectorAll('.in-krok'));
  var btnWstecz = D.querySelector('[data-krok-wstecz]');
  var btnDalej  = D.querySelector('[data-krok-dalej]');

  function pokaz(numer) {
    biezacy = Math.min(ILE_KROKOW, Math.max(1, numer));

    ekrany.forEach(function (e) {
      if (Number(e.dataset.ekran) === biezacy) e.dataset.widoczny = 'tak';
      else delete e.dataset.widoczny;
    });

    pozycje.forEach(function (k) {
      var nr = Number(k.dataset.krok);
      k.dataset.stan = nr < biezacy ? 'zrobiony' : (nr === biezacy ? 'biezacy' : 'oczekuje');
      if (nr === biezacy) k.setAttribute('aria-current', 'step');
      else k.removeAttribute('aria-current');
    });

    /* Krok zapisu niesie jedną czynność: przerwanie, a po zakończeniu — dalej.
       Powrót w trakcie zapisu nie ma sensu, więc kontrolka poboczna odpada. */
    btnWstecz.hidden = biezacy === KROK_INSTALACJI;
    btnWstecz.textContent = ETYKIETY_WSTECZ[biezacy - 1];
    btnDalej.textContent = ETYKIETY_DALEJ[biezacy - 1];
    btnDalej.className = 'dn-btn dn-btn--atrament';

    odswiezPodsumowanie();
    odswiezPolaPrzewijane();

    if (biezacy === KROK_INSTALACJI) uruchomPrzebieg();
    else zatrzymajPrzebieg();
    if (biezacy > KROK_INSTALACJI) ustawPostep(1);
  }

  /* ── 6 · BRAMKA ZGODY LICENCYJNEJ ─────────────────────────────────────── */

  var zgoda = D.querySelector('[data-zgoda]');
  var bladZgody = D.querySelector('[data-blad-zgody]');

  /* Pojawienie się komunikatu zmienia wysokość bloku warunków, więc jego
     wskazanie przewijania trzeba przeliczyć. */
  function ukryjBladZgody() {
    bladZgody.hidden = true;
    zgoda.removeAttribute('aria-invalid');
    odswiezPolaPrzewijane();
  }

  function pokazBladZgody() {
    bladZgody.hidden = false;
    zgoda.setAttribute('aria-invalid', 'true');
    zgoda.focus();
    odswiezPolaPrzewijane();
  }

  zgoda.addEventListener('change', function () {
    if (zgoda.checked) ukryjBladZgody();
  });

  /* ── 7 · POLE PRZEWIJANE ────────────────────────────────────────────────
     Blok dłuższy niż jego okno niesie własne wskazanie. Suwak szyny
     odwzorowuje dwie miary bloku: udział widocznej części w całości
     (wysokość) i miejsce przewinięcia (położenie); ściemnienie krawędzi mówi,
     z której strony treść biegnie dalej. Mechanizm niesie jeden blok —
     warunki licencji na kroku 2. */

  var polaPrzewijane = Array.prototype.slice.call(D.querySelectorAll('[data-przewijane]'));

  function odswiezPrzewijanie(pole) {
    var tresc = pole.querySelector('[data-przewijana]');
    var suwak = pole.querySelector('[data-suwak]');
    if (!tresc.clientHeight) return;   /* krok schowany — nie ma czego mierzyć */

    /* Najpierw rozstrzygnięcie o szynie: bierze ona miejsce z prawej strony
       treści, więc miary suwaka trzeba zdjąć już po tym rozstrzygnięciu. */
    pole.setAttribute('data-zapas', tresc.scrollHeight - tresc.clientHeight > 1 ? 'tak' : 'nie');

    var widoczne = tresc.clientHeight;
    var calosc = tresc.scrollHeight;
    var udzial = Math.min(1, widoczne / calosc);
    var zapas = calosc - widoczne;
    var polozenie = zapas > 0 ? tresc.scrollTop / zapas : 0;

    suwak.style.height = (udzial * 100) + '%';
    suwak.style.top = (polozenie * (1 - udzial) * 100) + '%';

    pole.setAttribute('data-poczatek', tresc.scrollTop <= 1 ? 'tak' : 'nie');
    pole.setAttribute('data-koniec', zapas - tresc.scrollTop <= 1 ? 'tak' : 'nie');
  }

  function odswiezPolaPrzewijane() { polaPrzewijane.forEach(odswiezPrzewijanie); }

  polaPrzewijane.forEach(function (pole) {
    pole.querySelector('[data-przewijana]').addEventListener('scroll', function () {
      odswiezPrzewijanie(pole);
    });
  });
  window.addEventListener('resize', odswiezPolaPrzewijane);

  function zamknij() {
    if (window.opener && !window.opener.closed) { window.close(); return; }
    if (D.referrer && history.length > 1) { history.back(); return; }
    location.href = '../przeplyw/przeplyw-wejscia.html';
  }

  btnDalej.addEventListener('click', function () {
    if (biezacy === ILE_KROKOW) { zamknij(); return; }
    if (biezacy === KROK_INSTALACJI && postep < 1) { zamknij(); return; }   /* Przerwij */
    if (biezacy === KROK_ZGODY && !zgoda.checked) { pokazBladZgody(); return; }
    pokaz(biezacy + 1);
  });

  btnWstecz.addEventListener('click', function () {
    if (ETYKIETY_WSTECZ[biezacy - 1] === 'Anuluj') zamknij();
    else pokaz(biezacy - 1);
  });

  D.querySelector('[data-zamknij]').addEventListener('click', zamknij);

  /* ── 8 · OTWARCIE OKNA ────────────────────────────────────────────────────
     Okno otwiera się grafiką przestrzenną sygnetu na całej swojej
     powierzchni, po której zejściu przechodzi w zwykły układ kreatora.
     Rzecz dotyczy OTWARCIA DOKUMENTU, nie wejścia w krok: powrót na powitanie
     przyciskiem „Wstecz” niczego nie odtwarza, bo nie jest otwarciem okna.
     Stan otwarcia niesie `data-otwarcie` na oknie; układ rozstrzyga styl
     (wejscie.css §10).

     Ruchu nie prowadzi skrypt ani przeglądarka: grafika jest gotowym plikiem
     wyrenderowanym przestrzennie, grającym raz i stającym na ostatniej
     klatce. Skrypt rozstrzyga cztery rzeczy — który plik wczytać, kiedy
     obraz osiada na klatce końcowej, kiedy ujęcie schodzi i jak je pominąć.

     ŹRÓDŁO OBRAZU. Warianty ciemny i jasny są osobnymi renderami, nie tym
     samym obrazem przebarwionym; wybiera je motyw czynny. Przy ograniczonym
     ruchu nie wczytuje się żaden — animacji, której się nie odtworzy, nie
     wolno ściągać.

     POMIJANIE. Trzy sekundy to długo, więc pierwsze wskazanie i pierwszy
     klawisz domykają ujęcie natychmiast: obraz staje na klatce końcowej,
     a okno przechodzi w układ kreatora. Ognisko klawiatury zostaje
     nietknięte — po domknięciu stoi tam, gdzie stałoby, gdyby ujęcia
     nie było. */

  var KATALOG_OTWARCIA = '../../zasoby/marka/otwarcie/';
  /* Długość gotowej animacji: 72 klatki po 24 na sekundę. Nie jest to żeton
     ruchu, tylko miara pliku — ruch jest w nim zapisany, nie odtwarzany
     regułami. Z żetonów bierze się zejście ujęcia. */
  var KLATEK_OTWARCIA = 72;
  var KLATEK_NA_SEKUNDE = 24;
  var CZAS_OTWARCIA = KLATEK_OTWARCIA / KLATEK_NA_SEKUNDE * 1000;

  var warstwaOtwarcia = D.querySelector('[data-otwarcie-warstwa]');
  var obrazOtwarcia = D.querySelector('[data-otwarcie-obraz]');
  var bezRuchu = window.matchMedia('(prefers-reduced-motion: reduce)');
  var zegarOtwarcia = null;
  var obrazOsiadl = false;

  /* Żeton czasu stoi w sekundach albo w milisekundach. */
  function czasZetonu(nazwa) {
    var wartosc = getComputedStyle(D.documentElement).getPropertyValue(nazwa).trim();
    if (!wartosc) return 0;
    return /ms$/.test(wartosc) ? parseFloat(wartosc) : parseFloat(wartosc) * 1000;
  }

  function plikOtwarcia(koncowy) {
    var motyw = D.documentElement.dataset.theme === 'light' ? 'jasny' : 'ciemny';
    return KATALOG_OTWARCIA + 'otwarcie-' + motyw + (koncowy ? '-koniec' : '') + '.webp';
  }

  function ustalOtwarcie() {
    zegarOtwarcia = null;
    delete okno.dataset.otwarcie;
    odswiezPolaPrzewijane();
  }

  /* Klatka końcowa — płaski znak widziany na wprost — jest stanem, w którym
     ujęcie ma stanąć. Podstawia się ją jawnie, zamiast ufać, że dekoder
     dojdzie do niej co do milisekundy: pomiar odtwarzania daje 2965 ms
     w jednym silniku i 3038 ms w drugim, więc na samym zegarze pliku znak
     spoczynkowy raz bywa widziany, a raz nie. */
  function ustawKlatkeKoncowa() {
    if (obrazOsiadl) return;
    obrazOsiadl = true;
    obrazOtwarcia.src = plikOtwarcia(true);
  }

  /* Osiadanie: po odtworzeniu obraz staje na klatce końcowej i stoi jeden
     takt warstwy, żeby płaski znak dało się zobaczyć, zanim wejdzie kreator.
     Bez tego taktu przenikanie zaczyna się w chwili lądowania. */
  function osiadzOtwarcie() {
    if (okno.dataset.otwarcie !== 'gra') return;
    ustawKlatkeKoncowa();
    zegarOtwarcia = setTimeout(domknijOtwarcie, czasZetonu('--dn-czas-3'));
  }

  /* Domknięcie oddaje okno kreatorowi od razu — zejście gasi już tylko samo
     ujęcie, które przestało cokolwiek zasłaniać. */
  function domknijOtwarcie() {
    if (okno.dataset.otwarcie !== 'gra') return;
    clearTimeout(zegarOtwarcia);
    D.removeEventListener('pointerdown', domknijOtwarcie);
    D.removeEventListener('keydown', domknijOtwarcie);
    ustawKlatkeKoncowa();
    okno.dataset.otwarcie = 'gasnie';
    zegarOtwarcia = setTimeout(ustalOtwarcie, czasZetonu('--dn-czas-2'));
  }

  function otworzOkno() {
    if (!warstwaOtwarcia || !obrazOtwarcia || okno.dataset.otwarcie !== 'gra') return;

    /* Ruch ograniczony: fazy otwarcia nie ma wcale. Styl trzyma kreator na
       wierzchu od pierwszej klatki, skrypt zdejmuje sam stan, a obraz zostaje
       bez źródła — nic się nie wczytuje. */
    if (bezRuchu.matches) { delete okno.dataset.otwarcie; return; }

    /* Odliczanie zaczyna się, gdy animacja rzeczywiście rusza, a nie gdy
       zapadła decyzja o jej wczytaniu. Plik, który nie doszedł, nie może
       zatrzymać okna na pustym polu. */
    function odlicz() {
      if (okno.dataset.otwarcie !== 'gra') return;
      zegarOtwarcia = setTimeout(osiadzOtwarcie, CZAS_OTWARCIA);
    }
    obrazOtwarcia.addEventListener('load', odlicz, { once: true });
    obrazOtwarcia.addEventListener('error', domknijOtwarcie, { once: true });

    D.addEventListener('pointerdown', domknijOtwarcie);
    D.addEventListener('keydown', domknijOtwarcie);

    /* Klatka końcowa idzie z góry do zapasu przeglądarki: pominięcie ujęcia
       ma stawiać obraz od razu, a nie czekać na wczytanie. */
    new Image().src = plikOtwarcia(true);
    obrazOtwarcia.src = plikOtwarcia(false);
  }

  odswiezWariant();
  pokaz(1);
  odswiezPolaPrzewijane();
  otworzOkno();
})();
