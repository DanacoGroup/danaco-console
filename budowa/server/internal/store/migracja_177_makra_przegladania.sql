-- Migracja 177 zakłada tabelę makr przeglądania, przechowującą kroki automatyzacji jako kolumnę JSON oraz stan nagrywania i odwołanie do przebiegu automatyki.

CREATE TABLE makro_przegladania (
    id                        INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny  TEXT    NOT NULL UNIQUE,
    okno                      TEXT    NOT NULL,
    nazwa                     TEXT    NOT NULL,
    kroki_json                TEXT,
    nagrywanie                INTEGER NOT NULL DEFAULT 0 CHECK(nagrywanie IN (0,1)),
    automatyka_zewnetrzna_id  TEXT,
    utworzono                 TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_makro_przegladania_okno ON makro_przegladania(okno, utworzono DESC, id);
