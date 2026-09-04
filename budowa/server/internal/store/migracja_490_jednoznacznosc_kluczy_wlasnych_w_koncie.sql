-- Klucz, który wpisuje Operator, ma być jednoznaczny w jego koncie, nie w całej
-- instalacji. Skrót autozamiany, skrót tekstowy profilu oraz zasięg granicy
-- Wykonawcy i polityki pivota to nazwy wybierane ręcznie:
-- więz obejmujący całą tabelę odbierał kontu młodszemu nazwę zajętą wcześniej
-- przez konto starsze, a zapis pod tą nazwą sięgał cudzego wiersza.
--
-- Więz wbudowany w tabelę zdejmuje się wyłącznie przez przebudowę; jednoznaczność
-- wraca jako wskaźnik po wyrażeniu, wzorem migracji 485. Wskazanie puste wchodzi
-- przez COALESCE jako zero, bo wartości NULL w indeksie są sobie nierówne
-- i dwa wiersze zastane uszłyby za różne.
--
-- Rodzaj znacznika Studia zostaje z więzem obejmującym całą tabelę: kolumna
-- `nazwa` jest celem klucza obcego z `znakowanie_studio`, a SQLite wymaga dla
-- takiego celu wskaźnika jednoznacznego dokładnie po tej kolumnie.
--
-- Kontrola przejazdu pilnuje liczebności każdej przebudowanej tabeli: niezgodność
-- wywraca się o CHECK i cofa transakcję migracji.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

PRAGMA foreign_keys = off;

-- ── autozamiana znaku Studia ────────────────────────────────────────────────

CREATE TABLE autozamiana_znaku_studio_nowa (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    skrot          TEXT    NOT NULL,
    zamiennik      TEXT    NOT NULL,
    czynna         INTEGER NOT NULL DEFAULT 1,
    fabryczna      INTEGER NOT NULL DEFAULT 0,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id       INTEGER
);

INSERT INTO autozamiana_znaku_studio_nowa
    (id, skrot, zamiennik, czynna, fabryczna, utworzono, zaktualizowano, konto_id)
SELECT id, skrot, zamiennik, czynna, fabryczna, utworzono, zaktualizowano, konto_id
  FROM autozamiana_znaku_studio;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 490, 'autozamiana_znaku_studio',
       (SELECT COUNT(*) FROM autozamiana_znaku_studio)
     - (SELECT COUNT(*) FROM autozamiana_znaku_studio_nowa);

DROP TABLE autozamiana_znaku_studio;
ALTER TABLE autozamiana_znaku_studio_nowa RENAME TO autozamiana_znaku_studio;

CREATE UNIQUE INDEX idx_autozamiana_znaku_studio_skrot
    ON autozamiana_znaku_studio (skrot, COALESCE(konto_id, 0));

-- ── granica Wykonawcy przeglądania ──────────────────────────────────────────

CREATE TABLE granica_wykonawcy_przegladania_nowa (
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    zasieg                  TEXT    NOT NULL,
    zasieg_id               TEXT    NOT NULL DEFAULT '',
    max_krokow              INTEGER NOT NULL,
    max_czas_sekund         INTEGER NOT NULL,
    domeny_dozwolone_json   TEXT,
    domeny_zablokowane_json TEXT,
    potwierdzaj_wyslanie    INTEGER NOT NULL DEFAULT 1 CHECK(potwierdzaj_wyslanie IN (0,1)),
    zaktualizowano          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id                INTEGER
);

INSERT INTO granica_wykonawcy_przegladania_nowa
    (id, zasieg, zasieg_id, max_krokow, max_czas_sekund, domeny_dozwolone_json,
     domeny_zablokowane_json, potwierdzaj_wyslanie, zaktualizowano, konto_id)
SELECT id, zasieg, zasieg_id, max_krokow, max_czas_sekund, domeny_dozwolone_json,
       domeny_zablokowane_json, potwierdzaj_wyslanie, zaktualizowano, konto_id
  FROM granica_wykonawcy_przegladania;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 490, 'granica_wykonawcy_przegladania',
       (SELECT COUNT(*) FROM granica_wykonawcy_przegladania)
     - (SELECT COUNT(*) FROM granica_wykonawcy_przegladania_nowa);

DROP TABLE granica_wykonawcy_przegladania;
ALTER TABLE granica_wykonawcy_przegladania_nowa
    RENAME TO granica_wykonawcy_przegladania;

CREATE UNIQUE INDEX idx_granica_wykonawcy_przegladania_zasieg
    ON granica_wykonawcy_przegladania (zasieg, zasieg_id, COALESCE(konto_id, 0));

-- ── polityka pivota tłumaczenia ─────────────────────────────────────────────

CREATE TABLE polityka_pivota_nowa (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    zasieg         TEXT    NOT NULL DEFAULT 'global',
    zasieg_id      TEXT    NOT NULL DEFAULT '',
    jezyk_domyslny TEXT,
    zaktualizowano INTEGER NOT NULL,
    konto_id       INTEGER
);

INSERT INTO polityka_pivota_nowa
    (id, zasieg, zasieg_id, jezyk_domyslny, zaktualizowano, konto_id)
SELECT id, zasieg, zasieg_id, jezyk_domyslny, zaktualizowano, konto_id
  FROM polityka_pivota;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 490, 'polityka_pivota',
       (SELECT COUNT(*) FROM polityka_pivota)
     - (SELECT COUNT(*) FROM polityka_pivota_nowa);

DROP TABLE polityka_pivota;
ALTER TABLE polityka_pivota_nowa RENAME TO polityka_pivota;

CREATE UNIQUE INDEX idx_polityka_pivota_zasieg
    ON polityka_pivota (zasieg, zasieg_id, COALESCE(konto_id, 0));

-- ── skrót tekstowy profilu ──────────────────────────────────────────────────

CREATE TABLE skrot_tekstowy_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    profil_kod               TEXT    NOT NULL DEFAULT '',
    skrot                    TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL,
    opis                     TEXT,
    pola_json                TEXT    NOT NULL DEFAULT '[]',
    czynny                   INTEGER NOT NULL DEFAULT 1,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL,
    konto_id                 INTEGER
);

INSERT INTO skrot_tekstowy_nowa
    (id, identyfikator_zewnetrzny, profil_kod, skrot, tresc, opis, pola_json,
     czynny, utworzono, zaktualizowano, konto_id)
SELECT id, identyfikator_zewnetrzny, profil_kod, skrot, tresc, opis, pola_json,
       czynny, utworzono, zaktualizowano, konto_id
  FROM skrot_tekstowy;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 490, 'skrot_tekstowy',
       (SELECT COUNT(*) FROM skrot_tekstowy)
     - (SELECT COUNT(*) FROM skrot_tekstowy_nowa);

DROP TABLE skrot_tekstowy;
ALTER TABLE skrot_tekstowy_nowa RENAME TO skrot_tekstowy;

CREATE INDEX idx_skrot_tekstowy_wykaz ON skrot_tekstowy (profil_kod, skrot);

CREATE UNIQUE INDEX idx_skrot_tekstowy_profil_skrot
    ON skrot_tekstowy (profil_kod, skrot, COALESCE(konto_id, 0));

PRAGMA foreign_keys = on;
