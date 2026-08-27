-- Migracja 338 dodaje tabelę licencji zasobów wciągniętych z katalogów
-- zewnętrznych, powiązaną kluczem obcym z zasobem.

CREATE TABLE licencja_zasobu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    zasob_id                 INTEGER NOT NULL UNIQUE
                                     REFERENCES zasob_design(id) ON DELETE CASCADE,
    dostawca                 TEXT    NOT NULL,
    identyfikator_u_dostawcy TEXT    NOT NULL,
    licencja                 TEXT,
    autor                    TEXT,
    odsylacz                 TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Odczyt idzie po zasobie; wyszukanie po dostawcy służy sprawdzeniu, czy ten
-- sam materiał nie został wciągnięty drugi raz.
CREATE INDEX idx_licencja_zasobu_design_dostawca
    ON licencja_zasobu_design(dostawca, identyfikator_u_dostawcy);
