/**
 * Panel Preview Window — własny katalog treści. Ten sam wzorzec co
 * `moduly/studio/tresci.ts`: jedyne miejsce tego panelu z tekstem widocznym
 * dla Operatora, osobno od katalogu współdzielonego, bo panel mieszka
 * wyłącznie w plikach tego katalogu.
 */

export const trescPodgladu = {
  brakOkna: 'Nie udało się otworzyć okna modułu.',

  brakDokumentuTytul: 'Brak dokumentu w podglądzie',
  brakDokumentuOpis: 'W tym oknie nie ma jeszcze otwartego dokumentu.',

  renderowanie: 'Przygotowywanie podglądu…',

  brakStronTytul: 'Brak stron do pokazania',
  brakStronOpis: 'Brak stron do wyświetlenia.',

  odmowaRenderu: 'Nie udało się przygotować podglądu',
  odmowaUkladu: 'Nie udało się odczytać ustawień strony',
  brakOpisu: 'Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.',

  format: 'Format wydania',
  strona: 'Strona',
  ciagly: 'Ciągły',
  powiekszenie: 'Powiększenie podglądu',
  poprzednia: 'Poprzednia',
  nastepna: 'Następna',

  podzielEkran: 'Podziel ekran',
  eksportuj: 'Profil wydania',
  profilDomyslny: 'Bez profilu wydania',
  profilBrak: 'Brak profili wydania.',
  profilOdmowa: 'Nie udało się wczytać profili wydania',
  profilWczytywanie: 'Wczytywanie profili wydania…',

  wyslijDoLibrary: 'Wyślij do Library',
  wyslijNiegotowe: 'Niedostępne w tej wersji.',
} as const;
