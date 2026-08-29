/**
 * Panel Preview Window — własny katalog treści. Ten sam wzorzec co
 * `moduly/studio/tresci.ts`: jedyne miejsce tego panelu z tekstem widocznym
 * dla Operatora, osobno od katalogu współdzielonego, bo panel mieszka
 * wyłącznie w plikach tego katalogu.
 *
 * Zdania złożone stoją tutaj, nie w kodzie: katalog przekazuje się tłumaczowi
 * bez dostępu do kodu. Kod podstawia wartości pod {nawiasy} i wybiera formę
 * liczebnika, której polszczyzna nie zapisze jednym wzorcem.
 */

import type { FormyLiczebnika } from '../liczebnik.ts';

export const trescPodgladu = {
  panel: {
    tytul: 'Preview Window',
    zamknijKarte: 'Zamknij kartę',
  },

  brakOkna: 'Nie udało się otworzyć okna modułu.',

  brakDokumentuTytul: 'Brak dokumentu w podglądzie',
  brakDokumentuOpis: 'Podgląd pokaże dokument otwarty w tym oknie, gdy tylko taki będzie.',

  renderowanie: 'Przygotowywanie podglądu…',

  /* Opis nie powtarza tytułu: tytuł nazywa stan, opis mówi, co go zdejmie. */
  brakStronTytul: 'Brak stron do pokazania',
  brakStronOpis: 'Podgląd pojawi się po zapisaniu treści dokumentu.',

  odmowaRenderu: 'Nie udało się przygotować podglądu',
  odmowaUkladu: 'Nie udało się odczytać ustawień strony',
  odmowaZmianyUkladu: 'Nie udało się zmienić ustawień strony',
  odmowaNosnikow: 'Nie udało się wczytać formatów nośnika',
  odmowaStrony: 'Nie udało się pobrać obrazu strony',

  format: 'Format wydania',
  strona: 'Strona',
  ciagly: 'Ciągły',
  powiekszenie: 'Powiększenie podglądu',
  poprzednia: 'Poprzednia',
  nastepna: 'Następna',
  stronaZe: '{numer} z {ile}',

  nosnik: 'Format nośnika',
  nosnikWczytywanie: 'Wczytywanie formatów…',
  nosnikBrak: 'Bez formatu nośnika',
  orientacja: 'Orientacja strony',
  pionowa: 'Pionowa',
  pozioma: 'Pozioma',
  ukladZmieniany: 'Zmiana ustawień strony…',

  stronaWczytywanie: 'Wczytywanie strony…',
  stronaBezObrazu: 'Rdzeń nie odłożył obrazu tej strony.',
  bezTytulu: 'Dokument bez tytułu',
  numerStrony: '— {numer} —',

  podzielEkran: 'Podziel ekran',
  eksportuj: 'Profil wydania',
  profilDomyslny: 'Bez profilu wydania',
  profilBrak: 'Brak profili wydania.',
  profilOdmowa: 'Nie udało się wczytać profili wydania',
  profilWczytywanie: 'Wczytywanie profili wydania…',

  zapiszProfil: 'Zapisz profil',
  zapiszProfilNazwa: 'Nazwa profilu wydania',
  zapiszProfilZnak: 'Nazwa nowego profilu',
  zapiszProfilZatwierdz: 'Zapisz',
  zapiszProfilPorzuc: 'Porzuć',
  zapiszProfilBiegnie: 'Zapisywanie profilu…',
  zapiszProfilOdmowa: 'Nie udało się zapisać profilu wydania',
  zapiszProfilBezNazwy: 'Podaj nazwę profilu, żeby go zapisać.',

  wyslijDoLibrary: 'Wyślij do Library',
  wyslijNiegotowe:
    'Wysyłka podglądu do Library nie jest jeszcze gotowa — brak dla niej czynności w umowie z rdzeniem.',

  stronyForma: { jedna: 'strona', kilka: 'strony', wiele: 'stron' } satisfies FormyLiczebnika,
} as const;
