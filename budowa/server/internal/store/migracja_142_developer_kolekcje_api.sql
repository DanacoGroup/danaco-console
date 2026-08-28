-- Migracja 142 zakłada tabelę kolekcji zapytań okna API Client, przechowującą zapytania i środowiska jako tekst JSON wraz z importem OpenAPI.

CREATE TABLE developer_kolekcja_api (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    kod        TEXT    NOT NULL UNIQUE,
    okno_kod   TEXT    NOT NULL,
    nazwa      TEXT    NOT NULL,
    zapytania  TEXT    NOT NULL DEFAULT '[]',
    srodowiska TEXT,
    zmieniono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_developer_kolekcja_okno ON developer_kolekcja_api(okno_kod, zmieniono DESC);
