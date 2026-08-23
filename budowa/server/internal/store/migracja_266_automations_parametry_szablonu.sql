-- Migracja 266 — parametry szablonu przepływu uzupełniane przy zastosowaniu.
--
-- Parametr jest wierszem, nie polem zapisu strukturalnego szablonu: formularz
-- wpięcia buduje się z niego pole po polu, a zastosowanie szablonu musi umieć
-- nazwać parametr, dla którego nie podano wartości (`missingParameters`).
CREATE TABLE parametr_szablonu_automatyki (
    szablon_id       INTEGER NOT NULL REFERENCES szablon_automatyki(id) ON DELETE CASCADE,
    nazwa            TEXT    NOT NULL,
    etykieta         TEXT,
    wymagany         INTEGER NOT NULL DEFAULT 0 CHECK(wymagany IN (0,1)),
    wartosc_domyslna TEXT,
    kolejnosc        INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (szablon_id, nazwa)
);
