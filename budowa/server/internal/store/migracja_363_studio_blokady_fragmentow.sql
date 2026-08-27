-- Migracja 363 dodaje blokady fragmentów dokumentu chroniące treść przed
-- zmianą modelu oraz blokady wzorcowe szablonu pisma.

CREATE TABLE blokada_fragmentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    powod                    TEXT,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(zasieg IN ('model','everyone')),
    -- Rodzaj wskazuje, kto blokadę założył; zdejmuje ją wyłącznie użytkownik
    -- dokumentu.
    zalozyl_rodzaj           TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(zalozyl_rodzaj IN ('uzytkownik','model')),
    zalozyl_agent_kod        TEXT,
    zalozyl_agent_nazwa      TEXT,
    z_szablonu               INTEGER NOT NULL DEFAULT 0,
    szablon_kod              TEXT,
    przesuniecia             INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
-- Indeks po zakresie wspiera częste sprawdzenie, czy zakres żądania dotyka
-- blokady dokumentu przed każdą zmianą treści.
CREATE INDEX idx_blokada_fragmentu_studio_zakres
    ON blokada_fragmentu_studio(dokument_id, zakres_od, zakres_do);

-- Blokady wzorcowe szablonu pisma leżą osobno od blokad dokumentu i są
-- kopiowane do każdego dokumentu założonego z szablonu.
CREATE TABLE blokada_szablonu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    szablon_kod              TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    powod                    TEXT,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(zasieg IN ('model','everyone')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
CREATE INDEX idx_blokada_szablonu_studio_szablon
    ON blokada_szablonu_studio(szablon_kod, zakres_od);
