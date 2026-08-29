/**
 * Panel Preview Window — własny katalog treści. Ten sam wzorzec co
 * `moduly/studio/tresci.ts`: jedyne miejsce tego panelu z tekstem widocznym
 * dla Operatora, osobno od katalogu współdzielonego, bo panel mieszka
 * wyłącznie w plikach tego katalogu.
 */

export const trescPodgladu = {
  brakOkna: 'Rdzeń nie założył okna modułu.',

  brakDokumentuTytul: 'Brak dokumentu w podglądzie',
  brakDokumentuOpis: 'Rdzeń nie otworzył jeszcze żadnego dokumentu w tym oknie.',

  renderowanie: 'Renderowanie podglądu…',

  brakStronTytul: 'Brak stron do pokazania',
  brakStronOpis: 'Rdzeń nie zwrócił żadnej strony podglądu.',

  odmowaRenderu: 'Rdzeń odmówił wygenerowania podglądu',
  odmowaUkladu: 'Rdzeń odmówił odczytu ustawień strony',
  brakOpisu: 'Rdzeń nie podał powodu odmowy.',

  format: 'Format wydania',
  strona: 'Strona',
  ciagly: 'Ciągły',
  powiekszenie: 'Powiększenie podglądu',
  poprzednia: 'Poprzednia',
  nastepna: 'Następna',

  podzielEkran: 'Podziel ekran',
  eksportuj: 'Profil wydania',
  profilDomyslny: 'Bez profilu wydania',
  profilBrak: 'Rdzeń nie zgłosił żadnego profilu wydania.',
  profilOdmowa: 'Rdzeń odmówił wykazu profili wydania',
  profilWczytywanie: 'Wczytywanie profili wydania…',

  wyslijDoLibrary: 'Wyślij do Library',
  wyslijNiegotowe: 'Wejdzie osobnym zakresem prac.',
} as const;
