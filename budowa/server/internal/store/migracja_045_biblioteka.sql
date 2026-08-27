-- Migracja 045 tworzy trwałość modułu biblioteki: plik repozytorium wiedzy, jego
-- wersje, etykiety i kolekcje zasobów, z treścią przechowywaną na dysku, nie
-- w bazie.

-- Tabela plik_biblioteki przechowuje plik repozytorium wiedzy wraz z metadanymi,
-- ścieżką na dysku i odwołaniem do treści bieżącej wersji.
CREATE TABLE plik_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    sciezka                  TEXT,
    mime_type                TEXT,
    rozmiar_bajtow           INTEGER,
    projekt_id               TEXT,
    modul_zrodlowy_id        TEXT,
    -- Suma kontrolna i treść powtarzają się tu z bieżącej wersji, by uniknąć złączenia przy liście plików.
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

-- Tabela wersja_pliku_biblioteki przechowuje każdą wersję pliku jako osobny
-- wiersz historii z własnym autorem i sumą kontrolną, nie jako licznik.
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
-- Indeks porządkuje wersje pliku od najnowszej, ponieważ komenda
-- library.version.list zwraca historię wersji w tej kolejności.
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
-- Indeks wspiera filtrowanie plików po etykiecie, ponieważ komenda
-- library.file.list przyjmuje etykiety jako warunek żądania.
CREATE INDEX idx_etykieta_pliku_biblioteki_etykieta ON etykieta_pliku_biblioteki(etykieta, plik_id);

-- Tabela kolekcja_biblioteki przechowuje kolekcję zasobów jako byt z własną
-- nazwą i opisem, tworzony osobną komendą library.collection.create.
CREATE TABLE kolekcja_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kolekcja_biblioteki_nazwa ON kolekcja_biblioteki(nazwa);

-- Tabela przypisanie_kolekcji_biblioteki wiąże plik z kolekcją; library.collection.assign
-- przypisuje wiele plików naraz, a indeks w obu kierunkach wspiera też
-- filtrowanie plików po kolekcji.
CREATE TABLE przypisanie_kolekcji_biblioteki (
    kolekcja_id  INTEGER NOT NULL REFERENCES kolekcja_biblioteki(id) ON DELETE CASCADE,
    plik_id      INTEGER NOT NULL REFERENCES plik_biblioteki(id) ON DELETE CASCADE,
    utworzono    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (kolekcja_id, plik_id)
);
CREATE INDEX idx_przypisanie_kolekcji_biblioteki_plik ON przypisanie_kolekcji_biblioteki(plik_id, kolekcja_id);
