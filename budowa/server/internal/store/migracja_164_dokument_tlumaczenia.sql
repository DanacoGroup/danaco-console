-- Migracja 164 zakłada tabelę dokumentu wniesionego do tłumaczenia wraz z segmentami niosącymi miejsce w strukturze pliku, numer strony i nazwę stylu.

CREATE TABLE dokument_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    sciezka                  TEXT    NOT NULL,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('docx','pdf','pptx','xlsx','odt','markdown','html')),
    liczba_stron             INTEGER,
    uzyto_ocr                INTEGER NOT NULL DEFAULT 0 CHECK(uzyto_ocr IN (0,1)),
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_dokument_tlumaczenia_okno ON dokument_tlumaczenia(okno_id, zaktualizowano DESC);

CREATE TABLE segment_dokumentu_tlumaczenia (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    dokument_id  INTEGER NOT NULL REFERENCES dokument_tlumaczenia(id) ON DELETE CASCADE,
    kolejnosc    INTEGER NOT NULL,
    tresc        TEXT    NOT NULL,
    sciezka_wezla TEXT,
    strona       INTEGER,
    styl         TEXT,
    UNIQUE(dokument_id, kolejnosc)
);
