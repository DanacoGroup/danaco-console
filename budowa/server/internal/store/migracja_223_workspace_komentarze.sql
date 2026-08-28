-- Migracja 223 zakłada tabelę komentarzy przy bytach projektu, wspólną dla wszystkich rodzajów bytu, z przywołaniami znakiem małpy w osobnej kolumnie.

CREATE TABLE komentarz_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    rodzaj_bytu              TEXT    NOT NULL
                                     CHECK(rodzaj_bytu IN ('task','note','file','memoryEntry',
                                                           'instructions','comment','project')),
    byt                      TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL,
    komentarz_nadrzedny      TEXT    NOT NULL DEFAULT '',
    rodzaj_autora            TEXT    NOT NULL DEFAULT 'operator'
                                     CHECK(rodzaj_autora IN ('operator','agent')),
    autor                    TEXT    NOT NULL DEFAULT '',
    -- Przywołania rozdzielone znakiem nowego wiersza.
    przywolania              TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_komentarz_projektu_cel ON komentarz_projektu(projekt_id, rodzaj_bytu, byt, id);
CREATE INDEX idx_komentarz_projektu_watek ON komentarz_projektu(komentarz_nadrzedny);
