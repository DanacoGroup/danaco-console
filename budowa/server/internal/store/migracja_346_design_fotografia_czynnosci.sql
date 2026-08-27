-- Migracja 346 dodaje łańcuch czynności edycji zasobu obrazowego modułu
-- Design wraz z zapisanymi nastawami każdej czynności.

CREATE TABLE czynnosc_fotografii_design (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    zasob_id        INTEGER NOT NULL REFERENCES zasob_design(id) ON DELETE CASCADE,
    zasob_zrodla_id INTEGER          REFERENCES zasob_design(id) ON DELETE SET NULL,
    komenda         TEXT    NOT NULL,
    nastawy_json    TEXT,
    policzone_przez TEXT,
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Odczyt idzie po zasobie wynikowym, czyli łańcuchem wstecz, oraz po źródle,
-- czyli tym, co z zasobu powstało dalej.
CREATE INDEX idx_czynnosc_fotografii_design_zasob
    ON czynnosc_fotografii_design(zasob_id, id);
CREATE INDEX idx_czynnosc_fotografii_design_zrodlo
    ON czynnosc_fotografii_design(zasob_zrodla_id, id);
