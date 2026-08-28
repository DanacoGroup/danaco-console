/**
 * Kod modułu Translate i kody jego okien operacyjnych, którymi moduł zgłasza się rdzeniowi serwera aplikacji.
 */

/**
 * Kod modułu Translate zarejestrowany w rejestrze modułów rdzenia serwera tej aplikacji Danaco Console.
 */
export const KOD_MODULU = 'translate';

/**
 * Kody okien operacyjnych budowanych przez ten moduł Translate w tej budowie aplikacji Danaco Console.
 */
export const KODY_OKIEN = {
  zrodlo: 'source-panel',
  panele: 'translation-panels',
  glosariusz: 'glossary-manager',
  pamiec: 'translation-memory-panel',
  formaty: 'format-studio',
  jakosc: 'qa-review-center',
  warsztat: 'translation-workshop',
} as const;
