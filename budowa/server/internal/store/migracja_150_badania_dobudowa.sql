-- Migracja 150 dobudowuje moduł Research o byty brakujące wobec kontraktu:
-- katalogowanie źródeł, adnotacje lektury, książkę kodów, sprzeczności,
-- ślad prowenancji, odkrycia, monitory i wersje raportu. Uzasadnienie w docs.

ALTER TABLE zrodlo_badania ADD COLUMN stan_lektury TEXT NOT NULL DEFAULT 'unread';
ALTER TABLE zrodlo_badania ADD COLUMN etap_indeks INTEGER;
ALTER TABLE zrodlo_badania ADD COLUMN identyfikator TEXT;
ALTER TABLE zrodlo_badania ADD COLUMN csl_json TEXT;
ALTER TABLE zrodlo_badania ADD COLUMN tresc TEXT;
ALTER TABLE zrodlo_badania ADD COLUMN stron INTEGER;
ALTER TABLE zrodlo_badania ADD COLUMN warstwa_tekstu INTEGER NOT NULL DEFAULT 0;
ALTER TABLE zrodlo_badania ADD COLUMN wycofanie_stan TEXT;
ALTER TABLE zrodlo_badania ADD COLUMN wycofanie_adres TEXT;
ALTER TABLE zrodlo_badania ADD COLUMN wycofanie_sprawdzono TEXT;

CREATE TABLE etykieta_zrodla_badania (
    zrodlo_id INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    etykieta  TEXT    NOT NULL,
    PRIMARY KEY (zrodlo_id, etykieta)
);

CREATE TABLE kolekcja_zrodla_badania (
    zrodlo_id INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    kolekcja  TEXT    NOT NULL,
    PRIMARY KEY (zrodlo_id, kolekcja)
);

CREATE TABLE pytanie_zrodla_badania (
    zrodlo_id   INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    pytanie_kod TEXT    NOT NULL,
    PRIMARY KEY (zrodlo_id, pytanie_kod)
);

CREATE TABLE zalacznik_zrodla_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zrodlo_id                INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('fulltext','snapshot','note','data')),
    plik_biblioteki_id       TEXT,
    sciezka                  TEXT,
    rozmiar_bajtow           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zalacznik_zrodla_badania ON zalacznik_zrodla_badania(zrodlo_id, utworzono DESC);

CREATE TABLE streszczenie_zrodla_badania (
    zrodlo_id   INTEGER NOT NULL PRIMARY KEY REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    abstrakt    TEXT,
    tezy        TEXT NOT NULL DEFAULT '[]',
    metodologia TEXT,
    wnioski     TEXT NOT NULL DEFAULT '[]',
    utworzono   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE tabela_zrodla_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zrodlo_id                INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    podpis                   TEXT,
    naglowki                 TEXT    NOT NULL DEFAULT '[]',
    wiersze                  TEXT    NOT NULL DEFAULT '[]',
    kotwica_strona           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_tabela_zrodla_badania ON tabela_zrodla_badania(zrodlo_id, id);

CREATE TABLE adnotacja_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zrodlo_id                INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('highlight','note','bookmark')),
    cytat                    TEXT,
    komentarz                TEXT,
    kolor                    TEXT,
    kotwica_rodzaj           TEXT    NOT NULL DEFAULT 'page'
                                     CHECK(kotwica_rodzaj IN ('page','selector','timestamp')),
    kotwica_strona           INTEGER,
    kotwica_od               INTEGER,
    kotwica_do               INTEGER,
    kotwica_selektor         TEXT,
    kotwica_czas_ms          INTEGER,
    ustalenie_kod            TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_adnotacja_badania_zrodlo ON adnotacja_badania(zrodlo_id, utworzono DESC);

ALTER TABLE ustalenie_badania ADD COLUMN rodzaj TEXT;
ALTER TABLE ustalenie_badania ADD COLUMN waga TEXT;
ALTER TABLE ustalenie_badania ADD COLUMN notatka TEXT;
ALTER TABLE ustalenie_badania ADD COLUMN wymaga_potwierdzenia INTEGER NOT NULL DEFAULT 0;
ALTER TABLE ustalenie_badania ADD COLUMN adnotacja_kod TEXT;
ALTER TABLE ustalenie_badania ADD COLUMN kotwica_rodzaj TEXT;
ALTER TABLE ustalenie_badania ADD COLUMN kotwica_strona INTEGER;
ALTER TABLE ustalenie_badania ADD COLUMN kotwica_od INTEGER;
ALTER TABLE ustalenie_badania ADD COLUMN kotwica_do INTEGER;
ALTER TABLE ustalenie_badania ADD COLUMN kotwica_selektor TEXT;
ALTER TABLE ustalenie_badania ADD COLUMN kotwica_czas_ms INTEGER;

CREATE TABLE kod_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    nadrzedny_kod            TEXT
);
CREATE INDEX idx_kod_badania_okno ON kod_badania(okno, nazwa);

CREATE TABLE kod_ustalenia_badania (
    ustalenie_id INTEGER NOT NULL REFERENCES ustalenie_badania(id) ON DELETE CASCADE,
    kod_id       INTEGER NOT NULL REFERENCES kod_badania(id) ON DELETE CASCADE,
    PRIMARY KEY (ustalenie_id, kod_id)
);
CREATE INDEX idx_kod_ustalenia_badania_kod ON kod_ustalenia_badania(kod_id, ustalenie_id);

CREATE TABLE sprzecznosc_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    streszczenie             TEXT    NOT NULL,
    roznica_liczbowa         TEXT,
    rozstrzygnieta           INTEGER NOT NULL DEFAULT 0,
    ustalenie_rozstrzygajace TEXT,
    uzasadnienie             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sprzecznosc_badania_okno ON sprzecznosc_badania(okno, utworzono DESC);

CREATE TABLE ustalenie_sprzecznosci_badania (
    sprzecznosc_id INTEGER NOT NULL REFERENCES sprzecznosc_badania(id) ON DELETE CASCADE,
    ustalenie_kod  TEXT    NOT NULL,
    PRIMARY KEY (sprzecznosc_id, ustalenie_kod)
);

CREATE TABLE weryfikacja_ustalenia_badania (
    ustalenie_id INTEGER NOT NULL PRIMARY KEY REFERENCES ustalenie_badania(id) ON DELETE CASCADE,
    werdykt      TEXT    NOT NULL
                         CHECK(werdykt IN ('supported','unsupported','contradicted','inconclusive')),
    uzasadnienie TEXT    NOT NULL DEFAULT '',
    sprawdzono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE watek_ustalen_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    automatyczny             INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE ustalenie_watku_badania (
    watek_id      INTEGER NOT NULL REFERENCES watek_ustalen_badania(id) ON DELETE CASCADE,
    ustalenie_kod TEXT    NOT NULL,
    PRIMARY KEY (watek_id, ustalenie_kod)
);

CREATE TABLE prowenancja_badania (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    ustalenie_kod  TEXT    NOT NULL,
    o_czasie       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    aktor          TEXT    NOT NULL DEFAULT 'operator'
                           CHECK(aktor IN ('operator','assistant','model','core')),
    czynnosc       TEXT    NOT NULL,
    zrodlo_kod     TEXT,
    kotwica_strona INTEGER
);
CREATE INDEX idx_prowenancja_badania ON prowenancja_badania(ustalenie_kod, o_czasie);

ALTER TABLE przestrzen_badania ADD COLUMN odbiorca TEXT;
ALTER TABLE przestrzen_badania ADD COLUMN protokol TEXT;
ALTER TABLE przestrzen_badania ADD COLUMN granice TEXT;
ALTER TABLE przestrzen_badania ADD COLUMN notatka TEXT;

CREATE TABLE pytanie_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tekst                    TEXT    NOT NULL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE wynik_odkrycia_badania (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    klucz             TEXT    NOT NULL UNIQUE,
    okno              TEXT    NOT NULL,
    tytul             TEXT    NOT NULL,
    adres             TEXT,
    autorzy           TEXT    NOT NULL DEFAULT '[]',
    rok               INTEGER,
    dostawca          TEXT    NOT NULL DEFAULT '',
    identyfikator     TEXT,
    fragment          TEXT,
    otwarty_dostep    INTEGER,
    duplikat          INTEGER NOT NULL DEFAULT 0,
    odrzucony         INTEGER NOT NULL DEFAULT 0,
    powod_odrzucenia  TEXT,
    zrodlo_kod        TEXT,
    utworzono         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wynik_odkrycia_badania_okno ON wynik_odkrycia_badania(okno, utworzono DESC);

CREATE TABLE monitor_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL CHECK(rodzaj IN ('topic','feed')),
    zapytanie                TEXT,
    adres                    TEXT,
    interwal_minut           INTEGER,
    wlaczony                 INTEGER NOT NULL DEFAULT 1,
    oczekujace               INTEGER NOT NULL DEFAULT 0,
    odswiezono_o             TEXT
);
CREATE INDEX idx_monitor_badania_okno ON monitor_badania(okno, identyfikator_zewnetrzny);

CREATE TABLE styl_cytowania_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    wlasny                   INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE szablon_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    tytuly_sekcji            TEXT    NOT NULL DEFAULT '[]',
    wlasny                   INTEGER NOT NULL DEFAULT 1
);

ALTER TABLE raport_badania ADD COLUMN przypisy_umiejscowienie TEXT;
ALTER TABLE raport_badania ADD COLUMN przypisy_skrocone INTEGER NOT NULL DEFAULT 0;

CREATE TABLE wersja_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    etykieta                 TEXT,
    migawka                  TEXT    NOT NULL DEFAULT '[]',
    liczba_sekcji            INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wersja_raportu_badania ON wersja_raportu_badania(raport_id, utworzono DESC);

CREATE TABLE komentarz_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    sekcja_kod               TEXT,
    watek_kod                TEXT,
    tresc                    TEXT    NOT NULL,
    cytat                    TEXT,
    rozstrzygniety           INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_komentarz_raportu_badania ON komentarz_raportu_badania(raport_id, utworzono DESC);

CREATE TABLE blok_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    sekcja_kod               TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('matrix','timeline','chart','evidenceTable')),
    podpis                   TEXT,
    naglowki                 TEXT    NOT NULL DEFAULT '[]',
    wiersze                  TEXT    NOT NULL DEFAULT '[]',
    ustalenia                TEXT    NOT NULL DEFAULT '[]',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_blok_raportu_badania ON blok_raportu_badania(raport_id, sekcja_kod);

-- Eksport powstaje na nowo, bo więz CHECK na kolumnie format wymieniał pięć wartości, a kontrakt niesie już osiem.
CREATE TABLE eksport_raportu_badania_nowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('pdf','docx','markdown','html','txt',
                                                      'pptx','xlsx','latex')),
    cel                      TEXT    NOT NULL DEFAULT 'download'
                                     CHECK(cel IN ('download','library','studio','roundtable')),
    sciezka_docelowa         TEXT,
    plik_biblioteki_id       TEXT,
    sciezka_wyniku           TEXT,
    rozmiar_bajtow           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
INSERT INTO eksport_raportu_badania_nowy
    (id, identyfikator_zewnetrzny, raport_id, format, sciezka_docelowa,
     plik_biblioteki_id, sciezka_wyniku, rozmiar_bajtow, utworzono)
SELECT id, identyfikator_zewnetrzny, raport_id, format, sciezka_docelowa,
       plik_biblioteki_id, sciezka_wyniku, rozmiar_bajtow, utworzono
FROM eksport_raportu_badania;
DROP TABLE eksport_raportu_badania;
ALTER TABLE eksport_raportu_badania_nowy RENAME TO eksport_raportu_badania;
CREATE INDEX idx_eksport_raportu_badania_raport ON eksport_raportu_badania(raport_id, utworzono DESC);

CREATE TABLE szablon_eksportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('pdf','docx','markdown','html','txt',
                                                      'pptx','xlsx','latex')),
    cel                      TEXT,
    zawartosc                TEXT    NOT NULL DEFAULT '{}'
);

CREATE TABLE udostepnienie_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    adres                    TEXT    NOT NULL,
    wygasa_o                 TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_udostepnienie_raportu_badania ON udostepnienie_raportu_badania(raport_id, utworzono DESC);
