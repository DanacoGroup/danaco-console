-- Migracja 263 — mapowanie danych między krokami
-- (`automation.workflow.variables.set`, pole `mappings`).
--
-- Kroki wskazywane są kodem kroku, nie kluczem wiersza `krok_automatyki`:
-- Workflow Builder zapisuje mapowanie równocześnie z krokami, a zapis kroków
-- podmienia wiersze w całości. Klucz obcy do wiersza kasowałby mapowanie przy
-- każdym zapisie definicji. Mapowanie do kroku nieistniejącego jest
-- zastrzeżeniem oddawanym w wyniku, nie odmową schematu.
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
