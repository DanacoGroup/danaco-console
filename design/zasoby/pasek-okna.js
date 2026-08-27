/* Pasek narzędzi okna roboczego pokazuje zestaw narzędzi właściwy karcie bieżącej wyłącznie dla kart modułowych i zwija nadmiar pozycji przy zwężeniu okna.

   Odpowiedniki pozycji w menu nadmiaru powstają z pozycji paska — znacznik
   prototypu nie powtarza ich ręcznie. Kolejny prototyp wpina ten plik i podaje
   same narzędzia paska.

   Wymaga: `karty-okna.js` (zdarzenie `dn-karta-zmiana`), `pasek-okna.css`.

   Kontrakt znacznika:
     .dn-pasek-okna[data-modul]                  pasek okna roboczego
     [data-pasek-modul-nazwa]                    nazwa modułu karty bieżącej
     .dn-pasek-okna-grupa > [data-narzedzie]     narzędzie modułu
     .dn-pasek-okna-nadmiar                      pojemnik pozycji zwiniętych
     [data-nadmiar="tak"]                        pozycja przeniesiona pod „…"
   ============================================================================ */
(function () {
  'use strict';

  function q(s, k) { return (k || document).querySelector(s); }
  function qq(s, k) { return Array.prototype.slice.call((k || document).querySelectorAll(s)); }

  var pasek = q('.dn-pasek-okna');
  if (!pasek) { return; }

  var nadmiar = q('.dn-pasek-okna-nadmiar', pasek);
  var menuNadmiaru = nadmiar && q('[data-menu-tresc]', nadmiar);
  var narzedzia = qq('.dn-pasek-okna-grupa > [data-narzedzie]', pasek);
  var rozdzielacze = qq('.dn-pasek-okna-rozdzielacz', pasek);

  function nazwa(b) { return b.getAttribute('data-etykietka') || b.getAttribute('aria-label') || ''; }

  /* ── Odpowiedniki pozycji w menu nadmiaru ─────────────────────────────────
     Jeden wiersz na narzędzie, powstaje raz przy wpięciu pliku. Wiersz nosi
     nazwę i znak narzędzia oraz kieruje kliknięcie do pozycji paska, więc
     zwinięte narzędzie działa tak samo jak rozwinięte. */
  function zbudujOdpowiedniki() {
    if (!menuNadmiaru) { return; }
    narzedzia.forEach(function (b) {
      var wiersz = document.createElement('button');
      wiersz.className = 'sta-menu-poz';
      wiersz.type = 'button';
      wiersz.setAttribute('role', 'menuitem');
      wiersz.setAttribute('data-lustro', b.getAttribute('data-narzedzie'));
      wiersz.hidden = true;
      var znak = b.querySelector('svg');
      if (znak) {
        var kopia = znak.cloneNode(true);
        kopia.setAttribute('class', 'dn-menu-poz-ikona');
        wiersz.appendChild(kopia);
      }
      wiersz.appendChild(document.createTextNode(nazwa(b)));
      wiersz.addEventListener('click', function () { b.click(); });
      menuNadmiaru.appendChild(wiersz);
    });
  }

  /* ── Zwijanie nadmiaru ───────────────────────────────────────────────────
     Pasek nigdy nie łamie się do drugiego wiersza: narzędzia schodzą pod „…"
     od prawej, a rozdzielacz bez sąsiada z prawej strony gaśnie razem z nimi. */
  /* ── Zwijanie nadmiaru ───────────────────────────────────────────────────
     Pasek nigdy nie łamie się do drugiego wiersza. Zwijanie prowadzi pomiar:
     dopóki przycisk „…" nie mieści się w obrysie paska, kolejne narzędzie od
     prawej schodzi pod niego. Suma szerokości nie wystarcza — odstępy i
     rozdzielacze zmieniają się razem ze składem paska. */
  function nieMiesciSie() {
    return nadmiar.getBoundingClientRect().right > pasek.getBoundingClientRect().right - 1;
  }

  function porzadkujRozdzielacze() {
    rozdzielacze.forEach(function (r) {
      var poPrawej = false;
      for (var n = r.nextElementSibling; n; n = n.nextElementSibling) {
        if (n.classList.contains('dn-pasek-okna-nadmiar') || n.classList.contains('dn-pasek-okna-odstep')) { continue; }
        if (n.getAttribute('data-nadmiar') !== 'tak' && !n.hidden && n.getBoundingClientRect().width > 0) { poPrawej = true; break; }
      }
      if (poPrawej) { r.removeAttribute('data-nadmiar'); } else { r.setAttribute('data-nadmiar', 'tak'); }
    });
  }

  function przelicz() {
    if (!nadmiar || pasek.hidden) { return; }
    narzedzia.forEach(function (p) { p.removeAttribute('data-nadmiar'); });
    rozdzielacze.forEach(function (p) { p.removeAttribute('data-nadmiar'); });
    nadmiar.hidden = false;
    porzadkujRozdzielacze();
    for (var i = narzedzia.length - 1; i >= 0 && nieMiesciSie(); i -= 1) {
      narzedzia[i].setAttribute('data-nadmiar', 'tak');
      porzadkujRozdzielacze();
    }
    var zwiniete = narzedzia.filter(function (p) { return p.getAttribute('data-nadmiar') === 'tak'; });
    nadmiar.hidden = zwiniete.length === 0;
    if (menuNadmiaru) {
      qq('[data-lustro]', menuNadmiaru).forEach(function (l) {
        var zrodlo = q('[data-narzedzie="' + l.getAttribute('data-lustro') + '"]', pasek);
        l.hidden = !zrodlo || zrodlo.getAttribute('data-nadmiar') !== 'tak';
      });
    }
  }

  /* ── Obecność paska ──────────────────────────────────────────────────────
     Karta Centrum dowodzenia jest pulpitem okna, nie modułem — paska nie ma
     wcale. Karta modułowa niesie nazwę modułu i bierze jego zestaw narzędzi. */
  function ustawWidocznosc(rodzaj, modul) {
    pasek.hidden = rodzaj === 'centrum';
    if (!pasek.hidden && modul) {
      pasek.setAttribute('data-modul', modul);
      var nazwaModulu = q('[data-pasek-modul-nazwa]', pasek);
      if (nazwaModulu) { nazwaModulu.textContent = modul; }
    }
    przelicz();
  }

  qq('[data-narzedzie]', pasek).forEach(function (b) {
    b.addEventListener('click', function () {
      if (!window.dnToast) { return; }
      window.dnToast(nazwa(b), 'Narzędzie modułu ' + (pasek.getAttribute('data-modul') || '')
        + ' — działa na treści karty bieżącej.', 'informacja');
    });
  });

  document.addEventListener('dn-karta-zmiana', function (e) {
    ustawWidocznosc(e.detail.rodzaj, e.detail.modul);
  });

  window.addEventListener('resize', przelicz);

  zbudujOdpowiedniki();
  przelicz();
})();
