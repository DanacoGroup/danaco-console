/* Okno robocze obsługuje maksymalizację przez przycisk pasma, skrót klawiszowy i klawisz Escape, z pamięcią stanu dla każdego okna oraz przełącznikami izolacji.

   Wymaga: `prototyp.js` (`dnToast`, `dnOglos`), `okno-robocze.css`,
   `izolacja.css`.

   Kontrakt znacznika:
     .sta-powloka[data-okno][data-okno-max]      powłoka okna roboczego
     [data-okno-max-przelacz]                    przycisk pasma
     [data-izolacja-przelacz="incognito|galaz"]  przełączniki zakładania okna
     [data-okno-nowe][data-tworzy-sesje]        wyrzutnia: pozycja szyny, kafel
   Skrót: Ctrl+Shift+Enter przełącza maksymalizację, Esc z niej wychodzi.
   ============================================================================ */
(function () {
  'use strict';

  function q(s, k) { return (k || document).querySelector(s); }
  function qq(s, k) { return Array.prototype.slice.call((k || document).querySelectorAll(s)); }
  function toast(t, o, r, ms) { if (window.dnToast) { window.dnToast(t, o, r || 'informacja', ms || 3200); } }
  function oglos(t) { if (window.dnOglos) { window.dnOglos(t); } }

  var powloka = q('.sta-powloka');
  if (!powloka) { return; }

  var pamiec = {};

  function ustawMaks(wl) {
    var okno = powloka.getAttribute('data-okno') || 'okno-1';
    powloka.setAttribute('data-okno-max', wl ? 'tak' : 'nie');
    pamiec[okno] = wl;
    qq('[data-okno-max-przelacz]').forEach(function (b) {
      b.setAttribute('aria-pressed', wl ? 'true' : 'false');
    });
    oglos(wl ? 'Okno robocze zmaksymalizowane.' : 'Okno robocze przywrócone.');
  }

  qq('[data-okno-max-przelacz]').forEach(function (b) {
    b.addEventListener('click', function () {
      ustawMaks(powloka.getAttribute('data-okno-max') !== 'tak');
    });
  });

  /* ── Zakładanie okna roboczego ─────────────────────────────────────────────
     Szyna nawigacji jest wyrzutnią: pozycja otwiera nowe okno robocze na karcie
     Centrum dowodzenia. Okno powstaje izolowane — własna gałąź plików jest
     stanem domyślnym, incognito włącza Operator. Sesja powstaje przy pierwszym
     poleceniu, więc otwarcie okna nie zakłada jeszcze sesji; moduł niesie
     własność `data-tworzy-sesje`, która mówi, czy polecenie w nim ją zakłada. */
  qq('[data-okno-nowe]').forEach(function (b) {
    b.addEventListener('click', function () {
      var co = b.getAttribute('data-okno-nowe') || b.getAttribute('data-modul') || 'Okno robocze';
      var tworzy = b.getAttribute('data-tworzy-sesje');
      var opis = 'Nowe okno robocze otwiera się na karcie Centrum dowodzenia, w izolacji: własna gałąź plików.';
      if (tworzy === 'tak') {
        opis += ' Moduł ' + co + ' zakłada sesję przy pierwszym poleceniu.';
      } else if (tworzy === 'nie') {
        opis += ' Moduł ' + co + ' nie zakłada sesji — pracuje na zasobach konta.';
      } else {
        opis += ' Środowisko i moduł wybierasz na tej karcie.';
      }
      toast(co, opis, 'informacja', 4200);
      oglos(co + ' — nowe okno robocze.');
    });
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Enter' && e.ctrlKey && e.shiftKey) {
      e.preventDefault();
      ustawMaks(powloka.getAttribute('data-okno-max') !== 'tak');
      return;
    }
    if (e.key === 'Escape' && powloka.getAttribute('data-okno-max') === 'tak') {
      ustawMaks(false);
    }
  });

  /* ── Izolacja okna: dwa niezależne przełączniki ─────────────────────────── */

  var przelaczniki = {};
  qq('[data-izolacja-przelacz]').forEach(function (b) {
    przelaczniki[b.getAttribute('data-izolacja-przelacz')] = b;
  });

  function stan(b) { return b && b.getAttribute('aria-checked') === 'true'; }

  function ustawIzolacje() {
    var incognito = stan(przelaczniki.incognito);
    var galaz = stan(przelaczniki.galaz);
    var rodzaj = incognito ? 'incognito' : (galaz ? 'galaz' : 'brak');
    powloka.setAttribute('data-izolacja', rodzaj);
  }

  Object.keys(przelaczniki).forEach(function (klucz) {
    var b = przelaczniki[klucz];
    b.addEventListener('click', function () {
      var wl = stan(b);
      b.setAttribute('aria-checked', wl ? 'false' : 'true');
      /* Incognito bez izolacji plików nie ma sensu, więc włącza gałąź;
         rozłączenie zostaje w rękach Operatora. */
      if (klucz === 'incognito' && !wl && przelaczniki.galaz) {
        przelaczniki.galaz.setAttribute('aria-checked', 'true');
      }
      ustawIzolacje();
      if (klucz === 'incognito') {
        toast('Incognito', wl
          ? 'Okno widzi historię i pamięć sesji konta.'
          : 'Okno nie widzi historii, pamięci ani notatek innych sesji. Praca zapisuje się we własnej gałęzi.', 'informacja', 3600);
      } else {
        toast('Własna gałąź', wl
          ? 'Zmiany plików trafiają do gałęzi głównej.'
          : 'Zmiany plików trafiają do własnej gałęzi i czekają na scalenie.', 'informacja', 3600);
      }
    });
  });

  ustawIzolacje();
})();
