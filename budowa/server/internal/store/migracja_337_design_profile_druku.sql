-- Migracja 337 dodaje profile wydania do druku z ustawieniami spadu,
-- znaczników, przestrzeni barw i nośnika.

CREATE TABLE profil_druku_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    przestrzen_barw          TEXT    NOT NULL
                                     CHECK(przestrzen_barw IN ('cmyk','rgb','grayscale','spot')),
    norma                    TEXT    CHECK(norma IS NULL OR
                                           norma IN ('pdfX1a','pdfX3','pdfX4','pdfA2b','brak')),
    spad_mm                  REAL,
    znaczniki_ciecia         INTEGER NOT NULL DEFAULT 0 CHECK(znaczniki_ciecia IN (0,1)),
    znaczniki_pasowania      INTEGER NOT NULL DEFAULT 0 CHECK(znaczniki_pasowania IN (0,1)),
    pasek_barw               INTEGER NOT NULL DEFAULT 0 CHECK(pasek_barw IN (0,1)),
    rozdzielczosc            INTEGER,
    profil_icc               TEXT,
    nadruk_czerni            INTEGER NOT NULL DEFAULT 0 CHECK(nadruk_czerni IN (0,1)),
    nosnik                   TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Wykaz profili wydania do druku czyta się dla danego okna, od ostatnio
-- zmienianego profilu wydania do druku.
CREATE INDEX idx_profil_druku_design_okno
    ON profil_druku_design(okno, zaktualizowano DESC, id DESC);
