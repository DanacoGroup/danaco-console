-- Migracja 143 zakłada tabelę połączeń bazodanowych okna Data Console, niosącą odwołanie do poświadczenia w sejfie oraz nastawę tylko do odczytu.

CREATE TABLE developer_polaczenie_danych (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    okno_kod      TEXT    NOT NULL,
    nazwa         TEXT    NOT NULL,
    silnik        TEXT    NOT NULL CHECK(silnik IN ('postgres','mysql','sqlite')),
    host          TEXT,
    port          INTEGER,
    baza          TEXT    NOT NULL,
    uzytkownik    TEXT,
    poswiadczenie TEXT,
    tylko_odczyt  INTEGER NOT NULL DEFAULT 0,
    zmieniono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_developer_polaczenie_okno ON developer_polaczenie_danych(okno_kod, nazwa);
