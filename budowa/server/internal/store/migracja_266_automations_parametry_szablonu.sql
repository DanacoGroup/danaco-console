-- Migracja 266 wprowadza parametry szablonu przepływu uzupełniane przy zastosowaniu: parametr
-- jest wierszem, a nie polem zapisu strukturalnego szablonu, aby formularz wpięcia budował się
-- z niego pole po polu.
CREATE TABLE parametr_szablonu_automatyki (
    szablon_id       INTEGER NOT NULL REFERENCES szablon_automatyki(id) ON DELETE CASCADE,
    nazwa            TEXT    NOT NULL,
    etykieta         TEXT,
    wymagany         INTEGER NOT NULL DEFAULT 0 CHECK(wymagany IN (0,1)),
    wartosc_domyslna TEXT,
    kolejnosc        INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (szablon_id, nazwa)
);
