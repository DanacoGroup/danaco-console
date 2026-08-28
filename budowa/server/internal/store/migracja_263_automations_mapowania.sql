-- Migracja 263 wprowadza mapowanie danych między krokami automatyki: krok jest wskazywany
-- własnym kodem, nie kluczem obcym, ponieważ zapis kroków podmienia wiersze w całości.
CREATE TABLE mapowanie_danych_automatyki (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok_z        TEXT    NOT NULL,
    sciezka_z     TEXT    NOT NULL,
    krok_do       TEXT    NOT NULL,
    pole_do       TEXT    NOT NULL,
    szablon       TEXT,
    kolejnosc     INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_mapowanie_danych_automatyki_wlasciciel
    ON mapowanie_danych_automatyki(automatyka_id, kolejnosc, id);
