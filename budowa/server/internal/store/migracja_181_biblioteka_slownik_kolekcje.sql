-- Migracja 181 zakłada tabele słownika etykiet i relacji tezaurusa modułu Library oraz rozszerza kolekcje o rodzica i regułę kolekcji inteligentnej.

-- Zakłada tabelę etykieta_slownika_biblioteki niosącą byt etykiety niezależny od tego, czy dziś nosi ją jakikolwiek zasób.
CREATE TABLE etykieta_slownika_biblioteki (
    nazwa      TEXT PRIMARY KEY,
    -- Barwa jest kodem żetonu interfejsu, nie wartością szesnastkową przeżywającą zmianę palety.
    barwa      TEXT,
    utworzono  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Zakłada tabelę relacja_tezaurusa_biblioteki wiążącą dwie etykiety relacją modelu SKOS zapisaną jednym kierunkiem.
CREATE TABLE relacja_tezaurusa_biblioteki (
    etykieta_zrodlowa TEXT NOT NULL,
    etykieta_docelowa TEXT NOT NULL,
    rodzaj            TEXT NOT NULL
                      CHECK(rodzaj IN ('nadrzedna','podrzedna','pokrewna')),
    utworzono         TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (etykieta_zrodlowa, etykieta_docelowa, rodzaj)
);
CREATE INDEX idx_relacja_tezaurusa_biblioteki_cel
    ON relacja_tezaurusa_biblioteki(etykieta_docelowa, rodzaj);

-- Rozszerza kolekcja_biblioteki o rodzica dającego hierarchię oraz regułę, z której wynika skład kolekcji inteligentnej.
ALTER TABLE kolekcja_biblioteki
    ADD COLUMN rodzic_id INTEGER REFERENCES kolekcja_biblioteki(id) ON DELETE SET NULL;
ALTER TABLE kolekcja_biblioteki ADD COLUMN regula_kod TEXT;

CREATE INDEX idx_kolekcja_biblioteki_rodzic ON kolekcja_biblioteki(rodzic_id, nazwa);
