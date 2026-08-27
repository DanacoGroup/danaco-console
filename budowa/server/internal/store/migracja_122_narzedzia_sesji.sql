-- Tworzy tabelę narzedzie_sesji, przechowującą narzędzia dołożone doraźnie do sesji, żyjące do jej zakończenia i usuwane kaskadą razem z sesją.

CREATE TABLE narzedzie_sesji (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    sesja_id       INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,

    -- Nazwa pełna ze źródłem, będąca tożsamością dołożenia w obrębie sesji.
    nazwa_pelna    TEXT    NOT NULL,
    -- Nazwa skrócona, którą Operator faktycznie wpisuje w komendzie.
    nazwa_skrocona TEXT    NOT NULL,
    -- Opis pełnym zdaniem, aby wykaz dołożeń nie cierpiał na brak opisu narzędzia.
    opis           TEXT    NOT NULL,

    rodzaj         TEXT    NOT NULL DEFAULT 'tool'
                           CHECK (rodzaj IN ('tool', 'action')),
    grupa          TEXT    NOT NULL DEFAULT '',
    -- Przedrostek źródła pozycji, pokazujący Operatorowi, skąd pozycja pochodzi.
    zrodlo_pozycji TEXT    NOT NULL DEFAULT '',

    zrodlo         TEXT    NOT NULL DEFAULT 'slashCommand'
                           CHECK (zrodlo IN ('slashCommand', 'assistant')),
    dolozono       INTEGER NOT NULL DEFAULT 0
);

-- Tworzy unikalny indeks jednego dołożenia na sesję i nazwę pełną, uniemożliwiający powstanie drugiego wiersza dla powtórzonej komendy.
CREATE UNIQUE INDEX idx_narzedzie_sesji_nazwa ON narzedzie_sesji(sesja_id, nazwa_pelna);

-- Tworzy indeks odczytu dołożeń danej sesji w kolejności dokładania, wspierający wykaz narzędzi sesji.
CREATE INDEX idx_narzedzie_sesji_wykaz ON narzedzie_sesji(sesja_id, dolozono, id);
