-- Migracja 193 zakłada tabele ustaleń analizy oraz rejestru dowodów, rozdzielone, ponieważ rejestr odpowiada na pytanie, czym poparto twierdzenie.

CREATE TABLE debata_ustalenie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('argumentMining','speechAct','fallacy',
                                                      'steelman','dedup','factCheck','tone')),
    tura                     TEXT    NOT NULL DEFAULT '',
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    wezel                    TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    tresc                    TEXT    NOT NULL,
    pewnosc                  REAL    NOT NULL DEFAULT -1.0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_ustalenie_okno ON debata_ustalenie(okno, id);

CREATE TABLE debata_dowod (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tura                     TEXT    NOT NULL DEFAULT '',
    wypowiedz                TEXT    NOT NULL,
    twierdzenie              TEXT    NOT NULL,
    zrodlo                   TEXT,
    adres                    TEXT,
    poparte                  INTEGER NOT NULL DEFAULT 0 CHECK(poparte IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_dowod_okno ON debata_dowod(okno, id);
