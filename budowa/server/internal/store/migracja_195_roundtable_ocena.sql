-- Migracja 195 zakłada tabele rubryk oceny, kryteriów, oceny Operatora oraz werdyktów modeli-sędziów debaty, rozdzielonych po autorze i skutku.

CREATE TABLE debata_rubryka (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL DEFAULT '',
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_rubryka_okno ON debata_rubryka(okno, id);

CREATE TABLE debata_kryterium (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rubryka                  TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    waga                     REAL    NOT NULL DEFAULT 0,
    opis                     TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_kryterium_rubryka ON debata_kryterium(rubryka, kolejnosc, id);

CREATE TABLE debata_ocena (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL CHECK(rodzaj IN ('star','pairwise')),
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    -- Gwiazdki poza zakresem 1..5 znaczą brak oceny w skali; zero jest wartością domyślną porównania.
    gwiazdki                 INTEGER NOT NULL DEFAULT 0,
    wskazana                 TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_ocena_okno ON debata_ocena(okno, id);

CREATE TABLE debata_werdykt (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rubryka                  TEXT    NOT NULL,
    sedzia                   TEXT    NOT NULL,
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    punkty_json              TEXT    NOT NULL DEFAULT '{}',
    wynik                    REAL    NOT NULL DEFAULT 0,
    uzasadnienie             TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_werdykt_okno ON debata_werdykt(okno, id);
