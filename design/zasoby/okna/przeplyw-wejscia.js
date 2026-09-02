/* Skrypt prowadzi mechanikę okna przepływu wejścia: ekran startowy, nawiązanie połączenia bez przycisku potwierdzenia, uwierzytelnienie i przejście do przygotowania środowiska.
   Cel przejścia po połączeniu niesie znacznik panelu (`data-po-polaczeniu`).
   Zakładanie konta sprawdza przed wysłaniem budowę adresu, siłę hasła i zgodność
   haseł; login zajęty rozstrzyga rdzeń kodem `conflict`, progu prób tu nie ma.
   Usterki zbierają się w JEDEN komunikat nad formularzem, a pole, którego
   dotyczą, niesie `aria-invalid`.
   Odsłony osiągalne adresem: ?rejestracja=login-zajety | email-bledny |
   haslo-slabe | hasla-rozne
   ============================================================================ */
(function () {
'use strict';

/* Okna powstają z montażu, więc przy wczytaniu tego pliku nie ma jeszcze czego
   obsługiwać. Mechanika rusza dopiero na zdarzenie `wejscie-gotowe`.

   Gotowość rozstrzyga PUSTE miejsce montażu, nie obecność panelu: w podglądzie
   stoi jeszcze okno pisane ręcznie, więc panel istnieje od pierwszej chwili
   i mylił ten warunek. */
var miejsceMontazu = document.querySelector('[data-wejscie-okno]');
if (miejsceMontazu && !miejsceMontazu.firstElementChild) {
  document.addEventListener('wejscie-gotowe', uruchom, { once: true });
} else {
  uruchom();
}

function uruchom() {

var scena = document.querySelector('[data-uruchomienie-scena]');
if (!scena) return;
var warstwa = scena.querySelector('[data-uruchomienie]');
var pole = warstwa && warstwa.querySelector('[data-ekran-startowy]');
if (!warstwa || !pole) return;

/* Ile trwa jeden etap łączenia. Cztery etapy mieszczą się w czasie, w którym
   użytkownik jeszcze czeka, a nie zaczyna się zastanawiać, czy coś się zacięło. */
var ETAP_MS = 900;
/* Chwila po ostatnim etapie, żeby dało się zobaczyć, że wszystko doszło do
   końca, zanim okno ustąpi następnemu. */
var DOMKNIECIE_MS = 700;

var biegi = [];
function zatrzymaj() { biegi.forEach(clearTimeout); biegi = []; }

var ZNAK_GOTOWY = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" ' +
  'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 6 9 17l-5-5"/></svg>';

/* Znak etapu idzie za stanem: ptaszek dla zrobionego, tętno dla trwającego,
   numer dla czekającego. Znak zostawiony z poprzedniego stanu kłamie —
   ptaszek przy „w toku” mówi, że rzecz jest skończona. */
var ODMIANA = { gotowy: 'dn-krok--poprawny', pracuje: 'dn-krok--pracuje', blad: 'dn-krok--wstrzymany' };

function oznacz(krok, stan, meta, numer) {
  krok.setAttribute('data-stan', stan);
  Object.keys(ODMIANA).forEach(function (k) { krok.classList.toggle(ODMIANA[k], k === stan); });
  var m = krok.querySelector('.dn-krok-meta');
  if (m) m.textContent = meta;
  var z = krok.querySelector('.dn-krok-znak');
  if (!z) return;
  if (stan === 'gotowy') z.innerHTML = ZNAK_GOTOWY;
  else if (stan === 'pracuje') z.innerHTML = '<span class="pt-tetno" aria-hidden="true"></span>';
  else z.textContent = String(numer);
}

/* Etapy łączenia zapalają się po kolei, ostatni domyka przejście do etapu
   wskazanego przez panel. */
function polacz() {
  var panel = scena.querySelector('[data-po-polaczeniu][data-widok-aktywny="tak"]');
  if (!panel) return;
  var kroki = [].slice.call(panel.querySelectorAll('.dn-krok'));
  if (!kroki.length) return;

  kroki.forEach(function (k, i) {
    biegi.push(setTimeout(function () {
      kroki.forEach(function (w, j) {
        if (j < i) oznacz(w, 'gotowy', 'gotowe', j + 1);
        else if (j === i) oznacz(w, 'pracuje', 'w toku', j + 1);
        else oznacz(w, 'oczekuje', '—', j + 1);
      });
    }, ETAP_MS * i));
  });
  biegi.push(setTimeout(function () {
    kroki.forEach(function (k, j) { oznacz(k, 'gotowy', 'gotowe', j + 1); });
  }, ETAP_MS * kroki.length));

  var cel = panel.dataset.poPolaczeniu;
  if (cel && window.dnPrzejdz) {
    biegi.push(setTimeout(function () {
      window.dnPrzejdz(cel, 'etap');
    }, ETAP_MS * kroki.length + DOMKNIECIE_MS));
  }
}

/* Okno czeka za animacją. Stan niesie scena, nie okno — okno nie wie, że coś
   grało przed nim, i nie musi wiedzieć. */
function zaslon() {
  zatrzymaj();
  scena.dataset.uruchamianie = 'tak';
  warstwa.hidden = false;
  if (pole.ekranStartowy) pole.ekranStartowy.odtworz();
}

function odslon() {
  delete scena.dataset.uruchamianie;
  warstwa.hidden = true;
  polacz();
}

scena.addEventListener('ekran-startowy-koniec', odslon);

var odtworz = document.getElementById('btn-odtworz');
if (odtworz) odtworz.addEventListener('click', zaslon);

/* Gdy użytkownik woli ruch wyłączony, okno staje od razu — animacja startowa
   nie niesie treści, więc jej pominięcie niczego nie zabiera. Przebieg
   łączenia biegnie dalej, bo on treść niesie. */
if (window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
  odslon();
} else {
  zaslon();
}

/* ── Formularz uwierzytelnienia ──────────────────────────────────────────── */

function wszystkie(sel) { return Array.prototype.slice.call(document.querySelectorAll(sel)); }

  /* Odsłonięcie hasła — kontrolka zmienia typ pola i własną etykietę. */
  wszystkie('[data-odsloniecie]').forEach(function (b) {
    b.addEventListener('click', function () {
      var pole = document.getElementById(b.getAttribute('data-odsloniecie'));
      if (!pole) return;
      var odkryte = pole.type === 'text';
      pole.type = odkryte ? 'password' : 'text';
      b.setAttribute('aria-pressed', odkryte ? 'false' : 'true');
      b.setAttribute('aria-label', odkryte ? 'Pokaż hasło' : 'Ukryj hasło');
    });
  });

/* Kod potwierdzający — przejście między polami i wklejenie całości. Wiązanie
   idzie GRUPAMI: każdy zestaw `.au-kod` rządzi się sam. Jeden wspólny wykaz
   przeskakiwałby kursorem między oknami, a pola drugiego zestawu stały
   niepodłączone, bo nie miały znacznika. */
wszystkie('.au-kod').forEach(function (grupa) {
  var polaKodu = Array.prototype.slice.call(grupa.querySelectorAll('.au-kod-pole'));
  if (!polaKodu.length) return;

  /* Kod jest sześciocyfrowy — tak nazywa go okno i tak zapowiada pole
     (`inputmode="numeric"`). Przyjmowanie liter i zamiana ich na wersaliki
     przeczyła obu tym zapowiedziom. */
  function rozsyp(tekst) {
    var czysty = String(tekst || '').replace(/[^0-9]/g, '');
    polaKodu.forEach(function (q, j) { q.value = czysty[j] || ''; });
    polaKodu[Math.min(czysty.length, polaKodu.length - 1)].focus();
  }

  polaKodu.forEach(function (p, i) {
    p.addEventListener('input', function () {
      p.value = p.value.replace(/[^0-9]/g, '').slice(0, 1);
      if (p.value && polaKodu[i + 1]) polaKodu[i + 1].focus();
    });
    p.addEventListener('keydown', function (e) {
      if (e.key === 'Backspace' && !p.value && polaKodu[i - 1]) polaKodu[i - 1].focus();
    });
    p.addEventListener('paste', function (e) {
      e.preventDefault();
      rozsyp((e.clipboardData || window.clipboardData).getData('text'));
    });
  });

  var stopka = grupa.parentNode.querySelector('[data-wklej-kod]');
  if (!stopka) return;
  stopka.addEventListener('click', function () {
    if (!navigator.clipboard || !navigator.clipboard.readText) {
      if (window.dnToast) window.dnToast('Schowek niedostępny', 'Wpisz kod ręcznie w sześciu polach.', 'informacja');
      return;
    }
    navigator.clipboard.readText().then(rozsyp).catch(function () {
      if (window.dnToast) window.dnToast('Schowek niedostępny', 'Wpisz kod ręcznie w sześciu polach.', 'informacja');
    });
  });
});

/* Miara postępu z danej na znaczniku. Wartość jest daną, nie wyglądem, więc
   stoi w `data-wartosc`, a szerokość ustawia mechanika. */
var miary = document.querySelectorAll('.dn-postep-wartosc[data-wartosc]');
for (var m = 0; m < miary.length; m++) {
  miary[m].style.width = miary[m].dataset.wartosc + '%';
}

/* ── Odliczanie czasu ────────────────────────────────────────────────────── */

/* Każdy węzeł z `data-odliczanie` niesie czas w sekundach i sam się wypisuje.
   Napis „09:12", który nie ubywa, jest gorszy od braku napisu: obiecuje odmierzanie,
   którego nie ma. Jeden zegar na całe okno — nie dwanaście osobnych. */
function naZegar(sekundy) {
  var m = Math.floor(sekundy / 60), s = sekundy % 60;
  return (m < 10 ? '0' : '') + m + ':' + (s < 10 ? '0' : '') + s;
}

function odliczaj() {
  var pola = document.querySelectorAll('[data-odliczanie]');
  for (var i = 0; i < pola.length; i++) {
    var e = pola[i];
    var zostalo = parseInt(e.dataset.odliczanie, 10);
    if (isNaN(zostalo)) continue;
    /* Wypisywane są WSZYSTKIE liczniki, także w odsłonach ukrytych — inaczej
       licznik byłby pusty przez sekundę po pokazaniu odsłony. Ubywa natomiast
       tylko licznik widoczny: czas, który schodzi za plecami, doprowadza do
       tego, że użytkownik zastaje zero, choć odsłonę zobaczył przed chwilą. */
    var napis = e.dataset.odliczaniePostac === 'sekundy' ? String(zostalo) : naZegar(zostalo);
    if (e.textContent !== napis) e.textContent = napis;
    if (!e.offsetParent) continue;
    if (zostalo <= 0) { poZerze(e); continue; }
    e.dataset.odliczanie = String(zostalo - 1);
  }
}

/* Co się dzieje po dojściu do zera, rozstrzyga odsłona, w której licznik stoi.
   Ważność kodu wygasa — to stan, którego kompozycja nie ma, więc licznik staje
   i nic nie udaje. */
function poZerze(e) {
  var panel = e.closest('.we-panel');
  if (!panel) return;
  var widok = panel.dataset.widok;
  if (widok === 'w-blad') {
    /* Serwer nie odpowiadał; po odliczonym czasie okno ponawia próbę samo —
       tak, jak zapowiada baner. Licznik rusza od nowa razem z przebiegiem. */
    e.dataset.odliczanie = '15';
    window.dnPrzelaczWidok('w-laczenie', 'wariant');
    polacz();
    return;
  }
}

setInterval(odliczaj, 1000);
odliczaj();

/* ── Zakładanie konta: sprawdzenie przed wysłaniem ───────────────────────── */

/* Formularze, które ustawiają hasło. Oba sprawdzają to samo — różnią się
   wyłącznie tym, co jeszcze mają do sprawdzenia poza samym hasłem. */
var FORMULARZE = [
  { widok: 'logowanie', naglowek: 'usterki.naglowekLogowanie', zbiorczyBrak: 'brakDanych',
    wymagane: [{ id: 'log-login', usterka: 'brakLoginu' }, { id: 'log-haslo', usterka: 'brakHasla' }] },
  { widok: 'logowanie-blad', naglowek: 'usterki.naglowekLogowanie', zbiorczyBrak: 'brakDanych',
    wymagane: [{ id: 'blad-login', usterka: 'brakLoginu' }, { id: 'blad-haslo', usterka: 'brakHasla' }] },
  { widok: 'odzyskiwanie-adres', naglowek: 'usterki.naglowekKod',
    wymagane: [{ id: 'odz-email', usterka: 'brakAdresu' }], email: 'odz-email' },
  { widok: 'rejestracja', naglowek: 'usterki.naglowekKonto',
    login: 'rej-login', email: 'rej-email', haslo: 'rej-haslo', haslo2: 'rej-haslo-2' },
  { widok: 'odzyskiwanie-haslo', naglowek: 'usterki.naglowekHaslo',
    haslo: 'odz-haslo', haslo2: 'odz-haslo-2' }
];

/* Każda usterka nazywa rzecz i mówi, co z nią zrobić. Kolejność jest
   kolejnością pól w formularzu — lista czyta się z góry na dół tak samo, jak
   wzrok wraca do pól. */
/* Rozpoznania usterek. Same POLA stoją tutaj — to kontrakt z formularzem,
   nie treść. Napisy bierze się z katalogu, bo są tekstem dla użytkownika. */
var USTERKI = {
  'brakLoginu':   { pola: ['log-login', 'blad-login'] },
  'brakHasla':    { pola: ['log-haslo', 'blad-haslo'] },
  'brakAdresu':   { pola: ['odz-email'] },
  'brakDanych':   { pola: ['log-login', 'log-haslo', 'blad-login', 'blad-haslo'] },
  'login-zajety': { klucz: 'loginZajety', pola: ['rej-login'] },
  'email-bledny': { klucz: 'emailBledny', pola: ['rej-email', 'odz-email'] },
  'haslo-slabe':  { klucz: 'hasloSlabe',  pola: ['rej-haslo', 'odz-haslo'] },
  'hasla-rozne':  { klucz: 'haslaRozne',  pola: ['rej-haslo', 'rej-haslo-2', 'odz-haslo', 'odz-haslo-2'] }
};

/* Katalog treści okna wiąże napisy interfejsu z ich źródłem zewnętrznym: ani jeden napis wyświetlany operatorowi nie stoi zapisany wprost w tym pliku. */
var N = window.DanacoNarzedzia.zwiaz((window.DanacoWejscie || {}).tresci || {});
function tekst(sciezka) { return N.tekst(sciezka); }
function usterkaTekst(k, co) { return tekst('usterki.' + (USTERKI[k].klucz || k) + '.' + co); }

/* Zwraca pola, których dotyczy usterka logowania, ograniczone do tych spośród nich, które rzeczywiście występują w danym formularzu okna. */
function polaUsterki(f, k) {
  return USTERKI[k].pola.filter(function (id) { return pola(f).indexOf(id) !== -1; });
}

var KOLEJNOSC = ['brakDanych', 'brakLoginu', 'brakHasla', 'brakAdresu', 'login-zajety', 'email-bledny', 'haslo-slabe', 'hasla-rozne'];

var ZNAK_USTERKI = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" ' +
  'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' +
  '<path d="M10.3 3.6 2.5 17a2 2 0 0 0 1.7 3h15.6a2 2 0 0 0 1.7-3L13.7 3.6a2 2 0 0 0-3.4 0"/>' +
  '<path d="M12 9v4"/><path d="M12 17h.01"/></svg>';

function pola(f) {
  var lista = [f.login, f.email, f.haslo, f.haslo2].filter(Boolean);
  (f.wymagane || []).forEach(function (w) { if (lista.indexOf(w.id) === -1) lista.push(w.id); });
  return lista;
}

/* Zdejmowany jest wyłącznie komunikat postawiony przez sprawdzenie. Odsłona
   bywa niesie własny baner — jak „Nie rozpoznano danych logowania" — i ten
   należy do niej, nie do nas. */
function wyczysc(f) {
  var panel = document.querySelector('.we-panel[data-widok="' + f.widok + '"]');
  var kom = panel && panel.querySelector('.we-komunikaty');
  if (kom) {
    var moje = kom.querySelectorAll('[data-usterka-formularza]');
    for (var i = 0; i < moje.length; i++) moje[i].parentNode.removeChild(moje[i]);
  }
  pola(f).forEach(function (id) {
    var e = document.getElementById(id);
    if (e) e.removeAttribute('aria-invalid');
  });
}

function pokazUsterki(f, klucze) {
  wyczysc(f);
  var panel = document.querySelector('.we-panel[data-widok="' + f.widok + '"]');
  var komunikaty = panel && panel.querySelector('.we-komunikaty');
  if (!komunikaty || !klucze.length) return;
  klucze = KOLEJNOSC.filter(function (k) { return klucze.indexOf(k) !== -1; });

  var tresc;
  if (klucze.length === 1) {
    tresc = '<b>' + usterkaTekst(klucze[0], 'glowa') + '</b>' + usterkaTekst(klucze[0], 'tresc');
  } else {
    /* Przy kilku usterkach wykaz niesie same rozpoznania, bez wskazówek: co
       zrobić, mówi zaznaczone pole, a wskazówki powtórzone cztery razy
       wypchnęłyby formularz poza okno. Bez liczebnika, bo „dwie / trzy / cztery
       rzeczy” to trzy odmiany do utrzymania i trzy okazje do pomyłki. */
    tresc = '<b>' + tekst(f.naglowek) + ' ' + tekst('usterki.wiele') + '</b>' +
            '<ul class="dn-alert-lista">' +
            klucze.map(function (k) { return '<li>' + usterkaTekst(k, 'glowa') + '</li>'; }).join('') +
            '</ul>';
  }

  var el = document.createElement('div');
  el.className = 'dn-alert dn-alert--wstega dn-alert--blad';
  el.setAttribute('role', 'alert');
  el.setAttribute('data-usterka-formularza', '');
  el.innerHTML = '<span class="dn-alert-znak" aria-hidden="true">' + ZNAK_USTERKI +
                 '</span><span class="dn-alert-tresc">' + tresc + '</span>';
  komunikaty.appendChild(el);

  klucze.forEach(function (k) {
    polaUsterki(f, k).forEach(function (id) {
      var pole = document.getElementById(id);
      if (pole) pole.setAttribute('aria-invalid', 'true');
    });
  });
  var pierwsze = document.getElementById(polaUsterki(f, klucze[0])[0]);
  if (pierwsze) pierwsze.focus();
}

/* Budowa adresu, nie jego istnienie. Że adres istnieje, rozstrzyga dopiero kod
   wysłany na niego w kroku potwierdzenia — tutaj wyłapuje się to, co widać bez
   wysyłania: brak znaku @, brak nazwy, brak domeny, spacja w środku. */
function adresPoprawny(v) {
  return /^[^\s@]+@[^\s@.]+(\.[^\s@.]+)+$/.test(v.trim());
}

/* ── Siła hasła ──────────────────────────────────────────────────────────── */

/* Cztery warunki. Jedna reguła dla miernika i dla sprawdzenia przed wysłaniem —
   dwie osobne rozjechałyby się przy pierwszej zmianie wymagań. */
function warunkiHasla(v) {
  return {
    dlugosc: v.length >= 12,
    wielkosc: /[a-ząćęłńóśźż]/.test(v) && /[A-ZĄĆĘŁŃÓŚŹŻ]/.test(v),
    cyfra: /[0-9]/.test(v),
    znak: /[^0-9A-Za-zĄĆĘŁŃÓŚŹŻąćęłńóśźż]/.test(v)
  };
}
function hasloSpelnia(v) {
  var w = warunkiHasla(v);
  return w.dlugosc && w.wielkosc && w.cyfra && w.znak;
}

var SILA_PUSTE = 'Siła hasła zostanie oceniona podczas wpisywania';
var SILA_OPISY = ['Hasło nie spełnia żadnego z warunków',
                  'Hasło słabe — spełnia jeden warunek',
                  'Hasło dostateczne — spełnia dwa warunki',
                  'Hasło dobre — spełnia trzy warunki',
                  'Hasło mocne — spełnia wszystkie cztery warunki'];

/* Miernik wiąże się z polem przez `data-sila-dla`, a nie przez sąsiedztwo
   w drzewie. Sąsiedztwo bywa różne w różnych oknach i cicho się rozjeżdża —
   drugi miernik w tym prototypie stał martwy właśnie dlatego. */
function zalozMierniki() {
  var mierniki = document.querySelectorAll('[data-sila-dla]');
  for (var i = 0; i < mierniki.length; i++) {
    (function (blok) {
      var pole = document.getElementById(blok.dataset.silaDla);
      if (!pole) return;
      var opis = blok.querySelector('[data-sila-opis]');
      function odswiez() {
        var v = pole.value;
        var w = warunkiHasla(v);
        var stopien = 0;
        Object.keys(w).forEach(function (k) {
          var el = blok.querySelector('[data-warunek="' + k + '"]');
          if (el) el.setAttribute('data-spelniony', w[k] ? 'tak' : 'nie');
          if (w[k]) stopien += 1;
        });
        blok.setAttribute('data-stopien', String(stopien));
        if (opis) opis.textContent = v ? SILA_OPISY[stopien] : SILA_PUSTE;
      }
      pole.addEventListener('input', odswiez);
      odswiez();
    })(mierniki[i]);
  }
}
zalozMierniki();

function sprawdz(f) {
  var w = function (id) { return id ? document.getElementById(id) : null; };
  var login = w(f.login), email = w(f.email), haslo = w(f.haslo), haslo2 = w(f.haslo2);
  var braki = [];
  /* Pole puste rozstrzyga się przed wszystkim innym: „nieprawidłowy adres"
     przy pustym polu mówi nieprawdę — nic nie zostało wpisane. */
  var puste = (f.wymagane || []).filter(function (r) {
    var pole = w(r.id);
    return pole && !pole.value.trim();
  });
  /* Wszystkie pola puste to JEDNA sprawa — formularz nie został wypełniony.
     Rozbicie jej na osobne rozpoznania mnożyło wiersze komunikatu i wypychało
     treść poza okno, nie mówiąc przy tym nic więcej. */
  if (f.zbiorczyBrak && puste.length === (f.wymagane || []).length && puste.length > 1) {
    return [f.zbiorczyBrak];
  }
  puste.forEach(function (r) { braki.push(r.usterka); });
  if (braki.length) return braki;
  if (email && !adresPoprawny(email.value)) braki.push('email-bledny');
  if (haslo && !hasloSpelnia(haslo.value)) braki.push('haslo-slabe');
  if (haslo && haslo2 && haslo.value !== haslo2.value) braki.push('hasla-rozne');
  return braki;
}

FORMULARZE.forEach(function (f) {
  var panel = document.querySelector('.we-panel[data-widok="' + f.widok + '"]');
  if (!panel) return;

  var pas = document.querySelector('.we-pas[data-widok="' + f.widok + '"]');
  var glowna = pas && pas.querySelector('.dn-btn--sygnal');
  if (glowna) {
    /* Przechwycenie w fazie łapania: przepływ podglądu słucha kliknięć na
       dokumencie, więc zatrzymanie musi nastąpić, zanim tam dojdzie. */
    glowna.addEventListener('click', function (e) {
      var braki = sprawdz(f);
      if (braki.length) {
        e.preventDefault();
        e.stopPropagation();
        pokazUsterki(f, braki);
        return;
      }
    }, true);
  }

  pola(f).forEach(function (id) {
    var pole = document.getElementById(id);
    if (pole) pole.addEventListener('input', function () {
      if (pole.getAttribute('aria-invalid') === 'true') wyczysc(f);
    });
  });
});

/* Odsłona wskazana adresem — żeby dało się obejrzeć komunikat bez wpisywania. */
var zAdresu = new URLSearchParams(location.search).get('rejestracja');
if (zAdresu && USTERKI[zAdresu]) {
  var f = FORMULARZE.filter(function (x) { return x.widok === 'rejestracja'; })[0];
  if (window.dnPrzelaczWidok) {
    window.dnPrzejdz && window.dnPrzejdz('uwierzytelnienie', 'etap');
    window.dnPrzelaczWidok('rejestracja', 'stan');
  }
  pokazUsterki(f, [zAdresu]);
}
}
})();
