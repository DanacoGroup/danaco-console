/**
 * Panel Diff/Grep — katalog treści. Jedyne miejsce w tym pliku panelu z tekstem
 * widocznym dla Operatora; poza nim żaden łańcuch w `diff.ts` nie czyta się jak
 * zdanie. Ten sam układ co `moduly/studio/tresci.ts`, osobny plik zgodnie
 * z umową panelu — treść własna tego panelu nie wchodzi do katalogu wspólnego.
 */

export const tresc = {
  brakDokumentu: 'Rdzeń nie zgłosił jeszcze dokumentu w tym oknie.',

  porownanie: {
    etykietaBazowa: 'Wersja odniesienia',
    etykietaCel: 'Wersja porównywana',
    zastepczaWersja: 'identyfikator wersji…',
    porownaj: 'Porównaj',
    ladowanie: 'Wczytywanie porównania…',
    brakRoznic: 'Rdzeń nie zwrócił żadnej różnicy dla wskazanych wersji.',
    dodanych: 'dodanych',
    usunietych: 'usuniętych',
    zmienionych: 'zmienionych',
  },

  grep: {
    zastepczaTresc: 'Grep — wzorzec (regex)…',
    etykieta: 'Grep',
    regex: 'Wyrażenie regularne',
    naglowek: 'Trafienia wzorca',
    brakTrafien: 'Rdzeń nie znalazł dopasowań wzorca w porównywanych wersjach.',
  },

  adnotacja: {
    etykieta: 'Adnotacja',
    zastepczaTresc: 'Treść adnotacji…',
    dodaj: 'Dodaj',
  },

  zmiany: {
    naglowek: 'Zmiany do decyzji',
    ladowanie: 'Wczytywanie zmian…',
    brak: 'Rdzeń nie zgłosił żadnej zmiany oczekującej na decyzję w tym dokumencie.',
    akceptuj: 'Akceptuj zmianę',
    odrzuc: 'Odrzuć',
  },

  odmowa: {
    porownanie: 'Rdzeń odmówił porównania wersji',
    zmiany: 'Rdzeń odmówił wykazu zmian',
    decyzja: 'Rdzeń odmówił decyzji o zmianie',
    adnotacja: 'Rdzeń odmówił założenia adnotacji',
    brakOpisu: 'Rdzeń nie podał powodu odmowy.',
  },
} as const;
