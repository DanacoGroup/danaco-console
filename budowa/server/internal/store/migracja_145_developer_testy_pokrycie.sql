-- Migracja 145 zakłada tabele wyników testów i pokrycia kodu przebiegu budowania okna Build Output, wypełniane rozbiorem logu przy jego domknięciu.

CREATE TABLE developer_wynik_testu (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    budowanie_kod TEXT    NOT NULL,
    zestaw        TEXT,
    nazwa         TEXT    NOT NULL,
    stan          TEXT    NOT NULL CHECK(stan IN ('passed','failed','skipped')),
    czas_ms       INTEGER,
    tresc         TEXT,
    sciezka       TEXT,
    wiersz        INTEGER
);
CREATE INDEX idx_developer_wynik_testu_budowanie ON developer_wynik_testu(budowanie_kod, stan);

CREATE TABLE developer_pokrycie (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    budowanie_kod        TEXT    NOT NULL,
    sciezka              TEXT    NOT NULL,
    instrukcje           INTEGER NOT NULL DEFAULT 0,
    pokryte              INTEGER NOT NULL DEFAULT 0,
    procent              INTEGER NOT NULL DEFAULT 0,
    wiersze_bez_pokrycia TEXT,
    UNIQUE(budowanie_kod, sciezka)
);
CREATE INDEX idx_developer_pokrycie_budowanie ON developer_pokrycie(budowanie_kod);
