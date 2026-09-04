-- Domknięcie kroków 490–492. Ścieżka nagrania mowy, odcisk wpisu schowka, klucz
-- wyniku odkrycia badania, adres poświadczenia automatyki, zasięg zasady
-- retencji pamięci oraz adres fragmentu wiedzy i wiersza podobieństwa obrazu
-- opisują pracę jednego Operatora. Więz obejmujący całą tabelę mieszał tę pracę
-- między kontami: ten sam plik indeksowany w dwóch kontach nadpisywał się
-- wzajemnie zamiast stać w dwóch wierszach.
--
-- Postępowanie i uzasadnienie wskaźnika po wyrażeniu — jak w kroku 490.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

PRAGMA foreign_keys = off;

-- ── nagranie mowy ───────────────────────────────────────────────────────────

CREATE TABLE nagranie_mowy_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL,
    typ_tresci               TEXT    NOT NULL,
    rozmiar_bajtow           INTEGER NOT NULL,
    dlugosc_ms               INTEGER,
    sesja_kod                TEXT,
    okno_kod                 TEXT,
    trwale                   INTEGER NOT NULL DEFAULT 0,
    utworzono                INTEGER NOT NULL,
    wygasa                   INTEGER,
    konto_id                 INTEGER
);

INSERT INTO nagranie_mowy_nowa
    (id, identyfikator_zewnetrzny, sciezka, typ_tresci, rozmiar_bajtow, dlugosc_ms,
     sesja_kod, okno_kod, trwale, utworzono, wygasa, konto_id)
SELECT id, identyfikator_zewnetrzny, sciezka, typ_tresci, rozmiar_bajtow, dlugosc_ms,
       sesja_kod, okno_kod, trwale, utworzono, wygasa, konto_id
  FROM nagranie_mowy;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'nagranie_mowy',
       (SELECT COUNT(*) FROM nagranie_mowy) - (SELECT COUNT(*) FROM nagranie_mowy_nowa);

DROP TABLE nagranie_mowy;
ALTER TABLE nagranie_mowy_nowa RENAME TO nagranie_mowy;

CREATE INDEX idx_nagranie_mowy_okno ON nagranie_mowy(okno_kod, utworzono DESC);
CREATE INDEX idx_nagranie_mowy_wygasa ON nagranie_mowy(wygasa);

CREATE UNIQUE INDEX idx_nagranie_mowy_sciezka
    ON nagranie_mowy (sciezka, COALESCE(konto_id, 0));

-- ── wpis schowka ────────────────────────────────────────────────────────────

CREATE TABLE wpis_schowka_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL DEFAULT 'text',
    tresc                    TEXT    NOT NULL,
    odcisk                   TEXT    NOT NULL,
    zajawka                  TEXT,
    rozmiar_bajtow           INTEGER NOT NULL DEFAULT 0,
    przypiety                INTEGER NOT NULL DEFAULT 0,
    wrazliwy                 INTEGER NOT NULL DEFAULT 0,
    okno_zrodlowe            TEXT,
    utworzono                INTEGER NOT NULL,
    uzyto                    INTEGER,
    postac_json              TEXT,
    konto_id                 INTEGER
);

INSERT INTO wpis_schowka_nowa
    (id, identyfikator_zewnetrzny, rodzaj, tresc, odcisk, zajawka, rozmiar_bajtow,
     przypiety, wrazliwy, okno_zrodlowe, utworzono, uzyto, postac_json, konto_id)
SELECT id, identyfikator_zewnetrzny, rodzaj, tresc, odcisk, zajawka, rozmiar_bajtow,
       przypiety, wrazliwy, okno_zrodlowe, utworzono, uzyto, postac_json, konto_id
  FROM wpis_schowka;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'wpis_schowka',
       (SELECT COUNT(*) FROM wpis_schowka) - (SELECT COUNT(*) FROM wpis_schowka_nowa);

DROP TABLE wpis_schowka;
ALTER TABLE wpis_schowka_nowa RENAME TO wpis_schowka;

CREATE INDEX idx_wpis_schowka_wykaz
    ON wpis_schowka(przypiety DESC, utworzono DESC, id DESC);
CREATE INDEX idx_wpis_schowka_rodzaj ON wpis_schowka(rodzaj, utworzono DESC);

CREATE UNIQUE INDEX idx_wpis_schowka_odcisk
    ON wpis_schowka (odcisk, COALESCE(konto_id, 0));

-- ── wynik odkrycia badania ──────────────────────────────────────────────────

CREATE TABLE wynik_odkrycia_badania_nowa (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    klucz            TEXT    NOT NULL,
    okno             TEXT    NOT NULL,
    tytul            TEXT    NOT NULL,
    adres            TEXT,
    autorzy          TEXT    NOT NULL DEFAULT '[]',
    rok              INTEGER,
    dostawca         TEXT    NOT NULL DEFAULT '',
    identyfikator    TEXT,
    fragment         TEXT,
    otwarty_dostep   INTEGER,
    duplikat         INTEGER NOT NULL DEFAULT 0,
    odrzucony        INTEGER NOT NULL DEFAULT 0,
    powod_odrzucenia TEXT,
    zrodlo_kod       TEXT,
    utworzono        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id         INTEGER
);

INSERT INTO wynik_odkrycia_badania_nowa
    (id, klucz, okno, tytul, adres, autorzy, rok, dostawca, identyfikator, fragment,
     otwarty_dostep, duplikat, odrzucony, powod_odrzucenia, zrodlo_kod, utworzono, konto_id)
SELECT id, klucz, okno, tytul, adres, autorzy, rok, dostawca, identyfikator, fragment,
       otwarty_dostep, duplikat, odrzucony, powod_odrzucenia, zrodlo_kod, utworzono, konto_id
  FROM wynik_odkrycia_badania;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'wynik_odkrycia_badania',
       (SELECT COUNT(*) FROM wynik_odkrycia_badania)
     - (SELECT COUNT(*) FROM wynik_odkrycia_badania_nowa);

DROP TABLE wynik_odkrycia_badania;
ALTER TABLE wynik_odkrycia_badania_nowa RENAME TO wynik_odkrycia_badania;

CREATE INDEX idx_wynik_odkrycia_badania_okno
    ON wynik_odkrycia_badania(okno, utworzono DESC);

CREATE UNIQUE INDEX idx_wynik_odkrycia_badania_klucz
    ON wynik_odkrycia_badania (klucz, COALESCE(konto_id, 0));

-- ── poświadczenie automatyki ────────────────────────────────────────────────

CREATE TABLE poswiadczenie_automatyki_nowa (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    odwolanie      TEXT    NOT NULL,
    nazwa          TEXT    NOT NULL,
    zasieg         TEXT,
    zasieg_id      TEXT,
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id       INTEGER
);

INSERT INTO poswiadczenie_automatyki_nowa
    (id, odwolanie, nazwa, zasieg, zasieg_id, zaktualizowano, konto_id)
SELECT id, odwolanie, nazwa, zasieg, zasieg_id, zaktualizowano, konto_id
  FROM poswiadczenie_automatyki;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'poswiadczenie_automatyki',
       (SELECT COUNT(*) FROM poswiadczenie_automatyki)
     - (SELECT COUNT(*) FROM poswiadczenie_automatyki_nowa);

DROP TABLE poswiadczenie_automatyki;
ALTER TABLE poswiadczenie_automatyki_nowa RENAME TO poswiadczenie_automatyki;

CREATE INDEX idx_poswiadczenie_automatyki_nazwa ON poswiadczenie_automatyki(nazwa, id);

CREATE UNIQUE INDEX idx_poswiadczenie_automatyki_odwolanie
    ON poswiadczenie_automatyki (odwolanie, COALESCE(konto_id, 0));

CREATE UNIQUE INDEX idx_poswiadczenie_automatyki_zasieg
    ON poswiadczenie_automatyki (nazwa, zasieg, zasieg_id, COALESCE(konto_id, 0));

-- ── zasada retencji pamięci ─────────────────────────────────────────────────

CREATE TABLE zasada_retencji_pamieci_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zasieg                   TEXT    NOT NULL,
    zasieg_kod               TEXT    NOT NULL DEFAULT '',
    profil_kod               TEXT    NOT NULL DEFAULT '',
    dni_wygasania            INTEGER NOT NULL DEFAULT 0,
    wrazliwe_domyslnie       INTEGER NOT NULL DEFAULT 0,
    wzorce_json              TEXT    NOT NULL DEFAULT '[]',
    czynna                   INTEGER NOT NULL DEFAULT 1,
    zaktualizowano           INTEGER NOT NULL,
    konto_id                 INTEGER
);

INSERT INTO zasada_retencji_pamieci_nowa
    (id, identyfikator_zewnetrzny, zasieg, zasieg_kod, profil_kod, dni_wygasania,
     wrazliwe_domyslnie, wzorce_json, czynna, zaktualizowano, konto_id)
SELECT id, identyfikator_zewnetrzny, zasieg, zasieg_kod, profil_kod, dni_wygasania,
       wrazliwe_domyslnie, wzorce_json, czynna, zaktualizowano, konto_id
  FROM zasada_retencji_pamieci;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'zasada_retencji_pamieci',
       (SELECT COUNT(*) FROM zasada_retencji_pamieci)
     - (SELECT COUNT(*) FROM zasada_retencji_pamieci_nowa);

DROP TABLE zasada_retencji_pamieci;
ALTER TABLE zasada_retencji_pamieci_nowa RENAME TO zasada_retencji_pamieci;

CREATE UNIQUE INDEX idx_zasada_retencji_pamieci_zasieg
    ON zasada_retencji_pamieci (zasieg, zasieg_kod, profil_kod, COALESCE(konto_id, 0));

-- ── podobieństwo obrazu ─────────────────────────────────────────────────────

CREATE TABLE podobienstwo_obrazu_nowa (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    sciezka      TEXT    NOT NULL,
    odcisk       TEXT    NOT NULL,
    model        TEXT    NOT NULL,
    pytanie      TEXT    NOT NULL,
    podobienstwo REAL    NOT NULL,
    utworzono    INTEGER NOT NULL,
    konto_id     INTEGER
);

INSERT INTO podobienstwo_obrazu_nowa
    (id, sciezka, odcisk, model, pytanie, podobienstwo, utworzono, konto_id)
SELECT id, sciezka, odcisk, model, pytanie, podobienstwo, utworzono, konto_id
  FROM podobienstwo_obrazu;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'podobienstwo_obrazu',
       (SELECT COUNT(*) FROM podobienstwo_obrazu)
     - (SELECT COUNT(*) FROM podobienstwo_obrazu_nowa);

DROP TABLE podobienstwo_obrazu;
ALTER TABLE podobienstwo_obrazu_nowa RENAME TO podobienstwo_obrazu;

CREATE UNIQUE INDEX idx_podobienstwo_obrazu_adres
    ON podobienstwo_obrazu (sciezka, model, pytanie, COALESCE(konto_id, 0));

-- ── fragment wiedzy ─────────────────────────────────────────────────────────

CREATE TABLE fragment_wiedzy_nowa (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    zakres     TEXT    NOT NULL CHECK (zakres IN ('library', 'history', 'workspace')),
    zrodlo     TEXT    NOT NULL,
    zrodlo_kod TEXT    NOT NULL DEFAULT '',
    kolejnosc  INTEGER NOT NULL,
    tresc      TEXT    NOT NULL,
    model      TEXT    NOT NULL,
    wymiar     INTEGER NOT NULL CHECK (wymiar > 0),
    wektor     BLOB    NOT NULL,
    utworzono  INTEGER NOT NULL,
    konto_id   INTEGER
);

INSERT INTO fragment_wiedzy_nowa
    (id, zakres, zrodlo, zrodlo_kod, kolejnosc, tresc, model, wymiar, wektor,
     utworzono, konto_id)
SELECT id, zakres, zrodlo, zrodlo_kod, kolejnosc, tresc, model, wymiar, wektor,
       utworzono, konto_id
  FROM fragment_wiedzy;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 493, 'fragment_wiedzy',
       (SELECT COUNT(*) FROM fragment_wiedzy) - (SELECT COUNT(*) FROM fragment_wiedzy_nowa);

DROP TABLE fragment_wiedzy;
ALTER TABLE fragment_wiedzy_nowa RENAME TO fragment_wiedzy;

CREATE INDEX idx_fragment_wiedzy_zakres_model ON fragment_wiedzy (model, zakres);
CREATE INDEX idx_fragment_wiedzy_zrodlo ON fragment_wiedzy (zakres, zrodlo_kod);

CREATE UNIQUE INDEX idx_fragment_wiedzy_adres
    ON fragment_wiedzy (zakres, zrodlo_kod, kolejnosc, model, COALESCE(konto_id, 0));

PRAGMA foreign_keys = on;
