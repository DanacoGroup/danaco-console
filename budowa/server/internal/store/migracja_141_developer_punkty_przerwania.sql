-- Migracja 141 zakłada tabelę punktów przerwania okna Run and Debug, wiążąc każdy punkt z oknem i plikiem oraz rozdzielając warunek, warunek trafień i treść wpisu.

CREATE TABLE developer_punkt_przerwania (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    okno_kod      TEXT    NOT NULL,
    sciezka       TEXT    NOT NULL,
    wiersz        INTEGER NOT NULL,
    rodzaj        TEXT    NOT NULL DEFAULT 'line'
                          CHECK(rodzaj IN ('line','conditional','logpoint','exception')),
    warunek       TEXT,
    warunek_trafien TEXT,
    wpis          TEXT,
    zweryfikowany INTEGER NOT NULL DEFAULT 0,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno_kod, sciezka, wiersz)
);
CREATE INDEX idx_developer_punkt_okno ON developer_punkt_przerwania(okno_kod, sciezka, wiersz);
