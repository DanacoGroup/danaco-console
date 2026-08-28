-- Migracja 165 zakłada tabele zasobów lokalizacyjnych i ich kluczy oraz kwestii napisów wraz z indeksami wspierającymi odczyt okna i panelu.

-- Zakłada tabelę zasob_lokalizacji niosącą plik kluczy wniesiony do okna w jednym z ośmiu obsługiwanych formatów lokalizacji.
CREATE TABLE zasob_lokalizacji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    sciezka                  TEXT    NOT NULL,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('json','yaml','properties','androidXml',
                                                      'iosStrings','iosStringsdict','resx','gettextPo')),
    jezyk_zrodlowy           TEXT,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_zasob_lokalizacji_okno ON zasob_lokalizacji(okno_id, zaktualizowano DESC);

CREATE TABLE klucz_lokalizacji (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    zasob_id       INTEGER NOT NULL REFERENCES zasob_lokalizacji(id) ON DELETE CASCADE,
    klucz          TEXT    NOT NULL,
    tresc          TEXT    NOT NULL,
    -- Znaczniki podstawienia z treści, rozdzielone nową linią, bez osobnej tabeli.
    znaczniki      TEXT,
    kontekst       TEXT,
    zrzut_zasob_id TEXT,
    formy_mnogie   TEXT,
    kolejnosc      INTEGER NOT NULL DEFAULT 0,
    UNIQUE(zasob_id, klucz)
);

-- Zakłada tabelę kwestia_napisow niosącą wiersz napisów przypisany do panelu tłumaczenia albo do materiału źródłowego okna.
CREATE TABLE kwestia_napisow (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_id    INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    panel_id   INTEGER          REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    kolejnosc  INTEGER NOT NULL,
    poczatek_ms INTEGER NOT NULL,
    koniec_ms   INTEGER NOT NULL,
    tresc      TEXT    NOT NULL,
    mowca      TEXT
);
CREATE INDEX idx_kwestia_napisow_okno ON kwestia_napisow(okno_id, kolejnosc);
CREATE INDEX idx_kwestia_napisow_panel ON kwestia_napisow(panel_id, kolejnosc);
