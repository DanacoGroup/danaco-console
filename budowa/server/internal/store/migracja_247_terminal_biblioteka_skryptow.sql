-- Migracja 247 zakłada tabele biblioteki skryptów i migawek modułu Terminal
-- wraz z pełną historią wersji treści, przechowywaną odrębnie od pozycji bieżącej.

CREATE TABLE terminal_skrypt (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    rodzaj         TEXT    NOT NULL DEFAULT 'script'
                           CHECK(rodzaj IN ('script','snippet')),
    -- Powłoka, dla której treść napisano; ona wyznacza narzędzie analizy.
    powloka        TEXT    NOT NULL,
    -- Treść bieżąca wraz z nagłówkiem deklaracji parametrów.
    tresc          TEXT    NOT NULL DEFAULT '',
    -- Znaczniki porządkujące wykaz, rozdzielone przecinkiem.
    znaczniki      TEXT    NOT NULL DEFAULT '',
    -- Skrót przywołujący snippet z palety poleceń.
    alias          TEXT,
    -- Numer wersji bieżącej; rośnie przy każdym zapisie treści.
    wersja         INTEGER NOT NULL DEFAULT 1 CHECK(wersja >= 1),
    -- Chwila ostatniego uruchomienia; pusta, póki pozycji nie uruchamiano.
    uruchomiono    TEXT,
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_terminal_skrypt_rodzaj ON terminal_skrypt(rodzaj, nazwa);
CREATE UNIQUE INDEX idx_terminal_skrypt_alias ON terminal_skrypt(alias)
    WHERE alias IS NOT NULL AND alias <> '';

-- Usunięcie pozycji biblioteki kasuje kaskadowo wszystkie jej wersje, ponieważ wersja bez pozycji nie ma nazwy ani powłoki.
CREATE TABLE terminal_skrypt_wersja (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    skrypt_kod  TEXT    NOT NULL REFERENCES terminal_skrypt(kod) ON DELETE CASCADE,
    wersja      INTEGER NOT NULL CHECK(wersja >= 1),
    tresc       TEXT    NOT NULL,
    zapisano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(skrypt_kod, wersja)
);
