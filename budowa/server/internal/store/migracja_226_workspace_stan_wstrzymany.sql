-- Migracja 226 — stan projektu „wstrzymany".
--
-- Kontrakt zna trzy stany projektu (`active`, `paused`, `archived`), a schemat
-- z migracji 035 dopuszczał dwa. Skutek był taki, że komenda ustawiająca stan
-- kończyła się odmową o warunku kolumny — zdaniem o schemacie, a nie o pracy
-- Operatora. Opracowanie modułu wymaga stanu wstrzymanego wprost (rozdz. 6.2:
-- „projekt tymczasowo nieaktywny, dane zachowane"), więc brak był w schemacie,
-- nie w żądaniu.
--
-- SQLite nie zmienia warunku kolumny w miejscu, dlatego tabela powstaje na nowo
-- i przejmuje wiersze. Więzy obce wskazujące `projekt(id)` odtwarzają się same,
-- bo klucze główne idą przepisane bez zmiany wartości.

PRAGMA foreign_keys = off;

CREATE TABLE projekt_nowy (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT,
    stan           TEXT    NOT NULL DEFAULT 'active'
                           CHECK(stan IN ('active','paused','archived')),
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO projekt_nowy (id, kod, nazwa, opis, stan, utworzono, zaktualizowano)
SELECT id, kod, nazwa, opis, stan, utworzono, zaktualizowano FROM projekt;

DROP TABLE projekt;
ALTER TABLE projekt_nowy RENAME TO projekt;

PRAGMA foreign_keys = on;
