-- Migracja 347 dodaje nastawy warsztatu fotografii modułu Design jako
-- zapisany zestaw czynności stosowany do zasobów okna.

CREATE TABLE nastawa_fotografii_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    czynnosci_json           TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno, nazwa)
);

CREATE INDEX idx_nastawa_fotografii_design_okno
    ON nastawa_fotografii_design(okno, nazwa);
