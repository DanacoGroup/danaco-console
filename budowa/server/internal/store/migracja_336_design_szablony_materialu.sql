-- Migracja 336 dodaje szablony materiału marketingowego wraz z ich
-- warstwami oraz indeksami wspierającymi wykaz.

CREATE TABLE szablon_materialu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('social','banner','presentation',
                                                      'print','email')),
    szerokosc                REAL    NOT NULL,
    wysokosc                 REAL    NOT NULL,
    opis                     TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE warstwa_szablonu_materialu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    szablon_id               INTEGER NOT NULL
                                     REFERENCES szablon_materialu_design(id) ON DELETE CASCADE,
    zasob_id                 TEXT,
    x                        REAL,
    y                        REAL,
    szerokosc                REAL,
    wysokosc                 REAL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    zablokowana              INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja                TEXT
);

-- Wykaz szablonów materiału marketingowego czyta się dla danego okna, od
-- ostatnio zmienianego szablonu.
CREATE INDEX idx_szablon_materialu_design_okno
    ON szablon_materialu_design(okno, zaktualizowano DESC, id DESC);

-- Warstwy szablonu materiału czyta się zawsze kompletem jednego szablonu,
-- w zachowanej kolejności wyrysu.
CREATE INDEX idx_warstwa_szablonu_materialu_design_szablon
    ON warstwa_szablonu_materialu_design(szablon_id, kolejnosc, id);
