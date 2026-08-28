-- Migracja 174 zakłada tabelę kolejki czytania, niosącą stan przeczytania pozycji oraz znacznik czasu przypomnienia w formacie ISO 8601.

CREATE TABLE pozycja_czytania_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    notatka                  TEXT,
    przeczytana              INTEGER NOT NULL DEFAULT 0 CHECK(przeczytana IN (0,1)),
    przypomnienie            TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pozycja_czytania_okno ON pozycja_czytania_przegladania(okno, utworzono DESC, id);
