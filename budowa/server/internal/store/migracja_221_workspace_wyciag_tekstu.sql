-- Migracja 221 zakłada tabelę wyciag_tekstu_projektu niosącą treść wydobytą z pliku biblioteki projektu wraz ze sposobem jej wydobycia.

CREATE TABLE wyciag_tekstu_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    plik                     TEXT    NOT NULL,
    sposob                   TEXT    NOT NULL
                                     CHECK(sposob IN ('text','pdf','docx','ocr')),
    tresc                    TEXT    NOT NULL DEFAULT '',
    liczba_znakow            INTEGER NOT NULL DEFAULT 0,
    -- Języki rozpoznania rozdzielone znakiem nowego wiersza; puste przy odczycie warstwy tekstowej.
    jezyki                   TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(projekt_id, plik)
);
CREATE INDEX idx_wyciag_tekstu_projekt ON wyciag_tekstu_projektu(projekt_id);
