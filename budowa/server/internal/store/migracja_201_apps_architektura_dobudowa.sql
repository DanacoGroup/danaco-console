-- Migracja 201 dobudowuje historię wersji układu Architecture Designera oraz adnotacje projektowe przypinane do komponentu albo do zależności kanwy.

-- Zakłada tabelę wersja_architektury_apps niosącą liczbę komponentów i zdanie o różnicy wobec wersji poprzedniej.
CREATE TABLE wersja_architektury_apps (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    architektura_id    INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    wersja             INTEGER NOT NULL,
    liczba_komponentow INTEGER NOT NULL DEFAULT 0,
    -- Zdanie o różnicy wobec wersji poprzedniej; pierwsza wersja go nie ma.
    roznica            TEXT,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (architektura_id, wersja)
);
CREATE INDEX idx_wersja_architektury_apps ON wersja_architektury_apps(architektura_id, wersja DESC);

-- Zakłada tabelę adnotacja_architektury_apps niosącą notatkę przypiętą do komponentu, do zależności albo wolną na kanwie.
CREATE TABLE adnotacja_architektury_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    architektura_id          INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    komponent_kod            TEXT,
    zaleznosc_z              TEXT,
    zaleznosc_do             TEXT,
    tresc                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Jedno przypięcie albo żadne; zależność wymaga obu końców, bo pół krawędzi nie wskazuje niczego.
    CHECK (
        (komponent_kod IS NOT NULL AND zaleznosc_z IS NULL AND zaleznosc_do IS NULL)
        OR (komponent_kod IS NULL AND zaleznosc_z IS NOT NULL AND zaleznosc_do IS NOT NULL)
        OR (komponent_kod IS NULL AND zaleznosc_z IS NULL AND zaleznosc_do IS NULL)
    )
);
CREATE INDEX idx_adnotacja_architektury_apps ON adnotacja_architektury_apps(architektura_id, id);
