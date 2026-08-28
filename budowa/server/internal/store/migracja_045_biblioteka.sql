-- Migracja 045 — trwałość modułu Library: repozytorium wiedzy, jego wersje,
-- etykiety oraz kolekcje zasobów (Library Explorer, Versioning Panel,
-- Tags & Collections).
--
-- Treść pliku trzyma dysk, nie baza: kolumna `tresc_odwolanie` niesie odwołanie
-- do pliku na dysku, baza nie dostaje kolumny BLOB. Ten sam wzorzec powtarza się
-- na dwóch poziomach — plik biblioteki i każda jego wersja — bo
-- `library.version.restore` musi umieć przywrócić treść sprzed zmiany, więc
-- treść poprzednich wersji musi przeżyć nadpisanie bieżącej.
--
-- Wersja jest własnym bytem, nie polem licznika. `library.version.list` zwraca
-- listę `LibraryVersion` z własnym `id`, autorem i sumą kontrolną — to nie jest
-- rosnący numer przy pliku, tylko osobny wiersz historii, bo każda wersja niesie
-- własną treść i własnego autora (kontrakt: `LibraryVersion.author`).
--
-- Etykiety i kolekcje mają rozłączne tabele — to dwie różne prawdy o pliku:
-- etykieta jest wolnym tekstem bez własnej tożsamości, kolekcja jest bytem
-- z nazwą i opisem tworzonym osobną komendą `library.collection.create`. Stąd
-- etykieta żyje jako wiersz w tabeli złącznikowej z gołym tekstem, a przypisanie
-- do kolekcji odwołuje się do wiersza `kolekcja_biblioteki`.
--
-- Podgląd (`LibraryPreview`) nie ma własnej tabeli. To widok obliczany w locie
-- z bieżącej wersji pliku (rodzaj podglądu wynika z `mime_type`, treść
-- z `tresc_odwolanie` wersji) — trwały byt tu jest jeden: wersja pliku.

-- ── Plik repozytorium wiedzy ───────────────────────────────────────────────
CREATE TABLE plik_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    sciezka                  TEXT,
    mime_type                TEXT,
    rozmiar_bajtow           INTEGER,
    projekt_id               TEXT,
    modul_zrodlowy_id        TEXT,
    -- Suma kontrolna i odwołanie do treści opisują bieżącą wersję. Trzymamy je
    -- też tutaj (obok wiersza w `wersja_pliku_biblioteki`), bo `library.file.list`
    -- i `library.file.search` czytają listę plików bez dołączania wersji —
    -- powtórzenie kolumn oszczędza złączenie na ścieżce najczęstszej.
    wersja_biezaca_id        INTEGER,
    suma_kontrolna           TEXT,
    tresc_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- library.file.list i library.file.search filtrują po projekcie i sortują od
-- najnowszych; library.file.search dodatkowo przeszukuje nazwę.
CREATE INDEX idx_plik_biblioteki_projekt ON plik_biblioteki(projekt_id, zaktualizowano DESC);
CREATE INDEX idx_plik_biblioteki_nazwa ON plik_biblioteki(nazwa);

-- ── Wersja pliku — Versioning Panel ────────────────────────────────────────
CREATE TABLE wersja_pliku_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    plik_id                  INTEGER NOT NULL REFERENCES plik_biblioteki(id) ON DELETE CASCADE,
    etykieta                 TEXT,
    autor                    TEXT,
    rozmiar_bajtow           INTEGER,
    suma_kontrolna           TEXT,
    tresc_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- library.version.list zwraca wersje pliku od najnowszej.
CREATE INDEX idx_wersja_pliku_biblioteki_plik ON wersja_pliku_biblioteki(plik_id, utworzono DESC, id DESC);

-- ── Etykieta pliku — Tags & Collections ────────────────────────────────────
-- Etykieta to wolny tekst bez własnej tożsamości (kontrakt: `LibraryFile.tags
-- []string`), więc para (plik, etykieta) jest kluczem — bez surogatu.
CREATE TABLE etykieta_pliku_biblioteki (
    plik_id    INTEGER NOT NULL REFERENCES plik_biblioteki(id) ON DELETE CASCADE,
    etykieta   TEXT    NOT NULL,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (plik_id, etykieta)
);
-- library.file.list filtruje po etykiecie (Tags []string w żądaniu).
CREATE INDEX idx_etykieta_pliku_biblioteki_etykieta ON etykieta_pliku_biblioteki(etykieta, plik_id);

-- ── Kolekcja zasobów — Tags & Collections ──────────────────────────────────
CREATE TABLE kolekcja_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kolekcja_biblioteki_nazwa ON kolekcja_biblioteki(nazwa);

-- ── Przypisanie pliku do kolekcji ──────────────────────────────────────────
-- library.collection.assign przypisuje wiele plików naraz do jednej kolekcji;
-- library.file.list filtruje po CollectionId — stąd indeks w obu kierunkach
-- (klucz główny okrywa kolekcja→plik, drugi indeks plik→kolekcja).
CREATE TABLE przypisanie_kolekcji_biblioteki (
    kolekcja_id  INTEGER NOT NULL REFERENCES kolekcja_biblioteki(id) ON DELETE CASCADE,
    plik_id      INTEGER NOT NULL REFERENCES plik_biblioteki(id) ON DELETE CASCADE,
    utworzono    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (kolekcja_id, plik_id)
);
CREATE INDEX idx_przypisanie_kolekcji_biblioteki_plik ON przypisanie_kolekcji_biblioteki(plik_id, kolekcja_id);
