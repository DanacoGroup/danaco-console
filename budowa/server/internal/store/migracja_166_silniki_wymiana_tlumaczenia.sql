-- Migracja tworzy nastawy silników i wymianę zewnętrzną modułu Translate: profile
-- silników, politykę tłumaczenia pivotowego, pakiety przekazania, przebieg
-- pakietowy i most do dokumentu.

-- Tabela profil_silnika_tlumaczenia trzyma profil silnika tłumaczenia wraz z
-- kanałami modelu w kolejności pierwszeństwa wyboru.
CREATE TABLE profil_silnika_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    dziedzina                TEXT,
    zasieg                   TEXT    NOT NULL DEFAULT 'global',
    zasieg_id                TEXT,
    adaptacyjny              INTEGER NOT NULL DEFAULT 0 CHECK(adaptacyjny IN (0,1)),
    zasieg_pamieci           TEXT    CHECK(zasieg_pamieci IS NULL OR
                                           zasieg_pamieci IN ('card','project','team')),
    temperatura              REAL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_profil_silnika_zasieg ON profil_silnika_tlumaczenia(zasieg, zasieg_id);

CREATE TABLE kanal_profilu_silnika (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id INTEGER NOT NULL REFERENCES profil_silnika_tlumaczenia(id) ON DELETE CASCADE,
    kanal_kod TEXT    NOT NULL,
    kolejnosc INTEGER NOT NULL,
    UNIQUE(profil_id, kanal_kod)
);

-- Tabela polityka_pivota trzyma politykę tłumaczenia pivotowego jedną na zasięg
-- konfiguracji, z parami języków jako jej dzieckiem.
CREATE TABLE polityka_pivota (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    zasieg          TEXT    NOT NULL DEFAULT 'global',
    zasieg_id       TEXT    NOT NULL DEFAULT '',
    jezyk_domyslny  TEXT,
    zaktualizowano  INTEGER NOT NULL,
    UNIQUE(zasieg, zasieg_id)
);

CREATE TABLE para_pivota (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    polityka_id  INTEGER NOT NULL REFERENCES polityka_pivota(id) ON DELETE CASCADE,
    jezyk_zrodla TEXT    NOT NULL,
    jezyk_celu   TEXT    NOT NULL,
    jezyk_pivota TEXT    NOT NULL,
    UNIQUE(polityka_id, jezyk_zrodla, jezyk_celu)
);

-- Tabela pakiet_przekazania trzyma pakiet przekazania wykonawcy wraz ze stanem
-- cyklu życia: złożony, przekazany, zwrot przyjęty.
CREATE TABLE pakiet_przekazania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    -- Zawartości i panele pakietu są bytem bez własnej tożsamości: nikt nie adresuje pojedynczej pozycji.
    zawartosci               TEXT    NOT NULL,
    panele                   TEXT    NOT NULL,
    instrukcje               TEXT,
    sciezka                  TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'built'
                                     CHECK(stan IN ('built','sent','returned')),
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_pakiet_przekazania_okno ON pakiet_przekazania(okno_id, utworzono DESC);

-- Tabele zlecenie_pakietu_tlumaczenia i pozycja_pakietu_tlumaczenia trzymają
-- własną kolejkę przebiegu pakietowego modułu Translate.
CREATE TABLE zlecenie_pakietu_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    utworzono                INTEGER NOT NULL
);

CREATE TABLE pozycja_pakietu_tlumaczenia (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    zlecenie_id INTEGER NOT NULL REFERENCES zlecenie_pakietu_tlumaczenia(id) ON DELETE CASCADE,
    panel_kod   TEXT    NOT NULL,
    operacja    TEXT    NOT NULL
                        CHECK(operacja IN ('translate','qualityCheck','proofread','export')),
    stan        TEXT    NOT NULL DEFAULT 'pending'
                        CHECK(stan IN ('pending','done','error')),
    szczegol    TEXT,
    zaktualizowano INTEGER NOT NULL
);
CREATE INDEX idx_pozycja_pakietu_zlecenie ON pozycja_pakietu_tlumaczenia(zlecenie_id, id);

-- Tabela most_tlumaczenia trzyma dwustronne wiązanie okna Translate z dokumentem,
-- z którego wzięto materiał do tłumaczenia.
CREATE TABLE most_tlumaczenia (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_id       INTEGER NOT NULL UNIQUE REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    dokument_kod  TEXT    NOT NULL,
    zakotwiczenie TEXT,
    zaktualizowano INTEGER NOT NULL
);
