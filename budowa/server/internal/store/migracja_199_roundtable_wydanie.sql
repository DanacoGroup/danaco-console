-- Migracja 199 zakłada tabele szablonów moderacji debaty oraz artefaktów wydanych z debaty, niosących odwołanie do bajtów w magazynie treści.

CREATE TABLE debata_szablon (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    format                   TEXT    NOT NULL DEFAULT 'free'
                                     CHECK(format IN ('free','structured','oxford',
                                                      'roundRobin','delphi','expertPanel')),
    granica_tur              INTEGER NOT NULL DEFAULT 0,
    granica_czasu_ms         INTEGER NOT NULL DEFAULT 0,
    kolejnosc_glosu          TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_szablon_nazwa ON debata_szablon(nazwa);

CREATE TABLE debata_artefakt (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Rodzaj mówi, co wydano: zapis debaty, graf argumentów albo odsłuch.
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('transkrypt','graf','odsluch')),
    format                   TEXT    NOT NULL,
    -- Odwołanie względne magazynu treści, zawsze z ukośnikiem `/`.
    odwolanie                TEXT    NOT NULL,
    rozmiar                  INTEGER NOT NULL DEFAULT 0,
    suma_kontrolna           TEXT    NOT NULL DEFAULT '',
    -- Długość nagrania w milisekundach; zero przy artefaktach niebędących
    -- odsłuchem.
    dlugosc_ms               INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_artefakt_okno ON debata_artefakt(okno, id DESC);
