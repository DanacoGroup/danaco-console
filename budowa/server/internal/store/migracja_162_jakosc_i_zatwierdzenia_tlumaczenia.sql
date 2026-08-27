-- Migracja 162 rozszerza termin słownika o stan i dziedzinę, zakłada profile kontroli jakości z regułami niezgodności oraz obieg zatwierdzeń panelu.
ALTER TABLE termin_slownika ADD COLUMN stan TEXT
    CHECK(stan IS NULL OR stan IN ('approved','candidate','preferred','forbidden'));
ALTER TABLE termin_slownika ADD COLUMN dziedzina TEXT;
CREATE INDEX idx_termin_slownika_stan ON termin_slownika(stan, dziedzina);

-- ── Profil kontroli jakości ─────────────────────────────────────────────────
-- Zasięg (`zasieg`) niesie wartości `ConfigScope` kontraktu wprost; `zasieg_id`
-- wskazuje byt zasięgu i zostaje pusty przy zasięgu bez bytu (global).
CREATE TABLE profil_qa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    zasieg                   TEXT    NOT NULL DEFAULT 'global',
    zasieg_id                TEXT,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_profil_qa_zasieg ON profil_qa(zasieg, zasieg_id);

CREATE TABLE profil_qa_kontrola (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id INTEGER NOT NULL REFERENCES profil_qa(id) ON DELETE CASCADE,
    rodzaj    TEXT    NOT NULL
                      CHECK(rodzaj IN ('number','date','currency','placeholder','length','omission')),
    waga      TEXT    NOT NULL CHECK(waga IN ('hint','warning','error')),
    wlaczona  INTEGER NOT NULL DEFAULT 1 CHECK(wlaczona IN (0,1)),
    UNIQUE(profil_id, rodzaj)
);

-- Zakłada tabelę zatwierdzenie_panelu niosącą historię etapów obiegu wraz z autorem i uwagą, uzupełnioną migawką na panelu.
ALTER TABLE panel_tlumaczenia ADD COLUMN etap_zatwierdzenia TEXT
    CHECK(etap_zatwierdzenia IS NULL OR
          etap_zatwierdzenia IN ('translation','proofreading','approved','rejected'));
ALTER TABLE panel_tlumaczenia ADD COLUMN zatwierdzil TEXT;
ALTER TABLE panel_tlumaczenia ADD COLUMN zatwierdzono INTEGER;

CREATE TABLE zatwierdzenie_panelu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    panel_id                 INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    etap                     TEXT    NOT NULL
                                     CHECK(etap IN ('translation','proofreading','approved','rejected')),
    autor                    TEXT    NOT NULL,
    uwaga                    TEXT,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_zatwierdzenie_panelu_panel ON zatwierdzenie_panelu(panel_id, utworzono DESC);
