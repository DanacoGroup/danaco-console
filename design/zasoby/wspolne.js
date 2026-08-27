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
