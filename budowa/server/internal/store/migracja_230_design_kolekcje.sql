-- Migracja 230 zakłada tabele kolekcji zasobów modułu Design oraz przypisań
-- zasobów do kolekcji, jako byt odrębny od etykiety zasobu.

CREATE TABLE kolekcja_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Indeks wspiera odczyt kolekcji okna w porządku od najpóźniej zaktualizowanej do najstarszej pozycji wykazu.
CREATE INDEX idx_kolekcja_design_okno ON kolekcja_design(okno, zaktualizowano DESC, id DESC);

CREATE TABLE pozycja_kolekcji_design (
    kolekcja_id INTEGER NOT NULL REFERENCES kolekcja_design(id) ON DELETE CASCADE,
    zasob_id    TEXT    NOT NULL,
    kolejnosc   INTEGER NOT NULL DEFAULT 0,
    utworzono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (kolekcja_id, zasob_id)
);
-- Indeks wspiera zawężanie wykazu kolekcji okna do tych kolekcji, które zawierają wskazany zasób projektu.
CREATE INDEX idx_pozycja_kolekcji_design_zasob ON pozycja_kolekcji_design(zasob_id);
