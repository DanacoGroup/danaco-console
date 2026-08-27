-- Migracja zakłada trwałość pętli koordynator-wykonawca: zlecenie przekazania okna wraz
-- z kompletem kontekstu i dziennik wykonania akcji okna.

-- Jeden wiersz na jedno przekazanie okna; pozycja kolejki wiąże zlecenie z jedynym
-- silnikiem kolejek, bez dublowania.
CREATE TABLE zlecenie_przekazania (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    sesja_id               INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,
    okno_zrodlowe_id       INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    okno_docelowe_id       INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    pozycja_kolejki_id     INTEGER NOT NULL REFERENCES pozycja_kolejki(id) ON DELETE CASCADE,
    polecenie              TEXT    NOT NULL,
    -- Komplet kontekstu w całości, zapis JSON; puste, gdy przekazanie nie niesie kompletu kontekstu.
    komplet_kontekstu      TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zlecenie_przekazania_docelowe ON zlecenie_przekazania(okno_docelowe_id, utworzono);
CREATE INDEX idx_zlecenie_przekazania_zrodlowe ON zlecenie_przekazania(okno_zrodlowe_id, utworzono);
CREATE INDEX idx_zlecenie_przekazania_pozycja ON zlecenie_przekazania(pozycja_kolejki_id);

-- Katalog akcji mówi, jakie akcje istnieją; ta tabela mówi, kiedy i z jakim skutkiem
-- konkretne okno je wykonało.
CREATE TABLE log_akcji_okna (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    akcja_id               TEXT    NOT NULL,
    parametry              TEXT,
    wynik                  TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_log_akcji_okna_okno ON log_akcji_okna(okno_komunikacji_id, utworzono);
CREATE INDEX idx_log_akcji_okna_akcja ON log_akcji_okna(akcja_id);
