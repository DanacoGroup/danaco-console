-- Migracja 193 — ustalenia analizy i rejestr dowodów.
--
-- Ustalenie analizy zostaje w zapisie, bo `roundtable.analysis.run` kosztuje
-- wywołanie kanału modelu. Odczyt wyniku bez powtórnego wywołania jest tu
-- warunkiem użyteczności: transkrypt z ustaleniami (`includeAnalysis`) wydaje
-- się długo po tym, jak analiza przebiegła.
--
-- Rejestr dowodów jest osobny od ustaleń, bo ma inne pytanie: nie „co model
-- rozpoznał", tylko „czym poparto twierdzenie". Kolumna `poparte` niesie
-- odpowiedź, a brak źródła przy `poparte = 0` jest właśnie tym, co
-- `unsupportedOnly` wyciąga na wierzch.

CREATE TABLE debata_ustalenie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('argumentMining','speechAct','fallacy',
                                                      'steelman','dedup','factCheck','tone')),
    tura                     TEXT    NOT NULL DEFAULT '',
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    wezel                    TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    tresc                    TEXT    NOT NULL,
    pewnosc                  REAL    NOT NULL DEFAULT -1.0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_ustalenie_okno ON debata_ustalenie(okno, id);

CREATE TABLE debata_dowod (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tura                     TEXT    NOT NULL DEFAULT '',
    wypowiedz                TEXT    NOT NULL,
    twierdzenie              TEXT    NOT NULL,
    zrodlo                   TEXT,
    adres                    TEXT,
    poparte                  INTEGER NOT NULL DEFAULT 0 CHECK(poparte IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_dowod_okno ON debata_dowod(okno, id);
