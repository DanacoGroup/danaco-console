// Wspólna warstwa makiet Danaco Console łączy przełącznik motywu z zapasowym mechanizmem obsługi poleceń natywnych, gdy system operacyjny ich nie dostarcza.
(function () {
  var r = document.documentElement;
  function ustaw(motyw) {
    r.dataset.theme = motyw;
    try { localStorage.setItem('dn-motyw', motyw); } catch (e) {}
    document.querySelectorAll('[data-motyw-ikona]').forEach(function (el) {
      el.querySelectorAll('.ik-slonce').forEach(function (i) { i.style.display = motyw === 'dark' ? '' : 'none'; });
      el.querySelectorAll('.ik-ksiezyc').forEach(function (i) { i.style.display = motyw === 'dark' ? 'none' : ''; });
    });
  }
  var zap = null;
  try { zap = localStorage.getItem('dn-motyw'); } catch (e) {}
  ustaw(zap || (matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'));
  document.addEventListener('click', function (e) {
    var p = e.target.closest('[data-przelacz-motyw]');
    if (p) ustaw(r.dataset.theme === 'dark' ? 'light' : 'dark');
  });
  // Zapas dla przeglądarek bez Invoker Commands (commandfor/command)
  if (!('commandForElement' in HTMLButtonElement.prototype)) {
    document.addEventListener('click', function (e) {
      var b = e.target.closest('button[commandfor]');
      if (!b) return;
      var cel = document.getElementById(b.getAttribute('commandfor'));
      if (!cel) return;
      var cmd = b.getAttribute('command');
      if (cmd === 'show-modal' && cel.showModal) cel.showModal();
      else if (cmd === 'close' && cel.close) cel.close();
      else if (cmd === 'toggle-popover' && cel.togglePopover) cel.togglePopover();
    });
  }
})();

/* Obszar, który się przewija, musi dać się osiągnąć z klawiatury. Bez tego treść
   dłuższa niż pojemnik jest dostępna wyłącznie myszą — czytnik ekranu i nawigacja
   tabulatorem zatrzymują się na krawędzi i reszta wpisu przestaje istnieć.

   Warunek jest mierzony, nie zakładany: znacznik dostaje ognisko tylko wtedy, gdy
   naprawdę przewija się w pionie ORAZ nie ma w środku niczego ogniskowalnego.
   Ryczałtowe `tabindex` na każdym pojemniku dokładałoby przystanki tabulatora
   w miejscach, gdzie użytkownik i tak dojdzie do treści przyciskiem lub odnośnikiem. */
(function () {
  var OGNISKOWALNE = 'a[href],button,input,select,textarea,[tabindex]:not([tabindex="-1"])';

  function udostepnijPrzewijane() {
    document.querySelectorAll('*').forEach(function (el) {
      if (el.hasAttribute('tabindex')) return;
      var styl = getComputedStyle(el);
      var przewijalne = (el.scrollHeight > el.clientHeight + 1 &&
                         (styl.overflowY === 'auto' || styl.overflowY === 'scroll')) ||
                        (el.scrollWidth > el.clientWidth + 1 &&
                         (styl.overflowX === 'auto' || styl.overflowX === 'scroll'));
      if (!przewijalne) return;
      if (el.querySelector(OGNISKOWALNE)) return;
      el.setAttribute('tabindex', '0');
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', udostepnijPrzewijane);
  } else {
    udostepnijPrzewijane();
  }
  /* Odsłonięcie karty albo panelu zmienia to, co się przewija. */
  document.addEventListener('click', function () { setTimeout(udostepnijPrzewijane, 0); });
})();
