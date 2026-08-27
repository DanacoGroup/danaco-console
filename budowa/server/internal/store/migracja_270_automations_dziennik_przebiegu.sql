-- Migracja 270 wprowadza dziennik przebiegu automatyki trwały w bazie, aby przeglądarka logów
-- pokazywała pełny zapis zdarzeń pojedynczego uruchomienia, także sprzed otwarcia okna.
CREATE TABLE log_przebiegu (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    przebieg_id INTEGER NOT NULL REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    krok_kod    TEXT,
    poziom      TEXT    NOT NULL DEFAULT 'informacja'
                        CHECK(poziom IN ('informacja','ostrzezenie','blad')),
    tresc       TEXT    NOT NULL,
    chwila      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- automation.execution.log czyta wiersze przebiegu w kolejności czasu,
-- zawężając po poziomie i po kroku.
CREATE INDEX idx_log_przebiegu_czas
    ON log_przebiegu(przebieg_id, id);
