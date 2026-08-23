-- Migracja 225 — pozycje kalendarza projektu wciągnięte z iCal.
--
-- Kalendarz projektu składa się z dwóch źródeł: zadań z terminem (te leżą
-- w `zadanie_projektu` i drugiego wiersza nie potrzebują) oraz wydarzeń
-- wciągniętych plikiem `.ics` bez zakładania zadań. Ta tabela trzyma to drugie
-- źródło — bez niej wciągnięcie z `asTasks: false` nie miałoby gdzie osiąść
-- i kalendarz wracałby pusty mimo udanego wciągnięcia.
--
-- `uid_zewnetrzny` jest identyfikatorem wydarzenia w źródle. Warunek UNIQUE
-- na parze czyni powtórne wciągnięcie tego samego pliku aktualizacją, a nie
-- podwojeniem kalendarza.

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
    uid_zewnetrzny           TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(projekt_id, uid_zewnetrzny)
);
CREATE INDEX idx_pozycja_kalendarza_projekt ON pozycja_kalendarza_projektu(projekt_id, poczatek_ms);
