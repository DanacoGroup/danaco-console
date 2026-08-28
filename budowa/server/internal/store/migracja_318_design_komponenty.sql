-- Migracja 318 dodaje komponenty makiety wraz z wariantami oraz ich
-- instancjami rozmieszczonymi na kompozycjach.

CREATE TABLE komponent_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    warianty_json            TEXT,
    zestaw_zetonow_kod       TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_komponent_design_okno ON komponent_design(okno, id);

CREATE TABLE instancja_komponentu_design (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    komponent_id  INTEGER NOT NULL REFERENCES komponent_design(id) ON DELETE CASCADE,
    kompozycja_id INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    warstwa_kod   TEXT    NOT NULL UNIQUE,
    wariant       TEXT,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_instancja_komponentu_design_komponent
    ON instancja_komponentu_design(komponent_id, id);
