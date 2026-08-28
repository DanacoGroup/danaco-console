-- Migracja 260 wprowadza wersje definicji automatyki: wiersz niesie migawkę kroków zapisaną
-- przy każdym zapisie definicji, co daje źródło do porównania wersji i przywrócenia wcześniejszej.
CREATE TABLE wersja_automatyki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    automatyka_id  INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    wersja         INTEGER NOT NULL,
    kroki          TEXT    NOT NULL,
    opublikowana   INTEGER NOT NULL DEFAULT 0 CHECK(opublikowana IN (0,1)),
    autor          TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(automatyka_id, wersja)
);
-- Indeks porządkuje wersje definicji automatyki malejąco po numerze wersji, aby odczyt
-- najnowszej wersji następował bez sortowania w czasie zapytania.
CREATE INDEX idx_wersja_automatyki_czas
    ON wersja_automatyki(automatyka_id, wersja DESC);
