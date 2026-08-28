-- Migracja 172 zakłada tabelę monitorów zmian strony, przechowującą odwołanie do treści odniesienia, jej długość oraz stan ostatniego sprawdzenia.

CREATE TABLE monitor_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    selektor                 TEXT,
    interwal_sekund          INTEGER NOT NULL DEFAULT 3600,
    prog_zmiany              INTEGER,
    kanal_powiadomienia      TEXT,
    wlaczony                 INTEGER NOT NULL DEFAULT 1 CHECK(wlaczony IN (0,1)),
    stan                     TEXT    NOT NULL DEFAULT 'pending',
    odniesienie_odwolanie    TEXT,
    odniesienie_dlugosc      INTEGER,
    sprawdzono               TEXT,
    zmieniono                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_monitor_przegladania_okno ON monitor_przegladania(okno, utworzono DESC, id);
