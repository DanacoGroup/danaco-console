-- Migracja zakłada trwałość modułu deweloperskiego: wersję pliku zakładaną przy zapisie
-- oraz dziennik przebiegów budowania.

-- Kolumna kodu okna niesie identyfikator okna komunikacji w postaci kontraktu, bez klucza
-- obcego, bo okno bywa bytem samej pamięci nadzorcy.
CREATE TABLE developer_wersja_pliku (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    kod        TEXT    NOT NULL UNIQUE,
    okno_kod   TEXT    NOT NULL,
    sciezka    TEXT    NOT NULL,
    tresc      TEXT    NOT NULL,
    rozmiar    INTEGER NOT NULL DEFAULT 0,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_developer_wersja_sciezka ON developer_wersja_pliku(okno_kod, sciezka, utworzono DESC);

-- Kolumna argumentów przechowuje jeden wiersz rozdzielony znakiem nowej linii, nie
-- tabelą podrzędną bez odbiorcy.
CREATE TABLE developer_budowanie (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    okno_kod    TEXT    NOT NULL,
    zadanie     TEXT    NOT NULL,
    argumenty   TEXT,
    stan        TEXT    NOT NULL DEFAULT 'running'
                        CHECK(stan IN ('running','succeeded','failed','stopped')),
    kod_wyjscia INTEGER,
    log         TEXT,
    uruchomiono TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono  TEXT
);
CREATE INDEX idx_developer_budowanie_okno ON developer_budowanie(okno_kod, uruchomiono DESC);
CREATE INDEX idx_developer_budowanie_stan ON developer_budowanie(stan, uruchomiono DESC);
