-- Migracja 365 wprowadza znakowanie fragmentów: wyróżnienie barwą, znacznik
-- własny operatora i propozycję z brzmieniem, w jednej tabeli znakowań.

CREATE TABLE rodzaj_znacznika_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                    TEXT    NOT NULL UNIQUE,
    nazwa_widoczna           TEXT    NOT NULL,
    barwa                    TEXT,
    -- Rodzaju fabrycznego nie usuwa się — odpowiada odmowa nazywająca powód.
    fabryczny                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE znakowanie_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('highlight','mark','suggestion')),
    -- Tożsamością autora jest kod agenta; rodzaj autora zostaje grubym rozróżnieniem modelu i użytkownika.
    autor_rodzaj             TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor_rodzaj IN ('uzytkownik','model')),
    autor_agent_kod          TEXT,
    autor_agent_nazwa        TEXT,
    autor_agent_wersja       TEXT,
    autor_podagent_kod       TEXT,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    barwa                    TEXT,
    znacznik_nazwa           TEXT REFERENCES rodzaj_znacznika_studio(nazwa) ON DELETE SET NULL,
    tresc                    TEXT,
    -- Brzmienie proponowane i zastane stoją oba, by pokazać zmianę i umieć ją cofnąć.
    brzmienie_proponowane    TEXT,
    brzmienie_zastane        TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','accepted','rejected','resolved')),
    czynnosc_kod             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
-- Wykaz znakowań idzie w kolejności WYSTĄPIENIA W TREŚCI, nie zapisu — Operator
-- przechodzi dokument od początku do końca, a nie od najnowszego znakowania.
CREATE INDEX idx_znakowanie_studio_dokument
    ON znakowanie_studio(dokument_id, zakres_od, id);
CREATE INDEX idx_znakowanie_studio_filtr
    ON znakowanie_studio(dokument_id, rodzaj, autor_rodzaj, stan);
CREATE INDEX idx_znakowanie_studio_agent
    ON znakowanie_studio(dokument_id, autor_agent_kod);

-- Rodzaje fabryczne wchodzą tu, a nie w kod rdzenia, żeby wykaz był jeden
-- i żeby dało się zmienić ich barwę, gdy operator tego zażąda.
INSERT INTO rodzaj_znacznika_studio (nazwa, nazwa_widoczna, barwa, fabryczny) VALUES
    ('doSprawdzenia', 'Do sprawdzenia', '#f5a623', 1),
    ('wymagaZrodla',  'Wymaga źródła',  '#d0021b', 1),
    ('gotowe',        'Gotowe',         '#2f9e44', 1);
