-- Migracja zakłada rejestr przenosin plików torem zdalnym, bo przy zasięgu zdalnym
-- proces i pliki bywają na różnych maszynach.

CREATE TABLE zdalne_przeniesienie (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id          INTEGER NOT NULL REFERENCES host_zdalny(id) ON DELETE CASCADE,
    okno_id          TEXT    NOT NULL DEFAULT '',
    -- Kierunek ruchu bajtów względem maszyny rdzenia.
    kierunek         TEXT    NOT NULL CHECK(kierunek IN ('wyslanie','pobranie')),
    sciezka_zrodlowa TEXT    NOT NULL,
    sciezka_docelowa TEXT    NOT NULL,
    -- Rozmiar w bajtach zmierzony lokalnie; zero, gdy pomiar się nie powiódł, bez unieważnienia.
    rozmiar          INTEGER NOT NULL DEFAULT 0,
    przeniesiono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zdalne_przeniesienie_host ON zdalne_przeniesienie(host_id, przeniesiono);
