-- Migracja 220 zakłada tabelę tablicy wizualnej projektu, przechowującą scenę płótna jako jeden zapis JSON nierozbierany przez rdzeń.

CREATE TABLE tablica_wizualna_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL DEFAULT 'Tablica',
    scena                    TEXT    NOT NULL DEFAULT '{}',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_tablica_wizualna_projekt ON tablica_wizualna_projektu(projekt_id, id);
