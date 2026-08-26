/* ============================================================================
   PRZEPŁYW WEJŚCIA — mechanika okna

   Uruchomienie programu: po kliknięciu ikony na monitorze gra ekran startowy,
   a okno stoi za nim niewidoczne. Dopiero koniec animacji je odsłania — tak
   samo, jak powłoka w produkcie pokazuje okno dopiero po zdarzeniu
   `ekran-startowy-koniec`.

   Nawiązanie połączenia nie kończy się przyciskiem. Etapy zaznaczają się same,
   a po ostatnim okno przechodzi dalej — nie ma tu czego wybierać, więc nie ma
   czego klikać. Jedyną czynnością pozostaje zamknięcie programu.

   Cel przejścia niesie znacznik panelu (`data-po-polaczeniu`), nie ten skrypt.

   Zakładanie konta sprawdza się przed wysłaniem, nie po. Cztery rzeczy mogą
   pójść źle i każda mówi, co zrobić: login zajęty, adres e-mail bez właściwej
   budowy, hasło poniżej wymagań, hasła różne.

   Usterki zbierają się w JEDEN komunikat. Cztery osobne banery zepchnęłyby
   formularz poza okno, a użytkownik i tak czyta je jako jedną listę tego, co
   ma poprawić. Komunikat staje nad formularzem, a pole, którego dotyczy, niesie
   `aria-invalid` — czytnik ekranu dowiaduje się tego samego co oko.

   Odsłony osiągalne adresem, do których przepływ sam nie doprowadzi:
       ?rejestracja=login-zajety | email-bledny | haslo-slabe | hasla-rozne

   Podgląd odtwarza całość na żądanie: „Odtwórz przebieg” zaczyna od animacji,
   a nie od gotowego okna.
   ============================================================================ */
(function () {
'use strict';

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
function oznacz(krok, stan, meta, numer) {
  krok.setAttribute('data-stan', stan);
  var m = krok.querySelector('.we-krok-meta');
  if (m) m.textContent = meta;
  var z = krok.querySelector('.we-krok-znak');
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
  var kroki = [].slice.call(panel.querySelectorAll('.we-krok'));
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

/* ── Zakładanie konta: sprawdzenie przed wysłaniem ───────────────────────── */

/* Formularze, które ustawiają hasło. Oba sprawdzają to samo — różnią się
   wyłącznie tym, co jeszcze mają do sprawdzenia poza samym hasłem. */
var FORMULARZE = [
  { widok: 'rejestracja',    login: 'rej-login', email: 'rej-email', haslo: 'rej-haslo', haslo2: 'rej-haslo-2' },
  { widok: 'odzyskiwanie-haslo', haslo: 'odz-haslo', haslo2: 'odz-haslo-2' }
];

/* Loginy zajęte. W produkcie odpowiada na to serwer; w podglądzie musi stać
   nazwa, na której da się to zobaczyć — bez niej odsłona jest nieosiągalna. */
var ZAJETE = ['operator', 'admin', 'danaco', 'konsola'];

/* Każda usterka nazywa rzecz i mówi, co z nią zrobić. Kolejność jest
   kolejnością pól w formularzu — lista czyta się z góry na dół tak samo, jak
   wzrok wraca do pól. */
var USTERKI = {
  'login-zajety': {
    glowa: 'Ten login jest już zajęty.',
    tresc: 'Wybierz inny login. Adres e-mail może pozostać bez zmian.',
    pola: ['rej-login']
  },
  'email-bledny': {
    glowa: 'Nieprawidłowy adres e-mail.',
    tresc: 'Sprawdź, czy adres zawiera znak @ oraz nazwę domeny, na przykład nazwa@firma.pl.',
    pola: ['rej-email']
  },
  'haslo-slabe': {
    glowa: 'Hasło nie spełnia wymagań.',
    tresc: 'Spełnij wszystkie cztery warunki podane pod polem hasła.',
    pola: ['rej-haslo']
  },
  'hasla-rozne': {
    glowa: 'Hasła nie są zgodne.',
    tresc: 'Wpisz to samo hasło w obu polach.',
    pola: ['rej-haslo', 'rej-haslo-2']
  }
};
var KOLEJNOSC = ['login-zajety', 'email-bledny', 'haslo-slabe', 'hasla-rozne'];

var ZNAK_USTERKI = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" ' +
  'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' +
  '<path d="M10.3 3.6 2.5 17a2 2 0 0 0 1.7 3h15.6a2 2 0 0 0 1.7-3L13.7 3.6a2 2 0 0 0-3.4 0"/>' +
  '<path d="M12 9v4"/><path d="M12 17h.01"/></svg>';

function pola(f) {
  return [f.login, f.email, f.haslo, f.haslo2].filter(Boolean);
}

function wyczysc(f) {
  var panel = document.querySelector('.we-panel[data-widok="' + f.widok + '"]');
  var kom = panel && panel.querySelector('.we-komunikaty');
  if (kom) kom.innerHTML = '';
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
    var u = USTERKI[klucze[0]];
    tresc = '<b>' + u.glowa + '</b>' + u.tresc;
  } else {
    /* Przy kilku usterkach wykaz niesie same rozpoznania, bez wskazówek: co
       zrobić, mówi zaznaczone pole, a wskazówki powtórzone cztery razy
       wypchnęłyby formularz poza okno. Bez liczebnika, bo „dwie / trzy / cztery
       rzeczy” to trzy odmiany do utrzymania i trzy okazje do pomyłki. */
    tresc = '<b>' + f.naglowek + ' Popraw zaznaczone dane.</b>' +
            '<ul class="we-alarm-lista">' +
            klucze.map(function (k) { return '<li>' + USTERKI[k].glowa + '</li>'; }).join('') +
            '</ul>';
  }

  var el = document.createElement('div');
  el.className = 'we-alarm we-alarm--blad';
  el.setAttribute('role', 'alert');
  el.innerHTML = ZNAK_USTERKI + '<span>' + tresc + '</span>';
  komunikaty.appendChild(el);

  klucze.forEach(function (k) {
    USTERKI[k].pola.forEach(function (id) {
      var pole = document.getElementById(id);
      if (pole) pole.setAttribute('aria-invalid', 'true');
    });
  });
  var pierwsze = document.getElementById(USTERKI[klucze[0]].pola[0]);
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
  if (login && ZAJETE.indexOf(login.value.trim().toLowerCase()) !== -1) braki.push('login-zajety');
  if (email && !adresPoprawny(email.value)) braki.push('email-bledny');
  if (haslo && !hasloSpelnia(haslo.value)) braki.push('haslo-slabe');
  if (haslo && haslo2 && haslo.value !== haslo2.value) braki.push('hasla-rozne');
  return braki;
}

FORMULARZE.forEach(function (f) {
  var panel = document.querySelector('.we-panel[data-widok="' + f.widok + '"]');
  if (!panel) return;
  f.naglowek = f.widok === 'rejestracja' ? 'Nie można utworzyć konta.' : 'Nie można ustawić hasła.';

  var pas = document.querySelector('.we-pas[data-widok="' + f.widok + '"]');
  var glowna = pas && pas.querySelector('.dn-btn--sygnal');
  if (glowna) {
    /* Przechwycenie w fazie łapania: przepływ podglądu słucha kliknięć na
       dokumencie, więc zatrzymanie musi nastąpić, zanim tam dojdzie. */
    glowna.addEventListener('click', function (e) {
      var braki = sprawdz(f);
      if (!braki.length) { wyczysc(f); return; }
      e.preventDefault();
      e.stopPropagation();
      pokazUsterki(f, braki);
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
  var f = FORMULARZE[0];
  f.naglowek = 'Nie można utworzyć konta.';
  if (window.dnPrzelaczWidok) {
    window.dnPrzejdz && window.dnPrzejdz('uwierzytelnienie', 'etap');
    window.dnPrzelaczWidok('rejestracja', 'stan');
  }
  pokazUsterki(f, [zAdresu]);
}
})();
