-- Migracja 316 dodaje ścieżki wektorowe modułu Design wraz z symbolami
-- wielokrotnego użycia i ich członkami.

CREATE TABLE sciezka_wektorowa_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    warstwa_kod              TEXT,
    nazwa                    TEXT,
    wezly_json               TEXT    NOT NULL,
    zamknieta                INTEGER NOT NULL DEFAULT 0,
    wypelnienie_json         TEXT,
    obrys_json               TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sciezka_wektorowa_design_kompozycja
    ON sciezka_wektorowa_design(kompozycja_id, kolejnosc, id);
CREATE INDEX idx_sciezka_wektorowa_design_warstwa
    ON sciezka_wektorowa_design(warstwa_kod);

-- Symbol jest definicją wielokrotnego użycia złożoną ze ścieżek i warstw,
-- a jego członkostwo leży w osobnej tabeli.
CREATE TABLE symbol_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_symbol_design_kompozycja ON symbol_design(kompozycja_id, id);

CREATE TABLE czlonek_symbolu_design (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol_id   INTEGER NOT NULL REFERENCES symbol_design(id) ON DELETE CASCADE,
    rodzaj      TEXT    NOT NULL CHECK(rodzaj IN ('sciezka','warstwa')),
    czlonek_kod TEXT    NOT NULL,
    kolejnosc   INTEGER NOT NULL DEFAULT 0,
    UNIQUE(symbol_id, rodzaj, czlonek_kod)
);
CREATE INDEX idx_czlonek_symbolu_design_symbol ON czlonek_symbolu_design(symbol_id, kolejnosc, id);
