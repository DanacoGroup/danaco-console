-- Migracja 184 zakłada tabele polityk przechowywania zasobów oraz zadań utrwalenia archiwalnego modułu Library wraz z zapisem skutku walidacji.

CREATE TABLE polityka_retencji_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zasieg                   TEXT    NOT NULL,
    zasieg_id                TEXT,
    dni_przechowywania       INTEGER NOT NULL,
    czynnosc                 TEXT    NOT NULL
                             CHECK(czynnosc IN ('przeglad','archiwizacja','usuniecie')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_polityka_retencji_biblioteki_zasieg
    ON polityka_retencji_biblioteki(zasieg, zasieg_id);

CREATE TABLE zadanie_utrwalenia_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    plik_kod                 TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL CHECK(rodzaj IN ('pdfa','bagit','premis')),
    plik_wynikowy_kod        TEXT,
    profil                   TEXT,
    poprawne                 INTEGER NOT NULL DEFAULT 0 CHECK(poprawne IN (0,1)),
    -- Zapis walidacji: co sprawdzono i z jakim wynikiem, inaczej wynik byłby werdyktem bez uzasadnienia.
    raport                   TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zadanie_utrwalenia_biblioteki_plik
    ON zadanie_utrwalenia_biblioteki(plik_kod, utworzono DESC);
