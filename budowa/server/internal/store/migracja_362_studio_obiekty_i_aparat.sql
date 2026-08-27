-- Migracja tworzy tabele obiektu osadzonego, elementu aparatu dokumentu i pola
-- dokumentu Studio: obiekt ma własny wiersz niezależny od drzewa postaci, aparat
-- łączy trzynaście rodzajów w jednej tabeli, a pole liczy się przy podglądzie.

CREATE TABLE obiekt_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('image','shape','icon','textbox','logo','chart')),
    zrodlo                   TEXT
                                     CHECK(zrodlo IS NULL OR zrodlo IN ('file','coreAsset','designModule',
                                           'photoBank','libraryFile','web','drawn')),
    zasob_kod                TEXT,
    design_wezel_kod         TEXT,
    biblioteka_plik_kod      TEXT,
    adres_zrodla             TEXT,
    zakotwiczenie            TEXT    NOT NULL DEFAULT 'paragraph'
                                     CHECK(zakotwiczenie IN ('character','paragraph','page')),
    zakotwiczenie_pozycja    INTEGER NOT NULL DEFAULT 0,
    warstwa                  INTEGER NOT NULL DEFAULT 0,
    tekst_zastepczy          TEXT,
    tekst_wewnetrzny         TEXT,
    podpis                   TEXT,
    postac_json              TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_obiekt_dokumentu_studio_dokument
    ON obiekt_dokumentu_studio(dokument_id, zakotwiczenie_pozycja, id);
CREATE INDEX idx_obiekt_dokumentu_studio_rodzaj
    ON obiekt_dokumentu_studio(dokument_id, rodzaj);

CREATE TABLE element_aparatu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('toc','figureIndex','tableIndex','footnote','endnote',
                                           'caption','bookmark','crossReference','hyperlink','citation',
                                           'bibliography','indexEntry','index')),
    kotwica_od               INTEGER NOT NULL DEFAULT 0,
    kotwica_do               INTEGER NOT NULL DEFAULT 0,
    numer                    TEXT,
    etykieta                 TEXT,
    tresc                    TEXT,
    cel_kod                  TEXT,
    cel_adres                TEXT,
    -- Kolejność liczy numerację w obrębie rodzaju, niezależnie od porządku zapisu klucza wiersza.
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    nieswiezy                INTEGER NOT NULL DEFAULT 0,
    dane_json                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_element_aparatu_studio_dokument
    ON element_aparatu_studio(dokument_id, rodzaj, kolejnosc, id);
CREATE INDEX idx_element_aparatu_studio_kotwica
    ON element_aparatu_studio(dokument_id, kotwica_od);
CREATE INDEX idx_element_aparatu_studio_nieswieze
    ON element_aparatu_studio(dokument_id, nieswiezy);

CREATE TABLE pole_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('pageNumber','pageCount','date','time','documentTitle',
                                           'documentAuthor','documentProperty','calculated','templateField')),
    kotwica                  INTEGER NOT NULL DEFAULT 0,
    format                   TEXT,
    wyrazenie                TEXT,
    nazwa_wlasciwosci        TEXT,
    wartosc                  TEXT,
    nieswieze                INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pole_dokumentu_studio_dokument
    ON pole_dokumentu_studio(dokument_id, kotwica, id);
