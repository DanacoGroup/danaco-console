import type { NazwaIkony } from '../ikony/ikony';

/**
 * Pozycje listwy ustawień — strefa trzecia strony głównej. Wyłącznie treść
 * listwy, bez elementów i bez stylu.
 *
 * Strefa ma najniższą wagę wizualną i jest stale dostępna. Żadna pozycja nie
 * jest bramą ani przełącznikiem stanu — każda otwiera widok i każda jest
 * klikalna zawsze.
 */

/** Kod pozycji ustawień przekazywany w zdarzeniu wyboru. */
export type KodUstawienia =
  | 'konfiguracja'
  | 'dostepy'
  | 'modele'
  | 'ustawienia'
  | 'punkty-izolacji'
  | 'mobile'
  | 'aod';

export interface PozycjaUstawienia {
  kod: KodUstawienia;
  /** Etykieta widoczna obok ikony. */
  nazwa: string;
  /** Rozwinięcie nazwy — czytane przez technologie wspomagające. */
  wyjasnienie: string;
  ikona: NazwaIkony;
}

export const POZYCJE_USTAWIEN: readonly PozycjaUstawienia[] = [
  {
    kod: 'konfiguracja',
    nazwa: 'Okno konfiguracji',
    wyjasnienie: 'Ustawienia platformy i zasięgi konfiguracji',
    ikona: 'ustawienia',
  },
  {
    kod: 'dostepy',
    nazwa: 'Dostępy i katalog roboczy',
    wyjasnienie: 'Punkty dostępu, nadania dla okien rozmowy i katalog roboczy modelu',
    ikona: 'klodka',
  },
  {
    kod: 'modele',
    nazwa: 'Modele, konta i tożsamość',
    wyjasnienie: 'Konta modeli i code CLI, ustawienia per model i per konto, zasady i tożsamość modelu',
    ikona: 'agent',
  },
  {
    kod: 'ustawienia',
    nazwa: 'Ustawienia',
    wyjasnienie:
      'Uwierzytelnianie operatora — konta, kategorie i dokumenty tożsamości modelu mieszkają w oknie „Modele, konta i tożsamość”',
    ikona: 'uzytkownik',
  },
  {
    kod: 'punkty-izolacji',
    nazwa: 'Punkty izolacji',
    wyjasnienie: 'Kontekst, zakres techniczny, poziom zapisu, oś rozstrzygania i polityka efektywna izolacji',
    ikona: 'tarcza',
  },
  {
    kod: 'mobile',
    nazwa: 'Mobile',
    wyjasnienie: 'Widok nadzoru na telefonie i tablecie',
    ikona: 'karta-przegladarki',
  },
  {
    kod: 'aod',
    nazwa: 'Always On Display',
    wyjasnienie: 'Pływający awatar obecny w całej platformie',
    ikona: 'oko',
  },
];

/**
 * Pozycja „Okno konfiguracji" wyjęta z wykazu pod skrót w pasku aplikacji.
 *
 * Wyjęta, a nie przepisana: skrót prowadzi tam, gdzie pozycja listwy. Druga
 * definicja tej samej pozycji rozjechałaby się przy pierwszej zmianie nazwy
 * albo ikony.
 */
export const POZYCJA_KONFIGURACJI: PozycjaUstawienia =
  POZYCJE_USTAWIEN.find((pozycja) => pozycja.kod === 'konfiguracja') ?? POZYCJE_USTAWIEN[0]!;

/**
 * Kody, które wchodzą na listwę ustawień strony głównej.
 *
 * Wykaz jest zamknięty i liczy trzy pozycje, bo tyle wymienia tabela
 * „Zawartość listwy ustawień" opracowania (rozdz. 3.4): Okno konfiguracji,
 * Mobile, Always On Display. Pozostałe pozycje — uwierzytelnianie operatora,
 * punkty izolacji, dostępy oraz modele i tożsamość — są zakresami Okna
 * konfiguracji, a nie bytami obok niego; własne wejście na listwie dawałoby
 * dwie drogi do tej samej rzeczy. Pełny wykaz niesie dalej menu aplikacji
 * w pasku (`aplikacja/menu-aplikacji.ts`), gdzie hierarchii stref nie ma.
 */
export const KODY_STREFY_TRZECIEJ: readonly KodUstawienia[] = [
  'konfiguracja',
  'mobile',
  'aod',
];
