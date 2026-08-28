-- Tworzy tabelę zasada_przechowywania, ustalającą jedną zasadę retencji historii rozmowy dla zakresu okna, sesji albo zasięgu globalnego, wygrywającą od najwęższego zakresu.

CREATE TABLE zasada_przechowywania (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    zakres           TEXT    NOT NULL
                             CHECK (zakres IN ('window','session','global')),
    zakres_kod       TEXT,
    dni_trzymania    INTEGER CHECK (dni_trzymania    IS NULL OR dni_trzymania    > 0),
    pozycje_trzymane INTEGER CHECK (pozycje_trzymane IS NULL OR pozycje_trzymane > 0),
    utworzono        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- Zakres `global` ma być JEDEN, nie jeden na każdy pusty kod.
    CHECK ((zakres = 'global' AND zakres_kod IS NULL)
        OR (zakres <> 'global' AND zakres_kod IS NOT NULL AND zakres_kod <> ''))
);

-- Jedna zasada na zakres — druga byłaby drugą prawdą o tym samym oknie.
-- Wyrażenie COALESCE zamiast pary kolumn, bo NULL nie równa się NULL i zwykły
-- indeks UNIQUE przepuściłby dowolną liczbę wierszy `global`.
CREATE UNIQUE INDEX idx_zasada_przechowywania_zakres
    ON zasada_przechowywania (zakres, COALESCE(zakres_kod, ''));
