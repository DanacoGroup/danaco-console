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
    brakOkna: 'Nie udało się otworzyć okna modułu.',
    oczekiwanieDokumentu: 'To okno nie ma jeszcze przypisanego dokumentu.',
    wczytywanie: 'Wczytywanie historii wersji…',
  },

  odmowa: {
    lista: 'Nie udało się wczytać historii wersji',
    przywroc: 'Nie udało się przywrócić wersji.',
    galaz: 'Nie udało się założyć gałęzi.',
    etykieta: 'Nie udało się zmienić etykiety.',
    eksport: 'Nie udało się wyeksportować historii.',
  },

  pusto: {
    tytul: 'Brak wersji dokumentu',
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
    /* Zdania złożone stoją w katalogu, jak wszystkie pozostałe, a nie w kodzie:
       katalog przekazuje się tłumaczowi bez dostępu do kodu, więc zdanie zszyte
       w wyrażeniu z kodu byłoby dla niego niewidoczne. Kod podstawia wyłącznie
       wartości pod {nawiasy} i wybiera formę liczebnika, której polszczyzna
       nie da się zapisać jednym wzorcem. */
    wyeksportowano: 'Wyeksportowano {liczba} {rzeczownik} ({rozmiar}).',
    galazZalozona: 'Założono gałąź „{nazwa}”.',
    licznikWersjiJedna: '1 wersja',
    licznikWersji: '{liczba} {rzeczownik}',
  },

  /** Formy rzeczownika „wersja” po liczebniku; polszczyzna wymaga trzech. */
  odmianaWersji: {
    jedna: 'wersję',
    kilka: 'wersje',
    wiele: 'wersji',
  },
} as const;

/** Podstawia wartości pod {nawiasy} we wzorcu z katalogu. */
function zloz(wzorzec: string, dane: Record<string, string | number>): string {
  return wzorzec.replace(/\{(\w+)\}/g, (calosc, klucz: string) => {
    const wartosc = dane[klucz];
    return wartosc === undefined ? calosc : String(wartosc);
  });
}

/** Forma dopełniacza po liczebniku: 2–4 (poza nastoma) bierze „wersje”, reszta „wersji”. */
function odmianaWersjiPo(liczba: number): string {
  const ostatniaCyfra = liczba % 10;
  const ostatnieDwieCyfry = liczba % 100;
  const jestNastolatkiem = ostatnieDwieCyfry >= 12 && ostatnieDwieCyfry <= 14;
  return ostatniaCyfra >= 2 && ostatniaCyfra <= 4 && !jestNastolatkiem
    ? tresciRepo.odmianaWersji.kilka
    : tresciRepo.odmianaWersji.wiele;
}

/** Znacznik nagłówka panelu: liczba wersji zwrócona przez rdzeń, np. „7 wersji”. */
export function znacznikWersji(liczba: number): string {
  if (liczba === 1) return tresciRepo.komunikat.licznikWersjiJedna;
  return zloz(tresciRepo.komunikat.licznikWersji, { liczba, rzeczownik: odmianaWersjiPo(liczba) });
}

/** Komunikat po eksporcie: liczba spakowanych wersji i rozmiar archiwum. */
export function komunikatEksportu(wpisy: number, rozmiar: string): string {
  const rzeczownik = wpisy === 1 ? tresciRepo.odmianaWersji.jedna : odmianaWersjiPo(wpisy);
  return zloz(tresciRepo.komunikat.wyeksportowano, { liczba: wpisy, rzeczownik, rozmiar });
}

/** Komunikat po założeniu gałęzi — nazwa jest tekstem wpisanym przez Operatora. */
export function komunikatGalezi(nazwa: string): string {
  return zloz(tresciRepo.komunikat.galazZalozona, { nazwa });
}
