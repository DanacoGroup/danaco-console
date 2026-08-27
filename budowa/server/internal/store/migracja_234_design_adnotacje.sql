-- Migracja 234 zakłada tabelę adnotacji kompozycji Design Board wraz z wątkami,
-- przechowywanych niezależnie od zapisu układu warstw.

CREATE TABLE adnotacja_kompozycji_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    warstwa_id               TEXT,
    nadrzedna_id             TEXT,
    autor                    TEXT,
    tresc                    TEXT    NOT NULL,
    zamknieta                INTEGER NOT NULL DEFAULT 0 CHECK(zamknieta IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks wspiera odczyt adnotacji kompozycji w kolejności powstawania, zawężony do wątków niezamkniętych.
CREATE INDEX idx_adnotacja_kompozycji_design_kompozycja
    ON adnotacja_kompozycji_design(kompozycja_id, zamknieta, id);
