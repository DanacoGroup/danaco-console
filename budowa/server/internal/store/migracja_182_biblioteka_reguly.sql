-- Migracja 182 zakłada tabelę reguł repozytorium modułu Library, obejmującą kolekcję inteligentną, regułę napływu i folder obserwowany z warunkiem w JSON.

CREATE TABLE regula_biblioteki (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT NOT NULL UNIQUE,
    rodzaj               TEXT    NOT NULL
                         CHECK(rodzaj IN ('kolekcja','naplyw','obserwacja')),
    nazwa                TEXT    NOT NULL,
    -- Warunek w zapisie JSON; pusty obiekt znaczy regułę bez zawężenia.
    warunek              TEXT    NOT NULL DEFAULT '{}',
    kolekcja_docelowa_kod TEXT,
    -- Ścieżka folderu obserwowanego; wypełniona wyłącznie dla rodzaju
    -- `obserwacja`.
    sciezka_obserwowana  TEXT,
    czynna               INTEGER NOT NULL DEFAULT 1 CHECK(czynna IN (0,1)),
    ostatnie_przeliczenie TEXT,
    utworzono            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_regula_biblioteki_rodzaj ON regula_biblioteki(rodzaj, nazwa);

ALTER TABLE przypisanie_kolekcji_biblioteki
    ADD COLUMN zrodlo TEXT NOT NULL DEFAULT 'reczne'
        CHECK(zrodlo IN ('reczne','regula'));
