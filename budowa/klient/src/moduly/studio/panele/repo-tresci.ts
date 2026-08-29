/**
 * Panel Session Repository — katalog treści widocznych dla Operatora. Ten sam
 * wzorzec co `moduly/studio/tresci.ts`, trzymany osobno w katalogu paneli, żeby
 * ten panel nie dzielił jednego pliku treści z sześcioma równolegle pisanymi.
 *
 * Zdania złożone stoją w katalogu, nie w kodzie: katalog przekazuje się
 * tłumaczowi bez dostępu do kodu, więc zdanie zszyte w wyrażeniu byłoby dla
 * niego niewidoczne. Kod podstawia wartości pod {nawiasy} i wybiera formę
 * liczebnika, której polszczyzna nie zapisze jednym wzorcem.
 */

import { formaPo } from '../liczebnik.ts';

export const tresciRepo = {
  panel: {
    tytul: 'Session Repository',
  },

  karty: {
    etykietaPrzelacznika: 'Widok repozytorium sesji',
    wersje: 'Wersje',
    galezie: 'Gałęzie',
    dziennik: 'Dziennik',
    kopie: 'Kopie zapasowe',
  },

  stany: {
    brakOkna: 'Nie udało się otworzyć okna modułu.',
    oczekiwanieDokumentu: 'To okno nie ma jeszcze przypisanego dokumentu.',
    wczytywanie: 'Wczytywanie historii wersji…',
    wczytywanieGalezi: 'Wczytywanie gałęzi dokumentu…',
    wczytywanieDziennika: 'Wczytywanie dziennika czynności…',
    wczytywanieKopii: 'Wczytywanie kopii zapasowych…',
  },

  odmowa: {
    lista: 'Nie udało się wczytać historii wersji',
    przywroc: 'Nie udało się przywrócić wersji.',
    galaz: 'Nie udało się założyć gałęzi.',
    etykieta: 'Nie udało się zmienić etykiety.',
    eksport: 'Nie udało się wyeksportować historii.',
    odwolanie: 'Nie udało się utworzyć odwołania do wersji.',
    zalozycielska: 'Nie udało się wrócić do wersji założycielskiej.',
    galezie: 'Nie udało się wczytać gałęzi',
    scalenie: 'Nie udało się scalić gałęzi.',
    dziennik: 'Nie udało się wczytać dziennika czynności',
    cofniecie: 'Nie udało się cofnąć czynności.',
    ponowienie: 'Nie udało się ponowić czynności.',
    kopie: 'Nie udało się wczytać kopii zapasowych',
    kopiaZalozenie: 'Nie udało się założyć kopii zapasowej.',
    kopiaPrzywrocenie: 'Nie udało się przywrócić kopii zapasowej.',
  },

  pusto: {
    tytul: 'Brak wersji dokumentu',
    opis: 'Wersje pojawią się tu po pierwszym zapisie w dokumencie.',
    galezieTytul: 'Brak gałęzi',
    galezieOpis: 'Gałąź zakładasz przyciskiem „Rozgałęź” przy wybranej wersji.',
    dziennikTytul: 'Dziennik jest pusty',
    dziennikOpis: 'Czynności pojawią się tu po pierwszej zmianie w dokumencie.',
    kopieTytul: 'Brak kopii zapasowych',
    kopieOpis: 'Kopię zakładasz przyciskiem „Załóż kopię”; pozostałe powstają samoczynnie.',
    filtrTytul: 'Brak pozycji dla tego wyboru',
    filtrOpis: 'Zmień zawężenie u góry, żeby zobaczyć pozostałe pozycje.',
  },

  wiersz: {
    autorModel: 'model',
    autorUzytkownik: 'użytkownik',
    autorNieznany: 'autor nieznany',
    kluczowa: 'wersja kluczowa',
    zalozycielska: 'wersja założycielska',
    numer: 'Wersja {numer}',
  },

  akcje: {
    biezaca: 'Wersja bieżąca',
    przywroc: 'Przywróć',
    przywracanie: 'Przywracanie…',
    rozgalez: 'Rozgałęź',
    wiecej: 'Więcej',
    anuluj: 'Anuluj',
  },

  filtr: {
    szereg: 'Szereg wersji',
    wszystkie: 'Wszystkie',
    operator: 'Operatora',
    autozapis: 'Samoczynne',
    tylkoKluczowe: 'Tylko kluczowe',
    autor: 'Autor czynności',
    stan: 'Stan wpisu',
    stanCzynne: 'Stojące',
    stanCofniete: 'Cofnięte',
    niezapisane: 'Tylko z niezapisanymi zmianami',
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
    odwolanie: 'Utwórz odwołanie',
    odwolanieWBiegu: 'Tworzenie odwołania…',
  },

  stopka: {
    eksportuj: 'Eksportuj historię',
    eksportowanie: 'Eksportowanie…',
    zalozycielska: 'Wróć do wersji założycielskiej',
    zalozycielskaWBiegu: 'Powrót…',
  },

  galezie: {
    punktStartowy: 'Punkt startowy: {wersja}',
    scalona: 'scalona',
    scal: 'Scal',
    scalanie: 'Scalanie…',
    docelowa: 'Gałąź docelowa',
    konfliktyTytul: 'Scalenie wstrzymane — fragmenty sporne',
    konfliktNumer: 'Fragment {numer}',
    stronaScalana: 'Zostaw z gałęzi scalanej',
    stronaDocelowa: 'Zostaw z gałęzi docelowej',
    stronaObie: 'Zostaw obie treści',
    scalPonownie: 'Scal z tymi rozstrzygnięciami',
    porzucKonflikty: 'Porzuć scalanie',
  },

  dziennik: {
    cofnij: 'Cofnij',
    cofanie: 'Cofanie…',
    ponow: 'Ponów',
    ponawianie: 'Ponawianie…',
    cofnieta: 'cofnięta',
    zaleznosc: 'stoi pod późniejszą czynnością',
    pominieto: 'Pominięto:',
  },

  kopie: {
    zaloz: 'Załóż kopię',
    zakladanie: 'Zakładanie kopii…',
    przywroc: 'Przywróć',
    doNowego: 'Przywróć do nowego dokumentu',
    przywracanie: 'Przywracanie…',
    niezapisane: 'niesie zmiany niezapisane',
    nieudana: 'zapis kopii się nie powiódł',
    powodInterval: 'zapis samoczynny',
    powodEvent: 'zdarzenie okna',
    powodPrzedNieodwracalna: 'przed czynnością nieodwracalną',
    powodReczny: 'polecenie Operatora',
  },

  rodzajCzynnosci: {
    textEdit: 'zmiana treści',
    formatChange: 'zmiana postaci',
    styleChange: 'zmiana stylu',
    pageChange: 'zmiana strony',
    listChange: 'zmiana listy',
    tableChange: 'zmiana tabeli',
    objectChange: 'zmiana obiektu',
    apparatusChange: 'zmiana aparatu dokumentu',
    markupChange: 'zmiana znakowania',
    proposalAccept: 'przyjęcie propozycji',
    importChange: 'wniesienie treści',
    fieldChange: 'zmiana pola',
    nieznany: 'czynność',
  },

  niegotowe: {
    podgladPorownanie:
      'Podgląd wersji i porównanie dwóch wersji prowadzą okna Preview Window i Diff/Grep Panel — z tej strefy jeszcze się ich nie otwiera.',
    trescWlasnaKonfliktu:
      'Wpisanie własnej treści w miejsce spornego fragmentu nie jest jeszcze dostępne; do wyboru stoją treści obu gałęzi.',
    brakCeluScalania:
      'Scalenie wymaga drugiej gałęzi jako docelowej. Wykaz gałęzi nie niesie nazwy gałęzi głównej, więc do niej scalić się stąd nie da.',
  },

  komunikat: {
    przywrocono: 'Dokument przywrócony do wskazanej wersji.',
    etykietaZapisana: 'Etykieta wersji zapisana.',
    wyeksportowano: 'Wyeksportowano {liczba} {rzeczownik} ({rozmiar}).',
    galazZalozona: 'Założono gałąź „{nazwa}”.',
    licznikWersjiJedna: '1 wersja',
    licznikWersji: '{liczba} {rzeczownik}',
    odwolanieUtworzone: 'Odwołanie do wersji: {odwolanie}',
    wrocono: 'Dokument wrócił do wersji założycielskiej.',
    scalono: 'Gałąź „{nazwa}” scalona.',
    konflikty: 'Scalenie wstrzymane: {liczba} {rzeczownik} do rozstrzygnięcia.',
    cofnieto: 'Cofnięto {liczba} {rzeczownik}.',
    cofnietoNic: 'Żadna czynność nie została cofnięta.',
    ponowiono: 'Ponowiono {liczba} {rzeczownik}.',
    ponowionoNic: 'Żadna czynność nie została ponowiona.',
    kopiaZalozona: 'Kopia zapasowa założona.',
    kopiaNieudana: 'Kopia zapasowa nie została zapisana.',
    kopiaPrzywrocona: 'Kopia zapasowa przywrócona do dokumentu.',
    kopiaPrzywroconaNowy: 'Kopia zapasowa przywrócona jako nowy dokument „{nazwa}”.',
    kopiaPrzywroconaNowyBezNazwy: 'Kopia zapasowa przywrócona jako nowy dokument bez nazwy.',
    szeregi: 'Operatora: {operatora} · samoczynnych: {samoczynnych}',
    niezapisaneKopie: 'Kopii ze zmianami niezapisanymi: {liczba}',
    wpisyDziennika: '{pokazane} z {wszystkich}',
  },

  /** Formy rzeczownika „wersja” po liczebniku; polszczyzna wymaga trzech. */
  odmianaWersji: {
    jedna: 'wersję',
    kilka: 'wersje',
    wiele: 'wersji',
  },

  odmianaCzynnosci: {
    jedna: 'czynność',
    kilka: 'czynności',
    wiele: 'czynności',
  },

  odmianaKonfliktow: {
    jedna: 'fragment sporny',
    kilka: 'fragmenty sporne',
    wiele: 'fragmentów spornych',
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

/** Komunikat po scaleniu gałęzi — nazwa pochodzi z wykazu gałęzi rdzenia. */
export function komunikatScalenia(nazwa: string): string {
  return zloz(tresciRepo.komunikat.scalono, { nazwa });
}

/** Komunikat o wstrzymanym scaleniu: ile fragmentów spornych czeka na rozstrzygnięcie. */
export function komunikatKonfliktow(liczba: number): string {
  const rzeczownik = formaPo(liczba, tresciRepo.odmianaKonfliktow);
  return zloz(tresciRepo.komunikat.konflikty, { liczba, rzeczownik });
}

/** Komunikat po cofnięciu czynności dziennika; zero mówi wprost, że nic nie zeszło. */
export function komunikatCofniecia(liczba: number): string {
  if (liczba === 0) return tresciRepo.komunikat.cofnietoNic;
  const rzeczownik = formaPo(liczba, tresciRepo.odmianaCzynnosci);
  return zloz(tresciRepo.komunikat.cofnieto, { liczba, rzeczownik });
}

/** Komunikat po ponowieniu czynności dziennika. */
export function komunikatPonowienia(liczba: number): string {
  if (liczba === 0) return tresciRepo.komunikat.ponowionoNic;
  const rzeczownik = formaPo(liczba, tresciRepo.odmianaCzynnosci);
  return zloz(tresciRepo.komunikat.ponowiono, { liczba, rzeczownik });
}

/** Odwołanie do wersji w obrębie platformy, wydane przez rdzeń. */
export function komunikatOdwolania(odwolanie: string): string {
  return zloz(tresciRepo.komunikat.odwolanieUtworzone, { odwolanie });
}

/** Komunikat po przywróceniu kopii do nowego dokumentu — nazwę nadaje rdzeń. */
export function komunikatKopiiNowyDokument(nazwa: string): string {
  return zloz(tresciRepo.komunikat.kopiaPrzywroconaNowy, { nazwa });
}

/** Rozbicie wykazu wersji na szereg Operatora i szereg zapisów samoczynnych. */
export function opisSzeregow(operatora: number, samoczynnych: number): string {
  return zloz(tresciRepo.komunikat.szeregi, { operatora, samoczynnych });
}

/** Ile kopii zapasowych niesie zmiany niezapisane w dokumencie. */
export function opisNiezapisanych(liczba: number): string {
  return zloz(tresciRepo.komunikat.niezapisaneKopie, { liczba });
}

/** Znacznik dziennika: ile wpisów widać z ilu, gdy rdzeń przyciął wykaz. */
export function opisWpisow(pokazane: number, wszystkich: number): string {
  return zloz(tresciRepo.komunikat.wpisyDziennika, { pokazane, wszystkich });
}

/** Podpis wiersza wersji: numer liczony od najstarszej w wykazie. */
export function podpisWersji(numer: number): string {
  return zloz(tresciRepo.wiersz.numer, { numer });
}

/** Podpis punktu startowego gałęzi — wersja, od której gałąź odeszła. */
export function podpisPunktuStartowego(wersja: string): string {
  return zloz(tresciRepo.galezie.punktStartowy, { wersja });
}

/** Podpis fragmentu spornego w wykazie konfliktów scalania. */
export function opisKonfliktu(numer: number): string {
  return zloz(tresciRepo.galezie.konfliktNumer, { numer });
}
