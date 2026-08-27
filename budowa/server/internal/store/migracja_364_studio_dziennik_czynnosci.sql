-- Migracja 364 dodaje odwracalny dziennik czynności dokumentu wraz
-- z zależnościami między czynnościami.

CREATE TABLE czynnosc_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    kolejnosc                INTEGER NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('textEdit','formatChange','styleChange','pageChange',
                                           'listChange','tableChange','objectChange','apparatusChange',
                                           'markupChange','proposalAccept','importChange')),
    -- Rodzaj autora rozróżnia człowieka i wykonawcę; tożsamość niesie kod
    -- agenta modułu Agents.
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
CREATE INDEX idx_czynnosc_dokumentu_studio_dokument
    ON czynnosc_dokumentu_studio(dokument_id, kolejnosc DESC);
CREATE INDEX idx_czynnosc_dokumentu_studio_autor
    ON czynnosc_dokumentu_studio(dokument_id, autor_rodzaj, autor_agent_kod);
CREATE INDEX idx_czynnosc_dokumentu_studio_stan
    ON czynnosc_dokumentu_studio(dokument_id, stan, kolejnosc);

CREATE TABLE zaleznosc_czynnosci_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    czynnosc_id              INTEGER NOT NULL
                                     REFERENCES czynnosc_dokumentu_studio(id) ON DELETE CASCADE,
    podstawa_id              INTEGER NOT NULL
                                     REFERENCES czynnosc_dokumentu_studio(id) ON DELETE CASCADE,
    powod                    TEXT,
    UNIQUE(czynnosc_id, podstawa_id),
    -- Czynność nie może stać na sobie samej, bo cofnięcie stałoby się
    -- niemożliwe do naprawienia.
    CHECK(czynnosc_id <> podstawa_id)
);
CREATE INDEX idx_zaleznosc_czynnosci_studio_podstawa
    ON zaleznosc_czynnosci_studio(podstawa_id, czynnosc_id);
