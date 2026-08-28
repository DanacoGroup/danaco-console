-- Migracja 279 dodaje tabele bramki dołączenia, grupy kroków oraz kroku
-- wycofującego dla układu zależności automatyki.

CREATE TABLE orkiestracja_bramka (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok          TEXT    NOT NULL,
    regula        TEXT    NOT NULL DEFAULT 'all'
                          CHECK (regula IN ('all', 'any', 'count')),
    -- Liczba torów wymagana przy regule licznikowej; przy pozostałych zero.
    licznik       INTEGER NOT NULL DEFAULT 0 CHECK (licznik >= 0),
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok),
    -- Reguła licznikowa bez liczby torów otwierałaby się zawsze albo nigdy.
    CHECK (regula <> 'count' OR licznik > 0)
);

CREATE TABLE orkiestracja_grupa (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    nazwa         TEXT    NOT NULL,
    rodzaj        TEXT    NOT NULL DEFAULT 'parallel'
                          CHECK (rodzaj IN ('sequential', 'parallel', 'conditional')),
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX idx_orkiestracja_grupa_automatyka
    ON orkiestracja_grupa(automatyka_id, nazwa);

CREATE TABLE orkiestracja_grupa_krok (
    grupa_id  INTEGER NOT NULL REFERENCES orkiestracja_grupa(id) ON DELETE CASCADE,
    krok      TEXT    NOT NULL,
    kolejnosc INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (grupa_id, krok)
);

CREATE TABLE orkiestracja_kompensacja (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok          TEXT    NOT NULL,
    krok_wycofu   TEXT    NOT NULL,
    zapisano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (automatyka_id, krok),
    -- Krok wycofujący sam siebie nie wycofuje niczego.
    CHECK (krok <> krok_wycofu)
);
