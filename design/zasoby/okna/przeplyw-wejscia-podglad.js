/* Plik dostarcza rusztowanie podglądu przepływu wejścia: przejścia między etapami, dołączenie ramy aplikacji po wydaniu tokenu oraz symulację cyklu uzgodnienia protokołu. ============================================================================
   OPRAWA PODGLĄDU PRZEPŁYWU WEJŚCIA

   Rusztowanie, na którym ogląda się okna przepływu: przejścia między etapami
   i stanami, dołączenie ramy aplikacji po wydaniu tokenu oraz symulacja cyklu
   uzgodnienia protokołu z dziennikiem komunikatów.

   Nic z tego nie należy do produktu. Mechanika samych okien stoi w
   `zasoby/okna/przeplyw-wejscia.js`; tutaj wywołuje się ją przez `dnPrzejdz`.
   ============================================================================ */
(function () {
  'use strict';
  var D = document;
  function wszystkie(s, k) { return Array.prototype.slice.call((k || D).querySelectorAll(s)); }

  /* ── 1 · PRZEJŚCIA MIĘDZY ETAPAMI I STANAMI ─────────────────────────────
     `data-przelacz` niesie wyłącznie zakładka (role="tab") — prototyp.js
     ustawia na takich elementach `aria-selected`, a ten atrybut wolno nosić
     tylko zakładce. Kontrolki wewnątrz okna prowadzą do stanu przez
     `data-idz`, więc nie dostają cudzej roli. */

  var PODPISY = {
    laczenie: 'etap 1 · connection.hello → connection.hello.ack',
    uwierzytelnienie: 'etap 2 · rejestracja · potwierdzenie · logowanie · odzyskiwanie',
    przygotowanie: 'etap 3 · między uwierzytelnieniem a stroną główną'
  };

  function odswiezZakladki(grupa) {
    /* Roving tabindex: w kolejności fokusu stoi wyłącznie zakładka wybrana,
       między zakładkami przenoszą strzałki. prototyp.js ustawia to przy
       starcie; po przełączeniu trzeba odświeżyć. */
    wszystkie('[role="tablist"]').forEach(function (lista) {
      var zakladki = wszystkie('[role="tab"]', lista);
      if (grupa && !zakladki.some(function (z) { return z.getAttribute('data-grupa') === grupa; })) return;
      zakladki.forEach(function (z) {
        z.setAttribute('tabindex', z.getAttribute('aria-selected') === 'true' ? '0' : '-1');
      });
    });
  }

  /* Pigułki „Zaloguj się / Zarejestruj się" należą do okna, nie do rusztowania.
     prototyp.js zeruje `aria-selected` wszystkim przełącznikom grupy „stan”,
     a stan „dane niezgodne” czy „odzyskiwanie” nie ma własnej pigułki — bez
     tego kroku żadna nie zostawała wskazana. Panel niesie swój wybór
     w `data-pigulka`. */
  function odswiezPigulki() {
    wszystkie('[data-pigulka]').forEach(function (panel) {
      var chce = panel.getAttribute('data-pigulka');
      wszystkie('.dn-zakladki--pigulki [role="tab"]', panel).forEach(function (t) {
        var akt = t.getAttribute('data-przelacz') === chce;
        t.setAttribute('aria-selected', akt ? 'true' : 'false');
        t.setAttribute('tabindex', akt ? '0' : '-1');
        t.classList.toggle('dn-zakladka--wybrana', akt);
      });
    });
  }

  function idz(nazwa, grupa) {
    if (!window.dnPrzelaczWidok) return;
    window.dnPrzelaczWidok(nazwa, grupa);
    odswiezZakladki(grupa);
    odswiezPigulki();
    if (grupa === 'etap') poEtapie(nazwa);
  }

  /* Szew dla mechaniki okna: przejście między etapami wywołuje się z zewnątrz
     tak samo, jak wywołuje je kliknięcie. Skrypt okna nie zna wnętrza podglądu
     i nie musi go znać. */
  window.dnPrzejdz = idz;

  function poEtapie(nazwa) {
    var podpis = D.querySelector('[data-podpis-etapu]');
    if (podpis && PODPISY[nazwa]) podpis.textContent = PODPISY[nazwa];
    if (nazwa === 'przygotowanie') { zamontujRame(); odswiezSkladniki(); }
  }

  /* Etap trzeci składa okno z kopii znacznika, a kopia niesie martwe płótna bez
     wiązania. Składniki trzeba więc założyć jeszcze raz — i przez chwilę
     ponawiać, bo powłoka wstawia treść własnym rytmem. */
  function odswiezSkladniki() {
    var prob = 0;
    (function zaloz() {
      wszystkie('[data-powloki]').forEach(function (pole) {
        if (!pole.powloki && window.DanacoPowloki) window.DanacoPowloki.zaloz(pole);
      });
      if (++prob < 30) setTimeout(zaloz, 60);
    })();
  }

  D.addEventListener('click', function (e) {
    var el = e.target.closest('[data-idz]');
    if (!el) return;
    if (el.tagName === 'A') e.preventDefault();
    idz(el.getAttribute('data-idz'), el.getAttribute('data-idz-grupa') || null);
  });

  /* Zakładki przełącza prototyp.js; tutaj domyka się skutki uboczne. Nasłuch
     prototypu wpina się dopiero na DOMContentLoaded, a ten skrypt wykonuje się
     wcześniej — bez odłożenia czytalibyśmy `aria-selected` sprzed przełączenia. */
  D.addEventListener('click', function (e) {
    var z = e.target.closest('[data-przelacz][data-grupa]');
    if (!z) return;
    var grupa = z.getAttribute('data-grupa');
    setTimeout(function () {
      odswiezZakladki(grupa);
      odswiezPigulki();
      if (grupa === 'etap') poEtapie(z.getAttribute('data-przelacz'));
    }, 0);
  });

  /* ── 2 · RAMA APLIKACJI PO WYDANIU TOKENU ───────────────────────────────
     Etapy 1 i 2 nie wpinają `rama.css` ani `rama.js`. Rama wchodzi dopiero
     tutaj: arkusz, powłoka (szyna, belka, wstążka, pasek stanu) i jej
     zachowanie dołączane są skryptem, a okno przygotowania przenosi się
     do gniazda powłoki. */

  var ramaWToku = false;

  function wczytajArkusz(sciezka, po) {
    var l = D.createElement('link');
    l.rel = 'stylesheet';
    l.href = sciezka;
    l.addEventListener('load', po);
    l.addEventListener('error', po);
    D.head.appendChild(l);
  }
  function wczytajSkrypt(sciezka, po) {
    var s = D.createElement('script');
    s.src = sciezka;
    s.onload = po;
    s.onerror = po;
    D.body.appendChild(s);
  }

  function zamontujRame() {
    if (ramaWToku || D.getElementById('dn-powloka-montaz')) return;
    var gniazdo = D.getElementById('gniazdo-ramy');
    var oprawa = D.getElementById('oprawa-przygotowania');
    if (!gniazdo || !oprawa) return;
    ramaWToku = true;

    /* powloka.js czyta treść okna z <template id="dn-tresc-okna">
       i wstawia ją między szynę nawigacji a pasek stanu. */
    var tpl = D.createElement('template');
    tpl.id = 'dn-tresc-okna';
    tpl.innerHTML = oprawa.innerHTML;
    D.body.appendChild(tpl);

    var montaz = D.createElement('div');
    montaz.className = 'pt-powloka pt-powloka--jednolita';
    montaz.id = 'dn-powloka-montaz';
    gniazdo.replaceChildren(montaz);

    wczytajArkusz('../../zasoby/rama.css', function () {
      wczytajArkusz('../../zasoby/okna-modalne.css', function () {
        wczytajSkrypt('../../zasoby/powloka.js', function () {
          wczytajSkrypt('../../zasoby/okna-modalne.js', function () {
            wczytajSkrypt('../../zasoby/rama.js', function () {
              if (window.dnToast) {
                window.dnToast('Token wydany urządzeniu',
                  'Rama aplikacji jest od tej chwili dostępna.', 'sukces');
              }
            });
          });
        });
      });
    });
  }

  /* ── 3 · SYMULACJA CYKLU UZGODNIENIA PROTOKOŁU (ETAP 1) ─────────────────── */

  var log        = D.getElementById('log-uzgodnienia');
  var postep     = D.getElementById('postep-uzgodnienia');
  var pasek      = postep.querySelector('.dn-postep-wartosc');
  var licznik    = D.querySelector('[data-postep-licznik]');
  var btnOdtworz = D.getElementById('btn-odtworz');
  var etykietaOdtworz = D.querySelector('[data-etykieta-odtworz]');
  var plakietka  = D.getElementById('plakietka-stanu');
  var fazy       = wszystkie('.pt-faza');

  /* kolejność pól: wersja protokołu · identyfikator urządzenia · token */
  var HELLO = [
    { typ: 'zachęta', tekst: '» connection.hello' },
    { typ: 'para', klucz: 'wersja protokołu',    wartosc: 'dnp/1.0' },
    { typ: 'para', klucz: 'identyfikator urz.',  wartosc: 'urzadzenie-operatora-01' }
  ];
  var HELLO_OGON = [
    { typ: 'para', klucz: 'klient',              wartosc: 'Danaco Console 2.0 · powłoka Tauri' },
    { typ: 'para', klucz: 'kanał',               wartosc: 'WebSocket · adres z pakietu klienckiego' }
  ];

  var PRZEBIEGI = {
    'w-laczenie': {
      hello: [ { typ: 'para', klucz: 'token dostępu', wartosc: 'brak — pierwsze połączenie tego urządzenia' } ],
      ack: [
        { typ: 'zachęta', tekst: '« connection.hello.ack' },
        { typ: 'para', klucz: 'wersja serwera',   wartosc: 'Danaco Console 2.0' },
        { typ: 'para', klucz: 'faza wdrożenia',   wartosc: 'budowa' },
        { typ: 'para', klucz: 'uwierzytelnienie', wartosc: 'brak ważnego tokenu' },
        { typ: 'ok',   tekst: '  kanał otwarty — przejście do etapu 2: okno rejestracji i logowania' }
      ],
      stanFaz: ['gotowa', 'gotowa', 'gotowa', 'pracuje'],
      metaFaz: ['otwarty', 'dnp/1.0', 'brak tokenu', 'etap 2'],
      plakietka: 'kanał otwarty',
      klasaPlakietki: 'dn-plakietka dn-plakietka--sygnal',
      tetno: true
    },
    'w-token': {
      hello: [ { typ: 'para', klucz: 'token dostępu', wartosc: 'obecny — token urządzenia' } ],
      ack: [
        { typ: 'zachęta', tekst: '« connection.hello.ack' },
        { typ: 'para', klucz: 'wersja serwera',   wartosc: 'Danaco Console 2.0' },
        { typ: 'para', klucz: 'faza wdrożenia',   wartosc: 'produkcyjna' },
        { typ: 'para', klucz: 'uwierzytelnienie', wartosc: 'token ważny — urządzenie rozpoznane' },
        { typ: 'para', klucz: 'zapamiętana karta', wartosc: 'brak' },
        { typ: 'ok',   tekst: '  krok rejestracji i logowania pominięty — przejście do etapu 3' }
      ],
      stanFaz: ['gotowa', 'gotowa', 'gotowa', 'gotowa'],
      metaFaz: ['otwarty', 'dnp/1.0', 'token ważny', 'etap 3'],
      plakietka: 'urządzenie rozpoznane',
      klasaPlakietki: 'dn-plakietka dn-plakietka--sukces',
      tetno: false
    },
    'w-blad': {
      hello: [ { typ: 'para', klucz: 'token dostępu', wartosc: 'obecny — token urządzenia' } ],
      ack: [
        { typ: 'blad',  tekst: '× brak odpowiedzi na connection.hello' },
        { typ: 'para',  klucz: 'stan kanału',  wartosc: 'zamknięty' },
        { typ: 'para',  klucz: 'przyczyna',    wartosc: 'serwer nieosiągalny / adres nieprawidłowy / przerwanie sieci' },
        { typ: 'uwaga', tekst: '  oczekiwanie na ponowienie — automatyczne albo kontrolką „Spróbuj ponownie”' }
      ],
      stanFaz: ['blad', 'oczekuje', 'oczekuje', 'oczekuje'],
      metaFaz: ['zamknięty', 'przerwane', 'brak', 'brak'],
      plakietka: 'brak połączenia',
      klasaPlakietki: 'dn-plakietka dn-plakietka--blad',
      tetno: false
    }
  };

  function wierszLogu(poz, ukryty) {
    var w = D.createElement('span');
    w.className = 'pt-log-wiersz';
    if (ukryty) w.setAttribute('data-ukryty', 'tak');
    if (poz.typ === 'para') {
      var k = D.createElement('span');
      k.className = 'pt-log-klucz';
      k.textContent = '  ' + poz.klucz;
      var v = D.createElement('span');
      v.textContent = poz.wartosc;
      w.appendChild(k); w.appendChild(v);
      return w;
    }
    var t = D.createElement('span');
    if (poz.typ === 'zachęta')    t.className = 'pt-znak-zachety';
    else if (poz.typ === 'ok')    t.className = 'pt-wynik-ok';
    else if (poz.typ === 'blad')  t.className = 'pt-wynik-blad';
    else if (poz.typ === 'uwaga') t.className = 'pt-wynik-uwaga';
    else                          t.className = 'pt-komentarz';
    t.textContent = poz.tekst;
    w.appendChild(t);
    return w;
  }

  function zlozPrzebieg(nazwa) {
    var p = PRZEBIEGI[nazwa] || PRZEBIEGI['w-laczenie'];
    return HELLO.concat(p.hello).concat(HELLO_OGON).concat(p.ack);
  }

  function ustawFazy(nazwa) {
    var p = PRZEBIEGI[nazwa] || PRZEBIEGI['w-laczenie'];
    fazy.forEach(function (f, i) {
      f.setAttribute('data-stan-fazy', p.stanFaz[i]);
      var meta = f.querySelector('[data-meta-fazy]');
      if (meta) meta.textContent = p.metaFaz[i];
      var znak = f.querySelector('.pt-faza-znak');
      if (!znak) return;
      if (p.stanFaz[i] === 'gotowa') {
        znak.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 6 9 17l-5-5"/></svg>';
      } else if (p.stanFaz[i] === 'blad') {
        znak.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>';
      } else {
        znak.textContent = String(i + 1);
      }
    });
  }

  function ustawPlakietke(nazwa) {
    var p = PRZEBIEGI[nazwa] || PRZEBIEGI['w-laczenie'];
    plakietka.className = p.klasaPlakietki;
    plakietka.innerHTML = '';
    if (p.tetno) {
      var k = D.createElement('span');
      k.className = 'pt-tetno';
      k.setAttribute('aria-hidden', 'true');
      plakietka.appendChild(k);
    } else {
      var ik = D.createElement('span');
      ik.setAttribute('aria-hidden', 'true');
      ik.innerHTML = p.klasaPlakietki.indexOf('blad') !== -1
        ? '<svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>'
        : '<svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/></svg>';
      plakietka.appendChild(ik.firstChild);
    }
    plakietka.appendChild(D.createTextNode(' ' + p.plakietka));
  }

  function narysujLog(nazwa, ukryte) {
    log.innerHTML = '';
    zlozPrzebieg(nazwa).forEach(function (poz) { log.appendChild(wierszLogu(poz, ukryte)); });
  }

  function ustawPostep(wartosc) {
    pasek.style.width = wartosc + '%';
    if (licznik) licznik.textContent = Math.round(wartosc) + '%';
  }

  var odtwarzanie = null;

  function pokazWariant(nazwa) {
    if (odtwarzanie) return;      /* w trakcie odtwarzania log buduje sekwencja */
    narysujLog(nazwa, false);
    ustawFazy(nazwa);
    ustawPlakietke(nazwa);
    ustawPostep(nazwa === 'w-blad' ? 35 : 100);
  }

  D.addEventListener('click', function (e) {
    var p = e.target.closest('[data-przelacz][data-grupa="wariant"]');
    if (!p) return;
    pokazWariant(p.getAttribute('data-przelacz'));
  });

  function zatrzymajOdtwarzanie() {
    if (!odtwarzanie) return;
    odtwarzanie.forEach(function (id) { clearTimeout(id); });
    odtwarzanie = null;
  }

  function odtworz(nazwaDocelowa) {
    zatrzymajOdtwarzanie();
    odtwarzanie = [];

    var ograniczonyRuch = matchMedia('(prefers-reduced-motion: reduce)').matches;
    var krokCzas = ograniczonyRuch ? 40 : 150;

    narysujLog(nazwaDocelowa, true);
    ustawPostep(0);
    fazy.forEach(function (f, i) {
      f.setAttribute('data-stan-fazy', i === 0 ? 'pracuje' : 'oczekuje');
      var meta = f.querySelector('[data-meta-fazy]');
      if (meta) meta.textContent = i === 0 ? 'w toku' : 'oczekuje';
      var znak = f.querySelector('.pt-faza-znak');
      if (znak) znak.textContent = String(i + 1);
    });
    plakietka.className = 'dn-plakietka dn-plakietka--sygnal';
    plakietka.innerHTML = '<span class="pt-tetno" aria-hidden="true"></span>';
    plakietka.appendChild(D.createTextNode(' uzgadnianie…'));

    if (etykietaOdtworz) etykietaOdtworz.textContent = 'Odtwarzam…';
    btnOdtworz.setAttribute('aria-busy', 'true');

    var wiersze = wszystkie('.pt-log-wiersz', log);
    wiersze.forEach(function (w, i) {
      odtwarzanie.push(setTimeout(function () {
        w.removeAttribute('data-ukryty');
        ustawPostep(Math.min(100, ((i + 1) / wiersze.length) * 100));
        log.scrollTop = log.scrollHeight;
      }, krokCzas * (i + 1)));
    });

    var czasCalosci = krokCzas * (wiersze.length + 1);
    [0, 1, 2, 3].forEach(function (i) {
      odtwarzanie.push(setTimeout(function () {
        ustawFazy(nazwaDocelowa);
        fazy.forEach(function (f, j) {
          if (j > i) {
            f.setAttribute('data-stan-fazy', 'oczekuje');
            var meta = f.querySelector('[data-meta-fazy]');
            if (meta) meta.textContent = 'oczekuje';
            var znak = f.querySelector('.pt-faza-znak');
            if (znak) znak.textContent = String(j + 1);
          }
        });
      }, (czasCalosci / 4) * (i + 1)));
    });

    odtwarzanie.push(setTimeout(function () {
      odtwarzanie = null;
      if (etykietaOdtworz) etykietaOdtworz.textContent = 'Odtwórz przebieg';
      btnOdtworz.removeAttribute('aria-busy');
      idz(nazwaDocelowa, 'wariant');
      pokazWariant(nazwaDocelowa);

      if (!window.dnToast) return;
      if (nazwaDocelowa === 'w-blad') {
        window.dnToast('Błąd połączenia', 'Sprawdź połączenie sieciowe i spróbuj ponownie.', 'ostrzezenie');
      } else if (nazwaDocelowa === 'w-token') {
        window.dnToast('Rozpoznano zaufane urządzenie', 'Uruchomienie bez logowania.', 'sukces');
      } else {
        window.dnToast('Połączono z serwerem', 'Za chwilę otworzy się okno logowania.', 'informacja');
      }
    }, czasCalosci + 200));
  }

  btnOdtworz.addEventListener('click', function () {
    var aktywna = D.querySelector('[data-przelacz][data-grupa="wariant"][aria-selected="true"]');
    odtworz(aktywna ? aktywna.getAttribute('data-przelacz') : 'w-laczenie');
  });

  var btnPonow = D.getElementById('btn-ponow');
  if (btnPonow) {
    btnPonow.addEventListener('click', function () {
      if (window.dnToast) {
        window.dnToast('Ponawianie połączenia', 'Etap 1 uruchomiony od nawiązania kanału WebSocket.', 'informacja');
      }
      idz('w-laczenie', 'wariant');
      odtworz('w-laczenie');
    });
  }

  narysujLog('w-laczenie', false);
  ustawFazy('w-laczenie');
  ustawPlakietke('w-laczenie');
  ustawPostep(100);

  poEtapie('laczenie');
  odswiezPigulki();
})();
