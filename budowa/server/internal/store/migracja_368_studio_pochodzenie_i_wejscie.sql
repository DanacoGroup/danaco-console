-- Migracja 368 dodaje pochodzenie fragmentów dokumentu, autozamianę znaków
-- oraz zapamiętaną postać malarza formatów.

CREATE TABLE pochodzenie_fragmentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('web','libraryFile','research','clipboard',
                                           'template','importedFile')),
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    adres_zrodla             TEXT,
    biblioteka_plik_kod      TEXT,
    wersja_zrodla            TEXT,
    tytul_zrodla             TEXT,
    siegnieto                TEXT,
    autor_rodzaj             TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor_rodzaj IN ('uzytkownik','model')),
    autor_agent_kod          TEXT,
    autor_agent_nazwa        TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
CREATE INDEX idx_pochodzenie_fragmentu_studio_dokument
    ON pochodzenie_fragmentu_studio(dokument_id, zakres_od, id);
CREATE INDEX idx_pochodzenie_fragmentu_studio_rodzaj
    ON pochodzenie_fragmentu_studio(dokument_id, rodzaj);

CREATE TABLE autozamiana_znaku_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    skrot                    TEXT    NOT NULL UNIQUE,
    zamiennik                TEXT    NOT NULL,
    czynna                   INTEGER NOT NULL DEFAULT 1,
    fabryczna                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Zasady fabryczne autozamiany obejmują znaki interpunkcyjne niedostępne
-- z klawiatury i wchodzą do bazy, aby dało się je wyłączyć.
INSERT INTO autozamiana_znaku_studio (skrot, zamiennik, fabryczna) VALUES
    ('--',   '–', 1),
    ('---',  '—', 1),
    ('(c)',  '©', 1),
    ('(r)',  '®', 1),
    ('(tm)', '™', 1),
    ('...',  '…', 1),
    ('->',   '→', 1),
    ('<-',   '←', 1),
    ('+-',   '±', 1),
    ('!=',   '≠', 1),
    ('<=',   '≤', 1),
    ('>=',   '≥', 1);

CREATE TABLE postac_malarza_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    dokument_kod             TEXT,
    postac_znaku_json        TEXT,
    postac_akapitu_json      TEXT,
    styl_nazwany             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Pobrana postać wygasa, bo malarz jest narzędziem jednej czynności.
    wygasa                   TEXT
);
CREATE INDEX idx_postac_malarza_studio_okno
    ON postac_malarza_studio(okno, utworzono DESC, id DESC);
