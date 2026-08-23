-- Migracja 185 — moduł Library: udostępnienia odnośnikiem i nasłuchy zewnętrzne.
--
-- Token udostępnienia leży w kolumnie jawnie, nie jako skrót. To jest wybór
-- zgodny z zasadą jawności kluczy platformy: Operator ma móc odczytać wystawiony
-- odnośnik i przekazać go powtórnie, a nie wystawiać nowy, bo pierwszy da się
-- wyłącznie sprawdzić. Zawężeniem dostępu jest termin i odwołanie, nie
-- nieodczytywalność.
--
-- Odwołanie udostępnienia jest znacznikiem czasu, nie skasowaniem wiersza:
-- „odnośnik przestaje działać, wpis zostaje w dzienniku audytu" — wiersz
-- odwołany świadczy, że odnośnik istniał.
--
-- Nasłuch i jego zdarzenia stoją w dwóch tabelach, bo zdarzenie jest wartością
-- z wykazu, a nie tekstem: warunek CHECK po stronie wiersza zdarzenia wychwytuje
-- literówkę przy zapisie, czego lista sklejona w jednej kolumnie zrobić nie może.

CREATE TABLE udostepnienie_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT NOT NULL UNIQUE,
    zasieg                   TEXT NOT NULL CHECK(zasieg IN ('plik','kolekcja')),
    cel_kod                  TEXT NOT NULL,
    token                    TEXT NOT NULL UNIQUE,
    wygasa                   TEXT,
    odwolano                 TEXT,
    utworzono                TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_udostepnienie_biblioteki_cel ON udostepnienie_biblioteki(cel_kod, utworzono DESC);

CREATE TABLE webhook_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    adres                    TEXT    NOT NULL,
    sekret                   TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    ostatnie_zgloszenie      TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE zdarzenie_webhooka_biblioteki (
    webhook_id INTEGER NOT NULL REFERENCES webhook_biblioteki(id) ON DELETE CASCADE,
    zdarzenie  TEXT    NOT NULL
               CHECK(zdarzenie IN ('plik_dodany','plik_zmieniony','plik_zarchiwizowany',
                                   'plik_przywrocony','wersja_dolozona','regula_zadzialala')),
    PRIMARY KEY (webhook_id, zdarzenie)
);
