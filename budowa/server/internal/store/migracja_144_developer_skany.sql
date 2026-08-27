-- Migracja 144 zakłada dwie tabele: przebiegi skanowania bezpieczeństwa okna Dev Tools oraz ich znaleziska, rozdzielone ze względu na osobne zapytania.

CREATE TABLE developer_skan (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    okno_kod    TEXT    NOT NULL,
    rodzaje     TEXT    NOT NULL,
    stan        TEXT    NOT NULL DEFAULT 'running'
                        CHECK(stan IN ('running','succeeded','failed','stopped')),
    znalezisk   INTEGER,
    uruchomiono TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono  TEXT
);
CREATE INDEX idx_developer_skan_okno ON developer_skan(okno_kod, uruchomiono DESC);

CREATE TABLE developer_znalezisko (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    skan_kod      TEXT    NOT NULL REFERENCES developer_skan(kod) ON DELETE CASCADE,
    rodzaj        TEXT    NOT NULL CHECK(rodzaj IN ('dependencies','secrets','code','licenses')),
    waga          TEXT    NOT NULL CHECK(waga IN ('error','warning','info')),
    tytul         TEXT    NOT NULL,
    opis          TEXT,
    sciezka       TEXT,
    wiersz        INTEGER,
    regula        TEXT,
    cve           TEXT,
    pakiet        TEXT,
    wersja_naprawy TEXT
);
CREATE INDEX idx_developer_znalezisko_skan ON developer_znalezisko(skan_kod, waga, rodzaj);
