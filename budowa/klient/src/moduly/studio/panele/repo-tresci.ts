/**
 * Panel Session Repository — katalog treści widocznych dla Operatora. Ten sam
 * wzorzec co `moduly/studio/tresci.ts`, trzymany osobno w katalogu paneli, żeby
 * ten panel nie dzielił jednego pliku treści z sześcioma równolegle pisanymi.
 */

export const tresciRepo = {
  panel: {
    tytul: 'Session Repository',
  },

  stany: {
    brakOkna: 'Rdzeń nie założył okna modułu.',
    oczekiwanieDokumentu: 'Rdzeń nie przypisał jeszcze dokumentu temu oknu.',
    wczytywanie: 'Wczytywanie historii wersji…',
  },

  odmowa: {
    lista: 'Rdzeń odmówił historii wersji',
    przywroc: 'Rdzeń odmówił przywrócenia wersji.',
    galaz: 'Rdzeń odmówił założenia gałęzi.',
    etykieta: 'Rdzeń odmówił zmiany etykiety.',
    eksport: 'Rdzeń odmówił eksportu historii.',
    brakOpisu: 'Rdzeń nie podał powodu odmowy.',
  },

  pusto: {
    tytul: 'Rdzeń nie zgłosił żadnej wersji dokumentu',
    opis: 'Wersje pojawią się tu po pierwszym zapisie w dokumencie.',
  },

  wiersz: {
    autorModel: 'model',
    autorUzytkownik: 'użytkownik',
    autorNieznany: 'autor nieznany',
    kluczowa: 'wersja kluczowa',
  },

  akcje: {
    biezaca: 'Wersja bieżąca',
    przywroc: 'Przywróć',
    przywracanie: 'Przywracanie…',
    rozgalez: 'Rozgałęź',
    wiecej: 'Więcej',
    anuluj: 'Anuluj',
  },

  formularzGalaz: {
    etykieta: 'Nazwa gałęzi',
    zastepcza: 'np. redakcja równoległa',
    zaloz: 'Załóż gałąź',
    zakladanie: 'Zakładanie…',
  },

  formularzEtykieta: {
    etykieta: 'Nazwa wersji',
    kluczowa: 'Wersja kluczowa',
    zapisz: 'Zapisz etykietę',
    zapisywanie: 'Zapisywanie…',
  },

  stopka: {
    eksportuj: 'Eksportuj historię',
    eksportowanie: 'Eksportowanie…',
  },

  komunikat: {
    przywrocono: 'Dokument przywrócony do wskazanej wersji.',
    etykietaZapisana: 'Etykieta wersji zapisana.',
  },
} as const;

/** Odmiana rzeczownika „wersja” po liczebniku — forma dopełniacza (2–4 / 5+). */
function odmianaWersjiPo(liczba: number): 'wersje' | 'wersji' {
  const ostatniaCyfra = liczba % 10;
  const ostatnieDwieCyfry = liczba % 100;
  const jestNastolatkiem = ostatnieDwieCyfry >= 12 && ostatnieDwieCyfry <= 14;
  return ostatniaCyfra >= 2 && ostatniaCyfra <= 4 && !jestNastolatkiem ? 'wersje' : 'wersji';
}

/** Znacznik nagłówka panelu: liczba wersji zwrócona przez rdzeń, np. „7 wersji”. */
export function znacznikWersji(liczba: number): string {
  return liczba === 1 ? '1 wersja' : `${liczba} ${odmianaWersjiPo(liczba)}`;
}

/** Komunikat po eksporcie: liczba spakowanych wersji i rozmiar archiwum. */
export function komunikatEksportu(wpisy: number, rozmiar: string): string {
  const rzeczownik = wpisy === 1 ? 'wersję' : odmianaWersjiPo(wpisy);
  return `Wyeksportowano ${wpisy} ${rzeczownik} (${rozmiar}).`;
}

/** Komunikat po założeniu gałęzi — nazwa jest tekstem wpisanym przez Operatora. */
export function komunikatGalezi(nazwa: string): string {
  return `Założono gałąź „${nazwa}”.`;
}
