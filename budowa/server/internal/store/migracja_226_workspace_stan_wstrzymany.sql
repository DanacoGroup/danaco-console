-- Migracja 226 rozszerza dopuszczalne stany projektu o stan wstrzymany,
-- odtwarzając tabelę projekt z nowym warunkiem kolumny stan.

PRAGMA foreign_keys = off;

CREATE TABLE projekt_nowy (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT,
    stan           TEXT    NOT NULL DEFAULT 'active'
                           CHECK(stan IN ('active','paused','archived')),
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO projekt_nowy (id, kod, nazwa, opis, stan, utworzono, zaktualizowano)
SELECT id, kod, nazwa, opis, stan, utworzono, zaktualizowano FROM projekt;

DROP TABLE projekt;
ALTER TABLE projekt_nowy RENAME TO projekt;

PRAGMA foreign_keys = on;
