-- Migracja 250 — obserwacje plików modułu Terminal (okno Task & Schedule).
--
-- Obserwacja uruchamia polecenie karty przy zmianie plików pasujących do
-- wzorca. Wyzwalacz plikowy istniał dotąd wyłącznie w automatykach, czyli poza
-- powłoką: nie dało się powiedzieć „po każdej zmianie w tym katalogu zbuduj
-- projekt W TEJ karcie, w jej katalogu i jej środowisku”.
--
-- `licznik` i `wyzwolono` są dziennikiem, nie stanem żywym — po restarcie rdzenia
-- obserwacja nie biegnie (przy montażu przechodzi w `stopped`), ale liczba
-- dotychczasowych wyzwoleń zostaje. Bez niej Operator nie odróżniłby obserwacji
-- założonej i niedziałającej od takiej, która po prostu nie miała czego złapać.
--
-- Klucza obcego do `terminal_karta` nie ma z tego samego powodu, dla którego nie
-- ma go karta do okna: karta bywa bytem pamięci rdzenia bez wiersza, a wpis
-- obserwacji nie ma prawa nie powstać z tego powodu.

CREATE TABLE terminal_obserwacja (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    okno_kod    TEXT    NOT NULL,
    karta_kod   TEXT    NOT NULL,
    -- Wzorzec ścieżek objętych obserwacją, liczony względem katalogu karty.
    wzorzec     TEXT    NOT NULL,
    -- Polecenie uruchamiane po zmianie, w powłoce i katalogu karty.
    polecenie   TEXT    NOT NULL,
    -- Tłumienie powtórzeń w milisekundach; puste bierze wartość rdzenia.
    tlumienie   INTEGER CHECK(tlumienie IS NULL OR tlumienie >= 0),
    -- Czy obserwacja obejmuje podkatalogi.
    rekurencyjnie INTEGER NOT NULL DEFAULT 0 CHECK(rekurencyjnie IN (0,1)),
    stan        TEXT    NOT NULL DEFAULT 'active'
                        CHECK(stan IN ('active','stopped','failed')),
    licznik     INTEGER NOT NULL DEFAULT 0 CHECK(licznik >= 0),
    wyzwolono   TEXT,
    powod       TEXT    NOT NULL DEFAULT '',
    zalozono    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_terminal_obserwacja_okno ON terminal_obserwacja(okno_kod, zalozono DESC);
CREATE INDEX idx_terminal_obserwacja_stan ON terminal_obserwacja(stan, zalozono DESC);
