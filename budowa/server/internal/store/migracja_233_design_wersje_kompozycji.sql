-- Migracja 233 zakłada tabele nazwanych wersji kompozycji Design Board wraz
-- z migawką warstw, przechowywanych niezależnie od bieżącego układu.

CREATE TABLE wersja_kompozycji_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    nazwa                    TEXT,
    uzasadnienie             TEXT,
    liczba_warstw            INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks wspiera odczyt nazwanych wersji kompozycji w porządku od najświeższej do najstarszej pozycji.
CREATE INDEX idx_wersja_kompozycji_design_kompozycja
    ON wersja_kompozycji_design(kompozycja_id, utworzono DESC, id DESC);

CREATE TABLE warstwa_wersji_kompozycji_design (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    wersja_id             INTEGER NOT NULL REFERENCES wersja_kompozycji_design(id) ON DELETE CASCADE,
    identyfikator_warstwy TEXT    NOT NULL,
    zasob_id              TEXT,
    x                     REAL,
    y                     REAL,
    szerokosc             REAL,
    wysokosc              REAL,
    kolejnosc             INTEGER NOT NULL DEFAULT 0,
    zablokowana           INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja             TEXT
);
CREATE INDEX idx_warstwa_wersji_kompozycji_design_wersja
    ON warstwa_wersji_kompozycji_design(wersja_id, kolejnosc, id);
