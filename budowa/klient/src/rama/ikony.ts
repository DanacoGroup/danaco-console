/**
 * Rama aplikacji — zestaw znaków. Godło przeniesione z prototypu bez zmiany
 * kształtu; pozostałe rysunki wzięte wprost z paska stanu Centrum dowodzenia.
 * Klasę i wymiar nakłada miejsce użycia, nie sam znak.
 */

export const ikony = {
  godlo:
    '<svg viewBox="0 0 96 96" aria-hidden="true"><path fill="currentColor" d="M12 26 H24 L44 48 L24 70 H12 L32 48 Z"/><path fill="currentColor" d="M40 26 H52 L72 48 L52 70 H40 L60 48 Z"/><circle class="kropka" cx="83" cy="63.5" r="6.5"/></svg>',
  karty:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="2" y="5" width="20" height="14" rx="2"/><path d="M2 10h20"/></svg>',
  slonce:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4"/></svg>',
  ksiezyc:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 14.5A8.5 8.5 0 0 1 9.5 4a8.5 8.5 0 1 0 10.5 10.5Z"/></svg>',
  modul:
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>',
} as const;

/** Nazwa znaku z zestawu, wskazująca jeden z gotowych rysunków dostępnych w tym pliku pod postacią klucza. */
export type NazwaZnaku = keyof typeof ikony;
