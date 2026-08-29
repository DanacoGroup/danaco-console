/**
 * Moduł Studio — katalog treści, jedyne miejsce z tekstem widocznym dla
 * Operatora w tym module. Ten sam wzorzec co w `rama/tresci.ts`.
 */

export type WezelTresci = string | number | { [klucz: string]: WezelTresci };

export const tresci = {
  okno: {
    etykieta: 'Okno robocze — Studio',
  },

  pasmo: {
    etykietaKart: 'Karty okna roboczego',
    nowe: 'Nowe okno pomocnicze',
    maksymalizuj: 'Maksymalizuj okno robocze',
  },

  karty: {
    editor: 'Studio Editor',
    tools: 'Tools',
    diff: 'Diff',
    repo: 'Repo',
    preview: 'Preview',
    pliki: 'Pliki',
    plan: 'Plan',
  },

  wstazka: {
    etykieta: 'Wstążka okna roboczego — narzędzia modułu Studio',
    bezSesji: 'Sesja bez nazwy',
    ukladEtykieta: 'Układ okna roboczego',
    oknoKomunikacji: 'Okno komunikacji',
    podzialPionowy: 'Podział pionowy',
    wersjeEtykieta: 'Wersje dokumentu',
    zapiszWersje: 'Zapisz wersję',
    porownajWersje: 'Porównaj wersje',
    podgladWydruku: 'Podgląd wydruku',
    dostosuj: 'Dostosuj wstążkę okna roboczego',
    operacjeEtykieta: 'Operacje modułu',
    uruchomOperacje: 'Uruchom operację',
    wyslijDoLibrary: 'Wyślij do Library',
    brakKomendyLibrary: 'Wysyłanie dokumentu do Library nie jest dostępne w tej wersji.',
  },

  szyna: {
    etykieta: 'Dokumenty sesji',
    nowyDokument: '+ Nowy dokument',
    filtry: 'Filtr, sortowanie, grupowanie',
    brakDokumentow: 'Brak dokumentów w tej sesji.',
    bezTytulu: 'Dokument bez tytułu',
    nazwaNowego: 'Nowy dokument',
    brakPorzadkowania: 'Filtrowanie i sortowanie nie jest dostępne w tej wersji.',
    odmowaBezOpisu: 'Nie udało się wykonać czynności. Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.',
  },

  czat: {
    tytul: 'Chat Window',
    etykietaKontekst: 'Kontekst',
    brakParametrow: 'Parametry tej rozmowy nie zostały ustalone.',
    brakKontekstu: 'Brak przypiętego kontekstu.',
    brakWiadomosci: 'Brak wiadomości w tej rozmowie.',
    etykietaTresci: 'Pole polecenia',
    zastepczaTresc: 'Opisz operację na dokumencie lub zaznaczeniu…',
    wyslij: 'Do kolejki',
    wczytywanie: 'Wczytywanie rozmowy…',
    odmowaWykazu: 'Nie udało się wczytać rozmowy',
    stanBlad: 'zakończona błędem',
    stanZatrzymana: 'zatrzymana',
    /* Nazwa nadawcy wpisu, kluczowana rolą z kontraktu. */
    role: {
      user: 'Operator',
      assistant: 'Inteligencja',
      system: 'Platforma',
      tool: 'Narzędzie',
    },
  },

  /* Bez obietnicy terminu: „będzie dostępna w kolejnej wersji" mówi
     o harmonogramie zespołu, a nie o programie, który Operator ma przed sobą. */
  panelNiegotowy: {
    tytul: 'Niedostępne w tej wersji',
    opis: 'W tej wersji programu ta karta nie ma czynności.',
  },

  dokument: {
    tytul: 'Studio Editor',
    nazwaNowego: 'Dokument bez tytułu',
    zakladanie: 'Otwieranie okna roboczego…',
    etykietaTresci: 'Treść dokumentu',
    zapisz: 'Zapisz',
    zapisywanie: 'Zapisywanie…',
    zapisany: 'Zapisano',
    wersja: 'wersja',
    bezWersji: 'bez wersji',
  },

  /* Pasek narzędziowy edytora — strefa `dn-edytor-pasek` prototypu. Etykiety
     nazywają czynność, bo znak sam jej nie nazywa czytnikowi ekranu. */
  pasekEdytora: {
    etykieta: 'Narzędzia tekstu',
    pogrubienie: 'Pogrubienie',
    kursywa: 'Kursywa',
    podkreslenie: 'Podkreślenie',
    naglowek: 'Nagłówek',
    stylAkapitu: 'Styl akapitu',
    lista: 'Lista punktowana',
    tabela: 'Tabela',
    kod: 'Blok kodu',
    cytat: 'Cytat',
    szukaj: 'Znajdź i zamień',
    /* Czynność bez pokrycia w rdzeniu nazywa niegotowość, zamiast udawać
       działanie — zasada zero blokad z `komponenty.css`. */
    zapowiedziane: 'Ta czynność wejdzie w kolejnym wydaniu.',
    bezZaznaczenia: 'Zaznacz fragment tekstu, żeby zmienić jego postać.',
  },

  /* Pas stanu edytora — strefa `st-status` prototypu. Wartości pochodzą
     z dokumentu rdzenia; żadna nie jest wpisana na sztywno. */
  stanEdytora: {
    /* Formy odmienne: „1 słowo”, „2 słowa”, „5 słów” — liczba wybiera formę,
       więc jeden napis nie wystarcza. */
    slowoJedna: 'słowo',
    slowoKilka: 'słowa',
    slowoWiele: 'słów',
    zapisany: 'zapisano',
    zapisywanie: 'zapisywanie…',
    niezapisany: 'niezapisane zmiany',
    wersja: 'wersja',
    bezWersji: 'bez wersji',
    /* Postać edycji nazwana wprost: pole niesie tekst, nie widok składu.
       „WYSIWYG” z prototypu opisuje postać, której edytor jeszcze nie ma. */
    postac: 'Tekst',
  },

  dokumentOdmowa: {
    sesja: 'Nie udało się założyć sesji',
    kanal: 'Żaden model nie jest dostępny',
    okno: 'Nie udało się otworzyć okna modułu',
    dokument: 'Nie udało się założyć dokumentu',
    zapis: 'Nie udało się zapisać dokumentu',
  },

  odmowa: {
    brakOpisu: 'Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.',
  },
} as const;
