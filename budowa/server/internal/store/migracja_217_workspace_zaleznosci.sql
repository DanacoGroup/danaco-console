-- Migracja 217 zakłada tabelę zależności między zadaniami projektu, wiążącą poprzednika i następnika z warunkiem unikalności pary.

CREATE TABLE zaleznosc_zadan_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    poprzednik               TEXT    NOT NULL,
    nastepnik                TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL DEFAULT 'finishToStart'
                                     CHECK(rodzaj IN ('finishToStart','blocks')),
    odstep_minut             INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(projekt_id, poprzednik, nastepnik)
);
CREATE INDEX idx_zaleznosc_zadan_poprzednik ON zaleznosc_zadan_projektu(poprzednik);
CREATE INDEX idx_zaleznosc_zadan_nastepnik ON zaleznosc_zadan_projektu(nastepnik);
