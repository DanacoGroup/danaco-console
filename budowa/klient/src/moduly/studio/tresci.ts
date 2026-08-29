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
    bezSesji: '(bez nazwy sesji)',
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
    brakKomendyLibrary: 'Kontrakt nie niesie komendy odkładającej dokument Studia do Library.',
  },

  szyna: {
    etykieta: 'Dokumenty sesji',
    nowyDokument: '+ Nowy dokument',
    filtry: 'Filtr, sortowanie, grupowanie',
    brakDokumentow: 'Ta sesja nie ma jeszcze żadnego dokumentu.',
    bezTytulu: '(dokument bez tytułu)',
    nazwaNowego: 'Nowy dokument',
    brakPorzadkowania: 'Filtrowanie i sortowanie nie jest jeszcze dostępne.',
    odmowaBezOpisu: 'Nie udało się wykonać czynności; powód nie został podany.',
  },

  czat: {
    tytul: 'Chat Window',
    etykietaKontekst: 'Kontekst',
    brakParametrow: 'Parametry tej rozmowy nie są jeszcze ustalone.',
    brakKontekstu: 'Brak przypiętego kontekstu.',
    brakWiadomosci: 'Brak wiadomości w tej rozmowie.',
    etykietaTresci: 'Pole polecenia',
    zastepczaTresc: 'Opisz operację na dokumencie lub zaznaczeniu…',
    wyslij: 'Do kolejki',
  },

  panelNiegotowy: {
    tytul: 'Panel jeszcze nie powstał',
    opis: 'Ten panel wejdzie osobnym zakresem prac.',
  },

  dokument: {
    tytul: 'Studio Editor',
    nazwaNowego: 'Dokument bez tytułu',
    zakladanie: 'Otwieranie okna roboczego…',
    etykietaTresci: 'Treść dokumentu',
    zapisz: 'Zapisz',
    zapisywanie: 'Zapisywanie…',
    zapisany: 'Zapisane w repozytorium sesji',
    wersja: 'wersja',
    bezWersji: 'bez wersji',
  },

  dokumentOdmowa: {
    sesja: 'Nie udało się założyć sesji',
    kanal: 'Nie ma żadnego kanału modelu',
    okno: 'Nie udało się otworzyć okna modułu',
    dokument: 'Nie udało się założyć dokumentu',
    zapis: 'Nie udało się zapisać dokumentu',
  },

  odmowa: {
    brakOpisu: 'Powód nie został podany.',
  },
} as const;
