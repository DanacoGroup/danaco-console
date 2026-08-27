-- Migracja 197 zakłada tabele macierzy decyzyjnej wariantów i kryteriów z wagami, gdzie ocena wariantu w kryteriach idzie jednym polem JSON.

CREATE TABLE debata_macierz (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_macierz_okno ON debata_macierz(okno, id DESC);

CREATE TABLE debata_macierz_kryterium (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    macierz                  TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    waga                     REAL    NOT NULL DEFAULT 0,
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_macierz_kryterium ON debata_macierz_kryterium(macierz, kolejnosc, id);

CREATE TABLE debata_macierz_wariant (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    macierz                  TEXT    NOT NULL,
    etykieta                 TEXT    NOT NULL,
    oceny_json               TEXT    NOT NULL DEFAULT '{}',
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_macierz_wariant ON debata_macierz_wariant(macierz, kolejnosc, id);
