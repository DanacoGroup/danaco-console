/**
 * Panel Diff/Grep — katalog treści. Jedyne miejsce panelu z tekstem widocznym
 * dla Operatora; poza nim żaden łańcuch w `diff.ts` nie czyta się jak zdanie.
 * Ten sam układ co `moduly/studio/tresci.ts`, osobny plik zgodnie z umową
 * panelu — treść własna panelu nie wchodzi do katalogu wspólnego.
 */

export const tresc = {
  panel: { tytul: 'Diff/Grep Panel' },

  brakDokumentu: 'Brak dokumentu w tym oknie.',

  porownanie: {
    etykietaBazowa: 'Wersja odniesienia',
    etykietaCel: 'Wersja porównywana',
    zastepczaWersja: 'identyfikator wersji…',
    porownaj: 'Porównaj',
    ladowanie: 'Wczytywanie porównania…',
    brakRoznic: 'Wskazane wersje nie różnią się.',
    /* Prototyp stawia w tym miejscu listę wersji do wyboru. Wykaz wersji
       prowadzi rodzina studio.repository.*, poza tym panelem, więc wersję
       Operator wskazuje identyfikatorem, a zdanie mówi, skąd go wziąć. */
    skadWersje: 'Wersję wskazujesz jej identyfikatorem; wykaz wersji prowadzi karta Repozytorium.',
    /* Kształt zapisu z prototypu: „Statystyka: +12 … · −48 … · 3 …”, ze znakiem
       kierunku przed liczbą. Jednostką są fragmenty różnicy, bo tyle podaje
       serwer — prototyp liczył słowa, których rdzeń nie zwraca. */
    statystykaEtykieta: 'Statystyka:',
    znakDodania: '+',
    znakUsuniecia: '−',
    jednostkaFragmenty: { jedna: 'fragment', kilka: 'fragmenty', wiele: 'fragmentów' },
    jednostkaZmian: { jedna: 'zmieniony fragment', kilka: 'zmienione fragmenty', wiele: 'zmienionych fragmentów' },
  },

  tryb: {
    etykieta: 'Zakres porównania',
    scalony: 'scalony',
    postac: 'postać',
    wyglad: 'wygląd',
    zrodlo: 'materiał wejściowy',
  },

  postac: {
    ladowanie: 'Wczytywanie różnic postaci…',
    brak: 'Postać wskazanych wersji nie różni się.',
    statystykaEtykieta: 'Postać:',
    jednostkaCech: { jedna: 'cecha', kilka: 'cechy', wiele: 'cech' },
    jednostkaZmienionych: { jedna: 'zmieniona cecha', kilka: 'zmienione cechy', wiele: 'zmienionych cech' },
  },

  wyglad: {
    ladowanie: 'Wczytywanie różnic wyglądu…',
    brak: 'Wygląd wskazanych wersji nie różni się.',
    wymaganeWersje: 'Porównanie wyglądu wymaga wskazania obu wersji: odniesienia i porównywanej.',
    obszar: 'Strona',
    udzial: 'różnica',
    nakladki: 'Nakładki różnicy odłożone w magazynie:',
    jednostkaNakladek: { jedna: 'nakładka', kilka: 'nakładki', wiele: 'nakładek' },
  },

  zrodlo: {
    ladowanie: 'Wczytywanie zestawienia z materiałem wejściowym…',
    brak: 'Dokument nie odbiega od materiału wejściowego.',
    nieodczytany: 'Nie udało się odczytać materiału wejściowego dla tego dokumentu.',
  },

  grep: {
    /* Nazwa panelu „Diff/Grep” zostaje — to nazwa własna produktu. Ale
       podpowiedź w polu i etykieta pola wyboru mówią Operatorowi, co ma zrobić,
       więc idą polszczyzną: „grep” jest nazwą narzędzia uniksowego, „regex”
       skrótem technicznym. Rozstrzygnięcie Właściciela z 29.08.2026: treść
       niepoprawna podlega poprawieniu także wtedy, gdy stoi tak w prototypie. */
    zastepczaTresc: 'Wzorzec wyszukiwania…',
    etykieta: 'Wyszukiwanie wzorca',
    regex: 'Traktuj jako wyrażenie regularne',
    naglowek: 'Trafienia wzorca',
    brakTrafien: 'Nie znaleziono dopasowań wzorca w porównywanych wersjach.',
    tylkoTekst: 'Wzorca szuka się w porównaniu scalonym; pozostałe zakresy go nie niosą.',
  },

  adnotacja: {
    etykieta: 'Adnotacja',
    zastepczaTresc: 'Treść adnotacji…',
    dodaj: 'Dodaj',
    naglowek: 'Adnotacje',
    ladowanie: 'Wczytywanie adnotacji…',
    autorOperator: 'Operator',
    autorWykonawca: 'Wykonawca',
  },

  przeniesienie: {
    etykieta: 'Przenieś fragment',
    wTrakcie: 'Przenoszenie…',
    wymaganaWersja: 'Przeniesienie fragmentu wymaga wskazania wersji odniesienia — to z niej fragment jest brany.',
    pominiete: 'Pominięto fragmentów:',
  },

  sledzenie: {
    etykieta: 'Śledzenie zmian',
    nieznane: 'Stan śledzenia zmian jest znany dopiero po jego ustawieniu.',
    wlaczone: 'Śledzenie zmian jest włączone — każdy zapis odkłada zmianę do decyzji.',
    wylaczone: 'Śledzenie zmian jest wyłączone — zapis nadpisuje treść.',
  },

  zmiany: {
    naglowek: 'Zmiany do decyzji',
    ladowanie: 'Wczytywanie zmian…',
    brak: 'W tym dokumencie nie ma zmian oczekujących na decyzję.',
    akceptuj: 'Akceptuj zmianę',
    odrzuc: 'Odrzuć',
    autorOperator: 'Operator',
    autorWykonawca: 'Wykonawca',
  },

  raport: {
    etykieta: 'Format raportu',
    wydaj: 'Eksportuj raport zmian',
    wTrakcie: 'Wydawanie raportu…',
    wymaganaWersja: 'Raport zmian wymaga wskazania wersji odniesienia.',
    gotowy: 'Raport zmian odłożony w magazynie:',
    formaty: { pdf: 'PDF', docx: 'DOCX', txt: 'Tekst zwykły', markdown: 'Markdown' },
    jednostkaBajtow: { jedna: 'bajt', kilka: 'bajty', wiele: 'bajtów' },
  },

  odmowa: {
    porownanie: 'Nie udało się porównać wersji',
    postac: 'Nie udało się porównać postaci wersji',
    wyglad: 'Nie udało się porównać wyglądu wersji',
    zrodlo: 'Nie udało się zestawić dokumentu z materiałem wejściowym',
    zmiany: 'Nie udało się wczytać wykazu zmian',
    decyzja: 'Nie udało się zapisać decyzji o zmianie',
    sledzenie: 'Nie udało się przestawić śledzenia zmian',
    adnotacja: 'Nie udało się założyć adnotacji',
    adnotacje: 'Nie udało się wczytać adnotacji',
    przeniesienie: 'Nie udało się przenieść fragmentu',
    raport: 'Nie udało się wydać raportu zmian',
  },
} as const;
