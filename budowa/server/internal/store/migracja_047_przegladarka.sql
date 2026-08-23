-- Migracja 047 — trwałość modułu Browser: migawka strony, źródło zebrane
-- w toku przeglądania i notatka powiązana ze źródłem.
--
-- Migawka jest tabelą historii nawigacji, nie osobnym bytem „odwiedzona strona”.
-- `browser.navigate` i `browser.snapshot.get` oba zwracają `BrowserSnapshot` —
-- każde wywołanie jest nowym zrzutem stanu strony w danej chwili. Wiersz
-- `migawka_strony` wstawiany przy każdym z tych wywołań, z porządkiem po
-- `utworzono`, jest już historią nawigacji: druga tabela „odwiedzona strona”
-- niosłaby te same fakty (adres, tytuł, czas) drugi raz.
--
-- Treść obszerna trafia do pliku, baza trzyma odwołanie. Migawka strony bywa duża
-- (pełny HTML, długi tekst renderowany) — kolumny `tekst_odwolanie`
-- i `zrodlo_odwolanie` niosą odwołanie do pliku, nie treść wprost. `ScreenshotRef`
-- z kontraktu jest już odwołaniem po stronie klienta, więc kolumna
-- `zrzut_odwolanie` zapisuje ją bez przekształcenia.
--
-- Okno jest kolumną tekstową, nie więzem obcym: `windowId` modułu Browser jest
-- oknem operacyjnym, nie oknem komunikacji, więc więzu do `okno_komunikacji` tu
-- nie ma.
--
-- „Źródło” tego modułu nie jest „źródłem” modułu Research. Kontrakt niesie oba
-- pod nazwą zaczynającą się od „Source”, ale kształty się rozjeżdżają:
-- `ResearchSource` ma `kind` (web/document/note/dataset), `credibility`
-- i `libraryFileId` — jest oceną wiarygodności zasobu badawczego, niekoniecznie
-- strony WWW. `BrowserSource` ma `snapshotId` i `key` — jest odciskiem strony
-- zebranym w toku przeglądania, zawsze powiązanym z migawką tego samego okna.
-- Inny kształt i różny cykl życia dają osobną tabelę `zrodlo_przegladania`,
-- a nie współdzielenie `zrodlo_badania`.
--
-- `zrodlo_zewnetrzny_id` w notatce jest wartością tekstową, nie więzem obcym:
-- `SourceId` kontraktu jest identyfikatorem zewnętrznym, opcjonalnym; notatka bez
-- źródła (dotycząca całej strony, nie jednego zebranego odcinka) musi zostać
-- zapisywalna, więc więzu NOT NULL / REFERENCES tu nie ma.

-- ── Migawka strony — Screenshot/Snapshot pane ─────────────────────────────────
CREATE TABLE migawka_strony (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    tekst_odwolanie          TEXT,
    zrodlo_odwolanie         TEXT,
    zrzut_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_migawka_strony_okno ON migawka_strony(okno, utworzono DESC, id);

-- ── Źródło zebrane w toku przeglądania — Sources drawer ───────────────────────
CREATE TABLE zrodlo_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    migawka_zewnetrzna_id    TEXT,
    kluczowe                 INTEGER NOT NULL DEFAULT 0 CHECK(kluczowe IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zrodlo_przegladania_okno ON zrodlo_przegladania(okno, utworzono DESC, id);

-- ── Notatka powiązana ze źródłem — Notes panel ────────────────────────────────
CREATE TABLE notatka_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    zrodlo_zewnetrzny_id     TEXT,
    tresc                    TEXT    NOT NULL DEFAULT '',
    cytat                    TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_notatka_przegladania_okno ON notatka_przegladania(okno, utworzono DESC, id);
