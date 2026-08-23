-- Migracja 220 — tablica wizualna projektu (płótno, mapa myśli).
--
-- Scena leży jednym zapisem JSON i rdzeń jej nie rozbiera. Kształt sceny —
-- kartki, strzałki, grupy, osadzenia — należy do widoku, który ją rysuje;
-- rozbiór na wiersze związałby schemat bazy z rysunkiem interfejsu i każda
-- zmiana kształtu kartki byłaby migracją.

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
