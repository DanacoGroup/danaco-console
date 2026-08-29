/**
 * Panel Diff/Grep — katalog treści. Jedyne miejsce w tym pliku panelu z tekstem
 * widocznym dla Operatora; poza nim żaden łańcuch w `diff.ts` nie czyta się jak
 * zdanie. Ten sam układ co `moduly/studio/tresci.ts`, osobny plik zgodnie
 * z umową panelu — treść własna tego panelu nie wchodzi do katalogu wspólnego.
 */

export const tresc = {
  brakDokumentu: 'Brak dokumentu w tym oknie.',

  porownanie: {
    etykietaBazowa: 'Wersja odniesienia',
    etykietaCel: 'Wersja porównywana',
    zastepczaWersja: 'identyfikator wersji…',
    porownaj: 'Porównaj',
    ladowanie: 'Wczytywanie porównania…',
    brakRoznic: 'Wskazane wersje nie różnią się.',
    /* Kształt zapisu z prototypu: „Statystyka: +12 … · −48 … · 3 …”, ze znakiem
       kierunku przed liczbą. Jednostką są fragmenty różnicy, bo tyle podaje
       serwer — prototyp liczył słowa, których rdzeń nie zwraca. */
    statystykaEtykieta: 'Statystyka:',
    znakDodania: '+',
    znakUsuniecia: '−',
    jednostkaFragmenty: { jedna: 'fragment', kilka: 'fragmenty', wiele: 'fragmentów' },
    jednostkaZmian: { jedna: 'zmieniony fragment', kilka: 'zmienione fragmenty', wiele: 'zmienionych fragmentów' },
  },

  grep: {
    /* Nazwa panelu „Diff/Grep” zostaje — to nazwa własna produktu. Ale
       podpowiedź w polu i etykieta pola wyboru mówią Operatorowi, co ma zrobić,
       więc idą polszczyzną: „grep” jest nazwą narzędzia uniksowego, „regex”
       skrótem technicznym. Rozstrzygnięcie Właściciela z 29.08.2026: treść
       niepoprawna podlega poprawieniu także wtedy, gdy stoi tak w prototypie. */
    zastepczaTresc: 'Wzorzec wyszukiwania…',
    etykieta: 'Grep',
    regex: 'Traktuj jako wyrażenie regularne',
    naglowek: 'Trafienia wzorca',
    brakTrafien: 'Nie znaleziono dopasowań wzorca w porównywanych wersjach.',
  },

  adnotacja: {
    etykieta: 'Adnotacja',
    zastepczaTresc: 'Treść adnotacji…',
    dodaj: 'Dodaj',
  },

  zmiany: {
    naglowek: 'Zmiany do decyzji',
    ladowanie: 'Wczytywanie zmian…',
    brak: 'W tym dokumencie nie ma zmian oczekujących na decyzję.',
    akceptuj: 'Akceptuj zmianę',
    odrzuc: 'Odrzuć',
  },

  odmowa: {
    porownanie: 'Nie udało się porównać wersji',
    zmiany: 'Nie udało się wczytać wykazu zmian',
    decyzja: 'Nie udało się zapisać decyzji o zmianie',
    adnotacja: 'Nie udało się założyć adnotacji',
  },
} as const;
