-- Migracja zakłada tabelę powiązań kolejki z oknami, ekspertem, projektem albo
-- automatyką, wspólną dla wszystkich czterech rodzajów bytu.

CREATE TABLE powiazanie_kolejki (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id    INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    -- Cztery rodzaje bytów wymienione wprost przez kontrakt `queue.link`.
    rodzaj        TEXT    NOT NULL
                          CHECK(rodzaj IN ('okno','ekspert','projekt','automatyka')),
    -- Identyfikator bytu w kształcie, w jakim niesie go kontrakt.
    byt           TEXT    NOT NULL,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Powiązanie powtórzone nie jest drugim faktem — jest tym samym faktem.
    UNIQUE(kolejka_id, rodzaj, byt)
);

CREATE INDEX idx_powiazanie_kolejki_kolejka ON powiazanie_kolejki(kolejka_id, rodzaj);
CREATE INDEX idx_powiazanie_kolejki_byt ON powiazanie_kolejki(rodzaj, byt);
