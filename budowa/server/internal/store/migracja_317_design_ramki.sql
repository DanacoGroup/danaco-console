-- Migracja 317 dodaje ramki makiety, więzy responsywne warstw oraz siatki
-- układu modułu Design wraz z indeksami.

CREATE TABLE ramka_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    szerokosc                REAL    NOT NULL,
    wysokosc                 REAL    NOT NULL,
    x                        REAL,
    y                        REAL,
    nastawa_urzadzenia       TEXT,
    siatka_json              TEXT,
    uklad_json               TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_ramka_design_kompozycja ON ramka_design(kompozycja_id, id);

CREATE TABLE warstwa_ramki_design (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    ramka_id    INTEGER NOT NULL REFERENCES ramka_design(id) ON DELETE CASCADE,
    warstwa_kod TEXT    NOT NULL UNIQUE,
    kolejnosc   INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_warstwa_ramki_design_ramka ON warstwa_ramki_design(ramka_id, kolejnosc, id);

-- Więz responsywny mówi, co warstwa robi przy zmianie rozmiaru ramki; kotwica
-- pozioma i pionowa idą jednym wierszem.
CREATE TABLE wiez_ramki_design (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    ramka_id    INTEGER NOT NULL REFERENCES ramka_design(id) ON DELETE CASCADE,
    warstwa_kod TEXT    NOT NULL,
    poziomo     TEXT    NOT NULL,
    pionowo     TEXT    NOT NULL,
    UNIQUE(ramka_id, warstwa_kod)
);
CREATE INDEX idx_wiez_ramki_design_ramka ON wiez_ramki_design(ramka_id, id);

-- Siatka kompozycji obowiązuje całe płótno projektu, w odróżnieniu od siatki
-- ramki leżącej przy danym ekranie.
CREATE TABLE siatka_kompozycji_design (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kompozycja_id  INTEGER NOT NULL UNIQUE REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    siatka_json    TEXT    NOT NULL,
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
