import type { NazwaIkony } from '../ikony/ikony';

/** Pozycje listwy ustawień strefy trzeciej strony głównej niosą wyłącznie treść, bez elementów i bez stylu, każda klikalna zawsze i żadna nie jest bramą ani przełącznikiem stanu. */
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

/** Pozycja okna konfiguracji jest wyjęta z wykazu, nie przepisana, pod skrót w pasku aplikacji, żeby druga definicja nie rozjechała się przy zmianie nazwy albo ikony. */
export const POZYCJA_KONFIGURACJI: PozycjaUstawienia =
  POZYCJE_USTAWIEN.find((pozycja) => pozycja.kod === 'konfiguracja') ?? POZYCJE_USTAWIEN[0]!;

/** Kody wchodzące na listwę ustawień strony głównej są zamkniętym wykazem trzech pozycji; pozostałe pozycje są zakresami okna konfiguracji, nie osobnymi bytami obok niego. */
export const KODY_STREFY_TRZECIEJ: readonly KodUstawienia[] = [
  'konfiguracja',
  'mobile',
  'aod',
];
