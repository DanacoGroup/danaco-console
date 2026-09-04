-- Ciąg dalszy kroku 490. Nastawa, warstwa izolacji i zasób pamięci mają klucz
-- złożony z poziomu zasięgu i klucza wpisanego przez Operatora, a rejestracja
-- powiadomień — z urządzenia i adresu w kanale. Więz obejmujący całą tabelę
-- kazał tym kluczom być jednoznacznymi w instalacji: konto młodsze nie mogło
-- zapisać nastawy pod kluczem zajętym przez konto starsze, a upsert pod takim
-- kluczem sięgał cudzego wiersza.
--
-- Postępowanie i uzasadnienie wskaźnika po wyrażeniu — jak w kroku 490.
-- Wskaźniki zwykłe przebudowanych tabel odtwarza się w niezmienionej postaci.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

PRAGMA foreign_keys = off;

-- ── nastawa platformy ───────────────────────────────────────────────────────

CREATE TABLE ustawienie_nowa (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom_zasiegu_id INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu     TEXT    NOT NULL DEFAULT '',
    os                TEXT    NOT NULL DEFAULT 'platform'
                              REFERENCES os_zasiegu(kod) ON UPDATE CASCADE,
    klucz_osi         TEXT    NOT NULL DEFAULT '',
    klucz             TEXT    NOT NULL,
    wartosc           TEXT,
    rodzaj_wartosci   TEXT    NOT NULL DEFAULT 'tekst'
                              CHECK(rodzaj_wartosci IN ('tekst','liczba','logiczna','json')),
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id          INTEGER
);

INSERT INTO ustawienie_nowa
    (id, poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz, wartosc,
     rodzaj_wartosci, zaktualizowano, konto_id)
SELECT id, poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz, wartosc,
       rodzaj_wartosci, zaktualizowano, konto_id
  FROM ustawienie;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 491, 'ustawienie',
       (SELECT COUNT(*) FROM ustawienie) - (SELECT COUNT(*) FROM ustawienie_nowa);

DROP TABLE ustawienie;
ALTER TABLE ustawienie_nowa RENAME TO ustawienie;

CREATE INDEX idx_ustawienie_klucz ON ustawienie(klucz);
CREATE INDEX idx_ustawienie_os ON ustawienie(os, klucz_osi, klucz);

CREATE UNIQUE INDEX idx_ustawienie_adres
    ON ustawienie (poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz,
                   COALESCE(konto_id, 0));

-- ── warstwa izolacji ────────────────────────────────────────────────────────

CREATE TABLE warstwa_izolacji_nowa (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom_zasiegu_id INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu     TEXT    NOT NULL DEFAULT '',
    warstwa           TEXT    NOT NULL CHECK(warstwa IN ('default','session')),
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id          INTEGER
);

INSERT INTO warstwa_izolacji_nowa
    (id, poziom_zasiegu_id, klucz_zasiegu, warstwa, zaktualizowano, konto_id)
SELECT id, poziom_zasiegu_id, klucz_zasiegu, warstwa, zaktualizowano, konto_id
  FROM warstwa_izolacji;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 491, 'warstwa_izolacji',
       (SELECT COUNT(*) FROM warstwa_izolacji)
     - (SELECT COUNT(*) FROM warstwa_izolacji_nowa);

DROP TABLE warstwa_izolacji;
ALTER TABLE warstwa_izolacji_nowa RENAME TO warstwa_izolacji;

CREATE UNIQUE INDEX idx_warstwa_izolacji_adres
    ON warstwa_izolacji (poziom_zasiegu_id, klucz_zasiegu, COALESCE(konto_id, 0));

-- ── zasób pamięci ───────────────────────────────────────────────────────────

CREATE TABLE zasob_pamieci_nowa (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom          TEXT    NOT NULL
                            CHECK(poziom IN ('globalna','srodowisko','modul','projekt','sesja')),
    klucz_zasiegu   TEXT    NOT NULL DEFAULT '',
    klucz           TEXT    NOT NULL,
    tresc           TEXT,
    tresc_odwolanie TEXT,
    waga            INTEGER NOT NULL DEFAULT 0,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id        INTEGER
);

INSERT INTO zasob_pamieci_nowa
    (id, poziom, klucz_zasiegu, klucz, tresc, tresc_odwolanie, waga, utworzono,
     zaktualizowano, konto_id)
SELECT id, poziom, klucz_zasiegu, klucz, tresc, tresc_odwolanie, waga, utworzono,
       zaktualizowano, konto_id
  FROM zasob_pamieci;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 491, 'zasob_pamieci',
       (SELECT COUNT(*) FROM zasob_pamieci) - (SELECT COUNT(*) FROM zasob_pamieci_nowa);

DROP TABLE zasob_pamieci;
ALTER TABLE zasob_pamieci_nowa RENAME TO zasob_pamieci;

CREATE INDEX idx_zasob_pamieci_poziom ON zasob_pamieci(poziom, klucz_zasiegu);

CREATE UNIQUE INDEX idx_zasob_pamieci_adres
    ON zasob_pamieci (poziom, klucz_zasiegu, klucz, COALESCE(konto_id, 0));

-- ── rejestracja powiadomień urządzenia ──────────────────────────────────────

CREATE TABLE urzadzenie_powiadomien_nowa (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    urzadzenie_id        INTEGER NOT NULL REFERENCES urzadzenie(id) ON DELETE CASCADE,
    kanal                TEXT    NOT NULL DEFAULT 'polaczenie'
                                 CHECK(kanal IN ('polaczenie')),
    klucz_kanalu         TEXT    NOT NULL CHECK(TRIM(klucz_kanalu) <> ''),
    etykieta             TEXT,
    aktywne              INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),
    zarejestrowano       TEXT    NOT NULL
                                 DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    wyrejestrowano       TEXT,
    ostatnio_dostarczono TEXT,
    konto_id             INTEGER,
    CHECK((aktywne = 1) = (wyrejestrowano IS NULL))
);

INSERT INTO urzadzenie_powiadomien_nowa
    (id, urzadzenie_id, kanal, klucz_kanalu, etykieta, aktywne, zarejestrowano,
     wyrejestrowano, ostatnio_dostarczono, konto_id)
SELECT id, urzadzenie_id, kanal, klucz_kanalu, etykieta, aktywne, zarejestrowano,
       wyrejestrowano, ostatnio_dostarczono, konto_id
  FROM urzadzenie_powiadomien;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 491, 'urzadzenie_powiadomien',
       (SELECT COUNT(*) FROM urzadzenie_powiadomien)
     - (SELECT COUNT(*) FROM urzadzenie_powiadomien_nowa);

DROP TABLE urzadzenie_powiadomien;
ALTER TABLE urzadzenie_powiadomien_nowa RENAME TO urzadzenie_powiadomien;

CREATE INDEX idx_urzadzenie_powiadomien_czynne
    ON urzadzenie_powiadomien(kanal, klucz_kanalu) WHERE aktywne = 1;
CREATE INDEX idx_urzadzenie_powiadomien_urzadzenie
    ON urzadzenie_powiadomien(urzadzenie_id);

CREATE UNIQUE INDEX idx_urzadzenie_powiadomien_adres
    ON urzadzenie_powiadomien (urzadzenie_id, kanal, klucz_kanalu,
                               COALESCE(konto_id, 0));

PRAGMA foreign_keys = on;
