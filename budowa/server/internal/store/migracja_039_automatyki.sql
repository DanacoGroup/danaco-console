-- Migracja zakłada tabele definicji automatyki, jej kroków i zależności między krokami,
-- jako trwałość modułu automatyzacji.

-- Automatyka jest komponentem własnym, widocznym niezależnie od karty sesji, w której
-- Operator akurat pracuje.
CREATE TABLE automatyka (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    czynna                   INTEGER NOT NULL DEFAULT 1 CHECK(czynna IN (0,1)),
    -- Wersja rośnie przy każdym zapisie, by dało się powiedzieć, którą wersję Operator ogląda.
    wersja                   INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_automatyka_nazwa ON automatyka(nazwa, id);

-- Wartości kolumny rodzaju są wartościami kontraktu wprost, więc przekład wiersza na
-- krok kontraktu nie potrzebuje słownika.
CREATE TABLE krok_automatyki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id            INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    nazwa                    TEXT,
    rodzaj                   TEXT    NOT NULL DEFAULT 'command'
                                     CHECK(rodzaj IN ('command','model','condition','branch','wait')),
    komenda                  TEXT,
    parametry                TEXT,
    warunek                  TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (automatyka_id, identyfikator_zewnetrzny)
);
CREATE INDEX idx_krok_automatyki_kolejnosc ON krok_automatyki(automatyka_id, kolejnosc, id);

-- Więz pierwotny obejmuje parę kroków, więc ten sam łuk nie powstanie dwa razy; pętli
-- własnej schemat nie dopuszcza wprost.
CREATE TABLE zaleznosc_kroku_automatyki (
    automatyka_id  INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok_z         TEXT    NOT NULL,
    krok_do        TEXT    NOT NULL,
    rodzaj         TEXT    NOT NULL DEFAULT 'sequential'
                           CHECK(rodzaj IN ('sequential','parallel','conditional')),
    warunek        TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok_z, krok_do),
    CHECK (krok_z <> krok_do)
);
CREATE INDEX idx_zaleznosc_kroku_do ON zaleznosc_kroku_automatyki(automatyka_id, krok_do);
