-- Migracja 221 — wskaźnik treści plików projektu.
--
-- Wiersz niesie treść wydobytą z pliku biblioteki projektu wraz ze sposobem,
-- którym ją wydobyto. Treść leży tutaj, a nie w pliku obok materiału, bo
-- wyszukiwanie po słowach idzie zapytaniem do bazy — jednym, nad wszystkimi
-- bytami projektu naraz.
--
-- Plik wskazywany jest ścieżką względną katalogu roboczego projektu — tym samym
-- identyfikatorem, którym plik wychodzi z `workspace.library.list`. Warunek
-- UNIQUE na parze pilnuje jednego wyciągu na plik: powtórne wydobycie
-- nadpisuje, a nie dokłada drugiego zdania o tym samym pliku.

CREATE TABLE wyciag_tekstu_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    plik                     TEXT    NOT NULL,
    sposob                   TEXT    NOT NULL
                                     CHECK(sposob IN ('text','pdf','docx','ocr')),
    tresc                    TEXT    NOT NULL DEFAULT '',
    liczba_znakow            INTEGER NOT NULL DEFAULT 0,
    -- Języki rozpoznania rozdzielone znakiem nowego wiersza; puste przy odczycie
    -- warstwy tekstowej, bo żaden język nie brał wtedy udziału.
    jezyki                   TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(projekt_id, plik)
);
CREATE INDEX idx_wyciag_tekstu_projekt ON wyciag_tekstu_projektu(projekt_id);
