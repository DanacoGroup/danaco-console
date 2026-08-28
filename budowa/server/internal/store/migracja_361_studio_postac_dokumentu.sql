-- Migracja 361 tworzy trwałą postać dokumentu modułu studio: arkusz stylów
-- i sekcje, dotąd stan wyłącznie kliencki, teraz zapisywany razem z treścią.

CREATE TABLE postac_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    dokument_id              INTEGER NOT NULL UNIQUE
                                     REFERENCES dokument_studio(id) ON DELETE CASCADE,
    postac_json              TEXT    NOT NULL DEFAULT '{}',
    nastawy_strony_json      TEXT,
    wersja_postaci           INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Tabela styl_nazwany_studio przechowuje styl nazwany dokumentu; nazwa jest
-- jego jedynym identyfikatorem, a styl_nadrzedny wiąże dziedziczenie nazwą,
-- nie kluczem.
CREATE TABLE styl_nazwany_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    nazwa_widoczna           TEXT,
    rodzaj                   TEXT    NOT NULL DEFAULT 'paragraph'
                                     CHECK(rodzaj IN ('paragraph','character','table','list')),
    styl_nadrzedny           TEXT,
    styl_nastepny            TEXT,
    postac_znaku_json        TEXT,
    postac_akapitu_json      TEXT,
    fabryczny                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(dokument_id, nazwa)
);
CREATE INDEX idx_styl_nazwany_studio_dokument ON styl_nazwany_studio(dokument_id, nazwa);
CREATE INDEX idx_styl_nazwany_studio_nadrzedny ON styl_nazwany_studio(dokument_id, styl_nadrzedny);

-- Tabela sekcja_dokumentu_studio przechowuje sekcję z własnymi nastawami
-- strony, nagłówkami i stopkami jednym zapisem JSON, żeby dokument
-- z załącznikiem poziomym pozostał jednym dokumentem.
CREATE TABLE sekcja_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    tytul                    TEXT,
    zakres_od                INTEGER NOT NULL DEFAULT 0,
    zakres_do                INTEGER NOT NULL DEFAULT 0,
    rozpoczecie              TEXT    NOT NULL DEFAULT 'continuous'
                                     CHECK(rozpoczecie IN ('continuous','newPage','evenPage',
                                                           'oddPage','newColumn')),
    nastawy_strony_json      TEXT,
    naglowki_json            TEXT,
    numeracja_json           TEXT,
    znak_wodny_json          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sekcja_dokumentu_studio_dokument
    ON sekcja_dokumentu_studio(dokument_id, kolejnosc, id);
