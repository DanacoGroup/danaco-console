-- Migracja 339 dodaje strony szablonu materiału modułu Design wraz z ich
-- warstwami dla publikacji wielostronicowych.

CREATE TABLE strona_szablonu_materialu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    szablon_id               INTEGER NOT NULL
                                     REFERENCES szablon_materialu_design(id) ON DELETE CASCADE,
    numer                    INTEGER NOT NULL,
    nazwa                    TEXT,
    UNIQUE(szablon_id, numer)
);

CREATE INDEX idx_strona_szablonu_materialu_design_szablon
    ON strona_szablonu_materialu_design(szablon_id, numer, id);

CREATE TABLE warstwa_strony_szablonu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    strona_id                INTEGER NOT NULL
                                     REFERENCES strona_szablonu_materialu_design(id)
                                     ON DELETE CASCADE,
    zasob_id                 TEXT,
    x                        REAL,
    y                        REAL,
    szerokosc                REAL,
    wysokosc                 REAL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    zablokowana              INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja                TEXT
);

CREATE INDEX idx_warstwa_strony_szablonu_design_strona
    ON warstwa_strony_szablonu_design(strona_id, kolejnosc, id);
