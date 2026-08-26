/* ============================================================================
   KREATOR INSTALACJI — zestaw ikon

   Rysunek znaku należy do zestawu, nie do składnika: ta sama ikona bywa
   potrzebna w kilku miejscach i musi wszędzie wyglądać tak samo. Składnik
   podaje jej klasę i rolę w dostępności, nigdy ścieżek.

   Brak ikony jest zgłoszeniem do Właściciela, nie powodem do narysowania
   własnej — patrz zasada 7 w zasoby/ARCHITEKTURA.md.
   ============================================================================ */
(function () {
'use strict';
var K = (window.DanacoKreator = window.DanacoKreator || {});

K.ikony = {
  pobranie: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3v12"/><path d="m7 10 5 5 5-5"/><path d="M5 21h14"/></svg>',
  godlo: '<svg viewBox="0 0 96 96"><path fill="currentColor" d="M12 26 H24 L44 48 L24 70 H12 L32 48 Z"/><path fill="currentColor" d="M40 26 H52 L72 48 L52 70 H40 L60 48 Z"/><circle class="kropka" cx="83" cy="63.5" r="6.5"/></svg>',
  zwin: '<svg viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"><path d="M2 7h10"/></svg>',
  rozwin: '<svg viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round"><rect x="2.2" y="2.2" width="9.6" height="9.6" rx="1.2"/></svg>',
  zamknij: '<svg viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"><path d="M3 3l8 8M11 3l-8 8"/></svg>',
  godloDuze: '<svg viewBox="0 0 96 96"><path fill="currentColor" d="M12 26 H24 L44 48 L24 70 H12 L32 48 Z"/><path fill="currentColor" d="M40 26 H52 L72 48 L52 70 H40 L60 48 Z"/><circle class="kropka" cx="83" cy="63.5" r="6.5"/></svg>',
  wifi: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M2.4 9.6a15 15 0 0 1 19.2 0"/><path d="M5.6 13.2a10 10 0 0 1 12.8 0"/><path d="M8.8 16.8a5 5 0 0 1 6.4 0"/><circle cx="12" cy="20.4" r="1"/></svg>',
  ptaszek: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m5 13 4 4L19 7"/></svg>',
  wynikGotowe: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"/><path d="m8 12 3 3 5-6"/></svg>',
  wynikOstrzezenia: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3.8 2.6 20h18.8Z"/><path d="M12 10v4"/><path d="M12 17.2h.01"/></svg>',
  ostrzezenie: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3.8 2.6 20h18.8Z"/><path d="M12 10v4"/><path d="M12 17.2h.01"/></svg>'
};
})();
