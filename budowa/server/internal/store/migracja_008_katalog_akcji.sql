-- Migracja zakłada tabelę katalogu akcji sterowanego danymi: panel akcji powstaje
-- z rejestru wierszy, nie ze zmiany w kodzie.

CREATE TABLE akcja (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    ikona                  TEXT    NOT NULL DEFAULT '',
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    komenda                TEXT    NOT NULL,
    warunek_dostepnosci    TEXT    NOT NULL DEFAULT ''
                                   CHECK(warunek_dostepnosci IN ('','srodowisko','sesja',
                                                                 'okno','kolejka','kanal')),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    aktywna                INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_akcja_zasieg ON akcja(poziom_zasiegu_id, klucz_zasiegu, kolejnosc);
CREATE INDEX idx_akcja_komenda ON akcja(komenda);
