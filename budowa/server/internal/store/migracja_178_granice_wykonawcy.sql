-- Migracja 178 zakłada tabelę granic działania Wykonawcy na stronie, kluczowaną parą zasięgu i jego wskazania, z domenami dozwolonymi i zablokowanymi w JSON.

CREATE TABLE granica_wykonawcy_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    zasieg                   TEXT    NOT NULL,
    zasieg_id                TEXT    NOT NULL DEFAULT '',
    max_krokow               INTEGER NOT NULL,
    max_czas_sekund          INTEGER NOT NULL,
    domeny_dozwolone_json    TEXT,
    domeny_zablokowane_json  TEXT,
    potwierdzaj_wyslanie     INTEGER NOT NULL DEFAULT 1 CHECK(potwierdzaj_wyslanie IN (0,1)),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(zasieg, zasieg_id)
);
