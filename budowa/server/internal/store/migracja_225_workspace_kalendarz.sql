-- Migracja 225 zakłada tabelę pozycji kalendarza projektu dla wydarzeń
-- wciągniętych z pliku iCal, odrębną od zadań z terminem własnym.

CREATE TABLE pozycja_kalendarza_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    tytul                    TEXT    NOT NULL,
    poczatek_ms              INTEGER NOT NULL DEFAULT 0,
    koniec_ms                INTEGER NOT NULL DEFAULT 0,
    caly_dzien               INTEGER NOT NULL DEFAULT 0,
    kamien_milowy            INTEGER NOT NULL DEFAULT 0,
    regula_powtarzalnosci    TEXT    NOT NULL DEFAULT '',
    -- Identyfikator wydarzenia w źródle; powtórne wciągnięcie aktualizuje wiersz.
    uid_zewnetrzny           TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(projekt_id, uid_zewnetrzny)
);
CREATE INDEX idx_pozycja_kalendarza_projekt ON pozycja_kalendarza_projektu(projekt_id, poczatek_ms);
