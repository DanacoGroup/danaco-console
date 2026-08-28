-- Migracja 171 zakłada tabelę zakładek okna przeglądarki, oddzieloną od źródeł przeglądania i niosącą folder, etykiety w kolumnie JSON oraz notatkę.

CREATE TABLE zakladka_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    folder                   TEXT,
    etykiety_json            TEXT,
    notatka                  TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zakladka_przegladania_okno ON zakladka_przegladania(okno, utworzono DESC, id);
