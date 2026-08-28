-- Migracja 295 dodaje magazyn nagrań mowy przechowywanych na dysku jako
-- pliki, z odnośnikiem używanym do transkrypcji i odsłuchu.

CREATE TABLE nagranie_mowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL UNIQUE,
    typ_tresci               TEXT    NOT NULL,
    rozmiar_bajtow           INTEGER NOT NULL,
    dlugosc_ms               INTEGER,
    sesja_kod                TEXT,
    okno_kod                 TEXT,
    trwale                   INTEGER NOT NULL DEFAULT 0,
    utworzono                INTEGER NOT NULL,
    wygasa                   INTEGER
);
CREATE INDEX idx_nagranie_mowy_okno ON nagranie_mowy(okno_kod, utworzono DESC);
CREATE INDEX idx_nagranie_mowy_wygasa ON nagranie_mowy(wygasa);
