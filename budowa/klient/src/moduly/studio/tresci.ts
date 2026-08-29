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
    brakDokumentow: 'Rdzeń nie zgłosił żadnego dokumentu w tej sesji.',
    bezTytulu: '(dokument bez tytułu)',
    nazwaNowego: 'Nowy dokument',
    brakPorzadkowania: 'Rdzeń nie podaje pól, po których szyna mogłaby filtrować i sortować.',
    odmowaBezOpisu: 'Rdzeń odmówił i nie podał powodu.',
  },

  czat: {
    tytul: 'Chat Window',
    etykietaKontekst: 'Kontekst',
    brakParametrow: 'Rdzeń nie podał jeszcze parametrów tej rozmowy.',
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
    zakladanie: 'Zakładanie okna roboczego w rdzeniu…',
    etykietaTresci: 'Treść dokumentu',
    zapisz: 'Zapisz',
    zapisywanie: 'Zapisywanie…',
    zapisany: 'Zapisane w repozytorium sesji',
    wersja: 'wersja',
    bezWersji: 'bez wersji',
  },

  dokumentOdmowa: {
    sesja: 'Rdzeń nie założył sesji',
    kanal: 'Rdzeń nie zgłosił żadnego kanału modelu',
    okno: 'Rdzeń nie założył okna modułu',
    dokument: 'Rdzeń nie założył dokumentu',
    zapis: 'Rdzeń odmówił zapisu dokumentu',
  },

  odmowa: {
    brakOpisu: 'Rdzeń nie podał powodu odmowy.',
  },
} as const;
