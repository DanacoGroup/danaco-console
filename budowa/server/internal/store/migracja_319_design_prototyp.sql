-- Migracja 319 dodaje graf prototypu modułu Design: połączenia przejść
-- między ramkami wraz z wyzwalaczem i czasem.

CREATE TABLE polaczenie_prototypu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    ramka_od_kod             TEXT    NOT NULL,
    ramka_do_kod             TEXT    NOT NULL,
    wyzwalacz                TEXT    NOT NULL,
    przejscie                TEXT    NOT NULL,
    czas_ms                  INTEGER,
    warstwa_kod              TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_polaczenie_prototypu_design_kompozycja
    ON polaczenie_prototypu_design(kompozycja_id, id);
CREATE INDEX idx_polaczenie_prototypu_design_od
    ON polaczenie_prototypu_design(ramka_od_kod);
