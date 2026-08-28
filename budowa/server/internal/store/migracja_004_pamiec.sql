-- Migracja 004 — obszar pamięci wielopoziomowej.
--
-- Schemat obszaru pamięci należy wyłącznie do tego pliku. IF NOT EXISTS jest tu
-- potrzebne: migracja musi przejść także na bazie, w której obie tabele już stoją.
--
-- Poziomy pamięci: pięć wartości, najwęższą jest sesja — okno komunikacji nie
-- jest poziomem pamięci. Treść obszerna trafia do pliku (`tresc_odwolanie`).

CREATE TABLE IF NOT EXISTS zasob_pamieci (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom                 TEXT    NOT NULL
                                   CHECK(poziom IN ('globalna','srodowisko','modul','projekt','sesja')),
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    klucz                  TEXT    NOT NULL,
    tresc                  TEXT,
    tresc_odwolanie        TEXT,
    waga                   INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom, klucz_zasiegu, klucz)
);

CREATE INDEX IF NOT EXISTS idx_zasob_pamieci_poziom ON zasob_pamieci(poziom, klucz_zasiegu);

CREATE TABLE IF NOT EXISTS konfiguracja_pamieci_sesji (
    sesja_id               INTEGER PRIMARY KEY REFERENCES sesja(id) ON DELETE CASCADE,
    poziomy_wlaczone       TEXT    NOT NULL DEFAULT 'globalna,srodowisko,modul,projekt,sesja',
    zapis_wlaczony         INTEGER NOT NULL DEFAULT 1 CHECK(zapis_wlaczony IN (0,1)),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
