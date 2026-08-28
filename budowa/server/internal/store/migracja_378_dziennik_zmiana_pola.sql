-- Migracja 378 przebudowuje tabelę czynności dokumentu, dodając rodzaj
-- zmiany pola do słownika czynności dziennika.

PRAGMA foreign_keys = OFF;

CREATE TABLE czynnosc_dokumentu_studio_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    kolejnosc                INTEGER NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('textEdit','formatChange','styleChange','pageChange',
                                           'listChange','tableChange','objectChange','apparatusChange',
                                           'markupChange','proposalAccept','importChange','fieldChange')),
    autor_rodzaj             TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor_rodzaj IN ('uzytkownik','model')),
    autor_agent_kod          TEXT,
    autor_agent_nazwa        TEXT,
    autor_agent_wersja       TEXT,
    autor_podagent_kod       TEXT,
    opis                     TEXT    NOT NULL,
    zakres_od                INTEGER,
    zakres_do                INTEGER,
    stan_przed_json          TEXT,
    stan_po_json             TEXT,
    zmiana_sledzona_kod      TEXT,
    zadanie_kod              TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'active'
                                     CHECK(stan IN ('active','reverted')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(dokument_id, kolejnosc)
);

INSERT INTO czynnosc_dokumentu_studio_nowa
    (id, identyfikator_zewnetrzny, dokument_id, kolejnosc, rodzaj, autor_rodzaj,
     autor_agent_kod, autor_agent_nazwa, autor_agent_wersja, autor_podagent_kod,
     opis, zakres_od, zakres_do, stan_przed_json, stan_po_json,
     zmiana_sledzona_kod, zadanie_kod, stan, utworzono, zaktualizowano)
SELECT id, identyfikator_zewnetrzny, dokument_id, kolejnosc, rodzaj, autor_rodzaj,
       autor_agent_kod, autor_agent_nazwa, autor_agent_wersja, autor_podagent_kod,
       opis, zakres_od, zakres_do, stan_przed_json, stan_po_json,
       zmiana_sledzona_kod, zadanie_kod, stan, utworzono, zaktualizowano
  FROM czynnosc_dokumentu_studio;

DROP TABLE czynnosc_dokumentu_studio;
ALTER TABLE czynnosc_dokumentu_studio_nowa RENAME TO czynnosc_dokumentu_studio;

CREATE INDEX IF NOT EXISTS idx_czynnosc_dokumentu_studio_dokument
    ON czynnosc_dokumentu_studio(dokument_id, kolejnosc DESC);
CREATE INDEX IF NOT EXISTS idx_czynnosc_dokumentu_studio_autor
    ON czynnosc_dokumentu_studio(dokument_id, autor_rodzaj, autor_agent_kod);
CREATE INDEX IF NOT EXISTS idx_czynnosc_dokumentu_studio_stan
    ON czynnosc_dokumentu_studio(dokument_id, stan, kolejnosc);

PRAGMA foreign_keys = ON;
