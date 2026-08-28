-- Przebudowuje tabele zasob_design i etykieta_zasobu_design, aby warunek CHECK rodzaju zasobu dopuścił dokument, dźwięk, film i archiwum obok wartości istniejących.

CREATE TABLE zasob_design_nowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('image','vector','composition',
                                                      'document','audio','video','archive')),
    format                   TEXT,
    uri                      TEXT,
    prompt_id                INTEGER REFERENCES prompt_design(id) ON DELETE SET NULL,
    wariant_zasobu_id        TEXT,
    ulubiony                 INTEGER NOT NULL DEFAULT 0 CHECK(ulubiony IN (0,1)),
    szerokosc                INTEGER,
    wysokosc                 INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO zasob_design_nowy
    (id, identyfikator_zewnetrzny, okno, nazwa, rodzaj, format, uri, prompt_id,
     wariant_zasobu_id, ulubiony, szerokosc, wysokosc, utworzono)
SELECT id, identyfikator_zewnetrzny, okno, nazwa, rodzaj, format, uri, prompt_id,
       wariant_zasobu_id, ulubiony, szerokosc, wysokosc, utworzono
FROM zasob_design;

CREATE TABLE etykieta_zasobu_design_nowa (
    zasob_id   INTEGER NOT NULL REFERENCES zasob_design_nowy(id) ON DELETE CASCADE,
    etykieta   TEXT    NOT NULL,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (zasob_id, etykieta)
);

INSERT INTO etykieta_zasobu_design_nowa (zasob_id, etykieta, utworzono)
SELECT zasob_id, etykieta, utworzono FROM etykieta_zasobu_design;

DROP TABLE etykieta_zasobu_design;
DROP TABLE zasob_design;

ALTER TABLE zasob_design_nowy RENAME TO zasob_design;
ALTER TABLE etykieta_zasobu_design_nowa RENAME TO etykieta_zasobu_design;

-- Indeksy giną razem ze starymi tabelami — odtwarzamy je co do znaku:
-- `design.asset.list` filtruje po oknie i po etykiecie,
-- sortując od najnowszych.
CREATE INDEX idx_zasob_design_okno ON zasob_design(okno, utworzono DESC, id DESC);
CREATE INDEX idx_zasob_design_prompt ON zasob_design(prompt_id);
CREATE INDEX idx_etykieta_zasobu_design_etykieta ON etykieta_zasobu_design(etykieta, zasob_id);
