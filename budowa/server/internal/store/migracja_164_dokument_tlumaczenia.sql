-- Migracja 164 — dokument wniesiony do tłumaczenia (`translate.document.load`,
-- `translate.document.render`, `translate.document.layout.compare`) razem
-- z jego segmentami.
--
-- Dokument jest bytem trwałym, bo kontrakt adresuje go identyfikatorem
-- w dwóch kolejnych komendach: `document.render` i `document.layout.compare`
-- przyjmują `documentId` wydany przy wczytaniu. Bez wiersza identyfikator
-- byłby napisem, którego rdzeń przy następnym wywołaniu nie rozpozna.
--
-- Segment dokumentu to nie to samo co segment okna (migracja 161). Segment
-- okna jest kawałkiem tekstu źródłowego; segment dokumentu niesie dodatkowo
-- miejsce w strukturze pliku: ścieżkę węzła (akapit, komórka, kształt), numer
-- strony i nazwę stylu. To one pozwalają złożyć dokument z powrotem
-- z zachowaniem układu, i to ich brak sprawiłby, że `document.render` oddałby
-- goły tekst zamiast dokumentu.
--
-- `uzyto_ocr` jest własnością wczytania, nie dokumentu na dysku: ten sam plik
-- wczytany dwa razy — raz z warstwy tekstowej, raz z rozpoznania pisma — daje
-- dwa różne materiały i model ma prawo wiedzieć, na którym pracuje.

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
