-- Migracja 223 — komentarze przy bytach projektu.
--
-- Jedna tabela na komentarze wszystkich bytów: cel opisuje para rodzaj + byt,
-- a nie osobna tabela per byt. Wątek składa wskazanie komentarza nadrzędnego —
-- tak samo jak strona wiki składa hierarchię wskazaniem strony nadrzędnej.
--
-- Przywołania znakiem małpy leżą kolumną obok treści: rozpoznaje się je raz,
-- przy zapisie, bo panel powiadomień pyta o nie częściej, niż komentarz się
-- zmienia.

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
