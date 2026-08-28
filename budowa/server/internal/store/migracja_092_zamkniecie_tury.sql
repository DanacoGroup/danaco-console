-- Migracja zakłada zamknięcie tury jako zdarzenie wykonawcze, zapisujące, jakie zdarzenie
-- zamknęło turę i jaki stan wyprowadzono.

CREATE TABLE zamkniecie_tury (
    id                   INTEGER PRIMARY KEY,
    chwila               INTEGER NOT NULL,              -- epoka w milisekundach
    okno_kod             TEXT    NOT NULL,
    wiadomosc_kod        TEXT    NOT NULL UNIQUE,       -- jedna tura = jedno zamknięcie
    podtyp               TEXT    NOT NULL DEFAULT '',   -- dosłowne `subtype` zdarzenia result
    blad                 INTEGER NOT NULL DEFAULT 0
                         CHECK (blad IN (0, 1)),        -- dosłowne `is_error`
    -- Stan wiadomości wyprowadzony ze zdarzenia, zapisany dla odtwarzalności bez odgadywania reguły.
    stan_nadany          TEXT    NOT NULL
                         CHECK (stan_nadany IN ('zakonczona','zatrzymana','bledna')),
    koszt_usd            REAL    NOT NULL DEFAULT 0,
    tury                 INTEGER NOT NULL DEFAULT 0,
    czas_ms              INTEGER NOT NULL DEFAULT 0,
    konto                TEXT    NOT NULL DEFAULT '',   -- kod konta, nigdy sekret
    sesja_cli            TEXT    NOT NULL DEFAULT '',   -- identyfikator rozmowy CLI
    -- Wykaz typów linii przechwyconych w turze oraz liczba linii nieczytelnych.
    typy_zdarzen         TEXT    NOT NULL DEFAULT '',
    linie_nierozpoznane  INTEGER NOT NULL DEFAULT 0
);

-- Odczyt idzie po oknie jako ślad rozmowy oraz po czasie jako ślad dnia pracy, dwoma
-- osobnymi indeksami odczytu.
CREATE INDEX idx_zamkniecie_okno   ON zamkniecie_tury (okno_kod, chwila);
CREATE INDEX idx_zamkniecie_chwila ON zamkniecie_tury (chwila);
