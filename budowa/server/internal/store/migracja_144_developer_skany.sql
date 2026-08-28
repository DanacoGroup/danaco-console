-- Migracja 144 — przebiegi skanowania i ich znaleziska (Developer, zakładka
-- Bezpieczeństwo w Dev Tools).
--
-- Skan i znalezisko rozdzielono na dwie tabele, bo pytanie o nie zadaje się
-- osobno: `developer.scan.run` zakłada przebieg i oddaje jego nagłówek,
-- a `developer.scan.result.list` czyta znaleziska z filtrem po rodzaju i wadze,
-- często dla kilku przebiegów naraz. Jedna tabela kazałaby powtarzać nagłówek
-- przy każdym znalezisku.
--
-- Rodzaje skanu przebiegu zapisujemy jako tekst rozdzielony przecinkiem.
-- Rodzajów jest cztery i są zamkniętym słownikiem kontraktu (ScanKind); tabela
-- pośrednia na cztery wartości byłaby złożonością bez odbiorcy.
--
-- Znalezisko nie ma stanu „przyjęte/odrzucone”. Skan jest pomiarem stanu
-- repozytorium w danej chwili, a nie listą zadań: kolejny przebieg zakłada nowe
-- znaleziska, a poprzednie zostają śladem tamtego pomiaru.
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
