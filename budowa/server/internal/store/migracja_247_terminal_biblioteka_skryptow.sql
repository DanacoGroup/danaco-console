-- Migracja 247 — biblioteka skryptów i snippetów modułu Terminal
-- (okno Script Library).
--
-- Dlaczego wersje mają własną tabelę, a nie kolumnę `wersja` na pozycji.
-- Kontrakt komendy `terminal.script.save` mówi wprost: zapis zakłada KOLEJNĄ
-- WERSJĘ, a numer nadaje rdzeń. Pozycja z jedną kolumną treści traciłaby
-- poprzednie brzmienie przy każdym zapisie, a skrypt uruchamiany na maszynach
-- Operatora jest dokładnie tym rodzajem treści, do której trzeba móc wrócić po
-- nieudanej poprawce. Tabela pozycji niesie więc treść BIEŻĄCĄ (żeby wykaz
-- czytał się jednym zapytaniem), a tabela wersji — pełny ślad.
--
-- `alias` jest skrótem przywołującym snippet z palety poleceń, więc musi być
-- jednoznaczny w obrębie biblioteki. Warunek UNIQUE stoi na wyrażeniu, nie na
-- kolumnie wprost: aliasu nie ma większość pozycji, a NULL w SQLite nie zderza
-- się z NULL, więc wiele pozycji bez aliasu współistnieje bez przeszkody.

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

-- Wersje pozycji. Kasowanie kaskadowe jest z zamysłu: `terminal.script.remove`
-- usuwa pozycję WRAZ ZE WSZYSTKIMI JEJ WERSJAMI, a wersja bez pozycji nie ma
-- nazwy ani powłoki i nie da się jej pokazać.
CREATE TABLE terminal_skrypt_wersja (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    skrypt_kod  TEXT    NOT NULL REFERENCES terminal_skrypt(kod) ON DELETE CASCADE,
    wersja      INTEGER NOT NULL CHECK(wersja >= 1),
    tresc       TEXT    NOT NULL,
    zapisano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(skrypt_kod, wersja)
);
